package handlers

import (
	"go_final_project/checkfuncs"
	"net/http"

	"github.com/personaljeezus/go_final_project/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func PostHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tasks models.Tasks
		if err := c.BindJSON(&tasks); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка сериализации"})
			return
		}
		if err := checkfuncs.CheckPostTask(c, db, &tasks); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Task check fail"})
			return
		}
	}
}
