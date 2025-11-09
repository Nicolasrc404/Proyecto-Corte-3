package server

import (
	"encoding/json"
	"net/http"
	"time"
)

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Path       string `json:"path"`
	Message    string `json:"message"`
	Timestamp  string `json:"timestamp"`
}

func (s *Server) HandleError(w http.ResponseWriter, statusCode int, path string, cause error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := ErrorResponse{
		StatusCode: statusCode,
		Path:       path,
		Message:    cause.Error(),
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(resp)
	s.logger.Error(statusCode, path, cause)
}
