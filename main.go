//go:generate goagen bootstrap -d gitlab.com/letto/letto_backend/design

package main

import (
	"log"
	"os"
	"path"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/spf13/afero"

	"gitlab.com/letto/letto_backend/app"
	"gitlab.com/letto/letto_backend/controllers"
	"gitlab.com/letto/letto_backend/exec"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panicf("%s\n", err)
	}

	// Local filesystem Afero wrapper
	fs := afero.NewOsFs()

	// Prepare execution environments
	err = exec.PrepareJsRunner(path.Join(cwd, "exec", "js"))
	if err != nil {
		// The JsRunner failed during preparation
		log.Panicf("%s\n", err)
	}
	jsRunner := exec.NewJsRunner(fs)

	// Create service
	service := goa.New("letto")
	//s3 := adapters.NewS3("letto")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount controllers
	//app.MountWorkflowController(service, controllers.NewWorkflowController(service, &s3))
	app.MountTriggersController(service, controllers.NewTriggersController(service, fs, jsRunner))

	// Start service
	if err := service.ListenAndServe(":9292"); err != nil {
		service.LogError("startup", "err", err)
	}

}
