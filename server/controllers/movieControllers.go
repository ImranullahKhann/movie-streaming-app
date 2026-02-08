package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ImranullahKhann/movie-streaming-app/server/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"context"
	"time"
)

type MovieController struct {
	movieCollection *mongo.Collection
}

func NewMovieController(collection *mongo.Collection) *MovieController {
	return &MovieController{movieCollection: collection}
}

func (mc *MovieController) GetMovies(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	cursor, err := mc.movieCollection.Find(ctx, bson.D{})		

	if err != nil {
		c.JSON(500, gin.H{"error": "Can't access the database", "details": err})
		return
	}

	var movies []models.Movie
	if err = cursor.All(context.Background(), &movies); err != nil {
		c.JSON(500, gin.H{"error": "Can't read data", "details": err})
		return
	}

	c.JSON(200, gin.H{
		"movies": movies,
	})
}