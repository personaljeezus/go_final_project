package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/personaljeezus/go_final_project/models"
)

func (h *Handlers) PostHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tasks models.Tasks
		if err := c.BindJSON(&tasks); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка сериализации"})
			return
		}
		if _, err := h.Store.CheckPostTask(&tasks); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Task check fail"})
			return
		}
		c.JSON(http.StatusOK, &tasks.ID)
	}
}
