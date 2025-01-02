package server

import (
	"encoding/json"
	"net/http"
)

type ResponseOk struct {
	Msg  string
	Code uint
	Data any
}

func RespondOK(data any, w http.ResponseWriter, _ *http.Request) {
	resp := ResponseOk{
		Msg:  "success",
		Code: 0000,
		Data: data,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
