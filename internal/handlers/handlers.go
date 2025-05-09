package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/personaljeezus/go_final_project/internal/service"
	"github.com/personaljeezus/go_final_project/models"
)

type Handlers struct {
	Store *service.TaskService
}

func NewHandler(store *service.TaskService) *Handlers {
	return &Handlers{Store: store}
}
func (h *Handlers) PostHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tasks models.Tasks
		if err := c.BindJSON(&tasks); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка сериализации"})
			return
		}
		taskID, err := h.Store.Serv.CheckPostTask(&tasks)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Task check fail"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": taskID})
	}
}

func (h *Handlers) GetTaskByID(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")

		t, err := h.Store.Serv.GetSingleTask(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting task"})
			return
		}
		c.JSON(http.StatusOK, t)
	}
}

func (h *Handlers) GetTasksHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tasks, err := h.Store.Serv.GetTasks()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "GetTasks func fail"})
		}
		c.JSON(http.StatusOK, gin.H{"tasks": tasks})
	}
}
func (h *Handlers) DeleteHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Осутствует идентификатор задачи"})
			return
		}
		err := h.Store.Serv.DeleteTask(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении задачи"})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}
func (h *Handlers) NextDate() gin.HandlerFunc {
	return func(c *gin.Context) {
		nowParam := c.Query("now")
		dateParam := c.Query("date")
		repeatParam := c.Query("repeat")

		if nowParam == "" || dateParam == "" || repeatParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Поля параметров пусты"})
			return
		}

		now, err := time.Parse(models.DateLayout, nowParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка парсинга"})
			return
		}

		date, err := time.Parse(models.DateLayout, dateParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка парсинга"})
			return
		}

		nextDate, err := h.Store.NextWeekday(now, date.Format(models.DateLayout), repeatParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.String(http.StatusOK, nextDate)
	}
}
func (h *Handlers) DoneHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "id not found"})
		}
		task, err := h.Store.Serv.GetTask(id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}
		now := time.Now()
		currentTime, err := time.Parse(models.DateLayout, task.Date)
		if err != nil {
			return
		}
		newDate, err := h.Store.NextWeekday(now, currentTime.Format(models.DateLayout), task.Repeat)
		if err != nil {
			return
		}
		if task.Repeat != "" {
			err := h.Store.Serv.UpdateTaskDate(&task, newDate)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
			c.JSON(http.StatusOK, gin.H{})
		} else {
			if err := h.Store.Serv.DeleteTask(id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			} else {
				c.JSON(http.StatusOK, gin.H{})
			}
		}
	}
}

func (h *Handlers) PutHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.TasksInput
		if err := c.BindJSON(&input); err != nil {
			log.Printf("Ошибка сериализации: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка сериализации"})
			return
		}
		if _, err := h.Store.Serv.UpdateTask(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "input check failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}
