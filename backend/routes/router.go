package routes

import (
	"h3-travel/controllers"
	middlewares "h3-travel/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.POST("/signup", controllers.SignUp)
		api.POST("/login", controllers.Login)

		travel := api.Group("/travels")
		travel.GET("", controllers.GetTravels)
		travel.GET("/:id", controllers.GetTravel)
		travel.Use(middlewares.AdminMiddleware())
		{
			travel.POST("", controllers.CreateTravel)
			travel.PUT("/:id", controllers.UpdateTravel)
			travel.DELETE("/:id", controllers.DeleteTravel)
		}

		orders := api.Group("/orders")
		orders.Use(middlewares.JWTMiddleware())
		{
			orders.POST("", controllers.CreateOrder)
			orders.GET("/user", controllers.GetUserOrders)
			orders.PUT("/:id/cancel", controllers.CancelOrder)
		}
	}

	return r
}
