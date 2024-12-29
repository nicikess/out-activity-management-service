package domain

import (
	"math"
	"time"

	"github.com/google/uuid"
)

// RunStatus represents the current state of a run
type RunStatus string

const (
	RunStatusActive    RunStatus = "active"
	RunStatusPaused    RunStatus = "paused"
	RunStatusCompleted RunStatus = "completed"
)

// Coordinate represents a GPS coordinate with timestamp
type Coordinate struct {
	Latitude  float64   `json:"latitude" bson:"latitude"`
	Longitude float64   `json:"longitude" bson:"longitude"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

// RunStats contains statistics for a run
type RunStats struct {
	Distance    float64 `json:"distance" bson:"distance"`        // Total distance in meters
	Duration    int64   `json:"duration" bson:"duration"`        // Total duration in seconds
	AveragePace float64 `json:"averagePace" bson:"average_pace"` // Average pace in meters per second
}

// Run represents a running activity
type Run struct {
	ID        uuid.UUID    `json:"id" bson:"_id"`
	UserID    uuid.UUID    `json:"userId" bson:"user_id"`
	StartTime time.Time    `json:"startTime" bson:"start_time"`
	EndTime   *time.Time   `json:"endTime,omitempty" bson:"end_time,omitempty"`
	Status    RunStatus    `json:"status" bson:"status"`
	Route     []Coordinate `json:"route" bson:"route"`
	Stats     RunStats     `json:"stats" bson:"stats"`
}

// NewRun creates a new Run instance with the given user ID and initial coordinate
func NewRun(userID uuid.UUID, initialCoordinate Coordinate) *Run {
	return &Run{
		ID:        uuid.New(),
		UserID:    userID,
		StartTime: time.Now(),
		Status:    RunStatusActive,
		Route:     []Coordinate{initialCoordinate},
		Stats:     RunStats{},
	}
}

// AddCoordinate adds a new coordinate to the route and updates stats
func (r *Run) AddCoordinate(coord Coordinate) error {
	if r.Status != RunStatusActive {
		return ErrInvalidRunStatus
	}

	// Add coordinate to route
	r.Route = append(r.Route, coord)

	// Update stats if we have at least two coordinates
	if len(r.Route) >= 2 {
		prev := r.Route[len(r.Route)-2]
		distance := calculateDistance(prev, coord)
		r.Stats.Distance += distance

		// Update duration
		r.Stats.Duration = int64(coord.Timestamp.Sub(r.StartTime).Seconds())

		// Update average pace (m/s)
		if r.Stats.Duration > 0 {
			r.Stats.AveragePace = r.Stats.Distance / float64(r.Stats.Duration)
		}
	}

	return nil
}

// Pause transitions the run to paused state
func (r *Run) Pause() error {
	if r.Status != RunStatusActive {
		return ErrInvalidRunStatus
	}
	r.Status = RunStatusPaused
	return nil
}

// Resume transitions the run from paused to active state
func (r *Run) Resume() error {
	if r.Status != RunStatusPaused {
		return ErrInvalidRunStatus
	}
	r.Status = RunStatusActive
	return nil
}

// End completes the run and finalizes stats
func (r *Run) End() error {
	if r.Status == RunStatusCompleted {
		return ErrInvalidRunStatus
	}

	now := time.Now()
	r.EndTime = &now
	r.Status = RunStatusCompleted

	// Finalize duration calculation
	r.Stats.Duration = int64(now.Sub(r.StartTime).Seconds())

	// Final average pace calculation
	if r.Stats.Duration > 0 {
		r.Stats.AveragePace = r.Stats.Distance / float64(r.Stats.Duration)
	}

	return nil
}

// IsActive checks if the run is currently active
func (r *Run) IsActive() bool {
	return r.Status == RunStatusActive
}

// IsUserAuthorized checks if the given user ID matches the run's user
func (r *Run) IsUserAuthorized(userID uuid.UUID) bool {
	return r.UserID == userID
}

// calculateDistance calculates the distance between two coordinates in meters
func calculateDistance(c1, c2 Coordinate) float64 {
	// Implementation of the Haversine formula
	const earthRadius = 6371000 // Earth's radius in meters

	lat1 := toRadians(c1.Latitude)
	lon1 := toRadians(c1.Longitude)
	lat2 := toRadians(c2.Latitude)
	lon2 := toRadians(c2.Longitude)

	dlat := lat2 - lat1
	dlon := lon2 - lon1

	a := sinSquared(dlat/2) +
		cos(lat1)*cos(lat2)*sinSquared(dlon/2)
	c := 2 * asin(sqrt(a))

	return earthRadius * c
}

// Helper functions for calculateDistance
func toRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func sinSquared(x float64) float64 {
	sin := math.Sin(x)
	return sin * sin
}

func cos(x float64) float64 {
	return math.Cos(x)
}

func asin(x float64) float64 {
	return math.Asin(x)
}

func sqrt(x float64) float64 {
	return math.Sqrt(x)
}
