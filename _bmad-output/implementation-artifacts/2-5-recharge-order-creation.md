# Story 2.5: 充值订单创建

Status: ready-for-dev

## Story

**作为** 系统
**我希望** 创建充值订单并生成唯一订单号
**以便** 追踪支付流程和后续对账

## Acceptance Criteria

- [ ] AC1: 订单号格式：RECH + 年月日时分秒 + 10位随机字符串（如：RECH20260124150000AbCd1234Ef）
- [ ] AC2: 订单号全局唯一（数据库唯一索引）
- [ ] AC3: 订单初始状态为 `pending`
- [ ] AC4: 记录：user_id, amount, payment_method, payment_channel
- [ ] AC5: 设置订单过期时间（created_at + expire_minutes）
- [ ] AC6: 订单创建时间精确到毫秒

## Tasks / Subtasks

- [ ] Task 1: 创建 `backend/ent/schema/recharge_order.go` Schema
- [ ] Task 2: 实现订单号生成算法
- [ ] Task 3: 实现订单创建 Service 方法
- [ ] Task 4: 创建订单 Handler

## Dev Notes

### 数据库表

使用 `recharge_orders` 表存储

### 订单号生成

订单号生成使用时间戳+随机串

```go
func GenerateOrderNo() string {
    return fmt.Sprintf("RECH%s%s", 
        time.Now().Format("20060102150405"),
        randomString(10))
}
```

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-2.5]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
