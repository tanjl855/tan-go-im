package middle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/tanjl855/tan_go_im/pkg/im_log"
	"net/http"
	"time"
)

func LogMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		last := time.Now().Sub(start)
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		clientIP := c.GetHeader("X-Forwarded-For")
		if len(clientIP) == 0 {
			clientIP = c.ClientIP()
		}
		if raw != "" {
			path = path + "?" + raw
		}
		method := c.Request.Method
		statusCode := c.Writer.Status()
		if method != http.MethodHead {
			log.Info(fmt.Sprintf("|STATUS: %d	|Latency: %v	|Client ip: %s	|method: %s	|path: %s	",
				statusCode,
				last,
				clientIP,
				method,
				path))
		}
	}
}
