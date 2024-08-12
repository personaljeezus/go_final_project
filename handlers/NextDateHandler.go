package handlers

import (
	"go_final_project/checkfuncs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NextDateHandler(c *gin.Context) {
	nowParam := c.Query("now")
	dateParam := c.Query("date")
	repeatParam := c.Query("repeat")

	if nowParam == "" || dateParam == "" || repeatParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Поля параметров пусты"})
		return
	}

	now, err := time.Parse(checkfuncs.Layout, nowParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка парсинга"})
		return
	}

	date, err := time.Parse(checkfuncs.Layout, dateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка парсинга"})
		return
	}

	nextDate, err := checkfuncs.NextWeekday(now, date.Format(checkfuncs.Layout), repeatParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, nextDate)
}
