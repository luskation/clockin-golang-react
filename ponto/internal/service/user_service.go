package service

import (
	"context"

	"github.com/luskation/ponto/internal/apperr"
	"github.com/luskation/ponto/internal/domain"
	"github.com/luskation/ponto/internal/repository"
	"golang.org/x/crypto/bcrypt"
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
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperr.Internal("erro ao gerar hash")
	}
	u.Password = string(hashed)
	return s.repo.Create(ctx, u)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, apperr.NotFound("usuário não encontrado")
	}
	return u, nil
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

func (s *UserService) UpdatePassword(ctx context.Context, id, newPassword string) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return apperr.NotFound("usuário não encontrado")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperr.Internal("erro ao gerar hash")
	}
	return s.repo.UpdatePassword(ctx, id, string(hashed))
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return apperr.NotFound("usuário não encontrado")
	}
	return s.repo.Delete(ctx, id)
}
