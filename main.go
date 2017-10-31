//go:generate goagen bootstrap -d gitlab.com/letto/letto_backend/design

package main

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/spf13/afero"

	"gitlab.com/letto/letto_backend/app"
	"gitlab.com/letto/letto_backend/controllers"
	"gitlab.com/letto/letto_backend/exec/js"
)

func main() {
	// Local filesystem Afero wrapper
	fs := afero.NewOsFs()

	// TODO: make the following variables a configuration. These
	//   3 variables depend on the `docker-compose` config.
	// NB: if you need to run the Go app outside of Docker, you'll
	//     have to change the variables below.

	// `hostDataDir` defines where the `data` directory is contained
	// on the host.
	hostDataDir := "/home/ubuntu/letto_data"

	// `appDataDir` indicates where the `data` directory will be
	// for the Go app.
	appDataDir := "/tmp/data"

	// TODO: make the traces dir configurable in `Trace`
	//appTracesDir := "/tmp/traces"

	// `execDataDir` specified where the `data` directory will be
	// in the context of the execution container.
	execDataDir := "/usr/src/app/data"

	// Prepare execution environments
	jsRunner, err := js.NewRunner(fs, hostDataDir, appDataDir, execDataDir)
	if err != nil {
		panic(err)
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
