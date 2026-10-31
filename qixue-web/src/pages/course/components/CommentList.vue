<template>
	<div class="comment-section">
		<div class="comment-title" v-if="title">{{ title }}</div>
		<div class="comment-list">
			<div class="empty-tip" v-if="list.length === 0 && !loading">
				<van-empty description="暂无讨论" image="/static/empty/comment.png" :image-size="emptySize" />
			</div>
			<div class="comment-item" v-for="(item, index) in list" :key="item.id || index">
				<div class="comment-header">
					<div class="comment-user">
						<van-image :src="static((item.user && item.user.avatar) || '/static/default-avatar.png')" round :width="36 + 'px'" :height="36 + 'px'" fit="cover" />
						<div class="comment-user-info">
							<div class="comment-name"><span>{{ (item.user && item.user.nick_name) || '匿名' }}</span><coin-badge :coins="item.user ? item.user.coins : null" /></div>
							<div class="comment-time">{{ date(item.created_at, 'yyyy-mm-dd') }}</div>
						</div>
					</div>
					<div class="comment-actions">
						<div class="comment-reply-btn" @click.stop="$emit('reply', item, item)">
							<van-icon name="chat-o" size="14" color="#7f848c" />
							<span>回复</span>
						</div>
					</div>
				</div>
				<div class="comment-content">{{ item.content || '' }}</div>
				<div class="comment-replies" v-if="item.children && item.children.length > 0">
					<div class="reply-item" v-for="(reply, replyIndex) in item.children" :key="reply.id || replyIndex">
						<div class="reply-header">
							<div class="reply-name"><span>{{ (reply.user && reply.user.nick_name) || '匿名' }}</span><coin-badge :coins="reply.user ? reply.user.coins : null" /></div>
							<span class="reply-time">{{ date(reply.created_at, 'yyyy-mm-dd') }}</span>
							<div class="comment-reply-btn" v-if="replyOnChild" @click.stop="$emit('reply', reply, item)">
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
			<div class="comment-loading" v-if="loading">
				<van-loading type="circular" />
				<span>加载中...</span>
			</div>
			<div class="comment-no-more" v-if="hasMore === false && list.length > 0">
				<span>没有更多了</span>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	name: 'CommentList',
	props: {
		list: { type: Array, default: () => [] },
		loading: { type: Boolean, default: false },
		hasMore: { type: Boolean, default: true },
		title: { type: String, default: '' },
		emptySize: { type: Number, default: 80 },
		replyOnChild: { type: Boolean, default: false },
	},
	emits: ['reply'],
	methods: {
		replyToName(reply, root) {
			if (reply.parent && reply.parent.user && reply.parent.user.nick_name) {
				return reply.parent.user.nick_name
			}
			return (root && root.user && root.user.nick_name) || ''
		},
	},
}
</script>

<style lang="scss" scoped>
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
</style>
