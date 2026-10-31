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
				<audio v-if="audioSrc" :id="'audioPlayer'" :src="audioSrc" :poster="static(info.course.thumb)" @play="onAudioPlay"
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
						<span class="chapter-nav-label">下一节</span>
						<span class="chapter-nav-title">{{ nextChapter ? nextChapter.title : '没有了' }}</span>
					</div>
				</div>
			</div>

			<StudyCompleteCard
				v-if="info.paper_id && examPassed"
				title="学习与测评已完成"
				desc="您已完成本章节学习并通过测评，可查看答卷详情"
				action="查看测评"
				@action="goViewExamResult"
			/>
			<StudyCompleteCard
				v-else-if="info.paper_id && examPending"
				pending
				icon="clock"
				title="测评待批改"
				desc="您已提交测评，请等待教师批改"
				action="查看答卷"
				@action="goViewExamResult"
			/>
			<StudyCompleteCard
				v-else-if="info.paper_id"
				pending
				icon="clock"
				title="内容学习完毕"
				desc="恭喜您完成本章节学习，快去完成测评作业吧"
				action="完成测评作业"
				@action="goToExam"
			/>
			</div>

			<div class="chapter-main">
			<van-tabs v-if="!isDesktop" v-model:active="chapterTabIndex" color="#4285f4" title-active-color="#353535" title-inactive-color="#7f848c" line-width="22" line-height="3" @change="(i) => switchChapterTab({ index: Number(i) })">
				<van-tab v-for="(t, i) in chapterTabs" :key="i" :title="t.name" :name="i" />
			</van-tabs>
			<div v-else class="chapter-tabs">
				<div :class="{ on: detailTab === 'intro' }" @click="switchChapterTab(0)">介绍</div>
				<div :class="{ on: detailTab === 'comment' }" @click="switchChapterTab(1)">评论（{{ commentTotal }}）</div>
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

			<CommentList
				v-if="showChapterComment"
				:list="commentList"
				:loading="commentLoading"
				:has-more="commentHasMore"
				:title="'讨论' + (commentTotal > 0 ? '（' + commentTotal + '）' : '')"
				reply-on-child
				@reply="replyComment"
			/>
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
				<span>下一节</span>
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
import { autoCacheMedia, localPlaybackSrc } from '@/common/mediaCache'
import { getStorage, removeStorage, setStorage } from '@/common/storage'
import { createAudioContext, pageScrollTo } from '@/common/dom'
import { navigateTo, setPageTitle } from '@/common/nav'
	import { isDesktop as desktopOn } from '@/pc/desktop'
	import CommentInputPopup from './components/commentInputPopup.vue'
	import CommentDock from './components/CommentDock.vue'
	import CourseVideoPlayer from './components/CourseVideoPlayer.vue'
	import CommentList from './components/CommentList.vue'
	import StudyCompleteCard from './components/StudyCompleteCard.vue'

	const chapterDetailCache = {}
	const warmedMediaUrls = {}

	export default {
		components: {
			CommentInputPopup,
			CommentDock,
			CourseVideoPlayer,
			CommentList,
			StudyCompleteCard
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
				showCommentPopup: false,
				desktopComment: '',
				audioCurrentTime: 0,
				audioSrc: '',
				audioBlobUrl: '',
				pendingAudioPlay: false,
				audioFallbackUsed: false,
				audioDuration: 0,
				audioLoading: false,
				audioProgressPercent: 0,
				audioProgressTimer: null,
				isDragging: false,
				videoCurrentTime: 0,
				videoDuration: 0,
				videoProgressTimer: null,
				videoInitialTime: 0,
				commentPage: 1,
				commentLimit: 50,
				commentHasMore: true,
				commentLoading: false,
				replyCommentId: null,
				replyUserId: null,
				replyUserName: null,
				userPassPaperId: null,
				isPeopleCheck: 0,
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
			this.sampleVideoProgress();
			this.saveAudioProgress();
			this.saveVideoProgress();
		},
		onUnload() {
			this.sampleVideoProgress();
			this.saveAudioProgress();
			this.saveVideoProgress();
			this.pauseAudio();
			this.pauseVideo();
			this.releaseAudioBlob();
		},
		onReachBottom() {
			if (this.detailTab === 'comment' && this.commentList.length && this.commentHasMore && !this.commentLoading && this.info && this.info.id) {
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
			showChapterNav() {
				const uc = this.userChapter || {};
				return Number(uc.un_lock) === 1 || Number(uc.has_end_chapter) === 1 || Number(uc.is_end) === 1;
			},
			examPassed() {
				return !!(this.info && this.info.paper_id && this.userChapter && Number(this.userChapter.is_end) === 1);
			},
			examPending() {
				return !!(this.info && this.info.paper_id && Number(this.isPeopleCheck) === 1 && this.userPassPaperId);
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
			resetMediaState() {
				this.stopProgressTimer();
				this.stopVideoProgressTimer();
				this.isPlaying = false;
				this.audioCurrentTime = 0;
				this.audioDuration = 0;
				this.audioLoading = false;
				this.audioProgressPercent = 0;
				this.audioProgressInited = false;
				this.releaseAudioBlob();
				this.audioSrc = '';
				this.pendingAudioPlay = false;
				this.audioFallbackUsed = false;
				this.isDragging = false;
				this.videoCurrentTime = 0;
				this.videoDuration = 0;
				this.videoInitialTime = 0;
				this.prefetchToken = '';
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
				this.isPeopleCheck = Number(payload.isPeopleCheck) === 1 ? 1 : 0;
				this.prevChapter = payload.prevChapter || null;
				this.nextChapter = payload.nextChapter || null;
				if (payload.commentTotal != null && payload.commentTotal !== '') {
					const total = Number(payload.commentTotal);
					if (!Number.isNaN(total) && total >= 0) this.commentTotal = total;
				}
				if (full) {
					this.commentList = [];
					this.commentPage = 1;
					this.commentHasMore = true;
					this.commentReqId += 1;
				}
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
						this.bindLocalAudio();
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
			getChapterDetail(refresh = false) {
				if (!this.chapterId) return;
				if (refresh) {
					this.fetchChapterDetail(true);
					return;
				}
				this.pageLoading = true;
				this.pageError = false;
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
					refresh: 1,
					prefetch: 1
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
			toggleAudio() {
				if (!this.info || !this.info.media_url) {
					this.msg('音频资源不存在');
					return;
				}
				if (!document.getElementById('audioPlayer')) {
					this.pendingAudioPlay = true;
					return;
				}
				const audioContext = createAudioContext('audioPlayer', this);
				if (this.isPlaying) {
					audioContext.pause();
				} else {
					if (this.audioCurrentTime > 0 && this.audioCurrentTime < this.audioDuration) {
						audioContext.seek(this.audioCurrentTime);
					}
					audioContext.play();
				}
			},
			pauseAudio() {
				if (this.isPlaying) {
					const audioContext = createAudioContext('audioPlayer', this);
					audioContext.pause();
				}
			},
			onAudioPlay() {
				this.isPlaying = true;
				this.audioLoading = false;
				this.startProgressTimer();
				const url = this.static(this.info && this.info.media_url);
				if (url) autoCacheMedia(url, { kind: 'audio', title: (this.info && this.info.title) || '' });
			},
			onAudioPause() {
				this.isPlaying = false;
				this.stopProgressTimer();
				this.saveAudioProgress();
			},
			onAudioTimeUpdate(e) {
				if (this.isDragging) return;
				if (e.detail && e.detail.currentTime !== undefined) {
					this.audioCurrentTime = e.detail.currentTime;
					if (this.audioDuration > 0) {
						this.audioProgressPercent = (this.audioCurrentTime / this.audioDuration) * 100;
					}
					if (!this.audioDuration && e.detail.duration) {
						this.audioDuration = e.detail.duration;
					}
				}
			},
			onProgressChanging(e) {
				this.isDragging = true;
				const percent = e.detail.value;
				if (this.audioDuration > 0) {
					const newTime = (percent / 100) * this.audioDuration;
					this.audioCurrentTime = newTime;
					this.audioProgressPercent = percent;
				}
			},
			onProgressChange(e) {
				this.isDragging = false;
				const percent = e.detail.value;
				if (this.audioDuration > 0) {
					const seekTime = (percent / 100) * this.audioDuration;
					this.audioCurrentTime = seekTime;
					this.audioProgressPercent = percent;
					const audioContext = createAudioContext('audioPlayer', this);
					audioContext.seek(seekTime);
				}
			},
			onAudioEnded() {
				this.isPlaying = false;
				this.stopProgressTimer();
				this.clearAudioProgress();
			},
			onAudioWaiting() {
				this.audioLoading = true;
			},
			onAudioCanPlay(e) {
				this.audioLoading = false;
				if (e.detail && e.detail.duration) {
					this.audioDuration = e.detail.duration;
				} else if (this.info && this.info.media && this.info.media.duration) {
					this.audioDuration = this.info.media.duration;
				}
				this.$nextTick(() => {
					this.restoreAudioProgress();
				});
				this.prefetchNext();
			},
			onAudioError(e) {
				const remote = this.static(this.info && this.info.media_url);
				if (!this.audioFallbackUsed && remote && this.audioSrc && this.audioSrc !== remote) {
					this.audioFallbackUsed = true;
					this.audioLoading = false;
					this.isPlaying = false;
					this.releaseAudioBlob();
					this.audioSrc = remote;
					return;
				}
				this.audioLoading = false;
				this.isPlaying = false;
				console.error('音频播放错误:', e);
				this.msg('音频加载失败，请稍后重试');
			},
			startProgressTimer() {
				this.stopProgressTimer();
				this.audioProgressTimer = setInterval(() => {
					if (this.isPlaying && this.audioCurrentTime > 0) {
						this.saveAudioProgress();
					}
				}, 5000);
			},
			stopProgressTimer() {
				if (this.audioProgressTimer) {
					clearInterval(this.audioProgressTimer);
					this.audioProgressTimer = null;
				}
			},
			formatTime(seconds) {
				if (!seconds || isNaN(seconds)) return '00:00';
				const mins = Math.floor(seconds / 60);
				const secs = Math.floor(seconds % 60);
				return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
			},
			saveAudioProgress() {
				if (!this.chapterId) return;
				if (!this.audioCurrentTime || this.audioCurrentTime <= 0) return;
				if (!this.audioDuration || this.audioDuration <= 0) {
					if (this.info && this.info.media && this.info.media.duration) {
						this.audioDuration = this.info.media.duration;
					} else {
						return;
					}
				}
				if ((this.audioCurrentTime / this.audioDuration) >= 0.95) {
					this.clearAudioProgress();
					return;
				}
				const storageKey = `chapter_progress_${this.chapterId}`;
				try {
					setStorage(storageKey, {
						currentTime: this.audioCurrentTime,
						duration: this.audioDuration,
						updateTime: Date.now()
					});
				} catch (e) {
				}
			},
			restoreAudioProgress() {
				if (!this.chapterId) return;
				if (!this.audioDuration || this.audioDuration <= 0) {
					if (this.info && this.info.media && this.info.media.duration) {
						this.audioDuration = this.info.media.duration;
					} else {
						setTimeout(() => {
							this.restoreAudioProgress();
						}, 500);
						return;
					}
				}
				const storageKey = `chapter_progress_${this.chapterId}`;
				let savedProgress = null;
				try {
					savedProgress = getStorage(storageKey);
				} catch (e) {
				}
				if (savedProgress && savedProgress.currentTime && savedProgress.duration) {
					const { currentTime, duration } = savedProgress;
					if (currentTime > 0 && duration > 0 && (currentTime / duration) < 0.95) {
						if (this.audioDuration !== duration) {
							const progressPercent = (currentTime / duration) * 100;
							this.audioCurrentTime = (progressPercent / 100) * this.audioDuration;
						} else {
							this.audioCurrentTime = currentTime;
						}
						this.audioProgressPercent = (this.audioCurrentTime / this.audioDuration) * 100;
						setTimeout(() => {
							try {
								const audioContext = createAudioContext('audioPlayer', this);
								if (audioContext && typeof audioContext.seek === 'function') {
									audioContext.seek(this.audioCurrentTime);
								}
							} catch (e) {
							}
						}, 500);
					} else {
						this.clearAudioProgress();
					}
				} else if (this.userChapter && this.userChapter.progress && this.audioDuration > 0) {
					this.audioCurrentTime = (this.userChapter.progress / 100) * this.audioDuration;
					this.audioProgressPercent = this.userChapter.progress;
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
			clearAudioProgress() {
				if (!this.chapterId) return;
				const storageKey = `chapter_progress_${this.chapterId}`;
				removeStorage(storageKey);
			},
			releaseAudioBlob() {
				if (!this.audioBlobUrl) return;
				URL.revokeObjectURL(this.audioBlobUrl);
				this.audioBlobUrl = '';
			},
			async bindLocalAudio() {
				const remote = this.static(this.info && this.info.media_url);
				const chapterId = this.chapterId;
				let local = '';
				try {
					local = remote ? await localPlaybackSrc(remote, 'audio') : '';
				} catch (e) {}
				if (String(chapterId) !== String(this.chapterId)) {
					if (local && String(local).startsWith('blob:')) URL.revokeObjectURL(local);
					return;
				}
				if (!local) {
					this.audioSrc = remote || '';
					this.playPendingAudio(chapterId);
					return;
				}
				if (this.isPlaying) {
					if (String(local).startsWith('blob:')) URL.revokeObjectURL(local);
					return;
				}
				this.releaseAudioBlob();
				if (String(local).startsWith('blob:')) this.audioBlobUrl = local;
				this.audioSrc = local;
				this.playPendingAudio(chapterId);
			},
			playPendingAudio(chapterId) {
				if (!this.pendingAudioPlay) return;
				this.pendingAudioPlay = false;
				this.$nextTick(() => {
					if (String(chapterId) !== String(this.chapterId)) return;
					const audioContext = createAudioContext('audioPlayer', this);
					if (this.audioCurrentTime > 0 && this.audioCurrentTime < this.audioDuration) {
						audioContext.seek(this.audioCurrentTime);
					}
					audioContext.play();
				});
			},
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
			onVideoPlay() {
				this.startVideoProgressTimer();
			},
			onVideoCanPlay() {
				this.prefetchNext();
			},
			onVideoPause() {
				this.stopVideoProgressTimer();
				this.saveVideoProgress();
			},
			onVideoEnded() {
				this.stopVideoProgressTimer();
				this.clearVideoProgress();
			},
			sampleVideoProgress() {
				const player = this.$refs.courseVideo;
				if (player && typeof player.emitProgress === 'function') player.emitProgress();
			},
			onVideoTimeUpdate(payload) {
				const cur = payload && payload.currentTime != null ? Number(payload.currentTime) : this.videoCurrentTime;
				const dur = payload && payload.duration != null ? Number(payload.duration) : this.videoDuration;
				if (Number.isFinite(cur) && cur >= 0) this.videoCurrentTime = cur;
				if (Number.isFinite(dur) && dur > 0) this.videoDuration = dur;
			},
			pauseVideo() {
				if (this.info && this.info.type === 1 && this.$refs.courseVideo) {
					this.$refs.courseVideo.pause();
				}
			},
			startVideoProgressTimer() {
				this.stopVideoProgressTimer();
				this.videoProgressTimer = setInterval(() => {
					if (this.videoCurrentTime > 0) {
						this.saveVideoProgress();
					}
				}, 5000);
			},
			stopVideoProgressTimer() {
				if (this.videoProgressTimer) {
					clearInterval(this.videoProgressTimer);
					this.videoProgressTimer = null;
				}
			},
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
			saveVideoProgress() {
				if (!this.chapterId || !this.videoCurrentTime || this.videoDuration <= 0) return;
				if ((this.videoCurrentTime / this.videoDuration) >= 0.95) {
					this.clearVideoProgress();
					return;
				}
				const storageKey = `chapter_progress_${this.chapterId}`;
				setStorage(storageKey, {
					currentTime: this.videoCurrentTime,
					duration: this.videoDuration,
					updateTime: Date.now()
				});
			},
			clearVideoProgress() {
				if (!this.chapterId) return;
				const storageKey = `chapter_progress_${this.chapterId}`;
				removeStorage(storageKey);
			},
			goToExam() {
				if (!this.courseId && !this.chapterId) {
					this.msg('参数错误');
					return;
				}
				this.needStudyRefresh = true;
				const examUrl = `/pages/exam/exam?id=${this.info.paper_id || ''}`;
				navigateTo(examUrl);
			},
			goViewExamResult() {
				const id = this.userPassPaperId;
				if (!id) {
					this.msg('暂无测评记录，请下拉刷新重试');
					return;
				}
				this.needStudyRefresh = true;
				navigateTo(`/pages/exam/examDetail?id=${id}`);
			},
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
			toggleLike(item) {
			},
			viewReplies(item) {
			},
			switchChapterTab(e) {
				const index = typeof e === 'object' && e !== null ? e.index : e
				this.detailTab = Number(index) === 1 ? 'comment' : 'intro'
				if (this.detailTab === 'comment' && !this.commentList.length && this.chapterId) {
					this.getCommentList(true);
				}
			},
			showCommentInput() {
				if (this.isDesktop) {
					this.switchChapterTab(1);
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
			handleCommentClose() {
				this.showCommentPopup = false;
				this.desktopComment = '';
				this.clearReply();
			},
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
			getCommentList(reset = true) {
				if (!this.chapterId) return;
				const isReset = !!reset;
				if (isReset) {
					this.commentPage = 1;
					this.commentList = [];
					this.commentHasMore = true;
					this.commentLoading = false;
					this.commentTotal = 0;
				}
				if (this.commentLoading || !this.commentHasMore) return;

				const reqId = ++this.commentReqId;
				const reqChapterId = this.chapterId;
				const pageNum = this.commentPage;
				this.commentLoading = true;
				this.$http.post('/comment/getChapterCommentList', {
					chapter_id: Number(reqChapterId) || reqChapterId,
					course_id: Number(this.courseId) || this.courseId,
					page: pageNum,
					limit: this.commentLimit
				}).then(res => {
					if (reqId !== this.commentReqId || String(reqChapterId) !== String(this.chapterId)) return;
					this.commentLoading = false;
					const payload = (res && res.data) || {};
					const list = Array.isArray(payload.list) ? payload.list : [];
					this.commentList = isReset ? list.slice() : this.commentList.concat(list);
					const hasMore = payload.hasMore === true || payload.hasMore === 1 || payload.hasMore === 'true';
					this.commentHasMore = hasMore && list.length > 0;
					if (typeof payload.total === 'number') {
						this.commentTotal = payload.total;
					}
					this.fillCommentsIfShort();
				}).catch(err => {
					if (reqId !== this.commentReqId) return;
					this.commentLoading = false;
					if (isReset) this.commentList = [];
					console.log(err);
					this.msg((err && (err.message || err.msg)) || '加载评论失败');
				});
			},
			loadMoreComments() {
				if (this.commentLoading || !this.commentHasMore) return;
				this.commentPage += 1;
				this.getCommentList(false);
			},
			fillCommentsIfShort() {
				if (this.detailTab !== 'comment' || this.commentLoading || !this.commentHasMore) return;
				this.$nextTick(() => {
					const gap = document.documentElement.scrollHeight - window.innerHeight;
					if (gap < 80) this.loadMoreComments();
				});
			},
			replyComment(item, rootItem) {
				const root = rootItem || item;
				this.replyCommentId = root.id;
				this.replyUserId = item.user && item.user.id;
				this.replyUserName = (item.user && item.user.nick_name) || '匿名';
				if (this.isDesktop) {
					this.switchChapterTab(1);
					this.focusDesktop();
					return;
				}
				this.showCommentPopup = true;
			},
			goNeighbor(chapter) {
				if (!chapter || !chapter.id) return;
				if (chapter.un_lock == 0 && Number(this.userChapter && this.userChapter.is_end) !== 1) {
					this.msg(this.info && this.info.paper_id ? '请先完成测评' : '请先完成本节学习');
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
			goToCourse() {
				const courseId = this.courseId
					|| (this.info && (this.info.course_id || (this.info.course && this.info.course.id)));
				if (!courseId) {
					this.msg('课程信息缺失');
					return;
				}
				this.courseId = courseId;
				this.navTo(`/pages/course/courseDetail?id=${courseId}&tab=1`);
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
		display: flex;
		align-items: center;
		justify-content: center;
		line-height: 1;
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
	}
</style>