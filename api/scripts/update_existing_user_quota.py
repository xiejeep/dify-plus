#!/usr/bin/env python3
"""
更新现有用户默认额度的脚本

使用方法:
    python update_existing_user_quota.py --quota 50.0
    或
    python update_existing_user_quota.py --quota 50.0 --confirm

参数:
    --quota: 新的默认额度值（美元）
    --confirm: 确认执行更新（不加此参数只会显示预览）
"""

import argparse
import sys
import os

# 添加项目根目录到 Python 路径
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from extensions.ext_database import db
from models.account_money_extend import AccountMoneyExtend
from app import create_app


def update_user_quota(new_quota: float, confirm: bool = False):
    """更新所有用户的默认额度"""
    
    app = create_app()
    
    with app.app_context():
        # 查询所有用户的额度记录
        account_moneys = db.session.query(AccountMoneyExtend).all()
        
        if not account_moneys:
            print("❌ 没有找到任何用户额度记录")
            return
        
        print(f"📊 找到 {len(account_moneys)} 个用户额度记录")
        print("\n当前额度情况:")
        print("-" * 60)
        print(f"{'账户ID':<40} {'当前总额度':<10} {'已使用额度':<10}")
        print("-" * 60)
        
        total_updated = 0
        for account_money in account_moneys:
            account_id_str = str(account_money.account_id)[:8] + "..."
            print(f"{account_id_str:<40} {account_money.total_quota:<10} {account_money.used_quota:<10}")
            if float(account_money.total_quota) != new_quota:
                total_updated += 1
        
        print("-" * 60)
        print(f"📈 需要更新的记录数: {total_updated}")
        print(f"🎯 新的默认额度: {new_quota} 美元")
        
        if not confirm:
            print("\n⚠️  这是预览模式，如需执行更新，请添加 --confirm 参数")
            print("示例: python update_existing_user_quota.py --quota 50.0 --confirm")
            return
        
        # 执行更新
        print(f"\n🔄 开始更新用户额度...")
        
        try:
            updated_count = db.session.query(AccountMoneyExtend).update(
                {AccountMoneyExtend.total_quota: new_quota}
            )
            db.session.commit()
            
            print(f"✅ 成功更新 {updated_count} 个用户的额度为 {new_quota} 美元")
            
        except Exception as e:
            db.session.rollback()
            print(f"❌ 更新失败: {str(e)}")
            return False
        
        return True


def main():
    parser = argparse.ArgumentParser(description='更新现有用户的默认额度')
    parser.add_argument('--quota', type=float, required=True, 
                       help='新的默认额度值（美元）')
    parser.add_argument('--confirm', action='store_true', 
                       help='确认执行更新（不加此参数只显示预览）')
    
    args = parser.parse_args()
    
    if args.quota <= 0:
        print("❌ 额度值必须大于0")
        sys.exit(1)
    
    print(f"🚀 用户额度更新工具")
    print(f"目标额度: {args.quota} 美元")
    print(f"确认模式: {'是' if args.confirm else '否（预览模式）'}")
    print()
    
    success = update_user_quota(args.quota, args.confirm)
    
    if success and args.confirm:
        print("\n🎉 更新完成！新注册的用户将自动获得新的默认额度。")
        print("💡 提示：需要重启 API 和 Worker 服务以确保配置生效")


if __name__ == "__main__":
    main() 