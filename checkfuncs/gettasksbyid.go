package checkfuncs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckID(string) (c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pole id pustoe"})
	}
	return
}
