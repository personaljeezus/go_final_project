package checkfuncs

import (
	"database/sql"
	"go_final_project/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func CheckPostTask(c *gin.Context, db *sqlx.DB, task *models.Tasks) error {
	now := time.Now()
	if task.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поле title пустое"})
		return nil
	}
	if task.Date == "" {
		task.Date = now.Format(Layout)
	}
	_, err := time.Parse(Layout, task.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parse fail"})
		return nil
	}
	if task.Date < now.Format(Layout) {
		if task.Repeat == "" {
			task.Date = now.Format(Layout)
		} else {
			newDate, err := NextWeekday(now, task.Date, task.Repeat)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка расчёта следующей даты"})
				return nil
			}
			task.Date = newDate
		}
	}
	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Ошибка добавления данных в базу"})
	}
	task.ID, err = res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Ошибка - невозможно получить id"})
		return nil
	}
	c.JSON(http.StatusOK, gin.H{"id": task.ID})
	return nil
}
