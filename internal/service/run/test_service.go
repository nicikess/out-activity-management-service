package run

import (
	"context"
	"github.com/google/uuid"
	"github.com/nicikess/out-run-management-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// MockRepository is a mock implementation of RunRepository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, run *domain.Run) error {
	args := m.Called(ctx, run)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Run, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Run), args.Error(1)
}

func (m *MockRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*domain.Run, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Run), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, run *domain.Run) error {
	args := m.Called(ctx, run)
	return args.Error(0)
}

func (m *MockRepository) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Run, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Run), args.Error(1)
}

func (m *MockRepository) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func createTestCoordinate() domain.Coordinate {
	return domain.Coordinate{
		Latitude:  40.7128,
		Longitude: -74.0060,
		Timestamp: time.Now(),
	}
}

func TestService_StartRun(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	t.Run("Successfully start new run", func(t *testing.T) {
		userID := uuid.New()
		coord := createTestCoordinate()

		// Expect check for active run
		mockRepo.On("GetActiveByUserID", ctx, userID).Return(nil, domain.ErrRunNotFound)

		// Expect run creation
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Run")).Return(nil)

		run, err := service.StartRun(ctx, userID, coord)
		require.NoError(t, err)
		assert.NotNil(t, run)
		assert.Equal(t, userID, run.UserID)
		assert.Equal(t, domain.RunStatusActive, run.Status)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail when active run exists", func(t *testing.T) {
		userID := uuid.New()
		coord := createTestCoordinate()
		existingRun := domain.NewRun(userID, coord)

		mockRepo.On("GetActiveByUserID", ctx, userID).Return(existingRun, nil)

		run, err := service.StartRun(ctx, userID, coord)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrActiveRunExists, err)
		assert.Nil(t, run)

		mockRepo.AssertExpectations(t)
	})
}

func TestService_PauseRun(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	t.Run("Successfully pause run", func(t *testing.T) {
		userID := uuid.New()
		run := domain.NewRun(userID, createTestCoordinate())

		mockRepo.On("GetByID", ctx, run.ID).Return(run, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Run")).Return(nil)

		updatedRun, err := service.PauseRun(ctx, run.ID, userID)
		require.NoError(t, err)
		assert.Equal(t, domain.RunStatusPaused, updatedRun.Status)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail with unauthorized user", func(t *testing.T) {
		userID := uuid.New()
		wrongUserID := uuid.New()
		run := domain.NewRun(userID, createTestCoordinate())

		mockRepo.On("GetByID", ctx, run.ID).Return(run, nil)

		_, err := service.PauseRun(ctx, run.ID, wrongUserID)
		assert.Equal(t, domain.ErrUnauthorized, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestService_GetActiveRun(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	t.Run("Successfully get active run", func(t *testing.T) {
		userID := uuid.New()
		activeRun := domain.NewRun(userID, createTestCoordinate())

		mockRepo.On("GetActiveByUserID", ctx, userID).Return(activeRun, nil)

		run, err := service.GetActiveRun(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, activeRun.ID, run.ID)
		assert.Equal(t, domain.RunStatusActive, run.Status)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Return error when no active run", func(t *testing.T) {
		userID := uuid.New()

		mockRepo.On("GetActiveByUserID", ctx, userID).Return(nil, domain.ErrRunNotFound)

		run, err := service.GetActiveRun(ctx, userID)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrRunNotFound, err)
		assert.Nil(t, run)

		mockRepo.AssertExpectations(t)
	})
}

func TestService_AddCoordinate(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	t.Run("Successfully add coordinate", func(t *testing.T) {
		userID := uuid.New()
		run := domain.NewRun(userID, createTestCoordinate())
		newCoord := createTestCoordinate()

		mockRepo.On("GetByID", ctx, run.ID).Return(run, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Run")).Return(nil)

		updatedRun, err := service.AddCoordinate(ctx, run.ID, userID, newCoord)
		require.NoError(t, err)
		assert.Len(t, updatedRun.Route, 2)
		assert.Equal(t, newCoord, updatedRun.Route[1])

		mockRepo.AssertExpectations(t)
	})
}
