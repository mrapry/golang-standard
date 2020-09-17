package usecase

import (
	"context"
	"fmt"
	"master-service/internal/modules/book/domain"
	"master-service/internal/modules/book/repository"
	"master-service/pkg/shared"

	"github.com/mrapry/go-lib/codebase/interfaces"
	"github.com/mrapry/go-lib/golibshared"
	"github.com/mrapry/go-lib/logger"
	"github.com/mrapry/go-lib/sdk"
	"github.com/mrapry/go-lib/tracer"
	"github.com/spf13/cast"
	"go.uber.org/zap/zapcore"
)

//bookUsecaseImpl structure
type bookUsecaseImpl struct {
	repo      *repository.Repository
	sdk       sdk.SDK
	validator interfaces.Validator
}

// NewBookUsecase create new member usecase
func NewBookUsecase(repo *repository.Repository, sdk sdk.SDK, validator interfaces.Validator) BookUsecase {
	return &bookUsecaseImpl{repo: repo, sdk: sdk, validator: validator}
}

func (uc *bookUsecaseImpl) FindAll(ctx context.Context, filter *domain.Filter) (listBook []*domain.Book, meta *golibshared.Meta, err error) {
	opName := "book_usecase.find_all"

	tracer := tracer.StartTrace(ctx, opName)
	defer tracer.Finish(nil)
	ctx = tracer.Context()

	//get data book from repository
	repoRes := <-uc.repo.Book.FindAll(ctx, filter)
	if repoRes.Error != nil {
		logger.Log(zapcore.ErrorLevel, repoRes.Error.Error(), opName, "find_book")
		return nil, nil, repoRes.Error
	}

	// transform data to struct
	book := repoRes.Data.([]*domain.Book)

	// count member
	countRes := <-uc.repo.Book.Count(ctx, filter)
	if countRes.Error != nil {
		logger.Log(zapcore.ErrorLevel, countRes.Error.Error(), opName, "count_book")
		return nil, nil, countRes.Error
	}

	// transform data to struct
	total := countRes.Data.(int64)

	// set meta
	meta = golibshared.NewMeta(cast.ToInt64(filter.Page), cast.ToInt64(filter.Limit), total)

	return book, meta, nil
}

func (uc *bookUsecaseImpl) FindByID(ctx context.Context, ID string) (*domain.Book, error) {
	opName := "book_usecase.find_by_id"
	tracer := tracer.StartTrace(ctx, opName)
	defer tracer.Finish()

	var (
		err error
	)

	// get book
	repoRes := <-uc.repo.Book.FindByID(ctx, ID)
	if repoRes.Error != nil {
		err = fmt.Errorf(shared.ErrorDataNotFound, "book")
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "find_book")
		return nil, err
	}

	// transform data to struct
	book := repoRes.Data.(*domain.Book)

	return book, nil
}

func (uc *bookUsecaseImpl) Create(ctx context.Context, data *domain.Book) (*domain.Book, error) {
	opName := "book_usecase.create"
	tracer := tracer.StartTrace(ctx, opName)
	defer tracer.Finish(nil)
	ctx = tracer.Context()

	var (
		err error
	)

	//set is active true for new book
	data.IsActive = true

	// get book by name
	repoResCheckName := <-uc.repo.Book.Find(ctx, domain.Book{Name: data.Name})
	if repoResCheckName.Error == nil {
		err = fmt.Errorf(shared.ErrorDataAlreadyExist, "book")
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "find_book")
		return nil, err
	}

	//save data to db
	repoRes := <-uc.repo.Book.Insert(ctx, data)
	if repoRes.Error != nil {
		logger.Log(zapcore.ErrorLevel, repoRes.Error.Error(), opName, "insert_book")
		return nil, repoRes.Error
	}

	// transform data to struct
	book := repoRes.Data.(*domain.Book)

	return book, nil
}

func (uc *bookUsecaseImpl) Update(ctx context.Context, data *domain.Book, ID string) (*domain.Book, error) {
	opName := "book_usecase.update"
	tracer := tracer.StartTrace(ctx, opName)
	defer tracer.Finish(nil)
	ctx = tracer.Context()

	var (
		err error
	)

	// get book by id
	repoResCheckID := <-uc.repo.Book.FindByID(ctx, ID)
	if repoResCheckID.Error != nil {
		err = fmt.Errorf(shared.ErrorDataNotFound, "book")
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "find_book")
		return nil, err
	}

	book := repoResCheckID.Data.(*domain.Book)
	book.Name = data.Name

	//update data to db
	repoRes := <-uc.repo.Book.Save(ctx, book)
	if repoRes.Error != nil {
		logger.Log(zapcore.ErrorLevel, repoRes.Error.Error(), opName, "update_book")
		return nil, repoRes.Error
	}

	// transform data to struct
	bookUpdate := repoRes.Data.(*domain.Book)
	return bookUpdate, nil
}

func (uc *bookUsecaseImpl) RemoveByID(ctx context.Context, ID string) error {
	opName := "book_usecase.remove_by_id"

	var (
		err error
	)

	// get book
	repoResCheckID := <-uc.repo.Book.FindByID(ctx, ID)
	if repoResCheckID.Error != nil {
		err = fmt.Errorf(shared.ErrorDataNotFound, "book")
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "find_book")
		return err
	}

	// transform data to struct
	book := repoResCheckID.Data.(*domain.Book)
	if !book.IsActive {
		err = fmt.Errorf(shared.ErrorDataNotActived, "book")
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "find_book")
		return err
	}

	book.IsActive = false

	//update data to db
	repoRes := <-uc.repo.Book.Save(ctx, book)
	if repoRes.Error != nil {
		logger.Log(zapcore.ErrorLevel, repoRes.Error.Error(), opName, "update_book")
		return repoRes.Error
	}

	return nil
}

func (uc *bookUsecaseImpl) RestoreByID(ctx context.Context, ID string) error {
	opName := "book_usecase.restore_by_id"
	var (
		err error
	)

	// get book
	repoResCheckID := <-uc.repo.Book.FindByID(ctx, ID)
	if repoResCheckID.Error != nil {
		err = fmt.Errorf(shared.ErrorDataNotFound, "book")
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "find_book")
		return err
	}

	// transform data to struct
	book := repoResCheckID.Data.(*domain.Book)
	if book.IsActive {
		err = fmt.Errorf(shared.ErrorDataIsActived, "book")
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "find_book")
		return err
	}

	book.IsActive = true

	//update data to db
	repoRes := <-uc.repo.Book.Save(ctx, book)
	if repoRes.Error != nil {
		logger.Log(zapcore.ErrorLevel, repoRes.Error.Error(), opName, "update_book")
		return repoRes.Error
	}

	return nil
}
