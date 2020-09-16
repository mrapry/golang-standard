package main

const moduleMainTemplate = `package {{clean $.module}}

import (
	"{{$.GoModules}}/internal/modules/{{$.module}}/delivery/graphqlhandler"
	"{{$.GoModules}}/internal/modules/{{$.module}}/delivery/grpchandler"
	"{{$.GoModules}}/internal/modules/{{$.module}}/delivery/resthandler"
	"{{$.GoModules}}/internal/modules/{{$.module}}/delivery/workerhandler"
	"{{$.LibraryAddress}}/pkg/codebase/factory/dependency"
	"{{$.LibraryAddress}}/pkg/codebase/interfaces"
	"{{$.LibraryAddress}}/pkg/codebase/factory/types"
)

const (
	// Name module name
	Name types.Module = "{{clean (upper $.module)}}"
)

// Module model
type Module struct {
	restHandler    *resthandler.RestHandler
	grpcHandler    *grpchandler.GRPCHandler
	graphqlHandler *graphqlhandler.GraphQLHandler

	workerHandlers map[types.Worker]interfaces.WorkerHandler
}

// NewModule module constructor
func NewModule(deps dependency.Dependency) *Module {
	var mod Module
	mod.restHandler = resthandler.NewRestHandler(deps.GetMiddleware())
	mod.grpcHandler = grpchandler.NewGRPCHandler(deps.GetMiddleware())
	mod.graphqlHandler = graphqlhandler.NewGraphQLHandler(deps.GetMiddleware())

	mod.workerHandlers = map[types.Worker]interfaces.WorkerHandler{
		types.Kafka:           workerhandler.NewKafkaHandler(),
		types.Scheduler:       workerhandler.NewCronHandler(),
		types.RedisSubscriber: workerhandler.NewRedisHandler(),
		types.TaskQueue:       workerhandler.NewTaskQueueHandler(),
	}

	return &mod
}

// RestHandler method
func (m *Module) RestHandler() interfaces.EchoRestHandler {
	return m.restHandler
}

// GRPCHandler method
func (m *Module) GRPCHandler() interfaces.GRPCHandler {
	return m.grpcHandler
}

// GraphQLHandler method
func (m *Module) GraphQLHandler() interfaces.GraphQLHandler {
	return m.graphqlHandler
}

// WorkerHandler method
func (m *Module) WorkerHandler(workerType types.Worker) interfaces.WorkerHandler {
	return m.workerHandlers[workerType]
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
