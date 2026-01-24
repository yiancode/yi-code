# Story 2.3: 充值金额范围验证

Status: ready-for-dev

## Story

**作为** 系统
**我希望** 验证用户输入的充值金额在允许范围内
**以便** 防止无效或异常金额的订单

## Acceptance Criteria

- [ ] AC1: 前端验证：金额 ≥ min_amount 且 ≤ max_amount
- [ ] AC2: 后端验证：金额范围校验
- [ ] AC3: 金额不在范围内时显示错误提示
- [ ] AC4: 提交按钮在金额无效时禁用
- [ ] AC5: 错误提示明确说明允许范围

## Tasks / Subtasks

- [ ] Task 1: 实现前端金额验证逻辑
- [ ] Task 2: 实现后端金额验证逻辑
- [ ] Task 3: 实现错误提示展示

## Dev Notes

### 前端验证

前端 computed 属性计算金额有效性

### 后端验证

后端在创建订单时校验金额范围

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-2.3]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
