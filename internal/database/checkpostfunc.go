package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/personaljeezus/go_final_project/internal/service"
	"github.com/personaljeezus/go_final_project/models"
)

func (t TaskStorage) CheckPostTask(task *models.Tasks) (error, int64) {
	now := time.Now()
	if task.Title == "" {
		return errors.New("Поле id пустое"), 0
	}
	if task.Date == "" {
		task.Date = now.Format(models.Layout)
	}
	_, err := time.Parse(models.Layout, task.Date)
	if err != nil {
		return errors.New("Неверный формат даты"), 0
	}
	if task.Date < now.Format(models.Layout) {
		if task.Repeat == "" {
			task.Date = now.Format(models.Layout)
		} else {
			newDate, err := service.NextWeekday(now, task.Date, task.Repeat)
			if err != nil {
				return errors.New("Ошибка при расчёте следующей даты"), 0
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
		return errors.New("Ошибка добавления данных в бд"), 0
	}
	task.ID, err = res.LastInsertId()
	if err != nil {
		return errors.New("Не удается получить id"), 0
	}
	return err, task.ID
}
