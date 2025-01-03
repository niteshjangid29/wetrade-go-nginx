package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ZerodhaUser struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserID    primitive.ObjectID `json:"user_id"`
	ZerodhaID string             `json:"zerodha_id"`
	ApiKey    string             `json:"api_key"`
	ApiSecret string             `json:"api_secret"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
