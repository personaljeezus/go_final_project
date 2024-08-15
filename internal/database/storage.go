package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/personaljeezus/go_final_project/internal/service"
	"github.com/personaljeezus/go_final_project/models"
)

type TaskStorage struct {
	db *sqlx.DB
}

func NewTaskStorage(db *sqlx.DB) *TaskStorage {
	return &TaskStorage{db: db}
}

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
func (t TaskStorage) GetTasks() ([]map[string]string, error) {
	rows, err := t.db.Query("SELECT id, title, date, comment, repeat FROM scheduler ORDER BY date LIMIT ?", models.limit)
	if err != nil {
		return nil, errors.New("Ошибка при выполнении запроса")
	}
	defer rows.Close()
	tasks := make([]map[string]string, 0)
	for rows.Next() {
		var id int64
		var title, date, comment, repeat string
		if err := rows.Scan(&id, &title, &date, &comment, &repeat); err != nil {
			return nil, errors.New("Ошибка при чтении данных")
		}
		task := map[string]string{
			"id":      fmt.Sprint(id),
			"title":   title,
			"date":    date,
			"comment": comment,
			"repeat":  repeat,
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("Ошибка получения строк")
	}
	return tasks, nil
}
func (t TaskStorage) GetSingleTask(id string) (map[string]string, error) {
	var task models.Tasks
	err := t.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(
		&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No rows found for ID: %s", id)
			return nil, errors.New("Задание не найдено")
		}
		log.Printf("QueryRow error: %v", err)
		return nil, errors.New("Ошибка выполнения запроса")
	}
	taskMap := map[string]string{
		"id":      fmt.Sprintf("%d", task.ID),
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	}
	return taskMap, nil
}
func (t *TaskStorage) InputCheck(input *models.TasksInput) (int64, error) {
	now := time.Now()
	if input.ID == "" {
		return 0, errors.New("id requires")
	}
	id, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		return 0, errors.New("wrong id type")
	}
	if id > 1000000 {
		return 0, errors.New("id is too big")
	}
	if input.Title == "" {
		return 0, errors.New("title field missing, type smth")
	}
	if input.Date == "" {
		input.Date = now.Format(models.Layout)
	}
	realTime, err := time.Parse(models.Layout, input.Date)
	if err != nil {
		return 0, errors.New("input date parsing errored")
	}
	input.Date = realTime.Format(models.Layout)
	if input.Date < now.Format(models.Layout) && input.Repeat == "" {
		input.Date = now.Format(models.Layout)
	}
	if input.Date < now.Format(models.Layout) && input.Repeat != "" {
		newDate, err := service.NextWeekday(now, now.Format(models.Layout), input.Repeat)
		if err != nil {
			return 0, errors.New("next day calculation failed")
		}
		input.Date = newDate
	}
	task := models.Tasks{
		ID:      id,
		Date:    input.Date,
		Title:   input.Title,
		Comment: input.Comment,
		Repeat:  input.Repeat,
	}
	res, err := t.db.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return 0, errors.New("db exec fail")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, errors.New("rowsaffected err")
	}
	if rowsAffected == 0 {
		return 0, errors.New("rowsaffected = 0")
	}
	return res.RowsAffected()
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
