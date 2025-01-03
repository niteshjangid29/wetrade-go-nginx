package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wetrade/routes"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())

	routes.StockRoutes(r)
	routes.UserRoutes(r)
	routes.ContactRoutes(r)

	// print Hello world route for testing
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	r.POST("/place-order", func(c *gin.Context) {
		var requestBody struct {
			APIKey       string `json:"api_key"`
			APISecret    string `json:"api_secret"`
			RequestToken string `json:"request_token"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		kc := kiteconnect.New(requestBody.APIKey)

		data, err := kc.GenerateSession(requestBody.RequestToken, requestBody.APISecret)
		if err != nil {
			log.Println("Error generating session: ", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		kc.SetAccessToken(data.AccessToken)

		log.Println("Access Token: ", data.AccessToken)

		log.Println("Fetching GTTs...")
		orders, err := kc.GetGTTs()
		if err != nil {
			log.Fatalf("Error getting GTTs: %v", err)
		}
		log.Printf("gtt: %v", orders)

		log.Println("Placing GTT...")
		// Place GTT

		ordersArr := []struct {
			Tradingsymbol string
			LastPrice     float64
			TriggerValue  float64
			Quantity      float64
			LimitPrice    float64
		}{
			{
				Tradingsymbol: "NTPC",
				LastPrice:     370,
				TriggerValue:  371,
				Quantity:      3,
				LimitPrice:    371,
			},
			{
				Tradingsymbol: "SBIN",
				LastPrice:     816,
				TriggerValue:  820,
				Quantity:      2,
				LimitPrice:    820,
			},
			// {
			// 	Tradingsymbol: "INFY",
			// 	LastPrice:     1902,
			// 	TriggerValue:  1910,
			// 	Quantity:      1,
			// 	LimitPrice:    1910,
			// },
			// {
			// 	Tradingsymbol: "GSPL",
			// 	LastPrice:     332,
			// 	TriggerValue:  335,
			// 	Quantity:      1,
			// 	LimitPrice:    335,
			// },
			// {
			// 	Tradingsymbol: "GAIL",
			// 	LastPrice:     192.5,
			// 	TriggerValue:  195,
			// 	Quantity:      2,
			// 	LimitPrice:    195,
			// },
			// {
			// 	Tradingsymbol: "LT",
			// 	LastPrice:     3603,
			// 	TriggerValue:  3620,
			// 	Quantity:      2,
			// 	LimitPrice:    3620,
			// },
			// {
			// 	Tradingsymbol: "NTPC",
			// 	LastPrice:     370,
			// 	TriggerValue:  371,
			// 	Quantity:      3,
			// 	LimitPrice:    371,
			// },
			// {
			// 	Tradingsymbol: "SBIN",
			// 	LastPrice:     816,
			// 	TriggerValue:  820,
			// 	Quantity:      2,
			// 	LimitPrice:    820,
			// },
			// {
			// 	Tradingsymbol: "INFY",
			// 	LastPrice:     1902,
			// 	TriggerValue:  1910,
			// 	Quantity:      1,
			// 	LimitPrice:    1910,
			// },
			// {
			// 	Tradingsymbol: "GSPL",
			// 	LastPrice:     332,
			// 	TriggerValue:  335,
			// 	Quantity:      1,
			// 	LimitPrice:    335,
			// },
			// {
			// 	Tradingsymbol: "GAIL",
			// 	LastPrice:     192.5,
			// 	TriggerValue:  195,
			// 	Quantity:      2,
			// 	LimitPrice:    195,
			// },
			// {
			// 	Tradingsymbol: "LT",
			// 	LastPrice:     3603,
			// 	TriggerValue:  3620,
			// 	Quantity:      2,
			// 	LimitPrice:    3620,
			// },
		}

		for _, order := range ordersArr {
			gttResp, err := kc.PlaceGTT(kiteconnect.GTTParams{
				Tradingsymbol:   order.Tradingsymbol,
				Exchange:        "NSE",
				LastPrice:       order.LastPrice,
				TransactionType: kiteconnect.TransactionTypeBuy,
				Trigger: &kiteconnect.GTTSingleLegTrigger{
					TriggerParams: kiteconnect.TriggerParams{
						TriggerValue: order.TriggerValue,
						Quantity:     order.Quantity,
						LimitPrice:   order.LimitPrice,
					},
				},
			})
			if err != nil {
				log.Fatalf("error placing gtt: %v", err)
			}

			log.Println("placed GTT trigger_id = ", gttResp.TriggerID)
		}
	})

	// Get User Margins (this is your Zerodha API call)
	// r.GET("/margins", func(c *gin.Context) {
	// 	margins, err := kc.GetUserMargins()
	// 	if err != nil {
	// 		c.JSON(500, gin.H{
	// 			"error": fmt.Sprintf("Error getting margins: %v", err),
	// 		})
	// 		return
	// 	}

	// 	// Respond with the margin data
	// 	c.JSON(200, gin.H{
	// 		"margins": margins,
	// 	})
	// })

	// Example route to get user details (for testing)
	// r.GET("/user", func(c *gin.Context) {
	// 	user, err := kc.GetUserProfile()
	// 	if err != nil {
	// 		c.JSON(500, gin.H{
	// 			"error": fmt.Sprintf("Error getting user profile: %v", err),
	// 		})
	// 		return
	// 	}

	// 	c.JSON(200, gin.H{
	// 		"user": user,
	// 	})
	// })

	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "Hello, World!",
	// 	})
	// })

	r.Run(":8000")
}
