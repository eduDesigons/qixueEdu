package member

import (
	"github.com/gin-gonic/gin"

	adminmember "qixue-api/internal/service/adminmember"
)

type User struct{ svc *adminmember.Service }

func NewUser(svc *adminmember.Service) *User {
	return &User{svc: svc}
}

func (h *User) GetUserList(c *gin.Context) { h.svc.GetUserList(c) }

func (h *User) GetUserDetail(c *gin.Context) { h.svc.GetUserDetail(c) }

func (h *User) CreateUser(c *gin.Context) { h.svc.CreateUser(c) }

func (h *User) UpdateUser(c *gin.Context) { h.svc.UpdateUser(c) }

func (h *User) SetLock(c *gin.Context) { h.svc.SetLock(c) }

func (h *User) DeleteUser(c *gin.Context) { h.svc.DeleteUser(c) }

func (h *User) GetUserCourseList(c *gin.Context) { h.svc.GetUserCourseList(c) }

func (h *User) GetUserPathList(c *gin.Context) { h.svc.GetUserPathList(c) }

func (h *User) GetUserCourseChapterList(c *gin.Context) { h.svc.GetUserCourseChapterList(c) }

func (h *User) GetUserCreditRecordList(c *gin.Context) { h.svc.GetUserCreditRecordList(c) }

func (h *User) GetUserCoinsRecordList(c *gin.Context) { h.svc.GetUserCoinsRecordList(c) }

func (h *User) GetUserPaperList(c *gin.Context) { h.svc.GetUserPaperList(c) }

func (h *User) GetUserCertList(c *gin.Context) { h.svc.GetUserCertList(c) }
