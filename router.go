package groot

import (
	"net/http"
)

type Router struct {
	// Configurable Handler to be used when no route matches.
	NotFoundHandler http.Handler
	// Configurable Handler to be used to recover panic in routes.
	PanicHandler http.Handler
	// Routes to be matched, in order.
	routes []*Route
}
