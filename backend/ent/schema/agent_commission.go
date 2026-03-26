package schema

import (
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// AgentCommission holds the schema definition for the AgentCommission entity.
type AgentCommission struct {
	ent.Schema
}

func (AgentCommission) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "agent_commissions"},
	}
}

func (AgentCommission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (AgentCommission) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("agent_id"),
		field.Int64("user_id"),
		field.Int64("order_id").
			Optional().
			Nillable(),
		field.String("source_type").
			MaxLen(20).
			NotEmpty(),
		field.Float("source_amount").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0),
		field.Float("commission_rate").
			SchemaType(map[string]string{dialect.Postgres: "decimal(5,4)"}).
			Default(0),
		field.Float("commission_amount").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0),
		field.String("status").
			MaxLen(20).
			Default("pending"),
		field.Time("settled_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
	}
}

func (AgentCommission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("agent", User.Type).
			Ref("agent_commissions_as_agent").
			Field("agent_id").
			Required().
			Unique(),
		edge.From("user", User.Type).
			Ref("agent_commissions_as_user").
			Field("user_id").
			Required().
			Unique(),
		edge.From("order", PaymentOrder.Type).
			Ref("agent_commissions").
			Field("order_id").
			Unique(),
	}
}

func (AgentCommission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("agent_id"),
		index.Fields("user_id"),
		index.Fields("order_id"),
		index.Fields("status"),
		index.Fields("created_at"),
	}
}
