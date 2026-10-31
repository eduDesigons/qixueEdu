<template>
  <div class="pc-learn">
    <aside class="pc-learn__side">
      <section class="pc-learn__tasks" v-if="tasks.length">
        <header>
          <h2>待完成</h2>
          <span>{{ tasks.length }}</span>
        </header>
        <button type="button" v-for="item in tasks" :key="item.chapterId || item.paperId" @click="openTask(item)">
          <strong>{{ item.title }}</strong>
          <em>
            <span v-if="isGradeTask(item)" class="pc-learn__grade">待批改</span>
            <span v-if="item.desc">{{ item.desc }}</span>
          </em>
        </button>
      </section>
      <rank-card title="学习时长 Top5" kind="user" v-model:period="rankPeriod.duration" :list="durationList" :loading="rankLoading && !hasRankData" @tip="$emit('tip', 'duration')" @history="$emit('history', 'duration')" />
      <rank-card title="学习课程 Top5" kind="course" v-model:period="rankPeriod.course" :list="courseList" :loading="rankLoading && !hasRankData" @tip="$emit('tip', 'course')" @history="$emit('history', 'course')" />
      <rank-card title="积分 Top5" kind="user" v-model:period="rankPeriod.coins" :list="coinList" :loading="rankLoading && !hasRankData" @tip="$emit('tip', 'coins')" @history="$emit('history', 'coins')" />
    </aside>

    <div class="pc-learn__main">
      <section class="pc-learn__stats">
        <h2>学习概况</h2>
        <div class="pc-learn__stat-grid">
          <div>
            <strong>{{ showDuration(studyData.time) || '0分钟' }}</strong>
            <span>学习时长</span>
          </div>
          <div>
            <strong>{{ studyData.courseNum || 0 }}</strong>
            <span>学完课程</span>
          </div>
          <div>
            <strong>{{ studyData.chapterNum || 0 }}</strong>
            <span>学完课节</span>
          </div>
          <div>
            <strong>{{ studyData.todayNum || 0 }}</strong>
            <span>今日完成</span>
          </div>
        </div>
      </section>

      <section class="pc-learn__panel">
        <header>
          <h2>最近学习</h2>
          <button type="button" @click="navTo('/pages/member/lesson')">查看全部</button>
        </header>
        <div class="pc-learn__empty" v-if="!ready">正在加载…</div>
        <div class="pc-learn__empty" v-else-if="!courses.length">
          还没有正在学习的课程
          <button type="button" @click="navTo('/pages/index/course')">去选课</button>
        </div>
        <div v-else>
          <article v-for="item in courses" :key="item.id" @click="$emit('continue', item)">
            <img :src="static(item.course && item.course.thumb)" alt="" />
            <div>
              <h3>{{ (item.course && item.course.title) || '课程已下架' }}</h3>
              <p>{{ continueText(item) }}</p>
            </div>
            <span>{{ item && Number(item.is_end) === 1 ? '查看' : '继续' }}</span>
          </article>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import RankCard from '@/pages/index/components/RankCard.vue'

export default {
  name: 'LearnDesktop',
  components: { RankCard },
  props: {
    studyData: { type: Object, default: () => ({}) },
    courses: { type: Array, default: () => [] },
    tasks: { type: Array, default: () => [] },
    ready: Boolean,
    rankPeriod: { type: Object, default: () => ({}) },
    rankLoading: Boolean,
    hasRankData: Boolean,
    durationList: { type: Array, default: () => [] },
    courseList: { type: Array, default: () => [] },
    coinList: { type: Array, default: () => [] },
  },
  emits: ['continue', 'tip', 'history'],
  methods: {
    isGradeTask(item) {
      return Number(item && item.isPeopleCheck) === 1 && Number(item.userPaperId) > 0
    },
    openTask(item) {
      if (this.isGradeTask(item)) {
        this.navTo('/pages/exam/examDetail?id=' + item.userPaperId)
        return
      }
      if (item && item.paperId) {
        this.navTo('/pages/exam/exam?id=' + item.paperId)
      }
    },
    continueText(item) {
      const ch = item && item.continueChapter
      if (item && Number(item.is_end) === 1) return '已学完，查看课程'
      if (ch && ch.title) return (Number(item.has_end_chapter) > 0 ? '续学 · ' : '开始 · ') + ch.title
      if (item && item.has_end_chapter > 0) return '已学完' + item.has_end_chapter + '讲'
      return '开始学习'
    },
  },
}
</script>

<style scoped>
.pc-learn {
  display: grid;
  grid-template-columns: 360px minmax(0, 1fr);
  gap: 16px;
  align-items: start;
  padding: 20px 0 48px;
}
.pc-learn__side {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.pc-learn__main {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}
.pc-learn__stats,
.pc-learn__tasks,
.pc-learn__panel {
  border-radius: 12px;
  background: #fff;
}
.pc-learn__stats {
  padding: 20px 22px 18px;
}
.pc-learn__stats h2,
.pc-learn__tasks h2,
.pc-learn__panel h2 {
  margin: 0;
  font-size: calc(18px * var(--font-scale, 1));
  font-weight: 650;
  color: #1a1a1a;
}
.pc-learn__stat-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  margin-top: 16px;
}
.pc-learn__stat-grid > div {
  padding: 16px 14px 14px;
  border-radius: 10px;
  background: #f6f8fb;
}
.pc-learn__stat-grid strong {
  display: block;
  color: #1a1a1a;
  font-size: calc(22px * var(--font-scale, 1));
  font-weight: 700;
  line-height: 1.25;
}
.pc-learn__stat-grid span {
  display: block;
  margin-top: 6px;
  color: #8b93a1;
  font-size: calc(13px * var(--font-scale, 1));
}
.pc-learn__tasks { padding: 16px 16px 8px; }
.pc-learn__tasks header,
.pc-learn__panel header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}
.pc-learn__tasks header span {
  color: #4285f4;
  font-size: calc(14px * var(--font-scale, 1));
}
.pc-learn__tasks button {
  display: flex;
  width: 100%;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
  padding: 12px 0;
  border-top: 1px solid #f0f2f5;
  text-align: left;
  cursor: pointer;
}
.pc-learn__tasks button strong {
  color: #1a1a1a;
  font-size: calc(14px * var(--font-scale, 1));
  font-weight: 500;
}
.pc-learn__tasks button em {
  color: #8b93a1;
  font-size: calc(12px * var(--font-scale, 1));
  font-style: normal;
}
.pc-learn__grade {
  display: inline-block;
  margin-right: 6px;
  color: #ff3359;
  font-style: normal;
  font-weight: 500;
}
.pc-learn__tasks button:hover strong { color: #4285f4; }
.pc-learn__panel { padding: 18px 20px 8px; }
.pc-learn__panel header button {
  color: #8b93a1;
  font-size: calc(14px * var(--font-scale, 1));
  cursor: pointer;
}
.pc-learn__panel header button:hover { color: #4285f4; }
article {
  display: grid;
  grid-template-columns: 168px minmax(0, 1fr) auto;
  gap: 14px;
  align-items: center;
  padding: 14px 0;
  border-top: 1px solid #f0f2f5;
  cursor: pointer;
}
article img {
  width: 168px;
  height: 105px;
  border-radius: 8px;
  object-fit: cover;
  background: #eef0f4;
  display: block;
}
article h3 {
  margin: 0;
  color: #1a1a1a;
  font-size: calc(15px * var(--font-scale, 1));
  font-weight: 500;
  line-height: 1.45;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
article:hover h3 { color: #4285f4; }
article p {
  margin: 6px 0 0;
  color: #9ca3af;
  font-size: calc(13px * var(--font-scale, 1));
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
article span {
  color: #4285f4;
  font-size: calc(14px * var(--font-scale, 1));
  white-space: nowrap;
}
.pc-learn__empty {
  padding: 8px 0 16px;
  color: #8b93a1;
  font-size: calc(14px * var(--font-scale, 1));
}
.pc-learn__empty button {
  margin-left: 12px;
  color: #4285f4;
  cursor: pointer;
}
.pc-learn__side :deep(.rank-card) {
  margin-top: 0;
  padding: 16px 14px 8px;
  border-radius: 12px;
  box-shadow: none;
}
.pc-learn__side :deep(.rank-card-title) {
  font-size: calc(15px * var(--font-scale, 1));
  line-height: 1.3;
}
.pc-learn__side :deep(.rank-card-switch-item) {
  min-width: 46px;
  height: 26px;
  line-height: 26px;
  font-size: calc(12px * var(--font-scale, 1));
}
.pc-learn__side :deep(.rank-item) {
  min-height: 44px;
  padding: 4px 2px;
}
.pc-learn__side :deep(.rank-name) {
  font-size: calc(13px * var(--font-scale, 1));
}
.pc-learn__side :deep(.rank-value),
.pc-learn__side :deep(.rank-value-num) {
  font-size: calc(13px * var(--font-scale, 1));
}
</style>
