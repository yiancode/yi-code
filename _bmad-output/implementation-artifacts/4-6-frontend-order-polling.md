# Story 4.6: 前端订单状态轮询

Status: ready-for-dev

## Story

**作为** 前端应用
**我希望** 定期轮询订单状态
**以便** 在用户支付成功后及时跳转

## Acceptance Criteria

- [ ] AC1: 支付中页面每3秒查询一次订单状态
- [ ] AC2: 最多轮询40次（共2分钟）
- [ ] AC3: 状态变为 `paid` 时跳转到成功页面
- [ ] AC4: 状态变为 `failed` 或 `expired` 时跳转到失败页面
- [ ] AC5: 轮询期间页面显示loading指示器
- [ ] AC6: 页面离开时停止轮询

## Tasks / Subtasks

- [ ] Task 1: 实现轮询逻辑（setInterval）
- [ ] Task 2: 实现状态判断和页面跳转
- [ ] Task 3: 实现组件卸载时清理定时器

## Dev Notes

### 轮询实现

```typescript
const pollInterval = setInterval(() => {
  // 查询订单状态
}, 3000);

onUnmounted(() => {
  clearInterval(pollInterval);
});
```

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-4.6]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
