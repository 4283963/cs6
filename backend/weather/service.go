package weather

import (
	"log"
	"sync"
	"time"
)

const (
	LowLightThreshold = 20.0
	AdvanceMinutes    = 30
)

type WeatherStatus struct {
	Illuminance float64   `json:"illuminance"`
	Condition   string    `json:"condition"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WeatherService struct {
	mu       sync.RWMutex
	status   WeatherStatus
	trigger  chan struct{}
	listener chan WeatherStatus
}

var GlobalWeatherService *WeatherService

func Init() {
	GlobalWeatherService = &WeatherService{
		status: WeatherStatus{
			Illuminance: 1000,
			Condition:   "normal",
			UpdatedAt:   time.Now(),
		},
		listener: make(chan WeatherStatus, 100),
	}
	log.Println("Weather service initialized (default: 1000 lux, normal)")
}

func (w *WeatherService) GetStatus() WeatherStatus {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.status
}

func (w *WeatherService) UpdateStatus(illuminance float64, condition string) WeatherStatus {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.status.Illuminance = illuminance
	w.status.Condition = condition
	w.status.UpdatedAt = time.Now()

	log.Printf("🌤️  Weather updated: %.1f lux, condition: %s", illuminance, condition)

	select {
	case w.listener <- w.status:
	default:
	}

	return w.status
}

func (w *WeatherService) IsLowLight() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.status.Illuminance < LowLightThreshold
}

func (w *WeatherService) GetListener() <-chan WeatherStatus {
	return w.listener
}

func (w *WeatherService) SimulateStorm() {
	w.UpdateStatus(15, "暴雨")
}

func (w *WeatherService) SimulateNormal() {
	w.UpdateStatus(1000, "normal")
}
