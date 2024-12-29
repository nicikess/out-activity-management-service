package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/nicikess/out-run-management-service/internal/domain"
	"github.com/nicikess/out-run-management-service/internal/service/run"
	"github.com/nicikess/out-run-management-service/pkg/generated"
)

type RunHandler struct {
	service *run.Service
	logger  *zap.Logger
}

func NewRunHandler(service *run.Service, logger *zap.Logger) *RunHandler {
	return &RunHandler{
		service: service,
		logger:  logger,
	}
}

// getUserIDFromContext gets the user ID from the context
// In a real application, this would come from your auth middleware
func (h *RunHandler) getUserIDFromContext(r *http.Request) (uuid.UUID, error) {
	// TODO: Implement proper auth middleware
	// This is a placeholder - you should get the user ID from your auth system
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	return uuid.Parse(userID)
}

// StartRun handles the creation of a new run
func (h *RunHandler) StartRun(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		h.respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req generated.StartRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	coord := domain.Coordinate{
		Latitude:  float64(req.InitialLocation.Latitude),
		Longitude: float64(req.InitialLocation.Longitude),
		Timestamp: time.Now(),
	}

	run, err := h.service.StartRun(ctx, userID, coord)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrActiveRunExists):
			h.respondWithError(w, http.StatusConflict, "active run already exists")
		default:
			h.respondWithError(w, http.StatusInternalServerError, "failed to start run")
		}
		return
	}

	h.respondWithJSON(w, http.StatusCreated, run)
}

// GetRun handles retrieving a specific run
func (h *RunHandler) GetRun(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		h.respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	runID, err := uuid.Parse(chi.URLParam(r, "runId"))
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid run ID")
		return
	}

	run, err := h.service.GetRun(ctx, runID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRunNotFound):
			h.respondWithError(w, http.StatusNotFound, "run not found")
		case errors.Is(err, domain.ErrUnauthorized):
			h.respondWithError(w, http.StatusUnauthorized, "unauthorized")
		default:
			h.respondWithError(w, http.StatusInternalServerError, "failed to get run")
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, run)
}

// GetActiveRun handles retrieving the user's active run
func (h *RunHandler) GetActiveRun(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		h.respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	run, err := h.service.GetActiveRun(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRunNotFound):
			h.respondWithError(w, http.StatusNotFound, "no active run found")
		default:
			h.respondWithError(w, http.StatusInternalServerError, "failed to get active run")
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, run)
}

// PauseRun handles pausing an active run
func (h *RunHandler) PauseRun(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		h.respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	runID, err := uuid.Parse(chi.URLParam(r, "runId"))
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid run ID")
		return
	}

	run, err := h.service.PauseRun(ctx, runID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRunNotFound):
			h.respondWithError(w, http.StatusNotFound, "run not found")
		case errors.Is(err, domain.ErrUnauthorized):
			h.respondWithError(w, http.StatusUnauthorized, "unauthorized")
		case errors.Is(err, domain.ErrInvalidRunStatus):
			h.respondWithError(w, http.StatusConflict, "run cannot be paused")
		default:
			h.respondWithError(w, http.StatusInternalServerError, "failed to pause run")
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, run)
}

// ResumeRun handles resuming a paused run
func (h *RunHandler) ResumeRun(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		h.respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	runID, err := uuid.Parse(chi.URLParam(r, "runId"))
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid run ID")
		return
	}

	run, err := h.service.ResumeRun(ctx, runID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRunNotFound):
			h.respondWithError(w, http.StatusNotFound, "run not found")
		case errors.Is(err, domain.ErrUnauthorized):
			h.respondWithError(w, http.StatusUnauthorized, "unauthorized")
		case errors.Is(err, domain.ErrInvalidRunStatus):
			h.respondWithError(w, http.StatusConflict, "run cannot be resumed")
		default:
			h.respondWithError(w, http.StatusInternalServerError, "failed to resume run")
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, run)
}

// EndRun handles ending a run
func (h *RunHandler) EndRun(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		h.respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	runID, err := uuid.Parse(chi.URLParam(r, "runId"))
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid run ID")
		return
	}

	run, err := h.service.EndRun(ctx, runID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRunNotFound):
			h.respondWithError(w, http.StatusNotFound, "run not found")
		case errors.Is(err, domain.ErrUnauthorized):
			h.respondWithError(w, http.StatusUnauthorized, "unauthorized")
		case errors.Is(err, domain.ErrInvalidRunStatus):
			h.respondWithError(w, http.StatusConflict, "run cannot be ended")
		default:
			h.respondWithError(w, http.StatusInternalServerError, "failed to end run")
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, run)
}

// Helper methods for response handling

func (h *RunHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *RunHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		h.logger.Error("Failed to marshal response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
