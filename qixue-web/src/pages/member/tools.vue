<template>
	<div class="tools-page">
		<div class="tools-item" v-for="(item, index) in list" :key="index" @click="openTool(item)">
			<img class="tools-item__icon" :src="toolIcon(item)" style="object-fit:cover" @error="onToolIconError(item)" />
			<div class="tools-item__main">
				<span class="tools-item__title">{{ item.title }}</span>
				<span class="tools-item__desc">{{ item.desc }}</span>
			</div>
			<van-icon name="arrow" size="16" color="#b9bdc2" />
		</div>
		<div class="tools-loading" v-if="loading && !list.length">
			<div class="tools-loading__card" v-for="n in 3" :key="n"></div>
		</div>
		<div class="tools-empty" v-if="!loading && !list.length">
			<van-empty description="暂无工具" />
		</div>
	</div>
</template>

<script>
export default {
	data() {
		return {
			list: [],
			loading: true,
			brokenIcons: {}
		}
	},
	onShow() {
		this.loadList()
	},
	methods: {
		loadList() {
			this.loading = true
			this.$http.get('/sys/getTools').then((res) => {
				this.list = Array.isArray(res.data) ? res.data : []
			}).catch(() => {
				this.list = []
				this.msg('加载失败')
			}).finally(() => {
				this.loading = false
			})
		},
		toolIcon(item) {
			const src = this.static(item && item.icon)
			if (src && !this.brokenIcons[src]) return src
			return '/static/img/me/tools.svg'
		},
		onToolIconError(item) {
			const src = this.static(item && item.icon)
			if (!src || this.brokenIcons[src]) return
			this.brokenIcons = Object.assign({}, this.brokenIcons, { [src]: true })
		},
		openTool(item) {
			const url = String((item && item.url) || '').trim()
			if (!/^https?:\/\//i.test(url)) {
				this.msg('链接无效')
				return
			}
			window.open(url, '_blank')
		}
	}
}
</script>

<style lang="scss" scoped>
.tools-page {
	min-height: 100vh;
	box-sizing: border-box;
	padding: calc(24 * var(--rpx)) calc(30 * var(--rpx)) calc(calc(32 * var(--rpx)) + env(safe-area-inset-bottom));
	background: #f8f9fa;
}

.tools-loading {
	display: flex;
	flex-direction: column;
	gap: calc(16 * var(--rpx));
}

.tools-loading__card {
	height: calc(140 * var(--rpx));
	border-radius: calc(20 * var(--rpx));
	background: #fff;
}

.tools-item {
	display: flex;
	align-items: center;
	margin-bottom: calc(16 * var(--rpx));
	padding: calc(28 * var(--rpx));
	background: #fff;
	border-radius: calc(20 * var(--rpx));
}

.tools-item__icon {
	width: calc(88 * var(--rpx));
	height: calc(88 * var(--rpx));
	flex-shrink: 0;
	border-radius: calc(18 * var(--rpx));
	background: #f4f6f8;
}

.tools-item__main {
	flex: 1;
	min-width: 0;
	margin: 0 calc(20 * var(--rpx));
}

.tools-item__title {
	display: block;
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: #1a1a1a;
}

.tools-item__desc {
	display: block;
	margin-top: calc(8 * var(--rpx));
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	line-height: 1.5;
	color: #8a9199;
}

.tools-empty {
	padding-top: calc(120 * var(--rpx));
}
</style>
