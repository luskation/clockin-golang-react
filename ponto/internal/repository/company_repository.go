package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luskation/ponto/internal/domain"
)

type CompanyRepository struct {
	db *pgxpool.Pool
}

func NewCompanyRepository(db *pgxpool.Pool) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create(ctx context.Context, c *domain.Company) error {
	query := `INSERT INTO companies (name, cnpj) VALUES ($1, $2)
              RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, c.Name, c.CNPJ).
		Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *CompanyRepository) GetByID(ctx context.Context, id string) (*domain.Company, error) {
	query := `SELECT id, name, cnpj, created_at, updated_at FROM companies WHERE id = $1`
	c := &domain.Company{}
	err := r.db.QueryRow(ctx, query, id).
		Scan(&c.ID, &c.Name, &c.CNPJ, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CompanyRepository) List(ctx context.Context, page, limit int) ([]domain.Company, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM companies`).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `SELECT id, name, cnpj, created_at, updated_at
	          FROM companies ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var companies []domain.Company
	for rows.Next() {
		var c domain.Company
		if err := rows.Scan(&c.ID, &c.Name, &c.CNPJ, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, err
		}
		companies = append(companies, c)
	}
	return companies, total, nil
}

func (r *CompanyRepository) Update(ctx context.Context, c *domain.Company) error {
	query := `UPDATE companies SET name = $1, cnpj = $2, updated_at = NOW() WHERE id = $3
	          RETURNING created_at, updated_at`
	return r.db.QueryRow(ctx, query, c.Name, c.CNPJ, c.ID).
		Scan(&c.CreatedAt, &c.UpdatedAt)
}

func (r *CompanyRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM companies WHERE id = $1`, id)
	return err
}
