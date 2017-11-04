// Code generated by goagen v1.3.0, DO NOT EDIT.
//
// API "letto": Application Controllers
//
// Command:
// $ goagen
// --design=gitlab.com/letto/letto_backend/design
// --out=$(GOPATH)/src/gitlab.com/letto/letto_backend
// --version=v1.3.0

package app

import (
	"context"
	"github.com/goadesign/goa"
	"net/http"
)

// initService sets up the service encoders, decoders and mux.
func initService(service *goa.Service) {
	// Setup encoders and decoders
	service.Encoder.Register(goa.NewJSONEncoder, "application/json")
	service.Decoder.Register(goa.NewJSONDecoder, "application/json")

	// Setup default encoder and decoder
	service.Encoder.Register(goa.NewJSONEncoder, "*/*")
	service.Decoder.Register(goa.NewJSONDecoder, "*/*")
}

// HealthController is the controller interface for the Health actions.
type HealthController interface {
	goa.Muxer
	Health(*HealthHealthContext) error
}

// MountHealthController "mounts" a Health resource controller on the given service.
func MountHealthController(service *goa.Service, ctrl HealthController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewHealthHealthContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Health(rctx)
	}
	service.Mux.Handle("GET", "/api/_ah/health", ctrl.MuxHandler("health", h, nil))
	service.LogInfo("mount", "ctrl", "Health", "action", "Health", "route", "GET /api/_ah/health")
}

// TriggersController is the controller interface for the Triggers actions.
type TriggersController interface {
	goa.Muxer
	Webhook(*WebhookTriggersContext) error
}

// MountTriggersController "mounts" a Triggers resource controller on the given service.
func MountTriggersController(service *goa.Service, ctrl TriggersController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewWebhookTriggersContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Webhook(rctx)
	}
	service.Mux.Handle("GET", "/api/triggers/webhook/*group", ctrl.MuxHandler("webhook", h, nil))
	service.LogInfo("mount", "ctrl", "Triggers", "action", "Webhook", "route", "GET /api/triggers/webhook/*group")
	service.Mux.Handle("POST", "/api/triggers/webhook/*group", ctrl.MuxHandler("webhook", h, nil))
	service.LogInfo("mount", "ctrl", "Triggers", "action", "Webhook", "route", "POST /api/triggers/webhook/*group")
}

// WorkflowController is the controller interface for the Workflow actions.
type WorkflowController interface {
	goa.Muxer
	Create(*CreateWorkflowContext) error
	Delete(*DeleteWorkflowContext) error
	List(*ListWorkflowContext) error
	Read(*ReadWorkflowContext) error
	Update(*UpdateWorkflowContext) error
}

// MountWorkflowController "mounts" a Workflow resource controller on the given service.
func MountWorkflowController(service *goa.Service, ctrl WorkflowController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewCreateWorkflowContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*CreateWorkflowPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Create(rctx)
	}
	service.Mux.Handle("POST", "/api/workflows", ctrl.MuxHandler("create", h, unmarshalCreateWorkflowPayload))
	service.LogInfo("mount", "ctrl", "Workflow", "action", "Create", "route", "POST /api/workflows")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewDeleteWorkflowContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Delete(rctx)
	}
	service.Mux.Handle("DELETE", "/api/workflows/:workflowID", ctrl.MuxHandler("delete", h, nil))
	service.LogInfo("mount", "ctrl", "Workflow", "action", "Delete", "route", "DELETE /api/workflows/:workflowID")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewListWorkflowContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.List(rctx)
	}
	service.Mux.Handle("GET", "/api/workflows/", ctrl.MuxHandler("list", h, nil))
	service.LogInfo("mount", "ctrl", "Workflow", "action", "List", "route", "GET /api/workflows/")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewReadWorkflowContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Read(rctx)
	}
	service.Mux.Handle("GET", "/api/workflows/:workflowID", ctrl.MuxHandler("read", h, nil))
	service.LogInfo("mount", "ctrl", "Workflow", "action", "Read", "route", "GET /api/workflows/:workflowID")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewUpdateWorkflowContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*UpdateWorkflowPayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Update(rctx)
	}
	service.Mux.Handle("PUT", "/api/workflows/:workflowID", ctrl.MuxHandler("update", h, unmarshalUpdateWorkflowPayload))
	service.LogInfo("mount", "ctrl", "Workflow", "action", "Update", "route", "PUT /api/workflows/:workflowID")
}

// unmarshalCreateWorkflowPayload unmarshals the request body into the context request data Payload field.
func unmarshalCreateWorkflowPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &createWorkflowPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}

// unmarshalUpdateWorkflowPayload unmarshals the request body into the context request data Payload field.
func unmarshalUpdateWorkflowPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &updateWorkflowPayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}
