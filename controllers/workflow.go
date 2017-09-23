package controllers

import (
	"fmt"

	"github.com/goadesign/goa"
	"github.com/rchampourlier/letto_go/adapters"
	"github.com/rchampourlier/letto_go/app"
	"github.com/rchampourlier/letto_go/util"
)

// WorkflowController implements the workflow resource.
type WorkflowController struct {
	*goa.Controller
	adapter adapters.Adapter
}

// NewWorkflowController creates a workflow controller.
func NewWorkflowController(service *goa.Service, adapter adapters.Adapter) *WorkflowController {
	return &WorkflowController{
		Controller: service.NewController("WorkflowController"),
		adapter:    adapter,
	}
}

// TODO: implement
func validateGroup(group string) error {
	return nil
}

func parameterize(s string) string {
	return util.Parameterize(s, '-')
}

func genPath(group string, name string) string {
	pName := parameterize(name)
	return fmt.Sprintf("/workflows/%s/%s.js", group, pName)
}

// Create creates a new workflow and stores it on the chosen adapter
// if the specified parameters are valid.
func (c *WorkflowController) Create(ctx *app.CreateWorkflowContext) error {
	// WorkflowController_Create: start_implement

	name := ctx.Payload.Name
	group := ctx.Payload.Group
	path := genPath(group, name)

	c.adapter.CreateObject(path, ctx.Payload.Source)
	ctx.ResponseData.Header().Set("Location", path)

	// WorkflowController_Create: end_implement

	return ctx.Created()
}

// Delete runs the delete action.
func (c *WorkflowController) Delete(ctx *app.DeleteWorkflowContext) error {
	// WorkflowController_Delete: start_implement

	// Put your logic here

	// WorkflowController_Delete: end_implement
	res := &app.LettoWorkflow{}
	return ctx.OK(res)
}

// List runs the list action.
func (c *WorkflowController) List(ctx *app.ListWorkflowContext) error {
	// WorkflowController_List: start_implement

	// Put your logic here

	// WorkflowController_List: end_implement
	res := &app.LettoWorkflowList{}
	return ctx.OK(res)
}

// Read runs the read action.
func (c *WorkflowController) Read(ctx *app.ReadWorkflowContext) error {
	// WorkflowController_Read: start_implement

	// Put your logic here

	// WorkflowController_Read: end_implement
	res := &app.LettoWorkflow{}
	return ctx.OK(res)
}

// Update runs the update action.
func (c *WorkflowController) Update(ctx *app.UpdateWorkflowContext) error {
	// WorkflowController_Update: start_implement

	// Put your logic here

	// WorkflowController_Update: end_implement
	res := &app.LettoWorkflow{}
	return ctx.OK(res)
}
