package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func (h *Handlers) GetTaskByID(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			log.Fatal("id field missing")
		}
		t, err := h.Store.GetSingleTask(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting task"})
		}
		c.JSON(http.StatusOK, t)
	}
}
