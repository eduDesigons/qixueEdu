<template>
	<div v-if="visible || placeholder" class="coin-lv" :class="{ 'is-empty': !visible, 'is-lg': size === 'lg' }">
		<img
			v-for="(mark, index) in marks"
			:key="mark.type + index"
			class="coin-lv__icon"
			:style="{ animationDelay: (index * 0.28) + 's' }"
			:src="static(mark.icon)"
			alt=""
		/>
		<span v-if="rank && rank.crownExtra" class="coin-lv__more">+{{ rank.crownExtra }}</span>
		<span v-if="!visible && placeholder" class="coin-lv__empty">{{ placeholder }}</span>
	</div>
</template>

<script>
import { qqRank } from '@/common/coinLevel'

export default {
	name: 'CoinBadge',
	props: {
		coins: { default: null },
		placeholder: { type: String, default: '' },
		size: { type: String, default: '' },
	},
	computed: {
		rank() {
			if (this.coins === null || this.coins === undefined || this.coins === '') return null
			return qqRank(this.coins, this.appInfo && this.appInfo.coinConfig)
		},
		marks() {
			return (this.rank && this.rank.marks) || []
		},
		visible() {
			return this.marks.length > 0
		},
	},
}
</script>

<style scoped>
.coin-lv {
	display: inline-flex;
	align-items: center;
	flex-wrap: wrap;
	max-width: 148px;
	margin-left: 6px;
	vertical-align: middle;
	flex: none;
	gap: 1px;
}
.coin-lv__icon {
	width: 16px;
	height: 16px;
	object-fit: contain;
	display: block;
	flex: none;
	transform-origin: center;
	animation: coin-lv-twinkle 2.6s ease-in-out infinite;
}
.coin-lv.is-lg {
	max-width: none;
	gap: 2px;
}
.coin-lv.is-lg .coin-lv__icon {
	width: 28px;
	height: 28px;
}
.coin-lv__more {
	margin-left: 2px;
	font-size: 11px;
	line-height: 1;
	font-weight: 700;
	color: #e0a000;
}
.coin-lv.is-empty {
	max-width: none;
	margin-left: 0;
}
.coin-lv__empty {
	font-size: 14px;
	color: #9aa3af;
}

@keyframes coin-lv-twinkle {
	0%, 100% {
		transform: scale(1) rotate(0deg);
		filter: drop-shadow(0 0 0 rgba(255, 196, 64, 0));
	}
	45% {
		transform: scale(1.06) rotate(-3deg);
		filter: brightness(1.08) drop-shadow(0 0 2px rgba(255, 196, 64, 0.45));
	}
	70% {
		transform: scale(1.02) rotate(2deg);
		filter: brightness(1.02) drop-shadow(0 0 1px rgba(255, 196, 64, 0.2));
	}
}

@media (prefers-reduced-motion: reduce) {
	.coin-lv__icon {
		animation: none;
	}
}
</style>
