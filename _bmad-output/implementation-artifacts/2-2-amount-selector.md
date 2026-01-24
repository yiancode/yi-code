# Story 2.2: 充值金额选择器

Status: ready-for-dev

## Story

**作为** 普通用户
**我希望** 通过快捷按钮或自定义输入选择充值金额
**以便** 方便快速地选择常用金额或精确输入

## Acceptance Criteria

- [ ] AC1: 显示快捷金额按钮（从配置获取，如：10、50、100、200、500元）
- [ ] AC2: 点击快捷按钮选中对应金额，按钮高亮显示
- [ ] AC3: 支持自定义金额输入框
- [ ] AC4: 输入自定义金额时取消快捷按钮选中状态
- [ ] AC5: 金额输入只允许数字和小数点，最多2位小数
- [ ] AC6: 显示充值金额范围提示（如：最小1元，最大1000元）

## Tasks / Subtasks

- [ ] Task 1: 创建 `AmountSelector.vue` 组件
- [ ] Task 2: 实现快捷金额按钮
- [ ] Task 3: 实现自定义金额输入
- [ ] Task 4: 实现输入校验（正则表达式）

## Dev Notes

### 组件路径

前端组件路径：`src/components/user/recharge/AmountSelector.vue`

### 双向绑定

使用 v-model 双向绑定金额

### UX 设计

参考 `_bmad-output/planning-artifacts/ux-design.md` 的金额选择器设计

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-2.2]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
