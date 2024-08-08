package handlers

import (
	"go_final_project/checkfuncs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func DeleteHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		checkfuncs.CheckID(id)
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
		return
	}
}
