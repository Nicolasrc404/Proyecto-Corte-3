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

type AlchemistHandler struct {
	Repo      *repository.AlchemistRepository
	Log       func(status int, path string, start time.Time)
	HandleErr func(w http.ResponseWriter, statusCode int, path string, cause error)
}

func NewAlchemistHandler(repo *repository.AlchemistRepository,
	handleErr func(http.ResponseWriter, int, string, error),
	log func(int, string, time.Time)) *AlchemistHandler {
	return &AlchemistHandler{Repo: repo, HandleErr: handleErr, Log: log}
}

func (h *AlchemistHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	alchs, err := h.Repo.FindAll()
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	resp := make([]*api.AlchemistResponseDto, 0, len(alchs))
	for _, a := range alchs {
		resp = append(resp, &api.AlchemistResponseDto{
			ID:        int(a.ID),
			Name:      a.Name,
			Age:       a.Age,
			Specialty: a.Specialty,
			Rank:      a.Rank,
			CreatedAt: a.CreatedAt.Format(time.RFC3339),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": resp})
	h.Log(http.StatusOK, r.URL.Path, start)
}

func (h *AlchemistHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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
		h.HandleErr(w, http.StatusNotFound, r.URL.Path, errors.New("alchemist not found"))
		return
	}
	resp := &api.AlchemistResponseDto{
		ID:        int(a.ID),
		Name:      a.Name,
		Age:       a.Age,
		Specialty: a.Specialty,
		Rank:      a.Rank,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": resp})
	h.Log(http.StatusOK, r.URL.Path, start)
}

func (h *AlchemistHandler) Create(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req api.AlchemistRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}
	if req.Name == "" {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, errors.New("name required"))
		return
	}
	a := &models.Alchemist{
		Name:      req.Name,
		Age:       int(req.Age),
		Specialty: req.Specialty,
		Rank:      req.Rank,
	}
	a, err := h.Repo.Save(a)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	resp := &api.AlchemistResponseDto{
		ID:        int(a.ID),
		Name:      a.Name,
		Age:       a.Age,
		Specialty: a.Specialty,
		Rank:      a.Rank,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": resp})
	h.Log(http.StatusCreated, r.URL.Path, start)
}

func (h *AlchemistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	a, err := h.Repo.FindById(id)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	if a == nil {
		h.HandleErr(w, http.StatusNotFound, r.URL.Path, errors.New("alchemist not found"))
		return
	}
	if err := h.Repo.Delete(a); err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AlchemistHandler) Edit(w http.ResponseWriter, r *http.Request) {
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
		h.HandleErr(w, http.StatusNotFound, r.URL.Path, errors.New("alchemist not found"))
		return
	}

	var req api.AlchemistEditRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}

	if req.Name != nil {
		a.Name = *req.Name
	}
	if req.Age != nil {
		a.Age = int(*req.Age)
	}
	if req.Specialty != nil {
		a.Specialty = *req.Specialty
	}
	if req.Rank != nil {
		a.Rank = *req.Rank
	}

	if _, err := h.Repo.Save(a); err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}

	resp := &api.AlchemistResponseDto{
		ID:        int(a.ID),
		Name:      a.Name,
		Age:       a.Age,
		Specialty: a.Specialty,
		Rank:      a.Rank,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": resp})
	h.Log(http.StatusAccepted, r.URL.Path, start)
}
