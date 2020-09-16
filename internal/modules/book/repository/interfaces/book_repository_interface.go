package interfaces

import (
	"context"
	"master-service/internal/modules/book/domain"

	"github.com/mrapry/go-lib/golibshared"
)

// BookRepository abstract interface
type BookRepository interface {
	FindAll(ctx context.Context, filter *domain.Filter) <-chan golibshared.Result
	Count(ctx context.Context, filter *domain.Filter) <-chan golibshared.Result
	Find(ctx context.Context, obj domain.Book) <-chan golibshared.Result
	Save(ctx context.Context, data *domain.Book) <-chan golibshared.Result
	Insert(ctx context.Context, newData *domain.Book) <-chan golibshared.Result
}
