package handler

import (
	"context"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// UserUsageReportHandler handles user usage report API endpoints
type UserUsageReportHandler struct {
	reportService  *service.UserUsageReportService
	settingService *service.SettingService
	userRepo       service.UserRepository
}

// NewUserUsageReportHandler creates a new UserUsageReportHandler
func NewUserUsageReportHandler(
	reportService *service.UserUsageReportService,
	settingService *service.SettingService,
	userRepo service.UserRepository,
) *UserUsageReportHandler {
	return &UserUsageReportHandler{
		reportService:  reportService,
		settingService: settingService,
		userRepo:       userRepo,
	}
}

// checkEmailBound checks if user's email is bound and valid
func (h *UserUsageReportHandler) checkEmailBound(ctx context.Context, userID int64) bool {
	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return false
	}
	return user.Email != "" && !strings.HasSuffix(user.Email, ".invalid")
}

// GetConfigResponse represents the response for get config endpoint
type GetConfigResponse struct {
	Enabled       bool   `json:"enabled"`
	Schedule      string `json:"schedule"`
	Timezone      string `json:"timezone"`
	GlobalEnabled bool   `json:"global_enabled"`
	EmailBound    bool   `json:"email_bound"`
}

// GetConfig returns user's usage report configuration
// GET /api/v1/user/usage-report/config
func (h *UserUsageReportHandler) GetConfig(c *gin.Context) {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Get global config
	globalConfig, err := h.settingService.GetUsageReportConfig(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// Get user config
	config, err := h.reportService.GetUserReportConfig(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// Check if email is bound
	emailBound := h.checkEmailBound(c.Request.Context(), subject.UserID)

	response.Success(c, GetConfigResponse{
		Enabled:       config.Enabled,
		Schedule:      config.Schedule,
		Timezone:      config.Timezone,
		GlobalEnabled: globalConfig.GlobalEnabled,
		EmailBound:    emailBound,
	})
}

// UpdateConfigRequest represents the request to update usage report config
type UpdateConfigRequest struct {
	Enabled  *bool   `json:"enabled"`
	Schedule *string `json:"schedule"`
	Timezone *string `json:"timezone"`
}

// UpdateConfig updates user's usage report configuration
// PUT /api/v1/user/usage-report/config
func (h *UserUsageReportHandler) UpdateConfig(c *gin.Context) {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// Check if global enabled
	globalConfig, err := h.settingService.GetUsageReportConfig(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	if !globalConfig.GlobalEnabled && req.Enabled != nil && *req.Enabled {
		response.ErrorFrom(c, service.ErrUsageReportNotEnabled)
		return
	}

	// Check email binding before enabling the report
	if req.Enabled != nil && *req.Enabled {
		if !h.checkEmailBound(c.Request.Context(), subject.UserID) {
			response.ErrorFrom(c, service.ErrUsageReportEmailNotBound)
			return
		}
	}

	updateReq := &service.UpdateUserUsageReportConfigRequest{
		Enabled:  req.Enabled,
		Schedule: req.Schedule,
		Timezone: req.Timezone,
	}

	config, err := h.reportService.UpdateUserReportConfig(c.Request.Context(), subject.UserID, updateReq)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// Check if email is bound (for response)
	emailBound := h.checkEmailBound(c.Request.Context(), subject.UserID)

	// Return complete config response (same structure as GetConfig)
	response.Success(c, GetConfigResponse{
		Enabled:       config.Enabled,
		Schedule:      config.Schedule,
		Timezone:      config.Timezone,
		GlobalEnabled: globalConfig.GlobalEnabled,
		EmailBound:    emailBound,
	})
}

// SendTestReport sends a test usage report to the user
// POST /api/v1/user/usage-report/test
func (h *UserUsageReportHandler) SendTestReport(c *gin.Context) {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	if err := h.reportService.SendTestReport(c.Request.Context(), subject.UserID); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Test report sent successfully"})
}
