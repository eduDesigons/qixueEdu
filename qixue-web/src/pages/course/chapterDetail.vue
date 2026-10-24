<template>
	<div class="chapter-detail">
		<div class="chapter-skeleton" v-if="!info && pageLoading">
			<div class="skeleton-player"></div>
			<div class="skeleton-line skeleton-title"></div>
			<div class="skeleton-line"></div>
			<div class="skeleton-line skeleton-short"></div>
		</div>
		<div class="chapter-load-fail" v-else-if="!info && pageError" @click="getChapterDetail(false)">
			<span>加载失败，点击重试</span>
		</div>

		<template v-else-if="info">
			<div class="video-player" v-if="info.type === 1 && info.media_url">
				<CourseVideoPlayer
					ref="courseVideo"
					:src="info.media_url"
					:poster="info.media && info.media.thumb ? info.media.thumb : ''"
					:title="info.title || ''"
					:initial-time="videoInitialTime"
					video-id="videoPlayer"
					@play="onVideoPlay"
					@pause="onVideoPause"
					@timeupdate="onVideoTimeUpdate"
					@ended="onVideoEnded"
					@canplay="onVideoCanPlay"
				/>
			</div>
			<div class="media-empty" v-else-if="info.type === 1">
				<van-icon name="play-circle-o" size="40" color="#b9bdc2" />
				<span>暂无课件</span>
			</div>

			<div class="audio-player" v-if="info.type === 2 && info.media_url">
				<div class="audio-player-bg" :style="{ backgroundImage: info.course && info.course.thumb ? `url(${static(info.course.thumb)})` : 'none' }"></div>
				<div class="audio-player-content">
				<div class="audio-cover-wrapper" @click="toggleAudio">
					<div class="audio-cover-container" :class="{ 'audio-cover-rotating': isPlaying }">
						<img :src="static(info.course.thumb)" class="audio-cover-img"  style="object-fit:cover" />
					</div>
					<div class="audio-cover-mask" v-if="!isPlaying">
						<van-icon name="play-circle" size="48" color="#fff"
							class="audio-play-icon" />
						<div class="audio-loading" v-if="audioLoading">
							<van-loading color="#fff" type="circular" />
						</div>
					</div>
					<div class="audio-loading" v-if="audioLoading && isPlaying">
						<van-loading color="#fff" type="circular" />
					</div>
				</div>
				<div class="audio-content">
					<div class="audio-title-row">
						<div class="audio-title">{{ info.title || '讲述' }}</div>
						<div class="audio-subtitle" v-if="info.course && info.course.title">{{ info.course.title }}</div>
					</div>
					<div class="audio-progress-wrapper" v-if="audioDuration > 0 || (info && info.media && info.media.duration > 0)">
						<input
							type="range"
							class="audio-slider"
							:value="audioProgressPercent"
							min="0"
							max="100"
							step="0.1"
							@input="onProgressChanging({ detail: { value: Number($event.target.value) } })"
							@change="onProgressChange({ detail: { value: Number($event.target.value) } })"
						/>
					</div>
					<div class="audio-time-row">
						<span class="audio-time-current">{{ formatTime(audioCurrentTime) || '0:00' }}</span>
						<span class="audio-time-total">{{ formatTime(audioDuration) || formatTime(info.media.duration) || '0:00' }}</span>
					</div>
				</div>
				<audio :id="'audioPlayer'" :src="static(info.media_url)" :poster="static(info.course.thumb)" @play="onAudioPlay"
					@pause="onAudioPause" @ended="onAudioEnded" @timeupdate="onAudioTimeUpdate" @canplay="onAudioCanPlay"
					@waiting="onAudioWaiting" @error="onAudioError"></audio>
				</div>
			</div>
			<div class="media-empty" v-else-if="info.type === 2">
				<van-icon name="volume-o" size="40" color="#b9bdc2" />
				<span>暂无课件</span>
			</div>

			<div class="chapter-aside">
			<div class="chapter-info">
				<div class="chapter-info-title">{{ info.title || '' }}</div>
				<div class="chapter-info-tip">
					<span>{{ info.course && info.course.title || '' }}</span>
					<div class="chapter-info-enter" @click="goToCourse">
						<span>进入课程</span>
						<van-icon name="arrow" size="10" color="#4285f4" />
					</div>
				</div>
				<div class="chapter-nav" v-if="showChapterNav && (prevChapter || nextChapter)">
					<div class="chapter-nav-btn" :class="{ disabled: !prevChapter }" @click="goNeighbor(prevChapter)">
						<span class="chapter-nav-label">上一节</span>
						<span class="chapter-nav-title">{{ prevChapter ? prevChapter.title : '没有了' }}</span>
					</div>
					<div class="chapter-nav-btn" :class="{ disabled: !nextChapter }" @click="goNeighbor(nextChapter)">
						<span class="chapter-nav-label">{{ nextChapter && nextChapter.un_lock == 0 ? '解锁下一节' : '下一节' }}</span>
						<span class="chapter-nav-title">{{ nextChapter ? nextChapter.title : '没有了' }}</span>
					</div>
				</div>
			</div>

			<div class="study-complete-module" v-if="info.paper_id && examPassed">
				<div class="study-complete-content">
					<div class="study-complete-icon">
						<van-icon name="success" size="24" color="#4285f4" />
					</div>
					<div class="study-complete-text">
						<div class="study-complete-title">学习与测评已完成</div>
						<div class="study-complete-desc">您已完成本章节学习并通过测评，可查看答卷详情</div>
					</div>
				</div>
				<div class="study-complete-btn" @click="goViewExamResult">
					<span>查看测评</span>
					<van-icon name="arrow" size="14" color="#fff" />
				</div>
			</div>
			<div class="study-complete-module study-complete-module--pending" v-else-if="info.paper_id && examReady">
				<div class="study-complete-content">
					<div class="study-complete-icon">
						<van-icon name="clock" size="24" color="#ff3359" />
					</div>
					<div class="study-complete-text">
						<div class="study-complete-title">内容学习完毕</div>
						<div class="study-complete-desc">恭喜您完成本章节学习，快去完成测评作业吧</div>
					</div>
				</div>
				<div class="study-complete-btn" @click="goToExam">
					<span>完成测评作业</span>
					<van-icon name="arrow" size="14" color="#fff" />
				</div>
			</div>
			<div class="study-complete-module study-complete-module--locked" v-else-if="info.paper_id" @click="goToExam">
				<div class="study-complete-content">
					<div class="study-complete-icon">
						<van-icon name="lock" size="24" color="#b9bdc2" />
					</div>
					<div class="study-complete-text">
						<div class="study-complete-title">课后测评</div>
						<div class="study-complete-desc">{{ info.type === 3 ? '读完正文后即可参加测评作业' : '观看至九成后即可参加测评作业' }}</div>
					</div>
				</div>
				<div class="study-complete-btn">
					<span>完成测评作业</span>
					<van-icon name="arrow" size="14" color="#fff" />
				</div>
			</div>
			</div>

			<div class="chapter-main">
			<van-tabs v-if="!isDesktop" v-model:active="chapterTabIndex" color="#4285f4" title-active-color="#353535" title-inactive-color="#7f848c" line-width="22" line-height="3" @change="(i) => switchChapterTab({ index: Number(i) })">
				<van-tab v-for="(t, i) in chapterTabs" :key="i" :title="t.name" :name="i" />
			</van-tabs>
			<div v-else class="chapter-tabs">
				<div :class="{ on: detailTab === 'intro' }" @click="detailTab = 'intro'">介绍</div>
				<div :class="{ on: detailTab === 'comment' }" @click="detailTab = 'comment'">评论（{{ commentTotal }}）</div>
			</div>
			<div class="chapter-content" v-if="showChapterIntro && info.content">
				<div class="content-text">
					<div class="qx-parse" style="user-select:text" v-html="richHtml(info.content)"></div>
				</div>
			</div>
			<div class="media-empty content-empty" v-else-if="showChapterIntro && info.type === 3">
				<van-icon name="description" size="40" color="#b9bdc2" />
				<span>暂无正文</span>
			</div>
			<div class="chapter-intro-empty" v-else-if="showChapterIntro">本节暂无文字介绍</div>
			<div class="text-complete-wrap" v-if="showChapterIntro && info.type === 3 && userChapter.has_end_chapter != 1">
				<div class="text-complete-btn" @click="completeTextChapter">学完本节</div>
			</div>

			<div class="comment-section" v-if="showChapterComment">
				<div class="comment-title">讨论{{ commentTotal > 0 ? '（' + commentTotal + '）' : '' }}</div>
				<div class="comment-list">
					<div class="empty-tip" v-if="commentList.length === 0 && !commentLoading">
						<van-empty description="暂无讨论" image="/static/empty/comment.png" :image-size="80" />
					</div>
					<div class="comment-item" v-for="(item, index) in commentList" :key="item.id || index">
						<div class="comment-header">
							<div class="comment-user">
								<van-image :src="static(item.user && item.user.avatar || '/static/default-avatar.png')" round :width="36 + 'px'" :height="36 + 'px'" fit="cover" />
								<div class="comment-user-info">
									<div class="comment-name"><span>{{ item.user && item.user.nick_name || '匿名' }}</span><coin-badge :coins="item.user ? item.user.coins : null" /></div>
									<div class="comment-time">{{ date(item.created_at, 'yyyy-mm-dd') }}</div>
								</div>
							</div>
							<div class="comment-actions">
								<div class="comment-reply-btn" @click.stop="replyComment(item, item)">
									<van-icon name="chat-o" size="14" color="#7f848c" />
									<span>回复</span>
								</div>
							</div>
						</div>
						<div class="comment-content">{{ item.content || '' }}</div>
						<div class="comment-replies" v-if="item.children && item.children.length > 0">
							<div class="reply-item" v-for="(reply, replyIndex) in item.children"
								:key="reply.id || replyIndex">
								<div class="reply-header">
									<div class="reply-name"><span>{{ reply.user && reply.user.nick_name || '匿名' }}</span><coin-badge :coins="reply.user ? reply.user.coins : null" /></div>
									<span class="reply-time">{{ date(reply.created_at, 'yyyy-mm-dd') }}</span>
									<div class="comment-reply-btn" @click.stop="replyComment(reply, item)">
										<span>回复</span>
									</div>
								</div>
								<div class="reply-content">
									<span class="reply-to" v-if="replyToName(reply, item)">回复 @{{ replyToName(reply, item) }}：</span>
									<span>{{ reply.content || '' }}</span>
								</div>
							</div>
						</div>
					</div>
					<div class="comment-loading" v-if="commentLoading">
						<van-loading type="circular" />
						<span>加载中...</span>
					</div>
					<div class="comment-no-more" v-if="commentHasMore === false && commentList.length > 0">
						<span>没有更多了</span>
					</div>
				</div>
			</div>
			</div>
		</template>

		<div class="chapter-footer" v-if="info && showChapterComment">
			<div class="footer-collect" :class="{ active: isCollect }" @click="toggleCollect">
				<van-icon :name="isCollect ? 'star' : 'star-o'" size="22"
					:color="isCollect ? '#4285f4' : '#404040'" />
				<span>{{ isCollect ? '已收藏' : '收藏' }}</span>
			</div>
			<comment-dock
				v-if="isDesktop"
				ref="commentDock"
				v-model="desktopComment"
				:reply-name="replyUserName || ''"
				:placeholder="replyUserName ? `回复 @${replyUserName}` : '我来说一说…'"
				@send="publishDesktop"
				@cancel="clearReply"
			/>
			<div class="footer-input" v-else @click="showCommentInput">
				<van-icon name="chat-o" size="18" color="#404040" />
				<span>我来说一说…</span>
			</div>
			<div class="footer-next" v-if="showChapterNav && nextChapter" :class="{ locked: nextChapter.un_lock == 0 }" @click="goNeighbor(nextChapter)">
				<span>{{ nextChapter.un_lock == 0 ? '解锁下一节' : '下一节' }}</span>
			</div>
		</div>

		<div class="chapter-switching" v-if="switching">
			<van-loading color="#fff" type="circular" />
			<span>正在打开…</span>
		</div>

		<CommentInputPopup v-if="!isDesktop" :show="showCommentPopup" :placeholder="replyUserName ? `回复 @${replyUserName}：` : '我来说一说…'"
			@close="handleCommentClose" @publish="handleCommentPublish" />
	</div>
</template>

<script>
import { getStorage, removeStorage, setStorage } from '@/common/storage'
import { createAudioContext, createSelectorQuery, pageScrollTo, showModal } from '@/common/dom'
import { navigateBack, navigateTo, setPageTitle } from '@/common/nav'
	import { isDesktop as desktopOn } from '@/pc/desktop'
	import CommentInputPopup from './components/commentInputPopup.vue'
	import CommentDock from './components/CommentDock.vue'
	import CourseVideoPlayer from './components/CourseVideoPlayer.vue'

	const chapterDetailCache = {}
	const warmedMediaUrls = {}

	export default {
		components: {
			CommentInputPopup,
			CommentDock,
			CourseVideoPlayer
		},
		data() {
			return {
				chapterId: null,
				courseId: null,
				info: null,
				userChapter: {
					has_end_chapter: 0,
					is_end: 0
				},
				commentList: [],
				isCollect: false,
				isPlaying: false,
				showCommentPopup: false, // 是否显示评论输入弹窗
				desktopComment: '',
				audioCurrentTime: 0, // 当前播放时间（秒）
				audioDuration: 0, // 音频总时长（秒）
				audioLoading: false, // 音频加载状态
				audioProgressPercent: 0, // 播放进度百分比
				audioProgressTimer: null, // 播放进度更新定时器
				isDragging: false, // 是否正在拖动进度条
				videoCurrentTime: 0, // 视频当前播放时间（秒）
				videoDuration: 0, // 视频总时长（秒）
				videoProgressTimer: null, // 视频播放进度更新定时器
				videoInitialTime: 0, // 视频初始播放位置（从进度恢复）
				commentPage: 1, // 评论列表页码
				commentLimit: 10, // 每页数量
				commentHasMore: true, // 是否还有更多评论
				commentLoading: false, // 评论加载状态
				replyCommentId: null, // 当前回复的评论ID
				replyUserId: null, // 当前回复的用户ID
				replyUserName: null, // 当前回复的用户名
				/** 最近一次通过的章节测评答卷 id（chapterDetail 接口 user_pass_paper_id） */
				userPassPaperId: null,
				pageLoading: false,
				pageError: false,
				prevChapter: null,
				nextChapter: null,
				commentTotal: 0,
				commentPublishing: false,
				detailTab: 'intro',
				audioProgressInited: false,
				needStudyRefresh: false,
				prefetchToken: '',
				detailReqId: 0,
				commentReqId: 0,
				switching: false,
				studyMarking: false,
				textReadTimer: null,
			};
		},
		onLoad(options) {
			if (options.id) {
				this.chapterId = options.id;
			}
			if (options.courseId) {
				this.courseId = options.courseId;
			}
			this.getChapterDetail(false);
		},
		onShow() {
			if (this.needStudyRefresh && this.chapterId) {
				this.needStudyRefresh = false;
				this.getChapterDetail(true);
			}
		},
		onHide() {
			this.saveAudioProgress();
			this.saveVideoProgress();
		},
		onUnload() {
			if (this.textReadTimer) {
				clearTimeout(this.textReadTimer);
				this.textReadTimer = null;
			}
			this.saveAudioProgress();
			this.saveVideoProgress();
			this.pauseAudio();
			this.pauseVideo();
		},
		onPageScroll() {
			this.scheduleTextReadCheck();
		},
		onReachBottom() {
			if (this.info && this.info.type === 3) {
				this.completeTextChapter();
			}
			if (this.commentHasMore && !this.commentLoading && this.info && this.info.id) {
				this.loadMoreComments();
			}
		},
		computed: {
			isDesktop() {
				return desktopOn.value
			},
			showChapterIntro() {
				return this.detailTab === 'intro'
			},
			showChapterComment() {
				return this.detailTab === 'comment'
			},
			chapterTabIndex() {
				return this.detailTab === 'comment' ? 1 : 0
			},
			chapterTabs() {
				return [
					{ name: '介绍' },
					{ name: '评论（' + (Number(this.commentTotal) || 0) + '）' },
				]
			},
			hasReachedStudyThreshold() {
				const uc = this.userChapter || {};
				return Number(uc.has_end_chapter) === 1 || Number(uc.is_end) === 1;
			},
			showChapterNav() {
				return this.hasReachedStudyThreshold;
			},
			examPassed() {
				return !!(this.info && this.info.paper_id && this.userChapter && Number(this.userChapter.is_end) === 1);
			},
			examReady() {
				return this.hasReachedStudyThreshold;
			},
			examWatchTip() {
				return this.info && this.info.type === 3 ? '请先读完正文' : '请先观看至九成';
			}
		},
		methods: {
			cacheKey(id) {
				return String(id || '');
			},
			getCachedChapter(id) {
				return chapterDetailCache[this.cacheKey(id)] || null;
			},
			setCachedChapter(id, data) {
				if (!id || !data) return;
				chapterDetailCache[this.cacheKey(id)] = data;
			},
			syncCachedUserChapter() {
				if (!this.chapterId || !this.userChapter) return;
				const cached = this.getCachedChapter(this.chapterId);
				if (!cached) return;
				cached.userChapter = this.userChapter;
				this.setCachedChapter(this.chapterId, cached);
			},
			resetMediaState() {
				this.stopProgressTimer();
				this.stopVideoProgressTimer();
				this.isPlaying = false;
				this.audioCurrentTime = 0;
				this.audioDuration = 0;
				this.audioLoading = false;
				this.audioProgressPercent = 0;
				this.audioProgressInited = false;
				this.isDragging = false;
				this.videoCurrentTime = 0;
				this.videoDuration = 0;
				this.videoInitialTime = 0;
				this.prefetchToken = '';
				this.studyMarking = false;
				if (this.textReadTimer) {
					clearTimeout(this.textReadTimer);
					this.textReadTimer = null;
				}
			},
			scrollChapterTop() {
				pageScrollTo({
					scrollTop: 0,
					duration: 0
				});
			},
			syncChapterUrl() {
				if (!this.chapterId) return;
				try {
					const url = new URL(window.location.href);
					url.searchParams.set('id', String(this.chapterId));
					if (this.courseId) {
						url.searchParams.set('courseId', String(this.courseId));
					}
					window.history.replaceState(window.history.state, '', url.pathname + url.search + url.hash);
				} catch (e) {}
			},
			applyChapterPayload(data, full) {
				const payload = data || {};
				const nextUc = payload.userChapter || { has_end_chapter: 0, is_end: 0 };
				if (!full && Number(this.userChapter && this.userChapter.has_end_chapter) === 1) {
					nextUc.has_end_chapter = 1;
				}
				this.userChapter = nextUc;
				this.isCollect = payload.isCollect || false;
				this.userPassPaperId = payload.userPassPaperId || null;
				this.prevChapter = payload.prevChapter || null;
				this.nextChapter = payload.nextChapter || null;
				if (!this.courseId && payload.info && payload.info.course_id) {
					this.courseId = payload.info.course_id;
				}
				if (!full) return;
				this.detailTab = 'intro';
				const nextInfo = payload.info || null;
				const resume = nextInfo && Number(nextInfo.type) === 1
					? this.peekVideoProgress(nextInfo.id)
					: { currentTime: 0, duration: 0 };
				this.info = nextInfo;
				if (!this.info) return;
				this.chapterId = this.info.id;
				this.switching = false;
				setPageTitle(this.info.title || '章节详情');
				this.syncChapterUrl();
				this.resetMediaState();
				this.videoInitialTime = resume.currentTime;
				this.videoCurrentTime = resume.currentTime;
				this.videoDuration = resume.duration;
				this.scrollChapterTop();
				this.$nextTick(() => {
					if (this.info && this.info.media_url && this.info.type === 2) {
						this.initAudioProgress();
					}
					if (this.info && this.info.type === 3) {
						this.prefetchNext();
					}
				});
			},
			fetchChapterDetail(refresh = false) {
				const reqId = ++this.detailReqId;
				const reqChapterId = this.chapterId;
				this.$http.get('/course/chapterDetail', {
					id: reqChapterId,
					refresh: refresh ? 1 : 0
				}).then(res => {
					if (reqId !== this.detailReqId) return;
					if (String(this.chapterId) !== String(reqChapterId)) return;
					this.pageLoading = false;
					this.pageError = false;
					this.switching = false;
					const data = res.data || {};
					if (data.info && data.info.id) {
						this.setCachedChapter(data.info.id, data);
					}
					const sameChapter = this.info && data.info && String(this.info.id) === String(data.info.id);
					this.applyChapterPayload(data, !sameChapter);
				}).catch(err => {
					if (reqId !== this.detailReqId) return;
					console.log(err);
					this.pageLoading = false;
					this.switching = false;
					if (!refresh || !this.info) {
						this.pageError = true;
						this.msg('获取章节详情失败');
					}
				});
			},
			// 获取章节详情
			getChapterDetail(refresh = false) {
				if (!this.chapterId) return;
				if (refresh) {
					this.fetchChapterDetail(true);
					return;
				}
				this.pageLoading = true;
				this.pageError = false;
				this.getCommentList(true);
				const cached = this.getCachedChapter(this.chapterId);
				if (cached && cached.info) {
					this.pageLoading = false;
					this.switching = false;
					this.applyChapterPayload(cached, true);
					this.fetchChapterDetail(false);
					return;
				}
				this.fetchChapterDetail(false);
			},
			prefetchNext() {
				const next = this.nextChapter;
				if (!next || !next.id) return;
				const id = this.cacheKey(next.id);
				if (this.prefetchToken === id) return;
				this.prefetchToken = id;
				const cached = this.getCachedChapter(id);
				if (cached && cached.info) {
					this.warmupChapterMedia(cached);
					return;
				}
				this.$http.get('/course/chapterDetail', {
					id: next.id,
					refresh: 1
				}).then(res => {
					const data = res.data || {};
					if (data.info && data.info.id) {
						this.setCachedChapter(data.info.id, data);
						this.warmupChapterMedia(data);
					}
				}).catch(() => {});
			},
			warmupChapterMedia(data) {
				const info = data && data.info;
				if (!info || !info.media_url) return;
				if (info.type !== 1 && info.type !== 2) return;
				const url = this.static(info.media_url);
				if (!url || warmedMediaUrls[url]) return;
				warmedMediaUrls[url] = 1;
				fetch(url, {
					method: 'GET',
					headers: { Range: 'bytes=0-524287' },
					signal: typeof AbortSignal !== 'undefined' && AbortSignal.timeout ? AbortSignal.timeout(15000) : undefined,
				}).catch(() => {})
			},
			// 音频播放控制
			toggleAudio() {
				if (!this.info || !this.info.media_url) {
					this.msg('音频资源不存在');
					return;
				}
				const audioContext = createAudioContext('audioPlayer', this);
				if (this.isPlaying) {
					audioContext.pause();
				} else {
					// 如果有保存的播放进度，先跳转到该位置
					if (this.audioCurrentTime > 0 && this.audioCurrentTime < this.audioDuration) {
						audioContext.seek(this.audioCurrentTime);
					}
					audioContext.play();
				}
			},
			// 暂停音频
			pauseAudio() {
				if (this.isPlaying) {
					const audioContext = createAudioContext('audioPlayer', this);
					audioContext.pause();
				}
			},
			// 音频播放事件
			onAudioPlay() {
				this.isPlaying = true;
				this.audioLoading = false;
				// 开始记录播放进度
				this.startProgressTimer();
			},
			// 音频暂停事件
			onAudioPause() {
				this.isPlaying = false;
				// 停止进度更新
				this.stopProgressTimer();
				// 保存播放进度
				this.saveAudioProgress();
			},
			// 音频播放进度更新
			onAudioTimeUpdate(e) {
				// 如果正在拖动，不更新进度条
				if (this.isDragging) return;
				if (e.detail && e.detail.currentTime !== undefined) {
					this.audioCurrentTime = e.detail.currentTime;
					if (this.audioDuration > 0) {
						this.audioProgressPercent = (this.audioCurrentTime / this.audioDuration) * 100;
					}
					// 如果音频时长还未设置，尝试从事件中获取
					if (!this.audioDuration && e.detail.duration) {
						this.audioDuration = e.detail.duration;
					}
					this.tryMarkStudyByPlaybackProgress(this.audioCurrentTime, this.audioDuration);
				}
			},
			// 进度条拖动中
			onProgressChanging(e) {
				this.isDragging = true;
				const percent = e.detail.value;
				if (this.audioDuration > 0) {
					const newTime = (percent / 100) * this.audioDuration;
					this.audioCurrentTime = newTime;
					this.audioProgressPercent = percent;
				}
			},
			// 进度条拖动结束
			onProgressChange(e) {
				this.isDragging = false;
				const percent = e.detail.value;
				if (this.audioDuration > 0) {
					const seekTime = (percent / 100) * this.audioDuration;
					this.audioCurrentTime = seekTime;
					this.audioProgressPercent = percent;
					// 跳转到指定位置
					const audioContext = createAudioContext('audioPlayer', this);
					audioContext.seek(seekTime);
					this.tryMarkStudyByPlaybackProgress(this.audioCurrentTime, this.audioDuration);
				}
			},
			// 音频播放完成
			onAudioEnded() {
				this.isPlaying = false;
				this.stopProgressTimer();
				this.clearAudioProgress();
				this.markStudyComplete();
			},
			// 音频加载中
			onAudioWaiting() {
				this.audioLoading = true;
			},
			// 音频可以播放
			onAudioCanPlay(e) {
				this.audioLoading = false;
				if (e.detail && e.detail.duration) {
					this.audioDuration = e.detail.duration;
				} else if (this.info && this.info.media && this.info.media.duration) {
					// 如果事件中没有时长，使用媒体信息中的时长
					this.audioDuration = this.info.media.duration;
				}
				// 恢复播放进度（延迟一下确保时长已设置）
				this.$nextTick(() => {
					this.restoreAudioProgress();
				});
				this.prefetchNext();
			},
			// 音频错误
			onAudioError(e) {
				this.audioLoading = false;
				this.isPlaying = false;
				console.error('音频播放错误:', e);
				this.msg('音频加载失败，请稍后重试');
			},
			// 开始进度更新定时器
			startProgressTimer() {
				this.stopProgressTimer();
				this.audioProgressTimer = setInterval(() => {
					// 每5秒保存一次播放进度
					if (this.isPlaying && this.audioCurrentTime > 0) {
						this.saveAudioProgress();
					}
				}, 5000);
			},
			// 停止进度更新定时器
			stopProgressTimer() {
				if (this.audioProgressTimer) {
					clearInterval(this.audioProgressTimer);
					this.audioProgressTimer = null;
				}
			},
			// 格式化时间（秒转 mm:ss）
			formatTime(seconds) {
				if (!seconds || isNaN(seconds)) return '00:00';
				const mins = Math.floor(seconds / 60);
				const secs = Math.floor(seconds % 60);
				return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
			},
			// 保存播放进度
			saveAudioProgress() {
				if (!this.chapterId) return;
				// 确保有当前时间
				if (!this.audioCurrentTime || this.audioCurrentTime <= 0) return;
				// 确保有音频时长
				if (!this.audioDuration || this.audioDuration <= 0) {
					// 如果还没有时长，尝试从媒体信息获取
					if (this.info && this.info.media && this.info.media.duration) {
						this.audioDuration = this.info.media.duration;
					} else {
						return;
					}
				}
				// 如果播放进度超过95%，认为已经听完，不保存
				if ((this.audioCurrentTime / this.audioDuration) >= 0.95) {
					this.clearAudioProgress();
					return;
				}
				// 保存到本地存储
				const storageKey = `chapter_progress_${this.chapterId}`;
				try {
					setStorage(storageKey, {
						currentTime: this.audioCurrentTime,
						duration: this.audioDuration,
						updateTime: Date.now()
					});
				} catch (e) {
					// 保存失败，静默处理
				}
			},
			// 恢复播放进度
			restoreAudioProgress() {
				if (!this.chapterId) return;
				// 确保有音频时长
				if (!this.audioDuration || this.audioDuration <= 0) {
					// 如果还没有时长，尝试从媒体信息获取
					if (this.info && this.info.media && this.info.media.duration) {
						this.audioDuration = this.info.media.duration;
					} else {
						// 如果还是没有，延迟重试
						setTimeout(() => {
							this.restoreAudioProgress();
						}, 500);
						return;
					}
				}
				// 从本地存储获取播放进度
				const storageKey = `chapter_progress_${this.chapterId}`;
				let savedProgress = null;
				try {
					savedProgress = getStorage(storageKey);
				} catch (e) {
					// 读取失败，静默处理
				}
				if (savedProgress && savedProgress.currentTime && savedProgress.duration) {
					const { currentTime, duration } = savedProgress;
					// 如果保存的进度小于总时长的95%，则恢复播放位置
					// 如果已经播放到95%以上，认为已经听完，不恢复
					if (currentTime > 0 && duration > 0 && (currentTime / duration) < 0.95) {
						// 如果当前时长和保存的时长不一致，使用比例来恢复进度
						if (this.audioDuration !== duration) {
							const progressPercent = (currentTime / duration) * 100;
							this.audioCurrentTime = (progressPercent / 100) * this.audioDuration;
						} else {
							this.audioCurrentTime = currentTime;
						}
						this.audioProgressPercent = (this.audioCurrentTime / this.audioDuration) * 100;
						// 跳转到保存的播放位置（延迟执行确保音频已加载）
						setTimeout(() => {
							try {
								const audioContext = createAudioContext('audioPlayer', this);
								if (audioContext && typeof audioContext.seek === 'function') {
									audioContext.seek(this.audioCurrentTime);
								}
							} catch (e) {
								// 跳转失败，静默处理
							}
						}, 500);
					} else {
						// 播放完成，清除进度
						this.clearAudioProgress();
					}
				} else if (this.userChapter && this.userChapter.progress && this.audioDuration > 0) {
					// 如果本地存储没有，尝试从用户章节信息中获取播放进度
					this.audioCurrentTime = (this.userChapter.progress / 100) * this.audioDuration;
					this.audioProgressPercent = this.userChapter.progress;
					// 跳转到保存的播放位置（延迟执行确保音频已加载）
					setTimeout(() => {
						try {
							const audioContext = createAudioContext('audioPlayer', this);
							if (audioContext && typeof audioContext.seek === 'function') {
								audioContext.seek(this.audioCurrentTime);
							}
						} catch (e) {
							console.error('跳转播放位置失败:', e);
						}
					}, 500);
				}
			},
			// 清除音频播放进度
			clearAudioProgress() {
				if (!this.chapterId) return;
				const storageKey = `chapter_progress_${this.chapterId}`;
				removeStorage(storageKey);
			},
			/**
			 * 播放进度达到 90% 即上报学完；测评入口只看 has_end_chapter / is_end
			 */
			tryMarkStudyByPlaybackProgress(currentTime, duration) {
				if (!duration || duration <= 0 || currentTime == null || isNaN(currentTime)) return;
				if (currentTime / duration < 0.9) return;
				this.markStudyComplete();
			},
			// 标记学习完成
			markStudyComplete() {
				if (!this.chapterId || this.studyMarking) return;
				if (!this.info || !this.userChapter) return;
				if (this.userChapter.is_end == 1 || this.userChapter.has_end_chapter == 1) return;
				if (Number(this.userChapter.un_lock) !== 1) return;
				this.studyMarking = true;
				this.$http.post('/learn/learnEnd', {
					chapter_id: this.chapterId
				}).then(res => {
					const next = res.data || this.userChapter || {};
					next.has_end_chapter = 1;
					this.userChapter = next;
					this.syncCachedUserChapter();
					if (Number(this.userChapter.is_end) === 1) {
						this.offerNextChapter();
					}
				}).catch(() => {
				}).finally(() => {
					this.studyMarking = false;
				});
			},
			scheduleTextReadCheck() {
				if (!this.info || this.info.type !== 3) return;
				if (!this.userChapter || this.userChapter.has_end_chapter == 1) return;
				if (this.textReadTimer) clearTimeout(this.textReadTimer);
				this.textReadTimer = setTimeout(() => {
					this.textReadTimer = null;
					this.checkTextChapterRead();
				}, 200);
			},
			checkTextChapterRead() {
				if (!this.info || this.info.type !== 3 || !this.info.content) return;
				if (!this.userChapter || this.userChapter.has_end_chapter == 1) return;
				createSelectorQuery().in(this).select('.chapter-content').boundingClientRect((rect) => {
					if (!rect) return;
					const viewportH = (typeof window !== 'undefined' && window.innerHeight) ? window.innerHeight : 0;
					if (rect.bottom <= viewportH + 80) {
						this.completeTextChapter();
					}
				}).exec();
			},
			completeTextChapter() {
				if (!this.info || this.info.type !== 3) return;
				this.markStudyComplete();
			},
			// 初始化音频播放进度
			initAudioProgress() {
				if (this.audioProgressInited) return;
				this.audioProgressInited = true;
				if (this.info && this.info.media && this.info.media.duration) {
					this.audioDuration = this.info.media.duration;
				}
				this.$nextTick(() => {
					this.restoreAudioProgress();
				});
			},
			// 视频播放事件
			onVideoPlay() {
				// 开始记录播放进度
				this.startVideoProgressTimer();
			},
			onVideoCanPlay() {
				this.prefetchNext();
			},
			onVideoPause() {
				// 停止进度更新
				this.stopVideoProgressTimer();
				// 保存播放进度
				this.saveVideoProgress();
			},
			onVideoEnded() {
				this.stopVideoProgressTimer();
				this.clearVideoProgress();
				this.markStudyComplete();
			},
			// 视频播放进度更新（由 CourseVideoPlayer 的 timeupdate 传入）
			onVideoTimeUpdate(payload) {
				if (payload && payload.currentTime !== undefined) this.videoCurrentTime = payload.currentTime;
				if (payload && payload.duration !== undefined) this.videoDuration = payload.duration;
				const cur = payload && payload.currentTime !== undefined ? payload.currentTime : this.videoCurrentTime;
				const dur = payload && payload.duration !== undefined ? payload.duration : this.videoDuration;
				this.tryMarkStudyByPlaybackProgress(cur, dur);
			},
			// 暂停视频
			pauseVideo() {
				if (this.info && this.info.type === 1 && this.$refs.courseVideo) {
					this.$refs.courseVideo.pause();
				}
			},
			// 开始视频进度更新定时器
			startVideoProgressTimer() {
				this.stopVideoProgressTimer();
				this.videoProgressTimer = setInterval(() => {
					// 每5秒保存一次播放进度
					if (this.videoCurrentTime > 0) {
						this.saveVideoProgress();
					}
				}, 5000);
			},
			// 停止视频进度更新定时器
			stopVideoProgressTimer() {
				if (this.videoProgressTimer) {
					clearInterval(this.videoProgressTimer);
					this.videoProgressTimer = null;
				}
			},
			// 初始化视频播放进度（从存储读取，传给 CourseVideoPlayer 的 initialTime）
			peekVideoProgress(chapterId) {
				if (!chapterId) return { currentTime: 0, duration: 0 };
				const savedProgress = getStorage(`chapter_progress_${chapterId}`);
				if (savedProgress && savedProgress.currentTime != null && savedProgress.duration) {
					const { currentTime, duration } = savedProgress;
					if (currentTime > 0 && duration > 0 && (currentTime / duration) < 0.95) {
						return { currentTime, duration };
					}
				}
				return { currentTime: 0, duration: 0 };
			},
			// 保存视频播放进度
			saveVideoProgress() {
				if (!this.chapterId || !this.videoCurrentTime || this.videoDuration <= 0) return;
				// 如果播放进度超过95%，认为已经看完，不保存
				if ((this.videoCurrentTime / this.videoDuration) >= 0.95) {
					this.clearVideoProgress();
					return;
				}
				// 保存到本地存储
				const storageKey = `chapter_progress_${this.chapterId}`;
				setStorage(storageKey, {
					currentTime: this.videoCurrentTime,
					duration: this.videoDuration,
					updateTime: Date.now()
				});
			},
			// 清除视频播放进度
			clearVideoProgress() {
				if (!this.chapterId) return;
				const storageKey = `chapter_progress_${this.chapterId}`;
				removeStorage(storageKey);
			},
			// 跳转到测评作业页面
			goToExam() {
				if (!this.hasReachedStudyThreshold) {
					this.msg(this.examWatchTip + '再参加测评');
					return;
				}
				if (!this.courseId && !this.chapterId) {
					this.msg('参数错误');
					return;
				}
				this.needStudyRefresh = true;
				// 跳转到测评页面，传递课程ID和章节ID
				const examUrl = `/pages/exam/exam?id=${this.info.paper_id || ''}`;
				navigateTo(examUrl);
			},
			/** 已通过后查看答卷（与 exam 页「查看详情」一致） */
			goViewExamResult() {
				const id = this.userPassPaperId;
				if (!id) {
					this.msg('暂无测评记录，请下拉刷新重试');
					return;
				}
				this.needStudyRefresh = true;
				navigateTo(`/pages/exam/examDetail?id=${id}`);
			},
			// 收藏/取消收藏
			toggleCollect() {
				if (!this.chapterId) return;
				this.$http.post('/learn/collectChapter', {
					chapter_id: this.chapterId
				}).then(res => {
					this.isCollect = !this.isCollect;
					this.msg(this.isCollect ? '收藏成功' : '取消收藏');
				}).catch(err => {
					console.log(err);
					this.msg(err.message || '操作失败');
				});
			},
			// 点赞/取消点赞
			toggleLike(item) {
				// 实现点赞逻辑
			},
			// 查看回复
			viewReplies(item) {
				// 跳转到回复页面
			},
			switchChapterTab(e) {
				const index = typeof e === 'object' && e !== null ? e.index : e
				this.detailTab = Number(index) === 1 ? 'comment' : 'intro'
			},
			// 显示评论输入
			showCommentInput() {
				if (this.isDesktop) {
					this.detailTab = 'comment';
					this.focusDesktop();
					return;
				}
				this.showCommentPopup = true;
			},
			focusDesktop() {
				this.$nextTick(() => {
					const dock = this.$refs.commentDock
					if (dock && dock.focus) dock.focus()
				})
			},
			publishDesktop() {
				const text = String(this.desktopComment || '').trim()
				if (!text) {
					this.msg('请输入评论内容')
					return
				}
				this.handleCommentPublish(text)
			},
			clearReply() {
				this.replyCommentId = null
				this.replyUserId = null
				this.replyUserName = null
			},
			// 关闭评论输入弹窗
			handleCommentClose() {
				this.showCommentPopup = false;
				this.desktopComment = '';
				this.clearReply();
			},
			// 发布评论
			handleCommentPublish(content) {
				if (!this.chapterId) {
					this.msg('参数错误');
					return;
				}
				if (this.commentPublishing) return;
				this.commentPublishing = true;
				const params = {
					chapter_id: this.chapterId,
					course_id: this.courseId,
					content: content,
				};

				if (this.replyCommentId) {
					params.pid = this.replyCommentId;
				} else {
					params.pid = 0;
				}

				this.$http.post('/comment/chapterComment', params).then(res => {
					this.msg('发布成功');
					this.handleCommentClose();
					this.getCommentList(true);
				}).catch(err => {
					console.log(err);
					this.msg(err.message || '发布失败');
				}).finally(() => {
					this.commentPublishing = false;
				});
			},
			// 获取评论列表
			getCommentList(reset = true) {
				if (!this.chapterId) return;
				if (reset) {
					this.commentPage = 1;
					this.commentList = [];
					this.commentHasMore = true;
					this.commentLoading = false;
					this.commentTotal = 0;
				}
				if (this.commentLoading || !this.commentHasMore) return;

				this.commentLoading = true;
				const reqId = ++this.commentReqId;
				const reqChapterId = this.chapterId;
				this.$http.post('/comment/getChapterCommentList', {
					chapter_id: reqChapterId,
					course_id: this.courseId,
					page: this.commentPage,
					limit: this.commentLimit
				}).then(res => {
					if (reqId !== this.commentReqId) return;
					this.commentLoading = false;
					this.commentList = [...this.commentList, ...res.data.list];
					this.commentHasMore = res.data.hasMore;
					if (typeof res.data.total === 'number') {
						this.commentTotal = res.data.total;
					}
				}).catch(err => {
					if (reqId !== this.commentReqId) return;
					this.commentLoading = false;
					console.log(err);
					this.msg('加载评论失败');
				});
			},
			// 加载更多评论
			loadMoreComments() {
				this.commentPage++;
				this.getCommentList(false);
			},
			// 回复评论
			replyComment(item, rootItem) {
				const root = rootItem || item;
				this.replyCommentId = root.id;
				this.replyUserId = item.user && item.user.id;
				this.replyUserName = (item.user && item.user.nick_name) || '匿名';
				if (this.isDesktop) {
					this.detailTab = 'comment';
					this.focusDesktop();
					return;
				}
				this.showCommentPopup = true;
			},
			replyToName(reply, root) {
				if (reply.parent && reply.parent.user && reply.parent.user.nick_name) {
					return reply.parent.user.nick_name;
				}
				return (root && root.user && root.user.nick_name) || '';
			},
			goNeighbor(chapter) {
				if (!chapter || !chapter.id) return;
				if (chapter.un_lock == 0) {
					if (Number(this.userChapter && this.userChapter.is_end) !== 1) {
						this.msg(this.info && this.info.paper_id ? '请先完成测评' : '请先完成本节学习');
						return;
					}
					showModal({
						title: '解锁下一节',
						content: chapter.title ? `「${chapter.title}」尚未解锁，是否现在解锁？` : '该章节尚未解锁，是否现在解锁？',
						confirmText: '解锁',
						success: (res) => {
							if (!res.confirm) return;
							this.$http.post('/learn/unlockChapter', {
								chapter_id: chapter.id
							}).then(() => {
								this.msg('解锁成功');
								chapter.un_lock = 1;
								this.enterChapter(chapter);
							}).catch((err) => {
								const message = (err && err.message) || '';
								if (message.indexOf('已经解锁') !== -1) {
									chapter.un_lock = 1;
									this.enterChapter(chapter);
									return;
								}
								this.msg(message || '请先完成上一章节');
							});
						}
					});
					return;
				}
				this.enterChapter(chapter);
			},
			enterChapter(chapter) {
				if (!chapter || !chapter.id) return;
				if (String(chapter.id) === String(this.chapterId) && this.info) return;
				this.pauseVideo();
				this.pauseAudio();
				this.saveVideoProgress();
				this.saveAudioProgress();
				this.chapterId = chapter.id;
				this.scrollChapterTop();
				this.syncChapterUrl();
				this.switching = !this.getCachedChapter(chapter.id);
				this.getChapterDetail(false);
			},
			offerNextChapter() {
				if (!this.nextChapter || !this.nextChapter.id) return;
				if (Number(this.userChapter && this.userChapter.is_end) !== 1) return;
				const locked = this.nextChapter.un_lock == 0;
				showModal({
					title: '本节已学完',
					content: locked
						? `下一节「${this.nextChapter.title}」尚未解锁，是否解锁并学习？`
						: `是否进入下一节「${this.nextChapter.title}」？`,
					confirmText: locked ? '解锁学习' : '下一节',
					success: (res) => {
						if (res.confirm) this.goNeighbor(this.nextChapter);
					}
				});
			},
			// 跳转到课程
			goToCourse() {
				if (!this.courseId) return;
				const pages = getCurrentPages();
				const prev = pages.length > 1 ? pages[pages.length - 2] : null;
				const route = (prev && prev.route) || '';
				if (route.indexOf('pages/course/courseDetail') !== -1) {
					navigateBack();
					return;
				}
				this.navTo(`/pages/course/courseDetail?id=${this.courseId}&tab=1`);
			},
			// 显示时长（秒转 mm:ss）
			showDuration(seconds) {
				if (!seconds || isNaN(seconds)) return '00:00';
				return this.formatTime(seconds);
			},
			// 显示文件大小
			showSize(bytes) {
				if (!bytes || isNaN(bytes)) return '0B';
				if (bytes < 1024) {
					return bytes + 'B';
				} else if (bytes < 1024 * 1024) {
					return (bytes / 1024).toFixed(2) + 'KB';
				} else if (bytes < 1024 * 1024 * 1024) {
					return (bytes / (1024 * 1024)).toFixed(2) + 'MB';
				} else {
					return (bytes / (1024 * 1024 * 1024)).toFixed(2) + 'GB';
				}
			}
		}
	};
</script>

<style lang="scss" scoped>
	.chapter-detail {
		padding-bottom: calc(calc(140 * var(--rpx)) + env(safe-area-inset-bottom));
		background: #fff;
		min-height: 100vh;
	}

	.chapter-skeleton {
		padding: 0 0 calc(40 * var(--rpx));
	}

	.skeleton-player {
		width: 100%;
		height: calc(422 * var(--rpx));
		background: #1a1a1a;
	}

	.skeleton-line {
		height: calc(28 * var(--rpx));
		margin: calc(28 * var(--rpx)) calc(40 * var(--rpx)) 0;
		border-radius: calc(8 * var(--rpx));
		background: #f0f2f5;
	}

	.skeleton-title {
		width: 70%;
		height: calc(40 * var(--rpx));
	}

	.skeleton-short {
		width: 40%;
	}

	.chapter-load-fail {
		display: flex;
		align-items: center;
		justify-content: center;
		height: calc(422 * var(--rpx));
		background: #111;
		color: #fff;
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	}

	.chapter-switching {
		position: fixed;
		left: 0;
		right: 0;
		top: 0;
		bottom: 0;
		z-index: 200;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: calc(16 * var(--rpx));
		background: rgba(0, 0, 0, 0.45);
		color: #fff;
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	}

	.media-empty {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: calc(422 * var(--rpx));
		background: #111;
		gap: calc(16 * var(--rpx));
		color: #b9bdc2;
		font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
		margin-bottom: calc(24 * var(--rpx));
	}

	.chapter-nav {
		display: flex;
		gap: calc(20 * var(--rpx));
		margin-top: calc(24 * var(--rpx));
	}

	.chapter-nav-btn {
		flex: 1;
		min-width: 0;
		padding: calc(12 * var(--rpx)) calc(20 * var(--rpx));
		border-radius: calc(16 * var(--rpx));
		background: #f6f7fb;
		color: #4285f4;
		text-align: left;
	}

	.chapter-nav-label {
		display: block;
		font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
		color: #7f848c;
		line-height: 1.4;
	}

	.chapter-nav-title {
		display: block;
		margin-top: calc(4 * var(--rpx));
		font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
		color: #404040;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.chapter-nav-btn.disabled {
		opacity: 0.7;
	}

	.chapter-nav-btn.disabled .chapter-nav-title {
		color: #b9bdc2;
	}

	.chapter-header {
		padding: calc(28 * var(--rpx)) calc(40 * var(--rpx)) calc(20 * var(--rpx));
	}

	.chapter-title {
		font-size: calc(calc(44 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 600;
		color: #353535;
		line-height: 1.4;
		margin-bottom: calc(16 * var(--rpx));
	}

	.chapter-subtitle {
		display: flex;
		align-items: center;
		justify-content: space-between;
		font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
		color: #7f848c;
	}

	.chapter-enter-btn {
		display: flex;
		align-items: center;
		padding: calc(8 * var(--rpx)) calc(16 * var(--rpx));
		background: #f3f8ff;
		border-radius: calc(26 * var(--rpx));
		color: #4285f4;
		font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	}

	.chapter-cover {
		width: 100%;
		margin-bottom: calc(40 * var(--rpx));
	}

	.cover-img {
		width: 100%;
		display: block;
	}

	.audio-player {
		position: relative;
		display: flex;
		align-items: flex-start;
		padding: calc(32 * var(--rpx));
		margin: 0 calc(32 * var(--rpx)) calc(40 * var(--rpx));
		border-radius: calc(24 * var(--rpx));
		box-shadow: 0 calc(8 * var(--rpx)) calc(24 * var(--rpx)) rgba(66, 133, 244, 0.12);
		overflow: hidden;
		min-height: calc(180 * var(--rpx));
	}

	.audio-player-bg {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-size: cover;
		background-position: center;
		background-repeat: no-repeat;
		filter: blur(calc(0 * var(--rpx)));
		transform: scale(1.1);
		opacity: 1;
		z-index: 0;
	}

	.audio-player-content {
		position: relative;
		z-index: 1;
		display: flex;
		align-items: flex-start;
		width: 100%;
		background: rgba(255, 255, 255, 0.6);
		border-radius: calc(20 * var(--rpx));
		padding: calc(32 * var(--rpx));
		margin: calc(-32 * var(--rpx));
		backdrop-filter: blur(calc(5 * var(--rpx)));
	}

	/* 专辑封面 */
	.audio-cover-wrapper {
		position: relative;
		width: calc(120 * var(--rpx));
		height: calc(120 * var(--rpx));
		flex-shrink: 0;
		margin-right: calc(24 * var(--rpx));
	}

	.audio-cover-container {
		width: calc(120 * var(--rpx));
		height: calc(120 * var(--rpx));
		border-radius: 50%;
		overflow: hidden;
		box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(0, 0, 0, 0.15);
		transition: transform 0.3s ease;
	}

	.audio-cover-rotating {
		animation: rotate 20s linear infinite;
	}

	@keyframes rotate {
		from {
			transform: rotate(0deg);
		}
		to {
			transform: rotate(360deg);
		}
	}

	.audio-cover-img {
		width: 100%;
		height: 100%;
	}

	.audio-cover-mask {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.audio-play-icon {
		opacity: 0.95;
	}

	.audio-loading {
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
	}

	/* 播放内容区域 */
	.audio-content {
		flex: 1;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		min-width: 0;
		overflow: hidden;
		width: 0;
	}

	.audio-title-row {
		margin-bottom: calc(20 * var(--rpx));
	}

	.audio-title {
		font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 600;
		color: #1a1a1a;
		margin-bottom: calc(8 * var(--rpx));
		line-height: 1.4;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.audio-subtitle {
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #7f848c;
		line-height: 1.4;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	/* 进度条 */
	.audio-progress-wrapper {
		margin: calc(20 * var(--rpx)) 0 calc(12 * var(--rpx));
		width: 100%;
		overflow: hidden;
		position: relative;
	}

	.audio-slider {
		width: 100%;
		box-sizing: border-box;
		margin: 0;
		padding: 0;
		flex-shrink: 0;
	}

	/* 时间显示 */
	.audio-time-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
		color: #999;
	}

	.audio-time-current {
		font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
		color: #666;
	}

	.audio-time-total {
		font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
		color: #999;
	}

	.video-player {
		width: 100%;
		margin-bottom: calc(40 * var(--rpx));
	}

	.chapter-aside {
		display: block;
	}

	.van-tabs {
		background: #fff;
		padding: 15px 46px 7px;
	}

	.van-tabs :deep(.van-tab) {
		font-size: calc(17px * var(--font-scale, 1));
	}

	.chapter-tabs {
		display: flex;
		gap: 4px;
		padding: 0 12px;
		background: #fff;
		border-bottom: 1px solid #f0f0f0;
	}

	.chapter-tabs > div {
		padding: 14px 16px 12px;
		border-bottom: 2px solid transparent;
		color: #6b7280;
		font-size: calc(15px * var(--font-scale, 1));
	}

	.chapter-tabs > div.on {
		color: #4285f4;
		font-weight: 650;
		border-bottom-color: #4285f4;
	}

	.chapter-main {
		display: block;
	}

	.chapter-info {
		padding: calc(32 * var(--rpx)) calc(40 * var(--rpx));
		border-bottom: calc(1 * var(--rpx)) solid #f0f0f0;
	}

	/* 学习完成模块样式 */
	.study-complete-module {
		margin: calc(40 * var(--rpx)) calc(40 * var(--rpx));
		padding: calc(32 * var(--rpx));
		background: linear-gradient(135deg, #f0f7ff 0%, #e8f4ff 100%);
		border-radius: calc(16 * var(--rpx));
		border: calc(1 * var(--rpx)) solid #d0e7ff;
	}

	.study-complete-content {
		display: flex;
		align-items: flex-start;
		margin-bottom: calc(24 * var(--rpx));
	}

	.study-complete-icon {
		margin-right: calc(20 * var(--rpx));
		flex-shrink: 0;
	}

	.study-complete-text {
		flex: 1;
	}

	.study-complete-title {
		font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 500;
		color: #404040;
		margin-bottom: calc(8 * var(--rpx));
	}

	.study-complete-desc {
		font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
		color: #7f848c;
		line-height: 1.5;
	}

	.study-complete-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		padding: calc(20 * var(--rpx)) 0;
		background: linear-gradient(135deg, #5c86ff, #8bbeff);
		border-radius: calc(84 * var(--rpx));
		color: #fff;
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 500;
	}

	.study-complete-btn span {
		margin-right: calc(8 * var(--rpx));
	}

	/* 待测评：整体主色 #ff3359 */
	.study-complete-module--pending {
		background: linear-gradient(135deg, #fff5f7 0%, #ffe8ee 100%);
		border-color: #ffc9d6;
	}

	.study-complete-module--pending .study-complete-btn {
		background: linear-gradient(135deg, #ff3359, #ff6b82);
	}

	.study-complete-module--locked {
		background: #f5f6f8;
		border-color: #e6e8eb;
	}

	.study-complete-module--locked .study-complete-title,
	.study-complete-module--locked .study-complete-desc {
		color: #9aa0a6;
	}

	.study-complete-module--locked .study-complete-btn {
		background: #c8ccd2;
	}

	.chapter-info-title {
		font-size: calc(calc(40 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 500;
		color: #404040;
		margin-bottom: calc(20 * var(--rpx));
	}

	.chapter-info-subtitle {
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #7f848c;
		margin-bottom: calc(12 * var(--rpx));
	}

	.chapter-info-tip {
		display: flex;
		align-items: center;
		justify-content: space-between;
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #7f848c;
	}

	.chapter-info-enter {
		display: flex;
		align-items: center;
		color: #4285f4;
	}

	.chapter-intro-empty {
		padding: 48px 24px;
		color: #9ca3af;
		font-size: 14px;
		text-align: center;
	}

	.chapter-content {
		padding: calc(40 * var(--rpx));
	}

	.content-title {
		font-size: calc(calc(36 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 500;
		color: #404040;
		margin-bottom: calc(24 * var(--rpx));
	}

	.content-text {
		font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
		line-height: 1.8;
		color: #404040;
		word-break: break-word;
	}

	.content-empty {
		height: calc(280 * var(--rpx));
		background: #f7f8fa;
		margin: 0 0 calc(8 * var(--rpx));
	}

	.text-complete-wrap {
		padding: calc(8 * var(--rpx)) calc(40 * var(--rpx)) calc(24 * var(--rpx));
	}

	.text-complete-btn {
		height: calc(80 * var(--rpx));
		border-radius: calc(40 * var(--rpx));
		display: flex;
		align-items: center;
		justify-content: center;
		background: linear-gradient(135deg, #5c86ff, #8bbeff);
		color: #fff;
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 500;
	}

	.comment-section {
		padding: calc(40 * var(--rpx)) 0;
		margin-top: calc(40 * var(--rpx));
		border-top: calc(20 * var(--rpx)) solid #f6f7fb;
	}

	.comment-title {
		font-size: calc(calc(36 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 500;
		color: #404040;
		padding: 0 calc(40 * var(--rpx)) calc(20 * var(--rpx));
		margin-bottom: calc(40 * var(--rpx));
	}

	.comment-list {
		padding: 0 calc(40 * var(--rpx));
	}

	.empty-tip {
		text-align: center;
		padding: calc(48 * var(--rpx)) 0 calc(24 * var(--rpx));
		color: #999;
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	}

	.comment-item {
		padding-bottom: calc(32 * var(--rpx));
		margin-bottom: calc(38 * var(--rpx));
		border-bottom: calc(1 * var(--rpx)) solid #f0f0f0;
	}

	.comment-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: calc(20 * var(--rpx));
	}

	.comment-user {
		display: flex;
		align-items: center;
		flex: 1;
	}

	.comment-avatar {
		width: calc(70 * var(--rpx));
		height: calc(70 * var(--rpx));
		border-radius: 50%;
		margin-right: calc(16 * var(--rpx));
	}

	.comment-user-info {
		padding-left: calc(16 * var(--rpx));
		flex: 1;
	}

	.comment-name {
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 500;
		color: #404040;
		margin-bottom: calc(8 * var(--rpx));
	}

	.comment-time {
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #7f848c;
		transform: scale(0.83);
		transform-origin: 0;
	}

	.comment-like {
		padding: calc(8 * var(--rpx));
	}

	.comment-content {
		padding-left: calc(96 * var(--rpx));
		font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
		line-height: 1.5;
		color: #404040;
		word-break: break-word;
		margin-bottom: calc(16 * var(--rpx));
	}

	.comment-actions {
		display: flex;
		align-items: center;
		gap: calc(20 * var(--rpx));
	}

	.comment-reply-btn {
		display: flex;
		align-items: center;
		gap: calc(6 * var(--rpx));
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #7f848c;
	}

	.comment-replies {
		margin-left: calc(96 * var(--rpx));
		margin-top: calc(20 * var(--rpx));
		padding: calc(20 * var(--rpx));
		background: #f5f5f5;
		border-radius: calc(8 * var(--rpx));
	}

	.reply-item {
		margin-bottom: calc(16 * var(--rpx));
	}

	.reply-item:last-child {
		margin-bottom: 0;
	}

	.reply-header {
		display: flex;
		align-items: center;
		gap: calc(12 * var(--rpx));
		margin-bottom: calc(8 * var(--rpx));
	}

	.reply-header .comment-reply-btn {
		margin-left: auto;
		font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	}

	.reply-name {
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #4285f4;
		font-weight: 500;
	}

	.reply-time {
		font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
		color: #b9bdc2;
	}

	.reply-content {
		font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
		color: #404040;
		line-height: 1.6;
	}

	.reply-to {
		color: #4285f4;
	}

	.comment-reply-more {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: calc(8 * var(--rpx));
		padding: calc(16 * var(--rpx)) 0;
		margin-top: calc(16 * var(--rpx));
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #7f848c;
	}

	.comment-loading {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: calc(12 * var(--rpx));
		padding: calc(40 * var(--rpx)) 0;
		font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
		color: #7f848c;
	}

	.comment-no-more {
		text-align: center;
		padding: calc(40 * var(--rpx)) 0;
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #b9bdc2;
	}

	.chapter-footer {
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
		z-index: 100;
		display: flex;
		align-items: center;
		min-height: calc(108 * var(--rpx));
		padding-right: calc(24 * var(--rpx));
		padding-bottom: env(safe-area-inset-bottom);
		background: #fff;
		box-shadow: 0 0 calc(30 * var(--rpx)) rgba(233, 234, 237, 0.65);
	}

	.footer-collect {
		display: flex;
		flex-direction: column;
		align-items: center;
		width: calc(124 * var(--rpx));
		padding-top: calc(2 * var(--rpx));
	}

	.footer-collect span {
		margin-top: calc(6 * var(--rpx));
		font-size: calc(calc(18 * var(--rpx)) * var(--font-scale, 1));
		color: #404040;
	}

	.footer-input {
		flex: 1;
		display: flex;
		align-items: center;
		height: calc(72 * var(--rpx));
		padding-left: calc(30 * var(--rpx));
		background: #f6f7fb;
		border-radius: calc(36 * var(--rpx));
		color: #b9bdc2;
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	}

	.footer-input span {
		margin-left: calc(16 * var(--rpx));
	}

	.footer-next {
		margin-left: calc(16 * var(--rpx));
		padding: 0 calc(28 * var(--rpx));
		height: calc(72 * var(--rpx));
		line-height: calc(72 * var(--rpx));
		border-radius: calc(36 * var(--rpx));
		background: #4285f4;
		color: #fff;
		font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
		flex-shrink: 0;
	}

	.footer-next.locked {
		background: #fff;
		color: #4285f4;
		border: calc(1 * var(--rpx)) solid #4285f4;
		line-height: calc(70 * var(--rpx));
	}
</style>