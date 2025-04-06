package httpd

import (
	"context"
	"github.com/christfirst/uriel/internal/dto"
)

func (h *Handler) GetViewer(ctx context.Context) (dto.ViewerResponse, error) {
	viewer, err := h.service.GetViewerData(ctx, getJWT(ctx).Token.Subject())
	if err != nil {
		return dto.ViewerResponse{}, err
	}

	return viewer, err
}
