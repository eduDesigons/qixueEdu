import { request } from '@/utils'

export default {
  getInfo: type => request.get('/admin/pms/config/getInfo', { params: { type } }),
  saveInfo: (data, type) => request.post(`/admin/pms/config/saveInfo?type=${type}`, data),
}
