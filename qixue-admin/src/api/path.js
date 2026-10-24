import { request } from '@/utils'

export default {
  getPathCategories: () => request.get('/admin/course/path/getPathCategories'),
  pathList: data => request.post('/admin/course/path/getPathList', data),
  getPathDetail: id => request.get(`/admin/course/path/getPathDetail?id=${id}`),
  pathDelete: id => request.get(`/admin/course/path/pathDelete?id=${id}`),
  pathSave: data => request.post('/admin/course/path/pathSave', data),

  // 关联课程
  pathCourse: data => request.post(`/admin/course/path/getPathCourseList`, data),
  pathCourseListSave: data => request.post(`/admin/course/path/pathCourseListSave`, data),
  pathCourseDelete: id => request.get(`/admin/course/path/pathCourseDelete?id=${id}`),
  pathCourseSearch: data => request.post(`/admin/course/path/getPathCourseSearch`, data),

  // 学员
  pathUser: data => request.post(`/admin/course/path/getUserPathList`, data),
  userPathDelete: id => request.get(`/admin/course/path/userPathDelete?id=${id}`), // 移除学员
  userPathSearch: data => request.post(`/admin/course/path/getUserPathSearch`, data),
  userPathListSave: data => request.post(`/admin/course/path/userPathListSave`, data), // 学员批量加入
  userPathApplyList: data => request.post(`/admin/course/path/getUserPathApplyList`, data), // 获取申请列表
  userPathApplyStatusChange: data => request.post(`/admin/course/path/userPathApplyStatusChange`, data), // 修改申请状态
  userPathApplyDelete: id => request.get(`/admin/course/path/userPathApplyDelete?id=${id}`), // 删除申请

}
