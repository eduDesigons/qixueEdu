<template>
  <AppPage show-footer class="dash-page">
    <div class="dash">
      <header class="page-mast dash-span-12">
        <h1 class="page-mast__title">工作台</h1>
        <p class="page-mast__sub">今日学情与待办</p>
      </header>
      <header class="hero dash-span-12">
        <LayoutAurora
          class="hero__aurora"
          preset="hero"
          :colors="verseAuroraColors"
          :params="verseAuroraParams"
          :opacity="1"
        />
        <blockquote v-if="verse" class="hero__verse">
          <div class="hero__quote">
            <span class="hero__mark hero__mark--l" aria-hidden="true">“</span>
            <p>{{ verse }}</p>
            <span class="hero__mark hero__mark--r" aria-hidden="true">”</span>
          </div>
          <span class="hero__rule" aria-hidden="true" />
        </blockquote>
        <div v-else class="hero__verse hero__verse--empty">
          <div class="hero__quote">
            <p>静候一句良言</p>
          </div>
          <span class="hero__rule" aria-hidden="true" />
        </div>
      </header>

      <NSpin :show="statsLoading" class="dash-span-12">
        <section v-if="visibleTodos.length" class="panel panel--todo">
          <div class="panel__head">
            <h3>待办事项</h3>
            <em v-if="todoTotal">含 {{ todoTotal }} 项待处理</em>
          </div>
          <div class="todo">
            <button
              v-if="canAccess('/member/courseApply')"
              type="button"
              class="todo__item todo__item--coral"
              :class="{ 'is-on': stats.apply_pending_count > 0 }"
              @click="goTo('/member/courseApply?tab=course')"
            >
              <span class="todo__ico"><i class="i-fe:inbox" /></span>
              <div class="todo__meta">
                <strong><NNumberAnimation :from="0" :to="stats.apply_pending_count" /></strong>
                <b>入课待审</b>
                <small>学员入课申请</small>
              </div>
            </button>
            <button
              v-if="canAccess('/member/courseApply')"
              type="button"
              class="todo__item todo__item--gold"
              :class="{ 'is-on': stats.path_apply_pending_count > 0 }"
              @click="goTo('/member/courseApply?tab=path')"
            >
              <span class="todo__ico"><i class="i-fe:layers" /></span>
              <div class="todo__meta">
                <strong><NNumberAnimation :from="0" :to="stats.path_apply_pending_count" /></strong>
                <b>入系待审</b>
                <small>系列加入申请</small>
              </div>
            </button>
            <button
              v-if="canAccess('/paper/assessment')"
              type="button"
              class="todo__item todo__item--paper"
              :class="{ 'is-on': stats.paper_pending_count > 0 }"
              @click="goTo('/paper/assessment')"
            >
              <span class="todo__ico"><i class="i-fe:edit-3" /></span>
              <div class="todo__meta">
                <strong><NNumberAnimation :from="0" :to="stats.paper_pending_count" /></strong>
                <b>作业待批改</b>
                <small>需人工批改</small>
              </div>
            </button>
            <button
              v-if="canAccess('/member/feedback')"
              type="button"
              class="todo__item todo__item--feedback"
              :class="{ 'is-on': stats.feedback_pending_count > 0 }"
              @click="goTo('/member/feedback?status=0')"
            >
              <span class="todo__ico"><i class="i-fe:message-square" /></span>
              <div class="todo__meta">
                <strong><NNumberAnimation :from="0" :to="stats.feedback_pending_count" /></strong>
                <b>意见建议</b>
                <small>待回复</small>
              </div>
            </button>
            <button
              v-if="canAccessComment"
              type="button"
              class="todo__item todo__item--teal"
              :class="{ 'is-on': (stats.course_comment_count + stats.chapter_comment_count) > 0 }"
              @click="goComment"
            >
              <span class="todo__ico"><i class="i-fe:message-circle" /></span>
              <div class="todo__meta">
                <strong>
                  <NNumberAnimation :from="0" :to="stats.course_comment_count" />
                  <span class="todo__split">/</span>
                  <NNumberAnimation :from="0" :to="stats.chapter_comment_count" />
                </strong>
                <b>今日评论</b>
                <small>课程 / 课节</small>
              </div>
            </button>
          </div>
        </section>

        <div class="sec">今日学情</div>
        <section class="bento">
          <component
            :is="canAccess('/course/list') ? 'button' : 'div'"
            type="button"
            class="bento__feat"
            @click="canAccess('/course/list') && goTo('/course/list')"
          >
            <span>完成课节</span>
            <div class="bento__feat-num">
              <strong><NNumberAnimation :from="0" :to="stats.today_chapter_end_count" /></strong>
              <em>课节</em>
            </div>
            <svg v-if="sparkChapter.d" class="bento__spark" viewBox="0 0 280 72" preserveAspectRatio="none">
              <defs>
                <linearGradient id="dashSpark" x1="0" y1="0" x2="1" y2="0">
                  <stop offset="0%" stop-color="rgba(var(--primary-color), 0.72)" />
                  <stop offset="100%" stop-color="rgb(var(--primary-color))" />
                </linearGradient>
                <linearGradient id="dashSparkFill" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stop-color="rgb(var(--primary-color))" stop-opacity="0.28" />
                  <stop offset="100%" stop-color="rgb(var(--primary-color))" stop-opacity="0" />
                </linearGradient>
              </defs>
              <path :d="sparkChapter.area" fill="url(#dashSparkFill)" />
              <path :d="sparkChapter.d" fill="none" stroke="url(#dashSpark)" stroke-width="2.4" stroke-linecap="round" />
            </svg>
          </component>
          <div class="bento__mid">
            <component
              :is="canAccess('/course/list') ? 'button' : 'div'"
              type="button"
              class="bento__tile bento__tile--gold bento__tile--spark"
              @click="canAccess('/course/list') && goTo('/course/list')"
            >
              <span>完成课程</span>
              <strong><NNumberAnimation :from="0" :to="stats.today_course_end_count" /></strong>
              <em>课程</em>
              <svg v-if="sparkCourse.d" class="bento__spark" viewBox="0 0 280 72" preserveAspectRatio="none">
                <defs>
                  <linearGradient id="dashSparkCourse" x1="0" y1="0" x2="1" y2="0">
                    <stop offset="0%" stop-color="rgb(201 151 42 / 0.72)" />
                    <stop offset="100%" stop-color="#c9972a" />
                  </linearGradient>
                  <linearGradient id="dashSparkCourseFill" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stop-color="#c9972a" stop-opacity="0.28" />
                    <stop offset="100%" stop-color="#c9972a" stop-opacity="0" />
                  </linearGradient>
                </defs>
                <path :d="sparkCourse.area" fill="url(#dashSparkCourseFill)" />
                <path :d="sparkCourse.d" fill="none" stroke="url(#dashSparkCourse)" stroke-width="2.4" stroke-linecap="round" />
              </svg>
            </component>
            <component
              :is="canAccess('/paper/exam') ? 'button' : 'div'"
              type="button"
              class="bento__tile bento__tile--slate bento__tile--spark"
              @click="canAccess('/paper/exam') && goTo('/paper/exam?pass=1')"
            >
              <span>测试通过</span>
              <strong><NNumberAnimation :from="0" :to="stats.today_paper_pass_count" /></strong>
              <em>{{ stats.today_paper_pass_rate }}% · {{ stats.today_paper_count }} 次提交</em>
              <svg v-if="sparkPaper.d" class="bento__spark" viewBox="0 0 280 72" preserveAspectRatio="none">
                <defs>
                  <linearGradient id="dashSparkPaper" x1="0" y1="0" x2="1" y2="0">
                    <stop offset="0%" stop-color="rgb(93 109 130 / 0.72)" />
                    <stop offset="100%" stop-color="#5d6d82" />
                  </linearGradient>
                  <linearGradient id="dashSparkPaperFill" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stop-color="#5d6d82" stop-opacity="0.28" />
                    <stop offset="100%" stop-color="#5d6d82" stop-opacity="0" />
                  </linearGradient>
                </defs>
                <path :d="sparkPaper.area" fill="url(#dashSparkPaperFill)" />
                <path :d="sparkPaper.d" fill="none" stroke="url(#dashSparkPaper)" stroke-width="2.4" stroke-linecap="round" />
              </svg>
            </component>
          </div>
          <div class="bento__side">
            <component
              :is="canJumpUser ? 'button' : 'div'"
              type="button"
              class="bento__tile bento__tile--gold"
              @click="canJumpUser && goTo('/member/users')"
            >
              <span>学分</span>
              <strong><NNumberAnimation :from="0" :to="stats.today_credit_sum" /></strong>
              <em>今日发放</em>
            </component>
            <component
              :is="canJumpUser ? 'button' : 'div'"
              type="button"
              class="bento__tile bento__tile--teal"
              @click="canJumpUser && goTo('/member/users')"
            >
              <span>今日发放积分</span>
              <strong><NNumberAnimation :from="0" :to="stats.today_coins_sum" /></strong>
              <em>积分</em>
            </component>
          </div>
        </section>
      </NSpin>

      <div class="dash__row dash-span-12">
        <section class="panel panel--chart">
          <div class="panel__head">
            <h3>学习趋势</h3>
            <p>最近 14 天</p>
          </div>
          <VChart
            v-if="lineOption"
            class="chart"
            :option="lineOption"
            autoresize
            :style="{ width: '100%', height: '280px' }"
          />
          <div v-else-if="!chartLoading" class="chart chart--empty">
            暂无趋势数据
          </div>
          <NSpin v-else :show="chartLoading" class="chart chart--spin" />
        </section>
        <section class="panel panel--credit">
          <div class="panel__head">
            <h3>发放积分</h3>
            <p>最近 14 天</p>
            <button
              v-if="canJumpUser"
              type="button"
              class="panel__more"
              @click="goTo('/member/users')"
            >
              查看全部
            </button>
            <span v-else>Top 8</span>
          </div>
          <NSpin :show="userRankLoading">
            <div v-if="userPodium.length" class="credit">
              <div class="podium">
                <div
                  v-for="item in userPodium"
                  :key="item.user_id"
                  class="podium__col"
                  :class="[`is-${item.rank}`, { 'is-link': canJumpUser }]"
                  @click="openUserDetail(item.user_id)"
                >
                  <div class="podium__avatar">
                    <NAvatar round :size="item.rank === 1 ? 54 : 40" :src="$assetUrl(item.avatar)" />
                  </div>
                  <small>{{ item.nick_name || '未设置昵称' }}</small>
                  <b>{{ item.rank }}</b>
                </div>
              </div>
              <div class="credit__list">
                <div
                  v-for="item in userRankTop"
                  :key="item.user_id"
                  class="credit__row"
                  :class="{ 'is-link': canJumpUser }"
                  @click="openUserDetail(item.user_id)"
                >
                  <em :class="rankNumClass(item.rank)">{{ item.rank }}</em>
                  <span>{{ item.nick_name || '未设置昵称' }}</span>
                  <span class="bar"><span :style="{ width: `${userRankPercent(item)}%` }" /></span>
                  <strong>{{ formatCount(item.credit) }} 分</strong>
                </div>
              </div>
            </div>
            <div v-else-if="!userRankLoading" class="empty">暂无数据</div>
          </NSpin>
        </section>
      </div>

      <div class="dash__row dash__row--2 dash-span-12">
        <section class="panel panel--enroll">
          <div class="panel__head">
            <h3>课程学习中</h3>
            <button
              v-if="canJumpCourse"
              type="button"
              class="panel__more"
              @click="goTo('/course/list')"
            >
              查看全部
            </button>
            <span v-else>Top 10</span>
          </div>
          <NSpin :show="courseRankLoading">
            <div v-if="courseRankShort.length" class="rank">
              <div
                v-for="item in courseRankShort"
                :key="item.course_id"
                class="rank__row"
                :class="{ 'is-link': canJumpCourse }"
                @click="openCourseDetail(item.course_id)"
              >
                <em :class="rankNumClass(item.rank)">{{ padRank(item.rank) }}</em>
                <div>
                  <b :title="item.title">{{ item.title }}</b>
                  <span class="bar"><span :style="{ width: `${courseRankPercent(item)}%` }" /></span>
                </div>
                <strong>{{ formatCount(item.enroll_count) }}<small>人</small></strong>
              </div>
            </div>
            <div v-else-if="!courseRankLoading" class="empty">暂无数据</div>
          </NSpin>
        </section>
        <section class="panel panel--chapter">
          <div class="panel__head">
            <h3>近 30 天课节</h3>
            <button
              v-if="canJumpChapter"
              type="button"
              class="panel__more"
              @click="goTo('/course/list')"
            >
              查看全部
            </button>
            <span v-else>Top 8</span>
          </div>
          <NSpin :show="chapterRankLoading">
            <div v-if="chapterRankShort.length" class="rank">
              <div
                v-for="item in chapterRankShort"
                :key="`${item.course_id}-${item.chapter_id}`"
                class="rank__row"
                :class="{ 'is-link': canJumpChapter }"
                @click="openChapterDetail(item.course_id, item.chapter_id)"
              >
                <em :class="rankNumClass(item.rank)">{{ padRank(item.rank) }}</em>
                <div>
                  <b :title="item.chapter_title">{{ item.chapter_title }}</b>
                  <small>{{ item.course_title }}</small>
                  <span class="bar"><span :style="{ width: `${chapterRankPercent(item)}%` }" /></span>
                </div>
                <strong>{{ formatCount(item.complete_count) }}<small>人</small></strong>
              </div>
            </div>
            <div v-else-if="!chapterRankLoading" class="empty">暂无数据</div>
          </NSpin>
        </section>
      </div>

      <div class="dash-span-12 cdn-block">
        <section class="cdn-shell">
          <header class="cdn-shell__head">
            <div class="cdn-shell__titles">
              <div class="sec">CDN 流量</div>
            </div>
          </header>

          <NSpin :show="cdnLoading">
            <template v-if="cdnData?.enabled">
              <div class="cdn-section">
                <div class="cdn-kpi">
                  <article class="cdn-card cdn-card--bw-today">
                    <div class="cdn-card__ico"><i class="i-fe:trending-up" /></div>
                    <span>今日流量</span>
                    <strong>{{ formatBytes(cdnData.today?.bandwidth) }}</strong>
                    <em>{{ cdnData.today?.date || '今天' }}</em>
                  </article>
                  <article class="cdn-card cdn-card--req-today">
                    <div class="cdn-card__ico"><i class="i-fe:zap" /></div>
                    <span>今日请求</span>
                    <strong>{{ formatCount(cdnData.today?.requests) }}</strong>
                    <em>次</em>
                  </article>
                  <article class="cdn-card cdn-card--bw-month">
                    <div class="cdn-card__ico"><i class="i-fe:bar-chart-2" /></div>
                    <span>30 天流量</span>
                    <strong>{{ formatBytes(cdnData.month?.total_bandwidth) }}</strong>
                    <em>{{ cdnData.month?.from }} ~ {{ cdnData.month?.to }}</em>
                  </article>
                  <article class="cdn-card cdn-card--req-month">
                    <div class="cdn-card__ico"><i class="i-fe:layers" /></div>
                    <span>30 天请求</span>
                    <strong>{{ formatCount(cdnData.month?.total_requests) }}</strong>
                    <em>次</em>
                  </article>
                </div>

                <div class="cdn-grid">
                  <section class="panel panel--cdn panel--cdn-bw">
                    <div class="panel__head">
                      <h3>今日流量</h3>
                      <p>按小时 · 对比昨/前日</p>
                    </div>
                    <VChart
                      v-if="cdnTodayBandwidthOption"
                      class="chart"
                      :option="cdnTodayBandwidthOption"
                      autoresize
                      :style="{ width: '100%', height: '220px' }"
                    />
                    <div v-else class="chart chart--empty">暂无数据</div>
                  </section>

                  <section class="panel panel--cdn panel--cdn-req">
                    <div class="panel__head">
                      <h3>今日请求</h3>
                      <p>按小时 · 对比昨/前日</p>
                    </div>
                    <VChart
                      v-if="cdnTodayRequestOption"
                      class="chart"
                      :option="cdnTodayRequestOption"
                      autoresize
                      :style="{ width: '100%', height: '220px' }"
                    />
                    <div v-else class="chart chart--empty">暂无数据</div>
                  </section>
                </div>
              </div>

              <div class="cdn-section">
                <div class="cdn-section__label">
                  <i class="i-fe:calendar" />
                  <span>近 30 天</span>
                  <em>按天</em>
                </div>
                <div class="cdn-grid">
                  <section class="panel panel--cdn panel--cdn-bw-month">
                    <div class="panel__head">
                      <h3>流量趋势</h3>
                      <p>近 30 天 · 按天</p>
                    </div>
                    <VChart
                      v-if="cdnBandwidthOption"
                      class="chart"
                      :option="cdnBandwidthOption"
                      autoresize
                      :style="{ width: '100%', height: '250px' }"
                    />
                    <div v-else class="chart chart--empty">暂无数据</div>
                  </section>

                  <section class="panel panel--cdn panel--cdn-req-month">
                    <div class="panel__head">
                      <h3>请求趋势</h3>
                      <p>近 30 天 · 按天</p>
                    </div>
                    <VChart
                      v-if="cdnRequestOption"
                      class="chart"
                      :option="cdnRequestOption"
                      autoresize
                      :style="{ width: '100%', height: '250px' }"
                    />
                    <div v-else class="chart chart--empty">暂无数据</div>
                  </section>
                </div>
              </div>
            </template>
            <div v-else-if="!cdnLoading" class="cdn-empty">
              <i class="i-fe:hard-drive cdn-empty__ico" />
              <p>{{ cdnData?.message || 'CDN 统计暂不可用' }}</p>
            </div>
          </NSpin>
        </section>
      </div>
    </div>
  </AppPage>
</template>

<script setup>
import { BarChart, LineChart } from 'echarts/charts'
import { GridComponent, LegendComponent, TooltipComponent } from 'echarts/components'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { NAvatar, NNumberAnimation, NSpin } from 'naive-ui'
import { computed, onActivated, onMounted, reactive, ref } from 'vue'
import VChart from 'vue-echarts'
import { useRouter } from 'vue-router'
import api from '@/api'
import { LayoutAurora } from '@/layouts/components'
import { usePermissionStore } from '@/store'

defineOptions({ name: 'Home' })

use([CanvasRenderer, LineChart, BarChart, GridComponent, TooltipComponent, LegendComponent])

const router = useRouter()
const permissionStore = usePermissionStore()
const verse = ref('')
const verseAuroraColors = ['#7eb8bd', '#FFFFFF', '#FFFFFF']
const verseAuroraParams = {
  speed: 14,
  shape: 'checks',
  shapeScale: 11,
  rotation: -5,
  proportion: 48,
  softness: 100,
  distortion: 22,
  swirl: 13,
  swirlIterations: 9,
  scale: 0.55,
  offsetX: 0,
  offsetY: 55,
}
const chartLoading = ref(true)
const statsLoading = ref(true)
const cdnLoading = ref(true)
const courseRankLoading = ref(true)
const chapterRankLoading = ref(true)
const userRankLoading = ref(true)
const chartData = ref(null)
const cdnData = ref(null)
const courseRankList = ref([])
const chapterRankList = ref([])
const userRankList = ref([])

const stats = reactive({
  apply_pending_count: 0,
  path_apply_pending_count: 0,
  paper_pending_count: 0,
  today_course_end_count: 0,
  today_chapter_end_count: 0,
  today_paper_count: 0,
  today_paper_pass_count: 0,
  today_paper_pass_rate: 0,
  today_credit_sum: 0,
  today_coins_sum: 0,
  course_comment_count: 0,
  chapter_comment_count: 0,
  feedback_pending_count: 0,
  paper_pass_rate_all: 0,
})

function normalizePath(path = '') {
  if (!path)
    return ''
  const pure = String(path).split('?')[0].split('#')[0]
  if (pure.length > 1 && pure.endsWith('/'))
    return pure.slice(0, -1)
  return pure
}

function matchAccessPath(routePath, target) {
  if (!routePath || !target)
    return false
  if (routePath === target)
    return true
  if (!routePath.includes(':'))
    return false
  const pattern = routePath
    .split('/')
    .map((seg) => {
      if (seg.startsWith(':'))
        return '[^/]+'
      return seg.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
    })
    .join('/')
  return new RegExp(`^${pattern}$`).test(target)
}

function canAccess(path, fallbacks = []) {
  const target = normalizePath(path)
  if (!target)
    return false
  const routes = permissionStore.accessRoutes || []
  const hit = routes.some(item => matchAccessPath(normalizePath(item.path), target))
  if (hit)
    return true
  return fallbacks.some((fb) => {
    const fallback = normalizePath(fb)
    return fallback && routes.some(item => matchAccessPath(normalizePath(item.path), fallback))
  })
}

const canAccessComment = computed(() =>
  canAccess('/comment/courseComment') || canAccess('/comment/chapterComment'),
)
const visibleTodos = computed(() => ([
  canAccess('/member/courseApply'),
  canAccess('/paper/assessment'),
  canAccess('/member/feedback'),
  canAccessComment.value,
].filter(Boolean)))
const canJumpUser = computed(() => canAccess('/member/user/0', ['/member/users']))
const canJumpCourse = computed(() => canAccess('/course/detail/0', ['/course/list']))
const canJumpChapter = computed(() => canAccess('/course/0/chapterDetail/0', ['/course/list']))

const todoTotal = computed(() => {
  let total = 0
  if (canAccess('/member/courseApply')) {
    total += stats.apply_pending_count || 0
    total += stats.path_apply_pending_count || 0
  }
  if (canAccess('/paper/assessment'))
    total += stats.paper_pending_count || 0
  if (canAccess('/member/feedback'))
    total += stats.feedback_pending_count || 0
  return total
})
const courseRankShort = computed(() => courseRankList.value.slice(0, 10))
const chapterRankShort = computed(() => chapterRankList.value.slice(0, 8))
const userPodium = computed(() => userRankList.value.slice(0, 3))
const userRankTop = computed(() => userRankList.value.slice(0, 8))

function buildSpark(key) {
  const arr = (chartData.value?.line?.[key] || []).map(n => Number(n) || 0)
  if (arr.length < 2)
    return { d: '', area: '' }
  const max = Math.max(...arr, 1)
  const w = 280
  const h = 72
  const pts = arr.map((v, i) => {
    const x = (i / (arr.length - 1)) * w
    const y = h - 8 - (v / max) * (h - 16)
    return [x.toFixed(1), y.toFixed(1)]
  })
  const d = pts.map((p, i) => `${i ? 'L' : 'M'}${p[0]} ${p[1]}`).join(' ')
  return { d, area: `${d} L${w} ${h} L0 ${h} Z` }
}

const sparkChapter = computed(() => buildSpark('chapter_end'))
const sparkCourse = computed(() => buildSpark('course_end'))
const sparkPaper = computed(() => buildSpark('paper_pass'))

function goTo(path) {
  const pathname = String(path || '').split('?')[0]
  if (path && canAccess(pathname))
    router.push(path)
}

function goComment() {
  const courseOk = canAccess('/comment/courseComment')
  const chapterOk = canAccess('/comment/chapterComment')
  const suffix = '?today=1'
  if (courseOk && chapterOk) {
    goTo((stats.chapter_comment_count > stats.course_comment_count
      ? '/comment/chapterComment'
      : '/comment/courseComment') + suffix)
    return
  }
  if (courseOk)
    goTo(`/comment/courseComment${suffix}`)
  else if (chapterOk)
    goTo(`/comment/chapterComment${suffix}`)
}

function openUserDetail(userId) {
  if (!userId || !canJumpUser.value)
    return
  router.push(`/member/user/${userId}`)
}

function openCourseDetail(courseId) {
  if (!courseId || !canJumpCourse.value)
    return
  router.push(`/course/detail/${courseId}`)
}

function openChapterDetail(courseId, chapterId) {
  if (!courseId || !chapterId || !canJumpChapter.value)
    return
  router.push(`/course/${courseId}/chapterDetail/${chapterId}`)
}

const courseRankMax = computed(() => Math.max(...courseRankShort.value.map(i => i.enroll_count), 1))
const chapterRankMax = computed(() => Math.max(...chapterRankShort.value.map(i => i.complete_count), 1))
const userRankMax = computed(() => Math.max(...userRankTop.value.map(i => i.credit), 1))
function courseRankPercent(item) {
  return Math.max(8, Math.round((item.enroll_count / courseRankMax.value) * 100))
}
function chapterRankPercent(item) {
  return Math.max(8, Math.round((item.complete_count / chapterRankMax.value) * 100))
}
function userRankPercent(item) {
  return Math.max(8, Math.round((item.credit / userRankMax.value) * 100))
}
function asRankList(data) {
  if (Array.isArray(data))
    return data
  if (data && typeof data === 'object')
    return Object.values(data)
  return []
}
function padRank(rank) {
  return String(rank).padStart(2, '0')
}
function formatCount(n) {
  return Number(n || 0).toLocaleString('zh-CN')
}
function formatBytes(bytes) {
  const n = Number(bytes) || 0
  if (n < 1024)
    return `${n} B`
  const units = ['KB', 'MB', 'GB', 'TB', 'PB']
  let v = n
  let i = -1
  do {
    v /= 1024
    i++
  } while (v >= 1024 && i < units.length - 1)
  return `${v >= 100 ? v.toFixed(0) : v >= 10 ? v.toFixed(1) : v.toFixed(2)} ${units[i]}`
}
function formatBandwidthAxis(bytes) {
  const n = Number(bytes) || 0
  if (n <= 0)
    return '0'
  const gb = n / (1024 ** 3)
  if (gb >= 1024)
    return `${(gb / 1024).toFixed(1)}T`
  if (gb >= 10)
    return `${gb.toFixed(0)}G`
  if (gb >= 1)
    return `${gb.toFixed(1)}G`
  return `${(n / (1024 ** 2)).toFixed(0)}M`
}

function cdnDates() {
  return (cdnData.value?.month?.dates || []).map(d => String(d).slice(5))
}

function cdnHourAxis() {
  return {
    type: 'category',
    boundaryGap: false,
    axisLine: { show: false },
    axisTick: { show: false },
    axisLabel: { color: '#8a9697', fontSize: 11 },
    data: cdnData.value?.today?.hours || [],
  }
}

function cdnBaseTooltip(kind = 'count') {
  return {
    trigger: 'axis',
    backgroundColor: 'rgba(31, 42, 43, 0.92)',
    borderWidth: 0,
    extraCssText: 'border-radius: 12px; padding: 8px 12px; box-shadow: 0 10px 24px rgba(0,0,0,.18);',
    textStyle: { color: '#fff', fontSize: 12 },
    axisPointer: { lineStyle: { color: themeColor(0.35) } },
    formatter(params) {
      if (!Array.isArray(params) || !params.length)
        return ''
      const head = params[0].axisValueLabel || params[0].axisValue
      const lines = params.map((p) => {
        const val = kind === 'bytes' ? formatBytes(p.value) : formatCount(p.value)
        return `${p.marker}${p.seriesName}  ${val}`
      })
      return [head, ...lines].join('<br/>')
    },
  }
}

function cdnAxisCategory() {
  return {
    type: 'category',
    boundaryGap: false,
    data: cdnDates(),
    axisLine: { show: false },
    axisTick: { show: false },
    axisLabel: { color: '#8a9697', fontSize: 11 },
  }
}
function rankNumClass(rank) {
  if (rank === 1)
    return 'is-gold'
  if (rank === 2)
    return 'is-silver'
  if (rank === 3)
    return 'is-bronze'
  return ''
}

function themeRgb() {
  if (typeof document === 'undefined')
    return '49, 108, 114'
  return getComputedStyle(document.documentElement).getPropertyValue('--primary-color').trim() || '49, 108, 114'
}

function themeColor(alpha) {
  const rgb = themeRgb()
  return alpha == null ? `rgb(${rgb})` : `rgba(${rgb}, ${alpha})`
}

const lineOption = computed(() => {
  const line = chartData.value?.line
  if (!line?.dates?.length)
    return null
  const primary = themeColor()
  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(31, 42, 43, 0.92)',
      borderWidth: 0,
      extraCssText: 'border-radius: 12px; padding: 8px 12px; box-shadow: 0 10px 24px rgba(0,0,0,.18);',
      textStyle: { color: '#fff', fontSize: 12 },
      axisPointer: { lineStyle: { color: themeColor(0.35) } },
    },
    legend: {
      data: ['完成课节', '完成课程', '测试通过', '发放学分'],
      top: 0,
      right: 0,
      itemWidth: 10,
      itemHeight: 6,
      textStyle: { color: '#6b7879', fontSize: 11 },
    },
    color: [primary, '#c9972a', '#5d6d82', '#8b9bb4'],
    grid: { left: 4, right: 8, bottom: 8, top: 28, containLabel: true },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: line.dates.map(d => d.slice(5)),
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { color: '#8a9697', fontSize: 11 },
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      splitLine: { lineStyle: { color: themeColor(0.08), type: 'dashed' } },
      axisLabel: { color: '#8a9697', fontSize: 11 },
    },
    series: [
      {
        name: '完成课节',
        type: 'line',
        smooth: true,
        showSymbol: false,
        lineStyle: { width: 2.4 },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: themeColor(0.24) },
              { offset: 1, color: themeColor(0.02) },
            ],
          },
        },
        emphasis: { focus: 'series' },
        data: line.chapter_end,
      },
      { name: '完成课程', type: 'line', smooth: true, showSymbol: false, lineStyle: { width: 2.2 }, emphasis: { focus: 'series' }, data: line.course_end },
      { name: '测试通过', type: 'line', smooth: true, showSymbol: false, lineStyle: { width: 2.2 }, emphasis: { focus: 'series' }, data: line.paper_pass },
      { name: '发放学分', type: 'line', smooth: true, showSymbol: false, lineStyle: { width: 2 }, emphasis: { focus: 'series' }, data: line.credit_sum },
    ],
  }
})

const cdnTodayBandwidthOption = computed(() => {
  const today = cdnData.value?.today
  if (!today?.hours?.length)
    return null
  const yLabel = today.yesterday_date ? `昨日 ${String(today.yesterday_date).slice(5)}` : '昨日'
  const bLabel = today.day_before_date ? `前日 ${String(today.day_before_date).slice(5)}` : '前日'
  return {
    legend: {
      top: 0,
      right: 8,
      itemWidth: 10,
      itemHeight: 6,
      textStyle: { color: '#6b7879', fontSize: 11 },
      data: ['今日', yLabel, bLabel],
    },
    tooltip: cdnBaseTooltip('bytes'),
    grid: { left: 8, right: 8, bottom: 8, top: 28, containLabel: true },
    color: ['#2f7f86', '#8aabb0', '#c5d4d6'],
    xAxis: cdnHourAxis(),
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: 'rgba(47, 127, 134, 0.1)', type: 'dashed' } },
      axisLabel: { color: '#8a9697', fontSize: 11, formatter: formatBandwidthAxis },
    },
    series: [
      {
        name: '今日',
        type: 'line',
        smooth: true,
        showSymbol: false,
        connectNulls: false,
        lineStyle: { width: 2.8, color: '#2f7f86' },
        itemStyle: { color: '#2f7f86' },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(47, 127, 134, 0.28)' },
              { offset: 1, color: 'rgba(47, 127, 134, 0.02)' },
            ],
          },
        },
        z: 3,
        data: today.bandwidth_hourly || [],
      },
      {
        name: yLabel,
        type: 'line',
        smooth: true,
        showSymbol: false,
        lineStyle: { width: 2, type: 'dashed', color: '#7a9ea3' },
        itemStyle: { color: '#7a9ea3' },
        z: 2,
        data: today.bandwidth_yesterday_hourly || [],
      },
      {
        name: bLabel,
        type: 'line',
        smooth: true,
        showSymbol: false,
        lineStyle: { width: 1.8, type: 'dotted', color: '#b0c4c7' },
        itemStyle: { color: '#b0c4c7' },
        z: 1,
        data: today.bandwidth_day_before_hourly || [],
      },
    ],
  }
})

const cdnTodayRequestOption = computed(() => {
  const today = cdnData.value?.today
  if (!today?.hours?.length)
    return null
  const yLabel = today.yesterday_date ? `昨日 ${String(today.yesterday_date).slice(5)}` : '昨日'
  const bLabel = today.day_before_date ? `前日 ${String(today.day_before_date).slice(5)}` : '前日'
  return {
    legend: {
      top: 0,
      right: 8,
      itemWidth: 10,
      itemHeight: 6,
      textStyle: { color: '#6b7879', fontSize: 11 },
      data: ['今日', yLabel, bLabel],
    },
    tooltip: cdnBaseTooltip('count'),
    grid: { left: 8, right: 8, bottom: 8, top: 28, containLabel: true },
    color: ['#c47a3a', '#d4a574', '#e6d0b5'],
    xAxis: cdnHourAxis(),
    yAxis: {
      type: 'value',
      minInterval: 1,
      splitLine: { lineStyle: { color: 'rgba(196, 122, 58, 0.1)', type: 'dashed' } },
      axisLabel: { color: '#8a9697', fontSize: 11 },
    },
    series: [
      {
        name: '今日',
        type: 'line',
        smooth: true,
        showSymbol: false,
        connectNulls: false,
        lineStyle: { width: 2.8, color: '#c47a3a' },
        itemStyle: { color: '#c47a3a' },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(196, 122, 58, 0.26)' },
              { offset: 1, color: 'rgba(196, 122, 58, 0.02)' },
            ],
          },
        },
        z: 3,
        data: today.requests_hourly || [],
      },
      {
        name: yLabel,
        type: 'line',
        smooth: true,
        showSymbol: false,
        lineStyle: { width: 2, type: 'dashed', color: '#d4a574' },
        itemStyle: { color: '#d4a574' },
        z: 2,
        data: today.requests_yesterday_hourly || [],
      },
      {
        name: bLabel,
        type: 'line',
        smooth: true,
        showSymbol: false,
        lineStyle: { width: 1.8, type: 'dotted', color: '#e6d0b5' },
        itemStyle: { color: '#e6d0b5' },
        z: 1,
        data: today.requests_day_before_hourly || [],
      },
    ],
  }
})

const cdnBandwidthOption = computed(() => {
  const month = cdnData.value?.month
  if (!month?.dates?.length)
    return null
  return {
    tooltip: cdnBaseTooltip('bytes'),
    grid: { left: 8, right: 8, bottom: 8, top: 20, containLabel: true },
    color: ['#498276'],
    xAxis: {
      ...cdnAxisCategory(),
      boundaryGap: true,
      axisLabel: { color: '#8a9697', fontSize: 10, interval: 2 },
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: 'rgba(73, 130, 118, 0.1)', type: 'dashed' } },
      axisLabel: { color: '#8a9697', fontSize: 11, formatter: formatBandwidthAxis },
    },
    series: [
      {
        name: '流量',
        type: 'bar',
        barMaxWidth: 16,
        itemStyle: {
          borderRadius: [7, 7, 2, 2],
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(73, 130, 118, 0.92)' },
              { offset: 1, color: 'rgba(73, 130, 118, 0.35)' },
            ],
          },
        },
        emphasis: {
          itemStyle: { color: '#498276' },
        },
        data: month.bandwidth || [],
      },
    ],
  }
})

const cdnRequestOption = computed(() => {
  const month = cdnData.value?.month
  if (!month?.dates?.length)
    return null
  return {
    tooltip: cdnBaseTooltip('count'),
    grid: { left: 8, right: 8, bottom: 8, top: 20, containLabel: true },
    color: ['#5d6d82'],
    xAxis: {
      ...cdnAxisCategory(),
      boundaryGap: true,
      axisLabel: { color: '#8a9697', fontSize: 10, interval: 2 },
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      splitLine: { lineStyle: { color: 'rgba(93, 109, 130, 0.1)', type: 'dashed' } },
      axisLabel: { color: '#8a9697', fontSize: 11 },
    },
    series: [
      {
        name: '请求',
        type: 'bar',
        barMaxWidth: 16,
        itemStyle: {
          borderRadius: [7, 7, 2, 2],
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(93, 109, 130, 0.92)' },
              { offset: 1, color: 'rgba(93, 109, 130, 0.32)' },
            ],
          },
        },
        emphasis: {
          itemStyle: { color: '#5d6d82' },
        },
        data: month.requests || [],
      },
    ],
  }
})

function applyStatsPatch(data) {
  Object.keys(stats).forEach((key) => {
    if (data[key] !== undefined)
      stats[key] = data[key]
  })
}

function loadDashboard() {
  cdnLoading.value = true
  statsLoading.value = true
  chartLoading.value = true
  courseRankLoading.value = true
  chapterRankLoading.value = true
  userRankLoading.value = true
  let pendingDone = false
  let todayDone = false
  const finishStats = () => {
    if (pendingDone && todayDone)
      statsLoading.value = false
  }

  api.home.getPendingStats().then((res) => {
    applyStatsPatch(res.data || {})
  }).catch((err) => {
    console.error('获取待办统计失败:', err)
  }).finally(() => {
    pendingDone = true
    finishStats()
  })

  api.home.getTodayStats().then((res) => {
    applyStatsPatch(res.data || {})
  }).catch((err) => {
    console.error('获取今日学情失败:', err)
  }).finally(() => {
    todayDone = true
    finishStats()
  })

  api.home.getChartData().then((res) => {
    chartData.value = res.data || null
  }).catch((err) => {
    console.error('获取图表数据失败:', err)
  }).finally(() => {
    chartLoading.value = false
  })

  api.home.getCourseRank().then((res) => {
    courseRankList.value = asRankList(res.data)
  }).catch((err) => {
    console.error('获取课程排行失败:', err)
  }).finally(() => {
    courseRankLoading.value = false
  })

  api.home.getChapterRank().then((res) => {
    chapterRankList.value = asRankList(res.data)
  }).catch((err) => {
    console.error('获取课节排行失败:', err)
  }).finally(() => {
    chapterRankLoading.value = false
  })

  api.home.getCreditRank().then((res) => {
    userRankList.value = asRankList(res.data)
  }).catch((err) => {
    console.error('获取近两周积分排行失败:', err)
  }).finally(() => {
    userRankLoading.value = false
  })

  api.home.getCdnTraffic().then((res) => {
    cdnData.value = res.data || null
  }).catch((err) => {
    console.error('获取 CDN 流量失败:', err)
    cdnData.value = { enabled: false, message: 'CDN 流量获取失败' }
  }).finally(() => {
    cdnLoading.value = false
  })
}

onMounted(() => {
  // 一言优先请求，避免被看板重查询挤占
  api.home.getVerse().then((res) => {
    verse.value = res.data?.verse || ''
  }).catch((err) => {
    console.error('获取一言失败:', err)
  })
  loadDashboard()
})

let skipActivateDash = true
onActivated(() => {
  if (skipActivateDash) {
    skipActivateDash = false
    return
  }
  loadDashboard()
})
</script>

<style lang="scss" scoped>
.dash {
  display: grid;
  grid-template-columns: repeat(12, minmax(0, 1fr));
  column-gap: 16px;
  row-gap: 16px;
  margin: 0;
}

.dash > .page-mast {
  padding: 4px 0 0;
}

.dash-span-12,
.dash :deep(.n-spin-container) {
  grid-column: 1 / -1;
}

.dash :deep(.n-spin-content) {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.hero {
  --lf-x: 50%;
  --lf-y: 45%;
  --lf-nx: 0.5;
  --lf-ny: 0.45;
  --apple-ease: cubic-bezier(0.32, 0.72, 0, 1);
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 208px;
  margin-bottom: 0;
  padding: 44px 44px 40px;
  border-radius: 22px;
  border: 1px solid var(--ctl-glass-border, rgb(255 255 255 / 74%));
  background: rgb(252 252 252 / 92%);
  backdrop-filter: none;
  -webkit-backdrop-filter: none;
  box-shadow:
    inset 0 1px 0 rgb(255 255 255 / 86%),
    0 0 0 1px rgba(var(--primary-color), 0.05);
  transition: box-shadow 0.4s var(--apple-ease);

  &::before {
    content: '';
    position: absolute;
    inset: 0;
    z-index: 1;
    pointer-events: none;
    border-radius: inherit;
    background:
      radial-gradient(200px 160px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 42%), transparent 62%),
      radial-gradient(320px 240px at var(--lf-x) var(--lf-y), rgba(var(--primary-color), 0.22), transparent 70%);
    mix-blend-mode: soft-light;
    opacity: 0;
    transition: opacity 0.35s ease;
  }

  &.is-lit {
    box-shadow:
      inset 0 1px 0 rgb(255 255 255 / 80%),
      inset calc((var(--lf-nx) - 0.5) * 36px) calc((var(--lf-ny) - 0.5) * 28px) 70px rgb(255 255 255 / 28%);

    &::before {
      opacity: 1;
    }
  }
}

.hero__aurora {
  border-radius: inherit;
  z-index: 0;
}

.hero__verse {
  position: relative;
  z-index: 2;
  width: 100%;
  max-width: 760px;
  margin: 0;
  padding: 0;
  border: 0;
  text-align: center;
}

.hero__quote {
  position: relative;
  display: inline-block;
  max-width: 100%;
  padding: 12px 40px 8px;
}

.hero__mark {
  position: absolute;
  font-family: var(--font-serif, Georgia, 'Songti SC', 'Noto Serif SC', serif);
  font-size: 72px;
  line-height: 1;
  color: #d4b56a;
  opacity: 0.42;
  pointer-events: none;
  user-select: none;

  &--l {
    top: -10px;
    left: 0;
  }

  &--r {
    right: 0;
    bottom: -28px;
  }
}

.hero__quote p {
  position: relative;
  z-index: 1;
  margin: 0;
  font-family: 'Songti SC', 'Noto Serif SC', 'Source Han Serif SC', 'STSong', Georgia, serif;
  font-size: 22px;
  font-weight: 500;
  line-height: 1.9;
  letter-spacing: 0.08em;
  color: #243f43;
  text-shadow: 0 1px 0 rgb(255 255 255 / 70%);
}

.hero__verse--empty .hero__quote p {
  color: #7a8e91;
  font-size: 16px;
  letter-spacing: 0.12em;
}

.hero__rule {
  display: block;
  width: 72px;
  height: 1.5px;
  margin: 18px auto 0;
  border-radius: 999px;
  background: linear-gradient(90deg, transparent, #c9972a 30%, #e0c06a 50%, #c9972a 70%, transparent);
  opacity: 0.85;
}

.panel {
  --lf-x: 88%;
  --lf-y: 16%;
  --lf-nx: 0.88;
  --lf-ny: 0.16;
  --apple-ease: cubic-bezier(0.32, 0.72, 0, 1);
  position: relative;
  overflow: hidden;
  margin-bottom: 0;
  padding: 16px 22px;
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  border-radius: 20px;
  background:
    radial-gradient(150px 120px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 75%), transparent 62%),
    radial-gradient(240px 180px at var(--lf-x) var(--lf-y), rgba(var(--primary-color), 0.16), transparent 70%),
    radial-gradient(220px 110px at 100% 0%, rgba(var(--primary-color), 0.1), transparent 68%),
    linear-gradient(165deg, rgb(255 255 255 / 72%) 0%, rgb(244 250 249 / 70%) 100%);
  backdrop-filter: none;
  -webkit-backdrop-filter: none;
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 45%);
  transition: box-shadow 0.4s var(--apple-ease), border-color 0.3s ease;

  &.is-lit {
    border-color: rgba(var(--primary-color), 0.16);
    box-shadow:
      inset 0 1px 0 rgb(255 255 255 / 55%),
      0 0 0 3px rgba(var(--primary-color), 0.08);
  }

  &::before {
    content: '';
    position: absolute;
    right: -48px;
    top: -56px;
    width: 160px;
    height: 160px;
    border-radius: 50%;
    background: rgba(var(--primary-color), 0.06);
    pointer-events: none;
    transition:
      transform 0.45s cubic-bezier(0.22, 1, 0.36, 1),
      opacity 0.35s ease;
    will-change: transform;
  }

  &::after {
    content: '';
    position: absolute;
    right: 36px;
    bottom: -42px;
    width: 88px;
    height: 88px;
    border-radius: 50%;
    background: rgba(var(--primary-color), 0.14);
    pointer-events: none;
    transition:
      transform 0.45s cubic-bezier(0.22, 1, 0.36, 1),
      opacity 0.35s ease;
    will-change: transform;
  }

  &:hover::before,
  &.is-lit::before {
    transform:
      translate(
        calc((var(--lf-nx) - 0.5) * -32px),
        calc((var(--lf-ny) - 0.5) * 24px)
      )
      scale(1.1);
  }

  &:hover::after,
  &.is-lit::after {
    transform:
      translate(
        calc((var(--lf-nx) - 0.5) * 24px),
        calc((var(--lf-ny) - 0.5) * -20px)
      )
      scale(1.14);
  }

  > * {
    position: relative;
    z-index: 1;
  }
}

.panel--todo {
  background:
    radial-gradient(150px 120px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 72%), transparent 62%),
    radial-gradient(240px 180px at var(--lf-x) var(--lf-y), rgb(201 123 99 / 14%), transparent 70%),
    radial-gradient(240px 120px at 100% 0%, rgb(201 123 99 / 8%), transparent 70%),
    linear-gradient(180deg, #fffdfb 0%, #fff 60%);

  &::before {
    background: rgb(201 123 99 / 6%);
  }

  &::after {
    background: rgb(201 151 42 / 8%);
  }
}

.panel--chart {
  background:
    radial-gradient(150px 120px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 75%), transparent 62%),
    radial-gradient(240px 180px at var(--lf-x) var(--lf-y), rgba(var(--primary-color), 0.22), transparent 70%),
    radial-gradient(240px 120px at 100% 0%, rgba(var(--primary-color), 0.22), transparent 68%),
    linear-gradient(180deg, #eef7f6 0%, #fff 58%);
}

.panel--credit {
  background:
    radial-gradient(150px 120px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 75%), transparent 62%),
    radial-gradient(240px 180px at var(--lf-x) var(--lf-y), rgb(201 151 42 / 18%), transparent 70%),
    radial-gradient(240px 120px at 100% 0%, rgb(201 151 42 / 20%), transparent 68%),
    linear-gradient(180deg, #fbf6ea 0%, #fff 58%);

  &::before {
    background: rgb(201 151 42 / 12%);
  }

  &::after {
    background: rgb(224 192 106 / 18%);
  }

  .panel__head span {
    color: #c9972a;
    background: rgb(201 151 42 / 12%);
  }
}

.panel--enroll {
  background:
    radial-gradient(150px 120px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 75%), transparent 62%),
    radial-gradient(240px 180px at var(--lf-x) var(--lf-y), rgba(var(--primary-color), 0.18), transparent 70%),
    radial-gradient(240px 120px at 100% 0%, rgba(var(--primary-color), 0.16), transparent 68%),
    linear-gradient(180deg, #eef6f2 0%, #fff 58%);

  &::before {
    background: rgb(73 130 118 / 8%);
  }

  &::after {
    background: rgba(var(--primary-color), 0.16);
  }
}

.panel--chapter {
  background:
    radial-gradient(150px 120px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 75%), transparent 62%),
    radial-gradient(240px 180px at var(--lf-x) var(--lf-y), rgb(139 155 180 / 18%), transparent 70%),
    radial-gradient(240px 120px at 100% 0%, rgb(139 155 180 / 18%), transparent 68%),
    linear-gradient(180deg, #f3f5fa 0%, #fff 58%);

  &::before {
    background: rgb(139 155 180 / 12%);
  }

  &::after {
    background: rgb(139 155 180 / 10%);
  }

  .panel__head span {
    color: #5d6d82;
    background: rgb(139 155 180 / 14%);
  }
}

.sec {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 2px 2px 4px;
  font-family: var(--font-display);
  font-size: 13px;
  font-weight: 600;
  letter-spacing: -0.01em;
  color: #5f7172;

  &::before {
    content: '';
    width: 6px;
    height: 6px;
    border-radius: 2px;
    background: linear-gradient(135deg, rgb(var(--primary-color)), rgba(var(--primary-color), 0.72));
    box-shadow: 0 0 0 3px rgba(var(--primary-color), 0.1);
  }
}

.cdn-block {
  display: flex;
  flex-direction: column;
}

.cdn-shell {
  position: relative;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  gap: 18px;
  padding: 18px 18px 16px;
  border-radius: 22px;
  border: 1px solid rgb(47 127 134 / 10%);
  background:
    radial-gradient(420px 180px at 0% 0%, rgb(47 127 134 / 8%), transparent 60%),
    radial-gradient(360px 160px at 100% 0%, rgb(196 122 58 / 7%), transparent 55%),
    linear-gradient(180deg, #f7fbfb 0%, #fff 42%);
  box-shadow: inset 0 1px 0 #fff;
}

.cdn-shell__head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 12px;
  padding: 0 2px;
}

.cdn-shell__titles {
  display: flex;
  flex-direction: column;
  gap: 4px;

  .sec {
    margin: 0;
  }
}

.cdn-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.cdn-section__label {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  width: fit-content;
  padding: 4px 10px 4px 8px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
  color: #3d4f51;
  background: rgb(255 255 255 / 86%);
  border: 1px solid rgb(47 127 134 / 10%);

  i {
    font-size: 14px;
    color: #2f7f86;
  }

  em {
    margin-left: 2px;
    padding: 1px 7px;
    border-radius: 999px;
    font-style: normal;
    font-size: 11px;
    font-weight: 500;
    color: #2f7f86;
    background: rgb(47 127 134 / 10%);
  }
}

.cdn-kpi {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.cdn-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.cdn-card {
  position: relative;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 6px;
  min-height: 112px;
  padding: 16px 18px;
  border-radius: 18px;
  border: 1px solid transparent;
  box-shadow: inset 0 1px 0 #fff;
  transition: transform 0.18s ease, box-shadow 0.18s ease;

  &::before {
    content: '';
    position: absolute;
    right: -24px;
    top: -28px;
    width: 88px;
    height: 88px;
    border-radius: 50%;
    pointer-events: none;
  }

  .cdn-card__ico {
    position: absolute;
    right: 14px;
    top: 14px;
    z-index: 1;
    display: grid;
    place-items: center;
    width: 32px;
    height: 32px;
    border-radius: 10px;
    font-size: 15px;
    background: rgb(255 255 255 / 72%);
    border: 1px solid rgb(255 255 255 / 80%);
  }

  span {
    position: relative;
    z-index: 1;
    font-size: 13px;
    color: #6b7879;
  }

  strong {
    position: relative;
    z-index: 1;
    margin-top: 2px;
    font-size: 26px;
    font-weight: 700;
    letter-spacing: -0.03em;
    font-variant-numeric: tabular-nums;
    line-height: 1.15;
  }

  em {
    position: relative;
    z-index: 1;
    margin-top: auto;
    font-style: normal;
    font-size: 12px;
    color: #8a9697;
  }

  &:hover {
    transform: translateY(-2px);
    box-shadow:
      inset 0 1px 0 #fff,
      0 12px 28px rgb(31 42 43 / 7%);
  }
}

.cdn-card--bw-today {
  border-color: rgb(47 127 134 / 14%);
  background:
    radial-gradient(140px 90px at 88% 12%, rgb(47 127 134 / 18%), transparent 70%),
    linear-gradient(180deg, #eef7f6 0%, #fff 62%);

  &::before { background: rgb(47 127 134 / 10%); }
  strong { color: #2f7f86; }
  .cdn-card__ico { color: #2f7f86; }
}

.cdn-card--req-today {
  border-color: rgb(196 122 58 / 16%);
  background:
    radial-gradient(140px 90px at 88% 12%, rgb(196 122 58 / 16%), transparent 70%),
    linear-gradient(180deg, #fbf4ec 0%, #fff 62%);

  &::before { background: rgb(196 122 58 / 12%); }
  strong { color: #c47a3a; }
  .cdn-card__ico { color: #c47a3a; }
}

.cdn-card--bw-month {
  border-color: rgb(73 130 118 / 16%);
  background:
    radial-gradient(140px 90px at 88% 12%, rgb(73 130 118 / 16%), transparent 70%),
    linear-gradient(180deg, #eef6f3 0%, #fff 62%);

  &::before { background: rgb(73 130 118 / 10%); }
  strong { color: #498276; }
  .cdn-card__ico { color: #498276; }
}

.cdn-card--req-month {
  border-color: rgb(93 109 130 / 16%);
  background:
    radial-gradient(140px 90px at 88% 12%, rgb(93 109 130 / 14%), transparent 70%),
    linear-gradient(180deg, #f2f4f8 0%, #fff 62%);

  &::before { background: rgb(93 109 130 / 10%); }
  strong { color: #5d6d82; }
  .cdn-card__ico { color: #5d6d82; }
}

.panel--cdn {
  margin: 0;
  border-radius: 18px;
  border: 1px solid rgb(255 255 255 / 70%);
  box-shadow: inset 0 1px 0 #fff;
  overflow: hidden;
}

.panel--cdn-bw {
  background:
    radial-gradient(160px 120px at 100% 0%, rgb(47 127 134 / 14%), transparent 68%),
    linear-gradient(180deg, #eef7f6 0%, #fff 58%);

  .panel__head span,
  .panel__head p {
    color: #2f7f86;
  }
}

.panel--cdn-req {
  background:
    radial-gradient(160px 120px at 100% 0%, rgb(196 122 58 / 14%), transparent 68%),
    linear-gradient(180deg, #fbf4ec 0%, #fff 58%);

  .panel__head span,
  .panel__head p {
    color: #c47a3a;
  }
}

/* 近 30 天趋势：对齐 KPI「30 天流量 / 30 天请求」 */
.panel--cdn-bw-month {
  background:
    radial-gradient(160px 120px at 100% 0%, rgb(73 130 118 / 14%), transparent 68%),
    linear-gradient(180deg, #eef6f3 0%, #fff 58%);

  .panel__head span,
  .panel__head p {
    color: #498276;
  }
}

.panel--cdn-req-month {
  background:
    radial-gradient(160px 120px at 100% 0%, rgb(93 109 130 / 14%), transparent 68%),
    linear-gradient(180deg, #f2f4f8 0%, #fff 58%);

  .panel__head span,
  .panel__head p {
    color: #5d6d82;
  }
}

.cdn-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  min-height: 140px;
  padding: 28px 16px;
  border-radius: 18px;
  border: 1px dashed rgba(var(--primary-color), 0.16);
  text-align: center;
  color: #8a9697;
  font-size: 13px;
  background: rgb(255 255 255 / 65%);

  p {
    margin: 0;
    max-width: 36em;
    line-height: 1.5;
  }
}

.cdn-empty__ico {
  font-size: 22px;
  opacity: 0.55;
  color: #2f7f86;
}

@media (max-width: 1100px) {
  .cdn-kpi {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .cdn-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .cdn-shell {
    padding: 14px;
    border-radius: 18px;
  }

  .cdn-shell__head {
    flex-direction: column;
    align-items: flex-start;
  }

  .cdn-kpi {
    grid-template-columns: 1fr;
  }
}

.panel__head {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0 0 8px;
  padding: 0 12px;

  h3 {
    margin: 0;
    font-family: var(--font-display);
    font-size: 13px;
    font-weight: 560;
    letter-spacing: -0.01em;
    color: var(--ctl-text-soft, #5a6b6d);
  }

  p {
    margin: 0;
    font-size: 12px;
    color: #8a9697;
  }

  span {
    margin-left: auto;
    padding: 2px 8px;
    border-radius: 999px;
    font-size: 11px;
    color: rgb(var(--primary-color));
    background: rgba(var(--primary-color), 0.08);
  }

  em {
    margin-left: auto;
    padding: 3px 10px;
    border-radius: 999px;
    font-style: normal;
    font-size: 12px;
    color: #c97b63;
    background: rgb(201 123 99 / 12%);
    box-shadow: inset 0 1px 0 rgb(255 255 255 / 40%);
  }
}

.panel__more {
  appearance: none;
  margin-left: auto;
  padding: 2px 8px;
  border: 0;
  border-radius: 999px;
  background: rgba(var(--primary-color), 0.08);
  color: rgb(var(--primary-color));
  font-size: 11px;
  font-weight: 560;
  cursor: pointer;

  &:hover {
    background: rgba(var(--primary-color), 0.14);
  }

  &:active {
    transform: scale(0.97);
  }
}

.todo {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 12px;
}

.todo__item {
  --lf-x: 78%;
  --lf-y: 22%;
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  gap: 14px;
  min-height: 92px;
  padding: 16px 16px 16px 14px;
  border: 1px solid transparent;
  border-radius: 16px;
  text-align: left;
  color: inherit;
  cursor: default;
  box-shadow: inset 0 1px 0 #fff;
  transition: transform 0.18s ease, box-shadow 0.18s ease, border-color 0.18s ease;

  &::before,
  &::after {
    content: '';
    position: absolute;
    border-radius: 50%;
    pointer-events: none;
  }

  &::before {
    right: -28px;
    top: -36px;
    width: 92px;
    height: 92px;
  }

  &::after {
    right: 22px;
    bottom: -28px;
    width: 54px;
    height: 54px;
  }
}

.todo__ico,
.todo__meta {
  position: relative;
  z-index: 1;
}

.todo__ico {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  border-radius: 14px;
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 80%);

  i {
    font-size: 18px;
  }
}

.todo__meta {
  display: flex;
  min-width: 0;
  flex-direction: column;
}

.todo__item strong {
  display: flex;
  align-items: baseline;
  font-size: 26px;
  font-weight: 720;
  line-height: 1.1;
  letter-spacing: -0.03em;
  font-variant-numeric: tabular-nums;
}

.todo__item b {
  margin-top: 4px;
  font-size: 13px;
  font-weight: 650;
  color: #2d3a3b;
}

.todo__item small {
  margin-top: 2px;
  font-size: 12px;
  color: #8a9697;
}

button.todo__item {
  cursor: pointer;
  appearance: none;

  &:hover {
    transform: translateY(-2px);
  }

  &:active {
    transform: translateY(0) scale(0.99);
  }
}

.todo__item--coral {
  border-color: rgb(201 123 99 / 14%);
  background:
    radial-gradient(120px 90px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 70%), transparent 62%),
    radial-gradient(180px 140px at var(--lf-x) var(--lf-y), rgb(201 123 99 / 16%), transparent 70%),
    radial-gradient(140px 70px at 100% 0%, rgb(201 123 99 / 18%), transparent 68%),
    linear-gradient(180deg, #fff6f2 0%, #fff 62%);

  &::before { background: rgb(201 123 99 / 12%); }
  &::after { background: rgb(201 123 99 / 10%); }

  .todo__ico {
    background:
      linear-gradient(180deg, rgb(255 255 255 / 80%), transparent 62%),
      rgb(201 123 99 / 16%);

    i { color: #c97b63; }
  }

  strong { color: #c97b63; }

  &:hover {
    border-color: rgb(201 123 99 / 32%);
    box-shadow: 0 12px 22px rgb(201 123 99 / 14%);
  }

  &.is-on {
    background:
      radial-gradient(120px 90px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 70%), transparent 62%),
      radial-gradient(180px 140px at var(--lf-x) var(--lf-y), rgb(201 123 99 / 20%), transparent 70%),
      radial-gradient(140px 70px at 100% 0%, rgb(201 123 99 / 26%), transparent 68%),
      linear-gradient(180deg, #ffeee8 0%, #fff 62%);
  }
}

.todo__item--gold {
  border-color: rgb(201 151 42 / 16%);
  background:
    radial-gradient(120px 90px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 70%), transparent 62%),
    radial-gradient(180px 140px at var(--lf-x) var(--lf-y), rgb(201 151 42 / 16%), transparent 70%),
    radial-gradient(140px 70px at 100% 0%, rgb(201 151 42 / 20%), transparent 68%),
    linear-gradient(180deg, #fbf6ea 0%, #fff 62%);

  &::before { background: rgb(201 151 42 / 14%); }
  &::after { background: rgb(224 192 106 / 16%); }

  .todo__ico {
    background:
      linear-gradient(180deg, rgb(255 255 255 / 80%), transparent 62%),
      rgb(201 151 42 / 16%);

    i { color: #c9972a; }
  }

  strong { color: #c9972a; }

  &:hover {
    border-color: rgb(201 151 42 / 36%);
    box-shadow: 0 12px 22px rgb(201 151 42 / 16%);
  }

  &.is-on {
    background:
      radial-gradient(120px 90px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 70%), transparent 62%),
      radial-gradient(180px 140px at var(--lf-x) var(--lf-y), rgb(201 151 42 / 22%), transparent 70%),
      radial-gradient(140px 70px at 100% 0%, rgb(201 151 42 / 28%), transparent 68%),
      linear-gradient(180deg, #f8efd8 0%, #fff 62%);
  }
}

.todo__item--paper {
  border-color: rgb(139 155 180 / 16%);
  background:
    radial-gradient(120px 90px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 70%), transparent 62%),
    radial-gradient(180px 140px at var(--lf-x) var(--lf-y), rgb(139 155 180 / 16%), transparent 70%),
    radial-gradient(140px 70px at 100% 0%, rgb(139 155 180 / 20%), transparent 68%),
    linear-gradient(180deg, #f3f5fa 0%, #fff 62%);

  &::before { background: rgb(139 155 180 / 14%); }
  &::after { background: rgb(139 155 180 / 10%); }

  .todo__ico {
    background:
      linear-gradient(180deg, rgb(255 255 255 / 80%), transparent 62%),
      rgb(139 155 180 / 18%);

    i { color: #5d6d82; }
  }

  strong { color: #5d6d82; }

  &:hover {
    border-color: rgb(139 155 180 / 34%);
    box-shadow: 0 12px 22px rgb(93 109 130 / 12%);
  }

  &.is-on {
    background:
      radial-gradient(120px 90px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 70%), transparent 62%),
      radial-gradient(180px 140px at var(--lf-x) var(--lf-y), rgb(139 155 180 / 22%), transparent 70%),
      radial-gradient(140px 70px at 100% 0%, rgb(139 155 180 / 28%), transparent 68%),
      linear-gradient(180deg, #e8edf5 0%, #fff 62%);
  }
}

.todo__item--feedback {
  border-color: rgb(66 134 245 / 14%);
  background:
    radial-gradient(120px 90px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 70%), transparent 62%),
    radial-gradient(180px 140px at var(--lf-x) var(--lf-y), rgb(66 134 245 / 16%), transparent 70%),
    radial-gradient(140px 70px at 100% 0%, rgb(66 134 245 / 22%), transparent 68%),
    linear-gradient(180deg, #eef4ff 0%, #fff 62%);

  &::before { background: rgb(66 134 245 / 12%); }
  &::after { background: rgb(66 134 245 / 10%); }

  .todo__ico {
    background:
      linear-gradient(180deg, rgb(255 255 255 / 80%), transparent 62%),
      rgb(66 134 245 / 14%);

    i { color: #4286f5; }
  }

  strong { color: #4286f5; }

  &:hover {
    border-color: rgb(66 134 245 / 32%);
    box-shadow: 0 12px 22px rgb(66 134 245 / 12%);
  }

  &.is-on {
    background:
      radial-gradient(120px 90px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 70%), transparent 62%),
      radial-gradient(180px 140px at var(--lf-x) var(--lf-y), rgb(66 134 245 / 22%), transparent 70%),
      radial-gradient(140px 70px at 100% 0%, rgb(66 134 245 / 28%), transparent 68%),
      linear-gradient(180deg, #e4eeff 0%, #fff 62%);
  }
}

.todo__item--teal {
  border-color: rgb(49 108 114 / 12%);
  background:
    radial-gradient(120px 90px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 70%), transparent 62%),
    radial-gradient(180px 140px at var(--lf-x) var(--lf-y), rgb(126 196 200 / 18%), transparent 70%),
    radial-gradient(140px 70px at 100% 0%, rgb(126 196 200 / 24%), transparent 68%),
    linear-gradient(180deg, #eef7f6 0%, #fff 62%);

  &::before { background: rgb(49 108 114 / 10%); }
  &::after { background: rgb(126 196 200 / 16%); }

  .todo__ico {
    background:
      linear-gradient(180deg, rgb(255 255 255 / 80%), transparent 62%),
      rgb(49 108 114 / 12%);

    i { color: #316c72; }
  }

  strong { color: #316c72; }

  &:hover {
    border-color: rgb(49 108 114 / 28%);
    box-shadow: 0 12px 22px rgb(49 108 114 / 10%);
  }

  &.is-on {
    background:
      radial-gradient(120px 90px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 70%), transparent 62%),
      radial-gradient(180px 140px at var(--lf-x) var(--lf-y), rgb(126 196 200 / 24%), transparent 70%),
      radial-gradient(140px 70px at 100% 0%, rgb(126 196 200 / 30%), transparent 68%),
      linear-gradient(180deg, #e4f3f2 0%, #fff 62%);
  }
}

.todo__split {
  margin: 0 3px;
  opacity: 0.38;
  font-weight: 500;
}

.todo__link {
  padding: 0;
  border: 0;
  background: none;
  color: inherit;
  font: inherit;
  cursor: pointer;

  &:hover {
    text-decoration: underline;
  }
}

.bento {
  display: grid;
  grid-template-columns: repeat(12, minmax(0, 1fr));
  gap: 16px;
  margin-bottom: 0;
}

.bento__feat {
  grid-column: span 5;
  width: 100%;
  min-width: 0;
}

.bento__mid {
  grid-column: span 4;
}

.bento__side {
  grid-column: span 3;
}

.bento__feat,
.bento__tile {
  position: relative;
  overflow: hidden;
  border: 1px solid rgba(var(--primary-color), 0.08);
  border-radius: 20px;
  background: #fff;
  box-shadow: inset 0 1px 0 #fff;
}

.bento__feat {
  --lf-x: 78%;
  --lf-y: 22%;
  --lf-nx: 0.78;
  --lf-ny: 0.22;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: flex-start;
  min-height: 184px;
  padding: 18px 20px 12px;
  text-align: left;
  background:
    radial-gradient(150px 120px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 72%), transparent 62%),
    radial-gradient(220px 160px at var(--lf-x) var(--lf-y), rgba(var(--primary-color), 0.24), transparent 70%),
    radial-gradient(180px 80px at 90% 0%, rgba(var(--primary-color), 0.28), transparent 70%),
    linear-gradient(180deg, #e8f4f3 0%, #fff 58%);
  transition: box-shadow 0.35s ease;

  > span {
    position: relative;
    z-index: 2;
    flex-shrink: 0;
    font-size: 13px;
    color: #6b7879;
  }

  &::before {
    content: '';
    position: absolute;
    right: -36px;
    top: -48px;
    width: 140px;
    height: 140px;
    border-radius: 50%;
    background: rgba(var(--primary-color), 0.06);
    pointer-events: none;
    transition:
      transform 0.45s cubic-bezier(0.22, 1, 0.36, 1),
      opacity 0.35s ease;
    will-change: transform;
  }

  &::after {
    content: '';
    position: absolute;
    right: 28px;
    bottom: -36px;
    width: 78px;
    height: 78px;
    border-radius: 50%;
    background: rgba(var(--primary-color), 0.16);
    pointer-events: none;
    transition:
      transform 0.45s cubic-bezier(0.22, 1, 0.36, 1),
      opacity 0.35s ease;
    will-change: transform;
  }

  &:hover::before,
  &.is-lit::before {
    transform:
      translate(
        calc((var(--lf-nx) - 0.5) * -28px),
        calc((var(--lf-ny) - 0.5) * 20px)
      )
      scale(1.1);
  }

  &:hover::after,
  &.is-lit::after {
    transform:
      translate(
        calc((var(--lf-nx) - 0.5) * 22px),
        calc((var(--lf-ny) - 0.5) * -18px)
      )
      scale(1.14);
  }

  &.is-lit {
    box-shadow:
      inset 0 1px 0 #fff,
      inset calc((var(--lf-nx) - 0.5) * 24px) calc((var(--lf-ny) - 0.5) * 18px) 42px rgb(255 255 255 / 34%);
  }
}

.bento__feat-num {
  position: relative;
  z-index: 2;
  display: flex;
  flex-shrink: 0;
  align-items: baseline;
  gap: 8px;
  margin-top: 8px;

  strong {
    font-size: 42px;
    font-weight: 720;
    letter-spacing: -0.04em;
    color: rgb(var(--primary-color));
    font-variant-numeric: tabular-nums;
  }

  em {
    font-style: normal;
    font-size: 13px;
    color: #8a9697;
  }
}

.bento__spark {
  position: absolute;
  inset: auto 0 0 0;
  height: 88px;
  width: 100%;
  z-index: 1;
  display: block;
  pointer-events: none;
}

.bento__tile--spark {
  justify-content: flex-start;
  min-height: 132px;
  padding-bottom: 52px;

  span,
  strong,
  em {
    z-index: 2;
  }

  .bento__spark {
    height: 56px;
  }
}

.bento__mid {
  display: grid;
  gap: 12px;
}

.bento__side {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  min-height: 184px;

  .bento__tile {
    justify-content: center;
    min-width: 0;
  }
}

.bento__tile {
  --lf-x: 82%;
  --lf-y: 20%;
  --lf-nx: 0.82;
  --lf-ny: 0.2;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 16px 18px;
  text-align: left;
  cursor: default;
  color: inherit;
  background:
    radial-gradient(120px 90px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 72%), transparent 62%),
    radial-gradient(180px 140px at var(--lf-x) var(--lf-y), rgba(var(--primary-color), 0.2), transparent 70%),
    radial-gradient(140px 70px at 96% -8%, rgba(var(--primary-color), 0.24), transparent 70%),
    linear-gradient(180deg, #f3faf9 0%, #fff 60%);
  transition: box-shadow 0.35s ease;

  &::before {
    content: '';
    position: absolute;
    right: -28px;
    top: -36px;
    width: 100px;
    height: 100px;
    border-radius: 50%;
    background: rgba(var(--primary-color), 0.06);
    pointer-events: none;
    transition:
      transform 0.4s cubic-bezier(0.22, 1, 0.36, 1),
      opacity 0.35s ease;
    will-change: transform;
  }

  &::after {
    content: '';
    position: absolute;
    right: 18px;
    bottom: -28px;
    width: 58px;
    height: 58px;
    border-radius: 50%;
    background: rgba(var(--primary-color), 0.16);
    pointer-events: none;
    transition:
      transform 0.4s cubic-bezier(0.22, 1, 0.36, 1),
      opacity 0.35s ease;
    will-change: transform;
  }

  &:hover::before,
  &.is-lit::before {
    transform:
      translate(
        calc((var(--lf-nx) - 0.5) * -22px),
        calc((var(--lf-ny) - 0.5) * 16px)
      )
      scale(1.12);
  }

  &:hover::after,
  &.is-lit::after {
    transform:
      translate(
        calc((var(--lf-nx) - 0.5) * 18px),
        calc((var(--lf-ny) - 0.5) * -14px)
      )
      scale(1.16);
  }

  &.is-lit {
    box-shadow:
      inset 0 1px 0 #fff,
      inset calc((var(--lf-nx) - 0.5) * 20px) calc((var(--lf-ny) - 0.5) * 16px) 36px rgb(255 255 255 / 34%);
  }

  span,
  strong,
  em {
    position: relative;
    z-index: 1;
  }

  span {
    font-size: 13px;
    color: #6b7879;
  }

  strong {
    margin-top: 6px;
    font-size: 28px;
    font-weight: 700;
    letter-spacing: -0.03em;
    color: #1f2a2b;
    font-variant-numeric: tabular-nums;
  }

  em {
    margin-top: 4px;
    font-style: normal;
    font-size: 12px;
    color: #8a9697;
  }
}

.bento__tile--gold {
  border-color: rgb(201 151 42 / 14%);
  background:
    radial-gradient(120px 90px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 72%), transparent 62%),
    radial-gradient(180px 140px at var(--lf-x) var(--lf-y), rgb(201 151 42 / 18%), transparent 70%),
    radial-gradient(140px 70px at 96% -8%, rgb(201 151 42 / 22%), transparent 70%),
    linear-gradient(180deg, #fbf6ea 0%, #fff 60%);

  &::before { background: rgb(201 151 42 / 12%); }
  &::after { background: rgb(224 192 106 / 18%); }

  strong { color: #c9972a; }
}

.bento__tile--slate {
  border-color: rgb(139 155 180 / 16%);
  background:
    radial-gradient(120px 90px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 72%), transparent 62%),
    radial-gradient(180px 140px at var(--lf-x) var(--lf-y), rgb(139 155 180 / 18%), transparent 70%),
    radial-gradient(140px 70px at 96% -8%, rgb(139 155 180 / 22%), transparent 70%),
    linear-gradient(180deg, #f3f5fa 0%, #fff 60%);

  &::before { background: rgb(139 155 180 / 14%); }
  &::after { background: rgb(139 155 180 / 12%); }

  strong { color: #5d6d82; }
}

.bento__tile--teal {
  border-color: rgba(var(--primary-color), 0.14);
  background:
    radial-gradient(120px 90px at var(--lf-x) var(--lf-y), rgb(255 255 255 / 72%), transparent 62%),
    radial-gradient(180px 140px at var(--lf-x) var(--lf-y), rgba(var(--primary-color), 0.22), transparent 70%),
    radial-gradient(140px 70px at 96% -8%, rgba(var(--primary-color), 0.26), transparent 70%),
    linear-gradient(180deg, #e8f4f3 0%, #fff 60%);

  &::before { background: rgba(var(--primary-color), 0.08); }
  &::after { background: rgba(var(--primary-color), 0.18); }

  strong { color: rgb(var(--primary-color)); }
}

button.bento__feat {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: flex-start;
  cursor: pointer;
  width: 100%;
  text-align: left;
  appearance: none;
  color: inherit;
  font: inherit;

  &:hover {
    box-shadow: 0 0 0 3px rgba(var(--primary-color), 0.1);
  }

  &:active {
    transform: scale(0.99);
  }
}

button.bento__tile {
  cursor: pointer;
  transition: border-color 0.18s ease, transform 0.18s ease, box-shadow 0.18s ease;

  &:hover {
    transform: none;
  }

  &:active {
    transform: scale(0.98);
  }
}

button.bento__tile--gold:hover {
  border-color: rgb(201 151 42 / 36%);
  box-shadow: 0 0 0 3px rgb(201 151 42 / 14%);
}

button.bento__tile--slate:hover {
  border-color: rgb(139 155 180 / 36%);
  box-shadow: 0 0 0 3px rgb(93 109 130 / 12%);
}

button.bento__tile--teal:hover {
  border-color: rgba(var(--primary-color), 0.28);
  box-shadow: 0 0 0 3px rgba(var(--primary-color), 0.12);
}

.dash__row {
  display: grid;
  grid-template-columns: repeat(12, minmax(0, 1fr));
  gap: 16px;
  margin-bottom: 0;
}

.dash__row > .panel--chart {
  grid-column: span 8;
}

.dash__row > .panel--credit {
  grid-column: span 4;
}

.dash__row--2 > * {
  grid-column: span 6;
}

.panel--chart,
.panel--credit {
  margin-bottom: 0;
}

.chart {
  position: relative;
  z-index: 1;
  width: 100% !important;
  height: 280px !important;
  min-height: 280px;
}

.chart--spin,
.chart--empty {
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart--empty {
  color: #8a9697;
  font-size: 13px;
}

.credit {
  display: grid;
  grid-template-columns: 158px 1fr;
  gap: 14px;
  align-items: center;
}

.podium {
  display: flex;
  align-items: flex-end;
  justify-content: center;
  gap: 8px;
  padding: 4px 0 10px;
}

.podium__avatar {
  position: relative;
}

.podium__col {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 0;

  &.is-link {
    cursor: pointer;
    border-radius: 12px;
    padding: 4px 6px;
    margin: -4px -6px;
    transition: background 0.18s ease, transform 0.18s ease;

    &:hover {
      background: rgba(var(--primary-color), 0.05);
    }

    &:active {
      transform: scale(0.97);
    }
  }

  :deep(.n-avatar) {
    box-shadow: 0 0 0 3px #fff, 0 0 0 5px rgba(var(--primary-color), 0.14);
  }

  small {
    max-width: 76px;
    margin-top: 7px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 11px;
    color: #6b7879;
  }

  b {
    margin-top: 8px;
    width: 40px;
    height: 30px;
    border-radius: 8px 8px 4px 4px;
    font-size: 12px;
    line-height: 30px;
    text-align: center;
    color: #fff;
    background: linear-gradient(180deg, #a8aeb6, #8b929b);
    box-shadow: inset 0 1px 0 rgb(255 255 255 / 35%);
  }

  &.is-1 {
    order: 2;

    .podium__avatar::after {
      content: '';
      position: absolute;
      inset: -5px;
      border-radius: 50%;
      border: 2px solid rgb(201 151 42 / 55%);
      pointer-events: none;
    }

    :deep(.n-avatar) {
      box-shadow: 0 0 0 3px #fff, 0 0 0 5px rgb(201 151 42 / 45%);
    }

    b {
      width: 46px;
      height: 50px;
      line-height: 50px;
      background: linear-gradient(180deg, #ead089, #c9972a);
      box-shadow: inset 0 1px 0 rgb(255 255 255 / 45%);
    }
  }

  &.is-2 {
    order: 1;

    b {
      background: linear-gradient(180deg, #c5ccd4, #8f99a6);
    }
  }

  &.is-3 {
    order: 3;

    b {
      background: linear-gradient(180deg, #d09a62, #b87333);
    }
  }
}

.credit__list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.credit__row {
  display: grid;
  grid-template-columns: 16px minmax(0, 1fr) 56px 52px;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  margin: 0 -6px;
  border-radius: 10px;
  font-size: 12px;
  transition:
    background 0.18s ease,
    transform 0.18s cubic-bezier(0.32, 0.72, 0, 1);

  &:hover {
    background: rgba(var(--primary-color), 0.04);
  }

  &:active {
    transform: scale(0.995);
  }

  &.is-link {
    cursor: pointer;

    &:hover {
      background: rgba(var(--primary-color), 0.08);
    }
  }

  em {
    font-style: normal;
    font-weight: 700;
    color: #9aa4a5;

    &.is-gold { color: #c9972a; }
    &.is-silver { color: #7b8490; }
    &.is-bronze { color: #b87333; }
  }

  > span:not(.bar) {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: #2d3a3b;
  }

  strong {
    text-align: right;
    font-weight: 600;
    color: #6b7879;
    font-variant-numeric: tabular-nums;
  }

  &:nth-child(1) .bar span {
    background: linear-gradient(90deg, #e0c06a, #c9972a);
  }
}

.rank__row:nth-child(1) .bar span {
  background: linear-gradient(90deg, #e0c06a, #c9972a);
}

.rank__row:nth-child(2) .bar span {
  background: linear-gradient(90deg, #b7c4c8, #7b8490);
}

.rank__row:nth-child(3) .bar span {
  background: linear-gradient(90deg, #d09a62, #b87333);
}

.bar {
  display: block;
  height: 5px;
  overflow: hidden;
  border-radius: 999px;
  background: rgba(var(--primary-color), 0.08);

  span {
    display: block;
    height: 100%;
    border-radius: inherit;
    background: linear-gradient(90deg, rgba(var(--primary-color), 0.65), rgb(var(--primary-color)));
    transition: width 0.45s cubic-bezier(0.32, 0.72, 0, 1);
  }
}

.rank__row {
  display: grid;
  grid-template-columns: 28px 1fr auto;
  gap: 10px;
  align-items: center;
  padding: 9px 8px;
  margin: 0 -6px;
  border-radius: 12px;
  border-bottom: 0;
  transition:
    background 0.18s ease,
    transform 0.18s cubic-bezier(0.32, 0.72, 0, 1);

  &:hover {
    background: rgba(var(--primary-color), 0.05);
  }

  &:active {
    transform: scale(0.995);
  }

  &.is-link {
    cursor: pointer;

    &:hover {
      background: rgba(var(--primary-color), 0.09);
    }
  }

  > em {
    font-style: normal;
    font-size: 12px;
    font-weight: 700;
    color: #9aa4a5;
    font-variant-numeric: tabular-nums;

    &.is-gold { color: #c9972a; }
    &.is-silver { color: #7b8490; }
    &.is-bronze { color: #b87333; }
  }

  b {
    display: block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 13px;
    font-weight: 600;
    color: #2d3a3b;
  }

  small {
    display: block;
    margin-top: 2px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: #8a9697;
  }

  .bar {
    margin-top: 6px;
  }

  strong {
    font-size: 13px;
    font-variant-numeric: tabular-nums;

    small {
      margin-left: 2px;
      display: inline;
    }
  }
}

.empty {
  padding: 28px 0;
  text-align: center;
  color: #8a9697;
  font-size: 13px;
}

@media (max-width: 1100px) {
  .todo {
    grid-template-columns: 1fr 1fr;
  }

  .bento__feat,
  .bento__mid,
  .bento__side,
  .dash__row > .panel--chart,
  .dash__row > .panel--credit,
  .dash__row--2 > * {
    grid-column: span 6;
  }
}

@media (max-width: 768px) {
  .bento__feat,
  .bento__mid,
  .bento__side,
  .dash__row > .panel--chart,
  .dash__row > .panel--credit,
  .dash__row--2 > * {
    grid-column: 1 / -1;
  }
  .hero {
    min-height: 168px;
    padding: 28px 18px 24px;
  }

  .hero__quote {
    padding: 8px 24px 4px;
  }

  .hero__mark {
    font-size: 52px;
  }

  .hero__quote p {
    font-size: 17px;
    letter-spacing: 0.05em;
  }

  .bento {
    grid-template-columns: 1fr;
  }
}
</style>
