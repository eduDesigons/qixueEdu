import { request } from '@/utils'

export default {
  paperDetail: id => request.get(`/admin/course/paper/getPaperDetail?id=${id}`),
  paperLogs: data => request.post('/admin/course/paper/getPaperLogs', data),
  paperDelete: id => request.get(`/admin/course/paper/paperDelete?id=${id}`),
  paperSave: data => request.post('/admin/course/paper/paperSave', data),
  examList: data => request.post('/admin/course/paper/getExamList', data),
  assessmentList: data => request.post('/admin/course/paper/getAssessmentList', data),
  getUserPaperDetail: id => request.get(`/admin/course/paper/getUserPaperDetail?id=${id}`),
  paperGrade: data => request.post('/admin/course/paper/paperGrade', data),
}
