package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"rri/task-back/internal/api/dto"
	"rri/task-back/internal/repository"
	"rri/task-back/internal/service"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

// getUserIDFromContext извлекает ID пользователя из контекста (установлен middleware аутентификации)
func getUserIDFromContext(r *http.Request) string {
	return "faba5b90-0237-41f2-84d5-db5d83096412"
	//return r.Context().Value("userID").(string)
}

// respondWithError отправляет JSON с ошибкой и соответствующим статусом
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(dto.ErrorResponse{Error: message, Code: code})
}

// CreateTask обрабатывает POST /api/tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	task, err := h.taskService.CreateTask(r.Context(), userID, &req)
	if err != nil {
		// Обработка ошибок валидации и доступа
		var validationErr interface{ Errors() map[string]string } // пример для validator
		if errors.As(err, &validationErr) {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else if errors.Is(err, repository.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "Resource not found")
		} else if err.Error() == "access denied: editor or owner role required" {
			respondWithError(w, http.StatusForbidden, "Insufficient permissions")
		} else {
			// Логируем внутреннюю ошибку
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// GetTask обрабатывает GET /api/tasks/{id}
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing task id")
		return
	}

	task, err := h.taskService.GetTask(r.Context(), taskID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "Task not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// UpdateTask обрабатывает PUT /api/tasks/{id}
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing task id")
		return
	}

	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	task, err := h.taskService.UpdateTask(r.Context(), taskID, userID, &req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "Task not found")
		} else if err.Error() == "access denied: editor or owner role required" {
			respondWithError(w, http.StatusForbidden, "Insufficient permissions")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// DeleteTask обрабатывает DELETE /api/tasks/{id}
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing task id")
		return
	}

	err := h.taskService.DeleteTask(r.Context(), taskID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "Task not found")
		} else if err.Error() == "access denied: editor or owner role required" {
			respondWithError(w, http.StatusForbidden, "Insufficient permissions")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// MoveTask обрабатывает POST /api/tasks/{id}/move
func (h *TaskHandler) MoveTask(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing task id")
		return
	}

	var req dto.MoveTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	task, err := h.taskService.MoveTask(r.Context(), taskID, userID, &req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "Task not found")
		} else if err.Error() == "access denied: editor or owner role required" {
			respondWithError(w, http.StatusForbidden, "Insufficient permissions")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
