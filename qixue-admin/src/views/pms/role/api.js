/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:29:27
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { request } from '@/utils'

export default {
  create: data => request.post('/admin/pms/role/addRole', data),
  read: params => request.post('/admin/pms/role/roleList', params),
  update: data => request.post(`/admin/pms/role/editRole`, data),
  delete: id => request.post(`/admin/pms/role/delRole`, { id }),

  getAllPermissionTree: () => request.get('/admin/pms/role/getMenuTree'),
  getAllUsers: data => request.post('/admin/pms/admin/adminList', data),
  addRoleUsers: data => request.post(`/admin/pms/admin/addRoleToAdmin`, data),
  removeRoleUsers: data => request.post(`/admin/pms/admin/delRoleFromAdmin`, data),
}
