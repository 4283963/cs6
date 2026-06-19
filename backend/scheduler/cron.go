package scheduler

import (
	"fmt"
	"log"
	"streetlight-controller/database"
	"streetlight-controller/models"
	"streetlight-controller/weather"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron              *cron.Cron
	jobs              map[uint][]cron.EntryID
	advancedToday     map[uint]bool
	advancedTodayDate string
	mu                sync.Mutex
	StatusCh          chan string
}

var GlobalScheduler *Scheduler

func Init() {
	GlobalScheduler = &Scheduler{
		cron:              cron.New(cron.WithSeconds()),
		jobs:              make(map[uint][]cron.EntryID),
		advancedToday:     make(map[uint]bool),
		advancedTodayDate: time.Now().Format("2006-01-02"),
		StatusCh:          make(chan string, 100),
	}
	GlobalScheduler.cron.Start()
	go GlobalScheduler.startDailyReset()
	log.Println("Scheduler initialized successfully")
}

func (s *Scheduler) startDailyReset() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		today := time.Now().Format("2006-01-02")
		if today != s.advancedTodayDate {
			s.advancedToday = make(map[uint]bool)
			s.advancedTodayDate = today
			log.Printf("🔄 Daily reset: cleared advanced-on records for new day %s", today)
		}
		s.mu.Unlock()
	}
}

func (s *Scheduler) checkAndResetDaily() {
	s.mu.Lock()
	defer s.mu.Unlock()
	today := time.Now().Format("2006-01-02")
	if today != s.advancedTodayDate {
		s.advancedToday = make(map[uint]bool)
		s.advancedTodayDate = today
	}
}

func (s *Scheduler) triggerOn(schedule *models.Schedule, isAdvanced bool) {
	s.checkAndResetDaily()

	s.mu.Lock()
	if isAdvanced {
		if s.advancedToday[schedule.ID] {
			s.mu.Unlock()
			log.Printf("ℹ️  Schedule %d already turned on advanced today, skipping", schedule.ID)
			return
		}
		s.advancedToday[schedule.ID] = true
	} else {
		if s.advancedToday[schedule.ID] {
			s.mu.Unlock()
			log.Printf("ℹ️  Schedule %d already turned on advanced today, skipping regular on-time", schedule.ID)
			return
		}
	}
	s.mu.Unlock()

	var msg string
	if isAdvanced {
		msg = fmt.Sprintf("[%s] ⚡ 暴雨低光照，群组 [%s] 路灯提前半小时开启",
			time.Now().Format("15:04"), schedule.GroupName)
	} else {
		msg = fmt.Sprintf("[%s] 群组 [%s] 路灯已开启", schedule.OnTime, schedule.GroupName)
	}
	log.Println(msg)
	select {
	case s.StatusCh <- msg:
	default:
	}
}

func (s *Scheduler) triggerOff(schedule *models.Schedule) {
	msg := fmt.Sprintf("[%s] 群组 [%s] 路灯已关闭", schedule.OffTime, schedule.GroupName)
	log.Println(msg)
	select {
	case s.StatusCh <- msg:
	default:
	}
}

func (s *Scheduler) PreCheck(schedule *models.Schedule) {
	if weather.GlobalWeatherService.IsLowLight() {
		weatherStatus := weather.GlobalWeatherService.GetStatus()
		log.Printf("🌧️  Pre-check for %s: low light detected (%.1f lux, %s), triggering advanced on",
			schedule.GroupName, weatherStatus.Illuminance, weatherStatus.Condition)
		s.triggerOn(schedule, true)
	} else {
		weatherStatus := weather.GlobalWeatherService.GetStatus()
		log.Printf("☀️  Pre-check for %s: normal light (%.1f lux), will wait for regular on-time %s",
			schedule.GroupName, weatherStatus.Illuminance, schedule.OnTime)
	}
}

func (s *Scheduler) TriggerPreCheckByID(id uint) error {
	s.mu.Lock()
	_, exists := s.jobs[id]
	s.mu.Unlock()

	if !exists {
		return fmt.Errorf("schedule %d not found", id)
	}

	var schedule models.Schedule
	result := database.DB.First(&schedule, id)
	if result.Error != nil {
		return result.Error
	}

	s.PreCheck(&schedule)
	return nil
}

func (s *Scheduler) AddSchedule(schedule *models.Schedule) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entryIDs, exists := s.jobs[schedule.ID]; exists {
		for _, entryID := range entryIDs {
			s.cron.Remove(entryID)
		}
		delete(s.jobs, schedule.ID)
		delete(s.advancedToday, schedule.ID)
		log.Printf("Old schedule %d removed before adding new one", schedule.ID)
	}

	onHour, onMin, _ := parseTime(schedule.OnTime)
	preCheckHour, preCheckMin := calculateAdvanceTime(onHour, onMin, weather.AdvanceMinutes)

	preCheckSpec := fmt.Sprintf("0 %d %d * * ?", preCheckMin, preCheckHour)
	preCheckID, err := s.cron.AddFunc(preCheckSpec, func() {
		log.Printf("🔍 Running pre-check for schedule %d (%s), 30min before on-time",
			schedule.ID, schedule.GroupName)
		s.PreCheck(schedule)
	})
	if err != nil {
		return fmt.Errorf("failed to add pre-check schedule: %w", err)
	}

	onSpec := fmt.Sprintf("0 %d %d * * ?", onMin, onHour)
	onID, err := s.cron.AddFunc(onSpec, func() {
		s.triggerOn(schedule, false)
	})
	if err != nil {
		s.cron.Remove(preCheckID)
		return fmt.Errorf("failed to add on schedule: %w", err)
	}

	offHour, offMin, _ := parseTime(schedule.OffTime)
	offSpec := fmt.Sprintf("0 %d %d * * ?", offMin, offHour)
	offID, err := s.cron.AddFunc(offSpec, func() {
		s.triggerOff(schedule)
	})
	if err != nil {
		s.cron.Remove(preCheckID)
		s.cron.Remove(onID)
		return fmt.Errorf("failed to add off schedule: %w", err)
	}

	s.jobs[schedule.ID] = []cron.EntryID{preCheckID, onID, offID}
	log.Printf("✅ Schedule added/updated - ID: %d, Group: %s", schedule.ID, schedule.GroupName)
	log.Printf("   Pre-check: %02d:%02d (cron: %s)", preCheckHour, preCheckMin, preCheckSpec)
	log.Printf("   On-time:   %s (cron: %s)", schedule.OnTime, onSpec)
	log.Printf("   Off-time:  %s (cron: %s)", schedule.OffTime, offSpec)
	return nil
}

func calculateAdvanceTime(hour, min, advanceMin int) (int, int) {
	totalMin := hour*60 + min - advanceMin
	if totalMin < 0 {
		totalMin += 24 * 60
	}
	return totalMin / 60, totalMin % 60
}

func (s *Scheduler) RemoveSchedule(id uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entryIDs, exists := s.jobs[id]; exists {
		for _, entryID := range entryIDs {
			s.cron.Remove(entryID)
			log.Printf("  - Cron entry %d removed for schedule %d", entryID, id)
		}
		delete(s.jobs, id)
		log.Printf("✅ Schedule %d completely removed from memory", id)
	} else {
		log.Printf("⚠️  Schedule %d not found in memory (may have been removed already)", id)
	}
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

func parseTime(timeStr string) (hour, min int, err error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid time format")
	}
	fmt.Sscanf(parts[0], "%d", &hour)
	fmt.Sscanf(parts[1], "%d", &min)
	return hour, min, nil
}
