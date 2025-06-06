# 数据库清空脚本 (PowerShell 版本)
# 用于解决 dify-plus 初始化失败问题

Write-Host "===========================================" -ForegroundColor Green
Write-Host "Dify-Plus 数据库清空脚本 (PowerShell)" -ForegroundColor Green
Write-Host "===========================================" -ForegroundColor Green

# 数据库连接参数（根据 docker-compose 配置）
$DB_HOST = "localhost"
$DB_PORT = "5432"
$DB_NAME = "dify"
$DB_USER = "postgres"
$DB_PASSWORD = "difyai123456"

# 检查是否有 Docker 容器运行
Write-Host "检查数据库容器状态..." -ForegroundColor Yellow

$containerCheck = docker ps | Select-String "docker-db-1"
if (-not $containerCheck) {
    Write-Host "❌ 错误: 数据库容器未运行，请先启动 dify-plus 服务" -ForegroundColor Red
    Write-Host "   运行命令: docker-compose -f docker-compose.dify-plus.yaml up -d" -ForegroundColor Yellow
    exit 1
}

Write-Host "✅ 数据库容器正在运行" -ForegroundColor Green

# 确认操作
Write-Host ""
Write-Host "⚠️  警告: 此操作将清空所有系统表数据，包括：" -ForegroundColor Yellow
Write-Host "   - 用户账户数据" -ForegroundColor Yellow
Write-Host "   - 权限和角色数据" -ForegroundColor Yellow
Write-Host "   - API 配置数据" -ForegroundColor Yellow
Write-Host "   - 操作记录数据" -ForegroundColor Yellow
Write-Host "   - 其他系统配置数据" -ForegroundColor Yellow
Write-Host ""

$confirm = Read-Host "确认要继续吗？(输入 yes 继续)"

if ($confirm -ne "yes") {
    Write-Host "操作已取消" -ForegroundColor Yellow
    exit 0
}

Write-Host ""
Write-Host "开始清空数据库..." -ForegroundColor Green

# 使用 docker exec 连接到数据库容器执行 SQL
Write-Host "正在执行数据库清空操作..." -ForegroundColor Yellow

try {
    # 执行清空操作
    Get-Content .\clear_database.sql | docker exec -i docker-db-1 psql -U $DB_USER -d $DB_NAME
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "✅ 数据库清空成功！" -ForegroundColor Green
        Write-Host ""
        Write-Host "📝 后续步骤：" -ForegroundColor Cyan
        Write-Host "1. 重启 admin-server 容器：" -ForegroundColor White
        Write-Host "   docker-compose -f docker-compose.dify-plus.yaml restart admin-server" -ForegroundColor Gray
        Write-Host ""
        Write-Host "2. 访问管理中心进行初始化：" -ForegroundColor White
        Write-Host "   http://localhost:8081" -ForegroundColor Gray
        Write-Host ""
        Write-Host "3. 如果仍有问题，可以重启所有服务：" -ForegroundColor White
        Write-Host "   docker-compose -f docker-compose.dify-plus.yaml restart" -ForegroundColor Gray
    } else {
        Write-Host "❌ 数据库清空失败，请检查错误信息" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "❌ 执行过程中发生错误: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
} 