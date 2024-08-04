package main

import (
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
	db, err := dbCheck()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	r := gin.Default()
	api := r.Group("/api")
	{
		api.GET("/nextdate", nextDateHandler)
		api.GET("/tasks", getTasksHandler)
		api.POST("/task", PostHandler)
		api.GET("/task", getTaskByID)
		api.PUT("/task", PutHandler)
		api.DELETE("/task", DeleteHandler)
		api.POST("/task/done", DoneHandler)
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
		c.JSON(404, gin.H{"message": "err"})
	}
}
