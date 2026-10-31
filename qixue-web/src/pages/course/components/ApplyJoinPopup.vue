<template>
	<van-popup
		:show="show"
		position="bottom"
		round
		:close-on-click-overlay="true"
		:safe-area-inset-bottom="true"
		@update:show="onShow"
	>
		<div class="apply-join-popup">
			<div class="apply-join-popup__head">
				<div class="apply-join-popup__title">申请加入{{ isPath ? '系列课程' : '课程' }}</div>
				<div class="apply-join-popup__name" v-if="title">{{ title }}</div>
				<div class="apply-join-popup__close" @click="close">
					<van-icon name="cross" size="16" color="#7f848c" />
				</div>
			</div>
			<div class="apply-join-popup__body">
				<textarea
					v-model="reason"
					class="apply-join-popup__textarea"
					placeholder="请说明申请加入的原因，例如希望提升相关技能、对该内容感兴趣等"
					maxlength="500"
				></textarea>
				<div class="apply-join-popup__count">{{ reason.length }}/500</div>
				<div class="apply-join-popup__tip">申请理由至少 10 个字，有助于审核通过</div>
			</div>
			<div class="apply-join-popup__foot">
				<div class="comBtn" :class="{ disabled: submitting || !canSubmit }" @click="submit">
					<span>{{ submitting ? '提交中...' : '提交申请' }}</span>
				</div>
			</div>
		</div>
	</van-popup>
</template>

<script>
export default {
	name: 'ApplyJoinPopup',
	props: {
		show: { type: Boolean, default: false },
		courseId: { type: [Number, String], default: null },
		pathId: { type: [Number, String], default: null },
		title: { type: String, default: '' }
	},
	emits: ['update:show', 'submitted'],
	data() {
		return {
			reason: '',
			submitting: false
		}
	},
	computed: {
		isPath() {
			return !!this.pathId
		},
		canSubmit() {
			return this.reason.trim().length >= 10
		}
	},
	watch: {
		show(val) {
			if (val) {
				this.reason = ''
				this.submitting = false
			}
		}
	},
	methods: {
		onShow(val) {
			this.$emit('update:show', val)
		},
		close() {
			this.$emit('update:show', false)
		},
		submit() {
			if (!this.canSubmit) {
				this.msg('申请理由至少需要10个字符')
				return
			}
			if (!this.pathId && !this.courseId) {
				this.msg('参数错误，请重试')
				return
			}
			if (this.submitting) return
			this.submitting = true
			const req = this.pathId
				? this.$http.post('/path/applyPath', {
					id: this.pathId,
					reason: this.reason.trim()
				})
				: this.$http.post('/learn/applyCourse', {
					id: this.courseId,
					reason: this.reason.trim()
				})
			req.then(() => {
				this.msg('申请提交成功，请等待审核')
				this.close()
				this.$emit('submitted')
			}).catch((err) => {
				this.msg((err && err.message) || '提交失败')
			}).finally(() => {
				this.submitting = false
			})
		}
	}
}
</script>

<style lang="scss" scoped>
.apply-join-popup {
	padding: calc(32 * var(--rpx)) calc(32 * var(--rpx)) calc(24 * var(--rpx));
	background: #fff;
}

.apply-join-popup__head {
	position: relative;
	padding-right: calc(48 * var(--rpx));
	margin-bottom: calc(24 * var(--rpx));
}

.apply-join-popup__title {
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: #1a1a1a;
	line-height: 1.4;
}

.apply-join-popup__name {
	margin-top: calc(8 * var(--rpx));
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
	line-height: 1.5;
}

.apply-join-popup__close {
	position: absolute;
	top: 0;
	right: 0;
	width: calc(48 * var(--rpx));
	height: calc(48 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
}

.apply-join-popup__body {
	position: relative;
}

.apply-join-popup__textarea {
	width: 100%;
	min-height: calc(220 * var(--rpx));
	padding: calc(24 * var(--rpx));
	padding-bottom: calc(56 * var(--rpx));
	background: #f8f9fa;
	border: calc(2 * var(--rpx)) solid #e9ecef;
	border-radius: calc(16 * var(--rpx));
	font-size: max(16px, calc(calc(28 * var(--rpx)) * var(--font-scale, 1)));
	line-height: 1.7;
	color: #1a1a1a;
	outline: none;
	resize: none;

	&:focus {
		border-color: #4285f4;
	}

	&::placeholder {
		color: #adb5bd;
	}
}

.apply-join-popup__count {
	position: absolute;
	right: calc(24 * var(--rpx));
	bottom: calc(72 * var(--rpx));
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	color: #adb5bd;
}

.apply-join-popup__tip {
	margin-top: calc(16 * var(--rpx));
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #8b93a1;
	line-height: 1.5;
}

.apply-join-popup__foot {
	margin-top: calc(28 * var(--rpx));
}
</style>
