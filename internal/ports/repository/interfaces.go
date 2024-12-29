package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nicikess/out-run-management-service/internal/domain"
)

// RunRepository defines the interface for run data persistence
type RunRepository interface {
	// Create stores a new run
	Create(ctx context.Context, run *domain.Run) error

	// GetByID retrieves a run by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Run, error)

	// GetActiveByUserID retrieves the active run for a user, if any
	GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*domain.Run, error)

	// Update updates an existing run
	Update(ctx context.Context, run *domain.Run) error

	// ListByUserID retrieves all runs for a user
	ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Run, error)

	// Close closes the repository connection
	Close(ctx context.Context) error
}
