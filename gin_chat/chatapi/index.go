package chatapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetIndex(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "index is here !!",
	})
}
