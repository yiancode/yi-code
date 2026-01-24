# Story 3.1: 接收微信支付回调

Status: ready-for-dev

## Story

**作为** 系统
**我希望** 接收微信支付平台的异步回调通知
**以便** 获知用户支付结果

## Acceptance Criteria

- [ ] AC1: 提供 POST `/api/v1/webhook/wechat/payment` 回调接口
- [ ] AC2: 回调接口无需JWT认证
- [ ] AC3: 记录所有回调请求到 `payment_callbacks` 表
- [ ] AC4: 记录内容：请求头、请求体、接收时间
- [ ] AC5: 响应格式符合微信支付规范

## Tasks / Subtasks

- [ ] Task 1: 创建 `backend/ent/schema/payment_callback.go` Schema
- [ ] Task 2: 创建回调 Handler
- [ ] Task 3: 注册公开路由（无需认证）
- [ ] Task 4: 实现回调日志记录

## Dev Notes

### 数据库表

参考 `_bmad-output/planning-artifacts/epics.md#数据库需求` 的 payment_callbacks 表设计

### 响应格式

微信支付要求的响应格式：
- 成功: `{"code": "SUCCESS", "message": ""}`
- 失败: `{"code": "FAIL", "message": "失败原因"}`

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-3.1]
- [微信支付回调文档](https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_5.shtml)

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
