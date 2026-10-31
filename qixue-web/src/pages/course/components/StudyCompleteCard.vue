<template>
	<div class="study-complete-module" :class="{ 'study-complete-module--pending': pending }">
		<div class="study-complete-content">
			<div class="study-complete-icon">
				<van-icon :name="icon" size="24" :color="pending ? '#ff3359' : '#4285f4'" />
			</div>
			<div class="study-complete-text">
				<div class="study-complete-title">{{ title }}</div>
				<div class="study-complete-desc" v-if="desc">{{ desc }}</div>
			</div>
		</div>
		<div class="study-complete-btn" v-if="action" @click="$emit('action')">
			<span>{{ action }}</span>
			<van-icon name="arrow" size="14" color="#fff" />
		</div>
		<slot />
	</div>
</template>

<script>
export default {
	name: 'StudyCompleteCard',
	props: {
		pending: { type: Boolean, default: false },
		icon: { type: String, default: 'success' },
		title: { type: String, default: '' },
		desc: { type: String, default: '' },
		action: { type: String, default: '' },
	},
	emits: ['action'],
}
</script>

<style lang="scss" scoped>
.study-complete-module {
	margin: calc(40 * var(--rpx)) calc(40 * var(--rpx));
	padding: calc(32 * var(--rpx));
	background: linear-gradient(135deg, #f0f7ff 0%, #e8f4ff 100%);
	border-radius: calc(16 * var(--rpx));
	border: calc(1 * var(--rpx)) solid #d0e7ff;
}

.study-complete-content {
	display: flex;
	align-items: flex-start;
	margin-bottom: calc(24 * var(--rpx));
}

.study-complete-icon {
	margin-right: calc(20 * var(--rpx));
	flex-shrink: 0;
}

.study-complete-text {
	flex: 1;
}

.study-complete-title {
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #404040;
	margin-bottom: calc(8 * var(--rpx));
}

.study-complete-desc {
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
	line-height: 1.5;
}

.study-complete-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	padding: calc(20 * var(--rpx)) 0;
	background: linear-gradient(135deg, #5c86ff, #8bbeff);
	border-radius: calc(84 * var(--rpx));
	color: #fff;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	line-height: 1;
}

.study-complete-btn span {
	margin-right: calc(8 * var(--rpx));
	line-height: 1;
}

.study-complete-module--pending {
	background: linear-gradient(135deg, #fff5f7 0%, #ffe8ee 100%);
	border-color: #ffc9d6;
}

.study-complete-module--pending .study-complete-btn {
	background: linear-gradient(135deg, #ff3359, #ff6b82);
}
</style>
