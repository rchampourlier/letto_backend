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
	"unicode/utf8"
)

// An automation workflow (default view)
//
// Identifier: application/letto.workflow+json; view=default
type LettoWorkflow struct {
	// Path to the workflow
	Path string `form:"path" json:"path" xml:"path"`
}

// Validate validates the LettoWorkflow media type instance.
func (mt *LettoWorkflow) Validate() (err error) {
	if mt.Path == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "path"))
	}
	if utf8.RuneCountInString(mt.Path) < 1 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.path`, mt.Path, utf8.RuneCountInString(mt.Path), 1, true))
	}
	return
}

// An automation workflow (full view)
//
// Identifier: application/letto.workflow+json; view=full
type LettoWorkflowFull struct {
	// Path to the workflow
	Path string `form:"path" json:"path" xml:"path"`
	// Source code of the workflow
	SourceCode string `form:"source_code" json:"source_code" xml:"source_code"`
}

// Validate validates the LettoWorkflowFull media type instance.
func (mt *LettoWorkflowFull) Validate() (err error) {
	if mt.Path == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "path"))
	}
	if mt.SourceCode == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "source_code"))
	}
	if utf8.RuneCountInString(mt.Path) < 1 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.path`, mt.Path, utf8.RuneCountInString(mt.Path), 1, true))
	}
	if utf8.RuneCountInString(mt.SourceCode) < 1 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.source_code`, mt.SourceCode, utf8.RuneCountInString(mt.SourceCode), 1, true))
	}
	return
}

// LettoWorkflowCollection is the media type for an array of LettoWorkflow (default view)
//
// Identifier: application/letto.workflow+json; type=collection; view=default
type LettoWorkflowCollection []*LettoWorkflow

// Validate validates the LettoWorkflowCollection media type instance.
func (mt LettoWorkflowCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// LettoWorkflowCollection is the media type for an array of LettoWorkflow (full view)
//
// Identifier: application/letto.workflow+json; type=collection; view=full
type LettoWorkflowFullCollection []*LettoWorkflowFull

// Validate validates the LettoWorkflowFullCollection media type instance.
func (mt LettoWorkflowFullCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}
