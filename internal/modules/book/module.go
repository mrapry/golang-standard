package book

import (
	"master-service/internal/modules/book/delivery/resthandler"
	"master-service/internal/modules/book/repository"
	"master-service/internal/modules/book/usecase"

	"github.com/mrapry/go-lib/codebase/factory/dependency"
	"github.com/mrapry/go-lib/codebase/factory/types"
	"github.com/mrapry/go-lib/codebase/interfaces"
)

const (
	// Name module name
	Name types.Module = "Book"
)

// Module model
type Module struct {
	restHandler *resthandler.RestHandler
}

// NewModule module constructor
func NewModule(deps dependency.Dependency) *Module {
	repo := repository.NewRepository(deps.GetMongoDatabase().ReadDB(), deps.GetMongoDatabase().WriteDB())
	uc := usecase.NewBookUsecase(repo, deps.GetSDK(), deps.GetValidator(), deps.GetRedisPool().Store())

	var mod Module
	mod.restHandler = resthandler.NewRestHandler(uc, deps.GetMiddleware(), deps.GetValidator())

	return &mod
}

// RestHandler method
func (m *Module) RestHandler() interfaces.EchoRestHandler {
	return m.restHandler
}

// GRPCHandler method
func (m *Module) GRPCHandler() interfaces.GRPCHandler {
	return nil
}

// GraphQLHandler method
func (m *Module) GraphQLHandler() interfaces.GraphQLHandler {
	return nil
}

// WorkerHandler method
func (m *Module) WorkerHandler(workerType types.Worker) interfaces.WorkerHandler {
	return nil
}

// Name get module name
func (m *Module) Name() types.Module {
	return Name
}
