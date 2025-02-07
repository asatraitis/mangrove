package handler

import (
	"encoding/json"
	"net/http"

	"github.com/asatraitis/mangrove/internal/dto"
)

func sendErrResponse[T any](w http.ResponseWriter, err *dto.ResponseError, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(
		dto.Response[T]{
			Response: nil,
			Error:    err,
		},
	)
}

func getReqIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
