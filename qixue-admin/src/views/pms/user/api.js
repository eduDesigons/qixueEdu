/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:29:51
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { request } from '@/utils'

export default {
  create: data => request.post('/admin/pms/admin/addAdmin', data),
  read: data => request.post('/admin/pms/admin/adminList', data),
  update: data => request.post(`/admin/pms/admin/editAdmin`, data),
  delete: id => request.post(`/admin/pms/admin/delAdmin`, { id }),
  resetPwd: data => request.post(`/admin/pms/admin/resetPassword`, data),

  getAllRoles: () => request.get('/admin/pms/role/getAllRole'),
}
