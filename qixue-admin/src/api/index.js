/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/04 22:50:38
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { request } from '@/utils'
import cert from './cert.js'
import chapter from './chapter.js'
import common from './common.js'
import config from './config.js'
import course from './course.js'
import decoration from './decoration.js'
import feedback from './feedback.js'
import home from './home.js'
import message from './message.js'
import ops from './ops.js'
import paper from './paper.js'
import path from './path.js'
import user from './user.js'
import area from './area.js'
import verse from './verse.js'

function putBlob(url, blob, contentType, onProgress) {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    xhr.open('PUT', url)
    xhr.timeout = 0
    if (contentType)
      xhr.setRequestHeader('Content-Type', contentType)
    xhr.upload.onprogress = (e) => {
      onProgress?.(e.loaded || 0)
    }
    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300)
        resolve()
      else
        reject(new Error(xhr.responseText?.slice(0, 200) || `上传失败 HTTP ${xhr.status}`))
    }
    xhr.onerror = () => reject(new Error('网络错误'))
    xhr.ontimeout = () => reject(new Error('上传超时'))
    xhr.send(blob)
  })
}

async function putFileMultipart(token, file, onProgress) {
  if (!Array.isArray(token?.parts) || !token.parts.length)
    throw new Error('未获取到分片上传地址')
  const total = Number(file.size) || 0
  const partSize = Number(token.partSize) || (8 * 1024 * 1024)
  const parts = token.parts
  const concurrency = Math.min(Math.max(Number(token.concurrency) || 6, 1), 8)
  const loadedMap = Array.from({ length: parts.length }, () => 0)
  const report = () => {
    if (typeof onProgress !== 'function' || total <= 0)
      return
    const loaded = loadedMap.reduce((sum, n) => sum + n, 0)
    onProgress(Math.min(99, Math.round((loaded / total) * 100)), loaded, total)
  }

  let cursor = 0
  const runWorker = async () => {
    while (cursor < parts.length) {
      const index = cursor++
      const part = parts[index]
      const start = (part.partNumber - 1) * partSize
      const blob = new Blob([file.slice(start, start + partSize)])
      await putBlob(part.putUrl, blob, '', (loaded) => {
        loadedMap[index] = loaded
        report()
      })
      loadedMap[index] = blob.size
      report()
    }
  }

  report()
  await Promise.all(Array.from({ length: Math.min(concurrency, parts.length) }, runWorker))
}

export default {
  getUser: () => request.get('/admin/login/getUserInfo'),
  logout: () => request.post('/admin/login/logout', {}, { needTip: false }),
  switchCurrentRole: role => request.post(`/auth/current-role/switch/${role}`),
  getRolePermissions: () => request.get('/admin/login/getUserPermissions'),

  upload: (file, cateId, maxEdge = 0) => {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('cateId', cateId)
    if (maxEdge > 0)
      formData.append('maxEdge', String(maxEdge))
    // 不手动设置 Content-Type，让浏览器自动设置（包含 boundary）
    return request({
      url: '/admin/media/media/upload',
      method: 'POST',
      data: formData,
      timeout: 0,
    })
  },
  /**
   * 视频直传：服务端只签发分片 URL，浏览器直传 S3，不经 API 中转文件
   */
  initVideoDirectUpload: ({ fileName, fileSize, mimeType, duration = 0, thumbFile = null }) => {
    const fd = new FormData()
    fd.append('fileName', fileName)
    fd.append('fileSize', String(fileSize))
    fd.append('mimeType', mimeType)
    fd.append('duration', String(duration))
    if (thumbFile) {
      fd.append('thumb', thumbFile, 'thumb.jpg')
    }
    return request.post('/admin/media/media/initVideoDirectUpload', fd, {
      timeout: 0,
    })
  },

  uploadVideoToS3: async (token, file, onProgress) => {
    try {
      await putFileMultipart(token, file, onProgress)
      await request.post('/admin/media/media/completeVideoDirectUpload', {
        objectKey: token.objectKey,
        uploadId: token.uploadId,
      }, { timeout: 0 })
      onProgress?.(100, file.size, file.size)
    }
    catch (err) {
      if (token?.uploadId && token.objectKey) {
        request.post('/admin/media/media/abortVideoDirectUpload', {
          objectKey: token.objectKey,
          uploadId: token.uploadId,
        }, { timeout: 0, needTip: false }).catch(() => {})
      }
      throw err
    }
  },

  saveUploaded: payload => request.post('/admin/media/media/saveUploaded', payload, { timeout: 0 }),
  mediaCateList: type => request.post('/admin/media/mediaCate/list', { type }),
  mediaList: data => request.post('/admin/media/media/list', data),
  videoTus: file => request.post('/admin/media/media/videoUts', { uploadLength: file.size.toString(), filename: file.name }),
  videoTusSuccess: (videoId, cateId) => request.post('/admin/media/media/videoTusSuccess', { videoId, cateId }),
  cert,
  course,
  chapter,
  paper,
  path,
  user,
  message,
  feedback,
  ops,
  common,
  decoration,
  config,
  home,
  area,
  verse,
}
