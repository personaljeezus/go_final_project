package handlers

import (
	"database/sql"
	"go_final_project/checkfuncs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func DoneHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		log.Printf("ID: %s", id)
		task, err := checkfuncs.GetTaskByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}
		if task.Repeat != "" {
			err := checkfuncs.UpdateTaskDate(db, &task)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task date"})
				return
			}
			c.JSON(http.StatusOK, gin.H{})
		} else {
			if err := checkfuncs.DeleteTaskByID(db, id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
				return
			} else {
				c.JSON(http.StatusOK, gin.H{})
			}
		}
	}
}
