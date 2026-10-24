<template>
	<van-popup :show="show" position="center" :close-on-click-overlay="false" style="background:transparent" @update:show="(v)=>{if(!v)handleClose()}">
		<div class="notice-popup-wrapper">
			<div class="notice-popup-container">
				<div class="notice-popup-main">
					<div class="notice-popup-username">{{ userInfo.nick_name || '您好' }}<coin-badge v-if="userInfo && userInfo.id" :coins="userInfo.coins" /></div>
					<div class="notice-popup-company" v-if="notice.title">
						<span class="notice-popup-company-name">{{ notice.title }}</span>
					</div>
					<div class="notice-popup-text" v-if="notice.content">
						<div class="qx-parse" style="user-select:text" v-html="richHtml(notice.content)"></div>
					</div>
				</div>
                <div class="notice-popup-button" @click="handleLearn"> 开始学习 </div>
			</div>
		</div>
	</van-popup>
</template>

<script>
export default {
	name: 'NoticePopup',
    props: {
        notice: {
            type: Object,
            default: () => {}
        }
    },
	data() {
		return {
			show: false,
		};
	},
	methods: {
		handleClose() {
			this.show = false;
		},
		handleLearn() {
			this.handleClose();
			this.$emit('confirm');
		}
	}
};
</script>

<style lang="scss" scoped>
.notice-popup-wrapper {
	width: calc(640 * var(--rpx));
	padding: calc(12 * var(--rpx));
	background: hsla(0, 0%, 100%, 0.2);
	border: calc(2 * var(--rpx)) solid hsla(0, 0%, 100%, 0.2);
	border-radius: calc(48 * var(--rpx)) calc(108 * var(--rpx)) calc(48 * var(--rpx)) calc(48 * var(--rpx));
	position: relative;
	box-sizing: border-box;
}

.notice-popup-wrapper::before {
	content: '';
	position: absolute;
	left: calc(-12 * var(--rpx));
	top: calc(-6 * var(--rpx));
	width: calc(220 * var(--rpx));
	height: calc(162 * var(--rpx));
	background: url('/static/img/welcome-hello.png') no-repeat center;
	background-size: calc(220 * var(--rpx)) calc(162 * var(--rpx));
	z-index: 20;
}

.notice-popup-wrapper::after {
	content: '';
	position: absolute;
	right: calc(-48 * var(--rpx));
	top: calc(-34 * var(--rpx));
	width: calc(262 * var(--rpx));
	height: calc(250 * var(--rpx));
	background: url('/static/img/welcome-star.png') no-repeat center;
	background-size: calc(262 * var(--rpx)) calc(250 * var(--rpx));
	z-index: 20;
}

.notice-popup-container {
	position: static;
	box-sizing: border-box;
	padding: calc(140 * var(--rpx)) calc(60 * var(--rpx)) calc(72 * var(--rpx));
	background: #fff url('/static/img/welcome-bg-pc.png') no-repeat top;
	background-size: calc(616 * var(--rpx)) auto;
	border-radius: calc(48 * var(--rpx)) calc(108 * var(--rpx)) calc(48 * var(--rpx)) calc(48 * var(--rpx));
}

.notice-popup-main {
	padding: 0;
	margin-bottom: calc(64 * var(--rpx));
}

.notice-popup-close {
	position: absolute;
	width: calc(72 * var(--rpx));
	height: calc(72 * var(--rpx));
	background:url(data:image/png;base64,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) no-repeat 50%/cover;
	opacity: 0.7;
    bottom: calc(-60 * var(--rpx));
		top: auto;
		right: 50%;
		opacity: 1;
		margin-right: calc(-36 * var(--rpx));
}

.notice-popup-close:active {
	opacity: 1;
}

.notice-popup-username {
	max-width: 100%;
	line-height: calc(74 * var(--rpx));
	font-size: calc(calc(52 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: #404040;
	display: flex;
	align-items: center;
	justify-content: center;
	flex-wrap: wrap;
	overflow: visible;
	word-break: keep-all;
	white-space: normal;
}

.notice-popup-company {
	display: flex;
	align-items: center;
	min-width: 0;
	max-width: 100%;
	margin-top: calc(12 * var(--rpx));
}

.notice-popup-company-name {
	max-width: 100%;
	line-height: calc(40 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 400;
	color: #404040;
	overflow: hidden;
	word-break: keep-all;
	white-space: nowrap;
	text-overflow: ellipsis;
}

.notice-popup-text {
	position: relative;
	margin-top: calc(36 * var(--rpx));
	line-height: calc(48 * var(--rpx));
	width: 100%;
	max-width: 100%;
	box-sizing: border-box;
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 400;
	color: #404040;
	max-height: 40vh;
	overflow-x: hidden;
	overflow-y: auto;
}

.notice-popup-text::before {
	content: '';
	position: absolute;
	left: 0;
	top: calc(-32 * var(--rpx));
	width: calc(60 * var(--rpx));
	height: calc(8 * var(--rpx));
	background: #4285f4;
}

.notice-popup-button {
	width: calc(440 * var(--rpx));
	height: calc(92 * var(--rpx));
	margin: 0 auto;
	background: linear-gradient(270deg, #5ec7ff, #5d9aff);
	border-radius: calc(20 * var(--rpx));
	font-weight: 500;
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	color: #fff;
	display: flex;
	align-items: center;
	justify-content: center;
	text-align: center;
	box-sizing: border-box;
}

.notice-popup-button:active {
	background: linear-gradient(270deg, #40bcfd, #4288fd);
}

/* 移动端适配 */
@media screen and (max-width: 1023px) {
	.notice-popup-wrapper {
		width: calc(590 * var(--rpx));
	}

	.notice-popup-wrapper::before {
		top: calc(6 * var(--rpx));
		width: calc(200 * var(--rpx));
		height: calc(148 * var(--rpx));
		background-size: calc(200 * var(--rpx)) calc(148 * var(--rpx));
	}

	.notice-popup-wrapper::after {
		right: calc(-32 * var(--rpx));
	}

	.notice-popup-container {
		padding: calc(140 * var(--rpx)) calc(40 * var(--rpx)) calc(60 * var(--rpx));
	}

	.notice-popup-main {
		margin-bottom: calc(54 * var(--rpx));
	}

	.notice-popup-text {
		letter-spacing: calc(-1 * var(--rpx));
	}

	.notice-popup-button {
		width: calc(400 * var(--rpx));
		height: calc(80 * var(--rpx));
	}
}
</style>

<style>
.van-popup:has(.notice-popup-wrapper) {
	overflow: visible;
}
.notice-popup-text img,
.notice-popup-text table {
	max-width: 100%;
}
</style>
