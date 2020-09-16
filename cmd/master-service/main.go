package main

import (
	"fmt"
	"runtime/debug"

	"github.com/mrapry/go-lib/config"
	service "master-service/internal"
	"github.com/mrapry/go-lib/codebase/app"
)

const (
	serviceName = "master-service"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\x1b[31;1mFailed to start %s service: %v\x1b[0m\n", serviceName, r)
			fmt.Printf("Stack trace: \n%s\n", debug.Stack())
		}
	}()

	cfg := config.Init(fmt.Sprintf("cmd/%s/", serviceName))
	defer cfg.Exit()

	srv := service.NewService(serviceName, cfg)
	app.New(srv).Run()
}

