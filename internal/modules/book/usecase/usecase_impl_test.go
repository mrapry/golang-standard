package usecase

import (
	"context"
	"fmt"
	"master-service/internal/modules/book/domain"
	"master-service/internal/modules/book/repository"
	bookRepoMock "master-service/internal/modules/book/repository/interfaces/mock"
	pkgMock "master-service/pkg/mock/mocks"
	"master-service/pkg/shared"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"

	authMock "github.com/mrapry/go-lib/sdk/auth-service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mrapry/go-lib/golibshared"
	"github.com/mrapry/go-lib/sdk"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ctx              context.Context
	bookUsecaseMocks BookUsecase
	bookRepoMocks    *bookRepoMock.BookRepository
	validatorMocks   *pkgMock.Validator
	storageMocks     *pkgMock.Storage
	authMocks        *authMock.ServiceAuth
)

func bookUsecaseMock() {
	//set context
	ctx = context.Background()

	// set repo
	mongodb := &mongo.Database{}
	repo := repository.NewRepository(mongodb, mongodb)
	bookRepoMocks = &bookRepoMock.BookRepository{}
	repo.Book = bookRepoMocks

	// set service
	authMocks = &authMock.ServiceAuth{}
	sdk := sdk.NewSDK(
		sdk.SetAuthService(authMocks),
	)

	//set validator mock
	validatorMocks = &pkgMock.Validator{}

	storageMocks = &pkgMock.Storage{}

	// set usecase
	bookUsecaseMocks = NewBookUsecase(repo, sdk, validatorMocks, storageMocks)

}

func TestNewBookUsecase(t *testing.T) {
	testName := shared.SetTestcaseName(1, "new book usecase")

	t.Run(testName, func(t *testing.T) {
		bookUsecaseMock()

		// set usecase
		usecase := bookUsecaseMocks

		assert.NotNil(t, usecase)
	})

}

func Test_bookUsecaseImpl_FindAll(t *testing.T) {
	testCase := map[string]struct {
		wantError bool
		findAll   *golibshared.Result
		count     *golibshared.Result
	}{
		shared.SetTestcaseName(1, "positive find all book"): {
			wantError: false,
			findAll: &golibshared.Result{Data: []*domain.Book{
				&domain.Book{},
			}},
			count: &golibshared.Result{Data: int64(1)},
		},
		shared.SetTestcaseName(2, "negative find all count book"): {
			wantError: true,
			findAll: &golibshared.Result{Data: []*domain.Book{
				&domain.Book{},
			}},
			count: &golibshared.Result{Error: fmt.Errorf(golibshared.ErrorGeneral)},
		},
		shared.SetTestcaseName(3, "negative find all find book"): {
			wantError: true,
			findAll:   &golibshared.Result{Error: fmt.Errorf(golibshared.ErrorGeneral)},
		},
	}
	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {
			bookUsecaseMock()
			if test.findAll != nil {
				result := shared.SetMockerySharedResult(*test.findAll)
				bookRepoMocks.On("FindAll", mock.Anything, mock.Anything).Return(result).Once()
			}

			if test.count != nil {
				result := shared.SetMockerySharedResult(*test.count)
				bookRepoMocks.On("Count", mock.Anything, mock.Anything).Return(result).Once()
			}

			// set usecase
			usecase := bookUsecaseMocks

			// run the usecase
			_, _, err := usecase.FindAll(ctx, &domain.Filter{})
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			bookRepoMocks.AssertExpectations(t)
		})
	}
}

func Test_bookUsecaseImpl_FindByID(t *testing.T) {
	type redis struct {
		book  *domain.Book
		Error error
	}
	testCase := map[string]struct {
		wantError     bool
		ID            string
		redis, redis2 *redis
		find          *golibshared.Result
	}{
		shared.SetTestcaseName(1, "positive find by id get from redis"): {
			wantError: false,
			ID:        "5f62fcee09cd352630be5237",
			redis:     &redis{book: &domain.Book{}},
		},
		shared.SetTestcaseName(2, "positive find by id get from database and save redis"): {
			wantError: false,
			ID:        "5f62fcee09cd352630be5237",
			redis:     &redis{Error: fmt.Errorf(golibshared.ErrorGeneral)},
			find:      &golibshared.Result{Data: &domain.Book{}},
			redis2:    &redis{book: &domain.Book{}},
		},

		shared.SetTestcaseName(3, "negative find by id"): {
			wantError: true,
			ID:        "12345",
			redis:     &redis{Error: fmt.Errorf(golibshared.ErrorGeneral)},
			find:      &golibshared.Result{Error: fmt.Errorf(golibshared.ErrorGeneral)},
		},
		shared.SetTestcaseName(3, "negative find by id"): {
			wantError: true,
			ID:        "",
			redis:     &redis{Error: fmt.Errorf(golibshared.ErrorGeneral)},
			find:      &golibshared.Result{Error: fmt.Errorf(golibshared.ErrorGeneral)},
		},
	}
	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {
			bookUsecaseMock()

			if test.redis != nil {
				var data string = mock.Anything
				var err error = test.redis.Error
				storageMocks.On("Get", mock.Anything, mock.Anything).Return(data, err).Once()
			}

			if test.find != nil {
				result := shared.SetMockerySharedResult(*test.find)
				bookRepoMocks.On("FindByID", mock.Anything, mock.Anything).Return(result).Once()
			}

			if test.redis2 != nil {
				var data string = mock.Anything
				var err error = test.redis2.Error
				storageMocks.On("Set", mock.Anything, mock.Anything, mock.Anything, 100*time.Minute).Return(data, err).Once()
			}

			// set usecase
			usecase := bookUsecaseMocks

			// run the usecase
			_, err := usecase.FindByID(ctx, test.ID)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			bookRepoMocks.AssertExpectations(t)
		})
	}

}

func Test_bookUsecaseImpl_Create(t *testing.T) {
	testCase := map[string]struct {
		wantError   bool
		dataUsecase *domain.Book
		find        *golibshared.Result
		insert      *golibshared.Result
	}{
		shared.SetTestcaseName(1, "negative name already taken"): {
			wantError: true,
			dataUsecase: &domain.Book{
				Name: gofakeit.Name(),
				ISBN: gofakeit.SSN(),
			},
			find: &golibshared.Result{Data: &domain.Book{}},
		},
		shared.SetTestcaseName(2, "positive name not already taken"): {
			wantError: false,
			dataUsecase: &domain.Book{
				Name: gofakeit.Name(),
				ISBN: gofakeit.SSN(),
			},
			find:   &golibshared.Result{Error: fmt.Errorf(golibshared.ErrorGeneral)},
			insert: &golibshared.Result{Data: &domain.Book{}},
		},
		shared.SetTestcaseName(3, "data book is nil"): {
			wantError:   true,
			dataUsecase: &domain.Book{},
			find:        &golibshared.Result{Error: fmt.Errorf(golibshared.ErrorGeneral)},
			insert:      &golibshared.Result{Error: fmt.Errorf(golibshared.ErrorGeneral)},
		},
	}
	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {
			bookUsecaseMock()

			if test.find != nil {
				result := shared.SetMockerySharedResult(*test.find)
				bookRepoMocks.On("Find", mock.Anything, mock.Anything).Return(result).Once()
			}

			if test.insert != nil {
				result := shared.SetMockerySharedResult(*test.insert)
				bookRepoMocks.On("Insert", mock.Anything, mock.Anything).Return(result).Once()
			}

			// set usecase
			usecase := bookUsecaseMocks

			// run the usecase
			_, err := usecase.Create(ctx, test.dataUsecase)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			bookRepoMocks.AssertExpectations(t)
		})
	}

}

func Test_bookUsecaseImpl_Update(t *testing.T) {
	type redis struct {
		book  *domain.Book
		Error error
	}
	testCase := map[string]struct {
		wantError   bool
		dataUsecase *domain.Book
		ID          string
		findByID    *golibshared.Result
		update      *golibshared.Result
		redis       *redis
	}{
		shared.SetTestcaseName(1, "negative ID is not found"): {
			wantError: true,
			dataUsecase: &domain.Book{
				Name: gofakeit.Name(),
				ISBN: gofakeit.SSN(),
			},
			findByID: &golibshared.Result{Error: fmt.Errorf(golibshared.ErrorGeneral)},
		},
		shared.SetTestcaseName(2, "positive ID found"): {
			wantError: false,
			dataUsecase: &domain.Book{
				Name: gofakeit.Name(),
				ISBN: gofakeit.SSN(),
			},
			ID:       "5f62fcee09cd352630be5237",
			findByID: &golibshared.Result{Data: &domain.Book{}},
			update:   &golibshared.Result{Data: &domain.Book{}},
			redis:    &redis{book: &domain.Book{}},
		},
		shared.SetTestcaseName(3, "data book is null"): {
			wantError:   true,
			dataUsecase: &domain.Book{},
			ID:          "5f62fcee09cd352630be5237",
			findByID:    &golibshared.Result{Data: &domain.Book{}},
			update:      &golibshared.Result{Error: fmt.Errorf(golibshared.ErrorGeneral)},
		},
	}
	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {
			bookUsecaseMock()

			if test.findByID != nil {
				result := shared.SetMockerySharedResult(*test.findByID)
				bookRepoMocks.On("FindByID", mock.Anything, mock.Anything).Return(result).Once()
			}

			if test.update != nil {
				result := shared.SetMockerySharedResult(*test.update)
				bookRepoMocks.On("Save", mock.Anything, mock.Anything).Return(result).Once()
			}

			if test.redis != nil {
				var data string = mock.Anything
				var err error = test.redis.Error
				storageMocks.On("Delete", mock.Anything, mock.Anything).Return(data, err).Once()
			}

			// set usecase
			usecase := bookUsecaseMocks

			// run the usecase
			_, err := usecase.Update(ctx, test.dataUsecase, test.ID)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			bookRepoMocks.AssertExpectations(t)
		})
	}
}

func Test_bookUsecaseImpl_RemoveByID(t *testing.T) {
	type redis struct {
		book  *domain.Book
		Error error
	}
	testCase := map[string]struct {
		wantError bool
		ID        string
		findByID  *golibshared.Result
		remove    *golibshared.Result
		redis     *redis
	}{
		shared.SetTestcaseName(0, "Positive case, ID ditemukan dan data tidak dalam keadaan non active"): {
			wantError: false,
			ID:        "5f62fcee09cd352630be5237",
			findByID: &golibshared.Result{Data: &domain.Book{
				Name:     gofakeit.Name(),
				ISBN:     gofakeit.SSN(),
				IsActive: true,
			}},
			remove: &golibshared.Result{Data: &domain.Book{}},
			redis:  &redis{book: &domain.Book{}},
		},
		shared.SetTestcaseName(1, "Data sudah dalam keadaan non active"): {
			wantError: true,
			ID:        "5f62fcee09cd352630be5237",
			findByID: &golibshared.Result{Data: &domain.Book{
				Name:     gofakeit.Name(),
				ISBN:     gofakeit.SSN(),
				IsActive: false,
			}},
		},
		shared.SetTestcaseName(2, "ID tidak ditemukan"): {
			wantError: true,
			ID:        "5f62fcee09cd352630be5237",
			findByID:  &golibshared.Result{Error: fmt.Errorf(golibshared.ErrorGeneral)},
		},
		shared.SetTestcaseName(3, "Invalid format ID"): {
			wantError: true,
			ID:        "01w313!!!!!",
			findByID:  &golibshared.Result{Error: fmt.Errorf(golibshared.ErrorGeneral)},
		},
		shared.SetTestcaseName(4, "Error Save data"): {
			wantError: true,
			ID:        "5f62fcee09cd352630be5237",
			findByID: &golibshared.Result{Data: &domain.Book{
				Name:     gofakeit.Name(),
				ISBN:     gofakeit.SSN(),
				IsActive: true,
			}},
			remove: &golibshared.Result{Error: fmt.Errorf(golibshared.ErrorGeneral)},
		},
	}
	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {
			bookUsecaseMock()

			if test.findByID != nil {
				result := shared.SetMockerySharedResult(*test.findByID)
				bookRepoMocks.On("FindByID", mock.Anything, mock.Anything).Return(result).Once()
			}

			if test.remove != nil {
				result := shared.SetMockerySharedResult(*test.remove)
				bookRepoMocks.On("Save", mock.Anything, mock.Anything).Return(result).Once()
			}

			if test.redis != nil {
				var data string = mock.Anything
				var err error = test.redis.Error
				storageMocks.On("Delete", mock.Anything, mock.Anything).Return(data, err).Once()
			}

			// set usecase
			usecase := bookUsecaseMocks

			// run the usecase
			err := usecase.RemoveByID(ctx, test.ID)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			bookRepoMocks.AssertExpectations(t)
		})
	}
}

func Test_bookUsecaseImpl_RestoreByID(t *testing.T) {
	type redis struct {
		book  *domain.Book
		Error error
	}
	testCase := map[string]struct {
		wantError bool
		ID        string
		findByID  *golibshared.Result
		restore   *golibshared.Result
		redis     *redis
	}{
		shared.SetTestcaseName(0, "Positive case, ID ditemukan dan data tidak dalam keadaan non active"): {
			wantError: false,
			ID:        "5f62fcee09cd352630be5237",
			findByID: &golibshared.Result{Data: &domain.Book{
				Name:     gofakeit.Name(),
				ISBN:     gofakeit.SSN(),
				IsActive: false,
			}},
			restore: &golibshared.Result{Data: &domain.Book{}},
			redis:   &redis{book: &domain.Book{}},
		},
		shared.SetTestcaseName(1, "Data sudah dalam keadaan non active"): {
			wantError: true,
			ID:        "5f62fcee09cd352630be5237",
			findByID: &golibshared.Result{Data: &domain.Book{
				Name:     gofakeit.Name(),
				ISBN:     gofakeit.SSN(),
				IsActive: true,
			}},
		},
		shared.SetTestcaseName(2, "ID tidak ditemukan"): {
			wantError: true,
			ID:        "5f62fcee09cd352630be5237",
			findByID:  &golibshared.Result{Error: fmt.Errorf(golibshared.ErrorGeneral)},
		},
		shared.SetTestcaseName(3, "Invalid format ID"): {
			wantError: true,
			ID:        "01w313!!!!!",
			findByID:  &golibshared.Result{Error: fmt.Errorf(golibshared.ErrorGeneral)},
		},
		shared.SetTestcaseName(4, "Error Save data"): {
			wantError: true,
			ID:        "5f62fcee09cd352630be5237",
			findByID: &golibshared.Result{Data: &domain.Book{
				Name:     gofakeit.Name(),
				ISBN:     gofakeit.SSN(),
				IsActive: false,
			}},
			restore: &golibshared.Result{Error: fmt.Errorf(golibshared.ErrorGeneral)},
		},
	}
	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {
			bookUsecaseMock()

			if test.findByID != nil {
				result := shared.SetMockerySharedResult(*test.findByID)
				bookRepoMocks.On("FindByID", mock.Anything, mock.Anything).Return(result).Once()
			}

			if test.restore != nil {
				result := shared.SetMockerySharedResult(*test.restore)
				bookRepoMocks.On("Save", mock.Anything, mock.Anything).Return(result).Once()
			}

			if test.redis != nil {
				var data string = mock.Anything
				var err error = test.redis.Error
				storageMocks.On("Delete", mock.Anything, mock.Anything).Return(data, err).Once()
			}

			// set usecase
			usecase := bookUsecaseMocks

			// run the usecase
			err := usecase.RestoreByID(ctx, test.ID)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			bookRepoMocks.AssertExpectations(t)
		})
	}
}
