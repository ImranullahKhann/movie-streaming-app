package main

import (
	cont "github.com/ImranullahKhann/movie-streaming-app/server/controllers"
	db "github.com/ImranullahKhann/movie-streaming-app/server/database"
	"github.com/ImranullahKhann/movie-streaming-app/server/middleware"
	"github.com/ImranullahKhann/movie-streaming-app/server/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	router := gin.Default()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	for _, k := range []string{"ACCESS_SECRET", "REFRESH_SECRET"} {
		if os.Getenv(k) == "" {
			log.Fatalf("%s not set", k)
		}
	}

	rds := store.NewRedis()

	dbClient, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FRONTEND_ORIGIN")},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	mc := cont.NewMovieController(db.OpenCollection(dbClient, "movies"))
	uc := cont.NewUserController(db.OpenCollection(dbClient, "users"), rds)

	router.GET("/movies", mc.GetMovies)
	router.GET("/movies/:imdbID", mc.GetMovie)
	router.POST("/movies/", mc.AddMovie)

	router.POST("/register/", middleware.AuthMiddleware(rds), uc.RegisterUser)
	router.POST("/login/", uc.LoginUser)
	router.GET("/logout/", middleware.AuthMiddleware(rds), uc.LogoutUser)
	router.GET("/token/refresh", uc.RefreshTokens)

	router.Run() // listens on 8080 by default
}
