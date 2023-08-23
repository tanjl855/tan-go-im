package chatapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "index is here !!",
	})
}
