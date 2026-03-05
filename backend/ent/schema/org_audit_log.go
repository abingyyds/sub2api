package schema

import (
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// OrgAuditLog holds the schema definition for the OrgAuditLog entity.
// Note: No SoftDeleteMixin — audit logs are immutable and cannot be deleted.
type OrgAuditLog struct {
	ent.Schema
}

func (OrgAuditLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "org_audit_logs"},
	}
}

func (OrgAuditLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("org_id"),
		field.Int64("user_id"),
		field.Int64("member_id").Optional().Nillable(),
		field.Int64("project_id").Optional().Nillable(),
		field.Int64("usage_log_id").Optional().Nillable(),
		field.String("action").MaxLen(50).NotEmpty(),
		field.String("model").MaxLen(100).Optional().Nillable(),
		field.String("audit_mode").MaxLen(20).Default(service.OrgAuditModeMetadata),
		field.String("request_summary").Optional().Nillable().SchemaType(map[string]string{dialect.Postgres: "text"}),
		field.String("request_content").Optional().Nillable().SchemaType(map[string]string{dialect.Postgres: "text"}),
		field.String("response_summary").Optional().Nillable().SchemaType(map[string]string{dialect.Postgres: "text"}),
		field.JSON("keywords", []string{}).Optional(),
		field.Bool("flagged").Default(false),
		field.String("flag_reason").Optional().Nillable().SchemaType(map[string]string{dialect.Postgres: "text"}),
		field.Int("input_tokens").Optional().Nillable(),
		field.Int("output_tokens").Optional().Nillable(),
		field.Float("cost_usd").Optional().Nillable().SchemaType(map[string]string{dialect.Postgres: "decimal(20,10)"}),
		field.String("ip_address").MaxLen(45).Optional().Nillable(),
		field.String("user_agent").MaxLen(512).Optional().Nillable(),
		field.JSON("detail", map[string]interface{}{}).Optional(),
		field.Time("created_at").Immutable().Default(time.Now).SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
	}
}

func (OrgAuditLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("org_id"),
		index.Fields("user_id"),
		index.Fields("project_id"),
		index.Fields("created_at"),
		index.Fields("org_id", "flagged"),
		index.Fields("action"),
	}
}
