package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func (h *Handlers) GetTasksHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tasks, err := h.Store.GetTasks()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "GetTasks func fail"})
		}
		c.JSON(http.StatusOK, gin.H{"tasks": tasks})
		return
	}
}
