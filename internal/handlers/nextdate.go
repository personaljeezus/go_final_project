package handlers

import (
	"net/http"
	"time"

	"github.com/personaljeezus/go_final_project/internal/service"
	"github.com/personaljeezus/go_final_project/models"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) nextDate() gin.HandlerFunc {
	return func(c *gin.Context) {
		nowParam := c.Query("now")
		dateParam := c.Query("date")
		repeatParam := c.Query("repeat")

		if nowParam == "" || dateParam == "" || repeatParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Поля параметров пусты"})
			return
		}

		now, err := time.Parse(models.Layout, nowParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка парсинга"})
			return
		}

		date, err := time.Parse(models.Layout, dateParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка парсинга"})
			return
		}

		nextDate, err := service.NextWeekday(now, date.Format(models.Layout), repeatParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.String(http.StatusOK, nextDate)
	}
}
