<template>
  <section class="pc-hero">
    <button class="pc-hero__back" type="button" @click="navTo('/pages/index/course')">返回课程</button>
    <div class="pc-hero__body">
      <img :src="static(info.thumb || info.cover || info.img_url)" alt="" />
      <div>
        <em v-if="info.category && info.category.name">{{ info.category.name }}</em>
        <h1>{{ info.title }}</h1>
        <p>{{ info.desc }}</p>
        <ul>
          <li><strong>{{ chapters.length || info.chapter_count || 0 }}</strong><span>讲</span></li>
          <li><strong>{{ info.credit || 0 }}</strong><span>学分</span></li>
          <li><strong>{{ info.type_name || '课程' }}</strong><span>形式</span></li>
          <li><strong>{{ info.view_count || 0 }}</strong><span>浏览</span></li>
        </ul>
        <div class="pc-hero__actions">
          <button v-if="isJoin" type="button" @click="$emit('primary')">{{ label || '继续学习' }}</button>
          <button v-else-if="applyPending" type="button" disabled>审核中，请耐心等待</button>
          <button v-else-if="!chapters.length" type="button" disabled>暂无课节</button>
          <button v-else type="button" @click="$emit('join')">{{ canJoin ? '加入课程' : '申请加入' }}</button>
          <button class="pc-hero__ghost" type="button" @click="$emit('collect')">{{ isCollect ? '已收藏' : '收藏' }}</button>
        </div>
        <small v-if="isJoin && tip">{{ tip }}</small>
        <small v-else-if="!isJoin && !canJoin && info.min_need_credit">要求最低学分：{{ info.min_need_credit }}</small>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  name: 'CourseHero',
  props: {
    info: { type: Object, required: true },
    chapters: { type: Array, default: () => [] },
    isJoin: Boolean,
    canJoin: Boolean,
    applyPending: Boolean,
    isCollect: Boolean,
    label: String,
    tip: String,
  },
  emits: ['primary', 'join', 'collect'],
}
</script>

<style scoped>
.pc-hero { margin-bottom: 8px; padding: 8px 0 12px; }
.pc-hero__back { margin-bottom: 12px; color: #4286f5; font-size: calc(14px * var(--font-scale, 1)); }
.pc-hero__body {
  display: grid;
  grid-template-columns: 360px minmax(0, 1fr);
  gap: 28px;
  padding: 20px;
  border-radius: 12px;
  background: #fff;
}
img { width: 100%; height: 220px; object-fit: cover; border-radius: 10px; background: #e8edf3; }
em { color: #fa8919; font-size: calc(13px * var(--font-scale, 1)); font-style: normal; }
h1 { margin: 8px 0; font-size: calc(26px * var(--font-scale, 1)); line-height: 1.35; font-weight: 700; }
p { margin: 0; color: #6b7280; font-size: calc(14px * var(--font-scale, 1)); line-height: 1.6; }
ul { display: flex; gap: 28px; margin: 18px 0; padding: 0; list-style: none; }
li strong, li span { display: block; }
li strong { font-size: calc(18px * var(--font-scale, 1)); }
li span { margin-top: 2px; color: #8b93a1; font-size: calc(12px * var(--font-scale, 1)); }
.pc-hero__actions { display: flex; gap: 10px; }
button[disabled] { background: #e5e7eb; color: #9ca3af; }
.pc-hero__actions > button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 40px;
  padding: 0 22px;
  border-radius: 8px;
  background: #4286f5;
  color: #fff;
  font-size: calc(15px * var(--font-scale, 1));
  line-height: 1;
}
.pc-hero__ghost { background: #f3f5f8 !important; color: #353535 !important; }
small { display: block; margin-top: 10px; color: #8b93a1; }
</style>
