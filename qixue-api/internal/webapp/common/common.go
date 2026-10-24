package common

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/response"
	"qixue-api/internal/service/bunny"
)

type Handler struct {
	uploader *bunny.Uploader
}

func New() *Handler {
	return &Handler{uploader: bunny.NewUploader(bunny.LoadConfig())}
}

func (h *Handler) Upload(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil || fh == nil {
		response.Fail(c, 1, "文件不存在")
		return
	}
	maxEdge, _ := strconv.Atoi(c.PostForm("maxEdge"))
	result, err := h.uploader.HandleFormFile(fh, "", maxEdge)
	if err != nil {
		response.Fail(c, 1, err.Error())
		return
	}
	response.OK(c, "ok", result)
}
