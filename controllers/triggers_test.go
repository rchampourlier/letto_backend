package controllers_test

/*import (
"context"
"testing"
¨"github.com/goadesign/goa"
"github.com/rchampourlier/letto_go/adapters"
"github.com/rchampourlier/letto_go/app"
"github.com/rchampourlier/letto_go/app/test"
"github.com/rchampourlier/letto_go/controllers")*/

import (
	"testing"
)

// To be tested:
//   - body and headers are correctly extracted and passed to JS,
//   - JS is correctly executed.
func TestTriggersWebhook(t *testing.T) {
	/*var (
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
	}*/
}