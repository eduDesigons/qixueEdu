<template>
	<CourseDesktop
		v-if="isDesktop"
		:categories="categoryTree"
		:level1="currentLevel1"
		:level2="currentLevel2Id"
		:sort-type="sortType"
		:course-type="courseType"
		:list="courseList"
		:loading="loading"
		:no-more="noMore"
		:keyword="searchKeyword"
		@level1="selectLevel1"
		@level2="selectLevel2"
		@sort="setSortType"
		@type="selectCourseType"
		@more="loadMore"
		@open="goDetail"
		@clear="clearSearch"
		@reset="resetFilters"
	/>
	<div v-else class="course-page">

		<div class="custom-navbar">

			<div class="navbar-status-bar" :style="{ height: statusBarHeight + 'px' }"></div>

			<div class="navbar-content">

				<div class="navbar-left">
					<span class="navbar-left-text">课程资源</span>
				</div>

				<div class="navbar-center">
					<div v-if="showSearch" class="navbar-search">
						<van-search
							v-model="searchKeyword"
							placeholder="搜索课程名称"
							shape="round"
							background="transparent"
							:show-action="!!(searchKeyword && searchKeyword.trim())"
							style="width:100%;--van-search-content-background:#f5f5f5"
							@search="handleSearch"
							@clear="clearSearch"
							@focus="handleSearchFocus"
							@blur="handleSearchBlur"
						>
							<template #action>
								<div @click="handleSearchAction">搜索</div>
							</template>
						</van-search>
					</div>
				</div>

				<div class="navbar-right" v-if="!showSearch" @click="toggleSearch">
					<van-icon name="search" size="20" color="#333" />
				</div>
			</div>
		</div>

		<div class="navbar-placeholder" :style="{ height: navbarHeight + 'px' }"></div>

		<div class="filter-bar" id="filterBar">
			<div class="filter-item" :class="{active: showCategoryPopup || isCategorySelected}" @click="toggleCategoryPopup">
				<span>{{ currentCategoryName }}</span>
				<span class="filter-icon" :class="{up: showCategoryPopup}">
					▽
				</span>
			</div>
			<div class="filter-item" :class="{active: sortType === 'hot'}" @click="setSortType('hot')">
				<span>最热</span>
			</div>
			<div class="filter-item" :class="{active: sortType === 'new'}" @click="setSortType('new')">
				<span>最新</span>
			</div>
			<div class="filter-item" :class="{active: showFilterPopup || courseType !== 0}" @click="toggleFilterPopup">
				<span>{{ courseTypeName }}</span>
				<span class="filter-icon" :class="{up: showFilterPopup}">▽</span>
			</div>
		</div>

		<div class="course-list qx-scroll-y" @scroll="onCourseScroll">
			<div class="course-item" v-for="(course, index) in courseList" :key="course.id || index" @click="goDetail(course)">
				<div class="course-cover">
					<img :src="static(course.cover || course.thumb)" class="cover-img" loading="lazy"  style="object-fit:cover" />
					<div class="course-tag" v-if="course.category && course.category.name">{{ course.category.name }}</div>
				</div>
				<div class="course-info">
					<div class="course-title">{{ course.title }}</div>
					<div class="course-author">{{ course.author || course.desc }}</div>
					<div class="course-source">
						<span class="source-text">{{ course.credit }} 学分</span>
					</div>
				</div>
			</div>
			<div class="course-skel" v-if="loading && courseList.length === 0">
				<div class="course-skel-item" v-for="n in 4" :key="n">
					<div class="course-skel-cover"></div>
					<div>
						<div class="course-skel-line"></div>
						<div class="course-skel-line course-skel-line--short"></div>
					</div>
				</div>
			</div>
			<div class="course-empty" v-else-if="!loading && courseList.length === 0">
				<van-empty description="没有找到相关课程" image="/static/empty/course.png" :image-size="120" />
			</div>
			<div class="loading-more" v-if="loading && courseList.length > 0">
				<span>加载中...</span>
			</div>
			<div class="no-more" v-if="noMore && courseList.length > 0">
				<span>没有更多了</span>
			</div>
		</div>

		<div class="category-popup" v-if="showCategoryPopup">
			<div class="category-mask" @click="showCategoryPopup = false"></div>
			<div class="category-content">
				<div class="category-columns">

					<div class="category-column first">

						<div class="category-item"
							:class="{active: currentLevel1 === -1}"
							@click="selectLevel1(-1)">
							全部
						</div>
						<div class="category-item"
							v-for="(item, index) in categoryTree"
							:key="item.id || index"
							:class="{active: currentLevel1 === index}"
							@click="selectLevel1(index)">
							{{ item.name }}
						</div>
					</div>

					<div class="category-column second">

						<div class="category-item"
							:class="{active: currentLevel2Id === 'all'}"
							@click="selectLevel2({id: 'all', name: '全部', isAll: true})">
							全部
							<span class="check" v-if="currentLevel2Id === 'all'">✓</span>
						</div>
						<div class="category-item"
							v-for="(item, index) in categoryLevel2"
							:key="item.id || index"
							:class="{active: currentLevel2Id === item.id}"
							@click="selectLevel2(item)">
							{{ item.name }}
							<span class="check" v-if="currentLevel2Id === item.id">✓</span>
						</div>
					</div>
				</div>
			</div>
		</div>

		<div class="filter-popup" v-if="showFilterPopup">
			<div class="filter-mask" @click="showFilterPopup = false"></div>
			<div class="filter-content">
				<div class="filter-options">
					<div class="filter-option-item"
						:class="{active: courseType === 0}"
						@click="selectCourseType(0)">
						全部
						<span class="check" v-if="courseType === 0">✓</span>
					</div>
					<div class="filter-option-item"
						:class="{active: courseType === 1}"
						@click="selectCourseType(1)">
						视频课
						<span class="check" v-if="courseType === 1">✓</span>
					</div>
					<div class="filter-option-item"
						:class="{active: courseType === 2}"
						@click="selectCourseType(2)">
						音频课
						<span class="check" v-if="courseType === 2">✓</span>
					</div>
					<div class="filter-option-item"
						:class="{active: courseType === 3}"
						@click="selectCourseType(3)">
						图文课
						<span class="check" v-if="courseType === 3">✓</span>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import { getStorage, setStorage } from '@/common/storage'
import { getSystemInfoSync, upx2px } from '@/common/dom'
import { emit, on } from '@/common/bus'
import { navigateTo } from '@/common/nav'
import { defineAsyncComponent } from 'vue'
import { isDesktop as desktopOn } from '@/pc/desktop'
export default {
	name: 'course',
	components: {
		CourseDesktop: defineAsyncComponent(() => import('@/pc/CourseDesktop.vue')),
	},
		data() {
			return {
			navIndex: 0,
			sortType: 'hot',
			showCategoryPopup: false,
			showFilterPopup: false,
			currentLevel1: -1,
			currentLevel2Id: null,
			currentCategory: null,
			courseType: 0,
			loading: false,
			noMore: false,
			page: 1,
			categoryTree: [],
			courseList: [],
			showSearch: false,
			searchKeyword: '',
			searchFocus: false,
			statusBarHeight: 0,
			listReqId: 0
		};
	},
	computed: {
		isDesktop() {
			return desktopOn.value
		},
		categoryLevel2() {
			if (this.currentLevel1 === -1) return [];
			if (this.categoryTree.length === 0) return [];
			const currentCategory = this.categoryTree[this.currentLevel1];
			if (!currentCategory) return [];
			return currentCategory.children || [];
		},
		currentCategoryName() {
			if (this.currentCategory && !this.currentCategory.isAll) {
				return this.currentCategory.name;
			}
			if (this.currentLevel2Id === 'all' && this.currentLevel1 >= 0 && this.categoryTree.length > 0 && this.categoryTree[this.currentLevel1]) {
				return this.categoryTree[this.currentLevel1].name;
			}
			if (this.currentLevel1 === -1) {
				return '全部分类';
			}
			if (this.currentLevel2Id === null && this.categoryTree.length > 0 && this.currentLevel1 >= 0 && this.categoryTree[this.currentLevel1]) {
				return this.categoryTree[this.currentLevel1].name;
			}
			if (this.categoryTree.length > 0 && this.currentLevel1 >= 0 && this.categoryTree[this.currentLevel1]) {
				return this.categoryTree[this.currentLevel1].name;
			}
			return '全部分类';
		},
		courseTypeName() {
			const typeMap = {
				0: '筛选',
				1: '视频课',
				2: '音频课',
				3: '图文课'
			};
			return typeMap[this.courseType] || '筛选';
		},
		isCategorySelected() {
			return this.currentLevel1 !== -1;
		},
		navbarHeight() {
			const contentHeight = upx2px(88);
			return this.statusBarHeight + contentHeight;
		}
	},
	onLoad() {
		const systemInfo = getSystemInfoSync();
		this.statusBarHeight = systemInfo.statusBarHeight || 0;

		this.restoreCourseListCache();
		this.getCourseCategories();
		on('coursePageSearch', (keyword) => {
			this.searchKeyword = keyword || ''
			this.showSearch = !!this.searchKeyword
			emit('courseSearchSync', this.searchKeyword)
			this.refreshList()
		});
		on('coursePageChangeCategory', (data) => {
			console.log('coursePageChangeCategory', data);

			const cateId = data && typeof data === 'object' ? data.cateId : data;
			if (this.categoryTree.length > 0) {
				this.handleCategoryChange(cateId);
			} else {
				const checkInterval = setInterval(() => {
					if (this.categoryTree.length > 0) {
						clearInterval(checkInterval);
						this.handleCategoryChange(cateId);
					}
				}, 100);
				// 5秒后清除定时器，避免无限等待
				setTimeout(() => clearInterval(checkInterval), 5000);
			}
		});
	},
		methods: {
		toggleSearch() {
			if (!this.showSearch) {
				this.showSearch = true
				this.$nextTick(() => {
					this.searchFocus = true
				})
			}
		},
		closeSearch() {
			if (this.showSearch) {
				this.showSearch = false
				this.clearSearch()
			}
		},
		handleSearch(value) {
			if (value && value.trim()) {
				this.refreshList()
			}
		},
		handleSearchAction() {
			if (this.searchKeyword && this.searchKeyword.trim()) {
				this.refreshList()
			}
		},
		handleSearchFocus() {
			this.searchFocus = true
		},
		handleSearchBlur() {
			// 延迟一下，避免点击清除按钮时失焦
			setTimeout(() => {
				this.searchFocus = false
				if (!this.searchKeyword || !this.searchKeyword.trim()) {
					this.showSearch = false
				}
			}, 200)
		},
		clearSearch() {
			this.searchKeyword = ''
			emit('courseSearchSync', '')
			this.refreshList()
		},
		resetFilters() {
			this.searchKeyword = ''
			emit('courseSearchSync', '')
			this.currentLevel1 = -1
			this.currentLevel2Id = 'all'
			this.currentCategory = { id: null, name: '全部', isAll: true }
			this.courseType = 0
			this.sortType = 'hot'
			this.showCategoryPopup = false
			this.showFilterPopup = false
			this.refreshList()
		},
		getCourseCategories() {
			this.restoreCategoryCache()
			this.currentLevel1 = -1;
			this.currentLevel2Id = 'all';
			this.currentCategory = {id: null, name: '全部', isAll: true};
			this.refreshList(!this.courseList.length);
			let cachedAt = 0
			try {
				cachedAt = Number(getStorage('qixue_course_cates_at') || 0)
			} catch (e) {}
			if (this.categoryTree.length && Date.now() - cachedAt < 30 * 60 * 1000) return
			this.$http.get('/course/getCourseCategories')
			.then(res => {
				this.categoryTree = res.data || [];
				try {
					setStorage('qixue_course_cates', this.categoryTree)
					setStorage('qixue_course_cates_at', Date.now())
				} catch (e) {}
			}).catch(err => {
				console.log(err);
			});
		},
		restoreCourseListCache() {
			try {
				const cached = getStorage('qixue_course_list')
				if (cached && cached.length) {
					this.courseList = cached
				}
			} catch (e) {}
		},
		restoreCategoryCache() {
			try {
				const cached = getStorage('qixue_course_cates')
				if (cached && cached.length) {
					this.categoryTree = cached
					return
				}
				const homeBoot = getStorage('qixue_home_boot')
				if (homeBoot && homeBoot.categories && homeBoot.categories.length) {
					this.categoryTree = homeBoot.categories
				}
			} catch (e) {}
		},
		handleCategoryChange(cateId) {
			const level1CategoryIndex = this.categoryTree.findIndex(item => item.id === cateId);
			if (level1CategoryIndex !== -1) {
				this.currentLevel1 = level1CategoryIndex;
				this.currentLevel2Id = 'all';
				const level1Category = this.categoryTree[level1CategoryIndex];
				if (level1Category) {
					this.currentCategory = {id: level1Category.id, name: level1Category.name, isAll: true};
				}
				this.refreshList();
			}
		},
		toggleCategoryPopup() {
			this.showCategoryPopup = !this.showCategoryPopup;
			if (this.showCategoryPopup) {
				this.showFilterPopup = false;
			}
		},
		toggleFilterPopup() {
			this.showFilterPopup = !this.showFilterPopup;
			if (this.showFilterPopup) {
				this.showCategoryPopup = false;
			}
		},
		selectCourseType(type) {
			this.courseType = type;
			this.showFilterPopup = false;
			this.refreshList();
		},
		selectLevel1(index) {
			if (index === -1) {
				this.currentLevel1 = -1;
				this.currentLevel2Id = 'all';
				this.currentCategory = {id: null, name: '全部', isAll: true};
				this.showCategoryPopup = false;
				this.refreshList();
				return;
			}
			if (index < 0 || index >= this.categoryTree.length) return;
			this.currentLevel1 = index;
			this.currentLevel2Id = null;
			this.currentCategory = null;
		},
		selectLevel2(item) {
			if (!item) return;
			const currentLevel1Category = this.currentLevel1 >= 0 ? this.categoryTree[this.currentLevel1] : null;

			if (item.isAll || item.id === 'all') {
				if (this.currentLevel1 === -1) {
					this.currentCategory = {id: null, name: '全部', isAll: true};
				} else if (currentLevel1Category) {
					this.currentCategory = {id: currentLevel1Category.id, name: currentLevel1Category.name, isAll: true};
				}
				this.currentLevel2Id = 'all';
			} else {
				this.currentCategory = item;
				this.currentLevel2Id = item.id;
			}
			this.showCategoryPopup = false;
			this.refreshList();
		},
		setSortType(type) {
			this.sortType = type;
			this.refreshList();
		},
		refreshList(clear = true) {
			this.page = 1;
			this.noMore = false;
			if (clear) {
				this.courseList = [];
			}
			this.getCourseList();
		},
		getCourseList() {
			const reqId = ++this.listReqId;
			this.loading = true;

			const params = {
				page: this.page,
				limit: this.isDesktop ? 12 : 10,
				sortType: this.sortType
			};

			if (this.currentLevel1 === -1) {
				// 选中一级分类的"全部"，不传分类参数或传null
				params.cateId1 = null;
				params.cateId2 = null;
			} else {
				const level1Category = this.categoryTree[this.currentLevel1];
				if (level1Category) {
					params.cateId1 = level1Category.id;

					if (this.currentLevel2Id === 'all' || this.currentLevel2Id === null) {
						params.cateId2 = null;
					} else if (this.currentCategory && !this.currentCategory.isAll) {
						params.cateId2 = this.currentCategory.id;
					} else {
						params.cateId2 = null;
					}
				}
			}

			// 设置课程类型参数 type（0表示全部，不传参数）
			if (this.courseType !== 0) {
				params.type = this.courseType;
			}

			if (this.searchKeyword && this.searchKeyword.trim()) {
				params.title = this.searchKeyword.trim();
			}

			this.$http.post('/course/getCourseList', params).then(res => {
				if (reqId !== this.listReqId) return;
				const newData = res.data.list || res.data || [];
				if (this.page === 1) {
					this.courseList = newData;
					if (this.currentLevel1 === -1 && this.courseType === 0 && this.sortType === 'hot' && !this.searchKeyword) {
						try {
							setStorage('qixue_course_list', newData)
						} catch (e) {}
					}
				} else {
					this.courseList = this.courseList.concat(newData);
				}
				this.noMore = res.data.hasMore === false;
				this.loading = false;
			}).catch(err => {
				if (reqId !== this.listReqId) return;
				console.log(err);
				this.loading = false;
				this.msg('获取课程列表失败');
			});
		},
		onCourseScroll(e) {
			const el = e.target
			if (!el || this.loading || this.noMore) return
			if (el.scrollHeight - el.scrollTop - el.clientHeight < 80) this.loadMore()
		},
		loadMore() {
			if (this.loading || this.noMore) return;
			this.page++;
			this.getCourseList();
		},
		goDetail(course) {
			navigateTo(`/pages/course/courseDetail?id=${course.id}`);
		}
	}
}
</script>

<style lang="scss">
.qx-route-pages-index-course {
	background: #fff url('/static/bg/com.png') no-repeat 0 0 / contain fixed;
}
.course-page {
	display: flex;
	flex-direction: column;
	width: 100%;
	height: calc(100vh - var(--window-bottom));
}

.navbar-placeholder {
	width: 100%;
	flex-shrink: 0;
}

.custom-navbar {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	z-index: 999;
	background: transparent;
}

.navbar-status-bar {
	width: 100%;
	background: transparent;
}

.navbar-content {
	display: flex;
	align-items: center;
	height: calc(88 * var(--rpx));
	padding: 0 calc(32 * var(--rpx));
	background: transparent;
}

.navbar-left {
	flex-shrink: 0;
	margin-right: calc(24 * var(--rpx));
}

.navbar-left-text {
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #333;
}

.navbar-center {
	flex: 1;
	display: flex;
	align-items: center;
	justify-content: center;
	overflow: hidden;
}

.navbar-title {
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #333;
}

.navbar-search {
	width: 100%;
	animation: slideIn 0.3s ease-in-out;
}

.navbar-right {
	flex-shrink: 0;
	margin-left: calc(24 * var(--rpx));
	padding: calc(8 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
}

@keyframes slideIn {
	from {
		opacity: 0;
		transform: translateX(calc(-20 * var(--rpx)));
	}
	to {
		opacity: 1;
		transform: translateX(0);
	}
}

.filter-bar {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: calc(24 * var(--rpx)) calc(32 * var(--rpx));
}

.filter-item {
	display: flex;
	align-items: center;
	margin-right: calc(48 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #666;

	&.active {
		color: #4285f4;
	}

	.arrow {
		margin-left: calc(8 * var(--rpx));
		font-size: calc(calc(20 * var(--rpx)) * var(--font-scale, 1));
		transition: transform 0.3s;

		&.up {
			transform: rotate(0deg);
		}
	}

	.filter-icon {
		margin-left: calc(8 * var(--rpx));
		font-size: calc(calc(20 * var(--rpx)) * var(--font-scale, 1));
		transition: transform 0.3s;

		&.up {
			transform: rotate(180deg);
		}
	}
}

.filter-popup {
	position: fixed;
	top: calc(var(--status-bar-height) + calc(160 * var(--rpx)));
	left: 0;
	right: 0;
	bottom: 0;
	z-index: 999;
}

.filter-mask {
	position: absolute;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(0, 0, 0, 0.5);
}

.filter-content {
	background: #fff;
	position: relative;
	z-index: 1;
	padding: calc(24 * var(--rpx));
}

.filter-options {
	display: flex;
	flex-direction: column;
}

.filter-option-item {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: calc(24 * var(--rpx)) 0;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #333;
	border-bottom: calc(1 * var(--rpx)) solid #f0f0f0;

	&:last-child {
		border-bottom: none;
	}

	&.active {
		color: #4285f4;
	}

	.check {
		color: #4285f4;
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	}
}

.category-popup {
	position: fixed;
	top: calc(var(--status-bar-height) + calc(160 * var(--rpx)));
	left: 0;
	right: 0;
	bottom: 0;
	z-index: 999;
}

.category-mask {
	position: absolute;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(0, 0, 0, 0.5);
}

.category-content {
	background: #fff;
	position: relative;
	z-index: 1;
}

.category-columns {
	display: flex;
	height: calc(600 * var(--rpx));
}

.category-column {
	overflow-y: auto;

	&.first {
		width: calc(200 * var(--rpx));
		background: #f8f9fc;
		flex-shrink: 0;
	}

	&.second {
		flex: 1;
		background: #fff;
	}
}

.category-item {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: calc(28 * var(--rpx)) calc(24 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #333;

	&.active {
		color: #4285f4;
		background: rgba(66, 133, 244, 0.05);
	}

	.check {
		color: #4285f4;
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	}
}

.course-page .course-list {
	flex: 1;
	overflow-y: auto;
}

.course-page .course-item {
	display: flex;
	padding: calc(12 * var(--rpx)) calc(20 * var(--rpx));
}

.course-page .course-cover {
	position: relative;
	width: calc(180 * var(--rpx));
		height: calc(270 * var(--rpx));
		border-radius: calc(12 * var(--rpx));
		overflow: hidden;
		flex-shrink: 0;
}

.course-page .cover-img {
	width: 100%;
	height: 100%;
	border-radius: calc(12 * var(--rpx));
		overflow: hidden;
}

.course-page .course-tag {
	position: absolute;
	top: calc(4 * var(--rpx));
	left: calc(4 * var(--rpx));
	padding: calc(6 * var(--rpx)) calc(12 * var(--rpx));
	background: rgba(0, 0, 0, 0.6);
	color: #fff;
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	border-radius: calc(6 * var(--rpx));
}

.course-page .course-info {
	flex: 1;
	min-width: 0;
	padding: 0 calc(20 * var(--rpx));
	display: flex;
	flex-direction: column;
	justify-content: space-between;
}

.course-page .course-title {
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 500;
		color: #1a1a1a;
		line-height: 1.4;
		overflow: hidden;
		display: -webkit-box;
		-webkit-line-clamp: 4;
		line-clamp: 4;
		-webkit-box-orient: vertical;
}

.course-page .course-author {
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

.course-page .course-source {
	display: flex;
	align-items: center;
	margin-top: calc(16 * var(--rpx));

	.source-text {
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #999;
	}

	.source-label {
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #999;
		margin-left: calc(8 * var(--rpx));
	}
}

.course-page .course-empty {
	padding: calc(80 * var(--rpx)) 0 calc(40 * var(--rpx));
}

.course-page .course-skel-item {
	display: flex;
	padding: calc(12 * var(--rpx)) calc(20 * var(--rpx));
}

.course-page .course-skel-cover,
.course-page .course-skel-line {
	background: #eef0f4;
	animation: course-skel-pulse 1.2s ease-in-out infinite;
}

.course-page .course-skel-cover {
	width: calc(180 * var(--rpx));
	height: calc(270 * var(--rpx));
	border-radius: calc(12 * var(--rpx));
	flex-shrink: 0;
}

.course-page .course-skel-line {
	flex: 1;
	height: calc(28 * var(--rpx));
	margin: calc(16 * var(--rpx)) calc(20 * var(--rpx)) 0;
	border-radius: calc(6 * var(--rpx));
}

.course-page .course-skel-line--short {
	width: 40%;
	height: calc(22 * var(--rpx));
}

@keyframes course-skel-pulse {
	0%, 100% { opacity: 1; }
	50% { opacity: 0.55; }
}

.course-page .loading-more,
.course-page .no-more {
	text-align: center;
	padding: calc(32 * var(--rpx));
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: #999;
}
</style>
