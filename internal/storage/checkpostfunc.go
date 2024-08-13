package service

import (
	"database/sql"
	"errors"
	"go_final_project/internal/models"
	. "go_final_project/internal/storage"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (t TaskStorage) CheckPostTask(task *models.Tasks) error {
	now := time.Now()
	if task.Title == "" {
		return errors.New("Поле id пустое")
	}
	if task.Date == "" {
		task.Date = now.Format(Layout)
	}
	_, err := time.Parse(Layout, task.Date)
	if err != nil {
		return errors.New("Неверный формат даты")
	}
	if task.Date < now.Format(Layout) {
		if task.Repeat == "" {
			task.Date = now.Format(Layout)
		} else {
			newDate, err := NextWeekday(now, task.Date, task.Repeat)
			if err != nil {
				return errors.New("Ошибка при расчёте следующей даты")
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
		return errors.New("Ошибка добавления данных в бд")
	}
	task.ID, err = res.LastInsertId()
	if err != nil {
		return errors.New("Не удается получить id")
	}
	c.JSON(http.StatusOK, gin.H{"id": task.ID})
	return nil
}
