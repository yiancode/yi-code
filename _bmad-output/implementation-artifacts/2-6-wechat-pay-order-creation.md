# Story 2.6: 微信支付订单创建

Status: ready-for-dev

## Story

**作为** 系统
**我希望** 调用微信支付API创建支付订单
**以便** 获取支付参数供用户完成支付

## Acceptance Criteria

- [ ] AC1: 调用微信支付下单API（Native或JSAPI）
- [ ] AC2: 传递参数：商户订单号、金额（分）、商品描述、notify_url
- [ ] AC3: 保存返回的 prepay_id
- [ ] AC4: API调用超时设置为30秒
- [ ] AC5: API调用失败时记录详细错误日志
- [ ] AC6: 返回给前端：订单号、支付参数

## Tasks / Subtasks

- [ ] Task 1: 实现 `WeChatPayService.CreateOrder()` 方法
- [ ] Task 2: 实现金额转换（元 → 分）
- [ ] Task 3: 实现 Native 支付下单
- [ ] Task 4: 实现 JSAPI 支付下单
- [ ] Task 5: 实现错误处理和日志记录

## Dev Notes

### 微信支付SDK

使用微信支付Go SDK：`github.com/wechatpay-apiv3/wechatpay-go`

### 金额转换

金额转换：元 → 分（乘100）

```go
amountInCents := int64(amount * 100)
```

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-2.6]
- [Source: docs/微信支付Go-SDK集成指南.md]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
