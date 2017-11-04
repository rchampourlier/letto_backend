package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("triggers", func() {
	BasePath("/triggers")

	Action("webhook", func() {
		Description("Receive incoming webhook")
		Routing(GET("/webhook/*group"), POST("/webhook/*group"))
		Response(OK)
	})
})

var _ = Resource("workflow", func() {

	BasePath("/workflows")
	DefaultMedia(Workflow)

	Action("list", func() {
		Description("List workflows")
		Routing(GET("/"))
		Response(OK, CollectionOf(Workflow))
	})

	Action("create", func() {
		Description("Create a new workflow")
		Routing(POST(""))
		Payload(WorkflowPayload, func() {
			Required("path", "source_code")
		})
		Response(Created, `\A/workflows/[^/]+/[\w-]+\z`)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("read", func() {
		Description("Read a workflow by ID")
		Routing(GET("/:workflowID"))
		Params(func() {
			Param("workflowID", String, "Workflow ID")
		})
		Response(OK)
	})

	Action("update", func() {
		Description("Update an existing workflow by ID")
		Routing(PUT("/:workflowID"))
		Params(func() {
			Param("workflowID", String, "Workflow ID")
		})
		Payload(WorkflowPayload, func() {
			Required("path", "source_code")
		})
		Response(NoContent)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})

	Action("delete", func() {
		Description("Delete a workflow by ID")
		Routing(DELETE("/:workflowID"))
		Params(func() {
			Param("workflowID", String, "Workflow ID")
		})
		Response(OK)
	})
})
