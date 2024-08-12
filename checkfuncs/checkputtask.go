package checkfuncs

import (
	"errors"
	"go_final_project/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func InputCheck(c *gin.Context, db *sqlx.DB, input *models.TasksInput) error {
	now := time.Now()
	if input.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поле id пустое"})
		return nil
	}
	id, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат id"})
		return nil
	}
	if id > 1000000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат id"})
		return nil
	}
	if input.Title == "" {
		c.JSON(http.StatusAccepted, gin.H{"error": "Поле заголовка пустое"})
		return nil
	}
	if input.Date == "" {
		input.Date = now.Format(Layout)
	}
	realTime, err := time.Parse(Layout, input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка парсинга даты"})
		return nil
	}
	input.Date = realTime.Format(Layout)
	if input.Date < now.Format(Layout) && input.Repeat == "" {
		input.Date = now.Format(Layout)
	}
	if input.Date < now.Format(Layout) && input.Repeat != "" {
		newDate, err := NextWeekday(now, now.Format(Layout), input.Repeat)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка расчёта даты"})
			return nil
		}
		input.Date = newDate
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
		return errors.New("db exec fail")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения количества затронутых строк"})
		return errors.New("rowsaffected err")
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return errors.New("rowsaffected = 0")
	}
	c.JSON(http.StatusOK, gin.H{})
	return nil
}
