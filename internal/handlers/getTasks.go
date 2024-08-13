package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetTasksHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var limit int = 50
		rows, err := db.Query("SELECT id, title, date, comment, repeat FROM scheduler ORDER BY date LIMIT ?", limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при выполнении запроса"})
			return
		}
		defer rows.Close()
		tasks := make([]map[string]string, 0)
		for rows.Next() {
			var id int64
			var title, date, comment, repeat string
			if err := rows.Scan(&id, &title, &date, &comment, &repeat); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении данных"})
				return
			}
			task := map[string]string{
				"id":      fmt.Sprint(id),
				"title":   title,
				"date":    date,
				"comment": comment,
				"repeat":  repeat,
			}
			tasks = append(tasks, task)
		}
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка rows"})
			return
		}
		if tasks == nil {
			tasks = make([]map[string]string, 0)
		}
		c.JSON(http.StatusOK, gin.H{"tasks": tasks})
		return
	}
}
