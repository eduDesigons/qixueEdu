<template>
	<div class="exam-detail-page">
		<div v-if="paper.id || paper.title">
		<div class="exam-detail-content">
			<!-- 考试基本信息 -->
			<div class="exam-info-card">
				<div class="info-header">
					<div class="info-title">{{ paper.title || '' }}</div>
					<div class="info-status" :class="isPeopleCheck ? 'status-pending' : getStatusClass(isPass)">
						<span>{{ isPeopleCheck ? '需要人工批改' : (isPass ? '已通过' : '未通过') }}</span>
					</div>
				</div>
				<div class="score-card" v-if="!isPeopleCheck">
					<div class="score-item">
						<span class="score-label">总分</span>
						<span class="score-number">{{ paper.score || 0 }}</span>
					</div>
					<div class="score-divider"></div>
					<div class="score-item">
						<span class="score-label">得分</span>
						<span class="score-number" :class="getTotalScoreClass()">{{ totalScore }}</span>
					</div>
					<div class="score-divider"></div>
					<div class="score-item">
						<span class="score-label">及格分</span>
						<span class="score-number">{{ paper.pass_score || 0 }}</span>
					</div>
				</div>
			</div>
			
			<!-- 题目列表 -->
			<div class="questions-list">
				<div 
					class="question-item" 
					v-for="(question, index) in questions" 
					:key="question.id || index"
				>
					<!-- 完型填空题 -->
					<div v-if="question.type === 'completion'" class="question-content">
						<div class="question-title-text"><span :selectable="true" :user-select="true">{{ index + 1 }}. {{ question.title }}</span></div>
						<div class="completion-answers">
							<div class="answer-section">
								<div class="answer-tags">
									<div 
										class="answer-tag" 
										v-for="(answer, idx) in getUserCompletionAnswers(question)" 
										:key="idx"
										:class="getCompletionAnswerClass(question, idx)"
									>
										{{ answer || '未作答' }}
									</div>
									<div v-if="getUserCompletionAnswers(question).length === 0" class="no-answer-tag">
										未作答
									</div>
								</div>
							</div>
							<div 
								v-if="isPass && question.answer" 
								class="answer-section correct-answer-section"
							>
								<div class="answer-label">正确答案</div>
								<div class="answer-tags">
									<div 
										class="answer-tag correct-answer-tag" 
										v-for="(answer, idx) in getCorrectCompletionAnswers(question)" 
										:key="idx"
									>
										{{ answer }}
									</div>
								</div>
							</div>
						</div>
					</div>
					
					<!-- 单选题 -->
					<div v-else-if="question.type === 'single'" class="question-content">
						<div class="question-title-text"><span :selectable="true" :user-select="true">{{ index + 1 }}. {{ question.title }}</span></div>
						<div class="options-list">
							<div 
								class="option-item" 
								v-for="(option, optIndex) in question.options" 
								:key="optIndex"
								:class="getOptionClass(question, optIndex)"
							>
								<div class="option-index">{{ String.fromCharCode(65 + optIndex) }}</div>
								<span class="option-text">{{ option }}</span>
								<div 
									v-if="isUserAnswer(question, optIndex)" 
									class="option-badge user-answer"
								>
									我的
								</div>
								<div 
									v-if="isPass && isCorrectAnswer(question, optIndex)" 
									class="option-badge correct-answer-badge"
								>
									正确
								</div>
							</div>
						</div>
					</div>
					
					<!-- 多选题 -->
					<div v-else-if="question.type === 'multiple'" class="question-content">
						<div class="question-title-text"><span :selectable="true" :user-select="true">{{ index + 1 }}. {{ question.title }}</span></div>
						<div class="options-list">
							<div 
								class="option-item" 
								v-for="(option, optIndex) in question.options" 
								:key="optIndex"
								:class="getOptionClass(question, optIndex)"
							>
								<div class="option-index">{{ String.fromCharCode(65 + optIndex) }}</div>
								<span class="option-text">{{ option }}</span>
								<div 
									v-if="isUserAnswer(question, optIndex)" 
									class="option-badge user-answer"
								>
									我的
								</div>
								<div 
									v-if="isPass && isCorrectAnswer(question, optIndex)" 
									class="option-badge correct-answer-badge"
								>
									正确
								</div>
							</div>
						</div>
					</div>
					
					<!-- 文本题 -->
					<div v-else-if="question.type === 'text'" class="question-content">
						<div class="question-title-text"><span :selectable="true" :user-select="true">{{ index + 1 }}. {{ question.title }}</span></div>
						<div class="text-answer-wrapper">
							<div class="text-answer-content" :class="getTextAnswerClass(question)">
								<span  :selectable="true" :user-select="true">{{ getUserTextAnswer(question) || '未作答' }}</span>
							</div>
							<div 
								v-if="isPass && question.answer" 
								class="correct-answer-text"
							>
								<div class="correct-answer-label">正确答案</div>
								<div class="correct-answer-content">
									<span  :selectable="true" :user-select="true">{{ question.answer }}</span>
								</div>
							</div>
						</div>
					</div>
					
					<!-- 上传图片题 -->
					<div v-else-if="question.type === 'uploadImage'" class="question-content">
						<div class="question-title-text"><span :selectable="true" :user-select="true">{{ index + 1 }}. {{ question.title }}</span></div>
						<div class="upload-answer-wrapper">
							<div class="upload-images" v-if="getUserUploadAnswer(question).length > 0">
								<div 
									class="upload-image-item"
									v-for="(img, idx) in getUserUploadAnswer(question)"
									:key="idx"
									@click="previewImage(img)"
								>
									<img :src="static(img)" class="upload-img" style="object-fit:cover" />
								</div>
							</div>
							<div v-else class="no-answer">未作答</div>
						</div>
					</div>
					
					<!-- 上传文档题 -->
					<div v-else-if="question.type === 'uploadDocument'" class="question-content">
						<div class="question-title-text"><span :selectable="true" :user-select="true">{{ index + 1 }}. {{ question.title }}</span></div>
						<div class="upload-answer-wrapper">
							<div class="upload-files" v-if="getUserUploadAnswer(question).length > 0">
								<div 
									class="upload-file-item"
									v-for="(file, idx) in getUserUploadAnswer(question)"
									:key="idx"
									@click="openDocument(file)"
								>
									<div class="upload-file-info">
										<van-icon name="description" size="32" color="#4285f4" />
										<div class="upload-file-detail">
											<span class="upload-file-name">{{ getFileName(file) }}</span>
										</div>
									</div>
								</div>
							</div>
							<div v-else class="no-answer">未作答</div>
						</div>
					</div>
					
					<!-- 答题结果 -->
					<div class="question-result">
						<!-- 上传类型题目，需要人工批改时显示 -->
						<div v-if="isPeopleCheck && (question.type === 'uploadImage' || question.type === 'uploadDocument')" class="result-status status-pending">
							<van-icon name="info-o" size="18" color="#ff9800" />
							<span>需要人工批改</span>
						</div>
						<!-- 其他题目类型，显示正确/错误和得分 -->
						<div v-else class="result-status" :class="getQuestionResultClass(question)">
							<van-icon 
								:name="isQuestionCorrect(question) ? 'checkmark-circle' : 'close-circle'" 
								size="18" 
								:color="isQuestionCorrect(question) ? '#52c41a' : '#ff4d4f'"
							 />
							<span>{{ isQuestionCorrect(question) ? '正确' : '错误' }}</span>
							<span class="result-score">
								<span :class="getQuestionScoreClass(question)">{{ getQuestionScore(question) }}</span>/{{ question.score || 0 }}分
							</span>
						</div>
						<div v-if="getGradeRemark(question)" class="grade-remark">
							<div class="grade-remark-label">批改备注</div>
							<span class="grade-remark-text" :selectable="true" :user-select="true">{{ getGradeRemark(question) }}</span>
						</div>
					</div>
				</div>
			</div>

		</div>
		<div class="exam-detail-footer" v-if="canRetake">
			<div class="exam-detail-retake" @click="retakeExam">重新作答</div>
		</div>
		</div>
		
		<div class="exam-load-fail" v-else-if="pageError" @click="getExamDetail">
			<span>加载失败，点击重试</span>
		</div>
		<!-- 加载中 -->
		<div v-else class="exam-detail-skel">
			<div class="exam-detail-skel-card"></div>
			<div class="exam-detail-skel-line"></div>
			<div class="exam-detail-skel-line exam-detail-skel-line--short"></div>
			<div class="exam-detail-skel-card exam-detail-skel-card--q"></div>
			<div class="exam-detail-skel-card exam-detail-skel-card--q"></div>
		</div>
	</div>
</template>

<script>
import { showToast } from '@/common/toast'
import { previewImage } from '@/common/dom'
import { redirectTo, setPageTitle } from '@/common/nav'
export default {
		data() {
			return {
				id: null,
			loading: false,
			pageError: false,
			examDetail: {
				paper: {},
				questions: [],
				is_pass: false,
				has_passed: false,
				is_people_check: false
			}
			}
		},
		computed: {
			// 获取题目列表
			questions() {
				return this.examDetail.questions || [];
			},
			// 获取试卷信息
			paper() {
				return this.examDetail.paper || {};
			},
		// 是否通过
		isPass() {
			return this.examDetail.is_pass === true || this.examDetail.is_pass === 1;
		},
		// 是否需要人工批改
		isPeopleCheck() {
			return this.examDetail.is_people_check === true || this.examDetail.is_people_check === 1;
		},
		hasPassed() {
			return this.examDetail.has_passed === true || this.examDetail.has_passed === 1 || this.isPass;
		},
		canRetake() {
			return !!(this.paper.id && !this.isPeopleCheck && !this.hasPassed);
		},
			// 计算总得分
			totalScore() {
				let total = 0;
				this.questions.forEach(question => {
					if (question.user_answer && question.user_answer.score !== undefined) {
						total += question.user_answer.score || 0;
					}
				});
				return total;
			}
		},
		onLoad(options) {
			if (options.id) {
				this.id = options.id;
				this.getExamDetail();
			}
		},
		methods: {
			// 获取考试详情
			getExamDetail() {
				if (!this.id) return;
				this.loading = true;
				this.pageError = false;
				this.$http.get('/paper/getPaperLogInfo', { id: this.id }).then(res => {
					this.loading = false;
					if (res.code === 200 || res.data) {
						const data = res.data || res;
					this.examDetail = {
						paper: data.paper || {},
						questions: (data.questions || []).sort((a, b) => (a.sort || 0) - (b.sort || 0)),
						is_pass: data.is_pass || false,
						has_passed: data.has_passed || data.is_pass || false,
						is_people_check: data.is_people_check || false
					};
					if (this.paper.title) {
						setPageTitle(this.paper.title);
					}
					} else {
						this.pageError = true;
						this.msg(res.msg || '获取考试详情失败');
					}
				}).catch(err => {
					this.loading = false;
					this.pageError = true;
					console.log('获取考试详情失败:', err);
					this.msg('获取考试详情失败');
				});
			},
			retakeExam() {
				if (this.hasPassed) {
					this.msg('已通过，无需重新作答');
					return;
				}
				const paperId = this.paper && this.paper.id;
				if (!paperId) {
					this.msg('试卷信息不完整');
					return;
				}
				redirectTo('/pages/exam/exam?id=' + paperId);
			},
		
		// 获取用户完型填空答案列表
		getUserCompletionAnswers(question) {
			const userAnswer = this.getUserAnswerValue(question);
			if (!userAnswer) {
				return [];
			}
			// answer_content 是数组
			if (Array.isArray(userAnswer)) {
				return userAnswer;
			}

			if (typeof userAnswer === 'string') {
				return userAnswer ? [userAnswer] : [];
			}
			return [];
		},
		
		// 获取正确答案列表（完型填空）
		getCorrectCompletionAnswers(question) {
			if (!question.answer || !Array.isArray(question.answer)) {
				return [];
			}
			return question.answer;
		},
		
		// 获取完型填空答案样式类
		getCompletionAnswerClass(question, blankIndex) {
			if (!this.isPass) {
				return 'answer-default';
			}
			const userAnswers = this.getUserCompletionAnswers(question);
			const correctAnswers = this.getCorrectCompletionAnswers(question);
			const userAnswer = userAnswers[blankIndex] || '';
			const correctAnswer = correctAnswers[blankIndex] || '';
			if (userAnswer === correctAnswer) {
				return 'answer-correct';
			}
			return 'answer-incorrect';
		},
		
		// 获取用户答案值（从 user_answer 对象中获取 answer_content 字段）
		getUserAnswerValue(question) {
			if (!question.user_answer) return null;
			// user_answer 是一个对象，需要获取其中的 answer_content 字段
			return question.user_answer.answer_content !== undefined 
				? question.user_answer.answer_content 
				: (question.user_answer.answer || question.user_answer);
		},
		
		// 判断是否是用户答案（单选/多选）
		isUserAnswer(question, optIndex) {
			const userAnswer = this.getUserAnswerValue(question);
			if (!userAnswer) return false;
			if (question.type === 'single') {
				return Array.isArray(userAnswer) 
					? userAnswer[0] === optIndex 
					: userAnswer === optIndex;
			} else if (question.type === 'multiple') {
				return Array.isArray(userAnswer) && userAnswer.includes(optIndex);
			}
			return false;
		},
		
		// 判断是否是正确答案（单选/多选）
		isCorrectAnswer(question, optIndex) {
			if (!question.answer || !this.isPass) return false;
			if (question.type === 'single') {
				return Array.isArray(question.answer) 
					? question.answer[0] === optIndex 
					: question.answer === optIndex;
			} else if (question.type === 'multiple') {
				return Array.isArray(question.answer) && question.answer.includes(optIndex);
			}
			return false;
		},
		
		// 获取选项样式类
		getOptionClass(question, optIndex) {
			const isUser = this.isUserAnswer(question, optIndex);
			const isCorrect = this.isCorrectAnswer(question, optIndex);
			
			if (!this.isPass) {
				return isUser ? 'option-user' : '';
			}
			
			if (isUser && isCorrect) {
				return 'option-correct';
			} else if (isUser && !isCorrect) {
				return 'option-incorrect';
			} else if (!isUser && isCorrect) {
				return 'option-missed';
			}
			return '';
		},
		
		// 获取用户文本答案
		getUserTextAnswer(question) {
			const userAnswer = this.getUserAnswerValue(question);
			// answer_content 可能是字符串或数组
			if (typeof userAnswer === 'string') {
				return userAnswer;
			}
			if (Array.isArray(userAnswer) && userAnswer.length > 0) {
				// 如果是数组，取第一个元素或拼接
				return userAnswer.join(' ');
			}
			return '';
		},
		
		// 获取文本答案样式类
		getTextAnswerClass(question) {
			if (!this.isPass) {
				return 'text-answer-default';
			}
			return this.isQuestionCorrect(question) ? 'text-answer-correct' : 'text-answer-incorrect';
		},
		
		// 获取用户上传答案
		getUserUploadAnswer(question) {
			const userAnswer = this.getUserAnswerValue(question);
			if (!userAnswer) return [];
			return Array.isArray(userAnswer) ? userAnswer : [];
		},
		
		// 预览图片
		previewImage(url) {
			const fullUrl = this.static(url)
			previewImage({
				urls: [fullUrl],
				current: fullUrl
			});
		},
		
		// 从URL中提取文件名
		getFileName(url) {
			if (!url) return '未知文件';
			// 从URL中提取文件名
			const parts = url.split('/');
			const fileName = parts[parts.length - 1];
			// 如果包含查询参数，去掉
			const name = fileName.split('?')[0];
			return name || '未知文件';
		},
		
		// 打开文档
		openDocument(fileUrl) {
			const url = this.static(fileUrl)
			if (!url) {
				showToast('文件地址无效');
				return;
			}
			
			// H5环境：新标签页打开
			if (typeof window !== 'undefined' && window.open) {
				window.open(url, '_blank');
			} else {
				// 降级方案：使用 location.href
				window.location.href = url;
			}
			
		},
		
		// 判断题目是否正确
		isQuestionCorrect(question) {
			// 从 user_answer 对象中获取 is_correct
			if (question.user_answer && question.user_answer.is_correct !== undefined) {
				return question.user_answer.is_correct === true || question.user_answer.is_correct === 1;
			}
			return question.is_correct === true || question.is_correct === 1;
		},
		
		// 获取题目得分
		getQuestionScore(question) {
			// 从 user_answer 对象中获取 score
			if (question.user_answer && question.user_answer.score !== undefined) {
				return question.user_answer.score || 0;
			}
			return question.user_score || question.score || 0;
		},
		getGradeRemark(question) {
			const ua = question && question.user_answer
			if (!ua) return ''
			return String(ua.grade_remark || ua.gradeRemark || '').trim()
		},
		
		// 获取题目得分样式类
		getQuestionScoreClass(question) {
			const score = this.getQuestionScore(question);
			const totalScore = question.score || 0;
			if (score === totalScore) return 'score-full';
			if (score > 0) return 'score-partial';
			return 'score-zero';
		},
		
		// 获取题目结果样式类
		getQuestionResultClass(question) {
			return this.isQuestionCorrect(question) ? 'result-correct' : 'result-incorrect';
		},
		
		// 获取状态样式类
		getStatusClass(isPass) {
			return isPass ? 'status-pass' : 'status-fail';
		},
		
		// 获取分数样式类
		getScoreClass(score) {
			if (score >= 90) return 'score-excellent';
			if (score >= 80) return 'score-good';
			if (score >= 60) return 'score-pass';
			return 'score-fail';
		},
		
		// 获取总得分样式类
		getTotalScoreClass() {
			const totalScore = this.totalScore;
			const paperScore = this.paper.score || 0;
			const passScore = this.paper.pass_score || 0;
			
			if (totalScore >= paperScore * 0.9) return 'score-excellent';
			if (totalScore >= paperScore * 0.8) return 'score-good';
			if (totalScore >= passScore) return 'score-pass';
			return 'score-fail';
		}
	}
}
</script>

<style lang="scss" scoped>
.exam-detail-page {
	min-height: 100vh;
	background: #f6f7fb url('/static/bg/test-bg.png') no-repeat 0 bottom /contain fixed;
}

.exam-detail-content {
	padding: calc(20 * var(--rpx)) calc(30 * var(--rpx)) calc(calc(140 * var(--rpx)) + env(safe-area-inset-bottom));
}

.exam-info-card {
	background: #fff url('/static/bg/test-bg.png') no-repeat right bottom / contain;
	border-radius: calc(16 * var(--rpx));
	padding: calc(30 * var(--rpx));
	margin-bottom: calc(20 * var(--rpx));
	box-shadow: 0 calc(2 * var(--rpx)) calc(12 * var(--rpx)) rgba(76, 139, 244, 0.08);
	position: relative;
	overflow: hidden;
	
	&::before {
		content: '';
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(255, 255, 255, 0.7);
		z-index: 0;
	}
	
	> * {
		position: relative;
		z-index: 1;
	}
}

.info-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: calc(24 * var(--rpx));
}

.info-title {
	flex: 1;
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
	overflow: hidden;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	line-clamp: 2;
	-webkit-box-orient: vertical;
}

.info-status {
	display: inline-flex;
	align-items: center;
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	padding: calc(6 * var(--rpx)) calc(18 * var(--rpx));
	border-radius: calc(20 * var(--rpx));
	margin-left: calc(16 * var(--rpx));
	font-weight: 500;
	flex-shrink: 0;
	
	&.status-pass {
		background: linear-gradient(135deg, #6bc5ff 0%, #4c8bf4 100%);
		color: #ffffff;
	}
	
	&.status-fail {
		background: #ffebee;
		color: #f44336;
	}
	
	&.status-pending {
		background: linear-gradient(135deg, #ffb74d 0%, #ff9800 100%);
		color: #ffffff;
	}
}

.score-card {
	display: flex;
	align-items: center;
	justify-content: space-around;
	background: linear-gradient(135deg, #e8f2ff 0%, #d0e7ff 100%);
	border-radius: calc(16 * var(--rpx));
	padding: calc(28 * var(--rpx)) calc(20 * var(--rpx));
	border: calc(1 * var(--rpx)) solid rgba(76, 139, 244, 0.1);
}

.score-item {
	display: flex;
	flex-direction: column;
	align-items: center;
	flex: 1;
}

.score-label {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
	margin-bottom: calc(8 * var(--rpx));
}

.score-number {
	font-size: calc(calc(44 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 700;
	line-height: 1;
	color: #353535;
	
	&.score-excellent {
		color: #52c41a;
	}
	
	&.score-good {
		color: #4c8bf4;
	}
	
	&.score-pass {
		color: #faad14;
	}
	
	&.score-fail {
		color: #ff4d4f;
	}
}

.score-divider {
	width: calc(1 * var(--rpx));
	height: calc(50 * var(--rpx));
	background: rgba(76, 139, 244, 0.2);
	margin: 0 calc(20 * var(--rpx));
}

.info-time {
	display: flex;
	align-items: center;
	gap: calc(6 * var(--rpx));
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #999;
}

.questions-list {
	padding-bottom: calc(40 * var(--rpx));
}

.question-item {
	background: #fff url('/static/bg/test-bg.png') no-repeat right bottom / contain;
	border-radius: calc(16 * var(--rpx));
	padding: calc(30 * var(--rpx));
	margin-bottom: calc(20 * var(--rpx));
	box-shadow: 0 calc(2 * var(--rpx)) calc(12 * var(--rpx)) rgba(76, 139, 244, 0.06);
	position: relative;
	overflow: hidden;
	
	&::before {
		content: '';
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(255, 255, 255, 0.7);
		z-index: 0;
	}
	
	> * {
		position: relative;
		z-index: 1;
	}
}

.question-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: calc(24 * var(--rpx));
	padding-bottom: calc(20 * var(--rpx));
	border-bottom: calc(1 * var(--rpx)) solid #f0f0f0;
}

.question-number {
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
}

.question-score-info {
	display: flex;
	align-items: center;
	gap: calc(8 * var(--rpx));
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
}

.question-score-label {
	color: #7f848c;
}

.question-score-value {
	font-weight: 500;
	
	&.score-full {
		color: #52c41a;
	}
	
	&.score-partial {
		color: #faad14;
	}
	
	&.score-zero {
		color: #ff4d4f;
	}
}

.question-content {
	margin-bottom: calc(20 * var(--rpx));
}

.question-title-text {
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	color: #353535;
	line-height: 1.6;
	margin-bottom: calc(20 * var(--rpx));
	font-weight: 500;
}

.question-type-tag {
	padding: calc(4 * var(--rpx)) calc(12 * var(--rpx));
	border-radius: calc(4 * var(--rpx));
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	color: #fff;
	flex-shrink: 0;
	
	&.completion-tag {
		background: #722ed1;
	}
	
	&.single-tag {
		background: #4c8bf4;
	}
	
	&.multiple-tag {
		background: #13c2c2;
	}
	
	&.text-tag {
		background: #fa8c16;
	}
	
	&.upload-tag {
		background: #eb2f96;
	}
}

.completion-answers {
	margin-top: 0;
}

.answer-section {
	margin-bottom: calc(16 * var(--rpx));
	
	&:last-child {
		margin-bottom: 0;
	}
}

.answer-label {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #52c41a;
	margin-bottom: calc(8 * var(--rpx));
}

.answer-tags {
	display: flex;
	flex-wrap: wrap;
	gap: calc(12 * var(--rpx));
	align-items: center;
}

.answer-tag {
	display: inline-flex;
	align-items: center;
	justify-content: center;
	padding: calc(14 * var(--rpx)) calc(28 * var(--rpx));
	border-radius: calc(12 * var(--rpx));
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	line-height: 1.4;
	transition: all 0.3s;
	
	&.answer-default {
		background: #f5f5f5;
		color: #7f848c;
		border: calc(1 * var(--rpx)) solid #e8e8e8;
	}
	
	&.answer-correct {
		background: linear-gradient(135deg, #f6ffed 0%, rgba(82, 196, 26, 0.1) 100%);
		color: #52c41a;
		border: calc(1 * var(--rpx)) solid #b7eb8f;
		font-weight: 500;
		box-shadow: 0 calc(2 * var(--rpx)) calc(8 * var(--rpx)) rgba(82, 196, 26, 0.15);
	}
	
	&.answer-incorrect {
		background: linear-gradient(135deg, #fff1f0 0%, rgba(255, 77, 79, 0.1) 100%);
		color: #ff4d4f;
		border: calc(1 * var(--rpx)) solid #ffccc7;
		font-weight: 500;
		box-shadow: 0 calc(2 * var(--rpx)) calc(8 * var(--rpx)) rgba(255, 77, 79, 0.15);
	}
}

.correct-answer-tag {
	background: linear-gradient(135deg, #f6ffed 0%, rgba(82, 196, 26, 0.1) 100%);
	color: #52c41a;
	border: calc(1 * var(--rpx)) solid #b7eb8f;
	font-weight: 500;
	box-shadow: 0 calc(2 * var(--rpx)) calc(8 * var(--rpx)) rgba(82, 196, 26, 0.15);
}

.no-answer-tag {
	display: inline-flex;
	align-items: center;
	justify-content: center;
	padding: calc(12 * var(--rpx)) calc(24 * var(--rpx));
	border-radius: calc(8 * var(--rpx));
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	background: #fafafa;
	color: #999;
	border: calc(1 * var(--rpx)) dashed #e8e8e8;
}

.options-list {
	display: flex;
	flex-direction: column;
	gap: calc(16 * var(--rpx));
}

.option-item {
	display: flex;
	align-items: center;
	padding: calc(24 * var(--rpx)) calc(20 * var(--rpx));
	border-radius: calc(16 * var(--rpx));
	border: calc(2 * var(--rpx)) solid #e8e8e8;
	background: #fafafa;
	position: relative;
	transition: all 0.3s;
	box-shadow: 0 calc(2 * var(--rpx)) calc(8 * var(--rpx)) rgba(0, 0, 0, 0.04);
	
	&.option-user {
		border-color: #4c8bf4;
		background: linear-gradient(135deg, #e8f2ff 0%, rgba(76, 139, 244, 0.05) 100%);
		box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(76, 139, 244, 0.15);
	}
	
	&.option-correct {
		border-color: #52c41a;
		background: linear-gradient(135deg, #f6ffed 0%, rgba(82, 196, 26, 0.05) 100%);
		box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(82, 196, 26, 0.15);
	}
	
	&.option-incorrect {
		border-color: #ff4d4f;
		background: linear-gradient(135deg, #fff1f0 0%, rgba(255, 77, 79, 0.05) 100%);
		box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(255, 77, 79, 0.15);
	}
	
	&.option-missed {
		border-color: #52c41a;
		background: linear-gradient(135deg, #f6ffed 0%, rgba(82, 196, 26, 0.05) 100%);
		border-style: dashed;
		box-shadow: 0 calc(2 * var(--rpx)) calc(8 * var(--rpx)) rgba(82, 196, 26, 0.1);
	}
}

.option-index {
	width: calc(52 * var(--rpx));
	height: calc(52 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
	border-radius: 50%;
	background: linear-gradient(135deg, #ffffff 0%, #f5f7fa 100%);
	color: #4c8bf4;
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	margin-right: calc(16 * var(--rpx));
	flex-shrink: 0;
	box-shadow: 0 calc(2 * var(--rpx)) calc(8 * var(--rpx)) rgba(76, 139, 244, 0.15);
	border: calc(2 * var(--rpx)) solid rgba(76, 139, 244, 0.2);
}

.option-item.option-correct .option-index {
	background: linear-gradient(135deg, #52c41a 0%, #73d13d 100%);
	color: #fff;
	border-color: #52c41a;
	box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(82, 196, 26, 0.3);
}

.option-item.option-incorrect .option-index {
	background: linear-gradient(135deg, #ff4d4f 0%, #ff7875 100%);
	color: #fff;
	border-color: #ff4d4f;
	box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(255, 77, 79, 0.3);
}

.option-item.option-missed .option-index {
	background: linear-gradient(135deg, #52c41a 0%, #73d13d 100%);
	color: #fff;
	border-color: #52c41a;
	box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(82, 196, 26, 0.3);
}

.option-text {
	flex: 1;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #353535;
	line-height: 1.6;
}

.option-badge {
	padding: calc(6 * var(--rpx)) calc(16 * var(--rpx));
	border-radius: calc(12 * var(--rpx));
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	margin-left: calc(12 * var(--rpx));
	font-weight: 500;
	
	&.user-answer {
		background: linear-gradient(135deg, #e8f2ff 0%, rgba(76, 139, 244, 0.1) 100%);
		color: #4c8bf4;
		border: calc(1 * var(--rpx)) solid rgba(76, 139, 244, 0.3);
	}
	
	&.correct-answer-badge {
		background: linear-gradient(135deg, #f6ffed 0%, rgba(82, 196, 26, 0.1) 100%);
		color: #52c41a;
		border: calc(1 * var(--rpx)) solid rgba(82, 196, 26, 0.3);
	}
}

.text-answer-wrapper {
	margin-top: 0;
}

.text-answer-content {
	padding: calc(24 * var(--rpx));
	border-radius: calc(16 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	line-height: 1.8;
	min-height: calc(120 * var(--rpx));
	transition: all 0.3s;
	
	&.text-answer-default {
		background: #fafafa;
		color: #353535;
		border: calc(2 * var(--rpx)) solid #e8e8e8;
	}
	
	&.text-answer-correct {
		background: linear-gradient(135deg, #f6ffed 0%, rgba(82, 196, 26, 0.05) 100%);
		color: #52c41a;
		border: calc(2 * var(--rpx)) solid #b7eb8f;
		box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(82, 196, 26, 0.15);
	}
	
	&.text-answer-incorrect {
		background: linear-gradient(135deg, #fff1f0 0%, rgba(255, 77, 79, 0.05) 100%);
		color: #ff4d4f;
		border: calc(2 * var(--rpx)) solid #ffccc7;
		box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(255, 77, 79, 0.15);
	}
}

.correct-answer-text {
	margin-top: calc(24 * var(--rpx));
	padding: calc(24 * var(--rpx));
	background: linear-gradient(135deg, #f6ffed 0%, rgba(82, 196, 26, 0.05) 100%);
	border-radius: calc(16 * var(--rpx));
	border: calc(2 * var(--rpx)) solid #b7eb8f;
	box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(82, 196, 26, 0.15);
}

.correct-answer-label {
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: #52c41a;
	margin-bottom: calc(12 * var(--rpx));
	font-weight: 500;
}

.correct-answer-content {
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #52c41a;
	line-height: 1.6;
}

.upload-answer-wrapper {
	margin-top: 0;
}

.upload-images {
	display: flex;
	flex-wrap: wrap;
	gap: calc(16 * var(--rpx));
}

.upload-image-item {
	position: relative;
	width: calc(200 * var(--rpx));
	height: calc(200 * var(--rpx));
	border-radius: calc(16 * var(--rpx));
	overflow: hidden;
	box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(0, 0, 0, 0.1);
	transition: all 0.3s;
	
	&:active {
		transform: scale(0.95);
	}
}

.upload-img {
	width: 100%;
	height: 100%;
}

.upload-files {
	display: flex;
	flex-direction: column;
	gap: calc(16 * var(--rpx));
}

.upload-file-item {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: calc(24 * var(--rpx));
	background: #fafafa;
	border-radius: calc(16 * var(--rpx));
	border: calc(2 * var(--rpx)) solid #e8e8e8;
	transition: all 0.3s;
	
	&:active {
		transform: translateY(calc(-2 * var(--rpx)));
		box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(0, 0, 0, 0.1);
	}
}

.upload-file-info {
	display: flex;
	align-items: center;
	flex: 1;
	margin-right: calc(16 * var(--rpx));
}

.upload-file-detail {
	flex: 1;
	margin-left: calc(20 * var(--rpx));
	display: flex;
	flex-direction: column;
	gap: calc(8 * var(--rpx));
}

.upload-file-name {
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #353535;
	font-weight: 500;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.upload-file-url {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #999;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.no-answer {
	padding: calc(20 * var(--rpx));
	background: #fafafa;
	border-radius: calc(12 * var(--rpx));
	color: #999;
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	text-align: center;
}

.question-result {
	margin-top: calc(20 * var(--rpx));
	padding-top: calc(20 * var(--rpx));
	border-top: calc(1 * var(--rpx)) solid #f0f0f0;
}

.grade-remark {
	margin-top: calc(16 * var(--rpx));
	padding: calc(16 * var(--rpx)) calc(20 * var(--rpx));
	background: #f7f8fa;
	border-radius: calc(12 * var(--rpx));
	border-left: calc(6 * var(--rpx)) solid #316c72;
}

.grade-remark-label {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #316c72;
	font-weight: 600;
	margin-bottom: calc(8 * var(--rpx));
}

.grade-remark-text {
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: #333;
	line-height: 1.6;
	white-space: pre-wrap;
	word-break: break-word;
}

.result-status {
	display: flex;
	align-items: center;
	gap: calc(10 * var(--rpx));
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	
	&.result-correct {
		color: #52c41a;
	}
	
	&.result-incorrect {
		color: #ff4d4f;
	}
}

.result-score {
	margin-left: auto;
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	
	.result-status.result-correct & {
		color: #52c41a;
	}
	
	.result-status.result-incorrect & {
		color: #ff4d4f;
	}
	
	.score-full {
		color: #52c41a;
		font-weight: 600;
	}
	
	.score-partial {
		color: #faad14;
		font-weight: 600;
	}
	
	.score-zero {
		color: #ff4d4f;
		font-weight: 600;
	}
}

.exam-detail-skel {
	padding: calc(30 * var(--rpx));
}

.exam-detail-skel-card,
.exam-detail-skel-line {
	background: #eef0f4;
	animation: exam-detail-skel-pulse 1.2s ease-in-out infinite;
}

.exam-detail-skel-card {
	height: calc(220 * var(--rpx));
	border-radius: calc(16 * var(--rpx));
}

.exam-detail-skel-card--q {
	height: calc(160 * var(--rpx));
	margin-top: calc(24 * var(--rpx));
}

.exam-detail-skel-line {
	height: calc(28 * var(--rpx));
	margin-top: calc(24 * var(--rpx));
	border-radius: calc(6 * var(--rpx));
}

.exam-detail-skel-line--short {
	width: 46%;
}

@keyframes exam-detail-skel-pulse {
	0%, 100% { opacity: 1; }
	50% { opacity: 0.55; }
}

.exam-load-fail {
	padding: calc(160 * var(--rpx)) calc(40 * var(--rpx));
	text-align: center;
	color: #7f848c;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
}

.exam-detail-footer {
	position: fixed;
	left: 0;
	right: 0;
	bottom: 0;
	z-index: 200;
	padding: calc(16 * var(--rpx)) calc(30 * var(--rpx)) calc(calc(16 * var(--rpx)) + env(safe-area-inset-bottom));
	background: #fff;
	box-shadow: 0 calc(-4 * var(--rpx)) calc(16 * var(--rpx)) rgba(0, 0, 0, 0.04);
}

.exam-detail-retake {
	height: calc(80 * var(--rpx));
	border-radius: calc(40 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
	background: linear-gradient(135deg, #5c86ff, #8bbeff);
	color: #fff;
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
}
</style>
