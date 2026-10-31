import { request } from '@/utils'

export default {
  courseEnrollList: data => request.post('/admin/ops/ops/courseEnrollList', data),
  pathEnrollList: data => request.post('/admin/ops/ops/pathEnrollList', data),
  learnList: data => request.post('/admin/ops/ops/learnList', data),
  creditList: data => request.post('/admin/ops/ops/creditList', data),
  coinsList: data => request.post('/admin/ops/ops/coinsList', data),
  chart: data => request.post('/admin/ops/ops/chart', data),
  inspectKinds: data => request.post('/admin/ops/ops/inspectKinds', data),
  inspect: data => request.post('/admin/ops/ops/inspect', data, { timeout: 90000 }),
  repair: data => request.post('/admin/ops/ops/repair', data, { timeout: 90000 }),
}
