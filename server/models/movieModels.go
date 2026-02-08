package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Converting a Go type to BSON is called marshalling, while the reverse process is called unmarshalling

type Movie struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ImdbID      string             `bson:"imdb_id" json:"imdbID"`
	Title       string             `bson:"title" json:"title" validate:"required,min=2"`
	PosterPath  string             `bson:"poster_path" json:"posterPath" validate:"required,url"`
	YoutubeId   string             `bson:"youtube_id" json:"youtubeId" validate:"required,url"`
	genre       Genre              `json:"genre" validate"dive,required"`
	AdminReview string             `bson:"admin_review" json:"adminReview" validate"max=128"`
	Ranking     Ranking            `json:"ranking" validate:"dive"`
}

type Genre struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GenreID   int                `bson:"genre_id" json"genreId" validate:"required"`
	GenreName string             `bson:"genre_name" json:"genreName" validate:"required,max=64"`
}

type Ranking struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RankingValue int                `bson:"ranking_value" json:"rankingValue" validate="required"`
	RankingName  string             `bson:"ranking_name" json:"rankingName" validate="required,max=24"`
}
