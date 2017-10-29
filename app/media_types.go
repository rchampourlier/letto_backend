// Code generated by goagen v1.3.0, DO NOT EDIT.
//
// API "letto": Application Media Types
//
// Command:
// $ goagen
// --design=gitlab.com/letto/letto_backend/design
// --out=$(GOPATH)/src/gitlab.com/letto/letto_backend
// --version=v1.3.0

package app

import (
	"github.com/goadesign/goa"
	uuid "github.com/satori/go.uuid"
)

// An automation workflow (default view)
//
// Identifier: application/letto.workflow+json; view=default
type LettoWorkflow struct {
	// API href for reading a workflow
	Href string `form:"href" json:"href" xml:"href"`
	// Name for the workflow
	Name string `form:"name" json:"name" xml:"name"`
	// Workflow UUID
	UUID uuid.UUID `form:"uuid" json:"uuid" xml:"uuid"`
}

// Validate validates the LettoWorkflow media type instance.
func (mt *LettoWorkflow) Validate() (err error) {

	if mt.Href == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "href"))
	}
	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	return
}

// An automation workflow (full view)
//
// Identifier: application/letto.workflow+json; view=full
type LettoWorkflowFull struct {
	// API href for reading a workflow
	Href string `form:"href" json:"href" xml:"href"`
	// Name for the workflow
	Name string `form:"name" json:"name" xml:"name"`
	// Source code for the workflow
	Source string `form:"source" json:"source" xml:"source"`
	// Workflow UUID
	UUID uuid.UUID `form:"uuid" json:"uuid" xml:"uuid"`
}

// Validate validates the LettoWorkflowFull media type instance.
func (mt *LettoWorkflowFull) Validate() (err error) {

	if mt.Href == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "href"))
	}
	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if mt.Source == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "source"))
	}
	return
}

// An automation workflow (link view)
//
// Identifier: application/letto.workflow+json; view=link
type LettoWorkflowLink struct {
	// API href for reading a workflow
	Href string `form:"href" json:"href" xml:"href"`
}

// Validate validates the LettoWorkflowLink media type instance.
func (mt *LettoWorkflowLink) Validate() (err error) {
	if mt.Href == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "href"))
	}
	return
}

// LettoWorkflow_list media type (default view)
//
// Identifier: application/letto.workflow_list+json; view=default
type LettoWorkflowList struct {
	// Links to related resources
	Links *LettoWorkflowListLinks `form:"links,omitempty" json:"links,omitempty" xml:"links,omitempty"`
}

// LettoWorkflow_listLinks contains links to related resources of LettoWorkflow_list.
type LettoWorkflowListLinks struct {
}
