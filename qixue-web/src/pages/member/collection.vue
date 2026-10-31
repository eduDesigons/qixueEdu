<template>
	<div class="collection-page">

		<div class="collection-tabs">
			<van-tabs v-model:active="activeTab" color="#4285f4" line-width="22" line-height="3" @change="(i) => switchTab({ index: i })">
				<van-tab v-for="(t, i) in tabsList" :key="i" :title="t.name" :name="i" />
			</van-tabs>
		</div>

		<mescroll-body
			v-if="activeTab === 0"
			ref="courseMescroll"
			@init="mescrollInit"
			@down="downCallback"
			@up="upCallback"
			:down="downOption"
			:up="upOption"
		>
			<div class="collection-list">
				<div
					class="collection-item course-item"
					v-for="(item, index) in courseList"
					:key="item.id || index"
					@click="goToCourseDetail(item)"
				>
					<div class="course-cover">
						<img :src="static((item.course && item.course.thumb) || '')" class="cover-img" loading="lazy"  style="object-fit:cover" />
					</div>
					<div class="course-info">
						<div class="course-title">{{ (item.course && item.course.title) || '课程已下架' }}</div>
						<div class="course-desc">
							{{ (item.course && item.course.desc) || '' }}
						</div>
						<div class="course-meta">
							<span class="course-credit">{{ formatDate(item.created_at, 'YYYY-MM-DD') }}</span>
						</div>
					</div>
				</div>
			</div>
		</mescroll-body>

		<mescroll-body
			v-if="activeTab === 1"
			ref="chapterMescroll"
			@init="mescrollInit"
			@down="downCallback"
			@up="upCallback"
			:down="downOption"
			:up="upOption"
		>
			<div class="collection-list">
				<div
					class="collection-item chapter-item"
					v-for="(item, index) in chapterList"
					:key="item.id || index"
					@click="goToChapterDetail(item)"
				>
					<div class="chapter-cover">
						<img :src="static(item.course && item.course.thumb || item.cover || '/static/img/video-img-1.png')" class="cover-img" loading="lazy"  style="object-fit:cover" />
					</div>
					<div class="chapter-info">
						<div class="chapter-title">{{ item.chapter.title || '' }}</div>
						<div class="chapter-desc">
							<span class="chapter-course" v-if="item.course">{{ item.course.title || '' }}</span>
						</div>
					</div>
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
			activeTab: 0,
			tabsList: [
				{ name: '课程' },
				{ name: '章节' }
			],
			courseList: [],
			chapterList: [],
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
					tip: '暂无收藏数据'
				}
			}
		}
	},
	onLoad() {
		this.restoreCollectionCache();
	},
	methods: {
		restoreCollectionCache() {
			try {
				const cached = getStorage('qixue_my_like_course');
				if (cached && cached.length) {
					this.courseList = cached;
				}
			} catch (e) {}
		},
		switchTab(e) {
			this.activeTab = e.index;
			this.$nextTick(() => {
				if (this.mescroll) {
					this.mescroll.resetUpScroll();
				}
			});
		},
		downCallback() {
			this.mescroll.resetUpScroll();
		},
		upCallback(page) {
			if (this.activeTab === 0) {
				this.getCourseList(page);
			} else {
				this.getChapterList(page);
			}
		},
		getCourseList(page) {
			this.$http.get('/user/myLikeCourse', {
				page: page.num,
				limit: page.size
			}).then(res => {
				if (res.code === 200) {
					let list = res.data.list || []

					if (page.num === 1) {
						this.courseList = list;
						try {
							setStorage('qixue_my_like_course', list);
						} catch (e) {}
					} else {
						this.courseList = this.courseList.concat(list);
					}
					this.mescroll.endSuccess(list.length, list.length >= page.size);
				} else {
					this.msg(res.msg || '获取课程列表失败');
					this.mescroll.endErr();
				}
			}).catch(err => {
				console.log('获取课程列表失败:', err);
				this.msg('获取课程列表失败');
				this.mescroll.endErr();
			});
		},
		getChapterList(page) {
			this.$http.get('/user/myLikeChapter', {
				page: page.num,
				limit: page.size
			}).then(res => {
				if (res.code === 200) {
					let list =res.data.list || []

					if (page.num === 1) {
						this.chapterList = [];
					}
					this.chapterList = this.chapterList.concat(list);

					this.mescroll.endSuccess(list.length, list.length >= page.size);
				} else {
					this.msg(res.msg || '获取章节列表失败');
					this.mescroll.endErr();
				}
			}).catch(err => {
				console.log('获取章节列表失败:', err);
				this.msg('获取章节列表失败');
				this.mescroll.endErr();
			});
		},
		goToCourseDetail(item) {
			const courseId = item.course_id || (item.course && item.course.id);
			if (!courseId) {
				this.msg('课程已下架');
				return;
			}
			this.navTo(`/pages/course/courseDetail?id=${courseId}`);
		},
		goToChapterDetail(item) {
			const chapterId = item.chapter_id || (item.chapter && item.chapter.id);
			if (!chapterId) {
				this.msg('章节已下架');
				return;
			}
			const courseId = item.course_id || (item.course && item.course.id) || (item.chapter && item.chapter.course_id);
			const url = courseId
				? `/pages/course/chapterDetail?id=${chapterId}&courseId=${courseId}`
				: `/pages/course/chapterDetail?id=${chapterId}`;
			this.navTo(url);
		},
		showDuration(seconds) {
			if (!seconds) return '00:00';
			const h = Math.floor(seconds / 3600);
			const m = Math.floor((seconds % 3600) / 60);
			const s = Math.floor(seconds % 60);
			if (h > 0) {
				return `${h.toString().padStart(2, '0')}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
			}
			return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
		}
	}
}
</script>

<style lang="scss" scoped>
.collection-page {
	min-height: 100vh;
	background: #f6f7fb;
}

.collection-tabs {
	background: #fff;
	position: sticky;
	top: 0;
	z-index: 10;
}

.collection-list {
	padding: calc(20 * var(--rpx)) calc(30 * var(--rpx));
}

.collection-item {
	background: #fff;
	border-radius: calc(16 * var(--rpx));
	margin-bottom: calc(20 * var(--rpx));
	overflow: hidden;
}

.course-item {
	display: flex;
	padding: calc(24 * var(--rpx));
	align-items: center;
}

.course-cover {
	position: relative;
	width: calc(180 * var(--rpx));
		height: calc(270 * var(--rpx));
		border-radius: calc(12 * var(--rpx));
	overflow: hidden;
	flex-shrink: 0;
}

.cover-img {
	width: 100%;
	height: 100%;
}

.course-tag {
	position: absolute;
	top: calc(8 * var(--rpx));
	left: calc(8 * var(--rpx));
	padding: calc(4 * var(--rpx)) calc(12 * var(--rpx));
	background: rgba(0, 0, 0, 0.6);
	color: #fff;
	font-size: calc(calc(20 * var(--rpx)) * var(--font-scale, 1));
	border-radius: calc(6 * var(--rpx));
}

.course-info {
	flex: 1;
	min-width: 0;
	height: 100%;
	margin-left: calc(24 * var(--rpx));
	display: flex;
	flex-direction: column;
	justify-content: space-between;
}

.course-title {
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
	line-height: 1.4;
	overflow: hidden;
	display: -webkit-box;
	-webkit-line-clamp: 4;
	line-clamp: 4;
	-webkit-box-orient: vertical;
	margin-bottom: calc(12 * var(--rpx));
}

.course-desc {
	display: -webkit-box;
	-webkit-box-orient: vertical;
	-webkit-line-clamp: 4;
	line-clamp: 4;
	overflow: hidden;
	word-break: break-word;
	overflow-wrap: anywhere;
	white-space: normal;
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #9ca3af;
	margin-bottom: calc(12 * var(--rpx));
}

.course-lesson {
	margin-right: calc(8 * var(--rpx));
}

.course-author {
	color: #7f848c;
}

.course-meta {
	display: flex;
	align-items: center;
}

.course-credit {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #999;
}

.chapter-item {
	display: flex;
	padding: calc(24 * var(--rpx));
	align-items: center;
}

.chapter-cover {
	position: relative;
	width: calc(120 * var(--rpx));
	height: calc(120 * var(--rpx));
	border-radius: calc(12 * var(--rpx));
	overflow: hidden;
	flex-shrink: 0;
	display: flex;
	align-items: center;
	justify-content: center;
	background: #f5f5f5;
}

.chapter-type-icon {
	position: absolute;
	top: 50%;
	left: 50%;
	transform: translate(-50%, -50%);
	z-index: 2;
}

.chapter-info {
	flex: 1;
	margin-left: calc(24 * var(--rpx));
	display: flex;
	flex-direction: column;
	justify-content: space-between;
	min-height: calc(120 * var(--rpx));
}

.chapter-title {
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
	line-height: 1.4;
	overflow: hidden;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	line-clamp: 2;
	-webkit-box-orient: vertical;
	margin-bottom: calc(12 * var(--rpx));
}

.chapter-desc {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
	margin-bottom: calc(12 * var(--rpx));
}

.chapter-course {
	color: #7f848c;
}

.chapter-meta {
	display: flex;
	align-items: center;
}

.chapter-duration {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #999;
}
</style>