package domain

import "errors"

var (
	// ErrRunNotFound indicates that the requested run doesn't exist
	ErrRunNotFound = errors.New("run not found")

	// ErrActiveRunExists indicates that user already has an active run
	ErrActiveRunExists = errors.New("active run already exists")

	// ErrInvalidRunStatus indicates invalid run status transition
	ErrInvalidRunStatus = errors.New("invalid run status")

	// ErrUnauthorized indicates that the user is not authorized to access the run
	ErrUnauthorized = errors.New("unauthorized access to run")
)

// IsNotFoundError checks if the error is a not found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrRunNotFound)
}

// IsConflictError checks if the error is a conflict error
func IsConflictError(err error) bool {
	return errors.Is(err, ErrActiveRunExists) || errors.Is(err, ErrInvalidRunStatus)
}

// IsUnauthorizedError checks if the error is an unauthorized error
func IsUnauthorizedError(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}
