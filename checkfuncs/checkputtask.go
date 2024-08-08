package checkfuncs

import (
	"go_final_project/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func InputCheck(c *gin.Context, input *models.TasksInput) error {
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
	now := time.Now()
	if input.Date == "" && input.Repeat == "" {
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
	return nil
}
