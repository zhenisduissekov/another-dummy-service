package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPort(t *testing.T) {
	t.Parallel()

	portId := "id"
	portCode := "code"
	portName := "name"
	portCity := "city"
	portCountry := "country"

	t.Run("valid", func(t *testing.T) {
		port, err := NewPort(portId, portName, portCode, portCity, portCountry, nil, nil, nil, "", "", nil)
		require.NoError(t, err)

		require.Equal(t, portId, port.Id())
		require.Equal(t, portCode, port.Code())
		require.Equal(t, portName, port.Name())
		require.Equal(t, portCountry, port.Country())
	})

	t.Run("missing port id", func(t *testing.T) {
		_, err := NewPort("", portName, portCode, portCity, portCountry, nil, nil, nil, "", "", nil)
		require.Error(t, err)
	})

	t.Run("missing port name", func(t *testing.T) {
		_, err := NewPort(portId, "", portCode, portCity, portCountry, nil, nil, nil, "", "", nil)
		require.Error(t, err)
	})

	t.Run("missing port city", func(t *testing.T) {
		_, err := NewPort(portId, portName, portCode, "", portCountry, nil, nil, nil, "", "", nil)
		require.Error(t, err)
	})

	t.Run("missing port country", func(t *testing.T) {
		_, err := NewPort(portId, portName, portCode, portCity, "", nil, nil, nil, "", "", nil)
		require.Error(t, err)
	})
}
