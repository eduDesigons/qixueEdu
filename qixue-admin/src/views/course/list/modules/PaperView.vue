<template>
  <div class="paper-view">
    <div v-if="paperId" class="cus-scroll paper-view__body">
      <div class="editor-panel-head">
        <span>题目列表</span>
        <div class="editor-panel-head__extra">
          <span>分数线</span>
          <span class="paper-view__score">
            <NNumberAnimation :from="0" :to="paperInfo.pass_score" />
          </span>
          <span>/ 总分</span>
          <span class="paper-view__score">
            <NNumberAnimation :from="0" :to="totalScore" />
          </span>
        </div>
      </div>

      <div v-if="questions.length === 0" class="empty-state">
        <i class="i-material-symbols:menu-outline-rounded paper-view__empty-ico" />
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
        >
          <div
            v-for="(question, index) in questions"
            :key="question._uniqueKey || index"
            class="question-item"
          >
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
                class="question-upload"
              >
                图片附件题（学生需上传图片）
              </div>
              <div
                v-else-if="question.type === 'uploadDocument'"
                class="question-upload"
              >
                文档附件题（学生需上传文档）
              </div>
            </div>
          </div>
        </VueDraggable>
      </div>
    </div>
    <div v-else class="empty-state">
      还未设置题目
    </div>
  </div>
</template>

<script setup>
import { NNumberAnimation, NTag } from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'
import api from '@/api/index.js'

const QUESTION_TYPE_META = {
  single: { name: '单选题', tagType: 'info' },
  multiple: { name: '多选题', tagType: 'success' },
  completion: { name: '填空题', tagType: 'warning' },
  text: { name: '简答题', tagType: 'primary' },
  uploadImage: { name: '上传图片', tagType: 'info' },
  uploadDocument: { name: '上传文档', tagType: 'default' },
}

function getQuestionTypeMeta(type) {
  return QUESTION_TYPE_META[type] || { name: '未知题型', tagType: 'default' }
}

const props = defineProps({
  paperId: {
    type: [String, Number],
    default: null,
  },
})

const questions = ref([])
const paperInfo = ref({})

async function loadPaper() {
  if (!props.paperId) {
    questions.value = []
    paperInfo.value = {}
    return
  }
  try {
    const { data } = await api.paper.paperDetail(props.paperId)
    questions.value = data.questions || []
    paperInfo.value = data || {}
  }
  catch (error) {
    $message.error(error?.message || '加载试卷失败')
  }
}

watch(() => props.paperId, () => {
  loadPaper()
}, { immediate: true })

defineExpose({
  reload: loadPaper,
})

const totalScore = computed(() => {
  return questions.value.reduce((total, question) => {
    return total + (question.score || 0)
  }, 0)
})

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
</script>

<style lang="scss" scoped>
.paper-view {
  display: flex;
  flex: 1;
  min-height: 0;
  padding: 4px 0 8px;
}

.paper-view__body {
  flex: 1;
  min-width: 0;
  overflow-y: auto;
}

.paper-view__score {
  font-weight: 600;
  color: rgb(var(--primary-color));
}

.paper-view__empty-ico {
  font-size: 48px;
  margin-bottom: 16px;
}

.questions-container {
  min-height: 220px;
  padding: 4px;
  border-radius: 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
}

.draggable-container {
  min-height: 100px;
}

.question-item {
  border: 0;
  border-radius: 12px;
  margin-bottom: 8px;
  padding: 12px 14px;
  cursor: grab;
  background: rgb(255 255 255 / 55%);
  box-shadow: none;
  transition:
    transform 0.22s var(--ease-apple, cubic-bezier(0.32, 0.72, 0, 1)),
    background 0.2s ease;
}

.question-item:hover {
  background: rgb(255 255 255 / 78%);
}

.question-item:active {
  cursor: grabbing;
  transform: scale(0.985);
}

.question-item.ghost {
  opacity: 0.55;
  background: rgba(var(--primary-color), 0.05);
}

.question-item.chosen {
  background: rgba(var(--primary-color), 0.08);
  box-shadow:
    inset 0 1px 0 rgb(255 255 255 / 50%),
    0 0 0 3px rgba(var(--primary-color), 0.12);
}

.question-item-header {
  display: flex;
  align-items: flex-start;
  margin-bottom: 10px;
}

.question-item-title {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  min-width: 0;
  flex: 1;
}

.question-item-title h4 {
  margin: 0;
  font-family: var(--font-display);
  font-size: 15px;
  font-weight: 560;
  letter-spacing: -0.015em;
  line-height: 1.45;
  word-break: break-word;
  color: var(--ctl-text, #1a2332);
}

.question-score {
  font-size: 12px;
  font-weight: 550;
  letter-spacing: -0.01em;
  color: var(--ctl-text-soft, #5b6b7c);
  margin-bottom: 8px;
}

.question-options {
  margin-top: 8px;
}

.question-option {
  padding: 6px 0;
  letter-spacing: -0.01em;
  color: var(--ctl-text, #1a2332);
}

.question-answer {
  margin-top: 8px;
  line-height: 1.6;
  letter-spacing: -0.01em;
  color: var(--ctl-text-soft, #5b6b7c);
}

.question-upload {
  margin-top: 8px;
  padding: 10px 12px;
  background: rgba(var(--primary-color), 0.1);
  border: 1px solid rgba(var(--primary-color), 0.22);
  border-radius: 12px;
  letter-spacing: -0.01em;
  color: rgb(var(--primary-color));
}
</style>
