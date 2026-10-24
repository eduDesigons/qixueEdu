package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func OK(c *gin.Context, message string, data any) {
	write(c, 200, message, data)
}

func Fail(c *gin.Context, code int, message string) {
	write(c, code, message, nil)
}

func FailData(c *gin.Context, code int, message string, data any) {
	write(c, code, message, data)
}

func write(c *gin.Context, code int, message string, data any) {
	c.JSON(http.StatusOK, Body{Code: code, Message: message, Data: data})
}
