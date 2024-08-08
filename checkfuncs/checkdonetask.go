package checkfuncs

import (
	"go_final_project/models"
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
	_, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	return err
}
func GetTaskByID(db *sqlx.DB, id string) (models.Tasks, error) {
	var task models.Tasks
	err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(
		&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	return task, err
}
