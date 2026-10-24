import { request } from '@/utils'

export default {
  getUserList: data => request.post('/admin/member/user/getUserList', data),
  getUserDetail: id => request.get(`/admin/member/user/getUserDetail?id=${id}`),
  getUserPathList: data => request.get('/admin/member/user/getUserPathList', { params: data }),
  getUserCourseList: data => request.get('/admin/member/user/getUserCourseList', { params: data }),
  getUserCourseChapterList: data => request.get('/admin/member/user/getUserCourseChapterList', { params: data }),
  getUserCreditRecordList: data => request.post('/admin/member/user/getUserCreditRecordList', data),
  getUserCoinsRecordList: data => request.post('/admin/member/user/getUserCoinsRecordList', data),
  getUserPaperList: data => request.post('/admin/member/user/getUserPaperList', data),
  // 新增用户
  createUser: data => request.post('/admin/member/user/createUser', data),
  // 更新用户
  updateUser: data => request.post('/admin/member/user/updateUser', data),
  setLock: data => request.post('/admin/member/user/setLock', data),
  userDelete: id => request.get(`/admin/member/user/deleteUser?id=${id}`),
  getUserCertList: data => request.post('/admin/member/user/getUserCertList', data),
}
