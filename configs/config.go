package configs

import (
	"context"
	"os"

	"github.com/mrapry/go-lib/codebase/factory/dependency"
	"github.com/mrapry/go-lib/codebase/interfaces"
	"github.com/mrapry/go-lib/config"
	"github.com/mrapry/go-lib/config/database"
	"github.com/mrapry/go-lib/middleware"
	"github.com/mrapry/go-lib/sdk"
	auth_service "github.com/mrapry/go-lib/sdk/auth-service"
	"github.com/mrapry/go-lib/validator"
)

// LoadConfigs load selected dependency configuration in this service
func LoadConfigs(baseCfg *config.Config) (deps dependency.Dependency) {

	loadAdditionalEnv()
	baseCfg.LoadFunc(func(ctx context.Context) []interfaces.Closer {

		// sdk
		authSdk := auth_service.NewAuthService(os.Getenv("AUTH_SERVICE_HOST"), os.Getenv("AUTH_SERVICE_BASIC_AUTH"))
		sdkDeps := sdk.NewSDK(
			sdk.SetAuthService(authSdk),
		)

		mongoDeps := database.InitMongoDB(ctx)

		// inject all service dependencies
		deps = dependency.InitDependency(
			dependency.SetMiddleware(middleware.NewMiddleware(sdkDeps.AuthService())),
			dependency.SetValidator(validator.NewValidator()),
			dependency.SetSDK(sdkDeps),
			dependency.SetMongoDatabase(mongoDeps),
			// ... add more dependencies
		)

		return []interfaces.Closer{mongoDeps} // throw back to config for close connection when application shutdown
	})

	return deps
}
