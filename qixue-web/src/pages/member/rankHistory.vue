<template>
	<div class="rank-history">
		<div class="rank-history-tabs">
			<div
				class="rank-history-tab"
				:class="{ active: period === 'week' }"
				@click="switchPeriod('week')"
			>周榜</div>
			<div
				class="rank-history-tab"
				:class="{ active: period === 'month' }"
				@click="switchPeriod('month')"
			>月榜</div>
		</div>
		<div v-if="loading" class="rank-history-empty">加载中...</div>
		<div v-else-if="!list.length" class="rank-history-empty">暂无历史榜单</div>
		<div v-else>
			<div class="rank-history-block" v-for="item in list" :key="item.periodKey">
				<div class="rank-history-label">{{ item.periodLabel }}</div>
				<div class="rank-history-card">
					<div
						class="rank-item"
						:class="{ me: row.isMe }"
						v-for="row in item.list"
						:key="row.id || row.rank"
						@click="onItem(row)"
					>
						<div class="rank-no" :class="'rank-no--' + row.rank">{{ row.rank }}</div>
						<van-image
							v-if="type !== 'course'"
							:src="row.avatar ? static(row.avatar) : ''"
							customStyle="margin-right: calc(16 * var(--rpx));"
						 round :width="32 + 'px'" :height="32 + 'px'" fit="cover" />
						<div class="rank-name">
							<span>{{ row.name }}</span>
							<span v-if="row.isMe" class="rank-me">（我）</span>
						</div>
						<coin-badge v-if="type !== 'course'" :coins="row.coins" />
						<div class="rank-value">
							<span class="rank-value-num">{{ formatValue(row.value) }}</span>
							<span class="rank-value-unit">{{ row.unit }}</span>
						</div>
					</div>
					<div v-if="!item.list.length" class="rank-history-empty rank-history-empty--inner">暂无数据</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import { setPageTitle } from '@/common/nav'
const titles = {
	duration: '学习时长榜',
	course: '课程学习榜',
	coins: '积分榜'
}

export default {
	data() {
		return {
			type: 'duration',
			period: 'week',
			loading: false,
			list: []
		}
	},
	onLoad(query) {
		this.type = query.type || 'duration'
		setPageTitle(titles[this.type] || '历史榜单')
		this.loadHistory()
	},
	methods: {
		switchPeriod(period) {
			if (this.period === period) return
			this.period = period
			this.loadHistory()
		},
		formatValue(value) {
			const n = Number(value) || 0
			if (n >= 10000) {
				const w = n / 10000
				return (w >= 10 ? Math.round(w) : Math.round(w * 10) / 10) + 'w'
			}
			return String(n)
		},
		onItem(row) {
			if (this.type === 'course' && row.id) {
				this.navTo('/pages/course/courseDetail?id=' + row.id)
			}
		},
		loadHistory() {
			this.loading = true
			this.$http.get('/rank/history', {
				type: this.type,
				period: this.period
			}).then(res => {
				this.list = (res.data && res.data.list) || []
			}).catch(() => {
				this.list = []
			}).finally(() => {
				this.loading = false
			})
		}
	}
}
</script>

<style lang="scss" scoped>
.rank-history {
	min-height: 100vh;
	padding: calc(24 * var(--rpx)) calc(30 * var(--rpx)) calc(40 * var(--rpx));
	background: #f6f7fb;
}

.rank-history-tabs {
	display: flex;
	padding: calc(6 * var(--rpx));
	margin-bottom: calc(24 * var(--rpx));
	background: #fff;
	border-radius: calc(40 * var(--rpx));
}

.rank-history-tab {
	flex: 1;
	height: calc(64 * var(--rpx));
	line-height: calc(64 * var(--rpx));
	text-align: center;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
	border-radius: calc(32 * var(--rpx));
}

.rank-history-tab.active {
	background: #4286f5;
	color: #fff;
	font-weight: 500;
}

.rank-history-empty {
	padding: calc(80 * var(--rpx)) 0;
	text-align: center;
	color: #c0c4cc;
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
}

.rank-history-empty--inner {
	padding: calc(40 * var(--rpx)) 0;
}

.rank-history-label {
	margin: calc(8 * var(--rpx)) calc(8 * var(--rpx)) calc(16 * var(--rpx));
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
}

.rank-history-card {
	margin-bottom: calc(24 * var(--rpx));
	padding: calc(12 * var(--rpx)) calc(16 * var(--rpx));
	background: #fff;
	border-radius: calc(20 * var(--rpx));
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
</style>
