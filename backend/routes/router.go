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

		voyage := api.Group("/travels")
		voyage.GET("", controllers.GetTravels)
		voyage.GET("/:id", controllers.GetVoyage)
		voyage.Use(middlewares.AdminMiddleware())
		{
			voyage.POST("", controllers.CreateVoyage)
			voyage.PUT("/:id", controllers.UpdateVoyage)
			voyage.DELETE("/:id", controllers.DeleteVoyage)
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
