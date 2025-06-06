-- 创建邮箱验证码表 (PostgreSQL版本)
CREATE TABLE IF NOT EXISTS sys_email_verifications (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    email VARCHAR(255) NOT NULL,
    code VARCHAR(6) NOT NULL,
    type INTEGER NOT NULL DEFAULT 1,
    expired_at TIMESTAMP NOT NULL,
    used BOOLEAN NOT NULL DEFAULT FALSE,
    ip VARCHAR(45) DEFAULT NULL
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_email ON sys_email_verifications(email);
CREATE INDEX IF NOT EXISTS idx_sys_email_verifications_deleted_at ON sys_email_verifications(deleted_at);

-- 验证表创建成功
SELECT 'sys_email_verifications table created successfully' as result; 