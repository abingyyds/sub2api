package schema

import (
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Group holds the schema definition for the Group entity.
type Group struct {
	ent.Schema
}

func (Group) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "groups"},
	}
}

func (Group) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

func (Group) Fields() []ent.Field {
	return []ent.Field{
		// 唯一约束通过部分索引实现（WHERE deleted_at IS NULL），支持软删除后重用
		// 见迁移文件 016_soft_delete_partial_unique_indexes.sql
		field.String("name").
			MaxLen(100).
			NotEmpty(),
		field.String("description").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "text"}),
		field.Float("rate_multiplier").
			SchemaType(map[string]string{dialect.Postgres: "decimal(10,4)"}).
			Default(1.0),
		field.Float("display_rate_multiplier").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "decimal(10,4)"}).
			Comment("仅用于页面价格展示的倍率，未设置时回退到 rate_multiplier"),
		field.Bool("is_exclusive").
			Default(false),
		field.String("status").
			MaxLen(20).
			Default(service.StatusActive),

		// Subscription-related fields (added by migration 003)
		field.String("platform").
			MaxLen(50).
			Default(service.PlatformAnthropic),
		field.String("subscription_type").
			MaxLen(20).
			Default(service.SubscriptionTypeStandard),
		field.Float("daily_limit_usd").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}),
		field.Float("weekly_limit_usd").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}),
		field.Float("monthly_limit_usd").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}),
		field.Int("default_validity_days").
			Default(30),

		// 图片生成计费配置（antigravity 和 gemini 平台使用）
		field.Float("image_price_1k").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}),
		field.Float("image_price_2k").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}),
		field.Float("image_price_4k").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}),

		// Claude Code 客户端限制 (added by migration 029)
		field.Bool("claude_code_only").
			Default(false).
			Comment("是否仅允许 Claude Code 客户端"),
		field.Int64("fallback_group_id").
			Optional().
			Nillable().
			Comment("非 Claude Code 请求降级使用的分组 ID"),

		// 模型路由配置 (added by migration 040)
		field.JSON("model_routing", map[string][]int64{}).
			Optional().
			SchemaType(map[string]string{dialect.Postgres: "jsonb"}).
			Comment("模型路由配置：模型模式 -> 优先账号ID列表"),

		// 模型路由开关 (added by migration 041)
		field.Bool("model_routing_enabled").
			Default(false).
			Comment("是否启用模型路由配置"),

		// 套餐价格（分），>0 时分组自动显示为可购买套餐 (added by migration 061)
		field.Int("price_fen").
			Default(0).
			Comment("套餐价格（分），大于0时自动显示为可购买套餐"),

		// 是否上架到购买页面 (added by migration 062)
		field.Bool("listed").
			Default(false).
			Comment("是否上架到购买页面展示"),

		// 套餐自定义特性列表 (added by migration 063)
		field.JSON("plan_features", []string{}).
			Optional().
			SchemaType(map[string]string{dialect.Postgres: "jsonb"}).
			Comment("套餐自定义特性列表，上架时展示在购买页面"),

		// 分组卡片展示字段 (added by migration 073)
		field.JSON("tags", []string{}).
			Optional().
			SchemaType(map[string]string{dialect.Postgres: "jsonb"}).
			Comment("自定义标签列表，如「官方 API」「逆向」「推荐」「暂不可用」等"),

		// 是否在模型广场展示 (added by migration 087)
		field.Bool("model_plaza_visible").
			Default(true).
			Comment("是否在模型广场展示该分组"),
		field.String("display_price").
			Default("").
			SchemaType(map[string]string{dialect.Postgres: "text"}).
			Comment("展示价格文案，如「6 块 / 1 美元」"),
		field.String("display_discount").
			Default("").
			SchemaType(map[string]string{dialect.Postgres: "text"}).
			Comment("展示折扣文案，如「8.3折」"),
	}
}

func (Group) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("api_keys", APIKey.Type),
		edge.To("redeem_codes", RedeemCode.Type),
		edge.To("subscriptions", UserSubscription.Type),
		edge.To("org_subscriptions", OrgSubscription.Type),
		edge.To("org_projects", OrgProject.Type),
		edge.To("usage_logs", UsageLog.Type),
		edge.From("accounts", Account.Type).
			Ref("groups").
			Through("account_groups", AccountGroup.Type),
		edge.From("allowed_users", User.Type).
			Ref("allowed_groups").
			Through("user_allowed_groups", UserAllowedGroup.Type),
		// 注意：fallback_group_id 直接作为字段使用，不定义 edge
		// 这样允许多个分组指向同一个降级分组（M2O 关系）
	}
}

func (Group) Indexes() []ent.Index {
	return []ent.Index{
		// name 字段已在 Fields() 中声明 Unique()，无需重复索引
		index.Fields("status"),
		index.Fields("platform"),
		index.Fields("subscription_type"),
		index.Fields("is_exclusive"),
		index.Fields("deleted_at"),
	}
}
