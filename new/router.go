//go:build wasm

package preveltekit

import (
	"syscall/js"
)

// Router handles client-side routing
type Router struct {
	componentStore *Store[Component]
	routes         []Route
	id             string
	notFound       func()
	currentPath    *Store[string]
	beforeNav      func(from, to string) bool // return false to cancel navigation
	linksSetup     bool                       // tracks if click listener is already registered
	clickFn        js.Func                    // retained to prevent GC
	popstateFn     js.Func                    // retained to prevent GC
}

// NewRouter creates a new router instance with a component store, routes, and ID
func NewRouter(componentStore *Store[Component], routes []Route, id string) *Router {
	return &Router{
		componentStore: componentStore,
		routes:         routes,
		id:             id,
		currentPath:    newInternal(""),
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

// Start initializes the router and handles the current URL
func (r *Router) Start() {
	js.Global().Get("console").Call("log", "[DEBUG] Router.Start called")
	js.Global().Get("console").Call("log", "[DEBUG] componentStore is nil:", r.componentStore == nil)
	js.Global().Get("console").Call("log", "[DEBUG] routes count:", len(r.routes))

	// Bind component store to DOM container for reactive updates
	containerID := componentContainers[r.id]
	if containerID == "" {
		containerID = "component-root" // fallback
	}
	js.Global().Get("console").Call("log", "[DEBUG] binding router", r.id, "to container:", containerID)
	bindComponentStore(r.componentStore, containerID)

	// Handle initial route
	path := js.Global().Get("location").Get("pathname").String()
	js.Global().Get("console").Call("log", "[DEBUG] initial path:", path)
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
		if hasPrefix(hrefStr, "http://") || hasPrefix(hrefStr, "https://") ||
			hasPrefix(hrefStr, "//") || hasPrefix(hrefStr, "mailto:") ||
			hasPrefix(hrefStr, "tel:") {
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
		if hrefStr == "#" || (hasPrefix(hrefStr, "#") && len(hrefStr) > 1) {
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
	js.Global().Get("console").Call("log", "[DEBUG] handleRoute called with path:", path)

	// Normalize path
	if path == "" {
		path = "/"
	}
	if path != "/" && hasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	r.currentPath.Set(path)

	// Find matching route (most specific first)
	var bestMatch *Route
	bestSpecificity := -1

	for i := range r.routes {
		route := &r.routes[i]
		_, specificity, ok := matchRoute(route.Path, path)
		js.Global().Get("console").Call("log", "[DEBUG] checking route:", route.Path, "specificity:", specificity, "ok:", ok)
		if ok && specificity > bestSpecificity {
			bestMatch = route
			bestSpecificity = specificity
		}
	}

	js.Global().Get("console").Call("log", "[DEBUG] bestMatch found:", bestMatch != nil)
	if bestMatch != nil {
		js.Global().Get("console").Call("log", "[DEBUG] bestMatch.Path:", bestMatch.Path)
		js.Global().Get("console").Call("log", "[DEBUG] bestMatch.Component is nil:", bestMatch.Component == nil)
	}

	if bestMatch != nil && bestMatch.Component != nil {
		js.Global().Get("console").Call("log", "[DEBUG] setting componentStore")
		r.componentStore.Set(bestMatch.Component)
	} else if r.notFound != nil {
		r.notFound()
	}
}

// matchRoute matches a path against a route pattern
// Returns params, specificity score, and whether it matched
func matchRoute(pattern, path string) (map[string]string, int, bool) {
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

	// Handle wildcard prefix patterns like */suffix
	if hasPrefix(pattern, "*/") {
		suffix := pattern[2:]
		if path == "/"+suffix {
			return params, 2, true
		}
		// Match /{segment}/{suffix}
		pathSegs := splitPath(path)
		if len(pathSegs) >= 2 && pathSegs[len(pathSegs)-1] == suffix {
			return params, 2, true
		}
		return nil, 0, false
	}

	// Standard segment-based matching
	patternSegs := splitPath(pattern)
	pathSegs := splitPath(path)

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
	if !hasSuffix(current, "/") {
		// Remove last segment for relative resolution
		if idx := lastIndexByte(current, '/'); idx >= 0 {
			current = current[:idx+1]
		}
	}

	path := current + href

	// Clean up ../ segments
	if containsDotDot(path) {
		segments := splitPathAll(path)
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
		path = "/" + joinPath(clean)
	}

	// Clean double slashes
	path = cleanDoubleSlash(path)

	return path
}

// --- Inline string helpers (avoid strings package) ---

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func trimSlashes(s string) string {
	start, end := 0, len(s)
	for start < end && s[start] == '/' {
		start++
	}
	for end > start && s[end-1] == '/' {
		end--
	}
	return s[start:end]
}

func splitPath(s string) []string {
	s = trimSlashes(s)
	if s == "" {
		return nil
	}
	n := 1
	for i := 0; i < len(s); i++ {
		if s[i] == '/' {
			n++
		}
	}
	parts := make([]string, 0, n)
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == '/' {
			if start < i {
				parts = append(parts, s[start:i])
			}
			start = i + 1
		}
	}
	return parts
}

func splitPathAll(s string) []string {
	if s == "" {
		return nil
	}
	n := 1
	for i := 0; i < len(s); i++ {
		if s[i] == '/' {
			n++
		}
	}
	parts := make([]string, 0, n)
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == '/' {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	return parts
}

func joinPath(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	n := len(parts) - 1
	for _, p := range parts {
		n += len(p)
	}
	b := make([]byte, 0, n)
	for i, p := range parts {
		b = append(b, p...)
		if i < len(parts)-1 {
			b = append(b, '/')
		}
	}
	return string(b)
}

func lastIndexByte(s string, c byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func containsDotDot(s string) bool {
	for i := 0; i+1 < len(s); i++ {
		if s[i] == '.' && s[i+1] == '.' {
			if i+2 >= len(s) || s[i+2] == '/' {
				return true
			}
		}
	}
	return false
}

func cleanDoubleSlash(s string) string {
	hasDouble := false
	for i := 0; i+1 < len(s); i++ {
		if s[i] == '/' && s[i+1] == '/' {
			hasDouble = true
			break
		}
	}
	if !hasDouble {
		return s
	}
	b := make([]byte, 0, len(s))
	prev := byte(0)
	for i := 0; i < len(s); i++ {
		if s[i] == '/' && prev == '/' {
			continue
		}
		b = append(b, s[i])
		prev = s[i]
	}
	return string(b)
}
