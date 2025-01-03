package controllers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/wetrade/config"
	"github.com/wetrade/database"
	"github.com/wetrade/models"
	generate "github.com/wetrade/tokens"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.UserData(database.UserCollection, config.LoadConfig())
var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panicln("Error hashing password:", err)
	}
	return string(bytes)
}

func ComparePassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Email or Password is Incorrect"
		valid = false
	}
	return valid, msg
}

func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": validationErr.Error()})
			return
		}

		// Convert email to lowercase
		user.Email = strings.ToLower(user.Email)

		// Check if user with the same email already exists
		existingUser := models.User{}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)

		if err == nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "User with this email already exists"})
			return
		} else if err != mongo.ErrNoDocuments {
			log.Println("Error checking if user exists", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong. Please try again after sometime"})
			return
		}

		user.ID = primitive.NewObjectID()
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		password := HashPassword(user.Password)
		user.Password = password
		user.Role = "user" //default role is user

		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			log.Println("Error registering user into the database", insertErr)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong. Please try again after sometime"})
			return
		}

		defer cancel()
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "User registered successfully",
		})
	}
}

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": strings.ToLower(user.Email)}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Email or Password is incorrect"})
			return
		}

		isValid, msg := ComparePassword(user.Password, foundUser.Password)
		defer cancel()
		if !isValid {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": msg})
			return
		}

		token, _ := generate.GenerateToken(foundUser.FirstName, foundUser.Email, foundUser.ID.Hex(), foundUser.Role)
		defer cancel()

		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "User logged in successfully",
			"token":   token,
			"role":    foundUser.Role,
		})
	}
}

func GetUserDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("id")

		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User Id is required"})
			return
		}

		objID, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid User Id"})
			return
		}

		var user models.User

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching user details"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User details fetched successfully",
			"user":    user,
		})
	}
}
