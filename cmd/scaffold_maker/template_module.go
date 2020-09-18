package main

//ModuleTemplate module template
const ModuleTemplate = `package {{$.module}}

import (
	"{{$.GoModules}}/internal/modules/{{$.module}}/delivery/resthandler"
	"{{$.GoModules}}/internal/modules/{{$.module}}/repository"
	"{{$.GoModules}}/internal/modules/{{$.module}}/usecase"

	"{{$.LibraryAddress}}/codebase/factory/dependency"
	"{{$.LibraryAddress}}/codebase/factory/types"
	"{{$.LibraryAddress}}/codebase/interfaces"
)

const (
	// Name module name
	Name types.Module = "{{clean (upper $.module)}}"
)

// Module model
type Module struct {
	restHandler *resthandler.RestHandler
}

// NewModule module constructor
func NewModule(deps dependency.Dependency) *Module {
	repo := repository.NewRepository(deps.GetMongoDatabase().ReadDB(), deps.GetMongoDatabase().WriteDB())
	uc := usecase.New{{clean (upper $.module)}}Usecase(repo, deps.GetSDK(), deps.GetValidator(), deps.GetRedisPool().Store())

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
`
const defaultFile = `package {{$.packageName}}`

func defaultDataSource(fileName string) []byte {
	return loadTemplate(defaultFile, map[string]string{"packageName": fileName})
}
