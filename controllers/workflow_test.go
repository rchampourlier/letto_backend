package controllers_test

import (
	"context"
	"testing"

	"github.com/goadesign/goa"
	"gitlab.com/letto/letto_backend/adapters"
	"gitlab.com/letto/letto_backend/app"
	"gitlab.com/letto/letto_backend/app/test"
	"gitlab.com/letto/letto_backend/controllers"
)

// Test cases:
//   - S01. label missing
//   - S02. label invalid
//   - S03. source missing
//   - S04. source invalid
//   - S05. filter missing
//   - S06. filter invalid
//   - S07. params valid, adapter failure
//   - H01. params valid, adapter success
//   - H02. params valid, adapter success, with a name that has already been used
//
// Happy cases:
//   - WorkflowCreated response
//   - A workflow must have been created at the expected path
//   - "Location" header contains the URL for the created workflow
//
// Sad cases
//   - WorkflowBadRequest response
//   - Workflow must not be created
//   - Response contains error message
func TestCreateWorkflow(t *testing.T) {
	var (
		service = goa.New("letto-test")
		adapter = adapters.NewInMemory()
		ctrl    = controllers.NewWorkflowController(service, &adapter)
		ctx     = context.Background()
	)

	// Happy cases

	// H01
	payload := &app.CreateWorkflowPayload{
		Group:  "test-group",
		Name:   "Test name",
		Source: "Test source"}

	// WorkflowCreated response
	r := test.CreateWorkflowCreated(t, ctx, service, ctrl, payload)

	// A workflow must have been created at the expected path
	if count, _ := adapter.Count(); count != 1 {
		t.Errorf("expected 1 created workflow, found %d\n", count)
	}

	// "Location" header contains the URL for the created workflow
	loc := r.Header().Get("Location")
	expStr := "/workflows/test-group/test-name.js"
	if loc != expStr {
		t.Errorf("invalid location: got %s, expected %s", loc, expStr)
	}
}
