package inmem

import (
	"fmt"
	"time"

	"github.com/zhenisduissekov/another-dummy-service/internal/domain"
)

type InmemPort struct {
	Id          string
	Name        string
	Code        string
	City        string
	Country     string
	Alias       []string
	Regions     []string
	Coordinates []float64
	Province    string
	Timezone    string
	Unlocs      []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *InmemPort) Copy() *InmemPort {
	if p == nil {
		return nil
	}

	return &InmemPort{
		Id:          p.Id,
		Name:        p.Name,
		Code:        p.Code,
		City:        p.City,
		Country:     p.Country,
		Alias:       append([]string(nil), p.Alias...),
		Regions:     append([]string(nil), p.Regions...),
		Coordinates: append([]float64(nil), p.Coordinates...),
		Province:    p.Province,
		Timezone:    p.Timezone,
		Unlocs:      p.Unlocs,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func portStoreToDomain(port *InmemPort) (*domain.Port, error) {
	if port == nil {
		return nil, fmt.Errorf("store port is nil")
	}

	return domain.NewPort(
		port.Id,
		port.Name,
		port.Code,
		port.City,
		port.Country,
		append([]string(nil), port.Alias...),
		append([]string(nil), port.Regions...),
		append([]float64(nil), port.Coordinates...),
		port.Province,
		port.Timezone,
		append([]string(nil), port.Unlocs...),
	)
}

func portDomainToStore(p *domain.Port) *InmemPort {
	return &InmemPort{
		Id:          p.Id(),
		Name:        p.Name(),
		Code:        p.Code(),
		City:        p.City(),
		Country:     p.Country(),
		Alias:       append([]string(nil), p.Alias()...),
		Regions:     append([]string(nil), p.Regions()...),
		Coordinates: append([]float64(nil), p.Coordinates()...),
		Province:    p.Province(),
		Timezone:    p.Timezone(),
		Unlocs:      append([]string(nil), p.Unlocs()...),
	}
}
