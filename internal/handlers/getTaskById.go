package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func (h *Handlers) GetTaskByID(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")

		t, err := h.Store.GetSingleTask(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting task"})
			return
		}
		c.JSON(http.StatusOK, t)
	}
}
