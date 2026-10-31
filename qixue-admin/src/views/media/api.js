/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:29:51
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import api from '@/api'
import { request } from '@/utils'

export default {
  cateList: type => request.post('/admin/media/mediaCate/list', { type }),
  mediaList: data => request.post(`/admin/media/media/list`, data),
  upload: (file, cateId) => api.upload(file, cateId),
  videoTus: file => api.videoTus(file),
  videoTusSuccess: (videoId, cateId) => api.videoTusSuccess(videoId, cateId),
}
