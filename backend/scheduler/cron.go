package scheduler

import (
	"fmt"
	"log"
	"streetlight-controller/models"
	"strings"
	"sync"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron     *cron.Cron
	jobs     map[uint][]cron.EntryID
	mu       sync.Mutex
	StatusCh chan string
}

var GlobalScheduler *Scheduler

func Init() {
	GlobalScheduler = &Scheduler{
		cron:     cron.New(cron.WithSeconds()),
		jobs:     make(map[uint][]cron.EntryID),
		StatusCh: make(chan string, 100),
	}
	GlobalScheduler.cron.Start()
	log.Println("Scheduler initialized successfully")
}

func (s *Scheduler) AddSchedule(schedule *models.Schedule) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[schedule.ID]; exists {
		s.RemoveSchedule(schedule.ID)
	}

	onHour, onMin, _ := parseTime(schedule.OnTime)
	onSpec := fmt.Sprintf("0 %d %d * * ?", onMin, onHour)
	onID, err := s.cron.AddFunc(onSpec, func() {
		msg := fmt.Sprintf("[%s] 群组 [%s] 路灯已开启", schedule.OnTime, schedule.GroupName)
		log.Println(msg)
		select {
		case s.StatusCh <- msg:
		default:
		}
	})
	if err != nil {
		return fmt.Errorf("failed to add on schedule: %w", err)
	}

	offHour, offMin, _ := parseTime(schedule.OffTime)
	offSpec := fmt.Sprintf("0 %d %d * * ?", offMin, offHour)
	offID, err := s.cron.AddFunc(offSpec, func() {
		msg := fmt.Sprintf("[%s] 群组 [%s] 路灯已关闭", schedule.OffTime, schedule.GroupName)
		log.Println(msg)
		select {
		case s.StatusCh <- msg:
		default:
		}
	})
	if err != nil {
		s.cron.Remove(onID)
		return fmt.Errorf("failed to add off schedule: %w", err)
	}

	s.jobs[schedule.ID] = []cron.EntryID{onID, offID}
	log.Printf("Schedule added for group %s: ON=%s, OFF=%s", schedule.GroupName, schedule.OnTime, schedule.OffTime)
	return nil
}

func (s *Scheduler) RemoveSchedule(id uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entryIDs, exists := s.jobs[id]; exists {
		for _, entryID := range entryIDs {
			s.cron.Remove(entryID)
		}
		delete(s.jobs, id)
		log.Printf("Schedule %d removed", id)
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
