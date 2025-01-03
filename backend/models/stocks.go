package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Stock struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Symbol          string             `json:"symbol"`
	LastPrice       float64            `json:"last_price"`
	TriggerPrice    float64            `json:"trigger_price"`
	LimitPrice      float64            `json:"limit_price"`
	Quantity        int                `json:"quantity"`
	Exchange        string             `json:"exchange"`
	TransactionType string             `json:"transaction_type"`
	Date            time.Time          `json:"date"`
	CreatedBy       string             `json:"created_by"`
	IsVerified      bool               `json:"is_verified" default:"false"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}
