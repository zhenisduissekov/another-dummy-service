package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zhenisduissekov/another-dummy-service/internal/domain"
	"io"
)

func portHttpToDomain(port *Port) (*domain.Port, error) {
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

func readPorts(ctx context.Context, r io.Reader, portChan chan Port) error {
	decoder := json.NewDecoder(r)

	t, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("failed to read opening delimiter: %w", err)
	}

	if t != json.Delim('{') {
		return fmt.Errorf("expected string, got %v", t)
	}

	for decoder.More() {
		if ctx.Err != nil {
			return ctx.Err()
		}

		t, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("failed to read opening delimiter: %w", err)
		}

		portId, ok := t.(string)
		if !ok {
			return fmt.Errorf("expected string, got %v", t)
		}

		var port Port
		err = decoder.Decode(&port)
		if err != nil {
			return fmt.Errorf("failed to decode port: %w", err)
		}

		port.Id = portId
		portChan <- port
	}

	return nil
}
