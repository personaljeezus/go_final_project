package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/personaljeezus/go_final_project/internal/utils"
	"github.com/personaljeezus/go_final_project/models"
)

type TaskStorage struct {
	db *sqlx.DB
}

func NewTaskStorage(db *sqlx.DB) *TaskStorage {
	return &TaskStorage{db: db}
}

func (t TaskStorage) CheckPostTask(task *models.Tasks) (int64, error) {
	if err := utils.CheckPostLogic(task); err != nil {
		errors.New("213")
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
	rows, err := t.db.Query("SELECT id, title, date, comment, repeat FROM scheduler ORDER BY date LIMIT ?", models.Limit)
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
			return nil, err
		}
		log.Printf("QueryRow error: %v", err)
		return nil, err
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
func (t *TaskStorage) UpdateTask(input *models.TasksInput) (int64, error) {
	if err := utils.CheckUpdateLogic(input); err != nil {
		errors.New("213")
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
		return 0, err
	}
	if rowsAffected == 0 {
		return 0, err
	}
	return res.RowsAffected()
}
func (t *TaskStorage) UpdateTaskDate(task *models.Tasks, newDate string) error {
	res, err := t.db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", newDate, task.ID)
	if err != nil {
		return errors.New("db exec fail")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return err
	}
	return err
}
func (t TaskStorage) DeleteTask(id string) error {
	res, err := t.db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return err
	}
	return nil
}
func (t TaskStorage) GetTask(id string) (models.Tasks, error) {
	var task models.Tasks
	err := t.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(
		&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, err
		} else {
			return task, err
		}
	}
	return task, err
}
