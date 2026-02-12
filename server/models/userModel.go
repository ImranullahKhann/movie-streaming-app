package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type User struct {
	ID              bson.ObjectID `bson:"_id,omitempty"`
	FirstName       string        `bson:"first_name" json:"first_name" validate:"required"`
	LastName        string        `bson:"last_name" json:"last_name" validate:"required"`
	Email           string        `bson:"email" json:"email" validate:"required,email"`
	Password        string        `bson:"password" json:"password" validate:"required,min=8"`
	Role            string        `bson:"role" json:"role" validate:"required"`
	CreatedAt       time.Time     `bson:"created_at" json:"created_at" validate:"required"`
	UpdatedAt       time.Time     `bson:"updated_at" json:"updated_at" validate:"required"`
	FavouriteGenres []Genre       `bson:"favourite_genres" json:"favourite_genres" validate:"dive,required"`
}

// Data Transfer Object
type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
type UserResponse struct {
	UserId          string  `json:"user_id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Email           string  `json:"email"`
	Role            string  `json:"role"`
	Token           string  `json:"token"`
	RefreshToken    string  `json:"refresh_token"`
	FavouriteGenres []Genre `json:"favourite_genres"`
}
