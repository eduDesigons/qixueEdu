<template>
	<div class="pannel">
		<mescroll-uni @init="mescrollInit" :height="height" :down="downOption" @down="downCallback" :up="upOption"
			@up="upCallback" @emptyclick="emptyClick">
			<div class="swiper-list">
				<course-item v-for="item in list" :info="item"></course-item>
			</div>
		</mescroll-uni>
	</div>
</template>

<script>
	import MescrollMixin from "@/uni_modules/mescroll-uni/components/mescroll-uni/mescroll-mixins.js";
	import MescrollMoreItemMixin from "@/uni_modules/mescroll-uni/components/mescroll-uni/mixins/mescroll-more-item.js";
	export default {
		mixins: [MescrollMixin, MescrollMoreItemMixin], // 注意此处还需使用MescrollMoreItemMixin (必须写在MescrollMixin后面)
		data() {
			return {
				downOption: {
					auto: false // 不自动加载 (mixin已处理第一个tab触发downCallback)
				},
				cursor: null,
				upOption: {
					auto: false, // 不自动加载
				},
				list: [], //列表数据
			}
		},
		props: {
			height: [Number, String], // mescroll的高度
			disableScroll: Boolean, // 是否禁止滚动, 默认false
			req: {
				type: Object,
				default: () => ({})
			},
			filterSet: {
				type: Object,
				default: () => ({})
			}
		},
		methods: {
			/*上拉加载的回调: 其中page.num:当前页 从1开始, page.size:每页数据条数,默认10 */
			upCallback(page) {
				//联网加载数据
				let data = {}
				if (page.num == 1) {
					this.cursor = null
				}
				this.$http.post('/webapp/home/getUserList', {
					...this.req,
					...this.filterSet,
					size: page.size,
					cursor: this.cursor
				}).then(res => {
					//设置列表数据
					if (page.num == 1) this.list = []; //如果是第一页需手动制空列表
					this.list = this.list.concat(res.data.list); //追加新数据
					this.cursor = res.data.cursor; //游标
					this.mescroll.endSuccess(res.data.count, !!res.data.cursor);
				}).catch(() => {
					//联网失败, 结束加载
					this.mescroll.endErr();
				})
			}
		}
	}
</script>

<style lang="scss" scoped>
	.swiper-list {
		padding: calc(20 * var(--rpx));
		padding-bottom: calc(120 * var(--rpx));
	}
</style>