package router

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Route struct {
	method     string
	pathRegexp *regexp.Regexp
	handler    http.Handler
}

// Handler sets a handler for the route.
func (this *Route) Handler(handler http.Handler) *Route {
	this.handler = handler
	return this
}

// HandlerFunc sets a handler function for the route.
func (this *Route) HandlerFunc(f func(http.ResponseWriter, *http.Request)) *Route {
	return this.Handler(http.HandlerFunc(f))
}

func (this *Route) PathPrefix(template string) *Route {
	r := regexp.MustCompile(`{[^/#?()\.\\]+}`)

	tmpRegx := r.ReplaceAllStringFunc(template, func(m string) string {
		m = m[1 : len(m)-1]
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m)
	})

	this.pathRegexp = regexp.MustCompile("^" + tmpRegx + ".*$")
	return this
}

// Path defines URL template for the route.
// It accepts a template with zero or more URL variables enclosed by {}.
func (this *Route) Path(template string) *Route {
	r := regexp.MustCompile(`{[^/#?()\.\\]+}`)

	tmpRegx := r.ReplaceAllStringFunc(template, func(m string) string {
		m = m[1 : len(m)-1]
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m)
	})

	this.pathRegexp = regexp.MustCompile("^" + tmpRegx + "$")
	return this
}

// Method sets a method for the route.
func (this *Route) Method(method string) *Route {
	this.method = strings.ToUpper(method)
	return this
}
