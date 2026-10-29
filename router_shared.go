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
	id             string
	params         []routeParam      // stores returned by Param
	currentParams  map[string]string // parameters of the current route
	routerJS                         // browser state (empty in SSR)
}

// routeParam is a store returned by Param for a parameter name.
type routeParam struct {
	name  string
	store *Store[string]
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
		id:             id,
	}
}

// Param returns a store with the value of the route parameter name in the
// current path, e.g. "42" for the route "/users/:id" at "/users/42" ("" if
// the current route has no such parameter). The router updates it before
// it shows the route's component. When only the parameter changes, the
// component stays the same: bind the store or use OnChange to update it.
func (r *Router) Param(name string) *Store[string] {
	for _, p := range r.params {
		if p.name == name {
			return p.store
		}
	}
	s := newWithID(r.id+".param."+name, r.currentParams[name])
	r.params = append(r.params, routeParam{name: name, store: s})
	return s
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
	var bestParams map[string]string
	bestSpecificity := -1

	for i := range r.routes {
		route := &r.routes[i]
		resolved := resolveRoute(r.basePath, route.Path)
		params, specificity, ok := matchRoute(resolved, path)
		if ok && specificity > bestSpecificity {
			bestMatch = route
			bestParams = params
			bestSpecificity = specificity
		}
	}

	// Update the parameter stores before the component is shown
	r.currentParams = bestParams
	for _, p := range r.params {
		if v := bestParams[p.name]; p.store.Get() != v {
			p.store.Set(v)
		}
	}

	if bestMatch != nil && bestMatch.Component != nil {
		r.componentStore.Set(bestMatch.Component)
	} else if r.notFound != nil {
		r.notFound()
	}
}
