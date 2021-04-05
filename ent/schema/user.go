package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
	"github.com/google/uuid"
)

type OAuthProvider string

const (
	GoogleProvider   OAuthProvider = "google"
	FacebookProvider OAuthProvider = "facebook"
)

func (op OAuthProvider) Values() []string {
	return []string{string(GoogleProvider), string(FacebookProvider)}
}

type User struct {
	ent.Schema
}

func (User) Config() ent.Config {
	return ent.Config{
		Table: "user",
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

func (u User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("uuid", uuid.UUID{}).Default(uuid.New).Unique(),
		field.String("email"),
		field.Bool("email_verified").Default(false),
		field.String("password").Optional(),
		field.String("provider_id").Optional(),
		field.Enum("provider_type").Values(OAuthProvider("").Values()...).Optional().Nillable(),
		field.String("first_name").Optional(),
		field.String("last_name").Optional(),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("token", Token.Type),
	}
}
