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

// OrgMember holds the schema definition for the OrgMember entity.
type OrgMember struct {
	ent.Schema
}

func (OrgMember) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "org_members"},
	}
}

func (OrgMember) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

func (OrgMember) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("org_id"),
		field.Int64("user_id"),
		field.String("role").
			MaxLen(20).
			Default(service.OrgMemberRoleMember),
		field.Float("monthly_quota_usd").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}),
		field.Float("daily_quota_usd").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}),
		field.Float("monthly_usage_usd").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,10)"}).
			Default(0),
		field.Float("daily_usage_usd").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,10)"}).
			Default(0),
		field.Time("monthly_window_start").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
		field.Time("daily_window_start").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
		field.String("status").
			MaxLen(20).
			Default(service.StatusActive),
		field.String("notes").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "text"}),
	}
}

func (OrgMember) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("organization", Organization.Type).
			Ref("members").
			Field("org_id").
			Unique().
			Required(),
		edge.From("user", User.Type).
			Ref("org_memberships").
			Field("user_id").
			Unique().
			Required(),
	}
}

func (OrgMember) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("org_id"),
		index.Fields("user_id"),
		index.Fields("org_id", "user_id"),
		index.Fields("deleted_at"),
	}
}
