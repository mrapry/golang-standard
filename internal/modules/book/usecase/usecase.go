package usecase

import (
	"context"
	"master-service/internal/modules/book/domain"

	"github.com/mrapry/go-lib/golibshared"
)

// BookUsecase abstract interface
type BookUsecase interface {
	FindAll(ctx context.Context, filter *domain.Filter) ([]*domain.Book, *golibshared.Meta, error)
	FindByID(ctx context.Context, ID string) (*domain.Book, error)
	Create(ctx context.Context, data *domain.Book) (*domain.Book, error)
	Update(ctx context.Context, data *domain.Book) (*domain.Book, error)
	RemoveByID(ctx context.Context, ID string) error
	RestoreByID(ctx context.Context, ID string) error
}
