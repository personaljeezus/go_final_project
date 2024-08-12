package checkfuncs

import (
	"database/sql"
	"go_final_project/models"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

func UpdateTaskDate(db *sqlx.DB, task *models.Tasks) error {
	now := time.Now()
	currentTime, err := time.Parse(Layout, task.Date)
	if err != nil {
		return err
	}
	newDate, err := NextWeekday(now, currentTime.Format(Layout), task.Repeat)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", newDate, task.ID)
	return err
}
func DeleteTaskByID(db *sqlx.DB, id string) error {
	res, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Ошибка получения результата запроса: %v", err)
		return nil
	}
	if rowsAffected == 0 {
		return nil
	}
	return nil
}
func GetTaskByID(db *sqlx.DB, id string) (models.Tasks, error) {
	var task models.Tasks
	err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(
		&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Задание не найдено: %s", id)
		} else {
			log.Printf("Ошибка бд: %v", err)
		}
		return task, err
	}
	log.Printf("Фулл таска: %+v", task)
	return task, err
}
