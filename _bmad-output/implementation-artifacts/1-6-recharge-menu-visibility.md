# Story 1.6: 充值菜单显示控制

Status: ready-for-dev

## Story

**作为** 普通用户
**我希望** 只有在充值功能启用时才看到充值入口
**以便** 避免看到无法使用的功能

## Acceptance Criteria

- [ ] AC1: 前端根据 `wechat_pay.enabled` 控制菜单显隐（简化为单一开关）
- [ ] AC2: GET `/api/v1/recharge/config` 接口返回 `enabled` 字段
- [ ] AC3: 启用时侧边栏显示「余额充值」菜单项
- [ ] AC4: 禁用时隐藏菜单项，直接访问路由跳转到首页

## Tasks / Subtasks

- [ ] Task 1: 后端 - 充值配置公开接口 (AC: 2)
  - [ ] 1.1 创建 `backend/internal/handler/recharge/handler.go`（RechargeHandler 结构体）
  - [ ] 1.2 实现 GET `/api/v1/recharge/config` 接口
  - [ ] 1.3 返回 enabled、min_amount、max_amount、default_amounts
  - [ ] 1.4 注册路由到 `backend/internal/server/routes/user.go`

- [ ] Task 2: 后端 - Wire 依赖注入配置 (AC: 2)
  - [ ] 2.1 在 `backend/internal/handler/wire.go` 注册 RechargeHandler
  - [ ] 2.2 在 Handlers 结构体中添加 Recharge 字段
  - [ ] 2.3 确保依赖注入链完整

- [ ] Task 3: 前端 - API 客户端 (AC: 1)
  - [ ] 3.1 创建 `frontend/src/api/recharge.ts`
  - [ ] 3.2 定义 RechargeConfig 接口和 getConfig 方法

- [ ] Task 4: 前端 - Pinia Store (AC: 1)
  - [ ] 4.1 创建 `frontend/src/stores/recharge.ts`
  - [ ] 4.2 实现 `fetchConfig()` 方法和 `isEnabled` getter
  - [ ] 4.3 在 `frontend/src/stores/index.ts` 导出

- [ ] Task 5: 前端 - 菜单控制 (AC: 3)
  - [ ] 5.1 在 AppSidebar.vue 中条件渲染充值菜单项
  - [ ] 5.2 使用 recharge store 的 isEnabled 状态

- [ ] Task 6: 前端 - 路由定义与守卫 (AC: 4)
  - [ ] 6.1 在 router/index.ts 添加充值相关路由
  - [ ] 6.2 实现充值路由守卫，禁用时跳转到首页

- [ ] Task 7: 前端 - i18n 国际化 (AC: 3)
  - [ ] 7.1 添加中文翻译 `nav.recharge`
  - [ ] 7.2 添加英文翻译 `nav.recharge`

- [ ] Task 8: 单元测试 (AC: 1-4)
  - [ ] 8.1 后端：测试 GetConfig 接口返回正确结构
  - [ ] 8.2 后端：测试 enabled 状态正确反映 WeChatPayService
  - [ ] 8.3 前端：测试菜单条件渲染逻辑

## Dev Notes

### 依赖关系

**前置条件**:
- Story 1.1（WeChatPayService 及 IsEnabled 方法已实现）
- Story 1.3（SettingService.GetRechargeSettings 已实现）

本 Story 实现前端与后端的联动，确保充值功能的可见性由配置统一控制。

### 设计原则

1. **单一开关控制**: 以 `wechat_pay.enabled` 作为充值功能的总开关
2. **无需认证**: 配置接口为公开接口，登录前即可获取（决定是否显示菜单）
3. **缓存优先**: 前端使用 localStorage 缓存减少闪烁
4. **优雅降级**: 接口失败时默认不显示充值入口

### 后端实现

#### 1. RechargeHandler 结构体

创建 `backend/internal/handler/recharge/handler.go`：

```go
package recharge

import (
    "github.com/Wei-Shaw/sub2api/internal/pkg/response"
    "github.com/Wei-Shaw/sub2api/internal/service"
    "github.com/gin-gonic/gin"
)

// RechargeHandler 充值相关接口处理器
type RechargeHandler struct {
    wechatPayService *service.WeChatPayService
    settingService   *service.SettingService
}

// NewRechargeHandler 创建充值处理器
func NewRechargeHandler(
    wechatPayService *service.WeChatPayService,
    settingService *service.SettingService,
) *RechargeHandler {
    return &RechargeHandler{
        wechatPayService: wechatPayService,
        settingService:   settingService,
    }
}

// RechargeConfigResponse 充值配置响应（公开接口）
type RechargeConfigResponse struct {
    Enabled        bool      `json:"enabled"`
    MinAmount      float64   `json:"min_amount"`
    MaxAmount      float64   `json:"max_amount"`
    DefaultAmounts []float64 `json:"default_amounts"`
}

// GetConfig 获取充值配置（无需认证）
// GET /api/v1/recharge/config
func (h *RechargeHandler) GetConfig(c *gin.Context) {
    // enabled 从 WeChatPayService 获取
    enabled := h.wechatPayService.IsEnabled()

    // 如果未启用，返回最小响应
    if !enabled {
        response.Success(c, RechargeConfigResponse{
            Enabled:        false,
            MinAmount:      0,
            MaxAmount:      0,
            DefaultAmounts: []float64{},
        })
        return
    }

    // 其他配置从 SettingService 获取
    settings, err := h.settingService.GetRechargeSettings(c.Request.Context())
    if err != nil {
        // 配置获取失败时仍返回 enabled 状态，使用默认值
        response.Success(c, RechargeConfigResponse{
            Enabled:        true,
            MinAmount:      1.0,
            MaxAmount:      1000.0,
            DefaultAmounts: []float64{10, 50, 100, 200, 500},
        })
        return
    }

    response.Success(c, RechargeConfigResponse{
        Enabled:        enabled,
        MinAmount:      settings.MinAmount,
        MaxAmount:      settings.MaxAmount,
        DefaultAmounts: settings.DefaultAmounts,
    })
}
```

#### 2. Wire 依赖注入

修改 `backend/internal/handler/wire.go`：

```go
// 在 import 中添加
import (
    // ...
    "github.com/Wei-Shaw/sub2api/internal/handler/recharge"
)

// 在 Handlers 结构体中添加
type Handlers struct {
    // ... 现有字段
    Recharge *recharge.RechargeHandler
}

// 在 ProvideHandlers 函数参数中添加
func ProvideHandlers(
    // ... 现有参数
    rechargeHandler *recharge.RechargeHandler,
) *Handlers {
    return &Handlers{
        // ... 现有字段
        Recharge: rechargeHandler,
    }
}

// 添加 Provider
var HandlerSet = wire.NewSet(
    // ... 现有 providers
    recharge.NewRechargeHandler,
    ProvideHandlers,
)
```

#### 3. 路由注册

修改 `backend/internal/server/routes/user.go` 或新建 `common.go`：

```go
// RegisterCommonRoutes 注册公共路由（无需认证）
func RegisterCommonRoutes(
    v1 *gin.RouterGroup,
    h *handler.Handlers,
) {
    // 充值配置（公开接口）
    rechargeGroup := v1.Group("/recharge")
    {
        rechargeGroup.GET("/config", h.Recharge.GetConfig)
    }
}
```

在 `backend/internal/server/routes/routes.go` 中调用：

```go
func RegisterRoutes(engine *gin.Engine, h *handler.Handlers, ...) {
    v1 := engine.Group("/api/v1")

    // 公共路由（无需认证）
    RegisterCommonRoutes(v1, h)

    // ... 其他路由注册
}
```

### 前端实现

#### 1. API 客户端

创建 `frontend/src/api/recharge.ts`：

```typescript
import request from '@/utils/request'

export interface RechargeConfig {
  enabled: boolean
  min_amount: number
  max_amount: number
  default_amounts: number[]
}

export const rechargeAPI = {
  /**
   * 获取充值配置（公开接口，无需认证）
   */
  getConfig(): Promise<RechargeConfig> {
    return request.get('/recharge/config')
  }
}

export default rechargeAPI
```

在 `frontend/src/api/index.ts` 中导出：

```typescript
// ... 现有导出
export { rechargeAPI } from './recharge'
```

#### 2. Pinia Store

创建 `frontend/src/stores/recharge.ts`：

```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { rechargeAPI, type RechargeConfig } from '@/api/recharge'

const CACHE_KEY = 'recharge_config_cached'
const CACHE_ENABLED_KEY = 'recharge_enabled_cached'

export const useRechargeStore = defineStore('recharge', () => {
  // 从缓存读取初始值，避免 UI 闪烁
  const readCachedEnabled = (): boolean => {
    try {
      const raw = localStorage.getItem(CACHE_ENABLED_KEY)
      return raw === 'true'
    } catch {
      return false
    }
  }

  const writeCachedEnabled = (value: boolean) => {
    try {
      localStorage.setItem(CACHE_ENABLED_KEY, value ? 'true' : 'false')
    } catch {
      // ignore localStorage failures
    }
  }

  const config = ref<RechargeConfig | null>(null)
  const loading = ref(false)
  const loaded = ref(false)
  const error = ref<string | null>(null)

  // 使用缓存的初始值
  const cachedEnabled = ref(readCachedEnabled())

  // 计算属性：是否启用充值
  const isEnabled = computed(() => {
    if (config.value !== null) {
      return config.value.enabled
    }
    // 未加载时使用缓存值
    return cachedEnabled.value
  })

  // 配置详情
  const minAmount = computed(() => config.value?.min_amount ?? 1)
  const maxAmount = computed(() => config.value?.max_amount ?? 1000)
  const defaultAmounts = computed(() => config.value?.default_amounts ?? [10, 50, 100, 200, 500])

  /**
   * 获取充值配置
   */
  async function fetchConfig(force = false): Promise<void> {
    if (loaded.value && !force) return
    if (loading.value) return

    loading.value = true
    error.value = null

    try {
      const data = await rechargeAPI.getConfig()
      config.value = data
      cachedEnabled.value = data.enabled
      writeCachedEnabled(data.enabled)
      loaded.value = true
    } catch (err) {
      console.error('[rechargeStore] Failed to fetch config:', err)
      error.value = err instanceof Error ? err.message : 'Unknown error'
      // 加载失败时保持缓存值，不改变 enabled 状态
      loaded.value = true
    } finally {
      loading.value = false
    }
  }

  /**
   * 重置状态（用于登出等场景）
   */
  function reset() {
    config.value = null
    loaded.value = false
    error.value = null
  }

  return {
    config,
    loading,
    loaded,
    error,
    isEnabled,
    minAmount,
    maxAmount,
    defaultAmounts,
    fetchConfig,
    reset
  }
})
```

在 `frontend/src/stores/index.ts` 中导出：

```typescript
// ... 现有导出
export { useRechargeStore } from './recharge'
```

#### 3. 菜单控制

修改 `frontend/src/components/layout/AppSidebar.vue`：

```vue
<script setup lang="ts">
import { computed, h, onMounted, ref, watch } from 'vue'
// ... 现有导入
import { useRechargeStore } from '@/stores'

// ... 现有代码

const rechargeStore = useRechargeStore()

// 钱包图标（用于充值菜单）
const WalletIcon = {
  render: () =>
    h(
      'svg',
      { fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor', 'stroke-width': '1.5' },
      [
        h('path', {
          'stroke-linecap': 'round',
          'stroke-linejoin': 'round',
          d: 'M21 12a2.25 2.25 0 00-2.25-2.25H15a3 3 0 11-6 0H5.25A2.25 2.25 0 003 12m18 0v6a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 18v-6m18 0V9M3 12V9m18 0a2.25 2.25 0 00-2.25-2.25H5.25A2.25 2.25 0 003 9m18 0V6a2.25 2.25 0 00-2.25-2.25H5.25A2.25 2.25 0 003 6v3'
        })
      ]
    )
}

// User navigation items (for regular users)
const userNavItems = computed(() => {
  const items = [
    { path: '/dashboard', label: t('nav.dashboard'), icon: DashboardIcon },
    { path: '/keys', label: t('nav.apiKeys'), icon: KeyIcon },
    { path: '/usage', label: t('nav.usage'), icon: ChartIcon, hideInSimpleMode: true },
    { path: '/subscriptions', label: t('nav.mySubscriptions'), icon: CreditCardIcon, hideInSimpleMode: true },
    { path: '/redeem', label: t('nav.redeem'), icon: GiftIcon, hideInSimpleMode: true },
    // 充值菜单项：根据 rechargeStore.isEnabled 控制
    ...(rechargeStore.isEnabled
      ? [{ path: '/recharge', label: t('nav.recharge'), icon: WalletIcon, hideInSimpleMode: true }]
      : []),
    { path: '/profile', label: t('nav.profile'), icon: UserIcon }
  ]
  return authStore.isSimpleMode ? items.filter(item => !item.hideInSimpleMode) : items
})

// Personal navigation items (for admin's "My Account" section)
const personalNavItems = computed(() => {
  const items = [
    { path: '/keys', label: t('nav.apiKeys'), icon: KeyIcon },
    { path: '/usage', label: t('nav.usage'), icon: ChartIcon, hideInSimpleMode: true },
    { path: '/subscriptions', label: t('nav.mySubscriptions'), icon: CreditCardIcon, hideInSimpleMode: true },
    { path: '/redeem', label: t('nav.redeem'), icon: GiftIcon, hideInSimpleMode: true },
    // 充值菜单项：根据 rechargeStore.isEnabled 控制
    ...(rechargeStore.isEnabled
      ? [{ path: '/recharge', label: t('nav.recharge'), icon: WalletIcon, hideInSimpleMode: true }]
      : []),
    { path: '/profile', label: t('nav.profile'), icon: UserIcon }
  ]
  return authStore.isSimpleMode ? items.filter(item => !item.hideInSimpleMode) : items
})

// 在 onMounted 中获取充值配置
onMounted(() => {
  if (isAdmin.value) {
    adminSettingsStore.fetch()
  }
  // 获取充值配置
  rechargeStore.fetchConfig()
})
</script>
```

#### 4. 路由定义与守卫

修改 `frontend/src/router/index.ts`：

```typescript
import { useRechargeStore } from '@/stores/recharge'

// 在 routes 数组中添加充值路由
const routes: RouteRecordRaw[] = [
  // ... 现有路由

  // ==================== User Routes ====================
  // ... 现有用户路由

  // 充值页面（需要认证）
  {
    path: '/recharge',
    name: 'Recharge',
    component: () => import('@/views/user/RechargeView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: 'Balance Recharge',
      titleKey: 'recharge.title',
      descriptionKey: 'recharge.description'
    }
  },

  // ... 其他路由
]

// 在 beforeEach 守卫中添加充值路由检查
router.beforeEach(async (to, _from, next) => {
  // ... 现有逻辑

  // 检查充值路由访问权限
  if (to.path.startsWith('/recharge')) {
    const rechargeStore = useRechargeStore()

    // 确保配置已加载
    if (!rechargeStore.loaded) {
      await rechargeStore.fetchConfig()
    }

    // 如果充值功能未启用，重定向到首页
    if (!rechargeStore.isEnabled) {
      next('/dashboard')
      return
    }
  }

  // All checks passed, allow navigation
  next()
})
```

#### 5. i18n 国际化

修改 `frontend/src/locales/zh-CN.json`：

```json
{
  "nav": {
    "recharge": "余额充值",
    // ... 其他翻译
  },
  "recharge": {
    "title": "余额充值",
    "description": "使用微信支付为账户充值"
  }
}
```

修改 `frontend/src/locales/en-US.json`：

```json
{
  "nav": {
    "recharge": "Balance Recharge",
    // ... 其他翻译
  },
  "recharge": {
    "title": "Balance Recharge",
    "description": "Top up your account balance via WeChat Pay"
  }
}
```

### 项目结构对齐

| 文件 | 作用 |
|------|------|
| `backend/internal/handler/recharge/handler.go` | 充值配置接口处理器 |
| `backend/internal/handler/wire.go` | Wire 依赖注入配置 |
| `backend/internal/server/routes/common.go` | 公共路由注册（或 user.go） |
| `frontend/src/api/recharge.ts` | API 客户端 |
| `frontend/src/stores/recharge.ts` | 充值状态管理 |
| `frontend/src/stores/index.ts` | Store 导出 |
| `frontend/src/components/layout/AppSidebar.vue` | 侧边栏菜单 |
| `frontend/src/router/index.ts` | 路由定义与守卫 |
| `frontend/src/locales/zh-CN.json` | 中文翻译 |
| `frontend/src/locales/en-US.json` | 英文翻译 |

### 测试场景

```go
// backend/internal/handler/recharge/handler_test.go
func TestGetConfig(t *testing.T) {
    // 场景1: 充值功能启用
    t.Run("RechargeEnabled", func(t *testing.T) {
        // Mock WeChatPayService.IsEnabled() 返回 true
        // Mock SettingService.GetRechargeSettings() 返回配置
        // 期望返回 enabled=true 和完整配置
    })

    // 场景2: 充值功能禁用
    t.Run("RechargeDisabled", func(t *testing.T) {
        // Mock WeChatPayService.IsEnabled() 返回 false
        // 期望返回 enabled=false 和空配置
    })

    // 场景3: 配置获取失败时使用默认值
    t.Run("SettingsError", func(t *testing.T) {
        // Mock WeChatPayService.IsEnabled() 返回 true
        // Mock SettingService.GetRechargeSettings() 返回错误
        // 期望返回 enabled=true 和默认配置值
    })
}
```

```typescript
// frontend/src/stores/__tests__/recharge.test.ts
describe('RechargeStore', () => {
  // 场景1: 成功获取配置
  it('should fetch config successfully', async () => {
    // Mock API 返回 enabled=true
    // 调用 fetchConfig()
    // 验证 isEnabled 为 true
  })

  // 场景2: API 失败时使用缓存
  it('should use cached value on API failure', async () => {
    // 设置 localStorage 缓存 enabled=true
    // Mock API 失败
    // 验证 isEnabled 仍为 true（来自缓存）
  })

  // 场景3: 禁用状态
  it('should return false when disabled', async () => {
    // Mock API 返回 enabled=false
    // 调用 fetchConfig()
    // 验证 isEnabled 为 false
  })
})
```

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-1.6] - 用户故事定义
- [Source: _bmad-output/implementation-artifacts/1-1-load-wechat-pay-config.md] - Story 1.1 定义
- [Source: _bmad-output/implementation-artifacts/1-3-admin-recharge-config.md] - Story 1.3 定义
- [Source: backend/internal/handler/admin/setting_handler.go] - 现有 Handler 模式参考
- [Source: backend/internal/server/routes/user.go] - 现有路由注册模式
- [Source: frontend/src/components/layout/AppSidebar.vue] - 现有侧边栏实现
- [Source: frontend/src/stores/adminSettings.ts] - 现有 Store 模式参考
- [Source: frontend/src/router/index.ts] - 现有路由守卫模式

### 设计决策

1. **公开接口**: `/api/v1/recharge/config` 无需认证，因为菜单显隐在用户登录前就需要决定
2. **LocalStorage 缓存**: 避免页面刷新时菜单闪烁
3. **优雅降级**: API 失败时使用缓存值，而非强制隐藏
4. **路由守卫**: 双重保护，即使用户手动访问 URL 也会被拦截
5. **简易模式兼容**: 充值菜单项在简易模式下隐藏（`hideInSimpleMode: true`）

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Debug Log References

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
