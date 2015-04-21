package router

import (
	"log"
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

// If a panic handler was defined it calls
func (this *Router) handlePanic(w http.ResponseWriter, r *http.Request) {
	if this.PanicHandler != nil {
		err := recover()
		if err != nil {
			log.Println(err)
			this.PanicHandler.ServeHTTP(w, r)
		}
	}
}

func (this *Router) parseRequestParams(r *http.Request, route *Route, matches []string) {
	r.ParseForm()
	if len(matches) > 0 && matches[0] == r.URL.Path {
		for i, name := range route.pathRegexp.SubexpNames() {
			if len(name) > 0 {
				r.Form.Add(name, matches[i])
			}
		}
	}
}

// NewRouter returns a new router instance.
func NewRouter() *Router {
	return &Router{}
}

// NewRoute registers an empty route.
func (this *Router) NewRoute() *Route {
	r := &Route{}
	r.method = "GET"
	this.routes = append(this.routes, r)
	return r
}

func (this *Router) Handle(path string, handler http.Handler) *Route {
	return this.NewRoute().Path(path).Handler(handler)
}

func (this *Router) HandleFunc(path string, f func(http.ResponseWriter,
	*http.Request)) *Route {
	return this.NewRoute().Path(path).HandlerFunc(f)
}

func (this *Router) FindMatchingRoute(r *http.Request) *Route {
	for _, route := range this.routes {
		if route.pathRegexp != nil && route.method == r.Method {
			if matches := route.pathRegexp.FindStringSubmatch(r.URL.Path); matches != nil {
				this.parseRequestParams(r, route, matches)
				return route
			}
		}
	}
	return nil
}

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer this.handlePanic(w, r)

	var handler http.Handler
	var matchingRoute *Route

	if matchingRoute = this.FindMatchingRoute(r); matchingRoute != nil {
		handler = matchingRoute.handler
	}

	if handler == nil {
		handler = this.NotFoundHandler
		if handler == nil {
			handler = http.NotFoundHandler()
		}
	}

	handler.ServeHTTP(w, r)
}
