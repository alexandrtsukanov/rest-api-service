package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"restapi-tasks/internal/database"
	"restapi-tasks/internal/models"
	"strconv"
	"strings"
)

type Handlers struct {
	service *database.TaskService
}

func (h *Handlers) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetAll()

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Getting tasks error")
		return
	}

	RespondWithJSON(w, http.StatusOK, tasks)
}

func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	id, idValidationErr := h.validateId(w, r)

	if idValidationErr != nil {
		return
	}

	task, err := h.service.GetById(id)

	if err != nil {
		if errors.Is(err, fmt.Errorf(database.NOT_FOUND_ERROR, id)) {
			RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondWithJSON(w, http.StatusOK, task)
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid task data")
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		RespondWithError(w, http.StatusBadRequest, "Task title cannot be empty")
		return
	}

	task, err := h.service.Create(input)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, task)
}

func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, idValidationErr := h.validateId(w, r)

	if idValidationErr != nil {
		return
	}

	var input models.UpdateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid task data")
		return
	}

	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		RespondWithError(w, http.StatusBadRequest, "Task title cannot be empty")
		return
	}

	task, err := h.service.Update(id, input)

	if err != nil {
		if errors.Is(err, fmt.Errorf(database.NOT_FOUND_ERROR, id)) {
			RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondWithJSON(w, http.StatusOK, task)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, idValidationErr := h.validateId(w, r)

	if idValidationErr != nil {
		return
	}

	removeErr := h.service.Delete(id)

	if removeErr != nil {
		if errors.Is(removeErr, fmt.Errorf(database.NOT_FOUND_ERROR, id)) {
			RespondWithError(w, http.StatusNotFound, removeErr.Error())
		} else {
			RespondWithError(w, http.StatusInternalServerError, removeErr.Error())
		}
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}

func (h *Handlers) validateId(w http.ResponseWriter, r *http.Request) (int, error) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idString := pathParts[0]

	id, err := strconv.Atoi(idString)

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return 0, err
	}

	return id, nil
}

func NewHandlers(service *database.TaskService) *Handlers {
	return &Handlers{service}
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(payload)
}

func RespondWithError(w http.ResponseWriter, statusCode int, messsage string) {
	RespondWithJSON(w, statusCode, map[string]string{"error": messsage})
}
