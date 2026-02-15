//go:build !wasm

package preveltekit

// Router handles client-side routing
type Router struct {
	componentStore *Store[Component]
	routes         []Route
	id             string
	basePath       string
	notFound       func()
	currentPath    *Store[string]
	beforeNav      func(from, to string) bool
}

// NewRouter creates a new router instance and registers the ID for SSR.
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
		id:             id,
		currentPath:    newWithID(id+".path", ""),
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

// Start initializes the router and handles the current URL (from SSRPath)
func (r *Router) Start() {
	// Get path from fake js.Global (set via SetSSRPath)
	path := jsGlobalFunc().Get("location").Get("pathname").String()
	if path == "" {
		path = "/"
	}
	// In SSR, the base path is always "/" since SSRPaths are root-relative
	r.basePath = "/"
	r.currentPath.Set(path)
	r.handleRoute(path)
}

// handleRoute matches the path and sets the component
func (r *Router) handleRoute(path string) {
	// Normalize path
	if path == "" {
		path = "/"
	}
	if path != "/" && len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

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

// SetupLinks intercepts link clicks (no-op for SSR)
func (r *Router) SetupLinks() {}

// Navigate programmatically navigates to a path (no-op for SSR)
func (r *Router) Navigate(path string) {}

// Replace navigates without adding to history (no-op for SSR)
func (r *Router) Replace(path string) {}
