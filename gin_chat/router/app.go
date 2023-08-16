package app

import (
	"gin_chat/chatapi"
	"github.com/gin-gonic/gin"
)

// Router handlers
func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/index", chatapi.GetIndex)
	return r
}
