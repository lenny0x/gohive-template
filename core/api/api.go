package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gohive/core/errors"
	"github.com/gohive/core/logger"
)

// Result API 统一响应结构
type Result struct {
	Success   bool   `json:"success"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	Timestamp int64  `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Result{
		Success:   true,
		Code:      "SUCCESS",
		Message:   "OK",
		Data:      data,
		Timestamp: time.Now().Unix(),
		RequestID: getRequestID(c),
	})
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	appErr, ok := errors.AsAppError(err)
	if !ok {
		appErr = errors.ErrInternal.WithCause(err)
		logger.Errorf("Unexpected error: %v", err)
	}

	if errors.IsSystemError(appErr) {
		logger.Errorf("System error: %v", appErr)
	}

	c.JSON(appErr.HTTPStatus, Result{
		Success:   false,
		Code:      appErr.Code,
		Message:   appErr.Message,
		Timestamp: time.Now().Unix(),
		RequestID: getRequestID(c),
	})
}

// Abort 中断请求并返回错误
func Abort(c *gin.Context, err error) {
	Error(c, err)
	c.Abort()
}

// getRequestID 获取请求ID
func getRequestID(c *gin.Context) string {
	if id, exists := c.Get("request_id"); exists {
		return id.(string)
	}
	return c.GetHeader("X-Request-ID")
}
