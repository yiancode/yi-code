# Story 4.3: 取消未支付订单

Status: ready-for-dev

## Story

**作为** 普通用户
**我希望** 取消未支付的充值订单
**以便** 放弃当前订单重新充值

## Acceptance Criteria

- [ ] AC1: POST `/api/v1/recharge/orders/:order_no/cancel` 接口
- [ ] AC2: 只能取消状态为 `pending` 的订单
- [ ] AC3: 取消后状态变为 `failed`
- [ ] AC4: 记录取消原因：用户主动取消
- [ ] AC5: 已支付/已过期/已取消的订单不可再取消

## Tasks / Subtasks

- [ ] Task 1: 创建取消订单 Handler
- [ ] Task 2: 实现取消订单 Service 方法
- [ ] Task 3: 实现并发安全的状态更新

## Dev Notes

### 并发控制

乐观锁或条件更新防止并发问题

```sql
UPDATE recharge_orders SET status = 'failed' WHERE order_no = ? AND status = 'pending'
```

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-4.3]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
