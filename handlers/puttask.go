package handlers

import (
	"go_final_project/checkfuncs"
	"go_final_project/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func PutHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.TasksInput
		if err := c.BindJSON(&input); err != nil {
			log.Printf("Ошибка сериализации: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка сериализации"})
			return
		}
		id, err := strconv.ParseInt(input.ID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат id"})
			return
		}
		if err := checkfuncs.InputCheck(c, &input); err != nil {
			return
		}
		task := models.Tasks{
			ID:      id,
			Date:    input.Date,
			Title:   input.Title,
			Comment: input.Comment,
			Repeat:  input.Repeat,
		}
		res, err := db.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
			task.Date, task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Ошибка обновления данных в базе данных"})
			return
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения количества затронутых строк"})
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}
