<template>
	<div class="coin-page">
		<mescroll-body
			ref="coinMescroll"
			@init="mescrollInit"
			@down="downCallback"
			@up="upCallback"
			:down="downOption"
			:up="upOption"
		>
			<div class="coin-head">
				<span class="coin-head-label">当前积分</span>
				<span class="coin-head-sum">{{ total }}</span>
				<div class="qq-now" v-if="rank && rank.level">
					<span class="qq-lv">我的等级 Lv.{{ rank.level }}</span>
					<coin-badge size="lg" :coins="total" />
					<span class="qq-lack">再获得 {{ rank.lack }} 积分升级</span>
				</div>
				<div v-else class="qq-lack">再获得 {{ rank ? rank.lack : 25 }} 积分升级</div>
			</div>
			<div class="coin-rules">
				<div class="coin-rules-title">如何获得积分</div>
				<div class="coin-rules-grid">
					<div class="coin-rule" v-for="rule in rules" :key="rule.name">
						<span class="coin-rule-name">{{ rule.name }}</span>
						<span class="coin-rule-sum">+{{ rule.sum }}</span>
					</div>
				</div>
				<div class="coin-rules-note">积分只发放、暂不消耗。有效评论每天最多 3 次。</div>
			</div>
			<div class="coin-list">
				<div class="coin-ledger-head">
					<span>来源</span>
					<span>说明</span>
					<span>时间</span>
					<span>积分</span>
				</div>
				<div
					class="coin-item"
					v-for="(item, index) in coinList"
					:key="item.id || index"
				>
					<div class="coin-item-left">
						<div class="coin-type">{{ coinLabel(item) }}</div>
						<div class="coin-time">{{ timeFrom(item.created_at) }}</div>
						<div class="coin-desc" v-if="coinRemark(item)">{{ coinRemark(item) }}</div>
					</div>
					<div class="coin-item-right">
						<span class="coin-amount" :class="{ 'add': item.sum > 0, 'sub': item.sum < 0 }">
							{{ item.sum > 0 ? '+' : '' }}{{ item.sum }}
						</span>
					</div>
				</div>
			</div>
		</mescroll-body>
	</div>
</template>

<script>
import { getStorage, setStorage } from '@/common/storage'
import MescrollMixin from "@/uni_modules/mescroll-uni/components/mescroll-uni/mescroll-mixins.js";
import { qqRank } from '@/common/coinLevel'

export default {
	mixins: [MescrollMixin],
	data() {
		return {
			coinList: [],
			total: 0,
			rules: [
				{ name: '每日首次进入', sum: 5 },
				{ name: '看完视频/音频', sum: 10 },
				{ name: '读完图文', sum: 5 },
				{ name: '有效评论', sum: 5 },
				{ name: '通过测评', sum: 10 },
				{ name: '课程结业', sum: 20 }
			],
			downOption: {
				auto: false
			},
			upOption: {
				auto: true,
				page: {
					num: 0,
					size: 10
				},
				noMoreSize: 4,
				empty: {
					tip: '暂无记录',
					icon: '/static/empty/empty.png'
				}
			}
		}
	},
	computed: {
		rank() {
			return qqRank(this.total, this.appInfo && this.appInfo.coinConfig)
		},
	},
	onLoad() {
		this.restoreCoinCache();
		this.loadTotal();
	},
	methods: {
		mescrollInit(mescroll) {
			this.mescroll = mescroll;
		},
		coinRemark(item) {
			return String(item.remark || '').replace(/\s*#[^#]+#/g, '').trim();
		},
		coinLabel(item) {
			const field = item.field || '';
			const map = {
				coins: '积分',
				coin: '积分',
				credit: '学分',
				credit1: '学分',
				credit2: '学分'
			};
			return map[field] || field || '积分变动';
		},
		restoreCoinCache() {
			try {
				const cached = getStorage('qixue_my_coins');
				if (cached && cached.length) {
					this.coinList = cached;
				}
			} catch (e) {}
		},
		loadTotal() {
			try {
				const stats = getStorage('qixue_me_stats');
				if (stats && stats.coins != null) this.total = stats.coins;
			} catch (e) {}
			this.$http.get('/user/myCount').then(res => {
				if (res.data && res.data.coins != null) {
					this.total = res.data.coins
					if (this.userInfo && this.userInfo.id) {
						this.$store.commit('setUserInfo', Object.assign({}, this.userInfo, { coins: this.total }))
					}
				}
			}).catch(() => {});
		},
		downCallback() {
			this.mescroll.resetUpScroll();
		},
		upCallback(page) {
			this.getCoinList(page);
		},
		getCoinList(page) {
			this.$http.get('/user/myCoinsLog', {
				page: page.num,
				limit: page.size
			}).then(res => {
				if (res.code === 200) {
					let list = res.data.list || []

					if (page.num === 1) {
						this.coinList = list;
						try {
							setStorage('qixue_my_coins', list);
						} catch (e) {}
					} else {
						this.coinList = this.coinList.concat(list);
					}
					this.mescroll.endSuccess(list.length, list.length >= page.size);
				} else {
					this.msg(res.msg || '获取积分记录失败');
					this.mescroll.endErr();
				}
			}).catch(err => {
				console.log('获取积分记录失败:', err);
				this.msg('获取积分记录失败');
				this.mescroll.endErr();
			});
		}
	}
}
</script>

<style lang="scss" scoped>
.coin-page {
	min-height: 100vh;
	background: #f6f7fb;
}

.coin-head,
.coin-rules {
	margin: calc(20 * var(--rpx)) calc(30 * var(--rpx)) 0;
	padding: calc(24 * var(--rpx));
	background: #ffffff;
	border-radius: calc(16 * var(--rpx));
}

.coin-head-label {
	display: block;
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
}

.coin-head-sum {
	display: block;
	margin-top: calc(8 * var(--rpx));
	font-size: calc(calc(56 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 700;
	line-height: 1.1;
	color: #4285f4;
}

.qq-now {
	display: flex;
	align-items: center;
	flex-wrap: wrap;
	gap: calc(12 * var(--rpx));
	margin-top: calc(20 * var(--rpx));
}

.qq-lv {
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 700;
	color: #e0a000;
}

.qq-lack {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
}

.coin-head > .qq-lack {
	display: block;
	margin-top: calc(12 * var(--rpx));
}

.coin-rules-title {
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
	margin-bottom: calc(16 * var(--rpx));
}

.coin-rules-grid {
	display: flex;
	flex-wrap: wrap;
	gap: calc(12 * var(--rpx));
}

.coin-rule {
	display: flex;
	align-items: center;
	justify-content: space-between;
	width: calc(50% - calc(6 * var(--rpx)));
	box-sizing: border-box;
	padding: calc(16 * var(--rpx)) calc(18 * var(--rpx));
	border-radius: calc(12 * var(--rpx));
	background: #f6f8fb;
}

.coin-rule-name {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #404040;
}

.coin-rule-sum {
	flex-shrink: 0;
	margin-left: calc(8 * var(--rpx));
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 650;
	color: #4285f4;
}

.coin-rules-note {
	margin-top: calc(16 * var(--rpx));
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	color: #9aa3af;
	line-height: 1.5;
}

.coin-list {
	padding: calc(20 * var(--rpx)) calc(30 * var(--rpx));
}

.coin-ledger-head {
	display: none;
}

.coin-item {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: calc(32 * var(--rpx)) calc(24 * var(--rpx));
	background: #ffffff;
	border-radius: calc(16 * var(--rpx));
	margin-bottom: calc(20 * var(--rpx));

	&:last-child {
		margin-bottom: 0;
	}
}

.coin-item-left {
	flex: 1;
	display: flex;
	flex-direction: column;
	gap: calc(12 * var(--rpx));
}

.coin-type {
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
	line-height: 1;
}

.coin-time {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	margin: calc(10 * var(--rpx)) 0;
	color: #b2b2b2;
	line-height: 1;
}

.coin-desc {
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
	line-height: 1.4;
	margin-top: calc(4 * var(--rpx));
}

.coin-item-right {
	flex-shrink: 0;
	margin-left: calc(24 * var(--rpx));
}

.coin-amount {
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	line-height: 1;

	&.add {
		color: #4caf50;
	}

	&.sub {
		color: #f44336;
	}
}
</style>
