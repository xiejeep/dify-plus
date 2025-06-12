from datetime import datetime, date
from extensions.ext_database import db
from .types import StringUUID


class AccountCheckinExtend(db.Model):
    """用户签到记录表"""
    __tablename__ = "account_checkin_extend"
    __table_args__ = (
        db.PrimaryKeyConstraint("id", name="account_checkin_pkey"),
        db.Index("idx_account_checkin_account_id", "account_id"),
        db.Index("idx_account_checkin_date", "checkin_date"),
        db.UniqueConstraint("account_id", "checkin_date", name="unique_account_checkin_date")
    )

    id = db.Column(StringUUID, server_default=db.text("uuid_generate_v4()"))
    account_id = db.Column(StringUUID, nullable=False)
    checkin_date = db.Column(db.Date, nullable=False)  # 签到日期
    reward_amount = db.Column(db.Numeric(16, 7), nullable=False)  # 奖励金额
    created_at = db.Column(db.DateTime, nullable=False, server_default=db.text("CURRENT_TIMESTAMP(0)"))

    @classmethod
    def is_checked_in_today(cls, account_id):
        """检查用户今天是否已经签到"""
        today = date.today()
        return db.session.query(cls).filter(
            cls.account_id == account_id,
            cls.checkin_date == today
        ).first() is not None

    @classmethod 
    def get_checkin_streak(cls, account_id):
        """获取用户连续签到天数"""
        checkin_records = db.session.query(cls).filter(
            cls.account_id == account_id
        ).order_by(cls.checkin_date.desc()).limit(30).all()
        
        if not checkin_records:
            return 0
            
        streak = 0
        today = date.today()
        current_date = today
        
        for record in checkin_records:
            if record.checkin_date == current_date:
                streak += 1
                current_date = current_date.replace(day=current_date.day - 1)
            else:
                break
                
        return streak

    @classmethod
    def create_checkin_record(cls, account_id, reward_amount=0.1):
        """创建签到记录"""
        today = date.today()
        
        # 检查今天是否已经签到
        if cls.is_checked_in_today(account_id):
            return None
            
        checkin_record = cls(
            account_id=account_id,
            checkin_date=today,
            reward_amount=reward_amount
        )
        db.session.add(checkin_record)
        return checkin_record 