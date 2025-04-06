package repository

import (
	"context"
	"github.com/your_org/uriel/internal/ent"
	"github.com/your_org/uriel/internal/ent/user"
)

type Repository struct {
	db *ent.Client
}

func NewRepository(db *ent.Client) *Repository {
	return &Repository{
		db: db,
	}
}

type CreateUserInput struct {
	Email        string
	Password     string
	AuthProvider string
	OAuthID      string
}

func (p Repository) CreateUser(ctx context.Context, input CreateUserInput) (ent.User, error) {
	user, err := p.db.User.Create().
		SetEmail(input.Email).
		SetPassword(input.Password).
		SetAuthProvider(input.AuthProvider).
		SetOauthID(input.OAuthID).
		Save(ctx)

	return *user, err
}

func (p *Repository) UserByEmailExists(ctx context.Context, email string) (bool, error) {
	return p.db.User.Query().Where(user.Email(email)).Exist(ctx)
}

func (p *Repository) GetUserByEmail(ctx context.Context, email string) (ent.User, error) {
	user, err := p.db.User.Query().Where(user.Email(email)).First(ctx)
	if err != nil {
		return ent.User{}, err
	}

	return *user, err
}

func (p *Repository) GetUserByID(ctx context.Context, id string) (ent.User, error) {
	user, err := p.db.User.Query().Where(user.ID(id)).First(ctx)
	if err != nil {
		return ent.User{}, err
	}

	return *user, err
}
