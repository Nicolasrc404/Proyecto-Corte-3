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
	Repo             *repository.MaterialRepository
	Dispatcher       AsyncDispatcher
	CurrentUser      func(*http.Request) string
	ReportAsyncError func(string, error)
	HandleErr        func(http.ResponseWriter, int, string, error)
	Log              func(int, string, time.Time)
}

func NewMaterialHandler(
	repo *repository.MaterialRepository,
	dispatcher AsyncDispatcher,
	currentUser func(*http.Request) string,
	reportAsyncError func(string, error),
	handleErr func(http.ResponseWriter, int, string, error),
	log func(int, string, time.Time),
) *MaterialHandler {
	return &MaterialHandler{
		Repo:             repo,
		Dispatcher:       dispatcher,
		CurrentUser:      currentUser,
		ReportAsyncError: reportAsyncError,
		HandleErr:        handleErr,
		Log:              log,
	}
}

func (h *MaterialHandler) userEmail(r *http.Request) string {
	if h.CurrentUser != nil {
		return h.CurrentUser(r)
	}
	return ""
}

func (h *MaterialHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	materials, err := h.Repo.FindAll()
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	resp := make([]*api.MaterialResponseDto, 0, len(materials))
	for _, m := range materials {
		resp = append(resp, &api.MaterialResponseDto{
			ID:        int(m.ID),
			Name:      m.Name,
			Category:  m.Category,
			Quantity:  m.Quantity,
			CreatedAt: m.CreatedAt.Format(time.RFC3339),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": resp})
	h.Log(http.StatusOK, r.URL.Path, start)
}

func (h *MaterialHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	id, err := strconv.Atoi(mux.Vars(r)["id"])
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
	resp := &api.MaterialResponseDto{
		ID:        int(m.ID),
		Name:      m.Name,
		Category:  m.Category,
		Quantity:  m.Quantity,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": resp})
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
	m := &models.Material{
		Name:     req.Name,
		Category: req.Category,
		Quantity: req.Quantity,
	}
	m, err := h.Repo.Save(m)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	if h.Dispatcher != nil {
		if err := h.Dispatcher.EnqueueAudit("create", "material", m.ID, h.userEmail(r), "Registro de material"); err != nil {
			h.ReportAsyncError(r.URL.Path, err)
		}
	}
	resp := &api.MaterialResponseDto{
		ID:        int(m.ID),
		Name:      m.Name,
		Category:  m.Category,
		Quantity:  m.Quantity,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": resp})
	h.Log(http.StatusCreated, r.URL.Path, start)
}

func (h *MaterialHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
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
	if h.Dispatcher != nil {
		if err := h.Dispatcher.EnqueueAudit("delete", "material", m.ID, h.userEmail(r), "Eliminación de material"); err != nil {
			h.ReportAsyncError(r.URL.Path, err)
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *MaterialHandler) Edit(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	id, err := strconv.Atoi(mux.Vars(r)["id"])
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

	var req api.MaterialEditRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}

	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Category != nil {
		m.Category = *req.Category
	}
	if req.Quantity != nil {
		m.Quantity = *req.Quantity
	}

	m, err = h.Repo.Save(m)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	if h.Dispatcher != nil {
		if err := h.Dispatcher.EnqueueAudit("update", "material", m.ID, h.userEmail(r), "Actualización de material"); err != nil {
			h.ReportAsyncError(r.URL.Path, err)
		}
	}

	resp := &api.MaterialResponseDto{
		ID:        int(m.ID),
		Name:      m.Name,
		Category:  m.Category,
		Quantity:  m.Quantity,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": resp})
	h.Log(http.StatusAccepted, r.URL.Path, start)
}
