package controllers

import (
	"context"
	"errors"
	"github.com/ImranullahKhann/movie-streaming-app/server/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"net/http"
	"time"
)

type MovieController struct {
	movieCollection *mongo.Collection
	userCollection  *mongo.Collection
	validate        *validator.Validate
}

func NewMovieController(movieCollection *mongo.Collection, userCollection *mongo.Collection) *MovieController {
	return &MovieController{
		movieCollection: movieCollection,
		userCollection:  userCollection,
		validate:        validator.New(),
	}
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
		c.JSON(500, gin.H{"error": "Couldn't write to database", "details": err})
		return
	}

	c.JSON(201, gin.H{"message": "Movie added successfully"})
}

func (mc *MovieController) GetRecommendedMovies(c *gin.Context) {
	userEmail, _ := c.Get("userEmail")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	var user models.User
	err := mc.userCollection.FindOne(ctx, bson.D{{Key: "email", Value: userEmail}}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		return
	}

	var genreNames []string
	for _, genre := range user.FavouriteGenres {
		genreNames = append(genreNames, genre.GenreName)
	}

	
	opts := options.Find().
		SetSort(bson.D{{Key: "rating", Value: -1}}).
		SetLimit(5)
	
	filter := bson.D{{
		Key: "genre.genre_name",
		Value: bson.D{{Key: "$in", Value: genreNames}},
	}}
	
	cursor, err := mc.movieCollection.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed", "details": err.Error()})
		return
	}

	var recommendedMovies []models.Movie
	if err = cursor.All(ctx, &recommendedMovies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't read data", "details": err.Error()})
		return
	}

	if recommendedMovies == nil {
		recommendedMovies = []models.Movie{}
	}

	c.JSON(http.StatusOK, gin.H{"recommendedMovies": recommendedMovies})
}