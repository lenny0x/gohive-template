package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gohive/core/api"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	api.Success(c, gin.H{
		"status": "ok",
	})
}
