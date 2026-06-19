package models

import "time"

type Schedule struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	GroupName string    `json:"group_name" gorm:"not null;uniqueIndex"`
	OnTime    string    `json:"on_time" gorm:"not null"`
	OffTime   string    `json:"off_time" gorm:"not null"`
	Status    string    `json:"status" gorm:"default:active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateScheduleRequest struct {
	GroupName string `json:"group_name" binding:"required"`
	OnTime    string `json:"on_time" binding:"required"`
	OffTime   string `json:"off_time" binding:"required"`
}
