package service

import (
	"context"
	"time"

	"github.com/luskation/ponto/internal/apperr"
	"github.com/luskation/ponto/internal/domain"
)

type timeEntryRepository interface {
	Create(ctx context.Context, e *domain.TimeEntry) error
	GetLastByUser(ctx context.Context, userID string) (*domain.TimeEntry, error)
	GetByID(ctx context.Context, id string) (*domain.TimeEntry, error)
	Update(ctx context.Context, e *domain.TimeEntry) error
	Delete(ctx context.Context, id string) error
	ListByUser(ctx context.Context, userID string, page, limit int) ([]domain.TimeEntry, int, error)
	ListAll(ctx context.Context, page, limit int) ([]domain.TimeEntry, int, error)
}

type TimeEntryService struct {
	repo timeEntryRepository
}

func NewTimeEntryService(repo timeEntryRepository) *TimeEntryService {
	return &TimeEntryService{repo: repo}
}

// RegisterEntry decide o tipo do registro a partir do último, em vez de o
// cliente informar clock_in/clock_out — evita que uma falha no app do usuário
// gere dois clock_in seguidos e distorça o histórico de horas trabalhadas.
func (s *TimeEntryService) RegisterEntry(ctx context.Context, userID string) (*domain.TimeEntry, error) {
	last, err := s.repo.GetLastByUser(ctx, userID)
	if err != nil {
		return nil, apperr.Internal("erro ao consultar último registro", err)
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
		return nil, apperr.Internal("erro ao registrar ponto", err)
	}
	return entry, nil
}

func (s *TimeEntryService) Update(ctx context.Context, requesterID, requesterRole string, e *domain.TimeEntry) error {
	existing, err := s.repo.GetByID(ctx, e.ID)
	if err != nil {
		return apperr.NotFound("registro de ponto não encontrado")
	}
	if existing.UserID != requesterID && requesterRole != string(domain.RoleAdmin) {
		return apperr.Forbidden("apenas o dono do registro ou um admin pode editá-lo")
	}
	if e.Type != domain.ClockIn && e.Type != domain.ClockOut {
		return apperr.BadRequest("type deve ser clock_in ou clock_out")
	}
	e.UserID = existing.UserID
	return s.repo.Update(ctx, e)
}

func (s *TimeEntryService) Delete(ctx context.Context, requesterID, requesterRole, id string) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperr.NotFound("registro de ponto não encontrado")
	}
	if existing.UserID != requesterID && requesterRole != string(domain.RoleAdmin) {
		return apperr.Forbidden("apenas o dono do registro ou um admin pode excluí-lo")
	}
	return s.repo.Delete(ctx, id)
}

func (s *TimeEntryService) ListByUser(ctx context.Context, userID string, page, limit int) ([]domain.TimeEntry, int, error) {
	return s.repo.ListByUser(ctx, userID, page, limit)
}

func (s *TimeEntryService) ListAll(ctx context.Context, page, limit int) ([]domain.TimeEntry, int, error) {
	return s.repo.ListAll(ctx, page, limit)
}
