package controllers

import (
	"context"
	"errors"
	"github.com/ImranullahKhann/movie-streaming-app/server/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func (mc *MovieController) GetMovie(c *gin.Context) {
	imdbID := c.Param("imdbID")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var movie models.Movie

	err := mc.movieCollection.FindOne(ctx,
		bson.D{{Key: "imdb_id", Value: imdbID}},
	).Decode(&movie)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(404, gin.H{"error": "Movie not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Can't read data", "details": err})
		return
	}

	c.JSON(200, gin.H{"movie": movie})
}
