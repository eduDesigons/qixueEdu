<template>
	<div class="path-page">
		<van-nav-bar
			:title="path.title"
			:border="false"
			left-arrow
			fixed
			placeholder
			@click-left="$router.back()"
			style="--van-nav-bar-background:transparent;--van-nav-bar-icon-color:#fff;--van-nav-bar-title-text-color:#fff"
		/>
		<div class="path-skel" v-if="pageLoading && !path.id">
			<div class="path-skel-banner"></div>
			<div class="path-skel-line"></div>
			<div class="path-skel-line path-skel-line--short"></div>
			<div class="path-skel-card"></div>
		</div>
		<div class="path-load-fail" v-else-if="pageError && !path.id" @click="getCoursePath">
			<span>加载失败，点击重试</span>
		</div>
		<div v-else>
		<div class="path-details-title">
			<div class="path-details">
				<img :src="static(path.thumb)" class="path-details-img"  style="object-fit:cover" />
				<div class="path-details-title-content">
					<div class="path-details-title-row">
						<span class="path-details-title-text">{{ path.title }}</span>
					</div>
					<span class="path-details-desc-text">{{ path.desc }}</span>
				</div>
			</div>
			

			<div class="path-details-info">
				<div>
					<div class="path-details-content">
						<span class="primary-font">{{ (path.courses && path.courses.length) || 0 }}</span>
						<span class="primary-desc">门</span>
					</div>
					<div class="path-details-info-desc">包含课程</div>
				</div>
				<div>
					<div class="path-details-content">
						<span class="primary-font">{{ path.min_need_credit }}</span>
						<span class="primary-desc">分</span>
					</div>
					<div class="path-details-info-desc">要求学分</div>
				</div>
				<div>
					<div class="path-details-content">
						<span class="primary-font">{{ path.users_count }}</span>
						<span class="primary-desc">人</span>
					</div>
					<div class="path-details-info-desc">正在学习</div>
				</div>
			</div>
			<div class="path-cert-tip path-cert-tip--link" v-if="path.my_cert && path.my_cert.image" @click="viewMyCert">
				已获得「{{ path.cert && path.cert.name || '结业证书' }}」，点击查看
			</div>
			<div class="path-cert-tip" v-else-if="path.cert_id">
				完成系列全部课程后，可获得「{{ path.cert && path.cert.name || '结业证书' }}」
			</div>
		</div>
		<div class="path-details-content-wrap">
				<div class="qx-parse" v-html="richHtml(path.content)"></div>
				<div class="content-title"> 学习路线 </div>
				<div class="content-path" steps="">
					<div class="steps-content-vertical" steps-content="steps-content-vertical">
						<div class="step-vertical" step-style-attr="step-css-default-vertical" v-for="item in path.courses" :key="item.id">
							<div step-title="step-title-vertical"
								class="step-title step-title-vertical"></div>
							<div step-icon="step-icon-position-vertical" class="step-icon-position-vertical">
								<div class="step-icon" step-icon="step-icon-vertical"></div>
							</div>
							<div step-line="step-line-vertical"
								class="step-line step-line-vertical"></div>
							<div step-content="step-content-vertical" class="step-content">
								<div class="step-header">
									<span class="content-path-title">{{ item.title }}</span>
									<span class="content-path-time">{{ item.credit}}学分</span>
								</div>
								<div class="content-path-desc">
									{{ item.desc }}
								</div>
								<!-- <div class="content-nav">
									<div>零基础学Java 等4门课 </div>
									<div class="icon-box rotate">
										<span class="iconfont rotate-icon" icon=""></span>
									</div>
								</div> -->
								<div class="lessons">
									<div class="lesson-container" @click="navTo('/pages/course/courseDetail?id='+item.id)">
										<div class="lesson">
											<div class="lesson-item">
												<div class="lesson-img">
													<img
														:src="static(item.thumb)" loading="lazy"
														alt="课程图片"  />
													<div class="lesson-triangle"></div>
													<!-- <div class="lesson-tag">
														<div class="lesson-tag-mini-text"> {{ item.category.name }} </div>
													</div> -->
												</div>
												<div class="lesson-info">
													<div class="lesson-title">{{ item.title }}</div>
													<div class="lesson-author">{{ item.category.name }}</div>
												</div>
											</div>
											<!-- <div class="lesson-harvest">
												<div class="lesson-harvest-html">
													<div><span class="harvest-strong">你将获得</span>：</div>
													<div class="harvest-list">
														<div class="harvest-item">全面掌握 Java 核心语法；</div>
														<div class="harvest-item">玩转 Java 常用类库及工具；</div>
														<div class="harvest-item">攻克面向对象、多线程等技术难点；</div>
														<div class="harvest-item">独立用 Java 编写一款小游戏。</div>
													</div>
												</div>
											</div> -->
										</div>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		
		<!-- 底部加入按钮 -->
		<div class="path-footer" v-if="!path.has_join">
			<div v-if="path.apply_pending || path.applyPending" class="path-footer-btn path-footer-btn--disabled">
				<span>审核中，请耐心等待</span>
			</div>
			<div v-else-if="path.show_apply" class="path-footer-btn" :class="{ 'loading': joinLoading }" @click="joinPath">
				<van-loading v-if="joinLoading" size="16" color="#fff" type="circular" />
				<span v-if="!joinLoading">申请加入系列课程</span>
				<span v-else>申请中...</span>
			</div>
			<div v-else class="path-footer-btn" :class="{ 'loading': joinLoading }" @click="joinPath">
				<van-loading v-if="joinLoading" size="16" color="#fff" type="circular" />
				<span v-if="!joinLoading">加入系列课程</span>
				<span v-else>加入中...</span>
			</div>
		</div>
		</div>
	</div>
</template>

<script>
import { getStorage, setStorage } from '@/common/storage'
import { previewImage } from '@/common/dom'
	export default {
		data() {
			return {
				id: null,
				path: {
					title: '',
					courses: [],
					min_need_credit: 0,
					users_count: 0,
					has_join: false,
					show_apply: false,
				},
				isCollect: false,
				joinLoading: false,
				pageLoading: true,
				pageError: false,
				pathLoaded: false,
				needRefresh: false
			}
		},
		onLoad(options) {
			this.id = options.id;
			this.restorePathCache();
			this.getCoursePath();
		},
		onShow() {
			if (this.needRefresh) {
				this.needRefresh = false;
				this.getCoursePath();
			}
		},
		methods: {
			pathCacheKey() {
				return 'qixue_path_' + this.id;
			},
			restorePathCache() {
				if (!this.id) return;
				try {
					const cached = getStorage(this.pathCacheKey());
					if (cached && cached.id) {
						this.path = cached;
						this.pageLoading = false;
						this.pathLoaded = true;
					}
				} catch (e) {}
			},
			getCoursePath() {
				if (!this.id) return;
				if (!this.path.id) {
					this.pageLoading = true;
					this.pageError = false;
				}
				this.$http.get('/path/getPathDetail', { id: this.id }).then(res => {
					this.path = res.data || this.path;
					this.isCollect = res.data.is_collect || res.data.isCollect || false;
					this.pathLoaded = true;
					this.pageError = false;
					try {
						setStorage(this.pathCacheKey(), this.path);
					} catch (e) {}
				}).catch(() => {
					if (!this.path.id) this.pageError = true;
				}).finally(() => {
					this.pageLoading = false;
				});
			},
			viewMyCert() {
				const imageUrl = this.static((this.path.my_cert && this.path.my_cert.image) || '');
				if (!imageUrl) {
					this.msg('暂无证书图片');
					return;
				}
				previewImage({
					urls: [imageUrl],
					current: imageUrl
				});
			},
			// 加入学习路线
			joinPath() {
				if (!this.id) return;
				// 防止连点
				if (this.joinLoading) return;
				
				if (this.path.apply_pending || this.path.applyPending) {
					this.msg('申请审核中，请耐心等待');
					return;
				}
				if (this.path.show_apply) {
					this.needRefresh = true;
					const title = encodeURIComponent(this.path.title || '');
					this.navTo('/pages/course/applyJoin?pathId=' + this.id + '&pathTitle=' + title);
					return;
				}
				
				// 开始加载
				this.joinLoading = true;
				
				this.$http.post('/path/joinPath', {
					id: this.id
				}).then(res => {
					this.msg('加入成功');
					this.path.has_join = true;
					this.path.show_apply = false;
				}).catch(err => {
					console.log(err);
					this.msg(err.message || '加入失败');
				}).finally(() => {
					// 结束加载
					this.joinLoading = false;
				});
			}
		}
	}
</script>

<style>
	.path-page {
		min-height: 100vh;
		background: linear-gradient(135deg, #5c86ff, #8bbeff);
		--26dcd839: #5c86ff;
	}
	
	.steps-content-vertical {
		padding-left: 28px;
	}

	/*! CSS Used from: Embedded */
	.step-vertical {
		position: relative;
		line-height: 21px;
		--6697b51c: 21px;
		--190e58b0: 13px;
		--6d69e864: -15px;
		--7810f2bc: 0px;
		--2624a800: -21px;
		--212099ba: -28px;
		--2f579da6: 13px;
	}

	.step-vertical .step-line {
		background-color: #e9eaed;
	}

	.step-vertical .step-line-vertical {
		position: absolute;
		top: var(--190e58b0);
		left: var(--6d69e864);
		width: 1px;
		height: 100%;
	}

	.step-vertical .step-icon {
		z-index: 10;
		display: inline-block;
		width: 7px;
		height: 7px;
		margin-top: 4px;
		background: #fff;
		border: 2px solid #e9eaed;
		border-radius: 50%;
	}

	.step-vertical .step-icon-position-vertical {
		position: absolute;
		top: var(--7810f2bc);
		left: var(--2624a800);
		z-index: 10;
	}

	.step-vertical .step-title {
		line-height: 18px;
		font-size: calc(13px * var(--font-scale, 1));
		font-weight: 400;
		color: #7f848c;
	}

	.step-vertical .step-title-vertical {
		box-sizing: border-box;
		position: absolute;
		top: var(--7810f2bc);
		left: var(--212099ba);
		width: var(--2f579da6);
		height: 100%;
		padding: 0 5px;
		word-break: break-all;
		text-align: center;
	}

	.step-vertical .step-content {
		overflow: hidden;
		font-size: calc(15px * var(--font-scale, 1));
		font-weight: 400;
	}

	/*! CSS Used from: Embedded */
	.lesson-container {
		margin: 0 -10px;
		margin-top: 20px;
	}

	.lesson {
		padding: 14px 12px 18px 12px;
		margin: 0 10px;
		border-radius: 10px;
		box-shadow: 0 0 15px 0 rgba(233, 234, 237, .65);
	}

	.lesson-item {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		background-color: #fff;
	}

	.lesson-img {
		position: relative;
		flex-shrink: 0;
		width: 82px;
		height: 82px;
		margin-right: 16px;
	}

	.lesson-img>img {
		width: 82px;
		height: 82px;
		border-radius: 4px;
	}

	.lesson-info {
		overflow: hidden;
		height: 82px;
		flex: 1;
		display: flex;
		flex-direction: column;
		justify-content: center;
	}

	.lesson-title {
		line-height: 21px;
		color: #404040;
		font-size: calc(18px * var(--font-scale, 1));
		font-weight: 500;
	}

	.lesson-author {
		margin-top: 10px;
		line-height: 17px;
		color: #7f848c;
		font-size: calc(12px * var(--font-scale, 1));
		font-weight: 400;
	}

	.lesson-tag {
		position: absolute;
		top: 0;
		left: -2px;
		width: 32px;
		height: 14px;
		background: #4285f4;
		border-radius: 4px 4px 4px 0;
		line-height: 14px;
		text-align: center;
		color: #fff;
		z-index: 10;
	}

	.lesson-tag-mini-text {
		margin: 0 -1px;
		transform: scale(.83);
		font-size: calc(12px * var(--font-scale, 1));
		font-weight: 400;
	}

	.lesson-triangle {
		width: 0;
		height: 0;
		position: absolute;
		top: 12px;
		left: -4px;
		border: 2px solid transparent;
		border-right: 2px solid #2261c8;
		z-index: 9;
		font-size: 0;
		font-weight: 400;
	}

	.lesson-harvest {
		margin-top: 30px;
		color: #7f848c;
		font-size: calc(12px * var(--font-scale, 1));
		font-weight: 400;
	}

	.lesson-harvest .harvest-strong {
		line-height: 18px;
		color: #404040;
		font-size: calc(13px * var(--font-scale, 1));
		font-weight: 500;
	}

	.lesson-harvest-html {
		margin-top: 3px;
		line-height: 23px;
		color: #404040;
	}

	.harvest-list {
		list-style: none;
		margin: 0;
		padding: 0;
	}

	.harvest-item {
		margin: 0;
		padding: 0;
	}

	/*! CSS Used from: Embedded */
	.nav-bar-title {
		transition: all .3s ease;
		overflow: hidden;
		word-break: keep-all;
		white-space: nowrap;
		text-overflow: ellipsis;
	}

	.nav-bar-placeholder,
	.nav-bar-shim {
		height: 44px;
	}

	.nav-bar-shim {
		position: fixed;
		left: 0;
		right: 0;
		z-index: 20000;
	}

	.nav-bar-header {
		position: relative;
		display: flex;
		height: 44px;
	}

	.nav-bar-title {
		box-sizing: border-box;
		width: 100%;
		padding: 0 80px;
		color: #353535;
		text-align: center;
		font-size: calc(19px * var(--font-scale, 1));
		font-weight: 500;
		line-height: 44px;
	}

	.nav-bar-end-container,
	.nav-bar-pre-container {
		position: absolute;
		top: 0;
		z-index: 10;
		display: flex;
		align-items: center;
		height: 44px;
		color: inherit;
	}

	.nav-bar-pre-container {
		left: 0;
		padding-left: 14px;
	}

	.nav-bar-end-container {
		right: 15px;
	}

	.nav-bar-back {
		margin-left: -5px;
		font-size: calc(22px * var(--font-scale, 1));
		color: #b9bdc2;
		line-height: 1;
		color: #fff;
	}

	/*! CSS Used from: Embedded */
	.path-details-nav [native-top-bar-icon],
	.path-details-nav [native-top-bar-title] {
		color: #fff;
	}

	.path-details-title {
		background: var(--11a63af6);
		padding: calc(20 * var(--rpx)) calc(30 * var(--rpx));
	}

	.path-details {
		display: flex;
		align-items: center;
		gap: calc(24 * var(--rpx));
	}

	.path-details-img {
		width: calc(200 * var(--rpx));
		height: calc(150 * var(--rpx));
		border-radius: calc(12 * var(--rpx));
		flex-shrink: 0;
	}

	.path-details-title-content {
		padding-left: calc(10 * var(--rpx));
		flex: 1;
		display: flex;
		flex-direction: column;
		justify-content: center;
		min-width: 0;
	}

	.path-details-title-row {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: calc(16 * var(--rpx));
		margin-bottom: calc(12 * var(--rpx));
	}

	.path-details-title-text {
		flex: 1;
		font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 500;
		color: #fff;
		line-height: 1.4;
		overflow: hidden;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
	}

	.path-details-collect-btn {
		flex-shrink: 0;
		width: calc(60 * var(--rpx));
		height: calc(60 * var(--rpx));
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 50%;
		background: rgba(255, 255, 255, 0.2);
	}

	.path-details-desc-text {
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: rgba(255, 255, 255, 0.9);
		line-height: 1.5;
		overflow: hidden;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
	}

	.path-details-title-box {
		box-sizing: border-box;
		width: 100%;
		height: 100%;
		background-size: cover;
	}

	.path-details-info {
		box-sizing: border-box;
		display: flex;
		justify-content: space-between;
		align-items: center;
		height: 100px;
		padding: 1px 15px 12px 15px;
		text-align: center;
		color: #fff;
	}

	.path-details-info>div {
		width: 115px;
	}

	.path-details-info-desc {
		margin-top: 3px;
		line-height: 17px;
		font-size: calc(12px * var(--font-scale, 1));
		font-weight: 400;
	}

	.path-cert-tip {
		margin: 0 15px 16px;
		padding: 10px 12px;
		border-radius: 8px;
		background: rgba(255, 255, 255, 0.18);
		color: #fff;
		font-size: calc(13px * var(--font-scale, 1));
		line-height: 20px;
	}
	.path-cert-tip--link {
		text-decoration: underline;
	}

	.primary-font {
		font-family: DIN;
		line-height: 28px;
		font-size: calc(24px * var(--font-scale, 1));
		font-weight: 400;
	}

	.path-details-content {
		display: flex;
		justify-content: center;
		align-items: flex-end;
	}

	.primary-desc {
		margin-left: 2px;
		line-height: 21px;
		transform: scale(.91);
		font-size: calc(12px * var(--font-scale, 1));
		font-weight: 400;
	}

	.path-details-content-wrap {
		padding: 20px 15px 113px 0;
		margin-top: -12px;
		border-radius: 12px 12px 0 0;
		background-color: #fff;
		padding-bottom: calc(113px + calc(128 * var(--rpx)) + env(safe-area-inset-bottom));
	}

	.content-title {
		position: relative;
		padding-left: 15px;
		line-height: 25px;
		color: #353535;
		font-size: calc(18px * var(--font-scale, 1));
		font-weight: 500;
	}

	.content-title:before {
		content: "";
		position: absolute;
		left: 0;
		top: 50%;
		width: 3px;
		height: 16px;
		background: var(--26dcd839);
		border-radius: 0 5px 5px 0;
		transform: translateY(-50%);
	}

	.content-path {
		margin-top: 26px;
		--cf5ca076: 28px;
	}

	.content-path [step-icon=step-icon-position-vertical] {
		left: -20px;
	}

	.content-path [step-icon=step-icon-position-vertical] [step-icon=step-icon-vertical] {
		width: 7px;
		height: 7px;
		border: 2px solid var(--26dcd839);
	}

	.content-path [step-content=step-content-vertical] {
		overflow: visible;
	}

	.content-path [step-style-attr=step-css-default-vertical]:last-child [step-line=step-line-vertical] {
		display: none;
	}

	.content-path-title {
		margin-right: 7px;
		color: #353535;
		line-height: 22px;
		font-size: calc(16px * var(--font-scale, 1));
		font-weight: 500;
	}

	.content-path-time {
		color: #7f848c;
		line-height: 17px;
		font-size: calc(12px * var(--font-scale, 1));
		font-weight: 400;
	}

	.content-path-desc {
		margin-top: 15px;
		line-height: 25px;
		color: #7f848c;
		font-size: calc(14px * var(--font-scale, 1));
		font-weight: 400;
		word-break: break-word;
		overflow-wrap: anywhere;
		white-space: normal;
	}

	.content-nav {
		display: flex;
		justify-content: space-between;
		padding: 12px 16px 12px 16px;
		margin-top: 13px;
		background: #f6f7fb;
		border-radius: 10px;
		line-height: 20px;
		color: #404040;
		font-size: calc(14px * var(--font-scale, 1));
		font-weight: 500;
	}

	.rotate-icon {
		color: #7f848c;
		font-size: calc(14px * var(--font-scale, 1));
		font-weight: 300;
	}

	.icon-box {
		display: flex;
		align-items: center;
		margin-left: 16px;
	}

	.rotate {
		transform: rotate(180deg);
	}

	.lessons {
		padding-bottom: 48px;
	}

	.step-header {
		display: flex;
		align-items: center;
	}

	/* 底部加入按钮 */
	.path-footer {
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
		padding: calc(20 * var(--rpx)) calc(30 * var(--rpx));
		padding-bottom: calc(calc(20 * var(--rpx)) + env(safe-area-inset-bottom));
		background: #fff;
		box-shadow: 0 calc(-2 * var(--rpx)) calc(10 * var(--rpx)) rgba(0, 0, 0, 0.1);
		z-index: 100;
	}

	.path-footer-btn {
		width: 100%;
		height: calc(88 * var(--rpx));
		display: flex;
		align-items: center;
		justify-content: center;
		background: linear-gradient(135deg, #5c86ff, #8bbeff);
		border-radius: calc(44 * var(--rpx));
		font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 500;
		color: #fff;
		transition: opacity 0.3s;
		gap: calc(12 * var(--rpx));
	}

	.path-footer-btn.loading,
	.path-footer-btn--disabled {
		opacity: 0.7;
		pointer-events: none;
	}

	.path-footer-btn span {
		line-height: 1;
	}

	.path-skel {
		padding: calc(40 * var(--rpx)) calc(30 * var(--rpx));
	}

	.path-skel-banner,
	.path-skel-line,
	.path-skel-card {
		background: rgba(255, 255, 255, 0.35);
		animation: path-skel-pulse 1.2s ease-in-out infinite;
	}

	.path-skel-banner {
		height: calc(280 * var(--rpx));
		border-radius: calc(16 * var(--rpx));
	}

	.path-skel-line {
		height: calc(28 * var(--rpx));
		margin-top: calc(24 * var(--rpx));
		border-radius: calc(6 * var(--rpx));
	}

	.path-skel-line--short {
		width: 50%;
	}

	.path-skel-card {
		height: calc(220 * var(--rpx));
		margin-top: calc(40 * var(--rpx));
		border-radius: calc(16 * var(--rpx));
	}

	.path-load-fail {
		padding: calc(160 * var(--rpx)) calc(40 * var(--rpx));
		text-align: center;
		color: #fff;
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	}

	@keyframes path-skel-pulse {
		0%, 100% { opacity: 1; }
		50% { opacity: 0.55; }
	}
</style>