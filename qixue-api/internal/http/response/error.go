package response

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type Biz struct {
	Code    int
	Message string
}

func (e *Biz) Error() string { return e.Message }

func BizErr(message string) error {
	return &Biz{Code: 1, Message: message}
}

func BizCode(code int, message string) error {
	return &Biz{Code: code, Message: message}
}

func WriteErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	var biz *Biz
	if errors.As(err, &biz) {
		Fail(c, biz.Code, biz.Message)
		return true
	}
	Fail(c, 500, "服务器发生异常")
	return true
}
