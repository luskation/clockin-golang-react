package service

import (
	"context"

	"github.com/luskation/ponto/internal/apperr"
	"github.com/luskation/ponto/internal/domain"
	"github.com/luskation/ponto/internal/repository"
)

type CompanyService struct {
	repo *repository.CompanyRepository
}

func NewCompanyService(repo *repository.CompanyRepository) *CompanyService {
	return &CompanyService{repo: repo}
}

func (s *CompanyService) Create(ctx context.Context, c *domain.Company) error {
	if c.Name == "" || c.CNPJ == "" {
		return apperr.BadRequest("name e cnpj são obrigatórios")
	}
	return s.repo.Create(ctx, c)
}

func (s *CompanyService) GetByID(ctx context.Context, id string) (*domain.Company, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperr.NotFound("empresa não encontrada")
	}
	return c, nil
}

func (s *CompanyService) List(ctx context.Context) ([]domain.Company, error) {
	return s.repo.List(ctx)
}

func (s *CompanyService) Update(ctx context.Context, c *domain.Company) error {
	if c.Name == "" || c.CNPJ == "" {
		return apperr.BadRequest("name e cnpj são obrigatórios")
	}
	if _, err := s.repo.GetByID(ctx, c.ID); err != nil {
		return apperr.NotFound("empresa não encontrada")
	}
	return s.repo.Update(ctx, c)
}

func (s *CompanyService) Delete(ctx context.Context, id string) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return apperr.NotFound("empresa não encontrada")
	}
	return s.repo.Delete(ctx, id)
}
