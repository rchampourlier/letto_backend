package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var Workflow = MediaType("application/letto.workflow+json", func() {
	Description("An automation workflow")
	Reference(WorkflowPayload)
	Attributes(func() {
		Attribute("path", String, "Path to the workflow")
		Attribute("source_code", String, "Source code of the workflow")
		Required("path", "source_code")
	})

	View("default", func() {
		Attribute("path")
	})

	View("full", func() {
		Attribute("path")
		Attribute("source_code")
	})
})
