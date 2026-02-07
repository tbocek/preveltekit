//go:build !wasm

package preveltekit

// pendingRouterIDs collects router IDs during SSR for container mapping
var pendingRouterIDs []string

// Router handles client-side routing
type Router struct {
	componentStore *Store[Component]
	routes         []Route
	id             string
	notFound       func()
	currentPath    *Store[string]
	beforeNav      func(from, to string) bool
}

// NewRouter creates a new router instance and registers the ID for SSR.
// Automatically registers all route components as options on the component store
// so SSR can pre-render all branches.
func NewRouter(componentStore *Store[Component], routes []Route, id string) *Router {
	// Register ID for SSR to discover container mapping
	pendingRouterIDs = append(pendingRouterIDs, id)
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
		currentPath:    NewWithID(id+".path", ""),
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

// Start initializes the router and handles the current URL (from SSRPath)
func (r *Router) Start() {
	// Get path from fake js.Global (set via SetSSRPath)
	path := jsGlobalFunc().Get("location").Get("pathname").String()
	if path == "" {
		path = "/"
	}
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
	var bestMatch *Route
	bestSpecificity := -1

	for i := range r.routes {
		route := &r.routes[i]
		_, specificity, ok := matchRouteSSR(route.Path, path)
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

// matchRouteSSR matches a path against a route pattern (SSR version)
func matchRouteSSR(pattern, path string) (map[string]string, int, bool) {
	params := make(map[string]string)

	// Handle root path
	if pattern == "/" {
		if path == "/" {
			return params, 100, true
		}
		return nil, 0, false
	}

	// Handle catch-all pattern
	if pattern == "*" || pattern == "**" {
		return params, 1, true
	}

	// Standard segment-based matching
	patternSegs := splitPathSSR(pattern)
	pathSegs := splitPathSSR(path)

	if len(patternSegs) != len(pathSegs) {
		return nil, 0, false
	}

	specificity := 0
	for i, seg := range patternSegs {
		if len(seg) > 0 && seg[0] == ':' {
			// Parameter segment
			paramName := seg[1:]
			params[paramName] = pathSegs[i]
			specificity += 5
		} else if seg == pathSegs[i] {
			// Exact match
			specificity += 10
		} else {
			return nil, 0, false
		}
	}

	return params, specificity, true
}

// splitPathSSR splits a path into segments
func splitPathSSR(s string) []string {
	// Trim leading/trailing slashes
	start, end := 0, len(s)
	for start < end && s[start] == '/' {
		start++
	}
	for end > start && s[end-1] == '/' {
		end--
	}
	s = s[start:end]
	if s == "" {
		return nil
	}

	// Count segments
	n := 1
	for i := 0; i < len(s); i++ {
		if s[i] == '/' {
			n++
		}
	}

	// Split
	parts := make([]string, 0, n)
	segStart := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == '/' {
			if segStart < i {
				parts = append(parts, s[segStart:i])
			}
			segStart = i + 1
		}
	}
	return parts
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
func Link(el *jsValue, router *Router) {}
