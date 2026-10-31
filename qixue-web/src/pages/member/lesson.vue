<template>
	<div>
		<div class="lesson-wap">
			<div class="lesson-wap-page">
				<div class="filter-bar">
					<div class="filter-warp" style="height: 44px;">
						<div class="filter-nav" style="height: 44px;">
							<div class="filter-nav-item" :class="{ active: activeType === 1 }" @click="switchType(1)">
								<span>最近学习</span>
							</div>
							<div class="filter-nav-item" :class="{ active: activeType === 2 }" @click="switchType(2)">
								<span>最近加入</span>
							</div>
							<div class="filter-nav-item" :class="{ active: activeType === 3 }" @click="switchType(3)">
								<span>已学完</span>
							</div>
						</div>
					</div>
				</div>
				<mescroll-body
					ref="mescrollRef"
					@init="mescrollInit"
					@down="downCallback"
					@up="upCallback"
					:down="downOption"
					:up="upOption"
				>
				<div id="lesson-list-wap" class="lesson-list-wap">
					<div class="course-item-wrap" v-for="(item, index) in courseList" :key="item.id || index">
						<div class="lesson-item-wap lesson-item">
							<div class="lesson-item-detail">
								<div class="lesson-image" @click="goToCourseDetail(item)">
									<img :src="static((item.course && item.course.thumb) || '')" loading="lazy"  style="object-fit:cover" />
								</div>
								<div class="lesson-item-disc">
									<div class="lesson-title" @click="goToCourseDetail(item)">{{ item.course && item.course.title || '' }}</div>
									<div class="lesson-hour">
										<span v-if="item.course">共{{ item.course.chapter_count || 0 }}讲 | 已学{{ item.has_end_chapter || 0 }}讲</span>
									</div>
									<div class="lesson-progress">
										<div class="lesson-progress-bar">
											<div com-progress="" class="com-progress">
												<div progress-wrap="" border="true" class="progress-wrap"
													style="height: 8px; border-radius: 4px; background-color: rgb(255, 255, 255); border-width: 2px;">
													<div progress="" class="progress"
														:style="`width: ${item.progress || 0}%; height: 4px; border-radius: 4px; background-color: rgba(66, 133, 244, 0.76);`">
													</div>
												</div>
												<div progress-desc="" class="progress-value" :data-ok="item.progress >= 100">{{ item.progress || 0 }}%</div>
											</div>
										</div>
										<div class="lesson-continue" @click="goContinue(item)">{{ continueBtn(item) }}</div>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
				</mescroll-body>
			</div>
		</div>
	</div>
</template>

<script>
import { getStorage } from '@/common/storage'
	import MescrollMixin from "@/uni_modules/mescroll-uni/components/mescroll-uni/mescroll-mixins.js";

	export default {
		mixins: [MescrollMixin],
		data() {
			return {
				activeType: 1,
				courseList: [],
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
						tip: '暂无课程数据',
						icon: '/static/empty/course.png'
					}
				}
			};
		},
		onLoad() {
			this.restoreRecentCache();
		},
		methods: {
			restoreRecentCache() {
				if (this.activeType !== 1) return;
				try {
					const cached = getStorage('qixue_learn_cache');
					const list = cached && cached.myCourse;
					if (!list || !list.length) return;
					this.courseList = list.map((item) => {
						const count = (item.course && item.course.chapter_count) || 0;
						return {
							...item,
							progress: count > 0 ? Math.round((item.has_end_chapter || 0) / count * 100) : 0
						};
					});
				} catch (e) {}
			},
			switchType(type) {
				if (this.activeType === type) return;
				this.activeType = type;
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
				this.getCourseList(page);
			},
			getCourseList(page) {
				this.$http.get('/learn/myCourse', {
					type: this.activeType,
					page: page.num,
					limit: page.size
				}).then(res => {
					if (res.code === 200) {
						let list = res.data.list || [];

						list.forEach(item => {
							if (item.course && item.course.chapter_count > 0) {
								item.progress = Math.round((item.has_end_chapter || 0) / item.course.chapter_count * 100);
							} else {
								item.progress = 0;
							}
						});

						if (page.num === 1) {
							this.courseList = [];
						}
						this.courseList = this.courseList.concat(list);

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
			goToCourseDetail(item) {
				const courseId = item.course_id || (item.course && item.course.id);
				if (!item.course) {
					this.msg('课程已下架');
					return;
				}
				if (courseId) {
					this.navTo(`/pages/course/courseDetail?id=${courseId}`);
				}
			},
			continueBtn(item) {
				if (Number(item.is_end) === 1) return '查看';
				if (Number(item.has_end_chapter) > 0) return '续学';
				return '开始';
			},
			goContinue(item) {
				const courseId = item.course_id || (item.course && item.course.id);
				const ch = item.continueChapter;
				if (ch && ch.id && Number(item.is_end) !== 1) {
					this.navTo(`/pages/course/chapterDetail?id=${ch.id}&courseId=${courseId}`);
					return;
				}
				this.goToCourseDetail(item);
			},
			formatTime(seconds) {
				if (!seconds) return '';
				const mins = Math.floor(seconds / 60);
				const secs = seconds % 60;
				return `${mins}:${secs.toString().padStart(2, '0')}`;
			}
		}
	}
</script>

<style lang="less">

	.paging {
		display: flex;
		align-items: center;
		width: 100%;
	}

	.loading {
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.progress-sort {
		transition: all .3s ease;
	}

	.lesson-wap .filter-nav-item>span {
		overflow: hidden;
		word-break: keep-all;
		white-space: nowrap;
		text-overflow: ellipsis;
	}

	.lesson-wap .filter-bar {
		position: relative;
		margin-bottom: -4px;
		background: #fff;
	}

	.lesson-wap .filter-warp {
		position: relative;
		padding: 4px 15px 0 15px;
		justify-content: space-between;
		z-index: 30;
	}

	.lesson-wap .filter-nav,
	.lesson-wap .filter-warp {
		display: flex;
		align-items: center;
		background: #fff;
	}

	.lesson-wap .filter-nav {
		gap: 16px;
	}

	.lesson-wap .filter-nav-item {
		display: flex;
		align-items: center;
		font-size: calc(14px * var(--font-scale, 1));
		font-weight: 500;
		line-height: 21px;
		padding: 0 calc(10 * var(--rpx));
		color: #404040;
	}

	.lesson-wap .filter-nav-item,
	.lesson-wap .filter-nav-item:nth-child(2) {
		justify-content: center;
	}

	.lesson-wap .filter-nav-item.active {
		color: #4285f4;
	}

	.progress-sort {
		position: relative;
		display: flex;
		align-items: center;
		line-height: 21px;
		font-size: calc(14px * var(--font-scale, 1));
		font-weight: 500;
		color: #4c4c4c;
		cursor: pointer;
	}

	.progress-sort>span {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		line-height: 1;
	}

	.progress-sort>span .ascend {
		line-height: 1;
		color: #7f848c;
		transform: rotate(180deg) scale(.5);
	}

	.progress-sort>span .descend {
		margin-top: -11px;
		line-height: 1;
		color: #7f848c;
		transform: scale(.5);
	}

	.filter-actions {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 17px;
	}

	.filter-actions-wrap {
		position: relative;
	}

	.filter-actions-wrap:before {
		content: "";
		position: absolute;
		top: 0;
		right: -8px;
		width: 1px;
		height: 14px;
		background: #d4d6da;
	}

	.filter-actions-wrap:last-child:before {
		display: none;
	}

	.filter-actions-item {
		position: relative;
		display: flex;
		align-items: center;
		font-size: calc(14px * var(--font-scale, 1));
		font-weight: 500;
		color: #4c4c4c;
		line-height: 1;
		cursor: pointer;
		-webkit-user-select: none;
		-moz-user-select: none;
		-ms-user-select: none;
		user-select: none;
	}

	.filter-actions-item i {
		margin-left: 2px;
		font-size: calc(12px * var(--font-scale, 1));
		color: #7f848c;
	}

	.NativeTopBar_title_CWhmq {
		transition: all .3s ease;
		overflow: hidden;
		word-break: keep-all;
		white-space: nowrap;
		text-overflow: ellipsis;
	}

	.NativeTopBar_placelholder_1riAA,
	.NativeTopBar_shim_2poDt {
		height: 44px;
	}

	.NativeTopBar_shim_2poDt {
		position: fixed;
		left: 0;
		right: 0;
		z-index: 20000;
	}

	.NativeTopBar_header_2ZNBh {
		position: relative;
		display: flex;
		height: 44px;
	}

	.NativeTopBar_title_CWhmq {
		box-sizing: border-box;
		width: 100%;
		padding: 0 80px;
		color: #353535;
		text-align: center;
		font-size: calc(19px * var(--font-scale, 1));
		font-weight: 500;
		line-height: 44px;
	}

	.NativeTopBar_end-container_2EY3l,
	.NativeTopBar_pre-container_2Hlp5 {
		position: absolute;
		top: 0;
		z-index: 10;
		display: flex;
		align-items: center;
		height: 44px;
		color: inherit;
	}

	.NativeTopBar_pre-container_2Hlp5 {
		left: 0;
		padding-left: 14px;
	}

	.NativeTopBar_end-container_2EY3l {
		right: 15px;
	}

	.NativeTopBar_back_k-qKP {
		margin-left: -5px;
		font-size: calc(22px * var(--font-scale, 1));
		color: #b9bdc2;
		line-height: 1;
		color: #fff;
	}

	.com-progress {
		display: flex;
		align-items: center;
	}

	.progress-wrap {
		box-sizing: border-box;
		position: relative;
		flex-grow: 1;
		flex-shrink: 0;
		overflow: hidden;
	}

	.progress-wrap[border=true] {
		border: 1px solid #f1f2f3;
	}

	.progress {
		flex-shrink: 0;
		position: absolute;
		left: -1px;
		top: 0;
		bottom: 0;
		transition: width .3s ease-in;
	}

	.progress-value {
		margin-left: 6px;
		flex-shrink: 0;
		line-height: 1;
		font-size: calc(12px * var(--font-scale, 1));
		font-weight: 500;
		color: #f36749;
	}

	.progress-value[data-ok=true] {
		color: #7f848c;
		font-weight: 400;
	}

	.lesson-hour,
	.lesson-title {
		overflow: hidden;
		word-break: keep-all;
		white-space: nowrap;
		text-overflow: ellipsis;
	}

	.lesson-item-disc {
		display: flex;
		flex-direction: column;
	}

	.lesson-item-wap {
		box-sizing: border-box;
		padding: 16px 0 17px 15px;
		background-color: #fff;
	}

	.lesson-item-detail {
		display: flex;
		padding-right: 15px;
	}

	.lesson-image {
		position: relative;
		flex-shrink: 0;
		width: 120px;
		height: 90px;
		margin-right: 13px;
		background-color: #f6f7fb;
		border-radius: 4px;
	}

	.lesson-image>img {
		width: 100%;
		height: 100%;
		border-radius: 4px;
		-o-object-fit: cover;
		object-fit: cover;
	}

	.lesson-time-box {
		position: absolute;
		bottom: 0;
		right: 0;
		width: 120px;
		height: 18px;
		background: linear-gradient(180deg, transparent, rgba(0, 0, 0, .5));
		border-radius: 0 0 4px 4px;
	}

	.lesson-total-time {
		position: absolute;
		bottom: 1px;
		right: 5px;
		line-height: 16px;
		font-size: calc(12px * var(--font-scale, 1));
		font-weight: 400;
		color: #fff;
		opacity: 1;
	}

	.lesson-item-disc {
		flex: 1;
		min-width: 0;
		padding-bottom: 4px;
		margin-top: 4px;
	}

	.lesson-title {
		width: 100%;
		font-size: calc(15px * var(--font-scale, 1));
		font-weight: 500;
		color: #404040;
		line-height: 20px;
	}

	.lesson-hour {
		margin-top: 4px;
		line-height: 16px;
		font-size: calc(12px * var(--font-scale, 1));
		font-weight: 400;
		color: #7f848c;
	}

	.lesson-hour>span {
		margin-right: 7px;
		font-feature-settings: "tnum";
	}

	.lesson-progress {
		flex: 1;
		display: flex;
		justify-content: space-between;
		align-items: flex-end;
	}

	.lesson-progress-bar {
		flex: 1;
		display: flex;
		max-width: 97px;
		padding-top: 2px;
	}

	.lesson-progress-bar [com-progress] {
		flex: 1;
		margin-right: 16px;
	}

	.lesson-progress-bar [progress-wrap] {
		flex: 1;
		height: 8px;
		border-color: rgba(185, 189, 194, .2);
	}

	.lesson-progress-bar [progress] {
		background: linear-gradient(90deg, #4f9cfa, #4285f4);
	}

	.lesson-progress-bar [progress-desc] {
		width: 31px;
	}

	.lesson-continue {
		flex-shrink: 0;
		margin-left: 8px;
		padding: 0 12px;
		height: 26px;
		line-height: 26px;
		border-radius: 13px;
		background: #4285f4;
		color: #fff;
		font-size: calc(12px * var(--font-scale, 1));
	}

	.lesson-extime,
	.lesson-progress-bar [progress-desc] {
		flex: none;
		font-size: calc(12px * var(--font-scale, 1));
		font-weight: 400;
		line-height: 16px;
		color: #7f848c;
	}

	.lesson-extime {
		margin-top: 2px;
		white-space: nowrap;
	}

	.lesson-list-loading,
	.paging-loading {
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.lesson-wap {
		overflow: hidden;
		height: 100%;
	}

	.lesson-wap-page {
		overflow: hidden;
		height: 100%;
	}

	.lesson-list-wap {
		overflow: auto;
		height: calc(var(--vh)*100 - 82px);
	}

	.lesson-list-loading {
		position: fixed;
		left: 0;
		top: 100px;
		width: 100%;
		height: 500px;
	}

	.paging-loading {
		width: 100%;
	}

	.lesson-hidden-icon {
		font-size: calc(21px * var(--font-scale, 1));
		color: #7f848c;
	}

	.is-end {
		width: 100%;
		padding: 30px 0 100px;
		background: #fff;
		line-height: 20px;
		text-align: center;
		font-size: calc(14px * var(--font-scale, 1));
		font-weight: 500;
		color: #b9bdc2;
	}

	.course-item-wrap {
		position: relative;
		display: flex;
		align-items: center;
	}

	.course-item-wrap:before {
		position: absolute;
		left: 15px;
		right: 0;
		bottom: 0;
		z-index: 1;
		content: "";
		height: 1px;
		background: #e8e8e8;
	}

	.lesson-item {
		flex: 1;
		max-width: 100%;
	}

	@-webkit-keyframes gk-rotate {
		0% {
			transform: rotate(0deg);
		}

		to {
			transform: rotate(1turn);
		}
	}

	@keyframes gk-rotate {
		0% {
			transform: rotate(0deg);
		}

		to {
			transform: rotate(1turn);
		}
	}

	@-webkit-keyframes gk-circular {
		0% {
			stroke-dasharray: 1, 200;
			stroke-dashoffset: 0;
		}

		50% {
			stroke-dasharray: 90, 150;
			stroke-dashoffset: -40;
		}

		to {
			stroke-dasharray: 90, 150;
			stroke-dashoffset: -120;
		}
	}

	@keyframes gk-circular {
		0% {
			stroke-dasharray: 1, 200;
			stroke-dashoffset: 0;
		}

		50% {
			stroke-dasharray: 90, 150;
			stroke-dashoffset: -40;
		}

		to {
			stroke-dasharray: 90, 150;
			stroke-dashoffset: -120;
		}
	}
</style>