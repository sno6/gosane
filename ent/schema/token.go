package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
	"github.com/google/uuid"
)

type Token struct {
	ent.Schema
}

func (Token) Config() ent.Config {
	return ent.Config{
		Table: "token",
	}
}

func (Token) Mixin() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

func (Token) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("uuid", uuid.UUID{}).Default(uuid.New).Unique(),
		field.String("refresh_token"),
		field.Time("access_expires_at"),
		field.Time("refresh_expires_at"),
	}
}

func (Token) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("token").Unique(),
	}
}
