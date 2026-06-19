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
			schedule.DELETE("/:id", handlers.DeleteSchedule)
		}
		api.GET("/status/stream", handlers.GetStatusStream)
	}

	return r
}
