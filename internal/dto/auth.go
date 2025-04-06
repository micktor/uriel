package dto

import "github.com/lestrrat-go/jwx/v2/jwt"

type RegisterRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	AuthProvider string `json:"-"`
	OAuthID      string `json:"-"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JWTAuth struct {
	Token  jwt.Token
	Claims map[string]interface{}
}
