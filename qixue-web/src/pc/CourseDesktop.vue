<template>
  <div class="pc-course">
    <aside>
      <button type="button" class="pc-nav__item" :class="{ on: level1 < 0 }" @click="$emit('level1', -1)">全部课程</button>
      <div v-for="(item, index) in categories" :key="item.id || index">
        <button type="button" class="pc-nav__item" :class="{ on: level1 === index }" @click="pickParent(index, item)">
          <span>{{ item.name }}</span>
          <i v-if="item.children && item.children.length" :class="{ open: level1 === index }"></i>
        </button>
        <div class="pc-sub" v-if="level1 === index && item.children && item.children.length">
          <button type="button" :class="{ on: level2 === 'all' || level2 == null }" @click="$emit('level2', { id: 'all', name: '全部', isAll: true })">全部</button>
          <button
            type="button"
            v-for="child in item.children"
            :key="child.id"
            :class="{ on: level2 === child.id }"
            @click="$emit('level2', child)"
          >{{ child.name }}</button>
        </div>
      </div>
    </aside>

    <section>
      <div class="pc-bar">
        <header>
          <h1>{{ title }}</h1>
          <div class="pc-sort">
            <button type="button" :class="{ on: sortType === 'hot' }" @click="pickSort('hot')">最热</button>
            <button type="button" :class="{ on: sortType === 'new' }" @click="pickSort('new')">最新</button>
          </div>
        </header>
        <div class="pc-types">
          <button type="button" v-for="item in types" :key="item.value" :class="{ on: courseType === item.value }" @click="pickType(item.value)">{{ item.label }}</button>
          <button type="button" class="pc-key" v-if="keyword" @click="$emit('clear')">{{ keyword }} ×</button>
        </div>
      </div>

      <div class="pc-grid" v-if="list.length">
        <article
          v-for="course in list"
          :key="course.id"
          tabindex="0"
          @click="$emit('open', course)"
          @keyup.enter="$emit('open', course)"
        >
          <div class="pc-cover">
            <img :src="static(course.cover || course.thumb)" alt="" />
            <em v-if="course.category && course.category.name">{{ course.category.name }}</em>
          </div>
          <h3>{{ course.title }}</h3>
          <p>{{ course.author || course.desc }}</p>
          <span>{{ course.credit || 0 }} 学分</span>
        </article>
      </div>

      <div class="pc-skel" v-else-if="loading">
        <i v-for="n in 6" :key="n"></i>
      </div>
      <div class="pc-empty" v-else>
        <p>没有找到相关课程</p>
        <button type="button" v-if="filtered" @click="$emit('reset')">查看全部课程</button>
      </div>

      <div ref="sentinel" class="pc-sentinel"></div>
      <button class="pc-more" type="button" v-if="list.length && !noMore" :disabled="loading" @click="$emit('more')">
        {{ loading ? '加载中…' : '加载更多' }}
      </button>
      <div class="pc-end" v-if="noMore && list.length">共 {{ list.length }} 门</div>
    </section>
  </div>
</template>

<script>
export default {
  name: 'CourseDesktop',
  props: {
    categories: { type: Array, default: () => [] },
    level1: { type: Number, default: -1 },
    level2: { default: 'all' },
    sortType: { type: String, default: 'hot' },
    courseType: { type: Number, default: 0 },
    list: { type: Array, default: () => [] },
    loading: Boolean,
    noMore: Boolean,
    keyword: { type: String, default: '' },
  },
  emits: ['level1', 'level2', 'sort', 'type', 'more', 'open', 'clear', 'reset'],
  computed: {
    parentName() {
      const item = this.categories[this.level1]
      return (item && item.name) || '分类'
    },
    title() {
      if (this.keyword) return '搜索「' + this.keyword + '」'
      if (this.level1 < 0) return '全部课程'
      const parent = this.categories[this.level1]
      const child = parent && parent.children && parent.children.find((row) => row.id === this.level2)
      if (child) return parent.name + ' / ' + child.name
      return this.parentName
    },
    filtered() {
      return !!(this.keyword || this.level1 >= 0 || this.courseType)
    },
  },
  watch: {
    list() { this.bindSentinel() },
    loading(value) { if (!value) this.bindSentinel() },
    noMore() { this.bindSentinel() },
  },
  mounted() {
    this.observer = new IntersectionObserver((entries) => {
      if (entries.some((entry) => entry.isIntersecting) && !this.loading && !this.noMore && this.list.length) {
        this.$emit('more')
      }
    }, { rootMargin: '240px' })
    this.bindSentinel()
  },
  beforeUnmount() {
    if (this.observer) this.observer.disconnect()
  },
  data() {
    return {
      types: [
        { label: '全部', value: 0 },
        { label: '视频课', value: 1 },
        { label: '音频课', value: 2 },
        { label: '图文课', value: 3 },
      ],
    }
  },
  methods: {
    pickParent(index, item) {
      this.$emit('level1', index)
      this.$emit('level2', { id: item.id, name: item.name, isAll: true })
    },
    pickSort(value) {
      if (this.sortType === value) return
      this.$emit('sort', value)
    },
    pickType(value) {
      if (this.courseType === value) return
      this.$emit('type', value)
    },
    bindSentinel() {
      this.$nextTick(() => {
        if (!this.observer) return
        this.observer.disconnect()
        if (this.$refs.sentinel && !this.noMore) this.observer.observe(this.$refs.sentinel)
      })
    },
  },
}
</script>

<style scoped>
.pc-course {
  display: grid;
  grid-template-columns: 220px minmax(0, 1fr);
  gap: 16px;
  align-items: start;
  min-height: calc(100vh - 64px);
  padding: 20px 0 48px;
}
aside {
  position: sticky;
  top: 80px;
  max-height: calc(100vh - 96px);
  overflow: auto;
  padding: 8px;
  border-radius: 12px;
  background: #fff;
}
.pc-nav__item,
.pc-sub button {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  min-height: 38px;
  padding: 8px 12px;
  border-radius: 8px;
  color: #404040;
  font-size: calc(14px * var(--font-scale, 1));
  text-align: left;
  cursor: pointer;
}
.pc-nav__item span {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.pc-nav__item i {
  flex: none;
  width: 6px;
  height: 6px;
  margin-left: 8px;
  border-right: 1.5px solid #c5cad3;
  border-bottom: 1.5px solid #c5cad3;
  transform: rotate(-45deg);
}
.pc-nav__item i.open { transform: rotate(45deg); }
.pc-nav__item:hover,
.pc-sub button:hover { background: #f6f7fb; }
.pc-nav__item.on,
.pc-sub button.on {
  background: #eef3fb;
  color: #4285f4;
  font-weight: 600;
}
.pc-sub { padding: 2px 0 6px 12px; }
.pc-sub button {
  min-height: 34px;
  color: #666;
  font-size: calc(13px * var(--font-scale, 1));
  font-weight: 400;
}
.pc-bar {
  padding: 18px 20px 16px;
  border-radius: 12px;
  background: #fff;
}
header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}
h1 {
  margin: 0;
  min-width: 0;
  color: #1a1a1a;
  font-size: calc(20px * var(--font-scale, 1));
  font-weight: 650;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.pc-sort {
  display: flex;
  flex: none;
  padding: 3px;
  border-radius: 8px;
  background: #f6f7fb;
}
.pc-sort button {
  height: 28px;
  padding: 0 14px;
  border-radius: 6px;
  color: #8b93a1;
  font-size: calc(13px * var(--font-scale, 1));
  cursor: pointer;
}
.pc-sort button.on {
  background: #fff;
  color: #1a1a1a;
  font-weight: 600;
}
.pc-types {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 14px;
}
.pc-types button {
  height: 32px;
  padding: 0 14px;
  border-radius: 16px;
  background: #f6f7fb;
  color: #666;
  font-size: calc(13px * var(--font-scale, 1));
  cursor: pointer;
}
.pc-types button.on {
  background: rgba(66, 133, 244, 0.1);
  color: #4285f4;
}
.pc-key {
  background: #eef3fb !important;
  color: #4285f4 !important;
}
.pc-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
  margin-top: 16px;
}
article {
  min-width: 0;
  padding: 10px 10px 14px;
  border-radius: 12px;
  background: #fff;
  cursor: pointer;
}
article:hover h3,
article:focus-visible h3 { color: #4285f4; }
article:focus-visible { outline: 2px solid #4285f4; outline-offset: 2px; }
.pc-cover {
  position: relative;
  overflow: hidden;
  aspect-ratio: 16 / 10;
  border-radius: 8px;
  background: #eef0f4;
}
.pc-cover img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
  transition: transform 0.25s ease;
}
article:hover .pc-cover img { transform: scale(1.04); }
.pc-cover em {
  position: absolute;
  top: 8px;
  left: 8px;
  max-width: calc(100% - 16px);
  padding: 2px 8px;
  border-radius: 4px;
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  font-size: calc(12px * var(--font-scale, 1));
  font-style: normal;
  line-height: 18px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
h3 {
  margin: 10px 2px 0;
  color: #1a1a1a;
  font-size: calc(15px * var(--font-scale, 1));
  font-weight: 500;
  line-height: 1.45;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
p {
  margin: 6px 2px 0;
  color: #9ca3af;
  font-size: calc(13px * var(--font-scale, 1));
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
article span {
  display: block;
  margin: 8px 2px 0;
  color: #9ca3af;
  font-size: calc(12px * var(--font-scale, 1));
}
.pc-skel {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
  margin-top: 16px;
}
.pc-skel i {
  display: block;
  height: 248px;
  border-radius: 12px;
  background: linear-gradient(90deg, #eef1f6, #f7f8fb, #eef1f6);
  background-size: 200% 100%;
  animation: pc-shine 1.2s infinite;
}
@keyframes pc-shine {
  0% { background-position: 100% 0; }
  100% { background-position: -100% 0; }
}
.pc-sentinel { height: 1px; }
.pc-more,
.pc-empty,
.pc-end {
  margin-top: 16px;
  color: #9ca3af;
  font-size: calc(14px * var(--font-scale, 1));
  text-align: center;
}
.pc-more {
  width: 100%;
  height: 40px;
  border-radius: 8px;
  background: #fff;
  color: #4285f4;
  cursor: pointer;
}
.pc-more:hover { background: #eef3fb; }
.pc-more:disabled { cursor: default; color: #9ca3af; }
.pc-empty {
  padding: 48px 16px;
  border-radius: 12px;
  background: #fff;
}
.pc-empty p { margin: 0; }
.pc-empty button {
  margin-top: 12px;
  color: #4285f4;
  font-size: calc(14px * var(--font-scale, 1));
  cursor: pointer;
}
</style>
