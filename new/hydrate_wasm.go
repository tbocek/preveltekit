//go:build js && wasm

package preveltekit

import (
	"reflect"
	"syscall/js"
)

// Track which if-blocks have been set up to avoid duplicates
var setupIfBlocks = make(map[string]bool)

// Track which each-blocks have been set up to avoid duplicates
var setupEachBlocks = make(map[string]bool)

// trimPrefix removes prefix from s if present
func trimPrefix(s, prefix string) string {
	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}

// splitFirst splits s on first occurrence of sep, returns (before, after, found)
func splitFirst(s, sep string) (string, string, bool) {
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			return s[:i], s[i+len(sep):], true
		}
	}
	return s, "", false
}

// containsChar checks if s contains the character c
func containsChar(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}

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

// HydrateConfig configures the hydration process.
type HydrateConfig struct {
	OutputDir string
	Children  map[string]Component
}

// Hydrate sets up DOM bindings for reactivity.
// The bindings JSON is passed via a global variable set by the CLI-generated code.
func Hydrate(app Component, opts ...func(*HydrateConfig)) {
	cfg := &HydrateConfig{
		OutputDir: "dist",
		Children:  make(map[string]Component),
	}
	for _, opt := range opts {
		opt(cfg)
	}

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
			if currentStore != nil {
				for _, route := range hr.Routes() {
					if route.Handler != nil {
						// Call handler to set CurrentComponent
						route.Handler(nil)
						// Read the component that was set
						if comp := currentStore.Get(); comp != nil {
							cfg.Children[route.Path] = comp
						}
					}
				}
			}
		}
	}

	hydrateWASM(app, cfg)
}

// WithOutputDir sets the output directory (no-op in WASM, but needed for API compatibility).
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

// hydrateWASM sets up DOM bindings from the embedded bindings JSON.
func hydrateWASM(app Component, cfg *HydrateConfig) {
	// Get bindings from global variable (set by CLI-generated code)
	bindingsJS := js.Global().Get("_preveltekit_bindings")
	if bindingsJS.IsUndefined() || bindingsJS.IsNull() {
		runLifecycle(app, cfg)
		select {}
		return
	}

	bindingsJSON := bindingsJS.String()
	bindings := parseBindings(bindingsJSON)
	if bindings == nil {
		runLifecycle(app, cfg)
		select {}
		return
	}

	// Build component map: "component" -> app, "basics" -> child, etc.
	components := map[string]Component{
		"component": app,
	}
	for path, child := range cfg.Children {
		name := trimPrefix(path, "/")
		components[name] = child
	}

	// Run OnCreate phase
	runOnCreate(app, cfg)

	// Collect event handlers from component Render() trees
	collectHandlers(app, "")
	for path, child := range cfg.Children {
		name := trimPrefix(path, "/")
		collectHandlers(child, name)
	}

	// Apply all bindings
	cleanup := &Cleanup{}
	applyBindings(bindings, components, cleanup)

	// Run OnMount phase
	runOnMount(app, cfg)

	// Keep WASM running
	select {}
}

// runOnCreate calls OnCreate and injects styles for app and all children.
func runOnCreate(app Component, cfg *HydrateConfig) {
	// App's OnCreate was already called in Hydrate before Routes() was called
	if hs, ok := app.(HasStyle); ok {
		InjectStyle("app", hs.Style())
	}
	// Initialize stores for all children first
	for _, child := range cfg.Children {
		initStores(child)
	}
	// Then call OnCreate
	for path, child := range cfg.Children {
		if oc, ok := child.(HasOnCreate); ok {
			oc.OnCreate()
		}
		if hs, ok := child.(HasStyle); ok {
			InjectStyle(trimPrefix(path, "/"), hs.Style())
		}
	}
}

// runOnMount calls OnMount for app and all children.
func runOnMount(app Component, cfg *HydrateConfig) {
	if om, ok := app.(HasOnMount); ok {
		om.OnMount()
	}
	for _, child := range cfg.Children {
		if om, ok := child.(HasOnMount); ok {
			om.OnMount()
		}
	}
}

// runLifecycle runs full lifecycle (OnCreate + styles + OnMount) for all components.
func runLifecycle(app Component, cfg *HydrateConfig) {
	runOnCreate(app, cfg)
	runOnMount(app, cfg)
}

// resolveStore resolves a store path like "component.Count" or "basics.Score" to a store pointer.
func resolveStore(storeID string, components map[string]Component) any {
	compName, fieldName, ok := splitFirst(storeID, ".")
	if !ok {
		return nil
	}

	comp, ok := components[compName]
	if !ok {
		return nil
	}

	rv := reflect.ValueOf(comp).Elem()
	field := rv.FieldByName(fieldName)
	if !field.IsValid() || field.IsNil() {
		return nil
	}

	return field.Interface()
}

// handlerMap stores element ID -> handler function mappings collected from Render()
var handlerMap = make(map[string]func())

// collectHandlers walks a component's Render() tree and extracts event handlers.
// It uses the same ID generation logic as SSR to ensure IDs match.
func collectHandlers(comp Component, prefix string) {
	node := comp.Render()
	ctx := &handlerCollectCtx{prefix: prefix}
	collectHandlersFromNode(node, ctx)
}

// handlerCollectCtx tracks state while collecting handlers (mirrors BuildContext counters)
type handlerCollectCtx struct {
	prefix       string
	eventCounter int
	bindCounter  int
	classCounter int
	attrCounter  int
	compCounter  int
}

func (ctx *handlerCollectCtx) nextEventID() string {
	id := "ev" + intToStr(ctx.eventCounter)
	ctx.eventCounter++
	return id
}

func (ctx *handlerCollectCtx) nextClassID() string {
	id := "cl" + intToStr(ctx.classCounter)
	ctx.classCounter++
	return id
}

func (ctx *handlerCollectCtx) fullElementID(localID string) string {
	if ctx.prefix == "" {
		return localID
	}
	return ctx.prefix + "_" + localID
}

// collectHandlersFromNode recursively walks a node tree collecting event handlers.
func collectHandlersFromNode(n Node, ctx *handlerCollectCtx) {
	if n == nil {
		return
	}

	switch node := n.(type) {
	case *Element:
		// Check for event attributes
		for _, attr := range node.Attrs {
			if ev, ok := attr.(*EventAttr); ok {
				localID := ctx.nextEventID()
				fullID := ctx.fullElementID(localID)
				// Store the handler function directly
				handlerMap[fullID] = ev.Handler
			}
		}
		// Recurse into children
		for _, child := range node.Children {
			collectHandlersFromNode(child, ctx)
		}

	case *HtmlNode:
		// Html nodes can contain embedded EventAttr
		for _, part := range node.Parts {
			if ev, ok := part.(*EventAttr); ok {
				localID := ctx.nextEventID()
				fullID := ctx.fullElementID(localID)
				handlerMap[fullID] = ev.Handler
			} else if childNode, ok := part.(Node); ok {
				collectHandlersFromNode(childNode, ctx)
			}
		}

	case *Fragment:
		for _, child := range node.Children {
			collectHandlersFromNode(child, ctx)
		}

	case *IfNode:
		// Collect from all branches (handlers exist in all possible paths)
		for _, branch := range node.Branches {
			for _, child := range branch.Children {
				collectHandlersFromNode(child, ctx)
			}
		}
		for _, child := range node.ElseNode {
			collectHandlersFromNode(child, ctx)
		}

	case *EachNode:
		// Each nodes have a body function - we can't easily extract handlers from it
		// since it requires an item. For now, skip (each items are re-rendered anyway)

	case *ClassIfNode:
		// ClassIfNode can have an OnClick handler attached via WithOnClick()
		if node.OnClick != nil {
			localID := ctx.nextEventID()
			fullID := ctx.fullElementID(localID)
			handlerMap[fullID] = node.OnClick
		} else {
			// No OnClick - still need to increment class counter to stay in sync with SSR
			ctx.nextClassID()
		}

	case *ComponentNode:
		// Nested component - recurse with new prefix
		compMarker := "comp" + intToStr(ctx.compCounter)
		ctx.compCounter++
		fullPrefix := ctx.fullElementID(compMarker)
		if comp, ok := node.Instance.(Component); ok {
			collectHandlers(comp, fullPrefix)
		}
	}
}

// getHandler looks up a handler by element ID from the collected handlers.
func getHandler(elementID string) func() {
	return handlerMap[elementID]
}

// bindTextDynamic sets up a text binding for any store type.
func bindTextDynamic(markerID string, store any, isHTML bool) {
	switch s := store.(type) {
	case *Store[string]:
		bindText(markerID, s, isHTML)
	case *Store[int]:
		bindText(markerID, s, isHTML)
	case *Store[bool]:
		bindText(markerID, s, isHTML)
	case *Store[float64]:
		bindText(markerID, s, isHTML)
	}
}

// bindText is a helper that calls BindText or BindHTML based on isHTML flag.
func bindText[T any](markerID string, store Bindable[T], isHTML bool) {
	if isHTML {
		BindHTML(markerID, store)
	} else {
		BindText(markerID, store)
	}
}

// bindInputDynamic sets up an input binding using reflection.
func bindInputDynamic(cleanup *Cleanup, elementID string, store any, bindType string) {
	switch s := store.(type) {
	case *Store[string]:
		if bindType == "checked" {
			// String store with checkbox doesn't make sense, skip
		} else {
			BindInputs(cleanup, []Inp{{elementID, s}})
		}
	case *Store[bool]:
		if bindType == "checked" {
			BindCheckboxes(cleanup, []Chk{{elementID, s}})
		}
	}
}

// bindEventDynamic sets up an event binding.
func bindEventDynamic(cleanup *Cleanup, elementID, event string, handler func()) {
	BindEvents(cleanup, []Evt{{elementID, event, handler}})
}

// bindShowIf sets up a show/hide binding based on condition.
func bindShowIf(sib HydrateShowIfBinding, components map[string]Component) {
	el := GetEl(sib.ElementID)
	if !ok(el) {
		return
	}

	store := resolveStore(sib.StoreID, components)
	if store == nil {
		return
	}

	// Function to evaluate the condition and update visibility
	updateVisibility := func() {
		var visible bool
		if sib.IsBool {
			if s, ok := store.(*Store[bool]); ok {
				visible = s.Get()
			}
		} else {
			// Comparison condition
			switch s := store.(type) {
			case *Store[int]:
				visible = compare(s.Get(), sib.Op, atoiSafe(sib.Operand))
			case *Store[string]:
				visible = compare(s.Get(), sib.Op, sib.Operand)
			case *Store[float64]:
				visible = compare(s.Get(), sib.Op, atofSafe(sib.Operand))
			}
		}

		if visible {
			el.Get("style").Call("removeProperty", "display")
		} else {
			el.Get("style").Set("display", "none")
		}
	}

	// Initial update
	updateVisibility()

	// Subscribe to changes
	subscribeToStore(store, updateVisibility)
}

// bindIfBlock sets up an if-block with reactive condition evaluation.
func bindIfBlock(ifb HydrateIfBlock, components map[string]Component) {
	// Skip if already set up (prevents duplicate setup from nested if-blocks)
	if setupIfBlocks[ifb.MarkerID] {
		return
	}
	setupIfBlocks[ifb.MarkerID] = true

	// Find the existing SSR content
	currentEl := FindExistingIfContent(ifb.MarkerID)

	// Track current cleanup for bindings
	currentCleanup := &Cleanup{}

	// Evaluate which branch is active and update content
	updateIfBlock := func() {
		var activeHTML string
		var activeBindings *HydrateBindings
		found := false

		for i := 0; i < len(ifb.Branches); i++ {
			branch := ifb.Branches[i]
			if evalCondition(branch, components) {
				activeHTML = branch.HTML
				activeBindings = branch.Bindings
				found = true
				break
			}
		}

		if !found {
			activeHTML = ifb.ElseHTML
			activeBindings = ifb.ElseBindings
		}

		currentEl = FindExistingIfContent(ifb.MarkerID)
		currentEl = ReplaceContent(ifb.MarkerID, currentEl, activeHTML)

		currentCleanup.Release()
		currentCleanup = &Cleanup{}

		if activeBindings != nil {
			clearBoundMarkers(activeBindings)
			applyBindings(activeBindings, components, currentCleanup)
		}
	}

	// Subscribe to store changes for all dependencies (deduplicated)
	seenDeps := make(map[string]bool)
	for _, dep := range ifb.Deps {
		if seenDeps[dep] {
			continue
		}
		seenDeps[dep] = true

		store := resolveStore(dep, components)
		if store == nil {
			// Try with component prefix if dep doesn't contain a dot
			if !containsChar(dep, '.') {
				store = resolveStore("component."+dep, components)
			}
		}
		if store != nil {
			subscribeToStore(store, updateIfBlock)
		}
	}

	// Call updateIfBlock to sync DOM with current state
	// This handles nested if-blocks where state may have changed after SSR
	updateIfBlock()
}

// clearBoundMarkers clears marker tracking for bindings that will be re-applied.
// This is needed when if-block content is replaced via ReplaceContent.
func clearBoundMarkers(bindings *HydrateBindings) {
	if bindings == nil {
		return
	}
	for _, tb := range bindings.TextBindings {
		ClearBoundMarker(tb.MarkerID)
	}
	// Recursively clear nested if-block markers
	for _, ifb := range bindings.IfBlocks {
		for _, branch := range ifb.Branches {
			if branch.Bindings != nil {
				clearBoundMarkers(branch.Bindings)
			}
		}
		if ifb.ElseBindings != nil {
			clearBoundMarkers(ifb.ElseBindings)
		}
		// Also clear the if-block's own setup status so it can be re-setup
		delete(setupIfBlocks, ifb.MarkerID)
	}
	// Clear each-block setup status
	for _, eb := range bindings.EachBlocks {
		delete(setupEachBlocks, eb.MarkerID)
	}
}

// applyBindings applies all bindings from a HydrateBindings struct to the DOM.
func applyBindings(bindings *HydrateBindings, components map[string]Component, cleanup *Cleanup) {
	// Text bindings
	for _, tb := range bindings.TextBindings {
		store := resolveStore(tb.StoreID, components)
		if store != nil {
			bindTextDynamic(tb.MarkerID, store, tb.IsHTML)
		}
	}

	// Input bindings
	for _, ib := range bindings.InputBindings {
		store := resolveStore(ib.StoreID, components)
		if store != nil {
			bindInputDynamic(cleanup, ib.ElementID, store, ib.BindType)
		}
	}

	// Event bindings
	for _, ev := range bindings.Events {
		handler := getHandler(ev.ElementID)
		if handler != nil {
			bindEventDynamic(cleanup, ev.ElementID, ev.Event, handler)
		}
	}

	// Class bindings
	for _, cb := range bindings.ClassBindings {
		bindClassBinding(cb, components)
	}

	// Nested if-blocks
	for _, ifb := range bindings.IfBlocks {
		bindIfBlock(ifb, components)
	}

	// ShowIf bindings
	for _, sib := range bindings.ShowIfBindings {
		bindShowIf(sib, components)
	}

	// Attr bindings (dynamic attributes like data-type)
	for _, ab := range bindings.AttrBindings {
		bindAttr(ab, components)
	}

	// Each block bindings (list iteration)
	for _, eb := range bindings.EachBlocks {
		bindEachBlock(eb, components)
	}

	// Router bindings
	for _, rb := range bindings.RouterBindings {
		bindRouter(rb, components)
	}
}

// bindClassBinding sets up a class toggle binding based on condition.
func bindClassBinding(cb HydrateClassBinding, components map[string]Component) {
	el := GetEl(cb.ElementID)
	if !ok(el) {
		return
	}

	// Resolve stores from dependencies
	var stores []any
	for _, dep := range cb.Deps {
		store := resolveStore(dep, components)
		if store != nil {
			stores = append(stores, store)
		}
	}

	if len(stores) == 0 {
		return
	}

	// Use first store for condition evaluation
	store := stores[0]

	// Function to evaluate the condition and toggle class
	updateClass := func() {
		var active bool

		// Check if it's a boolean store with no operator
		if cb.Op == "" {
			if s, ok := store.(*Store[bool]); ok {
				active = s.Get()
			}
		} else {
			// Comparison condition
			switch s := store.(type) {
			case *Store[int]:
				active = compare(s.Get(), cb.Op, atoiSafe(cb.Operand))
			case *Store[string]:
				active = compare(s.Get(), cb.Op, cb.Operand)
			case *Store[float64]:
				active = compare(s.Get(), cb.Op, atofSafe(cb.Operand))
			case *Store[bool]:
				// Boolean with operator (e.g., "== true")
				operandBool := cb.Operand == "true"
				active = compareBool(s.Get(), cb.Op, operandBool)
			}
		}

		classList := el.Get("classList")
		if active {
			classList.Call("add", cb.ClassName)
		} else {
			classList.Call("remove", cb.ClassName)
		}
	}

	// Initial update
	updateClass()

	// Subscribe to changes for all dependencies
	for _, s := range stores {
		subscribeToStore(s, updateClass)
	}
}

// bindAttr sets up a dynamic attribute binding.
func bindAttr(ab HydrateAttrBinding, components map[string]Component) {
	el := GetEl(ab.ElementID)
	if !ok(el) {
		// Try finding by data-attrbind attribute
		el = Document.Call("querySelector", "[data-attrbind=\""+ab.ElementID+"\"]")
		if !ok(el) {
			return
		}
	}

	// Collect stores for this binding
	var stores []any
	for _, storeID := range ab.StoreIDs {
		store := resolveStore(storeID, components)
		if store != nil {
			stores = append(stores, store)
		}
	}

	if len(stores) == 0 {
		return
	}

	// Function to update the attribute value
	updateAttr := func() {
		value := ab.Template
		for i, store := range stores {
			placeholder := "{" + intToStr(i) + "}"
			value = replaceAll(value, placeholder, storeToString(store))
		}
		el.Call("setAttribute", ab.AttrName, value)
	}

	// Initial update
	updateAttr()

	// Subscribe to changes
	for _, store := range stores {
		subscribeToStore(store, updateAttr)
	}
}

// bindEachBlock sets up a list iteration binding.
func bindEachBlock(eb HydrateEachBlock, components map[string]Component) {
	if eb.ListID == "" {
		return
	}

	// Check if already setup
	if setupEachBlocks[eb.MarkerID] {
		return
	}
	setupEachBlocks[eb.MarkerID] = true

	// Find the marker comment
	marker := FindComment(eb.MarkerID)
	if marker.IsNull() {
		return
	}

	// Resolve the list
	listAny := resolveStore(eb.ListID, components)
	if listAny == nil {
		return
	}

	// Get the component that owns this list for rendering
	compName, _, _ := splitFirst(eb.ListID, ".")
	comp, compOk := components[compName]
	if !compOk {
		return
	}

	// Get the parent element (the <ul> or container that holds the list items)
	// The marker comment is placed after all the <span> wrappers, inside the container
	parent := marker.Get("parentNode")

	// Extract item ID prefix from marker (e.g., "lists_e0" -> "lists_")
	itemIDPrefix := eb.MarkerID[:len(eb.MarkerID)-1]
	if len(itemIDPrefix) > 0 && itemIDPrefix[len(itemIDPrefix)-1] == 'e' {
		itemIDPrefix = itemIDPrefix[:len(itemIDPrefix)-1]
	}

	// Subscribe to list changes and re-render
	switch list := listAny.(type) {
	case *List[string]:
		bindListItems(list, parent, itemIDPrefix, escapeHTMLWasm)
	case *List[int]:
		bindListItems(list, parent, itemIDPrefix, intToStr)
	}

	_ = comp
}

// bindListItems sets up list rendering and subscribes to changes.
func bindListItems[T comparable](list *List[T], parent js.Value, itemIDPrefix string, format func(T) string) {
	renderItems := func(items []T) {
		var html string
		for i, item := range items {
			html += `<span id="` + itemIDPrefix + `_` + intToStr(i) + `"><li><span class="index">` + intToStr(i) + `</span> ` + format(item) + `</li></span>`
		}
		if !parent.IsNull() && parent.Get("nodeType").Int() == 1 {
			parent.Set("innerHTML", html)
		}
	}

	// Don't render initially - SSR already rendered the items
	// Just subscribe to changes for future updates
	list.OnChange(renderItems)
}

// bindRouter sets up page switching based on a component store.
func bindRouter(rb HydrateRouterBinding, components map[string]Component) {
	store := resolveStore(rb.StoreID, components)
	if store == nil {
		return
	}

	compStore, isCompStore := store.(*Store[Component])
	if !isCompStore {
		return
	}

	// Track currently visible page
	var currentPage string

	updateVisibility := func(comp Component) {
		newPage := componentNameWasm(comp)

		// Hide old page
		if currentPage != "" && currentPage != newPage {
			if oldEl := GetEl("page-" + currentPage); ok(oldEl) {
				oldEl.Get("style").Set("display", "none")
			}
		}
		// Hide notfound
		if notFoundEl := GetEl("page-notfound"); ok(notFoundEl) {
			notFoundEl.Get("style").Set("display", "none")
		}

		// Show new page
		newEl := GetEl("page-" + newPage)
		if ok(newEl) {
			newEl.Get("style").Call("removeProperty", "display")
			currentPage = newPage
		} else {
			// Show notfound if page doesn't exist
			if notFoundEl := GetEl("page-notfound"); ok(notFoundEl) {
				notFoundEl.Get("style").Call("removeProperty", "display")
			}
			currentPage = ""
		}
	}

	// Initial state - don't call updateVisibility since SSR already rendered correctly
	currentPage = componentNameWasm(compStore.Get())

	// Subscribe to changes
	compStore.OnChange(updateVisibility)
}

// componentNameWasm returns the lowercase type name of a component.
func componentNameWasm(c Component) string {
	if c == nil {
		return ""
	}
	t := reflect.TypeOf(c)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	name := t.Name()
	if len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z' {
		b := []byte(name)
		b[0] = b[0] + 32
		return string(b)
	}
	return name
}

// escapeHTMLWasm escapes HTML special characters
func escapeHTMLWasm(s string) string {
	var result []byte
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '&':
			result = append(result, []byte("&amp;")...)
		case '<':
			result = append(result, []byte("&lt;")...)
		case '>':
			result = append(result, []byte("&gt;")...)
		case '"':
			result = append(result, []byte("&quot;")...)
		default:
			result = append(result, s[i])
		}
	}
	return string(result)
}

// intToStr converts int to string without fmt
func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
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
			// Get the element type to check the actual store type
			elemType := ft.Elem()
			elemName := elemType.Name()
			pkgPath := elemType.PkgPath()

			// Check if it's from preveltekit package
			if pkgPath != "preveltekit" {
				continue
			}

			switch elemName {
			case "LocalStore":
				f.Set(reflect.ValueOf(&LocalStore{Store: New("")}))
			case "Store":
				// Need to determine the type parameter - check the field's type
				// Since we can't easily get generic params, use the Elem().Field approach
				// or just try to create based on common types
				// For now, create a Store[string] as default and let OnCreate override
				// Actually, we need to check the actual generic type
				// Use reflection to check the type parameter
				initGenericStore(f, ft)
			case "List":
				initGenericList(f, ft)
			}
		}
	}
}

// initGenericStore initializes a generic Store field based on its type parameter.
func initGenericStore(f reflect.Value, ft reflect.Type) {
	// Create a zero value of the pointer type and check what methods it has
	// This is tricky with generics - we need to use type assertions at runtime

	// Get the full type string which includes generic params in some cases
	// For WASM, the type shows as *preveltekit.Store without params
	// We need to infer from the field usage or use a different approach

	// Try to create each type and see which one is assignable
	stringStore := New("")
	if reflect.TypeOf(stringStore).AssignableTo(ft) {
		f.Set(reflect.ValueOf(stringStore))
		return
	}

	intStore := New(0)
	if reflect.TypeOf(intStore).AssignableTo(ft) {
		f.Set(reflect.ValueOf(intStore))
		return
	}

	boolStore := New(false)
	if reflect.TypeOf(boolStore).AssignableTo(ft) {
		f.Set(reflect.ValueOf(boolStore))
		return
	}

	float64Store := New(0.0)
	if reflect.TypeOf(float64Store).AssignableTo(ft) {
		f.Set(reflect.ValueOf(float64Store))
		return
	}
}

// initGenericList initializes a generic List field based on its type parameter.
func initGenericList(f reflect.Value, ft reflect.Type) {
	stringList := NewList[string]()
	if reflect.TypeOf(stringList).AssignableTo(ft) {
		f.Set(reflect.ValueOf(stringList))
		return
	}

	intList := NewList[int]()
	if reflect.TypeOf(intList).AssignableTo(ft) {
		f.Set(reflect.ValueOf(intList))
		return
	}
}

// replaceAll replaces all occurrences of old with new in s
func replaceAll(s, old, new string) string {
	if old == "" {
		return s
	}
	var result []byte
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result = append(result, new...)
			i += len(old)
		} else {
			result = append(result, s[i])
			i++
		}
	}
	return string(result)
}

// evalCondition evaluates a branch condition using structured data.
func evalCondition(branch HydrateIfBranch, components map[string]Component) bool {
	if branch.StoreID == "" {
		return false
	}

	store := resolveStore(branch.StoreID, components)
	if store == nil {
		return false
	}

	if branch.IsBool {
		if s, ok := store.(*Store[bool]); ok {
			return s.Get()
		}
		return false
	}

	// Compare based on operator
	switch s := store.(type) {
	case *Store[int]:
		return compare(s.Get(), branch.Op, atoiSafe(branch.Operand))
	case *Store[string]:
		return compare(s.Get(), branch.Op, branch.Operand)
	case *Store[float64]:
		return compare(s.Get(), branch.Op, atofSafe(branch.Operand))
	}

	return false
}

// compare compares two ordered values with the given operator.
func compare[T int | float64 | string](val T, op string, operand T) bool {
	switch op {
	case "==":
		return val == operand
	case "!=":
		return val != operand
	case ">=":
		return val >= operand
	case ">":
		return val > operand
	case "<=":
		return val <= operand
	case "<":
		return val < operand
	}
	return false
}

// compareBool compares two boolean values with the given operator.
func compareBool(val bool, op string, operand bool) bool {
	switch op {
	case "==":
		return val == operand
	case "!=":
		return val != operand
	}
	return false
}

func atoiSafe(s string) int {
	n := 0
	neg := false
	for i, c := range s {
		if c == '-' && i == 0 {
			neg = true
			continue
		}
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	if neg {
		return -n
	}
	return n
}

func atofSafe(s string) float64 {
	// Simple float parsing
	var result float64
	var decimal float64 = 1
	neg := false
	afterDot := false

	for i, c := range s {
		if c == '-' && i == 0 {
			neg = true
			continue
		}
		if c == '.' {
			afterDot = true
			continue
		}
		if c < '0' || c > '9' {
			break
		}
		if afterDot {
			decimal *= 10
			result += float64(c-'0') / decimal
		} else {
			result = result*10 + float64(c-'0')
		}
	}
	if neg {
		return -result
	}
	return result
}

// subscribeToStore subscribes a callback to store changes.
func subscribeToStore(store any, callback func()) {
	switch s := store.(type) {
	case *Store[int]:
		s.OnChange(func(_ int) { callback() })
	case *Store[string]:
		s.OnChange(func(_ string) { callback() })
	case *Store[bool]:
		s.OnChange(func(_ bool) { callback() })
	case *Store[float64]:
		s.OnChange(func(_ float64) { callback() })
	}
}

// storeToString returns the string representation of a store's current value.
func storeToString(store any) string {
	switch s := store.(type) {
	case *Store[string]:
		return s.Get()
	case *Store[int]:
		return intToStr(s.Get())
	case *Store[bool]:
		if s.Get() {
			return "true"
		}
		return "false"
	case *Store[float64]:
		// Simple float formatting (avoid fmt)
		return floatToStr(s.Get())
	}
	return ""
}

// floatToStr converts float64 to string without fmt package.
func floatToStr(f float64) string {
	if f == 0 {
		return "0"
	}
	neg := f < 0
	if neg {
		f = -f
	}
	// Integer part
	intPart := int(f)
	fracPart := f - float64(intPart)
	result := intToStr(intPart)
	// Fractional part (up to 6 digits)
	if fracPart > 0.0000001 {
		result += "."
		for i := 0; i < 6 && fracPart > 0.0000001; i++ {
			fracPart *= 10
			digit := int(fracPart)
			result += string(byte('0' + digit))
			fracPart -= float64(digit)
		}
	}
	if neg {
		return "-" + result
	}
	return result
}

// HydrateBindings is the JSON representation of bindings for WASM.
type HydrateBindings struct {
	TextBindings   []HydrateTextBinding   `json:"TextBindings"`
	Events         []HydrateEvent         `json:"Events"`
	IfBlocks       []HydrateIfBlock       `json:"IfBlocks"`
	EachBlocks     []HydrateEachBlock     `json:"EachBlocks"`
	InputBindings  []HydrateInputBinding  `json:"InputBindings"`
	ClassBindings  []HydrateClassBinding  `json:"ClassBindings"`
	ShowIfBindings []HydrateShowIfBinding `json:"ShowIfBindings"`
	AttrBindings   []HydrateAttrBinding   `json:"AttrBindings"`
	RouterBindings []HydrateRouterBinding `json:"RouterBindings"`
}

type HydrateTextBinding struct {
	MarkerID string `json:"marker_id"`
	StoreID  string `json:"store_id"`
	IsHTML   bool   `json:"is_html"`
}

type HydrateEvent struct {
	ElementID string   `json:"ElementID"`
	Event     string   `json:"Event"`
	Modifiers []string `json:"Modifiers"`
}

type HydrateIfBlock struct {
	MarkerID     string            `json:"MarkerID"`
	Branches     []HydrateIfBranch `json:"Branches"`
	ElseHTML     string            `json:"ElseHTML"`
	ElseBindings *HydrateBindings  `json:"ElseBindings,omitempty"`
	Deps         []string          `json:"Deps"`
}

type HydrateIfBranch struct {
	CondExpr string           `json:"CondExpr"`
	HTML     string           `json:"HTML"`
	Bindings *HydrateBindings `json:"Bindings,omitempty"`
	StoreID  string           `json:"store_id,omitempty"`
	Op       string           `json:"op,omitempty"`
	Operand  string           `json:"operand,omitempty"`
	IsBool   bool             `json:"is_bool,omitempty"`
}

type HydrateInputBinding struct {
	ElementID string `json:"element_id"`
	StoreID   string `json:"store_id"`
	BindType  string `json:"bind_type"`
}

type HydrateClassBinding struct {
	ElementID string   `json:"element_id"`
	ClassName string   `json:"class_name"`
	CondExpr  string   `json:"cond_expr"`
	Op        string   `json:"op"`
	Operand   string   `json:"operand"`
	Deps      []string `json:"deps"`
}

type HydrateShowIfBinding struct {
	ElementID string   `json:"element_id"`
	StoreID   string   `json:"store_id"`
	Op        string   `json:"op"`
	Operand   string   `json:"operand"`
	IsBool    bool     `json:"is_bool"`
	Deps      []string `json:"deps"`
}

type HydrateAttrBinding struct {
	ElementID string   `json:"element_id"`
	AttrName  string   `json:"attr_name"`
	Template  string   `json:"template"`
	StoreIDs  []string `json:"store_ids"`
}

type HydrateEachBlock struct {
	MarkerID string `json:"MarkerID"`
	ListID   string `json:"ListID"`
	ItemVar  string `json:"ItemVar"`
	IndexVar string `json:"IndexVar"`
}

type HydrateRouterBinding struct {
	StoreID string `json:"store_id"`
}
