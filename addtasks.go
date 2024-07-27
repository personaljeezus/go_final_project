package main

import (
	"database/sql"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Tasks struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
type Tsk struct {
	ts []Tasks `json:ts`
}

var mu sync.Mutex

func PostHandler(c *gin.Context) {
	var tasks Tasks
	today := time.Now()
	if err := c.BindJSON(&tasks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка сериализации"})
		return
	}
	if tasks.Title == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Title required"})
		return
	}
	if tasks.Date == "" && tasks.Repeat == "" {
		tasks.Date = today.Format("20060102")
	}
	realTime, err := time.Parse("20060102", tasks.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка парсинга даты, неверный формат даты"})
		return
	}
	tasks.Date = realTime.Format("20060102")

	if tasks.Date < today.Format("20060102") {
		tasks.Date = today.Format("20060102")
	}
	if tasks.Repeat != "" {
		newDate, err := nextWeekday(today, time.Now(), tasks.Repeat)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка преобразования даты"})
			return
		}
		tasks.Date = newDate.Format("20060102")
	}
	db, err := openDB()
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Невозможно открыть базу данных"})
		return
	}
	defer db.Close()

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
		c.JSON(500, gin.H{"error": "Ошибка - невозможно получить id"})
		return
	}
	c.JSON(200, gin.H{"id": tasks.ID})
}
func getTasksHandler(c *gin.Context) {
	db, err := openDB()
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Невозможно открыть базу данных"})
		return
	}
	defer db.Close()
	var limit int = 50
	rows, err := db.Query("SELECT id, title, date, comment, repeat FROM scheduler ORDER BY date LIMIT ?", limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при выполнении запроса"})
		return
	}
	defer rows.Close()
	tsk := Tsk{ts: make([]Tasks, 0)}
	for rows.Next() {
		var task Tasks
		if err := rows.Scan(&task.ID, &task.Title, &task.Date, &task.Comment, &task.Repeat); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении данных"})
			return
		}
		tsk.ts = append(tsk.ts, task)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке результатов запроса"})
		return
	}
	c.JSON(http.StatusOK, tsk)
}
func getTaskByID(c *gin.Context) {
	id := c.Query("id")
	var task Tasks
	db, err := openDB()
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Невозможно открыть базу данных"})
		return
	}
	defer db.Close()
	err = db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(
		&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не указан идентификатор"})
		}
		return
	}
	c.JSON(http.StatusOK, task)
}
func search(c *gin.Context) {
	searchParam := c.Query("search")
	if searchParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Параметр поиска не предоставлен"})
		return
	}
	search, err := time.Parse("20060102", searchParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты. Используйте YYYYMMDD"})
		return
	}
	db, err := openDB()
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Невозможно открыть базу данных"})
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?", search.Format("20060102"), 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при выполнении запроса"})
		return
	}
	defer rows.Close()
	tasks := make([]Tasks, 0)
	for rows.Next() {
		var task Tasks
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении данных"})
			return
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке результатов запроса"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}
func PutHandler(c *gin.Context) {
	var task Tasks
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка сериализации"})
		return
	}

	if task.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поле идентификатора пустое"})
		return
	}
	if task.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поле заголовок пустое"})
		return
	}
	today := time.Now()
	if task.Date == "" && task.Repeat == "" {
		task.Date = today.Format("20060102")
	}
	realTime, err := time.Parse("20060102", task.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка парсинга даты, неверный формат даты"})
		return
	}
	task.Date = realTime.Format("20060102")
	if task.Date < today.Format("20060102") {
		task.Date = today.Format("20060102")
	}
	if task.Repeat != "" {
		newDate, err := nextWeekday(today, realTime, task.Repeat)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка преобразования даты"})
			return
		}
		task.Date = newDate.Format("20060102")
	}
	db, err := openDB()
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Невозможно открыть базу данных"})
		return
	}
	defer db.Close()
	res, err := db.Exec("UPDATE scheduler SET Date = ?, Title = ?, Comment = ?, Repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Ошибка обновления данных в базе"})
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения количества затронутых строк"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}
	c.JSON(200, gin.H{"status": "Задача обновлена"})
}
func DeleteHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Осутствует идентификатор задачи"})
		return
	}
	db, err := openDB()
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Невозможно открыть базу данных"})
		return
	}
	defer db.Close()
	res, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Ошибка удаления данных из базы"})
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения количества затронутых строк"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Задача успешно удалена"})
}
