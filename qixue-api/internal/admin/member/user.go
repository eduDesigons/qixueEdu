package member

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	adminmember "qixue-api/internal/service/admin/member"
)

type User struct{ svc *adminmember.Service }

func NewUser(svc *adminmember.Service) *User {
	return &User{svc: svc}
}

func (h *User) writeErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	var be *adminmember.BizError
	if errors.As(err, &be) {
		response.Fail(c, 1, be.Message)
		return true
	}
	return response.WriteErr(c, err)
}

func (h *User) GetUserList(c *gin.Context) {
	var body struct {
		NickName string   `json:"nick_name"`
		Account  string   `json:"account"`
		Email    string   `json:"email"`
		AreaID   bind.Int `json:"area_id"`
		IsActive any      `json:"is_active"`
		IsLock   any      `json:"is_lock"`
		Page     bind.Int `json:"page"`
		Size     bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	data, err := h.svc.GetUserList(c.Request.Context(), middleware.UID(c), adminmember.ListQuery{
		NickName: body.NickName, Account: body.Account, Email: body.Email,
		AreaID: body.AreaID.Int(), IsActive: body.IsActive, IsLock: body.IsLock,
		Page: body.Page.Int(), Size: body.Size.Int(),
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *User) GetUserDetail(c *gin.Context) {
	detail, err := h.svc.GetUserDetail(c.Request.Context(), bind.Query(c, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "", detail)
}

func (h *User) CreateUser(c *gin.Context) {
	var body struct {
		NickName   string   `json:"nick_name"`
		Account    string   `json:"account"`
		Mobile     string   `json:"mobile"`
		AreaID     bind.Int `json:"area_id"`
		CategoryID bind.Int `json:"category_id"`
		Credit     bind.Int `json:"credit"`
		Coins      bind.Int `json:"coins"`
		Password   string   `json:"password"`
		Email      string   `json:"email"`
	}
	_ = c.ShouldBindJSON(&body)
	err := h.svc.CreateUser(c.Request.Context(), adminmember.CreateInput{
		NickName: body.NickName, Account: body.Account, Mobile: body.Mobile,
		AreaID: body.AreaID.Int(), CategoryID: body.CategoryID.Int(),
		Credit: body.Credit.Int(), Coins: body.Coins.Int(),
		Password: body.Password, Email: body.Email,
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "新增成功", []any{})
}

func (h *User) UpdateUser(c *gin.Context) {
	var raw map[string]any
	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	err := h.svc.UpdateUser(c.Request.Context(), raw)
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "编辑成功", []any{})
}

func (h *User) SetLock(c *gin.Context) {
	var body struct {
		ID     bind.Int `json:"id"`
		IsLock bind.Int `json:"is_lock"`
	}
	_ = c.ShouldBindJSON(&body)
	msg, err := h.svc.SetLock(c.Request.Context(), adminmember.SetLockInput{
		ID: body.ID.Int(), IsLock: body.IsLock.Int(),
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, msg, []any{})
}

func (h *User) DeleteUser(c *gin.Context) {
	err := h.svc.DeleteUser(c.Request.Context(), bind.Query(c, "id"))
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "删除成功", []any{})
}

func (h *User) GetUserCourseList(c *gin.Context) {
	data, err := h.svc.GetUserCourseList(c.Request.Context(), adminmember.UserSubListQuery{
		UserID: bind.Query(c, "userId"), Title: c.Query("title"),
		Page: bind.Query(c, "page"), Size: bind.Query(c, "size"),
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *User) GetUserPathList(c *gin.Context) {
	data, err := h.svc.GetUserPathList(c.Request.Context(), adminmember.UserSubListQuery{
		UserID: bind.Query(c, "userId"), Title: c.Query("title"),
		Page: bind.Query(c, "page"), Size: bind.Query(c, "size"),
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *User) GetUserCourseChapterList(c *gin.Context) {
	data, err := h.svc.GetUserCourseChapterList(c.Request.Context(), adminmember.UserSubListQuery{
		UserID: bind.Query(c, "userId"), Title: c.Query("title"),
		Page: bind.Query(c, "page"), Size: bind.Query(c, "size"),
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *User) GetUserCreditRecordList(c *gin.Context) {
	var body struct {
		UserID bind.Int `json:"user_id"`
		Page   bind.Int `json:"page"`
		Size   bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	data, err := h.svc.GetUserCreditRecordList(c.Request.Context(), adminmember.RecordListQuery{
		UserID: body.UserID.Int(), Page: body.Page.Int(), Size: body.Size.Int(),
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *User) GetUserCoinsRecordList(c *gin.Context) {
	var body struct {
		UserID bind.Int `json:"user_id"`
		Page   bind.Int `json:"page"`
		Size   bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	data, err := h.svc.GetUserCoinsRecordList(c.Request.Context(), adminmember.RecordListQuery{
		UserID: body.UserID.Int(), Page: body.Page.Int(), Size: body.Size.Int(),
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *User) GetUserPaperList(c *gin.Context) {
	var body struct {
		UserID        bind.Int `json:"user_id"`
		Type          string   `json:"type"`
		IsPass        any      `json:"is_pass"`
		IsPeopleCheck any      `json:"is_people_check"`
		Page          bind.Int `json:"page"`
		Size          bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	data, err := h.svc.GetUserPaperList(c.Request.Context(), adminmember.PaperListQuery{
		UserID: body.UserID.Int(), Type: body.Type,
		IsPass: body.IsPass, IsPeopleCheck: body.IsPeopleCheck,
		Page: body.Page.Int(), Size: body.Size.Int(),
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}

func (h *User) GetUserCertList(c *gin.Context) {
	var body struct {
		UserID bind.Int `json:"user_id"`
		Page   bind.Int `json:"page"`
		Size   bind.Int `json:"size"`
	}
	_ = c.ShouldBindJSON(&body)
	data, err := h.svc.GetUserCertList(c.Request.Context(), adminmember.RecordListQuery{
		UserID: body.UserID.Int(), Page: body.Page.Int(), Size: body.Size.Int(),
	})
	if h.writeErr(c, err) {
		return
	}
	response.OK(c, "ok", data)
}
