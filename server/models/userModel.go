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
	Password        string        `bson:"password" json:"password" validate:"required"`
	Role            string        `bson:"role" json:"role" validate:"required"`
	CreatedAt       time.Time     `bson:"created_at" json:"created_at" validate:"required"`
	UpdatedAt       time.Time     `bson:"updated_at" json:"updated_at" validate:"required"`
	FavouriteGenres []Genre       `bson:"favourite_genres" json:"favourite_genres" validate:"dive,required"`
}
