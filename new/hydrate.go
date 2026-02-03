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

// Component is the interface that all declarative components must implement.
type Component interface {
	Render() Node
}

// HasRoutes is implemented by components that define routes.
type HasRoutes interface {
	Routes() []StaticRoute
}

// HasStyle is implemented by components that have CSS styles.
type HasStyle interface {
	Style() string
}

// HasOnCreate is implemented by components with OnCreate lifecycle.
type HasOnCreate interface {
	OnCreate()
}

// HasOnMount is implemented by components with OnMount lifecycle.
type HasOnMount interface {
	OnMount()
}

// HasCurrentComponent is implemented by apps that have a CurrentComponent store for routing.
type HasCurrentComponent interface {
	GetCurrentComponent() *Store[Component]
}

// Hydrate is the main entry point for declarative components.
// In SSR mode (native build), it generates HTML files and outputs bindings.
func Hydrate(app Component) {
	children := make(map[string]Component)

	// Initialize app stores first
	initStores(app)

	// Call OnCreate to let app initialize its children
	if oc, ok := app.(HasOnCreate); ok {
		oc.OnCreate()
	}

	// Auto-discover children from Routes() by calling handlers and reading CurrentComponent
	if hr, ok := app.(HasRoutes); ok {
		if hcc, ok := app.(HasCurrentComponent); ok {
			currentStore := hcc.GetCurrentComponent()
			routes := hr.Routes()
			for _, route := range routes {
				if route.Handler != nil {
					// Call handler to set CurrentComponent
					route.Handler(nil)
					// Read the component that was set
					if comp := currentStore.Get(); comp != nil {
						children[route.Path] = comp
					}
				}
			}
			// Reset CurrentComponent to the default route (first route, typically "/" or "/basics")
			// so that SSR renders the correct initial component
			if len(routes) > 0 && routes[0].Handler != nil {
				routes[0].Handler(nil)
			}
		}
	}

	// SSR mode - generate HTML and bindings
	hydrateGenerateAll(app, children)
}

// hydrateGenerateAll generates HTML files for all routes.
func hydrateGenerateAll(app Component, children map[string]Component) {
	// Deduplicate children (multiple routes may point to same component, e.g., "/" and "/basics")
	// Use pointer address to detect duplicates
	seenChildren := make(map[uintptr]bool)
	uniqueChildren := make([]Component, 0)
	for _, child := range children {
		ptr := reflect.ValueOf(child).Pointer()
		if !seenChildren[ptr] {
			seenChildren[ptr] = true
			uniqueChildren = append(uniqueChildren, child)
		}
	}

	// Initialize stores and call OnCreate/OnMount for unique children only
	// (app was already initialized in Hydrate before Routes() was called)
	for _, child := range uniqueChildren {
		initStores(child)
	}
	for _, child := range uniqueChildren {
		if oc, ok := child.(HasOnCreate); ok {
			oc.OnCreate()
		}
		if om, ok := child.(HasOnMount); ok {
			om.OnMount()
		}
	}

	// Get app style
	var appStyle string
	if hs, ok := app.(HasStyle); ok {
		appStyle = hs.Style()
	}

	// Build store maps
	appStoreMap := buildStoreMap(app, "component")
	childStoreMaps := make(map[string]map[uintptr]string)
	for path, child := range children {
		name := strings.TrimPrefix(path, "/")
		childStoreMaps[name] = buildStoreMap(child, name)
	}

	// Collect all bindings (merged from all routes for single bindings.json)
	var allBindings CollectedBindings

	// Create output directory
	os.MkdirAll("dist", 0755)

	// Pre-render all children with unique prefixes
	childrenContent := make(map[string]string)
	childrenBindings := make(map[string]*CollectedBindings)
	allCollectedStyles := make(map[string]string)
	var childStyles string
	for path, child := range children {
		name := strings.TrimPrefix(path, "/")
		// Use prefix to ensure unique IDs across children
		// Also pass the child's store map so nested components can resolve dynamic props
		result := RenderHTMLWithContextFull(child.Render(),
			WithPrefixCtx(name),
			WithParentStoreMapCtx(childStoreMaps[name]),
		)
		childrenContent[name] = result.HTML

		// Collect child styles (from component's Style() method)
		if hs, ok := child.(HasStyle); ok {
			childStyles += hs.Style() + "\n"
		}

		// Merge collected styles from nested components
		for compName, style := range result.CollectedStyles {
			if _, exists := allCollectedStyles[compName]; !exists {
				allCollectedStyles[compName] = style
			}
		}

		// Resolve child bindings and store them for if-block branch inclusion
		if result.Bindings != nil {
			resolveBindings(result.Bindings, childStoreMaps[name], name, child)
			childrenBindings[name] = result.Bindings
		}
	}

	// Get routes from app
	var routes []StaticRoute
	if hr, ok := app.(HasRoutes); ok {
		routes = hr.Routes()
	}

	// Generate HTML file for each route
	for _, route := range routes {
		// Get the child name for this route
		childName := strings.TrimPrefix(route.Path, "/")
		if childName == "" {
			childName = "basics" // Default for "/" route
		}

		// Render app with this specific child's content
		html, bindings := RenderHTMLWithChildContent(app.Render(), childName, childrenContent, childrenBindings)

		// Resolve app bindings
		resolveBindings(bindings, appStoreMap, "component", app)

		// Build full HTML document
		fullHTML := buildHTMLDocument(html, appStyle, childStyles, allCollectedStyles)

		// Write HTML file
		htmlPath := filepath.Join("dist", route.HTMLFile)
		os.WriteFile(htmlPath, []byte(fullHTML), 0644)
		fmt.Fprintf(os.Stderr, "Generated: %s\n", htmlPath)

		// Merge bindings for this route into allBindings
		mergeBindings(&allBindings, bindings)
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

	// Merge attr bindings
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

	// Merge each blocks
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

	// Merge attr cond bindings
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

}

// buildStoreMap builds a map from store pointer addresses to field paths.
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

// initStores initializes all nil store fields in a component with default values.
func initStores(comp Component) {
	rv := reflect.ValueOf(comp).Elem()
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		f := rv.Field(i)
		ft := rt.Field(i).Type

		// Check if field is a nil pointer to a store type
		if f.Kind() == reflect.Ptr && f.IsNil() && f.CanSet() {
			ftStr := ft.String()
			switch {
			case ftStr == "*preveltekit.LocalStore":
				f.Set(reflect.ValueOf(&LocalStore{Store: New("")}))
			case strings.HasPrefix(ftStr, "*preveltekit.Store[string]"):
				f.Set(reflect.ValueOf(New("")))
			case strings.HasPrefix(ftStr, "*preveltekit.Store[int]"):
				f.Set(reflect.ValueOf(New(0)))
			case strings.HasPrefix(ftStr, "*preveltekit.Store[bool]"):
				f.Set(reflect.ValueOf(New(false)))
			case strings.HasPrefix(ftStr, "*preveltekit.Store[float64]"):
				f.Set(reflect.ValueOf(New(0.0)))
			case strings.HasPrefix(ftStr, "*preveltekit.List[string]"):
				f.Set(reflect.ValueOf(NewList[string]()))
			case strings.HasPrefix(ftStr, "*preveltekit.List[int]"):
				f.Set(reflect.ValueOf(NewList[int]()))
			}
		}
	}
}

// resolveBindings resolves store references in bindings to field paths.
func resolveBindings(bindings *CollectedBindings, storeMap map[uintptr]string, prefix string, comp Component) {
	// Resolve text bindings
	for i := range bindings.TextBindings {
		// Skip if already resolved (e.g., from child component bindings)
		if bindings.TextBindings[i].StoreID != "" {
			continue
		}
		if bindings.TextBindings[i].StoreRef != nil {
			addr := reflect.ValueOf(bindings.TextBindings[i].StoreRef).Pointer()
			if name, ok := storeMap[addr]; ok {
				bindings.TextBindings[i].StoreID = name
			}
		}
	}

	// Resolve input bindings
	for i := range bindings.InputBindings {
		// Skip if already resolved
		if bindings.InputBindings[i].StoreID != "" {
			continue
		}
		if bindings.InputBindings[i].StoreRef != nil {
			addr := reflect.ValueOf(bindings.InputBindings[i].StoreRef).Pointer()
			if name, ok := storeMap[addr]; ok {
				bindings.InputBindings[i].StoreID = name
			}
		}
	}

	// Resolve if-block conditions and recursively resolve nested bindings
	for i := range bindings.IfBlocks {
		for j := range bindings.IfBlocks[i].Branches {
			// Skip if already resolved
			if bindings.IfBlocks[i].Branches[j].StoreID != "" {
				// Still need to recursively resolve nested bindings
				if bindings.IfBlocks[i].Branches[j].Bindings != nil {
					resolveBindings(bindings.IfBlocks[i].Branches[j].Bindings, storeMap, prefix, comp)
				}
				continue
			}
			cond := bindings.IfBlocks[i].Branches[j].CondRef
			if cond != nil {
				if sc, ok := cond.(*StoreCondition); ok && sc.Store != nil {
					addr := reflect.ValueOf(sc.Store).Pointer()
					if name, ok := storeMap[addr]; ok {
						// Build expression with proper operand quoting
						operand := fmt.Sprintf("%v", sc.Operand)
						if !isNumeric(operand) && operand != "true" && operand != "false" {
							operand = `"` + operand + `"`
						}
						bindings.IfBlocks[i].Branches[j].CondExpr = name + ".Get() " + sc.Op + " " + operand

						// Add structured condition data for WASM evaluation
						bindings.IfBlocks[i].Branches[j].StoreID = name
						bindings.IfBlocks[i].Branches[j].Op = sc.Op
						bindings.IfBlocks[i].Branches[j].Operand = fmt.Sprintf("%v", sc.Operand)

						// Add to deps
						parts := strings.Split(name, ".")
						fieldName := parts[len(parts)-1]
						if prefix == "component" {
							bindings.IfBlocks[i].Deps = append(bindings.IfBlocks[i].Deps, fieldName)
						} else {
							bindings.IfBlocks[i].Deps = append(bindings.IfBlocks[i].Deps, name)
						}
					}
				}

				if bc, ok := cond.(*BoolCondition); ok && bc.Store != nil {
					addr := reflect.ValueOf(bc.Store).Pointer()
					if name, ok := storeMap[addr]; ok {
						bindings.IfBlocks[i].Branches[j].CondExpr = name + ".Get()"

						// Add structured condition data for WASM evaluation
						bindings.IfBlocks[i].Branches[j].StoreID = name
						bindings.IfBlocks[i].Branches[j].IsBool = true

						parts := strings.Split(name, ".")
						fieldName := parts[len(parts)-1]
						if prefix == "component" {
							bindings.IfBlocks[i].Deps = append(bindings.IfBlocks[i].Deps, fieldName)
						} else {
							bindings.IfBlocks[i].Deps = append(bindings.IfBlocks[i].Deps, name)
						}
					}
				}
			}

			// Recursively resolve nested bindings in this branch
			if bindings.IfBlocks[i].Branches[j].Bindings != nil {
				resolveBindings(bindings.IfBlocks[i].Branches[j].Bindings, storeMap, prefix, comp)
			}
		}

		// Recursively resolve else branch bindings
		if bindings.IfBlocks[i].ElseBindings != nil {
			resolveBindings(bindings.IfBlocks[i].ElseBindings, storeMap, prefix, comp)
		}
	}

	// Event handlers are resolved directly in WASM via collectHandlers
	// No need to resolve handler names here anymore

	// Resolve attr bindings
	for i := range bindings.AttrBindings {
		// Skip if already resolved
		if len(bindings.AttrBindings[i].StoreIDs) > 0 {
			continue
		}
		if len(bindings.AttrBindings[i].StoreRefs) > 0 {
			var storeIDs []string
			for _, storeRef := range bindings.AttrBindings[i].StoreRefs {
				if storeRef != nil {
					addr := reflect.ValueOf(storeRef).Pointer()
					if name, ok := storeMap[addr]; ok {
						storeIDs = append(storeIDs, name)
					}
				}
			}
			bindings.AttrBindings[i].StoreIDs = storeIDs
		}
	}

	// Resolve each block list references
	for i := range bindings.EachBlocks {
		// Skip if already resolved
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

	// Resolve AttrCond bindings (from HtmlNode.AttrIf())
	for i := range bindings.AttrCondBindings {
		// Skip if already resolved
		if len(bindings.AttrCondBindings[i].Deps) > 0 {
			continue
		}

		var deps []string

		// Resolve condition store
		if bindings.AttrCondBindings[i].CondStoreRef != nil {
			addr := reflect.ValueOf(bindings.AttrCondBindings[i].CondStoreRef).Pointer()
			if name, ok := storeMap[addr]; ok {
				deps = append(deps, name)
			}
		}

		// Resolve true value store (if dynamic)
		if bindings.AttrCondBindings[i].TrueStoreRef != nil {
			addr := reflect.ValueOf(bindings.AttrCondBindings[i].TrueStoreRef).Pointer()
			if name, ok := storeMap[addr]; ok {
				bindings.AttrCondBindings[i].TrueStoreID = name
				deps = append(deps, name)
			}
		}

		// Resolve false value store (if dynamic)
		if bindings.AttrCondBindings[i].FalseStoreRef != nil {
			addr := reflect.ValueOf(bindings.AttrCondBindings[i].FalseStoreRef).Pointer()
			if name, ok := storeMap[addr]; ok {
				bindings.AttrCondBindings[i].FalseStoreID = name
				deps = append(deps, name)
			}
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
// Load bindings first, then WASM
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
