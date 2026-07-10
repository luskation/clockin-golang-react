package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luskation/ponto/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	query := `INSERT INTO users (company_id, name, email, password, role) VALUES ($1, $2, $3, $4, $5)
              RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, u.CompanyID, u.Name, u.Email, u.Password, u.Role).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, company_id, name, email, password, role, created_at, updated_at FROM users WHERE id = $1`
	u := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).
		Scan(&u.ID, &u.CompanyID, &u.Name, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, company_id, name, email, password, role, created_at, updated_at FROM users WHERE email = $1`
	u := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).
		Scan(&u.ID, &u.CompanyID, &u.Name, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) List(ctx context.Context, page, limit int) ([]domain.User, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `SELECT id, company_id, name, email, password, role, created_at, updated_at
	          FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.CompanyID, &u.Name, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (r *UserRepository) Update(ctx context.Context, u *domain.User) error {
	query := `UPDATE users SET name = $1, email = $2, role = $3, updated_at = NOW() WHERE id = $4
	          RETURNING company_id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, u.Name, u.Email, u.Role, u.ID).
		Scan(&u.CompanyID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id, hashedPassword string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`, hashedPassword, id)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}
