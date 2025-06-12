-- 签到积分系统数据库表结构

-- 用户积分账户表
CREATE TABLE IF NOT EXISTS user_points_extend (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL UNIQUE,
    total_points DECIMAL(16, 7) DEFAULT 0.0 NOT NULL,
    available_points DECIMAL(16, 7) DEFAULT 0.0 NOT NULL,
    used_points DECIMAL(16, 7) DEFAULT 0.0 NOT NULL,
    created_at TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(0) NOT NULL,
    updated_at TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(0) NOT NULL
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_user_points_account_id ON user_points_extend(account_id);

-- 添加外键约束
-- ALTER TABLE user_points_extend ADD CONSTRAINT fk_user_points_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE;

-- 用户签到记录表
CREATE TABLE IF NOT EXISTS checkin_record_extend (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL,
    checkin_date DATE NOT NULL,
    points_earned DECIMAL(16, 7) NOT NULL,
    consecutive_days INTEGER DEFAULT 1 NOT NULL,
    is_bonus BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(0) NOT NULL
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_checkin_record_account_id ON checkin_record_extend(account_id);
CREATE INDEX IF NOT EXISTS idx_checkin_record_checkin_date ON checkin_record_extend(checkin_date);

-- 创建唯一约束，确保每个用户每天只能签到一次
CREATE UNIQUE INDEX IF NOT EXISTS uk_account_checkin_date ON checkin_record_extend(account_id, checkin_date);

-- 积分流水记录表
CREATE TABLE IF NOT EXISTS points_transaction_extend (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL,
    transaction_type VARCHAR(50) NOT NULL,  -- checkin, exchange, bonus, manual
    points_change DECIMAL(16, 7) NOT NULL,
    points_before DECIMAL(16, 7) NOT NULL,
    points_after DECIMAL(16, 7) NOT NULL,
    description VARCHAR(200),
    related_id UUID,
    created_at TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(0) NOT NULL
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_points_transaction_account_id ON points_transaction_extend(account_id);
CREATE INDEX IF NOT EXISTS idx_points_transaction_type ON points_transaction_extend(transaction_type);
CREATE INDEX IF NOT EXISTS idx_points_transaction_created_at ON points_transaction_extend(created_at);

-- 积分兑换记录表
CREATE TABLE IF NOT EXISTS points_exchange_extend (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL,
    exchange_type VARCHAR(50) NOT NULL,  -- quota
    points_cost DECIMAL(16, 7) NOT NULL,
    quota_amount DECIMAL(16, 7),
    status VARCHAR(20) DEFAULT 'completed' NOT NULL,  -- pending, completed, failed
    description VARCHAR(200),
    created_at TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(0) NOT NULL,
    updated_at TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(0) NOT NULL
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_points_exchange_account_id ON points_exchange_extend(account_id);
CREATE INDEX IF NOT EXISTS idx_points_exchange_status ON points_exchange_extend(status);

-- 积分配置表
CREATE TABLE IF NOT EXISTS points_config_extend (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    config_key VARCHAR(100) NOT NULL UNIQUE,
    config_value DECIMAL(16, 7) NOT NULL,
    description VARCHAR(200),
    created_at TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(0) NOT NULL,
    updated_at TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(0) NOT NULL
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_points_config_config_key ON points_config_extend(config_key);

-- 插入默认积分配置
INSERT INTO points_config_extend (config_key, config_value, description) VALUES
('daily_checkin_points', 10.0, '每日签到基础积分'),
('consecutive_bonus_days', 7.0, '连续签到奖励间隔天数'),
('consecutive_bonus_points', 50.0, '连续签到奖励积分'),
('points_to_quota_rate', 100.0, '积分兑换额度比例(100积分=1美元)')
ON CONFLICT (config_key) DO NOTHING;

-- 添加注释
COMMENT ON TABLE user_points_extend IS '用户积分账户表';
COMMENT ON COLUMN user_points_extend.account_id IS '关联用户账户ID';
COMMENT ON COLUMN user_points_extend.total_points IS '总积分(累计获得)';
COMMENT ON COLUMN user_points_extend.available_points IS '可用积分';
COMMENT ON COLUMN user_points_extend.used_points IS '已使用积分';

COMMENT ON TABLE checkin_record_extend IS '用户签到记录表';
COMMENT ON COLUMN checkin_record_extend.account_id IS '关联用户账户ID';
COMMENT ON COLUMN checkin_record_extend.checkin_date IS '签到日期';
COMMENT ON COLUMN checkin_record_extend.points_earned IS '签到获得积分';
COMMENT ON COLUMN checkin_record_extend.consecutive_days IS '连续签到天数';
COMMENT ON COLUMN checkin_record_extend.is_bonus IS '是否为奖励签到';

COMMENT ON TABLE points_transaction_extend IS '积分流水记录表';
COMMENT ON COLUMN points_transaction_extend.account_id IS '关联用户账户ID';
COMMENT ON COLUMN points_transaction_extend.transaction_type IS '交易类型';
COMMENT ON COLUMN points_transaction_extend.points_change IS '积分变化(正数增加,负数减少)';
COMMENT ON COLUMN points_transaction_extend.points_before IS '交易前积分';
COMMENT ON COLUMN points_transaction_extend.points_after IS '交易后积分';
COMMENT ON COLUMN points_transaction_extend.related_id IS '关联记录ID';

COMMENT ON TABLE points_exchange_extend IS '积分兑换记录表';
COMMENT ON COLUMN points_exchange_extend.account_id IS '关联用户账户ID';
COMMENT ON COLUMN points_exchange_extend.exchange_type IS '兑换类型';
COMMENT ON COLUMN points_exchange_extend.points_cost IS '消耗积分';
COMMENT ON COLUMN points_exchange_extend.quota_amount IS '兑换的额度金额(美元)';
COMMENT ON COLUMN points_exchange_extend.status IS '兑换状态';

COMMENT ON TABLE points_config_extend IS '积分配置表';
COMMENT ON COLUMN points_config_extend.config_key IS '配置键';
COMMENT ON COLUMN points_config_extend.config_value IS '配置值';
COMMENT ON COLUMN points_config_extend.description IS '配置描述'; 