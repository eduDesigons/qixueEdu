<template>
  <div class="pc-me-home">
    <section class="pc-me-card">
      <h2>学习概况</h2>
      <div class="pc-me-stats">
        <button type="button" @click="navTo('/pages/member/lesson')">
          <img src="/static/img/me/exam.png" alt="" />
          <em>
            <strong>{{ stats.courseCount || 0 }}</strong>
            <span>我的课程</span>
          </em>
        </button>
        <button type="button" @click="navTo('/pages/member/credit')">
          <svg viewBox="0 0 24 24" aria-hidden="true">
            <rect width="24" height="24" rx="7" fill="#FFE8D6" />
            <circle cx="12" cy="10" r="4" fill="#F5A623" />
            <path d="M9.4 13.2 8.2 19.2 12 17.2l3.8 2-1.2-6" fill="#F08A24" />
          </svg>
          <em>
            <strong>{{ stats.credit || 0 }}</strong>
            <span>我的学分</span>
          </em>
        </button>
        <button type="button" @click="navTo('/pages/member/coin')">
          <svg viewBox="0 0 24 24" aria-hidden="true">
            <rect width="24" height="24" rx="7" fill="#FFF4D6" />
            <circle cx="12" cy="12" r="6.2" fill="#F5B400" />
            <path d="M12 8.4v7.2M9.7 10.3h3.1a1.5 1.5 0 0 1 0 3h-2.2a1.5 1.5 0 0 0 0 3H14.4" fill="none" stroke="#fff" stroke-width="1.4" stroke-linecap="round" />
          </svg>
          <em>
            <strong>{{ stats.coins || 0 }}</strong>
            <span>我的积分</span>
          </em>
        </button>
      </div>
      <p v-if="unread > 0" class="pc-me-unread">
        你有 {{ unread > 99 ? '99+' : unread }} 条未读消息
        <button type="button" @click="navTo('/pages/member/message')">去查看</button>
      </p>
    </section>

    <section class="pc-me-card">
      <header class="pc-me-card__head">
        <div>
          <h2>学习热力</h2>
          <p>{{ heatTip }}</p>
        </div>
        <button type="button" class="pc-me-link" @click="navTo('/pages/member/history')">查看轨迹</button>
      </header>
      <div class="heat">
        <div class="heat__months">
          <span v-for="(label, index) in heatMonths" :key="'m' + index">{{ label }}</span>
        </div>
        <div class="heat__body">
          <div class="heat__days">
            <span>一</span>
            <span></span>
            <span>三</span>
            <span></span>
            <span>五</span>
            <span></span>
            <span></span>
          </div>
          <div class="heat__grid" @mouseleave="heatTip = heatSummary">
            <i
              v-for="cell in heatCells"
              :key="cell.key"
              :class="{ future: cell.future }"
              :style="cell.future ? null : { background: heatColor(cell.count) }"
              @mouseenter="heatTip = cell.tip"
            />
          </div>
        </div>
        <div class="heat__legend">
          <span>少</span>
          <i v-for="color in heatScale" :key="color" :style="{ background: color }" />
          <span>多</span>
        </div>
      </div>
    </section>

    <div class="pc-me-split">
      <section class="pc-me-card">
        <header class="pc-me-card__head">
          <div>
            <h2>学习曲线</h2>
            <p>{{ curveTip }}</p>
          </div>
        </header>
        <svg class="curve" viewBox="0 0 640 188" role="img" aria-label="近12周学完课节">
          <line
            v-for="line in curveGrid"
            :key="line.y"
            x1="32"
            :x2="628"
            :y1="line.y"
            :y2="line.y"
          />
          <span v-for="line in curveGrid" :key="'t' + line.y" x="0" :y="line.y + 4">{{ line.label }}</span>
          <path :d="curveArea" />
          <path class="curve__line" :d="curveLine" />
          <g v-for="point in curvePoints" :key="point.label" @mouseenter="curveTip = point.tip" @mouseleave="curveTip = curveSummary">
            <circle :cx="point.x" :cy="point.y" r="8" class="curve__hit" />
            <circle :cx="point.x" :cy="point.y" r="3.5" class="curve__dot" />
          </g>
          <span
            v-for="point in curveLabels"
            :key="'x' + point.label"
            :x="point.x"
            y="180"
            text-anchor="middle"
          >{{ point.label }}</span>
        </svg>
      </section>

      <section class="pc-me-card">
        <header class="pc-me-card__head">
          <div>
            <h2>一周分布</h2>
            <p>近 17 周各天学完的课节</p>
          </div>
        </header>
        <ul class="week">
          <li v-for="item in weekBars" :key="item.label">
            <span>{{ item.label }}</span>
            <b><i :style="{ width: item.pct + '%' }" /></b>
            <em>{{ item.count }}</em>
          </li>
        </ul>
      </section>
    </div>
  </div>
</template>

<script>
import { getStorage, setStorage } from '@/common/storage'
const WEEKS = 17
const CURVE_WEEKS = 12
const WEEK_LABELS = ['一', '二', '三', '四', '五', '六', '日']
const HEAT_COLORS = ['#eef1f6', '#d6e4ff', '#9dbfff', '#6aa1fa', '#4285f4']

function startOfDay(date) {
  return new Date(date.getFullYear(), date.getMonth(), date.getDate())
}

function addDays(date, days) {
  const next = new Date(date)
  next.setDate(next.getDate() + days)
  return next
}

function pad(value) {
  return String(value).padStart(2, '0')
}

function dayKey(date) {
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`
}

function mondayOf(date) {
  const day = startOfDay(date)
  return addDays(day, -((day.getDay() + 6) % 7))
}

export default {
  name: 'MeDesktop',
  props: {
    stats: { type: Object, default: () => ({}) },
    unread: { type: Number, default: 0 },
  },
  data() {
    return {
      dayMap: {},
      heatTip: '近 17 周学完课节',
      curveTip: '近 12 周，每周学完的课节',
      trendLoading: false,
      trendAt: 0,
    }
  },
  computed: {
    heatScale() {
      return HEAT_COLORS
    },
    windowStart() {
      return addDays(mondayOf(new Date()), -(WEEKS - 1) * 7)
    },
    heatCells() {
      const today = dayKey(startOfDay(new Date()))
      const cells = []
      for (let index = 0; index < WEEKS * 7; index += 1) {
        const date = addDays(this.windowStart, index)
        const key = dayKey(date)
        const future = key > today
        const count = future ? 0 : (this.dayMap[key] || 0)
        cells.push({
          key,
          count,
          future,
          month: date.getMonth(),
          tip: future ? '还没到这一天' : `${date.getMonth() + 1}月${date.getDate()}日 · 学完 ${count} 讲`,
        })
      }
      return cells
    },
    heatMonths() {
      return Array.from({ length: WEEKS }, (_, week) => {
        const cell = this.heatCells[week * 7]
        const prev = week ? this.heatCells[(week - 1) * 7] : null
        if (!cell || (prev && prev.month === cell.month)) return ''
        return `${cell.month + 1}月`
      })
    },
    heatMax() {
      return this.heatCells.reduce((max, cell) => Math.max(max, cell.count), 0)
    },
    heatTotal() {
      return this.heatCells.reduce((sum, cell) => sum + cell.count, 0)
    },
    heatSummary() {
      return this.heatTotal ? `近 17 周学完 ${this.heatTotal} 讲` : '近 17 周还没有学完的课节'
    },
    curveWeeks() {
      const endMonday = mondayOf(new Date())
      return Array.from({ length: CURVE_WEEKS }, (_, index) => {
        const start = addDays(endMonday, -(CURVE_WEEKS - 1 - index) * 7)
        let count = 0
        for (let day = 0; day < 7; day += 1) count += this.dayMap[dayKey(addDays(start, day))] || 0
        const end = addDays(start, 6)
        return {
          count,
          label: `${start.getMonth() + 1}/${start.getDate()}`,
          tip: `${start.getMonth() + 1}月${start.getDate()}日 – ${end.getMonth() + 1}月${end.getDate()}日 · ${count} 讲`,
        }
      })
    },
    curveMax() {
      return Math.max(1, ...this.curveWeeks.map((item) => item.count))
    },
    curvePoints() {
      const padL = 32
      const padR = 12
      const padT = 16
      const padB = 28
      const width = 640 - padL - padR
      const height = 188 - padT - padB
      return this.curveWeeks.map((item, index) => ({
        ...item,
        x: padL + (index / (CURVE_WEEKS - 1)) * width,
        y: padT + (1 - item.count / this.curveMax) * height,
      }))
    },
    curveLine() {
      return this.curvePoints.map((point, index) => `${index ? 'L' : 'M'}${point.x.toFixed(1)},${point.y.toFixed(1)}`).join(' ')
    },
    curveArea() {
      const points = this.curvePoints
      if (!points.length) return ''
      const base = 188 - 28
      return `${this.curveLine} L${points[points.length - 1].x.toFixed(1)},${base} L${points[0].x.toFixed(1)},${base} Z`
    },
    curveGrid() {
      const padT = 16
      const height = 188 - 16 - 28
      return [1, 0.5, 0].map((ratio) => ({
        y: padT + (1 - ratio) * height,
        label: String(Math.round(this.curveMax * ratio)),
      }))
    },
    curveLabels() {
      return this.curvePoints.filter((_, index) => index % 2 === 0)
    },
    curveSummary() {
      const total = this.curveWeeks.reduce((sum, item) => sum + item.count, 0)
      return total ? `近 12 周学完 ${total} 讲` : '近 12 周还没有学完的课节'
    },
    weekBars() {
      const counts = [0, 0, 0, 0, 0, 0, 0]
      this.heatCells.forEach((cell, index) => {
        if (!cell.future) counts[index % 7] += cell.count
      })
      const max = Math.max(1, ...counts)
      return WEEK_LABELS.map((label, index) => ({
        label,
        count: counts[index],
        pct: Math.round((counts[index] / max) * 100),
      }))
    },
  },
  mounted() {
    this.restoreTrend()
    this.heatTip = this.heatSummary
    this.curveTip = this.curveSummary
    this.loadTrend()
  },
  methods: {
    heatColor(count) {
      if (!count) return HEAT_COLORS[0]
      if (this.heatMax <= 1) return HEAT_COLORS[4]
      const step = this.heatMax / 4
      return HEAT_COLORS[Math.min(4, Math.ceil(count / step))]
    },
    restoreTrend() {
      try {
        const cached = getStorage('qixue_me_trend')
        if (cached && cached.days) this.dayMap = cached.days
      } catch (e) {}
    },
    loadTrend() {
      if (!this.$store.getters.hasLogin || this.trendLoading) return
      if (this.trendAt && Date.now() - this.trendAt < 20000) return
      this.trendLoading = true
      this.fetchTrend().finally(() => {
        this.trendLoading = false
        this.trendAt = Date.now()
        this.heatTip = this.heatSummary
        this.curveTip = this.curveSummary
      })
    },
    fetchTrend() {
      const cutoff = dayKey(this.windowStart)
      const days = {}
      const pull = (page) => this.$http.get('/learn/myHistory', { page, limit: 30 }).then((res) => {
        const list = (res.data && res.data.list) || []
        let older = false
        list.forEach((item) => {
          const date = String(item.date || item.happenedAt || '').slice(0, 10)
          if (!/^\d{4}-\d{2}-\d{2}$/.test(date)) return
          if (date < cutoff) {
            older = true
            return
          }
          days[date] = (days[date] || 0) + 1
        })
        const hasMore = !!(res.data && res.data.hasMore)
        if (!hasMore || !list.length || older || page >= 6) {
          this.dayMap = days
          try {
            setStorage('qixue_me_trend', { days })
          } catch (e) {}
          return
        }
        return pull(page + 1)
      })
      return pull(1).catch(() => {})
    },
  },
}
</script>

<style scoped>
.pc-me-home {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.pc-me-card {
  padding: 20px 22px 18px;
  border-radius: 12px;
  background: #fff;
}
h2 {
  margin: 0;
  color: #1a1a1a;
  font-size: calc(18px * var(--font-scale, 1));
  font-weight: 650;
}
.pc-me-stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin-top: 16px;
}
.pc-me-stats button {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px 16px;
  border: 1px solid #eef1f6;
  border-radius: 12px;
  background: #fff;
  text-align: left;
  cursor: pointer;
}
.pc-me-stats button:hover {
  border-color: #d6e4ff;
  background: #f7faff;
}
.pc-me-stats img,
.pc-me-stats svg {
  width: 40px;
  height: 40px;
  flex: none;
  object-fit: contain;
}
.pc-me-stats em { font-style: normal; }
.pc-me-stats strong {
  display: block;
  color: #1a1a1a;
  font-size: calc(26px * var(--font-scale, 1));
  font-weight: 700;
  line-height: 1.1;
}
.pc-me-stats span {
  display: block;
  margin-top: 6px;
  color: #8b93a1;
  font-size: calc(13px * var(--font-scale, 1));
}
.pc-me-unread {
  margin: 16px 0 0;
  color: #404040;
  font-size: calc(14px * var(--font-scale, 1));
}
.pc-me-unread button,
.pc-me-link {
  margin-left: 8px;
  padding: 0;
  border: 0;
  background: transparent;
  color: #4285f4;
  font-size: calc(14px * var(--font-scale, 1));
  cursor: pointer;
}
.pc-me-link { margin-left: 0; flex: none; }
.pc-me-card__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}
.pc-me-card__head p {
  margin: 6px 0 0;
  color: #8b93a1;
  font-size: calc(13px * var(--font-scale, 1));
}
.pc-me-split {
  display: grid;
  grid-template-columns: minmax(0, 1.45fr) minmax(260px, 0.8fr);
  gap: 12px;
}
.heat {
  width: 100%;
}
.heat__months,
.heat__grid {
  display: grid;
  grid-template-columns: repeat(17, minmax(0, 1fr));
  gap: 5px;
  width: 100%;
}
.heat__months {
  margin-left: 26px;
  min-height: 16px;
  color: #9aa3af;
  font-size: calc(12px * var(--font-scale, 1));
}
.heat__months span {
  white-space: nowrap;
  overflow: visible;
}
.heat__body {
  display: grid;
  grid-template-columns: 18px minmax(0, 1fr);
  gap: 8px;
  align-items: stretch;
}
.heat__days {
  display: grid;
  grid-template-rows: repeat(7, minmax(0, 1fr));
  gap: 5px;
  color: #9aa3af;
  font-size: calc(12px * var(--font-scale, 1));
  line-height: 1;
}
.heat__days span {
  display: flex;
  align-items: center;
}
.heat__grid {
  grid-auto-flow: column;
  grid-template-rows: repeat(7, auto);
}
.heat__grid i {
  display: block;
  width: 100%;
  aspect-ratio: 1;
  border-radius: 3px;
  background: #eef1f6;
}
.heat__grid i.future { background: transparent; }
.heat__legend {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 4px;
  margin-top: 10px;
  color: #9aa3af;
  font-size: calc(12px * var(--font-scale, 1));
}
.heat__legend i {
  width: 12px;
  height: 12px;
  border-radius: 2px;
}
.curve {
  display: block;
  width: 100%;
  height: auto;
}
.curve line { stroke: #eef1f6; stroke-width: 1; }
.curve span { fill: #9aa3af; font-size: calc(11px * var(--font-scale, 1)); }
.curve path { fill: rgba(66, 133, 244, 0.14); stroke: none; }
.curve .curve__line {
  fill: none;
  stroke: #4285f4;
  stroke-width: 2;
}
.curve__dot { fill: #4285f4; }
.curve__hit { fill: transparent; }
.week {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin: 4px 0 0;
  padding: 0;
  list-style: none;
}
.week li {
  display: grid;
  grid-template-columns: 16px minmax(0, 1fr) 28px;
  gap: 8px;
  align-items: center;
}
.week span,
.week em {
  color: #8b93a1;
  font-size: calc(12px * var(--font-scale, 1));
  font-style: normal;
}
.week em { text-align: right; color: #404040; }
.week b {
  display: block;
  height: 8px;
  border-radius: 99px;
  background: #eef1f6;
  overflow: hidden;
}
.week i {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: #4285f4;
}
</style>
