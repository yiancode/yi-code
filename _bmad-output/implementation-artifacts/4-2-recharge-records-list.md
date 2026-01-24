# Story 4.2: 充值记录列表

Status: ready-for-dev

## Story

**作为** 普通用户
**我希望** 查看我的充值记录列表
**以便** 了解历史充值情况

## Acceptance Criteria

- [ ] AC1: GET `/api/v1/recharge/orders` 接口
- [ ] AC2: 分页展示，每页10条，按创建时间倒序
- [ ] AC3: 支持筛选条件：状态（可选）、时间范围（可选）
- [ ] AC4: 返回字段：order_no, amount, status, created_at, paid_at
- [ ] AC5: 返回分页信息：total, page, page_size

## Tasks / Subtasks

- [ ] Task 1: 创建充值记录列表 Handler
- [ ] Task 2: 实现分页查询 Service 方法
- [ ] Task 3: 实现筛选条件处理
- [ ] Task 4: 创建前端充值记录页面

## Dev Notes

### 分页实现

分页使用 Offset + Limit

### 前端页面

前端页面路径：`src/views/user/recharge/RechargeRecordsView.vue`

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-4.2]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
