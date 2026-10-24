<template>
	<div>
		<div class="course-skeleton" v-if="!info && loading">
			<div class="skeleton-header">
				<div class="skeleton-line skeleton-title"></div>
				<div class="skeleton-line skeleton-short"></div>
			</div>
			<div class="skeleton-cover"></div>
			<div class="skeleton-line"></div>
			<div class="skeleton-line skeleton-short"></div>
		</div>
		<div class="course-load-fail" v-else-if="!info && pageError" @click="getCourseDetail(false)">
			<span>加载失败，点击重试</span>
		</div>
		<div class="course-intro" v-else-if="info">
			<CourseHero
				v-if="isDesktop"
				:info="info"
				:chapters="chapters"
				:is-join="isJoin"
				:can-join="canJoin"
				:apply-pending="applyPending"
				:is-collect="isCollect"
				:label="footerPrimaryLabel"
				:tip="footerPrimaryTip"
				@primary="onFooterPrimary"
				@join="joinCourse"
				@collect="toggleCollect"
			/>
			<div class="course-brief__header" v-if="!isDesktop">
				<div class="course-brief__left">
					<div class="course-brief__title">
						<span class="course-brief__category" v-if="info.category && info.category.name">{{ info.category.name }}</span>
						{{ info.title || '' }}
					</div>
					<div class="course-brief__subtitle">{{ info.desc || '' }}</div>
				</div>
				<div class="course-brief__right">
					<img :src="static(info.thumb || info.cover || info.img_url)"
						class="course-brief__img"  style="object-fit:cover" />
					<div class="course-brief__cover-modal">
						<div class="course-brief__play-num">
							<van-icon name="eye-o" size="16" color="#fff" />
							<span>{{ info.view_count || 0 }}</span>
						</div>
					</div>
				</div>
			</div>
			<div class="course-brief__center" v-if="!isDesktop">
				<div class="course-brief__item">
					<div class="course-brief__label">{{ chapters.length || info.chapter_count || 0 }} 讲</div>
					<div class="course-brief__tip">总讲数</div>
				</div>
				<div class="course-brief__item">
					<div class="course-brief__label">{{ info.credit || '0' }}</div>
					<div class="course-brief__tip">获得学分</div>
				</div>
				<div class="course-brief__item">
					<div class="course-brief__label">{{ info.type_name || '' }}</div>
					<div class="course-brief__tip">授课形式</div>
				</div>
			</div>

			<div class="course-cert-tip course-cert-tip--link" v-if="info.my_cert && info.my_cert.image" @click="viewMyCert">
				已获得「{{ info.cert && info.cert.name || '结业证书' }}」，点击查看
			</div>
			<div class="course-cert-tip" v-else-if="info.cert_id && !isEnd">
				完成本课程后，可获得「{{ info.cert && info.cert.name || '结业证书' }}」
			</div>

			<div class="study-complete-module" v-if="!isDesktop && isEnd && !showPaper">
				<div class="study-complete-content">
					<div class="study-complete-icon">
						<van-icon name="success" size="24" color="#52c41a" />
					</div>
					<div class="study-complete-text">
						<div class="study-complete-title">课程已结业</div>
						<div class="study-complete-desc" v-if="info.my_cert && info.my_cert.image">结业证书已发放，可随时查看</div>
						<div class="study-complete-desc" v-else-if="info.cert_id">证书生成中，可稍后到「我的证书」查看</div>
						<div class="study-complete-desc" v-else>恭喜完成本课程学习</div>
					</div>
				</div>
				<div class="study-complete-btn" v-if="info.my_cert && info.my_cert.image" @click="viewMyCert">
					<span>查看证书</span>
					<van-icon name="arrow" size="14" color="#fff" />
				</div>
				<div class="study-complete-btn" v-else @click="goMyCert">
					<span>我的证书</span>
					<van-icon name="arrow" size="14" color="#fff" />
				</div>
			</div>

			<div class="study-complete-module study-complete-module--pending" v-if="!isDesktop && showPaper">
				<div class="study-complete-content">
					<div class="study-complete-icon">
						<van-icon name="clock" size="24" color="#ff3359" />
					</div>
					<div class="study-complete-text">
						<div class="study-complete-title">课程学习完毕</div>
						<div class="study-complete-desc" v-if="info.paper_id">恭喜您完成本课程学习，快去完成考试吧</div>
					</div>
				</div>
				<div class="study-complete-btn" v-if="info.paper_id"
					@click="goExam()">
					<span>完成考试</span>
					<van-icon name="arrow" size="14" color="#fff" />
				</div>
				<div class="study-complete-btn" v-else @click="endCourse()">
					<span>领取学分</span>
					<van-icon name="arrow" size="14" color="#fff" />
				</div>
			</div>

			<div class="course-body">
				<van-tabs v-if="!isDesktop && courseTabs.length" v-model:active="activeTab" color="#4285f4" title-active-color="#353535" title-inactive-color="#7f848c" line-width="22" line-height="3" @change="(i) => switchTab({ index: Number(i) })">
					<van-tab v-for="(t, i) in courseTabs" :key="i" :title="t.name" :name="i" />
				</van-tabs>
				<div v-if="isDesktop" class="pc-course-tabs">
					<div :class="{ on: activeTab !== 2 && activeTab !== 3 }" @click="activeTab = 0">简介</div>
					<div :class="{ on: activeTab === 2 }" @click="activeTab = 2">课件</div>
					<div :class="{ on: activeTab === 3 }" @click="switchTab({ index: 3 })">评论（{{ commentTotal }}）</div>
				</div>
				<!-- 简介内容 -->
				<div class="intro-content" v-if="isDesktop ? activeTab !== 2 && activeTab !== 3 : activeTab === 0">
					<div class="intro-module">
						<div class="intro-module__title"><span>课程介绍</span></div>
						<div class="intro-module__item">
							<div class="article-typo">
								<div class="qx-parse" style="user-select:text" v-html="richHtml(info.content)"></div>
							</div>
						</div>
					</div>
				</div>

				<!-- 目录/章节列表 -->
				<div class="catalog-wrapper" v-if="isDesktop || activeTab === 1">
					<div class="catalog-list">
						<div class="catalog-chapter">
							<div class="catalog-chapter__title">课程 ({{ chapters.length }}讲) </div>
						</div>
						<div class="catalog-article-list">
							<div
								class="catalog-article-item"
								:class="{ 'is-current': isContinueChapter(chapter) }"
								:id="'catalog-ch-' + chapter.id"
								v-for="(chapter, index) in chapters"
								:key="chapter.id || index"
								@click="goChapterDetail(chapter)"
							>
								<span class="iconfont catalog-article-item__type-icon">
									<van-icon :name="chapterTypeIcon(chapter)" size="16" :color="isContinueChapter(chapter) ? '#4285f4' : '#b9bdc2'" />
								</span>
								<div class="catalog-article-item__info">
									<div class="catalog-article-item__title">
										{{ chapterIndexLabel(index) }} | {{ chapter.title }}
									</div>
									<div class="catalog-article-item__bottom">
										<span class="catalog-article-item__duration">{{ chapterTypeLabel(chapter) }}<span v-if="chapter.media && chapter.media.duration"> · 时长 {{ showDuration(chapter.media.duration) }}</span></span>
									</div>
								</div>
								<div class="catalog-article-item__status" v-if="isJoin">
									<div
										class="catalog-article-item__try"
										:class="{
											'is-pending': chapterStatus(chapter, index) === 'exam',
											'is-locked': chapterStatus(chapter, index) === 'wait'
										}"
										@click.stop="onChapterStatusClick(chapter, index)"
									>{{ chapterStatusText(chapter, index) }}</div>
								</div>
							</div>
						</div>
					</div>
				</div>

				<!-- 附件列表 -->
				<div class="attach-wrapper" v-if="activeTab === 2">
					<div class="attach-list">
						<div class="attach-item" v-for="(attach, index) in attachList" :key="attach.id || index">
							<img class="attach-item__icon" src="/static/img/file.png" />
							<div class="attach-item__info">
								<div class="attach-item__name">{{ attach.name || attach.file_name }}</div>
								<div class="attach-item__size" v-if="attach.size">{{ showSize(attach.size) }}</div>
							</div>
							<div class="attach-item__action" v-if="isJoin" @click="downloadAttach(attach)">
								<img class="attach-item__action_icon" src="/static/img/me/download.png" />
							</div>
						</div>
						<div class="empty-tip" v-if="attachList.length === 0">
							<van-empty description="暂无课件" image="/static/empty/attach.png" :image-size="100" />
						</div>
					</div>
				</div>

				<!-- 评论区域 -->
				<div class="comment-section" v-if="activeTab === 3">
					<!-- 评论列表 -->
					<div class="comment-list">
						<div class="empty-tip" v-if="commentList.length === 0 && !commentLoading">
							<van-empty description="暂无讨论" image="/static/empty/comment.png" :image-size="100" />
						</div>
						<div class="comment-item" v-for="(item, index) in commentList" :key="item.id || index">
							<div class="comment-header">
								<div class="comment-user">
									<van-image :src="static((item.user && item.user.avatar) || '/static/default-avatar.png')" round :width="36 + 'px'" :height="36 + 'px'" fit="cover" />
									<div class="comment-user-info">
										<div class="comment-name"><span>{{ (item.user && item.user.nick_name) || '匿名' }}</span><coin-badge :coins="item.user ? item.user.coins : null" /></div>
										<div class="comment-time">{{ date(item.created_at, 'yyyy-mm-dd') }}</div>
									</div>
								</div>
								<div class="comment-actions">
									<div class="comment-reply-btn" @click.stop="replyComment(item)">
										<van-icon name="chat-o" size="14" color="#7f848c" />
										<span>回复</span>
									</div>
								</div>
							</div>
							<div class="comment-content">{{ item.content || '' }}</div>
							<!-- 回复列表 -->
							<div class="comment-replies" v-if="item.children && item.children.length > 0">
								<div class="reply-item" v-for="(reply, replyIndex) in item.children"
									:key="reply.id || replyIndex">
									<div class="reply-header">
										<div class="reply-name"><span>{{ (reply.user && reply.user.nick_name) || '匿名' }}</span><coin-badge :coins="reply.user ? reply.user.coins : null" /></div>
										<span class="reply-time">{{ date(reply.created_at, 'yyyy-mm-dd') }}</span>
									</div>
									<div class="reply-content">
										<span class="reply-to" v-if="replyToName(reply, item)">回复
											@{{ replyToName(reply, item) }}：</span>
										<span>{{ reply.content || '' }}</span>
									</div>
								</div>
							</div>
						</div>
						<!-- 加载更多提示 -->
						<div class="comment-loading" v-if="commentLoading">
							<van-loading type="circular" />
							<span>加载中...</span>
						</div>
						<div class="comment-no-more" v-if="commentHasMore === false && commentList.length > 0">
							<span>没有更多了</span>
						</div>
					</div>
				</div>
				<div class="chapter-footer" v-if="activeTab === 3">
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
				</div>
			</div>
			<div class="course-footer" v-if="!isDesktop && activeTab !== 3">
				<div class="course-footer__left">
					<div class="course-footer__collect" :class="{ active: isCollect }" @click="toggleCollect">
						<span class="iconfont course-footer__icon">
							<van-icon :name="isCollect ? 'star' : 'star-o'" size="16"
								:color="isCollect ? '#4285f4' : '#404040'" />
						</span>
						<span>{{ isCollect ? '已收藏' : '收藏' }}</span>
					</div>
				</div>
				<div class="course-footer__right">
					<div class="course-btn course-footer__btn" v-if="isJoin" @click="onFooterPrimary">
						<span class="course-footer__label">{{ footerPrimaryLabel }}</span>
						<span class="course-footer__tip" v-if="footerPrimaryTip">{{ footerPrimaryTip }}</span>
					</div>
					<div class="course-btn course-footer__btn is-disabled" v-else-if="applyPending">
						<span class="course-footer__label">审核中，请耐心等待</span>
					</div>
					<div class="course-btn course-footer__btn is-disabled" v-else-if="chapterTotal <= 0">
						<span class="course-footer__label">暂无课节</span>
					</div>
					<div class="course-btn course-footer__btn" v-else @click="joinCourse()">
						<span class="course-footer__label">{{ canJoin ? '加入课程' : '申请加入' }}</span>
						<span class="course-footer__tip" v-if="!canJoin">要求最低学分：{{ info.min_need_credit }}</span>
					</div>
				</div>
			</div>
		</div>
		<!-- 评论输入弹窗 -->
		<CommentInputPopup v-if="!isDesktop" :show="showCommentPopup" :placeholder="replyUserName ? `回复 @${replyUserName}：` : '我来说一说…'"
			@close="handleCommentClose" @publish="handleCommentPublish" />
	</div>
</template>

<script>
import { hideLoading, showLoading } from '@/common/toast'
import { getStorage, setStorage } from '@/common/storage'
import { createSelectorQuery, pageScrollTo, previewImage, showModal } from '@/common/dom'
import { setPageTitle } from '@/common/nav'
	import { defineAsyncComponent } from 'vue'
	import { isDesktop as desktopOn } from '@/pc/desktop'
	import CommentInputPopup from './components/commentInputPopup.vue'
	import CommentDock from './components/CommentDock.vue'

	export default {
		components: {
			CommentInputPopup,
			CommentDock,
			CourseHero: defineAsyncComponent(() => import('@/pc/CourseHero.vue')),
		},
		data() {
			return {
				courseId: null,
				loading: false,
				info: null, // 课程信息
				chapters: [], // 章节列表
				isJoin: false, // 是否加入
				isEnd: false, // 是否课程结束
				canJoin: false, // 是否可以加入
				applyPending: false,
				showPaper: false, // 是否显示考试
				isCollect: false, // 是否收藏
				attachList: [], // 附件列表
				activeTab: 0, // 当前激活的tab: 0-简介, 1-目录, 2-课件, 3-评论
				commentTotal: 0,
				commentList: [], // 评论列表
				showCommentPopup: false, // 是否显示评论输入弹窗
				desktopComment: '',
				commentPage: 1, // 评论列表页码
				commentLimit: 10, // 每页数量
				commentHasMore: true, // 是否还有更多评论
				commentLoading: false, // 评论加载状态
				replyCommentId: null, // 当前回复的评论ID
				replyUserId: null, // 当前回复的用户ID
				replyUserName: null, // 当前回复的用户名
				pageError: false,
				detailLoaded: false,
				needStudyRefresh: false
			}
		},
		computed: {
			isDesktop() {
				return desktopOn.value
			},
			chapterTotal() {
				const listed = this.chapters && this.chapters.length
				if (listed) return listed
				return Number(this.info && this.info.chapter_count) || 0
			},
			continueChapter() {
				if (!this.isJoin || !this.chapters || !this.chapters.length) return null;
				const list = this.chapters;
				const inProgress = list.find((c) => c.un_lock == 1 && c.is_end != 1);
				if (inProgress) return inProgress;
				const unlockable = list.find((c, i) => {
					if (c.un_lock != 0) return false;
					if (i === 0) return true;
					return list[i - 1] && list[i - 1].is_end == 1;
				});
				if (unlockable) return unlockable;
				return list[list.length - 1];
			},
			footerPrimaryLabel() {
				if (this.showPaper && this.info && this.info.paper_id) return '完成考试';
				if (this.showPaper) return '领取学分';
				const ch = this.continueChapter;
				if (!ch) return '继续学习';
				const idx = this.chapters.findIndex((c) => String(c.id) === String(ch.id));
				if (this.isChapterLocked(ch)) return '解锁学习';
				const started = this.chapters.some((c) => c.is_end == 1);
				if (!started && idx === 0) return '开始学习';
				return '继续学习';
			},
			footerPrimaryTip() {
				if (this.showPaper) return '';
				const ch = this.continueChapter;
				if (!ch) return '';
				const idx = this.chapters.findIndex((c) => String(c.id) === String(ch.id));
				const num = this.chapterIndexLabel(idx < 0 ? 0 : idx);
				return num + ' | ' + (ch.title || '');
			},
			courseTabs() {
				return [
					{ name: '简介' },
					{ name: '目录' },
					{ name: '课件' },
					{ name: '评论（' + (Number(this.commentTotal) || 0) + '）' },
				]
			}
		},
		onLoad(options) {
			if (options.id) {
				this.courseId = options.id;
			}
			if (options.tab != null && options.tab !== '') {
				const tab = Number(options.tab);
				if (!isNaN(tab)) this.activeTab = tab;
			}
		},
		onShow() {
			if (!this.detailLoaded) {
				this.getCourseDetail(false);
				return;
			}
			if (this.needStudyRefresh) {
				this.needStudyRefresh = false;
				this.getCourseDetail(true);
			}
		},
		onReady() {},
		methods: {
			// 获取课程详情
			courseCacheKey() {
				return 'qixue_course_' + this.courseId;
			},
			restoreCourseCache() {
				if (this.info || !this.courseId) return;
				try {
					const cached = getStorage(this.courseCacheKey());
					if (cached && cached.info) {
						this.applyCoursePayload(cached, false);
					}
				} catch (e) {}
			},
			saveCourseCache(data) {
				if (!this.courseId || !data || !data.info) return;
				try {
					setStorage(this.courseCacheKey(), data);
				} catch (e) {}
			},
			getCourseDetail(refresh = false) {
				if (!this.courseId) return;
				if (!refresh) {
					this.restoreCourseCache();
				}
				if (!this.info) {
					this.loading = true;
					this.pageError = false;
				}
				this.$http.get('/course/courseDetail', {
					id: this.courseId,
					refresh: refresh ? 1 : 0
				}).then(res => {
					this.applyCoursePayload(res.data, refresh);
					if (!refresh) {
						this.saveCourseCache(res.data);
					}
				}).catch(err => {
					this.loading = false;
					console.log(err);
					if (!this.info) {
						this.pageError = true;
					}
					this.msg('获取课程详情失败');
				});
			},
			applyCoursePayload(data, refresh) {
				this.loading = false;
				this.pageError = false;
				this.detailLoaded = true;
				const nextInfo = (data && data.info) || null;
				if (refresh && this.info && nextInfo && String(this.info.id) === String(nextInfo.id)) {
					this.info.view_count = nextInfo.view_count;
					this.info.chapter_count = nextInfo.chapter_count;
					this.info.cert_id = nextInfo.cert_id;
					this.info.cert = nextInfo.cert;
					this.info.my_cert = nextInfo.my_cert;
					this.info.paper_id = nextInfo.paper_id;
				} else {
					this.info = nextInfo;
					setPageTitle((this.info && this.info.title) || '');
				}
				this.isJoin = (data && data.isJoin) || false;
				this.isEnd = (data && data.isEnd) || false;
				this.canJoin = (data && data.canJoin) || false;
				this.applyPending = !!(data && data.applyPending);
				this.isCollect = (data && data.isCollect) || false;
				this.showPaper = (data && data.showPaper) || false;
				if (!refresh || !this.attachList.length) {
					this.attachList = (data && data.attachList) || [];
				}
				if (data && data.chapters) {
					this.chapters = data.chapters;
					if (this.activeTab === 1) {
						this.scrollToContinueChapter();
					}
				} else {
					this.refreshChapters();
				}
				if (!this.commentList.length) {
					this.getCommentList(true);
				}
			},
			viewMyCert() {
				const imageUrl = this.static((this.info && this.info.my_cert && this.info.my_cert.image) || '');
				if (!imageUrl) {
					this.msg('暂无证书图片');
					return;
				}
				previewImage({
					urls: [imageUrl],
					current: imageUrl
				});
			},
			goMyCert() {
				this.navTo('/pages/member/cert');
			},
			refreshChapters() {
				if (!this.courseId) return;
				this.$http.get('/course/getCourseChapterList', {
					course_id: this.courseId
				}).then(res => {
					this.chapters = res.data || [];
					if (this.activeTab === 1) {
						this.scrollToContinueChapter();
					}
				});
			},
			goExam() {
				if (!this.info || !this.info.paper_id) return;
				this.needStudyRefresh = true;
				this.navTo('/pages/exam/exam?id=' + this.info.paper_id);
			},
			chapterIndexLabel(index) {
				return (index + 1).toString().padStart(2, '0');
			},
			chapterTypeIcon(chapter) {
				const type = chapter && Number(chapter.type);
				if (type === 2) return 'volume';
				if (type === 3) return 'file-text';
				return 'play-circle';
			},
			chapterTypeLabel(chapter) {
				const type = chapter && Number(chapter.type);
				if (type === 3) return '图文';
				if (type === 2) return '音频';
				return '视频';
			},
			isChapterLocked(chapter) {
				return Number(chapter && chapter.un_lock) !== 1;
			},
			canUnlockChapter(chapter, index) {
				if (!this.isJoin || !chapter || !this.isChapterLocked(chapter)) return false;
				const i = index != null
					? Number(index)
					: this.chapters.findIndex((c) => String(c.id) === String(chapter.id));
				if (Number(chapter.sort) === 1 || i <= 0) return true;
				const prev = this.chapters[i - 1];
				return !!(prev && Number(prev.is_end) === 1);
			},
			chapterStatus(chapter, index) {
				if (!chapter) return 'wait';
				if (this.isChapterLocked(chapter)) {
					return this.canUnlockChapter(chapter, index) ? 'lock' : 'wait';
				}
				if (Number(chapter.is_end) === 1) return 'done';
				if (chapter.paper_id && Number(chapter.has_end_chapter) === 1) return 'exam';
				return 'study';
			},
			chapterStatusText(chapter, index) {
				const status = this.chapterStatus(chapter, index);
				if (status === 'lock') return '解锁';
				if (status === 'wait') return '未解锁';
				if (status === 'done') return '已学完';
				if (status === 'exam') return '去测评';
				return '学习';
			},
			isContinueChapter(chapter) {
				return !!(this.continueChapter && chapter && String(this.continueChapter.id) === String(chapter.id));
			},
			scrollToContinueChapter() {
				const ch = this.continueChapter;
				if (!ch || this.activeTab !== 1) return;
				this.$nextTick(() => {
					const query = createSelectorQuery().in(this);
					query.select('#catalog-ch-' + ch.id).boundingClientRect();
					query.selectViewport().scrollOffset();
					query.exec((res) => {
						const rect = res && res[0];
						const scroll = res && res[1];
						if (!rect || !scroll) return;
						const top = (scroll.scrollTop || 0) + rect.top - 80;
						if (top > 0) {
							pageScrollTo({
								scrollTop: top,
								duration: 280
							});
						}
					});
				});
			},
			switchTab(e) {
				const index = typeof e === 'object' && e !== null ? e.index : e;
				this.activeTab = Number(index);
				if (this.activeTab === 1) {
					this.scrollToContinueChapter();
				}
				if (this.activeTab === 3) {
					this.getCommentList(true);
				}
			},
			onFooterPrimary() {
				if (this.showPaper && this.info && this.info.paper_id) {
					this.goExam();
					return;
				}
				if (this.showPaper) {
					this.endCourse();
					return;
				}
				this.enterContinueChapter();
			},
			enterContinueChapter() {
				if (!this.continueChapter) {
					this.switchTab(1);
					return;
				}
				this.goChapterDetail(this.continueChapter);
			},
			onChapterStatusClick(chapter, index) {
				const status = this.chapterStatus(chapter, index);
				if (status === 'wait') {
					this.msg('请先完成上一章节');
					return;
				}
				if (status === 'lock') {
					this.confirmUnlockAndEnter(chapter);
					return;
				}
				if (status === 'exam' && chapter.paper_id) {
					this.needStudyRefresh = true;
					this.navTo('/pages/exam/exam?id=' + chapter.paper_id);
					return;
				}
				this.goChapterDetail(chapter);
			},
			confirmUnlockAndEnter(chapter) {
				showModal({
					title: '解锁章节',
					content: chapter.title ? `「${chapter.title}」尚未解锁，是否现在解锁？` : '该章节尚未解锁，是否现在解锁？',
					confirmText: '解锁',
					success: (res) => {
						if (!res.confirm) return;
						this.doUnlockChapter(chapter, true);
					}
				});
			},
			doUnlockChapter(chapter, enter) {
				this.$http.post('/learn/unlockChapter', {
					chapter_id: chapter.id
				}).then(() => {
					chapter.un_lock = 1;
					this.msg('解锁成功');
					if (enter) {
						this.openChapter(chapter);
						return;
					}
					this.getCourseDetail(true);
				}).catch((err) => {
					const message = (err && err.message) || '';
					if (message.indexOf('已经解锁') !== -1) {
						chapter.un_lock = 1;
						if (enter) this.openChapter(chapter);
						return;
					}
					this.msg(message || '请先完成上一章节');
				});
			},
			openChapter(chapter) {
				this.needStudyRefresh = true;
				this.navTo(`/pages/course/chapterDetail?id=${chapter.id}&courseId=${this.courseId}`);
			},
			// 收藏/取消收藏
			toggleCollect() {
				if (!this.courseId) return;
				this.$http.post('/learn/collect', {
					course_id: this.courseId
				}).then(res => {
					this.isCollect = !this.isCollect;
					this.msg(this.isCollect ? '收藏成功' : '取消收藏');
				}).catch(err => {
					console.log(err);
					this.msg(err.message || '操作失败');
				});
			},
			// 跳转到章节详情
			goChapterDetail(chapter) {
				if (!this.isJoin) {
					this.msg('请先加入课程');
					return;
				}
				if (this.isChapterLocked(chapter)) {
					if (this.canUnlockChapter(chapter)) {
						this.confirmUnlockAndEnter(chapter);
						return;
					}
					this.msg('请先完成上一章节');
					return;
				}
				this.openChapter(chapter);
			},
			// 加入课程
			joinCourse() {
				if (this.chapterTotal <= 0) {
					this.msg('课程暂无课节，无法加入');
					return;
				}
				if (this.applyPending) {
					this.msg('申请审核中，请耐心等待');
					return;
				}
				if (this.canJoin) {
					this.$http.post('/learn/joinCourse', {
						course_id: this.info.id
					}).then(res => {
						this.msg('加入成功');
						this.isJoin = true;
						this.canJoin = false;
						this.activeTab = 1;
						this.refreshChapters();
					}).catch(err => {
						console.log(err);
						this.msg(err.message || '操作失败');
					});
				} else {
					//跳转申请页面
					this.needStudyRefresh = true;
					this.navTo('/pages/course/applyJoin?id=' + this.courseId + '&title=' + encodeURIComponent((this.info && this.info.title) || ''));
				}
			},
			// 下载附件
			downloadAttach(attach) {
				if (!attach.path) {
					this.msg('附件链接不存在');
					return;
				}
				const url = this.static(attach.path);
				const fileName = attach.name || attach.file_name || '附件';

				showLoading('下载中...');
				fetch(url).then((res) => {
					if (!res.ok) throw new Error('download failed');
					return res.blob();
				}).then((blob) => {
					const objectUrl = URL.createObjectURL(blob);
					const link = document.createElement('a');
					link.setAttribute('download', fileName);
					link.setAttribute('href', objectUrl);
					document.body.appendChild(link);
					link.click();
					document.body.removeChild(link);
					URL.revokeObjectURL(objectUrl);
					hideLoading();
				}).catch(() => {
					hideLoading();
					window.open(url, '_blank');
				});

			},
			// 显示评论输入
			showCommentInput() {
				if (!this.isJoin) {
					this.msg('请先加入课程');
					return;
				}
				if (this.isDesktop) {
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
				if (!this.courseId) {
					this.msg('参数错误');
					return;
				}
				if (!this.isJoin) {
					this.msg('请先加入课程');
					return;
				}
				// 调用接口发布评论或回复
				const params = {
					course_id: this.courseId,
					content: content,
				};

				// 如果是回复，添加回复相关参数
				if (this.replyCommentId) {
					params.pid = this.replyCommentId;
					params.reply_user_id = this.replyUserId;
				} else {
					params.pid = 0;
				}

				this.$http.post('/comment/courseComment', params).then(res => {
					this.msg('发布成功');
					this.handleCommentClose();
					this.getCommentList(true);
				}).catch(err => {
					console.log(err);
					this.msg(err.message || '发布失败');
				});
			},
			// 获取评论列表
			getCommentList(reset = true) {
				if (!this.courseId) return;
				if (reset) {
					this.commentPage = 1;
					this.commentList = [];
					this.commentHasMore = true;
				}
				if (this.commentLoading || !this.commentHasMore) return;

				this.commentLoading = true;
				this.$http.post('/comment/getCourseCommentList', {
					course_id: this.courseId,
					page: this.commentPage,
					limit: this.commentLimit
				}).then(res => {
					this.commentLoading = false;
					const list = res.data.list || [];
					this.commentList = this.commentList.concat(list);
					this.commentHasMore = res.data.hasMore;
					if (typeof res.data.total === 'number') {
						this.commentTotal = res.data.total;
					}
				}).catch(err => {
					this.commentLoading = false;
					console.log(err);
					this.msg('加载评论失败');
				})
			},
			// 加载更多评论
			loadMoreComments() {
				this.commentPage++;
				this.getCommentList(false);
			},
			// 回复评论
			replyToName(reply, root) {
				if (reply.parent && reply.parent.user && reply.parent.user.nick_name) {
					return reply.parent.user.nick_name;
				}
				return (root && root.user && root.user.nick_name) || '';
			},
			replyComment(item) {
				if (!this.isJoin) {
					this.msg('请先加入课程');
					return;
				}
				this.replyCommentId = item.id;
				this.replyUserId = item.user && item.user.id;
				this.replyUserName = (item.user && item.user.nick_name) || '匿名';
				if (this.isDesktop) {
					this.focusDesktop();
					return;
				}
				this.showCommentPopup = true;
			},
			// 完成课程
			endCourse() {
				if (!this.showPaper || this.isEnd) {
					return;
				}
				showLoading('请求中...')

				//课程结束
				this.$http.post('/learn/endCourse', {
					course_id: this.info.id
				}).then(res => {
					hideLoading();
					this.isEnd = true;
					this.showPaper = false;
					this.msg('课程结束');
					this.getCourseDetail(true);
				}).catch(err => {
					console.log(err);
					this.msg(err.message || '操作失败');
				}).finally(() => {
					hideLoading();
				})
			}
		},
		onReachBottom() {
			// 上拉加载更多评论
			if (this.activeTab === 3 && this.commentHasMore && !this.commentLoading && this.courseId) {
				this.loadMoreComments();
			}
		}
	}
</script>

<style lang="scss" scoped>
	.course-skeleton {
		padding: calc(30 * var(--rpx)) 0 calc(40 * var(--rpx));
		background: #fff;
		min-height: 100vh;
	}

	.skeleton-header {
		padding-bottom: calc(20 * var(--rpx));
	}

	.skeleton-cover {
		width: calc(200 * var(--rpx));
		height: calc(150 * var(--rpx));
		margin: 0 calc(40 * var(--rpx)) calc(20 * var(--rpx)) auto;
		border-radius: calc(12 * var(--rpx));
		background: #f0f2f5;
	}

	.skeleton-line {
		height: calc(28 * var(--rpx));
		margin: calc(24 * var(--rpx)) calc(40 * var(--rpx)) 0;
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

	.course-load-fail {
		padding: calc(160 * var(--rpx)) calc(40 * var(--rpx));
		text-align: center;
		color: #7f848c;
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	}

	/* 学习完成模块样式 */
	.study-complete-module {
		margin: calc(40 * var(--rpx)) calc(40 * var(--rpx));
		padding: calc(32 * var(--rpx));
		background: linear-gradient(135deg, #f0f7ff 0%, #e8f4ff 100%);
		border-radius: calc(16 * var(--rpx));
		border: calc(1 * var(--rpx)) solid #d0e7ff;
		animation: tip-in 0.38s cubic-bezier(0.32, 0.72, 0, 1) both;
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
		transition: transform 0.16s ease, opacity 0.16s ease;
	}

	.study-complete-btn:active {
		transform: scale(0.98);
		opacity: 0.92;
	}

	.study-complete-btn span {
		margin-right: calc(8 * var(--rpx));
	}

	/* 与章节详情「待测评」一致：主色 #ff3359 */
	.study-complete-module--pending {
		background: linear-gradient(135deg, #fff5f7 0%, #ffe8ee 100%);
		border-color: #ffc9d6;
	}

	.study-complete-module--pending .study-complete-btn {
		background: linear-gradient(135deg, #ff3359, #ff6b82);
	}

	.course-cert-tip {
		margin: calc(24 * var(--rpx)) calc(40 * var(--rpx)) 0;
		padding: calc(20 * var(--rpx)) calc(24 * var(--rpx));
		border-radius: calc(12 * var(--rpx));
		background: #f0f7ff;
		color: #3b6fd9;
		font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
		line-height: 1.5;
		animation: tip-in 0.34s cubic-bezier(0.32, 0.72, 0, 1) both;
	}

	.course-cert-tip--link {
		background: #f6ffed;
		color: #389e0d;
		text-decoration: underline;
	}

	.course-cert-tip--link:active {
		opacity: 0.85;
	}

	@keyframes tip-in {
		from {
			opacity: 0;
			transform: translateY(calc(12 * var(--rpx)));
		}
		to {
			opacity: 1;
			transform: none;
		}
	}

	.course-btn {
		box-sizing: border-box;
		display: flex;
		align-items: center;
		justify-content: center;
		text-align: center;
		line-height: 1;
		transition: all .3s ease;
	}

	.course-btn:active {
		transition: none;
	}

	.course-authority__wrap {
		box-sizing: border-box;
		display: flex;
		align-items: center;
		height: 44px;
		padding: 0 0 0 10px;
		background: #f4f9fe;
		border-radius: 6px;
		overflow: hidden;
	}

	.article-typo {
		overflow: hidden;
		line-height: 30px;
		font-weight: 400;
		color: #404040;
		word-break: break-word;
		overflow-wrap: anywhere;

		:deep(img),
		:deep(img ){
			max-width: 100% !important;
			height: auto !important;
		}

		:deep(p),
		:deep(span),
		:deep(div),
		:deep(section ){
			white-space: normal !important;
			word-break: break-word;
			overflow-wrap: anywhere;
			max-width: 100%;
		}

		.li-item {
			list-style: inherit;
			color: #404040;
			line-height: 28px;
			font-size: calc(15px * var(--font-scale, 1));
			font-weight: 400;
		}

		.ol-list {
			margin-left: 21px;
			list-style: decimal;
			font-weight: 400;
			margin-bottom: .625rem;

			&:last-child {
				margin-bottom: 0;
			}
		}

		.p-text {
			line-height: 28px;
			font-size: calc(15px * var(--font-scale, 1));
			font-weight: 400;
			text-align: justify;
			word-break: break-word;
			display: block;
			margin-bottom: .625rem;

			&:last-child {
				margin-bottom: 0;
			}

			&+.p-text {
				margin-top: 26px;
			}
		}
	}

	.course-brief__header {
		position: relative;
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0 20px 21px;
	}

	.course-brief__left {
		flex: 1;
		min-width: 0;
		overflow: hidden;
		padding-right: 18px;

		.course-brief__category {
			display: inline-block;
			height: 21px;
			padding: 1px 2px;
			margin-right: 6px;
			border: 1px solid rgba(250, 137, 25, .35);
			border-radius: 4px;
			white-space: nowrap;
			font-size: calc(13px * var(--font-scale, 1));
			font-weight: 500;
			color: #fa8919;
			vertical-align: middle;
		}

		.course-brief__title {
			margin-bottom: 8px;
			line-height: 25px;
			font-size: calc(18px * var(--font-scale, 1));
			font-weight: 500;
			color: #353535;
			word-break: break-word;
			overflow-wrap: anywhere;
		}

		.course-brief__subtitle {
			line-height: 20px;
			font-size: calc(13px * var(--font-scale, 1));
			font-weight: 400;
			color: #7f848c;
			overflow: hidden;
			word-break: break-word;
			overflow-wrap: anywhere;
			white-space: normal;
			display: -webkit-box;
			-webkit-line-clamp: 4;
			line-clamp: 4;
			-webkit-box-orient: vertical;
		}
	}

	.course-brief__right {
		display: flex;
		align-items: center;
		justify-content: center;
		overflow: hidden;
		position: relative;
		flex-shrink: 0;
		width: 90px;
		height: 90px;
		background: #f6f7fb;
		border-radius: 4px;

		.course-brief__img {
			width: 100%;
			height: 100%;
		}
	}

	.course-brief__right .course-brief__play-num {
		position: absolute;
		bottom: 0;
		right: 4px;
		display: flex;
		align-items: center;
		font-size: calc(12px * var(--font-scale, 1));
		font-weight: 400;
		color: #fff;
	}

	.course-brief__right .course-brief__play-num .course-brief__icon {
		margin-right: 3px;
		line-height: 1;
		font-size: calc(12px * var(--font-scale, 1));
		font-weight: 400;
	}

	.course-brief__right .course-brief__play-num span {
		display: block;
		line-height: 16px;
		transform: scale(.91);
	}

	.course-brief__cover-modal {
		position: absolute;
		left: 0;
		bottom: 0;
		width: 100%;
		height: 18px;
		background: linear-gradient(180deg, transparent, rgba(0, 0, 0, .5));
		border-radius: 0 0 4px 4px;
	}

	.course-brief__center {
		box-sizing: border-box;
		display: flex;
		align-items: center;
		justify-content: space-around;
		height: 75px;
		padding: 14px 11px 13px 10px;
		border-bottom: 8px solid #f6f7fb;
	}

	.course-brief__center .course-brief__item {
		flex: 1;
		text-align: center;
	}

	.course-brief__center .course-brief__item .course-brief__label {
		margin-bottom: 1px;
		line-height: 18px;
		font-size: calc(13px * var(--font-scale, 1));
		font-weight: 500;
		color: #404040;
	}

	.course-brief__center .course-brief__item .course-brief__tip {
		line-height: 17px;
		font-size: calc(12px * var(--font-scale, 1));
		font-weight: 400;
		color: #7f848c;
	}

	.course-footer .course-footer__collect,
	.course-footer .course-footer__share {
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.course-footer {
		z-index: 100;
		position: fixed;
		bottom: 0;
		left: 0;
		display: flex;
		width: 100%;
		min-height: 54px;
		background: #fff;
		box-shadow: 0 0 15px 0 rgba(233, 234, 237, .65);
		white-space: nowrap;
	}

	.course-footer .course-footer__share {
		flex-direction: column;
	}

	.course-footer .course-footer__collect {
		flex-direction: column;
		min-width: 30px;
		color: #404040;
	}

	.course-footer .course-footer__icon {
		margin-bottom: 4px;
		line-height: 1;
		font-size: calc(18px * var(--font-scale, 1));
		font-weight: 400;
		color: #404040;
	}

	.course-footer span {
		line-height: 14px;
		white-space: nowrap;
		font-size: calc(10px * var(--font-scale, 1));
		font-weight: 400;
	}

	.course-footer__collect.active {
		color: #4285f4;
	}

	.course-footer .course-footer__btn {
		display: flex;
		flex-direction: column;
		width: 100%;
		min-width: 0;
		height: 54px;
		background: rgba(66, 133, 244, .8);
		font-size: calc(16px * var(--font-scale, 1));
		font-weight: 500;
		color: #fff;
	}

	.course-footer .course-footer__btn.is-disabled {
		background: #c5c9d1;
		pointer-events: none;
	}

	.course-footer .course-footer__label {
		margin-bottom: 2px;
		line-height: 20px;
		font-size: calc(16px * var(--font-scale, 1));
		font-weight: 500;
		color: #fff;
	}

	.course-footer .course-footer__tip {
		max-width: 68vw;
		overflow: hidden;
		text-overflow: ellipsis;
		line-height: 14px;
		font-size: calc(11px * var(--font-scale, 1));
		font-weight: 400;
		color: rgba(255, 255, 255, .9);
	}

	.course-footer__left {
		box-sizing: border-box;
		display: flex;
		align-items: center;
		justify-content: space-evenly;
		flex-shrink: 0;
		width: 90px;
		padding-top: 1px;
		padding-left: 3px;
	}

	/* 附件列表 */
	.attach-wrapper {
		padding: calc(20 * var(--rpx));
	}

	.attach-list {
		background: #fff;
	}

	.attach-item {
		display: flex;
		align-items: center;
		padding: calc(24 * var(--rpx));
		margin-bottom: calc(16 * var(--rpx));
		background: #f8f9fa;
		border-radius: calc(12 * var(--rpx));
	}

	.attach-item__icon {
		width: calc(60 * var(--rpx));
		height: calc(60 * var(--rpx));
		margin-right: calc(16 * var(--rpx));
	}

	.attach-item__info {
		flex: 1;
	}

	.attach-item__name {
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
		color: #333;
		margin-bottom: calc(8 * var(--rpx));
	}

	.attach-item__size {
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #999;
	}

	.attach-item__action {}

	.attach-item__action_icon {
		width: calc(60 * var(--rpx));
		height: calc(60 * var(--rpx));
	}

	.empty-tip {
		text-align: center;
		padding: calc(80 * var(--rpx)) 0;
		color: #999;
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	}

	.course-footer__right {
		flex: 1;
	}

	.intro-module:after {
		content: "";
		display: block;
		position: absolute;
		left: 0;
		bottom: 0;
		right: 0;
		height: 1px;
		background: url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAACCAYAAAHuhoaxAAAAAXNSR0IArs4c6QAAAERlWElmTU0AKgAAAAgAAYdpAAQAAAABAAAAGgAAAAAAA6ABAAMAAAABAAEAAKACAAQAAAABAAAAAaADAAQAAAABAAAAAgAAAAC+gucuAAAAEElEQVQIHWN4+ertfwYQAAAcAAPAK4oUCwAAAABJRU5ErkJggg==") repeat-x;
		background-size: 1px 1px;
	}

	.intro-module:last-child:after {
		display: none;
	}

	.intro-content {
		padding: 27px 20px 0;
		overflow: hidden;
	}

	.intro-module__item {
		overflow: hidden;
		word-break: break-word;
		overflow-wrap: anywhere;
	}

	.intro-module {
		position: relative;
		margin-bottom: 33px;
		padding-bottom: 24px;
	}

	.intro-module:after {
		right: -20px;
	}

	.intro-module .intro-module__lecturer-info {
		display: flex;
		margin-bottom: 18px;
		padding-top: 6px;
	}

	.intro-module .intro-module__lecturer-info .intro-module__lecturer-cover {
		overflow: hidden;
		flex-shrink: 0;
		width: 54px;
		height: 54px;
		margin-right: 12px;
		border-radius: 100%;
	}

	.intro-module .intro-module__lecturer-info .intro-module__lecturer-other {
		padding-top: 3px;
	}

	.intro-module .intro-module__lecturer-info .intro-module__lecturer-other .intro-module__lecturer-name {
		margin-bottom: 5px;
		line-height: 24px;
		font-size: calc(17px * var(--font-scale, 1));
		font-weight: 500;
		color: #353535;
	}

	.intro-module .intro-module__lecturer-info .intro-module__lecturer-other .intro-module__lecturer-subtitle {
		line-height: 18px;
		font-size: calc(13px * var(--font-scale, 1));
		font-weight: 400;
		color: #7f848c;
	}

	.intro-module__title {
		margin-bottom: 12px;
		line-height: 24px;
		font-size: calc(17px * var(--font-scale, 1));
		font-weight: 500;
		color: #353535;
	}

	.intro-module__title,
	.intro-module__title span {
		position: relative;
	}

	.intro-module__title:before {
		content: "";
		position: absolute;
		top: -2px;
		left: -9px;
		width: 16px;
		height: 16px;
		background: linear-gradient(1turn, #f3f8ff, #b3cefb);
		border-radius: 100%;
		opacity: .5;
	}

	.course-intro {
		flex-grow: 1;
		height: 100%;
		padding: 15px 0 54px;
		background: #fff;
	}

	/* u-tabs 样式调整 */
	.van-tabs {
		background: #fff;
		padding: 15px 46px 7px;
	}

	.van-tabs :deep(.van-tab) {
		font-size: calc(17px * var(--font-scale, 1));
	}

	.course-authority {
		box-sizing: border-box;
		padding: 18px 20px 0;
	}

	.iconfont {
		font-family: iconfont !important;
		font-size: calc(16px * var(--font-scale, 1));
		font-style: normal;
	}

	.catalog-article-item,
	.catalog-list {
		position: relative;
	}

	.catalog-article-item:after,
	.catalog-list:after {
		content: "";
		display: block;
		position: absolute;
		left: 0;
		bottom: 0;
		right: 0;
		height: 1px;
		background: url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAACCAYAAAHuhoaxAAAAAXNSR0IArs4c6QAAAERlWElmTU0AKgAAAAgAAYdpAAQAAAABAAAAGgAAAAAAA6ABAAMAAAABAAEAAKACAAQAAAABAAAAAaADAAQAAAABAAAAAgAAAAC+gucuAAAAEElEQVQIHWN4+ertfwYQAAAcAAPAK4oUCwAAAABJRU5ErkJggg==") repeat-x;
		background-size: 1px 1px;
	}

	.catalog-article-item:last-child:after,
	.catalog-list:last-child:after {
		display: none;
	}

	.catalog-wrapper {
		margin: 2px 0 60px;
	}

	.catalog-wrapper .catalog-chapter__icon--open {
		transform: rotate(180deg);
	}

	.catalog-list:after {
		left: 43px;
	}

	.catalog-chapter {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		padding: 18px 15px 18px;
		transition: padding .3s linear;
	}

	.catalog-chapter__title {
		line-height: 17px;
		font-size: calc(14px * var(--font-scale, 1));
		font-weight: 500;
		color: #404040;
	}

	.catalog-chapter__icon {
		line-height: 1;
		font-size: calc(14px * var(--font-scale, 1));
		font-weight: 400;
		color: #b9bdc2;
		transition: all .3s ease;
	}

	.catalog-article-item {
		position: relative;
		display: flex;
		justify-content: space-between;
		padding: 16px 15px 16px 0;
		margin-left: 43px;
	}

	.catalog-article-item .catalog-article-item__type-icon {
		position: absolute;
		left: -28px;
		top: 17px;
		line-height: 1;
		font-size: calc(18px * var(--font-scale, 1));
		font-weight: 400;
		color: #b9bdc2;
	}

	.catalog-article-item .catalog-article-item__info {
		flex: 1;
		min-width: 0;
		margin-right: 20px;
	}

	.catalog-article-item .catalog-article-item__info .catalog-article-item__title {
		margin-bottom: 6px;
		overflow: hidden;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
		line-height: 20px;
		font-size: calc(14px * var(--font-scale, 1));
		font-weight: 400;
		color: #404040;
	}

	.catalog-article-item .catalog-article-item__info .catalog-article-item__title[is-last=true],
	.catalog-article-item.is-current .catalog-article-item__title {
		color: #4285f4;
	}

	.catalog-article-item.is-current .catalog-article-item__try {
		background: #e8f1ff;
	}

	.catalog-article-item .catalog-article-item__info .catalog-article-item__bottom {
		display: flex;
		align-items: center;
	}

	.catalog-article-item .catalog-article-item__info .catalog-article-item__bottom span.catalog-article-item__duration {
		line-height: 17px;
		font-size: calc(12px * var(--font-scale, 1));
		font-weight: 400;
		color: #7f848c;
	}

	.catalog-article-item .catalog-article-item__status {
		white-space: nowrap;
	}

	.catalog-article-item .catalog-article-item__status .catalog-article-item__try {
		padding: 3px 10px;
		line-height: 16px;
		border-radius: 12px;
		background: #f3f8ff;
		font-size: calc(12px * var(--font-scale, 1));
		font-weight: 500;
		color: #4285f4;
		transform: translateY(-1px);
	}

	.catalog-article-item .catalog-article-item__status .catalog-article-item__try.is-locked {
		background: #f4f5f7;
		color: #9aa0a6;
	}

	.catalog-article-item .catalog-article-item__status .catalog-article-item__try.is-pending {
		background: #fff5f7;
		color: #ff3359;
	}

	.catalog-article-list {
		transition: height .3s ease;
	}

	/* 评论区域样式 */
	.comment-section {
		padding: calc(40 * var(--rpx)) 0;
		margin-top: calc(40 * var(--rpx));
		border-top: calc(20 * var(--rpx)) solid #f6f7fb;
	}

	.comment-list {
		padding: 0 calc(40 * var(--rpx));
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

	.reply-name {
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 500;
		color: #404040;
	}

	.reply-time {
		font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
		color: #7f848c;
		transform: scale(0.83);
		transform-origin: 0;
	}

	.reply-content {
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
		line-height: 1.5;
		color: #404040;
		word-break: break-word;
	}

	.reply-to {
		color: #4285f4;
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
		display: flex;
		align-items: center;
		padding: calc(16 * var(--rpx)) calc(32 * var(--rpx));
		padding-bottom: calc(calc(16 * var(--rpx)) + env(safe-area-inset-bottom));
		background: #fff;
		box-shadow: 0 calc(-2 * var(--rpx)) calc(16 * var(--rpx)) rgba(0, 0, 0, 0.08);
		z-index: 100;
	}

	.footer-collect {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		margin-right: calc(32 * var(--rpx));
		gap: calc(4 * var(--rpx));
	}

	.footer-collect span {
		font-size: calc(calc(20 * var(--rpx)) * var(--font-scale, 1));
		color: #404040;
	}

	.footer-collect.active span {
		color: #4285f4;
	}

	.footer-input {
		flex: 1;
		display: flex;
		align-items: center;
		gap: calc(12 * var(--rpx));
		padding: calc(16 * var(--rpx)) calc(24 * var(--rpx));
		background: #f5f5f5;
		border-radius: calc(40 * var(--rpx));
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
		color: #7f848c;
	}
</style>