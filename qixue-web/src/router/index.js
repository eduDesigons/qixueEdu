import { createRouter, createWebHistory } from 'vue-router'

const modules = import.meta.glob('../pages/**/*.vue')

const pages = [
  ['pages/index/home', '发现', true, true],
  ['pages/index/course', '课程', true, true],
  ['pages/index/learn', '学习', true, true],
  ['pages/index/me', '我的', true, true],
  ['pages/index/login', '登录', true, false],
  ['pages/course/courseDetail', '课程详情', false, false],
  ['pages/course/chapterDetail', '课节详情', false, false],
  ['pages/course/pathDetail', '系列详情', true, false],
  ['pages/course/pathList', '系列课程', false, false],
  ['pages/exam/exam', '考试测评', false, false],
  ['pages/exam/examDetail', '考试测评详情', false, false],
  ['pages/member/message', '消息中心', false, false],
  ['pages/member/rankHistory', '历史榜单', false, false],
  ['pages/member/lesson', '我的课程', false, false],
  ['pages/member/collection', '我的收藏', false, false],
  ['pages/member/history', '学习轨迹', false, false],
  ['pages/member/profile', '基本资料', false, false],
  ['pages/member/cert', '我的证书', false, false],
  ['pages/member/set', '设置', false, false],
  ['pages/member/cache', '清理缓存', false, false],
  ['pages/member/fontSize', '字体大小', false, false],
  ['pages/member/credit', '我的学分', false, false],
  ['pages/member/coin', '我的积分', false, false],
  ['pages/member/exam', '我的考试', false, false],
  ['pages/member/test', '我的测评', false, false],
  ['pages/member/comment', '我的评论', false, false],
  ['pages/member/question', '常见问题', false, false],
  ['pages/member/tools', '工具库', false, false],
  ['pages/member/feedback', '问题/建议', false, false],
  ['pages/member/path', '系列课程', false, false],
  ['pages/member/about', '关于我们', false, false],
  ['pages/h5/h5', '', true, false],
  ['pages/h5/bannerDetail', '', true, false],
]

export const tabNames = pages.filter((item) => item[3]).map((item) => item[0].replace(/\//g, '-'))

const routes = pages.map(([path, title, custom, tab]) => {
  const name = path.replace(/\//g, '-')
  return {
    path: `/${path}`,
    name,
    meta: { title, custom, tab },
    component: () => modules[`../${path}.vue`]().then((mod) => {
      const comp = mod.default
      comp.name = name
      comp.__qxPage = true
      return comp
    }),
  }
})

routes.unshift({ path: '/', redirect: '/pages/index/home' })

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  },
})

export default router
