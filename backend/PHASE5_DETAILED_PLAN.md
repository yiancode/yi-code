# 阶段 5：Handler 和路由层实现 - 详细开发计划

## 概述

本阶段实现充值功能的 HTTP Handler 层和路由配置，包括用户充值接口和微信支付回调接口。

**前置条件：** 阶段 1-4 已完成（配置、数据库、SDK、Service 层）

---

## 📋 任务清单

### 5.1 创建 DTO 定义 ✅

**文件:** `internal/handler/dto/recharge.go`

**内容包括:**
```go
// 请求 DTO
- CreateRechargeOrderRequest  // 创建充值订单请求
- PaginationRequest            // 分页请求参数

// 响应 DTO
- CreateRechargeOrderResponse  // 创建充值订单响应
- RechargeOrderResponse        // 充值订单详情响应
- BalanceLogResponse           // 余额变动日志响应
- PaginationResponse           // 分页响应
- RechargeOrderListResponse    // 充值订单列表响应
- BalanceLogListResponse       // 余额日志列表响应
```

**注意事项:**
- 所有金额字段使用 `float64` 类型
- 时间字段使用 `time.Time` 类型
- 可选字段使用指针类型 `*string`, `*float64`
- JSON tag 使用 snake_case 命名
- 使用 `omitempty` 标记可选字段

---

### 5.2 创建 RechargeHandler ✅

**文件:** `internal/handler/recharge_handler.go`

**依赖注入:**
```go
type RechargeHandler struct {
    rechargeService *service.RechargeService
    balanceService  *service.BalanceService
}

func NewRechargeHandler(
    rechargeService *service.RechargeService,
    balanceService *service.BalanceService,
) *RechargeHandler
```

**实现接口 (7个):**

#### 1. CreateOrder - 创建充值订单
```go
POST /api/v1/recharge/orders
Content-Type: application/json
Authorization: Bearer <token>

Request Body:
{
  "amount": 100.00,
  "payment_channel": "native",  // "jsapi" | "native"
  "openid": "oXXXX"             // jsapi 时必需
}

Response 200:
{
  "order_no": "R20260124123456789",
  "amount": 100.00,
  "expired_at": "2026-01-24T14:34:56Z",
  "code_url": "weixin://wxpay/...",  // native 支付
  "prepay_id": "wx24...",             // jsapi 支付
  "jsapi_params": {                   // jsapi 支付
    "appId": "wx...",
    "timeStamp": "1234567890",
    "nonceStr": "...",
    "package": "prepay_id=...",
    "signType": "RSA",
    "paySign": "..."
  }
}
```

**实现步骤:**
1. 从 JWT 中获取 `subject.UserID`
2. 绑定并验证 JSON 请求体
3. 验证 `payment_channel` 必须是 "jsapi" 或 "native"
4. 如果是 jsapi，验证 `openid` 不为空
5. 构建 `service.CreateRechargeOrderRequest`
6. 添加 `UserIP` (从 `c.ClientIP()`) 和 `UserAgent` (从 header)
7. 调用 `rechargeService.CreateOrder()`
8. 转换为 DTO 响应并返回

#### 2. GetOrder - 查询订单详情
```go
GET /api/v1/recharge/orders/:order_no
Authorization: Bearer <token>

Response 200:
{
  "id": 123,
  "order_no": "R20260124123456789",
  "user_id": 456,
  "amount": 100.00,
  "actual_amount": 100.00,
  "currency": "CNY",
  "payment_method": "wechat",
  "payment_channel": "native",
  "status": "paid",
  "transaction_id": "4200001234567890",
  "paid_at": "2026-01-24T12:45:00Z",
  "expired_at": "2026-01-24T14:34:56Z",
  "created_at": "2026-01-24T12:34:56Z",
  "updated_at": "2026-01-24T12:45:00Z"
}
```

**实现步骤:**
1. 获取 JWT 中的 `subject.UserID`
2. 从 URL 参数获取 `order_no`
3. 调用 `rechargeService.GetOrder(ctx, orderNo)`
4. **验证订单所有权:** `order.UserID == subject.UserID`
5. 如果不匹配，返回 403 Forbidden
6. 转换为 DTO 响应并返回

#### 3. ListOrders - 获取订单列表
```go
GET /api/v1/recharge/orders?page=1&page_size=20
Authorization: Bearer <token>

Response 200:
{
  "orders": [...],
  "pagination": {
    "total": 50,
    "page": 1,
    "page_size": 20,
    "pages": 3
  }
}
```

**实现步骤:**
1. 获取 JWT 中的 `subject.UserID`
2. 绑定查询参数到 `dto.PaginationRequest`
3. 设置默认值: `page=1`, `page_size=20`
4. 构建 `pagination.PaginationParams`
5. 调用 `rechargeService.GetUserOrders()`
6. 转换订单列表和分页信息为 DTO
7. 返回响应

#### 4. CancelOrder - 取消订单
```go
POST /api/v1/recharge/orders/:order_no/cancel
Authorization: Bearer <token>

Response 200:
{
  "message": "Order cancelled successfully"
}
```

**实现步骤:**
1. 获取 JWT 中的 `subject.UserID`
2. 从 URL 参数获取 `order_no`
3. **先查询订单验证所有权**
4. 调用 `rechargeService.CancelOrder(ctx, orderNo)`
5. 返回成功消息

#### 5. QueryOrderStatus - 查询订单状态（轮询接口）
```go
GET /api/v1/recharge/orders/:order_no/status
Authorization: Bearer <token>

Response 200:
{
  "order_no": "R20260124123456789",
  "status": "paid",
  "paid_at": "2026-01-24T12:45:00Z",
  "transaction_id": "4200001234567890"
}
```

**用途:** 前端轮询查询订单支付状态

**实现步骤:**
1. 获取 JWT 中的 `subject.UserID`
2. 从 URL 参数获取 `order_no`
3. 调用 `rechargeService.GetOrder()`
4. 验证订单所有权
5. **只返回状态相关字段**（减少响应大小）

#### 6. ListBalanceLogs - 获取余额变动记录
```go
GET /api/v1/recharge/balance-logs?page=1&page_size=20
Authorization: Bearer <token>

Response 200:
{
  "logs": [
    {
      "id": 123,
      "user_id": 456,
      "change_type": "recharge",
      "amount": 100.00,
      "balance_before": 50.00,
      "balance_after": 150.00,
      "description": "充值到账",
      "related_order_no": "R20260124123456789",
      "created_at": "2026-01-24T12:45:00Z"
    }
  ],
  "pagination": {...}
}
```

**实现步骤:**
1. 获取 JWT 中的 `subject.UserID`
2. 绑定分页参数
3. 调用 `balanceService.GetBalanceLogs()`
4. 转换为 DTO 响应

#### 7. GetRechargeConfig - 获取充值配置
```go
GET /api/v1/recharge/config

Response 200:
{
  "enabled": true,
  "min_amount": 1.00,
  "max_amount": 10000.00,
  "order_expire_minutes": 120
}
```

**用途:** 前端显示充值金额限制

**实现步骤:**
1. 调用 `rechargeService.GetConfig()`
2. 直接返回配置信息

**注意:** 此接口需要认证但无需验证用户权限

---

### 5.3 创建 WebhookHandler ✅

**文件:** `internal/handler/webhook_handler.go`

**依赖注入:**
```go
type WebhookHandler struct {
    paymentService *service.PaymentService
}

func NewWebhookHandler(paymentService *service.PaymentService) *WebhookHandler
```

**实现接口:**

#### HandleWeChatPayment - 处理微信支付回调
```go
POST /api/v1/webhook/wechat/payment
Content-Type: application/json
Wechatpay-Signature: ...
Wechatpay-Timestamp: ...
Wechatpay-Nonce: ...
Wechatpay-Serial: ...

Request Body: (微信支付加密的通知数据)

Success Response 200:
{
  "code": "SUCCESS",
  "message": ""
}

Error Response 500:
{
  "code": "FAIL",
  "message": "具体错误信息"
}
```

**实现步骤:**
1. 调用 `paymentService.ProcessCallback(ctx, c.Request)`
2. Service 层会处理:
   - 签名验证
   - 解密通知数据
   - 验证订单状态
   - 幂等性检查
   - 充值到账
   - 记录回调日志
3. **成功时返回 {"code":"SUCCESS"}**
4. **失败时返回 {"code":"FAIL", "message":"..."}**

**重要提示:**
- **不需要认证中间件** (微信通过签名验证)
- 必须返回微信要求的格式
- 失败时微信会重试（最多10次）
- 返回 SUCCESS 后微信不再重试

---

### 5.4 更新 Handler 集成 ⚠️

#### 修改 `internal/handler/handler.go`
```go
type Handlers struct {
    Auth          *AuthHandler
    User          *UserHandler
    APIKey        *APIKeyHandler
    Usage         *UsageHandler
    Redeem        *RedeemHandler
    Recharge      *RechargeHandler    // 新增
    Webhook       *WebhookHandler     // 新增
    Subscription  *SubscriptionHandler
    Admin         *AdminHandlers
    Gateway       *GatewayHandler
    OpenAIGateway *OpenAIGatewayHandler
    Setting       *SettingHandler
}
```

#### 修改 `internal/handler/wire.go`

**添加到 ProvideHandlers 参数:**
```go
func ProvideHandlers(
    authHandler *AuthHandler,
    userHandler *UserHandler,
    apiKeyHandler *APIKeyHandler,
    usageHandler *UsageHandler,
    redeemHandler *RedeemHandler,
    rechargeHandler *RechargeHandler,        // 新增
    webhookHandler *WebhookHandler,          // 新增
    subscriptionHandler *SubscriptionHandler,
    adminHandlers *AdminHandlers,
    gatewayHandler *GatewayHandler,
    openaiGatewayHandler *OpenAIGatewayHandler,
    settingHandler *SettingHandler,
) *Handlers {
    return &Handlers{
        Auth:          authHandler,
        User:          userHandler,
        APIKey:        apiKeyHandler,
        Usage:         usageHandler,
        Redeem:        redeemHandler,
        Recharge:      rechargeHandler,      // 新增
        Webhook:       webhookHandler,        // 新增
        Subscription:  subscriptionHandler,
        Admin:         adminHandlers,
        Gateway:       gatewayHandler,
        OpenAIGateway: openaiGatewayHandler,
        Setting:       settingHandler,
    }
}
```

**添加到 ProviderSet:**
```go
var ProviderSet = wire.NewSet(
    // Top-level handlers
    NewAuthHandler,
    NewUserHandler,
    NewAPIKeyHandler,
    NewUsageHandler,
    NewRedeemHandler,
    NewRechargeHandler,     // 新增
    NewWebhookHandler,      // 新增
    NewSubscriptionHandler,
    NewGatewayHandler,
    NewOpenAIGatewayHandler,
    ProvideSettingHandler,
    // ... admin handlers ...
)
```

---

### 5.5 创建路由配置 ✅

**文件:** `internal/server/routes/recharge.go`

```go
package routes

import (
    "github.com/Wei-Shaw/sub2api/internal/handler"
    "github.com/Wei-Shaw/sub2api/internal/server/middleware"
    "github.com/gin-gonic/gin"
)

// RegisterRechargeRoutes 注册充值相关路由（需要认证）
func RegisterRechargeRoutes(
    v1 *gin.RouterGroup,
    h *handler.Handlers,
    jwtAuth middleware.JWTAuthMiddleware,
) {
    authenticated := v1.Group("")
    authenticated.Use(gin.HandlerFunc(jwtAuth))
    {
        recharge := authenticated.Group("/recharge")
        {
            // 充值配置
            recharge.GET("/config", h.Recharge.GetRechargeConfig)

            // 订单管理
            recharge.POST("/orders", h.Recharge.CreateOrder)
            recharge.GET("/orders", h.Recharge.ListOrders)
            recharge.GET("/orders/:order_no", h.Recharge.GetOrder)
            recharge.GET("/orders/:order_no/status", h.Recharge.QueryOrderStatus)
            recharge.POST("/orders/:order_no/cancel", h.Recharge.CancelOrder)

            // 余额变动记录
            recharge.GET("/balance-logs", h.Recharge.ListBalanceLogs)
        }
    }
}

// RegisterWebhookRoutes 注册 webhook 路由（无需认证）
func RegisterWebhookRoutes(
    v1 *gin.RouterGroup,
    h *handler.Handlers,
) {
    webhook := v1.Group("/webhook")
    {
        // 微信支付回调
        webhook.POST("/wechat/payment", h.Webhook.HandleWeChatPayment)
    }
}
```

#### 修改 `internal/server/router.go`

```go
func registerRoutes(...) {
    // 通用路由
    routes.RegisterCommonRoutes(r)

    // API v1
    v1 := r.Group("/api/v1")

    // 注册各模块路由
    routes.RegisterAuthRoutes(v1, h, jwtAuth, redisClient)
    routes.RegisterUserRoutes(v1, h, jwtAuth)
    routes.RegisterRechargeRoutes(v1, h, jwtAuth)  // 新增
    routes.RegisterWebhookRoutes(v1, h)             // 新增
    routes.RegisterAdminRoutes(v1, h, adminAuth)
    routes.RegisterGatewayRoutes(r, h, apiKeyAuth, apiKeyService, subscriptionService, opsService, cfg)
}
```

---

## 🔧 依赖注入配置（关键步骤）

### 问题：循环依赖

**循环关系:**
```
RechargeService → PaymentService → RechargeService
                                    (通过 SetRechargeService)
```

### 解决方案：使用组合 Provider

#### 步骤 1: 创建 wechatpay Provider

**文件:** `internal/pkg/wechatpay/wire.go`

```go
package wechatpay

import (
    "fmt"
    "github.com/Wei-Shaw/sub2api/internal/config"
    "github.com/google/wire"
)

// ProvideClient 创建微信支付客户端
func ProvideClient(cfg *config.Config) (*Client, error) {
    if !cfg.WeChatPay.Enabled {
        // 返回 nil 表示功能未启用
        return nil, nil
    }

    clientConfig := Config{
        AppID:      cfg.WeChatPay.AppID,
        MchID:      cfg.WeChatPay.MchID,
        APIKey:     cfg.WeChatPay.APIKey,
        SerialNo:   cfg.WeChatPay.SerialNo,
        PrivateKey: cfg.WeChatPay.PrivateKeyPath,
        NotifyURL:  cfg.WeChatPay.NotifyURL,
    }

    client, err := NewClient(clientConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to create wechat pay client: %w", err)
    }

    return client, nil
}

// ProviderSet is the Wire provider set for wechatpay
var ProviderSet = wire.NewSet(
    ProvideClient,
)
```

#### 步骤 2: 修改 Service Wire 配置

**文件:** `internal/service/wire.go`

**添加 import:**
```go
import (
    "context"
    "database/sql"
    "time"

    "github.com/Wei-Shaw/sub2api/internal/config"
    "github.com/Wei-Shaw/sub2api/internal/pkg/wechatpay"  // 新增
    "github.com/google/wire"
    "github.com/redis/go-redis/v9"
)
```

**添加 Provider 函数:**
```go
// ProvideRechargeConfig 提供充值配置
func ProvideRechargeConfig(cfg *config.Config) RechargeConfig {
    return RechargeConfig{
        Enabled:            cfg.Recharge.Enabled,
        MinAmount:          cfg.Recharge.MinAmount,
        MaxAmount:          cfg.Recharge.MaxAmount,
        DefaultAmounts:     cfg.Recharge.DefaultAmounts,
        OrderExpireMinutes: cfg.Recharge.OrderExpireMinutes,
    }
}

// RechargeAndPaymentServices 包含充值和支付服务（解决循环依赖）
type RechargeAndPaymentServices struct {
    RechargeService *RechargeService
    PaymentService  *PaymentService
}

// ProvideRechargeAndPaymentServices 创建充值和支付服务（处理循环依赖）
func ProvideRechargeAndPaymentServices(
    rechargeOrderRepo RechargeOrderRepository,
    balanceService *BalanceService,
    wechatPayClient *wechatpay.Client,
    callbackRepo PaymentCallbackRepository,
    config RechargeConfig,
) *RechargeAndPaymentServices {
    // 先创建 PaymentService
    paymentService := NewPaymentService(wechatPayClient, callbackRepo)

    // 创建 RechargeService
    rechargeService := NewRechargeService(
        rechargeOrderRepo,
        paymentService,
        balanceService,
        config,
    )

    // 设置反向引用（解决循环依赖）
    paymentService.SetRechargeService(rechargeService)

    return &RechargeAndPaymentServices{
        RechargeService: rechargeService,
        PaymentService:  paymentService,
    }
}

// ProvideRechargeService 从组合服务中提取 RechargeService
func ProvideRechargeService(services *RechargeAndPaymentServices) *RechargeService {
    return services.RechargeService
}

// ProvidePaymentService 从组合服务中提取 PaymentService
func ProvidePaymentService(services *RechargeAndPaymentServices) *PaymentService {
    return services.PaymentService
}
```

**添加到 ProviderSet:**
```go
var ProviderSet = wire.NewSet(
    // ... existing services ...

    // Recharge and payment services
    NewBalanceService,
    ProvideRechargeConfig,
    ProvideRechargeAndPaymentServices,
    ProvideRechargeService,
    ProvidePaymentService,
)
```

#### 步骤 3: 修改主 Wire 配置

**文件:** `cmd/server/wire.go`

**添加 import:**
```go
import (
    // ... existing imports ...
    "github.com/Wei-Shaw/sub2api/internal/pkg/wechatpay"  // 新增
)
```

**添加到 wire.Build:**
```go
func initializeApplication(buildInfo handler.BuildInfo) (*Application, error) {
    wire.Build(
        // Infrastructure layer ProviderSets
        config.ProviderSet,
        wechatpay.ProviderSet,  // 新增

        // Business layer ProviderSets
        repository.ProviderSet,
        service.ProviderSet,
        middleware.ProviderSet,
        handler.ProviderSet,

        // Server layer ProviderSet
        server.ProviderSet,

        // BuildInfo provider
        provideServiceBuildInfo,

        // Cleanup function provider
        provideCleanup,

        // Application struct
        wire.Struct(new(Application), "Server", "Cleanup"),
    )
    return nil, nil
}
```

---

## 🧪 测试验证

### 编译验证

```bash
# 1. 生成 Wire 代码
cd cmd/server
~/go/bin/wire

# 2. 编译 handler 层
go build -v ./internal/handler/...

# 3. 编译整个项目
go build -v ./...

# 4. 运行项目
go run cmd/server/main.go
```

### 手动测试

#### 1. 测试充值配置接口
```bash
curl -X GET http://localhost:8080/api/v1/recharge/config \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**期望响应:**
```json
{
  "enabled": true,
  "min_amount": 1.00,
  "max_amount": 10000.00,
  "order_expire_minutes": 120
}
```

#### 2. 测试创建充值订单 (Native 支付)
```bash
curl -X POST http://localhost:8080/api/v1/recharge/orders \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.00,
    "payment_channel": "native"
  }'
```

**期望响应:**
```json
{
  "order_no": "R20260124123456789",
  "amount": 100.00,
  "expired_at": "2026-01-24T14:34:56Z",
  "code_url": "weixin://wxpay/bizpayurl?pr=xxxxxxx"
}
```

#### 3. 测试查询订单列表
```bash
curl -X GET "http://localhost:8080/api/v1/recharge/orders?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 4. 测试查询订单状态
```bash
curl -X GET http://localhost:8080/api/v1/recharge/orders/R20260124123456789/status \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 5. 测试余额变动记录
```bash
curl -X GET "http://localhost:8080/api/v1/recharge/balance-logs?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## ⚠️ 常见问题和解决方案

### 问题 1: Wire 找不到 Repository Provider

**错误:**
```
no provider found for github.com/Wei-Shaw/sub2api/internal/service.RechargeOrderRepository
```

**解决:** 确保在 `internal/repository/wire.go` 的 `ProviderSet` 中添加了:
```go
NewRechargeOrderRepository,
NewBalanceLogRepository,
NewPaymentCallbackRepository,
```

### 问题 2: 循环依赖错误

**错误:**
```
cycle detected in dependency graph
```

**解决:** 使用本文档中的 `RechargeAndPaymentServices` 组合模式。

### 问题 3: 编译错误 - 方法不存在

**错误:**
```
h.rechargeService.GetConfig undefined
```

**原因:** Service 层方法未实现

**解决:** 检查 Service 层是否添加了 `GetConfig()` 方法：
```go
func (s *RechargeService) GetConfig() RechargeConfig {
    return s.config
}
```

### 问题 4: 微信支付回调失败

**可能原因:**
1. 签名验证失败 → 检查证书配置
2. 重复处理 → 检查幂等性逻辑
3. 订单不存在 → 检查订单号匹配

**调试步骤:**
1. 查看 `payment_callbacks` 表记录
2. 检查 `signature_valid` 字段
3. 查看 `process_status` 和 `process_message`

---

## ✅ 完成标准

阶段 5 完成需满足以下条件：

- [ ] 所有 DTO 结构定义正确
- [ ] RechargeHandler 7 个接口全部实现
- [ ] WebhookHandler 实现并正确返回微信要求的格式
- [ ] Handler Wire 配置正确
- [ ] 路由配置完成
- [ ] 依赖注入配置正确（解决循环依赖）
- [ ] Wire 代码生成成功
- [ ] 项目编译通过
- [ ] 所有接口手动测试通过
- [ ] 微信支付回调测试通过（使用微信支付模拟器）

---

## 📝 代码审查清单

### Handler 层
- [ ] 所有 Handler 方法都有清晰的注释说明接口路径和用途
- [ ] JWT 认证正确获取 `subject.UserID`
- [ ] 订单所有权验证正确实现
- [ ] 错误处理使用统一的 `response.ErrorFrom()`
- [ ] 成功响应使用 `response.Success()`
- [ ] 分页参数有默认值处理

### DTO 层
- [ ] 所有字段都有正确的 JSON tag
- [ ] 可选字段使用指针类型
- [ ] 使用 `omitempty` 标记可选字段
- [ ] Request DTO 使用 binding tag 验证

### 路由层
- [ ] 充值接口都在认证中间件后面
- [ ] Webhook 接口没有认证中间件
- [ ] 路由分组清晰合理

### 依赖注入
- [ ] Wire Provider 函数签名正确
- [ ] 循环依赖正确解决
- [ ] ProviderSet 包含所有新增的 Provider

---

## 下一步：阶段 6

完成阶段 5 后，继续 **阶段 6：依赖注入最终整合和测试**

主要任务：
- 运行 Wire 生成最终代码
- 集成测试所有接口
- 修复遗留问题
- 准备进入前端开发

---

**文档版本:** v1.0
**最后更新:** 2026-01-24
**状态:** 待开始
