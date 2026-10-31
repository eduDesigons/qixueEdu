<template>
	<div class="exam-page">

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
						<div class="info-restart" v-if="hasDraft" @click="restartExam">重新开始</div>
					</div>
			</div>
		</div>

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
					<ExamQuestion
						:question="currentQuestion"
						:user-answers="userAnswers"
						:option-popup="optionPopup"
						@blank="openOptionPopup"
						@select-single="selectSingleOption"
						@toggle-multiple="toggleMultipleOption"
						@update-text="onUpdateText"
						@blur-text="saveDraft"
						@preview-image="previewUploadImage"
						@remove-image="removeUploadImage"
						@choose-image="chooseUploadImage"
						@remove-file="removeUploadFile"
						@choose-file="chooseUploadFile"
					/>
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

						<div class="pass-status-tag" :class="{ 'pass': submitResult.is_pass, 'fail': !submitResult.is_pass }">
							<span>{{ submitResult.is_pass ? '已通过考试' : '未通过考试' }}</span>
						</div>

						<div class="question-stats">
							正确
								<span class="question-stats-number">{{ correctCount }}</span>
								题,错误
								<span class="question-stats-number">{{ incorrectCount }}</span>
								题
						</div>

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
			<ExamSheet title="答题卡" :questions="questions" :item-class="sheetItemClass" key-prefix="pc-" @pick="pickSheetQuestion" />
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
				<ExamSheet title="答题卡" :questions="questions" :item-class="sheetItemClass" @pick="pickSheetQuestion" />
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
import { navigateBack, redirectTo, setPageTitle, switchTab } from '@/common/nav'
import { isDesktop as desktopOn } from '@/pc/desktop'
import { pickImages, pickFiles, previewImages } from '@/common/media'
import ExamQuestion from './components/ExamQuestion.vue'
import ExamSheet from './components/ExamSheet.vue'
import { parseCompletionBlanks, isQuestionAnswered, getUploadAnswers } from './helpers'
	export default {
		components: { ExamQuestion, ExamSheet },
		data() {
			return {
			paperId: null,
			paperData: {},
			questions: [],
			currentPage: 'info',
			currentQuestionIndex: 0,
			userAnswers: {},
			submitting: false,
			optionPopup: {
				show: false,
				questionId: null,
				blankIndex: null,
				charIndex: null
			},
			optionUsedIndexes: {},
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
		currentQuestion() {
			if (this.currentQuestionIndex >= 0 && this.currentQuestionIndex < this.questions.length) {
				return this.questions[this.currentQuestionIndex]
			}
			return null
		},
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
		correctCount() {
			if (!this.submitResult.logs || this.submitResult.logs.length === 0) return 0
			return this.submitResult.logs.filter(log => log.is_correct).length
		},
		incorrectCount() {
			if (!this.submitResult.logs || this.submitResult.logs.length === 0) return 0
			return this.submitResult.logs.filter(log => !log.is_correct).length
		},
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
				this.restoreDraft(draft)
			}
			this.draftEnabled = true
			this.goToQuestion(draft ? this.currentQuestionIndex : 0)
		},
		restartExam() {
			this.clearDraft()
			this.initAnswers()
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
		isQuestionAnswered(q) {
			return isQuestionAnswered(q, this.userAnswers[q.id])
		},
		onUpdateText(question, value) {
			this.$set(this.userAnswers, question.id, value)
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

		parseCompletionBlanks(question) {
			return parseCompletionBlanks(question)
		},

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

		closeOptionPopup() {
			this.optionPopup.show = false
		},

		deleteCurrentChar() {
			if (!this.optionPopup) return
			const { questionId, blankIndex, charIndex } = this.optionPopup
			if (!questionId || blankIndex === null || charIndex === null) return

			const answers = this.userAnswers[questionId] || []
			if (!answers[blankIndex]) return

			let deleteIndex = -1
			for (let i = charIndex - 1; i >= 0; i--) {
				if (answers[blankIndex][i] && answers[blankIndex][i] !== '') {
					deleteIndex = i
					break
				}
			}

			if (deleteIndex === -1 && answers[blankIndex][charIndex] && answers[blankIndex][charIndex] !== '') {
				deleteIndex = charIndex
			}

			if (deleteIndex !== -1) {
				this.$set(answers[blankIndex], deleteIndex, '')
				if (this.optionUsedIndexes[questionId] && this.optionUsedIndexes[questionId][blankIndex]) {
					const blankUsedIndexes = { ...this.optionUsedIndexes[questionId][blankIndex] }
					delete blankUsedIndexes[deleteIndex]
					this.$set(this.optionUsedIndexes[questionId], blankIndex, blankUsedIndexes)
				}
				this.$nextTick(() => {
					this.optionPopup.charIndex = deleteIndex
				})
			}
		},

		isOptionUsedInCurrentBlank(optionIndex) {
			if (!this.optionPopup) return false
			const { questionId, blankIndex } = this.optionPopup
			if (!questionId || blankIndex === null) return false

			const usedIndexes = this.optionUsedIndexes[questionId] || {}
			const blankUsedIndexes = usedIndexes[blankIndex] || {}
			return Object.values(blankUsedIndexes).includes(optionIndex)
		},

		selectOptionForBlank(optionIndex) {
			if (!this.optionPopup) return
			const { questionId, blankIndex, charIndex } = this.optionPopup
			if (!questionId || blankIndex === null || charIndex === null) return

			const question = this.questions.find(q => String(q.id) === String(questionId))
			if (!question) return
			const options = Array.isArray(question.options) ? question.options : []
			const option = options[optionIndex]
			if (option == null || option === '') return

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

			const currentValue = answers[blankIndex][charIndex]
			if (currentValue) {
				if (this.optionUsedIndexes[questionId] && this.optionUsedIndexes[questionId][blankIndex]) {
					const blankUsedIndexes = { ...this.optionUsedIndexes[questionId][blankIndex] }
					delete blankUsedIndexes[charIndex]
					this.$set(this.optionUsedIndexes[questionId], blankIndex, blankUsedIndexes)
				}
			}

			if (!this.optionUsedIndexes[questionId]) {
				this.$set(this.optionUsedIndexes, questionId, {})
			}
			if (!this.optionUsedIndexes[questionId][blankIndex]) {
				this.$set(this.optionUsedIndexes[questionId], blankIndex, {})
			}
			const blankUsedIndexes = { ...this.optionUsedIndexes[questionId][blankIndex] }
			blankUsedIndexes[charIndex] = optionIndex
			this.$set(this.optionUsedIndexes[questionId], blankIndex, blankUsedIndexes)

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

		selectSingleOption(question, optionIndex) {
			this.$set(this.userAnswers, question.id, [optionIndex])
		},

		toggleMultipleOption(question, optionIndex) {
			const answers = this.userAnswers[question.id] || []
			if (!Array.isArray(answers)) {
				this.$set(this.userAnswers, question.id, [])
			}

			const currentAnswers = this.userAnswers[question.id]
			const index = currentAnswers.indexOf(optionIndex)

			if (index > -1) {
				currentAnswers.splice(index, 1)
				this.$set(this.userAnswers, question.id, [...currentAnswers])
			} else {
				this.$set(this.userAnswers, question.id, [...currentAnswers, optionIndex])
			}
		},

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

		doSubmit() {
			if (this.submitting) return

			this.submitting = true
			showLoading('提交中...')

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

				if (q.type === 'completion') {
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
					answerData.answer = Array.isArray(answer) ? answer : []
				} else if (q.type === 'text') {
					answerData.answer = typeof answer === 'string' ? answer : (answer || '')
				} else if (q.type === 'uploadImage' || q.type === 'uploadDocument') {
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

		showSubmitResult(data) {
			this.submitResult = {
				id: data.id,
				score: data.score || 0,
				is_pass: data.is_pass || false,
				is_people_check: data.is_people_check || false,
				logs: data.logs || []
			}
			this.currentPage = 'result'
		},

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

		retakeExam() {
			this.clearDraft()
			this.draftEnabled = false
			this.initAnswers()
			this.currentPage = 'info'
			this.currentQuestionIndex = 0
			this.sheetVisible = false
			this.unansweredVisible = false
		},

		getUploadAnswers(questionId) {
			return getUploadAnswers(this.userAnswers[questionId])
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

		}
	}
</script>

<style lang="scss" scoped>
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
	line-height: 1;
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

.info-restart {
	margin-top: calc(24 * var(--rpx));
	text-align: center;
	color: #8b93a1;
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	line-height: 1.4;
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
	display: inline-flex;
	align-items: center;
	justify-content: center;
	padding: calc(8 * var(--rpx)) calc(20 * var(--rpx));
	border-radius: calc(24 * var(--rpx));
	background: $primary-light;
	color: $primary-color;
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	line-height: 1;
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
	line-height: 1;
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
		line-height: 1;
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

.sheet-item.todo {
	background: #fff5f7;
	color: #ff3359;
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
	line-height: 1;
}

.unanswered-btn.ghost {
	background: #f5f7fa;
	color: #404040;
}

.unanswered-btn.primary {
	background: linear-gradient(135deg, #5c86ff, #8bbeff);
	color: #fff;
}

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
		line-height: 1;
	}
}

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
	line-height: 1;
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
		line-height: 1;
	}
}

.remaining-attempts {
	text-align: center;

	span {
		font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
		color: #999;
	}
}
</style>