package main

import (
	"go_final_project/database"
	"go_final_project/handlers"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("ENV_PATH")
	defaultPort := os.Getenv("PORT")
	if defaultPort == "" {
		defaultPort = "7540"
	}
	db, err := database.DbCheck()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db, err = database.OpenDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	r := gin.Default()
	api := r.Group("/api")
	{
		api.GET("/nextdate", handlers.NextDateHandler)
		api.GET("/tasks", handlers.GetTasksHandler(db))
		api.POST("/task", handlers.PostHandler(db))
		api.GET("/task", handlers.GetTaskByID(db))
		api.PUT("/task", handlers.PutHandler(db))
		api.DELETE("/task", handlers.DeleteHandler(db))
		api.POST("/task/done", handlers.DeleteHandler(db))
	}
	r.Static("/js", "./web/js")
	r.Static("/css", "./web/css")
	r.StaticFile("favicon.ico", "./web/favicon.ico")
	r.LoadHTMLFiles("./web/index.html", "./web/login.html")
	r.GET("/index.html", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.GET("/login.html", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	err = r.Run(":" + defaultPort)
	var c *gin.Context
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "err"})
	}
}
