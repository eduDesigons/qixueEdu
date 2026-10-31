package login

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/middleware"
	"qixue-api/internal/http/response"
	"qixue-api/internal/platform/cache"
	"qixue-api/internal/admin/service/login"
)

type Handler struct {
	svc *adminlogin.Service
	rdb *redis.Client
}

func New(svc *adminlogin.Service, rdb *redis.Client) *Handler {
	return &Handler{svc: svc, rdb: rdb}
}

func (h *Handler) Captcha(c *gin.Context) {
	view, err := h.svc.Captcha(c.Request.Context())
	if err != nil {
		response.Fail(c, 1, "验证码加载失败")
		return
	}
	response.OK(c, "", view)
}

type checkBody struct {
	Token string   `json:"token"`
	X     bind.Int `json:"x"`
	Y     bind.Int `json:"y"`
}

func (h *Handler) CheckCaptcha(c *gin.Context) {
	var body checkBody
	if err := c.ShouldBindJSON(&body); err != nil || body.Token == "" {
		response.Fail(c, 10003, "参数错误")
		return
	}
	pass, err := h.svc.CheckCaptcha(c.Request.Context(), body.Token, body.X.Int(), body.Y.Int())
	if errors.Is(err, adminlogin.ErrBadCaptcha) {
		response.Fail(c, 10003, "验证失败")
		return
	}
	if err != nil {
		response.Fail(c, 500, "服务器发生异常")
		return
	}
	response.OK(c, "", gin.H{"captchaVerification": pass})
}

type loginBody struct {
	Username            string `json:"username"`
	Password            string `json:"password"`
	CaptchaVerification string `json:"captchaVerification"`
}

func (h *Handler) Login(c *gin.Context) {
	var body loginBody
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, 0, "请输入用户名和密码")
		return
	}
	username := strings.TrimSpace(body.Username)
	if username == "" || body.Password == "" {
		response.Fail(c, 0, "请输入用户名和密码")
		return
	}
	if len(username) > 64 || len(body.Password) > 128 {
		response.Fail(c, 0, "账号密码错误")
		return
	}
	if body.CaptchaVerification == "" {
		response.Fail(c, 10003, "请完成滑动验证")
		return
	}
	token, err := h.svc.Login(c.Request.Context(), username, body.Password, body.CaptchaVerification)
	switch {
	case errors.Is(err, adminlogin.ErrBadCaptcha):
		response.Fail(c, 10003, "验证码已失效，请重新验证")
	case errors.Is(err, adminlogin.ErrBlocked):
		response.Fail(c, 0, "错误次数过多，请稍后再试")
	case errors.Is(err, adminlogin.ErrPassword):
		response.Fail(c, 1, "账号密码错误")
	case errors.Is(err, adminlogin.ErrDisabled):
		response.Fail(c, 1, "当前管理员已被锁定无法登录")
	case err != nil:
		response.Fail(c, 500, "服务器发生异常")
	default:
		response.OK(c, "", gin.H{"accessToken": token})
	}
}

func (h *Handler) UserInfo(c *gin.Context) {
	profile, err := h.svc.UserInfo(c.Request.Context(), middleware.UID(c))
	if errors.Is(err, adminlogin.ErrNotFound) {
		response.Fail(c, 1, "用户不存在")
		return
	}
	if err != nil {
		response.Fail(c, 500, "服务器发生异常")
		return
	}
	response.OK(c, "获取个人信息", profile)
}

func (h *Handler) Permissions(c *gin.Context) {
	nodes, err := h.svc.Permissions(c.Request.Context(), middleware.UID(c))
	if err != nil {
		response.Fail(c, 500, "服务器发生异常")
		return
	}
	response.OK(c, "", nodes)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	patch, err := decodePatch(c)
	if err != nil {
		response.Fail(c, 1, "参数错误")
		return
	}
	err = h.svc.UpdateProfile(c.Request.Context(), middleware.UID(c), patch)
	switch {
	case errors.Is(err, adminlogin.ErrEmpty):
		response.Fail(c, 1, "没有可更新的内容")
	case errors.Is(err, adminlogin.ErrNotFound):
		response.Fail(c, 1, "用户不存在")
	case errors.Is(err, adminlogin.ErrNickName):
		response.Fail(c, 1, "昵称最多 10 个字")
	case err != nil:
		response.Fail(c, 500, "服务器发生异常")
	default:
		response.OK(c, "资料修改成功", []any{})
	}
}

type passwordBody struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var body passwordBody
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, 1, "请输入原密码和新密码")
		return
	}
	if body.OldPassword == "" || body.NewPassword == "" {
		response.Fail(c, 1, "请输入原密码和新密码")
		return
	}
	if len(body.NewPassword) < 6 {
		response.Fail(c, 1, "新密码至少 6 位")
		return
	}
	err := h.svc.ChangePassword(c.Request.Context(), middleware.UID(c), body.OldPassword, body.NewPassword)
	switch {
	case errors.Is(err, adminlogin.ErrNotFound):
		response.Fail(c, 1, "用户不存在")
	case errors.Is(err, adminlogin.ErrPassword):
		response.Fail(c, 1, "原密码错误")
	case err != nil:
		response.Fail(c, 500, "服务器发生异常")
	default:
		response.OK(c, "密码修改成功", []any{})
	}
}

func (h *Handler) Logout(c *gin.Context) {
	if err := middleware.BlockToken(c, h.rdb, middleware.Token(c), middleware.TokenExp(c)); err != nil {
		response.Fail(c, 500, "服务器发生异常")
		return
	}
	cache.Drop(c.Request.Context(), h.rdb, "admin_"+strconv.Itoa(middleware.UID(c)))
	response.OK(c, "已退出登录", []any{})
}

func decodePatch(c *gin.Context) (adminlogin.ProfilePatch, error) {
	var raw map[string]json.RawMessage
	if err := c.ShouldBindJSON(&raw); err != nil {
		return adminlogin.ProfilePatch{}, err
	}
	var patch adminlogin.ProfilePatch
	if value, ok := raw["nickName"]; ok {
		text, err := decodeString(value)
		if err != nil {
			return adminlogin.ProfilePatch{}, err
		}
		patch.NickName = &text
		patch.HasNick = true
	}
	if value, ok := raw["gender"]; ok {
		var n *int
		if string(value) != "null" {
			var num int
			if err := json.Unmarshal(value, &num); err != nil {
				return adminlogin.ProfilePatch{}, err
			}
			n = &num
		}
		patch.Gender = n
		patch.HasGender = true
	}
	for _, item := range []struct {
		key  string
		set  func(*string)
		mark *bool
	}{
		{"address", func(s *string) { patch.Address = s }, &patch.HasAddress},
		{"email", func(s *string) { patch.Email = s }, &patch.HasEmail},
		{"avatar", func(s *string) { patch.Avatar = s }, &patch.HasAvatar},
	} {
		value, ok := raw[item.key]
		if !ok {
			continue
		}
		if string(value) == "null" {
			item.set(nil)
			*item.mark = true
			continue
		}
		text, err := decodeString(value)
		if err != nil {
			return adminlogin.ProfilePatch{}, err
		}
		item.set(&text)
		*item.mark = true
	}
	return patch, nil
}

func decodeString(raw json.RawMessage) (string, error) {
	var text string
	if err := json.Unmarshal(raw, &text); err != nil {
		return "", err
	}
	return text, nil
}
