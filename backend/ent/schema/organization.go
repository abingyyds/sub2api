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

// Organization holds the schema definition for the Organization entity.
type Organization struct {
	ent.Schema
}

func (Organization) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "organizations"},
	}
}

func (Organization) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

func (Organization) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(200).
			NotEmpty(),
		field.String("slug").
			MaxLen(100).
			NotEmpty(),
		field.String("description").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "text"}),
		field.Int64("owner_user_id"),
		field.String("billing_mode").
			MaxLen(20).
			Default(service.OrgBillingModeBalance),
		field.Float("balance").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0),
		field.Float("monthly_budget_usd").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}),
		field.Int("max_members").
			Default(50),
		field.Int("max_api_keys").
			Default(100),
		field.String("status").
			MaxLen(20).
			Default(service.OrgStatusActive),
		field.String("audit_mode").
			MaxLen(20).
			Default(service.OrgAuditModeMetadata),
	}
}

func (Organization) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("owned_organizations").
			Field("owner_user_id").
			Unique().
			Required(),
		edge.To("members", OrgMember.Type),
		edge.To("subscriptions", OrgSubscription.Type),
		edge.To("api_keys", APIKey.Type),
		edge.To("projects", OrgProject.Type),
		edge.To("audit_logs", OrgAuditLog.Type),
	}
}

func (Organization) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("slug"),
		index.Fields("owner_user_id"),
		index.Fields("status"),
		index.Fields("deleted_at"),
	}
}
