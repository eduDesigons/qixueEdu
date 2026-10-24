<template>
	<div class="exam-page">
		<!-- 试卷信息页 -->
		<div v-if="currentPage === 'info'" class="page-wrapper">
			<div class="exam-skeleton" v-if="loading && !paperData.id">
				<div class="skeleton-line skeleton-title"></div>
				<div class="skeleton-line"></div>
				<div class="skeleton-line skeleton-short"></div>
				<div class="skeleton-btn"></div>
			</div>
			<div class="exam-load-fail" v-else-if="pageError" @click="getPaperData">
				<span>加载失败，点击重试</span>
			</div>
			<div v-else class="exam-info-screen">
					<div class="info-card">
						<div class="info-title">{{ paperData.title || '测评信息' }}</div>
						<div class="info-item">
							<span class="info-label">总分：</span>
							<span class="info-value">{{ paperData.score || 0 }}分</span>
						</div>
						<div class="info-item">
							<span class="info-label">及格：</span>
							<span class="info-value">{{ paperData.pass_score || 0 }}分</span>
						</div>
						<div class="info-item">
							<span class="info-label">题量：</span>
							<span class="info-value">{{ questions.length }}题</span>
						</div>
						<div class="info-btn" @click="startExam">
							<span>{{ hasDraft ? '继续上次作答' : '开始答题' }}</span>
						</div>
					</div>
			</div>
		</div>
		
		<!-- 题目页 -->
		<div v-if="currentPage === 'question'" class="page-wrapper">
			<div class="question-screen" v-if="currentQuestion">
				<div class="question-header">
					<div class="question-progress">
						第 {{ currentQuestionIndex + 1 }} / {{ questions.length }} 题
					</div>
					<div class="question-header-right">
						<div class="question-score">分值：{{ currentQuestion.score }}分</div>
						<div class="question-sheet-btn" v-if="!isDesktop" @click="sheetVisible = true">答题卡</div>
					</div>
				</div>
				
				<div class="question-content">
					<!-- 完型填空题 -->
					<div v-if="currentQuestion.type === 'completion'" class="completion-question">
						<div class="question-title-wrapper">
							<template v-for="(part, partIndex) in getCompletionParts(currentQuestion)" :key="partIndex">
								<span v-if="part.type === 'text'" class="title-text">{{ part.text }}</span>
								<div v-else-if="part.type === 'blank'" class="blank-wrapper">
									<div 
										class="blank-input" 
										v-for="(char, charIndex) in part.chars" 
										:key="charIndex"
										:class="{ 'filled': char, 'active': isActiveBlank(part.blankIndex, charIndex) }"
										@click="openOptionPopup(currentQuestion, part.blankIndex, charIndex)"
									>
										{{ char || ' ' }}
									</div>
								</div>
							</template>
						</div>
					</div>
					
					<!-- 单选题 -->
					<div v-else-if="currentQuestion.type === 'single'" class="single-question">
						<div class="question-title-wrapper">
							<span class="question-type-tag single-tag">单选</span>
							<span class="question-title-text">{{ currentQuestion.title }}</span>
						</div>
						<div class="single-options">
							<div 
								class="single-option-item" 
								v-for="(option, optIndex) in currentQuestion.options" 
								:key="optIndex"
								:class="{ 'selected': getUserAnswers(currentQuestion.id) === optIndex }"
								@click="selectSingleOption(currentQuestion, optIndex)"
							>
								<div class="option-index">{{ String.fromCharCode(65 + optIndex) }}</div>
								<span class="option-text">{{ option }}</span>
								<div class="option-indicator" v-if="getUserAnswers(currentQuestion.id) === optIndex"></div>
							</div>
						</div>
					</div>
					
					<!-- 多选题 -->
					<div v-else-if="currentQuestion.type === 'multiple'" class="multiple-question">
						<div class="question-title-wrapper">
							<div class="question-type-tag multiple-tag">多选</div>
							<span class="question-title-text">{{ currentQuestion.title }}</span>
						</div>
						<div class="multiple-options">
							<div 
								class="multiple-option-item" 
								v-for="(option, optIndex) in currentQuestion.options" 
								:key="optIndex"
								:class="{ 'selected': isMultipleOptionSelected(currentQuestion.id, optIndex) }"
								@click="toggleMultipleOption(currentQuestion, optIndex)"
							>
								<div class="option-index">{{ String.fromCharCode(65 + optIndex) }}</div>
								<span class="option-text">{{ option }}</span>
								<div class="option-checkmark" v-if="isMultipleOptionSelected(currentQuestion.id, optIndex)">✓</div>
							</div>
						</div>
					</div>
					
					<!-- 文本题 -->
					<div v-else-if="currentQuestion.type === 'text'" class="text-question">
						<div class="question-title-text">{{ currentQuestion.title }}</div>
						<div class="text-input-wrapper">
							<textarea 
								v-model="userAnswers[currentQuestion.id]" 
								class="text-input"
								placeholder="请输入您的答案..."
								:maxlength="500"
								:auto-height="true"
								@blur="saveDraft"
							></textarea>
							<div class="text-count">{{ (userAnswers[currentQuestion.id] || '').length }}/500</div>
						</div>
					</div>
					
					<!-- 上传图片题 -->
					<div v-else-if="currentQuestion.type === 'uploadImage'" class="upload-question">
						<div class="question-title-wrapper">
							<div class="question-type-tag upload-tag">图片</div>
							<span class="question-title-text">{{ currentQuestion.title }}</span>
						</div>
						<div class="upload-tip">
							<span>最多上传 {{ currentQuestion.options || 1 }} 张图片</span>
						</div>
						<div class="upload-list">
							<div 
								class="upload-item"
								v-for="(item, idx) in (getUploadAnswers(currentQuestion.id) || [])"
								:key="idx"
							>
								<img 
									:src="static(item)" 
									class="upload-image"
									@click="previewUploadImage(currentQuestion.id, idx)"
								 style="object-fit:cover" />
								<div class="upload-delete" @click="removeUploadImage(currentQuestion.id, idx)">
									<van-icon name="cross" size="20" color="#fff" />
								</div>
							</div>
							<div 
								class="upload-add-btn"
								v-if="(getUploadAnswers(currentQuestion.id) || []).length < (currentQuestion.options || 1)"
								@click="chooseUploadImage(currentQuestion)"
							>
								<van-icon name="plus" size="32" color="#999" />
								<span class="upload-add-text">上传图片</span>
							</div>
						</div>
					</div>
					
					<!-- 上传文档题 -->
					<div v-else-if="currentQuestion.type === 'uploadDocument'" class="upload-question">
						<div class="question-title-wrapper">
							<div class="question-type-tag upload-tag">文档</div>
							<span class="question-title-text">{{ currentQuestion.title }}</span>
						</div>
						<div class="upload-tip">
							<span>最多上传 {{ currentQuestion.options || 1 }} 个文档</span>
						</div>
						<div class="upload-list">
							<div 
								class="upload-file-item"
								v-for="(item, idx) in (getUploadAnswers(currentQuestion.id) || [])"
								:key="idx"
							>
								<div class="upload-file-info">
									<van-icon name="description" size="32" color="#4285f4" />
									<div class="upload-file-detail">
										<span class="upload-file-name">{{ getFileName(item) }}</span>
									</div>
								</div>
								<div class="upload-file-delete" @click="removeUploadFile(currentQuestion.id, idx)">
									<van-icon name="cross" size="20" color="#999" />
								</div>
							</div>
							<div 
								class="upload-add-btn upload-file-btn"
								v-if="(getUploadAnswers(currentQuestion.id) || []).length < (currentQuestion.options || 1)"
								@click="chooseUploadFile(currentQuestion)"
							>
								<van-icon name="plus" size="32" color="#999" />
								<span class="upload-add-text">上传文档</span>
							</div>
						</div>
					</div>
				</div>
				
				<div class="question-footer">
					<div class="footer-btn prev-btn" v-if="currentQuestionIndex > 0" @click="goToQuestion(currentQuestionIndex - 1)">
						<van-icon name="arrow-left" size="16" color="#4285f4" />
						<span>上一题</span>
					</div>
					<div class="footer-btn submit-btn" v-if="currentQuestionIndex === questions.length - 1" @click="submitExam">
						<span>交卷</span>
					</div>
					<div class="footer-btn next-btn" v-if="currentQuestionIndex < questions.length - 1" @click="goToQuestion(currentQuestionIndex + 1)">
						<span>下一题</span>
						<van-icon name="arrow" size="16" color="#4285f4" />
					</div>
				</div>
			</div>
		</div>
		
		<!-- 提交结果页 -->
		<div v-if="currentPage === 'result'" class="page-wrapper">
			<div class="submit-result-screen">
				<div class="submit-result-content">
					<div v-if="(paperData.is_people_check == 1 || paperData.is_people_check === true) || (submitResult.is_people_check == 1 || submitResult.is_people_check === true)" class="result-tip">
						<van-icon name="info-o" size="80" color="#ff9800" />
						<span>等待人工批改</span>
						<span class="result-tip-desc">老师批改完成后可在考试记录里查看成绩</span>
						<div class="result-actions">
							<div class="result-btn primary-btn" v-if="submitResult.id" @click="viewResultDetail">
								<span>查看答卷</span>
							</div>
							<div class="result-btn secondary-btn" @click="leaveExam">
								<span>返回</span>
							</div>
						</div>
					</div>
					<div v-else class="result-info">
						<!-- 圆形分数显示 -->
						<div class="score-circle-wrapper">
							<van-circle
								:current-rate="scorePercent"
								:rate="scorePercent"
								:size="240"
								:stroke-width="12"
								color="#4285f4"
								layer-color="#e8f2ff"
								:speed="100"
							>
								<div class="score-circle-inner">
									<span class="score-number">{{ submitResult.score }}</span>
									<span class="score-unit">分</span>
								</div>
							</van-circle>
							<div class="score-circle-label">考试分数</div>
						</div>
						
						<!-- 通过状态标签 -->
						<div class="pass-status-tag" :class="{ 'pass': submitResult.is_pass, 'fail': !submitResult.is_pass }">
							<span>{{ submitResult.is_pass ? '已通过考试' : '未通过考试' }}</span>
						</div>
						
						<!-- 题目统计 -->
						<div class="question-stats">
							正确
								<span class="question-stats-number">{{ correctCount }}</span>
								题,错误
								<span class="question-stats-number">{{ incorrectCount }}</span>
								题
						</div>
						
						<!-- 按钮区域 -->
						<div class="result-actions">
							<div class="result-btn primary-btn" @click="viewResultDetail">
								<span>答题结果</span>
							</div>
							<div class="result-btn secondary-btn" @click="retakeExam" v-if="!submitResult.is_pass">
								<span>重新作答</span>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>

		<div v-if="isDesktop && currentPage === 'question' && questions.length" class="pc-exam-sheet">
			<div class="sheet-title">答题卡</div>
			<div class="sheet-grid">
				<div
					class="sheet-item"
					v-for="(q, i) in questions"
					:key="'pc-' + (q.id || i)"
					:class="sheetItemClass(i)"
					@click="pickSheetQuestion(i)"
				>{{ i + 1 }}</div>
			</div>
			<div class="sheet-legend">
				<span><span class="dot current"></span>当前</span>
				<span><span class="dot done"></span>已答</span>
				<span><span class="dot todo"></span>未答</span>
			</div>
		</div>
		
		<!-- 选项选择弹窗：无蒙版；customStyle 避免部分浏览器 flex:1 高度塌陷 -->
		<van-popup 
			:show="optionPopup.show" 
		 position="bottom" 
		 :close-on-click-overlay="false"
			:overlay="false"
			:zIndex="10090"
			:customStyle="{ flex: 'none', height: 'auto' }"
			@update:show="(v)=>{if(!v)closeOptionPopup()}"
			round="calc(20 * var(--rpx))"
		 :safe-area-inset-bottom="true"
		>
			<div class="option-popup-content">
				<div class="option-popup-list">
					<div 
						class="option-popup-item" 
						v-for="(option, index) in currentQuestionOptions" 
						:key="index"
						:class="{ 'used': isOptionUsedInCurrentBlank(index) }"
						@click="selectOptionForBlank(index)"
					>
						{{ option }}
					</div>
					<div class="option-popup-delete" @click="deleteCurrentChar">
						<van-icon name="revoke" size="24" color="#f56565" />
					</div>
				</div>
			</div>
		</van-popup>

		<van-popup
			:show="sheetVisible"
		 position="bottom"
			round="calc(20 * var(--rpx))"
		 :safe-area-inset-bottom="true"
			@update:show="(v)=>{if(!v)sheetVisible = false}"
		>
			<div class="sheet-panel">
				<div class="sheet-title">答题卡</div>
				<div class="sheet-grid">
					<div
						class="sheet-item"
						v-for="(q, i) in questions"
						:key="q.id || i"
						:class="sheetItemClass(i)"
						@click="pickSheetQuestion(i)"
					>{{ i + 1 }}</div>
				</div>
				<div class="sheet-legend">
					<span><span class="dot current"></span>当前</span>
					<span><span class="dot done"></span>已答</span>
					<span><span class="dot todo"></span>未答</span>
				</div>
			</div>
		</van-popup>

		<van-popup
			:show="unansweredVisible"
		 position="center"
			round="calc(20 * var(--rpx))"
		 :safe-area-inset-bottom="true"
			@update:show="(v)=>{if(!v)unansweredVisible = false}"
		>
			<div class="unanswered-panel">
				<div class="sheet-title">还有 {{ unansweredList.length }} 题未作答</div>
				<div class="sheet-grid">
					<div
						class="sheet-item todo"
						v-for="n in unansweredList"
						:key="n"
						@click="jumpUnanswered(n - 1)"
					>{{ n }}</div>
				</div>
				<div class="unanswered-actions">
					<div class="unanswered-btn ghost" @click="unansweredVisible = false">继续答题</div>
					<div class="unanswered-btn primary" @click="forceSubmitUnanswered">仍要提交</div>
				</div>
			</div>
		</van-popup>
	</div>
</template>

<script>
import { hideLoading, showLoading, showToast } from '@/common/toast'
import { getStorage, removeStorage, setStorage } from '@/common/storage'
import { showModal } from '@/common/dom'
import { navigateBack, redirectTo, setPageTitle, switchTab } from '@/common/nav'
import { isDesktop as desktopOn } from '@/pc/desktop'
import { pickImages, pickFiles, previewImages } from '@/common/media'
	export default {
		data() {
			return {
			paperId: null,
			paperData: {},
			questions: [],
			currentPage: 'info', // 'info' | 'question' | 'result'
			currentQuestionIndex: 0, // 当前题目索引
			userAnswers: {}, // 存储用户答案 {questionId: [[char1, char2], [char3, char4]]} 每个空是一个字符数组
			submitting: false,
			optionPopup: {
				show: false,
				questionId: null,
				blankIndex: null,
				charIndex: null
			},
			optionUsedIndexes: {}, // 存储每个空位每个字符位置对应的选项索引 {questionId: {blankIndex: {charIndex: optionIndex}}}
			submitResult: {
				id: null,
				score: 0,
				is_pass: false,
				is_people_check: false,
				logs: []
			},
			loading: false,
			pageError: false,
			hasDraft: false,
			draftEnabled: false,
			sheetVisible: false,
			unansweredVisible: false,
			unansweredList: []
			}
		},
		onLoad(options) {
			const id = (options && options.id) || (this.$route && this.$route.query && this.$route.query.id)
			if (id) {
				this.paperId = id
				this.getPaperData()
			}
		},
		onHide() {
			this.saveDraft()
		},
		onUnload() {
			this.saveDraft()
		},
		computed: {
		isDesktop() {
			return desktopOn.value
		},
		// 获取当前题目
		currentQuestion() {
			if (this.currentQuestionIndex >= 0 && this.currentQuestionIndex < this.questions.length) {
				return this.questions[this.currentQuestionIndex]
			}
			return null
		},
		// 获取当前题目的选项
		currentQuestionOptions() {
			const qid = this.optionPopup && this.optionPopup.questionId
			const question = this.questions.find(q => String(q.id) === String(qid)) || this.currentQuestion
			if (!question || question.type !== 'completion') return []
			const opts = question.options
			if (Array.isArray(opts)) return opts.filter(c => c !== '' && c != null)
			if (typeof opts === 'string') {
				try {
					const parsed = JSON.parse(opts)
					return Array.isArray(parsed) ? parsed.filter(c => c !== '' && c != null) : []
				} catch (e) {
					return []
				}
			}
			return []
		},
		// 计算正确题目数
		correctCount() {
			if (!this.submitResult.logs || this.submitResult.logs.length === 0) return 0
			return this.submitResult.logs.filter(log => log.is_correct).length
		},
		// 计算错误题目数
		incorrectCount() {
			if (!this.submitResult.logs || this.submitResult.logs.length === 0) return 0
			return this.submitResult.logs.filter(log => !log.is_correct).length
		},
		// 计算分数百分比（用于进度环）
		scorePercent() {
			const totalScore = this.paperData.score || 100
			return Math.min((this.submitResult.score / totalScore) * 100, 100)
		}
		},
		watch: {
			userAnswers: {
				deep: true,
				handler() {
					this.saveDraft()
				}
			},
			'$route.query.id'(id) {
				if (!id || String(id) === String(this.paperId)) return
				this.paperId = id
				this.paperData = {}
				this.questions = []
				this.currentPage = 'info'
				this.getPaperData()
			}
		},
		methods: {
		// 获取试卷数据
		getPaperData() {
			if (!this.paperId) return
			this.loading = true
			this.pageError = false
			this.$http.get('/paper/getPaperInfo', { id: this.paperId }).then(res => {
				this.loading = false
				if (res.data) {
					this.paperData = res.data.paper || {}
					if (res.data.is_people_check !== undefined) {
						this.paperData.is_people_check = res.data.is_people_check
					}
					if (this.paperData.title) {
						setPageTitle(this.paperData.title)
					}
					if (res.data.user_paper_id) {
						this.submitResult.id = res.data.user_paper_id
						this.submitResult.is_people_check = !!res.data.is_people_check
					}
					if (res.data.is_pass && res.data.user_paper_id) {
						redirectTo(`/pages/exam/examDetail?id=${res.data.user_paper_id}`)
						return
					}

					this.questions = (res.data.questions || []).sort((a, b) => a.sort - b.sort)
					this.questions.forEach(q => {
						if (q.type === 'completion' && !Array.isArray(q.options)) {
							if (typeof q.options === 'string') {
								try {
									q.options = JSON.parse(q.options)
								} catch (e) {
									q.options = []
								}
							} else {
								q.options = []
							}
						}
					})

					if (this.paperData.is_people_check == 1 || this.paperData.is_people_check === true) {
						this.currentPage = 'result'
						return
					}

					this.initAnswers()
					this.hasDraft = !!this.readDraft()
				}
			}).catch(err => {
				this.loading = false
				this.pageError = true
				console.log(err)
				this.msg('获取试卷信息失败')
			})
		},
		initAnswers() {
			this.userAnswers = {}
			this.optionUsedIndexes = {}
			this.questions.forEach(q => {
				if (q.type === 'completion') {
					const blanks = this.parseCompletionBlanks(q)
					const answers = blanks.map(blank => new Array(blank.charCount).fill(''))
					this.$set(this.userAnswers, q.id, answers)
				} else if (q.type === 'single' || q.type === 'multiple') {
					this.$set(this.userAnswers, q.id, [])
				} else if (q.type === 'text') {
					this.$set(this.userAnswers, q.id, '')
				} else if (q.type === 'uploadImage' || q.type === 'uploadDocument') {
					this.$set(this.userAnswers, q.id, [])
				}
			})
		},
		draftKey() {
			return 'exam_draft_' + this.paperId
		},
		readDraft() {
			if (!this.paperId) return null
			try {
				const draft = getStorage(this.draftKey())
				if (draft && draft.answers && typeof draft.answers === 'object') return draft
			} catch (e) {}
			return null
		},
		saveDraft() {
			if (!this.paperId || !this.draftEnabled || this.currentPage !== 'question') return
			try {
				setStorage(this.draftKey(), {
					answers: this.userAnswers,
					optionUsedIndexes: this.optionUsedIndexes,
					questionIndex: this.currentQuestionIndex
				})
				this.hasDraft = true
			} catch (e) {}
		},
		clearDraft() {
			try {
				removeStorage(this.draftKey())
			} catch (e) {}
			this.hasDraft = false
		},
		restoreDraft(draft) {
			const answers = draft.answers || {}
			Object.keys(answers).forEach((key) => {
				this.$set(this.userAnswers, key, answers[key])
			})
			this.optionUsedIndexes = draft.optionUsedIndexes || {}
			const idx = Number(draft.questionIndex)
			this.currentQuestionIndex = isNaN(idx) ? 0 : idx
		},
		startExam() {
			if (this.questions.length === 0) {
				this.msg('试卷暂无题目，请联系管理员')
				return
			}
			const draft = this.readDraft()
			if (draft) {
				showModal({
					title: '继续作答',
					content: '检测到未提交的草稿，是否继续上次作答？',
					confirmText: '继续',
					cancelText: '重新开始',
					success: (res) => {
						if (res.confirm) {
							this.restoreDraft(draft)
							this.draftEnabled = true
							this.goToQuestion(this.currentQuestionIndex)
						} else {
							this.clearDraft()
							this.initAnswers()
							this.draftEnabled = true
							this.goToQuestion(0)
						}
					}
				})
				return
			}
			this.draftEnabled = true
			this.goToQuestion(0)
		},
		goToQuestion(index) {
			if (index >= 0 && index < this.questions.length) {
				this.saveDraft()
				this.currentQuestionIndex = index
				this.currentPage = 'question'
				this.closeOptionPopup()
				this.sheetVisible = false
				this.unansweredVisible = false
				this.$nextTick(() => {
					this.autoOpenCompletionPopup()
				})
			}
		},
		findEmptyCompletionSlot(question, fromBlank = 0, fromChar = 0) {
			if (!question || question.type !== 'completion') return null
			const blanks = this.parseCompletionBlanks(question)
			const answers = this.userAnswers[question.id] || []
			for (let bi = fromBlank; bi < blanks.length; bi++) {
				const charCount = blanks[bi].charCount
				const chars = answers[bi] || []
				const start = bi === fromBlank ? fromChar : 0
				for (let ci = start; ci < charCount; ci++) {
					if (!chars[ci]) {
						return { blankIndex: bi, charIndex: ci }
					}
				}
			}
			return null
		},
		autoOpenCompletionPopup() {
			const q = this.currentQuestion
			if (!q || q.type !== 'completion') {
				this.closeOptionPopup()
				return
			}
			const pos = this.findEmptyCompletionSlot(q)
			if (pos) {
				this.openOptionPopup(q, pos.blankIndex, pos.charIndex, false)
			}
		},
		isActiveBlank(blankIndex, charIndex) {
			const q = this.currentQuestion
			return !!(
				this.optionPopup.show &&
				q &&
				String(this.optionPopup.questionId) === String(q.id) &&
				this.optionPopup.blankIndex === blankIndex &&
				this.optionPopup.charIndex === charIndex
			)
		},
		isQuestionAnswered(q) {
			const answer = this.userAnswers[q.id]
			if (q.type === 'completion') {
				if (!answer || answer.length === 0) return false
				return answer.every(blankAnswer =>
					Array.isArray(blankAnswer) && blankAnswer.every(char => char && char !== '')
				)
			}
			if (q.type === 'single' || q.type === 'multiple') {
				return Array.isArray(answer) && answer.length > 0
			}
			if (q.type === 'text') {
				return !!(answer && String(answer).trim() !== '')
			}
			if (q.type === 'uploadImage' || q.type === 'uploadDocument') {
				return Array.isArray(answer) && answer.length > 0
			}
			return false
		},
		getUnansweredIndexes() {
			const list = []
			this.questions.forEach((q, index) => {
				if (!this.isQuestionAnswered(q)) list.push(index)
			})
			return list
		},
		sheetItemClass(index) {
			const q = this.questions[index]
			const current = index === this.currentQuestionIndex
			const done = q && this.isQuestionAnswered(q)
			return {
				current: current,
				done: done && !current,
				todo: !done && !current
			}
		},
		pickSheetQuestion(index) {
			this.sheetVisible = false
			this.goToQuestion(index)
		},
		submitFromSheet() {
			this.sheetVisible = false
			this.submitExam()
		},
		jumpUnanswered(index) {
			this.unansweredVisible = false
			this.goToQuestion(index)
		},
		forceSubmitUnanswered() {
			this.unansweredVisible = false
			this.doSubmit()
		},
		
		// 解析填空题，找出所有空位
		parseCompletionBlanks(question) {
			const title = question.title
			const answerCount = question.answer_count || []
			const blanks = []
			const regex = /_+/g
			let match
			let blankIndex = 0
			
			while ((match = regex.exec(title)) !== null) {
				const charCount = answerCount[blankIndex] || 1
				blanks.push({
					start: match.index,
					end: match.index + match[0].length,
					charCount: charCount,
					index: blankIndex
				})
				blankIndex++
			}
			
			return blanks
		},
		
		// 获取填空题的显示部分（文本+空位）
		getCompletionParts(question) {
			const title = question.title
			const blanks = this.parseCompletionBlanks(question)
			const parts = []
			let lastIndex = 0
			
			blanks.forEach((blank, blankIndex) => {
				// 添加空白前的文本
				if (blank.start > lastIndex) {
					parts.push({
						type: 'text',
						text: title.substring(lastIndex, blank.start)
					})
				}
				
				// 添加空白
				const answers = this.userAnswers[question.id] || []
				const blankAnswers = answers[blankIndex] || []
				parts.push({
					type: 'blank',
					blankIndex: blankIndex,
					chars: blankAnswers.length > 0 ? [...blankAnswers] : new Array(blank.charCount).fill('')
				})
				
				lastIndex = blank.end
			})
			
			// 添加最后剩余的文本
			if (lastIndex < title.length) {
				parts.push({
					type: 'text',
					text: title.substring(lastIndex)
				})
			}
			
			return parts
		},
		
		// 打开选项弹窗
		openOptionPopup(question, blankIndex, charIndex, clearIfFilled = true) {
			if (!question) return
			const answers = this.userAnswers[question.id] || []
			if (clearIfFilled && answers[blankIndex] && answers[blankIndex][charIndex]) {
				this.$set(answers[blankIndex], charIndex, '')
				if (this.optionUsedIndexes[question.id] && this.optionUsedIndexes[question.id][blankIndex]) {
					const blankUsedIndexes = { ...this.optionUsedIndexes[question.id][blankIndex] }
					delete blankUsedIndexes[charIndex]
					this.$set(this.optionUsedIndexes[question.id], blankIndex, blankUsedIndexes)
				}
			}
			
			this.optionPopup = {
				show: true,
				questionId: question.id,
				blankIndex: blankIndex,
				charIndex: charIndex
			}
		},
		
		// 关闭选项弹窗
		closeOptionPopup() {
			this.optionPopup.show = false
		},
		
		// 删除上一个字符
		deleteCurrentChar() {
			if (!this.optionPopup) return
			const { questionId, blankIndex, charIndex } = this.optionPopup
			if (!questionId || blankIndex === null || charIndex === null) return
			
			const answers = this.userAnswers[questionId] || []
			if (!answers[blankIndex]) return
			
			// 从当前位置往前找，找到最后一个有字符的位置
			let deleteIndex = -1
			for (let i = charIndex - 1; i >= 0; i--) {
				if (answers[blankIndex][i] && answers[blankIndex][i] !== '') {
					deleteIndex = i
					break
				}
			}
			
			// 如果没找到上一个字符，且当前位置有字符，则删除当前位置
			if (deleteIndex === -1 && answers[blankIndex][charIndex] && answers[blankIndex][charIndex] !== '') {
				deleteIndex = charIndex
			}
			
			// 如果找到了要删除的位置
			if (deleteIndex !== -1) {
				this.$set(answers[blankIndex], deleteIndex, '')
				// 清除对应的选项索引记录
				if (this.optionUsedIndexes[questionId] && this.optionUsedIndexes[questionId][blankIndex]) {
					const blankUsedIndexes = { ...this.optionUsedIndexes[questionId][blankIndex] }
					delete blankUsedIndexes[deleteIndex]
					this.$set(this.optionUsedIndexes[questionId], blankIndex, blankUsedIndexes)
				}
				// 将焦点移动到删除的位置
				this.$nextTick(() => {
					this.optionPopup.charIndex = deleteIndex
				})
			}
		},
		
		// 判断选项索引在当前空位是否已被使用
		isOptionUsedInCurrentBlank(optionIndex) {
			if (!this.optionPopup) return false
			const { questionId, blankIndex } = this.optionPopup
			if (!questionId || blankIndex === null) return false
			
			const usedIndexes = this.optionUsedIndexes[questionId] || {}
			const blankUsedIndexes = usedIndexes[blankIndex] || {}
			// 检查该选项索引是否已被使用（在任意字符位置）
			return Object.values(blankUsedIndexes).includes(optionIndex)
		},
		
		// 选择选项填入空白
		selectOptionForBlank(optionIndex) {
			if (!this.optionPopup) return
			const { questionId, blankIndex, charIndex } = this.optionPopup
			if (!questionId || blankIndex === null || charIndex === null) return
			
			// 获取选项的字符值（id 可能是 number/string，统一转字符串比较）
			const question = this.questions.find(q => String(q.id) === String(questionId))
			if (!question) return
			const options = Array.isArray(question.options) ? question.options : []
			const option = options[optionIndex]
			if (option == null || option === '') return
			
			// 如果选项索引在当前空位已被使用，不允许重复选择
			if (this.isOptionUsedInCurrentBlank(optionIndex)) {
				this.msg('该选项在当前空位已被使用')
				return
			}
			
			const answers = this.userAnswers[questionId] || []
			if (!answers[blankIndex]) {
				const blanks = this.parseCompletionBlanks(question)
				const charCount = blanks[blankIndex].charCount
				this.$set(answers, blankIndex, new Array(charCount).fill(''))
			}
			
			// 如果当前位置已有值，先清除该位置对应的选项索引记录
			const currentValue = answers[blankIndex][charIndex]
			if (currentValue) {
				// 清除旧值的索引记录
				if (this.optionUsedIndexes[questionId] && this.optionUsedIndexes[questionId][blankIndex]) {
					const blankUsedIndexes = { ...this.optionUsedIndexes[questionId][blankIndex] }
					delete blankUsedIndexes[charIndex]
					this.$set(this.optionUsedIndexes[questionId], blankIndex, blankUsedIndexes)
				}
			}
			
			// 记录使用的选项索引（字符位置 -> 选项索引）
			if (!this.optionUsedIndexes[questionId]) {
				this.$set(this.optionUsedIndexes, questionId, {})
			}
			if (!this.optionUsedIndexes[questionId][blankIndex]) {
				this.$set(this.optionUsedIndexes[questionId], blankIndex, {})
			}
			// 记录当前字符位置对应的选项索引
			const blankUsedIndexes = { ...this.optionUsedIndexes[questionId][blankIndex] }
			blankUsedIndexes[charIndex] = optionIndex
			this.$set(this.optionUsedIndexes[questionId], blankIndex, blankUsedIndexes)
			
			// 填入选项字符值
			this.$set(answers[blankIndex], charIndex, option)
			
			const next = this.findEmptyCompletionSlot(question, blankIndex, charIndex + 1)
			if (next) {
				this.$nextTick(() => {
					this.optionPopup.show = true
					this.optionPopup.blankIndex = next.blankIndex
					this.optionPopup.charIndex = next.charIndex
				})
			} else {
				this.closeOptionPopup()
			}
		},
		

		getUserAnswers(questionId) {
			const answer = this.userAnswers[questionId]
			// 处理单选题：如果是数组，返回第一个元素；如果是数字，直接返回；否则返回null
			if (Array.isArray(answer)) {
				return answer.length > 0 ? answer[0] : null
			}

			return answer !== undefined && answer !== null ? answer : null
		},
		
		// 选择单选选项
		selectSingleOption(question, optionIndex) {
			// 单选题答案格式：数组，包含一个索引 [0]
			this.$set(this.userAnswers, question.id, [optionIndex])
		},
		
		// 判断多选题选项是否被选中
		isMultipleOptionSelected(questionId, optionIndex) {
			const answers = this.userAnswers[questionId] || []
			return Array.isArray(answers) && answers.includes(optionIndex)
		},
		
		// 切换多选题选项
		toggleMultipleOption(question, optionIndex) {
			const answers = this.userAnswers[question.id] || []
			if (!Array.isArray(answers)) {
				this.$set(this.userAnswers, question.id, [])
			}
			
			const currentAnswers = this.userAnswers[question.id]
			const index = currentAnswers.indexOf(optionIndex)
			
			if (index > -1) {
				// 已选中，取消选中
				currentAnswers.splice(index, 1)
				this.$set(this.userAnswers, question.id, [...currentAnswers])
			} else {
				// 未选中，添加到选中列表
				this.$set(this.userAnswers, question.id, [...currentAnswers, optionIndex])
			}
		},
		
		// 提交答案
		submitExam() {
			this.saveDraft()
			const unanswered = this.getUnansweredIndexes().map(i => i + 1)
			if (unanswered.length > 0) {
				this.unansweredList = unanswered
				this.unansweredVisible = true
				return
			}
			this.doSubmit()
		},
		
		// 执行提交
		doSubmit() {
			if (this.submitting) return
			
			this.submitting = true
			showLoading('提交中...')
			
			// 格式化答案数据
			const submitData = {
				id: this.paperId,
				answers: []
			}
			
			this.questions.forEach(q => {
				const answer = this.userAnswers[q.id]
				let answerData = {
					question_id: q.id,
					answer: answer
				}
				
				// 根据题型格式化答案（与组卷代码保持一致）
				if (q.type === 'completion') {
					// 完型填空：将每个空位的字符数组合并为字符串数组
					if (Array.isArray(answer)) {
						answerData.answer = answer.map(blankAnswer => 
							Array.isArray(blankAnswer) ? blankAnswer.join('') : (blankAnswer || '')
						)
					} else {
						answerData.answer = []
					}
				} else if (q.type === 'single') {

					if (Array.isArray(answer)) {
						answerData.answer = answer
					} else if (answer !== null && answer !== undefined) {

						answerData.answer = [answer]
					} else {
						answerData.answer = []
					}
				} else if (q.type === 'multiple') {
					// 多选题：答案已经是数组格式 [0, 1]
					answerData.answer = Array.isArray(answer) ? answer : []
				} else if (q.type === 'text') {
					// 简答题：答案应该是字符串
					answerData.answer = typeof answer === 'string' ? answer : (answer || '')
				} else if (q.type === 'uploadImage' || q.type === 'uploadDocument') {
					// 上传图片/文档：答案应该是数组格式，存储文件URL
					answerData.answer = Array.isArray(answer) ? answer : []
				}
				
				submitData.answers.push(answerData)
			})
			
			this.$http.post('/paper/submitPaper', submitData).then(res => {
				hideLoading()
				this.submitting = false
				if (res.data) {
					this.clearDraft()
					this.draftEnabled = false
					this.showSubmitResult(res.data)
				} else {
					this.msg('提交成功')
					this.clearDraft()
					setTimeout(() => {
						navigateBack()
					}, 1500)
				}
			}).catch(err => {
				hideLoading()
				this.submitting = false
				this.msg(err.message || '提交失败，请稍后重试')
			})
		},
		
		// 显示提交结果
		showSubmitResult(data) {
			this.submitResult = {
				id: data.id,
				score: data.score || 0,
				is_pass: data.is_pass || false,
				is_people_check: data.is_people_check || false,
				logs: data.logs || []
			}
			// 切换到提交结果页
			this.currentPage = 'result'
		},
		
		// 查看答题结果详情
		viewResultDetail() {
			if (!this.submitResult.id) {
				this.msg('暂无答卷')
				return
			}
			redirectTo(`/pages/exam/examDetail?id=${this.submitResult.id}`)
		},
		leaveExam() {
			if (window.history.length > 1) navigateBack(1)
			else switchTab('/pages/index/learn')
		},
		
		// 重新作答
		retakeExam() {
			this.clearDraft()
			this.draftEnabled = false
			this.initAnswers()
			this.currentPage = 'info'
			this.currentQuestionIndex = 0
			this.sheetVisible = false
			this.unansweredVisible = false
		},
		
		// 获取上传答案数组（图片或文档）
		getUploadAnswers(questionId) {
			const answer = this.userAnswers[questionId]
			return Array.isArray(answer) ? answer : []
		},
		
		chooseUploadImage(question) {
			const maxCount = question.options || 1
			const currentAnswers = this.getUploadAnswers(question.id)
			const remainingCount = maxCount - currentAnswers.length
			if (remainingCount <= 0) {
				showToast(`最多只能上传${maxCount}张图片`)
				return
			}
			pickImages({ count: remainingCount }).then((list) => {
				this.uploadFiles(question.id, list.map((item) => item.file))
			}).catch((err) => {
				if (err && err.message !== 'cancel') showToast('选择图片失败')
			})
		},
		removeUploadImage(questionId, index) {
			const answers = this.getUploadAnswers(questionId)
			if (answers[index]) {
				answers.splice(index, 1)
				this.$set(this.userAnswers, questionId, answers)
			}
		},
		previewUploadImage(questionId, index) {
			const answers = this.getUploadAnswers(questionId)
			if (answers[index]) {
				previewImages(answers.map((url) => this.static(url)), index)
			}
		},
		chooseUploadFile(question) {
			const maxCount = question.options || 1
			const currentAnswers = this.getUploadAnswers(question.id)
			const remainingCount = maxCount - currentAnswers.length
			if (remainingCount <= 0) {
				showToast(`最多只能上传${maxCount}个文档`)
				return
			}
			pickFiles({ count: remainingCount }).then((list) => {
				this.uploadFiles(question.id, list.map((item) => item.file))
			}).catch((err) => {
				if (err && err.message !== 'cancel') showToast('选择文件失败')
			})
		},
		removeUploadFile(questionId, index) {
			const answers = this.getUploadAnswers(questionId)
			if (answers[index]) {
				answers.splice(index, 1)
				this.$set(this.userAnswers, questionId, answers)
			}
		},
		uploadFiles(questionId, files) {
			showLoading('上传中...')
			const uploadPromises = files.map((file) => {
				return this.$http.upload(file).then((res) => {
					if (res.code === 200 && res.data) {
						if (typeof res.data === 'string') return res.data
						if (res.data.url) return res.data.url
						return res.data
					}
					throw new Error(res.message || res.msg || '上传失败')
				})
			})
			Promise.all(uploadPromises).then((urls) => {
				hideLoading()
				const currentAnswers = this.getUploadAnswers(questionId)
				this.$set(this.userAnswers, questionId, [...currentAnswers, ...urls])
				showToast('上传成功')
			}).catch((err) => {
				hideLoading()
				showToast(err.message || err.msg || '上传失败')
			})
		},
		
		// 从URL中提取文件名
		getFileName(url) {
			if (!url) return ''
			// 从URL中提取文件名
			const parts = url.split('/')
			const fileName = parts[parts.length - 1]
			// 如果包含查询参数，去掉
			const name = fileName.split('?')[0]
			return name || '未知文件'
		}
		}
	}
</script>

<style lang="scss" scoped>
// 主题色变量
$primary-color: #4286F5;
$primary-light: #e8f2ff;
$primary-dark: #357ae8;
$success-color: #4caf50;
$error-color: #f44336;
$warning-color: #ff9800;
$text-primary: #1a1a1a;
$text-secondary: #666;
$text-tertiary: #999;
$bg-primary: #ffffff;
$bg-secondary: #f5f7fa;
$border-color: #e5e7eb;
$shadow-sm: 0 calc(2 * var(--rpx)) calc(8 * var(--rpx)) rgba(0, 0, 0, 0.04);
$shadow-md: 0 calc(4 * var(--rpx)) calc(16 * var(--rpx)) rgba(0, 0, 0, 0.08);
$shadow-lg: 0 calc(8 * var(--rpx)) calc(24 * var(--rpx)) rgba(0, 0, 0, 0.12);

.exam-page {
	height: 100vh;
	background: url('/static/bg/exam-bg.png') no-repeat 0 bottom /contain fixed;
}

.page-wrapper {
	width: 100%;
	height: 100%;
	display: flex;
	flex-direction: column;
}

.exam-info-screen,
.question-screen {
	flex: 1;
	height: 100%;
	min-height: 0;
	overflow-y: auto;
	padding: calc(32 * var(--rpx));
	box-sizing: border-box;
	-webkit-overflow-scrolling: touch;
}

.info-card {
	background: #fff;
	border-radius: calc(24 * var(--rpx));
	padding: calc(56 * var(--rpx)) calc(40 * var(--rpx));
	margin-top: calc(40 * var(--rpx));
	position: relative;
	overflow: hidden;
	box-shadow: $shadow-sm;
}

.exam-skeleton {
	padding: calc(80 * var(--rpx)) calc(40 * var(--rpx));
}

.skeleton-line {
	height: calc(28 * var(--rpx));
	margin-bottom: calc(28 * var(--rpx));
	border-radius: calc(8 * var(--rpx));
	background: #f0f2f5;
}

.skeleton-title {
	width: 60%;
	height: calc(44 * var(--rpx));
	margin: 0 auto calc(48 * var(--rpx));
}

.skeleton-short {
	width: 40%;
}

.skeleton-btn {
	height: calc(84 * var(--rpx));
	margin-top: calc(48 * var(--rpx));
	border-radius: calc(46 * var(--rpx));
	background: #f0f2f5;
}

.exam-load-fail {
	padding: calc(160 * var(--rpx)) calc(40 * var(--rpx));
	text-align: center;
	color: #7f848c;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
}

.info-title {
	font-size: calc(calc(40 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 700;
	color: $text-primary;
	margin-bottom: calc(48 * var(--rpx));
	text-align: center;
	letter-spacing: calc(1 * var(--rpx));
	position: relative;
	
	&::after {
		content: '';
		position: absolute;
		bottom: calc(-20 * var(--rpx));
		left: 50%;
		transform: translateX(-50%);
		width: calc(80 * var(--rpx));
		height: calc(6 * var(--rpx));
		background: linear-gradient(90deg, $primary-color 0%, lighten($primary-color, 20%) 100%);
		border-radius: calc(3 * var(--rpx));
	}
}

.info-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: calc(32 * var(--rpx)) 0;
	border-bottom: calc(1 * var(--rpx)) solid $border-color;
	transition: all 0.3s;
	
	&:last-of-type {
		border-bottom: none;
	}
	
	&:hover {
		padding-left: calc(12 * var(--rpx));
	}
}

.info-label {
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	color: $text-secondary;
	font-weight: 500;
}

.info-value {
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 700;
	color: $primary-color;
	background: linear-gradient(135deg, $primary-color 0%, lighten($primary-color, 10%) 100%);
	-webkit-background-clip: text;
	-webkit-text-fill-color: transparent;
	background-clip: text;
}

.info-tip {
	display: flex;
	align-items: center;
	justify-content: center;
	margin-top: calc(48 * var(--rpx));
	padding: calc(24 * var(--rpx)) calc(32 * var(--rpx));
	background: linear-gradient(135deg, $primary-light 0%, rgba(66, 134, 245, 0.05) 100%);
	border-radius: calc(16 * var(--rpx));
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: $primary-color;
	border: calc(1 * var(--rpx)) solid rgba(66, 134, 245, 0.2);
	box-shadow: $shadow-sm;
	
	span {
		margin-left: calc(12 * var(--rpx));
		font-weight: 500;
	}
}

.info-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	margin-top: calc(48 * var(--rpx));
	height: calc(84 * var(--rpx));
	background: linear-gradient(135deg, #5c86ff, #8bbeff);
	border-radius: calc(46 * var(--rpx));
	color: #fff;
	font-size: calc(calc(34 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	box-shadow: 0 calc(12 * var(--rpx)) calc(40 * var(--rpx)) calc(-8 * var(--rpx)) rgba(92, 133, 255, .18);
	transition: all 0.3s;
	position: relative;
	overflow: hidden;
	
	&::before {
		content: '';
		position: absolute;
		top: 50%;
		left: 50%;
		width: 0;
		height: 0;
		border-radius: 50%;
		background: rgba(255, 255, 255, 0.2);
		transform: translate(-50%, -50%);
		transition: width 0.6s, height 0.6s;
	}
	
	&:active {
		transform: scale(0.98);
		box-shadow: 0 calc(4 * var(--rpx)) calc(16 * var(--rpx)) rgba(66, 134, 245, 0.25);
		
		&::before {
			width: calc(300 * var(--rpx));
			height: calc(300 * var(--rpx));
		}
	}
	
	span {
		margin-right: calc(12 * var(--rpx));
		position: relative;
		z-index: 1;
	}
	
	.van-icon {
		position: relative;
		z-index: 1;
	}
}

.question-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: calc(32 * var(--rpx));
	padding: calc(24 * var(--rpx)) calc(32 * var(--rpx));
	background: $bg-primary;
	border-radius: calc(20 * var(--rpx));
	box-shadow: $shadow-sm;
}

.question-header-right {
	display: flex;
	align-items: center;
	gap: calc(20 * var(--rpx));
}

.question-sheet-btn {
	padding: calc(8 * var(--rpx)) calc(20 * var(--rpx));
	border-radius: calc(24 * var(--rpx));
	background: $primary-light;
	color: $primary-color;
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
}

.question-progress {
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: $text-secondary;
	font-weight: 500;
	
	&::before {
		content: '';
		display: inline-block;
		width: calc(8 * var(--rpx));
		height: calc(8 * var(--rpx));
		background: $primary-color;
		border-radius: 50%;
		margin-right: calc(12 * var(--rpx));
		vertical-align: middle;
	}
}

.question-score {
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: $primary-color;
	background: linear-gradient(135deg, $primary-light 0%, rgba(66, 134, 245, 0.05) 100%);
	padding: calc(10 * var(--rpx)) calc(20 * var(--rpx));
	border-radius: calc(20 * var(--rpx));
	font-weight: 600;
	border: calc(1 * var(--rpx)) solid rgba(66, 134, 245, 0.2);
}

.question-content {
	background: $bg-primary;
	border-radius: calc(24 * var(--rpx));
	padding: calc(40 * var(--rpx));
	min-height: calc(400 * var(--rpx));
	margin-bottom: calc(32 * var(--rpx));
	box-shadow: $shadow-md;
	border: calc(1 * var(--rpx)) solid rgba(66, 134, 245, 0.08);
}

.question-title-text {
	font-size: calc(calc(34 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: $text-primary;
	line-height: 1.8;
	margin-bottom: calc(32 * var(--rpx));
	white-space: pre-wrap;
	letter-spacing: calc(0.5 * var(--rpx));
}

/* 完型填空样式 */
.completion-question .question-title-wrapper {
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #1a1a1a;
	line-height: 2;
	white-space: pre-wrap;
	display: block;
	margin-bottom: 0;
}

.title-text {
	display: inline;
}

.blank-wrapper {
	display: inline-flex;
	flex-wrap: wrap;
	align-items: center;
	gap: calc(8 * var(--rpx));
	margin: 0 calc(4 * var(--rpx));
	vertical-align: middle;
}

.blank-input {
	min-width: calc(64 * var(--rpx));
	height: calc(68 * var(--rpx));
	line-height: calc(68 * var(--rpx));
	text-align: center;
	background: $bg-primary;
	border: calc(2 * var(--rpx)) solid $primary-color;
	border-radius: calc(12 * var(--rpx));
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	color: $primary-color;
	font-weight: 600;
	box-shadow: $shadow-sm;
	transition: all 0.3s;
	
	&.filled {
		background: linear-gradient(135deg, $primary-light 0%, rgba(66, 134, 245, 0.1) 100%);
		color: $text-primary;
		border-color: $primary-dark;
	}

	&.active {
		border-color: $primary-dark;
		box-shadow: 0 0 0 calc(4 * var(--rpx)) rgba(66, 134, 245, 0.25);
		background: $primary-light;
	}
	
	&:active {
		background: $primary-light;
		transform: scale(0.95);
	}
}

/* 选项弹窗样式 */
.option-popup-content {
	padding: calc(40 * var(--rpx));
	max-height: 70vh;
	border-top-left-radius: calc(24 * var(--rpx));
	border-top-right-radius: calc(24 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
}

.option-popup-list {
	display: flex;
	flex-wrap: wrap;
	gap: calc(16 * var(--rpx));
	max-height: 50vh;
	overflow-y: auto;
	padding-bottom: calc(20 * var(--rpx));
	align-items: center;
	justify-content: center;
}

.option-popup-item {
	width: calc(80 * var(--rpx));
	height: calc(80 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
	background: linear-gradient(135deg, $bg-secondary 0%, $bg-primary 100%);
	border: calc(2 * var(--rpx)) solid $border-color;
	border-radius: calc(16 * var(--rpx));
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: $text-primary;
	text-align: center;
	transition: all 0.3s;
	box-shadow: $shadow-sm;
	
	&:active {
		background: linear-gradient(135deg, $primary-light 0%, rgba(66, 134, 245, 0.1) 100%);
		border-color: $primary-color;
		color: $primary-color;
		transform: scale(0.95);
		box-shadow: 0 calc(6 * var(--rpx)) calc(16 * var(--rpx)) rgba(66, 134, 245, 0.25);
	}
	
	&.used {
		background: #f0f0f0;
		color: #bbb;
		border-color: #e0e0e0;
		opacity: 0.7;
		cursor: not-allowed;
		
		&:active {
			transform: none;
			background: #f0f0f0;
			border-color: #e0e0e0;
			color: #bbb;
		}
	}
}

.option-popup-delete {
	width: calc(80 * var(--rpx));
	height: calc(80 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
	background: linear-gradient(135deg, #fee 0%, #fdd 100%);
	border: calc(2 * var(--rpx)) solid rgba(245, 101, 101, 0.4);
	border-radius: calc(16 * var(--rpx));
	transition: all 0.3s;
	box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(245, 101, 101, 0.2);
	
	&:active {
		background: linear-gradient(135deg, #fdd 0%, #fcc 100%);
		border-color: rgba(245, 101, 101, 0.6);
		transform: scale(0.95);
		box-shadow: 0 calc(6 * var(--rpx)) calc(16 * var(--rpx)) rgba(245, 101, 101, 0.3);
	}
}

/* 单选题样式 */
/* 选择题题目标题样式 */
.single-question .question-title-wrapper,
.multiple-question .question-title-wrapper {
	display: inline-block;
	align-items: flex-start;
	margin-bottom: calc(32 * var(--rpx));
}

.single-question .question-title-text,
.multiple-question .question-title-text {
	margin-bottom: 0;
	flex: 1;
}

.question-type-tag {
	padding: calc(6 * var(--rpx)) calc(16 * var(--rpx));
	border-radius: calc(6 * var(--rpx));
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	margin-right: calc(16 * var(--rpx));
	flex-shrink: 0;
	margin-top: calc(4 * var(--rpx));
	
	&.single-tag {
		background: #e8f4ff;
		color: #4285f4;
	}
	
	&.multiple-tag {
		background: #fff3e0;
		color: #ff9800;
	}
}

.single-options {
	display: flex;
	flex-direction: column;
	gap: calc(20 * var(--rpx));
}

.single-option-item {
	display: flex;
	align-items: center;
	padding: calc(32 * var(--rpx)) calc(28 * var(--rpx));
	background: $bg-primary;
	border: calc(2 * var(--rpx)) solid $border-color;
	border-radius: calc(20 * var(--rpx));
	transition: all 0.3s;
	position: relative;
	box-shadow: $shadow-sm;
	
	&.selected {
		background: linear-gradient(135deg, $primary-light 0%, rgba(66, 134, 245, 0.05) 100%);
		border-color: $primary-color;
		box-shadow: 0 calc(6 * var(--rpx)) calc(20 * var(--rpx)) rgba(66, 134, 245, 0.2);
		transform: translateY(calc(-2 * var(--rpx)));
	}
	
	&:active {
		transform: scale(0.98);
	}
}

.option-index {
	width: calc(52 * var(--rpx));
	height: calc(52 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
	background: $bg-secondary;
	border-radius: 50%;
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 700;
	color: $text-secondary;
	margin-right: calc(24 * var(--rpx));
	flex-shrink: 0;
	transition: all 0.3s;
	box-shadow: $shadow-sm;
}

.single-option-item.selected .option-index {
	background: linear-gradient(135deg, $primary-color 0%, $primary-dark 100%);
	color: #fff;
	box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(66, 134, 245, 0.3);
	transform: scale(1.1);
}

.option-text {
	flex: 1;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #404040;
	line-height: 1.8;
}

.option-indicator {
	width: calc(28 * var(--rpx));
	height: calc(28 * var(--rpx));
	border-radius: 50%;
	background: linear-gradient(135deg, $primary-color 0%, $primary-dark 100%);
	border: calc(4 * var(--rpx)) solid #fff;
	box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(66, 134, 245, 0.4);
	flex-shrink: 0;
	margin-left: calc(16 * var(--rpx));
	animation: pulse 2s infinite;
}

@keyframes pulse {
	0%, 100% {
		transform: scale(1);
		opacity: 1;
	}
	50% {
		transform: scale(1.1);
		opacity: 0.8;
	}
}

/* 多选题样式 */
.multiple-options {
	display: flex;
	flex-direction: column;
	gap: calc(20 * var(--rpx));
}

.multiple-option-item {
	display: flex;
	align-items: center;
	padding: calc(32 * var(--rpx)) calc(28 * var(--rpx));
	background: $bg-primary;
	border: calc(2 * var(--rpx)) solid $border-color;
	border-radius: calc(20 * var(--rpx));
	transition: all 0.3s;
	position: relative;
	box-shadow: $shadow-sm;
	
	&.selected {
		background: linear-gradient(135deg, #fff3e0 0%, rgba(255, 152, 0, 0.05) 100%);
		border-color: $warning-color;
		box-shadow: 0 calc(6 * var(--rpx)) calc(20 * var(--rpx)) rgba(255, 152, 0, 0.2);
		transform: translateY(calc(-2 * var(--rpx)));
	}
	
	&:active {
		transform: scale(0.98);
	}
}

.multiple-option-item.selected .option-index {
	background: linear-gradient(135deg, $warning-color 0%, darken($warning-color, 10%) 100%);
	color: #fff;
	box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(255, 152, 0, 0.3);
	transform: scale(1.1);
}

.option-checkmark {
	width: calc(32 * var(--rpx));
	height: calc(32 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
	background: #ff9800;
	border-radius: calc(6 * var(--rpx));
	color: #fff;
	font-size: calc(calc(20 * var(--rpx)) * var(--font-scale, 1));
	font-weight: bold;
	flex-shrink: 0;
	margin-left: calc(16 * var(--rpx));
}

/* 文本题样式 */
.text-input-wrapper {
	position: relative;
	background: $bg-secondary;
	border-radius: calc(20 * var(--rpx));
	padding: calc(32 * var(--rpx));
	min-height: calc(300 * var(--rpx));
	border: calc(2 * var(--rpx)) solid $border-color;
	box-shadow: $shadow-sm;
}

/* 上传题型样式 */
.upload-question {
	margin-bottom: calc(32 * var(--rpx));
}

.upload-tip {
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: $text-tertiary;
	margin-bottom: calc(24 * var(--rpx));
	font-weight: 500;
	padding: calc(16 * var(--rpx)) calc(24 * var(--rpx));
	background: linear-gradient(135deg, $primary-light 0%, rgba(66, 134, 245, 0.03) 100%);
	border-radius: calc(12 * var(--rpx));
	border-left: calc(4 * var(--rpx)) solid $primary-color;
}

.upload-list {
	display: flex;
	flex-wrap: wrap;
	gap: calc(20 * var(--rpx));
}

.upload-item {
	position: relative;
	width: calc(200 * var(--rpx));
	height: calc(200 * var(--rpx));
	border-radius: calc(16 * var(--rpx));
	overflow: hidden;
	box-shadow: $shadow-md;
	border: calc(2 * var(--rpx)) solid $border-color;
	transition: all 0.3s;
	
	&:active {
		transform: scale(0.95);
	}
}

.upload-image {
	width: 100%;
	height: 100%;
	border-radius: calc(16 * var(--rpx));
}

.upload-delete {
	position: absolute;
	top: calc(-10 * var(--rpx));
	right: calc(-10 * var(--rpx));
	width: calc(40 * var(--rpx));
	height: calc(40 * var(--rpx));
	background: #ff4757;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	box-shadow: 0 calc(2 * var(--rpx)) calc(8 * var(--rpx)) rgba(255, 71, 87, 0.3);
}

.upload-add-btn {
	width: calc(200 * var(--rpx));
	height: calc(200 * var(--rpx));
	border: calc(2 * var(--rpx)) dashed $primary-color;
	border-radius: calc(16 * var(--rpx));
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	background: linear-gradient(135deg, $primary-light 0%, rgba(66, 134, 245, 0.03) 100%);
	transition: all 0.3s;
	box-shadow: $shadow-sm;
	
	&:active {
		background: linear-gradient(135deg, rgba(66, 134, 245, 0.15) 0%, rgba(66, 134, 245, 0.08) 100%);
		border-color: $primary-dark;
		transform: scale(0.95);
	}
}

.upload-file-btn {
	height: calc(120 * var(--rpx));
}

.upload-add-text {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #999;
	margin-top: calc(12 * var(--rpx));
}

.upload-file-item {
	width: 100%;
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: calc(28 * var(--rpx));
	background: $bg-secondary;
	border-radius: calc(16 * var(--rpx));
	margin-bottom: calc(16 * var(--rpx));
	border: calc(2 * var(--rpx)) solid $border-color;
	box-shadow: $shadow-sm;
	transition: all 0.3s;
	
	&:active {
		transform: translateY(calc(-2 * var(--rpx)));
		box-shadow: $shadow-md;
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
	color: #404040;
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

.upload-file-delete {
	padding: calc(8 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
}

.text-input {
	width: 100%;
	min-height: calc(300 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #404040;
	line-height: 1.8;
}

.text-count {
	position: absolute;
	bottom: calc(16 * var(--rpx));
	right: calc(24 * var(--rpx));
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #999;
}

.question-footer {
	display: flex;
	justify-content: space-between;
	gap: calc(24 * var(--rpx));
}

.footer-btn {
	flex: 1;
	height: calc(88 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
	border-radius: calc(20 * var(--rpx));
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	transition: all 0.3s;
	box-shadow: $shadow-sm;
	
	&.prev-btn,
	&.next-btn {
		background: linear-gradient(135deg, $primary-light 0%, rgba(66, 134, 245, 0.05) 100%);
		color: $primary-color;
		border: calc(1 * var(--rpx)) solid rgba(66, 134, 245, 0.3);
		
		&:active {
			background: linear-gradient(135deg, rgba(66, 134, 245, 0.15) 0%, rgba(66, 134, 245, 0.1) 100%);
			transform: scale(0.98);
		}
	}
	
	&.submit-btn {
		background: linear-gradient(135deg, #5c86ff, #8bbeff);
		color: #fff;
		box-shadow: 0 calc(8 * var(--rpx)) calc(24 * var(--rpx)) rgba(66, 134, 245, 0.3);
		
		&:active {
			transform: scale(0.98);
			box-shadow: 0 calc(4 * var(--rpx)) calc(16 * var(--rpx)) rgba(66, 134, 245, 0.25);
		}
	}
	
	span {
		margin: 0 calc(12 * var(--rpx));
	}
}

.sheet-panel,
.unanswered-panel {
	padding: calc(32 * var(--rpx)) calc(32 * var(--rpx)) calc(40 * var(--rpx));
	background: #fff;
}

.unanswered-panel {
	width: calc(620 * var(--rpx));
	max-width: 86vw;
	border-radius: calc(20 * var(--rpx));
}

.sheet-title {
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: $text-primary;
	text-align: center;
	margin-bottom: calc(28 * var(--rpx));
}

.sheet-grid {
	display: flex;
	flex-wrap: wrap;
	gap: calc(16 * var(--rpx));
	margin-bottom: calc(24 * var(--rpx));
}

.sheet-item {
	width: calc(72 * var(--rpx));
	height: calc(72 * var(--rpx));
	border-radius: calc(16 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	background: #f5f7fa;
	color: #7f848c;
	border: calc(2 * var(--rpx)) solid transparent;
}

.sheet-item.done {
	background: #e8f1ff;
	color: $primary-color;
}

.sheet-item.current {
	background: $primary-color;
	color: #fff;
}

.sheet-item.todo {
	background: #fff5f7;
	color: #ff3359;
}

.sheet-legend {
	display: flex;
	justify-content: center;
	gap: calc(32 * var(--rpx));
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	color: $text-secondary;
	margin-bottom: calc(24 * var(--rpx));

	.dot {
		display: inline-block;
		width: calc(16 * var(--rpx));
		height: calc(16 * var(--rpx));
		border-radius: 50%;
		margin-right: calc(8 * var(--rpx));
		vertical-align: middle;
	}

	.dot.current { background: $primary-color; }
	.dot.done { background: #8bbeff; }
	.dot.todo { background: #ff3359; }
}

.sheet-submit {
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

.unanswered-actions {
	display: flex;
	gap: calc(16 * var(--rpx));
	margin-top: calc(8 * var(--rpx));
}

.unanswered-btn {
	flex: 1;
	height: calc(72 * var(--rpx));
	border-radius: calc(36 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
}

.unanswered-btn.ghost {
	background: #f5f7fa;
	color: #404040;
}

.unanswered-btn.primary {
	background: linear-gradient(135deg, #5c86ff, #8bbeff);
	color: #fff;
}

/* 提交结果页样式 */
.submit-result-screen {
	width: 100%;
	height: 100%;
	display: flex;
	flex-direction: column;
	background: #ffffff url('/static/bg/test-bg.png') no-repeat 0 bottom /contain fixed;
}

.submit-result-content {
	flex: 1;
	display: flex;
	flex-direction: column;
	padding: calc(60 * var(--rpx)) calc(30 * var(--rpx)) calc(40 * var(--rpx));
	overflow-y: auto;
}

.result-tip {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding: calc(60 * var(--rpx)) calc(40 * var(--rpx));
	
	span {
		margin-top: calc(40 * var(--rpx));
		font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
		color: #666;
		text-align: center;
		line-height: 1.6;
	}

	.result-tip-desc {
		margin-top: calc(16 * var(--rpx)) !important;
		font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1)) !important;
		color: #9ca3af !important;
	}
}

.result-info {
	flex: 1;
	display: flex;
	flex-direction: column;
	align-items: center;
}

/* 圆形分数显示 */
.score-circle-wrapper {
	display: flex;
	flex-direction: column;
	align-items: center;
	margin-bottom: calc(40 * var(--rpx));
}

.score-circle-wrapper :deep(.van-circle) {
	display: block;
	position: relative;
	margin: 0 auto;
}

/* van-circle 默认插槽不会包进 .van-circle__text，需自行绝对居中 */
.score-circle-inner {
	position: absolute;
	inset: 0;
	z-index: 1;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	text-align: center;
}

.score-number {
	font-size: calc(calc(72 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 700;
	color: #4285f4;
	line-height: 1;
	margin-bottom: calc(4 * var(--rpx));
}

.score-unit {
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #4285f4;
	line-height: 1;
}

.score-circle-label {
	margin-top: calc(20 * var(--rpx));
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: #999;
	text-align: center;
}

/* 通过状态标签 */
.pass-status-tag {
	display: inline-flex;
	align-items: center;
	justify-content: center;
	padding: calc(12 * var(--rpx)) calc(32 * var(--rpx));
	border-radius: calc(40 * var(--rpx));
	margin-bottom: calc(30 * var(--rpx));
	
	&.pass {
		background: linear-gradient(135deg, #6bc5ff 0%, #4285f4 100%);
		color: #ffffff;
	}
	
	&.fail {
		background: #ffebee;
		color: #f44336;
	}
	
	span {
		font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
		font-weight: 500;
	}
}

/* 题目统计 */
.question-stats {
	margin-top: calc(28 * var(--rpx));
    font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
    line-height: calc(32 * var(--rpx));
    font-weight: 400;
    color: #353535;
	.question-stats-number{
		color: #5c86ff;
		font-weight: 500;
	}
}

/* 按钮区域 */
.result-actions {
	width: 100%;
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: calc(24 * var(--rpx));
	margin-top: calc(80 * var(--rpx));
}

.result-btn {
	width: calc(500 * var(--rpx));
	height: calc(84 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
	border-radius: calc(46 * var(--rpx));
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	transition: all 0.3s;
	
	&.primary-btn {
		background: linear-gradient(135deg, #5c86ff, #8bbeff);
		color: #ffffff;
		box-shadow: 0 calc(12 * var(--rpx)) calc(40 * var(--rpx)) calc(-8 * var(--rpx)) rgba(92, 133, 255, .18);
		
		&:active {
			transform: scale(0.98);
			box-shadow: 0 calc(4 * var(--rpx)) calc(16 * var(--rpx)) rgba(66, 134, 245, 0.25);
		}
	}
	
	&.secondary-btn {
		background: #ffffff;
		border: calc(2 * var(--rpx)) solid #e8e8e8;
		color: #b2b2b2;
		
		&:active {
			background: #f5f5f5;
		}
	}
	
	span {
		letter-spacing: calc(1 * var(--rpx));
	}
}

/* 剩余次数 */
.remaining-attempts {
	text-align: center;
	
	span {
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #999;
	}
}
</style>