package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/personaljeezus/go_final_project/internal/service"
	"github.com/personaljeezus/go_final_project/models"
)

type TaskStorage struct {
	db *sqlx.DB
}

func NewTask(db *sqlx.DB) *TaskStorage {
	return &TaskStorage{db: db}
}

func (t *TaskStorage) UpdateTaskDate(task *models.Tasks) error {
	now := time.Now()
	currentTime, err := time.Parse(models.Layout, task.Date)
	if err != nil {
		return err
	}
	newDate, err := service.NextWeekday(now, currentTime.Format(models.Layout), task.Repeat)
	if err != nil {
		return err
	}
	res, err := t.db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", newDate, task.ID)
	if err != nil {
		return errors.New("db exec fail")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.New("rowsaffected err")
	}
	if rowsAffected == 0 {
		return errors.New("rowsaffected = 0")
	}
	return errors.New("task upd fail")
}
func (t TaskStorage) DeleteTask(id string) error {
	res, err := t.db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return errors.New("no result")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.New("rows affected error")
	}
	if rowsAffected == 0 {
		return errors.New("rows affected zero value")
	}
	return nil
}
func (t TaskStorage) GetTask(id string) (models.Tasks, error) {
	var task models.Tasks
	err := t.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(
		&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			errors.New("No rows")
		} else {
			errors.New("query fail")
		}
		return task, err
	}
	return task, err
}
