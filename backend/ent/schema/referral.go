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

// Referral holds the schema definition for the Referral entity.
type Referral struct {
	ent.Schema
}

func (Referral) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "referrals"},
	}
}

func (Referral) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (Referral) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("inviter_id"),
		field.Int64("invitee_id"),
		field.String("reward_status").
			MaxLen(20).
			Default("pending"),
		field.Float("reward_amount").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0),
		field.Time("rewarded_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
	}
}

func (Referral) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("inviter", User.Type).
			Ref("referrals_as_inviter").
			Field("inviter_id").
			Required().
			Unique(),
		edge.From("invitee", User.Type).
			Ref("referrals_as_invitee").
			Field("invitee_id").
			Required().
			Unique(),
	}
}

func (Referral) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("invitee_id").Unique(),
		index.Fields("inviter_id"),
		index.Fields("reward_status"),
	}
}
