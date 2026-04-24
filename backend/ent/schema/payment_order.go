package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PaymentOrder holds the schema definition for the PaymentOrder entity.
type PaymentOrder struct {
	ent.Schema
}

func (PaymentOrder) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "payment_orders"},
	}
}

// Fields of the PaymentOrder.
func (PaymentOrder) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("order_no").
			MaxLen(64).
			Unique().
			NotEmpty(),
		field.Int64("user_id"),
		field.String("plan_key").
			MaxLen(50).
			NotEmpty(),
		field.Int64("group_id"),
		field.Int("amount_fen"),
		field.Int("validity_days").
			Default(30),
		field.String("order_type").
			MaxLen(20).
			Default("subscription"),
		field.Float("balance_amount").
			Default(0),
		field.Int64("sub_site_id").
			Optional().
			Nillable().
			Comment("下单时所属分站，NULL 表示主站订单"),
		field.String("promo_code").
			MaxLen(32).
			Optional().
			Default("").
			Comment("使用的优惠码"),
		field.Int("discount_amount").
			Default(0).
			Comment("优惠码折扣金额(分)"),
		field.String("status").
			MaxLen(20).
			Default("pending"),
		field.String("pay_method").
			MaxLen(20).
			Default("wechat_native").
			Comment("wechat_native | alipay_native | epay_alipay | epay_wxpay"),
		field.String("wechat_transaction_id").
			MaxLen(64).
			Optional().
			Nillable().
			Comment("微信交易号"),
		field.String("alipay_trade_no").
			MaxLen(64).
			Optional().
			Nillable().
			Comment("支付宝交易号"),
		field.String("epay_trade_no").
			MaxLen(64).
			Optional().
			Nillable().
			Comment("易支付交易号"),
		field.String("invoice_company_name").
			MaxLen(255).
			Optional().
			Default("").
			Comment("发票抬头/公司名称"),
		field.String("invoice_tax_id").
			MaxLen(128).
			Optional().
			Default("").
			Comment("发票税号"),
		field.String("invoice_email").
			MaxLen(255).
			Optional().
			Default("").
			Comment("接收发票邮箱"),
		field.Text("invoice_remark").
			Optional().
			Default("").
			Comment("发票备注"),
		field.Time("invoice_requested_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}).
			Comment("发票申请时间"),
		field.Time("invoice_processed_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}).
			Comment("发票处理时间"),
		field.Text("code_url").
			Optional().
			Nillable(),
		field.Time("paid_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
		field.Time("expired_at").
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
	}
}

// Edges of the PaymentOrder.
func (PaymentOrder) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("payment_orders").
			Field("user_id").
			Unique().
			Required(),
		edge.To("agent_commissions", AgentCommission.Type),
	}
}

// Indexes of the PaymentOrder.
func (PaymentOrder) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("order_no"),
		index.Fields("status"),
		index.Fields("sub_site_id"),
		index.Fields("invoice_requested_at"),
	}
}
