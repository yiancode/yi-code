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

// usageReportTemplate is the pre-parsed email template (parsed once at package init)
var usageReportTemplate = template.Must(template.New("report").Parse(usageReportEmailTemplate))

// UserUsageReportConfig represents user's usage report configuration
type UserUsageReportConfig struct {
	Enabled  bool   `json:"enabled"`
	Schedule string `json:"schedule"` // HH:MM format
	Timezone string `json:"timezone"`
}

// UserUsageReportData represents the data for a usage report
type UserUsageReportData struct {
	TotalRequests      int64   `json:"total_requests"`
	TotalInputTokens   int64   `json:"total_input_tokens"`
	TotalOutputTokens  int64   `json:"total_output_tokens"`
	TotalCacheTokens   int64   `json:"total_cache_tokens"`
	CacheHitRate       float64 `json:"cache_hit_rate"`
	TotalCost          float64 `json:"total_cost"`
	ActualCost         float64 `json:"actual_cost"`
	ReportDate         string  `json:"report_date"`
	Username           string  `json:"username"`
	Email              string  `json:"email"`
	SiteName           string  `json:"site_name"`
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
	// Get user info
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get site name
	siteName := s.settingService.GetSiteName(ctx)

	// Calculate time range for the report date (full day in user's timezone)
	startTime := time.Date(reportDate.Year(), reportDate.Month(), reportDate.Day(), 0, 0, 0, 0, reportDate.Location())
	endTime := startTime.Add(24 * time.Hour)

	// Get usage stats
	stats, err := s.usageService.GetStatsByUser(ctx, userID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("get user stats: %w", err)
	}

	// Calculate cache hit rate
	var cacheHitRate float64
	totalTokens := stats.TotalInputTokens + stats.TotalOutputTokens + stats.TotalCacheTokens
	if totalTokens > 0 {
		cacheHitRate = float64(stats.TotalCacheTokens) / float64(totalTokens) * 100
	}

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
		ReportDate:        reportDate.Format("2006-01-02"),
		Username:          username,
		Email:             user.Email,
		SiteName:          siteName,
	}, nil
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
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background-color: #f5f5f5;
            margin: 0;
            padding: 20px;
            line-height: 1.6;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            border-radius: 12px;
            overflow: hidden;
            box-shadow: 0 4px 12px rgba(0,0,0,0.1);
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        .header h1 {
            margin: 0 0 8px 0;
            font-size: 24px;
            font-weight: 600;
        }
        .header p {
            margin: 0;
            opacity: 0.9;
            font-size: 14px;
        }
        .content {
            padding: 30px;
        }
        .greeting {
            font-size: 16px;
            color: #333;
            margin-bottom: 24px;
        }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 16px;
            margin-bottom: 24px;
        }
        .stat-card {
            background: linear-gradient(135deg, #f6f8fb 0%, #f1f3f6 100%);
            border-radius: 10px;
            padding: 20px;
            text-align: center;
        }
        .stat-card.highlight {
            background: linear-gradient(135deg, #667eea15 0%, #764ba215 100%);
            border: 1px solid #667eea30;
        }
        .stat-label {
            font-size: 12px;
            color: #666;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            margin-bottom: 4px;
        }
        .stat-value {
            font-size: 24px;
            font-weight: 700;
            color: #333;
        }
        .stat-value.cost {
            color: #667eea;
        }
        .detail-table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 24px;
        }
        .detail-table th,
        .detail-table td {
            padding: 12px 16px;
            text-align: left;
            border-bottom: 1px solid #eee;
        }
        .detail-table th {
            background-color: #f8f9fa;
            font-weight: 600;
            color: #666;
            font-size: 12px;
            text-transform: uppercase;
        }
        .detail-table td {
            color: #333;
        }
        .footer {
            background-color: #f8f9fa;
            padding: 20px 30px;
            text-align: center;
            color: #999;
            font-size: 12px;
            border-top: 1px solid #eee;
        }
        .footer a {
            color: #667eea;
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.SiteName}}</h1>
            <p>每日使用报告 · {{.ReportDate}}</p>
        </div>
        <div class="content">
            <p class="greeting">Hi {{.Username}}，以下是您昨日的使用概况：</p>

            <div class="stats-grid">
                <div class="stat-card">
                    <div class="stat-label">请求次数</div>
                    <div class="stat-value">{{.TotalRequests}}</div>
                </div>
                <div class="stat-card highlight">
                    <div class="stat-label">缓存命中率</div>
                    <div class="stat-value">{{printf "%.1f" .CacheHitRate}}%</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">总费用</div>
                    <div class="stat-value cost">${{printf "%.4f" .TotalCost}}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">实际扣费</div>
                    <div class="stat-value cost">${{printf "%.4f" .ActualCost}}</div>
                </div>
            </div>

            <table class="detail-table">
                <thead>
                    <tr>
                        <th>项目</th>
                        <th>数量</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td>输入 Tokens</td>
                        <td>{{.TotalInputTokens}}</td>
                    </tr>
                    <tr>
                        <td>输出 Tokens</td>
                        <td>{{.TotalOutputTokens}}</td>
                    </tr>
                    <tr>
                        <td>缓存 Tokens</td>
                        <td>{{.TotalCacheTokens}}</td>
                    </tr>
                </tbody>
            </table>
        </div>
        <div class="footer">
            <p>此邮件由系统自动发送，如需取消订阅，请在个人设置中关闭使用报告功能。</p>
        </div>
    </div>
</body>
</html>`
