<template>
	<div class="history-page">
		<div class="history-summary" v-if="total > 0">
			<span class="history-summary-num">{{ total }}</span>
			<span class="history-summary-label">讲已学完</span>
		</div>
		<mescroll-body
			ref="historyMescroll"
			@init="mescrollInit"
			@down="downCallback"
			@up="upCallback"
			:down="downOption"
			:up="upOption"
		>
			<div class="history-wrap">
				<div
					class="history-day"
					v-for="(group, gi) in groups"
					:key="group.date"
				>
					<div class="history-date">{{ group.label }}</div>
					<div class="history-steps">
						<div
							class="history-step"
							v-for="(item, ii) in group.items"
							:key="item.id || ii"
							@click="goItem(item)"
						>
							<span class="history-time">{{ item.time }}</span>
							<div class="history-rail">
								<div class="history-dot" :class="{ latest: gi === 0 && ii === 0 }"></div>
							</div>
							<div class="history-card">
								<span class="history-title">{{ item.title }}</span>
								<span class="history-course">{{ item.courseTitle }}</span>
							</div>
						</div>
					</div>
				</div>
			</div>
		</mescroll-body>
	</div>
</template>

<script>
import { getStorage, setStorage } from '@/common/storage'
import MescrollMixin from "@/uni_modules/mescroll-uni/components/mescroll-uni/mescroll-mixins.js";

export default {
	mixins: [MescrollMixin],
	data() {
		return {
			groups: [],
			total: 0,
			downOption: {
				auto: false
			},
			upOption: {
				auto: true,
				page: {
					num: 0,
					size: 20
				},
				noMoreSize: 4,
				empty: {
					tip: '还没有学完的课节',
					icon: '/static/empty/empty.png'
				}
			}
		}
	},
	onLoad() {
		this.restoreCache()
	},
	methods: {
		restoreCache() {
			try {
				const cached = getStorage('qixue_learn_history')
				if (cached && cached.groups && cached.groups.length) {
					this.groups = cached.groups
					this.total = cached.total || 0
				}
			} catch (e) {}
		},
		downCallback() {
			this.mescroll.resetUpScroll()
		},
		upCallback(page) {
			this.$http.get('/learn/myHistory', {
				page: page.num,
				limit: page.size
			}).then(res => {
				const rows = (res.data && res.data.list) || []
				if (page.num === 1) {
					this.groups = []
					this.total = (res.data && res.data.total) || 0
				}
				this.appendRows(rows)
				if (page.num === 1) {
					try {
						setStorage('qixue_learn_history', {
							groups: this.groups,
							total: this.total
						})
					} catch (e) {}
				}
				this.mescroll.endSuccess(rows.length, rows.length >= page.size)
			}).catch(() => {
				this.msg('获取学习轨迹失败')
				this.mescroll.endErr()
			})
		},
		appendRows(rows) {
			rows.forEach((item) => {
				const date = item.date || ''
				const last = this.groups[this.groups.length - 1]
				if (last && last.date === date) {
					last.items.push(item)
					return
				}
				this.groups.push({
					date: date,
					label: this.dateLabel(date),
					items: [item]
				})
			})
		},
		dateLabel(date) {
			if (!date) return ''
			const today = this.formatYmd(new Date())
			const y = new Date()
			y.setDate(y.getDate() - 1)
			const yesterday = this.formatYmd(y)
			if (date === today) return '今天'
			if (date === yesterday) return '昨天'
			const p = date.split('-')
			if (p.length === 3) {
				const year = Number(p[0])
				const nowYear = new Date().getFullYear()
				if (year === nowYear) {
					return Number(p[1]) + '月' + Number(p[2]) + '日'
				}
				return year + '年' + Number(p[1]) + '月' + Number(p[2]) + '日'
			}
			return date
		},
		formatYmd(d) {
			const m = d.getMonth() + 1
			const day = d.getDate()
			return d.getFullYear() + '-' + (m < 10 ? '0' : '') + m + '-' + (day < 10 ? '0' : '') + day
		},
		goItem(item) {
			if (item.chapterId) {
				const q = item.courseId
					? '/pages/course/chapterDetail?id=' + item.chapterId + '&courseId=' + item.courseId
					: '/pages/course/chapterDetail?id=' + item.chapterId
				this.navTo(q)
				return
			}
			if (item.courseId) {
				this.navTo('/pages/course/courseDetail?id=' + item.courseId)
			}
		}
	}
}
</script>

<style lang="scss" scoped>
.history-page {
	min-height: 100vh;
	background: #fff;
}

.history-summary {
	display: flex;
	align-items: baseline;
	padding: calc(24 * var(--rpx)) calc(32 * var(--rpx)) calc(8 * var(--rpx));
}

.history-summary-num {
	font-size: calc(calc(44 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 700;
	color: #4286f5;
	line-height: 1.1;
}

.history-summary-label {
	margin-left: calc(10 * var(--rpx));
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
}

.history-wrap {
	padding: calc(12 * var(--rpx)) 0 calc(48 * var(--rpx));
}

.history-day {
	padding-right: calc(30 * var(--rpx));
	margin-bottom: calc(48 * var(--rpx));
}

.history-date {
	width: calc(200 * var(--rpx));
	height: calc(52 * var(--rpx));
	margin-bottom: calc(28 * var(--rpx));
	background: #f6f7fb;
	border-radius: 0 calc(100 * var(--rpx)) calc(100 * var(--rpx)) 0;
	line-height: calc(52 * var(--rpx));
	text-align: center;
	color: #404040;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
}

.history-steps {
	padding-left: calc(150 * var(--rpx));
}

.history-step {
	position: relative;
	display: flex;
	align-items: flex-start;
	min-height: calc(88 * var(--rpx));
	padding-bottom: calc(28 * var(--rpx));
}

.history-step:last-child {
	padding-bottom: 0;
}

.history-time {
	position: absolute;
	left: calc(-150 * var(--rpx));
	top: calc(22 * var(--rpx));
	width: calc(100 * var(--rpx));
	text-align: right;
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
	line-height: calc(34 * var(--rpx));
}

.history-rail {
	position: absolute;
	left: calc(-28 * var(--rpx));
	top: 0;
	bottom: 0;
	width: calc(16 * var(--rpx));
}

.history-step:not(:last-child) .history-rail::after {
	content: '';
	position: absolute;
	left: calc(7 * var(--rpx));
	top: calc(36 * var(--rpx));
	bottom: calc(-28 * var(--rpx));
	width: calc(2 * var(--rpx));
	background: #e9eaed;
}

.history-dot {
	position: absolute;
	top: calc(28 * var(--rpx));
	left: 0;
	z-index: 1;
	width: calc(14 * var(--rpx));
	height: calc(14 * var(--rpx));
	background: #fff;
	border: calc(3 * var(--rpx)) solid #d8dbe0;
	border-radius: 50%;
	box-sizing: border-box;
}

.history-dot.latest {
	border-color: #4286f5;
	background: #4286f5;
}

.history-card {
	flex: 1;
	padding: calc(20 * var(--rpx)) calc(24 * var(--rpx));
	background: #f6f7fb;
	border-radius: calc(12 * var(--rpx));
}

.history-title {
	display: block;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #353535;
	line-height: calc(40 * var(--rpx));
}

.history-course {
	display: block;
	margin-top: calc(6 * var(--rpx));
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	color: #b9bdc2;
	line-height: calc(32 * var(--rpx));
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}
</style>
