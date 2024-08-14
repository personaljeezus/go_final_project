package database

import (
	"errors"
	"strconv"
	"time"

	"github.com/personaljeezus/go_final_project/internal/service"
	"github.com/personaljeezus/go_final_project/models"
)

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
