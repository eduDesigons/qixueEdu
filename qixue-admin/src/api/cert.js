import { request } from '@/utils'

export default {
  /** 证书列表（分页） */
  list: data => request.post('/admin/member/cert/list', data),
  /** 证书详情 */
  detail: id => request.get('/admin/member/cert/detail', { params: { id } }),
  /** 保存证书：name、params、template_image、preview_image、is_show；传 id 为编辑。预览图由前端合成上传 */
  save: data => request.post('/admin/member/cert/save', data),
  /** 上架 / 下架 */
  setShow: data => request.post('/admin/member/cert/setShow', data),
  /** 删除证书 */
  delete: id => request.post('/admin/member/cert/delete', { id }),
  /** 证书下拉 */
  options: () => request.get('/admin/member/cert/options'),
  /** 颁发 / 重新生成学员证书 */
  issue: data => request.post('/admin/member/cert/issue', data),
  /** 为已结业学员批量补发（绑定模板后确认） */
  backfill: data => request.post('/admin/member/cert/backfill', data),
  /** 已颁发列表（可传 course_id / path_id） */
  issuedList: data => request.post('/admin/member/cert/issuedList', data),
  /** 撤回学员证书 */
  revoke: data => request.post('/admin/member/cert/revoke', data),
}
