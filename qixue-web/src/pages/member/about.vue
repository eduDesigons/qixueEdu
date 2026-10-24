<template>
	<div class="about-page">
		<div class="about-scroll qx-scroll-y">
			<div class="about-content">
				<!-- 加载中 -->
				<div class="about-loading" v-if="loading">
					<van-loading type="circular" />
					<span>加载中...</span>
				</div>
				
				<!-- 富文本内容 -->
				<div class="about-rich-text" v-if="!loading && content">
					<div class="qx-parse" style="user-select:text" v-html="richHtml(content)"></div>
				</div>
				
				<!-- 空状态 -->
				<div class="about-empty" v-if="!loading && !content">
					<van-empty description="暂无内容" />
				</div>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	data() {
		return {
			content: '',
			loading: false
		}
	},
	onLoad() {
		this.getAboutContent();
	},
	methods: {
		// 获取关于我们内容
		getAboutContent() {
			this.loading = true;
			this.$http.get('/sys/about').then(res => {
				this.loading = false;
				this.content = res.data.content || '';
			}).catch(err => {
				this.loading = false;
				this.content = '';
			});
		}
	}
}
</script>

<style lang="scss" scoped>
.about-page {
	min-height: 100vh;
	background: #f6f7fb;
}

.about-scroll {
	height: 100vh;
}

.about-content {
	padding: calc(30 * var(--rpx));
	min-height: 100%;
}

.about-rich-text {
	background: #fff;
	border-radius: calc(16 * var(--rpx));
	padding: calc(40 * var(--rpx));
	line-height: 1.8;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #404040;
	word-break: break-word;
}

.about-loading {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	gap: calc(24 * var(--rpx));
	padding: calc(200 * var(--rpx)) 0;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
}

.about-empty {
	padding: calc(200 * var(--rpx)) 0;
}
</style>