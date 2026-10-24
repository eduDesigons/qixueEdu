<template>
	<div class="apply-join-page">

		<!-- 页面内容 -->
		<div class="page-content">
			<!-- 课程信息提示 -->
			<div class="course-info-tip" v-if="courseTitle || pathTitle">
				<div class="tip-text">
					<span class="tip-course">申请加入{{ courseTitle ? '课程' : '系列课程' }}</span>
					<span class="tip-title">{{ courseTitle || pathTitle }}</span>

				</div>
			</div>

			<!-- 申请理由输入 -->
			<div class="reason-section">
				<div class="section-header">
					<div class="section-title">申请理由</div>
					<div class="section-required">必填</div>
				</div>
				<div class="textarea-wrapper">
					<textarea v-model="reason" class="reason-textarea" :placeholder="placeholderText" :maxlength="500"
						:auto-height="true" :show-confirm-bar="false" :adjust-position="true"
						@input="handleTextareaInput"></textarea>
					<div class="textarea-count">{{ reason.length }}/500</div>
				</div>
				<div class="section-tip">
					<van-icon name="info-o" size="14" color="#7f848c" />
					<span>请简要说明您申请加入该课程的原因，有助于审核通过</span>
				</div>
			</div>
		</div>

		<!-- 底部提交按钮 -->
		<div class="footer-submit">
			<div class="comBtn" :class="{ 'disabled': submitting || !canSubmit }" @click="handleSubmit">
				<span v-if="!submitting">提交申请</span>
				<span v-else>提交中...</span>
			</div>
		</div>
	</div>
</template>

<script>
import { navigateBack, setPageTitle } from '@/common/nav'
export default {
	data() {
		return {
			courseId: null,
			courseTitle: '',
			pathId: null,
			pathTitle: '',
			reason: '',
			submitting: false,
			pageTitle: '申请加入课程'
		}
	},
	computed: {
		placeholderText() {
			return '请说明您申请加入该课程的原因，例如：希望通过学习提升相关技能、对该课程内容感兴趣等...'
		},
		canSubmit() {
			return this.reason.trim().length >= 10
		}
	},
	onLoad(options) {
		if (options.id) {
			this.courseId = options.id
		}
		if (options.title) {
			this.courseTitle = decodeURIComponent(options.title)
		}
		if (options.pathId) {
			this.pathId = options.pathId
		}
		if (options.pathTitle) {
			this.pathTitle = decodeURIComponent(options.pathTitle)
		}
		if (this.pathId) {
			setPageTitle('申请加入系列课程')
		}
	},
	methods: {
		handleTextareaInput(e) {
			this.reason = e.detail.value
		},
		// 提交申请
		handleSubmit() {
			// 验证输入
			if (!this.canSubmit) {
				this.msg('申请理由至少需要10个字符')
				return
			}

			if (!this.pathId && !this.courseId) {
				this.msg('参数错误，请重试')
				return
			}

			if (this.submitting) {
				return
			}

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
				setTimeout(() => {
					navigateBack()
				}, 1500)
			}).catch(err => {
				this.msg((err && err.message) || '提交失败')
			}).finally(() => {
				this.submitting = false
			})
		}
	}
}
</script>

<style lang="scss" scoped>
.apply-join-page {
	display: flex;
	flex-direction: column;
	min-height: 100vh;
	background: #f8f9fa;
}

.page-content {
	flex: 1;
	padding: calc(32 * var(--rpx));
	padding-bottom: calc(180 * var(--rpx)); // 为底部按钮留出空间
}

.course-info-tip {
	display: flex;
	align-items: flex-start;
	padding: calc(32 * var(--rpx));
	background: linear-gradient(135deg, $mainColor 0%, #d9effd 100%);
	border-radius: calc(20 * var(--rpx));
	margin-bottom: calc(32 * var(--rpx));
	box-shadow: 0 calc(8 * var(--rpx)) calc(24 * var(--rpx)) rgba(102, 126, 234, 0.15);
}

.tip-icon {
	margin-right: calc(20 * var(--rpx));
	flex-shrink: 0;
	width: calc(48 * var(--rpx));
	height: calc(48 * var(--rpx));
	background: rgba(255, 255, 255, 0.2);
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
}

.tip-text {
	flex: 1;
	display: flex;
	flex-direction: column;
}

.tip-title {
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: #fff;
	margin-bottom: calc(12 * var(--rpx));
}

.tip-course {
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: rgba(255, 255, 255, 0.9);
	line-height: 1.6;
}

.reason-section {
	background: #fff;
	border-radius: calc(20 * var(--rpx));
	padding: calc(40 * var(--rpx));
	box-shadow: 0 calc(4 * var(--rpx)) calc(16 * var(--rpx)) rgba(0, 0, 0, 0.04);
}

.section-header {
	display: flex;
	align-items: center;
	margin-bottom: calc(32 * var(--rpx));
}

.section-title {
	font-size: calc(calc(36 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: #1a1a1a;
}

.section-required {
	margin-left: calc(12 * var(--rpx));
	padding: calc(4 * var(--rpx)) calc(12 * var(--rpx));
	background: #fee;
	color: #f56565;
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	border-radius: calc(4 * var(--rpx));
}

.textarea-wrapper {
	position: relative;
	background: #f8f9fa;
	border-radius: calc(16 * var(--rpx));
	padding: calc(24 * var(--rpx));
	margin-bottom: calc(24 * var(--rpx));
	border: calc(2 * var(--rpx)) solid #e9ecef;
	transition: border-color 0.3s;

	&:focus-within {
		border-color: #4285f4;
	}
}

.reason-textarea {
	width: 100%;
	min-height: calc(300 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	line-height: 1.8;
	color: #1a1a1a;
	background: transparent;
	border: none;
	outline: none;

	&::placeholder {
		color: #adb5bd;
	}
}

.textarea-count {
	position: absolute;
	right: calc(24 * var(--rpx));
	bottom: calc(24 * var(--rpx));
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #adb5bd;
}

.section-tip {
	display: flex;
	align-items: flex-start;
	padding-top: calc(24 * var(--rpx));
	border-top: calc(1 * var(--rpx)) solid #e9ecef;
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #6c757d;
	line-height: 1.8;

	span {
		margin-left: calc(8 * var(--rpx));
		flex: 1;
	}
}

.footer-submit {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	padding: calc(24 * var(--rpx)) calc(32 * var(--rpx));
	padding-bottom: calc(calc(24 * var(--rpx)) + env(safe-area-inset-bottom));
	background: #fff;
	border-top: calc(1 * var(--rpx)) solid #e9ecef;
	box-shadow: 0 calc(-4 * var(--rpx)) calc(16 * var(--rpx)) rgba(0, 0, 0, 0.04);
	z-index: 10;
}
</style>