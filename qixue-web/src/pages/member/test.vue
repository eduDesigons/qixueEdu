<template>
	<div class="exam-page">
		<!-- 考试记录列表 -->
		<mescroll-body 
			ref="examMescroll"
			@init="mescrollInit"
			@down="downCallback"
			@up="upCallback"
			:down="downOption"
			:up="upOption"
		>
			<div class="exam-list">
				<div 
					class="exam-item" 
					v-for="(item, index) in examList" 
					:key="item.id || index"
					@click="goToExamDetail(item)"
				>
					<div class="exam-cover">
						<img :src="static(paperThumb(item))" class="cover-img" loading="lazy"  style="object-fit:cover" />
					</div>
					<div class="exam-info">
						<div class="exam-title-row">
							<div class="exam-title">{{ paperTitle(item) }}</div>
						</div>
						<div class="exam-course-name" v-if="courseTitle(item)">
							{{ courseTitle(item) }}
						</div>
						<div class="exam-meta">
							<div class="exam-score" v-if="isPending(item)">
								<span class="score-label">状态：</span>
								<span class="score-value score-pending">需要人工批改</span>
							</div>
							<div class="exam-score" v-else-if="item.score !== null && item.score !== undefined">
								<span class="score-label">得分：</span>
								<span class="score-value" :class="getScoreClass(item.score)">{{ item.score }}分</span>
							</div>
							<div class="exam-time">
								<van-icon name="clock-o" size="14" color="#999" />
								<span>{{ timeFrom(item.created_at) }}</span>
							</div>
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
			examList: [],
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
					tip: '暂无记录',
					icon: '/static/empty/empty.png'
				}
			}
		}
	},
	onLoad() {
		this.restoreTestCache();
	},
	methods: {
		restoreTestCache() {
			try {
				const cached = getStorage('qixue_my_test');
				if (cached && cached.length) {
					this.examList = cached;
				}
			} catch (e) {}
		},
		// 下拉刷新的回调
		downCallback() {
			this.mescroll.resetUpScroll();
		},
		// 上拉加载的回调
		upCallback(page) {
			// 加载考试记录列表
			this.getExamList(page);
		},
		// 获取考试记录列表
		getExamList(page) {
			this.$http.get('/user/myTest', {
				page: page.num,
				limit: page.size
			}).then(res => {
				if (res.code === 200) {
					let list = res.data.list || []
					
					// 设置列表数据
					if (page.num === 1) {
						this.examList = list;
						try {
							setStorage('qixue_my_test', list);
						} catch (e) {}
					} else {
						this.examList = this.examList.concat(list);
					}
					this.mescroll.endSuccess(list.length, res.data.hasMore);
				} else {
					this.msg(res.msg || '获取考试记录失败');
					this.mescroll.endErr();
				}
			}).catch(err => {
				console.log('获取考试记录失败:', err);
				this.msg('获取考试记录失败');
				this.mescroll.endErr();
			});
		},
		paperTitle(item) {
			return (item.paper && item.paper.chapter && item.paper.chapter.title)
				|| (item.paper && item.paper.title)
				|| item.title
				|| '未命名测评';
		},
		paperThumb(item) {
			return (item.paper && item.paper.course && item.paper.course.thumb) || '';
		},
		courseTitle(item) {
			return (item.paper && item.paper.course && item.paper.course.title) || '';
		},
		isPending(item) {
			return Number(item.is_people_check) === 1;
		},
		isPassed(item) {
			return Number(item.is_pass) === 1;
		},
		statusText(item) {
			if (this.isPending(item)) return '待批改';
			return this.isPassed(item) ? '已通过' : '未通过';
		},
		statusClass(item) {
			if (this.isPending(item)) return 'pending';
			return this.isPassed(item) ? 'pass' : 'fail';
		},
		canRetake(item) {
			const paperId = (item.paper && item.paper.id) || item.paper_id;
			return !this.isPending(item) && !this.isPassed(item) && !!paperId;
		},
		retakeExam(item) {
			const paperId = (item.paper && item.paper.id) || item.paper_id;
			if (!paperId) {
				this.msg('试卷信息不完整');
				return;
			}
			this.navTo(`/pages/exam/exam?id=${paperId}`);
		},
		goToExamDetail(item) {
			this.navTo(`/pages/exam/examDetail?id=${item.id}`);
		},
		getScoreClass(score) {
			if (score >= 90) return 'score-excellent';
			if (score >= 80) return 'score-good';
			if (score >= 60) return 'score-pass';
			return 'score-fail';
		}
	}
}
</script>

<style lang="scss" scoped>
.exam-page {
	min-height: 100vh;
	background: #f6f7fb;
}

.exam-list {
	padding: calc(20 * var(--rpx)) calc(30 * var(--rpx));
}

.exam-item {
	background: #fff;
	border-radius: calc(16 * var(--rpx));
	margin-bottom: calc(20 * var(--rpx));
	overflow: hidden;
	padding: calc(24 * var(--rpx));
	display: flex;
	align-items: flex-start;
}

.exam-cover {
	position: relative;
	width: calc(180 * var(--rpx));
	height: calc(120 * var(--rpx));
	border-radius: calc(12 * var(--rpx));
	overflow: hidden;
	flex-shrink: 0;
	display: flex;
	align-items: center;
	justify-content: center;
	background: #f5f5f5;
}

.cover-img {
	width: 100%;
	height: 100%;
}

.exam-status-badge {
	position: absolute;
	top: calc(8 * var(--rpx));
	right: calc(8 * var(--rpx));
	padding: calc(4 * var(--rpx)) calc(12 * var(--rpx));
	border-radius: calc(6 * var(--rpx));
	font-size: calc(calc(20 * var(--rpx)) * var(--font-scale, 1));
	color: #fff;
	background: #999;
	
	&.status-pending {
		background: #b9bdc2;
	}
	
	&.status-progress {
		background: #4285f4;
	}
	
	&.status-completed {
		background: #52c41a;
	}
	
	&.status-expired {
		background: #ff4d4f;
	}
}

.exam-info {
	flex: 1;
	margin-left: calc(24 * var(--rpx));
	display: flex;
	flex-direction: column;
}

.exam-title-row {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: calc(12 * var(--rpx));
}

.exam-title {
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
	flex: 1;
	overflow: hidden;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	line-clamp: 2;
	-webkit-box-orient: vertical;
}

.exam-course-name {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
	margin-bottom: calc(12 * var(--rpx));
	overflow: hidden;
	display: -webkit-box;
	-webkit-line-clamp: 1;
	line-clamp: 1;
	-webkit-box-orient: vertical;
}

.exam-meta {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: calc(12 * var(--rpx));
}

.exam-score {
	display: flex;
	align-items: center;
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	
	.score-label {
		color: #7f848c;
		margin-right: calc(8 * var(--rpx));
	}
	
	.score-value {
		font-weight: 500;
		
		&.score-excellent {
			color: #52c41a;
		}
		
		&.score-good {
			color: #4285f4;
		}
		
		&.score-pass {
			color: #faad14;
		}
		
		&.score-fail {
			color: #ff4d4f;
		}
		
		&.score-pending {
			color: #ff9800;
		}
	}
}

.exam-time {
	display: flex;
	align-items: center;
	gap: calc(6 * var(--rpx));
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #999;
}

.exam-desc {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
	line-height: 1.5;
	overflow: hidden;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	line-clamp: 2;
	-webkit-box-orient: vertical;
}
</style>