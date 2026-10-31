<template>
	<div>
		<div class="profile-info-wrap">
			<div class="user-avatar">
					<van-image v-if="userInfo.avatar" :src="static(userInfo.avatar)" @click="previewPics(userInfo.avatar)" round :width="100 + 'px'" :height="100 + 'px'" fit="cover" />
					<span class="qx-avatar" v-else @click="setAvatar()" :style="{width:100+'px',height:100+'px',fontSize:Math.round(100/2.5)+'px'}">{{ (userInfo.nick_name || '学' || '').slice(0,1) }}</span>
					<img class="user-avatar-change" @click="setAvatar()" src="/static/img/me/profile_photo_change_icon.png" alt="" />
				</div>
				<clipAvatar ref="clipAvatar" @cropImg="cropImg"></clipAvatar>
			<div class="profile-fields">
			<div class="profile-section-title profile-section-title--base">基本资料</div>
			<div class="profile-section-desc">修改账号请联系管理员</div>
			<div class="profile-cell profile-cell--border profile-cell--static">
				<div class="profile-cell-left" >
					<span>账号</span>
				</div>
				<div class="profile-cell-value" >
					<span>{{ userInfo.account }}</span>
				</div>
				<div class="profile-cell-right"></div>
			</div>
			<div class="profile-cell profile-cell--border" @click="handleEditField('nick_name')">
				<div class="profile-cell-left" >
					<span>姓名</span>
				</div>
				<div class="profile-cell-value" >
					<span>{{ userInfo.nick_name || '未设置' }}</span>
				</div>
				<div class="profile-cell-right">
					<van-icon name="arrow" size="18" color="#b9bdc2" />
				</div>
			</div>
			<div class="profile-cell profile-cell--border profile-cell--static">
				<div class="profile-cell-left">
					<span>积分等级</span>
				</div>
				<div class="profile-cell-value">
					<coin-badge :coins="userInfo.coins" placeholder="暂无等级" />
				</div>
				<div class="profile-cell-right"></div>
			</div>
			<div class="profile-cell" @click="handleEditField('mobile')">
				<div class="profile-cell-left">
					<span>联系方式</span>
				</div>
				<div class="profile-cell-value" >
					<span>{{ userInfo.mobile || '未设置' }}</span>
				</div>
				<div class="profile-cell-right">
					<van-icon name="arrow" size="18" color="#b9bdc2" />
				</div>
			</div>
			<div class="profile-section-title">其他信息</div>
			<div class="profile-cell profile-cell--border" @click="handleEditField('gender')">
				<div class="profile-cell-left" >
					<span>性别</span>
				</div>
				<div class="profile-cell-value">
					<span>{{ getGenderText(userInfo.gender) }}</span>
				</div>
				<div class="profile-cell-right">
					<van-icon name="arrow" size="18" color="#b9bdc2" />
				</div>
			</div>
			<div class="profile-cell profile-cell--border" @click="handleEditField('age')">
				<div class="profile-cell-left" >
					<span>年龄</span>
				</div>
				<div class="profile-cell-value">
					<span>{{ userInfo.age || '未设置' }}</span>
				</div>
				<div class="profile-cell-right">
					<van-icon name="arrow" size="18" color="#b9bdc2" />
				</div>
			</div>
			<div class="profile-cell profile-cell--border" @click="handleEditField('degree')">
				<div class="profile-cell-left" >
					<span>学历</span>
				</div>
				<div class="profile-cell-value">
					<span>{{ getDegreeText(userInfo.degree) }}</span>
				</div>
				<div class="profile-cell-right">
					<van-icon name="arrow" size="18" color="#b9bdc2" />
				</div>
			</div>
			<div class="profile-cell profile-cell--border" @click="handleEditField('profession')">
				<div class="profile-cell-left" >
					<span>专业</span>
				</div>
				<div class="profile-cell-value">
					<span>{{ userInfo.profession || '未设置' }}</span>
				</div>
				<div class="profile-cell-right">
					<van-icon name="arrow" size="18" color="#b9bdc2" />
				</div>
			</div>
			<div class="profile-cell profile-cell--border" @click="handleEditField('job')">
				<div class="profile-cell-left">
					<span>职业</span>
				</div>
				<div class="profile-cell-value">
					<span>{{ userInfo.job || '未设置' }}</span>
				</div>
				<div class="profile-cell-right">
					<van-icon name="arrow" size="18" color="#b9bdc2" />
				</div>
			</div>
			<div class="profile-cell profile-cell--border" @click="handleEditField('marital')">
				<div class="profile-cell-left" >
					<span>婚姻状况</span>
				</div>
				<div class="profile-cell-value">
					<span>{{ getMaritalText(userInfo.marital) }}</span>
				</div>
				<div class="profile-cell-right">
					<van-icon name="arrow" size="18" color="#b9bdc2" />
				</div>
			</div>
			<div class="profile-cell" @click="handleEditField('target')">
				<div class="profile-cell-left" >
					<span>学习目标</span>
				</div>
				<div class="profile-cell-value" >
					<span>{{ userInfo.target || '未设置' }}</span>
				</div>
				<div class="profile-cell-right">
					<van-icon name="arrow" size="18" color="#b9bdc2" />
				</div>
			</div>
			</div>
		</div>

		<van-popup :show="showPicker" position="bottom" round @update:show="(v)=>{if(!v)handlePickerCancel()}">
			<van-picker
				:columns="pickerColumns"
				@confirm="handlePickerConfirm"
				@cancel="handlePickerCancel"
			/>
		</van-popup>

		<van-popup :show="showInputPopup && currentField === 'profession'" position="bottom" @update:show="(v)=>{if(!v)handleInputCancel()}">
			<div class="input-popup">
				<div class="input-popup-title">编辑专业</div>
				<input class="input-popup-input" v-model="inputValue" placeholder="请输入专业" maxlength="20" />
				<div class="input-popup-btns">
					<div class="input-popup-btn" @click="handleInputCancel">取消</div>
					<div class="input-popup-btn input-popup-btn--primary" @click="handleInputConfirm">确定</div>
				</div>
			</div>
		</van-popup>

		<van-popup :show="showInputPopup && currentField === 'job'" position="bottom" @update:show="(v)=>{if(!v)handleInputCancel()}">
			<div class="input-popup">
				<div class="input-popup-title">编辑职业</div>
				<input class="input-popup-input" v-model="inputValue" placeholder="请输入职业" maxlength="20" />
				<div class="input-popup-btns">
					<div class="input-popup-btn" @click="handleInputCancel">取消</div>
					<div class="input-popup-btn input-popup-btn--primary" @click="handleInputConfirm">确定</div>
				</div>
			</div>
		</van-popup>

		<van-popup :show="showInputPopup && currentField === 'nick_name'" position="bottom" @update:show="(v)=>{if(!v)handleInputCancel()}">
			<div class="input-popup">
				<div class="input-popup-title">编辑姓名</div>
				<input class="input-popup-input" v-model="inputValue" placeholder="请输入姓名" maxlength="20" />
				<div class="input-popup-btns">
					<div class="input-popup-btn" @click="handleInputCancel">取消</div>
					<div class="input-popup-btn input-popup-btn--primary" @click="handleInputConfirm">确定</div>
				</div>
			</div>
		</van-popup>

		<van-popup :show="showInputPopup && currentField === 'mobile'" position="bottom" @update:show="(v)=>{if(!v)handleInputCancel()}">
			<div class="input-popup">
				<div class="input-popup-title">编辑联系方式</div>
				<input class="input-popup-input" v-model="inputValue" placeholder="请输入联系方式" maxlength="20" type="number" />
				<div class="input-popup-btns">
					<div class="input-popup-btn" @click="handleInputCancel">取消</div>
					<div class="input-popup-btn input-popup-btn--primary" @click="handleInputConfirm">确定</div>
				</div>
			</div>
		</van-popup>

		<van-popup :show="showInputPopup && currentField === 'target'" position="bottom" @update:show="(v)=>{if(!v)handleInputCancel()}">
			<div class="input-popup">
				<div class="input-popup-title">编辑学习目标</div>
				<textarea class="input-popup-textarea" v-model="inputValue" placeholder="请输入学习目标" maxlength="200"></textarea>
				<div class="input-popup-btns">
					<div class="input-popup-btn" @click="handleInputCancel">取消</div>
					<div class="input-popup-btn input-popup-btn--primary" @click="handleInputConfirm">确定</div>
				</div>
			</div>
		</van-popup>
	</div>
</template>

<script>
import { mapState } from 'vuex'
import clipAvatar from '@/components/clipAvatar/clipAvatar.vue'
import { pickImages } from '@/common/media'

export default {
	components: {
		clipAvatar
	},
	data() {
		return {
			showPicker: false,
			showInputPopup: false,
			currentField: '',
			inputValue: '',
			genderOptions: ['男', '女'],
			ageOptions: Array.from({ length: 70 }, (_, i) => (i + 18).toString()),
			degreeOptions: ['小学', '初中', '高中', '中专', '大专', '本科', '硕士', '博士'],
			maritalOptions: ['未婚', '已婚', '离异', '丧偶']
		}
	},
	computed: {
		...mapState(['userInfo']),
		pickerColumns() {
			const map = {
				gender: this.genderOptions,
				age: this.ageOptions,
				degree: this.degreeOptions,
				marital: this.maritalOptions,
			}
			return (map[this.currentField] || []).map((item) => ({ text: String(item), value: item }))
		},
	},
	onLoad() {
		this.getUserInfo()
	},
	methods: {
		getUserInfo() {
			if (this.$store.getters.hasLogin) {
				this.$store.dispatch('getUserInfo')
			}
		},
		setAvatar() {
			pickImages({ count: 1 }).then((list) => {
				const item = list[0]
				if (!item) return
				this.$refs.clipAvatar.open(item.url)
			}).catch(() => {})
		},
		cropImg(file) {
			this.$http.upload(file).then((res) => {
				const url = res && res.data && (res.data.url || res.data)
				if (!url) throw new Error('empty url')
				return this.$http.post('/user/edit', { avatar: url })
			}).then(() => {
				this.msg('头像更新成功')
				this.$store.dispatch('getUserInfo')
			}).catch(() => {
				this.msg('头像更新失败')
			})
		},
		handleEditField(field) {
			this.currentField = field
			if (['gender', 'age', 'degree', 'marital'].includes(field)) {
				this.showPicker = true
			} else if (['profession', 'job', 'target', 'nick_name', 'mobile'].includes(field)) {
				this.inputValue = this.userInfo[field] || ''
				this.showInputPopup = true
			}
		},
		handlePickerConfirm({ selectedOptions }) {
			const option = selectedOptions && selectedOptions[0]
			const value = option ? (option.value !== undefined ? option.value : option.text) : ''
			this.updateUserInfo({ [this.currentField]: value })
			this.showPicker = false
		},
		handlePickerCancel() {
			this.showPicker = false
			this.currentField = ''
		},
		handleInputConfirm() {
			if (!this.inputValue.trim()) {
				this.msg('请输入内容')
				return
			}
			if (this.currentField === 'mobile') {
				const mobile = this.inputValue.trim()
				if (!/^1[3-9]\d{9}$/.test(mobile)) {
					this.msg('请输入正确的手机号码')
					return
				}
			}
			this.updateUserInfo({ [this.currentField]: this.inputValue.trim() })
			this.showInputPopup = false
			this.inputValue = ''
		},
		handleInputCancel() {
			this.showInputPopup = false
			this.currentField = ''
			this.inputValue = ''
		},
		updateUserInfo(data) {
			this.$http.post('/user/edit', data).then(res => {
				if (res.code === 200) {
					this.msg('更新成功')
					this.$store.dispatch('getUserInfo')
				} else {
					this.msg(res.msg || '更新失败')
				}
			}).catch(err => {
				console.log('更新失败:', err)
				this.msg('更新失败')
			})
		},
		getGenderText(gender) {
			if (gender === 1 || gender === '1' || gender === '男') return '男'
			if (gender === 2 || gender === '2' || gender === '女') return '女'
			return '未设置'
		},
		getDegreeText(degree) {
			if (degree === 0 || degree === '0' || !degree) return '未设置'
			if (this.degreeOptions.indexOf(degree) !== -1) return degree
			const index = parseInt(degree, 10) - 1
			return this.degreeOptions[index] || '未设置'
		},
		getMaritalText(marital) {
			if (marital === 0 || marital === '0' || !marital) return '未设置'
			if (this.maritalOptions.indexOf(marital) !== -1) return marital
			const index = parseInt(marital, 10) - 1
			return this.maritalOptions[index] || '未设置'
		}
	}
}
</script>

<style lang="scss">
.profile-info-wrap .user-avatar {
width: calc(200 * var(--rpx));
height: calc(200 * var(--rpx));
border-radius: 50%;
position: relative;

.user-avatar-img {
	width: calc(200 * var(--rpx));
height: calc(200 * var(--rpx));
	border-radius: 50%;
	border: calc(2 * var(--rpx)) solid #b2b2b2;
}

.user-avatar-change {
	width: calc(50 * var(--rpx));
	position: absolute;
	bottom: calc(0 * var(--rpx));
	right: calc(0 * var(--rpx));
}
}

.profile-form-input {
	box-sizing: border-box;
	background-color: transparent;
	border: 0;
	outline: none;
	-webkit-appearance: textfield;
	appearance: textfield;
}

.profile-form-input[type=hidden] {
	flex: 0;
	width: 0;
	height: 0;
}

.profile-form-element {
	box-sizing: border-box;
	position: relative;
	display: flex;
	transition: all .3s ease;
}

.profile-select-label {
	flex: 1;
	display: flex;
	align-items: center;
	line-height: 1;
	-webkit-user-select: none;
	-moz-user-select: none;
	-ms-user-select: none;
	user-select: none;
}

.profile-select-arrow {
	display: block;
	position: absolute;
	right: 0;
	top: 0;
	transform-origin: center;
	transition: all .3s ease;
}

.profile-cell--border {
	position: relative;
	border-bottom: calc(1 * var(--rpx)) solid #fcfafa;
}

.profile-cell .profile-cell-value {
	flex: 1;
}

.profile-fields {
	display: contents;
}

.profile-cell {
	display: flex;
	align-items: center;
	padding: calc(10 * var(--rpx)) 0;
	line-height: 22px;
	color: #7f848c;
	font-size: calc(14px * var(--font-scale, 1));
	font-weight: 400;
}

.profile-cell-value {
	color: #353535;
}

.profile-cell-value {
	font-weight: 400;
	line-height: 20px;
}

.profile-cell-left {
	text-align: left;
	width: calc(200 * var(--rpx));
}

.profile-cell-icon {
	color: #b9bdc2;
	font-size: calc(16px * var(--font-scale, 1));
	font-weight: 400;
}

.profile-cell-right {
	height: 22px;
	line-height: 22px;
	display: flex;
	align-items: center;
}

.profile-popup {
	z-index: 2024;
}

.profile-info-wrap {
	box-sizing: border-box;
	padding: 30px 0 100px 16px;
}

.profile-info-wrap [can-edit] {
	padding-right: 15px;
}

.profile-info-wrap [can-edit] .profile-cell-value {
	color: #7f848c;
	font-size: calc(14px * var(--font-scale, 1));
	font-weight: 500;
}

.profile-info-wrap [can-edit] .profile-cell-value {
	color: #353535;
	font-size: calc(14px * var(--font-scale, 1));
	font-weight: 500;
}

.profile-section-title {
	margin: 30px 0 10px;
	line-height: 26px;
	font-size: calc(19px * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
}

.profile-section-title--base {
	margin: 32px 0 0;
}

.profile-cell--enterprise {
	align-items: flex-start;
	padding: 13px 0;
	line-height: 22px;
}

.profile-cell--enterprise .profile-cell-left {
	line-height: 22px;
}

.profile-section-desc {
	margin: 2px 0 10px;
	color: #b9bdc2;
	line-height: 19px;
	font-size: calc(13px * var(--font-scale, 1));
	font-weight: 400;
}

.profile-form-element--hidden {
	display: none;
}

.input-popup {
	padding: calc(40 * var(--rpx)) calc(32 * var(--rpx));
	background: #fff;
	border-radius: calc(24 * var(--rpx)) calc(24 * var(--rpx)) 0 0;
}

.input-popup-title {
	font-size: calc(calc(36 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
	margin-bottom: calc(32 * var(--rpx));
	text-align: center;
}

.input-popup-input {
	width: 100%;
	height: calc(88 * var(--rpx));
	padding: 0 calc(24 * var(--rpx));
	background: #f5f5f5;
	border-radius: calc(12 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #353535;
	margin-bottom: calc(32 * var(--rpx));
	box-sizing: border-box;
}

.input-popup-textarea {
	width: 100%;
	min-height: calc(200 * var(--rpx));
	padding: calc(24 * var(--rpx));
	background: #f5f5f5;
	border-radius: calc(12 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #353535;
	margin-bottom: calc(32 * var(--rpx));
	box-sizing: border-box;
}

.input-popup-btns {
	display: flex;
	gap: calc(24 * var(--rpx));
}

.input-popup-btn {
	flex: 1;
	height: calc(88 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
	background: #f5f5f5;
	border-radius: calc(12 * var(--rpx));
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	line-height: 1;
	color: #7f848c;
}

.input-popup-btn--primary {
	background: #4285f4;
	color: #fff;
}
</style>