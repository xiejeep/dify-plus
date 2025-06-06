#!/bin/bash

# 数据库清空脚本
# 用于解决 dify-plus 初始化失败问题

echo "=========================================="
echo "Dify-Plus 数据库清空脚本"
echo "=========================================="

# 数据库连接参数（根据 docker-compose 配置）
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="dify"
DB_USER="postgres"
DB_PASSWORD="difyai123456"

# 检查是否有 Docker 容器运行
echo "检查数据库容器状态..."
if ! docker ps | grep -q "docker-db-1"; then
    echo "❌ 错误: 数据库容器未运行，请先启动 dify-plus 服务"
    echo "   运行命令: docker-compose -f docker-compose.dify-plus.yaml up -d"
    exit 1
fi

echo "✅ 数据库容器正在运行"

# 确认操作
echo ""
echo "⚠️  警告: 此操作将清空所有系统表数据，包括："
echo "   - 用户账户数据"
echo "   - 权限和角色数据"
echo "   - API 配置数据"
echo "   - 操作记录数据"
echo "   - 其他系统配置数据"
echo ""
read -p "确认要继续吗？(输入 yes 继续): " confirm

if [ "$confirm" != "yes" ]; then
    echo "操作已取消"
    exit 0
fi

echo ""
echo "开始清空数据库..."

# 方法1: 使用 docker exec 连接到数据库容器执行 SQL
echo "正在执行数据库清空操作..."

docker exec -i docker-db-1 psql -U "$DB_USER" -d "$DB_NAME" < clear_database.sql

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ 数据库清空成功！"
    echo ""
    echo "📝 后续步骤："
    echo "1. 重启 admin-server 容器："
    echo "   docker-compose -f docker-compose.dify-plus.yaml restart admin-server"
    echo ""
    echo "2. 访问管理中心进行初始化："
    echo "   http://localhost:8081"
    echo ""
    echo "3. 如果仍有问题，可以重启所有服务："
    echo "   docker-compose -f docker-compose.dify-plus.yaml restart"
else
    echo "❌ 数据库清空失败，请检查错误信息"
    exit 1
fi 