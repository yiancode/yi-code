# Story 2.1: 充值页面与余额展示

Status: ready-for-dev

## Story

**作为** 普通用户
**我希望** 进入充值页面时看到当前账户余额
**以便** 了解当前余额并决定充值金额

## Acceptance Criteria

- [ ] AC1: 充值页面顶部显示当前账户余额（格式：¥XX.XX）
- [ ] AC2: 余额从用户信息接口获取
- [ ] AC3: 余额显示实时刷新（进入页面时重新获取）
- [ ] AC4: 页面加载时显示loading状态

## Tasks / Subtasks

- [ ] Task 1: 创建 `RechargeView.vue` 页面
- [ ] Task 2: 实现余额展示组件
- [ ] Task 3: 调用用户信息接口获取余额
- [ ] Task 4: 实现页面加载状态

## Dev Notes

### 接口复用

复用 `/api/v1/user/info` 接口获取余额（用户表已有 balance 字段）

### 页面路径

前端页面路径：`src/views/user/recharge/RechargeView.vue`

### UX 设计

参考 `_bmad-output/planning-artifacts/ux-design.md` 的充值首页设计

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-2.1]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
