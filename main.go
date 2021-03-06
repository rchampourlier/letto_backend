//go:generate goagen bootstrap -d gitlab.com/letto/letto_backend/design

package main

import (
	"os"
	"path"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/spf13/afero"

	"gitlab.com/letto/letto_backend/adapters"
	"gitlab.com/letto/letto_backend/app"
	"gitlab.com/letto/letto_backend/controllers"
	"gitlab.com/letto/letto_backend/exec/js"
	"gitlab.com/letto/letto_backend/services"
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
	//
	// It should be passed as the first argument to the application.
	// In a Docker context, this should be done in the `docker-compose.yml`
	// which will also determine the location of the corresponding
	// volume.
	//
	// NB: this directory cannot be determined using `os.Getenv("HOME")`
	//     or any approach relying on the execution environment. Indeed,
	//     the directory is the host's directory, while the application
	//     is to be executed by a Docker container, which will not
	//     have the correct path for the data dir.
	hostDataDir := os.Args[1]

	// `appDataDir` indicates where the `data` directory will be
	// for the Go app.
	appDataDir := "/tmp/data"

	// `appLogsDir` specifies where the logs will be stored.
	appLogsDir := "/tmp/logs"

	// `execDataDir` specified where the `data` directory will be
	// in the context of the execution container.
	execDataDir := "/usr/src/app/data"

	// Prepare execution environments
	jsRunner, err := js.NewRunner(fs, hostDataDir, appDataDir, execDataDir)
	if err != nil {
		panic(err)
	}

	// Event-bus and services
	eventBus := services.NewEventBus(fs, appLogsDir)
	serviceActivateTrigger := services.NewActivateTrigger(&eventBus)
	serviceActivateTrigger.StartConsuming()
	serviceExecuteWorkflows := services.NewExecuteWorkflows(&eventBus, jsRunner)
	serviceExecuteWorkflows.StartConsuming()

	// Create service
	service := goa.New("letto")
	workflowsAdapter := adapters.NewAferoFs(fs, path.Join(appDataDir, "workflows"))

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount controllers
	app.MountWorkflowController(service, controllers.NewWorkflowController(service, &workflowsAdapter))
	app.MountTriggersController(service, controllers.NewTriggersController(service, &eventBus, jsRunner))

	// Start service
	if err := service.ListenAndServe(":9292"); err != nil {
		service.LogError("startup", "err", err)
	}
}
