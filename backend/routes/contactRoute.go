package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wetrade/controllers"
)

func ContactRoutes(r *gin.Engine) {
	apiV1 := r.Group("/api/v1")
	apiV1.POST("/contact", controllers.CreateContact())
	apiV1.GET("/contact/:id", controllers.GetContactDetailsById())
	apiV1.GET("/contacts", controllers.GetAllContacts())
}
