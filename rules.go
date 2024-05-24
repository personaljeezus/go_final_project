package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

func nextWeekday(now, date time.Time, repeat string) (time.Time, error) {
	if repeat == "y" {
		return date.AddDate(1, 0, 0), nil
	}
	if strings.HasPrefix(repeat, "d ") {
		parts := strings.Split(repeat, " ")
		if len(parts) != 2 {
			return time.Time{}, fmt.Errorf("Неверный формат правила повторения: %s", repeat)
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return time.Time{}, fmt.Errorf("Неверное количество дней: %s", parts[1])
		}
		return date.AddDate(0, 0, days), nil
	}
	if repeat == "" {
		_, err := db.Exec("DELETE FROM scheduler WHERE repeat IS NULL")
		if err != nil {
			return time.Time{}, fmt.Errorf("Поле repeat пустое: %s", repeat)
		}
	}
	return time.Time{}, fmt.Errorf("Неверное правило повторения: %s", repeat)
}
