<template>
	<div>
		<van-nav-bar :title="info.name" left-arrow fixed placeholder @click-left="$router.back()" />
		<div class="content">
			<div class="qx-parse" style="user-select:text" v-html="richHtml(info.content)"></div>
		</div>
	</div>
</template>

<script>
	export default {
		data() {
			return {
				info: {
					name: '',
					content: ''
				}
			}
		},
		onLoad(e) {
			this.$http.get('/home/getBannerDetail?id='+e.id).then(res=>{
				this.info = res.data
			})
		}
	}
</script>

<style>
.content{
	padding: calc(30 * var(--rpx));
}
</style>
