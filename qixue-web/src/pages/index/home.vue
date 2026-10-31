<template>
	<HomeDesktop
		v-if="isDesktop"
		:loading="homeLoading"
		:banners="swiperList"
		:navs="navs"
		:notice="notice"
		:rank-tabs="rankTabs"
		:rank-tab-index="rankTabIndex"
		:rank-tags="rankTags"
		:rank-tag-index="rankTagIndex"
		:rank-list="rankList"
		:path-list="pathList"
		:hot-list="hotCourseList"
		:hot-loading="hotCourseLoading"
		:hot-has-more="hotCourseHasMore"
		:continue-study="continueStudy"
		@banner="bannerClick"
		@category="goCoursePage"
		@notice="openNotice"
		@rank-tab="rankTabClick"
		@rank-tag="rankTagClick"
		@more-hot="loadMoreHotCourse"
	/>
	<div v-else class="home-page">
		<div class="home-skel" v-if="homeLoading && !swiperList.length && !navs.length">
			<div class="home-skel-banner"></div>
			<div class="home-skel-nav">
				<div class="home-skel-nav-item" v-for="n in 4" :key="n"></div>
			</div>
			<div class="home-skel-card"></div>
		</div>
		<div class="banner" v-if="swiperList.length">
			<van-swipe class="home-swipe" :autoplay="swiperList.length > 1 ? 3000 : 0" indicator-color="#fff" @change="() => {}">
				<van-swipe-item v-for="(item, index) in swiperList" :key="index" @click="bannerClick(index)">
					<img :src="item.img_url" alt="" />
				</van-swipe-item>
			</van-swipe>
		</div>
		<div class="content-box">
			<div class="nav" v-if="navs.length">
				<div class="nav-item" v-for="nav in navs" :key="nav.id" @click="goCoursePage(nav.id)">
					<img :src="nav.icon && static(nav.icon)" class="nav-image" loading="lazy"  style="width:100%;height:auto;display:block" />
					<div class="nav-name">{{nav.name}}</div>
				</div>
			</div>
			<div class="home-continue" v-if="continueStudy && continueStudy.course" @click="openContinue">
				<img :src="static(continueStudy.course.thumb)" alt="" />
				<div>
					<strong>继续学习</strong>
					<span>{{ continueStudy.course.title }}</span>
					<em v-if="continueStudy.continueChapter && continueStudy.continueChapter.title">{{ continueStudy.continueChapter.title }}</em>
				</div>
			</div>
			<div class="notice" v-if="notice.id" @click="openNotice">
				<div class="notice-con">
					<div class="notice-text">
						{{notice.title}}
					</div>
					<van-icon name="arrow" color="#858a92" size="20" />
				</div>
			</div>

			<div class="section hot-rank-section" v-if="rankTabs.length">
				<div class="section-header">
					<div class="section-title">热学榜</div>
					<div class="section-subtitle" @click="navTo('/pages/index/course')">更多</div>
				</div>
				<div class="rank-tabs">
					<div class="rank-tab" :class="{active: rankTabIndex === index}" v-for="(tab, index) in rankTabs"
						:key="index" @click="rankTabClick(index)">
						{{ tab.name }}
					</div>
				</div>
				<div class="rank-tags qx-scroll-x">
					<div class="rank-tag" :class="{active: rankTagIndex === index}" v-for="(tag, index) in rankTags"
						:key="index" @click="rankTagClick(index)">
						{{ tag.name }}
					</div>
				</div>
				<div class="rank-list">
					<div class="rank-item" v-for="(item, index) in rankList" :key="index" @click="navTo('/pages/course/courseDetail?id='+item.id)">
						<div class="rank-cover">
							<img :src="static(item.thumb)" class="rank-img" loading="lazy"  style="object-fit:cover" />
							<div class="rank-num" :class="'rank-' + (index + 1)">{{ index + 1 }}</div>
						</div>
						<div class="rank-content">
							<div class="rank-title">{{ item.title }}</div>
							<div class="rank-author">{{ item.desc }}</div>
							<div class="rank-count">{{ item.credit || 0 }} 学分</div>
						</div>
					</div>
					<div class="rank-empty" v-if="!rankList.length">该分类暂无课程</div>
				</div>
			</div>

			<div class="section path-section" v-if="pathList.length">
				<div class="section-header">
					<div class="section-title">系列课程</div>
					<div class="section-more" @click="navTo('/pages/course/pathList')">更多</div>
				</div>
				<div class="course-grid">
					<div class="course-item" v-for="(item, index) in pathList" :key="index" @click="navTo('/pages/course/pathDetail?id='+item.id)">
						<div class="course-cover">
							<img :src="static(item.thumb)" class="course-img" loading="lazy"  style="object-fit:cover" />
						</div>
						<div class="course-info">
							<div class="course-title">{{ item.title }}</div>
							<div class="course-author">{{ item.desc }}</div>
						</div>
					</div>
				</div>
			</div>

			<div class="section recommend-section" v-if="hotCourseList.length">
				<div class="section-header">
					<div class="section-title">热门推荐</div>
					<div class="section-more" @click="navTo('/pages/index/course')">更多</div>
				</div>
				<div class="course-grid">
					<div class="course-item" v-for="(item, index) in hotCourseList" :key="index" @click="navTo('/pages/course/courseDetail?id='+item.id)">
						<div class="course-cover">
							<img :src="static(item.thumb)" class="course-img" loading="lazy"  style="object-fit:cover" />
							<div class="course-label" v-if="item.category && item.category.name">{{ item.category.name }}</div>
						</div>
						<div class="course-info">
							<div class="course-title">{{ item.title }}</div>
							<div class="course-author">{{ item.desc }}</div>
						</div>
					</div>
				</div>

				<div class="loading-more" v-if="hotCourseLoading">
					<span>加载中...</span>
				</div>
				<div class="no-more" v-if="!hotCourseHasMore && hotCourseList.length > 0">
					<span>没有更多了</span>
				</div>
			</div>
		</div>
	</div>

		<van-popup :show="!!homeAd.id" position="center" :close-on-click-overlay="false" style="background:transparent">
			<div class="ad-box">
				<img class="ad-img" :src="static(homeAd.img_url)" @click="homeAdClick()"  style="object-fit:contain" />
				<van-icon class="ad-close" name="cross" size="28" color="#fff" @click="homeAdClose()" />
			</div>
		</van-popup>
		<NoticePopup ref="noticePopup" :notice="notice" @confirm="confirmNotice" />
</template>

<script>
import { enterHome, leaveHome } from '@/common/pwa'
import { getStorage, setStorage } from '@/common/storage'
import { emit } from '@/common/bus'
import { switchTab } from '@/common/nav'
import { defineAsyncComponent } from 'vue'
import { isDesktop as desktopOn } from '@/pc/desktop'
import NoticePopup from './components/noticePopup.vue'
	export default {
		name: 'home',
		components: {
			NoticePopup,
			HomeDesktop: defineAsyncComponent(() => import('@/pc/HomeDesktop.vue')),
		},
		computed: {
			isDesktop() {
				return desktopOn.value
			},
		},
		data() {
			return {
				currentBanner: 0,
				swiperList: [],
				notice: {},
				homeAd: {},
				navs: [],
				rankTabIndex: 0,
				rankTagIndex: 0,
				rankTabs: [],
				rankTags: [],
				liveList: [],
				rankList: [],
				rankCache: {},
				pathList: [],
				hotCourseList: [],
				hotCoursePage: 1,
				hotCoursePageSize: 10,
				hotCourseHasMore: true,
				hotCourseLoading: false,
				homeLoading: false,
				continueStudy: null,
			};
		},
		onLoad() {
			this.initData()
		},
		onShow() {
			enterHome()
		},
		onHide() {
			leaveHome()
		},
		onPullDownRefresh() {
			this.initData()
		},
		onReachBottom() {
			if (this.hotCourseHasMore && !this.hotCourseLoading) {
				this.loadMoreHotCourse();
			}
		},
		methods: {
			applyHomeBoot(data, fromCache = false) {
				if (!data) return
				this.swiperList = (data.swiperList || []).map(item => ({
					...item,
					img_url: this.static(item.img_url)
				}))
				this.notice = data.notice || {}
				this.navs = data.categories || []
				this.rankTabs = this.navs
				this.pathList = data.pathList || []
				if (data.hotCourses) {
					this.hotCourseList = data.hotCourses.list || []
					this.hotCourseHasMore = data.hotCourses.hasMore !== false
					this.hotCoursePage = this.hotCourseList.length ? 2 : 1
				}
				if (this.navs.length) {
					this.rankTags = this.navs[0].children || []
					this.rankTabIndex = 0
					this.rankTagIndex = 0
				}
				if (!fromCache) {
					this.continueStudy = data.continueStudy || null
				}
				if (data.rankList && data.rankList.length) {
					this.rankList = data.rankList
					const first = this.navs[0] && this.navs[0].children && this.navs[0].children[0]
					if (first && first.id) this.rankCache[first.id] = data.rankList
				} else if (this.navs.length && !fromCache) {
					this.rankTabClick(0)
				}
				if (fromCache && this.navs.length && !this.rankList.length) {
					this.rankTabClick(0)
				}
				if (!fromCache) {
					this.initNotice()
				}
			},
			getHomeInfo() {
				this.restoreHomeCache()
				this.homeLoading = !this.swiperList.length && !this.navs.length
				this.$http.get('/home/boot').then(res => {
					this.applyHomeBoot(res.data)
					try {
						const cached = Object.assign({}, res.data)
						delete cached.continueStudy
						setStorage('qixue_home_boot', cached)
					} catch (e) {}
				}).catch(() => {
					if (!this.navs.length) {
						this.$http.get('/home/getHomeInfo').then(res => {
							this.swiperList = (res.data.swiperList || []).map(item => ({
								...item,
								img_url: this.static(item.img_url)
							}))
							this.notice = res.data.notice || {}
						})
						this.$http.get('/course/getCourseCategories').then(res => {
							this.navs = res.data || []
							this.rankTabs = this.navs
							this.initNotice()
							if (this.navs.length) this.rankTabClick(0)
						})
						this.getHotPathList()
						this.getHotCourseList()
					}
				}).finally(() => {
					this.homeLoading = false
				})
			},
			restoreHomeCache() {
				try {
					const cached = getStorage('qixue_home_boot')
					if (cached) {
						this.applyHomeBoot(cached, true)
					}
				} catch (e) {}
			},
			ensureNoticeContent() {
				if (!this.notice || !this.notice.id) return Promise.resolve()
				if (this.notice.content) return Promise.resolve()
				return this.$http.get('/home/getNoticeDetail', { id: this.notice.id }).then(res => {
					if (res.data) this.notice = Object.assign({}, this.notice, res.data)
				}).catch(() => {})
			},
			initNotice() {
				if (!this.notice || !this.notice.id) return
				const readId = getStorage('hasReadNoticeId')
				if (String(readId) === String(this.notice.id)) return
				this.ensureNoticeContent().then(() => {
					this.$nextTick(() => {
						if (this.$refs.noticePopup) this.$refs.noticePopup.show = true
					})
				})
			},
			openNotice() {
				this.ensureNoticeContent().then(() => {
					if (this.$refs.noticePopup) this.$refs.noticePopup.show = true
				})
			},
			confirmNotice() {
				if (!this.notice || !this.notice.id) return
				setStorage('hasReadNoticeId', String(this.notice.id))
			},
			bannerClick(e){
				const banner = this.swiperList[e]
				console.log(banner)
				if(!banner || banner.type == 1){
					return;
				}
				if(banner.type == 2){
					if(!banner.link_url) return;
					if(!banner.link_url.indexOf('http')){
						this.navTo(banner.link_url)
					}else{
						window.location.href = banner.link_url
						return
						this.navTo('/pages/h5/h5?url=' + encodeURIComponent(banner.link_url))
					}
					return;
				}
				if(banner.type == 3){
					this.navTo('/pages/h5/bannerDetail?id='+banner.id)
				}
			},
			getHotCourseList() {
				this.hotCoursePage = 1;
				this.hotCourseHasMore = true;
				this.hotCourseList = [];
				this.loadMoreHotCourse();
			},
			loadMoreHotCourse() {
				if (this.hotCourseLoading || !this.hotCourseHasMore) {
					return;
				}
				this.hotCourseLoading = true;
				this.$http.get('/home/getHotCourseList', {
					page: this.hotCoursePage,
					limit: this.hotCoursePageSize
				}).then(res => {
					const newData = res.data.list || [];
					if (newData.length > 0) {
						this.hotCourseList = this.hotCourseList.concat(newData);
						this.hotCoursePage++;
					}
					this.hotCourseHasMore = res.data.hasMore !== false;
					this.hotCourseLoading = false;
				}).catch(err => {
					console.log(err);
					this.hotCourseLoading = false;
				});
			},
			getHotPathList() {
				this.$http.get('/home/getHotPathList').then(res => {
					this.pathList = res.data || []
				})
			},
			initData() {
				this.getHomeInfo()
			},
			rankTabClick(rankTabIndex) {
				this.rankTabIndex = rankTabIndex
				const tab = this.navs[rankTabIndex]
				this.rankTags = (tab && tab.children) || []
				this.rankTagClick(0)
			},
			rankTagClick(rankTagIndex) {
				this.rankTagIndex = rankTagIndex
				const tab = this.navs[this.rankTabIndex]
				const tag = tab && tab.children && tab.children[rankTagIndex]
				if (!tag || !tag.id) {
					this.rankList = []
					return
				}
				if (this.rankCache[tag.id]) {
					this.rankList = this.rankCache[tag.id]
					return
				}
				this.$http.get('/home/getHotClassListByCateId', {id: tag.id})
				.then(res => {
					this.rankList = res.data || []
					this.rankCache[tag.id] = this.rankList
				})
				.catch(err => {
					console.log(err)
				})
			},
			openContinue() {
				const row = this.continueStudy || {}
				const course = row.course || {}
				const chapter = row.continueChapter
				if (chapter && chapter.id) {
					this.navTo('/pages/course/chapterDetail?id=' + chapter.id + '&courseId=' + (course.id || chapter.course_id || ''))
					return
				}
				if (course.id) {
					this.navTo('/pages/course/courseDetail?id=' + course.id)
				}
			},
			goCoursePage(cateId){
			switchTab('/pages/index/course')
			emit('coursePageChangeCategory', cateId)
		}
		}
	}
</script>
<style lang="scss">
	.qx-route-pages-index-home {
	background: #F6F7FB url('/static/bg/com.png') no-repeat 0 0 / contain fixed;
}
	.ad-box {
		padding: calc(40 * var(--rpx));
		display: flex;
		align-items: center;
		position: relative;

		.ad-img {
			max-width: 100%;
			border-radius: calc(20 * var(--rpx));
		}

		.ad-close {
			position: absolute;
			width: calc(60 * var(--rpx));
			left: calc(50% - calc(30 * var(--rpx)));
			bottom: calc(-60 * var(--rpx))
		}
	}

	.home-page {
		padding-bottom: calc(180 * var(--rpx));
	}

	.home-skel {
		padding: calc(80 * var(--rpx)) calc(20 * var(--rpx)) calc(20 * var(--rpx));
	}

	.home-skel-banner,
	.home-skel-nav-item,
	.home-skel-card {
		background: #eef0f4;
		animation: home-skel-pulse 1.2s ease-in-out infinite;
	}

	.home-skel-banner {
		height: calc(280 * var(--rpx));
		border-radius: calc(40 * var(--rpx));
	}

	.home-skel-nav {
		display: flex;
		justify-content: space-evenly;
		margin-top: calc(32 * var(--rpx));
	}

	.home-skel-nav-item {
		width: calc(96 * var(--rpx));
		height: calc(96 * var(--rpx));
		border-radius: calc(24 * var(--rpx));
	}

	.home-skel-card {
		height: calc(280 * var(--rpx));
		margin-top: calc(32 * var(--rpx));
		border-radius: calc(16 * var(--rpx));
	}

	@keyframes home-skel-pulse {
		0%, 100% { opacity: 1; }
		50% { opacity: 0.55; }
	}

	.banner {
		padding: calc(20 * var(--rpx));
		padding-top: calc(80 * var(--rpx));
	}

	.content-box {
		margin-top: calc(10 * var(--rpx));
		padding: 0 calc(20 * var(--rpx));
	}

	.nav {
		width: 100%;
		display: flex;
		justify-content: space-evenly;

		.nav-item {
			display: flex;
			flex-direction: column;
			align-items: center;
			padding: calc(20 * var(--rpx));

			.nav-image {
				width: 100%;
			}

			.nav-name {
				font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
				margin-top: calc(10 * var(--rpx));
				color: #404040;
			}
		}
	}

	.home-continue {
		display: flex;
		align-items: center;
		gap: calc(16 * var(--rpx));
		margin: calc(32 * var(--rpx)) 0;
		padding: calc(28 * var(--rpx)) calc(24 * var(--rpx));
		border-radius: calc(16 * var(--rpx));
		background: #fff;
		box-shadow: 0 1px 8px 0 rgba(233, 234, 237, .65);
	}
	.home-continue img {
		width: calc(96 * var(--rpx));
		height: calc(72 * var(--rpx));
		object-fit: cover;
		border-radius: calc(8 * var(--rpx));
		flex: none;
		background: #e8edf3;
	}
	.home-continue div {
		min-width: 0;
		display: flex;
		flex-direction: column;
	}
	.home-continue strong {
		color: #4285f4;
		font-size: calc(24 * var(--rpx));
	}
	.home-continue span,
	.home-continue em {
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		font-style: normal;
	}
	.home-continue span {
		margin-top: calc(4 * var(--rpx));
		color: #1a1a1a;
		font-size: calc(26 * var(--rpx));
	}
	.home-continue em {
		color: #8b93a1;
		font-size: calc(22 * var(--rpx));
	}

	.notice {
		margin: calc(32 * var(--rpx)) 0;
		position: relative;
		display: flex;
		height: auto;
		padding: calc(16 * var(--rpx)) calc(14 * var(--rpx)) calc(16 * var(--rpx)) calc(100 * var(--rpx));
		box-shadow: 0 1px 8px 0 rgba(233, 234, 237, .65);
		align-items: center;
		border-radius: calc(16 * var(--rpx));
		background: no-repeat url('/static/img/notice.png') left calc(10 * var(--rpx)) center / calc(80 * var(--rpx)) calc(80 * var(--rpx)),
			no-repeat url('/static/img/notice-bg.png') 100% / calc(240 * var(--rpx)) calc(96 * var(--rpx)) #fff;

		.notice-bg {
			width: calc(70 * var(--rpx));
			height: calc(48 * var(--rpx));
			position: absolute;
			left: 0;
			top: 0;
		}

		.notice-con {
			position: relative;
			box-sizing: border-box;
			display: flex;
			align-items: center;
			width: 100%;
			min-height: calc(66 * var(--rpx));
			height: auto;
		}

		.notice-text {
			box-sizing: border-box;
			position: relative;
			flex: 1;
			width: 0;
			padding: 0 10px 0;
			line-height: 16px;
			font-size: calc(12px * var(--font-scale, 1));
			font-weight: 400;
			color: #404040;
			word-break: break-all;
			overflow: hidden;
			display: -webkit-box;
			-webkit-line-clamp: 2;
			-webkit-box-orient: vertical;
			cursor: pointer;

			&::before {
				display: block;
				content: "";
				position: absolute;
				top: 50%;
				left: 0;
				width: calc(2 * var(--rpx));
				height: calc(56 * var(--rpx));
				background: #e9eaed;
				border-radius: calc(2 * var(--rpx));
				transform: translateY(-50%);
			}
		}
	}

	.section {
		background: #fff;
		border-radius: calc(16 * var(--rpx));
		padding: calc(28 * var(--rpx)) calc(24 * var(--rpx));
		margin-bottom: calc(24 * var(--rpx));
		box-shadow: 0 calc(2 * var(--rpx)) calc(16 * var(--rpx)) rgba(0, 0, 0, 0.04);
	}

	.section-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: calc(24 * var(--rpx));
	}

	.section-title {
		font-size: calc(calc(36 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 600;
		color: #1a1a1a;
	}

	.section-subtitle {
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #9ca3af;
	}

	.section-more {
		font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
		color: #9ca3af;
	}

	.live-section {
		padding-right: 0;
	}

	.live-list {
		display: flex;
		white-space: nowrap;
		padding-right: calc(24 * var(--rpx));
	}

	.live-item {
		display: inline-block;
		width: calc(320 * var(--rpx));
		margin-right: calc(24 * var(--rpx));
		border-radius: calc(12 * var(--rpx));
		overflow: hidden;
		box-shadow: 0 calc(4 * var(--rpx)) calc(16 * var(--rpx)) rgba(0, 0, 0, 0.08);
		background: #fff;
		flex-shrink: 0;
	}

	.live-cover {
		position: relative;
		width: 100%;
		height: calc(180 * var(--rpx));
		overflow: hidden;
	}

	.live-img {
		width: 100%;
		height: 100%;
	}

	.live-time {
		position: absolute;
		top: calc(16 * var(--rpx));
		left: 0;
		display: flex;
		align-items: center;
		padding: calc(6 * var(--rpx)) calc(16 * var(--rpx)) calc(6 * var(--rpx)) calc(12 * var(--rpx));
		border-radius: 0 calc(20 * var(--rpx)) calc(20 * var(--rpx)) 0;
		font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 500;

		&.purple {
			background: rgba(167, 139, 250, 0.9);
			color: #fff;
		}

		&.blue {
			background: rgba(66, 133, 244, 0.9);
			color: #fff;
		}

		.time-icon {
			margin-right: calc(6 * var(--rpx));
			font-size: calc(calc(20 * var(--rpx)) * var(--font-scale, 1));
		}
	}

	.live-tag {
		position: absolute;
		top: calc(16 * var(--rpx));
		right: 0;
		padding: calc(6 * var(--rpx)) calc(16 * var(--rpx));
		background: rgba(0, 0, 0, 0.5);
		color: #fff;
		font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
		border-radius: calc(20 * var(--rpx)) 0 0 calc(20 * var(--rpx));
	}

	.live-info {
		padding: calc(20 * var(--rpx));
	}

	.live-title {
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 500;
		color: #1a1a1a;
		line-height: 1.4;
		margin-bottom: calc(8 * var(--rpx));
		white-space: normal;
		overflow: hidden;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
	}

	.live-desc {
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #9ca3af;
		white-space: normal;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.rank-tabs {
		display:flex;align-items:stretch;box-sizing:border-box;min-height:calc(72 * var(--rpx));height:auto;padding:calc(4 * var(--rpx));border-radius:calc(32 * var(--rpx));background:#f6f7fb;
	}

	.rank-tab {
		flex: 1;
		display:flex;align-items:center;justify-content:center;box-sizing:border-box;padding:calc(12 * var(--rpx)) calc(8 * var(--rpx));text-align:center;color:#b9bdc2;line-height:1.3;font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));font-weight:500;

		&.active {
			border-radius:calc(32 * var(--rpx));background:#fff;color:#353535;
		}
	}

	.rank-tags {
		margin:15px calc(-20 * var(--rpx)) 0 calc(-20 * var(--rpx));
		display: flex;
		white-space: nowrap;
		margin-bottom: calc(32 * var(--rpx));
		padding-bottom: calc(8 * var(--rpx));
	}

	.rank-tag {
		display: inline-block;
		padding: calc(12 * var(--rpx)) calc(28 * var(--rpx));
		margin-right: calc(20 * var(--rpx));
		background: #f5f6fa;
		border-radius: calc(32 * var(--rpx));
		font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
		color: #666;
		transition: all 0.3s;
		border: calc(2 * var(--rpx)) solid transparent;

		&.active {
			background: rgba(66, 133, 244, 0.1);
			color: #4285f4;
			border-color: #4285f4;
		}
	}

	.rank-list {
		display: flex;
		flex-direction: column;
		gap: calc(28 * var(--rpx));
	}

	.rank-item {
		display: flex;
		margin-bottom: calc(10 * var(--rpx));
	}

	.rank-empty {
		padding: calc(48 * var(--rpx)) 0 calc(16 * var(--rpx));
		text-align: center;
		font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
		color: #9ca3af;
	}

	.rank-cover {
		position: relative;
		width: calc(180 * var(--rpx));
		height: calc(270 * var(--rpx));
		border-radius: calc(12 * var(--rpx));
		overflow: hidden;
		flex-shrink: 0;
	}

	.rank-img {
		width: 100%;
		height: 100%;
	}

	.rank-num {
		position: absolute;
		top: 0;
		left: 0;
		width: calc(40 * var(--rpx));
		height: calc(48 * var(--rpx));
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 700;
		color: #fff;
		background-size: contain;
		background-repeat: no-repeat;

		&.rank-1 {
			background: linear-gradient(135deg, #f59e0b, #fbbf24);
			border-radius: 0 0 calc(8 * var(--rpx)) 0;
		}

		&.rank-2 {
			background: linear-gradient(135deg, #94a3b8, #cbd5e1);
			border-radius: 0 0 calc(8 * var(--rpx)) 0;
		}

		&.rank-3 {
			background: linear-gradient(135deg, #f97316, #fb923c);
			border-radius: 0 0 calc(8 * var(--rpx)) 0;
		}
	}

	.rank-content {
		flex: 1;
		min-width: 0;
		margin-left: calc(24 * var(--rpx));
		display: flex;
		flex-direction: column;
		justify-content: flex-start;
		min-height: calc(270 * var(--rpx));
		height: auto;
		padding: calc(8 * var(--rpx)) 0;
	}

	.rank-title {
		flex-shrink: 0;
		font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 500;
		color: #1a1a1a;
		line-height: normal;
		white-space: normal;
		word-break: break-word;
		overflow: visible;
	}

	.rank-author {
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #9ca3af;
		overflow: hidden;
		word-break: break-word;
		overflow-wrap: anywhere;
		white-space: normal;
		display: -webkit-box;
		-webkit-line-clamp: 4;
		line-clamp: 4;
		-webkit-box-orient: vertical;
	}

	.rank-count {
		margin-top: auto;
		padding-top: calc(12 * var(--rpx));
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #9ca3af;
	}

	.path-list {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: calc(24 * var(--rpx));
	}

	.path-item {
		display: flex;
		align-items: center;
		padding: calc(24 * var(--rpx));
		background: #f8f9fc;
		border-radius: calc(12 * var(--rpx));
	}

	.path-icon {
		width: calc(80 * var(--rpx));
		height: calc(80 * var(--rpx));
		border-radius: calc(16 * var(--rpx));
		overflow: hidden;
		margin-right: calc(20 * var(--rpx));
		background: #e8f0fe;

		img {
			width: 100%;
			height: 100%;
		}
	}

	.path-info {
		flex: 1;
	}

	.path-title {
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 500;
		color: #1a1a1a;
		margin-bottom: calc(8 * var(--rpx));
	}

	.path-count {
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #9ca3af;
	}

	.path-arrow {
		width: calc(56 * var(--rpx));
		height: calc(56 * var(--rpx));
		border-radius: calc(12 * var(--rpx));
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 600;
		background: #e5e7eb;
		color: #6b7280;

		&.blue {
			background: #4285f4;
			color: #fff;
		}

		&.purple {
			background: #a78bfa;
			color: #fff;
		}

		&.orange {
			background: #f97316;
			color: #fff;
		}
	}

	.course-grid {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: calc(24 * var(--rpx));
	}

	.course-item {
		border-radius: calc(12 * var(--rpx));
		overflow: hidden;
	}

	.course-cover {
		position: relative;
		width: 100%;
		height: calc(190 * var(--rpx));
		border-radius: calc(12 * var(--rpx));
		overflow: hidden;
	}

	.course-img {
		width: 100%;
		height: 100%;
	}

	.course-label {
		position: absolute;
		top: calc(12 * var(--rpx));
		left: calc(12 * var(--rpx));
		padding: calc(6 * var(--rpx)) calc(12 * var(--rpx));
		background: rgba(0, 0, 0, 0.6);
		color: #fff;
		font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
		border-radius: calc(6 * var(--rpx));

		&.green {
			background: linear-gradient(135deg, #10b981, #34d399);
		}
	}

	.course-badge {
		position: absolute;
		top: calc(12 * var(--rpx));
		right: calc(12 * var(--rpx));
		padding: calc(6 * var(--rpx)) calc(12 * var(--rpx));
		background: linear-gradient(135deg, #10b981, #34d399);
		color: #fff;
		font-size: calc(calc(20 * var(--rpx)) * var(--font-scale, 1));
		border-radius: calc(6 * var(--rpx));
	}

	.course-info {
		padding: calc(16 * var(--rpx)) 0;
	}

	.course-title {
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 500;
		color: #1a1a1a;
		line-height: 1.4;
		margin-bottom: calc(8 * var(--rpx));
		overflow: hidden;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
	}

	.course-author {
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #9ca3af;
		overflow: hidden;
		word-break: break-word;
		overflow-wrap: anywhere;
		white-space: normal;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
	}

	.loading-more {
		text-align: center;
		padding: calc(40 * var(--rpx)) 0;
		font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
		color: #9ca3af;
	}

	.no-more {
		text-align: center;
		padding: calc(40 * var(--rpx)) 0;
		font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
		color: #9ca3af;
	}

	.interview-section {
		margin-bottom: 0;
	}
</style>