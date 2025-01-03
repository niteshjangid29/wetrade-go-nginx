package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wetrade/controllers"
)

func UserRoutes(r *gin.Engine) {
	apiV1 := r.Group("/api/v1")
	apiV1.POST("/register", controllers.RegisterUser())
	apiV1.POST("/login", controllers.LoginUser())
	apiV1.GET("/user/:id", controllers.GetUserDetails())
}
