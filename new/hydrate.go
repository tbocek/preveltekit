//go:build !js || !wasm

package preveltekit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
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

// HydrateConfig configures the hydration process.
type HydrateConfig struct {
	// OutputDir is the directory where HTML files are written (default: "dist")
	OutputDir string
	// Children maps route paths to child components
	Children map[string]Component
	// NestedComponents maps component names to factory functions
	NestedComponents map[string]func() Component
}

// Hydrate is the main entry point for declarative components.
// In SSR mode (native build), it generates HTML files and outputs bindings.
// In WASM mode, it sets up DOM bindings for reactivity.
func Hydrate(app Component, opts ...func(*HydrateConfig)) {
	cfg := &HydrateConfig{
		OutputDir:        "dist",
		Children:         make(map[string]Component),
		NestedComponents: make(map[string]func() Component),
	}
	for _, opt := range opts {
		opt(cfg)
	}

	// SSR mode - generate HTML and bindings
	hydrateSSR(app, cfg)
}

// WithOutputDir sets the output directory for HTML files.
func WithOutputDir(dir string) func(*HydrateConfig) {
	return func(cfg *HydrateConfig) {
		cfg.OutputDir = dir
	}
}

// WithChild registers a child component for a route path.
func WithChild(path string, comp Component) func(*HydrateConfig) {
	return func(cfg *HydrateConfig) {
		cfg.Children[path] = comp
	}
}

// WithNestedComponent registers a nested component type by name.
// The factory function should return a new instance with initialized stores.
func WithNestedComponent(name string, factory func() Component) func(*HydrateConfig) {
	return func(cfg *HydrateConfig) {
		cfg.NestedComponents[name] = factory
	}
}

// hydrateSSR handles the SSR phase - generating HTML and collecting bindings.
func hydrateSSR(app Component, cfg *HydrateConfig) {
	mode := os.Getenv("HYDRATE_MODE")

	switch mode {
	case "generate-all":
		// Generate all HTML files and output merged bindings
		hydrateGenerateAll(app, cfg)
	default:
		// Single route mode (for backwards compatibility or specific route rendering)
		hydrateSingleRoute(app, cfg)
	}
}

// hydrateGenerateAll generates a single SPA HTML file with all children embedded.
func hydrateGenerateAll(app Component, cfg *HydrateConfig) {
	// Initialize stores for all components
	initStores(app)
	for _, child := range cfg.Children {
		initStores(child)
	}

	// Call OnCreate for all components
	if oc, ok := app.(HasOnCreate); ok {
		oc.OnCreate()
	}
	for _, child := range cfg.Children {
		if oc, ok := child.(HasOnCreate); ok {
			oc.OnCreate()
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
	for path, child := range cfg.Children {
		name := strings.TrimPrefix(path, "/")
		childStoreMaps[name] = buildStoreMap(child, name)
	}

	// Collect all bindings
	var allBindings CollectedBindings

	// Create output directory
	os.MkdirAll(cfg.OutputDir, 0755)

	// Pre-render all children with unique prefixes
	childrenContent := make(map[string]string)
	childrenBindings := make(map[string]*CollectedBindings)
	var childStyles string
	for path, child := range cfg.Children {
		name := strings.TrimPrefix(path, "/")
		// Use prefix to ensure unique IDs across children, and pass nested components
		// Also pass the child's store map so nested components can resolve dynamic props
		childHTML, childBindings := RenderHTMLWithContext(child.Render(),
			WithPrefixCtx(name),
			WithNestedComponentsCtx(cfg.NestedComponents),
			WithParentStoreMapCtx(childStoreMaps[name]),
		)
		childrenContent[name] = childHTML

		// Collect child styles
		if hs, ok := child.(HasStyle); ok {
			childStyles += hs.Style() + "\n"
		}

		// Resolve child bindings and store them for if-block branch inclusion
		if childBindings != nil {
			resolveBindings(childBindings, childStoreMaps[name], name, child)
			childrenBindings[name] = childBindings
		}
	}

	// Collect nested component styles
	for _, factory := range cfg.NestedComponents {
		comp := factory()
		if hs, ok := comp.(HasStyle); ok {
			childStyles += hs.Style() + "\n"
		}
	}

	// Render app with all children content and bindings available
	html, bindings := RenderHTMLWithChildren(app.Render(), childrenContent, childrenBindings)

	// Resolve app bindings
	resolveBindings(bindings, appStoreMap, "component", app)
	mergeBindings(&allBindings, bindings)

	// Build full HTML document with all styles
	fullHTML := buildHTMLDocument(html, appStyle, childStyles)

	// Write single index.html
	htmlPath := filepath.Join(cfg.OutputDir, "index.html")
	os.WriteFile(htmlPath, []byte(fullHTML), 0644)
	fmt.Fprintf(os.Stderr, "Generated: %s\n", htmlPath)

	// Output merged bindings as JSON
	bindingsJSON, _ := json.Marshal(allBindings)
	fmt.Fprintf(os.Stderr, "DEBUG: allBindings has %d text, %d events, %d if-blocks\n",
		len(allBindings.TextBindings), len(allBindings.Events), len(allBindings.IfBlocks))
	if len(allBindings.IfBlocks) > 0 {
		for i, ifb := range allBindings.IfBlocks {
			fmt.Fprintf(os.Stderr, "DEBUG: IfBlock[%d] marker=%s branches=%d\n", i, ifb.MarkerID, len(ifb.Branches))
			for j, br := range ifb.Branches {
				textCount := 0
				eventCount := 0
				if br.Bindings != nil {
					textCount = len(br.Bindings.TextBindings)
					eventCount = len(br.Bindings.Events)
				}
				fmt.Fprintf(os.Stderr, "DEBUG:   Branch[%d] text=%d events=%d\n", j, textCount, eventCount)
			}
		}
	}
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

	seenClass := make(map[string]bool)
	for _, b := range dst.ClassBindings {
		seenClass[b.ElementID+":"+b.ClassName] = true
	}
	for _, b := range src.ClassBindings {
		key := b.ElementID + ":" + b.ClassName
		if !seenClass[key] {
			dst.ClassBindings = append(dst.ClassBindings, b)
			seenClass[key] = true
		}
	}

	seenShowIf := make(map[string]bool)
	for _, b := range dst.ShowIfBindings {
		seenShowIf[b.ElementID] = true
	}
	for _, b := range src.ShowIfBindings {
		if !seenShowIf[b.ElementID] {
			dst.ShowIfBindings = append(dst.ShowIfBindings, b)
			seenShowIf[b.ElementID] = true
		}
	}
}

// hydrateSingleRoute handles rendering a single route (original behavior).
func hydrateSingleRoute(app Component, cfg *HydrateConfig) {
	prerenderPath := os.Getenv("PRERENDER_PATH")
	outputBindings := os.Getenv("OUTPUT_BINDINGS") == "1"

	// Initialize stores for all components
	initStores(app)
	for _, child := range cfg.Children {
		initStores(child)
	}

	// Call OnCreate if implemented
	if oc, ok := app.(HasOnCreate); ok {
		oc.OnCreate()
	}

	// Initialize child components
	for _, child := range cfg.Children {
		if oc, ok := child.(HasOnCreate); ok {
			oc.OnCreate()
		}
	}

	// Build store-to-field map for app component
	appStoreMap := buildStoreMap(app, "component")

	// Determine which child component to render based on prerender path
	var slotContent string
	var childBindings *CollectedBindings
	var activeChildName string
	var activeChild Component

	for path, child := range cfg.Children {
		if prerenderPath == path {
			activeChild = child
			activeChildName = componentVarName(child)
			slotContent, childBindings = RenderHTML(child.Render())
			break
		}
	}

	// Default to first child for "/" or empty path
	if activeChild == nil && (prerenderPath == "/" || prerenderPath == "") {
		// Try to find a "basics" or first child
		for path, child := range cfg.Children {
			if strings.Contains(path, "basics") {
				activeChild = child
				activeChildName = componentVarName(child)
				slotContent, childBindings = RenderHTML(child.Render())
				break
			}
		}
		// If no "basics", just use first child
		if activeChild == nil {
			for _, child := range cfg.Children {
				activeChild = child
				activeChildName = componentVarName(child)
				slotContent, childBindings = RenderHTML(child.Render())
				break
			}
		}
	}

	// Render app with slot content
	var html string
	var bindings *CollectedBindings
	if slotContent != "" {
		html, bindings = RenderHTMLWithSlot(app.Render(), slotContent)
	} else {
		html, bindings = RenderHTML(app.Render())
	}

	// Resolve app bindings
	resolveBindings(bindings, appStoreMap, "component", app)

	// Resolve child bindings if present
	if childBindings != nil && activeChild != nil {
		childStoreMap := buildStoreMap(activeChild, activeChildName)
		resolveBindings(childBindings, childStoreMap, activeChildName, activeChild)

		// Merge child bindings into main bindings
		bindings.TextBindings = append(bindings.TextBindings, childBindings.TextBindings...)
		bindings.InputBindings = append(bindings.InputBindings, childBindings.InputBindings...)
		bindings.Events = append(bindings.Events, childBindings.Events...)
		bindings.IfBlocks = append(bindings.IfBlocks, childBindings.IfBlocks...)
		bindings.ClassBindings = append(bindings.ClassBindings, childBindings.ClassBindings...)
	}

	// Output HTML
	fmt.Print(html)

	// Output bindings if requested
	if outputBindings {
		bindingsJSON, _ := json.Marshal(bindings)
		fmt.Fprintf(os.Stderr, "BINDINGS:%s\n", bindingsJSON)
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

	// Resolve class bindings
	for i := range bindings.ClassBindings {
		// Skip if already resolved
		if bindings.ClassBindings[i].CondExpr != "" {
			continue
		}
		if bindings.ClassBindings[i].StoreRef != nil {
			addr := reflect.ValueOf(bindings.ClassBindings[i].StoreRef).Pointer()
			if name, ok := storeMap[addr]; ok {
				// Build expression based on whether it's a comparison or simple bool
				if bindings.ClassBindings[i].Op != "" {
					// StoreCondition with comparison
					operand := bindings.ClassBindings[i].Operand
					// Quote string operands
					if !isNumeric(operand) && operand != "true" && operand != "false" {
						operand = `"` + operand + `"`
					}
					bindings.ClassBindings[i].CondExpr = name + ".Get() " + bindings.ClassBindings[i].Op + " " + operand
				} else {
					// BoolCondition
					bindings.ClassBindings[i].CondExpr = name + ".Get()"
				}
				// Extract just the field name for deps
				parts := strings.Split(name, ".")
				fieldName := parts[len(parts)-1]
				if prefix == "component" {
					bindings.ClassBindings[i].Deps = []string{fieldName}
				} else {
					bindings.ClassBindings[i].Deps = []string{name}
				}
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

	// Resolve ShowIf bindings
	for i := range bindings.ShowIfBindings {
		// Skip if already resolved
		if bindings.ShowIfBindings[i].StoreID != "" {
			continue
		}
		if bindings.ShowIfBindings[i].StoreRef != nil {
			addr := reflect.ValueOf(bindings.ShowIfBindings[i].StoreRef).Pointer()
			if name, ok := storeMap[addr]; ok {
				bindings.ShowIfBindings[i].StoreID = name
				// Extract just the field name for deps
				parts := strings.Split(name, ".")
				fieldName := parts[len(parts)-1]
				if prefix == "component" {
					bindings.ShowIfBindings[i].Deps = []string{fieldName}
				} else {
					bindings.ShowIfBindings[i].Deps = []string{name}
				}
			}
		}
	}

	// Resolve event handlers
	for i := range bindings.Events {
		// Skip if already resolved (e.g., from child component bindings)
		if bindings.Events[i].HandlerID != "" {
			continue
		}
		if bindings.Events[i].HandlerRef != nil {
			fn := reflect.ValueOf(bindings.Events[i].HandlerRef)
			if fn.Kind() == reflect.Func {
				name := runtime.FuncForPC(fn.Pointer()).Name()
				// Extract method name from "main.(*Counter).Increment-fm"
				if idx := strings.LastIndex(name, "."); idx >= 0 {
					name = name[idx+1:]
				}
				name = strings.TrimSuffix(name, "-fm")
				bindings.Events[i].HandlerID = prefix + "." + name
			}
		}
		// Serialize args
		if len(bindings.Events[i].Args) > 0 {
			var argStrs []string
			for _, arg := range bindings.Events[i].Args {
				argStrs = append(argStrs, fmt.Sprintf("%v", arg))
			}
			bindings.Events[i].ArgsStr = strings.Join(argStrs, ", ")
		}
	}

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

// GenerateHTMLFiles generates HTML files for all routes.
// This is called by the CLI after collecting bindings.
func GenerateHTMLFiles(app Component, cfg *HydrateConfig, outputDir string) error {
	// Get routes from app
	routes, ok := app.(HasRoutes)
	if !ok {
		return fmt.Errorf("app component must implement HasRoutes")
	}

	// Get app style
	var appStyle string
	if hs, ok := app.(HasStyle); ok {
		appStyle = hs.Style()
	}

	// Generate HTML for each route
	for _, route := range routes.Routes() {
		prerenderPath := route.Path

		// Find child component for this route
		var childComp Component
		var childName string
		for path, child := range cfg.Children {
			if path == prerenderPath || (prerenderPath == "/" && strings.Contains(path, "basics")) {
				childComp = child
				childName = strings.TrimPrefix(path, "/")
				break
			}
		}

		// Render with slot content - use prefix to match bindings markers
		var html string
		if childComp != nil {
			slotContent, _ := RenderHTMLWithPrefix(childComp.Render(), childName)
			html, _ = RenderHTMLWithSlot(app.Render(), slotContent)
		} else {
			html, _ = RenderHTML(app.Render())
		}

		// Build full HTML document
		var childStyle string
		if childComp != nil {
			if hs, ok := childComp.(HasStyle); ok {
				childStyle = hs.Style()
			}
		}

		fullHTML := buildHTMLDocument(html, appStyle, childStyle)

		// Write to file
		htmlFile := filepath.Join(outputDir, route.HTMLFile)
		if err := os.MkdirAll(filepath.Dir(htmlFile), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(htmlFile, []byte(fullHTML), 0644); err != nil {
			return err
		}
	}

	return nil
}

func buildHTMLDocument(body, appStyle, childStyle string) string {
	var styles string
	if appStyle != "" {
		styles += "<style>" + appStyle + "</style>\n"
	}
	if childStyle != "" {
		styles += "<style>" + childStyle + "</style>\n"
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
