<template>
	<div class="book-item" @click="navTo('/pagesLearn/detail?id='+info.id)">
		<div class="book-item-box">

			<div class="book-item-info">
				<div class="tc txtl2">
					{{info.title || '名字'}}
				</div>
				<div class="dc txtl1">
					{{info.intro || '介绍'}}
				</div>
			</div>
			<img class="book-img" :src="info.poster ? static(info.poster) : '/static/img/video-img-2.png'"  style="object-fit:cover" />
		</div>
		<div class="book-item-bot">
			<div class="cate">
				#{{info.cate || '分类'}}
			</div>
			<div class="num">
				阅读量：{{info.view_num || 0}}次
			</div>
		</div>
	</div>
</template>

<script>
	export default {
		name: "book-item",
		props: {
			info: {
				type: Object,
				default: () => ({})
			}
		},
		data() {
			return {

			};
		}
	}
</script>

<style lang="scss" scoped>
	.book-item {
		border-bottom: calc(1 * var(--rpx)) solid #efefef;
		padding: calc(20 * var(--rpx)) 0;

		.book-item-box {
			display: flex;
			justify-content: space-between;

			.book-item-info {
				flex: 1;
				padding-right: calc(40 * var(--rpx));
				display: flex;
				flex-direction: column;
				justify-content: space-around;

				.tc {
					color: #232323;
					font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
					line-height: calc(40 * var(--rpx));
				}

				.dc {
					padding-top: calc(10 * var(--rpx));
					color: #777577;
					font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
				}
			}

			.book-img {
				width: calc(160 * var(--rpx));
				height: calc(160 * var(--rpx));
				border-radius: calc(20 * var(--rpx));
			}
		}

		.book-item-bot {
			padding-top: calc(20 * var(--rpx));
			display: flex;
			justify-content: space-between;
			font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
			color: #afafaf;
			.cate{
				color: $primary;
			}
		}
	}
</style>