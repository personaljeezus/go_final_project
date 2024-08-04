package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

var mu sync.Mutex

func PostHandler(c *gin.Context) {
	var tasks Tasks
	now := time.Now()
	if err := c.BindJSON(&tasks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка сериализации"})
		return
	}
	if tasks.Title == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Title required"})
		return
	}
	if tasks.Date == "" {
		tasks.Date = now.Format("20060102")
	}
	realTime, err := time.Parse("20060102", tasks.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка парсинга даты, неверный формат даты"})
		return
	}
	tasks.Date = realTime.Format("20060102")

	if tasks.Date < now.Format("20060102") && tasks.Repeat == "" {
		tasks.Date = now.Format("20060102")
	}
	if tasks.Date < now.Format("20060102") && tasks.Repeat != "" {
		newDate, err := nextWeekday(now, now.Format("20060102"), tasks.Repeat)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка преобразования даты"})
			return
		}
		tasks.Date = newDate
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
	tasks := make([]map[string]string, 0)
	for rows.Next() {
		var id int64
		var title, date, comment, repeat string
		if err := rows.Scan(&id, &title, &date, &comment, &repeat); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении данных"})
			return
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка rows"})
		return
	}
	if tasks == nil {
		tasks = make([]map[string]string, 0)
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}
func getTaskByID(c *gin.Context) {
	id := c.Query("id")
	var task Tasks
	db, err := openDB()
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Невозможно открыть бд"})
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
	taskMap := map[string]string{
		"id":      fmt.Sprintf("%d", task.ID),
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	}

	c.JSON(http.StatusOK, taskMap)
}

type TasksInput struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func PutHandler(c *gin.Context) {
	var input TasksInput
	if err := c.BindJSON(&input); err != nil {
		log.Printf("Ошибка сериализации: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка сериализации"})
		return
	}
	if input.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поле id пустое"})
		return
	}
	id, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат id"})
		return
	}
	if id > 1000000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат id"})
		return
	}
	if input.Title == "" {
		c.JSON(http.StatusAccepted, gin.H{"error": "Поле заголовка пустое"})
		return
	}
	now := time.Now()
	if input.Date == "" && input.Repeat == "" {
		input.Date = now.Format("20060102")
	}

	realTime, err := time.Parse("20060102", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка парсинга даты"})
		return
	}
	input.Date = realTime.Format("20060102")
	if input.Date < now.Format("20060102") && input.Repeat == "" {
		input.Date = now.Format("20060102")
	}
	if input.Date < now.Format("20060102") && input.Repeat != "" {
		newDate, err := nextWeekday(now, now.Format("20060102"), input.Repeat)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка расчёта даты"})
			return
		}
		input.Date = newDate
	}
	task := Tasks{
		ID:      id,
		Date:    input.Date,
		Title:   input.Title,
		Comment: input.Comment,
		Repeat:  input.Repeat,
	}
	db, err := openDB()
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Невозможно открыть бд"})
		return
	}
	defer db.Close()
	res, err := db.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Ошибка обновления данных в бд"})
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения количества затронутых строк"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
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
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Ошибка удаления данных из бд"})
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "rowsaffected err"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "rowsaffected - 0"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
func DoneHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id req"})
		return
	}
	log.Printf("ID: %s", id)
	db, err := openDB()
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Невозможно открыть бд"})
		return
	}
	defer db.Close()
	var task Tasks
	err = db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(
		&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Задание не найдено: %s", id)
			c.JSON(http.StatusNotFound, gin.H{})
		} else {
			log.Printf("Ошибка бд: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
		}
		return
	}
	log.Printf("Фулл таска: %+v", task)
	now := time.Now()
	if task.Repeat == "" {
		_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			log.Printf("Задача не удалилась: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления задачи"})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
	if task.Repeat != "" {
		currentTime, err := time.Parse("20060102", task.Date)
		if err != nil {
			log.Printf("Ошибка парсинга даты 301: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка парсинга даты"})
			return
		}
		newDate, err := nextWeekday(now, currentTime.Format("20060102"), task.Repeat)
		if err != nil {
			log.Printf("Расчет некстдейт фейлится: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка расчёта даты"})
			return
		}
		_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", newDate, id)
		if err != nil {
			log.Printf("Ошибка апдейта даты 313: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления даты задачи"})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}
