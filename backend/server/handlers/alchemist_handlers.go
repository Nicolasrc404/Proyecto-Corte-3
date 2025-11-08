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
	Repo       *repository.AlchemistRepository
	PeopleRepo repository.Repository[models.Person]
	Log        func(status int, path string, start time.Time)
	HandleErr  func(w http.ResponseWriter, statusCode int, path string, cause error)
}

func NewAlchemistHandler(repo *repository.AlchemistRepository,
	peopleRepo repository.Repository[models.Person],
	handleErr func(http.ResponseWriter, int, string, error),
	log func(int, string, time.Time)) *AlchemistHandler {

	return &AlchemistHandler{
		Repo:       repo,
		PeopleRepo: peopleRepo,
		HandleErr:  handleErr,
		Log:        log,
	}
}

func (h *AlchemistHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	alchs, err := h.Repo.FindAll()
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	resp := []*api.AlchemistResponseDto{}
	for _, a := range alchs {
		resp = append(resp, &api.AlchemistResponseDto{
			ID:        int(a.ID),
			Name:      a.Person.Name,
			Age:       a.Person.Age,
			Specialty: a.Specialty,
			Rank:      a.Rank,
			CreatedAt: a.CreatedAt.String(),
		})
	}
	json.NewEncoder(w).Encode(resp)
	h.Log(http.StatusOK, r.URL.Path, start)
}

func (h *AlchemistHandler) Create(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var req api.AlchemistRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}
	person := &models.Person{Name: req.Name, Age: int(req.Age)}
	person, err := h.PeopleRepo.Save(person)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	a := &models.Alchemist{
		PersonID:  person.ID,
		Specialty: req.Specialty,
		Rank:      req.Rank,
	}
	a, err = h.Repo.Save(a)
	if err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	resp := &api.AlchemistResponseDto{
		ID:        int(a.ID),
		Name:      person.Name,
		Age:       person.Age,
		Specialty: a.Specialty,
		Rank:      a.Rank,
		CreatedAt: a.CreatedAt.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
	h.Log(http.StatusCreated, r.URL.Path, start)
}

func (h *AlchemistHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}

	a, err := h.Repo.FindById(int(id))
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
		Name:      a.Person.Name,
		Age:       a.Person.Age,
		Specialty: a.Specialty,
		Rank:      a.Rank,
		CreatedAt: a.CreatedAt.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
	h.Log(http.StatusOK, r.URL.Path, start)
}

func (h *AlchemistHandler) Edit(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}

	a, err := h.Repo.FindById(int(id))
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
		a.Person.Name = *req.Name
	}
	if req.Age != nil {
		a.Person.Age = int(*req.Age)
	}
	if req.Specialty != nil {
		a.Specialty = *req.Specialty
	}
	if req.Rank != nil {
		a.Rank = *req.Rank
	}

	if _, err := h.PeopleRepo.Save(&a.Person); err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	if _, err := h.Repo.Save(a); err != nil {
		h.HandleErr(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}

	resp := &api.AlchemistResponseDto{
		ID:        int(a.ID),
		Name:      a.Person.Name,
		Age:       a.Person.Age,
		Specialty: a.Specialty,
		Rank:      a.Rank,
		CreatedAt: a.CreatedAt.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
	h.Log(http.StatusAccepted, r.URL.Path, start)
}

func (h *AlchemistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		h.HandleErr(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}

	a, err := h.Repo.FindById(int(id))
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
