package mongodb

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nicikess/out-run-management-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type mongoContainer struct {
	container testcontainers.Container
	uri       string
}

func setupMongoContainer(ctx context.Context) (*mongoContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:5.0",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("Waiting for connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "27017")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("mongodb://%s:%s", hostIP, mappedPort.Port())

	return &mongoContainer{
		container: container,
		uri:       uri,
	}, nil
}

func createTestRun(userID uuid.UUID) *domain.Run {
	return domain.NewRun(userID, domain.Coordinate{
		Latitude:  40.7128,
		Longitude: -74.0060,
		Timestamp: time.Now(),
	})
}

func TestNewRepository(t *testing.T) {
	ctx := context.Background()

	// Setup MongoDB container
	mongoC, err := setupMongoContainer(ctx)
	require.NoError(t, err)
	defer func() {
		if err := mongoC.container.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	// Test repository creation
	repo, err := NewRepository(ctx, mongoC.uri)
	require.NoError(t, err)
	defer repo.Close(ctx)

	assert.NotNil(t, repo)
}

func TestRepository_Create(t *testing.T) {
	ctx := context.Background()

	// Setup MongoDB container
	mongoC, err := setupMongoContainer(ctx)
	require.NoError(t, err)
	defer func() {
		if err := mongoC.container.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	// Create repository
	repo, err := NewRepository(ctx, mongoC.uri)
	require.NoError(t, err)
	defer repo.Close(ctx)

	// Create test run
	run := createTestRun(uuid.New())

	// Test creating run
	err = repo.Create(ctx, run)
	assert.NoError(t, err)

	// Verify run was created
	savedRun, err := repo.GetByID(ctx, run.ID)
	assert.NoError(t, err)
	assert.Equal(t, run.ID, savedRun.ID)
	assert.Equal(t, run.UserID, savedRun.UserID)
}

func TestRepository_GetByID(t *testing.T) {
	ctx := context.Background()

	// Setup MongoDB container
	mongoC, err := setupMongoContainer(ctx)
	require.NoError(t, err)
	defer func() {
		if err := mongoC.container.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	// Create repository
	repo, err := NewRepository(ctx, mongoC.uri)
	require.NoError(t, err)
	defer repo.Close(ctx)

	// Test cases
	t.Run("Existing run", func(t *testing.T) {
		run := createTestRun(uuid.New())
		err := repo.Create(ctx, run)
		require.NoError(t, err)

		savedRun, err := repo.GetByID(ctx, run.ID)
		assert.NoError(t, err)
		assert.Equal(t, run.ID, savedRun.ID)
	})

	t.Run("Non-existing run", func(t *testing.T) {
		_, err := repo.GetByID(ctx, uuid.New())
		assert.ErrorIs(t, err, domain.ErrRunNotFound)
	})
}

func TestRepository_GetActiveByUserID(t *testing.T) {
	ctx := context.Background()

	// Setup MongoDB container
	mongoC, err := setupMongoContainer(ctx)
	require.NoError(t, err)
	defer func() {
		if err := mongoC.container.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	// Create repository
	repo, err := NewRepository(ctx, mongoC.uri)
	require.NoError(t, err)
	defer repo.Close(ctx)

	// Test cases
	t.Run("User with active run", func(t *testing.T) {
		userID := uuid.New()
		run := createTestRun(userID)
		err := repo.Create(ctx, run)
		require.NoError(t, err)

		activeRun, err := repo.GetActiveByUserID(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, run.ID, activeRun.ID)
	})

	t.Run("User without active run", func(t *testing.T) {
		userID := uuid.New()
		_, err := repo.GetActiveByUserID(ctx, userID)
		assert.ErrorIs(t, err, domain.ErrRunNotFound)
	})
}

func TestRepository_Update(t *testing.T) {
	ctx := context.Background()

	// Setup MongoDB container
	mongoC, err := setupMongoContainer(ctx)
	require.NoError(t, err)
	defer func() {
		if err := mongoC.container.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	// Create repository
	repo, err := NewRepository(ctx, mongoC.uri)
	require.NoError(t, err)
	defer repo.Close(ctx)

	// Test cases
	t.Run("Update existing run", func(t *testing.T) {
		run := createTestRun(uuid.New())
		err := repo.Create(ctx, run)
		require.NoError(t, err)

		// Update run status
		run.Status = domain.RunStatusPaused
		err = repo.Update(ctx, run)
		assert.NoError(t, err)

		// Verify update
		updatedRun, err := repo.GetByID(ctx, run.ID)
		assert.NoError(t, err)
		assert.Equal(t, domain.RunStatusPaused, updatedRun.Status)
	})

	t.Run("Update non-existing run", func(t *testing.T) {
		run := createTestRun(uuid.New())
		err := repo.Update(ctx, run)
		assert.ErrorIs(t, err, domain.ErrRunNotFound)
	})
}

func TestRepository_ListByUserID(t *testing.T) {
	ctx := context.Background()

	// Setup MongoDB container
	mongoC, err := setupMongoContainer(ctx)
	require.NoError(t, err)
	defer func() {
		if err := mongoC.container.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	// Create repository
	repo, err := NewRepository(ctx, mongoC.uri)
	require.NoError(t, err)
	defer repo.Close(ctx)

	// Create test data
	userID := uuid.New()
	var runs []*domain.Run
	for i := 0; i < 5; i++ {
		run := createTestRun(userID)
		err := repo.Create(ctx, run)
		require.NoError(t, err)
		runs = append(runs, run)
		time.Sleep(time.Millisecond) // Ensure different timestamps
	}

	// Test cases
	t.Run("List with pagination", func(t *testing.T) {
		// Get first 2 runs
		firstPage, err := repo.ListByUserID(ctx, userID, 2, 0)
		assert.NoError(t, err)
		assert.Len(t, firstPage, 2)

		// Get next 2 runs
		secondPage, err := repo.ListByUserID(ctx, userID, 2, 2)
		assert.NoError(t, err)
		assert.Len(t, secondPage, 2)

		// Verify they're different runs
		assert.NotEqual(t, firstPage[0].ID, secondPage[0].ID)
	})

	t.Run("List for user with no runs", func(t *testing.T) {
		emptyUserID := uuid.New()
		runs, err := repo.ListByUserID(ctx, emptyUserID, 10, 0)
		assert.NoError(t, err)
		assert.Empty(t, runs)
	})
}
