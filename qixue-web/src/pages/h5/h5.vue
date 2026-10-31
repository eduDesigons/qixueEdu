<template>
	<web-view :src="webUrl"></web-view>
</template>

<script>
	export default {
		data() {
			return {
				webUrl: ''
			}
		},
		onLoad(option) {
			this.webUrl = decodeURIComponent(option.url) || ''
			console.log(this.webUrl)
		},
		onShow() {
		},
		onUnload() {
		}
	}
</script>

<style>
</style>
