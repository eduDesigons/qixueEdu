import { request } from '@/utils'

export default {
  getCourseCategories: () => request.get('/admin/course/course/getCourseCategories'),
  categoryList: type => request.get('/admin/course/category/list', { params: { type } }),
  categorySave: data => request.post('/admin/course/category/save', data),
  categoryDelete: data => request.post('/admin/course/category/delete', data),
  courseList: data => request.post('/admin/course/course/getCourseList', data),
  getCourseDetail: id => request.get(`/admin/course/course/getCourseDetail?id=${id}`),
  courseDelete: id => request.get(`/admin/course/course/courseDelete?id=${id}`),
  courseSave: data => request.post('/admin/course/course/courseSave', data),

  courseAttach: ({ courseId }) => request.get(`/admin/course/course/getCourseAttachList?courseId=${courseId}`),
  courseAttachDelete: id => request.get(`/admin/course/course/courseAttachDelete?id=${id}`),
  courseAttachSave: data => request.post('/admin/course/course/courseAttachSave', data),

  courseUser: data => request.post(`/admin/course/course/getCourseUserList`, data),
  courseUserDelete: id => request.get(`/admin/course/course/courseUserDelete?id=${id}`),
  courseUserSearch: data => request.post(`/admin/course/course/getCourseUserSearch`, data),
  courseUserListSave: data => request.post(`/admin/course/course/courseUserListSave`, data),
  courseUserApplyList: data => request.post(`/admin/course/course/getCourseUserApplyList`, data),
  courseUserApplyStatusChange: data => request.post(`/admin/course/course/courseUserApplyStatusChange`, data),
  courseUserApplyDelete: id => request.get(`/admin/course/course/courseUserApplyDelete?id=${id}`),

  courseComment: data => request.post(`/admin/course/comment/getCourseCommentList`, data),
  courseCommentReplies: data => request.get('/admin/course/comment/getCourseCommentReplies', { params: data }),
  courseCommentDelete: id => request.get(`/admin/course/comment/courseCommentDelete?id=${id}`),
  courseCommentBatchDelete: ids => request.post('/admin/course/comment/courseCommentBatchDelete', { ids }),
  courseCommentReply: data => request.post('/admin/course/comment/addCourseCommentReply', data),

  courseChapterComment: data => request.post(`/admin/course/comment/getCourseChapterCommentList`, data),
  courseChapterCommentReplies: data => request.get('/admin/course/comment/getChapterCommentReplies', { params: data }),
  courseChapterCommentDelete: id => request.get(`/admin/course/comment/courseChapterCommentDelete?id=${id}`),
  courseChapterCommentBatchDelete: ids => request.post('/admin/course/comment/courseChapterCommentBatchDelete', { ids }),
  courseChapterCommentReply: data => request.post('/admin/course/comment/addChapterCommentReply', data),
}
