package main

import (
	"github.com/goadesign/goa"
	"github.com/rchampourlier/letto_backend/app"
)

// WorkflowController implements the workflow resource.
type WorkflowController struct {
	*goa.Controller
}

// NewWorkflowController creates a workflow controller.
func NewWorkflowController(service *goa.Service) *WorkflowController {
	return &WorkflowController{Controller: service.NewController("WorkflowController")}
}

// Create runs the create action.
func (c *WorkflowController) Create(ctx *app.CreateWorkflowContext) error {
	// WorkflowController_Create: start_implement

	// Put your logic here

	// WorkflowController_Create: end_implement
	return nil
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
