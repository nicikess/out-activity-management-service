package run

import (
	"context"

	"github.com/google/uuid"
	"github.com/nicikess/out-run-management-service/internal/domain"
	"github.com/nicikess/out-run-management-service/internal/ports/repository"
)

type Service struct {
	repo repository.RunRepository
}

type serviceImpl struct {
	repo repository.RunRepository
}

func NewService(repo repository.RunRepository) *Service {
	return &Service{
		repo: repo,
	}
}

// StartRun creates a new run for the user
func (s *Service) StartRun(ctx context.Context, userID uuid.UUID, initialCoordinate domain.Coordinate) (*domain.Run, error) {
	// Check if user already has an active run
	if _, err := s.repo.GetActiveByUserID(ctx, userID); err == nil {
		return nil, domain.ErrActiveRunExists
	} else if err != domain.ErrRunNotFound {
		return nil, err
	}

	// Create new run
	run := domain.NewRun(userID, initialCoordinate)
	if err := s.repo.Create(ctx, run); err != nil {
		return nil, err
	}

	return run, nil
}

// PauseRun pauses an active run
func (s *Service) PauseRun(ctx context.Context, runID, userID uuid.UUID) (*domain.Run, error) {
	run, err := s.repo.GetByID(ctx, runID)
	if err != nil {
		return nil, err
	}

	// Verify user ownership
	if !run.IsUserAuthorized(userID) {
		return nil, domain.ErrUnauthorized
	}

	// Pause the run
	if err := run.Pause(); err != nil {
		return nil, err
	}

	// Update in repository
	if err := s.repo.Update(ctx, run); err != nil {
		return nil, err
	}

	return run, nil
}

// ResumeRun resumes a paused run
func (s *Service) ResumeRun(ctx context.Context, runID, userID uuid.UUID) (*domain.Run, error) {
	run, err := s.repo.GetByID(ctx, runID)
	if err != nil {
		return nil, err
	}

	// Verify user ownership
	if !run.IsUserAuthorized(userID) {
		return nil, domain.ErrUnauthorized
	}

	// Resume the run
	if err := run.Resume(); err != nil {
		return nil, err
	}

	// Update in repository
	if err := s.repo.Update(ctx, run); err != nil {
		return nil, err
	}

	return run, nil
}

// EndRun ends an active or paused run
func (s *Service) EndRun(ctx context.Context, runID, userID uuid.UUID) (*domain.Run, error) {
	run, err := s.repo.GetByID(ctx, runID)
	if err != nil {
		return nil, err
	}

	// Verify user ownership
	if !run.IsUserAuthorized(userID) {
		return nil, domain.ErrUnauthorized
	}

	// End the run
	if err := run.End(); err != nil {
		return nil, err
	}

	// Update in repository
	if err := s.repo.Update(ctx, run); err != nil {
		return nil, err
	}

	return run, nil
}

// AddCoordinate adds a new coordinate to an active run
func (s *Service) AddCoordinate(ctx context.Context, runID, userID uuid.UUID, coord domain.Coordinate) (*domain.Run, error) {
	run, err := s.repo.GetByID(ctx, runID)
	if err != nil {
		return nil, err
	}

	// Verify user ownership
	if !run.IsUserAuthorized(userID) {
		return nil, domain.ErrUnauthorized
	}

	// Add coordinate
	if err := run.AddCoordinate(coord); err != nil {
		return nil, err
	}

	// Update in repository
	if err := s.repo.Update(ctx, run); err != nil {
		return nil, err
	}

	return run, nil
}

// GetRun retrieves a run by ID
func (s *Service) GetRun(ctx context.Context, runID, userID uuid.UUID) (*domain.Run, error) {
	run, err := s.repo.GetByID(ctx, runID)
	if err != nil {
		return nil, err
	}

	// Verify user ownership
	if !run.IsUserAuthorized(userID) {
		return nil, domain.ErrUnauthorized
	}

	return run, nil
}

// GetActiveRun gets the user's active run if any
func (s *Service) GetActiveRun(ctx context.Context, userID uuid.UUID) (*domain.Run, error) {
	return s.repo.GetActiveByUserID(ctx, userID)
}

// ListRuns lists all runs for a user with pagination
func (s *Service) ListRuns(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Run, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 10 // default limit
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.ListByUserID(ctx, userID, limit, offset)
}
