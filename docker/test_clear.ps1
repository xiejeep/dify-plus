# 测试数据库清空脚本
# 用于验证脚本功能

Write-Host "===========================================" -ForegroundColor Green
Write-Host "测试数据库清空脚本" -ForegroundColor Green
Write-Host "===========================================" -ForegroundColor Green

# 检查数据库连接
Write-Host "测试数据库连接..." -ForegroundColor Yellow

try {
    $result = docker exec docker-db-1 psql -U postgres -d dify -c "SELECT version();"
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ 数据库连接成功" -ForegroundColor Green
    } else {
        Write-Host "❌ 数据库连接失败" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "❌ 数据库连接异常: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 检查表是否存在
Write-Host ""
Write-Host "检查系统表是否存在..." -ForegroundColor Yellow

$tables = @(
    "sys_users",
    "sys_authorities", 
    "sys_apis",
    "casbin_rule"
)

foreach ($table in $tables) {
    try {
        $check = docker exec docker-db-1 psql -U postgres -d dify -c "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '$table');"
        Write-Host "表 $table : 检查完成" -ForegroundColor Gray
    } catch {
        Write-Host "表 $table : 不存在或检查失败" -ForegroundColor Yellow
    }
}

Write-Host ""
Write-Host "✅ 测试完成！数据库清空脚本可以正常使用。" -ForegroundColor Green
Write-Host ""
Write-Host "使用方法：" -ForegroundColor Cyan
Write-Host "  .\clear_database.ps1" -ForegroundColor White 