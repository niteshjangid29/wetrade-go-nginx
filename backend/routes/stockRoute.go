package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wetrade/controllers"
	"github.com/wetrade/middleware"
)

func StockRoutes(r *gin.Engine) {
	apiV1 := r.Group("/api/v1")
	apiV1.POST("/create-stock", controllers.CreateStock())
	apiV1.POST("/create-multiple-stocks", controllers.CreateMultipleStocks())
	apiV1.GET("/stocks", middleware.Authenticate("owner"), controllers.GetAllStocks())
}
