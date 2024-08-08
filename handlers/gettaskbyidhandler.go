package handlers

import (
	"database/sql"
	"fmt"
	"go_final_project/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetTaskByID(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		var task models.Tasks
		err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(
			&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Не указан идентификатор"})
			}
			return
		}
		taskMap := map[string]string{
			"id":      fmt.Sprintf("%d", task.ID),
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		}
		c.JSON(http.StatusOK, taskMap)
		return
	}
}
