# Story 3.3: 订单状态更新与余额到账

Status: ready-for-dev

## Story

**作为** 系统
**我希望** 回调验证通过后更新订单状态并增加用户余额
**以便** 完成充值流程

## Acceptance Criteria

- [ ] AC1: 使用Redis分布式锁（key: recharge:callback:{order_no}，过期30秒）
- [ ] AC2: 检查订单状态是否为 pending，非pending直接返回SUCCESS
- [ ] AC3: 验证回调金额与订单金额一致
- [ ] AC4: 在同一数据库事务中：更新订单状态、记录transaction_id、增加用户余额、插入余额变动日志
- [ ] AC5: 事务成功后返回 SUCCESS
- [ ] AC6: 事务失败时回滚并返回 FAIL

## Tasks / Subtasks

- [ ] Task 1: 实现 Redis 分布式锁
- [ ] Task 2: 实现订单状态检查
- [ ] Task 3: 实现金额验证
- [ ] Task 4: 实现数据库事务（Ent Tx）
- [ ] Task 5: 实现用户余额更新（行锁）

## Dev Notes

### Redis 分布式锁

使用 SETNX + 过期时间

### 行锁

`SELECT balance FROM users WHERE id = ? FOR UPDATE`

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story-3.3]

## Dev Agent Record

### Agent Model Used

(待开发时填写)

### Completion Notes List

(待开发时填写)

### File List

(待开发时填写)
