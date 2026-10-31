import { request } from '@/utils'

export default {
  list: data => request.post('/admin/pms/verse/verseList', data),
  save: data => request.post('/admin/pms/verse/save', data),
  delete: data => request.post('/admin/pms/verse/delete', data),
}
