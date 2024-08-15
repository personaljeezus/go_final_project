package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/personaljeezus/go_final_project/internal/database"
	"github.com/personaljeezus/go_final_project/models"
)

type Handlers struct {
	Store *database.TaskStorage
}

func NewHandler(store *database.TaskStorage) *Handlers {
	return &Handlers{Store: store}
}
func (h *Handlers) PostHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tasks models.Tasks
		if err := c.BindJSON(&tasks); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка сериализации"})
			return
		}
		taskID, err := h.Store.CheckPostTask(&tasks)
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

		t, err := h.Store.GetSingleTask(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting task"})
			return
		}
		c.JSON(http.StatusOK, t)
	}
}

func (h *Handlers) GetTasksHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tasks, err := h.Store.GetTasks()
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
		err := h.Store.DeleteTask(id)
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

		now, err := time.Parse(models.Layout, nowParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка парсинга"})
			return
		}

		date, err := time.Parse(models.Layout, dateParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Ошибка парсинга"})
			return
		}

		nextDate, err := data.NextWeekday(now, date.Format(models.Layout), repeatParam)
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
		task, err := h.Store.GetTask(id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}
		if task.Repeat != "" {
			err := h.Store.UpdateTaskDate(&task)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
			c.JSON(http.StatusOK, gin.H{})
		} else {
			if err := h.Store.DeleteTask(id); err != nil {
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
		if _, err := h.Store.InputCheck(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "input check failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}
