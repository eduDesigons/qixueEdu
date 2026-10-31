<template>
	<div class="comment-page">

		<div class="comment-tabs">
			<van-tabs v-model:active="activeTab" color="#4285f4" line-width="22" line-height="3" @change="(i) => switchTab({ index: i })">
				<van-tab v-for="(t, i) in tabsList" :key="i" :title="t.name" :name="i" />
			</van-tabs>
		</div>

		<mescroll-body
			v-if="activeTab === 0"
			ref="courseMescroll"
			@init="mescrollInit"
			@down="downCallback"
			@up="upCallback"
			:down="downOption"
			:up="upOption"
		>
			<div class="comment-list">
				<div
					class="comment-item"
					v-for="(item, index) in courseCommentList"
					:key="item.id || index"
				>
					<div class="comment-header">
						<div class="comment-user">
							<van-image :src="static(item.user && item.user.avatar || '/static/default-avatar.png')" round :width="36 + 'px'" :height="36 + 'px'" fit="cover" />
							<div class="comment-user-info">
								<div class="comment-name"><span>{{ item.user && item.user.nick_name || '匿名' }}</span><coin-badge :coins="item.user ? item.user.coins : null" /></div>
							</div>
						</div>
						<span class="comment-time">{{ timeFrom(item.created_at) }}</span>
					</div>
					<div class="comment-content">{{ item.content || '' }}</div>

					<div v-if="item.course" class="comment-share-block" @click.stop="goToCourseDetail(item)">
						<div class="comment-share-block-inner">
							<img :src="static(item.course.thumb || '/static/img/video-img-1.png')" class="share-block-img" loading="lazy"  style="object-fit:cover" />
							<div class="share-block-info">
								<span class="share-block-title">{{ item.course.title || '' }}</span>
								<span class="share-block-subtitle" v-if="item.course.desc">{{ item.course.desc }}</span>
							</div>
						</div>
					</div>

					<div class="comment-replies" v-if="item.children && item.children.length > 0">
						<div class="reply-item" v-for="(reply, replyIndex) in item.children" :key="reply.id || replyIndex">
							<div class="reply-header">
								<div class="reply-name"><span>{{ reply.user && reply.user.nick_name || '匿名' }}</span><coin-badge :coins="reply.user ? reply.user.coins : null" /></div>
								<span class="reply-time">{{ timeFrom(reply.created_at) }}</span>
							</div>
							<div class="reply-content">
								<span class="reply-to" v-if="replyToName(reply, item)">回复 @{{ replyToName(reply, item) }}：</span>
								<span>{{ reply.content || '' }}</span>
							</div>
						</div>
					</div>
					<div class="comment-footer" v-if="item.reply_count > 0 || item.children && item.children.length > 0" @click.stop="viewReplies(item, 'course')">
						<span class="comment-count-text">共{{ (item.reply_count || (item.children && item.children.length) || 0) }}条评论 ></span>
					</div>
				</div>
			</div>
		</mescroll-body>

		<mescroll-body
			v-if="activeTab === 1"
			ref="chapterMescroll"
			@init="mescrollInit"
			@down="downCallback"
			@up="upCallback"
			:down="downOption"
			:up="upOption"
		>
			<div class="comment-list">
				<div
					class="comment-item"
					v-for="(item, index) in chapterCommentList"
					:key="item.id || index"
				>
					<div class="comment-header">
						<div class="comment-user">
							<van-image :src="static(item.user && item.user.avatar || '/static/default-avatar.png')" round :width="36 + 'px'" :height="36 + 'px'" fit="cover" />
							<div class="comment-user-info">
								<div class="comment-name"><span>{{ item.user && item.user.nick_name || '匿名' }}</span><coin-badge :coins="item.user ? item.user.coins : null" /></div>
							</div>
						</div>
						<span class="comment-time">{{ timeFrom(item.created_at) }}</span>
					</div>
					<div class="comment-content">{{ item.content || '' }}</div>

					<div v-if="item.chapter" class="comment-share-block" @click.stop="goToChapterDetail(item)">
						<div class="comment-share-block-inner">
							<img :src="static((item.course && item.course.thumb) || '/static/img/video-img-1.png')" class="share-block-img" loading="lazy"  style="object-fit:cover" />
							<div class="share-block-info">
								<span class="share-block-title">{{ item.chapter.title || '' }}</span>
								<span class="share-block-subtitle" v-if="item.course && item.course.title">{{ item.course.title }}</span>
							</div>
						</div>
					</div>

					<div class="comment-replies" v-if="item.children && item.children.length > 0">
						<div class="reply-item" v-for="(reply, replyIndex) in item.children" :key="reply.id || replyIndex">
							<div class="reply-header">
								<div class="reply-name"><span>{{ reply.user && reply.user.nick_name || '匿名' }}</span><coin-badge :coins="reply.user ? reply.user.coins : null" /></div>
								<span class="reply-time">{{ timeFrom(reply.created_at) }}</span>
							</div>
							<div class="reply-content">
								<span class="reply-to" v-if="replyToName(reply, item)">回复 @{{ replyToName(reply, item) }}：</span>
								<span>{{ reply.content || '' }}</span>
							</div>
						</div>
					</div>
					<div class="comment-footer" v-if="item.reply_count > 0 || item.children && item.children.length > 0" @click.stop="viewReplies(item, 'chapter')">
						<span class="comment-count-text">共{{ (item.reply_count || (item.children && item.children.length) || 0) }}条评论 ></span>
					</div>
				</div>
			</div>
		</mescroll-body>
	</div>
</template>

<script>
import MescrollMixin from "@/uni_modules/mescroll-uni/components/mescroll-uni/mescroll-mixins.js";

export default {
	mixins: [MescrollMixin],
	data() {
		return {
			activeTab: 0,
			tabsList: [
				{ name: '课程评论' },
				{ name: '章节评论' }
			],
			courseCommentList: [],
			chapterCommentList: [],
			downOption: {
				auto: false
			},
			upOption: {
				auto: true,
				page: {
					num: 0,
					size: 10
				},
				noMoreSize: 4,
				empty: {
					tip: '暂无评论数据'
				}
			}
		}
	},
	onLoad() {
	},
	methods: {
		switchTab(e) {
			this.activeTab = e.index;
			this.$nextTick(() => {
				if (this.mescroll) {
					this.mescroll.resetUpScroll();
				}
			});
		},
		downCallback() {
			this.mescroll.resetUpScroll();
		},
		upCallback(page) {
			if (this.activeTab === 0) {
				this.getCourseCommentList(page);
			} else {
				this.getChapterCommentList(page);
			}
		},
		getCourseCommentList(page) {
			this.$http.get('/user/myCourseComment', {
				page: page.num,
				limit: page.size
			}).then(res => {
				const payload = (res && res.data) || {}
				const list = Array.isArray(payload.list) ? payload.list : []
				if (page.num === 1) {
					this.courseCommentList = []
				}
				this.courseCommentList = this.courseCommentList.concat(list)
				const hasMore = payload.hasMore === true || list.length >= page.size
				this.mescroll.endSuccess(list.length, hasMore)
			}).catch(err => {
				console.log('获取课程评论失败:', err)
				this.msg((err && (err.message || err.msg)) || '获取课程评论失败')
				this.mescroll.endErr()
			})
		},
		getChapterCommentList(page) {
			this.$http.get('/user/myChapterComment', {
				page: page.num,
				limit: page.size
			}).then(res => {
				const payload = (res && res.data) || {}
				const list = Array.isArray(payload.list) ? payload.list : []
				if (page.num === 1) {
					this.chapterCommentList = []
				}
				this.chapterCommentList = this.chapterCommentList.concat(list)
				const hasMore = payload.hasMore === true || list.length >= page.size
				this.mescroll.endSuccess(list.length, hasMore)
			}).catch(err => {
				console.log('获取章节评论失败:', err)
				this.msg((err && (err.message || err.msg)) || '获取章节评论失败')
				this.mescroll.endErr()
			})
		},
		replyToName(reply, root) {
			if (reply && reply.parent && reply.parent.user && reply.parent.user.nick_name) {
				return reply.parent.user.nick_name;
			}
			return (root && root.user && root.user.nick_name) || '';
		},
		goToCourseDetail(item) {
			const courseId = item.course_id || (item.course && item.course.id);
			if (courseId) {
				this.navTo(`/pages/course/courseDetail?id=${courseId}`);
			}
		},
		goToChapterDetail(item) {
			const chapterId = item.chapter_id || (item.chapter && item.chapter.id);
			if (chapterId) {
				const courseId = (item.chapter && item.chapter.course_id) || (item.chapter && item.chapter.course && item.chapter.course.id);
				const url = courseId ? `/pages/course/chapterDetail?id=${chapterId}&courseId=${courseId}` : `/pages/course/chapterDetail?id=${chapterId}`;
				this.navTo(url);
			}
		},
		viewReplies(item, type) {
			if (type === 'course') {
				const courseId = item.course_id || (item.course && item.course.id);
				if (courseId) {
					this.navTo(`/pages/course/courseDetail?id=${courseId}&tab=3`);
				}
				return;
			}
			this.goToChapterDetail(item);
		}
	}
}
</script>

<style lang="scss" scoped>
.comment-page {
	min-height: 100vh;
	background: #ffffff;
}

.comment-tabs {
	background: #fff;
	position: sticky;
	top: 0;
	z-index: 10;
}

.comment-list {
	padding: calc(20 * var(--rpx)) calc(30 * var(--rpx));
}

.comment-item {
	padding: calc(32 * var(--rpx)) 0;
	border-bottom: calc(1 * var(--rpx)) solid #f0f0f0;

	&:last-child {
		border-bottom: none;
	}
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
	margin-left: calc(16 * var(--rpx));
	display: flex;
	align-items: center;
}

.comment-name {
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
	line-height: 1;
}

.comment-time {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #b2b2b2;
	line-height: 1;
	flex-shrink: 0;
}

.comment-content {
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	color: #353535;
	line-height: 1.6;
	margin-bottom: calc(20 * var(--rpx));
	word-break: break-word;
}

.comment-share-block {
	margin-top: calc(20 * var(--rpx));
	margin-bottom: calc(20 * var(--rpx));
	border-radius: calc(12 * var(--rpx));
	overflow: hidden;
	background: #f6f7fb;
}

.comment-share-block-inner {
	display: flex;
	align-items: center;
	padding: calc(24 * var(--rpx));
}

.share-block-img {
	width: calc(120 * var(--rpx));
	height: calc(120 * var(--rpx));
	border-radius: calc(8 * var(--rpx));
	flex-shrink: 0;
	margin-right: calc(20 * var(--rpx));
}

.share-block-info {
	flex: 1;
	display: flex;
	flex-direction: column;
	gap: calc(8 * var(--rpx));
}

.share-block-title {
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #404040;
	line-height: 1.4;
	overflow: hidden;
	display: -webkit-box;
	-webkit-line-clamp: 1;
	line-clamp: 1;
	-webkit-box-orient: vertical;
}

.share-block-subtitle {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
	line-height: 1.4;
	overflow: hidden;
	display: -webkit-box;
	-webkit-line-clamp: 1;
	line-clamp: 1;
	-webkit-box-orient: vertical;
}

.comment-share-nested {
	background: #ffffff;
	border-radius: calc(12 * var(--rpx));
	padding: calc(20 * var(--rpx));
}

.share-nested-header {
	margin-bottom: calc(16 * var(--rpx));
	padding-bottom: calc(12 * var(--rpx));
	border-bottom: calc(1 * var(--rpx)) solid #f0f0f0;
}

.share-nested-user {
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: #353535;
	font-weight: 500;
}

.share-nested-content {
	display: flex;
	align-items: center;
}

.comment-replies {
	margin-top: calc(20 * var(--rpx));
	padding: calc(20 * var(--rpx));
	background: #f5f5f5;
	border-radius: calc(8 * var(--rpx));
}

.reply-item {
	margin-bottom: calc(16 * var(--rpx));

	&:last-child {
		margin-bottom: 0;
	}
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
	color: #b2b2b2;
}

.reply-content {
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: #404040;
	line-height: 1.6;
	word-break: break-word;
}

.reply-to {
	color: #4285f4;
}

.comment-footer {
	margin-top: calc(20 * var(--rpx));
	text-align: right;
}

.comment-count-text {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #b2b2b2;
}
</style>
