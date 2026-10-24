import { request } from '@/utils'

export default {
  getAreaList: () => request.get('/admin/common/area/getAreaList'),
  getAreaTree: () => request.get('/admin/common/area/getAreaTree'),
  getSysInfo: () => request.get('/admin/common/sys/getSysInfo'),
}
