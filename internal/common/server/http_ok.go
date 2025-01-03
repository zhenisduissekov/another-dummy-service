package server

import (
	"encoding/json"
	"net/http"
	"time"
)

type ResponseOK struct {
	Message    string `json:"message"`    // Optional human-readable message
	HTTPStatus int    `json:"httpStatus"` // HTTP status code (e.g., 200, 201)
	Data       any    `json:"data"`       // The payload for the response
	Timestamp  string `json:"timestamp"`  // Time of the response (ISO 8601 format)
}

func RespondOK(data any, w http.ResponseWriter, r *http.Request) {
	resp := ResponseOK{
		Message:    "Request processed successfully.",
		HTTPStatus: http.StatusOK,
		Data:       data,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
