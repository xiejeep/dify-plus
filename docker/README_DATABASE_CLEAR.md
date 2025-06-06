# 数据库清空脚本使用说明

## 概述

当 dify-plus 初始化失败时，通常是由于数据库中存在不一致的数据导致的。这个脚本集合可以帮助您完全清空相关系统表，以便重新进行干净的初始化。

## 脚本文件

1. **clear_database.sql** - SQL 清空脚本
2. **clear_database.sh** - Linux/Mac 执行脚本
3. **clear_database.ps1** - Windows PowerShell 执行脚本

## 使用场景

### 适用情况
- 🔸 初始化过程中出现权限冲突错误
- 🔸 系统提示数据库状态不一致
- 🔸 管理员账户创建失败
- 🔸 权限和角色配置异常
- 🔸 需要完全重置系统到初始状态

### 不适用情况
- ❌ 生产环境中有重要业务数据
- ❌ 只需要修改部分配置的情况

## 使用方法

### Windows 用户

```powershell
# 1. 确保在 docker 目录下
cd E:\code\dify-plus\docker

# 2. 确保 dify-plus 服务正在运行
docker-compose -f docker-compose.dify-plus.yaml ps

# 3. 执行清空脚本
.\clear_database.ps1
```

### Linux/Mac 用户

```bash
# 1. 确保在 docker 目录下
cd /path/to/dify-plus/docker

# 2. 确保 dify-plus 服务正在运行
docker-compose -f docker-compose.dify-plus.yaml ps

# 3. 执行清空脚本
./clear_database.sh
```

### 手动执行 SQL（高级用户）

```bash
# 直接在数据库容器中执行 SQL 脚本
docker exec -i docker-db-1 psql -U postgres -d dify < clear_database.sql
```

## 清空的表列表

脚本会清空以下系统表：

### 核心系统表
- `sys_authorities` - 权限角色表
- `sys_apis` - API 配置表
- `casbin_rule` - Casbin 权限规则表
- `sys_users` - 用户表
- `sys_base_menus` - 菜单表

### 配置和记录表
- `sys_dictionaries` - 字典表
- `sys_dictionary_details` - 字典详情表
- `sys_operation_records` - 操作记录表
- `jwt_blacklists` - JWT 黑名单表

### 业务数据表
- `exa_file_upload_and_downloads` - 文件上传记录表
- `exa_file_chunks` - 文件块表
- `exa_customers` - 客户表

### 新增功能表
- `sys_email_verifications` - 邮箱验证表（用户注册功能）

## 执行后的操作步骤

### 1. 重启服务
```bash
# 重启管理服务器
docker-compose -f docker-compose.dify-plus.yaml restart admin-server

# 或者重启全部服务
docker-compose -f docker-compose.dify-plus.yaml restart
```

### 2. 重新初始化
1. 打开浏览器访问：`http://localhost:8081`
2. 系统会自动检测需要初始化
3. 按照提示完成系统初始化
4. 创建管理员账户

### 3. 验证功能
- 登录管理中心
- 检查用户注册功能：`http://localhost:8081/#/register`
- 确认各项系统功能正常

## 安全注意事项

⚠️ **重要警告**
- 此操作不可逆，会删除所有系统数据
- 执行前请确认没有重要的业务数据
- 建议在测试环境中先验证

✅ **安全特性**
- 脚本会检查容器运行状态
- 需要用户手动确认才能执行
- 使用 TRUNCATE ... CASCADE 确保数据完整性

## 故障排除

### 常见错误

1. **容器未运行**
```
❌ 错误: 数据库容器未运行
解决: docker-compose -f docker-compose.dify-plus.yaml up -d
```

2. **权限不足**
```
❌ 错误: permission denied
解决: 确保 Docker 有足够权限，或使用管理员权限运行
```

3. **表不存在**
```
❌ 错误: relation "xxx" does not exist
解决: 正常情况，表可能还未创建，继续初始化即可
```

### 联系支持
如果在使用过程中遇到问题，请提供：
- 错误信息截图
- Docker 容器运行状态
- 操作系统版本信息 