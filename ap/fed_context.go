package ap

import (
	"context"
	"gitlab.cs.fau.de/kissen/fed/db"
)

// The key we use to put stuff into the context
const _FED_CONTEXT_KEY = "FedContext"

type FedContext struct {
	// The URL scheme used in the request, e.g. "https".
	Scheme *string

	// The hostname used in the request, e.g. "example.com".
	Host *string

	// The base path the service is running under, e.g. /ap/.
	BasePath *string

	// Access to the database.
	Storage db.FedStorage

	// The user this request has authenticated as.
	Username *string
}

// Returns a new Context that contains an initialized FedContext.
// Never returns nil.
func WithFedContext(ctx context.Context) context.Context {
	if ctx.Value(_FED_CONTEXT_KEY) != nil {
		panic("context already carries value for _FED_CONTEXT_KEY")
	}

	empty := &FedContext{}
	return context.WithValue(ctx, _FED_CONTEXT_KEY, empty)
}

// Return the FedContext from the provided Context. Never returns nil.
func FromContext(ctx context.Context) *FedContext {
	value := ctx.Value(_FED_CONTEXT_KEY)

	if value == nil {
		panic("no value found for _FED_CONTEXT_KEY")
	}

	fctx, ok := value.(*FedContext)

	if !ok {
		panic("value for _FED_CONTEXT_KEY has wrong type")
	}

	return fctx
}

// Return a pointer pointing to argument s. This is fedutilful when setting
// optional/pointer members of the FedContext struct.
func Just(s string) *string {
	return &s
}
