import { request } from '@/utils'

export default {
  send: data => request.post('/admin/member/message/send', data),
  list: data => request.post('/admin/member/message/list', data),
  delete: id => request.post('/admin/member/message/delete', { id }),
}
