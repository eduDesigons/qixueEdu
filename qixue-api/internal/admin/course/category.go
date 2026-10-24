package course

import (
	"github.com/gin-gonic/gin"

	admincategory "qixue-api/internal/service/admincategory"
)

type Category struct{ svc *admincategory.Service }

func NewCategory(svc *admincategory.Service) *Category {
	return &Category{svc: svc}
}

func (h *Category) List(c *gin.Context) { h.svc.List(c) }

func (h *Category) Save(c *gin.Context) { h.svc.Save(c) }

func (h *Category) Delete(c *gin.Context) { h.svc.Delete(c) }

func (h *Category) Options(c *gin.Context) { h.svc.Options(c) }
