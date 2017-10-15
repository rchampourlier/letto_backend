//go:generate goagen bootstrap -d github.com/rchampourlier/letto_go/design

package main

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/spf13/afero"

	"github.com/rchampourlier/letto_go/adapters"
	"github.com/rchampourlier/letto_go/app"
	"github.com/rchampourlier/letto_go/controllers"
)

func main() {
	// Create service
	service := goa.New("letto")
	s3 := adapters.NewS3("letto")

	fs := afero.NewOsFs()

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount controllers
	app.MountWorkflowController(service, controllers.NewWorkflowController(service, &s3))
	app.MountTriggersController(service, controllers.NewTriggersController(service, fs))

	// Start service
	if err := service.ListenAndServe(":9292"); err != nil {
		service.LogError("startup", "err", err)
	}

}
