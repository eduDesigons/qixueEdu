import { request } from '@/utils'

export default {
  list: data => request.post('/admin/member/feedback/list', data),
  reply: data => request.post('/admin/member/feedback/reply', data),
}
