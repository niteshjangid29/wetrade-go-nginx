package controllers

import (
	"context"
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

var contactCollection *mongo.Collection = database.ContactData(database.ContactCollection, config.LoadConfig())

func CreateContact() gin.HandlerFunc {
	return func(c *gin.Context) {
		var contact models.Contact
		if err := c.BindJSON(&contact); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		contact.ID = primitive.NewObjectID()
		contact.CreatedAt = time.Now()

		_, err := contactCollection.InsertOne(context.Background(), contact)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Contact created successfully"})
	}
}

func GetContactDetailsById() gin.HandlerFunc {
	return func(c *gin.Context) {
		contactId := c.Param("id")

		if contactId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Contact ID is required"})
			return
		}

		objID, err := primitive.ObjectIDFromHex(contactId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Contact ID"})
			return
		}

		var contact models.Contact

		err = contactCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&contact)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching contact details"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Contact details fetched successfully",
			"contact": contact,
		})
	}
}

func GetAllContacts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var contacts []models.Contact
		cursor, err := contactCollection.Find(context.Background(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var contact models.Contact
			cursor.Decode(&contact)
			contacts = append(contacts, contact)
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "Contacts fetched successfully",
			"contacts": contacts,
		})
	}
}
