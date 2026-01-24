# Story 4.7: 支付中页面倒计时

Status: ready-for-dev

## Story

**作为** 普通用户
**我希望** 看到订单过期倒计时
**以便** 了解剩余支付时间

## Acceptance Criteria

- [ ] AC1: 显示格式：XX分XX秒
- [ ] AC2: 每秒更新一次
- [ ] AC3: 倒计时归零时跳转到失败页面
- [ ] AC4: 倒计时与后端过期时间同步

## Tasks / Subtasks

- [ ] Task 1: 创建 `OrderCountdown.vue` 组件
- [ ] Task 2: 实现倒计时计算逻辑
- [ ] Task 3: 实现倒计时归零跳转

## Dev Notes

### 组件路径

前端组件路径：`src/components/user/recharge/OrderCountdown.vue`

### 计算逻辑

根据 expired_at 计算剩余时间

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-4.7]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
