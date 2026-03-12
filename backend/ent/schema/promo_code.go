package schema

import (
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PromoCode holds the schema definition for the PromoCode entity.
//
// 购买优惠码：用户购买时使用，可获得折扣（满减或打折）
// 与 RedeemCode 不同，PromoCode 支持多次使用（有使用次数限制）
//
// 删除策略：硬删除
type PromoCode struct {
	ent.Schema
}

func (PromoCode) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promo_codes"},
	}
}

func (PromoCode) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			MaxLen(32).
			NotEmpty().
			Unique().
			Comment("优惠码"),
		field.Float("discount_amount").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0).
			Comment("折扣金额: fixed时为分, percentage时为百分比(如10=减10%)"),
		field.String("discount_type").
			MaxLen(20).
			Default("fixed").
			Comment("折扣类型: fixed(固定减免/分), percentage(百分比折扣)"),
		field.Int("min_order_amount").
			Default(0).
			Comment("最低订单金额(分), 0表示无门槛"),
		field.Int("max_uses").
			Default(0).
			Comment("最大使用次数，0表示无限制"),
		field.Int("used_count").
			Default(0).
			Comment("已使用次数"),
		field.String("status").
			MaxLen(20).
			Default(service.PromoCodeStatusActive).
			Comment("状态: active, disabled"),
		field.Time("expires_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}).
			Comment("过期时间，null表示永不过期"),
		field.String("notes").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "text"}).
			Comment("备注"),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
	}
}

func (PromoCode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("usage_records", PromoCodeUsage.Type),
	}
}

func (PromoCode) Indexes() []ent.Index {
	return []ent.Index{
		// code 字段已在 Fields() 中声明 Unique()，无需重复索引
		index.Fields("status"),
		index.Fields("expires_at"),
	}
}
