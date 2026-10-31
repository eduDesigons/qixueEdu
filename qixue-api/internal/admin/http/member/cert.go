package member

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	"qixue-api/internal/admin/service/cert"
)

type Cert struct{ svc *admincert.Service }

func NewCert(svc *admincert.Service) *Cert {
	return &Cert{svc: svc}
}

func (h *Cert) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	var blocked *admincert.DeleteBlockedError
	switch {
	case errors.As(err, &blocked),
		errors.Is(err, admincert.ErrMissingID),
		errors.Is(err, admincert.ErrNotFound),
		errors.Is(err, admincert.ErrEmptyName),
		errors.Is(err, admincert.ErrEmptyTemplate),
		errors.Is(err, admincert.ErrMissingTarget),
		errors.Is(err, admincert.ErrMissingScope),
		errors.Is(err, admincert.ErrCourseNoIssue),
		errors.Is(err, admincert.ErrPathNoIssue),
		errors.Is(err, admincert.ErrIssuedNotFound):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Cert) List(c *gin.Context) {
	var body struct {
		Name   string   `json:"name"`
		IsShow any      `json:"is_show"`
		Page   bind.Int `json:"page"`
		Size   bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	q := admincert.ListQuery{Name: body.Name, Page: body.Page.Int(), Size: body.Size.Int()}
	if bind.Present(body.IsShow) {
		n := bind.AsInt(body.IsShow)
		q.IsShow = &n
	}
	data, err := h.svc.List(c.Request.Context(), q)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *Cert) Detail(c *gin.Context) {
	view, err := h.svc.Detail(c.Request.Context(), bind.Query(c, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", view)
}

func (h *Cert) Options(c *gin.Context) {
	list, err := h.svc.Options(c.Request.Context())
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", list)
}

func (h *Cert) SetShow(c *gin.Context) {
	var body struct {
		ID     bind.Int `json:"id"`
		IsShow any      `json:"is_show"`
	}
	_ = c.ShouldBindJSON(&body)
	view, msg, err := h.svc.SetShow(c.Request.Context(), body.ID.Int(), bind.AsInt(body.IsShow))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, msg, view)
}

func (h *Cert) Save(c *gin.Context) {
	var body struct {
		ID            bind.Int        `json:"id"`
		Name          string          `json:"name"`
		TemplateImage string          `json:"template_image"`
		PreviewImage  string          `json:"preview_image"`
		Params        json.RawMessage `json:"params"`
		IsShow        any             `json:"is_show"`
	}
	_ = c.ShouldBindJSON(&body)
	in := admincert.SaveInput{
		ID: body.ID.Int(), Name: body.Name, TemplateImage: body.TemplateImage,
		PreviewImage: body.PreviewImage, Params: body.Params,
	}
	if body.IsShow != nil {
		n := bind.AsInt(body.IsShow)
		in.IsShow = &n
	}
	view, msg, err := h.svc.Save(c.Request.Context(), in)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, msg, view)
}

func (h *Cert) Delete(c *gin.Context) {
	var body struct {
		ID bind.Int `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	if h.writeErr(c, h.svc.Delete(c.Request.Context(), body.ID.Int())) {
		return
	}
	response.OK(c, "删除成功", nil)
}

func (h *Cert) Issue(c *gin.Context) {
	var body struct {
		UserID   bind.Int `json:"user_id"`
		PathID   bind.Int `json:"path_id"`
		CourseID bind.Int `json:"course_id"`
		Force    bool     `json:"force"`
	}
	_ = c.ShouldBindJSON(&body)
	view, msg, err := h.svc.Issue(c.Request.Context(), admincert.IssueInput{
		UserID: body.UserID.Int(), PathID: body.PathID.Int(), CourseID: body.CourseID.Int(), Force: body.Force,
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, msg, view)
}

func (h *Cert) Backfill(c *gin.Context) {
	var body struct {
		CourseID bind.Int `json:"course_id"`
		PathID   bind.Int `json:"path_id"`
	}
	_ = c.ShouldBindJSON(&body)
	count, err := h.svc.Backfill(c.Request.Context(), body.CourseID.Int(), body.PathID.Int())
	if h.writeErr(c, err) {
		return
	}
	if count <= 0 {
		response.OK(c, "暂无需要补发的学员", gin.H{"count": 0})
		return
	}
	response.OK(c, fmt.Sprintf("已为 %d 名学员补发证书", count), gin.H{"count": count})
}

func (h *Cert) IssuedList(c *gin.Context) {
	var body struct {
		Name     string   `json:"name"`
		CourseID bind.Int `json:"course_id"`
		PathID   bind.Int `json:"path_id"`
		Page     bind.Int `json:"page"`
		Size     bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	data, err := h.svc.IssuedList(c.Request.Context(), middleware.UID(c), admincert.IssuedListQuery{
		Name: body.Name, CourseID: body.CourseID.Int(), PathID: body.PathID.Int(),
		Page: body.Page.Int(), Size: body.Size.Int(),
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *Cert) Revoke(c *gin.Context) {
	var body struct {
		ID bind.Int `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	view, msg, err := h.svc.Revoke(c.Request.Context(), body.ID.Int())
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, msg, view)
}
