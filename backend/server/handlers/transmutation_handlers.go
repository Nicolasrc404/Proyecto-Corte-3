package handlers

import (
	"backend-avanzada/api"
	"backend-avanzada/models"
	"backend-avanzada/repository"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type TransmutationHandler struct {
	Repo      *repository.TransmutationRepository
	TaskQueue *interface {
		StartTask(int, time.Duration, func(*models.Kill) error, *models.Kill)
	} // placeholder para tu TaskQueue
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

	// Simula procesamiento asíncrono (puedes adaptar a tu TaskQueue real)
	go func() {
		time.Sleep(5 * time.Second)
		t.Status = "completada"
		t.Result = "Éxito: transmutación estable."
		h.Repo.Save(t)
	}()

	resp := &api.TransmutationResponseDto{
		ID:          int(t.ID),
		AlchemistID: t.AlchemistID,
		MaterialID:  t.MaterialID,
		Status:      t.Status,
		Result:      t.Result,
		CreatedAt:   t.CreatedAt.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
	h.Log(http.StatusCreated, r.URL.Path, start)
}
