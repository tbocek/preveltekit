package preveltekit

import "strings"

// Router handles client-side routing. In SSR it only resolves the path being
// pre-rendered; navigation and link interception run in WASM (router.go).
type Router struct {
	componentStore *Store[Component]
	routes         []Route
	basePath       string // detected at Start(), used to resolve relative route paths
	notFound       func()
	currentPath    *Store[string]
	beforeNav      func(from, to string) bool // return false to cancel navigation
	routerJS                                  // browser state (empty in SSR)
}

// NewRouter creates a new router instance with a component store, routes, and ID.
// Automatically registers all route components as options on the component store
// so SSR can pre-render all branches.
func NewRouter(componentStore *Store[Component], routes []Route, id string) *Router {
	// Register all route components as store options for pre-baked rendering
	for _, route := range routes {
		if route.Component != nil {
			componentStore.WithOptions(route.Component)
		}
	}
	return &Router{
		componentStore: componentStore,
		routes:         routes,
		currentPath:    newWithID(id+".path", ""),
	}
}

// NotFound sets the handler for unmatched routes
func (r *Router) NotFound(handler func()) {
	r.notFound = handler
}

// BeforeNavigate sets a callback that runs before each navigation
// Return false to cancel the navigation
func (r *Router) BeforeNavigate(fn func(from, to string) bool) {
	r.beforeNav = fn
}

// CurrentPath returns a store containing the current path
func (r *Router) CurrentPath() *Store[string] {
	return r.currentPath
}

// handleRoute matches the path and sets the component of the most specific route.
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
	// Each route's Path is resolved against the base path before matching.
	var bestMatch *Route
	bestSpecificity := -1

	for i := range r.routes {
		route := &r.routes[i]
		resolved := resolveRoute(r.basePath, route.Path)
		_, specificity, ok := matchRoute(resolved, path)
		if ok && specificity > bestSpecificity {
			bestMatch = route
			bestSpecificity = specificity
		}
	}

	if bestMatch != nil && bestMatch.Component != nil {
		r.componentStore.Set(bestMatch.Component)
	} else if r.notFound != nil {
		r.notFound()
	}
}
