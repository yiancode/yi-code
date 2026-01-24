# Story 4.4: 订单过期自动处理

Status: ready-for-dev

## Story

**作为** 系统
**我希望** 自动将超时未支付的订单标记为过期
**以便** 清理无效订单

## Acceptance Criteria

- [ ] AC1: 定时任务每分钟执行一次
- [ ] AC2: 扫描 `status = pending` 且 `expired_at < now()` 的订单
- [ ] AC3: 批量更新状态为 `expired`
- [ ] AC4: 每次最多处理100条
- [ ] AC5: 记录过期订单数量日志

## Tasks / Subtasks

- [ ] Task 1: 创建定时任务调度器
- [ ] Task 2: 实现过期订单扫描逻辑
- [ ] Task 3: 实现批量状态更新

## Dev Notes

### 定时任务

使用 cron 或 ticker 定时执行

### 批量更新

```sql
UPDATE recharge_orders SET status = 'expired' 
WHERE status = 'pending' AND expired_at < NOW() 
LIMIT 100
```

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-4.4]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
