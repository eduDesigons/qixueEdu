<template>
	<div>
		<div class="set-page">
			<div class="set-content">
				<div class="set-divider"></div>
				<div class="set-box">
					<div class="set-cell set-cell--border" @click="navTo('/pages/member/profile')">
						<div class="set-cell-left">
							<span>个人信息</span>
						</div>
						<div class="set-cell-value"></div>
						<div class="set-cell-right">
							<van-icon name="arrow" size="22" color="#b9bdc2" />
						</div>
					</div>
					<div class="set-cell set-cell--border" @click="navTo('/pages/member/fontSize')">
						<div class="set-cell-left">
							<span>字体大小</span>
						</div>
						<div class="set-cell-value set-cell-value--end">
							<span>{{ fontSizeText }}</span>
						</div>
						<div class="set-cell-right">
							<van-icon name="arrow" size="22" color="#b9bdc2" />
						</div>
					</div>
					<div class="set-cell set-cell--border" @click="handleChangePassword">
						<div class="set-cell-left">
							<span>修改密码</span>
						</div>
						<div class="set-cell-value"></div>
						<div class="set-cell-right">
							<van-icon name="arrow" size="22" color="#b9bdc2" />
						</div>
					</div>
					<div class="set-cell" @click="navTo('/pages/member/about')">
						<div class="set-cell-left">
							<span>关于我们</span>
						</div>
						<div class="set-cell-value"></div>
						<div class="set-cell-right">
							<van-icon name="arrow" size="22" color="#b9bdc2" />
						</div>
					</div>
				</div>
				
			</div>
			<div class="set-logout-btn" @click="handleLogout">
					<span>退出登录</span>
				</div>
		</div>

		<!-- 修改密码弹窗 -->
		<van-popup :show="showPasswordPopup" position="bottom" @update:show="(v)=>{if(!v)handlePasswordCancel()}" :round="20">
			<div class="password-popup">
				<div class="password-popup-title">修改密码</div>
				<div class="password-popup-form">
					<div class="password-form-item">
						<span class="password-form-label">旧密码</span>
						<input 
							class="password-form-input" 
							v-model="passwordForm.oldPassword" 
							placeholder="请输入旧密码" 
							password 
							maxlength="20" 
						/>
					</div>
					<div class="password-form-item">
						<span class="password-form-label">新密码</span>
						<input 
							class="password-form-input" 
							v-model="passwordForm.newPassword" 
							placeholder="请输入新密码" 
							password 
							maxlength="20" 
						/>
					</div>
					<div class="password-form-item">
						<span class="password-form-label">确认密码</span>
						<input 
							class="password-form-input" 
							v-model="passwordForm.confirmPassword" 
							placeholder="请再次输入新密码" 
							password 
							maxlength="20" 
						/>
					</div>
				</div>
				<div class="password-popup-btns">
					<div class="password-popup-btn" @click="handlePasswordCancel">取消</div>
					<div class="password-popup-btn password-popup-btn--primary" @click="handlePasswordConfirm">确定</div>
				</div>
			</div>
		</van-popup>
	</div>
</template>

<script>
import { showModal } from '@/common/dom'
import { fontScaleLabel } from '@/common/fontScale'

export default {
	data() {
		return {
			showPasswordPopup: false,
			passwordForm: {
				oldPassword: '',
				newPassword: '',
				confirmPassword: ''
			},
			fontSizeText: '标准'
		}
	},
	onShow() {
		this.fontSizeText = fontScaleLabel()
	},
	methods: {
		// 修改密码
		handleChangePassword() {
			this.passwordForm = {
				oldPassword: '',
				newPassword: '',
				confirmPassword: ''
			}
			this.showPasswordPopup = true
		},
		// 确认修改密码
		handlePasswordConfirm() {
			if (!this.passwordForm.oldPassword.trim()) {
				this.msg('请输入旧密码')
				return
			}
			if (!this.passwordForm.newPassword.trim()) {
				this.msg('请输入新密码')
				return
			}
			if (this.passwordForm.newPassword.length < 6) {
				this.msg('新密码长度不能少于6位')
				return
			}
			if (this.passwordForm.newPassword !== this.passwordForm.confirmPassword) {
				this.msg('两次输入的密码不一致')
				return
			}
			if (this.passwordForm.oldPassword === this.passwordForm.newPassword) {
				this.msg('新密码不能与旧密码相同')
				return
			}
			
			this.$http.post('/user/setPwd', {
				old_pwd: this.passwordForm.oldPassword,
				new_pwd: this.passwordForm.newPassword
			}).then(res => {
				if (res.code === 200) {
					this.msg('密码修改成功')
					this.handlePasswordCancel()
				} else {
					this.msg(res.msg || '密码修改失败')
				}
			}).catch(err => {
				this.msg((err && err.message) || '密码修改失败')
			})
		},
		// 取消修改密码
		handlePasswordCancel() {
			this.showPasswordPopup = false
			this.passwordForm = {
				oldPassword: '',
				newPassword: '',
				confirmPassword: ''
			}
		},
		// 退出登录
		handleLogout() {
			showModal({
				title: '提示',
				content: '确定要退出登录吗？',
				success: (res) => {
					if (res.confirm) {
						this.$store.commit('logout')
					}
				}
			})
		}
	}
}
</script>

<style lang="scss">
.set-page {
	flex: 1;
	background-color: #f6f7fb;
}

.set-content {
	background-color: #fff;
}

.set-divider {
	width: 100%;
	height: calc(20 * var(--rpx));
	background-color: #f6f7fb;
}

.set-box {
	box-sizing: border-box;
	padding-left: calc(30 * var(--rpx));
}

.set-cell {
	display: flex;
	align-items: center;
	line-height: calc(96 * var(--rpx));
	color: #7f848c;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 400;
}

.set-cell .set-cell-value {
	flex: 1;
}

.set-cell-value {
	color: #353535;
	font-weight: 400;
	line-height: 1.4;
	text-align: right;
	padding-right: calc(8 * var(--rpx));
}

.set-cell-value--end {
	text-align: right;
	padding-right: calc(8 * var(--rpx));
}

.set-cell-left {
	text-align: left;
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	color: #353535;
}

.set-cell-right {
	height: calc(96 * var(--rpx));
	line-height: calc(96 * var(--rpx));
	margin-right: calc(30 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: flex-end;
}

.set-cell--border {
	position: relative;
	border-bottom: calc(1 * var(--rpx)) solid #fcfafa;
}

.set-logout-btn {
	box-sizing: border-box;
	display: flex;
	align-items: center;
	justify-content: center;
	text-align: center;
	line-height: 1;
	height: calc(96 * var(--rpx));
	width: 90%;
	border-radius: calc(16 * var(--rpx));
	padding-top: calc(2 * var(--rpx));
	background-color: #fff;
	background: #f3f8ff;
	color: #4285f4;
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	transition: background .3s ease;
	position: fixed;
	bottom: calc(100 * var(--rpx));
	left: 5%;
	right: 0;
}

.set-logout-btn:active {
	background: #ebf3ff;
}

/* 修改密码弹窗样式 */
.password-popup {
	padding: calc(60 * var(--rpx)) calc(40 * var(--rpx)) calc(40 * var(--rpx));
	background: #fff;
	border-radius: calc(24 * var(--rpx));
	box-sizing: border-box;
}

.password-popup-title {
	font-size: calc(calc(36 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
	margin-bottom: calc(40 * var(--rpx));
	text-align: center;
}

.password-popup-form {
	margin-bottom: calc(40 * var(--rpx));
}

.password-form-item {
	margin-bottom: calc(32 * var(--rpx));
}

.password-form-item:last-child {
	margin-bottom: 0;
}

.password-form-label {
	display: block;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #353535;
	margin-bottom: calc(16 * var(--rpx));
}

.password-form-input {
	width: 100%;
	height: calc(88 * var(--rpx));
	padding: 0 calc(24 * var(--rpx));
	background: #f5f5f5;
	border-radius: calc(12 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #353535;
	box-sizing: border-box;
}

.password-popup-btns {
	display: flex;
	gap: calc(24 * var(--rpx));
}

.password-popup-btn {
	flex: 1;
	height: calc(88 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
	background: #f5f5f5;
	border-radius: calc(12 * var(--rpx));
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
}

.password-popup-btn--primary {
	background: #4285f4;
	color: #fff;
}
</style>