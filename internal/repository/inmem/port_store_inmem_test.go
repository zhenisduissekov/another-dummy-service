package inmem

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/zhenisduissekov/another-dummy-service/internal/domain"
	"testing"
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
