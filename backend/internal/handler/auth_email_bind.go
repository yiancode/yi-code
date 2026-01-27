package handler

import (
	"strings"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/ip"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

// 合成邮箱域名常量（用于识别微信/LinuxDo等OAuth用户）
const (
	LinuxDoSyntheticEmailDomain = "@linuxdo.invalid"
)

// BindEmailRequest 邮箱绑定请求
type BindEmailRequest struct {
	Email      string `json:"email" binding:"required,email"`
	VerifyCode string `json:"verify_code" binding:"required"`
}

// BindEmailResponse 邮箱绑定响应
type BindEmailResponse struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

// SendBindEmailCodeRequest 发送绑定邮箱验证码请求
type SendBindEmailCodeRequest struct {
	Email          string `json:"email" binding:"required,email"`
	TurnstileToken string `json:"turnstile_token"`
}

// SendBindEmailCode 发送绑定邮箱验证码（已登录用户）
// POST /api/v1/auth/send-bind-email-code
func (h *AuthHandler) SendBindEmailCode(c *gin.Context) {
	var req SendBindEmailCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数无效: "+err.Error())
		return
	}

	// 获取当前登录用户
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "用户未登录")
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))

	// 验证邮箱不是合成邮箱格式
	if isSyntheticEmail(email) {
		response.ErrorFrom(c, infraerrors.BadRequest("INVALID_EMAIL", "不允许使用该邮箱格式"))
		return
	}

	// 获取当前用户信息
	currentUser, err := h.userService.GetByID(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// 检查用户当前邮箱是否为合成邮箱（只有合成邮箱用户才需要绑定真实邮箱）
	if !isSyntheticEmail(currentUser.Email) {
		response.ErrorFrom(c, infraerrors.BadRequest("ALREADY_HAS_EMAIL", "您已绑定真实邮箱"))
		return
	}

	// 检查邮箱是否已被其他用户使用
	existingUser, err := h.userService.GetByEmail(c.Request.Context(), email)
	if err == nil && existingUser != nil && existingUser.ID != subject.UserID {
		response.ErrorFrom(c, infraerrors.Conflict("EMAIL_ALREADY_USED", "该邮箱已被其他用户使用"))
		return
	}

	// Turnstile 验证
	if err := h.authService.VerifyTurnstile(c.Request.Context(), req.TurnstileToken, ip.GetClientIP(c)); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// 发送验证码
	if err := h.authService.SendBindEmailCode(c.Request.Context(), email); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, SendVerifyCodeResponse{
		Message:   "验证码已发送",
		Countdown: 60,
	})
}

// BindEmail 绑定邮箱（已登录用户）
// POST /api/v1/auth/bind-email
func (h *AuthHandler) BindEmail(c *gin.Context) {
	var req BindEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数无效: "+err.Error())
		return
	}

	// 获取当前登录用户
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "用户未登录")
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	verifyCode := strings.TrimSpace(req.VerifyCode)

	// 验证邮箱不是合成邮箱格式
	if isSyntheticEmail(email) {
		response.ErrorFrom(c, infraerrors.BadRequest("INVALID_EMAIL", "不允许使用该邮箱格式"))
		return
	}

	// 获取当前用户信息
	currentUser, err := h.userService.GetByID(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// 检查用户当前邮箱是否为合成邮箱（只有合成邮箱用户才需要绑定真实邮箱）
	if !isSyntheticEmail(currentUser.Email) {
		response.ErrorFrom(c, infraerrors.BadRequest("ALREADY_HAS_EMAIL", "您已绑定真实邮箱"))
		return
	}

	// 检查邮箱是否已被其他用户使用
	existingUser, err := h.userService.GetByEmail(c.Request.Context(), email)
	if err == nil && existingUser != nil && existingUser.ID != subject.UserID {
		response.ErrorFrom(c, infraerrors.Conflict("EMAIL_ALREADY_USED", "该邮箱已被其他用户使用"))
		return
	}

	// 验证验证码
	if err := h.authService.VerifyEmailCode(c.Request.Context(), email, verifyCode); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// 更新用户邮箱
	if err := h.userService.UpdateEmail(c.Request.Context(), subject.UserID, email); err != nil {
		response.ErrorFrom(c, infraerrors.InternalServer("UPDATE_FAILED", "更新邮箱失败"))
		return
	}

	response.Success(c, BindEmailResponse{
		Email:   email,
		Message: "邮箱绑定成功",
	})
}

// isSyntheticEmail 判断是否为合成邮箱
func isSyntheticEmail(email string) bool {
	email = strings.ToLower(email)
	return strings.HasSuffix(email, WeChatSyntheticEmailDomain) ||
		strings.HasSuffix(email, LinuxDoSyntheticEmailDomain)
}
