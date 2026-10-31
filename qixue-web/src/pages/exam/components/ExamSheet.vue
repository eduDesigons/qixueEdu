<template>
	<div>
		<div class="sheet-title" v-if="title">{{ title }}</div>
		<div class="sheet-grid">
			<div
				class="sheet-item"
				v-for="(q, i) in questions"
				:key="(keyPrefix || '') + (q.id || i)"
				:class="itemClass(i)"
				@click="$emit('pick', i)"
			>{{ i + 1 }}</div>
		</div>
		<div class="sheet-legend" v-if="legend">
			<span><span class="dot current"></span>当前</span>
			<span><span class="dot done"></span>已答</span>
			<span><span class="dot todo"></span>未答</span>
		</div>
	</div>
</template>

<script>
export default {
	name: 'ExamSheet',
	props: {
		title: { type: String, default: '' },
		questions: { type: Array, default: () => [] },
		itemClass: { type: Function, default: () => ({}) },
		legend: { type: Boolean, default: true },
		keyPrefix: { type: String, default: '' },
	},
	emits: ['pick'],
}
</script>

<style lang="scss" scoped>
$primary-color: #4286F5;
$text-primary: #1a1a1a;
$text-secondary: #666;

.sheet-title {
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: $text-primary;
	text-align: center;
	margin-bottom: calc(28 * var(--rpx));
}

.sheet-grid {
	display: flex;
	flex-wrap: wrap;
	gap: calc(16 * var(--rpx));
	margin-bottom: calc(24 * var(--rpx));
}

.sheet-item {
	width: calc(72 * var(--rpx));
	height: calc(72 * var(--rpx));
	border-radius: calc(16 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	background: #f5f7fa;
	color: #7f848c;
	border: calc(2 * var(--rpx)) solid transparent;
}

.sheet-item.done {
	background: #e8f1ff;
	color: $primary-color;
}

.sheet-item.current {
	background: $primary-color;
	color: #fff;
}

.sheet-item.todo {
	background: #fff5f7;
	color: #ff3359;
}

.sheet-legend {
	display: flex;
	justify-content: center;
	gap: calc(32 * var(--rpx));
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	color: $text-secondary;
	margin-bottom: calc(24 * var(--rpx));

	.dot {
		display: inline-block;
		width: calc(16 * var(--rpx));
		height: calc(16 * var(--rpx));
		border-radius: 50%;
		margin-right: calc(8 * var(--rpx));
		vertical-align: middle;
	}

	.dot.current { background: $primary-color; }
	.dot.done { background: #8bbeff; }
	.dot.todo { background: #ff3359; }
}
</style>
