//go:build !js || !wasm

package preveltekit

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// HasRoutes is implemented by components that define routes.
type HasRoutes interface {
	Routes() []Route
}

// HasStyle is implemented by components that have scoped CSS styles.
type HasStyle interface {
	Style() string
}

// HasGlobalStyle is implemented by components that have unscoped global CSS styles.
// Global styles are emitted without any CSS scoping — useful for base/reset styles.
type HasGlobalStyle interface {
	GlobalStyle() string
}

// HasNew is implemented by components that can create fresh instances.
type HasNew interface {
	New() Component
}

// HasOnMount is implemented by components with OnMount lifecycle.
type HasOnMount interface {
	OnMount()
}

// Hydrate is the main entry point for declarative components.
// In SSR mode (native build), it generates HTML files and outputs bindings JSON.
// The ID-based system is used for stores and handlers, but bindings JSON is still
// needed for If-blocks, Each-blocks, and AttrCond bindings.
func Hydrate(app Component) {
	// First pass: discover all SSR paths
	// Create fresh instance to get routes
	if hn, ok := app.(HasNew); ok {
		app = hn.New()
	}

	var ssrPaths []Route
	if hr, ok := app.(HasRoutes); ok {
		for _, route := range hr.Routes() {
			if route.SSRPath != "" {
				ssrPaths = append(ssrPaths, route)
			}
		}
	}

	// Collect all bindings (merged from all routes for single bindings.bin)
	var allBindings CollectedBindings

	// Create output directory
	os.MkdirAll("dist", 0755)

	// Second pass: generate HTML for each SSR path with fresh state
	var savedStringsToRemove map[string]int
	for ri, route := range ssrPaths {
		// Reset global counters so each iteration starts from s0,
		// matching the single app.New() call in WASM.
		resetRegistries()

		// Only collect Html() strings on first iteration — one pass through
		// all component Render() methods is enough to capture all source-level
		// string literals. Multiple iterations would over-count.
		if ri == 1 {
			savedStringsToRemove = wasmStringsToRemove
			wasmStringsToRemove = nil
		}

		// Set the SSR path before lifecycle methods
		SetSSRPath(route.SSRPath)

		// Create fresh app instance
		var freshApp Component
		if hn, ok := app.(HasNew); ok {
			freshApp = hn.New()
		}

		// Call OnMount (creates router which reads path and sets component)
		if om, ok := freshApp.(HasOnMount); ok {
			om.OnMount()
		}

		// Render the full tree - router already set the correct component
		// Routes are now auto-registered as store options via NewRouter,
		// so no need to pass them through context.
		ctx := NewBuildContext()

		// Collect app global styles (unscoped)
		if hgs, ok := freshApp.(HasGlobalStyle); ok {
			if gs := hgs.GlobalStyle(); gs != "" {
				ctx.CollectedGlobalStyles["app"] = gs
			}
		}

		// Set app-level scope before rendering so all app HTML gets the class
		if hs, ok := freshApp.(HasStyle); ok {
			scopeAttr := GetOrCreateScope("app")
			ctx.ScopeAttr = scopeAttr
			ctx.CollectedStyles["app"] = scopeCSS(hs.Style(), scopeAttr)
		}

		html := nodeToHTML(freshApp.Render(), ctx)

		// Resolve store references to string IDs, then merge into shared bindings
		resolveBindings(ctx.Bindings)
		mergeBindings(&allBindings, ctx.Bindings)

		// Build full HTML document
		fullHTML := buildHTMLDocument(minifyHTML(html), ctx.CollectedGlobalStyles, ctx.CollectedStyles)

		// Write HTML file
		htmlPath := filepath.Join("dist", route.HTMLFile)
		os.WriteFile(htmlPath, []byte(fullHTML), 0644)
		fmt.Fprintf(os.Stderr, "Generated: %s\n", htmlPath)
	}

	// Write bindings as binary
	bindingsBin := encodeBindings(&allBindings)
	binPath := filepath.Join("dist", "bindings.bin")
	os.WriteFile(binPath, bindingsBin, 0644)
	fmt.Fprintf(os.Stderr, "Generated: %s\n", binPath)

	// Write strings to remove from WASM
	// Format: first line is expected count, rest is the raw string
	removeDir := filepath.Join("dist", "remove_from_wasm")
	os.RemoveAll(removeDir)
	os.MkdirAll(removeDir, 0755)
	// Use saved map (from first iteration) if available, else current
	if savedStringsToRemove == nil {
		savedStringsToRemove = wasmStringsToRemove
	}
	i := 0
	for s, count := range savedStringsToRemove {
		content := fmt.Sprintf("%d\n%s", count, s)
		os.WriteFile(filepath.Join(removeDir, fmt.Sprintf("%d", i)), []byte(content), 0644)
		i++
	}
	fmt.Fprintf(os.Stderr, "Generated: %d strings to remove from WASM\n", i)
}

// mergeBindings merges src bindings into dst, avoiding duplicates by marker/element ID.
func mergeBindings(dst, src *CollectedBindings) {
	// Track seen IDs to avoid duplicates
	seenText := make(map[string]bool)
	for _, b := range dst.TextBindings {
		seenText[b.MarkerID] = true
	}
	for _, b := range src.TextBindings {
		if !seenText[b.MarkerID] {
			dst.TextBindings = append(dst.TextBindings, b)
			seenText[b.MarkerID] = true
		}
	}

	seenInput := make(map[string]bool)
	for _, b := range dst.InputBindings {
		seenInput[b.StoreID] = true
	}
	for _, b := range src.InputBindings {
		if !seenInput[b.StoreID] {
			dst.InputBindings = append(dst.InputBindings, b)
			seenInput[b.StoreID] = true
		}
	}

	seenEvent := make(map[string]bool)
	for _, b := range dst.Events {
		seenEvent[b.ElementID+":"+b.Event] = true
	}
	for _, b := range src.Events {
		key := b.ElementID + ":" + b.Event
		if !seenEvent[key] {
			dst.Events = append(dst.Events, b)
			seenEvent[key] = true
		}
	}

	seenIf := make(map[string]bool)
	for _, b := range dst.IfBlocks {
		seenIf[b.MarkerID] = true
	}
	for _, b := range src.IfBlocks {
		if !seenIf[b.MarkerID] {
			dst.IfBlocks = append(dst.IfBlocks, b)
			seenIf[b.MarkerID] = true
		}
	}

	seenAttr := make(map[string]bool)
	for _, b := range dst.AttrBindings {
		seenAttr[b.ElementID] = true
	}
	for _, b := range src.AttrBindings {
		if !seenAttr[b.ElementID] {
			dst.AttrBindings = append(dst.AttrBindings, b)
			seenAttr[b.ElementID] = true
		}
	}

	seenEach := make(map[string]bool)
	for _, b := range dst.EachBlocks {
		seenEach[b.MarkerID] = true
	}
	for _, b := range src.EachBlocks {
		if !seenEach[b.MarkerID] {
			dst.EachBlocks = append(dst.EachBlocks, b)
			seenEach[b.MarkerID] = true
		}
	}

	seenAttrCond := make(map[string]bool)
	for _, b := range dst.AttrCondBindings {
		seenAttrCond[b.ElementID+":"+b.AttrName] = true
	}
	for _, b := range src.AttrCondBindings {
		key := b.ElementID + ":" + b.AttrName
		if !seenAttrCond[key] {
			dst.AttrCondBindings = append(dst.AttrCondBindings, b)
			seenAttrCond[key] = true
		}
	}

	seenCompBlock := make(map[string]bool)
	for _, cb := range dst.ComponentBlocks {
		seenCompBlock[cb.MarkerID] = true
	}
	for _, cb := range src.ComponentBlocks {
		if !seenCompBlock[cb.MarkerID] {
			dst.ComponentBlocks = append(dst.ComponentBlocks, cb)
			seenCompBlock[cb.MarkerID] = true
		}
	}
}

// getStoreID extracts the ID from any store type using the HasID interface.
// Returns empty string if the value is not a store or has no ID.
func getStoreID(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(HasID); ok {
		return s.ID()
	}
	return ""
}

// resolveBindings resolves store references in bindings using store IDs directly.
// Stores have auto-generated IDs from New(), so we use those directly.
func resolveBindings(bindings *CollectedBindings) {
	// Resolve text bindings - use store's ID directly
	for i := range bindings.TextBindings {
		if bindings.TextBindings[i].StoreID != "" {
			continue
		}
		if id := getStoreID(bindings.TextBindings[i].StoreRef); id != "" {
			bindings.TextBindings[i].StoreID = id
		}
	}

	// Resolve input bindings - use store's ID directly
	for i := range bindings.InputBindings {
		if bindings.InputBindings[i].StoreID != "" {
			continue
		}
		if id := getStoreID(bindings.InputBindings[i].StoreRef); id != "" {
			bindings.InputBindings[i].StoreID = id
		}
	}

	// Resolve if-block conditions using store IDs directly
	for i := range bindings.IfBlocks {
		for j := range bindings.IfBlocks[i].Branches {
			// Skip if already resolved
			if bindings.IfBlocks[i].Branches[j].StoreID != "" {
				if bindings.IfBlocks[i].Branches[j].Bindings != nil {
					resolveBindings(bindings.IfBlocks[i].Branches[j].Bindings)
				}
				continue
			}
			cond := bindings.IfBlocks[i].Branches[j].CondRef
			if cond != nil {
				if sc, ok := cond.(*StoreCondition); ok && sc.Store != nil {
					// Use store's ID directly instead of reflection
					storeID := getStoreID(sc.Store)
					if storeID != "" {
						bindings.IfBlocks[i].Branches[j].StoreID = storeID
						bindings.IfBlocks[i].Branches[j].Op = sc.Op
						bindings.IfBlocks[i].Branches[j].Operand = fmt.Sprintf("%v", sc.Operand)
						bindings.IfBlocks[i].Deps = append(bindings.IfBlocks[i].Deps, storeID)
					}
				}

				if bc, ok := cond.(*BoolCondition); ok && bc.Store != nil {
					// Use store's ID directly instead of reflection
					storeID := getStoreID(bc.Store)
					if storeID != "" {
						bindings.IfBlocks[i].Branches[j].StoreID = storeID
						bindings.IfBlocks[i].Branches[j].IsBool = true
						bindings.IfBlocks[i].Deps = append(bindings.IfBlocks[i].Deps, storeID)
					}
				}
			}

			if bindings.IfBlocks[i].Branches[j].Bindings != nil {
				resolveBindings(bindings.IfBlocks[i].Branches[j].Bindings)
			}
		}

		if bindings.IfBlocks[i].ElseBindings != nil {
			resolveBindings(bindings.IfBlocks[i].ElseBindings)
		}

		// Deduplicate deps
		seen := make(map[string]bool)
		unique := bindings.IfBlocks[i].Deps[:0]
		for _, d := range bindings.IfBlocks[i].Deps {
			if !seen[d] {
				seen[d] = true
				unique = append(unique, d)
			}
		}
		bindings.IfBlocks[i].Deps = unique
	}

	// Resolve attr bindings - use store IDs directly
	for i := range bindings.AttrBindings {
		if len(bindings.AttrBindings[i].StoreIDs) > 0 {
			continue
		}
		var storeIDs []string
		for _, storeRef := range bindings.AttrBindings[i].StoreRefs {
			if id := getStoreID(storeRef); id != "" {
				storeIDs = append(storeIDs, id)
			}
		}
		bindings.AttrBindings[i].StoreIDs = storeIDs
	}

	// Resolve each block list references - use list's ID directly
	for i := range bindings.EachBlocks {
		if bindings.EachBlocks[i].ListID != "" {
			continue
		}
		if bindings.EachBlocks[i].ListRef != nil {
			if id := getStoreID(bindings.EachBlocks[i].ListRef); id != "" {
				bindings.EachBlocks[i].ListID = id
			}
		}
	}

	// Resolve AttrCond bindings - use store IDs directly
	for i := range bindings.AttrCondBindings {
		if len(bindings.AttrCondBindings[i].Deps) > 0 {
			continue
		}

		var deps []string

		if id := getStoreID(bindings.AttrCondBindings[i].CondStoreRef); id != "" {
			deps = append(deps, id)
		}

		if id := getStoreID(bindings.AttrCondBindings[i].TrueStoreRef); id != "" {
			bindings.AttrCondBindings[i].TrueStoreID = id
			deps = append(deps, id)
		}

		if id := getStoreID(bindings.AttrCondBindings[i].FalseStoreRef); id != "" {
			bindings.AttrCondBindings[i].FalseStoreID = id
			deps = append(deps, id)
		}

		bindings.AttrCondBindings[i].Deps = deps
	}
}

func buildHTMLDocument(body string, collectedGlobalStyles, collectedStyles map[string]string) string {
	var styles string

	var allStyles string

	// Global styles first (unscoped)
	if len(collectedGlobalStyles) > 0 {
		keys := make([]string, 0, len(collectedGlobalStyles))
		for k := range collectedGlobalStyles {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			allStyles += collectedGlobalStyles[k] + "\n"
		}
	}

	// Scoped styles
	if len(collectedStyles) > 0 {
		keys := make([]string, 0, len(collectedStyles))
		for k := range collectedStyles {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			allStyles += collectedStyles[k] + "\n"
		}
	}

	if allStyles != "" {
		styles = "<style>" + minifyCSS(allStyles) + "</style>\n"
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
%s</head>
<body>
%s
<script src="wasm_exec.js"></script>
<script>
fetch("bindings.bin")
  .then(r => r.arrayBuffer())
  .then(buf => {
    window._preveltekit_bindings = new Uint8Array(buf);
    const go = new Go();
    return WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
      .then(result => go.run(result.instance));
  });
</script>
</body>
</html>`, styles, body)
}
