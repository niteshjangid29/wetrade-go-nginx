package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Email     string             `json:"email"`
	Phone     string             `json:"phone"`
	Password  string             `json:"password"`
	Role      string             `json:"role" default:"user"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
