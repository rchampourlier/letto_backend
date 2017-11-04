package design

import (
	. "github.com/goadesign/goa/design/apidsl"
)

// WorkflowPayload defines the data structure used in the create workflow request body.
// It is also the base type for the workflow media type used to render workflows.
var WorkflowPayload = Type("WorkflowPayload", func() {
	Attribute("path", func() {
		MinLength(1)
		Example("group/subgroup")
	})
	Attribute("source_code", func() {
		MinLength(1)
		Example("console.log('Hello world!');")
	})
})
