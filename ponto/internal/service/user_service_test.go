package service

import (
	"context"
	"errors"
	"testing"

	"github.com/luskation/ponto/internal/domain"
	"github.com/stretchr/testify/assert"
)

type fakeUserRepository struct {
	createFn         func(ctx context.Context, u *domain.User) error
	getByIDFn        func(ctx context.Context, id string) (*domain.User, error)
	getByEmailFn     func(ctx context.Context, email string) (*domain.User, error)
	listFn           func(ctx context.Context, page, limit int) ([]domain.User, int, error)
	updateFn         func(ctx context.Context, u *domain.User) error
	updatePasswordFn func(ctx context.Context, id, hashedPassword string) error
	deleteFn         func(ctx context.Context, id string) error
}

func (f *fakeUserRepository) Create(ctx context.Context, u *domain.User) error {
	return f.createFn(ctx, u)
}

func (f *fakeUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return f.getByIDFn(ctx, id)
}

func (f *fakeUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return f.getByEmailFn(ctx, email)
}

func (f *fakeUserRepository) List(ctx context.Context, page, limit int) ([]domain.User, int, error) {
	return f.listFn(ctx, page, limit)
}

func (f *fakeUserRepository) Update(ctx context.Context, u *domain.User) error {
	return f.updateFn(ctx, u)
}

func (f *fakeUserRepository) UpdatePassword(ctx context.Context, id, hashedPassword string) error {
	return f.updatePasswordFn(ctx, id, hashedPassword)
}

func (f *fakeUserRepository) Delete(ctx context.Context, id string) error {
	return f.deleteFn(ctx, id)
}

func validUser() *domain.User {
	return &domain.User{CompanyID: "company-1", Name: "Ana", Email: "ana@example.com", Password: "segredo123"}
}

func TestUserService_Create_MissingFields(t *testing.T) {
	s := NewUserService(&fakeUserRepository{})

	err := s.Create(context.Background(), &domain.User{})

	assert.Error(t, err)
}

func TestUserService_Create_InvalidEmail(t *testing.T) {
	s := NewUserService(&fakeUserRepository{})
	u := validUser()
	u.Email = "não-é-email"

	err := s.Create(context.Background(), u)

	assert.Error(t, err)
}

func TestUserService_Create_ShortPassword(t *testing.T) {
	s := NewUserService(&fakeUserRepository{})
	u := validUser()
	u.Password = "123"

	err := s.Create(context.Background(), u)

	assert.Error(t, err)
}

func TestUserService_Create_DefaultsRoleToEmployee(t *testing.T) {
	var created *domain.User
	repo := &fakeUserRepository{
		createFn: func(ctx context.Context, u *domain.User) error {
			created = u
			return nil
		},
	}
	s := NewUserService(repo)

	err := s.Create(context.Background(), validUser())

	assert.NoError(t, err)
	assert.Equal(t, domain.RoleEmployee, created.Role)
}

func TestUserService_Create_HashesPassword(t *testing.T) {
	var created *domain.User
	repo := &fakeUserRepository{
		createFn: func(ctx context.Context, u *domain.User) error {
			created = u
			return nil
		},
	}
	s := NewUserService(repo)
	plain := "segredo123"
	u := validUser()
	u.Password = plain

	err := s.Create(context.Background(), u)

	assert.NoError(t, err)
	assert.NotEqual(t, plain, created.Password)
	assert.NotEmpty(t, created.Password)
}

func TestUserService_Update_InvalidEmail(t *testing.T) {
	s := NewUserService(&fakeUserRepository{})

	err := s.Update(context.Background(), &domain.User{ID: "1", Name: "Ana", Email: "inválido"})

	assert.Error(t, err)
}

func TestUserService_Update_NotFound(t *testing.T) {
	repo := &fakeUserRepository{
		getByIDFn: func(ctx context.Context, id string) (*domain.User, error) {
			return nil, errors.New("not found")
		},
	}
	s := NewUserService(repo)

	err := s.Update(context.Background(), &domain.User{ID: "1", Name: "Ana", Email: "ana@example.com"})

	assert.Error(t, err)
}

func TestUserService_Delete_CannotDeleteSelf(t *testing.T) {
	repo := &fakeUserRepository{
		deleteFn: func(ctx context.Context, id string) error {
			t.Fatal("repo.Delete não deveria ser chamado ao tentar autoexcluir-se")
			return nil
		},
	}
	s := NewUserService(repo)

	err := s.Delete(context.Background(), "user-1", "user-1")

	assert.Error(t, err)
}

func TestUserService_Delete_NotFound(t *testing.T) {
	repo := &fakeUserRepository{
		getByIDFn: func(ctx context.Context, id string) (*domain.User, error) {
			return nil, errors.New("not found")
		},
	}
	s := NewUserService(repo)

	err := s.Delete(context.Background(), "admin-1", "user-2")

	assert.Error(t, err)
}

func TestUserService_Delete_Success(t *testing.T) {
	deleted := false
	repo := &fakeUserRepository{
		getByIDFn: func(ctx context.Context, id string) (*domain.User, error) {
			return &domain.User{ID: id}, nil
		},
		deleteFn: func(ctx context.Context, id string) error {
			deleted = true
			return nil
		},
	}
	s := NewUserService(repo)

	err := s.Delete(context.Background(), "admin-1", "user-2")

	assert.NoError(t, err)
	assert.True(t, deleted)
}
