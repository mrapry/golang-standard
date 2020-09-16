package main

const serviceMainTemplate = `package service

import (
	"{{$.GoModules}}/configs"

	"{{$.LibraryAddress}}/codebase/factory"
	"{{$.LibraryAddress}}/codebase/factory/dependency"
	"{{$.LibraryAddress}}/codebase/factory/types"
	"{{$.LibraryAddress}}/config"
)

//Service structure
type Service struct {
	deps    dependency.Dependency
	modules []factory.ModuleFactory
	name    types.Service
}

// NewService in this service
func NewService(serviceName string, cfg *config.Config) factory.ServiceFactory {
	deps := configs.LoadConfigs(cfg)

	modules := []factory.ModuleFactory{
		// name_of_module.NewModule(deps),
	}

	return &Service{
		deps:    deps,
		modules: modules,
		name:    types.Service(serviceName),
	}
}

// GetDependency method
func (s *Service) GetDependency() dependency.Dependency {
	return s.deps
}

// GetModules method
func (s *Service) GetModules() []factory.ModuleFactory {
	return s.modules
}

// Name method
func (s *Service) Name() types.Service {
	return s.name
}
`
