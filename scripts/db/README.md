# 数据库脚本目录

本目录包含数据库迁移脚本和维护工具，用于管理项目的数据库 schema 变更。

## 📋 目录

- [快速开始](#快速开始)
- [脚本列表](#脚本列表)
- [执行记录](#执行记录)
- [执行指南](#执行指南)
- [回滚操作](#回滚操作)
- [常见问题](#常见问题)
- [参考资料](#参考资料)

---

## 🚀 快速开始

### 连接数据库

```bash
# 方式 1: 使用环境变量
export PGPASSWORD='your_password'
psql -h host -p 5432 -U username -d database_name

# 方式 2: 直接在命令中指定密码
PGPASSWORD='your_password' psql -h host -p 5432 -U username -d database_name
```

### 执行迁移脚本

```bash
# 执行单个迁移
PGPASSWORD='your_password' psql -h host -p 5432 -U username -d database_name -f scripts/db/044_add_user_wechat_openid.sql

# 验证执行结果
PGPASSWORD='your_password' psql -h host -p 5432 -U username -d database_name -c "\d users"
```

---

## 📜 脚本列表

### 044_add_user_wechat_openid.sql

**用途**: 为 users 表添加微信 OpenID 字段

**功能**:
- 添加 `wechat_openid` 字段 (VARCHAR(64))
- 创建部分唯一索引，确保同一个微信账号只能绑定一个用户
- 支持软删除场景（仅对 `deleted_at IS NULL` 的记录生效）

**使用场景**:
- 微信公众号扫码登录
- 微信账号与系统用户绑定
- 防止重复绑定

**表结构变更**:
```sql
ALTER TABLE users ADD COLUMN wechat_openid VARCHAR(64) DEFAULT '';
```

**索引**:
```sql
CREATE UNIQUE INDEX idx_users_wechat_openid_unique
    ON users(wechat_openid)
    WHERE wechat_openid != '' AND deleted_at IS NULL;
```

---

## 📝 执行记录

### 2026-01-22 03:07:00: 添加微信 OpenID 字段

| 项目 | 内容 |
|------|------|
| **脚本** | 044_add_user_wechat_openid.sql |
| **问题** | 用户登录失败，错误: `column users.wechat_openid does not exist` |
| **原因** | 添加微信公众号登录功能后，数据库缺少对应字段 |
| **执行时间** | 2026-01-22 03:07:00 (CST) |
| **执行数据库** | 106.53.117.99:5432 / code_ai80_vip |
| **执行人** | yian |
| **状态** | ✅ 成功 |

**执行的 DDL**:
```sql
-- 添加 wechat_openid 字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS wechat_openid VARCHAR(64) DEFAULT '';

-- 创建唯一索引（仅对非空且未删除的记录）
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_wechat_openid_unique
    ON users(wechat_openid) WHERE wechat_openid != '' AND deleted_at IS NULL;
```

**验证结果**:
- ✅ 字段已添加
- ✅ 唯一索引已创建
- ✅ 登录功能恢复正常
- ✅ 不影响现有用户数据

**影响范围**:
- 表: `users`
- 新增字段: 1 个
- 新增索引: 1 个
- 影响行数: 0 (DDL 操作)

---

## 📖 执行指南

### 1. 开发环境

```bash
# 连接到开发数据库
PGPASSWORD='your_password' psql -h localhost -p 5432 -U username -d database_name -f scripts/db/xxx_migration.sql

# 查看执行结果
PGPASSWORD='your_password' psql -h localhost -p 5432 -U username -d database_name -c "\d table_name"
```

### 2. 生产环境（标准流程）

```bash
# 第一步：备份数据库
pg_dump -h host -U user -d database > backup_$(date +%Y%m%d_%H%M%S).sql

# 第二步：在测试环境验证
PGPASSWORD='test_password' psql -h test-host -p 5432 -U test_user -d test_db -f scripts/db/xxx_migration.sql

# 第三步：执行生产迁移
PGPASSWORD='prod_password' psql -h prod-host -p 5432 -U prod_user -d prod_db -f scripts/db/xxx_migration.sql

# 第四步：验证结果
PGPASSWORD='prod_password' psql -h prod-host -p 5432 -U prod_user -d prod_db -c "\d table_name"

# 第五步：验证业务功能
# 测试登录、注册等核心功能
```

### 3. Docker 环境

```bash
# 进入数据库容器
docker exec -it your_postgres_container psql -U username -d database_name

# 或直接执行 SQL 文件
docker exec -i your_postgres_container psql -U username -d database_name < scripts/db/xxx_migration.sql
```

### 4. 大表迁移（避免锁表）

对于数据量大的表，建议使用 `CONCURRENTLY` 选项：

```sql
-- 使用 CONCURRENTLY 避免长时间锁表
-- 注意：不能在事务中使用
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_name ON table_name(column_name);

-- 添加字段（小表可以直接添加，大表建议分批）
ALTER TABLE large_table ADD COLUMN IF NOT EXISTS new_column VARCHAR(64) DEFAULT '';

-- 对于超大表，可以考虑：
-- 1. 在业务低峰期执行
-- 2. 使用 pg_repack 重建表
-- 3. 创建新表，逐步迁移数据
```

---

## ⏮️ 回滚操作

### 回滚 044_add_user_wechat_openid.sql

```sql
-- 第一步：删除索引
DROP INDEX IF EXISTS idx_users_wechat_openid_unique;

-- 第二步：删除字段
ALTER TABLE users DROP COLUMN IF EXISTS wechat_openid;

-- 第三步：验证
\d users
```

### 回滚注意事项

1. **评估影响**: 回滚前确认是否有数据依赖该字段
2. **业务停机**: 如果业务代码已使用该字段，需要先回滚代码
3. **数据备份**: 如果字段中已有数据，回滚会导致数据丢失
4. **测试验证**: 在测试环境先执行回滚操作

---

## ⚠️ 注意事项

### 生产环境操作规范

1. ✅ **操作前必须备份数据库**
2. ✅ **在业务低峰期执行 DDL 操作**
3. ✅ **先在测试环境验证**
4. ✅ **记录执行时间、执行人和结果**
5. ✅ **准备回滚方案**
6. ✅ **评估执行时间和锁表影响**

### 索引创建最佳实践

```sql
-- 1. 小表（< 10万行）：直接创建
CREATE INDEX IF NOT EXISTS idx_name ON table_name(column_name);

-- 2. 中型表（10万 - 100万行）：使用 CONCURRENTLY
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_name ON table_name(column_name);

-- 3. 大表（> 100万行）：
--    - 在业务低峰期执行
--    - 监控执行进度
--    - 评估磁盘空间（索引约占数据大小的 20-30%）
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_name ON table_name(column_name);
```

### 性能影响评估

| 操作类型 | 锁级别 | 影响 | 建议 |
|---------|--------|------|------|
| ADD COLUMN (无默认值) | ACCESS EXCLUSIVE | 短暂锁表 | 可在业务时间执行 |
| ADD COLUMN (有默认值) | ACCESS EXCLUSIVE | 长时间锁表 | 低峰期执行 |
| CREATE INDEX | SHARE | 阻塞写操作 | 使用 CONCURRENTLY |
| CREATE INDEX CONCURRENTLY | SHARE UPDATE EXCLUSIVE | 不阻塞读写 | 推荐 |
| DROP COLUMN | ACCESS EXCLUSIVE | 短暂锁表 | 低峰期执行 |

---

## 🔍 常见问题

### Q1: 如何查看当前数据库的所有表？

```sql
\dt
-- 或
SELECT tablename FROM pg_tables WHERE schemaname = 'public';
```

### Q2: 如何查看表结构？

```sql
\d table_name
-- 或查看详细信息
\d+ table_name
```

### Q3: 如何查看所有索引？

```sql
\di
-- 或查看特定表的索引
SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'users';
```

### Q4: 迁移脚本执行失败如何处理？

```bash
# 1. 查看错误日志
tail -n 100 /var/log/postgresql/postgresql.log

# 2. 检查数据库状态
PGPASSWORD='password' psql -h host -U user -d database -c "SELECT version();"

# 3. 回滚到备份
psql -h host -U user -d database < backup_20260122.sql

# 4. 检查事务状态
SELECT * FROM pg_stat_activity WHERE state = 'idle in transaction';
```

### Q5: 如何检查索引是否创建成功？

```sql
-- 方式 1: 查看表索引
\d users

-- 方式 2: 查询系统表
SELECT
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE tablename = 'users'
  AND indexname = 'idx_users_wechat_openid_unique';

-- 方式 3: 验证索引生效
EXPLAIN SELECT * FROM users WHERE wechat_openid = 'test';
```

### Q6: 如何监控长时间运行的查询？

```sql
-- 查看当前活动的查询
SELECT
    pid,
    usename,
    application_name,
    state,
    query_start,
    now() - query_start AS duration,
    query
FROM pg_stat_activity
WHERE state != 'idle'
ORDER BY duration DESC;

-- 终止长时间运行的查询（慎用）
SELECT pg_terminate_backend(pid);
```

### Q7: 部分索引 (Partial Index) 的优势是什么？

部分索引只对满足特定条件的行建立索引，优势包括：

1. **节省存储空间**: 只索引需要的数据
2. **提高查询性能**: 索引更小，查询更快
3. **支持软删除**: 可以对 `deleted_at IS NULL` 建立唯一索引

```sql
-- 示例：只对未删除的记录建立唯一索引
CREATE UNIQUE INDEX idx_users_wechat_openid_unique
    ON users(wechat_openid)
    WHERE wechat_openid != '' AND deleted_at IS NULL;
```

---

## 📚 参考资料

### 官方文档
- [PostgreSQL 官方文档](https://www.postgresql.org/docs/)
- [PostgreSQL DDL 语句](https://www.postgresql.org/docs/current/ddl.html)
- [PostgreSQL 索引](https://www.postgresql.org/docs/current/indexes.html)
- [PostgreSQL 部分索引](https://www.postgresql.org/docs/current/indexes-partial.html)

### 工具
- [Goose 迁移工具](https://github.com/pressly/goose)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [pgAdmin](https://www.pgadmin.org/)
- [DBeaver](https://dbeaver.io/)

### 项目相关
- [项目迁移目录](../../backend/migrations/)
- [数据库 Schema 定义](../../backend/ent/schema/)
- [数据库配置](../../backend/config.yaml)

### 最佳实践
- [PostgreSQL Performance Optimization](https://wiki.postgresql.org/wiki/Performance_Optimization)
- [PostgreSQL Locking](https://www.postgresql.org/docs/current/explicit-locking.html)
- [Zero-Downtime Migrations](https://www.braintreepayments.com/blog/safe-operations-for-high-volume-postgresql/)

---

## 📞 联系方式

如有问题或建议，请联系：
- 技术负责人: yian20133213@gmail.com
- 项目仓库: [GitHub](https://github.com/yiancode/yi-code)

---

**最后更新**: 2026-01-22
**维护者**: yian
