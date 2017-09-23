package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("letto", func() {
	Title("The place for your complex workflows")
	Description("Go service providing the backend API for Letto")
	Host("localhost:9292")
	Scheme("http")
	BasePath("/api")

	ResponseTemplate(Created, func(pattern string) {
		Description("Resource created")
		Status(201)
		Headers(func() {
			Header("Location", String, "href to created resource", func() {
				Pattern(pattern)
			})
		})
	})
})

var _ = Resource("workflow", func() {

	BasePath("/workflows")
	DefaultMedia(WorkflowMedia)

	Action("list", func() {
		Description("List workflows")
		Routing(GET("/"))
		Response(OK, WorkflowListMedia)
	})

	Action("create", func() {
		Description("Create a new workflow")
		Routing(POST(""))
		Payload(func() {
			Member("source", String, "Source code to execute for this workflow")
			Member("name", String, "Name of the workflow")
			Member("group", String, "A way of grouping workflows together to be triggered by a specific endpoint's URL", func() {
				Pattern(`\A[\w-]+\z`)
			})
			Required("source", "name", "group")
		})
		Response(Created, `\A/workflows/[^/]+/[\w-]+\z`)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("read", func() {
		Description("Read a workflow by UUID")
		Routing(GET("/:workflowUUID"))
		Params(func() {
			Param("workflowUUID", String, "Workflow UUID")
		})
		Response(OK)
	})

	Action("update", func() {
		Description("Update an existing workflow by UUID")
		Routing(PUT("/:workflowUUID"))
		Params(func() {
			Param("workflowUUID", String, "Workflow UUID")
			// TODO: add workflow, need a type
		})
		Response(OK)
	})

	Action("delete", func() {
		Description("Delete a workflow by UUID")
		Routing(DELETE("/:workflowUUID"))
		Params(func() {
			Param("workflowUUID", String, "Workflow UUID")
		})
		Response(OK)
	})
})

var WorkflowMedia = MediaType("application/letto.workflow+json", func() {
	Description("An automation workflow")
	Attributes(func() {
		Attribute("uuid", UUID, "Workflow UUID")
		Attribute("href", String, "API href for reading a workflow")
		Attribute("name", String, "Name for the workflow")
		Attribute("source", String, "Source code for the workflow")
		Required("uuid", "href", "name", "source")
	})
	View("default", func() {
		Attribute("uuid")
		Attribute("href")
		Attribute("name")
	})
	View("full", func() {
		Attribute("uuid")
		Attribute("href")
		Attribute("name")
		Attribute("source")
	})
	View("link", func() {
		Attribute("href")
	})
})

var WorkflowListMedia = MediaType("application/letto.workflow_list+json", func() {
	ArrayOf(WorkflowMedia)
	View("default", func() {
		Attribute("links")
	})
})
