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
	HandleErr func(http.ResponseWriter, int, string, error)
	Log       func(int, string, time.Time)
}

func NewAuditHandler(
	repo *repository.AuditRepository,
	handleErr func(http.ResponseWriter, int, string, error),
	log func(int, string, time.Time),
) *AuditHandler {
	return &AuditHandler{Repo: repo, HandleErr: handleErr, Log: log}
}

func (h *AuditHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	audits, err := h.Repo.FindAll()
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}

	resp := make([]*api.AuditResponseDto, 0, len(audits))
	for _, a := range audits {
		resp = append(resp, &api.AuditResponseDto{
			ID:        int(a.ID),
			Action:    a.Action,
			Entity:    a.Entity,
			EntityID:  a.EntityID,
			UserEmail: a.UserEmail,
			CreatedAt: a.CreatedAt.Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": resp})
	h.Log(http.StatusOK, r.URL.Path, start)
}
