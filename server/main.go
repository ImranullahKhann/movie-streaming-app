package main

import (
	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator/v10"
	cont "github.com/ImranullahKhann/movie-streaming-app/server/controllers"
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

	dbClient, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	mc := cont.NewMovieController(db.OpenCollection(dbClient, "movies"))

	router.GET("/movies", mc.GetMovies)
	router.GET("/movie/:imdbID", mc.GetMovie)

	router.Run() // listens on 8080 by default
}
