//go:build !js || !wasm

package preveltekit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// HasRoutes is implemented by components that define routes.
type HasRoutes interface {
	Routes() []Route
}

// HasStyle is implemented by components that have CSS styles.
type HasStyle interface {
	Style() string
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

	// Collect all bindings (merged from all routes for single bindings.json)
	var allBindings CollectedBindings

	// Create output directory
	os.MkdirAll("dist", 0755)

	// Second pass: generate HTML for each SSR path with fresh state
	for _, route := range ssrPaths {
		// Set the SSR path before lifecycle methods
		SetSSRPath(route.SSRPath)

		// Create fresh app instance
		freshApp := app.(HasNew).New()

		// Call OnMount (creates router which reads path and sets component)
		if om, ok := freshApp.(HasOnMount); ok {
			om.OnMount()
		}

		// Get the registered router IDs (from OnMount -> NewRouter)
		routerIDs := GetPendingRouterIDs()
		routerID := ""
		if len(routerIDs) > 0 {
			routerID = routerIDs[0]
		}

		// Get app style
		var appStyle string
		if hs, ok := freshApp.(HasStyle); ok {
			appStyle = hs.Style()
		}

		// Build app store map
		appStoreMap := buildStoreMap(freshApp, "component")

		// Render the full tree - router already set the correct component
		result := RenderHTMLWithContextFull(freshApp.Render(),
			WithRouteGroupIDCtx(routerID),
		)

		// Resolve bindings (needed for If-blocks, Each-blocks, AttrConds)
		resolveBindings(result.Bindings, appStoreMap, "component", freshApp)

		// Build full HTML document
		fullHTML := buildHTMLDocument(result.HTML, appStyle, "", result.CollectedStyles)

		// Write HTML file
		htmlPath := filepath.Join("dist", route.HTMLFile)
		os.WriteFile(htmlPath, []byte(fullHTML), 0644)
		fmt.Fprintf(os.Stderr, "Generated: %s\n", htmlPath)

		// Merge bindings for this route into allBindings
		mergeBindings(&allBindings, result.Bindings)
	}

	// Output merged bindings as JSON
	bindingsJSON, _ := json.Marshal(allBindings)
	fmt.Fprintf(os.Stderr, "BINDINGS:%s\n", bindingsJSON)
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
		seenInput[b.ElementID] = true
	}
	for _, b := range src.InputBindings {
		if !seenInput[b.ElementID] {
			dst.InputBindings = append(dst.InputBindings, b)
			seenInput[b.ElementID] = true
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

	seenContainer := make(map[string]bool)
	for _, c := range dst.ComponentContainers {
		seenContainer[c.ID] = true
	}
	for _, c := range src.ComponentContainers {
		if !seenContainer[c.ID] {
			dst.ComponentContainers = append(dst.ComponentContainers, c)
			seenContainer[c.ID] = true
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

// buildStoreMap builds a map from store pointer addresses to field paths.
// DEPRECATED: Use getStoreID instead to get store IDs directly.
func buildStoreMap(comp Component, prefix string) map[uintptr]string {
	storeMap := make(map[uintptr]string)
	rv := reflect.ValueOf(comp).Elem()
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		f := rv.Field(i)
		if f.Kind() == reflect.Ptr && !f.IsNil() {
			storeMap[f.Pointer()] = prefix + "." + rt.Field(i).Name
		}
	}
	return storeMap
}

// componentVarName returns the variable name for a component (lowercase first letter of type name).
func componentVarName(comp Component) string {
	t := reflect.TypeOf(comp).Elem().Name()
	if len(t) == 0 {
		return "comp"
	}
	return strings.ToLower(t[:1]) + t[1:]
}

// resolveBindings resolves store references in bindings using store IDs directly.
// Stores have user-defined IDs set via New("myapp.Count", 0), so we use those directly.
func resolveBindings(bindings *CollectedBindings, storeMap map[uintptr]string, prefix string, comp Component) {
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
					resolveBindings(bindings.IfBlocks[i].Branches[j].Bindings, storeMap, prefix, comp)
				}
				continue
			}
			cond := bindings.IfBlocks[i].Branches[j].CondRef
			if cond != nil {
				if sc, ok := cond.(*StoreCondition); ok && sc.Store != nil {
					// Use store's ID directly instead of reflection
					storeID := getStoreID(sc.Store)
					if storeID != "" {
						operand := fmt.Sprintf("%v", sc.Operand)
						if !isNumeric(operand) && operand != "true" && operand != "false" {
							operand = `"` + operand + `"`
						}
						bindings.IfBlocks[i].Branches[j].CondExpr = storeID + ".Get() " + sc.Op + " " + operand
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
						bindings.IfBlocks[i].Branches[j].CondExpr = storeID + ".Get()"
						bindings.IfBlocks[i].Branches[j].StoreID = storeID
						bindings.IfBlocks[i].Branches[j].IsBool = true
						bindings.IfBlocks[i].Deps = append(bindings.IfBlocks[i].Deps, storeID)
					}
				}
			}

			if bindings.IfBlocks[i].Branches[j].Bindings != nil {
				resolveBindings(bindings.IfBlocks[i].Branches[j].Bindings, storeMap, prefix, comp)
			}
		}

		if bindings.IfBlocks[i].ElseBindings != nil {
			resolveBindings(bindings.IfBlocks[i].ElseBindings, storeMap, prefix, comp)
		}
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

	// Resolve each block list references - lists don't have IDs, use reflection fallback
	for i := range bindings.EachBlocks {
		if bindings.EachBlocks[i].ListID != "" {
			continue
		}
		if bindings.EachBlocks[i].ListRef != nil {
			addr := reflect.ValueOf(bindings.EachBlocks[i].ListRef).Pointer()
			if name, ok := storeMap[addr]; ok {
				bindings.EachBlocks[i].ListID = name
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

// isNumeric checks if a string represents a number.
func isNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i, c := range s {
		if c == '-' && i == 0 {
			continue
		}
		if c == '.' {
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func buildHTMLDocument(body, appStyle, childStyle string, collectedStyles map[string]string) string {
	var styles string
	if appStyle != "" {
		styles += "<style>" + appStyle + "</style>\n"
	}
	if childStyle != "" {
		styles += "<style>" + childStyle + "</style>\n"
	}
	// Add auto-collected styles from nested components
	if len(collectedStyles) > 0 {
		var nestedStyles string
		for _, style := range collectedStyles {
			nestedStyles += style + "\n"
		}
		if nestedStyles != "" {
			styles += "<style>" + nestedStyles + "</style>\n"
		}
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
fetch("bindings.json")
  .then(r => r.json())
  .then(bindings => {
    window._preveltekit_bindings = JSON.stringify(bindings);
    const go = new Go();
    return WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
      .then(result => go.run(result.instance));
  });
</script>
</body>
</html>`, styles, body)
}
