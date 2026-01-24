# Story 4.1: 查询单个订单详情

Status: ready-for-dev

## Story

**作为** 普通用户
**我希望** 查询某个充值订单的详细状态
**以便** 了解订单的支付进度

## Acceptance Criteria

- [ ] AC1: GET `/api/v1/recharge/orders/:order_no` 接口
- [ ] AC2: 返回字段：order_no, amount, status, payment_method, created_at, paid_at, expired_at
- [ ] AC3: 只能查询自己的订单
- [ ] AC4: 订单不存在时返回404

## Tasks / Subtasks

- [ ] Task 1: 创建订单查询 Handler
- [ ] Task 2: 实现订单详情 Service 方法
- [ ] Task 3: 添加用户权限校验（user_id 匹配）
- [ ] Task 4: 定义订单详情 DTO 响应结构

## Dev Notes

### 权限校验

查询时校验 user_id 匹配，使用 Ent 的 `Where` 条件

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-4.1]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
