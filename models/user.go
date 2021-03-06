package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	FirstName    string             `json:"Firstname" binding:"required" bson:"Firstname,omitempty"`
	LastName     string             `json:"Lastname" binding:"required" bson:"Lastname"`
	UserName     string             `bson:"Username" json:"Username" binding:"required"`
	Email        string             `bson:"Email" json:"Email" binding:"required"`
	Password     string             `bson:"Password" json:"Password" binding:"required"`
	CoursesTaken []string           `json:"Courses_taken,omitempty" bson:"Courses_taken,omitempty"`
	Roles        []Roles            `bson:"Roles,omitempty"`
	CreatedAt    time.Time          `json:"Created_at" bson:"Created_at"`
	UpdatedAt    time.Time          `json:"Updated_at,omitempty" bson:"Updated_at,omitempty"`
}

type Roles struct {
	Role string `json:"role" bson:"role,omitempty"`
	Db   string `json:"db" bson:"db,omitempty"`
}
