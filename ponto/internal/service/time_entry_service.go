package service

import (
	"context"
	"time"

	"github.com/luskation/ponto/internal/apperr"
	"github.com/luskation/ponto/internal/domain"
	"github.com/luskation/ponto/internal/repository"
)

type TimeEntryService struct {
	repo *repository.TimeEntryRepository
}

func NewTimeEntryService(repo *repository.TimeEntryRepository) *TimeEntryService {
	return &TimeEntryService{repo: repo}
}

// RegisterEntry bate o ponto, alternando entre clock_in e clock_out.
func (s *TimeEntryService) RegisterEntry(ctx context.Context, userID string) (*domain.TimeEntry, error) {
	last, err := s.repo.GetLastByUser(ctx, userID)
	if err != nil {
		return nil, apperr.Internal("erro ao consultar último registro")
	}

	nextType := domain.ClockIn
	if last != nil && last.Type == domain.ClockIn {
		nextType = domain.ClockOut
	}

	entry := &domain.TimeEntry{
		UserID:     userID,
		Type:       nextType,
		RecordedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, apperr.Internal("erro ao registrar ponto")
	}
	return entry, nil
}

func (s *TimeEntryService) ListByUser(ctx context.Context, userID string, page, limit int) ([]domain.TimeEntry, int, error) {
	return s.repo.ListByUser(ctx, userID, page, limit)
}

func (s *TimeEntryService) ListAll(ctx context.Context, page, limit int) ([]domain.TimeEntry, int, error) {
	return s.repo.ListAll(ctx, page, limit)
}
