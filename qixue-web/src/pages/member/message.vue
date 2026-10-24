<template>
	<div class="message-page">
		<div class="message-toolbar" v-if="list.length || unreadCount">
			<span class="message-toolbar-tip">未读 {{ unreadCount }} 条</span>
			<span class="message-toolbar-action" v-if="unreadCount" @click="readAll">全部已读</span>
		</div>
		<mescroll-body
			ref="msgMescroll"
			@init="mescrollInit"
			@down="downCallback"
			@up="upCallback"
			:down="downOption"
			:up="upOption"
		>
			<div class="message-list">
				<div
					class="message-item"
					:class="{ unread: !Number(item.is_read) }"
					v-for="(item, index) in list"
					:key="item.id || index"
					@click="openDetail(item)"
				>
					<div class="message-item-icon" :class="{ paper: item.type === 'paper' }">
						<van-icon :name="item.type === 'paper' ? 'description' : 'bell'" color="#fff" size="20" />
						<div v-if="!Number(item.is_read)" class="message-dot"></div>
					</div>
					<div class="message-item-body">
						<div class="message-item-head">
							<span class="message-item-title">{{ item.title || typeLabel(item.type) }}</span>
							<span class="message-item-time">{{ timeFrom(item.created_at) }}</span>
						</div>
						<div class="message-item-content">{{ plainText(item.content) }}</div>
					</div>
				</div>
			</div>
		</mescroll-body>

		<van-popup :show="showDetail" position="bottom" round="16" @update:show="(v)=>{if(!v)showDetail = false}">
			<div class="message-detail" v-if="current">
				<div class="message-detail-title">{{ current.title }}</div>
				<div class="message-detail-meta">
					<span>{{ typeLabel(current.type) }}</span>
					<span>{{ current.created_at }}</span>
				</div>
				<div class="message-detail-content">
					<div class="qx-parse" style="user-select:text" v-html="richHtml(current.content)"></div>
				</div>
				<div
					class="message-detail-link"
					v-if="detailLink"
					@click="goRelated"
				>{{ detailLink.text }}</div>
				<div class="message-detail-close" @click="showDetail = false">关闭</div>
			</div>
		</van-popup>
	</div>
</template>

<script>
import MescrollMixin from "@/uni_modules/mescroll-uni/components/mescroll-uni/mescroll-mixins.js";

export default {
	mixins: [MescrollMixin],
	data() {
		return {
			list: [],
			showDetail: false,
			current: null,
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
					tip: '暂无站内消息',
					icon: '/static/empty/empty.png'
				}
			}
		}
	},
	computed: {
		unreadCount() {
			return this.$store.state.unreadMessageCount || 0
		},
		detailLink() {
			const item = this.current
			if (!item) return null
			const type = item.related_type
			const id = item.related_id
			if (type === 'user_cert') {
				return { text: '查看证书', url: '/pages/member/cert' }
			}
			if ((type === 'user_paper' || type === 'paper' || item.type === 'paper') && id) {
				return { text: '查看成绩', url: '/pages/exam/examDetail?id=' + id }
			}
			if (type === 'course' && id) {
				return { text: '查看课程', url: '/pages/course/courseDetail?id=' + id }
			}
			if (type === 'chapter' && id) {
				return { text: '查看课节', url: '/pages/course/chapterDetail?id=' + id }
			}
			return null
		}
	},
	methods: {
		typeLabel(type) {
			return type === 'paper' ? '作业批改' : '系统通知'
		},
		goRelated() {
			if (!this.detailLink) return
			this.showDetail = false
			this.navTo(this.detailLink.url)
		},
		plainText(html) {
			if (!html) return ''
			return String(html)
				.replace(/<style[\s\S]*?<\/style>/gi, '')
				.replace(/<script[\s\S]*?<\/script>/gi, '')
				.replace(/<[^>]+>/g, '')
				.replace(/&nbsp;/gi, ' ')
				.replace(/&amp;/gi, '&')
				.replace(/&lt;/gi, '<')
				.replace(/&gt;/gi, '>')
				.replace(/\s+/g, ' ')
				.trim()
		},
		downCallback() {
			this.mescroll.resetUpScroll()
		},
		upCallback(page) {
			this.$http.get('/user/myMessages', {
				page: page.num,
				limit: page.size
			}).then(res => {
				const rows = (res.data && res.data.list) || []
				if (page.num === 1) {
					this.list = []
				}
				this.list = this.list.concat(rows)
				this.mescroll.endSuccess(rows.length, rows.length >= page.size)
			}).catch(() => {
				this.msg('获取消息失败')
				this.mescroll.endErr()
			})
		},
		refreshUnread() {
			this.$http.get('/user/unreadMessageCount').then(res => {
				const unread = (res.data && res.data.count) || 0
				this.$store.dispatch('applyMessageBadge', unread)
			}).catch(() => {})
		},
		openDetail(item) {
			this.current = item
			this.showDetail = true
			if (!Number(item.is_read)) {
				this.$http.post('/user/readMessage', { id: item.id }).then(() => {
					this.$set(item, 'is_read', 1)
					this.refreshUnread()
				}).catch(() => {})
			}
		},
		readAll() {
			this.$http.post('/user/readMessage', { id: 0 }).then(() => {
				this.list.forEach(item => {
					this.$set(item, 'is_read', 1)
				})
				this.refreshUnread()
				this.msg('已全部标记为已读')
			}).catch(() => {
				this.msg('操作失败')
			})
		}
	}
}
</script>

<style lang="scss" scoped>
.message-page {
	min-height: 100vh;
	background: #f6f7fb;
}

.message-toolbar {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: calc(20 * var(--rpx)) calc(30 * var(--rpx)) 0;
}

.message-toolbar-tip {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
}

.message-toolbar-action {
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: #4286f5;
}

.message-list {
	padding: calc(20 * var(--rpx)) calc(30 * var(--rpx));
}

.message-item {
	display: flex;
	padding: calc(28 * var(--rpx)) calc(24 * var(--rpx));
	background: #fff;
	border-radius: calc(16 * var(--rpx));
	margin-bottom: calc(20 * var(--rpx));

	&.unread {
		background: #f4f8ff;
	}
}

.message-item-icon {
	position: relative;
	width: calc(80 * var(--rpx));
	height: calc(80 * var(--rpx));
	border-radius: calc(20 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
	background: #3da4de;
	flex-shrink: 0;

	&.paper {
		background: #f5a623;
	}
}

.message-dot {
	position: absolute;
	top: calc(-6 * var(--rpx));
	right: calc(-6 * var(--rpx));
	width: calc(16 * var(--rpx));
	height: calc(16 * var(--rpx));
	border-radius: 50%;
	background: #f44336;
	border: calc(2 * var(--rpx)) solid #fff;
}

.message-item-body {
	flex: 1;
	margin-left: calc(20 * var(--rpx));
	min-width: 0;
}

.message-item-head {
	display: flex;
	align-items: flex-start;
	justify-content: space-between;
}

.message-item-title {
	flex: 1;
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
	line-height: calc(42 * var(--rpx));
	padding-right: calc(16 * var(--rpx));
}

.message-item-time {
	flex-shrink: 0;
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	color: #b9bdc2;
	line-height: calc(42 * var(--rpx));
}

.message-item-content {
	margin-top: calc(8 * var(--rpx));
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
	line-height: calc(36 * var(--rpx));
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.message-detail {
	padding: calc(40 * var(--rpx)) calc(32 * var(--rpx)) calc(60 * var(--rpx));
}

.message-detail-title {
	font-size: calc(calc(34 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: #353535;
	line-height: calc(48 * var(--rpx));
}

.message-detail-meta {
	display: flex;
	justify-content: space-between;
	margin-top: calc(16 * var(--rpx));
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #b9bdc2;
}

.message-detail-content {
	margin-top: calc(28 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #404040;
	line-height: 1.8;
	word-break: break-word;
}

.message-detail-link {
	margin-top: calc(32 * var(--rpx));
	height: calc(80 * var(--rpx));
	line-height: calc(80 * var(--rpx));
	text-align: center;
	background: #e8f1fe;
	color: #4286f5;
	border-radius: calc(40 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
}

.message-detail-close {
	margin-top: calc(48 * var(--rpx));
	height: calc(80 * var(--rpx));
	line-height: calc(80 * var(--rpx));
	text-align: center;
	background: #4286f5;
	color: #fff;
	border-radius: calc(40 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
}
</style>
