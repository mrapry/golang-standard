package service

import (
	"master-service/configs"
	"master-service/internal/modules/book"

	"github.com/mrapry/go-lib/codebase/factory"
	"github.com/mrapry/go-lib/codebase/factory/dependency"
	"github.com/mrapry/go-lib/codebase/factory/types"
	"github.com/mrapry/go-lib/config"
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
		book.NewModule(deps),
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
