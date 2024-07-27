package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

type Weekday int

const (
	layout         = "20060102"
	Monday Weekday = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func nextWeekday(now time.Time, date time.Time, repeat string) (time.Time, error) {
	if repeat == "y" {
		return date.AddDate(1, 0, 0), nil
	}
	if strings.HasPrefix(repeat, "d ") {
		parts := strings.Split(repeat, " ")
		if len(parts) != 2 {
			return time.Time{}, errors.New("Неверный формат правила повторения")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return time.Time{}, errors.New("Неверное количество дней")
		}
		return date.AddDate(0, 0, days), nil
	}
	if repeat == "" {
		db, err := openDB()
		if err != nil {
			return time.Time{}, err
		}
		_, err = db.Exec("DELETE FROM scheduler WHERE repeat IS NULL")
		if err != nil {
			return time.Time{}, fmt.Errorf("Поле repeat пустое: %s", repeat)
		}
	}
	return time.Time{}, fmt.Errorf("Неверное правило повторения: %s", repeat)
}
func nextDateHandler(c *gin.Context) {
	nowParam := c.Query("now")
	dateParam := c.Query("date")
	repeatParam := c.Query("repeat")

	if nowParam == "" || dateParam == "" || repeatParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Params fields are nil"})
		return
	}

	now, err := time.Parse("20060102", nowParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Now param parse error"})
		return
	}

	date, err := time.Parse("20060102", dateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Date param parse error"})
		return
	}

	nextDate, err := nextWeekday(now, date, repeatParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, nextDate.Format("20060102"))
}
