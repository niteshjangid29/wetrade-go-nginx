package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contact struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	FirstName  string             `json:"first_name"`
	LastName   string             `json:"last_name"`
	Email      string             `json:"email"`
	Phone      string             `json:"phone"`
	City       string             `json:"city"`
	State      string             `json:"state"`
	Investment int8               `json:"investment"`
	TradingExp int8               `json:"trading_exp"`
	Enroll     int8               `json:"enroll"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

// 1 - Below 1 Lac
// 2 - 1 Lac to 5 Lac
// 3 - 5 Lac to 10 Lac
// 4 - above 10 Lac

// 1 - Beginner
// 2 - Less than 1 year
// 3 - 1 to 3 years
// 4 - More than 3 years

// Want to Enroll in Upcoming Batch?
// 1 - Yes
// 2 - No
