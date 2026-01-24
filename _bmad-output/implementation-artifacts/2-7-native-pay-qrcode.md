# Story 2.7: Native支付二维码生成

Status: ready-for-dev

## Story

**作为** 普通用户（PC端）
**我希望** 看到支付二维码进行扫码支付
**以便** 在PC端使用微信扫码完成支付

## Acceptance Criteria

- [ ] AC1: Native支付返回 code_url
- [ ] AC2: 前端根据 code_url 生成二维码图片
- [ ] AC3: 二维码尺寸适中（200x200像素）
- [ ] AC4: 二维码下方显示「请使用微信扫码支付」提示
- [ ] AC5: 二维码生成时间 < 1秒

## Tasks / Subtasks

- [ ] Task 1: 创建 `QRCodeDisplay.vue` 组件
- [ ] Task 2: 集成 qrcode.js 或 vue-qrcode
- [ ] Task 3: 实现二维码展示样式

## Dev Notes

### 组件路径

前端组件路径：`src/components/user/recharge/QRCodeDisplay.vue`

### 二维码库

使用 qrcode.js 或 vue-qrcode 组件生成二维码

### UX 设计

参考 `_bmad-output/planning-artifacts/ux-design.md` 的支付中页面设计

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-2.7]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
