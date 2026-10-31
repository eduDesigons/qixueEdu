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
			<div class="coin-list">
				<div
					class="coin-item"
					v-for="(item, index) in coinList"
					:key="item.id || index"
				>
					<div class="coin-item-left">
						<div class="coin-type">{{ creditLabel(item) }}</div>
						<div class="coin-time">{{ timeFrom(item.created_at) }}</div>
						<div class="coin-desc">{{ item.remark }}</div>
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

export default {
	mixins: [MescrollMixin],
	data() {
		return {
			coinList: [],
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
	onLoad() {
		this.restoreCreditCache();
	},
	methods: {
		creditLabel(item) {
			const map = { credit: '学分', credit1: '学分', credit2: '学分', coins: '积分' };
			return map[item.field] || '学分变动';
		},
		restoreCreditCache() {
			try {
				const cached = getStorage('qixue_my_credit');
				if (cached && cached.length) {
					this.coinList = cached;
				}
			} catch (e) {}
		},
		downCallback() {
			this.mescroll.resetUpScroll();
		},
		upCallback(page) {
			this.getCoinList(page);
		},
		getCoinList(page) {
			this.$http.get('/user/myCreditLog', {
				page: page.num,
				limit: page.size
			}).then(res => {
				if (res.code === 200) {
					let list = res.data.list || []

					if (page.num === 1) {
						this.coinList = list;
						try {
							setStorage('qixue_my_credit', list);
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

.coin-list {
	padding: calc(20 * var(--rpx)) calc(30 * var(--rpx));
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
	margin: calc(10 * var(--rpx)) 0;
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
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
