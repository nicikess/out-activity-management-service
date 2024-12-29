package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewRun(t *testing.T) {
	// Setup
	userID := uuid.New()
	initialCoord := Coordinate{
		Latitude:  40.7128,
		Longitude: -74.0060,
		Timestamp: time.Now(),
	}

	// Execute
	run := NewRun(userID, initialCoord)

	// Assert
	assert.NotNil(t, run)
	assert.Equal(t, userID, run.UserID)
	assert.Equal(t, RunStatusActive, run.Status)
	assert.Equal(t, 1, len(run.Route))
	assert.Equal(t, initialCoord, run.Route[0])
	assert.Equal(t, float64(0), run.Stats.Distance)
	assert.Equal(t, int64(0), run.Stats.Duration)
}

func TestRun_AddCoordinate(t *testing.T) {
	// Setup
	run := createTestRun()

	newCoord := Coordinate{
		Latitude:  40.7130, // Slightly different from initial coordinate
		Longitude: -74.0065,
		Timestamp: time.Now().Add(1 * time.Minute),
	}

	// Execute
	err := run.AddCoordinate(newCoord)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(run.Route))
	assert.True(t, run.Stats.Distance > 0, "Distance should be calculated")
	assert.True(t, run.Stats.Duration > 0, "Duration should be calculated")
	assert.True(t, run.Stats.AveragePace > 0, "Average pace should be calculated")
}

func TestRun_AddCoordinate_InvalidStatus(t *testing.T) {
	// Setup
	run := createTestRun()
	run.Status = RunStatusCompleted

	newCoord := Coordinate{
		Latitude:  40.7130,
		Longitude: -74.0065,
		Timestamp: time.Now(),
	}

	// Execute
	err := run.AddCoordinate(newCoord)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidRunStatus, err)
}

func TestRun_StateTransitions(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus RunStatus
		operation     func(*Run) error
		expectedError error
		finalStatus   RunStatus
	}{
		{
			name:          "Active to Paused",
			initialStatus: RunStatusActive,
			operation:     (*Run).Pause,
			expectedError: nil,
			finalStatus:   RunStatusPaused,
		},
		{
			name:          "Paused to Active",
			initialStatus: RunStatusPaused,
			operation:     (*Run).Resume,
			expectedError: nil,
			finalStatus:   RunStatusActive,
		},
		{
			name:          "Cannot Pause Completed Run",
			initialStatus: RunStatusCompleted,
			operation:     (*Run).Pause,
			expectedError: ErrInvalidRunStatus,
			finalStatus:   RunStatusCompleted,
		},
		{
			name:          "Cannot Resume Active Run",
			initialStatus: RunStatusActive,
			operation:     (*Run).Resume,
			expectedError: ErrInvalidRunStatus,
			finalStatus:   RunStatusActive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			run := createTestRun()
			run.Status = tt.initialStatus

			err := tt.operation(run)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.finalStatus, run.Status)
			}
		})
	}
}

func TestRun_End(t *testing.T) {
	// Setup
	run := createTestRun()

	run.StartTime = time.Now().Add(-2 * time.Minute)

	// Add some coordinates to make it interesting
	run.AddCoordinate(Coordinate{
		Latitude:  40.7130,
		Longitude: -74.0065,
		Timestamp: time.Now().Add(1 * time.Minute),
	})

	// Execute
	err := run.End()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, RunStatusCompleted, run.Status)
	assert.NotNil(t, run.EndTime)
	assert.True(t, run.Stats.Duration > 0)
	assert.True(t, run.Stats.Distance > 0)
	assert.True(t, run.Stats.AveragePace > 0)
}

func TestRun_Authorization(t *testing.T) {
	// Setup
	correctUserID := uuid.New()
	wrongUserID := uuid.New()
	run := createTestRunWithUserID(correctUserID)

	// Assert
	assert.True(t, run.IsUserAuthorized(correctUserID))
	assert.False(t, run.IsUserAuthorized(wrongUserID))
}

func TestCalculateDistance(t *testing.T) {
	tests := []struct {
		name     string
		c1       Coordinate
		c2       Coordinate
		expected float64
		delta    float64
	}{
		{
			name: "Same point",
			c1: Coordinate{
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			c2: Coordinate{
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			expected: 0,
			delta:    0.1,
		},
		{
			name: "Short distance",
			c1: Coordinate{
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			c2: Coordinate{
				Latitude:  40.7130,
				Longitude: -74.0065,
			},
			expected: 50, // Approximately 50 meters
			delta:    5,  // Allow 5 meters deviation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance := calculateDistance(tt.c1, tt.c2)
			assert.InDelta(t, tt.expected, distance, tt.delta)
		})
	}
}

// Helper functions for tests

func createTestRun() *Run {
	return createTestRunWithUserID(uuid.New())
}

func createTestRunWithUserID(userID uuid.UUID) *Run {
	return NewRun(userID, Coordinate{
		Latitude:  40.7128,
		Longitude: -74.0060,
		Timestamp: time.Now(),
	})
}
