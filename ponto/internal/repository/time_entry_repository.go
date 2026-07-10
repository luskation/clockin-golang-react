package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luskation/ponto/internal/domain"
)

type TimeEntryRepository struct {
	db *pgxpool.Pool
}

func NewTimeEntryRepository(db *pgxpool.Pool) *TimeEntryRepository {
	return &TimeEntryRepository{db: db}
}

func (r *TimeEntryRepository) Create(ctx context.Context, e *domain.TimeEntry) error {
	query := `INSERT INTO time_entries (user_id, type, recorded_at, note) VALUES ($1, $2, $3, $4)
              RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, e.UserID, e.Type, e.RecordedAt, e.Note).
		Scan(&e.ID, &e.CreatedAt)
}

func (r *TimeEntryRepository) GetLastByUser(ctx context.Context, userID string) (*domain.TimeEntry, error) {
	query := `SELECT id, user_id, type, recorded_at, note, created_at FROM time_entries
              WHERE user_id = $1 ORDER BY recorded_at DESC LIMIT 1`
	e := &domain.TimeEntry{}
	err := r.db.QueryRow(ctx, query, userID).
		Scan(&e.ID, &e.UserID, &e.Type, &e.RecordedAt, &e.Note, &e.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return e, nil
}

func (r *TimeEntryRepository) ListByUser(ctx context.Context, userID string, page, limit int) ([]domain.TimeEntry, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM time_entries WHERE user_id = $1`, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `SELECT id, user_id, type, recorded_at, note, created_at FROM time_entries
              WHERE user_id = $1 ORDER BY recorded_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var entries []domain.TimeEntry
	for rows.Next() {
		var e domain.TimeEntry
		if err := rows.Scan(&e.ID, &e.UserID, &e.Type, &e.RecordedAt, &e.Note, &e.CreatedAt); err != nil {
			return nil, 0, err
		}
		entries = append(entries, e)
	}
	return entries, total, nil
}

func (r *TimeEntryRepository) ListAll(ctx context.Context, page, limit int) ([]domain.TimeEntry, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM time_entries`).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `SELECT id, user_id, type, recorded_at, note, created_at FROM time_entries
              ORDER BY recorded_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var entries []domain.TimeEntry
	for rows.Next() {
		var e domain.TimeEntry
		if err := rows.Scan(&e.ID, &e.UserID, &e.Type, &e.RecordedAt, &e.Note, &e.CreatedAt); err != nil {
			return nil, 0, err
		}
		entries = append(entries, e)
	}
	return entries, total, nil
}
