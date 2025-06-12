# 用户默认额度配置指南

## 概述

本文档说明如何配置新用户的默认总额度。系统已经重构为配置驱动的方式，避免了硬编码问题。

## 配置方式

### 1. 环境变量配置（推荐）

在 Docker 环境中，可以通过环境变量来设置默认额度：

```bash
# 在 docker-compose.dify-plus.yaml 中的 api 服务添加环境变量
environment:
  - DEFAULT_ACCOUNT_TOTAL_QUOTA=50.0  # 设置默认额度为50美元
```

### 2. 直接修改配置文件

在 `api/configs/dify_config.py` 中修改默认值：

```python
DEFAULT_ACCOUNT_TOTAL_QUOTA = os.environ.get("DEFAULT_ACCOUNT_TOTAL_QUOTA", "50.0")  # 修改默认值
```

## 配置说明

- **DEFAULT_ACCOUNT_TOTAL_QUOTA**: 新用户的默认总额度（单位：美元）
- **默认值**: 30.0 美元
- **生效时机**: 用户首次产生费用时系统自动创建额度记录

## 重构内容

### 修改前（硬编码）
```python
total_quota=15,  # TODO 初始总额度这里到时候默认15要改
```

### 修改后（配置驱动）
```python
total_quota=float(dify_config.DEFAULT_ACCOUNT_TOTAL_QUOTA),  # 从配置获取默认总额度
```

## 影响的文件

1. `api/configs/dify_config.py` - 配置定义
2. `api/configs/feature/__init__.py` - 功能配置
3. `api/events/event_handlers/update_account_money_when_messaeg_created_extend.py` - 消息创建事件
4. `api/tasks/extend/update_account_money_when_workflow_node_execution_created_extend.py` - 工作流执行任务
5. `api/controllers/console/workspace/account_extend.py` - 账户接口

## 使用示例

### Docker Compose 配置示例

```yaml
services:
  api:
    environment:
      # 设置默认额度为100美元
      - DEFAULT_ACCOUNT_TOTAL_QUOTA=100.0
      # 其他环境变量...
```

### 重启服务应用配置

```bash
cd docker
sudo docker compose -f docker-compose.dify-plus.yaml restart api worker
```

## 优势

1. **配置集中管理**: 所有额度相关配置集中在一个地方
2. **环境适配**: 不同环境可以设置不同的默认额度
3. **易于维护**: 避免了多处硬编码的维护问题
4. **动态调整**: 无需修改代码即可调整默认额度

## 更新现有用户额度

### 使用更新脚本

我们提供了一个便捷的脚本来更新现有用户的额度：

```bash
# 预览模式（查看会影响哪些用户，不实际更新）
cd api
python scripts/update_existing_user_quota.py --quota 50.0

# 执行更新（实际修改数据库）
python scripts/update_existing_user_quota.py --quota 50.0 --confirm
```

### 手动更新SQL

如果你更愿意使用SQL直接更新：

```sql
-- 查看当前所有用户额度
SELECT account_id, total_quota, used_quota FROM account_money_extend;

-- 更新所有用户的总额度为50美元
UPDATE account_money_extend SET total_quota = 50.0;
```

## 注意事项

1. 修改配置后需要重启相关服务
2. 已存在的用户额度记录不会自动更新，需要使用脚本或SQL手动更新
3. 配置值为字符串类型，系统会自动转换为浮点数
4. 建议在生产环境操作前先进行备份 