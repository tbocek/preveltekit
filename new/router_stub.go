//go:build !wasm

package preveltekit

// pendingRouterIDs collects router IDs during SSR for container mapping
var pendingRouterIDs []string

// Router handles client-side routing (stub for SSR - no-op)
type Router struct {
	routes      []Route
	id          string
	notFound    func()
	currentPath *Store[string]
	beforeNav   func(from, to string) bool
}

// NewRouter creates a new router instance and registers the ID for SSR
func NewRouter(componentStore *Store[Component], routes []Route, id string) *Router {
	// Register ID for SSR to discover container mapping
	pendingRouterIDs = append(pendingRouterIDs, id)
	return &Router{
		routes:      routes,
		id:          id,
		currentPath: New(""),
	}
}

// GetPendingRouterIDs returns and clears the pending router IDs
func GetPendingRouterIDs() []string {
	ids := pendingRouterIDs
	pendingRouterIDs = nil
	return ids
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
