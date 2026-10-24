package course

import (
	"github.com/gin-gonic/gin"

	admincourse "qixue-api/internal/service/admincourse"
)

type Course struct{ svc *admincourse.Service }

func NewCourse(svc *admincourse.Service) *Course {
	return &Course{svc: svc}
}

func (h *Course) GetCourseList(c *gin.Context) { h.svc.GetCourseList(c) }

func (h *Course) GetCourseDetail(c *gin.Context) { h.svc.GetCourseDetail(c) }

func (h *Course) CourseSave(c *gin.Context) { h.svc.CourseSave(c) }

func (h *Course) CourseDelete(c *gin.Context) { h.svc.CourseDelete(c) }

func (h *Course) GetCourseAttachList(c *gin.Context) { h.svc.GetCourseAttachList(c) }

func (h *Course) CourseAttachSave(c *gin.Context) { h.svc.CourseAttachSave(c) }

func (h *Course) CourseAttachDelete(c *gin.Context) { h.svc.CourseAttachDelete(c) }

func (h *Course) GetCourseUserList(c *gin.Context) { h.svc.GetCourseUserList(c) }

func (h *Course) GetCourseUserSearch(c *gin.Context) { h.svc.GetCourseUserSearch(c) }

func (h *Course) CourseUserListSave(c *gin.Context) { h.svc.CourseUserListSave(c) }

func (h *Course) CourseUserDelete(c *gin.Context) { h.svc.CourseUserDelete(c) }

func (h *Course) GetCourseUserApplyList(c *gin.Context) { h.svc.GetCourseUserApplyList(c) }

func (h *Course) CourseUserApplyStatusChange(c *gin.Context) { h.svc.CourseUserApplyStatusChange(c) }

func (h *Course) CourseUserApplyDelete(c *gin.Context) { h.svc.CourseUserApplyDelete(c) }
