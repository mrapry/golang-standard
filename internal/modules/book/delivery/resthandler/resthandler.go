package resthandler

import (
	"encoding/json"
	"io/ioutil"
	"master-service/internal/modules/book/domain"
	"master-service/internal/modules/book/usecase"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mrapry/go-lib/codebase/interfaces"
	"github.com/mrapry/go-lib/golibhelper"
	helper "github.com/mrapry/go-lib/golibhelper"
	"github.com/mrapry/go-lib/logger"
	"github.com/mrapry/go-lib/tracer"
	"github.com/mrapry/go-lib/wrapper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap/zapcore"
)

// RestHandler handler
type RestHandler struct {
	bookUsecase usecase.BookUsecase
	mw          interfaces.Middleware
	validator   interfaces.Validator
}

// NewRestHandler create new rest handler
func NewRestHandler(bookUsecase usecase.BookUsecase, mw interfaces.Middleware, validator interfaces.Validator) *RestHandler {
	return &RestHandler{
		bookUsecase: bookUsecase,
		mw:          mw,
		validator:   validator,
	}
}

// Mount handler with root "/"
// handling version in here
func (h *RestHandler) Mount(root *echo.Group) {
	v1Root := root.Group(helper.V1)

	book := v1Root.Group("/book", h.mw.HTTPBasicAuth(false))
	book.GET("", h.findAll)
	book.GET("/:id", h.findByID)
	book.POST("", h.create)
	book.PUT("/:id", h.update)
	book.DELETE("/:id", h.delete)
	book.PATCH("/:id", h.restore)
}

func (h *RestHandler) findAll(c echo.Context) error {
	opName := "book_resthandler.find_all"
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
	if err := h.validator.ValidateDocument("book/get_all", body); err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "validate_payload")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "Gagal memvalidasi filter", err.Error()).JSON(c.Response())
	}

	result, meta, err := h.bookUsecase.FindAll(ctx, &filter)
	if err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "usecase_error")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c.Response())
	}
	return wrapper.NewHTTPResponse(http.StatusOK, "Sukses mengambil data book", result, meta).JSON(c.Response())
}

func (h *RestHandler) findByID(c echo.Context) error {
	opName := "book_resthandler.find_by_id"
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

	result, err := h.bookUsecase.FindByID(ctx, id)
	if err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "usecase_error")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "Gagal mendapatkan data book", err.Error()).JSON(c.Response())
	}

	return wrapper.NewHTTPResponse(http.StatusOK, "Sukses", result).JSON(c.Response())
}

func (h *RestHandler) create(c echo.Context) error {
	opName := "book_resthandler.create"
	ctx := c.Request().Context()

	tracer := tracer.StartTrace(ctx, opName)
	defer tracer.Finish()
	ctx = tracer.Context()

	body, _ := ioutil.ReadAll(c.Request().Body)
	if err := h.validator.ValidateDocument("book/create", body); err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "validate_payload")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "Gagal dalam validasi data", err.Error()).JSON(c.Response())
	}

	var payload domain.Book
	if err := json.Unmarshal(body, &payload); err != nil {
		return wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c.Response())
	}

	result, err := h.bookUsecase.Create(ctx, &payload)
	if err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "usecase_error")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "Gagal menyimpan data book baru", err.Error()).JSON(c.Response())
	}

	return wrapper.NewHTTPResponse(http.StatusCreated, "Sukses menyimpan data user", result).JSON(c.Response())
}

func (h *RestHandler) update(c echo.Context) error {
	opName := "book_resthandler.update"
	ctx := c.Request().Context()

	tracer := tracer.StartTrace(ctx, opName)
	defer tracer.Finish()
	ctx = tracer.Context()

	body, _ := ioutil.ReadAll(c.Request().Body)
	if err := h.validator.ValidateDocument("book/update", body); err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "validate_payload")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "Gagal dalam validasi data", err.Error()).JSON(c.Response())
	}

	var payload domain.Book
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

	result, err := h.bookUsecase.Update(ctx, &payload, ID)
	if err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "usecase_error")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "Gagal update data book baru", err.Error()).JSON(c.Response())
	}

	return wrapper.NewHTTPResponse(http.StatusCreated, "Sukses update data user", result).JSON(c.Response())
}

func (h *RestHandler) delete(c echo.Context) error {
	opName := "book_resthandler.delete"
	ctx := c.Request().Context()

	tracer := tracer.StartTrace(ctx, opName)
	defer tracer.Finish()
	ctx = tracer.Context()

	id := c.Param("id")
	if id == "" {
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "id tidak boleh kosong").JSON(c.Response())
	}

	err := h.bookUsecase.RemoveByID(ctx, id)
	if err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "usecase_error")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "Gagal delete data book", err.Error()).JSON(c.Response())
	}

	return wrapper.NewHTTPResponse(http.StatusOK, "Sukses delete data").JSON(c.Response())
}

func (h *RestHandler) restore(c echo.Context) error {
	opName := "book_resthandler.restore"
	ctx := c.Request().Context()

	tracer := tracer.StartTrace(ctx, opName)
	defer tracer.Finish()
	ctx = tracer.Context()

	id := c.Param("id")
	if id == "" {
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "id tidak boleh kosong").JSON(c.Response())
	}

	err := h.bookUsecase.RestoreByID(ctx, id)
	if err != nil {
		logger.Log(zapcore.ErrorLevel, err.Error(), opName, "usecase_error")
		return wrapper.NewHTTPResponse(http.StatusBadRequest, "Gagal merestore data book", err.Error()).JSON(c.Response())
	}

	return wrapper.NewHTTPResponse(http.StatusOK, "Sukses restore data").JSON(c.Response())
}
