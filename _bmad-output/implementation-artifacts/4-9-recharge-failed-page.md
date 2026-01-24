# Story 4.9: 充值失败页面

Status: ready-for-dev

## Story

**作为** 普通用户
**我希望** 支付失败或过期后看到失败页面
**以便** 了解失败原因并可以重试

## Acceptance Criteria

- [ ] AC1: 显示失败图标和「充值失败」标题
- [ ] AC2: 显示失败原因（支付失败/订单过期/用户取消）
- [ ] AC3: 提供「重新充值」按钮
- [ ] AC4: 提供「查看充值记录」链接

## Tasks / Subtasks

- [ ] Task 1: 创建 `RechargeFailedView.vue` 页面
- [ ] Task 2: 实现失败页面布局（符合 UX 设计）
- [ ] Task 3: 根据订单状态显示不同失败原因

## Dev Notes

### 页面路径

前端页面路径：`src/views/user/recharge/RechargeFailedView.vue`

### 失败原因映射

- `failed` → 支付失败
- `expired` → 订单已过期
- 用户主动取消 → 已取消

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-4.9]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
