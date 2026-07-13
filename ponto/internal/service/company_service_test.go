package service

import (
	"context"
	"errors"
	"testing"

	"github.com/luskation/ponto/internal/domain"
	"github.com/stretchr/testify/assert"
)

type fakeCompanyRepository struct {
	createFn  func(ctx context.Context, c *domain.Company) error
	getByIDFn func(ctx context.Context, id string) (*domain.Company, error)
	listFn    func(ctx context.Context, page, limit int) ([]domain.Company, int, error)
	updateFn  func(ctx context.Context, c *domain.Company) error
	deleteFn  func(ctx context.Context, id string) error
}

func (f *fakeCompanyRepository) Create(ctx context.Context, c *domain.Company) error {
	return f.createFn(ctx, c)
}

func (f *fakeCompanyRepository) GetByID(ctx context.Context, id string) (*domain.Company, error) {
	return f.getByIDFn(ctx, id)
}

func (f *fakeCompanyRepository) List(ctx context.Context, page, limit int) ([]domain.Company, int, error) {
	return f.listFn(ctx, page, limit)
}

func (f *fakeCompanyRepository) Update(ctx context.Context, c *domain.Company) error {
	return f.updateFn(ctx, c)
}

func (f *fakeCompanyRepository) Delete(ctx context.Context, id string) error {
	return f.deleteFn(ctx, id)
}

func TestCompanyService_Create_InvalidCNPJ(t *testing.T) {
	repo := &fakeCompanyRepository{
		createFn: func(ctx context.Context, c *domain.Company) error {
			t.Fatal("repo.Create não deveria ser chamado com CNPJ inválido")
			return nil
		},
	}
	s := NewCompanyService(repo)

	err := s.Create(context.Background(), &domain.Company{Name: "Teste", CNPJ: "123"})

	assert.Error(t, err)
}

func TestCompanyService_Create_MissingFields(t *testing.T) {
	s := NewCompanyService(&fakeCompanyRepository{})

	err := s.Create(context.Background(), &domain.Company{Name: "", CNPJ: ""})

	assert.Error(t, err)
}

func TestCompanyService_Create_Success(t *testing.T) {
	var created *domain.Company
	repo := &fakeCompanyRepository{
		createFn: func(ctx context.Context, c *domain.Company) error {
			created = c
			return nil
		},
	}
	s := NewCompanyService(repo)

	err := s.Create(context.Background(), &domain.Company{Name: "Acme", CNPJ: "12.345.678/0001-90"})

	assert.NoError(t, err)
	assert.Equal(t, "Acme", created.Name)
}

func TestCompanyService_Update_InvalidCNPJ(t *testing.T) {
	s := NewCompanyService(&fakeCompanyRepository{})

	err := s.Update(context.Background(), &domain.Company{ID: "1", Name: "Acme", CNPJ: "invalido"})

	assert.Error(t, err)
}

func TestCompanyService_Update_NotFound(t *testing.T) {
	repo := &fakeCompanyRepository{
		getByIDFn: func(ctx context.Context, id string) (*domain.Company, error) {
			return nil, errors.New("not found")
		},
	}
	s := NewCompanyService(repo)

	err := s.Update(context.Background(), &domain.Company{ID: "1", Name: "Acme", CNPJ: "12.345.678/0001-90"})

	assert.Error(t, err)
}

func TestCompanyService_Delete_NotFound(t *testing.T) {
	repo := &fakeCompanyRepository{
		getByIDFn: func(ctx context.Context, id string) (*domain.Company, error) {
			return nil, errors.New("not found")
		},
	}
	s := NewCompanyService(repo)

	err := s.Delete(context.Background(), "1")

	assert.Error(t, err)
}
