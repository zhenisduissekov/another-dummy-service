package inmem

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zhenisduissekov/another-dummy-service/internal/domain"
)

type PortStore struct {
	data map[string]*InmemPort
	mu   sync.RWMutex
}

func NewPortStore() *PortStore {
	return &PortStore{
		data: make(map[string]*InmemPort),
	}
}

func (ps *PortStore) GetPort(_ context.Context, id string) (*domain.Port, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	storePort, exists := ps.data[id]
	if !exists {
		return nil, domain.ErrNotFound
	}

	domainPort, err := portStoreToDomain(storePort)
	if err != nil {
		return nil, fmt.Errorf("portStoreToDomain failed: %w", err)
	}

	return domainPort, nil
}

func (ps *PortStore) CountPorts(_ context.Context) (int, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return len(ps.data), nil
}

func (ps *PortStore) CreateOrUpdatePort(ctx context.Context, port *domain.Port) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if port == nil {
		return domain.ErrNil
	}

	storePort := portDomainToStore(port)

	ps.mu.Lock()
	defer ps.mu.Unlock()

	_, exists := ps.data[storePort.Id]
	if exists {
		return ps.updatePort(ctx, storePort)
	} else {
		return ps.createPort(ctx, storePort)
	}
}

func (ps *PortStore) createPort(ctx context.Context, storePort *InmemPort) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if storePort == nil {
		return domain.ErrNil
	}

	storePort.CreatedAt = time.Now()
	storePort.UpdatedAt = time.Now()

	ps.data[storePort.Id] = storePort

	return nil
}

func (ps *PortStore) updatePort(ctx context.Context, port *InmemPort) error {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if port == nil {
		return domain.ErrNil
	}

	port, exist := ps.data[port.Id]
	if !exist {
		return domain.ErrNotFound
	}

	storePortCopy := port.Copy()

	storePortCopy.Name = port.Name
	storePortCopy.Code = port.Code
	storePortCopy.City = port.City
	storePortCopy.Country = port.Country
	storePortCopy.Alias = append([]string(nil), port.Alias...)
	storePortCopy.Regions = append([]string(nil), port.Regions...)
	storePortCopy.Coordinates = append([]float64(nil), port.Coordinates...)
	storePortCopy.Province = port.Province
	storePortCopy.Timezone = port.Timezone
	storePortCopy.Unlocs = append([]string(nil), port.Unlocs...)

	storePortCopy.UpdatedAt = time.Now()

	ps.data[port.Id] = storePortCopy

	return nil
}

func (ps *PortStore) DeletePortById(ctx context.Context, id string) error {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Ensure thread safety
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// Check if the port exists, and delete if found
	if _, exists := ps.data[id]; !exists {
		return domain.ErrNotFound
	}
	delete(ps.data, id)
	return nil
}

func (ps *PortStore) DeleteAllPorts(ctx context.Context) error {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Ensure thread safety
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// Reinitialize the map to clear all ports
	ps.data = make(map[string]*InmemPort)

	return nil
}
