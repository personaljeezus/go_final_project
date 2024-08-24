package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/personaljeezus/go_final_project/models"
)

func CheckPostLogic(task *models.Tasks) error {
	now := time.Now()
	if task.Title == "" {
		return errors.New("Поле id пустое")
	}
	if task.Date == "" {
		task.Date = now.Format(models.DateLayout)
	}
	_, err := time.Parse(models.DateLayout, task.Date)
	if err != nil {
		return errors.New("Неверный формат даты")
	}
	if task.Date < now.Format(models.DateLayout) {
		if task.Repeat == "" {
			task.Date = now.Format(models.DateLayout)
		} else {
			newDate, err := s.Serv.NextWeekday(now, task.Date, task.Repeat)
			if err != nil {
				return errors.New("Ошибка при расчёте следующей даты")
			}
			task.Date = newDate
		}
	}
}
func CheckUpdateLogic(input *models.TasksInput) error {
	now := time.Now()
	if input.ID == "" {
		return errors.New("id requires")
	}
	id, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		return errors.New("wrong id type")
	}
	if id > 1000000 {
		return errors.New("id is too big")
	}
	if input.Title == "" {
		return errors.New("title field missing, type smth")
	}
	if input.Date == "" {
		input.Date = now.Format(models.DateLayout)
	}
	realTime, err := time.Parse(models.DateLayout, input.Date)
	if err != nil {
		return errors.New("input date parsing errored")
	}
	input.Date = realTime.Format(models.DateLayout)
	if input.Date < now.Format(models.DateLayout) && input.Repeat == "" {
		input.Date = now.Format(models.DateLayout)
	}
	if input.Date < now.Format(models.DateLayout) && input.Repeat != "" {
		newDate, err := NextWeekday(now, now.Format(models.DateLayout), input.Repeat)
		if err != nil {
			return 0, errors.New("next day calculation failed")
		}
		input.Date = newDate
	}

}
