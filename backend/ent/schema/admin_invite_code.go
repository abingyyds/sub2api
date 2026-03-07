package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

// AdminInviteCode holds the schema definition for the AdminInviteCode entity.
type AdminInviteCode struct {
	ent.Schema
}

// Fields of the AdminInviteCode.
func (AdminInviteCode) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("code").
			MaxLen(50).
			Unique().
			NotEmpty(),
		field.String("source_name").
			MaxLen(100).
			NotEmpty(),
		field.Int64("created_by"),
		field.Int("used_count").
			Default(0).
			NonNegative(),
		field.Int("max_uses").
			Optional().
			Nillable().
			Positive(),
		field.Bool("enabled").
			Default(true),
		field.Text("notes").
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the AdminInviteCode.
func (AdminInviteCode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("creator", User.Type).
			Ref("admin_invite_codes").
			Field("created_by").
			Unique().
			Required(),
	}
}

// Indexes of the AdminInviteCode.
func (AdminInviteCode) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code"),
		index.Fields("created_by"),
		index.Fields("enabled"),
	}
}
