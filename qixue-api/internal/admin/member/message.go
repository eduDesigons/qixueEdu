package member

import (
	"github.com/gin-gonic/gin"

	adminmessage "qixue-api/internal/service/adminmessage"
)

type Message struct{ svc *adminmessage.Service }

func NewMessage(svc *adminmessage.Service) *Message {
	return &Message{svc: svc}
}

func (h *Message) Send(c *gin.Context) { h.svc.Send(c) }

func (h *Message) List(c *gin.Context) { h.svc.List(c) }

func (h *Message) Delete(c *gin.Context) { h.svc.Delete(c) }
