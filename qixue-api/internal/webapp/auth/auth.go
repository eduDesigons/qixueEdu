package auth

import (
	"errors"

	"github.com/gin-gonic/gin"

	"qixue-api/internal/http/bind"
	"qixue-api/internal/http/response"
	"qixue-api/internal/service/webauth"
)

type Handler struct {
	svc *webauth.Service
}

func New(svc *webauth.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Captcha(c *gin.Context) {
	view, err := h.svc.Captcha(c.Request.Context())
	if err != nil {
		response.Fail(c, 1, "验证码加载失败")
		return
	}
	response.OK(c, "ok", view)
}

func (h *Handler) CheckCaptcha(c *gin.Context) {
	var body struct {
		Token string   `json:"token"`
		X     bind.Int `json:"x"`
		Y     bind.Int `json:"y"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Token == "" {
		response.Fail(c, 10003, "参数错误")
		return
	}
	pass, err := h.svc.CheckCaptcha(c.Request.Context(), body.Token, body.X.Int(), body.Y.Int())
	switch {
	case errors.Is(err, webauth.ErrCaptchaParam):
		response.Fail(c, 10003, err.Error())
	case errors.Is(err, webauth.ErrCaptchaFail):
		response.Fail(c, 10003, err.Error())
	case errors.Is(err, webauth.ErrServer):
		response.Fail(c, 500, err.Error())
	case err != nil:
		response.Fail(c, 500, "服务器发生异常")
	default:
		response.OK(c, "ok", gin.H{"captchaVerification": pass})
	}
}

func (h *Handler) Login(c *gin.Context) {
	var body struct {
		Username            string `json:"username"`
		Password            string `json:"password"`
		CaptchaVerification string `json:"captchaVerification"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, 0, "请填写用户名")
		return
	}
	result, err := h.svc.Login(c.Request.Context(), c.Request.Host, body.Username, body.Password, body.CaptchaVerification)
	switch {
	case errors.Is(err, webauth.ErrUsernameRequired),
		errors.Is(err, webauth.ErrPasswordRequired),
		errors.Is(err, webauth.ErrBlocked),
		errors.Is(err, webauth.ErrInactive),
		errors.Is(err, webauth.ErrLocked):
		response.Fail(c, 0, err.Error())
	case errors.Is(err, webauth.ErrNeedCaptcha),
		errors.Is(err, webauth.ErrCaptchaExpired):
		response.Fail(c, 10003, err.Error())
	case errors.Is(err, webauth.ErrBadCredential):
		need := result != nil && result.NeedCaptcha
		response.FailData(c, 0, err.Error(), gin.H{"needCaptcha": need})
	case errors.Is(err, webauth.ErrServer):
		response.Fail(c, 500, err.Error())
	case err != nil:
		response.Fail(c, 500, "服务器发生异常")
	default:
		response.OK(c, "登录成功", gin.H{"token": result.Token})
	}
}
