package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gohive/core/api"
	"github.com/gohive/core/config"
	"github.com/gohive/core/logger"
)

// Recovery 自定义错误恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				logger.Errorf("Panic recovered: %v\n%s", err, stack)

				resp := api.Result{
					Success:   false,
					Code:      "SYS_INTERNAL_ERROR",
					Message:   "Internal server error",
					Timestamp: time.Now().Unix(),
					RequestID: c.GetHeader("X-Request-ID"),
				}

				if config.IsDevelopment() {
					resp.Message = fmt.Sprintf("Panic: %v", err)
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
			}
		}()
		c.Next()
	}
}
