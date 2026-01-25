//go:build !wasm

package preveltekit

import (
	"os"
	"strings"
)

// Route defines a single route with path pattern and handler
type Route struct {
	Path    string
	Handler func(params map[string]string)
}

// Router handles client-side routing (stub for SSR)
type Router struct {
	routes      []Route
	notFound    func()
	currentPath *Store[string]
	beforeNav   func(from, to string) bool
}

// NewRouter creates a new router instance
func NewRouter() *Router {
	return &Router{
		currentPath: New(""),
	}
}

// Handle registers a route handler for a path pattern
func (r *Router) Handle(path string, handler func(params map[string]string)) {
	r.routes = append(r.routes, Route{Path: path, Handler: handler})
}

// RegisterRoutes registers multiple routes from StaticRoute definitions.
func (r *Router) RegisterRoutes(routes []StaticRoute) {
	for _, route := range routes {
		if route.Handler != nil {
			r.routes = append(r.routes, Route{Path: route.Path, Handler: route.Handler})
		}
	}
}

// NotFound sets the handler for unmatched routes
func (r *Router) NotFound(handler func()) {
	r.notFound = handler
}

// BeforeNavigate sets a callback that runs before each navigation
func (r *Router) BeforeNavigate(fn func(from, to string) bool) {
	r.beforeNav = fn
}

// CurrentPath returns a store containing the current path
func (r *Router) CurrentPath() *Store[string] {
	return r.currentPath
}

// Start initializes the router - for SSR, reads PRERENDER_PATH and calls matching handler
func (r *Router) Start() {
	path := os.Getenv("PRERENDER_PATH")
	if path == "" {
		path = "/"
	}
	r.handleRoute(path)
}

// handleRoute finds and executes the matching route handler
func (r *Router) handleRoute(path string) {
	// Normalize path
	if path == "" {
		path = "/"
	}
	if path != "/" && strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	r.currentPath.Set(path)

	// Find matching route (most specific first)
	var bestMatch *Route
	var bestParams map[string]string
	bestSpecificity := -1

	for i := range r.routes {
		route := &r.routes[i]
		params, specificity, ok := matchRouteSSR(route.Path, path)
		if ok && specificity > bestSpecificity {
			bestMatch = route
			bestParams = params
			bestSpecificity = specificity
		}
	}

	if bestMatch != nil {
		bestMatch.Handler(bestParams)
	} else if r.notFound != nil {
		r.notFound()
	}
}

// matchRouteSSR matches a path against a route pattern (SSR version)
func matchRouteSSR(pattern, path string) (map[string]string, int, bool) {
	params := make(map[string]string)

	if pattern == "/" {
		if path == "/" {
			return params, 100, true
		}
		return nil, 0, false
	}

	if pattern == "*" || pattern == "**" {
		return params, 1, true
	}

	if strings.HasPrefix(pattern, "*/") {
		suffix := pattern[2:]
		if path == "/"+suffix {
			return params, 2, true
		}
		pathSegs := strings.Split(strings.Trim(path, "/"), "/")
		if len(pathSegs) >= 2 && pathSegs[len(pathSegs)-1] == suffix {
			return params, 2, true
		}
		return nil, 0, false
	}

	patternSegs := strings.Split(strings.Trim(pattern, "/"), "/")
	pathSegs := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternSegs) != len(pathSegs) {
		return nil, 0, false
	}

	specificity := 0
	for i, seg := range patternSegs {
		if strings.HasPrefix(seg, ":") {
			paramName := seg[1:]
			params[paramName] = pathSegs[i]
			specificity += 5
		} else if seg == pathSegs[i] {
			specificity += 10
		} else {
			return nil, 0, false
		}
	}

	return params, specificity, true
}

// SetupLinks intercepts link clicks (no-op for SSR)
func (r *Router) SetupLinks() {}

// Navigate programmatically navigates to a path (no-op for SSR)
func (r *Router) Navigate(path string) {}

// Replace navigates without adding to history (no-op for SSR)
func (r *Router) Replace(path string) {}

// Navigate is a standalone function for simple navigation (no-op for SSR)
func Navigate(path string) {}

// Link sets up an anchor element for SPA navigation (no-op for SSR)
func Link(el jsValue, router *Router) {}
