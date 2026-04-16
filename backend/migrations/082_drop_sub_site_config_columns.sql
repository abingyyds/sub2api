-- 移除分站两个废弃的 JSON 配置字段
-- theme_config / custom_config 自始至终只在 HomeView 合并为 4 个键（hero_title / hero_description / cta_text / feature_tags），
-- 功能上两列重复、分站主不可编辑、键空间极窄；改由独立主题组件 + 分站基础字段（name/subtitle/home_content）覆盖。
-- 同步移除平台级默认值 setting。

ALTER TABLE sub_sites
    DROP COLUMN IF EXISTS theme_config,
    DROP COLUMN IF EXISTS custom_config;

DELETE FROM settings WHERE key = 'subsite_default_custom_config';
