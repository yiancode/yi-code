# 数据库脚本目录

本目录包含数据库迁移和维护脚本。

## 执行记录

### 2026-01-22: 添加微信 OpenID 字段

**问题**: 用户登录失败，错误提示 `column users.wechat_openid does not exist`

**原因**: 添加微信公众号登录功能后，数据库缺少 `wechat_openid` 字段

**解决方案**: 手动执行了 `044_add_user_wechat_openid.sql` 迁移脚本

**执行的 DDL**:
```sql
-- 添加 wechat_openid 字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS wechat_openid VARCHAR(64) DEFAULT '';

-- 创建唯一索引（仅对非空且未删除的记录）
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_wechat_openid_unique
    ON users(wechat_openid) WHERE wechat_openid != '' AND deleted_at IS NULL;
```

**执行时间**: 2026-01-22 03:07:00 (CST)

**执行数据库**:
- Host: 106.53.117.99:5432
- Database: code_ai80_vip

**结果**: ✅ 成功
- 字段已添加
- 唯一索引已创建
- 登录功能恢复正常

---

## 脚本说明

### 044_add_user_wechat_openid.sql

为 users 表添加微信 OpenID 字段，用于微信公众号账号绑定和扫码登录。

**功能**:
- 添加 `wechat_openid` 字段 (VARCHAR(64))
- 创建部分唯一索引，确保一个微信账号只能绑定一个用户

**使用场景**:
- 微信公众号扫码登录
- 微信账号绑定

---

## 维护指南

### 如何执行迁移脚本

1. **开发环境**:
```bash
# 连接到开发数据库
PGPASSWORD='your_password' psql -h localhost -p 5432 -U username -d database_name -f scripts/db/xxx_migration.sql
```

2. **生产环境**:
```bash
# 先备份数据库
pg_dump -h host -U user -d database > backup_$(date +%Y%m%d_%H%M%S).sql

# 执行迁移
PGPASSWORD='your_password' psql -h host -p 5432 -U user -d database -f scripts/db/xxx_migration.sql

# 验证结果
psql -h host -U user -d database -c "\d table_name"
```

### 注意事项

1. **生产环境操作前必须备份**
2. **在低峰期执行 DDL 操作**
3. **DDL 执行前先在测试环境验证**
4. **记录执行时间和结果**
5. **索引创建可能需要较长时间，注意表锁**

### 索引创建最佳实践

对于大表，创建索引时建议：
```sql
-- 使用 CONCURRENTLY 避免锁表（不能在事务中使用）
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_name ON table_name(column_name);
```

---

## 参考资料

- [PostgreSQL 官方文档](https://www.postgresql.org/docs/)
- [Goose 迁移工具](https://github.com/pressly/goose)
- [项目迁移目录](../backend/migrations/)
