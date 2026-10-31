package media

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/response"
	"qixue-api/internal/admin/service/media"
)

type Media struct{ svc *adminmedia.Service }

func NewMedia(db *gorm.DB) *Media {
	return &Media{svc: adminmedia.New(db)}
}

func (h *Media) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	var ce *adminmedia.ClientError
	switch {
	case errors.As(err, &ce),
		errors.Is(err, adminmedia.ErrNoFile),
		errors.Is(err, adminmedia.ErrBadParams),
		errors.Is(err, adminmedia.ErrSaveParams),
		errors.Is(err, adminmedia.ErrNotFound),
		errors.Is(err, adminmedia.ErrNeedFileMeta),
		errors.Is(err, adminmedia.ErrVideoTooBig):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Media) List(c *gin.Context) {
	var body struct {
		Type   string   `json:"type"`
		CateID any      `json:"cateId"`
		Title  string   `json:"title"`
		Page   bind.Int `json:"page"`
		Size   bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	q := adminmedia.ListQuery{
		Type: body.Type, Title: body.Title, Page: body.Page.Int(), Size: body.Size.Int(),
	}
	if bind.Present(body.CateID) {
		n := bind.AsInt(body.CateID)
		q.CateID = &n
	}
	data, err := h.svc.List(c.Request.Context(), q)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *Media) Upload(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil || fh == nil {
		response.Fail(c, 1, "文件不存在或无效")
		return
	}
	out, err := h.svc.Upload(c.Request.Context(), fh, bind.Query(c, "cateId"), bind.Query(c, "maxEdge"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", out)
}

func (h *Media) Destroy(c *gin.Context) {
	var body struct {
		ID any `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.ID == nil {
		body.ID = c.Query("id")
	}
	if h.writeErr(c, h.svc.Destroy(c.Request.Context(), bind.AsInt64(body.ID))) {
		return
	}
	response.OK(c, "删除成功", nil)
}

func (h *Media) SaveUploaded(c *gin.Context) {
	var body struct {
		URL      string   `json:"url"`
		Title    string   `json:"title"`
		Type     string   `json:"type"`
		CateID   bind.Int `json:"cate_id"`
		Thumb    string   `json:"thumb"`
		Duration int      `json:"duration"`
		Size     int64    `json:"size"`
		Vid      string   `json:"vid"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, 1, "落库失败: 参数错误")
		return
	}
	out, err := h.svc.SaveUploaded(c.Request.Context(), adminmedia.SaveUploadedInput{
		URL: body.URL, Title: body.Title, Type: body.Type, CateID: body.CateID.Int(),
		Thumb: body.Thumb, Duration: body.Duration, Size: body.Size, Vid: body.Vid,
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", out)
}

func (h *Media) InitVideoDirectUpload(c *gin.Context) {
	fileSize, _ := strconv.ParseInt(c.PostForm("fileSize"), 10, 64)
	duration, _ := strconv.Atoi(c.PostForm("duration"))
	in := adminmedia.InitVideoInput{
		FileName: c.PostForm("fileName"),
		MimeType: c.PostForm("mimeType"),
		FileSize: fileSize,
		Duration: duration,
	}
	if fh, err := c.FormFile("thumb"); err == nil && fh != nil {
		in.Thumb = fh
	}
	out, err := h.svc.InitVideoDirectUpload(in)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", out)
}

func (h *Media) CompleteVideoDirectUpload(c *gin.Context) {
	var body struct {
		ObjectKey string `json:"objectKey"`
		UploadID  string `json:"uploadId"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.ObjectKey == "" {
		body.ObjectKey = c.PostForm("objectKey")
		body.UploadID = c.PostForm("uploadId")
	}
	if h.writeErr(c, h.svc.CompleteVideoDirectUpload(body.ObjectKey, body.UploadID)) {
		return
	}
	response.OK(c, "ok", nil)
}

func (h *Media) AbortVideoDirectUpload(c *gin.Context) {
	var body struct {
		ObjectKey string `json:"objectKey"`
		UploadID  string `json:"uploadId"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.ObjectKey == "" {
		body.ObjectKey = c.PostForm("objectKey")
		body.UploadID = c.PostForm("uploadId")
	}
	h.svc.AbortVideoDirectUpload(body.ObjectKey, body.UploadID)
	response.OK(c, "ok", nil)
}
