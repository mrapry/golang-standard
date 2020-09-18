package main

//DeliveryRestTemplate constant template from deliveryResttemplate.
//{{$.GoModules}} = go module
//{{$.module}} = name of module (lowercase)
//{{clean (upper $.module)}} = name of module (capital)
const DeliveryRestTemplate = `package resthandler

import (
	"encoding/json"
	"io/ioutil"
	"{{$.GoModules}}/internal/modules/{{$.module}}/domain"
	"{{$.GoModules}}/internal/modules/{{$.module}}/usecase"
	"net/http"

	"github.com/labstack/echo"
	"{{$.LibraryAddress}}/codebase/interfaces"
	"{{$.LibraryAddress}}/golibhelper"
	helper "{{$.LibraryAddress}}/golibhelper"
	"{{$.LibraryAddress}}/logger"
	"{{$.LibraryAddress}}/tracer"
	"{{$.LibraryAddress}}/wrapper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap/zapcore"
)

// RestHandler handler
type RestHandler struct {
	{{$.module}}Usecase usecase.{{clean (upper $.module)}}Usecase
	mw          interfaces.Middleware
	validator   interfaces.Validator
}

// NewRestHandler create new rest handler
func NewRestHandler({{$.module}}Usecase usecase.{{clean (upper $.module)}}Usecase, mw interfaces.Middleware, validator interfaces.Validator) *RestHandler {
	return &RestHandler{
		{{$.module}}Usecase: {{$.module}}Usecase,
		mw:          mw,
		validator:   validator,
	}
}

// Mount handler with root "/"
// handling version in here
func (h *RestHandler) Mount(root *echo.Group) {
	v1Root := root.Group(helper.V1)

	{{$.module}} := v1Root.Group("/{{$.module}}", h.mw.HTTPBasicAuth(false))
	{{$.module}}.GET("", h.findAll)
	{{$.module}}.GET("/:id", h.findByID)
	{{$.module}}.POST("", h.create)
	{{$.module}}.PUT("/:id", h.update)
	{{$.module}}.DELETE("/:id", h.delete)
	{{$.module}}.PATCH("/:id", h.restore)
}

func (h *RestHandler) findAll(c echo.Context) error {
	opName := "{{$.module}}_resthandler.find_all"
	ctx := c.Request().Context()
	tracer := tracer.StartTrace(ctx, opName)
	defer tracer.Finish()
	ctx = tracer.Context()

	var filter domain.Filter
	if err := golibhelper.ParseFromQueryParam(c.Request().URL.Query(), &filter); err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "parse_query")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "Gagal memparsing filter", err).JSON(c.Response())
	}

	body, _ := json.Marshal(filter)
	if err := h.validator.ValidateDocument("{{$.module}}/get_all", body); err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "validate_payload")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "Gagal memvalidasi filter", err.Error()).JSON(c.Response())
	}

	result, meta, err := h.{{$.module}}Usecase.FindAll(ctx, &filter)
	if err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "usecase_error")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c.Response())
	}
	return wrapper.NewHTTPResponse(http.StatusOK, "Sukses mengambil data {{$.module}}", result, meta).JSON(c.Response())
}

func (h *RestHandler) findByID(c echo.Context) error {
	opName := "{{$.module}}_resthandler.find_by_id"
	ctx := c.Request().Context()

	tracer := tracer.StartTrace(ctx, opName)
	defer tracer.Finish()
	ctx = tracer.Context()

	id := c.Param("id")
	if id == "" {
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "id tidak boleh kosong").JSON(c.Response())
	}

	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "invalid id").JSON(c.Response())
	}

	result, err := h.{{$.module}}Usecase.FindByID(ctx, id)
	if err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "usecase_error")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "Gagal mendapatkan data {{$.module}}", err.Error()).JSON(c.Response())
	}

	return wrapper.NewHTTPResponse(http.StatusOK, "Sukses", result).JSON(c.Response())
}

func (h *RestHandler) create(c echo.Context) error {
	opName := "{{$.module}}_resthandler.create"
	ctx := c.Request().Context()

	tracer := tracer.StartTrace(ctx, opName)
	defer tracer.Finish()
	ctx = tracer.Context()

	body, _ := ioutil.ReadAll(c.Request().Body)
	if err := h.validator.ValidateDocument("{{$.module}}/create", body); err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "validate_payload")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "Gagal dalam validasi data", err.Error()).JSON(c.Response())
	}

	var payload domain.{{clean (upper $.module)}}
	if err := json.Unmarshal(body, &payload); err != nil {
		return wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c.Response())
	}

	result, err := h.{{$.module}}Usecase.Create(ctx, &payload)
	if err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "usecase_error")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "Gagal menyimpan data {{$.module}} baru", err.Error()).JSON(c.Response())
	}

	return wrapper.NewHTTPResponse(http.StatusCreated, "Sukses menyimpan data user", result).JSON(c.Response())
}

func (h *RestHandler) update(c echo.Context) error {
	opName := "{{$.module}}_resthandler.update"
	ctx := c.Request().Context()

	tracer := tracer.StartTrace(ctx, opName)
	defer tracer.Finish()
	ctx = tracer.Context()

	body, _ := ioutil.ReadAll(c.Request().Body)
	if err := h.validator.ValidateDocument("{{$.module}}/update", body); err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "validate_payload")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "Gagal dalam validasi data", err.Error()).JSON(c.Response())
	}

	var payload domain.{{clean (upper $.module)}}
	if err := json.Unmarshal(body, &payload); err != nil {
		return wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c.Response())
	}

	ID := c.Param("id")
	if ID == "" {
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "id tidak boleh kosong").JSON(c.Response())
	}
	_, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "invalid id").JSON(c.Response())
	}

	result, err := h.{{$.module}}Usecase.Update(ctx, &payload, ID)
	if err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "usecase_error")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "Gagal update data {{$.module}} baru", err.Error()).JSON(c.Response())
	}

	return wrapper.NewHTTPResponse(http.StatusCreated, "Sukses update data user", result).JSON(c.Response())
}

func (h *RestHandler) delete(c echo.Context) error {
	opName := "{{$.module}}_resthandler.delete"
	ctx := c.Request().Context()

	tracer := tracer.StartTrace(ctx, opName)
	defer tracer.Finish()
	ctx = tracer.Context()

	id := c.Param("id")
	if id == "" {
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "id tidak boleh kosong").JSON(c.Response())
	}

	err := h.{{$.module}}Usecase.RemoveByID(ctx, id)
	if err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "usecase_error")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "Gagal delete data {{$.module}}", err.Error()).JSON(c.Response())
	}

	return wrapper.NewHTTPResponse(http.StatusOK, "Sukses delete data").JSON(c.Response())
}

func (h *RestHandler) restore(c echo.Context) error {
	opName := "{{$.module}}_resthandler.restore"
	ctx := c.Request().Context()

	tracer := tracer.StartTrace(ctx, opName)
	defer tracer.Finish()
	ctx = tracer.Context()

	id := c.Param("id")
	if id == "" {
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "id tidak boleh kosong").JSON(c.Response())
	}

	err := h.{{$.module}}Usecase.RestoreByID(ctx, id)
	if err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "usecase_error")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "Gagal merestore data {{$.module}}", err.Error()).JSON(c.Response())
	}

	return wrapper.NewHTTPResponse(http.StatusOK, "Sukses restore data").JSON(c.Response())
}
`

//DeliveryRestTestTemplate constant template test from deliveryResttemplate.
//{{$.GoModules}} = go module
//{{$.module}} = name of module (lowercase)
//{{clean (upper $.module)}} = name of module (capital)
const DeliveryRestTestTemplate = `package resthandler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brianvoe/gofakeit"
	"github.com/integralist/go-findroot/find"
	"github.com/labstack/echo"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/mock"

	"{{$.GoModules}}/internal/modules/{{$.module}}/domain"
	{{$.module}}UsecaseMock "{{$.GoModules}}/internal/modules/{{$.module}}/usecase/mock"
	"{{$.GoModules}}/pkg/shared"
	"testing"

	"{{$.LibraryAddress}}/golibshared"
	"{{$.LibraryAddress}}/middleware"
	authMock "{{$.LibraryAddress}}/sdk/auth-service/mocks"
	"{{$.LibraryAddress}}/validator"

	"{{$.LibraryAddress}}/config"
)

var (
	ctx                  context.Context
	rest{{clean (upper $.module)}}HandlerMocks *RestHandler
	{{$.module}}UsecaseMocks     *{{$.module}}UsecaseMock.{{clean (upper $.module)}}Usecase
)

func restEcho{{clean (upper $.module)}}HandlerMocks() {
	ctx = context.Background()
	authMocks := &authMock.ServiceAuth{}

	// set middleware
	middleware := middleware.NewMiddleware(authMocks)

	// set usecase
	{{$.module}}UsecaseMocks = &{{$.module}}UsecaseMock.{{clean (upper $.module)}}Usecase{}

	// set root
	root, _ := find.Repo()

	env := config.Env{}
	env.JSONSchemaDir = fmt.Sprintf("%s/api/jsonschema", root.Path)
	config.SetEnv(env)

	// set json validator
	jsonValidator := validator.NewValidator()

	// set rest handler
	rest{{clean (upper $.module)}}HandlerMocks = &RestHandler{ {{$.module}}Usecase: {{$.module}}UsecaseMocks, mw: middleware, validator: jsonValidator}

}

func TestNewRestHandler(t *testing.T) {
	testName := shared.SetTestcaseName(1, "new rest {{$.module}} handler")

	t.Run(testName, func(t *testing.T) {
		restEcho{{clean (upper $.module)}}HandlerMocks()

		NewRestHandler(rest{{clean (upper $.module)}}HandlerMocks.{{$.module}}Usecase, rest{{clean (upper $.module)}}HandlerMocks.mw, rest{{clean (upper $.module)}}HandlerMocks.validator)
	})

}

func TestRestHandler_Mount(t *testing.T) {
	testName := shared.SetTestcaseName(1, "rest {{$.module}} handler mount")

	t.Run(testName, func(t *testing.T) {
		restEcho{{clean (upper $.module)}}HandlerMocks()

		// set rest handler
		restHandler := rest{{clean (upper $.module)}}HandlerMocks

		// set echo
		echoHandler := echo.New()
		groupEcho := echoHandler.Group(gofakeit.Word())

		restHandler.Mount(groupEcho)
	})
}

func TestRestHandler_findAll(t *testing.T) {
	type findAll struct {
		result []*domain.{{clean (upper $.module)}}
		meta   *golibshared.Meta
		err    error
	}
	testCase := map[string]struct {
		findAll *findAll
		query   string
	}{
		shared.SetTestcaseName(1, "positive rest {{$.module}} handler get all"): {
			findAll: &findAll{
				result: []*domain.{{clean (upper $.module)}}{},
				meta:   &golibshared.Meta{},
				err:    nil,
			},
		},
		shared.SetTestcaseName(2, "positive rest {{$.module}} handler get all usecase"): {
			findAll: &findAll{
				err: fmt.Errorf(golibshared.ErrorGeneral),
			},
			query: "limit=" + cast.ToString(1),
		},
		shared.SetTestcaseName(3, "positive rest {{$.module}} handler validate filter"): {
			query: "sort=" + gofakeit.Word(),
		},
		shared.SetTestcaseName(4, "positive rest {{$.module}} handler parsing filter"): {
			query: "limit=" + gofakeit.Word(),
		},
	}
	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {

			var (
				headers = map[string]string{}
			)

			// set HTTP mock
			c := golibshared.SetEchoHTTPMock(fmt.Sprintf("%s?%s", gofakeit.URL(), test.query), http.MethodGet, "", headers)

			restEcho{{clean (upper $.module)}}HandlerMocks()

			if test.findAll != nil {
				{{$.module}}UsecaseMocks.On("FindAll", mock.Anything, mock.Anything).Return(test.findAll.result, test.findAll.meta, test.findAll.err).Once()
			}

			// set rest handler
			restHandler := rest{{clean (upper $.module)}}HandlerMocks

			// set handler
			restHandler.findAll(c)

			{{$.module}}UsecaseMocks.AssertExpectations(t)
		})
	}

}

func TestRestHandler_findByID(t *testing.T) {
	type findByID struct {
		result *domain.{{clean (upper $.module)}}
		err    error
	}

	testCase := map[string]struct {
		findByID *findByID
		ID       string
	}{
		shared.SetTestcaseName(1, "positive rest {{$.module}} handler get by ID"): {
			findByID: &findByID{result: &domain.{{clean (upper $.module)}}{}},
			ID:       "5f62fcee09cd352630be5237",
		},
		shared.SetTestcaseName(2, "positive rest {{$.module}} handler get by ID  (ID not found)"): {
			findByID: &findByID{err: fmt.Errorf(golibshared.ErrorGeneral)},
			ID:       "5f62fcee09cd352630be5237",
		},
		shared.SetTestcaseName(3, "positive rest {{$.module}} handler get by ID null"): {
			ID: "",
		},
		shared.SetTestcaseName(4, "positive rest {{$.module}} handler get by ID invalid"): {
			ID: gofakeit.Word(),
		},
		shared.SetTestcaseName(5, "positive rest {{$.module}} handler get by ID (error get DB)"): {
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
			c := golibshared.SetEchoHTTPMock(gofakeit.URL(), http.MethodGet, "", headers)
			c.SetParamNames("id")
			c.SetParamValues(test.ID)

			restEcho{{clean (upper $.module)}}HandlerMocks()

			if test.findByID != nil {
				{{$.module}}UsecaseMocks.On("FindByID", mock.Anything, mock.Anything).Return(test.findByID.result, test.findByID.err).Once()
			}

			// set rest handler
			restHandler := rest{{clean (upper $.module)}}HandlerMocks

			// set handler
			restHandler.findByID(c)

			{{$.module}}UsecaseMocks.AssertExpectations(t)
		})
	}

}

func TestRestHandler_create(t *testing.T) {
	type register struct {
		result *domain.{{clean (upper $.module)}}
		err    error
	}

	type otherData struct {
		username string
	}

	testCase := map[string]struct {
		register  *register
		payload   *domain.{{clean (upper $.module)}}
		otherData *otherData
	}{
		shared.SetTestcaseName(1, "positive rest {{$.module}} handler create new {{$.module}}"): {
			payload: &domain.{{clean (upper $.module)}}{
				Name: "Payjo Suherman",
			},
			register: &register{result: &domain.{{clean (upper $.module)}}{}},
		},
		shared.SetTestcaseName(2, "positive rest {{$.module}} handler create new {{$.module}} handling usecase example name is alredy exist"): {
			payload: &domain.{{clean (upper $.module)}}{
				Name: gofakeit.Name(),
			},
			register: &register{err: fmt.Errorf(golibshared.ErrorGeneral)},
		},
		shared.SetTestcaseName(3, "negative rest {{$.module}} handler create new {{$.module}} validate body"): {
			payload: &domain.{{clean (upper $.module)}}{},
		},
		shared.SetTestcaseName(4, "negative rest {{$.module}} handler create new {{$.module}} unmarshal"): {
			payload: &domain.{{clean (upper $.module)}}{
				Name: gofakeit.Name(),
			},
			otherData: &otherData{
				username: "matmat",
			},
		},
		shared.SetTestcaseName(5, "negative rest {{$.module}} handler create new {{$.module}} validate body"): {
			payload: &domain.{{clean (upper $.module)}}{Name: "P"},
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

			restEcho{{clean (upper $.module)}}HandlerMocks()

			if test.register != nil {
				{{$.module}}UsecaseMocks.On("Create", mock.Anything, mock.Anything).Return(test.register.result, test.register.err).Once()
			}

			// set rest handler
			restHandler := rest{{clean (upper $.module)}}HandlerMocks

			// set handler
			restHandler.create(c)

			{{$.module}}UsecaseMocks.AssertExpectations(t)
		})
	}

}

func TestRestHandler_update(t *testing.T) {
	type update struct {
		result *domain.{{clean (upper $.module)}}
		err    error
	}

	type otherData struct {
		username string
	}

	testCase := map[string]struct {
		update    *update
		payload   *domain.{{clean (upper $.module)}}
		ID        string
		otherData *otherData
	}{
		shared.SetTestcaseName(1, "positive rest {{$.module}} handler update {{$.module}}"): {
			payload: &domain.{{clean (upper $.module)}}{
				Name: gofakeit.Name(),
			},
			ID:     "5f62fcee09cd352630be5237",
			update: &update{result: &domain.{{clean (upper $.module)}}{}},
		},
		shared.SetTestcaseName(2, "positive rest {{$.module}} handler update {{$.module}} usecase handling"): {
			payload: &domain.{{clean (upper $.module)}}{
				Name: gofakeit.Name(),
			},
			ID:     "5f62fcee09cd352630be5237",
			update: &update{err: fmt.Errorf(golibshared.ErrorGeneral)},
		},
		shared.SetTestcaseName(3, "positive rest {{$.module}} handler update {{$.module}} error ID Required"): {
			payload: &domain.{{clean (upper $.module)}}{
				Name: gofakeit.Name(),
			},
		},
		shared.SetTestcaseName(4, "positive rest {{$.module}} handler update {{$.module}} ID invalid"): {
			payload: &domain.{{clean (upper $.module)}}{
				Name: gofakeit.Name(),
			},
			ID: gofakeit.Name(),
		},

		shared.SetTestcaseName(5, "positive rest {{$.module}} handler update {{$.module}} ID atribut minlength"): {
			payload: &domain.{{clean (upper $.module)}}{
				Name: "P",
			},
		},
		shared.SetTestcaseName(5, "positive rest {{$.module}} handler update {{$.module}} ID unmarshal"): {
			payload: &domain.{{clean (upper $.module)}}{
				Name: gofakeit.Name(),
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

			restEcho{{clean (upper $.module)}}HandlerMocks()

			if test.update != nil {
				{{$.module}}UsecaseMocks.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(test.update.result, test.update.err).Once()
			}

			// set rest handler
			restHandler := rest{{clean (upper $.module)}}HandlerMocks

			// set handler
			restHandler.update(c)

			{{$.module}}UsecaseMocks.AssertExpectations(t)
		})
	}
}

func TestRestHandler_delete(t *testing.T) {
	type remove struct {
		result *domain.{{clean (upper $.module)}}
		err    error
	}
	testCase := map[string]struct {
		remove *remove
		ID     string
	}{
		shared.SetTestcaseName(1, "positive rest {{$.module}} handler remove by ID"): {
			remove: &remove{err: nil},
			ID:     "5f62fcee09cd352630be5237",
		},
		shared.SetTestcaseName(2, "positive rest {{$.module}} handler remove by ID  (ID not found)"): {
			remove: &remove{err: fmt.Errorf(golibshared.ErrorGeneral)},
			ID:     "5f62fcee09cd352630be5237",
		},
		shared.SetTestcaseName(3, "positive rest {{$.module}} handler remove by ID null"): {
			ID: "",
		},
		shared.SetTestcaseName(4, "positive rest {{$.module}} handler remove by ID invalid"): {
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
			c := golibshared.SetEchoHTTPMock(gofakeit.URL(), http.MethodGet, "", headers)
			c.SetParamNames("id")
			c.SetParamValues(test.ID)

			restEcho{{clean (upper $.module)}}HandlerMocks()

			if test.remove != nil {
				{{$.module}}UsecaseMocks.On("RemoveByID", mock.Anything, mock.Anything).Return(test.remove.result, test.remove.err).Once()
			}

			// set rest handler
			restHandler := rest{{clean (upper $.module)}}HandlerMocks

			// set handler
			restHandler.delete(c)

			{{$.module}}UsecaseMocks.AssertExpectations(t)
		})
	}
}

func TestRestHandler_restore(t *testing.T) {
	type restore struct {
		result *domain.{{clean (upper $.module)}}
		err    error
	}
	testCase := map[string]struct {
		restore *restore
		ID      string
	}{
		shared.SetTestcaseName(1, "positive rest {{$.module}} handler restore by ID"): {
			restore: &restore{err: nil},
			ID:      "5f62fcee09cd352630be5237",
		},
		shared.SetTestcaseName(2, "positive rest {{$.module}} handler restore by ID  (ID not found)"): {
			restore: &restore{err: fmt.Errorf(golibshared.ErrorGeneral)},
			ID:      "5f62fcee09cd352630be5237",
		},
		shared.SetTestcaseName(3, "positive rest {{$.module}} handler restore by ID null"): {
			ID: "",
		},
		shared.SetTestcaseName(4, "positive rest {{$.module}} handler restore by ID invalid"): {
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
			c := golibshared.SetEchoHTTPMock(gofakeit.URL(), http.MethodGet, "", headers)
			c.SetParamNames("id")
			c.SetParamValues(test.ID)

			restEcho{{clean (upper $.module)}}HandlerMocks()

			if test.restore != nil {
				{{$.module}}UsecaseMocks.On("RestoreByID", mock.Anything, mock.Anything).Return(test.restore.result, test.restore.err).Once()
			}

			// set rest handler
			restHandler := rest{{clean (upper $.module)}}HandlerMocks

			// set handler
			restHandler.restore(c)

			{{$.module}}UsecaseMocks.AssertExpectations(t)
		})
	}
}
`
