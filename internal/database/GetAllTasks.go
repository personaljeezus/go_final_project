package database

import (
	"errors"
	"fmt"
)

func (t TaskStorage) GetTasks() ([]map[string]string, error) {
	var limit int = 50
	rows, err := t.db.Query("SELECT id, title, date, comment, repeat FROM scheduler ORDER BY date LIMIT ?", limit)
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
	if tasks == nil {
		tasks = make([]map[string]string, 0)
	}
	return tasks, errors.New("Ошибка получения задач")
}
