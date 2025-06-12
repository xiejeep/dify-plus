#!/bin/bash

# Dify登录锁定解除脚本
# 作用：快速解除用户邮箱的登录锁定状态

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Redis配置
REDIS_CONTAINER="docker-redis-1"
REDIS_PASSWORD="difyai123456"

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查Docker容器是否运行
check_redis_container() {
    if ! docker ps | grep -q "$REDIS_CONTAINER"; then
        print_error "Redis容器 $REDIS_CONTAINER 未运行"
        print_info "请确保Dify服务正在运行"
        exit 1
    fi
    print_success "Redis容器检查通过"
}

# 显示所有被锁定的邮箱
show_locked_emails() {
    print_info "正在查找被锁定的邮箱..."
    
    locked_keys=$(docker exec $REDIS_CONTAINER redis-cli -a $REDIS_PASSWORD keys "*login_error_rate_limit*" 2>/dev/null)
    
    if [ -z "$locked_keys" ]; then
        print_success "当前没有被锁定的邮箱"
        return 0
    fi
    
    echo -e "\n${YELLOW}被锁定的邮箱列表：${NC}"
    echo "----------------------------------------"
    
    for key in $locked_keys; do
        email=$(echo $key | sed 's/login_error_rate_limit://')
        error_count=$(docker exec $REDIS_CONTAINER redis-cli -a $REDIS_PASSWORD get "$key" 2>/dev/null)
        ttl=$(docker exec $REDIS_CONTAINER redis-cli -a $REDIS_PASSWORD ttl "$key" 2>/dev/null)
        
        if [ "$ttl" -gt 0 ]; then
            hours=$((ttl / 3600))
            minutes=$(((ttl % 3600) / 60))
            echo -e "📧 ${BLUE}$email${NC}"
            echo -e "   错误次数: ${RED}$error_count${NC}"
            echo -e "   剩余时间: ${YELLOW}${hours}小时${minutes}分钟${NC}"
            echo "----------------------------------------"
        fi
    done
    echo ""
}

# 解除指定邮箱的锁定
unlock_email() {
    local email="$1"
    local key="login_error_rate_limit:$email"
    
    print_info "正在解除 $email 的登录锁定..."
    
    # 检查是否存在锁定记录
    exists=$(docker exec $REDIS_CONTAINER redis-cli -a $REDIS_PASSWORD exists "$key" 2>/dev/null)
    
    if [ "$exists" = "0" ]; then
        print_warning "邮箱 $email 没有被锁定"
        return 1
    fi
    
    # 获取错误次数和剩余时间
    error_count=$(docker exec $REDIS_CONTAINER redis-cli -a $REDIS_PASSWORD get "$key" 2>/dev/null)
    ttl=$(docker exec $REDIS_CONTAINER redis-cli -a $REDIS_PASSWORD ttl "$key" 2>/dev/null)
    
    echo -e "   当前错误次数: ${RED}$error_count${NC}"
    if [ "$ttl" -gt 0 ]; then
        hours=$((ttl / 3600))
        minutes=$(((ttl % 3600) / 60))
        echo -e "   原剩余时间: ${YELLOW}${hours}小时${minutes}分钟${NC}"
    fi
    
    # 删除锁定记录
    result=$(docker exec $REDIS_CONTAINER redis-cli -a $REDIS_PASSWORD del "$key" 2>/dev/null)
    
    if [ "$result" = "1" ]; then
        print_success "✅ 邮箱 $email 的登录锁定已成功解除！"
        print_info "用户现在可以正常登录了"
        return 0
    else
        print_error "解除锁定失败"
        return 1
    fi
}

# 解除所有邮箱的锁定
unlock_all_emails() {
    print_info "正在解除所有邮箱的登录锁定..."
    
    locked_keys=$(docker exec $REDIS_CONTAINER redis-cli -a $REDIS_PASSWORD keys "*login_error_rate_limit*" 2>/dev/null)
    
    if [ -z "$locked_keys" ]; then
        print_success "当前没有被锁定的邮箱"
        return 0
    fi
    
    count=0
    for key in $locked_keys; do
        email=$(echo $key | sed 's/login_error_rate_limit://')
        if unlock_email "$email"; then
            count=$((count + 1))
        fi
    done
    
    print_success "共解除了 $count 个邮箱的登录锁定"
}

# 显示帮助信息
show_help() {
    echo -e "${BLUE}Dify登录锁定解除脚本${NC}"
    echo ""
    echo "用法："
    echo "  $0                    # 交互式菜单"
    echo "  $0 list               # 显示所有被锁定的邮箱"
    echo "  $0 unlock <email>     # 解除指定邮箱的锁定"
    echo "  $0 unlock-all         # 解除所有邮箱的锁定"
    echo "  $0 help               # 显示帮助信息"
    echo ""
    echo "示例："
    echo "  $0 unlock user@example.com"
    echo "  $0 list"
    echo "  $0 unlock-all"
}

# 交互式菜单
interactive_menu() {
    while true; do
        echo ""
        echo -e "${BLUE}=== Dify登录锁定解除工具 ===${NC}"
        echo "1. 显示被锁定的邮箱"
        echo "2. 解除指定邮箱的锁定"
        echo "3. 解除所有邮箱的锁定"
        echo "4. 退出"
        echo ""
        
        # 检查是否在交互式终端中运行
        if [ -t 0 ]; then
            read -p "请选择操作 (1-4): " choice
        else
            print_error "此脚本需要在交互式终端中运行"
            print_info "请直接运行: ./unlock_login.sh"
            print_info "或使用命令行参数，如: ./unlock_login.sh list"
            exit 1
        fi
        
        case $choice in
            1)
                show_locked_emails
                ;;
            2)
                if [ -t 0 ]; then
                    read -p "请输入要解除锁定的邮箱: " email
                    if [ -n "$email" ]; then
                        unlock_email "$email"
                    else
                        print_error "邮箱不能为空"
                    fi
                else
                    print_error "无法在非交互式模式下读取邮箱输入"
                    exit 1
                fi
                ;;
            3)
                if [ -t 0 ]; then
                    read -p "确定要解除所有邮箱的锁定吗？(y/N): " confirm
                    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
                        unlock_all_emails
                    else
                        print_info "操作已取消"
                    fi
                else
                    print_error "无法在非交互式模式下读取确认输入"
                    exit 1
                fi
                ;;
            4)
                print_info "退出程序"
                exit 0
                ;;
            "")
                print_warning "请输入选项"
                ;;
            *)
                print_error "无效选择，请输入 1-4"
                ;;
        esac
        
        # 添加暂停，让用户看到结果
        if [ -t 0 ]; then
            echo ""
            read -p "按回车键继续..." dummy
        fi
    done
}

# 主函数
main() {
    print_info "Dify登录锁定解除脚本启动"
    
    # 检查Redis容器
    check_redis_container
    
    # 根据参数执行不同操作
    case "$1" in
        "list")
            show_locked_emails
            ;;
        "unlock")
            if [ -z "$2" ]; then
                print_error "请指定要解除锁定的邮箱"
                echo "用法: $0 unlock <email>"
                exit 1
            fi
            unlock_email "$2"
            ;;
        "unlock-all")
            unlock_all_emails
            ;;
        "help")
            show_help
            ;;
        "")
            interactive_menu
            ;;
        *)
            print_error "未知参数: $1"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@" 