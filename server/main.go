package main

import (
	cont "github.com/ImranullahKhann/movie-streaming-app/server/controllers"
	db "github.com/ImranullahKhann/movie-streaming-app/server/database"
	"github.com/gin-gonic/gin"
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
	uc := cont.NewUserController(db.OpenCollection(dbClient, "users"))

	router.GET("/movies", mc.GetMovies)
	router.GET("/movies/:imdbID", mc.GetMovie)
	router.POST("/movies/", mc.AddMovie)

	router.POST("/users/", uc.RegisterUser)

	router.Run() // listens on 8080 by default
}
