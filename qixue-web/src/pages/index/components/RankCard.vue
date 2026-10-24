<template>
	<div class="rank-card">
		<div class="rank-card-head">
			<div class="rank-card-title">
				<span>{{ title }}</span>
				<div class="rank-card-help" @click.stop="$emit('tip')">
					<van-icon name="question-o" size="16" color="#c0c4cc" />
				</div>
			</div>
			<div class="rank-card-switch">
				<div
					class="rank-card-switch-item"
					:class="{ active: period === 'week' }"
					@click="$emit('update:period', 'week')"
				>本周</div>
				<div
					class="rank-card-switch-item"
					:class="{ active: period === 'month' }"
					@click="$emit('update:period', 'month')"
				>本月</div>
			</div>
		</div>

		<div v-if="loading" class="rank-skel">
			<div class="rank-skel-row" v-for="n in 3" :key="n">
				<div class="rank-skel-no"></div>
				<div class="rank-skel-line"></div>
			</div>
		</div>
		<div v-else-if="!list.length" class="rank-card-empty">
			<span>暂无数据</span>
			<div class="rank-card-history" @click="$emit('history')">
				查看历史榜单
				<van-icon name="arrow" size="12" color="#4286f5" />
			</div>
		</div>
		<div v-else class="rank-list">
			<div
				class="rank-item"
				:class="{ me: item.isMe, course: kind === 'course' }"
				v-for="item in list"
				:key="item.id || item.rank"
				@click="onItem(item)"
			>
				<div class="rank-no" :class="'rank-no--' + item.rank">{{ item.rank }}</div>
				<van-image
					v-if="kind !== 'course'"
					:src="item.avatar ? static(item.avatar) : ''"
					customStyle="margin-right: calc(16 * var(--rpx));"
				 round :width="32 + 'px'" :height="32 + 'px'" fit="cover" />
				<div class="rank-name">
					<span>{{ item.name }}</span>
					<span v-if="item.isMe" class="rank-me">（我）</span>
				</div>
				<coin-badge v-if="kind !== 'course'" :coins="item.coins" />
				<div class="rank-value">
					<span class="rank-value-num">{{ formatValue(item.value) }}</span>
					<span class="rank-value-unit">{{ item.unit }}</span>
				</div>
			</div>
			<div class="rank-card-history rank-card-history--list" @click="$emit('history')">
				查看历史榜单
				<van-icon name="arrow" size="12" color="#4286f5" />
			</div>
		</div>
	</div>
</template>

<script>
export default {
	name: 'RankCard',
	props: {
		title: {
			type: String,
			default: ''
		},
		kind: {
			type: String,
			default: 'user'
		},
		period: {
			type: String,
			default: 'week'
		},
		list: {
			type: Array,
			default: () => []
		},
		loading: {
			type: Boolean,
			default: false
		}
	},
	methods: {
		formatValue(value) {
			const n = Number(value) || 0
			if (n >= 10000) {
				const w = n / 10000
				return (w >= 10 ? Math.round(w) : Math.round(w * 10) / 10) + 'w'
			}
			return String(n)
		},
		onItem(item) {
			if (this.kind === 'course' && item.id) {
				this.navTo('/pages/course/courseDetail?id=' + item.id)
			}
		}
	}
}
</script>

<style lang="scss" scoped>
.rank-card {
	margin-top: calc(24 * var(--rpx));
	padding: calc(28 * var(--rpx)) calc(24 * var(--rpx)) calc(16 * var(--rpx));
	background: #fff;
	border-radius: calc(24 * var(--rpx));
	box-shadow: 0 calc(4 * var(--rpx)) calc(16 * var(--rpx)) rgba(64, 64, 64, 0.04);
}

.rank-card-head {
	display: flex;
	align-items: center;
	justify-content: space-between;
}

.rank-card-title {
	display: flex;
	align-items: center;
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: #353535;
	line-height: calc(44 * var(--rpx));
}

.rank-card-help {
	margin-left: calc(8 * var(--rpx));
	padding: calc(4 * var(--rpx));
}

.rank-card-switch {
	display: flex;
	padding: calc(4 * var(--rpx));
	background: #f3f4f6;
	border-radius: calc(28 * var(--rpx));
}

.rank-card-switch-item {
	min-width: calc(88 * var(--rpx));
	height: calc(48 * var(--rpx));
	line-height: calc(48 * var(--rpx));
	text-align: center;
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
	border-radius: calc(24 * var(--rpx));
}

.rank-card-switch-item.active {
	background: #4286f5;
	color: #fff;
	font-weight: 500;
}

.rank-card-empty {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding: calc(56 * var(--rpx)) 0 calc(40 * var(--rpx));
	color: #c0c4cc;
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
}

.rank-card-history {
	display: flex;
	align-items: center;
	margin-top: calc(20 * var(--rpx));
	color: #4286f5;
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
}

.rank-card-history--list {
	justify-content: center;
	padding: calc(8 * var(--rpx)) 0 calc(12 * var(--rpx));
}

.rank-list {
	padding-top: calc(12 * var(--rpx));
}

.rank-item {
	display: flex;
	align-items: center;
	min-height: calc(88 * var(--rpx));
	padding: calc(8 * var(--rpx)) calc(12 * var(--rpx));
	border-radius: calc(12 * var(--rpx));
}

.rank-item.me {
	background: #f0f6ff;
}

.rank-no {
	width: calc(48 * var(--rpx));
	height: calc(48 * var(--rpx));
	margin-right: calc(12 * var(--rpx));
	line-height: calc(48 * var(--rpx));
	text-align: center;
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: #b9bdc2;
	flex-shrink: 0;
}

.rank-no--1,
.rank-no--2,
.rank-no--3 {
	color: #fff;
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
}

.rank-no--1 {
	background: #f5c452;
}

.rank-no--2 {
	background: #b7c4d4;
}

.rank-no--3 {
	background: #e0a070;
}

.rank-name {
	flex: 1;
	min-width: 0;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #353535;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.rank-item.me .rank-name {
	color: #4286f5;
}

.rank-me {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
}

.rank-value {
	flex-shrink: 0;
	margin-left: calc(16 * var(--rpx));
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
}

.rank-value-num {
	margin-right: calc(6 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: #404040;
}

.rank-skel {
	padding: calc(16 * var(--rpx)) 0 calc(8 * var(--rpx));
}

.rank-skel-row {
	display: flex;
	align-items: center;
	min-height: calc(80 * var(--rpx));
}

.rank-skel-no {
	width: calc(48 * var(--rpx));
	height: calc(48 * var(--rpx));
	margin-right: calc(16 * var(--rpx));
	border-radius: 50%;
	background: #eef0f4;
}

.rank-skel-line {
	flex: 1;
	height: calc(28 * var(--rpx));
	border-radius: calc(8 * var(--rpx));
	background: #eef0f4;
}
</style>
