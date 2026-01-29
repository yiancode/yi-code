# Story 1.1: 加载微信支付敏感配置

Status: in-progress

## Story

**作为** 系统运维人员
**我希望** 系统启动时自动从config.yaml加载微信支付敏感配置
**以便** 安全地管理支付凭证，不暴露到前端或管理界面

## Acceptance Criteria

- [ ] AC1: 系统启动时从 `config.yaml` 读取 `wechat_pay` 配置节
- [ ] AC2: 加载以下敏感配置项：`enabled`, `app_id`, `mch_id`, `api_v3_key`, `cert_serial_no`, `private_key_path`, `notify_url`
- [ ] AC3: 配置加载失败时记录错误日志，但不影响系统其他功能启动
- [ ] AC4: 私钥文件路径验证：文件存在且权限正确（600）
- [ ] AC5: 使用官方微信支付Go SDK初始化客户端

## Tasks / Subtasks

- [ ] Task 1: 定义 WeChatPayConfig 结构体 (AC: 1, 2)
  - [ ] 1.1 在 `backend/internal/config/config.go` 添加 `WeChatPayConfig` 结构体
  - [ ] 1.2 在 `Config` 结构体中添加 `WeChatPay WeChatPayConfig` 字段

- [ ] Task 2: 配置默认值和验证 (AC: 1, 3, 4)
  - [ ] 2.1 在 `setDefaults()` 函数中添加 `wechat_pay.*` 默认值
  - [ ] 2.2 在 `Validate()` 方法中添加 WeChatPay 配置验证逻辑
  - [ ] 2.3 验证私钥文件存在性（仅当 enabled=true 时）

- [ ] Task 3: 更新配置示例文件 (AC: 2)
  - [ ] 3.1 在 `deploy/config.example.yaml` 添加 `wechat_pay` 配置示例

- [ ] Task 4: 创建微信支付服务 (AC: 5)
  - [ ] 4.1 创建 `backend/internal/service/wechat_pay_service.go`
  - [ ] 4.2 实现 `WeChatPayService` 结构体和 `NewWeChatPayService` 构造函数
  - [ ] 4.3 实现微信支付客户端初始化方法

- [ ] Task 5: Wire 依赖注入 (AC: 5)
  - [ ] 5.1 在 `backend/internal/service/wire.go` 注册 `WeChatPayService` 提供者
  - [ ] 5.2 运行 `go generate ./...` 重新生成 wire 代码

- [ ] Task 6: 单元测试 (AC: 1-5)
  - [ ] 6.1 添加 WeChatPayConfig 配置加载测试
  - [ ] 6.2 添加配置验证测试（有效/无效场景）

## Dev Notes

### 架构约束

本项目使用以下技术栈：
- **配置管理**: Viper + Mapstructure
- **依赖注入**: Google Wire
- **ORM**: Ent
- **微信支付SDK**: `github.com/wechatpay-apiv3/wechatpay-go`

### 配置结构体定义

在 `backend/internal/config/config.go` 中添加：

```go
// WeChatPayConfig 微信支付敏感配置
// 安全要求：所有敏感字段仅从config.yaml加载，不暴露到API
type WeChatPayConfig struct {
    Enabled        bool   `mapstructure:"enabled"`          // 是否启用微信支付
    AppID          string `mapstructure:"app_id"`           // 微信应用ID（公众号/小程序）
    MchID          string `mapstructure:"mch_id"`           // 商户号
    APIv3Key       string `mapstructure:"api_v3_key"`       // APIv3密钥（32字符）
    CertSerialNo   string `mapstructure:"cert_serial_no"`   // 商户证书序列号
    PrivateKeyPath string `mapstructure:"private_key_path"` // 商户私钥文件路径
    NotifyURL      string `mapstructure:"notify_url"`       // 支付回调地址
}
```

在 `Config` 结构体中添加字段（位置参考现有字段顺序）：

```go
type Config struct {
    // ... 现有字段
    WeChatPay    WeChatPayConfig            `mapstructure:"wechat_pay"`
    // ... 其他字段
}
```

### 默认值设置

在 `setDefaults()` 函数中添加：

```go
// WeChatPay defaults
viper.SetDefault("wechat_pay.enabled", false)
viper.SetDefault("wechat_pay.app_id", "")
viper.SetDefault("wechat_pay.mch_id", "")
viper.SetDefault("wechat_pay.api_v3_key", "")
viper.SetDefault("wechat_pay.cert_serial_no", "")
viper.SetDefault("wechat_pay.private_key_path", "")
viper.SetDefault("wechat_pay.notify_url", "")
```

### 验证逻辑

在 `Validate()` 方法中添加（参考 `LinuxDoConnectConfig` 验证模式）：

```go
// WeChatPay validation
if c.WeChatPay.Enabled {
    if strings.TrimSpace(c.WeChatPay.AppID) == "" {
        return fmt.Errorf("wechat_pay.app_id is required when enabled")
    }
    if strings.TrimSpace(c.WeChatPay.MchID) == "" {
        return fmt.Errorf("wechat_pay.mch_id is required when enabled")
    }
    if strings.TrimSpace(c.WeChatPay.APIv3Key) == "" {
        return fmt.Errorf("wechat_pay.api_v3_key is required when enabled")
    }
    if len(c.WeChatPay.APIv3Key) != 32 {
        return fmt.Errorf("wechat_pay.api_v3_key must be exactly 32 characters")
    }
    if strings.TrimSpace(c.WeChatPay.CertSerialNo) == "" {
        return fmt.Errorf("wechat_pay.cert_serial_no is required when enabled")
    }
    if strings.TrimSpace(c.WeChatPay.PrivateKeyPath) == "" {
        return fmt.Errorf("wechat_pay.private_key_path is required when enabled")
    }
    // 验证私钥文件存在
    if _, err := os.Stat(c.WeChatPay.PrivateKeyPath); os.IsNotExist(err) {
        return fmt.Errorf("wechat_pay.private_key_path file does not exist: %s", c.WeChatPay.PrivateKeyPath)
    }
    if strings.TrimSpace(c.WeChatPay.NotifyURL) == "" {
        return fmt.Errorf("wechat_pay.notify_url is required when enabled")
    }
    if err := ValidateAbsoluteHTTPURL(c.WeChatPay.NotifyURL); err != nil {
        return fmt.Errorf("wechat_pay.notify_url invalid: %w", err)
    }
}
```

### 微信支付服务实现

创建 `backend/internal/service/wechat_pay_service.go`：

```go
package service

import (
    "context"
    "crypto/rsa"
    "fmt"
    "sync"

    "github.com/wechatpay-apiv3/wechatpay-go/core"
    "github.com/wechatpay-apiv3/wechatpay-go/core/option"
    "github.com/wechatpay-apiv3/wechatpay-go/utils"
    "your-project/internal/config"
    "your-project/internal/log"
)

// WeChatPayService 微信支付服务
type WeChatPayService struct {
    cfg        *config.Config
    client     *core.Client
    privateKey *rsa.PrivateKey
    mu         sync.RWMutex
    initialized bool
}

// NewWeChatPayService 创建微信支付服务
func NewWeChatPayService(cfg *config.Config) *WeChatPayService {
    svc := &WeChatPayService{
        cfg: cfg,
    }

    // 仅当启用时初始化客户端
    if cfg.WeChatPay.Enabled {
        if err := svc.initClient(); err != nil {
            // 初始化失败记录日志，但不阻止服务启动
            log.Error("Failed to initialize WeChat Pay client", "error", err)
        }
    }

    return svc
}

// initClient 初始化微信支付客户端
func (s *WeChatPayService) initClient() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // 加载商户私钥
    privateKey, err := utils.LoadPrivateKeyWithPath(s.cfg.WeChatPay.PrivateKeyPath)
    if err != nil {
        return fmt.Errorf("load private key failed: %w", err)
    }
    s.privateKey = privateKey

    // 创建微信支付客户端
    ctx := context.Background()
    client, err := core.NewClient(
        ctx,
        option.WithMerchantCredential(
            s.cfg.WeChatPay.MchID,
            s.cfg.WeChatPay.CertSerialNo,
            privateKey,
        ),
        option.WithWechatPayAutoAuthCipher(s.cfg.WeChatPay.APIv3Key),
    )
    if err != nil {
        return fmt.Errorf("create wechat pay client failed: %w", err)
    }

    s.client = client
    s.initialized = true
    log.Info("WeChat Pay client initialized successfully")
    return nil
}

// IsEnabled 检查微信支付是否启用
func (s *WeChatPayService) IsEnabled() bool {
    return s.cfg.WeChatPay.Enabled && s.initialized
}

// GetClient 获取微信支付客户端（线程安全）
func (s *WeChatPayService) GetClient() (*core.Client, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    if !s.initialized {
        return nil, fmt.Errorf("wechat pay client not initialized")
    }
    return s.client, nil
}

// GetPrivateKey 获取商户私钥（用于JSAPI签名）
func (s *WeChatPayService) GetPrivateKey() (*rsa.PrivateKey, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    if s.privateKey == nil {
        return nil, fmt.Errorf("private key not loaded")
    }
    return s.privateKey, nil
}

// GetConfig 获取微信支付配置（只读）
func (s *WeChatPayService) GetConfig() config.WeChatPayConfig {
    return s.cfg.WeChatPay
}
```

### Wire 依赖注入

在 `backend/internal/service/wire.go` 的 `ProviderSet` 中添加：

```go
var ProviderSet = wire.NewSet(
    // ... 现有提供者
    NewWeChatPayService,
)
```

### 配置示例文件

在 `deploy/config.example.yaml` 中添加：

```yaml
# 微信支付配置（敏感信息，仅配置文件存储）
wechat_pay:
  enabled: false                # 是否启用微信支付
  app_id: ""                    # 微信应用ID（公众号/小程序）
  mch_id: ""                    # 商户号
  api_v3_key: ""                # APIv3密钥（32字符，在商户平台设置）
  cert_serial_no: ""            # 商户证书序列号
  private_key_path: ""          # 商户私钥文件路径，如：/path/to/apiclient_key.pem
  notify_url: ""                # 支付回调地址，如：https://yourdomain.com/api/v1/webhook/wechat/payment
```

### 项目结构对齐

| 文件 | 作用 |
|------|------|
| `backend/internal/config/config.go` | 配置结构体定义、默认值、验证 |
| `backend/internal/service/wechat_pay_service.go` | 微信支付服务（新建） |
| `backend/internal/service/wire.go` | 依赖注入注册 |
| `deploy/config.example.yaml` | 配置示例 |

### References

- [Source: docs/微信支付Go-SDK集成指南.md] - SDK使用指南
- [Source: _bmad-output/planning-artifacts/epics.md#Story-1.1] - 用户故事定义
- [Source: backend/internal/config/config.go] - 现有配置模式参考（LinuxDoConnectConfig）
- [微信支付官方Go SDK](https://github.com/wechatpay-apiv3/wechatpay-go)

### 安全注意事项

1. **不要日志输出敏感信息**：APIv3Key、私钥内容等不应出现在日志中
2. **私钥文件权限**：生产环境私钥文件权限应为 600
3. **环境变量支持**：支持 `WECHAT_PAY_*` 环境变量覆盖配置

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Debug Log References

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
