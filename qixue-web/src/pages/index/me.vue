<template>
		<MeDesktop v-if="isDesktop" ref="meDesktop" :stats="stats" :unread="unreadCount" />
	<div v-else class="me-page topBg">
		<div class="user-center">
			<div class="user-info-content" @click="navTo('/pages/member/profile')">
				<div class="user-avatar">
					<div class="avatar-container">
						<div class="avatar">
							<van-image v-if="userInfo.avatar" :src="static(userInfo.avatar)" round :width="60 + 'px'" :height="60 + 'px'" fit="cover" />
							<span class="qx-avatar" v-else :style="{width:60+'px',height:60+'px',fontSize:Math.round(60/2.5)+'px'}">{{ (userInfo.nick_name || '').slice(0,1) }}</span>
						</div>
					</div>
				</div>
				<div class="user-info">
					<div class="user-info__inner">
						<div class="user-name">
							<span class="user-name__text">{{ userInfo.nick_name || '未登录' }}</span>
							<coin-badge v-if="userInfo && userInfo.id" :coins="userInfo.coins != null ? userInfo.coins : stats.coins" />
						</div>
						<div class="user-enterprise" v-if="userInfo.mobile">
							<span>{{ userInfo.mobile }}</span>
						</div>
					</div>
					<div class="user-msg" @click.stop="navTo('/pages/member/message')">
						<van-icon name="chat-o" size="22" color="#7f848c" />
						<van-badge v-if="unreadCount > 0" :content="unreadCount > 99 ? '99+' : unreadCount" style="position:absolute;top:-4px;right:-6px;z-index:2" />
					</div>
					<van-icon name="arrow" size="18" color="#b9bdc2" @click="navTo('/pages/member/profile')" />
				</div>
			</div>
			<div class="learn-stats">
				<div class="learn-stats__item" @click="navTo('/pages/member/lesson')">
					<div>{{ stats.courseCount || 0 }}</div>
					<span>我的课程</span>
				</div>
				<div class="learn-stats__item" @click="navTo('/pages/member/credit')">
					<div>{{ stats.credit || 0 }}</div>
					<span>我的学分</span>
				</div>
				<div class="learn-stats__item" @click="navTo('/pages/member/coin')">
					<div>{{ stats.coins || 0 }}</div>
					<span>我的积分</span>
				</div>
			</div>
		</div>

		<div class="my-content">
			<div class="section-title">我的内容</div>
			<div class="my-content__inner">
				<div class="content-item-wrap" @click="navTo('/pages/member/collection')">
					<div class="content-item">
						<div class="content-item__icon">
							<img src="/static/img/me/collect.png"  style="object-fit:cover" />
						</div>
						<div class="content-item__title">我的收藏</div>
					</div>
				</div>
				<div class="content-item-wrap" @click="navTo('/pages/member/exam')">
					<div class="content-item">
						<div class="content-item__icon">
							<img src="/static/img/me/test.png"  style="object-fit:cover" />
						</div>
						<div class="content-item__title">考试记录</div>
					</div>
				</div>
				
				<div class="content-item-wrap" @click="navTo('/pages/member/test')">
					<div class="content-item">
						<div class="content-item__icon">
							<img src="/static/img/me/note.png"  style="object-fit:cover" />
						</div>
						<div class="content-item__title">测评记录</div>
					</div>
				</div>
				<div class="content-item-wrap" @click="navTo('/pages/member/comment')">
					<div class="content-item">
						<div class="content-item__icon">
							<img src="/static/img/me/commit.png"  style="object-fit:cover" />
						</div>
						<div class="content-item__title">我的评论</div>
					</div>
				</div>
				<div class="content-item-wrap" @click="navTo('/pages/member/history')">
					<div class="content-item">
						<div class="content-item__icon">
							<img src="/static/img/me/history.png"  style="object-fit:cover" />
						</div>
						<div class="content-item__title">学习轨迹</div>
					</div>
				</div>
				<div class="content-item-wrap" @click="navTo('/pages/member/path')">
					<div class="content-item">
						<div class="content-item__icon">
							<img src="/static/img/me/path.png"  style="object-fit:cover" />
						</div>
						<div class="content-item__title">系列课程</div>
					</div>
				</div>
				<div class="content-item-wrap" @click="navTo('/pages/member/cert')">
					<div class="content-item">
						<div class="content-item__icon">
							<img src="/static/img/me/certificate.png"  style="object-fit:cover" />
						</div>
						<div class="content-item__title">我的证书</div>
					</div>
				</div>
				<div class="content-item-wrap" @click="navTo('/pages/member/message')">
					<div class="content-item">
						<div class="content-item__icon">
							<img src="/static/img/me/notice.png"  style="object-fit:cover" />
							<van-badge v-if="unreadCount > 0" :content="unreadCount > 99 ? '99+' : unreadCount" style="position:absolute;top:-4px;right:-6px;z-index:2" />
						</div>
						<div class="content-item__title">我的消息</div>
					</div>
				</div>
			</div>
		</div>

		<div class="common-func">
			<div class="common-func__title">常用功能</div>
			<div class="common-func__inner">
				<div class="common-item-wrap" @click="navTo('/pages/member/question')">
					<div class="common-item">
						<div class="common-item__icon common-item__icon--aibot">
							<img src="/static/img/me/question.png"  style="object-fit:cover" />
						</div>
						<div class="common-item__title">常见问题</div>
					</div>
				</div>
				<div class="common-item-wrap" @click="navTo('/pages/member/tools')">
					<div class="common-item">
						<div class="common-item__icon">
							<van-icon name="apps-o" size="26" color="#8a9199" />
						</div>
						<div class="common-item__title">工具库</div>
					</div>
				</div>
				<div class="common-item-wrap" @click="navTo('/pages/member/set')">
					<div class="common-item">
						<div class="common-item__icon">
							<img src="/static/img/me/set.png"  style="object-fit:cover" />
						</div>
						<div class="common-item__title">设置</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import { getStorage, setStorage } from '@/common/storage'
import { defineAsyncComponent } from 'vue'
import { isDesktop as desktopOn } from '@/pc/desktop'

export default {
	components: {
		MeDesktop: defineAsyncComponent(() => import('@/pc/MeDesktop.vue')),
	},
	data() {
		return {
			stats: {
				courseCount: 0,
				credit: 0,
				coins: 0
			}
		};
	},
	computed: {
		isDesktop() {
			return desktopOn.value
		},
		unreadCount() {
			return this.$store.state.unreadMessageCount || 0;
		}
	},
	onLoad() {
		this.restoreStats();
	},
	onShow() {
		if (this.$store.getters.hasLogin) {
			this.getUserStats();
			this.$nextTick(() => {
				const desk = this.$refs.meDesktop
				if (desk && desk.loadTrend) desk.loadTrend()
			})
		}
	},
	methods: {
		restoreStats() {
			try {
				const cached = getStorage('qixue_me_stats');
				if (cached) {
					this.stats = {
						courseCount: cached.courseCount || 0,
						credit: cached.credit || 0,
						coins: cached.coins || 0
					};
				}
			} catch (e) {}
		},
		getUserStats() {
			this.$http.get('/user/myCount').then(res => {
				if (res.data) {
					this.stats = {
						courseCount: res.data.courseCount || 0,
						credit: res.data.credit || 0,
						coins: res.data.coins != null ? res.data.coins : 0
					};
					if (this.userInfo && this.userInfo.id) {
						this.$store.commit('setUserInfo', Object.assign({}, this.userInfo, { coins: this.stats.coins }))
					}
					try {
						setStorage('qixue_me_stats', this.stats);
					} catch (e) {}
					const unread = res.data.unreadMessageCount || 0;
					this.$store.dispatch('applyMessageBadge', unread);
				}
			}).catch(err => {
				console.log('获取统计数据失败:', err);
			});
		}
	}
}
</script>

<style lang="scss">
.qx-route-pages-index-me {
	background: #F6F7FB url('/static/bg/com.png') no-repeat 0 0 / contain fixed;
	padding-bottom: calc(40 * var(--rpx));
}

.user-center {
	padding: calc(100 * var(--rpx)) calc(30 * var(--rpx)) calc(34 * var(--rpx)) calc(30 * var(--rpx));
	background-size: cover;
}

.user-info-content {
	display: flex;
	padding: 0 calc(30 * var(--rpx));
	align-items: center;
	overflow: visible;
}

.me-page .user-avatar {
	flex: none;
	width: auto;
	height: auto;
	border-radius: 0;
}

.avatar-container {
	position: relative;
	flex: none;
}

.me-page .avatar {
	display: flex;
	flex: none;
	width: 60px;
	height: 60px;
	box-sizing: content-box;
	border-radius: 50%;
	border: 2px solid #fff;
	overflow: hidden;
}

.me-page .avatar .van-image,
.me-page .avatar .qx-avatar,
.me-page .avatar .van-image {
	width: 60px;
	height: 60px;
	flex: none;
	border: 0;
	border-radius: 50%;
}

.user-info {
	display: flex;
	align-items: center;
	flex: 1;
	min-width: 0;
	padding-left: calc(20 * var(--rpx));
	overflow: visible;
}

.user-msg {
	position: relative;
	flex-shrink: 0;
	width: calc(56 * var(--rpx));
	height: calc(56 * var(--rpx));
	margin-right: calc(8 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
	overflow: visible;
}

.user-info__inner {
	flex: 1;
	min-width: 0;
}

.user-name {
	display: flex;
	align-items: center;
	margin-bottom: calc(8 * var(--rpx));
}

.user-name__text {
	font-size: calc(calc(40 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
	line-height: calc(56 * var(--rpx));
}

.user-enterprise span {
	overflow: hidden;
	word-break: keep-all;
	white-space: nowrap;
	text-overflow: ellipsis;
	line-height: calc(34 * var(--rpx));
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 400;
	color: #7f848c;
}

.learn-stats {
	display: flex;
	margin-top: calc(42 * var(--rpx));
}

.learn-stats__item {
	display: flex;
	align-items: center;
	justify-content: center;
	flex-direction: column;
	flex: 1;
	margin: calc(6 * var(--rpx)) 0 0 calc(6 * var(--rpx));
}

.learn-stats__item > div {
	line-height: calc(52 * var(--rpx));
	color: #404040;
	font-size: calc(calc(44 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 400;
	font-family: DIN;
}

.learn-stats__item > span {
	padding-top: calc(4 * var(--rpx));
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #7f848c;
	line-height: calc(36 * var(--rpx));
}

.me-page .section-title {
	margin: 0 calc(12 * var(--rpx));
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
	line-height: calc(44 * var(--rpx));
}

.my-content {
	position: relative;
	z-index: 1;
	margin: calc(20 * var(--rpx)) calc(30 * var(--rpx));
	box-sizing: border-box;
	padding: calc(30 * var(--rpx)) calc(20 * var(--rpx));
	background: #fff;
	border-radius: calc(24 * var(--rpx));
}

.my-content__inner {
	display: flex;
	flex-wrap: wrap;
}

.content-item-wrap {
	flex: 0 0 25%;
	min-width: 0;
	display: flex;
	justify-content: center;
}

.content-item {
	width: 100%;
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	margin-top: calc(38 * var(--rpx));
	overflow: visible;
}

.content-item__icon {
	position: relative;
	width: calc(56 * var(--rpx));
	height: calc(56 * var(--rpx));
	overflow: visible;
}

.content-item__icon img {
	width: 100%;
	height: 100%;
}

.content-item__title {
	width: 100%;
	box-sizing: border-box;
	padding: calc(16 * var(--rpx)) calc(6 * var(--rpx)) calc(10 * var(--rpx));
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 400;
	color: #404040;
	line-height: 1.35;
	text-align: center;
}

.common-func {
	box-sizing: border-box;
	padding: calc(30 * var(--rpx)) calc(20 * var(--rpx));
	background: #fff;
	border-radius: calc(24 * var(--rpx));
	margin: calc(24 * var(--rpx)) calc(30 * var(--rpx));
}

.common-func__title {
	margin: 0 calc(12 * var(--rpx));
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
	line-height: calc(44 * var(--rpx));
}

.common-func__inner {
	display: flex;
	flex-wrap: wrap;
}

.common-item-wrap {
	flex: 0 0 25%;
	min-width: 0;
	display: flex;
	justify-content: center;
}

.common-item {
	width: 100%;
	display: flex;
	flex-direction: column;
	align-items: center;
	margin-top: calc(40 * var(--rpx));
}

.common-item__icon {
	display: flex;
	align-items: center;
	justify-content: center;
	width: calc(56 * var(--rpx));
	height: calc(56 * var(--rpx));
	position: relative;
}

.common-item__icon img {
	width: calc(52 * var(--rpx));
	height: calc(52 * var(--rpx));
}

.common-item__icon--aibot::after {
	content: "";
	position: absolute;
	top: calc(-4 * var(--rpx));
	right: calc(-20 * var(--rpx));
	width: calc(36 * var(--rpx));
	height: calc(24 * var(--rpx));
}

.common-item__title {
	width: 100%;
	box-sizing: border-box;
	padding: calc(14 * var(--rpx)) calc(6 * var(--rpx)) calc(10 * var(--rpx));
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 400;
	color: #404040;
	line-height: 1.35;
	text-align: center;
}

.nav-set {
	width: calc(50 * var(--rpx));
	height: calc(50 * var(--rpx));
	padding: calc(20 * var(--rpx));
	border-radius: 50%;
}
</style>