package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/tanjl855/tan_go_im/servers/pool_server/internal/controller/http/ws"
)

func InitRouterWithOutAuth(g *gin.RouterGroup) {
	//初始化ws服务
	ws.Init()
	userG := g.Group("/ws")
	userG.GET("", ws.WS.WsHandler)
}
