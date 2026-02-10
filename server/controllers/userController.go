package controllers

import (
	"context"
	"github.com/ImranullahKhann/movie-streaming-app/server/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type UserController struct {
	userCollection *mongo.Collection
	validate       *validator.Validate
}

func NewUserController(collection *mongo.Collection) UserController {
	return UserController{userCollection: collection, validate: validator.New()}
}

func (uc *UserController) RegisterUser(c *gin.Context) {
	var newUser models.User

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	if err := uc.validate.Struct(newUser); err != nil {
		c.JSON(400, gin.H{"error": "Invalid field data", "details": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := uc.userCollection.InsertOne(ctx, newUser); err != nil {
		c.JSON(500, gin.H{"error": "Couldn't write to database", "details": err})
		return
	}

	c.JSON(201, gin.H{"message": "User created successfully"})
}
