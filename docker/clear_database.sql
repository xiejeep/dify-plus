-- 清空数据库脚本 - 用于解决初始化失败问题
-- 使用 TRUNCATE TABLE ... RESTART IDENTITY CASCADE 命令完全清空相关表

-- 禁用外键检查（如果需要）
SET session_replication_role = replica;

-- 清空权限角色表
TRUNCATE TABLE sys_authorities RESTART IDENTITY CASCADE;

-- 清空API表  
TRUNCATE TABLE sys_apis RESTART IDENTITY CASCADE;

-- 清空权限规则表
TRUNCATE TABLE casbin_rule RESTART IDENTITY CASCADE;

-- 清空用户表（如果需要完全重置）
TRUNCATE TABLE sys_users RESTART IDENTITY CASCADE;

-- 清空角色表
TRUNCATE TABLE sys_base_menus RESTART IDENTITY CASCADE;

-- 清空字典表
TRUNCATE TABLE sys_dictionaries RESTART IDENTITY CASCADE;
TRUNCATE TABLE sys_dictionary_details RESTART IDENTITY CASCADE;

-- 清空操作记录表
TRUNCATE TABLE sys_operation_records RESTART IDENTITY CASCADE;

-- 清空JWT黑名单表
TRUNCATE TABLE jwt_blacklists RESTART IDENTITY CASCADE;

-- 清空文件上传记录表（如果存在）
TRUNCATE TABLE exa_file_upload_and_downloads RESTART IDENTITY CASCADE;

-- 清空文件块表（如果存在）
TRUNCATE TABLE exa_file_chunks RESTART IDENTITY CASCADE;

-- 清空客户表（如果存在）
TRUNCATE TABLE exa_customers RESTART IDENTITY CASCADE;

-- 清空邮箱验证表（我们新增的）
TRUNCATE TABLE sys_email_verifications RESTART IDENTITY CASCADE;

-- 重新启用外键检查
SET session_replication_role = DEFAULT;

-- 显示完成信息
SELECT 'Database cleared successfully! All system tables have been truncated.' as result; 