package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/personaljeezus/go_final_project/models"
)

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
