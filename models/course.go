package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Name             string             `json:"name" binding:"required" bson:"Name,omitempty"`
	Author           string             `json:"author" binding:"required" bson:"Author"`
	Description      string             `json:"description" binding:"required" bson:"Description,omitempty"`
	CreatedAt        time.Time          `json:"created_at" bson:"Created_at"`
	UpdatedAt        time.Time          `json:"updated_at,omitempty" bson:"Updated_at,omitempty"`
	Section          []string           `json:"Section,omitempty" bson:"Section,omitempty"`
	StudentsEnrolled []string           `json:"StudentsEnrolled" bson:"StudentsEnrolled,omitempty"`
}

type Section struct {
	ID      primitive.ObjectID `bson:"_id" `
	Title   string             `json:"title" bson:"Title,omitempty"`
	Content []string           `json:"content" binding:"required" bson:"Content,omitempty"`
}

type Content struct {
	ID        primitive.ObjectID `bson:"_id" `
	Title     string             `json:"title" binding:"required" bson:"Title,omitempty"`
	Video     string             `json:"video" bson:"Content,omitempty"`
	Thumbnail string             `json:"thumbnail" bson:"Thumbnail,omitempty"`
}
