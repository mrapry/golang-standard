package resthandler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brianvoe/gofakeit"
	"github.com/integralist/go-findroot/find"
	"github.com/labstack/echo"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/mock"

	"master-service/internal/modules/book/domain"
	bookUsecaseMock "master-service/internal/modules/book/usecase/mock"
	"master-service/pkg/shared"
	"testing"

	"github.com/mrapry/go-lib/golibshared"
	"github.com/mrapry/go-lib/middleware"
	authMock "github.com/mrapry/go-lib/sdk/auth-service/mocks"
	"github.com/mrapry/go-lib/validator"

	"github.com/mrapry/go-lib/config"
)

var (
	ctx                  context.Context
	restBookHandlerMocks *RestHandler
	bookUsecaseMocks     *bookUsecaseMock.BookUsecase
)

func restEchoBookHandlerMocks() {
	ctx = context.Background()
	authMocks := &authMock.ServiceAuth{}

	// set middleware
	middleware := middleware.NewMiddleware(authMocks)

	// set usecase
	bookUsecaseMocks = &bookUsecaseMock.BookUsecase{}

	// set root
	root, _ := find.Repo()

	env := config.Env{}
	env.JSONSchemaDir = fmt.Sprintf("%s/api/jsonschema", root.Path)
	config.SetEnv(env)

	// set json validator
	jsonValidator := validator.NewValidator()

	// set rest handler
	restBookHandlerMocks = &RestHandler{bookUsecase: bookUsecaseMocks, mw: middleware, validator: jsonValidator}

}

func TestNewRestHandler(t *testing.T) {
	testName := shared.SetTestcaseName(1, "new rest book handler")

	t.Run(testName, func(t *testing.T) {
		restEchoBookHandlerMocks()

		NewRestHandler(restBookHandlerMocks.bookUsecase, restBookHandlerMocks.mw, restBookHandlerMocks.validator)
	})

}

func TestRestHandler_Mount(t *testing.T) {
	testName := shared.SetTestcaseName(1, "rest book handler mount")

	t.Run(testName, func(t *testing.T) {
		restEchoBookHandlerMocks()

		// set rest handler
		restHandler := restBookHandlerMocks

		// set echo
		echoHandler := echo.New()
		groupEcho := echoHandler.Group(gofakeit.Word())

		restHandler.Mount(groupEcho)
	})
}

func TestRestHandler_findAll(t *testing.T) {
	type findAll struct {
		result []*domain.Book
		meta   *golibshared.Meta
		err    error
	}
	testCase := map[string]struct {
		findAll *findAll
		query   string
	}{
		shared.SetTestcaseName(1, "positive rest book handler get all"): {
			findAll: &findAll{
				result: []*domain.Book{},
				meta:   &golibshared.Meta{},
				err:    nil,
			},
		},
		shared.SetTestcaseName(2, "positive rest book handler get all usecase"): {
			findAll: &findAll{
				err: fmt.Errorf(golibshared.ErrorGeneral),
			},
			query: "limit=" + cast.ToString(1),
		},
		shared.SetTestcaseName(3, "positive rest book handler validate filter"): {
			query: "sort=" + gofakeit.Word(),
		},
		shared.SetTestcaseName(4, "positive rest book handler parsing filter"): {
			query: "limit=" + gofakeit.Word(),
		},
	}
	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {

			var (
				headers = map[string]string{}
			)

			// set HTTP mock
			c := golibshared.SetEchoHTTPMock(fmt.Sprintf("%s?%s", gofakeit.URL(), test.query), http.MethodGet, ``, headers)

			restEchoBookHandlerMocks()

			if test.findAll != nil {
				bookUsecaseMocks.On("FindAll", mock.Anything, mock.Anything).Return(test.findAll.result, test.findAll.meta, test.findAll.err).Once()
			}

			// set rest handler
			restHandler := restBookHandlerMocks

			// set handler
			restHandler.findAll(c)

			bookUsecaseMocks.AssertExpectations(t)
		})
	}

}

func TestRestHandler_findByID(t *testing.T) {
	type findByID struct {
		result *domain.Book
		err    error
	}

	testCase := map[string]struct {
		findByID *findByID
		ID       string
	}{
		shared.SetTestcaseName(1, "positive rest book handler get by ID"): {
			findByID: &findByID{result: &domain.Book{}},
			ID:       "5f62fcee09cd352630be5237",
		},
		shared.SetTestcaseName(2, "positive rest book handler get by ID  (ID not found)"): {
			findByID: &findByID{err: fmt.Errorf(golibshared.ErrorGeneral)},
			ID:       "5f62fcee09cd352630be5237",
		},
		shared.SetTestcaseName(3, "positive rest book handler get by ID null"): {
			ID: "",
		},
		shared.SetTestcaseName(4, "positive rest book handler get by ID invalid"): {
			ID: gofakeit.Word(),
		},
		shared.SetTestcaseName(5, "positive rest book handler get by ID (error get DB)"): {
			ID:       "5f62fcee09cd352630be5237",
			findByID: &findByID{err: fmt.Errorf(golibshared.ErrorGeneral)},
		},
	}

	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {

			var (
				headers = map[string]string{}
			)

			// set http mock
			c := golibshared.SetEchoHTTPMock(gofakeit.URL(), http.MethodGet, ``, headers)
			c.SetParamNames("id")
			c.SetParamValues(test.ID)

			restEchoBookHandlerMocks()

			if test.findByID != nil {
				bookUsecaseMocks.On("FindByID", mock.Anything, mock.Anything).Return(test.findByID.result, test.findByID.err).Once()
			}

			// set rest handler
			restHandler := restBookHandlerMocks

			// set handler
			restHandler.findByID(c)

			bookUsecaseMocks.AssertExpectations(t)
		})
	}

}

func TestRestHandler_create(t *testing.T) {
	type register struct {
		result *domain.Book
		err    error
	}

	type otherData struct {
		username string
	}

	testCase := map[string]struct {
		register  *register
		payload   *domain.Book
		otherData *otherData
	}{
		shared.SetTestcaseName(1, "positive rest book handler create new book"): {
			payload: &domain.Book{
				Name: "Payjo Suherman",
				ISBN: gofakeit.SSN(),
			},
			register: &register{result: &domain.Book{}},
		},
		shared.SetTestcaseName(2, "positive rest book handler create new book handling usecase example name is alredy exist"): {
			payload: &domain.Book{
				Name: gofakeit.Name(),
				ISBN: gofakeit.SSN(),
			},
			register: &register{err: fmt.Errorf(golibshared.ErrorGeneral)},
		},
		shared.SetTestcaseName(3, "negative rest book handler create new book validate body"): {
			payload: &domain.Book{},
		},
		shared.SetTestcaseName(4, "negative rest book handler create new book unmarshal"): {
			payload: &domain.Book{
				Name: gofakeit.Name(),
				ISBN: gofakeit.SSN(),
			},
			otherData: &otherData{
				username: "matmat",
			},
		},
		shared.SetTestcaseName(5, "negative rest book handler create new book validate body"): {
			payload: &domain.Book{Name: "P", ISBN: gofakeit.SSN()},
		},
	}
	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {
			var (
				headers = map[string]string{}
			)

			headers[echo.HeaderContentType] = echo.MIMEApplicationJSON

			// set payload
			var payload string
			if test.payload != nil {
				payload = shared.CreateHttpRequestBodyMock(test.payload)
			}
			if test.otherData != nil {
				payload = shared.CreateHttpRequestBodyMock(test.otherData)
			}

			// set http mock
			c := golibshared.SetEchoHTTPMock(gofakeit.URL(), http.MethodPost, payload, headers)

			restEchoBookHandlerMocks()

			if test.register != nil {
				bookUsecaseMocks.On("Create", mock.Anything, mock.Anything).Return(test.register.result, test.register.err).Once()
			}

			// set rest handler
			restHandler := restBookHandlerMocks

			// set handler
			restHandler.create(c)

			bookUsecaseMocks.AssertExpectations(t)
		})
	}

}

func TestRestHandler_update(t *testing.T) {
	type update struct {
		result *domain.Book
		err    error
	}

	type otherData struct {
		username string
	}

	testCase := map[string]struct {
		update    *update
		payload   *domain.Book
		ID        string
		otherData *otherData
	}{
		shared.SetTestcaseName(1, "positive rest book handler update book"): {
			payload: &domain.Book{
				Name: gofakeit.Name(),
				ISBN: gofakeit.SSN(),
			},
			ID:     "5f62fcee09cd352630be5237",
			update: &update{result: &domain.Book{}},
		},
		shared.SetTestcaseName(2, "positive rest book handler update book usecase handling"): {
			payload: &domain.Book{
				Name: gofakeit.Name(),
				ISBN: gofakeit.SSN(),
			},
			ID:     "5f62fcee09cd352630be5237",
			update: &update{err: fmt.Errorf(golibshared.ErrorGeneral)},
		},
		shared.SetTestcaseName(3, "positive rest book handler update book error ID Required"): {
			payload: &domain.Book{
				Name: gofakeit.Name(),
				ISBN: gofakeit.SSN(),
			},
		},
		shared.SetTestcaseName(4, "positive rest book handler update book ID invalid"): {
			payload: &domain.Book{
				Name: gofakeit.Name(),
				ISBN: gofakeit.SSN(),
			},
			ID: gofakeit.Name(),
		},

		shared.SetTestcaseName(5, "positive rest book handler update book ID atribut minlength"): {
			payload: &domain.Book{
				Name: "P",
			},
		},
		shared.SetTestcaseName(5, "positive rest book handler update book ID unmarshal"): {
			payload: &domain.Book{
				Name: gofakeit.Name(),
				ISBN: gofakeit.SSN(),
			},
			otherData: &otherData{
				username: gofakeit.Name(),
			},
		},
	}
	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {
			var (
				headers = map[string]string{}
			)

			headers[echo.HeaderContentType] = echo.MIMEApplicationJSON

			// set payload
			var payload string
			if test.payload != nil {
				payload = shared.CreateHttpRequestBodyMock(test.payload)
			}
			if test.otherData != nil {
				payload = shared.CreateHttpRequestBodyMock(test.otherData)
			}

			// set http mock
			c := golibshared.SetEchoHTTPMock(gofakeit.URL(), http.MethodPost, payload, headers)
			c.SetParamNames("id")
			c.SetParamValues(test.ID)

			restEchoBookHandlerMocks()

			if test.update != nil {
				bookUsecaseMocks.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(test.update.result, test.update.err).Once()
			}

			// set rest handler
			restHandler := restBookHandlerMocks

			// set handler
			restHandler.update(c)

			bookUsecaseMocks.AssertExpectations(t)
		})
	}
}

func TestRestHandler_delete(t *testing.T) {
	type remove struct {
		result *domain.Book
		err    error
	}
	testCase := map[string]struct {
		remove *remove
		ID     string
	}{
		shared.SetTestcaseName(1, "positive rest book handler remove by ID"): {
			remove: &remove{err: nil},
			ID:     "5f62fcee09cd352630be5237",
		},
		shared.SetTestcaseName(2, "positive rest book handler remove by ID  (ID not found)"): {
			remove: &remove{err: fmt.Errorf(golibshared.ErrorGeneral)},
			ID:     "5f62fcee09cd352630be5237",
		},
		shared.SetTestcaseName(3, "positive rest book handler remove by ID null"): {
			ID: "",
		},
		shared.SetTestcaseName(4, "positive rest book handler remove by ID invalid"): {
			ID:     gofakeit.Word(),
			remove: &remove{err: fmt.Errorf(golibshared.ErrorGeneral)},
		},
	}
	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {

			var (
				headers = map[string]string{}
			)

			// set http mock
			c := golibshared.SetEchoHTTPMock(gofakeit.URL(), http.MethodGet, ``, headers)
			c.SetParamNames("id")
			c.SetParamValues(test.ID)

			restEchoBookHandlerMocks()

			if test.remove != nil {
				bookUsecaseMocks.On("RemoveByID", mock.Anything, mock.Anything).Return(test.remove.result, test.remove.err).Once()
			}

			// set rest handler
			restHandler := restBookHandlerMocks

			// set handler
			restHandler.delete(c)

			bookUsecaseMocks.AssertExpectations(t)
		})
	}
}

func TestRestHandler_restore(t *testing.T) {
	type restore struct {
		result *domain.Book
		err    error
	}
	testCase := map[string]struct {
		restore *restore
		ID      string
	}{
		shared.SetTestcaseName(1, "positive rest book handler restore by ID"): {
			restore: &restore{err: nil},
			ID:      "5f62fcee09cd352630be5237",
		},
		shared.SetTestcaseName(2, "positive rest book handler restore by ID  (ID not found)"): {
			restore: &restore{err: fmt.Errorf(golibshared.ErrorGeneral)},
			ID:      "5f62fcee09cd352630be5237",
		},
		shared.SetTestcaseName(3, "positive rest book handler restore by ID null"): {
			ID: "",
		},
		shared.SetTestcaseName(4, "positive rest book handler restore by ID invalid"): {
			ID:      gofakeit.Word(),
			restore: &restore{err: fmt.Errorf(golibshared.ErrorGeneral)},
		},
	}
	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {

			var (
				headers = map[string]string{}
			)

			// set http mock
			c := golibshared.SetEchoHTTPMock(gofakeit.URL(), http.MethodGet, ``, headers)
			c.SetParamNames("id")
			c.SetParamValues(test.ID)

			restEchoBookHandlerMocks()

			if test.restore != nil {
				bookUsecaseMocks.On("RestoreByID", mock.Anything, mock.Anything).Return(test.restore.result, test.restore.err).Once()
			}

			// set rest handler
			restHandler := restBookHandlerMocks

			// set handler
			restHandler.restore(c)

			bookUsecaseMocks.AssertExpectations(t)
		})
	}
}
