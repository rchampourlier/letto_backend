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
	Consumes("application/json")
	Produces("application/json")

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
