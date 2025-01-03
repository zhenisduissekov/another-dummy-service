package inmem

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/zhenisduissekov/another-dummy-service/internal/domain"
)

func TestPortStore_CreateOrUpdatePort(t *testing.T) {
	t.Parallel()

	store := NewPortStore()

	t.Run("create port", func(t *testing.T) {
		t.Parallel()

		randomPort := newRandomDomainPort(t)

		err := store.CreateOrUpdatePort(context.Background(), randomPort)
		require.NoError(t, err)

		port, err := store.GetPort(context.Background(), randomPort.Id())
		require.NoError(t, err)
		require.Equal(t, port, randomPort)
	})

	t.Run("update port", func(t *testing.T) {
		t.Parallel()

		randomPort := newRandomDomainPort(t)

		err := store.CreateOrUpdatePort(context.Background(), randomPort)
		require.NoError(t, err)

		beforeUpdatedPort, err := store.GetPort(context.Background(), randomPort.Id())
		require.NoError(t, err)
		require.Equal(t, beforeUpdatedPort, randomPort)

		err = beforeUpdatedPort.SetName("updated name")
		require.NoError(t, err)

		err = store.CreateOrUpdatePort(context.Background(), randomPort)
		require.NoError(t, err)

		updatedPort, err := store.GetPort(context.Background(), randomPort.Id())
		require.NoError(t, err)
		require.NotEqual(t, updatedPort, beforeUpdatedPort)
	})

	t.Run("delete port", func(t *testing.T) {
		t.Parallel()
		store := NewPortStore()

		randomPort := newRandomDomainPort(t)

		err := store.CreateOrUpdatePort(context.Background(), randomPort)
		require.NoError(t, err)

		port, err := store.GetPort(context.Background(), randomPort.Id())
		require.NoError(t, err)
		require.Equal(t, port, randomPort)

		err = store.DeletePortById(context.Background(), randomPort.Id())
		require.NoError(t, err)

		port, err = store.GetPort(context.Background(), randomPort.Id())
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("delete all ports", func(t *testing.T) {
		t.Parallel()

		store := NewPortStore()

		randomPort1 := newRandomDomainPort(t)
		randomPort2 := newRandomDomainPort(t)
		randomPort3 := newRandomDomainPort(t)

		err := store.CreateOrUpdatePort(context.Background(), randomPort1)
		require.NoError(t, err)
		err = store.CreateOrUpdatePort(context.Background(), randomPort2)
		require.NoError(t, err)
		err = store.CreateOrUpdatePort(context.Background(), randomPort3)
		require.NoError(t, err)

		port, err := store.GetPort(context.Background(), randomPort1.Id())
		require.NoError(t, err)
		require.Equal(t, port, randomPort1)
		port, err = store.GetPort(context.Background(), randomPort2.Id())
		require.NoError(t, err)
		require.Equal(t, port, randomPort2)
		port, err = store.GetPort(context.Background(), randomPort3.Id())
		require.NoError(t, err)
		require.Equal(t, port, randomPort3)

		count, err := store.CountPorts(context.Background())
		require.NoError(t, err)
		require.Equal(t, 3, count)

		err = store.DeleteAllPorts(context.Background())
		require.NoError(t, err)

		count, err = store.CountPorts(context.Background())
		require.NoError(t, err)
		require.Equal(t, 0, count)

		port, err = store.GetPort(context.Background(), randomPort1.Id())
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrNotFound)

		port, err = store.GetPort(context.Background(), randomPort2.Id())
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrNotFound)

		port, err = store.GetPort(context.Background(), randomPort3.Id())
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("nil port", func(t *testing.T) {
		t.Parallel()

		err := store.CreateOrUpdatePort(context.Background(), nil)
		require.ErrorIs(t, err, domain.ErrNil)
	})

}

func newRandomDomainPort(t *testing.T) *domain.Port {
	t.Helper()
	randomID := uuid.New().String()
	port, err := domain.NewPort(randomID, randomID, randomID, randomID, randomID, nil, nil, nil, randomID, randomID, nil)
	require.NoError(t, err)
	return port
}
