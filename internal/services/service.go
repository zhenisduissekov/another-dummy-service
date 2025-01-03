package services

import (
	"context"

	"github.com/zhenisduissekov/another-dummy-service/internal/domain"
)

type PortRepository interface {
	CreateOrUpdatePort(ctx context.Context, port *domain.Port) error
	CountPorts(ctx context.Context) (int, error)
	GetPort(ctx context.Context, id string) (*domain.Port, error)
	DeleteAllPorts(ctx context.Context) error
	DeletePortById(ctx context.Context, id string) error
}

type PortService struct {
	repo PortRepository
}

func NewPortService(repo PortRepository) PortService {
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

func (ps PortService) DeletePortById(ctx context.Context, id string) error {
	return ps.repo.DeletePortById(ctx, id)
}

func (ps PortService) DeleteAllPorts(ctx context.Context) error {
	return ps.repo.DeleteAllPorts(ctx)
}
