# 充值功能完整开发计划（阶段 6-13）

## 📊 项目进度概览

| 阶段 | 名称 | 状态 | 描述 |
|------|------|------|------|
| ✅ 1 | 配置和基础设施 | 已完成 | config.yaml, go.mod |
| ✅ 2 | 数据库设计 | 已完成 | Ent Schema, SQL Migration |
| ✅ 3 | 微信支付 SDK | 已完成 | wechatpay 包实现 |
| ✅ 4 | Service 层 | 已完成 | RechargeService, PaymentService, BalanceService |
| ⏳ 5 | Handler 和路由 | 待开始 | RechargeHandler, WebhookHandler, Routes |
| ⏳ 6 | 依赖注入整合 | 待开始 | Wire 最终配置 |
| ⏳ 7 | 前端类型定义 | 待开始 | TypeScript 接口 |
| ⏳ 8 | 前端状态管理 | 待开始 | Pinia Store, Vue 组件 |
| ⏳ 9 | 前端路由集成 | 待开始 | 路由配置，导航 |
| ⏳ 10 | 管理后台配置 | 待开始 | 支付配置页面 |
| ⏳ 11 | 定时任务 | 待开始 | 订单过期处理 |
| ⏳ 12 | 集成测试 | 待开始 | 完整流程测试 |
| ⏳ 13 | 生产部署 | 待开始 | 部署和监控 |

---

## 阶段 6：依赖注入最终整合和测试

### 目标
完成 Wire 依赖注入配置，解决所有编译问题，确保系统可运行。

### 主要任务

#### 6.1 Repository Wire 配置
**文件:** `internal/repository/wire.go`

```go
var ProviderSet = wire.NewSet(
    // ... existing repos ...

    // Recharge repositories
    NewRechargeOrderRepository,
    NewBalanceLogRepository,
    NewPaymentCallbackRepository,
)
```

#### 6.2 运行 Wire 生成
```bash
cd cmd/server
~/go/bin/wire
```

#### 6.3 解决编译问题
- 检查所有 import 路径
- 确认接口实现完整
- 验证类型匹配

#### 6.4 启动测试
```bash
go run cmd/server/main.go
```

#### 6.5 API 集成测试
使用 Postman/curl 测试所有接口：
- ✅ 获取充值配置
- ✅ 创建充值订单
- ✅ 查询订单状态
- ✅ 查询订单列表
- ✅ 取消订单
- ✅ 查询余额日志
- ✅ 微信支付回调

### 完成标准
- [ ] Wire 生成成功
- [ ] 项目编译通过
- [ ] 服务启动成功
- [ ] 所有 API 接口测试通过
- [ ] 微信支付回调测试通过

---

## 阶段 7：前端 API 客户端和类型定义

### 目标
为前端提供类型安全的 API 调用接口。

### 主要任务

#### 7.1 TypeScript 类型定义
**文件:** `frontend/src/types/recharge.ts`

```typescript
// 充值订单状态
export type OrderStatus = 'pending' | 'paid' | 'cancelled' | 'expired'

// 支付渠道
export type PaymentChannel = 'jsapi' | 'native'

// 创建订单请求
export interface CreateRechargeOrderRequest {
  amount: number
  payment_channel: PaymentChannel
  openid?: string
}

// 创建订单响应
export interface CreateRechargeOrderResponse {
  order_no: string
  amount: number
  expired_at: string
  prepay_id?: string
  code_url?: string
  jsapi_params?: WeChatJSAPIParams
}

// 微信 JSAPI 支付参数
export interface WeChatJSAPIParams {
  appId: string
  timeStamp: string
  nonceStr: string
  package: string
  signType: string
  paySign: string
}

// 充值订单
export interface RechargeOrder {
  id: number
  order_no: string
  user_id: number
  amount: number
  actual_amount?: number
  currency: string
  payment_method: string
  payment_channel?: string
  status: OrderStatus
  transaction_id?: string
  paid_at?: string
  expired_at: string
  created_at: string
  updated_at: string
}

// 余额变动日志
export interface BalanceLog {
  id: number
  user_id: number
  change_type: string
  amount: number
  balance_before: number
  balance_after: number
  description: string
  related_order_no?: string
  created_at: string
}

// 分页响应
export interface PaginatedResponse<T> {
  data: T[]
  pagination: {
    total: number
    page: number
    page_size: number
    pages: number
  }
}

// 充值配置
export interface RechargeConfig {
  enabled: boolean
  min_amount: number
  max_amount: number
  order_expire_minutes: number
}
```

#### 7.2 API 客户端
**文件:** `frontend/src/api/recharge.ts`

```typescript
import { apiClient } from './client'
import type {
  CreateRechargeOrderRequest,
  CreateRechargeOrderResponse,
  RechargeOrder,
  BalanceLog,
  PaginatedResponse,
  RechargeConfig
} from '@/types/recharge'

export const rechargeApi = {
  // 获取充值配置
  getConfig(): Promise<RechargeConfig> {
    return apiClient.get('/recharge/config')
  },

  // 创建充值订单
  createOrder(data: CreateRechargeOrderRequest): Promise<CreateRechargeOrderResponse> {
    return apiClient.post('/recharge/orders', data)
  },

  // 获取订单详情
  getOrder(orderNo: string): Promise<RechargeOrder> {
    return apiClient.get(`/recharge/orders/${orderNo}`)
  },

  // 查询订单状态
  queryOrderStatus(orderNo: string): Promise<{
    order_no: string
    status: string
    paid_at?: string
    transaction_id?: string
  }> {
    return apiClient.get(`/recharge/orders/${orderNo}/status`)
  },

  // 获取订单列表
  listOrders(page: number = 1, pageSize: number = 20): Promise<{
    orders: RechargeOrder[]
    pagination: any
  }> {
    return apiClient.get('/recharge/orders', {
      params: { page, page_size: pageSize }
    })
  },

  // 取消订单
  cancelOrder(orderNo: string): Promise<{ message: string }> {
    return apiClient.post(`/recharge/orders/${orderNo}/cancel`)
  },

  // 获取余额变动记录
  listBalanceLogs(page: number = 1, pageSize: number = 20): Promise<{
    logs: BalanceLog[]
    pagination: any
  }> {
    return apiClient.get('/recharge/balance-logs', {
      params: { page, page_size: pageSize }
    })
  }
}
```

### 完成标准
- [ ] 所有类型定义完整
- [ ] API 客户端方法完整
- [ ] TypeScript 类型检查通过
- [ ] 与后端接口一致

---

## 阶段 8：前端状态管理和组件

### 目标
实现充值功能的前端交互界面。

### 8.1 Pinia Store

**文件:** `frontend/src/stores/recharge.ts`

```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { rechargeApi } from '@/api/recharge'
import type { RechargeOrder, RechargeConfig } from '@/types/recharge'

export const useRechargeStore = defineStore('recharge', () => {
  // State
  const config = ref<RechargeConfig | null>(null)
  const orders = ref<RechargeOrder[]>([])
  const currentOrder = ref<RechargeOrder | null>(null)
  const loading = ref(false)

  // Getters
  const isEnabled = computed(() => config.value?.enabled ?? false)

  // Actions
  async function fetchConfig() {
    try {
      config.value = await rechargeApi.getConfig()
    } catch (error) {
      console.error('Failed to fetch recharge config:', error)
    }
  }

  async function createOrder(amount: number, channel: 'jsapi' | 'native') {
    loading.value = true
    try {
      const result = await rechargeApi.createOrder({
        amount,
        payment_channel: channel
      })
      return result
    } finally {
      loading.value = false
    }
  }

  async function fetchOrders(page: number = 1) {
    loading.value = true
    try {
      const result = await rechargeApi.listOrders(page, 10)
      orders.value = result.orders
      return result
    } finally {
      loading.value = false
    }
  }

  async function pollOrderStatus(orderNo: string, maxAttempts: number = 60) {
    for (let i = 0; i < maxAttempts; i++) {
      const status = await rechargeApi.queryOrderStatus(orderNo)
      if (status.status === 'paid') {
        return status
      }
      if (status.status === 'cancelled' || status.status === 'expired') {
        throw new Error(`Order ${status.status}`)
      }
      await new Promise(resolve => setTimeout(resolve, 2000))
    }
    throw new Error('Payment timeout')
  }

  return {
    config,
    orders,
    currentOrder,
    loading,
    isEnabled,
    fetchConfig,
    createOrder,
    fetchOrders,
    pollOrderStatus
  }
})
```

### 8.2 Vue 组件

#### RechargeDialog.vue - 充值对话框
```vue
<template>
  <el-dialog
    v-model="visible"
    title="余额充值"
    width="500px"
  >
    <div v-if="!qrCodeUrl">
      <!-- 金额选择 -->
      <div class="amount-selector">
        <div
          v-for="amount in [10, 50, 100, 200, 500]"
          :key="amount"
          class="amount-option"
          :class="{ active: selectedAmount === amount }"
          @click="selectedAmount = amount"
        >
          ¥{{ amount }}
        </div>
      </div>

      <!-- 自定义金额 -->
      <el-input
        v-model.number="customAmount"
        type="number"
        placeholder="自定义金额"
        :min="config?.min_amount"
        :max="config?.max_amount"
      />

      <!-- 支付方式选择 -->
      <el-radio-group v-model="paymentChannel">
        <el-radio label="native">扫码支付</el-radio>
        <el-radio label="jsapi">JSAPI支付</el-radio>
      </el-radio-group>
    </div>

    <!-- 二维码显示 -->
    <div v-else class="qrcode-container">
      <qrcode-vue :value="qrCodeUrl" :size="200" />
      <p>请使用微信扫码支付</p>
      <p class="amount">¥{{ orderAmount }}</p>
      <p class="countdown">订单将在 {{ countdown }} 后过期</p>
    </div>

    <template #footer>
      <el-button v-if="!qrCodeUrl" @click="visible = false">取消</el-button>
      <el-button
        v-if="!qrCodeUrl"
        type="primary"
        :loading="loading"
        @click="handleCreateOrder"
      >
        立即充值
      </el-button>
      <el-button v-else @click="handleCancel">取消订单</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import QrcodeVue from 'qrcode.vue'
import { useRechargeStore } from '@/stores/recharge'

const rechargeStore = useRechargeStore()

const visible = ref(false)
const selectedAmount = ref(100)
const customAmount = ref<number>()
const paymentChannel = ref<'native' | 'jsapi'>('native')
const qrCodeUrl = ref('')
const orderNo = ref('')
const countdown = ref(120)
const loading = ref(false)

const orderAmount = computed(() => customAmount.value || selectedAmount.value)

let countdownTimer: number | null = null
let pollTimer: number | null = null

async function handleCreateOrder() {
  try {
    loading.value = true
    const result = await rechargeStore.createOrder(
      orderAmount.value,
      paymentChannel.value
    )

    if (result.code_url) {
      qrCodeUrl.value = result.code_url
      orderNo.value = result.order_no
      startCountdown()
      startPolling()
    } else if (result.jsapi_params) {
      // 调用微信 JSAPI 支付
      await invokeWeChatPay(result.jsapi_params)
    }
  } catch (error: any) {
    ElMessage.error(error.message || '创建订单失败')
  } finally {
    loading.value = false
  }
}

function startCountdown() {
  countdown.value = 120
  countdownTimer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      handleCancel()
    }
  }, 1000) as any
}

async function startPolling() {
  try {
    await rechargeStore.pollOrderStatus(orderNo.value)
    ElMessage.success('充值成功！')
    visible.value = false
    // 刷新用户余额
  } catch (error: any) {
    ElMessage.error(error.message || '支付失败')
  }
}

function handleCancel() {
  if (countdownTimer) clearInterval(countdownTimer)
  if (pollTimer) clearInterval(pollTimer)
  qrCodeUrl.value = ''
  orderNo.value = ''
}

watch(visible, (val) => {
  if (!val) {
    handleCancel()
  }
})

defineExpose({ open: () => { visible.value = true } })
</script>
```

#### RechargeHistory.vue - 充值历史
```vue
<template>
  <div class="recharge-history">
    <el-table :data="orders" stripe>
      <el-table-column prop="order_no" label="订单号" width="180" />
      <el-table-column prop="amount" label="金额" width="100">
        <template #default="{ row }">
          ¥{{ row.amount.toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button
            v-if="row.status === 'pending'"
            link
            type="primary"
            @click="handleContinuePay(row)"
          >
            继续支付
          </el-button>
          <el-button
            v-if="row.status === 'pending'"
            link
            type="danger"
            @click="handleCancel(row.order_no)"
          >
            取消
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="currentPage"
      :page-size="10"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="fetchOrders"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRechargeStore } from '@/stores/recharge'
import type { RechargeOrder } from '@/types/recharge'

const rechargeStore = useRechargeStore()
const orders = ref<RechargeOrder[]>([])
const currentPage = ref(1)
const total = ref(0)

async function fetchOrders() {
  const result = await rechargeStore.fetchOrders(currentPage.value)
  orders.value = result.orders
  total.value = result.pagination.total
}

function getStatusType(status: string) {
  const types: Record<string, any> = {
    pending: 'warning',
    paid: 'success',
    cancelled: 'info',
    expired: 'danger'
  }
  return types[status] || 'info'
}

function getStatusText(status: string) {
  const texts: Record<string, string> = {
    pending: '待支付',
    paid: '已支付',
    cancelled: '已取消',
    expired: '已过期'
  }
  return texts[status] || status
}

onMounted(() => {
  fetchOrders()
})
</script>
```

### 完成标准
- [ ] Pinia Store 实现完整
- [ ] 充值对话框组件完成
- [ ] 充值历史组件完成
- [ ] 二维码支付流程可用
- [ ] 订单状态轮询正常
- [ ] 用户体验流畅

---

## 阶段 9：前端路由和导航集成

### 目标
将充值功能集成到现有前端系统。

### 主要任务

#### 9.1 添加路由
**文件:** `frontend/src/router/index.ts`

```typescript
{
  path: '/recharge',
  name: 'Recharge',
  component: () => import('@/views/recharge/Index.vue'),
  meta: { requiresAuth: true }
}
```

#### 9.2 导航菜单
在用户中心添加"余额充值"入口

#### 9.3 集成到用户中心
- 在个人资料页显示当前余额
- 添加"充值"按钮打开充值对话框
- 显示最近充值记录

### 完成标准
- [ ] 路由配置完成
- [ ] 导航入口添加
- [ ] 用户中心集成
- [ ] 跳转和返回流畅

---

## 阶段 10：管理后台支付配置页面

### 目标
提供管理员配置微信支付参数的界面。

### 主要任务

#### 10.1 后端接口
```go
// GET /api/v1/admin/payment/config
// PUT /api/v1/admin/payment/config
```

#### 10.2 前端页面
- 微信支付配置表单
- 配置测试功能
- 证书文件上传

### 配置项
- AppID
- MchID (商户号)
- API Key
- 证书序列号
- 私钥文件路径
- 回调 URL
- 启用/禁用开关
- 充值金额限制

### 完成标准
- [ ] 配置接口实现
- [ ] 前端表单完成
- [ ] 配置验证功能
- [ ] 测试功能可用

---

## 阶段 11：定时任务和后台服务

### 目标
实现订单过期自动处理。

### 11.1 订单过期处理

**实现位置:** 已在 `RechargeService.ExpireOrders()` 中实现

**创建定时任务:**
```go
// internal/service/recharge_scheduler.go

type RechargeScheduler struct {
    rechargeService *RechargeService
}

func (s *RechargeScheduler) Start() {
    ticker := time.NewTicker(5 * time.Minute)
    go func() {
        for range ticker.C {
            count, err := s.rechargeService.ExpireOrders(context.Background())
            if err != nil {
                log.Printf("Failed to expire orders: %v", err)
            } else {
                log.Printf("Expired %d orders", count)
            }
        }
    }()
}
```

**集成到 main.go:**
```go
scheduler := service.NewRechargeScheduler(rechargeService)
scheduler.Start()
```

### 完成标准
- [ ] 定时任务实现
- [ ] 集成到启动流程
- [ ] 日志记录完善
- [ ] 错误处理健壮

---

## 阶段 12：集成测试和文档

### 目标
全面测试充值功能，编写使用文档。

### 12.1 测试用例

#### 后端测试
```bash
# 正常流程测试
1. 创建充值订单
2. 模拟支付回调
3. 验证余额增加
4. 验证余额日志

# 异常流程测试
1. 重复回调测试（幂等性）
2. 订单过期测试
3. 取消订单测试
4. 并发充值测试
5. 签名验证失败测试
```

#### 前端测试
```bash
1. 充值对话框打开/关闭
2. 金额选择和验证
3. 二维码生成和显示
4. 支付状态轮询
5. 充值历史列表
6. 余额显示更新
```

### 12.2 文档编写

#### API 文档
- 所有接口的请求/响应示例
- 错误码说明
- 认证要求

#### 用户手册
- 如何充值
- 支付方式说明
- 常见问题

#### 管理员手册
- 如何配置微信支付
- 如何查看充值记录
- 如何处理充值问题

### 完成标准
- [ ] 所有测试用例通过
- [ ] API 文档完整
- [ ] 用户手册完成
- [ ] 管理员手册完成

---

## 阶段 13：生产环境部署

### 目标
将充值功能部署到生产环境。

### 13.1 部署前检查

```bash
# 1. 数据库迁移
psql -h your-host -U your-user -d your-db -f migrations/047_add_recharge_tables.sql

# 2. 配置文件
cat > /path/to/config.yaml <<EOF
wechat_pay:
  enabled: true
  app_id: "wx0b35f0f5c31fb07e"
  mch_id: "1654437140"
  api_key: "cda95734abaaf26d61bc98882be5878b"
  serial_no: "335F98DBC9AFFE51174ADBC5D999F68563FD949E"
  private_key_path: "/path/to/apiclient_key.pem"
  notify_url: "https://code.ai80.vip/api/v1/webhook/wechat/payment"

recharge:
  enabled: true
  min_amount: 1.00
  max_amount: 10000.00
  default_amounts: [10, 50, 100, 200, 500]
  order_expire_minutes: 120
EOF

# 3. 证书文件
cp apiclient_key.pem /path/to/certs/
chmod 600 /path/to/certs/apiclient_key.pem

# 4. 编译
go build -o sub2api cmd/server/main.go

# 5. 启动
./sub2api
```

### 13.2 微信支付配置

登录微信支付商户平台：
1. 设置 API 密钥
2. 上传 API 证书
3. 配置回调 URL: `https://code.ai80.vip/api/v1/webhook/wechat/payment`
4. 测试回调连通性

### 13.3 监控配置

```bash
# 1. 日志监控
tail -f /path/to/logs/app.log | grep -i "recharge\|payment"

# 2. 数据库监控
SELECT
  status,
  COUNT(*) as count,
  SUM(amount) as total_amount
FROM recharge_orders
WHERE created_at > NOW() - INTERVAL '1 day'
GROUP BY status;

# 3. 告警规则
- 充值订单支付失败率 > 5%
- 订单过期率 > 10%
- 回调处理失败 > 3次/小时
```

### 13.4 回滚计划

如果出现问题：
```bash
# 1. 停止服务
systemctl stop sub2api

# 2. 回滚代码
git checkout <previous-version>

# 3. 数据库回滚
# 注意：不要删除充值相关表，只需停用功能
# 在 config.yaml 中设置：
wechat_pay:
  enabled: false
recharge:
  enabled: false

# 4. 重启服务
systemctl start sub2api
```

### 完成标准
- [ ] 数据库迁移成功
- [ ] 配置文件正确
- [ ] 证书文件部署
- [ ] 服务启动成功
- [ ] 微信支付配置完成
- [ ] 生产环境测试通过
- [ ] 监控和告警配置
- [ ] 回滚计划准备

---

## 🎯 总结

### 完整功能清单

**后端 (Go):**
- [x] 配置管理
- [x] 数据库 Schema
- [x] 微信支付 SDK
- [x] Service 业务逻辑
- [ ] Handler HTTP 接口
- [ ] 依赖注入配置
- [ ] 定时任务

**前端 (Vue):**
- [ ] TypeScript 类型
- [ ] API 客户端
- [ ] Pinia Store
- [ ] 充值组件
- [ ] 路由集成

**管理:**
- [ ] 配置页面
- [ ] 文档编写
- [ ] 测试验证
- [ ] 生产部署

### 风险和注意事项

1. **支付安全:** 确保签名验证正确实现
2. **幂等性:** 防止重复充值
3. **并发安全:** 使用数据库锁
4. **错误处理:** 完善的错误日志
5. **监控告警:** 及时发现问题

### 预估时间

- 阶段 5-6: 2-3 天（后端接口完成）
- 阶段 7-9: 3-4 天（前端开发）
- 阶段 10-11: 1-2 天（管理和定时任务）
- 阶段 12: 1-2 天（测试和文档）
- 阶段 13: 1 天（部署）

**总计:** 约 8-12 天

---

**文档版本:** v1.0
**最后更新:** 2026-01-24
**下一步:** 开始阶段 5 实现
