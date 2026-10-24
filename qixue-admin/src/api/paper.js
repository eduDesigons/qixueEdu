import { request } from '@/utils'

export default {
  paperDetail: id => request.get(`/admin/course/paper/getPaperDetail?id=${id}`),
  paperLogs: data => request.post('/admin/course/paper/getPaperLogs', data),
  paperDelete: id => request.get(`/admin/course/paper/paperDelete?id=${id}`),
  paperSave: data => request.post('/admin/course/paper/paperSave', data),
  // 考试记录管理
  examList: data => request.post('/admin/course/paper/getExamList', data),
  // 测评记录管理
  assessmentList: data => request.post('/admin/course/paper/getAssessmentList', data),
  // 批改：获取用户答题详情
  getUserPaperDetail: id => request.get(`/admin/course/paper/getUserPaperDetail?id=${id}`),
  // 批改：提交批改
  paperGrade: data => request.post('/admin/course/paper/paperGrade', data),
}
