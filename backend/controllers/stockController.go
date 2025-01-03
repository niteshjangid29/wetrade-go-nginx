package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wetrade/config"
	"github.com/wetrade/database"
	"github.com/wetrade/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var stockCollection *mongo.Collection = database.StocksData(database.StockCollection, config.LoadConfig())

func CreateStock() gin.HandlerFunc {
	return func(c *gin.Context) {
		var stock models.Stock
		if err := c.BindJSON(&stock); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		stock.ID = primitive.NewObjectID()
		stock.CreatedAt = time.Now()

		_, err := stockCollection.InsertOne(context.Background(), stock)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Stock created successfully"})
	}
}

func CreateMultipleStocks() gin.HandlerFunc {
	return func(c *gin.Context) {
		var stocks []models.Stock
		if err := c.BindJSON(&stocks); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		interfaceStocks := make([]interface{}, len(stocks))
		for i, stock := range stocks {
			stock.ID = primitive.NewObjectID()
			stock.CreatedAt = time.Now()
			interfaceStocks[i] = stock
		}

		log.Println(interfaceStocks)

		_, err := stockCollection.InsertMany(context.Background(), interfaceStocks)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Stocks created successfully"})
	}
}

func GetAllStocks() gin.HandlerFunc {
	return func(c *gin.Context) {
		var stocks []models.Stock
		cursor, err := stockCollection.Find(context.Background(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		defer cursor.Close(context.Background())

		if err = cursor.All(context.Background(), &stocks); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, stocks)
	}
}
