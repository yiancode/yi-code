# Story 3.5: 充值成功通知

Status: ready-for-dev

## Story

**作为** 普通用户
**我希望** 充值成功后收到站内信通知
**以便** 及时了解充值结果

## Acceptance Criteria

- [ ] AC1: 支付成功后异步发送站内信
- [ ] AC2: 通知内容：充值金额、订单号、到账时间、当前余额
- [ ] AC3: 使用goroutine异步发送，不阻塞回调响应
- [ ] AC4: 发送失败时记录错误日志但不影响主流程

## Tasks / Subtasks

- [ ] Task 1: 创建通知内容模板
- [ ] Task 2: 实现异步发送逻辑
- [ ] Task 3: 调用现有通知服务

## Dev Notes

### 异步发送

使用 goroutine 或 channel 异步处理

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-3.5]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
