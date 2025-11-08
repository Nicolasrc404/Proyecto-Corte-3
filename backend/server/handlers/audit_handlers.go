package handlers

import (
	"backend-avanzada/api"
	"backend-avanzada/repository"
	"encoding/json"
	"net/http"
	"time"
)

type AuditHandler struct {
	Repo      *repository.AuditRepository
	Log       func(status int, path string, start time.Time)
	HandleErr func(w http.ResponseWriter, statusCode int, path string, cause error)
}

func NewAuditHandler(repo *repository.AuditRepository,
	handleErr func(http.ResponseWriter, int, string, error),
	log func(int, string, time.Time)) *AuditHandler {
	return &AuditHandler{Repo: repo, Log: log, HandleErr: handleErr}
}

func (h *AuditHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	audits, err := h.Repo.FindAll()
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	resp := []*api.AuditResponseDto{}
	for _, a := range audits {
		resp = append(resp, &api.AuditResponseDto{
			ID:        int(a.ID),
			Entity:    a.Entity,
			EntityID:  a.EntityID,
			Action:    a.Action,
			Timestamp: a.Timestamp,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
	h.Log(http.StatusOK, r.URL.Path, start)
}
