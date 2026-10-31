package member

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	"qixue-api/internal/admin/service/message"
)

type Message struct{ svc *adminmessage.Service }

func NewMessage(svc *adminmessage.Service) *Message {
	return &Message{svc: svc}
}

func (h *Message) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, adminmessage.ErrEmptyTitle),
		errors.Is(err, adminmessage.ErrTitleTooLong),
		errors.Is(err, adminmessage.ErrEmptyContent),
		errors.Is(err, adminmessage.ErrContentLong),
		errors.Is(err, adminmessage.ErrNoUsers),
		errors.Is(err, adminmessage.ErrUsersMissing),
		errors.Is(err, adminmessage.ErrNotFound):
		response.Fail(c, 1, err.Error())
		return true
	default:
		return response.WriteErr(c, err)
	}
}

func (h *Message) Send(c *gin.Context) {
	var body struct {
		Title       string `json:"title"`
		Content     string `json:"content"`
		Type        string `json:"type"`
		RelatedType any    `json:"related_type"`
		RelatedID   any    `json:"related_id"`
		UserIDs     any    `json:"user_ids"`
		UserID      any    `json:"user_id"`
	}
	_ = c.ShouldBindJSON(&body)

	userIDs := asInts(body.UserIDs)
	if len(userIDs) == 0 {
		if uid := bind.AsInt(body.UserID); uid > 0 {
			userIDs = []int{uid}
		}
	}
	count, err := h.svc.Send(c.Request.Context(), middleware.UID(c), adminmessage.SendInput{
		Title: body.Title, Content: body.Content, Type: body.Type,
		RelatedType: asOptionalString(body.RelatedType), RelatedID: asOptionalRelatedID(body.RelatedID),
		UserIDs: userIDs,
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "发送成功", gin.H{"count": count})
}

func (h *Message) List(c *gin.Context) {
	var body struct {
		Title    string   `json:"title"`
		NickName string   `json:"nick_name"`
		Type     string   `json:"type"`
		IsRead   any      `json:"is_read"`
		Page     bind.Int `json:"page"`
		Size     bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	q := adminmessage.ListQuery{
		Title: body.Title, NickName: body.NickName, Type: body.Type,
		Page: body.Page.Int(), Size: body.Size.Int(),
	}
	if bind.Present(body.IsRead) {
		n := bind.AsInt(body.IsRead)
		q.IsRead = &n
	}
	data, err := h.svc.List(c.Request.Context(), middleware.UID(c), q)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *Message) Delete(c *gin.Context) {
	var body struct {
		ID bind.Int `json:"id"`
	}
	_ = c.ShouldBindJSON(&body)
	if h.writeErr(c, h.svc.Delete(c.Request.Context(), body.ID.Int())) {
		return
	}
	response.OK(c, "删除成功", nil)
}

func asOptionalString(v any) *string {
	if v == nil {
		return nil
	}
	switch n := v.(type) {
	case string:
		s := strings.TrimSpace(n)
		if s == "" {
			return nil
		}
		return &s
	default:
		return nil
	}
}

func asOptionalRelatedID(v any) *int {
	if v == nil {
		return nil
	}
	id := bind.AsInt(v)
	if id <= 0 {
		return nil
	}
	return &id
}

func asInts(v any) []int {
	switch n := v.(type) {
	case []any:
		out := make([]int, 0, len(n))
		for _, item := range n {
			if id := bind.AsInt(item); id > 0 {
				out = append(out, id)
			}
		}
		return out
	case []int:
		return n
	case []float64:
		out := make([]int, 0, len(n))
		for _, item := range n {
			if id := int(item); id > 0 {
				out = append(out, id)
			}
		}
		return out
	default:
		return nil
	}
}
