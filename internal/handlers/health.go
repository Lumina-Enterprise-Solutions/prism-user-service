package handlers

import (
	"net/http"

	"github.com/Lumina-Enterprise-Solutions/prism-common-libs/pkg/database"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	db *database.PostgresDB
}

func NewHealthHandler(db *database.PostgresDB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "prism-user-service",
	})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	// Check database connection
	if err := h.db.DB.Exec("SELECT 1").Error; err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"reason": "database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}
