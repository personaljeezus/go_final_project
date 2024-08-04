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

func nextWeekday(now time.Time, date string, repeat string) (string, error) {
	parsedDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("Неверный формат даты")
	}

	if repeat == "y" {
		parsedDate = parsedDate.AddDate(1, 0, 0)
		for parsedDate.Before(now) {
			parsedDate = parsedDate.AddDate(1, 0, 0)
		}
		return parsedDate.Format("20060102"), nil
	}

	if strings.HasPrefix(repeat, "d ") {
		parts := strings.Split(repeat, " ")
		if len(parts) != 2 {
			return "", errors.New("Неверный формат правила повторения")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("Неверное количество дней")
		}
		parsedDate = parsedDate.AddDate(0, 0, days)
		for parsedDate.Before(now) {
			parsedDate = parsedDate.AddDate(0, 0, days)
		}
		return parsedDate.Format("20060102"), nil
	}

	return "", fmt.Errorf("Неверное правило повторения: %s", repeat)
}
func nextDateHandler(c *gin.Context) {
	nowParam := c.Query("now")
	dateParam := c.Query("date")
	repeatParam := c.Query("repeat")

	if nowParam == "" || dateParam == "" || repeatParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Поля параметров пусты"})
		return
	}

	now, err := time.Parse("20060102", nowParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка парсинга"})
		return
	}

	date, err := time.Parse("20060102", dateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка парсинга"})
		return
	}

	nextDate, err := nextWeekday(now, date.Format("20060102"), repeatParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, nextDate)
}
