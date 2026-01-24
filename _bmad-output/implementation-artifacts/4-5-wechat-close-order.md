# Story 4.5: 微信关闭订单API调用

Status: ready-for-dev

## Story

**作为** 系统
**我希望** 订单过期时调用微信关闭订单API
**以便** 防止用户在过期后仍能支付成功

## Acceptance Criteria

- [ ] AC1: 订单标记为expired后调用微信"关闭订单"API
- [ ] AC2: API调用失败时记录错误日志但不阻塞过期处理
- [ ] AC3: 使用goroutine异步调用，不阻塞主流程
- [ ] AC4: 记录关闭结果到订单备注字段

## Tasks / Subtasks

- [ ] Task 1: 实现微信关闭订单 Service 方法
- [ ] Task 2: 在过期处理中异步调用关单
- [ ] Task 3: 实现关单结果记录

## Dev Notes

### 微信SDK

调用微信支付SDK的关单方法

### 异步处理

使用 goroutine 异步处理避免超时

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-4.5]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
