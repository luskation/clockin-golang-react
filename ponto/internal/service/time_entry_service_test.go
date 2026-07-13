package service

import (
	"context"
	"errors"
	"testing"

	"github.com/luskation/ponto/internal/domain"
	"github.com/stretchr/testify/assert"
)

type fakeTimeEntryRepository struct {
	createFn        func(ctx context.Context, e *domain.TimeEntry) error
	getLastByUserFn func(ctx context.Context, userID string) (*domain.TimeEntry, error)
	getByIDFn       func(ctx context.Context, id string) (*domain.TimeEntry, error)
	updateFn        func(ctx context.Context, e *domain.TimeEntry) error
	deleteFn        func(ctx context.Context, id string) error
	listByUserFn    func(ctx context.Context, userID string, page, limit int) ([]domain.TimeEntry, int, error)
	listAllFn       func(ctx context.Context, page, limit int) ([]domain.TimeEntry, int, error)
}

func (f *fakeTimeEntryRepository) Create(ctx context.Context, e *domain.TimeEntry) error {
	return f.createFn(ctx, e)
}

func (f *fakeTimeEntryRepository) GetLastByUser(ctx context.Context, userID string) (*domain.TimeEntry, error) {
	return f.getLastByUserFn(ctx, userID)
}

func (f *fakeTimeEntryRepository) GetByID(ctx context.Context, id string) (*domain.TimeEntry, error) {
	return f.getByIDFn(ctx, id)
}

func (f *fakeTimeEntryRepository) Update(ctx context.Context, e *domain.TimeEntry) error {
	return f.updateFn(ctx, e)
}

func (f *fakeTimeEntryRepository) Delete(ctx context.Context, id string) error {
	return f.deleteFn(ctx, id)
}

func (f *fakeTimeEntryRepository) ListByUser(ctx context.Context, userID string, page, limit int) ([]domain.TimeEntry, int, error) {
	return f.listByUserFn(ctx, userID, page, limit)
}

func (f *fakeTimeEntryRepository) ListAll(ctx context.Context, page, limit int) ([]domain.TimeEntry, int, error) {
	return f.listAllFn(ctx, page, limit)
}

func TestTimeEntryService_RegisterEntry_FirstEntryIsClockIn(t *testing.T) {
	var created *domain.TimeEntry
	repo := &fakeTimeEntryRepository{
		getLastByUserFn: func(ctx context.Context, userID string) (*domain.TimeEntry, error) {
			return nil, nil
		},
		createFn: func(ctx context.Context, e *domain.TimeEntry) error {
			created = e
			return nil
		},
	}
	s := NewTimeEntryService(repo)

	entry, err := s.RegisterEntry(context.Background(), "user-1")

	assert.NoError(t, err)
	assert.Equal(t, domain.ClockIn, entry.Type)
	assert.Equal(t, domain.ClockIn, created.Type)
}

func TestTimeEntryService_RegisterEntry_TogglesToClockOut(t *testing.T) {
	repo := &fakeTimeEntryRepository{
		getLastByUserFn: func(ctx context.Context, userID string) (*domain.TimeEntry, error) {
			return &domain.TimeEntry{Type: domain.ClockIn}, nil
		},
		createFn: func(ctx context.Context, e *domain.TimeEntry) error {
			return nil
		},
	}
	s := NewTimeEntryService(repo)

	entry, err := s.RegisterEntry(context.Background(), "user-1")

	assert.NoError(t, err)
	assert.Equal(t, domain.ClockOut, entry.Type)
}

func TestTimeEntryService_RegisterEntry_TogglesToClockIn(t *testing.T) {
	repo := &fakeTimeEntryRepository{
		getLastByUserFn: func(ctx context.Context, userID string) (*domain.TimeEntry, error) {
			return &domain.TimeEntry{Type: domain.ClockOut}, nil
		},
		createFn: func(ctx context.Context, e *domain.TimeEntry) error {
			return nil
		},
	}
	s := NewTimeEntryService(repo)

	entry, err := s.RegisterEntry(context.Background(), "user-1")

	assert.NoError(t, err)
	assert.Equal(t, domain.ClockIn, entry.Type)
}

func TestTimeEntryService_Update_ForbiddenForOtherEmployee(t *testing.T) {
	repo := &fakeTimeEntryRepository{
		getByIDFn: func(ctx context.Context, id string) (*domain.TimeEntry, error) {
			return &domain.TimeEntry{ID: id, UserID: "owner-1"}, nil
		},
	}
	s := NewTimeEntryService(repo)

	err := s.Update(context.Background(), "other-user", "employee", &domain.TimeEntry{ID: "entry-1", Type: domain.ClockIn})

	assert.Error(t, err)
}

func TestTimeEntryService_Update_AllowedForOwner(t *testing.T) {
	updated := false
	repo := &fakeTimeEntryRepository{
		getByIDFn: func(ctx context.Context, id string) (*domain.TimeEntry, error) {
			return &domain.TimeEntry{ID: id, UserID: "owner-1"}, nil
		},
		updateFn: func(ctx context.Context, e *domain.TimeEntry) error {
			updated = true
			return nil
		},
	}
	s := NewTimeEntryService(repo)

	err := s.Update(context.Background(), "owner-1", "employee", &domain.TimeEntry{ID: "entry-1", Type: domain.ClockOut})

	assert.NoError(t, err)
	assert.True(t, updated)
}

func TestTimeEntryService_Update_AllowedForAdmin(t *testing.T) {
	updated := false
	repo := &fakeTimeEntryRepository{
		getByIDFn: func(ctx context.Context, id string) (*domain.TimeEntry, error) {
			return &domain.TimeEntry{ID: id, UserID: "owner-1"}, nil
		},
		updateFn: func(ctx context.Context, e *domain.TimeEntry) error {
			updated = true
			return nil
		},
	}
	s := NewTimeEntryService(repo)

	err := s.Update(context.Background(), "admin-1", "admin", &domain.TimeEntry{ID: "entry-1", Type: domain.ClockOut})

	assert.NoError(t, err)
	assert.True(t, updated)
}

func TestTimeEntryService_Update_InvalidType(t *testing.T) {
	repo := &fakeTimeEntryRepository{
		getByIDFn: func(ctx context.Context, id string) (*domain.TimeEntry, error) {
			return &domain.TimeEntry{ID: id, UserID: "owner-1"}, nil
		},
	}
	s := NewTimeEntryService(repo)

	err := s.Update(context.Background(), "owner-1", "employee", &domain.TimeEntry{ID: "entry-1", Type: "invalido"})

	assert.Error(t, err)
}

func TestTimeEntryService_Delete_ForbiddenForOtherEmployee(t *testing.T) {
	repo := &fakeTimeEntryRepository{
		getByIDFn: func(ctx context.Context, id string) (*domain.TimeEntry, error) {
			return &domain.TimeEntry{ID: id, UserID: "owner-1"}, nil
		},
		deleteFn: func(ctx context.Context, id string) error {
			t.Fatal("repo.Delete não deveria ser chamado sem permissão")
			return nil
		},
	}
	s := NewTimeEntryService(repo)

	err := s.Delete(context.Background(), "other-user", "employee", "entry-1")

	assert.Error(t, err)
}

func TestTimeEntryService_Delete_NotFound(t *testing.T) {
	repo := &fakeTimeEntryRepository{
		getByIDFn: func(ctx context.Context, id string) (*domain.TimeEntry, error) {
			return nil, errors.New("not found")
		},
	}
	s := NewTimeEntryService(repo)

	err := s.Delete(context.Background(), "user-1", "employee", "entry-1")

	assert.Error(t, err)
}
