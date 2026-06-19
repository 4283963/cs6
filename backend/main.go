package main

import (
	"log"
	"streetlight-controller/database"
	"streetlight-controller/models"
	"streetlight-controller/routes"
	"streetlight-controller/scheduler"
)

func main() {
	database.Init()
	scheduler.Init()

	loadExistingSchedules()

	r := routes.SetupRouter()

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func loadExistingSchedules() {
	var schedules []models.Schedule
	result := database.DB.Find(&schedules)
	if result.Error != nil {
		log.Printf("Warning: Failed to load existing schedules: %v", result.Error)
		return
	}

	for i := range schedules {
		err := scheduler.GlobalScheduler.AddSchedule(&schedules[i])
		if err != nil {
			log.Printf("Warning: Failed to load schedule %d: %v", schedules[i].ID, err)
		}
	}

	log.Printf("Loaded %d existing schedules", len(schedules))
}
