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

// OrgProject holds the schema definition for the OrgProject entity.
type OrgProject struct {
	ent.Schema
}

func (OrgProject) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "org_projects"},
	}
}

func (OrgProject) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

func (OrgProject) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("org_id"),
		field.String("name").MaxLen(200).NotEmpty(),
		field.String("description").Optional().Nillable().SchemaType(map[string]string{dialect.Postgres: "text"}),
		field.Int64("group_id").Optional().Nillable(),
		field.JSON("allowed_models", []string{}).Optional(),
		field.Float("monthly_budget_usd").Optional().Nillable().SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}),
		field.Float("monthly_usage_usd").SchemaType(map[string]string{dialect.Postgres: "decimal(20,10)"}).Default(0),
		field.Time("monthly_window_start").Optional().Nillable(),
		field.String("status").MaxLen(20).Default(service.OrgStatusActive),
	}
}

func (OrgProject) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("organization", Organization.Type).Ref("projects").Field("org_id").Unique().Required(),
		edge.From("group", Group.Type).Ref("org_projects").Field("group_id").Unique(),
		edge.To("api_keys", APIKey.Type),
	}
}

func (OrgProject) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("org_id"),
		index.Fields("org_id", "name"),
		index.Fields("deleted_at"),
	}
}
