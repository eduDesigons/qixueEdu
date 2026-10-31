<template>
	<div class="login-page">
		<login-aurora preset="hero" :colors="auroraColors" :params="auroraParams" />
		<div class="content">
		<div class="brand">{{ siteName }}</div>
		<div class="title">
			密码登录
		</div>
		<div class="password-login-form">
			<div class="item">
				<div class="input">
					<input class="input-long-text" placeholder="请输入帐号"
						v-model="passwordForm.username" />
				</div>
			</div>
			<div class="item">
				<div class="input">
					<input type="password" class="input-long-text" placeholder="请输入密码"
						v-model="passwordForm.password" />
				</div>
			</div>
		</div>

		<div class="comBtn" :class="{ dis: !(passwordForm.username && passwordForm.password) }"
			@click="passwordLoginHandler">
			立即登录
		</div>
		<div class="forgot-password" @click="handleForgotPassword">
			忘记密码
		</div>
		<slide-captcha ref="captcha" @success="onCaptchaSuccess" />
		</div>
	</div>
</template>

<script>
import { hideLoading, showLoading, showToast } from '@/common/toast'
import { reLaunch } from '@/common/nav'
	import SlideCaptcha from '@/components/SlideCaptcha.vue'
	import LoginAurora from '@/components/LoginAurora.vue'

	export default {
		components: { SlideCaptcha, LoginAurora },
		data() {
			return {
				auroraColors: ['#a8c8f8', '#e8f2ff', '#ffffff'],
				auroraParams: {
					speed: 72,
					shape: 'checks',
					shapeScale: 12,
					rotation: -3,
					proportion: 56,
					softness: 90,
					distortion: 24,
					swirl: 20,
					swirlIterations: 6,
					scale: 0.38,
					offsetY: 20,
				},
				passwordForm: {
					username: "",
					password: "",
				},
				needCaptcha: false,
			};
		},
		watch: {
			'passwordForm.username'() {
				this.needCaptcha = false
			},
		},
		methods: {
			syncAutofill() {
				const root = this.$el
				if (!root || !root.querySelectorAll) return
				const inputs = root.querySelectorAll('input')
				const username = inputs[0] && inputs[0].value
				const password = inputs[1] && inputs[1].value
				if (username) this.passwordForm.username = username
				if (password) this.passwordForm.password = password
			},
			credentials() {
				const username = String(this.passwordForm.username || '').trim()
				const password = String(this.passwordForm.password || '').trim()
				this.passwordForm.username = username
				this.passwordForm.password = password
				return { username, password }
			},
			passwordLoginHandler() {
				this.syncAutofill()
				const { username, password } = this.credentials()
				if (!username || !password) {
					showToast('请填写账号和密码')
					return
				}
				if (this.needCaptcha) {
					this.openCaptcha()
					return
				}
				this.submitLogin()
			},
			openCaptcha() {
				const captcha = this.$refs.captcha
				if (!captcha) {
					showToast('验证码未就绪，请重试')
					return
				}
				captcha.open()
			},
			onCaptchaSuccess(captchaVerification) {
				this.submitLogin(captchaVerification)
			},
			submitLogin(captchaVerification) {
				const { username, password } = this.credentials()
				if (!username || !password) {
					showToast('请填写账号和密码')
					return
				}
				showLoading("登录中...");
				const payload = { username, password }
				if (captchaVerification) payload.captchaVerification = captchaVerification
				this.$http.post('/auth/login', payload)
					.then((res) => {
						this.needCaptcha = false
						this.$store.commit('setToken', res.data.token);
						const sys = this.$store.state.appInfo && this.$store.state.appInfo.sysConfig
						if (!sys || !sys.assetDomain) this.$store.dispatch('getSysInfo')
						return this.$store.dispatch('getUserInfo').catch(() => {})
					})
					.then(() => {
						hideLoading();
						reLaunch('/pages/index/home')
					})
					.catch((err) => {
						hideLoading();
						const need = err && (err.code === 10003 || (err.data && err.data.needCaptcha))
						if (!need) return
						this.needCaptcha = true
						if (err.code === 10003) this.openCaptcha()
					});
			},
			handleForgotPassword() {
				showToast("请联系管理员重置密码");
			},
		}
	}
</script>
<style lang="scss" scoped>
	.login-page {
		position: relative;
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 28px 28px;
		box-sizing: border-box;
	}

	.content {
		position: relative;
		z-index: 1;
		width: min(360px, 100%);
		padding: 28px 20px 22px;
		display: flex;
		flex-direction: column;
		background: #fff;
		border-radius: 16px;
		box-shadow: 0 4px 40px rgba(0, 0, 0, 0.05);
	}

	.brand {
		padding-top: 0;
		font-size: calc(22px * var(--font-scale, 1));
		font-weight: 500;
		color: rgba(0, 0, 0, 0.9);
	}

	.title {
		padding: 4px 0 16px;
		font-size: calc(13px * var(--font-scale, 1));
		font-weight: 400;
		color: rgba(0, 0, 0, 0.45);
	}

	.password-login-form {

		.item {
			width: 100%;
			height: 44px;
			background: rgba(13, 16, 49, 0.04);
			border-radius: 10px;
			box-sizing: border-box;
			padding: 0 14px;
			display: flex;
			margin-bottom: 12px;

			.input {
				flex: 1;
				display: flex;
				flex-direction: row;
				justify-content: space-between;
				box-sizing: border-box;
				align-items: center;

				.input-long-text {
					width: 100%;
					height: 24px;
					box-sizing: border-box;
					outline: none;
					color: #333333;
					font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));

					&::placeholder {
						color: #bbbbbb;
					}
				}

				.input-text {
					max-width: calc(200 * var(--rpx));
					height: calc(36 * var(--rpx));
					float: left;
					box-sizing: border-box;
					outline: none;
					color: #333333;
					font-size: calc(calc(16 * var(--rpx)) * var(--font-scale, 1));

					&::placeholder {
						color: #bbbbbb;
					}
				}
			}
		}
	}

	.forgot-password {
		text-align: right;
		margin-top: 12px;
		font-size: calc(13px * var(--font-scale, 1));
		color: rgba(0, 0, 0, 0.45);
		cursor: pointer;
	}

	.comBtn {
		height: 44px;
		display: flex;
		align-items: center;
		justify-content: center;
		line-height: 1;
		margin-top: 8px;
		border-radius: 10px;
		background: $mainColor;
		&.dis {
			background: #dee2e6;
			color: #adb5bd;
		}
	}

	@media (min-width: 768px) {
		.login-page {
			padding: 24px 16px;
		}
		.content {
			width: 480px;
			padding: 60px;
		}
		.brand {
			font-size: calc(26px * var(--font-scale, 1));
		}
		.title {
			padding: 8px 0 28px;
			font-size: calc(14px * var(--font-scale, 1));
		}
		.password-login-form .item {
			height: 48px;
			margin-bottom: 16px;
			padding: 0 16px;
			border-radius: 12px;
		}
		.forgot-password {
			margin-top: 16px;
			font-size: calc(14px * var(--font-scale, 1));
		}
		.comBtn {
			height: 48px;
			margin-top: 12px;
			border-radius: 12px;
		}
	}

</style>