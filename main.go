//go:generate goagen bootstrap -d github.com/rchampourlier/letto_backend/design

package main

import (
	"log"
	"os"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/spf13/afero"

	//"github.com/rchampourlier/letto_backend/adapters"
	"github.com/rchampourlier/letto_backend/app"
	"github.com/rchampourlier/letto_backend/controllers"
	"github.com/rchampourlier/letto_backend/exec"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panicf("%s\n", err)
	}

	// Local filesystem Afero wrapper
	fs := afero.NewOsFs()

	// Prepare execution environments
	jsRunner := exec.NewJsRunner(fs)
	if err := jsRunner.Prepare(cwd); err != nil {
		log.Panicf("%s\n", err)
	}

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
