package groot

import (
	"net/http"
	"regexp"
)

type Route struct {
	// Specify the HTTP method (GET, POST, PUt, etc...)
	httpMethod string
	// Regexp to match route
	pathRegexp *regexp.Regexp
	// Route request handler
	handler http.Handler
}
