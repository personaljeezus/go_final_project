package handlers

import (
	"database/sql"
	"go_final_project/checkfuncs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func DoneHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if err := checkfuncs.CheckID(id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		task, err := checkfuncs.GetTaskByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		if task.Repeat == "" {
			if err := checkfuncs.DeleteTaskByID(db, id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
			} else {
				c.JSON(http.StatusOK, gin.H{})
			}
		} else {
			if err := checkfuncs.UpdateTaskDate(db, &task); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task date"})
			} else {
				c.JSON(http.StatusOK, gin.H{})
			}
		}
	}
}
