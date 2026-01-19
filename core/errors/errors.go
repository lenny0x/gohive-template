package errors

import (
	"fmt"
	"net/http"
)

// ErrorType 错误类型
type ErrorType string

const (
	// SystemError 系统错误 - 服务端内部错误，不应暴露详情给客户端
	SystemError ErrorType = "SYSTEM_ERROR"
	// BizError 业务错误 - 业务逻辑错误，可以展示给客户端
	BizError ErrorType = "BIZ_ERROR"
	// ValidationError 参数校验错误
	ValidationError ErrorType = "VALIDATION_ERROR"
	// AuthError 认证/授权错误
	AuthError ErrorType = "AUTH_ERROR"
)

// AppError 应用错误
type AppError struct {
	Type       ErrorType `json:"-"`
	Code       string    `json:"code"`
	Message    string    `json:"message"`
	Detail     string    `json:"-"` // 内部详情，不返回给客户端
	HTTPStatus int       `json:"-"`
	Cause      error     `json:"-"` // 原始错误
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %s (cause: %v)", e.Code, e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

// Is 实现 errors.Is 接口
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// WithCause 添加原始错误
func (e *AppError) WithCause(err error) *AppError {
	newErr := *e
	newErr.Cause = err
	return &newErr
}

// WithDetail 添加内部详情
func (e *AppError) WithDetail(detail string) *AppError {
	newErr := *e
	newErr.Detail = detail
	return &newErr
}

// WithMessage 覆盖错误消息
func (e *AppError) WithMessage(msg string) *AppError {
	newErr := *e
	newErr.Message = msg
	return &newErr
}

// ============================================================
// 系统错误定义 (5xx)
// ============================================================

var (
	// ErrInternal 内部服务器错误
	ErrInternal = &AppError{
		Type:       SystemError,
		Code:       "SYS_INTERNAL_ERROR",
		Message:    "Internal server error",
		HTTPStatus: http.StatusInternalServerError,
	}

	// ErrDatabase 数据库错误
	ErrDatabase = &AppError{
		Type:       SystemError,
		Code:       "SYS_DATABASE_ERROR",
		Message:    "Database error",
		HTTPStatus: http.StatusInternalServerError,
	}

	// ErrCache 缓存错误
	ErrCache = &AppError{
		Type:       SystemError,
		Code:       "SYS_CACHE_ERROR",
		Message:    "Cache error",
		HTTPStatus: http.StatusInternalServerError,
	}

	// ErrExternalService 外部服务错误
	ErrExternalService = &AppError{
		Type:       SystemError,
		Code:       "SYS_EXTERNAL_SERVICE_ERROR",
		Message:    "External service error",
		HTTPStatus: http.StatusBadGateway,
	}

	// ErrServiceUnavailable 服务不可用
	ErrServiceUnavailable = &AppError{
		Type:       SystemError,
		Code:       "SYS_SERVICE_UNAVAILABLE",
		Message:    "Service temporarily unavailable",
		HTTPStatus: http.StatusServiceUnavailable,
	}

	// ErrTimeout 请求超时
	ErrTimeout = &AppError{
		Type:       SystemError,
		Code:       "SYS_TIMEOUT",
		Message:    "Request timeout",
		HTTPStatus: http.StatusGatewayTimeout,
	}
)

// ============================================================
// 认证/授权错误定义 (401, 403)
// ============================================================

var (
	// ErrUnauthorized 未认证
	ErrUnauthorized = &AppError{
		Type:       AuthError,
		Code:       "AUTH_UNAUTHORIZED",
		Message:    "Unauthorized",
		HTTPStatus: http.StatusUnauthorized,
	}

	// ErrTokenExpired Token 过期
	ErrTokenExpired = &AppError{
		Type:       AuthError,
		Code:       "AUTH_TOKEN_EXPIRED",
		Message:    "Token has expired",
		HTTPStatus: http.StatusUnauthorized,
	}

	// ErrTokenInvalid Token 无效
	ErrTokenInvalid = &AppError{
		Type:       AuthError,
		Code:       "AUTH_TOKEN_INVALID",
		Message:    "Token is invalid",
		HTTPStatus: http.StatusUnauthorized,
	}

	// ErrForbidden 无权限
	ErrForbidden = &AppError{
		Type:       AuthError,
		Code:       "AUTH_FORBIDDEN",
		Message:    "Access forbidden",
		HTTPStatus: http.StatusForbidden,
	}
)

// ============================================================
// 参数校验错误定义 (400)
// ============================================================

var (
	// ErrBadRequest 请求参数错误
	ErrBadRequest = &AppError{
		Type:       ValidationError,
		Code:       "VALIDATION_BAD_REQUEST",
		Message:    "Bad request",
		HTTPStatus: http.StatusBadRequest,
	}

	// ErrInvalidParam 参数无效
	ErrInvalidParam = &AppError{
		Type:       ValidationError,
		Code:       "VALIDATION_INVALID_PARAM",
		Message:    "Invalid parameter",
		HTTPStatus: http.StatusBadRequest,
	}

	// ErrMissingParam 缺少参数
	ErrMissingParam = &AppError{
		Type:       ValidationError,
		Code:       "VALIDATION_MISSING_PARAM",
		Message:    "Missing required parameter",
		HTTPStatus: http.StatusBadRequest,
	}
)

// ============================================================
// 业务错误定义 (4xx)
// ============================================================

var (
	// ErrNotFound 资源不存在
	ErrNotFound = &AppError{
		Type:       BizError,
		Code:       "BIZ_NOT_FOUND",
		Message:    "Resource not found",
		HTTPStatus: http.StatusNotFound,
	}

	// ErrAlreadyExists 资源已存在
	ErrAlreadyExists = &AppError{
		Type:       BizError,
		Code:       "BIZ_ALREADY_EXISTS",
		Message:    "Resource already exists",
		HTTPStatus: http.StatusConflict,
	}

	// ErrConflict 资源冲突
	ErrConflict = &AppError{
		Type:       BizError,
		Code:       "BIZ_CONFLICT",
		Message:    "Resource conflict",
		HTTPStatus: http.StatusConflict,
	}

	// ErrTooManyRequests 请求过于频繁
	ErrTooManyRequests = &AppError{
		Type:       BizError,
		Code:       "BIZ_TOO_MANY_REQUESTS",
		Message:    "Too many requests",
		HTTPStatus: http.StatusTooManyRequests,
	}
)

// ============================================================
// 工厂函数
// ============================================================

// NewBizError 创建业务错误
func NewBizError(code, message string) *AppError {
	return &AppError{
		Type:       BizError,
		Code:       code,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

// NewBizErrorWithStatus 创建带 HTTP 状态码的业务错误
func NewBizErrorWithStatus(code, message string, httpStatus int) *AppError {
	return &AppError{
		Type:       BizError,
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// NewSystemError 创建系统错误
func NewSystemError(code, message string) *AppError {
	return &AppError{
		Type:       SystemError,
		Code:       code,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
	}
}

// NewValidationError 创建参数校验错误
func NewValidationError(code, message string) *AppError {
	return &AppError{
		Type:       ValidationError,
		Code:       code,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

// IsSystemError 判断是否为系统错误
func IsSystemError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == SystemError
	}
	return false
}

// IsBizError 判断是否为业务错误
func IsBizError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == BizError
	}
	return false
}

// AsAppError 转换为 AppError
func AsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}
