package controllers

import (
	"github.com/goadesign/goa"
	"gitlab.com/letto/letto_backend/adapters"
	"gitlab.com/letto/letto_backend/app"
	"gitlab.com/letto/letto_backend/util"
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

// List runs the list action.
func (c *WorkflowController) List(ctx *app.ListWorkflowContext) error {
	// WorkflowController_List: start_implement

	workflowPaths, err := c.adapter.ListObjectPaths()
	if err != nil {
		// TODO: should not panic here!
		panic(err)
	}
	workflowList := make(app.LettoWorkflowCollection, len(workflowPaths))
	for i, wp := range workflowPaths {
		workflowList[i] = &app.LettoWorkflow{
			Path: wp,
		}
	}

	// WorkflowController_List: end_implement
	return ctx.OK(workflowList)
}

// Create creates a new workflow and stores it on the chosen adapter
// if the specified parameters are valid.
func (c *WorkflowController) Create(ctx *app.CreateWorkflowContext) error {
	// WorkflowController_Create: start_implement

	path := ctx.Payload.Path
	c.adapter.CreateObject(path, ctx.Payload.SourceCode)
	ctx.ResponseData.Header().Set("Location", app.WorkflowHref(path))

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
	// TODO: should not panic!
	panic("Not implemented")

	// WorkflowController_Update: end_implement
	//return ctx.NoContent()
}
