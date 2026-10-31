<template>
  <CommonPage back fill>
    <template #action>
      <NButton type="primary" :loading="saving" :disabled="saving" @click="handleSave()">
        {{ chapterId ? '保存测评' : '保存试卷' }}
      </NButton>
      <NDropdown v-if="route.params.id" trigger="click" :options="moreOptions" @select="handleMoreSelect">
        <NButton secondary class="page-more-btn" aria-label="更多">
          <i class="i-fe:more-horizontal text-14" />
        </NButton>
      </NDropdown>
    </template>

    <div class="paper-editor-main">

      <div class="paper-editor-left">
        <div class="editor-panel-head">
          题型选择
        </div>
        <div class="question-type-list">
            <div
              v-for="type in questionTypes"
              :key="type.id"
              class="question-type-item"
              @click="addQuestion(type)"
            >
              <div class="question-type-icon">
                <i :class="type.icon" />
              </div>
              <div class="question-type-content">
                <div class="question-type-name">
                  {{ type.name }}
                </div>
                <div class="question-type-desc">
                  {{ type.desc }}
                </div>
              </div>
            </div>
        </div>
      </div>

      <div class="paper-editor-center">
        <div class="editor-panel-head">
          <span>题目列表</span>
          <div class="editor-panel-head__extra">
            <span>分数线</span>
            <NInputNumber
              v-model:value="passScore"
              :min="0"
              :max="totalScore"
              placeholder="分数线"
              size="small"
            />
            <span>/ 总分</span>
            <span class="color-primary font-bold">
              <n-number-animation :from="0" :to="totalScore" />
            </span>
          </div>
        </div>
        <div v-if="questions.length === 0" class="empty-state">
            <i class="i-material-symbols:menu-outline-rounded" style="font-size: 48px; margin-bottom: 16px;" />
            <p>暂无题目，请从左侧添加题目</p>
          </div>

          <div v-else class="questions-container">
            <VueDraggable
              v-model="questions"
              tag="div"
              class="draggable-container"
              :animation="150"
              ghost-class="ghost"
              chosen-class="chosen"
              filter=".question-item-delete"
              :prevent-on-filter="true"
            >
              <div
                v-for="(question, index) in questions"
                :key="question._uniqueKey || index"
                class="question-item"
                :class="{ active: selectIndex === index }"
                @click="openEditModal(index)"
              >
                <button
                  type="button"
                  class="question-item-delete"
                  aria-label="删除题目"
                  @click.stop="confirmRemoveQuestion(index)"
                >
                  <i class="i-material-symbols:delete" />
                </button>

                <div class="question-item-header">
                  <div class="question-item-title">
                    <NTag size="small" round :bordered="false" :type="getQuestionTypeMeta(question.type).tagType">
                      {{ getQuestionTypeMeta(question.type).name }}
                    </NTag>
                    <h4>{{ index + 1 }}. {{ question.title }}</h4>
                  </div>
                </div>

                <div class="question-item-content">
                  <div class="question-score">
                    分值: {{ question.score }}分
                  </div>

                  <div
                    v-if="question.type === 'single' || question.type === 'multiple'"
                    class="question-options"
                  >
                    <div
                      v-for="(option, optIndex) in question.options"
                      :key="optIndex"
                      class="question-option"
                      :class="{ 'color-primary font-bold underline': question.answer.includes(optIndex) }"
                    >
                      {{ String.fromCharCode(65 + optIndex) }}. {{ option }}
                    </div>
                  </div>

                  <div
                    v-else-if="question.type === 'completion'"
                    class="question-answer"
                  >
                    答案:
                    <span v-for="(answer, idx) in getCompletionAnswers(question)" :key="idx" class="mr-10 color-primary font-bold underline">
                      {{ answer || '（未填写）' }}
                    </span>
                  </div>

                  <div
                    v-else-if="question.type === 'text'"
                    class="question-answer"
                  >
                    <div v-if="question.answer && question.answer.trim()">
                      <div class="mb-2">
                        参考答案: <span class="color-primary font-bold">{{ question.answer }}</span>
                      </div>
                    </div>
                    <span v-else class="text-gray-400">
                      （未设置答案）
                    </span>
                  </div>

                  <div
                    v-else-if="question.type === 'uploadImage'"
                    class="question-upload color-primary"
                  >
                    图片附件题（学生需上传图片）
                  </div>
                  <div
                    v-else-if="question.type === 'uploadDocument'"
                    class="question-upload color-primary"
                  >
                    文档附件题（学生需上传文档）
                  </div>
                </div>
              </div>
            </VueDraggable>
          </div>
      </div>
    </div>

    <NModal
      v-model:show="editModal.show"
      preset="card"
      :style="{ width: '90%', maxWidth: '1200px' }"
      :title="editModal.title"
      size="large"
      class="question-edit-modal"
    >
      <div class="modal-content-wrapper">

        <div class="modal-preview">
          <div class="editor-panel-head">
            题目预览
          </div>
          <div class="modal-panel">
            <div v-if="selectedQuestion" class="question-preview">
              <div class="preview-header">
                <div class="preview-title-row">
                  <NTag size="small" round :bordered="false" :type="getQuestionTypeMeta(selectedQuestion.type).tagType">
                    {{ getQuestionTypeMeta(selectedQuestion.type).name }}
                  </NTag>
                  <h3>{{ selectedQuestion.title }}</h3>
                </div>
                <div class="preview-score">
                  分值: {{ selectedQuestion.score }}分
                </div>
              </div>

              <div class="preview-content">
                <div
                  v-if="selectedQuestion.type === 'single' || selectedQuestion.type === 'multiple'"
                  class="preview-options"
                >
                  <div
                    v-for="(option, optIndex) in selectedQuestion.options"
                    :key="optIndex"
                    class="preview-option"
                    :class="{ 'correct-answer': selectedQuestion.answer.includes(optIndex) }"
                  >
                    {{ String.fromCharCode(65 + optIndex) }}. {{ option }}
                  </div>
                </div>

                <div
                  v-else-if="selectedQuestion.type === 'completion'"
                  class="preview-answer"
                >
                  <div class="completion-preview">
                    <span v-for="(part, idx) in getCompletionPreviewParts(selectedQuestion)" :key="idx">
                      <span v-if="part.isBlank" class="blank-preview">{{ part.answer || '____' }}</span>
                      <span v-else>{{ part.text }}</span>
                    </span>
                  </div>
                </div>

                <div
                  v-else-if="selectedQuestion.type === 'text'"
                  class="preview-answer"
                >
                  <div v-if="selectedQuestion.answer && selectedQuestion.answer.trim()">
                    <div class="text-answer-preview mb-3">
                      <div class="preview-label">
                        参考答案：
                      </div>
                      <div class="preview-content">
                        {{ selectedQuestion.answer }}
                      </div>
                    </div>
                  </div>
                  <div v-else class="text-gray-400">
                    （未设置答案）
                  </div>
                </div>

                <div
                  v-else-if="selectedQuestion.type === 'uploadImage'"
                  class="preview-upload"
                >
                  图片附件题（学生需上传图片，最多{{ selectedQuestion.options }}份）
                </div>
                <div
                  v-else-if="selectedQuestion.type === 'uploadDocument'"
                  class="preview-upload"
                >
                  文档附件题（学生需上传文档，最多{{ selectedQuestion.options }}份）
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="modal-edit cus-scroll">
          <div class="editor-panel-head">
            题目属性编辑
          </div>
          <div class="modal-panel">
            <div v-if="selectedQuestion">
              <NForm class="settings-form" label-placement="top">
                <NFormItem label="题目名称">
                  <NInput
                    v-model:value="selectedQuestion.title"
                    type="textarea"
                    :autosize="{ minRows: 2, maxRows: 6 }"
                    placeholder="请输入题目名称（支持粘贴题目自动解析，格式：题目A.选项1B.选项2...）"
                    @update:value="updateQuestion('title', $event)"
                    @paste="handlePasteQuestion"
                  />
                </NFormItem>

                <NFormItem label="题目分数">
                  <NInputNumber
                    v-model:value="selectedQuestion.score"
                    :min="1"
                    @update:value="updateQuestion('score', $event)"
                  />
                </NFormItem>

                <template v-if="selectedQuestion.type === 'single' || selectedQuestion.type === 'multiple'">
                  <div class="editor-panel-head">
                    {{ selectedQuestion.type === 'single' ? '单选题' : '多选题' }}选项
                  </div>
                  <div class="option-list">
                    <div
                      v-for="(option, index) in selectedQuestion.options"
                      :key="index"
                      class="option-item"
                    >
                      <div class="option-content">
                        <NInput
                          v-model:value="selectedQuestion.options[index]"
                          placeholder="请输入选项内容"
                        />
                        <NSwitch
                          :value="selectedQuestion.answer.includes(index)"
                          :checked-value="true"
                          :unchecked-value="false"
                          @update:value="updateAnswer(index)"
                        >
                          <template #checked>
                            正确
                          </template>
                          <template #unchecked>
                            错误
                          </template>
                        </NSwitch>
                        <NButton type="warning" text @click="removeOption(index)">
                          <i class="i-material-symbols:close" />
                        </NButton>
                      </div>
                    </div>
                    <NButton secondary block @click="addOption">
                      添加选项
                    </NButton>
                  </div>
                </template>

                <template v-else-if="selectedQuestion.type === 'text'">
                  <div class="editor-panel-head">
                    参考答案
                  </div>
                  <NFormItem :show-label="false">
                    <NInput
                      v-model:value="textAnswer"
                      type="textarea"
                      placeholder="请输入参考答案"
                      :rows="3"
                      @update:value="updateTextAnswer($event)"
                    />
                  </NFormItem>
                </template>

                <template v-else-if="selectedQuestion.type === 'completion'">
                  <NAlert title="请使用连续的_（下划线）来代表一个填空">
                    <template #icon>
                      <i class="i-material-symbols:warning" />
                    </template>
                  </NAlert>
                  <div class="editor-panel-head">
                    填空答案
                  </div>
                  <div class="completion-edit-simple">
                    <div
                      v-for="(blank, blankIndex) in completionBlanks"
                      :key="blankIndex"
                      class="completion-blank-item"
                    >
                      <span class="blank-label">
                        填空{{ blankIndex + 1 }}
                      </span>
                      <NInput
                        v-model:value.trim="blank.answer"
                        placeholder="请输入答案"
                        @update:value="updateCompletionAnswer(blankIndex, $event)"
                      />
                    </div>
                    <div v-if="completionBlanks.length === 0" class="empty-blanks">
                      题目中未找到连续的 _（下划线）
                    </div>
                  </div>
                </template>

                <template v-else-if="selectedQuestion.type === 'uploadImage' || selectedQuestion.type === 'uploadDocument'">
                  <div class="editor-panel-head">
                    上传附件
                  </div>
                  <NFormItem label="最大上传份数">
                    <NInputNumber
                      v-model:value="selectedQuestion.options"
                      :min="1"
                      :max="9"
                    />
                  </NFormItem>
                </template>
              </NForm>
            </div>
          </div>
        </div>
      </div>
    </NModal>
  </CommonPage>
</template>

<script setup>
import {
  NAlert,
  NButton,
  NDropdown,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NModal,
  NSwitch,
  NTag,
} from 'naive-ui'

import { computed, nextTick, onMounted, ref } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api/index.js'
import { CommonPage } from '@/components/index.js'
import { useUnsavedGuard } from '@/composables'
import { useTabStore } from '@/store/index.js'

const questionTypes = [
  {
    name: '单选题',
    type: 'single',
    icon: 'i-material-symbols:radio-button-checked',
    desc: '选择唯一正确答案',
    defaultTitle: '单选题题目',
    options: ['选项1', '选项2'],
    answer: [],
    tagType: 'info',
  },
  {
    name: '多选题',
    type: 'multiple',
    icon: 'i-material-symbols:checklist',
    desc: '选择多个正确答案',
    defaultTitle: '多选题题目',
    options: ['选项1', '选项2'],
    answer: [],
    tagType: 'success',
  },
  {
    name: '填空题',
    type: 'completion',
    icon: 'i-material-symbols:text-fields',
    desc: '候选词填空',
    defaultTitle: '填空题__题目__',
    options: [],
    answer: [],
    tagType: 'warning',
  },
  {
    name: '简答题',
    type: 'text',
    icon: 'i-material-symbols:edit-note',
    desc: '简要回答',
    defaultTitle: '简答题题目',
    options: [],
    answer: '',
    tagType: 'primary',
  },
  {
    name: '上传图片',
    type: 'uploadImage',
    icon: 'i-material-symbols:image-arrow-up-outline-rounded',
    desc: '学生上传图片',
    defaultTitle: '上传图片题目',
    options: 1,
    answer: [],
    tagType: 'info',
  },
  {
    name: '上传文档',
    type: 'uploadDocument',
    icon: 'i-material-symbols:upload-file',
    desc: '学生上传文档',
    defaultTitle: '上传文档题目',
    options: 1,
    answer: [],
    tagType: 'default',
  },
]

function getQuestionTypeMeta(type) {
  const found = questionTypes.find(item => item.type === type)
  return found || { name: '未知题型', tagType: 'default' }
}

const selectIndex = ref(0)
const questions = ref([])
const route = useRoute()
const router = useRouter()
const tabStore = useTabStore()
const chapterId = ref(route.query.chapterId)

const totalScore = computed(() => {
  return questions.value.reduce((total, question) => {
    return total + question.score
  }, 0)
})

const passScore = ref(0)
const saving = ref(false)
const { snapshot, allowLeave } = useUnsavedGuard(() => ({
  questions: questions.value,
  passScore: passScore.value,
}))

function generateQuestionKey(question) {
  const content = question.title || ''
  const random = Math.random().toString(36).substring(2, 9)
  return `${content}_${random}_${Date.now()}`
}

function ensureQuestionKey(question) {
  if (!question._uniqueKey) {
    question._uniqueKey = generateQuestionKey(question)
  }
  return question
}

onMounted(async () => {
  if (route.params.id) {
    const { data } = await api.paper.paperDetail(route.params.id)
    questions.value = data.questions.map(q => ensureQuestionKey(q))
    passScore.value = data.pass_score
  }
  nextTick(snapshot)
})

async function handleSave() {
  if (saving.value)
    return
  if (questions.value.length === 0) {
    return window.$message.error('请配置题目')
  }
  if (passScore.value > totalScore.value) {
    return window.$message.error('分数线不能大于总分')
  }
  if (passScore.value < 0) {
    return window.$message.error('分数线不能小于0')
  }
  if (passScore.value === 0) {
    return window.$message.error('分数线未设置')
  }
  let hasError = false
  questions.value.forEach((question, index) => {
    const tip = `第${index + 1}题，`
    if (!question.title) {
      hasError = true
      return window.$message.error(`${tip}请填写题目`)
    }
    if (!question.score) {
      hasError = true
      return window.$message?.error(`${tip}请填写题目分数`)
    }

    switch (question.type) {
      case 'single':
      case 'multiple':
        if (!question.options.length) {
          hasError = true
          return window.$message?.error(`${tip}请填写题目选项`)
        }
        if (!question.answer.length) {
          hasError = true
          return window.$message?.error(`${tip}请填写题目答案`)
        }
        break
      case 'text': {
        if (!question.answer || typeof question.answer !== 'string' || !question.answer.trim()) {
          hasError = true
          return window.$message?.error(`${tip}请填写参考答案`)
        }
        break
      }
      case 'completion': {
        if (!question.title.includes('_')) {
          hasError = true
          return window.$message?.error(`${tip}请使用连续的下划线_设置填空`)
        }
        const blankMatches = question.title.match(/_+/g) || []
        const blankCount = blankMatches.length
        if (!Array.isArray(question.answer) || question.answer.length !== blankCount) {
          hasError = true
          return window.$message?.error(`${tip}答案数量与填空数量不匹配`)
        }
        if (question.answer.some(ans => !ans || ans.trim() === '')) {
          hasError = true
          return window.$message?.error(`${tip}请填写所有填空的答案`)
        }
        break
      }
    }
  })
  if (hasError)
    return
  saving.value = true
  try {
    await api.paper.paperSave({
      questions: questions.value,
      pass_score: passScore.value,
      score: totalScore.value,
      id: route.params.id,
      course_id: route.params.courseId,
      chapter_id: route.query.chapterId,
    })
    window.$message.success(route.params.id ? '更新成功' : '创建成功')
    allowLeave()
    const currentPath = route.path
    const hadOthers = tabStore.tabs.some(tab => tab.path !== currentPath)
    await tabStore.removeTab(currentPath)
    if (!hadOthers) {
      if (window.history.length > 1)
        router.back()
      else
        router.push('/')
    }
  }
  catch {
  }
  finally {
    saving.value = false
  }
}

function handleDelete() {
  api.paper.paperDelete(route.params.id).then(async () => {
    window.$message.success('删除成功')
    allowLeave()
    const currentPath = route.path
    const hadOthers = tabStore.tabs.some(tab => tab.path !== currentPath)
    await tabStore.removeTab(currentPath)
    if (!hadOthers) {
      if (window.history.length > 1)
        router.back()
      else
        router.push('/')
    }
  })
}

const moreOptions = [
  { label: chapterId.value ? '删除测评' : '删除试卷', key: 'delete' },
]

function handleMoreSelect(key) {
  if (key === 'delete')
    handleDelete()
}
const selectedQuestion = ref(null)

const editModal = ref({
  show: false,
  title: '编辑题目',
  questionIndex: -1,
})

const textAnswer = computed({
  get() {
    if (!selectedQuestion.value || selectedQuestion.value.type !== 'text') {
      return ''
    }

    if (typeof selectedQuestion.value.answer === 'object' && selectedQuestion.value.answer !== null) {
      return selectedQuestion.value.answer.referenceAnswer || ''
    }
    return selectedQuestion.value.answer || ''
  },
  set(value) {
    if (selectedQuestion.value && selectedQuestion.value.type === 'text') {
      selectedQuestion.value.answer = value
    }
  },
})

const completionBlanks = computed(() => {
  if (!selectedQuestion.value || selectedQuestion.value.type !== 'completion') {
    return []
  }

  const title = selectedQuestion.value.title || ''
  const answer = selectedQuestion.value.answer || []
  const blanks = []
  let currentBlank = null
  let blankIndex = 0

  for (let i = 0; i < title.length; i++) {
    if (title[i] === '_') {
      if (!currentBlank) {
        currentBlank = {
          startIndex: i,
          endIndex: i,
          answer: answer[blankIndex] || '',
          blankIndex,
        }
      }
      else {
        currentBlank.endIndex = i
      }
    }
    else {
      if (currentBlank) {
        blanks.push(currentBlank)
        blankIndex++
        currentBlank = null
      }
    }
  }

  if (currentBlank) {
    blanks.push(currentBlank)
  }

  return blanks
})

function addQuestion(type) {
  let answer
  if (type.type === 'text') {
    answer = typeof type.answer === 'string' ? type.answer : ''
  }
  else if (Array.isArray(type.answer)) {
    answer = [...type.answer]
  }
  else {
    answer = type.answer
  }

  const newQuestion = {
    type: type.type,
    title: type.defaultTitle,
    score: 5,
    options: Array.isArray(type.options) ? [...type.options] : type.options,
    answer,
  }
  ensureQuestionKey(newQuestion)

  const newIndex = questions.value.length
  questions.value.push(newQuestion)

  openEditModal(newIndex)
}

function setSelectIndex(index) {
  selectIndex.value = index
  selectedQuestion.value = questions.value[index] || null
}

function openEditModal(index) {
  if (index >= 0 && index < questions.value.length) {
    setSelectIndex(index)
    const typeName = getQuestionTypeMeta(questions.value[index]?.type).name
    editModal.value.questionIndex = index
    editModal.value.title = `编辑题目 - 第${index + 1}题（${typeName}）`
    editModal.value.show = true
  }
}

function confirmRemoveQuestion(index) {
  if (index < 0 || index >= questions.value.length)
    return
  const question = questions.value[index]
  const typeName = getQuestionTypeMeta(question?.type).name
  const title = (question?.title || '').trim() || '未命名题目'
  $dialog.warning({
    title: '删除题目',
    content: `确定删除第 ${index + 1} 题（${typeName}）「${title}」吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick() {
      removeQuestion(index)
    },
  })
}

function removeQuestion(index) {
  if (index < 0 || index >= questions.value.length)
    return
  const deletingOpen = editModal.value.show && editModal.value.questionIndex === index
  questions.value.splice(index, 1)
  if (deletingOpen) {
    editModal.value.show = false
  }
  else if (editModal.value.show && editModal.value.questionIndex > index) {
    editModal.value.questionIndex -= 1
    editModal.value.title = `编辑题目 - 第${editModal.value.questionIndex + 1}题（${getQuestionTypeMeta(questions.value[editModal.value.questionIndex]?.type).name}）`
  }
  if (questions.value.length > 0) {
    const newIndex = Math.min(index, questions.value.length - 1)
    setSelectIndex(newIndex)
  }
  else {
    selectIndex.value = -1
    selectedQuestion.value = null
  }
}

function getCompletionAnswers(question) {
  if (!question || question.type !== 'completion') {
    return []
  }
  const title = question.title || ''
  const answer = question.answer || []
  const answers = []
  let blankIndex = 0
  let inBlank = false

  for (let i = 0; i < title.length; i++) {
    if (title[i] === '_') {
      if (!inBlank) {
        answers.push(answer[blankIndex] || '')
        blankIndex++
        inBlank = true
      }
    }
    else {
      inBlank = false
    }
  }

  return answers
}

function getCompletionPreviewParts(question) {
  if (!question || question.type !== 'completion') {
    return []
  }
  const title = question.title || ''
  const answer = question.answer || []
  const parts = []
  let currentText = ''
  let blankIndex = 0
  let inBlank = false

  for (let i = 0; i < title.length; i++) {
    if (title[i] === '_') {
      if (!inBlank) {
        if (currentText) {
          parts.push({ isBlank: false, text: currentText })
          currentText = ''
        }
        parts.push({
          isBlank: true,
          answer: answer[blankIndex] || '',
          blankIndex,
        })
        blankIndex++
        inBlank = true
      }
    }
    else {
      if (inBlank) {
        inBlank = false
      }
      currentText += title[i]
    }
  }

  if (currentText) {
    parts.push({ isBlank: false, text: currentText })
  }

  return parts
}

function updateCompletionAnswer(blankIndex, value) {
  if (!selectedQuestion.value || selectedQuestion.value.type !== 'completion') {
    return
  }

  if (!Array.isArray(selectedQuestion.value.answer)) {
    selectedQuestion.value.answer = []
  }

  selectedQuestion.value.answer[blankIndex] = value.trim() || ''
}

function parsePastedQuestion(text) {
  if (!text || typeof text !== 'string') {
    return null
  }

  text = text.trim()

  text = text.replace(/^[\d一二三四五六七八九十]+[.、）)]\s*/, '')

  const optionMarkers = []
  const optionPattern = /([A-Z])[.、）)]/g
  let match = optionPattern.exec(text)
  while (match !== null) {
    optionMarkers.push({
      letter: match[1],
      position: match.index,
      fullMatchLength: match[0].length,
    })
    match = optionPattern.exec(text)
  }

  if (optionMarkers.length === 0) {
    return null
  }

  const questionText = text.substring(0, optionMarkers[0].position).trim()

  const options = []
  for (let i = 0; i < optionMarkers.length; i++) {
    const currentMarker = optionMarkers[i]
    const nextMarker = optionMarkers[i + 1]

    const contentStart = currentMarker.position + currentMarker.fullMatchLength

    const contentEnd = nextMarker ? nextMarker.position : text.length

    const optionText = text.substring(contentStart, contentEnd).trim()

    if (optionText) {
      options.push(optionText)
    }
  }

  if (options.length < 2) {
    return null
  }

  return {
    question: questionText,
    options,
  }
}

function handlePasteQuestion(event) {
  if (!selectedQuestion.value || (selectedQuestion.value.type !== 'single' && selectedQuestion.value.type !== 'multiple')) {
    return
  }

  const pastedText = (event.clipboardData || window.clipboardData).getData('text')

  if (!pastedText) {
    return
  }

  const parsed = parsePastedQuestion(pastedText)

  if (parsed && parsed.question && parsed.options.length > 0) {
    event.preventDefault()

    selectedQuestion.value.title = parsed.question

    selectedQuestion.value.options = parsed.options.filter(opt => opt)

    selectedQuestion.value.answer = []

    window.$message.success(`已自动解析：题目和${selectedQuestion.value.options.length}个选项`)
  }
}

function updateQuestion(field, value) {
  if (selectedQuestion.value) {
    if (field === 'title') {
      if (selectedQuestion.value.type === 'completion') {
        const title = value || ''
        const blankMatches = title.match(/_+/g) || []
        const blankCount = blankMatches.length
        if (!Array.isArray(selectedQuestion.value.answer) || selectedQuestion.value.answer.length !== blankCount) {
          selectedQuestion.value.answer = Array.from({ length: blankCount }, () => '')
        }
      }
      selectedQuestion.value[field] = value
    }
  }
}

function addOption() {
  selectedQuestion.value.options.push('')
}

function removeOption(index) {
  if (selectedQuestion.value && selectedQuestion.value.options) {
    selectedQuestion.value.options.splice(index, 1)
    if (selectedQuestion.value.answer.includes(index)) {
      selectedQuestion.value.answer = selectedQuestion.value.answer.filter(item => item !== index)
    }
    selectedQuestion.value.answer = selectedQuestion.value.answer.map((item) => {
      if (item > index) {
        return item - 1
      }
      return item
    })
  }
}

function updateAnswer(index) {
  if (selectedQuestion.value && selectedQuestion.value.answer) {
    if (selectedQuestion.value.type === 'single') {
      if (selectedQuestion.value.answer.includes(index)) {
        selectedQuestion.value.answer = []
      }
      else {
        selectedQuestion.value.answer = [index]
      }
    }
    if (selectedQuestion.value.type === 'multiple') {
      if (selectedQuestion.value.answer.includes(index)) {
        selectedQuestion.value.answer.splice(selectedQuestion.value.answer.indexOf(index), 1)
      }
      else {
        selectedQuestion.value.answer.push(index)
      }
    }
  }
}

function updateTextAnswer(value) {
  if (!selectedQuestion.value || selectedQuestion.value.type !== 'text') {
    return
  }
  selectedQuestion.value.answer = value || ''
}
</script>

<style lang="scss" scoped>
.paper-editor-main {
  display: flex;
  flex: 1;
  padding: 0;
  gap: 16px;
  min-height: 0;
  height: 100%;
}

.paper-editor-left {
  width: 248px;
  flex-shrink: 0;
  min-height: 0;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.question-type-list {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  overflow: auto;
  gap: 8px;
  padding: 8px;
  border-radius: 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
}

.question-type-item {
  display: flex;
  align-items: center;
  padding: 12px;
  border: 0;
  border-radius: 11px;
  cursor: pointer;
  transition:
    transform 0.22s var(--ease-apple, cubic-bezier(0.32, 0.72, 0, 1)),
    background 0.2s ease;
  background: rgb(255 255 255 / 55%);
  box-shadow: none;
}

.question-type-item:hover {
  background: rgb(255 255 255 / 78%);
  transform: none;
  box-shadow: none;
  border-color: transparent;
}

.question-type-item:active {
  transform: scale(0.985);
}

.question-type-icon {
  width: 40px;
  height: 40px;
  border-radius: 11px;
  background: rgba(var(--primary-color), 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  flex-shrink: 0;
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 45%);
  border: 1px solid rgba(var(--primary-color), 0.12);
}

.question-type-icon i {
  font-size: 20px;
  color: rgb(var(--primary-color));
}

.question-type-content {
  flex: 1;
}

.question-type-name {
  font-size: 15px;
  font-weight: 560;
  letter-spacing: -0.01em;
  margin-bottom: 4px;
  color: var(--text-primary, #1a2332);
}

.question-type-desc {
  font-size: 12px;
  letter-spacing: -0.01em;
  color: var(--text-secondary, #5b6b7c);
}

.paper-editor-center {
  flex: 1;
  min-width: 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.question-item {
  position: relative;
  border: 0;
  border-radius: 12px;
  margin-bottom: 8px;
  padding: 12px 40px 12px 14px;
  cursor: pointer;
  background: rgb(255 255 255 / 55%);
  box-shadow: none;
  transition:
    transform 0.22s var(--ease-apple, cubic-bezier(0.32, 0.72, 0, 1)),
    background 0.2s ease;
}

.question-item:hover {
  background: rgb(255 255 255 / 78%);
  box-shadow: none;
}

.question-item:active {
  transform: scale(0.985);
}

.question-item.active {
  background: rgba(var(--primary-color), 0.1);
  box-shadow: none;
}

.question-item.ghost {
  opacity: 0.55;
  background: rgba(var(--primary-color), 0.05);
}

.question-item.chosen {
  border-color: rgba(var(--primary-color), 0.28);
  background: rgba(var(--primary-color), 0.08);
  box-shadow:
    inset 0 1px 0 rgb(255 255 255 / 50%),
    0 0 0 3px rgba(var(--primary-color), 0.12);
}

.question-item-header {
  display: flex;
  align-items: flex-start;
  margin-bottom: 12px;
}

.question-item-title {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  min-width: 0;
  flex: 1;
  padding-right: 4px;
}

.question-item-title h4 {
  margin: 0;
  font-weight: 560;
  letter-spacing: -0.01em;
  line-height: 1.45;
  word-break: break-word;
}

.question-item-delete {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 2;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: 0;
  border-radius: 8px;
  background: rgb(255 255 255 / 88%);
  color: #d03050;
  cursor: pointer;
  font-size: 16px;
  box-shadow: 0 1px 4px rgb(0 0 0 / 8%);
  opacity: 0.72;
  transition: background 0.15s ease, color 0.15s ease, opacity 0.15s ease;
}

.question-item:hover .question-item-delete,
.question-item.active .question-item-delete {
  opacity: 1;
}

.question-item-delete:hover {
  background: rgba(208, 48, 80, 0.12);
  color: #b91c3c;
}

.question-score {
  font-size: 13px;
  letter-spacing: -0.01em;
  color: var(--text-secondary, #5b6b7c);
  margin-bottom: 8px;
}

.question-options {
  margin-top: 8px;
}

.question-option {
  padding: 4px 0;
}

.option-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.option-item {
  margin: 0;
}

.option-content {
  display: flex;
  gap: 12px;
  align-items: center;
}

.option-content :deep(.n-input) {
  flex: 1;
}

.empty-state {
  text-align: center;
  padding: 56px 20px;
  letter-spacing: -0.012em;
  font-weight: 500;
  color: var(--ctl-text-soft, #5b6b7c);
  border-radius: 16px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 40%);
}

.questions-container {
  min-height: 220px;
  flex: 1;
  overflow: auto;
  padding: 4px;
  border-radius: 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
}

.draggable-container {
  min-height: 100px;
}

.completion-content {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 8px;
}

.blank-display-inline {
  display: flex;
  align-items: center;
  gap: 4px;
}

.correct-answer {
  background: rgba(var(--primary-color), 0.12);
  padding: 2px 8px;
  border-radius: 999px;
  color: rgb(var(--primary-color));
  font-weight: 550;
  letter-spacing: -0.01em;
}

.question-edit-modal {
  :deep(.n-card) {
    height: 100%;
  }
}

.modal-content-wrapper {
  display: flex;
  gap: 16px;
  height: calc(80vh - 100px);
  min-height: 500px;
}

.modal-preview,
.modal-edit {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-panel {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 8px 12px 16px;
  border-radius: 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
}

.question-preview {
  padding: 8px 4px 12px;
}

.preview-header {
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
}

.preview-title-row {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-bottom: 8px;
}

.preview-header h3 {
  margin: 0;
  font-family: var(--font-display);
  font-size: 17px;
  font-weight: 600;
  letter-spacing: -0.022em;
  line-height: 1.45;
  word-break: break-word;
}

.preview-score {
  font-size: 13px;
  letter-spacing: -0.01em;
  color: var(--text-secondary, #5b6b7c);
}

.preview-content {
  margin-top: 16px;
}

.preview-options {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.preview-option {
  padding: 10px 12px;
  border: 0;
  border-radius: 10px;
  background: rgb(255 255 255 / 55%);
  box-shadow: none;
  transition:
    background 0.2s ease,
    transform 0.2s cubic-bezier(0.32, 0.72, 0, 1);
}

.preview-option.correct-answer {
  background: rgba(var(--primary-color), 0.1);
  color: rgb(var(--primary-color));
  font-weight: 550;
}

.preview-answer {
  padding: 12px 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  border-radius: 12px;
  min-height: 60px;
  letter-spacing: -0.01em;
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 40%);
}

.completion-preview {
  line-height: 2;
  font-size: 16px;
}

.blank-preview {
  display: inline-block;
  min-width: 30px;
  padding: 2px 8px;
  margin: 0 2px;
  background: rgba(var(--primary-color), 0.1);
  border: 1px solid rgba(var(--primary-color), 0.18);
  border-radius: 8px;
  color: rgb(var(--primary-color));
  font-weight: 600;
  text-align: center;
}

.text-answer-preview {
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.6;
}

.preview-label {
  font-size: 13px;
  letter-spacing: -0.01em;
  color: var(--text-secondary, #5b6b7c);
  margin-bottom: 8px;
  font-weight: 550;
}

.preview-content {
  padding: 12px 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  border-radius: 12px;
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.6;
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 40%);
}

.preview-upload {
  padding: 12px 14px;
  background: rgba(var(--primary-color), 0.1);
  border: 1px solid rgba(var(--primary-color), 0.22);
  border-radius: 12px;
  color: rgb(var(--primary-color));
  letter-spacing: -0.01em;
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 40%);
}

.completion-edit-simple {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 4px 0 8px;
}

.completion-blank-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.completion-blank-item .blank-label {
  min-width: 80px;
  font-size: 13px;
  letter-spacing: -0.01em;
  color: var(--text-secondary, #5b6b7c);
  font-weight: 550;
}

.completion-blank-item :deep(.n-input) {
  flex: 1;
}

.empty-blanks {
  padding: 20px;
  text-align: center;
  letter-spacing: -0.01em;
  color: var(--text-secondary, #5b6b7c);
}
</style>
