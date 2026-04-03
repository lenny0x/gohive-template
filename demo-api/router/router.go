package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gohive/demo-api/handler"
)

func Register(r *gin.Engine) {
	healthHandler := handler.NewHealthHandler()
	r.GET("/health", healthHandler.Check)

	// Admin API v1
	v1 := r.Group("/admin/v1")
	{
		// TODO: Add admin routes
		_ = v1
	}
}
