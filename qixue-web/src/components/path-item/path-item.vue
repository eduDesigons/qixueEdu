<template>
	<div class="daily">
		<div class="info">
			<div src="info-tip">
				<div class="daily-title f36">
					{{info.name || '系列名称'}}
				</div>
				<div class="daily-desc f24">
					{{info.intro || ''}}
				</div>
			</div>
			<img class="play-btn" src="/static/icon/line.png"  style="width:100%;height:auto;display:block" />
	
	
		</div>
		<img class="disk" :src="info.avatar ? static(info.avatar) : '/static/img/video-img-3.png'"  style="width:100%;height:auto;display:block" />
		<img class="disk disk-top" :src="info.avatar ? static(info.avatar) : '/static/img/video-img-3.png'"  style="width:100%;height:auto;display:block" />
	
	</div>
</template>

<script>
	export default {
		name: "path-item",
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
	
	.daily {
		position: relative;
		margin: calc(40 * var(--rpx)) calc(10 * var(--rpx));
		
		.info {
			padding: calc(10 * var(--rpx)) calc(30 * var(--rpx)) calc(10 * var(--rpx)) calc(220 * var(--rpx));
			display: flex;
			align-items: center;
			justify-content: space-between;
			// background-image: radial-gradient(circle, #ed717b, #f5939d, #fab3bd, #fcd4db, #fdf4f6);
			background: $blockBg;
			border-radius: calc(30 * var(--rpx));
			height: calc(120 * var(--rpx));
		
			.daily-title {
				color: $title;
			}
		
			.daily-desc {
				color: $desc;
			}
		
			.play-btn {
				width: calc(100 * var(--rpx));
				height: calc(100 * var(--rpx));
			}
		}
		
		.disk {
			position: absolute;
			left: calc(40 * var(--rpx));
			bottom: calc(0 * var(--rpx));
			width: calc(140 * var(--rpx));
			height: calc(140 * var(--rpx));
			border-radius: calc(30 * var(--rpx));
			border: calc(4 * var(--rpx)) solid #f2f2f2;
			box-shadow: calc(0 * var(--rpx)) calc(0 * var(--rpx)) calc(4 * var(--rpx)) calc(4 * var(--rpx)) #f2f2f2;
		
		}
		
		.disk-top {
			left: 0;
			bottom: 0;
			width: calc(160 * var(--rpx));
			height: calc(160 * var(--rpx));
		}
	}
		
</style>