package routes

import (
	"streetlight-controller/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		schedule := api.Group("/schedule")
		{
			schedule.POST("/create", handlers.CreateSchedule)
			schedule.GET("/list", handlers.GetSchedules)
			schedule.PUT("/:id", handlers.UpdateSchedule)
			schedule.DELETE("/:id", handlers.DeleteSchedule)
			schedule.POST("/:id/precheck", handlers.TriggerPreCheck)
		}
		weather := api.Group("/weather")
		{
			weather.GET("/current", handlers.GetWeather)
			weather.PUT("/update", handlers.UpdateWeather)
			weather.POST("/simulate/storm", handlers.SimulateStorm)
			weather.POST("/simulate/normal", handlers.SimulateNormal)
			weather.GET("/stream", handlers.GetWeatherStream)
		}
		api.GET("/status/stream", handlers.GetStatusStream)
	}

	return r
}
