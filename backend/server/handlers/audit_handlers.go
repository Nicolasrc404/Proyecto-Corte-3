package handlers

import (
	"backend-avanzada/api"
	"backend-avanzada/models"
	"backend-avanzada/repository"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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

// GET /audits
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
	json.NewEncoder(w).Encode(map[string]any{"data": resp})
	h.Log(http.StatusOK, r.URL.Path, start)
}

// GET /audits/{id}
func (h *AuditHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}

	a, err := h.Repo.FindById(id)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	if a == nil {
		h.HandleErr(w, http.StatusNotFound, r.URL.Path, errors.New("audit not found"))
		return
	}

	resp := &api.AuditResponseDto{
		ID:        int(a.ID),
		Action:    a.Action,
		Entity:    a.Entity,
		EntityID:  a.EntityID,
		UserEmail: a.UserEmail,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"data": resp})
	h.Log(http.StatusOK, r.URL.Path, start)
}

// POST /audits
func (h *AuditHandler) Create(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req api.AuditRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}

	a := &models.Audit{
		Action:    req.Action,
		Entity:    req.Entity,
		EntityID:  req.EntityID,
		UserEmail: req.UserEmail,
	}
	a, err := h.Repo.Save(a)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}

	resp := &api.AuditResponseDto{
		ID:        int(a.ID),
		Action:    a.Action,
		Entity:    a.Entity,
		EntityID:  a.EntityID,
		UserEmail: a.UserEmail,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"data": resp})
	h.Log(http.StatusCreated, r.URL.Path, start)
}

// PUT /audits/{id}
func (h *AuditHandler) Edit(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}

	a, err := h.Repo.FindById(id)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	if a == nil {
		h.HandleErr(w, http.StatusNotFound, r.URL.Path, errors.New("audit not found"))
		return
	}

	var req api.AuditEditRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}

	if req.Action != nil {
		a.Action = *req.Action
	}
	if req.Entity != nil {
		a.Entity = *req.Entity
	}
	if req.EntityID != nil {
		a.EntityID = *req.EntityID
	}
	if req.UserEmail != nil {
		a.UserEmail = *req.UserEmail
	}

	a, err = h.Repo.Save(a)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}

	resp := &api.AuditResponseDto{
		ID:        int(a.ID),
		Action:    a.Action,
		Entity:    a.Entity,
		EntityID:  a.EntityID,
		UserEmail: a.UserEmail,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]any{"data": resp})
	h.Log(http.StatusAccepted, r.URL.Path, start)
}

// DELETE /audits/{id}
func (h *AuditHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}

	a, err := h.Repo.FindById(id)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	if a == nil {
		h.HandleErr(w, http.StatusNotFound, r.URL.Path, errors.New("audit not found"))
		return
	}

	if err := h.Repo.Delete(a); err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
