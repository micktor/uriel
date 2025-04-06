package service

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/your_org/uriel/internal/config"
	"github.com/your_org/uriel/internal/repository"
)

type Service struct {
	config     *config.Config
	repository *repository.Repository
	tokenAuth  *jwtauth.JWTAuth
}
type ctxKey string

const JWTAuthCtxKey ctxKey = "uriel.jwt"

type JWTAuth struct {
	Token jwt.Token
}

func NewService(config *config.Config, repo *repository.Repository) *Service {
	jwtAuth := jwtauth.New(jwa.HS256.String(), []byte(config.JWTSecret), nil)
	return &Service{
		config:     config,
		repository: repo,
		tokenAuth:  jwtAuth,
	}
}

func (s *Service) GetTokenAuth() *jwtauth.JWTAuth {
	return s.tokenAuth
}
