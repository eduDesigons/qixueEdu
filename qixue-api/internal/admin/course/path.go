package course

import (
	"github.com/gin-gonic/gin"

	adminpath "qixue-api/internal/service/adminpath"
)

type Path struct{ svc *adminpath.Service }

func NewPath(svc *adminpath.Service) *Path {
	return &Path{svc: svc}
}

func (h *Path) GetPathCourseSearch(c *gin.Context) { h.svc.GetPathCourseSearch(c) }

func (h *Path) GetPathCourseList(c *gin.Context) { h.svc.GetPathCourseList(c) }

func (h *Path) PathCourseListSave(c *gin.Context) { h.svc.PathCourseListSave(c) }

func (h *Path) PathCourseDelete(c *gin.Context) { h.svc.PathCourseDelete(c) }

func (h *Path) GetPathCategories(c *gin.Context) { h.svc.GetPathCategories(c) }

func (h *Path) GetPathList(c *gin.Context) { h.svc.GetPathList(c) }

func (h *Path) GetPathDetail(c *gin.Context) { h.svc.GetPathDetail(c) }

func (h *Path) PathSave(c *gin.Context) { h.svc.PathSave(c) }

func (h *Path) PathDelete(c *gin.Context) { h.svc.PathDelete(c) }

func (h *Path) GetUserPathList(c *gin.Context) { h.svc.GetUserPathList(c) }

func (h *Path) GetUserPathSearch(c *gin.Context) { h.svc.GetUserPathSearch(c) }

func (h *Path) UserPathListSave(c *gin.Context) { h.svc.UserPathListSave(c) }

func (h *Path) UserPathDelete(c *gin.Context) { h.svc.UserPathDelete(c) }

func (h *Path) GetUserPathApplyList(c *gin.Context) { h.svc.GetUserPathApplyList(c) }

func (h *Path) UserPathApplyStatusChange(c *gin.Context) { h.svc.UserPathApplyStatusChange(c) }

func (h *Path) UserPathApplyDelete(c *gin.Context) { h.svc.UserPathApplyDelete(c) }
