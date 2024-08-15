package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/personaljeezus/go_final_project/internal/service"
	"github.com/personaljeezus/go_final_project/models"
)

func (t TaskStorage) CheckPostTask(task *models.Tasks) (int64, error) {
	now := time.Now()
	if task.Title == "" {
		return 0, errors.New("Поле id пустое")
	}
	if task.Date == "" {
		task.Date = now.Format(models.Layout)
	}
	_, err := time.Parse(models.Layout, task.Date)
	if err != nil {
		return 0, errors.New("Неверный формат даты")
	}
	if task.Date < now.Format(models.Layout) {
		if task.Repeat == "" {
			task.Date = now.Format(models.Layout)
		} else {
			newDate, err := service.NextWeekday(now, task.Date, task.Repeat)
			if err != nil {
				return 0, errors.New("Ошибка при расчёте следующей даты")
			}
			task.Date = newDate
		}
	}
	res, err := t.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, errors.New("Ошибка добавления данных в бд")
	}
	taskID, err := res.LastInsertId()
	if err != nil {
		return 0, errors.New("Не удается получить id")
	}
	task.ID = taskID
	return taskID, err
}
