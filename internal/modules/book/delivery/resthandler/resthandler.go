package resthandler

import (
	"encoding/json"
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

	book := v1Root.Group("/book")
	// book.GET("", h.hello)
	book.GET("", h.findAll)
}

func (h *RestHandler) hello(c echo.Context) error {
	return wrapper.NewHTTPResponse(http.StatusOK, "Hello, from service: master-service, module: book").JSON(c.Response())
}

func (h *RestHandler) findAll(c echo.Context) error {
	opName := "book_resthandler.get_all"
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
