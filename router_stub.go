//go:build !wasm

package preveltekit

// routerJS holds browser-only router state; nothing is needed in SSR.
type routerJS struct{}

// Start initializes the router and handles the current URL (from SSRPath)
func (r *Router) Start() {
	// Path set by Hydrate via SetSSRPath
	path := ssrPath
	if path == "" {
		path = "/"
	}
	// In SSR, the base path is always "/" since SSRPaths are root-relative
	r.basePath = "/"
	r.currentPath.Set(path)
	r.handleRoute(path)
}

// SetupLinks intercepts link clicks (no-op for SSR)
func (r *Router) SetupLinks() {}

// Navigate programmatically navigates to a path (no-op for SSR)
func (r *Router) Navigate(path string) {}

// Replace navigates without adding to history (no-op for SSR)
func (r *Router) Replace(path string) {}
