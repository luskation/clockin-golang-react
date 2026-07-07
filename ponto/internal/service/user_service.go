package service

import (
	"context"

	"github.com/luskation/ponto/internal/apperr"
	"github.com/luskation/ponto/internal/domain"
	"github.com/luskation/ponto/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, u *domain.User) error {
	if u.CompanyID == "" || u.Name == "" || u.Email == "" || u.Password == "" {
		return apperr.BadRequest("company_id, name, email e password são obrigatórios")
	}
	if u.Role == "" {
		u.Role = domain.RoleEmployee
	}
	return s.repo.Create(ctx, u)
}

func (s *UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperr.NotFound("usuário não encontrado")
	}
	return u, nil
}

func (s *UserService) List(ctx context.Context) ([]domain.User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) Update(ctx context.Context, u *domain.User) error {
	if u.Name == "" || u.Email == "" {
		return apperr.BadRequest("name e email são obrigatórios")
	}
	if _, err := s.repo.GetByID(ctx, u.ID); err != nil {
		return apperr.NotFound("usuário não encontrado")
	}
	return s.repo.Update(ctx, u)
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return apperr.NotFound("usuário não encontrado")
	}
	return s.repo.Delete(ctx, id)
}
