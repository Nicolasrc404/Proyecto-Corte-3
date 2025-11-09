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

type TransmutationHandler struct {
	Repo      *repository.TransmutationRepository
	Log       func(status int, path string, start time.Time)
	HandleErr func(w http.ResponseWriter, statusCode int, path string, cause error)
}

func NewTransmutationHandler(
	repo *repository.TransmutationRepository,
	handleErr func(http.ResponseWriter, int, string, error),
	log func(int, string, time.Time),
) *TransmutationHandler {
	return &TransmutationHandler{Repo: repo, HandleErr: handleErr, Log: log}
}

func (h *TransmutationHandler) Create(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req api.TransmutationRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}
	if req.AlchemistID == 0 || req.MaterialID == 0 {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, errors.New("invalid IDs"))
		return
	}

	t := &models.Transmutation{
		AlchemistID: req.AlchemistID,
		MaterialID:  req.MaterialID,
		Formula:     req.Formula,
		Status:      "en_proceso",
	}
	t, err := h.Repo.Save(t)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}

	// Simula procesamiento asíncrono
	go func(t *models.Transmutation) {
		time.Sleep(5 * time.Second)
		t.Status = "completada"
		t.Result = "Éxito: transmutación estable."
		h.Repo.Save(t)
	}(t)

	resp := &api.TransmutationResponseDto{
		ID:          int(t.ID),
		AlchemistID: t.AlchemistID,
		MaterialID:  t.MaterialID,
		Status:      t.Status,
		Result:      t.Result,
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": resp})
	h.Log(http.StatusCreated, r.URL.Path, start)
}

func (h *TransmutationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	transmutations, err := h.Repo.FindAll()
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	resp := []*api.TransmutationResponseDto{}
	for _, t := range transmutations {
		resp = append(resp, &api.TransmutationResponseDto{
			ID:          int(t.ID),
			AlchemistID: t.AlchemistID,
			MaterialID:  t.MaterialID,
			Status:      t.Status,
			Result:      t.Result,
			CreatedAt:   t.CreatedAt.Format(time.RFC3339),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": resp})
	h.Log(http.StatusOK, r.URL.Path, start)
}

func (h *TransmutationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}
	t, err := h.Repo.FindById(id)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	if t == nil {
		h.HandleErr(w, http.StatusNotFound, r.URL.Path, errors.New("transmutation not found"))
		return
	}
	resp := &api.TransmutationResponseDto{
		ID:          int(t.ID),
		AlchemistID: t.AlchemistID,
		MaterialID:  t.MaterialID,
		Status:      t.Status,
		Result:      t.Result,
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": resp})
	h.Log(http.StatusOK, r.URL.Path, start)
}

func (h *TransmutationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}

	t, err := h.Repo.FindById(id)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	if t == nil {
		h.HandleErr(w, http.StatusNotFound, r.URL.Path, errors.New("transmutation not found"))
		return
	}

	if err := h.Repo.Delete(t); err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TransmutationHandler) Edit(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}

	t, err := h.Repo.FindById(id)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	if t == nil {
		h.HandleErr(w, http.StatusNotFound, r.URL.Path, errors.New("transmutation not found"))
		return
	}

	var req api.TransmutationEditRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}

	if req.Formula != nil {
		t.Formula = *req.Formula
	}
	if req.Status != nil {
		t.Status = *req.Status
	}
	if req.Result != nil {
		t.Result = *req.Result
	}

	t, err = h.Repo.Save(t)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}

	resp := &api.TransmutationResponseDto{
		ID:          int(t.ID),
		AlchemistID: t.AlchemistID,
		MaterialID:  t.MaterialID,
		Status:      t.Status,
		Result:      t.Result,
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]any{"data": resp})
	h.Log(http.StatusAccepted, r.URL.Path, start)
}
