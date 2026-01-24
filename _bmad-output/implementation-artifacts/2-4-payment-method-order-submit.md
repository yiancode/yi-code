# Story 2.4: 支付方式选择与订单发起

Status: ready-for-dev

## Story

**作为** 普通用户
**我希望** 选择支付方式并发起充值订单
**以便** 开始支付流程

## Acceptance Criteria

- [ ] AC1: 显示支付方式选择器（当前仅支持微信支付）
- [ ] AC2: 微信支付图标和名称显示
- [ ] AC3: 点击「立即充值」按钮发起订单
- [ ] AC4: 按钮点击后显示loading状态，防止重复提交
- [ ] AC5: 订单创建成功后跳转到支付页面
- [ ] AC6: 订单创建失败显示错误提示

## Tasks / Subtasks

- [ ] Task 1: 创建 `PaymentMethodSelector.vue` 组件
- [ ] Task 2: 实现订单创建 API 调用
- [ ] Task 3: 实现提交状态管理
- [ ] Task 4: 实现页面跳转逻辑

## Dev Notes

### 组件路径

前端组件路径：`src/components/user/recharge/PaymentMethodSelector.vue`

### API 调用

调用 POST `/api/v1/recharge/orders` 接口

### API 客户端

API客户端：`src/api/recharge.ts`

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-2.4]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
