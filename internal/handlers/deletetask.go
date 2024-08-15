package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func (h *Handlers) DeleteHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Осутствует идентификатор задачи"})
			return
		}
		err := h.Store.DeleteTask(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении задачи"})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}
