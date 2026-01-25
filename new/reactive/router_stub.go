//go:build !wasm

package reactive

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

// Start initializes the router (no-op for SSR)
func (r *Router) Start() {}

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
