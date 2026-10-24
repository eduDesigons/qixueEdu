import { request } from '@/utils'

export default {
  chapterList: params => request.get('/admin/course/chapter/getChapterList', { params }),
  getChapterDetail: id => request.get(`/admin/course/chapter/getChapterDetail?id=${id}`),
  chapterDelete: id => request.get(`/admin/course/chapter/chapterDelete?id=${id}`),
  chapterSave: data => request.post('/admin/course/chapter/chapterSave', data),
  // 排序：direction 为 'up' 或 'down'
  chapterSort: (id, direction) => request.post('/admin/course/chapter/chapterSort', { id, direction }),

  // 学员
  chapterUser: data => request.post(`/admin/course/chapter/getChapterUserList`, data),
  unlockChapterUser: id => request.post('/admin/course/chapter/unlockChapterUser', { id }),
}
