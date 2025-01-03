package transport

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/zhenisduissekov/another-dummy-service/internal/common/server"
	"github.com/zhenisduissekov/another-dummy-service/internal/domain"
)

type PortService interface {
	GetPort(ctx context.Context, id string) (*domain.Port, error)
	CountPorts(ctx context.Context) (int, error)
	CreateOrUpdatePort(ctx context.Context, port *domain.Port) error
	DeleteAllPorts(ctx context.Context) error
	DeletePortById(ctx context.Context, id string) error
}

type HttpServer struct {
	service PortService
}

func NewHttpServer(service PortService) HttpServer {
	return HttpServer{
		service: service,
	}
}

func (h HttpServer) CountPorts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	count, err := h.service.CountPorts(ctx)
	if err != nil {
		server.RespondWithError(err, w, r)
		return
	}

	server.RespondOK(map[string]int{"count": count}, w, r)
}

func (h HttpServer) GetPort(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("id")

	port, err := h.service.GetPort(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.NotFound("port-not-found", err, w, r)
			return
		}
		server.RespondWithError(err, w, r)
		return
	}

	response := Port{
		Id:          port.Id(),
		Name:        port.Name(),
		Code:        port.Code(),
		City:        port.City(),
		Country:     port.Country(),
		Alias:       port.Alias(),
		Regions:     port.Regions(),
		Coordinates: port.Coordinates(),
		Province:    port.Province(),
		Unlocs:      port.Unlocs(),
		Timezone:    port.Timezone(),
	}

	server.RespondOK(response, w, r)
}

func (h HttpServer) UploadPorts(w http.ResponseWriter, r *http.Request) {

	portChan := make(chan Port)
	doneChan := make(chan struct{})
	errChan := make(chan error)

	go func() {
		err := readPorts(r.Context(), r.Body, portChan)
		if err != nil {
			errChan <- err
		} else {
			doneChan <- struct{}{}
		}
	}()
	portCounter := 0
	for {
		select {
		case <-r.Context().Done():
			log.Println("request context cancelled")
			return
		case <-doneChan:
			log.Println("finished reading ports")
			server.RespondOK(map[string]int{"total_ports": portCounter}, w, r)
			return
		case err := <-errChan:
			log.Printf("error while parsing port json: %+v", err)
			server.BadRequest("invalid json", err, w, r)
			return
		case port := <-portChan:
			portCounter++
			log.Printf("[%d] received port: %+v", portCounter, port)
			p, err := portHttpToDomain(&port)
			if err != nil {
				server.BadRequest("port-to-domain", err, w, r)
				return
			}

			err = h.service.CreateOrUpdatePort(r.Context(), p)
			if err != nil {
				server.RespondWithError(err, w, r)
				return
			}
		}
	}
}

func (h HttpServer) DeleteAllPorts(w http.ResponseWriter, r *http.Request) {
	deleteAll := r.URL.Query().Get("all") == "true"
	if !deleteAll {
		server.BadRequest("missing required parameter: all=true", nil, w, r)
		return
	}

	err := h.service.DeleteAllPorts(r.Context())
	if err != nil {
		server.InternalError("could not delete all ports", err, w, r)
		return
	}

	server.RespondOK("all ports deleted successfully", w, r)
}

func (h HttpServer) DeletePortsById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		server.BadRequest("missing port ID", nil, w, r)
		return
	}

	err := h.service.DeletePortById(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.NotFound("port not found", err, w, r)
			return
		}
		server.InternalError("could not delete port by id", err, w, r)
		return
	}

	server.RespondOK(fmt.Sprintf("deleted port[%s] successfully", id), w, r)
}
