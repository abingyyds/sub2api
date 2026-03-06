-- 057: 添加用户来源字段
-- 用于记录未通过邀请链接注册的用户是通过什么渠道了解到本平台的

ALTER TABLE users ADD COLUMN IF NOT EXISTS discovery_source VARCHAR(50);
