# Story 2.8: JSAPI支付调起

Status: ready-for-dev

## Story

**作为** 普通用户（微信内）
**我希望** 在微信内直接调起支付
**以便** 在微信浏览器中便捷支付

## Acceptance Criteria

- [ ] AC1: JSAPI支付返回前端调起参数（appId, timeStamp, nonceStr, package, signType, paySign）
- [ ] AC2: 前端调用 WeixinJSBridge 或微信JS-SDK调起支付
- [ ] AC3: 支付成功后跳转到成功页面
- [ ] AC4: 支付失败或取消后显示相应提示
- [ ] AC5: 非微信环境下隐藏JSAPI支付选项

## Tasks / Subtasks

- [ ] Task 1: 实现后端 JSAPI 签名参数生成
- [ ] Task 2: 实现前端微信环境检测
- [ ] Task 3: 实现 WeixinJSBridge 调用
- [ ] Task 4: 实现支付结果处理

## Dev Notes

### 微信环境检测

检测微信环境：`navigator.userAgent.indexOf('MicroMessenger') > -1`

### 调起支付

调用 `wx.chooseWXPay()` 或 `WeixinJSBridge.invoke('getBrandWCPayRequest', ...)`

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-2.8]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
