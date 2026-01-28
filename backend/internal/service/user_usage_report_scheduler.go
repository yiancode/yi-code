package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// UserUsageReportScheduler handles scheduled usage report sending
type UserUsageReportScheduler struct {
	reportService  *UserUsageReportService
	settingService *SettingService
	reportRepo     UserUsageReportRepository
	redisClient    *redis.Client

	stopCh chan struct{}
	wg     sync.WaitGroup
}

// NewUserUsageReportScheduler creates a new UserUsageReportScheduler
func NewUserUsageReportScheduler(
	reportService *UserUsageReportService,
	settingService *SettingService,
	reportRepo UserUsageReportRepository,
	redisClient *redis.Client,
) *UserUsageReportScheduler {
	return &UserUsageReportScheduler{
		reportService:  reportService,
		settingService: settingService,
		reportRepo:     reportRepo,
		redisClient:    redisClient,
		stopCh:         make(chan struct{}),
	}
}

// Start starts the scheduler
func (s *UserUsageReportScheduler) Start() {
	s.wg.Add(1)
	go s.run()
	log.Printf("[UsageReportScheduler] Started")
}

// Stop stops the scheduler
func (s *UserUsageReportScheduler) Stop() {
	close(s.stopCh)
	s.wg.Wait()
	log.Printf("[UsageReportScheduler] Stopped")
}

func (s *UserUsageReportScheduler) run() {
	defer s.wg.Done()

	// Check every minute
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.checkAndSendReports()
		}
	}
}

func (s *UserUsageReportScheduler) checkAndSendReports() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Check if globally enabled
	config, err := s.settingService.GetUsageReportConfig(ctx)
	if err != nil {
		log.Printf("[UsageReportScheduler] Error getting config: %v", err)
		return
	}

	if !config.GlobalEnabled {
		return
	}

	now := time.Now()
	currentHourMinute := now.Format("15:04")

	// Get users based on target scope
	users, err := s.getUsersToSend(ctx, config, now)
	if err != nil {
		log.Printf("[UsageReportScheduler] Error getting users: %v", err)
		return
	}

	// Send reports to eligible users
	for _, user := range users {
		// Determine send time based on scope
		var shouldSend bool
		if config.TargetScope == UsageReportScopeOptedIn {
			// Use user's configured schedule and timezone
			userConfig, err := s.reportRepo.GetUserReportConfig(ctx, user.ID)
			if err != nil {
				log.Printf("[UsageReportScheduler] Error getting user config for %d: %v", user.ID, err)
				continue
			}
			if !userConfig.Enabled {
				continue
			}

			// Load user's timezone
			loc, err := time.LoadLocation(userConfig.Timezone)
			if err != nil {
				loc = time.FixedZone("UTC+8", 8*3600) // Default to UTC+8
			}
			userNow := now.In(loc)
			userHourMinute := userNow.Format("15:04")
			shouldSend = userHourMinute == userConfig.Schedule
		} else {
			// Use global schedule for all/active_today modes
			shouldSend = currentHourMinute == config.GlobalSchedule
		}

		if !shouldSend {
			continue
		}

		// Check if already sent today (using Redis SetNX to prevent duplicate sends)
		// SetNX returns true if key was set (first time), false if key already exists
		sentKey := s.getSentKey(user.ID, now)
		isFirstSend, err := s.redisClient.SetNX(ctx, sentKey, "1", 25*time.Hour).Result()
		if err != nil {
			log.Printf("[UsageReportScheduler] Error checking sent status for user %d: %v", user.ID, err)
			continue
		}
		if !isFirstSend {
			// Already sent today, skip
			continue
		}

		// Generate and send report
		go s.sendReportToUser(user.ID, now)
	}
}

func (s *UserUsageReportScheduler) getUsersToSend(ctx context.Context, config *UsageReportConfig, now time.Time) ([]User, error) {
	return s.reportRepo.GetUsersForUsageReport(ctx, config.TargetScope, now)
}

func (s *UserUsageReportScheduler) sendReportToUser(userID int64, reportTime time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Acquire distributed lock to prevent concurrent sends
	lockKey := fmt.Sprintf("usage_report_lock:%d:%s", userID, reportTime.Format("2006-01-02"))
	locked, err := s.redisClient.SetNX(ctx, lockKey, "1", 5*time.Minute).Result()
	if err != nil || !locked {
		return
	}
	defer s.redisClient.Del(ctx, lockKey)

	// Generate report for yesterday
	yesterday := reportTime.AddDate(0, 0, -1)

	data, err := s.reportService.GenerateReportData(ctx, userID, yesterday)
	if err != nil {
		log.Printf("[UsageReportScheduler] Error generating report for user %d: %v", userID, err)
		// Remove sent key to allow retry
		s.redisClient.Del(ctx, s.getSentKey(userID, reportTime))
		return
	}

	// Only send if there was any usage
	if data.TotalRequests == 0 {
		log.Printf("[UsageReportScheduler] Skipping report for user %d: no usage", userID)
		return
	}

	if err := s.reportService.SendReport(ctx, userID, data); err != nil {
		log.Printf("[UsageReportScheduler] Error sending report for user %d: %v", userID, err)
		// Remove sent key to allow retry
		s.redisClient.Del(ctx, s.getSentKey(userID, reportTime))
		return
	}

	log.Printf("[UsageReportScheduler] Sent usage report to user %d", userID)
}

func (s *UserUsageReportScheduler) getSentKey(userID int64, t time.Time) string {
	return fmt.Sprintf("usage_report_sent:%d:%s", userID, t.Format("2006-01-02"))
}
