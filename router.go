//go:build wasm

package preveltekit

import (
	"strings"
	"syscall/js"
)

// Router handles client-side routing
type Router struct {
	componentStore *Store[Component]
	routes         []Route
	id             string
	basePath       string // detected at Start(), used to resolve relative route paths
	notFound       func()
	currentPath    *Store[string]
	beforeNav      func(from, to string) bool // return false to cancel navigation
	linksSetup     bool                       // tracks if click listener is already registered
	clickFn        js.Func                    // retained to prevent GC
	popstateFn     js.Func                    // retained to prevent GC
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
		id:             id,
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

// detectBasePath determines the base path by matching the current pathname
// against route SSRPaths. E.g., if pathname is "/preveltekit/manual" and a
// route has SSRPath "/manual", the base path is "/preveltekit".
func (r *Router) detectBasePath(pathname string) string {
	// Normalize: strip trailing slash for matching (unless root)
	norm := pathname
	if len(norm) > 1 && norm[len(norm)-1] == '/' {
		norm = norm[:len(norm)-1]
	}

	for _, route := range r.routes {
		ssr := route.SSRPath
		if ssr == "" {
			continue
		}
		if ssr == "/" {
			// Root route — base could be the entire pathname
			// Only use this if no other route matches more specifically
			continue
		}
		if strings.HasSuffix(norm, ssr) {
			base := norm[:len(norm)-len(ssr)]
			if base == "" {
				return "/"
			}
			return base
		}
	}
	// Fallback: the whole pathname is the base (root route matched)
	return norm
}

// Start initializes the router and handles the current URL
func (r *Router) Start() {
	// Handle initial route
	path := js.Global().Get("location").Get("pathname").String()
	r.basePath = r.detectBasePath(path)
	r.currentPath.Set(path)
	r.handleRoute(path)

	// Listen for popstate (back/forward)
	r.popstateFn = js.FuncOf(func(this js.Value, args []js.Value) any {
		path := js.Global().Get("location").Get("pathname").String()
		r.handleRoute(path)
		return nil
	})
	js.Global().Call("addEventListener", "popstate", r.popstateFn)

	// Intercept all link clicks for SPA navigation
	r.SetupLinks()
}

// SetupLinks intercepts clicks on all internal anchor elements for SPA navigation
// This is called automatically by Start(). Safe to call multiple times.
func (r *Router) SetupLinks() {
	if r.linksSetup {
		return
	}
	r.linksSetup = true

	r.clickFn = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		e := args[0]

		// Only handle left-click without modifiers
		if e.Get("button").Int() != 0 {
			return nil
		}
		if e.Get("ctrlKey").Bool() || e.Get("metaKey").Bool() ||
			e.Get("altKey").Bool() || e.Get("shiftKey").Bool() {
			return nil
		}

		// Find the anchor element (could be the target or a parent)
		target := e.Get("target")
		var anchor js.Value
		for !target.IsNull() && !target.IsUndefined() {
			tagName := target.Get("tagName")
			if !tagName.IsUndefined() && tagName.String() == "A" {
				anchor = target
				break
			}
			target = target.Get("parentElement")
		}

		if anchor.IsUndefined() || anchor.IsNull() {
			return nil
		}

		href := anchor.Call("getAttribute", "href")
		if href.IsNull() || href.IsUndefined() {
			return nil
		}
		hrefStr := href.String()

		// Skip external links
		if strings.HasPrefix(hrefStr, "http://") || strings.HasPrefix(hrefStr, "https://") ||
			strings.HasPrefix(hrefStr, "//") || strings.HasPrefix(hrefStr, "mailto:") ||
			strings.HasPrefix(hrefStr, "tel:") {
			return nil
		}

		// Skip links with external attribute
		if ext := anchor.Call("getAttribute", "external"); !ext.IsNull() {
			return nil
		}

		// Skip links with target="_blank"
		if tgt := anchor.Call("getAttribute", "target"); !tgt.IsNull() && tgt.String() == "_blank" {
			return nil
		}

		// Skip hash-only links
		if hrefStr == "#" || (strings.HasPrefix(hrefStr, "#") && len(hrefStr) > 1) {
			return nil
		}

		e.Call("preventDefault")

		// Resolve and navigate
		path := resolvePath(hrefStr)
		r.Navigate(path)

		return nil
	})
	js.Global().Get("document").Call("addEventListener", "click", r.clickFn)
}

// Navigate programmatically navigates to a path
func (r *Router) Navigate(path string) {
	currentPath := r.currentPath.Get()

	// Check beforeNav hook
	if r.beforeNav != nil && !r.beforeNav(currentPath, path) {
		return
	}

	js.Global().Get("history").Call("pushState", nil, "", path)
	r.handleRoute(path)
}

// Replace navigates without adding to history
func (r *Router) Replace(path string) {
	currentPath := r.currentPath.Get()

	if r.beforeNav != nil && !r.beforeNav(currentPath, path) {
		return
	}

	// Use replaceState then navigate to avoid adding to history
	js.Global().Get("history").Call("replaceState", nil, "", path)
	js.Global().Get("location").Call("replace", path)
}

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

// resolvePath resolves a relative or absolute href to an absolute path
func resolvePath(href string) string {
	if len(href) > 0 && href[0] == '/' {
		return href
	}

	if href == "" || href == "#" {
		return js.Global().Get("location").Get("pathname").String()
	}

	// Get current path
	current := js.Global().Get("location").Get("pathname").String()
	if !strings.HasSuffix(current, "/") {
		// Remove last segment for relative resolution
		if idx := strings.LastIndexByte(current, '/'); idx >= 0 {
			current = current[:idx+1]
		}
	}

	path := current + href

	// Preserve trailing slash (matters for relative resolution)
	trailingSlash := strings.HasSuffix(path, "/")

	// Clean up . and .. segments
	segments := strings.Split(path, "/")
	var clean []string
	for _, seg := range segments {
		if seg == ".." {
			if len(clean) > 0 {
				clean = clean[:len(clean)-1]
			}
		} else if seg != "" && seg != "." {
			clean = append(clean, seg)
		}
	}
	path = "/" + strings.Join(clean, "/")
	if trailingSlash && !strings.HasSuffix(path, "/") {
		path += "/"
	}

	return path
}

func cleanDoubleSlash(s string) string {
	if !strings.Contains(s, "//") {
		return s
	}
	for strings.Contains(s, "//") {
		s = strings.ReplaceAll(s, "//", "/")
	}
	return s
}
