package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	// WeChatSyntheticEmailDomain 是微信登录用户的合成邮箱后缀（RFC 保留域名）
	WeChatSyntheticEmailDomain = "@wechat-auth.invalid"
)

// wechatLoginResponse 微信服务器返回的验证结果
type wechatLoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data"` // 返回 WeChat ID
}

// WeChatAuthRequest 微信登录请求
type WeChatAuthRequest struct {
	Code string `form:"code" binding:"required"` // 验证码
}

// WeChatAuth 微信公众号验证码登录
// GET /api/v1/auth/oauth/wechat?code=xxx
func (h *AuthHandler) WeChatAuth(c *gin.Context) {
	var req WeChatAuthRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "验证码不能为空")
		return
	}

	// 检查微信登录是否启用
	cfg, err := h.getWeChatConfig(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// 使用验证码获取 WeChat ID
	wechatID, err := getWeChatIDByCode(c.Request.Context(), cfg, req.Code)
	if err != nil {
		log.Printf("[WeChat Auth] Failed to get wechat id: %v", err)
		response.ErrorFrom(c, infraerrors.BadRequest("WECHAT_AUTH_FAILED", err.Error()))
		return
	}

	// 使用合成邮箱作为唯一标识
	email := wechatSyntheticEmail(wechatID)
	username := "wechat_" + wechatID

	// 登录或注册
	token, user, err := h.authService.LoginOrRegisterOAuth(c.Request.Context(), email, username)
	if err != nil {
		log.Printf("[WeChat Auth] Login or register failed: %v", err)
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		User:        dto.UserFromService(user),
	})
}

// WeChatBind 微信账号绑定（已登录用户）
// GET /api/v1/auth/oauth/wechat/bind?code=xxx
func (h *AuthHandler) WeChatBind(c *gin.Context) {
	var req WeChatAuthRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "验证码不能为空")
		return
	}

	// 检查微信登录是否启用
	cfg, err := h.getWeChatConfig(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// 使用验证码获取 WeChat ID
	wechatID, err := getWeChatIDByCode(c.Request.Context(), cfg, req.Code)
	if err != nil {
		log.Printf("[WeChat Bind] Failed to get wechat id: %v", err)
		response.ErrorFrom(c, infraerrors.BadRequest("WECHAT_AUTH_FAILED", err.Error()))
		return
	}

	// 检查该 WeChat ID 是否已被绑定
	email := wechatSyntheticEmail(wechatID)
	exists, err := h.userService.ExistsByEmail(c.Request.Context(), email)
	if err != nil {
		log.Printf("[WeChat Bind] Failed to check email exists: %v", err)
		response.ErrorFrom(c, infraerrors.InternalServer("INTERNAL_ERROR", "检查绑定状态失败"))
		return
	}
	if exists {
		response.ErrorFrom(c, infraerrors.Conflict("WECHAT_ALREADY_BOUND", "该微信账号已被其他用户绑定"))
		return
	}

	// 绑定成功返回微信 ID（前端可用于显示绑定状态）
	response.Success(c, gin.H{
		"wechat_id": wechatID,
		"message":   "绑定成功",
	})
}

// getWeChatConfig 获取微信登录配置
func (h *AuthHandler) getWeChatConfig(ctx context.Context) (*service.WeChatConfig, error) {
	if h.settingSvc == nil {
		return nil, infraerrors.ServiceUnavailable("CONFIG_NOT_READY", "设置服务未就绪")
	}

	cfg, err := h.settingSvc.GetWeChatConfig(ctx)
	if err != nil {
		return nil, infraerrors.InternalServer("CONFIG_ERROR", "获取微信配置失败")
	}

	if !cfg.Enabled {
		return nil, infraerrors.NotFound("WECHAT_AUTH_DISABLED", "管理员未开启微信登录")
	}

	if strings.TrimSpace(cfg.ServerAddress) == "" {
		return nil, infraerrors.InternalServer("WECHAT_CONFIG_INVALID", "微信服务器地址未配置")
	}

	if strings.TrimSpace(cfg.ServerToken) == "" {
		return nil, infraerrors.InternalServer("WECHAT_CONFIG_INVALID", "微信服务器令牌未配置")
	}

	return cfg, nil
}

// getWeChatIDByCode 通过验证码获取 WeChat ID
func getWeChatIDByCode(ctx context.Context, cfg *service.WeChatConfig, code string) (string, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return "", errors.New("验证码不能为空")
	}

	// 向微信服务器发送请求
	url := fmt.Sprintf("%s/api/wechat/user?code=%s", strings.TrimRight(cfg.ServerAddress, "/"), code)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", cfg.ServerToken)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求微信服务器失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("微信服务器返回错误: %d", resp.StatusCode)
	}

	var result wechatLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if !result.Success {
		if result.Message != "" {
			return "", errors.New(result.Message)
		}
		return "", errors.New("验证失败")
	}

	if strings.TrimSpace(result.Data) == "" {
		return "", errors.New("验证码错误或已过期")
	}

	return strings.TrimSpace(result.Data), nil
}

// wechatSyntheticEmail 生成微信用户的合成邮箱
func wechatSyntheticEmail(wechatID string) string {
	wechatID = strings.TrimSpace(wechatID)
	if wechatID == "" {
		return ""
	}
	return "wechat-" + wechatID + WeChatSyntheticEmailDomain
}
