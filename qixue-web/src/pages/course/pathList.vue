<template>
	<div class="path-page">
		<mescroll-body
			ref="pathMescroll"
			@init="mescrollInit"
			@down="downCallback"
			@up="upCallback"
			:down="downOption"
			:up="upOption"
		>
			<div class="path-list">
				<div
					class="path-item"
					v-for="(item, index) in pathList"
					:key="item.id || index"
					@click="goToPathDetail(item)"
				>
					<div class="path-cover">
						<img :src="static(item.thumb || '/static/img/video-img-1.png')" class="path-img" loading="lazy"  style="object-fit:cover" />
					</div>
					<div class="path-info">
						<div class="path-title">{{ item.title || '学习路径' }}</div>
						<div class="path-desc">{{ item.desc || '' }}</div>
						<div class="path-meta">
							<span class="path-time" v-if="item.category">{{ item.category.name }}</span>
						</div>
					</div>
					<van-icon name="arrow" size="18" color="#b2b2b2" />
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
			pathList: [],
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
					tip: '暂无学习路径',
					icon: '/static/empty/empty.png'
				}
			}
		}
	},
	onLoad() {
		this.restorePathListCache();
	},
	methods: {
		restorePathListCache() {
			try {
				const cached = getStorage('qixue_path_list');
				if (cached && cached.length) {
					this.pathList = cached;
					return;
				}
				const boot = getStorage('qixue_home_boot');
				if (boot && boot.pathList && boot.pathList.length) {
					this.pathList = boot.pathList;
				}
			} catch (e) {}
		},
		downCallback() {
			this.mescroll.resetUpScroll();
		},
		upCallback(page) {
			this.getPathList(page);
		},
		getPathList(page) {
			this.$http.get('/path/getPathList', {
				page: page.num,
				limit: page.size
			}).then(res => {
				if (res.code === 200) {
					let list = res.data.list || []

					if (page.num === 1) {
						this.pathList = list;
						try {
							setStorage('qixue_path_list', list);
						} catch (e) {}
					} else {
						this.pathList = this.pathList.concat(list);
					}
					this.mescroll.endSuccess(list.length, list.length >= page.size);
				} else {
					this.msg(res.msg || '获取学习路径失败');
					this.mescroll.endErr();
				}
			}).catch(err => {
				console.log('获取学习路径失败:', err);
				this.msg('获取学习路径失败');
				this.mescroll.endErr();
			});
		},
		goToPathDetail(item) {
			const pathId = item.id;
			if (pathId) {
				this.navTo(`/pages/course/pathDetail?id=${pathId}`);
			}
		}
	}
}
</script>

<style lang="scss" scoped>
.path-page {
	min-height: 100vh;
	background: #f6f7fb;
}

.path-list {
	padding: calc(20 * var(--rpx)) calc(30 * var(--rpx));
}

.path-item {
	display: flex;
	align-items: center;
	padding: calc(24 * var(--rpx));
	background: #ffffff;
	border-radius: calc(16 * var(--rpx));
	margin-bottom: calc(20 * var(--rpx));

	&:last-child {
		margin-bottom: 0;
	}
}

.path-cover {
	width: calc(160 * var(--rpx));
	height: calc(120 * var(--rpx));
	border-radius: calc(12 * var(--rpx));
	overflow: hidden;
	flex-shrink: 0;
	background: #f5f5f5;
}

.path-img {
	width: 100%;
	height: 100%;
}

.path-info {
	flex: 1;
	margin-left: calc(24 * var(--rpx));
	display: flex;
	flex-direction: column;
	gap: calc(12 * var(--rpx));
}

.path-title {
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
	line-height: 1.4;
	overflow: hidden;
	display: -webkit-box;
	-webkit-line-clamp: 1;
	line-clamp: 1;
	-webkit-box-orient: vertical;
}

.path-desc {
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
	line-height: 1.4;
	overflow: hidden;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	line-clamp: 2;
	-webkit-box-orient: vertical;
}

.path-meta {
	display: flex;
	align-items: center;
	gap: calc(16 * var(--rpx));
	margin-top: calc(4 * var(--rpx));
}

.path-count {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #4285f4;
}

.path-time {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #b2b2b2;
}
</style>
