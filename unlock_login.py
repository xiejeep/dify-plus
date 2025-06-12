#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Dify登录锁定解除工具 (Python版本)
功能：快速解除用户邮箱的登录锁定状态，支持批量操作和详细统计
作者：AI Assistant
版本：1.0
"""

import redis
import argparse
import sys
import time
from datetime import datetime, timedelta
from typing import List, Dict, Optional
import json

# 配置
REDIS_HOST = "localhost"
REDIS_PORT = 6379
REDIS_PASSWORD = "difyai123456"
REDIS_DB = 0

# 颜色常量
class Colors:
    RED = '\033[0;31m'
    GREEN = '\033[0;32m'
    YELLOW = '\033[1;33m'
    BLUE = '\033[0;34m'
    PURPLE = '\033[0;35m'
    CYAN = '\033[0;36m'
    WHITE = '\033[1;37m'
    RESET = '\033[0m'
    BOLD = '\033[1m'

class LoginUnlocker:
    def __init__(self):
        """初始化Redis连接"""
        try:
            self.redis_client = redis.Redis(
                host=REDIS_HOST,
                port=REDIS_PORT,
                password=REDIS_PASSWORD,
                db=REDIS_DB,
                decode_responses=True
            )
            # 测试连接
            self.redis_client.ping()
            self.print_success("Redis连接成功")
        except redis.ConnectionError:
            self.print_error("无法连接到Redis服务")
            self.print_info("请确保Dify服务正在运行且Redis配置正确")
            sys.exit(1)
        except redis.AuthenticationError:
            self.print_error("Redis认证失败")
            self.print_info("请检查Redis密码配置")
            sys.exit(1)

    def print_info(self, message: str):
        """打印信息消息"""
        print(f"{Colors.BLUE}[INFO]{Colors.RESET} {message}")

    def print_success(self, message: str):
        """打印成功消息"""
        print(f"{Colors.GREEN}[SUCCESS]{Colors.RESET} {message}")

    def print_warning(self, message: str):
        """打印警告消息"""
        print(f"{Colors.YELLOW}[WARNING]{Colors.RESET} {message}")

    def print_error(self, message: str):
        """打印错误消息"""
        print(f"{Colors.RED}[ERROR]{Colors.RESET} {message}")

    def print_header(self, title: str):
        """打印标题"""
        print(f"\n{Colors.CYAN}{Colors.BOLD}=== {title} ==={Colors.RESET}")

    def get_locked_emails(self) -> Dict[str, Dict]:
        """获取所有被锁定的邮箱信息"""
        try:
            keys = self.redis_client.keys("login_error_rate_limit:*")
            locked_emails = {}
            
            for key in keys:
                email = key.replace("login_error_rate_limit:", "")
                error_count = self.redis_client.get(key)
                ttl = self.redis_client.ttl(key)
                
                if ttl > 0:  # 只显示还未过期的锁定
                    locked_emails[email] = {
                        'error_count': int(error_count) if error_count else 0,
                        'ttl': ttl,
                        'expire_time': datetime.now() + timedelta(seconds=ttl),
                        'key': key
                    }
            
            return locked_emails
        except Exception as e:
            self.print_error(f"获取锁定邮箱信息失败: {e}")
            return {}

    def format_time(self, seconds: int) -> str:
        """格式化时间显示"""
        if seconds <= 0:
            return "已过期"
        
        hours = seconds // 3600
        minutes = (seconds % 3600) // 60
        secs = seconds % 60
        
        if hours > 0:
            return f"{hours}小时{minutes}分钟{secs}秒"
        elif minutes > 0:
            return f"{minutes}分钟{secs}秒"
        else:
            return f"{secs}秒"

    def show_locked_emails(self, detailed: bool = False):
        """显示所有被锁定的邮箱"""
        self.print_info("正在查找被锁定的邮箱...")
        
        locked_emails = self.get_locked_emails()
        
        if not locked_emails:
            self.print_success("✅ 当前没有被锁定的邮箱")
            return
        
        self.print_header(f"被锁定的邮箱列表 (共{len(locked_emails)}个)")
        
        for i, (email, info) in enumerate(locked_emails.items(), 1):
            print(f"\n{Colors.WHITE}{i}.{Colors.RESET} 📧 {Colors.BLUE}{email}{Colors.RESET}")
            print(f"   错误次数: {Colors.RED}{info['error_count']}{Colors.RESET}")
            print(f"   剩余时间: {Colors.YELLOW}{self.format_time(info['ttl'])}{Colors.RESET}")
            
            if detailed:
                print(f"   过期时间: {Colors.PURPLE}{info['expire_time'].strftime('%Y-%m-%d %H:%M:%S')}{Colors.RESET}")
                print(f"   Redis键名: {Colors.CYAN}{info['key']}{Colors.RESET}")
            
            print("   " + "-" * 50)

    def unlock_email(self, email: str) -> bool:
        """解除指定邮箱的锁定"""
        key = f"login_error_rate_limit:{email}"
        
        self.print_info(f"正在解除 {Colors.BLUE}{email}{Colors.RESET} 的登录锁定...")
        
        # 检查是否存在锁定记录
        if not self.redis_client.exists(key):
            self.print_warning(f"邮箱 {email} 没有被锁定")
            return False
        
        # 获取锁定信息
        error_count = self.redis_client.get(key)
        ttl = self.redis_client.ttl(key)
        
        print(f"   当前错误次数: {Colors.RED}{error_count}{Colors.RESET}")
        if ttl > 0:
            print(f"   原剩余时间: {Colors.YELLOW}{self.format_time(ttl)}{Colors.RESET}")
        
        # 删除锁定记录
        try:
            result = self.redis_client.delete(key)
            if result:
                self.print_success(f"✅ 邮箱 {email} 的登录锁定已成功解除！")
                self.print_info("用户现在可以正常登录了")
                return True
            else:
                self.print_error("解除锁定失败")
                return False
        except Exception as e:
            self.print_error(f"解除锁定时出错: {e}")
            return False

    def unlock_all_emails(self) -> int:
        """解除所有邮箱的锁定"""
        locked_emails = self.get_locked_emails()
        
        if not locked_emails:
            self.print_success("当前没有被锁定的邮箱")
            return 0
        
        self.print_info(f"正在解除所有邮箱的登录锁定... (共{len(locked_emails)}个)")
        
        success_count = 0
        for email in locked_emails.keys():
            if self.unlock_email(email):
                success_count += 1
            print()  # 空行分隔
        
        self.print_success(f"共解除了 {success_count}/{len(locked_emails)} 个邮箱的登录锁定")
        return success_count

    def unlock_batch(self, emails: List[str]) -> Dict[str, bool]:
        """批量解除指定邮箱的锁定"""
        results = {}
        
        self.print_info(f"正在批量解除登录锁定... (共{len(emails)}个邮箱)")
        
        for email in emails:
            results[email] = self.unlock_email(email)
            print()  # 空行分隔
        
        success_count = sum(results.values())
        self.print_success(f"批量操作完成: {success_count}/{len(emails)} 个邮箱解除成功")
        
        return results

    def export_locked_emails(self, filename: str = None):
        """导出被锁定的邮箱信息到JSON文件"""
        if filename is None:
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            filename = f"locked_emails_{timestamp}.json"
        
        locked_emails = self.get_locked_emails()
        
        if not locked_emails:
            self.print_warning("没有被锁定的邮箱，无需导出")
            return
        
        # 准备导出数据
        export_data = {
            'export_time': datetime.now().isoformat(),
            'total_count': len(locked_emails),
            'emails': {}
        }
        
        for email, info in locked_emails.items():
            export_data['emails'][email] = {
                'error_count': info['error_count'],
                'remaining_seconds': info['ttl'],
                'expire_time': info['expire_time'].isoformat(),
                'redis_key': info['key']
            }
        
        try:
            with open(filename, 'w', encoding='utf-8') as f:
                json.dump(export_data, f, ensure_ascii=False, indent=2)
            
            self.print_success(f"✅ 锁定邮箱信息已导出到: {filename}")
            self.print_info(f"共导出 {len(locked_emails)} 个被锁定的邮箱")
        except Exception as e:
            self.print_error(f"导出失败: {e}")

    def get_statistics(self):
        """获取锁定统计信息"""
        locked_emails = self.get_locked_emails()
        
        if not locked_emails:
            self.print_success("✅ 当前没有被锁定的邮箱")
            return
        
        # 统计信息
        total_count = len(locked_emails)
        error_counts = [info['error_count'] for info in locked_emails.values()]
        ttls = [info['ttl'] for info in locked_emails.values()]
        
        avg_errors = sum(error_counts) / len(error_counts)
        avg_time_left = sum(ttls) / len(ttls)
        
        # 按剩余时间分组
        time_groups = {
            '1小时内': 0,
            '1-6小时': 0,
            '6-12小时': 0,
            '12-24小时': 0
        }
        
        for ttl in ttls:
            hours = ttl / 3600
            if hours <= 1:
                time_groups['1小时内'] += 1
            elif hours <= 6:
                time_groups['1-6小时'] += 1
            elif hours <= 12:
                time_groups['6-12小时'] += 1
            else:
                time_groups['12-24小时'] += 1
        
        self.print_header("锁定统计信息")
        print(f"📊 总锁定邮箱数: {Colors.RED}{total_count}{Colors.RESET}")
        print(f"📈 平均错误次数: {Colors.YELLOW}{avg_errors:.1f}{Colors.RESET}")
        print(f"⏰ 平均剩余时间: {Colors.CYAN}{self.format_time(int(avg_time_left))}{Colors.RESET}")
        
        print(f"\n{Colors.BOLD}按剩余时间分布:{Colors.RESET}")
        for time_range, count in time_groups.items():
            if count > 0:
                percentage = (count / total_count) * 100
                print(f"  {time_range}: {Colors.GREEN}{count}{Colors.RESET} ({percentage:.1f}%)")

    def interactive_menu(self):
        """交互式菜单"""
        while True:
            self.print_header("Dify登录锁定解除工具")
            print("1. 显示被锁定的邮箱")
            print("2. 显示详细锁定信息")
            print("3. 解除指定邮箱的锁定")
            print("4. 解除所有邮箱的锁定")
            print("5. 批量解除锁定")
            print("6. 显示统计信息")
            print("7. 导出锁定信息")
            print("8. 退出")
            
            try:
                choice = input(f"\n{Colors.WHITE}请选择操作 (1-8): {Colors.RESET}").strip()
                
                if choice == '1':
                    self.show_locked_emails()
                elif choice == '2':
                    self.show_locked_emails(detailed=True)
                elif choice == '3':
                    email = input("请输入要解除锁定的邮箱: ").strip()
                    if email:
                        self.unlock_email(email)
                    else:
                        self.print_error("邮箱不能为空")
                elif choice == '4':
                    confirm = input("确定要解除所有邮箱的锁定吗？(y/N): ").strip().lower()
                    if confirm in ['y', 'yes']:
                        self.unlock_all_emails()
                    else:
                        self.print_info("操作已取消")
                elif choice == '5':
                    emails_input = input("请输入要解除锁定的邮箱（用逗号分隔）: ").strip()
                    if emails_input:
                        emails = [email.strip() for email in emails_input.split(',') if email.strip()]
                        if emails:
                            self.unlock_batch(emails)
                        else:
                            self.print_error("请输入有效的邮箱地址")
                    else:
                        self.print_error("邮箱不能为空")
                elif choice == '6':
                    self.get_statistics()
                elif choice == '7':
                    filename = input("请输入导出文件名（留空使用默认名称）: ").strip()
                    self.export_locked_emails(filename if filename else None)
                elif choice == '8':
                    self.print_info("退出程序")
                    break
                else:
                    self.print_error("无效选择，请输入 1-8")
                
                input(f"\n{Colors.CYAN}按回车键继续...{Colors.RESET}")
                
            except KeyboardInterrupt:
                print(f"\n{Colors.YELLOW}用户中断操作{Colors.RESET}")
                break
            except Exception as e:
                self.print_error(f"操作出错: {e}")

def main():
    """主函数"""
    parser = argparse.ArgumentParser(
        description="Dify登录锁定解除工具",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例用法:
  python unlock_login.py                           # 交互式菜单
  python unlock_login.py --list                    # 显示被锁定的邮箱
  python unlock_login.py --unlock user@example.com # 解除指定邮箱的锁定
  python unlock_login.py --unlock-all              # 解除所有邮箱的锁定
  python unlock_login.py --stats                   # 显示统计信息
  python unlock_login.py --export                  # 导出锁定信息
        """
    )
    
    parser.add_argument('--list', action='store_true', help='显示所有被锁定的邮箱')
    parser.add_argument('--detailed', action='store_true', help='显示详细信息（与--list配合使用）')
    parser.add_argument('--unlock', metavar='EMAIL', help='解除指定邮箱的锁定')
    parser.add_argument('--unlock-all', action='store_true', help='解除所有邮箱的锁定')
    parser.add_argument('--batch', metavar='EMAILS', help='批量解除锁定（邮箱用逗号分隔）')
    parser.add_argument('--stats', action='store_true', help='显示统计信息')
    parser.add_argument('--export', metavar='FILENAME', nargs='?', const='', help='导出锁定信息到JSON文件')
    
    args = parser.parse_args()
    
    # 初始化解锁器
    unlocker = LoginUnlocker()
    
    # 根据参数执行对应操作
    if args.list:
        unlocker.show_locked_emails(detailed=args.detailed)
    elif args.unlock:
        unlocker.unlock_email(args.unlock)
    elif args.unlock_all:
        unlocker.unlock_all_emails()
    elif args.batch:
        emails = [email.strip() for email in args.batch.split(',') if email.strip()]
        if emails:
            unlocker.unlock_batch(emails)
        else:
            unlocker.print_error("请提供有效的邮箱地址")
    elif args.stats:
        unlocker.get_statistics()
    elif args.export is not None:
        filename = args.export if args.export else None
        unlocker.export_locked_emails(filename)
    else:
        # 默认进入交互式菜单
        unlocker.interactive_menu()

if __name__ == "__main__":
    main() 