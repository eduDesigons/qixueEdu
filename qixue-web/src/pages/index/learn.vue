<template>
	<LearnDesktop
		v-if="isDesktop"
		:study-data="studyData"
		:courses="myCourse"
		:tasks="myTask"
		:ready="hasInit"
		:rank-period="rankPeriod"
		:rank-loading="rankLoading"
		:has-rank-data="hasRankData"
		:duration-list="rankList('duration')"
		:course-list="rankList('course')"
		:coin-list="rankList('coins')"
		@continue="goContinue"
		@tip="showRankTip"
		@history="goRankHistory"
	/>
	<div v-else class="topBg">
		<div class="study-center-wap">
			<div class="layout">
				<div class="layout-content">
					<div class="study-header">
						<div class="study-time" v-if="studyData.time > 0">
							学习时长已达&nbsp;
							<span class="study-time-primary">{{ showDuration(studyData.time) }}</span>
							&nbsp;继续努力!
						</div>
						<div class="study-content">
							<div>
								<span class="study-content-time">
									<span class="study-content-primary">{{ studyData.courseNum }}</span>
									<span>门</span>
								</span>
								<span class="study-content-desc">学完课程</span>
							</div>
							<span class="study-vertical-line"></span>
							<div>
								<span class="study-content-time">
									<span class="study-content-primary">{{ studyData.chapterNum }}</span>
									<span>讲</span>
								</span>
								<span class="study-content-desc">学完课节</span>
							</div>
							<span class="study-vertical-line"></span>
							<div>
								<span class="study-content-time">
									<span class="study-content-primary">{{ studyData.todayNum }}</span>
									<span>讲</span>
								</span>
								<span class="study-content-desc">今日完成</span>
							</div>
						</div>
					</div>
				</div>
			</div>
			<div class="study-center-wap-content">
				<div class="layout-box-wap">
					<div class="layout-box-wap-header">
						<div class="layout-box-wap-header-title">
							<div title-under="" class="">最近学习</div>
						</div>
						<div class="layout-box-wap-header-all" @click="navTo('/pages/member/lesson')">查看全部</div>
					</div>
					<div class="layout-content-box-wap">
						<div class="studying-content" v-if="!hasInit">
							<div class="studying-content-item" v-for="n in 3" :key="'sk'+n">
								<div class="studying-img-box studying-skeleton"></div>
								<div class="studying-skel-line"></div>
								<div class="studying-skel-line studying-skel-line--short"></div>
							</div>
						</div>
						<div class="studying-empty" v-else-if="myCourse.length === 0">
							<div>
								<van-empty description="没有正在学习的课程，快去学习吧~" image="/static/empty/video.png" :image-size="100" />
								<div class="studying-empty-btn" @click="navTo('/pages/index/course')">去选课</div>
							</div>
						</div>
						<div class="studying-content" v-else>
							<div class="studying-content-item" v-for="item in myCourse" :key="item.id" @click="goContinue(item)">
								<div class="studying-img-box">
									<img
										:src="static((item.course && item.course.thumb) || '')" loading="lazy"
										alt="" />
								</div>
								<span class="studying-content-title">{{ (item.course && item.course.title) || '课程已下架' }}</span>
								<span class="studying-content-sub">{{ continueLabel(item) }}</span>
							</div>
						</div>
					</div>
				</div>
				<div class="layout-box-wap" v-if="myTask.length > 0">
					<div class="layout-box-wap-header">
						<div class="layout-box-wap-header-title">
							<div title-under="" class="">待完成任务({{ myTask.length }})</div>
						</div>
					</div>
					<div class="task-content-wrapper">
						<div class="task-list">
							<div class="task-item" v-for="item in myTask" :key="item.paperId" @click="navTo('/pages/exam/exam?id='+item.paperId)">
								<div class="task-item-icon">
									<img src="/static/img/me/test.png" width="30" height="30" alt="" style="width:30px;height:30px;" />
								</div>
								<div class="task-item-content">
									<div class="task-item-title">{{ item.title }}</div>
									<div class="task-item-desc" v-if="item.desc">{{ item.desc }}</div>
								</div>
								<div class="task-item-arrow">
									<van-icon name="arrow" size="14" color="#b9bdc2" />
								</div>
							</div>
						</div>
					</div>
				</div>
				<rank-card
					title="学习时长 Top5"
					kind="user"
					v-model:period="rankPeriod.duration"
					:list="rankList('duration')"
					:loading="rankLoading && !hasRankData"
					@tip="showRankTip('duration')"
					@history="goRankHistory('duration')"
				></rank-card>
				<rank-card
					title="学习课程 Top5"
					kind="course"
					v-model:period="rankPeriod.course"
					:list="rankList('course')"
					:loading="rankLoading && !hasRankData"
					@tip="showRankTip('course')"
					@history="goRankHistory('course')"
				></rank-card>
				<rank-card
					title="积分 Top5"
					kind="user"
					v-model:period="rankPeriod.coins"
					:list="rankList('coins')"
					:loading="rankLoading && !hasRankData"
					@tip="showRankTip('coins')"
					@history="goRankHistory('coins')"
				></rank-card>
			</div>
		</div>
	</div>
		<van-dialog
			:show="!!rankTip"
			title="榜单说明"
			:message="rankTip"
			confirm-button-text="知道了"
			:show-cancel-button="false"
			@confirm="rankTip = ''"
			@update:show="(v)=>{if(!v)rankTip=''}"
		/>
</template>

<script>
import { getStorage, setStorage } from '@/common/storage'
import { defineAsyncComponent } from 'vue'
import { isDesktop as desktopOn } from '@/pc/desktop'
import RankCard from './components/RankCard.vue'

export default {
	components: {
		RankCard,
		LearnDesktop: defineAsyncComponent(() => import('@/pc/LearnDesktop.vue')),
	},
	data() {
		return {
			hasInit: false,
			studyData: {
				courseNum: 0,
				chapterNum: 0,
				todayNum: 0,
				time: 0
			},
			myCourse: [],
			myTask: [],
			learnFetchedAt: 0,
			rankLoading: false,
			rankFetchedAt: 0,
			rankTip: '',
			rankPeriod: {
				duration: 'week',
				course: 'week',
				coins: 'week'
			},
			rankData: {
				duration: { week: { list: [] }, month: { list: [] } },
				course: { week: { list: [] }, month: { list: [] } },
				coins: { week: { list: [] }, month: { list: [] } },
				tips: {
					duration: '按全部学员本周/本月累计学习时长排名，完成课节后定时更新。',
					course: '按课程本周/本月累计学习人数排名，每日 00:00 更新。',
					coins: '按学员本周/本月获得的积分排名，定时更新。'
				}
			}
		}
	},
	computed: {
		isDesktop() {
			return desktopOn.value
		},
		hasRankData() {
			return ['duration', 'course', 'coins'].some((type) => this.rankList(type).length > 0)
		}
	},
	onLoad() {
		this.restoreLearnCache()
	},
	onShow() {
		const fresh = this.learnFetchedAt > 0 && Date.now() - this.learnFetchedAt < 45000
		if (!fresh) {
			this.getStudyData()
			this.getMyCourse()
			this.getMyTask()
		}
		this.getRankBoards(false)
	},
	methods: {
		restoreLearnCache() {
			try {
				const cached = getStorage('qixue_learn_cache')
				if (cached) {
					if (cached.studyData) this.studyData = cached.studyData
					if (cached.myCourse) {
						this.myCourse = cached.myCourse
						this.hasInit = true
					}
					if (cached.myTask) this.myTask = cached.myTask
				}
				const rank = getStorage('qixue_rank_boards')
				if (rank) {
					this.rankData = Object.assign({}, this.rankData, rank)
					this.rankFetchedAt = Number(getStorage('qixue_rank_boards_at') || 0)
				}
			} catch (e) {}
		},
		saveLearnCache() {
			try {
				setStorage('qixue_learn_cache', {
					studyData: this.studyData,
					myCourse: this.myCourse,
					myTask: this.myTask
				})
			} catch (e) {}
		},
		rankList(type) {
			const period = this.rankPeriod[type] || 'week'
			const board = this.rankData[type] && this.rankData[type][period]
			return (board && board.list) || []
		},
		showRankTip(type) {
			this.rankTip = (this.rankData.tips && this.rankData.tips[type]) || ''
		},
		goRankHistory(type) {
			this.navTo('/pages/member/rankHistory?type=' + type)
		},
		getRankBoards(force = false) {
			const coinsEmpty = this.rankList('coins').length === 0
			if (!force && !coinsEmpty && this.hasRankData && Date.now() - this.rankFetchedAt < 90000) {
				return
			}
			this.rankLoading = !this.hasRankData
			this.$http.get('/rank/boards').then(res => {
				if (res.data) {
					this.rankData = Object.assign({}, this.rankData, res.data)
					this.rankFetchedAt = Date.now()
					try {
						setStorage('qixue_rank_boards', res.data)
						setStorage('qixue_rank_boards_at', this.rankFetchedAt)
					} catch (e) {}
				}
			}).catch(() => {}).finally(() => {
				this.rankLoading = false
			})
		},
		getStudyData() {
			this.$http.get('/learn/learnNum').then(res => {
				this.studyData = res.data
				this.learnFetchedAt = Date.now()
				this.saveLearnCache()
			})
		},
		getMyCourse() {
			this.$http.get('/learn/myCourse').then(res => {
				this.myCourse = res.data.list || []
				this.learnFetchedAt = Date.now()
				this.saveLearnCache()
			}).catch(() => {}).finally(() => {
				this.hasInit = true
			})
		},
		getMyTask() {
			this.$http.get('/learn/myTask').then(res => {
				this.myTask = res.data || []
				this.learnFetchedAt = Date.now()
				this.saveLearnCache()
			})
		},
		continueLabel(item) {
			const ch = item && item.continueChapter
			if (item && Number(item.is_end) === 1) {
				return '已学完，查看课程'
			}
			if (ch && ch.title) {
				return (Number(item.has_end_chapter) > 0 ? '续学 · ' : '开始 · ') + ch.title
			}
			if (item && item.has_end_chapter > 0) {
				return '已学完' + item.has_end_chapter + '讲'
			}
			return '开始学习'
		},
		goContinue(item) {
			const courseId = item && item.course && item.course.id
			const ch = item && item.continueChapter
			if (!courseId) {
				this.msg('课程已下架')
				return
			}
			if (ch && ch.id && Number(item.is_end) !== 1) {
				this.navTo('/pages/course/chapterDetail?id=' + ch.id + '&courseId=' + courseId)
				return
			}
			this.navTo('/pages/course/courseDetail?id=' + courseId)
		}
	}
}
</script>

<style>
	.qx-route-pages-index-learn {
	background: #F6F7FB url('/static/bg/com.png') no-repeat 0 0 / contain fixed;
	padding-bottom: calc(160 * var(--rpx));
}
/*! CSS Used from: Embedded */
.HeaderWap_menu-icon_3hITh {
	display: flex;
	align-items: center;
	justify-content: center;
}

.HeaderWap_common-header-wap_1SFTT {
	position: relative;
	width: 100%;
	background: #fff;
}

.HeaderWap_shim_2DQKg {
	height: 52px;
}

.HeaderWap_header_3Lxv8 {
	box-sizing: border-box;
	top: 0;
	z-index: 2003;
	display: flex;
	align-items: center;
	width: 100%;
	height: 52px;
	padding: 0 15px;
	background: #fff;
}

.HeaderWap_header_3Lxv8 {
	position: fixed;
	left: 0;
}

.HeaderWap_logo_FeLmk {
	width: -webkit-fit-content;
	width: -moz-fit-content;
	width: fit-content;
	max-width: 145px;
	height: 20px;
}

.HeaderWap_logo_FeLmk img {
	width: 100%;
	height: 100%;
}

.HeaderWap_menu-icon_3hITh {
	width: 38px;
	height: 38px;
	color: #7f848c;
}

.HeaderWap_menu-icon_3hITh [icon] {
	position: relative;
	top: 1px;
	font-size: calc(18px * var(--font-scale, 1));
	font-weight: 500;
	line-height: 1;
	color: #7f848c;
}

.HeaderWap_placeholder_3xYrQ {
	flex-grow: 1;
}

.HeaderWap_account_RhqZK {
	flex-shrink: 0;
	display: flex;
	align-items: center;
	height: 30px;
	padding: 0 13px;
	background: #f3f8ff;
	border-radius: 15px;
	font-size: calc(13px * var(--font-scale, 1));
	font-weight: 500;
	color: #4285f4;
}

.HeaderWap_account_RhqZK .HeaderWap_line_3shPN {
	margin: 0 10px;
	line-height: 1;
}

.HeaderWap_account_RhqZK [icon] {
	line-height: 1;
	font-size: calc(14px * var(--font-scale, 1));
}

/*! CSS Used from: Embedded */
.layout,
.layout-content {
	display: flex;
	flex-direction: column;
}

.layout {
	flex: 1;
}

.layout-content {
	flex: 1;
}

/*! CSS Used from: Embedded */
.layout-box-wap {
	padding-top: 18px;
	margin-top: 12px;
	border-radius: 12px;
	background-color: #fff;
}

.layout-box-wap-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 12px;
	padding: 0 15px;
	line-height: 1.45;
}

.layout-box-wap-header-title {
	position: relative;
	display: flex;
	color: #353535;
	font-size: calc(16px * var(--font-scale, 1));
	font-weight: 500;
}

.layout-box-wap-header-all {
	flex-shrink: 0;
	color: #7f848c;
	font-size: calc(13px * var(--font-scale, 1));
	font-weight: 400;
}

/* 任务列表样式 */
.task-content-wrapper {
	padding: 15px 20px;
}

.task-list {
	display: flex;
	flex-direction: column;
	gap: 12px;
}

.task-item {
	display: flex;
	align-items: center;
	padding: 16px;
	background: #f8f9fa;
	border-radius: 12px;
	transition: all 0.3s;
	position: relative;
}

.task-item:active {
	background: #e9ecef;
	transform: scale(0.98);
}

.task-item-icon {
	border-radius: 8px;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-right: 12px;
	flex-shrink: 0;
}

.task-item-content {
	flex: 1;
	display: flex;
	flex-direction: column;
	min-width: 0;
}

.task-item-title {
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #404040;
	line-height: 1.5;
	margin-bottom: 6px;
	overflow: hidden;
	text-overflow: ellipsis;
	display: -webkit-box;
	-webkit-line-clamp: 1;
	line-clamp: 1;
	-webkit-box-orient: vertical;
}

.task-item-desc {
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
	line-height: 1.5;
	overflow: hidden;
	text-overflow: ellipsis;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	line-clamp: 2;
	-webkit-box-orient: vertical;
}

.task-item-arrow {
	margin-left: 12px;
	flex-shrink: 0;
	display: flex;
	align-items: center;
}

/*! CSS Used from: Embedded */
.studying-empty {
	display: flex;
	justify-content: center;
	align-items: center;
	min-height: 147px;
	height: auto;
	width: 100%;
	padding: 24px 16px;
	box-sizing: border-box;
	text-align: center;
}

.studying-empty-desc {
	margin-top: 22px;
	color: #b9bdc2;
	font-size: calc(14px * var(--font-scale, 1));
	font-weight: 400;
}

.studying-empty-btn {
	width: 150px;
	height: 40px;
	margin: 16px auto 0;
	background: #4285f4;
	border-radius: 10px;
	line-height: 40px;
	text-align: center;
	color: #fff;
	font-size: calc(15px * var(--font-scale, 1));
}

/*! CSS Used from: Embedded */
.study-center-wap {
	min-height: 100vh;
	padding-top: 46px;
}

.study-center-wap-content {
	margin: 0 15px 30px;
}

/*! CSS Used from: Embedded */

.study-header-top {
	min-height: 44px;
	height: auto;
	padding: 8px 18px;
	line-height: 1.4;
	text-align: right;
}

.study-header-top span {
	margin-left: 20px;
	color: #9ca3ac;
	font-size: calc(21px * var(--font-scale, 1));
	font-weight: 400;
}

.study-time {
	display: flex;
	align-items: center;
	margin: 3px 0 0 15px;
	color: #404040;
	font-size: calc(14px * var(--font-scale, 1));
	font-weight: 500;
}

.study-time span {
	color: #7f848c;
	font-size: calc(12px * var(--font-scale, 1));
	font-weight: 400;
}

.study-time-primary {
	color: #4285f4;
	font-size: calc(16px * var(--font-scale, 1));
	font-weight: 500;
}

.study-content {
	display: flex;
	justify-content: space-around;
	align-items: center;
	padding: 19px 15px 9px 15px;
	text-align: center;
	color: rgba(64, 64, 64, .8);
	font-size: calc(12px * var(--font-scale, 1));
	font-weight: 400;
}

.study-content>div {
	flex: 1;
}

.study-content-time {
	display: flex;
	justify-content: center;
	align-items: flex-end;
}

.study-content-time span:last-child {
	line-height: 22px;
}

.study-content-primary {
	margin-right: 2px;
	line-height: 28px;
	font-family: DIN;
	color: #404040;
	font-size: calc(22px * var(--font-scale, 1));
	font-weight: 400;
}

.study-content-desc {
	margin-top: 2px;
	font-size: calc(13px * var(--font-scale, 1));
	font-weight: 400;
}

.study-vertical-line {
	width: 1px;
	align-self: stretch;
	height: auto;
	margin: 8px 0;
	background-color: #e9eaed;
}

/*! CSS Used from: Embedded */
.iconfont {
	font-family: iconfont !important;
	font-size: calc(16px * var(--font-scale, 1));
	font-style: normal;
	-webkit-font-smoothing: antialiased;
	-moz-osx-font-smoothing: grayscale;
}

/*! CSS Used from: Embedded */
.studying-img-box {
	overflow: hidden;
	position: relative;
	height: 77px;
	border-radius: 4px;
}

.studying-img-box>img {
	width: 100%;
}

.studying-audio {
	box-sizing: border-box;
	position: absolute;
	right: 8px;
	top: 8px;
	width: 18px;
	height: 18px;
	padding-left: 1px;
	background-color: rgba(0, 0, 0, .2);
	border-radius: 50%;
	text-align: center;
	line-height: 18px;
	color: #fff;
}

.studying-audio-time,
.studying-audio span {
	font-size: calc(12px * var(--font-scale, 1));
	font-weight: 400;
}

.studying-audio-time {
	box-sizing: border-box;
	position: absolute;
	bottom: -2px;
	left: 0;
	height: 18px;
	width: 100%;
	padding: 0 5px;
	background: linear-gradient(180deg, hsla(0, 0%, 100%, 0), rgba(0, 0, 0, .5));
	text-align: right;
	color: #fff;
}

.studying-content {
	overflow-x: scroll;
	overflow-y: hidden;
	display: flex;
	padding: 20px 0;
}

.studying-content:after {
	content: "";
	border-right: 1px solid transparent;
}

.studying-content::-webkit-scrollbar {
	display: none;
}

.studying-content-item {
	flex-shrink: 0;
	flex-basis: 137px;
	margin-left: 12px;
	border-radius: 4px;
	line-height: 14px;
	color: #7f848c;
	font-size: calc(12px * var(--font-scale, 1));
	font-weight: 400;
}

.studying-content-item:last-child {
	margin-right: 15px;
}

.studying-content-item:first-child {
	margin-left: 15px;
}

.studying-content-title {
	margin: 12px 0 7px;
	line-height: 18px;
	color: #404040;
	font-size: calc(14px * var(--font-scale, 1));
	font-weight: 500;
	overflow: hidden;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	line-clamp: 2;
	-webkit-box-orient: vertical;
}

.studying-content-sub {
	display: block;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
	color: #4285f4;
}

.studying-skeleton {
	background: #eef0f4;
}

.studying-skel-line {
	height: 14px;
	margin-top: 12px;
	background: #eef0f4;
	border-radius: 4px;
}

.studying-skel-line--short {
	width: 60%;
	height: 12px;
	margin-top: 8px;
}
</style>