package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/personaljeezus/go_final_project/internal/database"
	"github.com/personaljeezus/go_final_project/models"
	_ "modernc.org/sqlite"
)

type TaskService struct {
	Serv *database.TaskStorage
}

func NewTaskService(serv *database.TaskStorage) *TaskService {
	return &TaskService{Serv: serv}
}
func (s TaskService) NextWeekday(now time.Time, date string, repeat string) (string, error) {
	parsedDate, err := time.Parse(models.DateLayout, date)
	if err != nil {
		return "", errors.New("Неверный формат даты")
	}
	if repeat == "y" {
		parsedDate = parsedDate.AddDate(1, 0, 0)
		for parsedDate.Before(now) {
			parsedDate = parsedDate.AddDate(1, 0, 0)
		}
		return parsedDate.Format(models.DateLayout), nil
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
		return parsedDate.Format(models.DateLayout), nil
	}

	return "", fmt.Errorf("Неверное правило повторения: %s", repeat)
}
