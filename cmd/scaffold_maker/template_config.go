package main

const configTemplate = `package configs

import (
	"context"
	"os"

	"{{.LibraryAddress}}/codebase/factory/dependency"
	"{{.LibraryAddress}}/codebase/interfaces"
	"{{.LibraryAddress}}/config"
	"{{.LibraryAddress}}/config/database"
	"{{.LibraryAddress}}/middleware"
	"{{.LibraryAddress}}/sdk"
	auth_service "{{.LibraryAddress}}/sdk/auth-service"
	"{{.LibraryAddress}}/validator"
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
`
const configLoadEnvTemplate = `package configs

// Environment additional in this service
type Environment struct {
	
}

var env Environment

// GetEnv get global additional environment
func GetEnv() Environment {
	return env
}

func loadAdditionalEnv() {
	
}
`
