import { request } from '@/utils'

export default {
  getVerse: () => request.get('/admin/home/dashboard/getVerse', { timeout: 8000 }),
  getStats: () => request.get('/admin/home/dashboard/getStats', { timeout: 15000 }),
  getPendingStats: () => request.get('/admin/home/dashboard/getPendingStats', { timeout: 15000 }),
  getTodayStats: () => request.get('/admin/home/dashboard/getTodayStats', { timeout: 15000 }),
  getChartData: () => request.get('/admin/home/dashboard/getChartData', { timeout: 20000 }),
  getCourseRank: () => request.get('/admin/home/dashboard/getCourseLearningRank', { timeout: 15000 }),
  getChapterRank: () => request.get('/admin/home/dashboard/getChapterLearningRank', { timeout: 30000 }),
  getCreditRank: () => request.get('/admin/home/dashboard/getChapterCompleteUserRank', { timeout: 15000 }),
  getCdnTraffic: () => request.get('/admin/home/dashboard/getCdnTraffic', { timeout: 20000 }),
}
