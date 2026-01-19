package errors

import "net/http"

// ============================================================
// 用户模块业务错误 (BIZ_USER_*)
// ============================================================

var (
	ErrUserNotFound = &AppError{
		Type:       BizError,
		Code:       "BIZ_USER_NOT_FOUND",
		Message:    "User not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrUserAlreadyExists = &AppError{
		Type:       BizError,
		Code:       "BIZ_USER_ALREADY_EXISTS",
		Message:    "User already exists",
		HTTPStatus: http.StatusConflict,
	}

	ErrUserDisabled = &AppError{
		Type:       BizError,
		Code:       "BIZ_USER_DISABLED",
		Message:    "User account is disabled",
		HTTPStatus: http.StatusForbidden,
	}

	ErrInvalidPassword = &AppError{
		Type:       BizError,
		Code:       "BIZ_INVALID_PASSWORD",
		Message:    "Invalid password",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrPasswordTooWeak = &AppError{
		Type:       BizError,
		Code:       "BIZ_PASSWORD_TOO_WEAK",
		Message:    "Password is too weak",
		HTTPStatus: http.StatusBadRequest,
	}
)

// ============================================================
// 订单模块业务错误 (BIZ_ORDER_*)
// ============================================================

var (
	ErrOrderNotFound = &AppError{
		Type:       BizError,
		Code:       "BIZ_ORDER_NOT_FOUND",
		Message:    "Order not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrOrderAlreadyPaid = &AppError{
		Type:       BizError,
		Code:       "BIZ_ORDER_ALREADY_PAID",
		Message:    "Order has already been paid",
		HTTPStatus: http.StatusConflict,
	}

	ErrOrderExpired = &AppError{
		Type:       BizError,
		Code:       "BIZ_ORDER_EXPIRED",
		Message:    "Order has expired",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrOrderCancelled = &AppError{
		Type:       BizError,
		Code:       "BIZ_ORDER_CANCELLED",
		Message:    "Order has been cancelled",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrInsufficientStock = &AppError{
		Type:       BizError,
		Code:       "BIZ_INSUFFICIENT_STOCK",
		Message:    "Insufficient stock",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrOrderAmountInvalid = &AppError{
		Type:       BizError,
		Code:       "BIZ_ORDER_AMOUNT_INVALID",
		Message:    "Order amount is invalid",
		HTTPStatus: http.StatusBadRequest,
	}
)

// ============================================================
// 支付模块业务错误 (BIZ_PAYMENT_*)
// ============================================================

var (
	ErrPaymentFailed = &AppError{
		Type:       BizError,
		Code:       "BIZ_PAYMENT_FAILED",
		Message:    "Payment failed",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrInsufficientBalance = &AppError{
		Type:       BizError,
		Code:       "BIZ_INSUFFICIENT_BALANCE",
		Message:    "Insufficient balance",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrPaymentTimeout = &AppError{
		Type:       BizError,
		Code:       "BIZ_PAYMENT_TIMEOUT",
		Message:    "Payment timeout",
		HTTPStatus: http.StatusRequestTimeout,
	}
)
