/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2024/04/01 15:52:04
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import axios from 'axios'
import { request } from '@/utils'

export default {
  getMenuTree: () => request.get('/admin/pms/menu/getMenuTree'),
  getButtons: ({ parentId }) => request.get(`/admin/pms/menu/getButtons?parentId=${parentId}`),
  getComponents: () => axios.get(`${import.meta.env.VITE_PUBLIC_PATH}components.json`),
  addPermission: data => request.post('/admin/pms/menu/addMenu', data),
  savePermission: data => request.post(`/admin/pms/menu/saveMenu`, data),
  deletePermission: id => request.post(`/admin/pms/menu/delMenu`, { id }),
}
