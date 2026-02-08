package main

import (
	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator/v10"
	db "github.com/ImranullahKhann/movie-streaming-app/server/database"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	router := gin.Default()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	db.ConnectDB()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	
	router.Run() // listens on 8080 by default
}
