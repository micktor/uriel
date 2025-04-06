package schema

import (
	"context"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/cuvva/cuvva-public-go/lib/ksuid"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			NotEmpty().
			Immutable().
			DefaultFunc(func() string {
				return ksuid.Generate(context.Background(), "user").String()
			}),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now),
		field.Time("deleted_at").
			Optional(),
		field.String("email").
			Unique().
			NotEmpty(),
		field.String("password").NotEmpty(),
		field.String("auth_provider").
			Optional(),
		field.String("oauth_id").
			Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{}
}
