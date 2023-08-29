package route

import (
	"github.com/gin-gonic/gin"
	"github.com/tanjl855/tan_go_im/servers/pool_server/internal/http/middle"
	"github.com/tanjl855/tan_go_im/servers/pool_server/internal/http/route/ws"
)

func InitRouters() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	//跨域和捕获panic
	r.Use(middle.Cors())
	r.Use(middle.LogMiddle())

	//sayHai
	r.HEAD("/", sayHi)
	r.GET("/", sayHi)

	g := r.Group("")

	//ws服务
	ws.InitRouterWithOutAuth(g)

	return r
}

func sayHi(ctx *gin.Context) {
	_, err := ctx.Writer.Write([]byte("hello tjl0.0"))
	if err != nil {
		return
	}
}
