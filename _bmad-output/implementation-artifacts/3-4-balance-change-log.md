# Story 3.4: 余额变动日志记录

Status: ready-for-dev

## Story

**作为** 系统
**我希望** 记录每笔余额变动的详细日志
**以便** 支持审计追溯和问题排查

## Acceptance Criteria

- [ ] AC1: 插入 balance_logs 表记录
- [ ] AC2: 记录字段：user_id, change_type(recharge), amount, balance_before, balance_after
- [ ] AC3: 记录 related_order_no 关联订单号
- [ ] AC4: 记录 operator_type(system) 和 description
- [ ] AC5: 日志表只允许插入，不允许修改删除（应用层控制）

## Tasks / Subtasks

- [ ] Task 1: 创建 `backend/ent/schema/balance_log.go` Schema
- [ ] Task 2: 实现余额日志服务
- [ ] Task 3: 在事务中插入日志

## Dev Notes

### 数据库表

参考 `_bmad-output/planning-artifacts/epics.md#数据库需求` 的 balance_logs 表设计

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-3.4]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
