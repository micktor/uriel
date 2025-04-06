package service

import (
	"context"
	"github.com/christfirst/uriel/internal/dto"
)

func (s *Service) GetViewerData(ctx context.Context, userID string) (dto.ViewerResponse, error) {
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		return dto.ViewerResponse{}, err
	}
	return dto.ViewerResponse{
		UserID: user.ID,
		Email:  user.Email,
	}, nil
}
