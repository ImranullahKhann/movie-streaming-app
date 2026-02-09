package controllers

import (
	"context"
	"errors"
	"github.com/ImranullahKhann/movie-streaming-app/server/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
	"github.com/go-playground/validator/v10"
)

type MovieController struct {
	movieCollection *mongo.Collection
	validate *validator.Validate	
}

func NewMovieController(collection *mongo.Collection) *MovieController {
	return &MovieController{movieCollection: collection, validate: validator.New()}
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
	if err = cursor.All(ctx, &movies); err != nil {
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

func (mc *MovieController) AddMovie(c *gin.Context) {
	var newMovie models.Movie

	if err := c.BindJSON(&newMovie); err != nil {
		c.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := mc.validate.Struct(newMovie)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid field data", "details": err.Error()})
		return
	}

	_, err = mc.movieCollection.InsertOne(ctx, newMovie)

	if err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong", "details": err})
		return
	}
	
	c.JSON(201, gin.H{"message": "Movie added successfully"})	
}