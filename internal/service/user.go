package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/your_org/uriel/internal/dto"
	"github.com/your_org/uriel/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

func (s *Service) OAuthLoginUser(ctx context.Context, request dto.RegisterRequest) (tokenString string, err error) {
	userExists, err := s.repository.UserByEmailExists(ctx, request.Email)
	if err != nil {
		return "", err
	}

	if !userExists {
		_, err = s.repository.CreateUser(ctx, repository.CreateUserInput{
			Email:        request.Email,
			AuthProvider: request.AuthProvider,
			OAuthID:      request.OAuthID,
		})
		if err != nil {
			return "", err
		}
	}

	user, err := s.repository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return "", err
	}

	_, tokenString, err = s.tokenAuth.Encode(map[string]interface{}{
		"sub": user.ID,
		"iss": "uriel",
		"exp": time.Now().Add(time.Hour * 24).UTC().Unix(),
		"iat": time.Now().UTC().Unix(),
	})
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) LoginUser(ctx context.Context, request dto.LoginRequest) (tokenString string, err error) {
	user, err := s.repository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return "", err
	}
	err = confirmHash(request.Password, user.Password)
	if err != nil {
		return "", err
	}

	_, tokenString, err = s.tokenAuth.Encode(map[string]interface{}{
		"sub": user.ID,
		"iss": "uriel",
		"exp": time.Now().Add(time.Hour * 24).UTC().Unix(),
		"iat": time.Now().UTC().Unix(),
	})
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type CreateUserRequest struct {
	Email    string
	Password string
}

// CreateUser creates a new user with the provided email and password, returning the user ID or an error.
func (s *Service) CreateUser(ctx context.Context, request CreateUserRequest) (string, error) {
	passwordHash, err := hash(request.Password)
	if err != nil {
		return "", err
	}

	userExists, err := s.repository.UserByEmailExists(ctx, request.Email)
	if err != nil {
		return "", err
	} else if userExists {
		return "", errors.New("user with this email already exists")
	}

	resp, err := s.repository.CreateUser(ctx, repository.CreateUserInput{
		Email:    request.Email,
		Password: passwordHash,
	})
	if err != nil {
		return "", err
	}
	return resp.ID, err
}

// hash hashes a password
func hash(password string) (string, error) {
	password = strings.TrimSpace(password)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}
	return string(hash), nil
}

// confirmHash verifies that a password matches the given hash
func confirmHash(password string, hash string) error {
	password = strings.TrimSpace(password)
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return err
	}
	return nil
}
