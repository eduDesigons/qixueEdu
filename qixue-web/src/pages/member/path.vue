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
						<img :src="static((item.path && item.path.thumb) || '')" class="path-img" loading="lazy"  style="object-fit:cover" />
					</div>
					<div class="path-info">
						<div class="path-title">{{ (item.path && item.path.title) || '系列已下架' }}</div>
						<div class="path-desc">{{ (item.path && item.path.desc) || '' }}</div>
						<div class="path-meta">
							<span class="path-tag done" v-if="item.is_end">已结业</span>
							<span class="path-tag cert" v-if="item.user_cert">已获证书</span>
							<span class="path-time" v-if="item.created_at">{{ (item.created_at || '').slice(0, 10) }}</span>
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
		this.restoreMyPathCache();
	},
	methods: {
		restoreMyPathCache() {
			try {
				const cached = getStorage('qixue_my_path');
				if (cached && cached.length) {
					this.pathList = cached;
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
			this.$http.get('/user/myPath', {
				page: page.num,
				limit: page.size
			}).then(res => {
				if (res.code === 200) {
					let list = res.data.list || []

					if (page.num === 1) {
						this.pathList = list;
						try {
							setStorage('qixue_my_path', list);
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
			const pathId = item.path_id || (item.path && item.path.id);
			if (!pathId || !item.path) {
				this.msg('系列已下架');
				return;
			}
			this.navTo(`/pages/course/pathDetail?id=${pathId}`);
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
	gap: calc(12 * var(--rpx));
	margin-top: calc(4 * var(--rpx));
	flex-wrap: wrap;
}

.path-tag {
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	padding: calc(2 * var(--rpx)) calc(10 * var(--rpx));
	border-radius: calc(6 * var(--rpx));
	line-height: 1.5;
}

.path-tag.done {
	color: #389e0d;
	background: #f6ffed;
}

.path-tag.cert {
	color: #d48806;
	background: #fffbe6;
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
