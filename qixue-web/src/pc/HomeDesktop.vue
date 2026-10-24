<template>
  <div class="pc-home">
    <div class="pc-hero" :class="{ 'is-plain': !loading && !banners.length }">
      <div class="pc-stage" v-if="loading || banners.length" @click="onStage">
        <img v-if="currentBanner" class="pc-banner" :src="currentBanner.img_url"  style="object-fit:cover" />
        <div v-else class="pc-stage__empty"></div>
        <div class="pc-dots" v-if="banners.length > 1">
          <i v-for="(item, index) in banners" :key="index" :class="{ on: index === bannerIndex }" @click.stop="bannerIndex = index"></i>
        </div>
      </div>
      <aside class="pc-side">
        <h2>{{ hello }}</h2>
        <coin-badge v-if="userInfo && userInfo.id" :coins="userInfo.coins" />
        <p>从课程库挑一门，或接着上次的进度学。</p>
        <div class="pc-side__actions">
          <button type="button" class="pc-side__go" @click="navTo('/pages/index/learn')">进入学习</button>
          <button type="button" class="pc-side__browse" @click="navTo('/pages/index/course')">浏览课程</button>
        </div>
      </aside>
    </div>

    <div class="pc-cats" v-if="navs.length">
      <button type="button" class="pc-cats__item" v-for="nav in navs" :key="nav.id" @click="$emit('category', nav.id)">
        <img v-if="nav.icon" :src="static(nav.icon)" class="pc-cats__icon"  style="object-fit:contain" />
        <span>{{ nav.name }}</span>
      </button>
      <button type="button" class="pc-cats__more" @click="navTo('/pages/index/course')">全部课程</button>
    </div>

    <div class="notice" v-if="notice && notice.id" @click="$emit('notice')">
      <div class="notice-con">
        <div class="notice-text">{{ notice.title }}</div>
        <van-icon name="arrow" color="#858a92" size="20" />
      </div>
    </div>

    <section class="pc-panel" v-if="rankTabs.length">
      <header>
        <h2>热学榜</h2>
        <button type="button" @click="navTo('/pages/index/course')">更多</button>
      </header>
      <div class="rank-tabs">
        <div class="rank-tab" :class="{ active: rankTabIndex === index }" v-for="(tab, index) in rankTabs" :key="tab.id || index" @click="$emit('rank-tab', index)">{{ tab.name }}</div>
      </div>
      <div class="rank-tags" v-if="rankTags.length">
        <div class="rank-tag" :class="{ active: rankTagIndex === index }" v-for="(tag, index) in rankTags" :key="tag.id || index" @click="$emit('rank-tag', index)">{{ tag.name }}</div>
      </div>
      <div class="pc-grid" v-if="rankList.length">
        <article class="pc-card" v-for="(item, index) in rankList" :key="item.id || index" @click="openCourse(item)">
          <div class="pc-card__cover">
            <img :src="static(item.thumb)" class="pc-card__img"  style="object-fit:cover" />
            <div class="rank-num" :class="'rank-' + (index + 1)">{{ index + 1 }}</div>
            <div class="pc-card__badge" v-if="item.category && item.category.name">{{ item.category.name }}</div>
          </div>
          <h3>{{ item.title }}</h3>
          <p>{{ item.desc }}</p>
          <em>{{ item.credit || 0 }} 学分</em>
        </article>
      </div>
      <div class="pc-empty" v-else-if="!loading">该分类暂无课程</div>
    </section>

    <section class="pc-panel" v-if="pathList.length">
      <header>
        <h2>系列课程</h2>
        <button type="button" @click="navTo('/pages/course/pathList')">更多</button>
      </header>
      <div class="pc-grid">
        <article class="pc-card" v-for="item in pathList" :key="item.id" @click="navTo('/pages/course/pathDetail?id=' + item.id)">
          <div class="pc-card__cover">
            <img :src="static(item.thumb)" class="pc-card__img"  style="object-fit:cover" />
          </div>
          <h3>{{ item.title }}</h3>
          <p>{{ item.desc }}</p>
        </article>
      </div>
    </section>

    <section class="pc-panel" v-if="hotList.length || hotLoading">
      <header>
        <h2>热门推荐</h2>
        <button type="button" @click="navTo('/pages/index/course')">更多</button>
      </header>
      <div class="pc-grid" v-if="hotList.length">
        <article class="pc-card" v-for="item in hotList" :key="item.id" @click="openCourse(item)">
          <div class="pc-card__cover">
            <img :src="static(item.thumb)" class="pc-card__img"  style="object-fit:cover" />
            <div class="pc-card__badge" v-if="item.category && item.category.name">{{ item.category.name }}</div>
          </div>
          <h3>{{ item.title }}</h3>
          <p>{{ item.desc }}</p>
        </article>
      </div>
      <div class="pc-empty" v-if="hotLoading">加载中...</div>
      <button type="button" class="pc-more" v-if="hotList.length && hotHasMore && !hotLoading" @click="$emit('more-hot')">加载更多</button>
      <div class="pc-empty" v-else-if="hotList.length && !hotHasMore">没有更多了</div>
    </section>
  </div>
</template>

<script>
export default {
  name: 'HomeDesktop',
  props: {
    loading: Boolean,
    banners: { type: Array, default: () => [] },
    navs: { type: Array, default: () => [] },
    notice: { type: Object, default: () => ({}) },
    rankTabs: { type: Array, default: () => [] },
    rankTabIndex: { type: Number, default: 0 },
    rankTags: { type: Array, default: () => [] },
    rankTagIndex: { type: Number, default: 0 },
    rankList: { type: Array, default: () => [] },
    pathList: { type: Array, default: () => [] },
    hotList: { type: Array, default: () => [] },
    hotLoading: Boolean,
    hotHasMore: Boolean,
  },
  emits: ['banner', 'category', 'notice', 'rank-tab', 'rank-tag', 'more-hot'],
  data() {
    return { bannerIndex: 0, timer: 0 }
  },
  computed: {
    currentBanner() {
      return this.banners[this.bannerIndex] || null
    },
    hello() {
      const name = this.userInfo && this.userInfo.nick_name
      return name ? name + '，今天继续学' : '今天学一点'
    },
  },
  watch: {
    banners() {
      this.bannerIndex = 0
      this.play()
    },
  },
  mounted() {
    this.play()
  },
  beforeUnmount() {
    if (this.timer) window.clearInterval(this.timer)
  },
  methods: {
    play() {
      if (this.timer) window.clearInterval(this.timer)
      if (this.banners.length < 2) return
      this.timer = window.setInterval(() => {
        this.bannerIndex = (this.bannerIndex + 1) % this.banners.length
      }, 4000)
    },
    onStage() {
      if (this.currentBanner) this.$emit('banner', this.bannerIndex)
    },
    openCourse(item) {
      if (!item || !item.id) return
      this.navTo('/pages/course/courseDetail?id=' + item.id)
    },
  },
}
</script>

<style scoped lang="scss">
.pc-home {
  padding: 20px 0 48px;
}
.pc-hero {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 320px;
  gap: 16px;
  align-items: stretch;
  min-height: 300px;
}
.pc-stage {
  position: relative;
  height: 300px;
  overflow: hidden;
  border-radius: 12px;
  background: #e8edf3;
  cursor: pointer;
}
.pc-stage__empty,
.pc-banner {
  width: 100%;
  height: 300px;
  display: block;
  object-fit: cover;
}
.pc-stage__empty {
  background: linear-gradient(100deg, #e8edf3 30%, #f7f8fb 50%, #e8edf3 70%);
  background-size: 200% 100%;
  animation: pc-stage-shine 1.1s ease-in-out infinite;
}
@keyframes pc-stage-shine {
  from { background-position: 100% 0; }
  to { background-position: -100% 0; }
}
.pc-hero.is-plain {
  grid-template-columns: 1fr;
  min-height: 0;
}
.pc-hero.is-plain .pc-side {
  flex-direction: row;
  align-items: center;
  gap: 24px;
  min-height: 0;
  padding: 22px 24px;
}
.pc-hero.is-plain .pc-side p {
  flex: 1;
  margin: 0;
}
.pc-hero.is-plain .pc-side__actions {
  flex-direction: row;
  flex-shrink: 0;
  margin-top: 0;
}
.pc-hero.is-plain .pc-side__go,
.pc-hero.is-plain .pc-side__browse {
  min-width: 120px;
  padding: 0 18px;
}
.pc-dots {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 14px;
  display: flex;
  justify-content: center;
  gap: 6px;
}
.pc-dots i {
  width: 7px;
  height: 7px;
  border-radius: 99px;
  background: rgba(255, 255, 255, 0.55);
  cursor: pointer;
}
.pc-dots i.on {
  width: 18px;
  background: #fff;
}
.pc-side {
  display: flex;
  flex-direction: column;
  min-height: 300px;
  padding: 22px 20px 18px;
  border-radius: 12px;
  background: #fff;
}
.pc-side h2 {
  margin: 0;
  font-size: calc(20px * var(--font-scale, 1));
  font-weight: 650;
  line-height: 1.35;
  color: #1a1a1a;
}
.pc-side p {
  margin: 10px 0 0;
  color: #8b93a1;
  font-size: calc(13px * var(--font-scale, 1));
  line-height: 1.7;
}
.pc-side__actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 18px;
}
.pc-side__go,
.pc-side__browse {
  height: 40px;
  border-radius: 8px;
  font-size: calc(15px * var(--font-scale, 1));
  cursor: pointer;
}
.pc-side__go {
  background: #4285f4;
  color: #fff;
}
.pc-side__go:hover { background: #3a78e0; }
.pc-side__browse {
  background: #f3f6fb;
  color: #4285f4;
}
.pc-side__browse:hover { background: #e7eefb; }
.pc-home > .notice { cursor: pointer; }
.pc-cats {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 16px;
  padding: 8px;
  overflow-x: auto;
  border-radius: 12px;
  background: #fff;
}
.pc-cats__item,
.pc-cats__more {
  display: flex;
  align-items: center;
  height: 40px;
  padding: 0 14px;
  border-radius: 8px;
  font-size: calc(14px * var(--font-scale, 1));
  white-space: nowrap;
  cursor: pointer;
}
.pc-cats__item {
  gap: 8px;
  color: #404040;
}
.pc-cats__item:hover {
  background: #f6f7fb;
  color: #4285f4;
}
.pc-cats__icon {
  width: 22px;
  height: 22px;
  object-fit: contain;
}
.pc-cats__more {
  margin-left: auto;
  color: #4285f4;
}
.pc-panel {
  margin-top: 16px;
  padding: 20px 22px 22px;
  border-radius: 12px;
  background: #fff;
}
.pc-panel header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.pc-panel h2 {
  margin: 0;
  font-size: calc(18px * var(--font-scale, 1));
  font-weight: 650;
  color: #1a1a1a;
}
.pc-panel header button {
  color: #8b93a1;
  font-size: calc(14px * var(--font-scale, 1));
  cursor: pointer;
}
.pc-panel header button:hover { color: #4285f4; }
.pc-panel .rank-tabs {
  width: fit-content;
  max-width: 100%;
}
.pc-panel .rank-tab {
  flex: none;
  min-width: 72px;
  cursor: pointer;
}
.pc-panel .rank-tags {
  margin: 14px 0 16px;
  padding: 0;
  overflow-x: auto;
}
.pc-panel .rank-tag { cursor: pointer; }
.pc-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 20px 16px;
}
.pc-card {
  min-width: 0;
  cursor: pointer;
}
.pc-card__cover {
  position: relative;
  overflow: hidden;
  aspect-ratio: 16 / 10;
  border-radius: 8px;
  background: #eef0f4;
}
.pc-card__img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
  transition: transform 0.25s ease;
}
.pc-card:hover .pc-card__img { transform: scale(1.04); }
.pc-card h3 {
  margin: 10px 0 0;
  color: #1a1a1a;
  font-size: calc(15px * var(--font-scale, 1));
  font-weight: 500;
  line-height: 1.45;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.pc-card:hover h3 { color: #4285f4; }
.pc-card p,
.pc-card em {
  margin: 6px 0 0;
  color: #9ca3af;
  font-size: calc(13px * var(--font-scale, 1));
  font-style: normal;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.pc-card__badge {
  position: absolute;
  top: 8px;
  left: 8px;
  max-width: calc(100% - 16px);
  padding: 2px 8px;
  border-radius: 4px;
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  font-size: calc(12px * var(--font-scale, 1));
  line-height: 18px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.pc-card .rank-num + .pc-card__badge {
  left: auto;
  right: 8px;
}
.pc-panel .rank-num {
  position: absolute;
  top: 0;
  left: 0;
  z-index: 1;
  width: 22px;
  height: 26px;
  font-size: calc(13px * var(--font-scale, 1));
}
.pc-empty {
  padding: 8px 0 4px;
  color: #9ca3af;
  font-size: calc(13px * var(--font-scale, 1));
  text-align: center;
}
.pc-more {
  display: block;
  width: 100%;
  margin-top: 16px;
  height: 40px;
  border-radius: 8px;
  background: #f6f7fb;
  color: #4285f4;
  font-size: calc(14px * var(--font-scale, 1));
  cursor: pointer;
}
.pc-more:hover { background: #eef3fb; }
</style>
