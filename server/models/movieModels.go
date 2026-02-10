package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Converting a Go type to BSON is called marshalling, while the reverse process is called unmarshalling

type Movie struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	ImdbID      string        `bson:"imdb_id" json:"imdb_id" validate:"required"`
	Title       string        `bson:"title" json:"title" validate:"required,min=2"`
	PosterPath  string        `bson:"poster_path" json:"poster_path" validate:"required,url"`
	YoutubeId   string        `bson:"youtube_id" json:"youtube_id" validate:"required"`
	Genres      []Genre       `json:"genre" validate"dive,required"`
	AdminReview string        `bson:"admin_review" json:"admin_review" validate"max=128"`
	Ranking     Ranking       `json:"ranking"`
}

type Genre struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	GenreID   int           `bson:"genre_id" json"genre_id" validate:"required"`
	GenreName string        `bson:"genre_name" json:"genre_name" validate:"required,max=64"`
}

type Ranking struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	RankingValue int           `bson:"ranking_value" json:"ranking_value" validate="required"`
	RankingName  string        `bson:"ranking_name" json:"ranking_name" validate="required,max=24"`
}
