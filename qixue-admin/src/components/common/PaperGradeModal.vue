<template>
  <NModal
    v-model:show="visible"
    preset="card"
    :title="readOnly ? '查看答题详情' : '批改试卷'"
    class="w-800px max-h-90vh"
    style="max-height: 90vh;"
    @after-leave="handleClose"
  >
    <div v-if="loading" class="py-40 text-center">
      <NSpin size="medium" />
    </div>
    <div v-else-if="detail" ref="gradeBody" class="grade-content overflow-y-auto" style="max-height: calc(90vh - 140px);">
      <div class="grade-meta info-block">
        <div class="grade-meta__main">
          <div class="entity-card__badges">
            <span
              class="entity-card__badge"
              :class="detail.user_paper?.is_people_check ? 'is-warn' : 'is-idle'"
            >
              {{ detail.user_paper?.is_people_check ? '待批改' : '已批改' }}
            </span>
            <span
              class="entity-card__badge"
              :class="detail.user_paper?.is_pass ? 'is-ok' : 'is-bad'"
            >
              {{ detail.user_paper?.is_pass ? '通过' : '未通过' }}
            </span>
          </div>
          <div class="grade-meta__ids">
            <button
              v-if="detail.user_paper?.user_id"
              type="button"
              class="entity-link"
              @click="goUser"
            >
              {{ detail.user_paper?.user?.nick_name || '学员' }}
            </button>
            <span v-else class="font-medium">{{ detail.user_paper?.user?.nick_name || '学员' }}</span>
            <button
              v-if="detail.paper?.course?.title"
              type="button"
              class="entity-card__chip"
              @click="goCourse"
            >
              {{ detail.paper.course.title }}
            </button>
            <button
              v-if="detail.paper?.chapter?.title"
              type="button"
              class="entity-card__chip"
              @click="goChapter"
            >
              {{ detail.paper.chapter.title }}
            </button>
          </div>
        </div>
        <div class="entity-card__metrics">
          <div class="entity-card__metric">
            <strong>{{ detail.user_paper?.score ?? '—' }}</strong>
            <span>成绩 / {{ detail.paper?.score ?? '—' }}</span>
          </div>
          <div class="entity-card__metric">
            <strong>{{ detail.paper?.pass_score ?? '—' }}</strong>
            <span>及格分</span>
          </div>
        </div>
      </div>
      <div class="grade-list">
      <div
        v-for="(q, idx) in detail.questions"
        :key="q.id"
        class="question-item"
        :class="{ 'question-item--need': !readOnly && q.user_answer?.need_people_check }"
      >
        <div class="mb-8 flex items-start justify-between">
          <h4 class="font-medium">
            {{ idx + 1 }}. {{ q.title }}
            <NTag v-if="!readOnly && q.user_answer?.need_people_check" type="warning" size="small" round :bordered="false" class="ml-8">需批改</NTag>
          </h4>
          <span class="text-12 grade-muted">分值：{{ q.score }}分</span>
        </div>
        <!-- 选项题 -->
        <div v-if="q.type === 'single' || q.type === 'multiple'" class="question-options text-14">
          <div v-for="(opt, oi) in (q.options || [])" :key="oi" class="mb-4">
            {{ String.fromCharCode(65 + oi) }}. {{ opt }}
            <NTag v-if="isUserAnswer(q, oi)" type="info" size="small" round :bordered="false" class="ml-8">用户选择</NTag>
            <NTag v-if="isCorrectAnswer(q, oi)" type="success" size="small" round :bordered="false" class="ml-4">正确答案</NTag>
          </div>
          <div class="mt-8 grade-muted">
            用户答案：{{ formatUserOptionAnswer(q) }}
          </div>
          <div class="grade-score-row grade-score-row--readonly">
            <span class="grade-score-row__label">得分</span>
            <strong>{{ q.user_answer?.score ?? 0 }}</strong>
            <span class="grade-score-row__unit">/ {{ q.score }} 分</span>
          </div>
        </div>
        <!-- 填空题 -->
        <div v-else-if="q.type === 'completion'" class="text-14">
          <div class="mb-4">
            <span class="grade-muted">用户答案：</span>
            <span>{{ formatArrayAnswer(q.user_answer?.answer_content) }}</span>
          </div>
          <div v-if="q.answer?.length" class="mb-4 grade-muted">
            参考答案：{{ formatArrayAnswer(q.answer) }}
          </div>
          <div v-if="!readOnly && q.user_answer?.need_people_check" class="grade-score-row">
            <span class="grade-score-row__label">得分</span>
            <NInputNumber
              v-model:value="scoreInputs[q.id]"
              :min="0"
              :max="q.score"
              size="small"
              class="grade-score-input"
            />
            <span class="grade-score-row__unit">/ {{ q.score }} 分</span>
            <button type="button" class="grade-score-chip" @click="scoreInputs[q.id] = q.score">满分</button>
            <button type="button" class="grade-score-chip" @click="scoreInputs[q.id] = 0">零分</button>
          </div>
          <div v-else class="grade-score-row grade-score-row--readonly">
            <span class="grade-score-row__label">得分</span>
            <strong>{{ q.user_answer?.score ?? 0 }}</strong>
            <span class="grade-score-row__unit">/ {{ q.score }} 分</span>
          </div>
        </div>
        <!-- 主观题 -->
        <div v-else-if="q.type === 'text'" class="text-14">
          <div class="mb-4">
            <span class="grade-muted">用户答案：</span>
            <div class="mt-4 grade-answer-well">{{ formatArrayAnswer(q.user_answer?.answer_content) || '（未作答）' }}</div>
          </div>
          <div v-if="q.answer?.length" class="mb-4 grade-muted">
            参考答案：{{ Array.isArray(q.answer) ? q.answer.join(' / ') : q.answer }}
          </div>
          <div v-if="!readOnly && q.user_answer?.need_people_check" class="grade-score-row">
            <span class="grade-score-row__label">得分</span>
            <NInputNumber
              v-model:value="scoreInputs[q.id]"
              :min="0"
              :max="q.score"
              size="small"
              class="grade-score-input"
            />
            <span class="grade-score-row__unit">/ {{ q.score }} 分</span>
            <button type="button" class="grade-score-chip" @click="scoreInputs[q.id] = q.score">满分</button>
            <button type="button" class="grade-score-chip" @click="scoreInputs[q.id] = 0">零分</button>
          </div>
          <div v-else class="grade-score-row grade-score-row--readonly">
            <span class="grade-score-row__label">得分</span>
            <strong>{{ q.user_answer?.score ?? 0 }}</strong>
            <span class="grade-score-row__unit">/ {{ q.score }} 分</span>
          </div>
        </div>
        <!-- 上传图片题 -->
        <div v-else-if="q.type === 'uploadImage'" class="text-14">
          <div class="mb-4">
            <span class="grade-muted">用户上传：</span>
            <div v-if="getUploadUrls(q).length" class="mt-4 flex flex-wrap gap-12">
              <NImage
                v-for="(url, i) in getUploadUrls(q)"
                :key="i"
                :src="getAssetUrl(url)"
                class="h-80 w-80 shrink-0 rounded-4 object-cover cursor-pointer"
                :img-props="{ style: 'border-radius: 4px; object-fit: cover;' }"
              >
                <template #placeholder>
                  <div class="h-80 w-80 flex items-center justify-center grade-thumb-ph rounded-12">
                    <NSpin size="small" />
                  </div>
                </template>
                <template #error>
                  <div class="h-80 w-80 flex items-center justify-center grade-thumb-ph rounded-12 text-12">
                    加载失败
                  </div>
                </template>
              </NImage>
            </div>
            <div v-else class="mt-4 grade-muted-soft">（未上传）</div>
          </div>
          <div v-if="!readOnly && q.user_answer?.need_people_check" class="grade-score-row">
            <span class="grade-score-row__label">得分</span>
            <NInputNumber
              v-model:value="scoreInputs[q.id]"
              :min="0"
              :max="q.score"
              size="small"
              class="grade-score-input"
            />
            <span class="grade-score-row__unit">/ {{ q.score }} 分</span>
            <button type="button" class="grade-score-chip" @click="scoreInputs[q.id] = q.score">满分</button>
            <button type="button" class="grade-score-chip" @click="scoreInputs[q.id] = 0">零分</button>
          </div>
          <div v-else class="grade-score-row grade-score-row--readonly">
            <span class="grade-score-row__label">得分</span>
            <strong>{{ q.user_answer?.score ?? 0 }}</strong>
            <span class="grade-score-row__unit">/ {{ q.score }} 分</span>
          </div>
        </div>
        <!-- 上传文档题 -->
        <div v-else-if="q.type === 'uploadDocument'" class="text-14">
          <div class="mb-4">
            <span class="grade-muted">用户上传：</span>
            <div v-if="getUploadUrls(q).length" class="mt-4 flex flex-col gap-8">
              <div
                v-for="(url, i) in getUploadUrls(q)"
                :key="i"
                class="grade-file-row flex items-center gap-8 px-12 py-8"
              >
                <i class="i-material-symbols:description text-20 text-primary" />
                <span class="min-w-0 flex-1 truncate">{{ getFileName(url, i) }}</span>
                <NButton size="tiny" type="primary" secondary @click="openFilePreview(url, i)">预览</NButton>
                <NButton
                  size="tiny"
                  tag="a"
                  :href="getAssetUrl(url)"
                  target="_blank"
                  rel="noopener" referrerpolicy="origin"
                  secondary
                >
                  下载
                </NButton>
              </div>
            </div>
            <div v-else class="mt-4 grade-muted-soft">（未上传）</div>
          </div>
          <div v-if="!readOnly && q.user_answer?.need_people_check" class="grade-score-row">
            <span class="grade-score-row__label">得分</span>
            <NInputNumber
              v-model:value="scoreInputs[q.id]"
              :min="0"
              :max="q.score"
              size="small"
              class="grade-score-input"
            />
            <span class="grade-score-row__unit">/ {{ q.score }} 分</span>
            <button type="button" class="grade-score-chip" @click="scoreInputs[q.id] = q.score">满分</button>
            <button type="button" class="grade-score-chip" @click="scoreInputs[q.id] = 0">零分</button>
          </div>
          <div v-else class="grade-score-row grade-score-row--readonly">
            <span class="grade-score-row__label">得分</span>
            <strong>{{ q.user_answer?.score ?? 0 }}</strong>
            <span class="grade-score-row__unit">/ {{ q.score }} 分</span>
          </div>
        </div>
        <!-- 其他题型 -->
        <div v-else class="text-14 grade-muted">
          <div>用户答案：{{ formatArrayAnswer(q.user_answer?.answer_content) || '-' }}</div>
          <div v-if="!readOnly && q.user_answer?.need_people_check" class="grade-score-row">
            <span class="grade-score-row__label">得分</span>
            <NInputNumber
              v-model:value="scoreInputs[q.id]"
              :min="0"
              :max="q.score"
              size="small"
              class="grade-score-input"
            />
            <span class="grade-score-row__unit">/ {{ q.score }} 分</span>
            <button type="button" class="grade-score-chip" @click="scoreInputs[q.id] = q.score">满分</button>
            <button type="button" class="grade-score-chip" @click="scoreInputs[q.id] = 0">零分</button>
          </div>
          <div v-else class="grade-score-row grade-score-row--readonly">
            <span class="grade-score-row__label">得分</span>
            <strong>{{ q.user_answer?.need_people_check ? (scoreInputs[q.id] ?? 0) : (q.user_answer?.score ?? 0) }}</strong>
            <span class="grade-score-row__unit">/ {{ q.score }} 分</span>
          </div>
        </div>
        <div v-if="q.user_answer?.need_people_check" class="mt-12">
          <template v-if="!readOnly">
            <div class="mb-8 grade-muted">批改备注</div>
            <NInput
              v-model:value="remarkInputs[q.id]"
              type="textarea"
              :autosize="{ minRows: 2, maxRows: 5 }"
              placeholder="选填，学员查看试卷时可见"
              maxlength="500"
              show-count
            />
          </template>
          <div v-else-if="q.user_answer?.grade_remark" class="text-14">
            <span class="grade-muted">批改备注：</span>
            <div class="mt-4 whitespace-pre-wrap grade-answer-well">{{ q.user_answer.grade_remark }}</div>
          </div>
        </div>
      </div>
      </div>
    </div>
    <template #footer>
      <NSpace justify="space-between">
        <NButton
          v-if="detail?.user_paper?.user_id"
          type="info"
          secondary
          @click="openSendMessage"
        >
          发送站内信
        </NButton>
        <NSpace v-else />
        <NSpace>
          <NButton v-if="readOnly" @click="visible = false">关闭</NButton>
          <template v-else>
            <NButton @click="visible = false">取消</NButton>
            <NButton
              v-if="hasNext"
              :loading="submitting"
              @click="handleSubmit(true)"
            >
              提交并下一条
            </NButton>
            <NButton type="primary" :loading="submitting" @click="handleSubmit(false)">
              提交批改
            </NButton>
          </template>
        </NSpace>
      </NSpace>
    </template>
  </NModal>
  <NModal
    v-model:show="previewVisible"
    preset="card"
    class="file-preview-modal"
    :class="{ 'is-immersive': docFullscreen }"
    :style="previewModalStyle"
    :z-index="3000"
  >
    <template #header>
      <div class="file-preview-head">
        <div class="file-preview-head__copy">
          <strong>{{ previewName || '文件预览' }}</strong>
        </div>
        <div class="file-preview-head__actions">
          <button
            v-if="previewKind === 'doc' && previewUrl"
            type="button"
            class="file-preview-head__btn"
            :title="docFullscreen ? '退出全屏' : '全屏'"
            @click="toggleDocFullscreen"
          >
            <i :class="docFullscreen ? 'i-fe:minimize' : 'i-fe:maximize'" />
          </button>
          <a
            v-if="previewUrl"
            class="file-preview-head__btn"
            :href="previewUrl"
            target="_blank"
            rel="noopener" referrerpolicy="origin"
            title="新窗口打开"
          >
            <i class="i-fe:external-link" />
          </a>
          <a
            v-if="previewUrl"
            class="file-preview-head__btn"
            :href="previewUrl"
            target="_blank"
            rel="noopener" referrerpolicy="origin"
            download
            title="下载"
          >
            <i class="i-fe:download" />
          </a>
        </div>
      </div>
    </template>
    <div v-if="previewKind === 'image'" class="file-preview-image">
      <img :src="previewUrl" :alt="previewName" >
    </div>
    <div v-else-if="previewKind === 'video'" class="file-preview-media">
      <video :src="previewUrl" controls preload="metadata" playsinline />
    </div>
    <div v-else-if="previewKind === 'audio'" class="file-preview-media file-preview-media--audio">
      <audio :src="previewUrl" controls preload="metadata" />
    </div>
    <OfficeDocViewer
      v-else-if="previewUrl"
      class="file-preview-doc"
      :src="previewUrl"
      :filename="previewName"
    />
  </NModal>
  <SendMessageModal ref="$msgModal" />
</template>

<script setup>
import { NModal, NSpin, NTag, NInputNumber, NInput, NButton, NSpace, NImage } from 'naive-ui'
import api from '@/api'
import { getAssetUrl } from '@/utils'
import { useRouter } from 'vue-router'
import OfficeDocViewer from './OfficeDocViewer.vue'
import SendMessageModal from './SendMessageModal.vue'

defineOptions({ name: 'PaperGradeModal' })

const router = useRouter()

const emit = defineEmits(['success'])

const $msgModal = ref(null)
const visible = ref(false)
const loading = ref(false)
const submitting = ref(false)
const readOnly = ref(false)
const hasNext = ref(false)
const detail = ref(null)
const scoreInputs = ref({})
const remarkInputs = ref({})
const previewVisible = ref(false)
const previewUrl = ref('')
const previewName = ref('')
const previewKind = ref('doc')
const docFullscreen = ref(false)
const gradeBody = ref(null)

const previewModalStyle = computed(() => {
  if (docFullscreen.value) {
    return {
      width: '100vw',
      maxWidth: '100vw',
      height: '100vh',
      margin: 0,
      borderRadius: 0,
    }
  }
  return {
    width: '96vw',
    maxWidth: '1920px',
  }
})

function toggleDocFullscreen() {
  docFullscreen.value = !docFullscreen.value
}

function goUser() {
  const id = detail.value?.user_paper?.user_id || detail.value?.user_paper?.user?.id
  if (id)
    router.push({ name: 'UserDetail', params: { id } })
}

function goCourse() {
  const id = detail.value?.paper?.course_id || detail.value?.paper?.course?.id
  if (id)
    router.push(`/course/detail/${id}`)
}

function goChapter() {
  const courseId = detail.value?.paper?.course_id || detail.value?.paper?.course?.id
  const chapterId = detail.value?.paper?.chapter_id || detail.value?.paper?.chapter?.id
  if (courseId && chapterId)
    router.push(`/course/${courseId}/chapterDetail/${chapterId}`)
}

function open(userPaperId, viewOnly = false, options = {}) {
  readOnly.value = viewOnly
  hasNext.value = Boolean(options.hasNext) && !viewOnly
  visible.value = true
  loading.value = true
  detail.value = null
  scoreInputs.value = {}
  remarkInputs.value = {}
  api.paper.getUserPaperDetail(userPaperId).then((res) => {
    detail.value = res.data
    const inputs = {}
    const remarks = {}
    for (const q of res.data.questions || []) {
      if (q.user_answer?.need_people_check) {
        inputs[q.id] = q.user_answer.score ?? 0
        remarks[q.id] = q.user_answer.grade_remark || ''
      }
    }
    scoreInputs.value = inputs
    remarkInputs.value = remarks
  }).catch((error) => {
    window.$message?.error(error?.message || '加载答卷失败')
  }).finally(() => {
    loading.value = false
    nextTick(() => {
      gradeBody.value?.querySelector('.question-item--need')?.scrollIntoView({
        block: 'center',
        behavior: 'smooth',
      })
    })
  })
}

function handleClose() {
  detail.value = null
  scoreInputs.value = {}
  remarkInputs.value = {}
  readOnly.value = false
  hasNext.value = false
}

function openSendMessage() {
  const userPaper = detail.value?.user_paper
  if (!userPaper?.user_id) return
  const courseTitle = detail.value?.paper?.course?.title || '课程'
  const chapterTitle = detail.value?.paper?.chapter?.title
  const score = userPaper.score
  const target = chapterTitle ? `「${courseTitle}」课节「${chapterTitle}」` : `「${courseTitle}」`
  $msgModal.value?.open({
    userId: userPaper.user_id,
    userName: userPaper.user?.nick_name,
    type: 'paper',
    title: '作业批改通知',
    content: `同学，你提交的${target}作业已批改${score != null ? `，成绩 ${score} 分` : ''}。`,
    relatedType: 'user_paper',
    relatedId: userPaper.id,
  })
}

function isUserAnswer(q, optionIndex) {
  const content = q.user_answer?.answer_content
  if (!Array.isArray(content)) return false
  return content.includes(optionIndex) || content.includes(String(optionIndex))
}

function isCorrectAnswer(q, optionIndex) {
  const ans = q.answer
  if (!Array.isArray(ans)) return false
  return ans.includes(optionIndex) || ans.includes(String(optionIndex))
}

function formatUserOptionAnswer(q) {
  const content = q.user_answer?.answer_content
  if (!Array.isArray(content)) return '-'
  return content.map(i => String.fromCharCode(65 + Number(i))).join(', ')
}

function formatArrayAnswer(arr) {
  if (!Array.isArray(arr)) return arr || '-'
  return arr.join(', ')
}

/** 解析上传题答案中的原始路径 */
function normalizeUploadPaths(content) {
  if (content == null || content === '')
    return []

  if (typeof content === 'string') {
    const trimmed = content.trim()
    if (!trimmed)
      return []
    if (trimmed.startsWith('[') || trimmed.startsWith('{')) {
      try {
        return normalizeUploadPaths(JSON.parse(trimmed))
      }
      catch {
        return [trimmed]
      }
    }
    return [trimmed]
  }

  if (Array.isArray(content)) {
    return content.flatMap((item) => {
      if (item == null || item === '')
        return []
      if (typeof item === 'string')
        return normalizeUploadPaths(item)
      const path = item?.url || item?.path || item?.visit_url || ''
      return path ? [path] : []
    })
  }

  if (typeof content === 'object') {
    const path = content.url || content.path || content.visit_url || ''
    return path ? [path] : []
  }

  return [String(content)]
}

/** 获取上传题的用户答案 URL 列表（原始路径，展示时再拼接域名） */
function getUploadUrls(q) {
  return normalizeUploadPaths(q.user_answer?.answer_content).filter(Boolean)
}

/** 从 URL 提取文件名 */
function getFileName(url, index) {
  if (!url) return `文档${index + 1}`
  try {
    const path = url.split('?')[0]
    const name = path.split('/').pop()
    return decodeURIComponent(name || `文档${index + 1}`)
  }
  catch {
    return `文档${index + 1}`
  }
}

function detectPreviewKind(url) {
  const path = String(url || '').split('?')[0].split('#')[0]
  if (/\.(jpe?g|png|gif|webp|bmp|svg)$/i.test(path))
    return 'image'
  if (/\.(mp4|webm|mov|m4v|ogv)$/i.test(path))
    return 'video'
  if (/\.(mp3|wav|m4a|aac|flac|ogg)$/i.test(path))
    return 'audio'
  return 'doc'
}

function openFilePreview(url, index) {
  previewUrl.value = getAssetUrl(url)
  previewName.value = getFileName(url, index)
  previewKind.value = detectPreviewKind(url)
  docFullscreen.value = false
  previewVisible.value = true
}

function handleSubmit(goNext = false) {
  if (readOnly.value || submitting.value)
    return
  if (!detail.value?.user_paper?.id) return
  const questionScores = []
  for (const q of detail.value.questions || []) {
    if (q.user_answer?.need_people_check && q.id in scoreInputs.value) {
      questionScores.push({
        question_id: q.id,
        score: Number(scoreInputs.value[q.id]) || 0,
        grade_remark: String(remarkInputs.value[q.id] ?? '').trim(),
      })
    }
  }
  submitting.value = true
  api.paper.paperGrade({
    user_paper_id: detail.value.user_paper.id,
    question_scores: questionScores,
  }).then(() => {
    window.$message?.success('批改成功')
    const id = detail.value.user_paper.id
    visible.value = false
    emit('success', { id, next: goNext === true })
  }).catch(() => {
    // 拦截器已提示失败原因
  }).finally(() => {
    submitting.value = false
  })
}

defineExpose({ open })
</script>

<style scoped>
.grade-meta {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin: 0 0 12px;
  padding: 12px 14px;
  border-radius: 14px;
}

.grade-meta__main {
  min-width: 0;
  flex: 1;
}

.grade-meta .entity-card__badges {
  margin-bottom: 6px;
}

.grade-meta__ids {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.grade-meta .entity-card__metrics {
  flex-shrink: 0;
  margin-top: 0;
}

.grade-list {
  display: flex;
  flex-direction: column;
  gap: 1px;
  padding: 4px;
  border-radius: 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
}

.question-item {
  border: 0;
  border-radius: 12px;
  padding: 14px 12px;
  background: rgb(255 255 255 / 55%);
  box-shadow: none;
}

.question-item--need {
  background: rgb(201 151 42 / 8%);
  box-shadow: inset 3px 0 0 #c9972a;
}

.grade-muted {
  color: var(--text-secondary, #5b6b7c);
  letter-spacing: -0.01em;
}

.grade-muted-soft {
  color: var(--ctl-text-soft, #8b97a5);
  letter-spacing: -0.01em;
}

.grade-answer-well {
  padding: 12px 14px;
  border-radius: 12px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.05));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  letter-spacing: -0.01em;
  color: var(--text-primary, #1a2332);
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 40%);
}

.grade-file-row {
  border-radius: 12px;
  border: 1px solid var(--ctl-field-border, rgba(var(--primary-color), 0.1));
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 45%);
}

.grade-thumb-ph {
  background: var(--ctl-fill, rgba(var(--primary-color), 0.06));
  color: var(--ctl-text-soft, #8b97a5);
  letter-spacing: -0.01em;
}

.grade-score-row {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  margin-top: 10px;
  min-height: 32px;
}

.grade-score-row__label {
  flex-shrink: 0;
  font-size: 13px;
  font-weight: 550;
  letter-spacing: -0.01em;
  color: var(--text-secondary, #5b6b7c);
}

.grade-score-row__unit {
  flex-shrink: 0;
  font-size: 12.5px;
  letter-spacing: -0.01em;
  color: var(--ctl-text-soft, #8b97a5);
}

.grade-score-chip {
  appearance: none;
  height: 24px;
  padding: 0 8px;
  border: 0;
  border-radius: 999px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.08));
  color: rgb(var(--primary-color));
  font-size: 11px;
  font-weight: 560;
  cursor: pointer;
}

.grade-score-chip:hover {
  background: var(--ctl-field-hover, rgba(var(--primary-color), 0.14));
}

.grade-score-row strong {
  font-size: 15px;
  font-weight: 600;
  letter-spacing: -0.02em;
  color: var(--text-primary, #1a2332);
  min-width: 1.5em;
}

.grade-score-input {
  width: 132px !important;
  min-width: 132px !important;
  flex-shrink: 0;
}

.file-preview-head {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  min-width: 0;
  padding-right: 8px;
}

.file-preview-head__copy {
  min-width: 0;
  flex: 1;
}

.file-preview-head__copy strong {
  display: block;
  overflow: hidden;
  font-family: var(--font-display, inherit);
  font-size: 15px;
  font-weight: 650;
  letter-spacing: -0.02em;
  color: #1a2332;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.file-preview-head__actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.file-preview-head__btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 0;
  border-radius: 999px;
  color: #1a2332;
  background: rgba(var(--primary-color), 0.08);
  text-decoration: none;
  cursor: pointer;
  transition:
    transform 0.2s cubic-bezier(0.32, 0.72, 0, 1),
    background 0.2s ease;
}

.file-preview-head__btn:hover {
  background: rgba(var(--primary-color), 0.14);
}

.file-preview-head__btn:active {
  transform: scale(0.94);
}

.file-preview-image {
  display: flex;
  align-items: center;
  justify-content: center;
  max-height: 70vh;
  overflow: auto;
}

.file-preview-image img {
  max-width: 100%;
  max-height: 70vh;
  object-fit: contain;
}

.file-preview-media {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 240px;
  max-height: 70vh;
  background: #0f1418;
  border-radius: 14px;
  overflow: hidden;
}

.file-preview-media video {
  width: 100%;
  max-height: 70vh;
  background: #000;
}

.file-preview-media--audio {
  min-height: 0;
  padding: 28px 16px;
  background: #f4f7f7;
}

.file-preview-media audio {
  width: min(640px, 100%);
}

.file-preview-doc {
  width: 100%;
  height: 88vh;
  min-height: 480px;
  overflow: hidden;
  background: #f4f7f7;
  border-radius: 14px;
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.1));
}
</style>

<style>
.file-preview-modal.is-immersive.n-modal,
.n-modal.file-preview-modal.is-immersive {
  max-width: 100vw !important;
  width: 100vw !important;
  height: 100vh !important;
  margin: 0 !important;
  border-radius: 0 !important;
}

.file-preview-modal.is-immersive .n-card,
.n-modal.file-preview-modal.is-immersive .n-card {
  height: 100vh !important;
  max-height: 100vh !important;
  border-radius: 0 !important;
  display: flex !important;
  flex-direction: column !important;
}

.file-preview-modal.is-immersive .n-card-header {
  flex-shrink: 0;
}

.file-preview-modal.is-immersive .n-card__content {
  flex: 1 !important;
  min-height: 0 !important;
  display: flex !important;
  flex-direction: column !important;
  padding: 0 !important;
}

.file-preview-modal.is-immersive .file-preview-doc {
  flex: 1;
  height: auto !important;
  min-height: 0 !important;
  border-radius: 0;
  border: 0;
}

.n-modal-container:has(.file-preview-modal.is-immersive) .n-modal-body-wrapper,
.n-modal-container:has(.file-preview-modal.is-immersive) .n-modal-scroll-content {
  padding: 0 !important;
  width: 100vw !important;
  max-width: 100vw !important;
  height: 100vh !important;
}
</style>
