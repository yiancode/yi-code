package service

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

var (
	ErrUsageReportEmailNotBound = infraerrors.BadRequest("EMAIL_NOT_BOUND", "please bind your email first to enable usage reports")
	ErrUsageReportNotEnabled    = infraerrors.BadRequest("USAGE_REPORT_NOT_ENABLED", "usage report is not enabled globally")
)

// Suggestion thresholds and configuration
const (
	CacheHitRateLowThreshold  = 50.0  // Below this is considered low cache hit rate
	CacheHitRateHighThreshold = 90.0  // Above this is considered excellent cache hit rate
	HighCostThreshold         = 100.0 // Daily cost above this is considered high
	CacheSavingRate           = 0.9   // Cache typically saves ~90% of token cost
)

// toFloat64 converts various types to float64 for template calculations
func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int64:
		return float64(val)
	case int:
		return float64(val)
	case string:
		var f float64
		if _, err := fmt.Sscanf(val, "%f", &f); err != nil {
			// Log conversion error for debugging
			// Note: Using fmt.Printf for now. Replace with proper logger in production.
			fmt.Printf("WARN: Failed to convert string '%s' to float64: %v\n", val, err)
			return 0
		}
		return f
	default:
		fmt.Printf("WARN: Unsupported type for float64 conversion: %T\n", v)
		return 0
	}
}

// usageReportTemplate is the pre-parsed email template (parsed once at package init)
var usageReportTemplate = template.Must(template.New("report").Funcs(template.FuncMap{
	"mulf": func(a, b interface{}) float64 {
		return toFloat64(a) * toFloat64(b)
	},
	"divf": func(a, b interface{}) float64 {
		fa := toFloat64(a)
		fb := toFloat64(b)
		if fb == 0 {
			// Log division by zero for debugging
			fmt.Printf("WARN: Division by zero detected in template: %v / %v\n", a, b)
			return 0
		}
		return fa / fb
	},
}).Parse(usageReportEmailTemplate))

// UserUsageReportConfig represents user's usage report configuration
type UserUsageReportConfig struct {
	Enabled  bool   `json:"enabled"`
	Schedule string `json:"schedule"` // HH:MM format
	Timezone string `json:"timezone"`
}

// HourlyUsageStat represents usage statistics for a specific hour
type HourlyUsageStat struct {
	Hour     int   `json:"hour"`
	Requests int64 `json:"requests"`
}

// ModelUsageStat represents usage statistics for a specific model
type ModelUsageStat struct {
	Model    string  `json:"model"`
	Requests int64   `json:"requests"`
	Cost     float64 `json:"cost"`
	Percent  float64 `json:"percent"`
}

// UserUsageReportData represents the data for a usage report
type UserUsageReportData struct {
	TotalRequests      int64             `json:"total_requests"`
	TotalInputTokens   int64             `json:"total_input_tokens"`
	TotalOutputTokens  int64             `json:"total_output_tokens"`
	TotalCacheTokens   int64             `json:"total_cache_tokens"`
	CacheHitRate       float64           `json:"cache_hit_rate"`
	TotalCost          float64           `json:"total_cost"`
	ActualCost         float64           `json:"actual_cost"`
	SavedCost          float64           `json:"saved_cost"`
	ReportDate         string            `json:"report_date"`
	Username           string            `json:"username"`
	Email              string            `json:"email"`
	SiteName           string            `json:"site_name"`
	SiteURL            string            `json:"site_url"`
	HourlyStats        []HourlyUsageStat `json:"hourly_stats"`
	ModelStats         []ModelUsageStat  `json:"model_stats"`
	PeakHour           int               `json:"peak_hour"`
	Suggestions        []string          `json:"suggestions"`
}

// UpdateUserUsageReportConfigRequest is the request to update user's report config
type UpdateUserUsageReportConfigRequest struct {
	Enabled  *bool   `json:"enabled"`
	Schedule *string `json:"schedule"`
	Timezone *string `json:"timezone"`
}

// UserUsageReportRepository defines the interface for user usage report operations
type UserUsageReportRepository interface {
	GetUserReportConfig(ctx context.Context, userID int64) (*UserUsageReportConfig, error)
	UpdateUserReportConfig(ctx context.Context, userID int64, config *UserUsageReportConfig) error
	GetUsersForUsageReport(ctx context.Context, scope string, now time.Time) ([]User, error)
	GetActiveUserIDs(ctx context.Context, date time.Time) ([]int64, error)
}

// UserUsageReportService handles user usage report functionality
type UserUsageReportService struct {
	userRepo       UserRepository
	usageService   *UsageService
	settingService *SettingService
	emailService   *EmailService
	reportRepo     UserUsageReportRepository
}

// NewUserUsageReportService creates a new UserUsageReportService
func NewUserUsageReportService(
	userRepo UserRepository,
	usageService *UsageService,
	settingService *SettingService,
	emailService *EmailService,
	reportRepo UserUsageReportRepository,
) *UserUsageReportService {
	return &UserUsageReportService{
		userRepo:       userRepo,
		usageService:   usageService,
		settingService: settingService,
		emailService:   emailService,
		reportRepo:     reportRepo,
	}
}

// GetUserReportConfig retrieves user's usage report configuration
func (s *UserUsageReportService) GetUserReportConfig(ctx context.Context, userID int64) (*UserUsageReportConfig, error) {
	return s.reportRepo.GetUserReportConfig(ctx, userID)
}

// UpdateUserReportConfig updates user's usage report configuration
func (s *UserUsageReportService) UpdateUserReportConfig(ctx context.Context, userID int64, req *UpdateUserUsageReportConfigRequest) (*UserUsageReportConfig, error) {
	// Get current user to check email
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if user has a valid email bound
	if !isValidEmail(user.Email) {
		return nil, ErrUsageReportEmailNotBound
	}

	// Get current config
	config, err := s.reportRepo.GetUserReportConfig(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.Enabled != nil {
		config.Enabled = *req.Enabled
	}
	if req.Schedule != nil {
		config.Schedule = *req.Schedule
	}
	if req.Timezone != nil {
		config.Timezone = *req.Timezone
	}

	// Validate schedule format (HH:MM)
	if config.Schedule != "" {
		if !isValidScheduleFormat(config.Schedule) {
			return nil, infraerrors.BadRequest("INVALID_SCHEDULE", "schedule must be in HH:MM format")
		}
	}

	// Update config
	if err := s.reportRepo.UpdateUserReportConfig(ctx, userID, config); err != nil {
		return nil, err
	}

	return config, nil
}

// GenerateReportData generates usage report data for a user for the given date
func (s *UserUsageReportService) GenerateReportData(ctx context.Context, userID int64, reportDate time.Time) (*UserUsageReportData, error) {
	// Validate report date - cannot be in the future
	now := time.Now()
	if reportDate.After(now) {
		return nil, infraerrors.BadRequest("INVALID_DATE", "report date cannot be in the future")
	}

	// Get user info
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get site name and URL
	siteName := s.settingService.GetSiteName(ctx)
	siteURL := ""
	if publicSettings, err := s.settingService.GetPublicSettings(ctx); err == nil {
		siteURL = publicSettings.APIBaseURL
	}

	// TODO: Performance Optimization - N+1 Query Issue
	// Current implementation makes 3+ separate database queries:
	// 1. GetStatsByUser, 2. GetUserUsageTrendByUserID, 3. GetUserModelStats
	// Consider creating a batch query method: GetUserReportDataBatch(ctx, userID, startTime, endTime)
	// This would significantly improve performance for bulk report generation.

	// Calculate time range for the report date (full day in user's timezone)
	startTime := time.Date(reportDate.Year(), reportDate.Month(), reportDate.Day(), 0, 0, 0, 0, reportDate.Location())
	endTime := startTime.Add(24 * time.Hour)

	// Get usage stats
	stats, err := s.usageService.GetStatsByUser(ctx, userID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("get user stats: %w", err)
	}

	// Get hourly trend data
	hourlyTrend, err := s.usageService.GetUserUsageTrendByUserID(ctx, userID, startTime, endTime, "hour")
	if err != nil {
		// Log error but don't fail, use empty data
		hourlyTrend = nil
	}

	// Get model stats
	modelStats, err := s.usageService.GetUserModelStats(ctx, userID, startTime, endTime)
	if err != nil {
		// Log error but don't fail, use empty data
		modelStats = nil
	}

	// Calculate cache hit rate
	var cacheHitRate float64
	totalTokens := stats.TotalInputTokens + stats.TotalOutputTokens + stats.TotalCacheTokens
	if totalTokens > 0 {
		cacheHitRate = float64(stats.TotalCacheTokens) / float64(totalTokens) * 100
	}

	// Calculate saved cost
	//
	// Estimation Method:
	// - Cache tokens typically cost ~10% of regular input tokens (90% savings)
	// - Formula: saved_cost = (cache_tokens / total_tokens) × total_cost × 0.9
	//
	// Example:
	//   Total tokens: 10,000 (8,000 input + 2,000 cache)
	//   Total cost: $1.00
	//   Cache ratio: 2,000 / 10,000 = 20%
	//   Saved cost: $1.00 × 20% × 0.9 = $0.18
	//
	// Note: This is a rough estimate. Actual savings may vary by model and pricing tier.
	savedCost := float64(0)
	if totalTokens > 0 && stats.TotalCost > 0 {
		cacheRatio := float64(stats.TotalCacheTokens) / float64(totalTokens)
		savedCost = stats.TotalCost * cacheRatio * CacheSavingRate
	}

	// Add any discount between total cost and actual charged amount
	if stats.TotalCost > stats.TotalActualCost {
		savedCost += (stats.TotalCost - stats.TotalActualCost)
	}

	// Convert hourly trend to HourlyUsageStat
	hourlyStats := make([]HourlyUsageStat, 24)
	for i := 0; i < 24; i++ {
		hourlyStats[i] = HourlyUsageStat{Hour: i, Requests: 0}
	}
	peakHour := 0
	maxRequests := int64(0)
	if hourlyTrend != nil {
		for _, point := range hourlyTrend {
			// Parse date string to get hour
			// Format from DB: "YYYY-MM-DD HH24:00" -> "2026-01-29 15:00"
			var t time.Time
			var err error

			// Try parsing with different formats
			formats := []string{
				"2006-01-02 15:04",      // PostgreSQL: YYYY-MM-DD HH24:00
				"2006-01-02 15:04:05",   // With seconds
				time.RFC3339,             // ISO 8601
			}

			for _, format := range formats {
				t, err = time.Parse(format, point.Date)
				if err == nil {
					break
				}
			}

			if err == nil {
				hour := t.Hour()
				if hour >= 0 && hour < 24 {
					hourlyStats[hour].Requests += point.Requests
					if hourlyStats[hour].Requests > maxRequests {
						maxRequests = hourlyStats[hour].Requests
						peakHour = hour
					}
				}
			}
		}
	}

	// Convert model stats to ModelUsageStat
	modelStatsList := make([]ModelUsageStat, 0, len(modelStats))
	for _, ms := range modelStats {
		percent := float64(0)
		if stats.TotalRequests > 0 {
			percent = float64(ms.Requests) / float64(stats.TotalRequests) * 100
		}
		modelStatsList = append(modelStatsList, ModelUsageStat{
			Model:    ms.Model,
			Requests: ms.Requests,
			Cost:     ms.Cost,
			Percent:  percent,
		})
	}

	// Generate suggestions
	suggestions := generateSuggestions(stats, cacheHitRate, peakHour, modelStatsList)

	username := user.Username
	if username == "" {
		username = strings.Split(user.Email, "@")[0]
	}

	return &UserUsageReportData{
		TotalRequests:     stats.TotalRequests,
		TotalInputTokens:  stats.TotalInputTokens,
		TotalOutputTokens: stats.TotalOutputTokens,
		TotalCacheTokens:  stats.TotalCacheTokens,
		CacheHitRate:      cacheHitRate,
		TotalCost:         stats.TotalCost,
		ActualCost:        stats.TotalActualCost,
		SavedCost:         savedCost,
		ReportDate:        reportDate.Format("2006-01-02"),
		Username:          username,
		Email:             user.Email,
		SiteName:          siteName,
		SiteURL:           siteURL,
		HourlyStats:       hourlyStats,
		ModelStats:        modelStatsList,
		PeakHour:          peakHour,
		Suggestions:       suggestions,
	}, nil
}

// generateSuggestions generates intelligent suggestions based on usage data
// TODO: Internationalization - Hard-coded Chinese Text
// Current implementation has Chinese strings hardcoded, making it impossible to provide
// English reports for international users. Consider:
// 1. Move suggestion templates to i18n resource files
// 2. Add user language preference field
// 3. Select suggestion text based on user.Language
// Note: Suggestion texts are currently in Chinese. For internationalization,
// consider moving these to i18n resource files.
func generateSuggestions(stats *UsageStats, cacheHitRate float64, peakHour int, modelStats []ModelUsageStat) []string {
	suggestions := make([]string, 0, 5)

	// Cache-related suggestions
	if cacheHitRate < CacheHitRateLowThreshold {
		suggestions = append(suggestions, "缓存命中率较低，建议优化提示词以提高缓存复用")
	} else if cacheHitRate >= CacheHitRateHighThreshold {
		suggestions = append(suggestions, "缓存命中率超高！继续保持这个好习惯")
	}

	// Peak hour suggestions (work hours: 9-18, late night: 22-6)
	if peakHour >= 9 && peakHour <= 18 {
		suggestions = append(suggestions, fmt.Sprintf("使用高峰在 %d:00，工作时间使用频繁", peakHour))
	} else if peakHour >= 22 || peakHour <= 6 {
		suggestions = append(suggestions, fmt.Sprintf("深夜 %d:00 是你的使用高峰，注意休息哦", peakHour))
	}

	// Cost-related suggestions
	if stats.TotalCost > HighCostThreshold {
		suggestions = append(suggestions, "本日费用较高，考虑使用更经济的模型")
	}

	// Model diversity suggestions
	if len(modelStats) == 1 {
		suggestions = append(suggestions, "尝试不同的模型，可能会有意外收获")
	} else if len(modelStats) > 3 {
		suggestions = append(suggestions, "你在探索多个模型，太棒了！")
	}

	// Request count suggestions
	if stats.TotalRequests > 1000 {
		suggestions = append(suggestions, "今天真高产！共调用了 "+fmt.Sprintf("%d", stats.TotalRequests)+" 次")
	}

	// Limit to 3-4 suggestions for email readability
	if len(suggestions) > 4 {
		suggestions = suggestions[:4]
	}

	return suggestions
}

// SendReport sends usage report email to a user
func (s *UserUsageReportService) SendReport(ctx context.Context, userID int64, data *UserUsageReportData) error {
	// Build email content
	subject := fmt.Sprintf("[%s] 使用报告 - %s", data.SiteName, data.ReportDate)
	body, err := s.buildReportEmailBody(data)
	if err != nil {
		return fmt.Errorf("build email body: %w", err)
	}

	// Send email
	if err := s.emailService.SendEmail(ctx, data.Email, subject, body); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

// SendTestReport sends a test usage report to the user
func (s *UserUsageReportService) SendTestReport(ctx context.Context, userID int64) error {
	// Get user info
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Check if user has a valid email bound
	if !isValidEmail(user.Email) {
		return ErrUsageReportEmailNotBound
	}

	// Get user's report config for timezone
	config, err := s.reportRepo.GetUserReportConfig(ctx, userID)
	if err != nil {
		return err
	}

	// Load user's configured timezone
	loc, err := time.LoadLocation(config.Timezone)
	if err != nil {
		loc = time.FixedZone("UTC+8", 8*3600) // Default to UTC+8
	}

	// Generate report for yesterday in user's timezone
	yesterday := time.Now().In(loc).AddDate(0, 0, -1)

	data, err := s.GenerateReportData(ctx, userID, yesterday)
	if err != nil {
		return err
	}

	return s.SendReport(ctx, userID, data)
}

// isValidEmail checks if email is a valid bound email (not synthetic)
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	// Check for synthetic/invalid email domains
	if strings.HasSuffix(email, ".invalid") {
		return false
	}
	return strings.Contains(email, "@")
}

// isValidScheduleFormat validates HH:MM format
func isValidScheduleFormat(schedule string) bool {
	if len(schedule) != 5 {
		return false
	}
	if schedule[2] != ':' {
		return false
	}
	hours := schedule[0:2]
	minutes := schedule[3:5]
	if hours < "00" || hours > "23" {
		return false
	}
	if minutes < "00" || minutes > "59" {
		return false
	}
	return true
}

// buildReportEmailBody builds the HTML email body for usage report
func (s *UserUsageReportService) buildReportEmailBody(data *UserUsageReportData) (string, error) {
	var buf bytes.Buffer
	if err := usageReportTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

const usageReportEmailTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
            background: #fafafa;
            padding: 20px;
            line-height: 1.6;
        }
        .container {
            max-width: 700px;
            margin: 0 auto;
            background: #ffffff;
            border-radius: 16px;
            overflow: hidden;
            box-shadow: 0 2px 12px rgba(251, 146, 60, 0.1);
        }
        .header {
            background: linear-gradient(135deg, #fb923c 0%, #f97316 100%);
            color: white;
            padding: 40px 32px 32px;
            position: relative;
            overflow: hidden;
        }
        .header::before {
            content: '';
            position: absolute;
            top: 20px;
            right: 30px;
            width: 48px;
            height: 48px;
            background-image: url('{{.SiteURL}}/email-logo.png');
            background-size: contain;
            background-repeat: no-repeat;
            background-position: center;
            opacity: 0.3;
        }
        .header-content {
            position: relative;
            z-index: 1;
        }
        .header h1 {
            font-size: 28px;
            font-weight: 800;
            margin-bottom: 4px;
            letter-spacing: -0.5px;
        }
        .header .subtitle {
            font-size: 14px;
            opacity: 0.95;
            font-weight: 500;
            margin-bottom: 12px;
        }
        .header .date {
            font-size: 13px;
            opacity: 0.85;
            font-weight: 400;
        }
        .content {
            padding: 32px;
        }
        .greeting {
            font-size: 16px;
            color: #1f2937;
            margin-bottom: 24px;
            font-weight: 500;
        }
        .greeting strong {
            color: #f97316;
            font-weight: 700;
        }

        /* Highlight Card */
        .highlight-card {
            background: linear-gradient(135deg, #fb923c 0%, #f97316 100%);
            border-radius: 16px;
            padding: 28px;
            margin-bottom: 24px;
            color: white;
            box-shadow: 0 8px 20px rgba(251, 146, 60, 0.25);
            display: table;
            width: 100%;
        }
        .highlight-left {
            display: table-cell;
            width: 55%;
            padding-right: 20px;
            vertical-align: middle;
        }
        .highlight-right {
            display: table-cell;
            width: 45%;
            padding-left: 20px;
            border-left: 2px solid rgba(255, 255, 255, 0.3);
            vertical-align: middle;
        }
        .highlight-label {
            font-size: 14px;
            opacity: 0.95;
            margin-bottom: 10px;
        }
        .highlight-value {
            font-size: 42px;
            font-weight: 900;
            letter-spacing: -2px;
            margin-bottom: 12px;
        }
        .highlight-note {
            font-size: 13px;
            opacity: 0.9;
            margin-bottom: 12px;
        }
        .highlight-tip {
            font-size: 12px;
            opacity: 0.85;
            margin-top: 8px;
            font-style: italic;
        }
        .progress-bar {
            background: rgba(255, 255, 255, 0.3);
            border-radius: 10px;
            height: 10px;
            overflow: hidden;
        }
        .progress-fill {
            height: 100%;
            background: white;
            border-radius: 10px;
            transition: width 0.6s ease;
            box-shadow: 0 0 10px rgba(255, 255, 255, 0.5);
        }
        .token-mini-grid {
            display: grid;
            grid-template-columns: repeat(3, 1fr);
            gap: 8px;
            margin-top: 8px;
        }
        .token-mini-card {
            background: rgba(255, 255, 255, 0.2);
            border-radius: 8px;
            padding: 10px;
            text-align: center;
        }
        .token-mini-label {
            font-size: 10px;
            opacity: 0.9;
            margin-bottom: 4px;
            text-transform: uppercase;
        }
        .token-mini-value {
            font-size: 16px;
            font-weight: 700;
        }

        /* Stats Grid */
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 14px;
            margin-bottom: 24px;
        }
        .stat-card {
            background: #fff7ed;
            border: 2px solid #fed7aa;
            border-radius: 12px;
            padding: 18px;
        }
        .stat-label {
            font-size: 12px;
            color: #9a3412;
            margin-bottom: 8px;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        .stat-value {
            font-size: 24px;
            font-weight: 800;
            color: #c2410c;
        }

        /* Two Column Layout */
        .two-column {
            display: table;
            width: 100%;
            margin-bottom: 24px;
        }
        .column-left {
            display: table-cell;
            width: 60%;
            padding-right: 12px;
            vertical-align: top;
        }
        .column-right {
            display: table-cell;
            width: 40%;
            padding-left: 12px;
            vertical-align: top;
        }

        /* Section */
        .section-title {
            font-size: 15px;
            font-weight: 700;
            color: #1f2937;
            margin-bottom: 14px;
            padding-bottom: 8px;
            border-bottom: 3px solid #fed7aa;
        }

        /* Hourly Chart */
        .hourly-chart {
            background: #fffbeb;
            border-radius: 12px;
            padding: 18px;
            margin-bottom: 20px;
        }
        .chart-bars {
            display: flex;
            align-items: flex-end;
            justify-content: space-between;
            height: 100px;
            gap: 1px;
            margin-bottom: 8px;
        }
        .chart-bar {
            flex: 1;
            background: linear-gradient(to top, #fb923c, #fed7aa);
            border-radius: 3px 3px 0 0;
            min-height: 2px;
        }
        .chart-labels {
            display: flex;
            justify-content: space-between;
            font-size: 10px;
            color: #92400e;
            font-weight: 600;
            padding: 0 2px;
        }

        /* Model Stats */
        .model-list {
            background: #fffbeb;
            border-radius: 12px;
            padding: 14px;
            margin-bottom: 20px;
        }
        .model-item {
            padding: 10px 0;
            border-bottom: 1px solid #fed7aa;
        }
        .model-item:last-child {
            border-bottom: none;
        }
        .model-name {
            font-size: 13px;
            font-weight: 700;
            color: #1f2937;
            margin-bottom: 6px;
        }
        .model-bar-container {
            height: 6px;
            background: #fed7aa;
            border-radius: 3px;
            overflow: hidden;
            margin-bottom: 4px;
        }
        .model-bar {
            height: 100%;
            background: linear-gradient(90deg, #fb923c, #f97316);
            border-radius: 3px;
        }
        .model-stats-text {
            font-size: 11px;
            color: #92400e;
        }

        /* Suggestions */
        .suggestions {
            background: linear-gradient(135deg, #fffbeb 0%, #fff7ed 100%);
            border-left: 4px solid #fb923c;
            border-radius: 12px;
            padding: 16px;
            margin-bottom: 20px;
        }
        .suggestion-item {
            padding: 6px 0;
            font-size: 13px;
            color: #78350f;
            line-height: 1.5;
        }
        .suggestion-item::before {
            content: '▸';
            color: #fb923c;
            font-weight: bold;
            margin-right: 6px;
        }

        /* Token Details */
        .token-grid {
            display: grid;
            grid-template-columns: repeat(3, 1fr);
            gap: 10px;
        }
        .token-card {
            background: #fffbeb;
            border-radius: 10px;
            padding: 14px;
            text-align: center;
            border: 2px solid #fed7aa;
        }
        .token-label {
            font-size: 11px;
            color: #92400e;
            margin-bottom: 6px;
            font-weight: 600;
            text-transform: uppercase;
        }
        .token-value {
            font-size: 16px;
            font-weight: 800;
            color: #c2410c;
        }

        /* Footer */
        .footer {
            background: #fff7ed;
            padding: 20px 32px;
            text-align: center;
            font-size: 12px;
            color: #92400e;
            border-top: 2px solid #fed7aa;
        }
        .footer a {
            color: #f97316;
            text-decoration: none;
            font-weight: 600;
        }

        /* Mobile */
        @media only screen and (max-width: 600px) {
            .two-column {
                display: block;
            }
            .column-left, .column-right {
                display: block;
                width: 100%;
                padding: 0;
            }
            .column-right {
                margin-top: 20px;
            }
            .highlight-card {
                display: block;
            }
            .highlight-left, .highlight-right {
                display: block;
                width: 100%;
                padding: 0;
                border-left: none;
            }
            .highlight-right {
                margin-top: 16px;
                padding-top: 16px;
                border-top: 2px solid rgba(255, 255, 255, 0.3);
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="header-content">
                <h1>{{.SiteName}}</h1>
                <div class="subtitle">🚌 AI编程巴士</div>
                <div class="date">每日使用报告 · {{.ReportDate}}</div>
            </div>
        </div>

        <div class="content">
            <p class="greeting">你好 <strong>{{.Username}}</strong>，昨天你通过 <strong>AI编程巴士</strong> 调用了顶尖模型 <strong>{{.TotalRequests}}</strong> 次</p>

            <!-- Cache Hit Rate -->
            <div class="highlight-card">
                <div class="highlight-left">
                    <div class="highlight-label">缓存命中率</div>
                    <div class="highlight-value">{{printf "%.1f" .CacheHitRate}}%</div>
                    {{if gt .SavedCost 0.01}}
                    <div class="highlight-note">对比其他没有缓存的 API 服务，你节省了 ${{printf "%.2f" .SavedCost}}</div>
                    {{end}}
                    <div class="progress-bar">
                        <div class="progress-fill" style="width: {{printf "%.0f" .CacheHitRate}}%"></div>
                    </div>
                    <div class="highlight-tip">💡 点击卡片查看详细分析</div>
                </div>
                <div class="highlight-right">
                    <div class="token-mini-grid">
                        <div class="token-mini-card">
                            <div class="token-mini-label">输入</div>
                            <div class="token-mini-value">{{.TotalInputTokens}}</div>
                        </div>
                        <div class="token-mini-card">
                            <div class="token-mini-label">输出</div>
                            <div class="token-mini-value">{{.TotalOutputTokens}}</div>
                        </div>
                        <div class="token-mini-card">
                            <div class="token-mini-label">缓存</div>
                            <div class="token-mini-value">{{.TotalCacheTokens}}</div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Quick Stats -->
            <div class="stats-grid">
                <div class="stat-card">
                    <div class="stat-label">请求次数</div>
                    <div class="stat-value">{{.TotalRequests}}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">总费用</div>
                    <div class="stat-value">${{printf "%.2f" .TotalCost}}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">实际扣费</div>
                    <div class="stat-value">${{printf "%.2f" .ActualCost}}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">节省金额</div>
                    <div class="stat-value">${{printf "%.2f" .SavedCost}}</div>
                </div>
            </div>

            <!-- Two Column Layout -->
            <div class="two-column">
                <!-- Left Column: Charts -->
                <div class="column-left">
                    <!-- Hourly Usage Chart -->
                    {{if .HourlyStats}}
                    <div class="section-title">24小时使用分布{{if gt .PeakHour 0}} (高峰: {{.PeakHour}}:00){{end}}</div>
                    <div class="hourly-chart">
                        <div class="chart-bars">
                            {{range .HourlyStats}}
                            {{$maxRequests := 1}}
                            {{range $.HourlyStats}}{{if gt .Requests $maxRequests}}{{$maxRequests = .Requests}}{{end}}{{end}}
                            {{$height := 0}}
                            {{if gt $maxRequests 0}}
                            {{$height = printf "%.0f" (mulf (divf (printf "%d" .Requests) (printf "%d" $maxRequests)) 100)}}
                            {{end}}
                            <div class="chart-bar" style="height: {{$height}}%;" title="{{.Hour}}:00 - {{.Requests}} 次"></div>
                            {{end}}
                        </div>
                        <div class="chart-labels">
                            <span>0</span>
                            <span>6</span>
                            <span>12</span>
                            <span>18</span>
                            <span>24</span>
                        </div>
                    </div>
                    {{end}}

                    <!-- Model Distribution -->
                    {{if .ModelStats}}
                    <div class="section-title">模型使用分布</div>
                    <div class="model-list">
                        {{range .ModelStats}}
                        <div class="model-item">
                            <div class="model-name">{{.Model}}</div>
                            <div class="model-bar-container">
                                <div class="model-bar" style="width: {{printf "%.0f" .Percent}}%"></div>
                            </div>
                            <div class="model-stats-text">
                                {{.Requests}} 次 · ${{printf "%.2f" .Cost}} · {{printf "%.1f" .Percent}}%
                            </div>
                        </div>
                        {{end}}
                    </div>
                    {{end}}
                </div>

                <!-- Right Column: Suggestions & Tokens -->
                <div class="column-right">
                    <!-- Suggestions -->
                    {{if .Suggestions}}
                    <div class="section-title">智能建议</div>
                    <div class="suggestions">
                        {{range .Suggestions}}
                        <div class="suggestion-item">{{.}}</div>
                        {{end}}
                    </div>
                    {{end}}

                    <!-- Token Details -->
                    <div class="section-title">Token 明细</div>
                    <div class="token-grid">
                        <div class="token-card">
                            <div class="token-label">输入</div>
                            <div class="token-value">{{.TotalInputTokens}}</div>
                        </div>
                        <div class="token-card">
                            <div class="token-label">输出</div>
                            <div class="token-value">{{.TotalOutputTokens}}</div>
                        </div>
                        <div class="token-card">
                            <div class="token-label">缓存</div>
                            <div class="token-value">{{.TotalCacheTokens}}</div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="footer">
            <p>这是你的每日使用报告，帮助你更好地了解 API 使用情况</p>
            <p style="margin-top: 8px;">不想收到邮件？在<a href="#">个人设置</a>中关闭即可</p>
        </div>
    </div>
</body>
</html>`
