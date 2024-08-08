package checkfuncs

import (
	"database/sql"
	"go_final_project/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func CheckPostTask(c *gin.Context, db *sqlx.DB, tasks *models.Tasks) error {
	now := time.Now()
	if tasks.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поле title пустое"})
		return nil
	}
	if tasks.Date == "" {
		tasks.Date = now.Format(Layout)
	}
	realTime, err := time.Parse(Layout, tasks.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parse failed"})
		return nil
	}
	tasks.Date = realTime.Format(Layout)

	if tasks.Date < now.Format(Layout) && tasks.Repeat == "" {
		tasks.Date = now.Format(Layout)
	}
	if tasks.Date < now.Format(Layout) && tasks.Repeat != "" {
		newDate, err := NextWeekday(now, now.Format(Layout), tasks.Repeat)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка расчёта следующей даты"})
			return nil
		}
		tasks.Date = newDate
	}
	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", tasks.Date),
		sql.Named("title", tasks.Title),
		sql.Named("comment", tasks.Comment),
		sql.Named("repeat", tasks.Repeat))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Ошибка добавления данных в базу"})
	}
	tasks.ID, err = res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Ошибка - невозможно получить id"})
		return nil
	}
	c.JSON(http.StatusOK, gin.H{"id": tasks.ID})
	return nil
}
