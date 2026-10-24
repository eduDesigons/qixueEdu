import { request } from '@/utils'

export default {
  list: () => request.post('/admin/pms/area/list'),
  save: data => request.post('/admin/pms/area/save', data),
  delete: data => request.post('/admin/pms/area/delete', data),
}
