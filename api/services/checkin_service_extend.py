from datetime import datetime, date, timedelta
from decimal import Decimal
from typing import Optional, Dict, Any, List
from sqlalchemy import and_, desc, func

from extensions.ext_database import db
from models.checkin_extend import (
    UserPointsExtend, 
    CheckinRecordExtend, 
    PointsTransactionExtend, 
    PointsExchangeExtend,
    PointsConfigExtend
)
from models.account_money_extend import AccountMoneyExtend


class CheckinServiceExtend:
    """签到服务类"""
    
    # 默认积分配置
    DEFAULT_CHECKIN_POINTS = 10  # 每日签到积分
    DEFAULT_BONUS_THRESHOLD = 7  # 连续签到奖励阈值（天数）
    DEFAULT_BONUS_POINTS = 50   # 连续签到奖励积分
    DEFAULT_POINTS_TO_QUOTA_RATE = 100  # 积分兑换额度比例（100积分=1美元）
    
    @classmethod
    def _get_config_value(cls, config_key: str, default_value: float) -> Decimal:
        """获取配置值"""
        config = db.session.query(PointsConfigExtend).filter(
            PointsConfigExtend.config_key == config_key
        ).first()
        
        if config:
            return config.config_value
        else:
            # 如果配置不存在，创建默认配置
            new_config = PointsConfigExtend(
                config_key=config_key,
                config_value=Decimal(str(default_value)),
                description=f"默认{config_key}配置"
            )
            db.session.add(new_config)
            db.session.commit()
            return Decimal(str(default_value))
    
    @classmethod
    def get_or_create_user_points(cls, account_id: str) -> UserPointsExtend:
        """获取或创建用户积分账户"""
        user_points = db.session.query(UserPointsExtend).filter(
            UserPointsExtend.account_id == account_id
        ).first()
        
        if not user_points:
            user_points = UserPointsExtend(
                account_id=account_id,
                total_points=Decimal('0'),
                available_points=Decimal('0'),
                used_points=Decimal('0')
            )
            db.session.add(user_points)
            db.session.commit()
        
        return user_points
    
    @classmethod
    def check_today_checkin(cls, account_id: str) -> bool:
        """检查今日是否已签到"""
        today = date.today()
        checkin_record = db.session.query(CheckinRecordExtend).filter(
            and_(
                CheckinRecordExtend.account_id == account_id,
                CheckinRecordExtend.checkin_date == today
            )
        ).first()
        
        return checkin_record is not None
    
    @classmethod
    def get_consecutive_days(cls, account_id: str) -> int:
        """获取连续签到天数"""
        today = date.today()
        yesterday = today - timedelta(days=1)
        
        # 检查昨天是否签到
        yesterday_checkin = db.session.query(CheckinRecordExtend).filter(
            and_(
                CheckinRecordExtend.account_id == account_id,
                CheckinRecordExtend.checkin_date == yesterday
            )
        ).first()
        
        if not yesterday_checkin:
            return 0  # 如果昨天没签到，连续天数重置为0
        
        return yesterday_checkin.consecutive_days
    
    @classmethod
    def daily_checkin(cls, account_id: str) -> Dict[str, Any]:
        """执行每日签到"""
        # 检查今日是否已签到
        if cls.check_today_checkin(account_id):
            return {
                'success': False,
                'message': '今日已签到，请明天再来！',
                'already_checked': True
            }
        
        # 获取配置
        checkin_points = cls._get_config_value('daily_checkin_points', cls.DEFAULT_CHECKIN_POINTS)
        bonus_threshold = int(cls._get_config_value('bonus_threshold', cls.DEFAULT_BONUS_THRESHOLD))
        bonus_points = cls._get_config_value('bonus_points', cls.DEFAULT_BONUS_POINTS)
        
        # 计算连续签到天数
        consecutive_days = cls.get_consecutive_days(account_id) + 1
        
        # 判断是否获得奖励
        is_bonus = consecutive_days % bonus_threshold == 0
        points_earned = checkin_points + (bonus_points if is_bonus else Decimal('0'))
        
        try:
            # 创建签到记录
            checkin_record = CheckinRecordExtend(
                account_id=account_id,
                checkin_date=date.today(),
                points_earned=points_earned,
                consecutive_days=consecutive_days,
                is_bonus=is_bonus
            )
            db.session.add(checkin_record)
            db.session.flush()  # 获取ID
            
            # 更新用户积分
            user_points = cls.get_or_create_user_points(account_id)
            points_before = user_points.available_points
            user_points.total_points += points_earned
            user_points.available_points += points_earned
            
            # 创建积分流水记录
            transaction = PointsTransactionExtend(
                account_id=account_id,
                transaction_type='checkin' if not is_bonus else 'bonus',
                points_change=points_earned,
                points_before=points_before,
                points_after=user_points.available_points,
                description=f"每日签到获得积分（连续{consecutive_days}天）" + ("，获得连续签到奖励！" if is_bonus else ""),
                related_id=checkin_record.id
            )
            db.session.add(transaction)
            
            db.session.commit()
            
            return {
                'success': True,
                'message': '签到成功！',
                'points_earned': float(points_earned),
                'consecutive_days': consecutive_days,
                'is_bonus': is_bonus,
                'total_points': float(user_points.available_points)
            }
            
        except Exception as e:
            db.session.rollback()
            return {
                'success': False,
                'message': f'签到失败：{str(e)}'
            }
    
    @classmethod
    def get_user_points_info(cls, account_id: str) -> Dict[str, Any]:
        """获取用户积分信息"""
        user_points = cls.get_or_create_user_points(account_id)
        
        # 获取今日是否已签到
        today_checked = cls.check_today_checkin(account_id)
        
        # 获取连续签到天数
        consecutive_days = cls.get_consecutive_days(account_id)
        if today_checked:
            consecutive_days += 1
        
        # 获取总签到天数
        total_checkin_days = db.session.query(func.count(CheckinRecordExtend.id)).filter(
            CheckinRecordExtend.account_id == account_id
        ).scalar() or 0
        
        return {
            'total_points': float(user_points.total_points),
            'available_points': float(user_points.available_points),
            'used_points': float(user_points.used_points),
            'today_checked': today_checked,
            'consecutive_days': consecutive_days,
            'total_checkin_days': total_checkin_days
        }
    
    @classmethod
    def get_checkin_history(cls, account_id: str, page: int = 1, page_size: int = 20) -> Dict[str, Any]:
        """获取签到历史"""
        offset = (page - 1) * page_size
        
        checkin_records = db.session.query(CheckinRecordExtend).filter(
            CheckinRecordExtend.account_id == account_id
        ).order_by(desc(CheckinRecordExtend.checkin_date)).offset(offset).limit(page_size).all()
        
        total = db.session.query(func.count(CheckinRecordExtend.id)).filter(
            CheckinRecordExtend.account_id == account_id
        ).scalar() or 0
        
        records = []
        for record in checkin_records:
            records.append({
                'checkin_date': record.checkin_date.strftime('%Y-%m-%d'),
                'points_earned': float(record.points_earned),
                'consecutive_days': record.consecutive_days,
                'is_bonus': record.is_bonus
            })
        
        return {
            'records': records,
            'total': total,
            'page': page,
            'page_size': page_size,
            'total_pages': (total + page_size - 1) // page_size
        }
    
    @classmethod
    def exchange_points_for_quota(cls, account_id: str, points_to_exchange: float) -> Dict[str, Any]:
        """积分兑换额度"""
        points_to_exchange = Decimal(str(points_to_exchange))
        
        if points_to_exchange <= 0:
            return {
                'success': False,
                'message': '兑换积分必须大于0'
            }
        
        # 获取用户积分
        user_points = cls.get_or_create_user_points(account_id)
        
        if user_points.available_points < points_to_exchange:
            return {
                'success': False,
                'message': '可用积分不足'
            }
        
        # 获取兑换比例
        exchange_rate = cls._get_config_value('points_to_quota_rate', cls.DEFAULT_POINTS_TO_QUOTA_RATE)
        quota_amount = points_to_exchange / exchange_rate
        
        try:
            # 扣除积分
            points_before = user_points.available_points
            user_points.available_points -= points_to_exchange
            user_points.used_points += points_to_exchange
            
            # 创建兑换记录
            exchange_record = PointsExchangeExtend(
                account_id=account_id,
                exchange_type='quota',
                points_cost=points_to_exchange,
                quota_amount=quota_amount,
                status='completed',
                description=f'积分兑换额度：{points_to_exchange}积分兑换{quota_amount}美元'
            )
            db.session.add(exchange_record)
            db.session.flush()
            
            # 创建积分流水记录
            transaction = PointsTransactionExtend(
                account_id=account_id,
                transaction_type='exchange',
                points_change=-points_to_exchange,
                points_before=points_before,
                points_after=user_points.available_points,
                description=f'积分兑换额度：消耗{points_to_exchange}积分，获得{quota_amount}美元',
                related_id=exchange_record.id
            )
            db.session.add(transaction)
            
            # 增加用户额度
            account_money = db.session.query(AccountMoneyExtend).filter(
                AccountMoneyExtend.account_id == account_id
            ).first()
            
            if account_money:
                account_money.total_quota += quota_amount
            else:
                # 如果用户额度记录不存在，创建新记录
                from configs import dify_config
                account_money = AccountMoneyExtend(
                    account_id=account_id,
                    total_quota=float(dify_config.DEFAULT_ACCOUNT_TOTAL_QUOTA) + float(quota_amount),
                    used_quota=Decimal('0')
                )
                db.session.add(account_money)
            
            db.session.commit()
            
            return {
                'success': True,
                'message': '兑换成功！',
                'points_used': float(points_to_exchange),
                'quota_gained': float(quota_amount),
                'remaining_points': float(user_points.available_points)
            }
            
        except Exception as e:
            db.session.rollback()
            return {
                'success': False,
                'message': f'兑换失败：{str(e)}'
            }
    
    @classmethod
    def get_exchange_history(cls, account_id: str, page: int = 1, page_size: int = 20) -> Dict[str, Any]:
        """获取兑换历史"""
        offset = (page - 1) * page_size
        
        exchange_records = db.session.query(PointsExchangeExtend).filter(
            PointsExchangeExtend.account_id == account_id
        ).order_by(desc(PointsExchangeExtend.created_at)).offset(offset).limit(page_size).all()
        
        total = db.session.query(func.count(PointsExchangeExtend.id)).filter(
            PointsExchangeExtend.account_id == account_id
        ).scalar() or 0
        
        records = []
        for record in exchange_records:
            records.append({
                'id': str(record.id),
                'exchange_type': record.exchange_type,
                'points_cost': float(record.points_cost),
                'quota_amount': float(record.quota_amount) if record.quota_amount else 0,
                'status': record.status,
                'description': record.description,
                'created_at': record.created_at.strftime('%Y-%m-%d %H:%M:%S')
            })
        
        return {
            'records': records,
            'total': total,
            'page': page,
            'page_size': page_size,
            'total_pages': (total + page_size - 1) // page_size
        }
    
    @classmethod
    def get_points_transactions(cls, account_id: str, page: int = 1, page_size: int = 20) -> Dict[str, Any]:
        """获取积分流水"""
        offset = (page - 1) * page_size
        
        transactions = db.session.query(PointsTransactionExtend).filter(
            PointsTransactionExtend.account_id == account_id
        ).order_by(desc(PointsTransactionExtend.created_at)).offset(offset).limit(page_size).all()
        
        total = db.session.query(func.count(PointsTransactionExtend.id)).filter(
            PointsTransactionExtend.account_id == account_id
        ).scalar() or 0
        
        records = []
        for transaction in transactions:
            records.append({
                'id': str(transaction.id),
                'transaction_type': transaction.transaction_type,
                'points_change': float(transaction.points_change),
                'points_before': float(transaction.points_before),
                'points_after': float(transaction.points_after),
                'description': transaction.description,
                'created_at': transaction.created_at.strftime('%Y-%m-%d %H:%M:%S')
            })
        
        return {
            'records': records,
            'total': total,
            'page': page,
            'page_size': page_size,
            'total_pages': (total + page_size - 1) // page_size
        } 