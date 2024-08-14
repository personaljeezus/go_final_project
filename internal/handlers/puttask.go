package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/personaljeezus/go_final_project/models"
)

type Handlers struct {
	Store *service.Service
}

func NewHandler(store *service.Service) *Handlers {
	return &Handlers{Store: store}
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
