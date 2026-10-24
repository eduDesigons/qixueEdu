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

  // 课件
  courseAttach: ({ courseId }) => request.get(`/admin/course/course/getCourseAttachList?courseId=${courseId}`),
  courseAttachDelete: id => request.get(`/admin/course/course/courseAttachDelete?id=${id}`),
  courseAttachSave: data => request.post('/admin/course/course/courseAttachSave', data),

  // 学员
  courseUser: data => request.post(`/admin/course/course/getCourseUserList`, data),
  courseUserDelete: id => request.get(`/admin/course/course/courseUserDelete?id=${id}`), // 移除学员
  courseUserSearch: data => request.post(`/admin/course/course/getCourseUserSearch`, data),
  courseUserListSave: data => request.post(`/admin/course/course/courseUserListSave`, data), // 学员批量加入课程
  courseUserApplyList: data => request.post(`/admin/course/course/getCourseUserApplyList`, data), // 获取申请列表
  courseUserApplyStatusChange: data => request.post(`/admin/course/course/courseUserApplyStatusChange`, data), // 修改申请状态
  courseUserApplyDelete: id => request.get(`/admin/course/course/courseUserApplyDelete?id=${id}`), // 删除申请

  // 课程评论
  courseComment: data => request.post(`/admin/course/comment/getCourseCommentList`, data),
  courseCommentReplies: data => request.get('/admin/course/comment/getCourseCommentReplies', { params: data }),
  courseCommentDelete: id => request.get(`/admin/course/comment/courseCommentDelete?id=${id}`),
  courseCommentBatchDelete: ids => request.post('/admin/course/comment/courseCommentBatchDelete', { ids }),
  courseCommentReply: data => request.post('/admin/course/comment/addCourseCommentReply', data),

  // 课节评论
  courseChapterComment: data => request.post(`/admin/course/comment/getCourseChapterCommentList`, data),
  courseChapterCommentReplies: data => request.get('/admin/course/comment/getChapterCommentReplies', { params: data }),
  courseChapterCommentDelete: id => request.get(`/admin/course/comment/courseChapterCommentDelete?id=${id}`),
  courseChapterCommentBatchDelete: ids => request.post('/admin/course/comment/courseChapterCommentBatchDelete', { ids }),
  courseChapterCommentReply: data => request.post('/admin/course/comment/addChapterCommentReply', data),
}
