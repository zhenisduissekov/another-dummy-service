package services

import (
	"context"
	"github.com/zhenisduissekov/another-dummy-service/internal/domain"
)

type PortRepository interface {
	CreateOrUpdatePort(ctx context.Context, port *domain.Port) error
	CountPorts(ctx context.Context) (int, error)
	GetPort(ctx context.Context, id string) (*domain.Port, error)
}

type PortService struct {
	repo PortRepository
}

func NewService(repo PortRepository) PortService {
	return PortService{
		repo: repo,
	}
}

func (ps PortService) GetPort(ctx context.Context, id string) (*domain.Port, error) {
	return ps.repo.GetPort(ctx, id)
}

func (ps PortService) CountPorts(ctx context.Context) (int, error) {
	return ps.repo.CountPorts(ctx)
}

func (ps PortService) CreateOrUpdatePort(ctx context.Context, port *domain.Port) error {
	return ps.repo.CreateOrUpdatePort(ctx, port)
}
