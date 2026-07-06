package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool cria um pool de conexões com o PostgreSQL.
// Usamos pool (não conexão única) porque Go é concorrente —
// múltiplas requisições podem acontecer ao mesmo tempo.
func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar pool: %w", err)
	}

	// Testa a conexão
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("banco inacessível: %w", err)
	}

	return pool, nil
}
