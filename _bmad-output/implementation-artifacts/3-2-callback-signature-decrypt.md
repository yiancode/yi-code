# Story 3.2: 回调签名验证与数据解密

Status: ready-for-dev

## Story

**作为** 系统
**我希望** 验证回调签名并解密回调数据
**以便** 确保回调来自微信支付平台且数据未被篡改

## Acceptance Criteria

- [ ] AC1: 使用微信支付平台证书验证签名
- [ ] AC2: 验证请求头中的 Wechatpay-Timestamp, Wechatpay-Nonce, Wechatpay-Signature
- [ ] AC3: 拒绝超过5分钟的请求（防重放攻击）
- [ ] AC4: 使用 APIv3 密钥解密回调数据（AEAD_AES_256_GCM）
- [ ] AC5: 签名验证失败时返回 FAIL 并记录日志
- [ ] AC6: 更新 payment_callbacks 表的 signature_valid 字段

## Tasks / Subtasks

- [ ] Task 1: 使用微信支付 SDK 验签方法
- [ ] Task 2: 实现时间戳检查（5分钟内）
- [ ] Task 3: 实现回调数据解密
- [ ] Task 4: 更新回调日志的验签结果

## Dev Notes

### 微信支付 SDK 验签

使用 `github.com/wechatpay-apiv3/wechatpay-go/core/notify`

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-3.2]
- [Source: docs/微信支付Go-SDK集成指南.md]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
