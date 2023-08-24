package im_rsp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 接口返回结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// HttpResponse setting gin.JSON
func HttpResponse(ctx *gin.Context, httpCode, code int, data interface{}) {
	switch data.(type) {
	case error:
		ctx.JSON(httpCode, Response{
			Code: code,
			Msg:  getMsg(code),
			Data: fmt.Sprint(data),
		})
	default:
		ctx.JSON(httpCode, Response{
			Code: code,
			Msg:  getMsg(code),
			Data: data,
		})
	}
}

func Success(ctx *gin.Context, data interface{}) {
	HttpResponse(ctx, http.StatusOK, 0, data)
	ctx.Abort()
}

func Failed(ctx *gin.Context, code int, err error) {
	HttpResponse(ctx, http.StatusOK, code, err)
	ctx.Abort()
}

func FailedWithData(ctx *gin.Context, code int, data interface{}) {
	HttpResponse(ctx, http.StatusOK, code, data)
	ctx.Abort()
}
