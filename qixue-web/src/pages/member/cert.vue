<template>
	<div class="cert-page">
		<div class="cert-scroll qx-scroll-y">
			<div class="cert-list">
				<div
					class="cert-item"
					v-for="(item, index) in certList"
					:key="item.id || index"
					:style="{ animationDelay: Math.min(index, 8) * 40 + 'ms' }"
					@click="viewCert(item)"
				>
					<div class="cert-cover">
						<img :src="static(item.image || item.cert_image || item.thumb || '/static/img/cert-default.png')" class="cover-img"  style="object-fit:cover" />
						<div class="cert-badge" v-if="item.is_valid">
							<van-icon name="success" size="16" color="#52c41a" />
						</div>
					</div>
					<div class="cert-info">
						<div class="cert-title">{{ item.title || item.cert_name || item.name || '' }}</div>
						<div class="cert-course-name" v-if="item.source_title || item.course">
							{{ item.source_type === 'path' ? '系列' : (item.source_type === 'course' ? '课程' : '') }}{{ item.source_type ? ' · ' : '' }}{{ item.source_title || (item.course && item.course.title) || '' }}
						</div>
						<div class="cert-meta">
							<div class="cert-time">
								<van-icon name="calendar-o" size="14" color="#999" />
								<span>{{ (item.issued_at || item.created_at || '').slice(0, 10) }}</span>
							</div>
							<div class="cert-number" v-if="item.cert_no || item.cert_number">
								<span>证书编号：{{ item.cert_no || item.cert_number }}</span>
							</div>
						</div>
						<div class="cert-view-hint">查看证书</div>
					</div>
				</div>

				<div class="cert-empty" v-if="!loading && certList.length === 0">
					<van-empty description="暂无证书" image="/static/empty/attach.png" />
					<div class="empty-tip">学完绑定证书的系列或课程即可获得</div>
					<div class="empty-btns">
						<div class="empty-btn" @click="goToPath">
							<span>我的系列</span>
						</div>
						<div class="empty-btn empty-btn--ghost" @click="goToLesson">
							<span>我的课程</span>
						</div>
					</div>
				</div>

				<div class="cert-loading" v-if="loading">
					<van-loading type="circular" />
					<span>加载中...</span>
				</div>
			</div>
		</div>

		<div class="cert-modal" v-if="previewVisible" @click="closePreview">
			<div class="cert-modal__panel" @click.stop>
				<div class="cert-modal__head">
					<span class="cert-modal__title">{{ previewTitle }}</span>
					<div class="cert-modal__close" @click="closePreview">
						<van-icon name="cross" size="18" color="#666" />
					</div>
				</div>
				<div class="cert-modal__body">
					<img
						class="cert-modal__img"
						:src="previewUrl"
						:show-menu-by-longpress="true"
						@longpress="onLongPressSave"
					 style="width:100%;height:auto;display:block" />
					<span class="cert-modal__tip">{{ saveTip }}</span>
				</div>
				<div class="cert-modal__foot">
					<div class="cert-modal__btn" @click="downloadCert">
						<van-icon name="down" size="18" color="#fff" />
						<span>下载证书</span>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import { hideLoading, showLoading } from '@/common/toast'
import { isDesktop as desktopOn } from '@/pc/desktop'

export default {
	data() {
		return {
			certList: [],
			loading: false,
			refreshing: false,
			previewVisible: false,
			previewUrl: '',
			previewTitle: ''
		}
	},
	computed: {
		isDesktop() {
			return desktopOn.value
		},
		saveTip() {
			return this.isDesktop ? '下载后可在新页面保存图片' : '长按图片可保存到相册'
		}
	},
	onLoad() {
		this.getCertList();
	},
	onShow() {
		if (this.certList.length) {
			this.getCertList();
		}
	},
	methods: {
		onRefresh() {
			this.refreshing = true;
			this.getCertList().finally(() => {
				this.refreshing = false;
			});
		},
		getCertList() {
			this.loading = true;
			return this.$http.get('/user/myCert').then(res => {
				this.loading = false;
				if (res.code === 200) {
					if (Array.isArray(res.data)) {
						this.certList = res.data;
					} else if (res.data && res.data.list) {
						this.certList = res.data.list;
					} else if (res.data && res.data.data) {
						this.certList = res.data.data;
					} else {
						this.certList = [];
					}
				} else {
					this.msg(res.msg || '获取证书列表失败');
					this.certList = [];
				}
			}).catch(err => {
				this.loading = false;
				console.log('获取证书列表失败:', err);
				this.msg('获取证书列表失败');
				this.certList = [];
			});
		},
		goToPath() {
			this.navTo('/pages/member/path');
		},
		goToLesson() {
			this.navTo('/pages/member/lesson');
		},
		viewCert(item) {
			const imageUrl = this.static(item.image || item.cert_image || item.thumb || '');
			if (!imageUrl) {
				this.msg('暂无证书图片');
				return;
			}
			this.previewUrl = imageUrl;
			this.previewTitle = item.title || item.cert_name || item.name || '结业证书';
			this.previewVisible = true;
		},
		closePreview() {
			this.previewVisible = false;
		},
		onLongPressSave() {
			this.downloadCert();
		},
		downloadCert() {
			const imageUrl = this.previewUrl;
			if (!imageUrl) {
				this.msg('证书下载地址不存在');
				return;
			}

			showLoading('保存中...');

			const a = document.createElement('a');
			a.href = imageUrl;
			a.download = (this.previewTitle || '证书') + '.png';
			a.target = '_blank';
			a.rel = 'noopener';
			document.body.appendChild(a);
			a.click();
			document.body.removeChild(a);
			hideLoading();
			this.msg(desktopOn.value ? '已在新页面打开，可右键保存图片' : '请在打开的图片中长按保存');

		}
	}
}
</script>

<style lang="scss" scoped>
.cert-page {
	min-height: 100vh;
	background: #f6f7fb;
}

.cert-scroll {
	height: 100vh;
}

.cert-list {
	padding: calc(20 * var(--rpx)) calc(30 * var(--rpx));
	min-height: 100%;
}

.cert-item {
	background: #fff;
	border-radius: calc(16 * var(--rpx));
	margin-bottom: calc(20 * var(--rpx));
	overflow: hidden;
	padding: calc(24 * var(--rpx));
	display: flex;
	align-items: flex-start;
	box-shadow: 0 calc(2 * var(--rpx)) calc(12 * var(--rpx)) rgba(28, 43, 46, 0.04);
	animation: cert-in 0.36s cubic-bezier(0.32, 0.72, 0, 1) both;
	transition: transform 0.18s ease, box-shadow 0.18s ease;
}

.cert-item:active {
	transform: scale(0.985);
	box-shadow: 0 calc(1 * var(--rpx)) calc(6 * var(--rpx)) rgba(28, 43, 46, 0.06);
}

.cert-cover {
	position: relative;
	width: calc(200 * var(--rpx));
	height: calc(140 * var(--rpx));
	border-radius: calc(12 * var(--rpx));
	overflow: hidden;
	flex-shrink: 0;
	display: flex;
	align-items: center;
	justify-content: center;
	background: linear-gradient(135deg, #5c86ff 0%, #8bbeff 100%);
	box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(0, 0, 0, 0.1);
}

.cover-img {
	width: 100%;
	height: 100%;
	transition: transform 0.28s cubic-bezier(0.32, 0.72, 0, 1);
}

.cert-item:active .cover-img {
	transform: scale(1.04);
}

.cert-badge {
	position: absolute;
	top: calc(8 * var(--rpx));
	right: calc(8 * var(--rpx));
	width: calc(32 * var(--rpx));
	height: calc(32 * var(--rpx));
	background: rgba(255, 255, 255, 0.9);
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
}

.cert-info {
	flex: 1;
	margin-left: calc(24 * var(--rpx));
	display: flex;
	flex-direction: column;
	min-width: 0;
}

.cert-title {
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
	margin-bottom: calc(12 * var(--rpx));
	overflow: hidden;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	line-clamp: 2;
	-webkit-box-orient: vertical;
}

.cert-course-name {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
	margin-bottom: calc(12 * var(--rpx));
	overflow: hidden;
	display: -webkit-box;
	-webkit-line-clamp: 1;
	line-clamp: 1;
	-webkit-box-orient: vertical;
}

.cert-meta {
	display: flex;
	flex-direction: column;
	gap: calc(8 * var(--rpx));
}

.cert-time {
	display: flex;
	align-items: center;
	gap: calc(6 * var(--rpx));
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #999;
}

.cert-number {
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	color: #b9bdc2;
}

.cert-view-hint {
	margin-top: calc(16 * var(--rpx));
	align-self: flex-start;
	padding: calc(6 * var(--rpx)) calc(16 * var(--rpx));
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	color: #4285f4;
	background: #f0f7ff;
	border-radius: calc(8 * var(--rpx));
}

.cert-empty {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding: calc(200 * var(--rpx)) 0;
}

.empty-tip {
	margin-top: calc(24 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
	margin-bottom: calc(40 * var(--rpx));
}

.empty-btns {
	display: flex;
	gap: calc(24 * var(--rpx));
	align-items: center;
	justify-content: center;
}

.empty-btn {
	display: inline-flex;
	align-items: center;
	justify-content: center;
	padding: calc(20 * var(--rpx)) calc(60 * var(--rpx));
	background: #4285f4;
	border-radius: calc(20 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	line-height: 1;
	color: #fff;
	transition: transform 0.16s ease, opacity 0.16s ease;
}

.empty-btn:active {
	transform: scale(0.97);
	opacity: 0.92;
}

.empty-btn--ghost {
	background: #fff;
	color: #4285f4;
	border: calc(1 * var(--rpx)) solid #4285f4;
}

@keyframes cert-in {
	from {
		opacity: 0;
		transform: translateY(calc(16 * var(--rpx)));
	}
	to {
		opacity: 1;
		transform: none;
	}
}

.cert-loading {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	gap: calc(24 * var(--rpx));
	padding: calc(200 * var(--rpx)) 0;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
}

.cert-modal {
	position: fixed;
	inset: 0;
	z-index: 1000;
	display: flex;
	align-items: center;
	justify-content: center;
	padding: calc(48 * var(--rpx)) calc(40 * var(--rpx));
	background: rgba(0, 0, 0, 0.55);
	box-sizing: border-box;
}

.cert-modal__panel {
	width: 100%;
	max-width: calc(680 * var(--rpx));
	max-height: 86vh;
	background: #fff;
	border-radius: calc(24 * var(--rpx));
	overflow: hidden;
	display: flex;
	flex-direction: column;
	box-shadow: 0 calc(16 * var(--rpx)) calc(48 * var(--rpx)) rgba(0, 0, 0, 0.18);
}

.cert-modal__head {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: calc(28 * var(--rpx)) calc(32 * var(--rpx)) calc(16 * var(--rpx));
	gap: calc(16 * var(--rpx));
}

.cert-modal__title {
	flex: 1;
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: #1c2b2e;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.cert-modal__close {
	width: calc(56 * var(--rpx));
	height: calc(56 * var(--rpx));
	border-radius: 50%;
	background: #f3f4f6;
	display: flex;
	align-items: center;
	justify-content: center;
	flex-shrink: 0;
}

.cert-modal__body {
	padding: calc(8 * var(--rpx)) calc(32 * var(--rpx)) calc(16 * var(--rpx));
	overflow-y: auto;
	display: flex;
	flex-direction: column;
	align-items: center;
}

.cert-modal__img {
	width: 100%;
	border-radius: calc(12 * var(--rpx));
	background: #f6f7fb;
}

.cert-modal__tip {
	margin-top: calc(20 * var(--rpx));
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #9aa3a7;
}

.cert-modal__foot {
	padding: calc(16 * var(--rpx)) calc(32 * var(--rpx)) calc(32 * var(--rpx));
}

.cert-modal__btn {
	height: calc(88 * var(--rpx));
	border-radius: calc(44 * var(--rpx));
	background: #4285f4;
	color: #fff;
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	display: flex;
	align-items: center;
	justify-content: center;
	line-height: 1;
	gap: calc(10 * var(--rpx));
}

.cert-modal__btn:active {
	opacity: 0.9;
}
</style>
