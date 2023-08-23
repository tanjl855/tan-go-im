package app

import (
	"gin_chat/chatapi"
	"gin_chat/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router handlers
func Router() *gin.Engine {
	r := gin.Default()

	// swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/index", chatapi.GetIndex)
	r.GET("/user/getUserList", chatapi.GetUserList)
	return r
}
