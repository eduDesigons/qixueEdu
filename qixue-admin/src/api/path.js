import { request } from '@/utils'

export default {
  getPathCategories: () => request.get('/admin/course/path/getPathCategories'),
  pathList: data => request.post('/admin/course/path/getPathList', data),
  getPathDetail: id => request.get(`/admin/course/path/getPathDetail?id=${id}`),
  pathDelete: id => request.get(`/admin/course/path/pathDelete?id=${id}`),
  pathSave: data => request.post('/admin/course/path/pathSave', data),

  pathCourse: data => request.post(`/admin/course/path/getPathCourseList`, data),
  pathCourseListSave: data => request.post(`/admin/course/path/pathCourseListSave`, data),
  pathCourseDelete: id => request.get(`/admin/course/path/pathCourseDelete?id=${id}`),
  pathCourseSearch: data => request.post(`/admin/course/path/getPathCourseSearch`, data),

  pathUser: data => request.post(`/admin/course/path/getUserPathList`, data),
  userPathDelete: id => request.get(`/admin/course/path/userPathDelete?id=${id}`),
  userPathSearch: data => request.post(`/admin/course/path/getUserPathSearch`, data),
  userPathListSave: data => request.post(`/admin/course/path/userPathListSave`, data),
  userPathApplyList: data => request.post(`/admin/course/path/getUserPathApplyList`, data),
  userPathApplyStatusChange: data => request.post(`/admin/course/path/userPathApplyStatusChange`, data),
  userPathApplyDelete: id => request.get(`/admin/course/path/userPathApplyDelete?id=${id}`),

}
