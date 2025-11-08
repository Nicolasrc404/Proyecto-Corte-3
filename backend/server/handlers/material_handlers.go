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

type MaterialHandler struct {
	Repo      *repository.MaterialRepository
	Log       func(status int, path string, start time.Time)
	HandleErr func(w http.ResponseWriter, statusCode int, path string, cause error)
}

func NewMaterialHandler(repo *repository.MaterialRepository,
	handleErr func(http.ResponseWriter, int, string, error),
	log func(int, string, time.Time)) *MaterialHandler {
	return &MaterialHandler{
		Repo:      repo,
		Log:       log,
		HandleErr: handleErr,
	}
}

func (h *MaterialHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	materials, err := h.Repo.FindAll()
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	resp := []*api.MaterialResponseDto{}
	for _, m := range materials {
		resp = append(resp, &api.MaterialResponseDto{
			ID:        int(m.ID),
			Name:      m.Name,
			Category:  m.Category,
			Quantity:  m.Quantity,
			CreatedAt: m.CreatedAt.String(),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
	h.Log(http.StatusOK, r.URL.Path, start)
}

func (h *MaterialHandler) Create(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req api.MaterialRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}
	if req.Name == "" {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, errors.New("name required"))
		return
	}
	m := &models.Material{Name: req.Name, Category: req.Category, Quantity: req.Quantity}
	m, err := h.Repo.Save(m)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	resp := &api.MaterialResponseDto{
		ID:        int(m.ID),
		Name:      m.Name,
		Category:  m.Category,
		Quantity:  m.Quantity,
		CreatedAt: m.CreatedAt.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
	h.Log(http.StatusCreated, r.URL.Path, start)
}

func (h *MaterialHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}
	m, err := h.Repo.FindById(id)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	if m == nil {
		h.HandleErr(w, http.StatusNotFound, r.URL.Path, errors.New("material not found"))
		return
	}
	if err := h.Repo.Delete(m); err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
