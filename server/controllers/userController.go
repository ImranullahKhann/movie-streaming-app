package controllers

import (
	"context"
	"errors"
	"github.com/ImranullahKhann/movie-streaming-app/server/models"
	"github.com/ImranullahKhann/movie-streaming-app/server/store"
	"github.com/ImranullahKhann/movie-streaming-app/server/utils"
	"github.com/ImranullahKhann/movie-streaming-app/server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
	"time"
)

type UserController struct {
	userCollection *mongo.Collection
	validate       *validator.Validate
	rds            *store.Redis
}

func NewUserController(collection *mongo.Collection, redisClient *store.Redis) UserController {
	return UserController{userCollection: collection, validate: validator.New(), rds: redisClient}
}

func (uc *UserController) RegisterUser(c *gin.Context) {
	var newUser models.User

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()
	hash, err := utils.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}
	newUser.Password = hash

	if err = uc.validate.Struct(newUser); err != nil {
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

func (uc *UserController) LoginUser(c *gin.Context) {
	var loginInfo models.UserLogin

	if err := c.BindJSON(&loginInfo); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := uc.validate.Struct(loginInfo); err != nil {
		c.JSON(400, gin.H{"error": "Invalid field data", "details": err.Error()})
		return
	}

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := uc.userCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: loginInfo.Email}},
	).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(404, gin.H{"error": "No such user"})
			return
		}
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}

	if !utils.VerifyPassword(loginInfo.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	toks, err := utils.IssueTokens(loginInfo.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue tokens"})
		return
	}
	if err := utils.Persist(c, uc.rds, toks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not persist tokens"})
		return
	}
	utils.SetAuthCookies(c, toks)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (uc *UserController) LogoutUser(c *gin.Context) {
	acc, _ := c.Cookie("access_token")
	ref, _ := c.Cookie("refresh_token")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if acc != "" {
		if claims, err := utils.ParseAccess(acc); err != nil {
			_ = uc.rds.DelJTI(ctx, "access:"+claims.ID)
		}
	}
	if ref != "" {
		if claims, err := utils.ParseRefresh(ref); err != nil {
			_ = uc.rds.DelJTI(ctx, "refresh:"+claims.ID)
		}
	}
	utils.ClearAuthCookies(c)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (uc *UserController) RefreshTokens(c *gin.Context) {
	ref, err := middleware.MustCookie(c, "access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}
	claims, err := utils.ParseRefresh(ref)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := uc.rds.GetUserByJTI(ctx, "refresh:"+claims.ID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh revoked"})
		return
	}

	_ = uc.rds.DelJTI(ctx, "refresh:"+claims.ID)

	toks, err := utils.IssueTokens(claims.Subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue new tokens"})
		return
	}
	if err := utils.Persist(ctx, uc.rds, toks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not persist new tokens"})
		return
	}

	utils.SetAuthCookies(c, toks)
	c.JSON(http.StatusCreated, gin.H{"ok": true})
}