//go:build js && wasm

package preveltekit

import (
	"syscall/js"
)

// Track which if-blocks have been set up to avoid duplicates
var setupIfBlocks = make(map[string]bool)

// Track which each-blocks have been set up to avoid duplicates
var setupEachBlocks = make(map[string]bool)

// Track which route-blocks have been set up to avoid duplicates
var setupRouteBlocks = make(map[string]bool)

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

// Hydrate sets up DOM bindings for reactivity.
// With the ID-based system, stores and handlers register themselves with unique IDs.
// We still need the bindings JSON for If-blocks, Each-blocks, and AttrCond bindings.
func Hydrate(app Component) {
	// Create fresh app instance with initialized stores (this registers stores)
	if hn, ok := app.(HasNew); ok {
		app = hn.New()
	}

	// Discover children from routes and initialize them (registers their stores)
	children := make(map[string]Component)
	if hr, ok := app.(HasRoutes); ok {
		for _, route := range hr.Routes() {
			if route.Component != nil {
				children[route.Path] = route.Component
			}
		}
	}

	// Call Render() on all components to register handlers via WithOn().
	// Must recurse into nested ComponentNodes so their WithOn calls fire too.
	renderRecursive(app)
	for _, child := range children {
		renderRecursive(child)
	}

	// Use the full hydration which handles If-blocks, Each-blocks, etc.
	hydrateWASM(app, children)
}

// renderRecursive calls Render() on a component and recursively on any
// nested ComponentNodes. This ensures all WithOn() calls fire to register
// handlers in the global registry.
func renderRecursive(comp Component) {
	node := comp.Render()
	walkNodeForComponents(node)
}

// walkNodeForComponents walks a Node tree and calls renderRecursive on
// any ComponentNode instances found.
func walkNodeForComponents(n Node) {
	if n == nil {
		return
	}
	switch node := n.(type) {
	case *HtmlNode:
		for _, part := range node.Parts {
			if child, ok := part.(Node); ok {
				walkNodeForComponents(child)
			}
		}
	case *Fragment:
		for _, child := range node.Children {
			walkNodeForComponents(child)
		}
	case *IfNode:
		for _, branch := range node.Branches {
			for _, child := range branch.Children {
				walkNodeForComponents(child)
			}
		}
		for _, child := range node.ElseNode {
			walkNodeForComponents(child)
		}
	case *ComponentNode:
		if comp, ok := node.Instance.(Component); ok {
			renderRecursive(comp)
		}
	}
}

// hydrateWASM sets up DOM bindings from the embedded bindings JSON.
func hydrateWASM(app Component, children map[string]Component) {

	// Get bindings from global variable (set by CLI-generated code)
	bindingsJS := js.Global().Get("_preveltekit_bindings")
	if bindingsJS.IsUndefined() || bindingsJS.IsNull() {
		runLifecycle(app, children)
		select {}
		return
	}

	bindingsJSON := bindingsJS.String()

	bindings := parseBindings(bindingsJSON)
	if bindings == nil {
		runLifecycle(app, children)
		select {}
		return
	}

	// Build component map: "component" -> app, "basics" -> child, etc.
	components := map[string]Component{
		"component": app,
	}
	for path, child := range children {
		name := trimPrefix(path, "/")
		components[name] = child
	}

	// Run OnCreate phase
	runOnCreate(app, children)

	// Apply all bindings
	cleanup := &Cleanup{}
	applyBindings(bindings, components, cleanup)

	// Run OnMount for all children FIRST to initialize their stores (e.g., CurrentTab = "home")
	// But NOT app.OnMount yet - that starts the router which changes CurrentComponent
	for _, child := range children {
		if om, ok := child.(HasOnMount); ok {
			om.OnMount()
		}
	}

	// NOW run app's OnMount which starts the router
	// Router will handle component store binding for SPA navigation
	if om, ok := app.(HasOnMount); ok {
		om.OnMount()
	}

	// Apply route blocks AFTER OnMount, since the router's path store
	// is created during OnMount -> NewRouter -> New(id+".path", "")
	for _, rb := range bindings.RouteBlocks {
		bindRouteBlock(rb, components)
	}

	// Keep WASM running
	select {}
}

// runOnCreate injects styles for app and all children.
// Note: New() was already called to initialize stores.
func runOnCreate(app Component, children map[string]Component) {
	if hs, ok := app.(HasStyle); ok {
		InjectStyle("app", hs.Style())
	}
	// Inject styles for all children
	for path, child := range children {
		if hs, ok := child.(HasStyle); ok {
			InjectStyle(trimPrefix(path, "/"), hs.Style())
		}
	}
}

// runOnMount calls OnMount for app and all children.
func runOnMount(app Component, children map[string]Component) {
	if om, ok := app.(HasOnMount); ok {
		om.OnMount()
	}
	for _, child := range children {
		if om, ok := child.(HasOnMount); ok {
			om.OnMount()
		}
	}
}

// runLifecycle runs full lifecycle (OnCreate + styles + OnMount) for all components.
func runLifecycle(app Component, children map[string]Component) {
	runOnCreate(app, children)
	runOnMount(app, children)
}

// resolveStore looks up a store by ID from the global registry.
// Stores register themselves via New(0) which auto-generates an ID and adds them to storeRegistry.
func resolveStore(storeID string, components map[string]Component) any {
	store := GetStore(storeID)
	if store == nil {
		js.Global().Get("console").Call("log", "[DEBUG] resolveStore: store not found in registry:", storeID)
	}
	return store
}

// getHandler looks up a handler by ID from the global handler registry.
// Handlers are registered via WithOn() which calls RegisterHandler().
func getHandler(elementID string) func() {
	return GetHandler(elementID)
}

// bindTextDynamic sets up a text binding for any store type.
func bindTextDynamic(markerID string, store any, isHTML bool) {
	js.Global().Get("console").Call("log", "[DEBUG] bindTextDynamic:", markerID, "store is nil:", store == nil)
	switch s := store.(type) {
	case *Store[string]:
		bindText(markerID, s, isHTML)
	case *Store[int]:
		js.Global().Get("console").Call("log", "[DEBUG] bindTextDynamic: binding int store, current value:", s.Get())
		bindText(markerID, s, isHTML)
	case *Store[bool]:
		bindText(markerID, s, isHTML)
	case *Store[float64]:
		bindText(markerID, s, isHTML)
	default:
		js.Global().Get("console").Call("log", "[DEBUG] bindTextDynamic: unknown store type")
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

	// Track which branch is currently active (-1 = else, 0+ = branch index)
	currentBranchIdx := -2 // -2 = not yet determined

	// Evaluate which branch is active and update content
	updateIfBlock := func() {
		var activeHTML string
		var activeBindings *HydrateBindings
		activeBranchIdx := -1 // -1 = else branch

		for i := 0; i < len(ifb.Branches); i++ {
			branch := ifb.Branches[i]
			if evalCondition(branch, components) {
				activeHTML = branch.HTML
				activeBindings = branch.Bindings
				activeBranchIdx = i
				break
			}
		}

		if activeBranchIdx == -1 {
			activeHTML = ifb.ElseHTML
			activeBindings = ifb.ElseBindings
		}

		// Skip re-rendering if the active branch hasn't changed
		// This prevents overwriting list content when only the list items changed
		if currentBranchIdx == activeBranchIdx && currentBranchIdx != -2 {
			return
		}
		currentBranchIdx = activeBranchIdx

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
	subscribedAny := false
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
			subscribedAny = true
		}
	}

	// Call updateIfBlock to sync DOM with current state.
	// But only if we successfully subscribed to at least one dep store.
	// If no deps could be resolved (e.g., internal stores with no ID),
	// trust the SSR-rendered content rather than replacing it with a wrong branch.
	if subscribedAny {
		updateIfBlock()
	}
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
	// Recursively clear route-block markers
	for _, rb := range bindings.RouteBlocks {
		for _, branch := range rb.Branches {
			if branch.Bindings != nil {
				clearBoundMarkers(branch.Bindings)
			}
		}
		delete(setupRouteBlocks, rb.MarkerID)
	}

	// NOTE: Do NOT clear setupEachBlocks here. Each-blocks subscribe to list.OnChange
	// which persists across if-block changes. Clearing would cause duplicate subscriptions.
	// The list callbacks will re-find the marker when needed.
}

// bindRouteBlock sets up a route-block with pre-baked HTML swap on path change.
// Modeled on bindIfBlock but with path-matching instead of store-condition evaluation.
func bindRouteBlock(rb HydrateRouteBlock, components map[string]Component) {
	if setupRouteBlocks[rb.MarkerID] {
		return
	}
	setupRouteBlocks[rb.MarkerID] = true

	// Find existing SSR content (the <span> before the comment marker)
	currentEl := FindExistingIfContent(rb.MarkerID)

	currentCleanup := &Cleanup{}

	currentBranchPath := ""
	firstCall := true

	updateRouteBlock := func() {
		js.Global().Get("console").Call("log", "[ROUTE] updateRouteBlock called, firstCall:", firstCall, "currentBranchPath:", currentBranchPath)

		// Get current path from the path store
		pathStore := resolveStore(rb.PathStoreID, components)
		path := ""
		if ps, ok2 := pathStore.(*Store[string]); ok2 {
			path = ps.Get()
		}
		if path == "" {
			path = js.Global().Get("location").Get("pathname").String()
		}
		if path == "" {
			path = "/"
		}
		js.Global().Get("console").Call("log", "[ROUTE] resolved path:", path)

		// Find matching branch using route matching
		var activeHTML string
		var activeBindings *HydrateBindings
		activePath := ""
		bestSpecificity := -1

		for _, branch := range rb.Branches {
			_, specificity, ok3 := matchRoute(branch.Path, path)
			js.Global().Get("console").Call("log", "[ROUTE] branch:", branch.Path, "match:", ok3, "specificity:", specificity, "hasBindings:", branch.Bindings != nil)
			if branch.Bindings != nil {
				js.Global().Get("console").Call("log", "[ROUTE]   branch events:", len(branch.Bindings.Events), "text:", len(branch.Bindings.TextBindings))
			}
			if ok3 && specificity > bestSpecificity {
				activeHTML = branch.HTML
				activeBindings = branch.Bindings
				activePath = branch.Path
				bestSpecificity = specificity
			}
		}

		js.Global().Get("console").Call("log", "[ROUTE] activePath:", activePath, "activeHTML len:", len(activeHTML), "hasActiveBindings:", activeBindings != nil)

		// Skip if same route is already active (but not on first call)
		if activePath == currentBranchPath && !firstCall {
			js.Global().Get("console").Call("log", "[ROUTE] SKIP: same route already active")
			return
		}

		if firstCall {
			firstCall = false
			currentBranchPath = activePath
			js.Global().Get("console").Call("log", "[ROUTE] FIRST CALL: skipping DOM swap, applying bindings only")
			// On initial load, SSR already rendered the correct HTML.
			// Don't swap the DOM, but DO apply the branch's bindings
			// so that event listeners, text bindings, etc. are wired up.
			if activeBindings != nil {
				js.Global().Get("console").Call("log", "[ROUTE] applying initial bindings, events:", len(activeBindings.Events))
				clearBoundMarkers(activeBindings)
				applyBindings(activeBindings, components, currentCleanup)
			} else {
				js.Global().Get("console").Call("log", "[ROUTE] WARNING: no bindings for initial route")
			}
			return
		}

		js.Global().Get("console").Call("log", "[ROUTE] SWAP: replacing DOM, old:", currentBranchPath, "new:", activePath)
		currentBranchPath = activePath

		// Swap HTML (same mechanism as IfBlock)
		currentEl = FindExistingIfContent(rb.MarkerID)
		currentEl = ReplaceContent(rb.MarkerID, currentEl, activeHTML)

		// Release old bindings, apply new
		currentCleanup.Release()
		currentCleanup = &Cleanup{}

		if activeBindings != nil {
			js.Global().Get("console").Call("log", "[ROUTE] applying swap bindings, events:", len(activeBindings.Events))
			clearBoundMarkers(activeBindings)
			applyBindings(activeBindings, components, currentCleanup)
		} else {
			js.Global().Get("console").Call("log", "[ROUTE] WARNING: no bindings for swapped route")
		}
	}

	// Subscribe to the router's currentPath store
	pathStore := resolveStore(rb.PathStoreID, components)
	if pathStore != nil {
		subscribeToStore(pathStore, updateRouteBlock)
	}

	// Initial sync (handles SSR → WASM handoff)
	updateRouteBlock()
}

// applyBindings applies all bindings from a HydrateBindings struct to the DOM.
func applyBindings(bindings *HydrateBindings, components map[string]Component, cleanup *Cleanup) {
	js.Global().Get("console").Call("log", "[DEBUG] applyBindings called")
	js.Global().Get("console").Call("log", "[DEBUG] bindings.Events count:", len(bindings.Events))
	js.Global().Get("console").Call("log", "[DEBUG] bindings.TextBindings count:", len(bindings.TextBindings))

	// Text bindings
	for _, tb := range bindings.TextBindings {
		js.Global().Get("console").Call("log", "[DEBUG] text binding:", tb.MarkerID, "storeID:", tb.StoreID)
		store := resolveStore(tb.StoreID, components)
		if store != nil {
			bindTextDynamic(tb.MarkerID, store, tb.IsHTML)
		} else {
			js.Global().Get("console").Call("log", "[DEBUG] text binding: store not resolved for", tb.StoreID)
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
		js.Global().Get("console").Call("log", "[DEBUG] binding event:", ev.ElementID, ev.Event, "handler found:", handler != nil)
		if handler != nil {
			// Wrap handler to add debug logging
			elementID := ev.ElementID
			event := ev.Event
			originalHandler := handler
			wrappedHandler := func() {
				js.Global().Get("console").Call("log", "[DEBUG] Handler executing for:", elementID, event)
				originalHandler()
				js.Global().Get("console").Call("log", "[DEBUG] Handler completed for:", elementID, event)
			}
			bindEventDynamic(cleanup, ev.ElementID, ev.Event, wrappedHandler)
		}
	}

	// Nested if-blocks
	for _, ifb := range bindings.IfBlocks {
		bindIfBlock(ifb, components)
	}

	// Attr bindings (dynamic attributes like data-type)
	for _, ab := range bindings.AttrBindings {
		bindAttr(ab, components)
	}

	// AttrCond bindings (conditional attributes from HtmlNode.AttrIf())
	for _, acb := range bindings.AttrCondBindings {
		bindAttrCondBinding(acb, components)
	}

	// Each block bindings (list iteration)
	for _, eb := range bindings.EachBlocks {
		bindEachBlock(eb, components)
	}

	// NOTE: Route blocks are NOT applied here. They require the router's path store
	// which is created during OnMount (after applyBindings runs). Route blocks are
	// applied explicitly after OnMount in hydrateWASM.

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

// bindAttrCondBinding sets up a conditional attribute binding from HtmlNode.AttrIf().
func bindAttrCondBinding(acb HydrateAttrCondBinding, components map[string]Component) {
	el := GetEl(acb.ElementID)
	if !ok(el) {
		return
	}

	if len(acb.Deps) == 0 {
		return
	}

	// Resolve the condition store (first dep is always the condition)
	condStore := resolveStore(acb.Deps[0], components)
	if condStore == nil {
		return
	}

	// Resolve true/false value stores if they're dynamic
	var trueStore, falseStore any
	if acb.TrueStoreID != "" {
		trueStore = resolveStore(acb.TrueStoreID, components)
	}
	if acb.FalseStoreID != "" {
		falseStore = resolveStore(acb.FalseStoreID, components)
	}

	// Function to evaluate condition and update attribute
	updateAttr := func() {
		// Evaluate condition
		var active bool
		if acb.IsBool {
			if s, ok := condStore.(*Store[bool]); ok {
				active = s.Get()
			}
		} else if acb.Op != "" {
			switch s := condStore.(type) {
			case *Store[int]:
				active = compare(s.Get(), acb.Op, atoiSafe(acb.Operand))
			case *Store[string]:
				active = compare(s.Get(), acb.Op, acb.Operand)
			case *Store[float64]:
				active = compare(s.Get(), acb.Op, atofSafe(acb.Operand))
			case *Store[bool]:
				operandBool := acb.Operand == "true"
				active = compareBool(s.Get(), acb.Op, operandBool)
			}
		}

		// Determine value to use
		var value string
		if active {
			if trueStore != nil {
				value = storeToString(trueStore)
			} else {
				value = acb.TrueValue
			}
		} else {
			if falseStore != nil {
				value = storeToString(falseStore)
			} else {
				value = acb.FalseValue
			}
		}

		// Special handling for class attribute - toggle instead of replace
		if acb.AttrName == "class" {
			classList := el.Get("classList")
			if active && acb.TrueValue != "" {
				classList.Call("add", acb.TrueValue)
			} else if acb.TrueValue != "" {
				classList.Call("remove", acb.TrueValue)
			}
			if !active && acb.FalseValue != "" {
				classList.Call("add", acb.FalseValue)
			} else if acb.FalseValue != "" {
				classList.Call("remove", acb.FalseValue)
			}
		} else {
			// For other attributes, set the value directly
			if value != "" {
				el.Call("setAttribute", acb.AttrName, value)
			} else {
				el.Call("removeAttribute", acb.AttrName)
			}
		}
	}

	// Initial update
	updateAttr()

	// Subscribe to condition store changes
	subscribeToStore(condStore, updateAttr)

	// Subscribe to value store changes if dynamic
	if trueStore != nil {
		subscribeToStore(trueStore, updateAttr)
	}
	if falseStore != nil {
		subscribeToStore(falseStore, updateAttr)
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

	// Extract item ID prefix from marker (e.g., "lists_e0" -> "lists_")
	itemIDPrefix := eb.MarkerID[:len(eb.MarkerID)-1]
	if len(itemIDPrefix) > 0 && itemIDPrefix[len(itemIDPrefix)-1] == 'e' {
		itemIDPrefix = itemIDPrefix[:len(itemIDPrefix)-1]
	}

	// Subscribe to list changes and re-render
	// Pass markerID so we can re-acquire the parent each time (handles DOM replacement from if-blocks)
	switch list := listAny.(type) {
	case *List[string]:
		bindListItems(list, eb.MarkerID, itemIDPrefix, escapeHTMLWasm)
	case *List[int]:
		bindListItems(list, eb.MarkerID, itemIDPrefix, intToStr)
	}

	_ = comp
}

// bindListItems sets up list rendering and subscribes to changes.
func bindListItems[T comparable](list *List[T], markerID string, itemIDPrefix string, format func(T) string) {
	renderItems := func(items []T) {
		var html string
		for i, item := range items {
			html += `<span id="` + itemIDPrefix + `_` + intToStr(i) + `"><li><span class="index">` + intToStr(i) + `</span> ` + format(item) + `</li></span>`
		}
		// Add marker comment at the end so it survives innerHTML replacement
		html += `<!--` + markerID + `-->`
		// Re-acquire parent from marker each time to handle DOM replacement (e.g., from if-blocks)
		marker := FindComment(markerID)
		if marker.IsNull() {
			return
		}
		parent := marker.Get("parentNode")
		if !parent.IsNull() && parent.Get("nodeType").Int() == 1 {
			parent.Set("innerHTML", html)
		}
	}

	// Don't render initially - SSR already rendered the items
	// Just subscribe to changes for future updates
	list.OnChange(renderItems)
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

// isStore checks if a value is a Store type (used for counter synchronization).
func isStore(v any) bool {
	switch v.(type) {
	case *Store[string], *Store[int], *Store[bool], *Store[float64]:
		return true
	}
	return false
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
	TextBindings        []HydrateTextBinding        `json:"TextBindings"`
	Events              []HydrateEvent              `json:"Events"`
	IfBlocks            []HydrateIfBlock            `json:"IfBlocks"`
	EachBlocks          []HydrateEachBlock          `json:"EachBlocks"`
	InputBindings       []HydrateInputBinding       `json:"InputBindings"`
	AttrBindings        []HydrateAttrBinding        `json:"AttrBindings"`
	AttrCondBindings    []HydrateAttrCondBinding    `json:"AttrCondBindings"`
	ComponentContainers []HydrateComponentContainer `json:"ComponentContainers,omitempty"`
	RouteBlocks         []HydrateRouteBlock         `json:"RouteBlocks,omitempty"`
}

// HydrateComponentContainer maps a route group ID to its DOM container ID
type HydrateComponentContainer struct {
	ID          string `json:"ID"`
	ContainerID string `json:"ContainerID"`
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

type HydrateAttrBinding struct {
	ElementID string   `json:"element_id"`
	AttrName  string   `json:"attr_name"`
	Template  string   `json:"template"`
	StoreIDs  []string `json:"store_ids"`
}

// HydrateAttrCondBinding represents a conditional attribute binding for WASM.
// Used by HtmlNode.AttrIf() for conditional attribute values.
type HydrateAttrCondBinding struct {
	ElementID    string   `json:"element_id"`
	AttrName     string   `json:"attr_name"`
	TrueValue    string   `json:"true_value"`
	FalseValue   string   `json:"false_value,omitempty"`
	TrueStoreID  string   `json:"true_store_id,omitempty"`
	FalseStoreID string   `json:"false_store_id,omitempty"`
	Op           string   `json:"op,omitempty"`
	Operand      string   `json:"operand,omitempty"`
	IsBool       bool     `json:"is_bool,omitempty"`
	Deps         []string `json:"deps,omitempty"`
}

type HydrateEachBlock struct {
	MarkerID string `json:"MarkerID"`
	ListID   string `json:"ListID"`
	ItemVar  string `json:"ItemVar"`
	IndexVar string `json:"IndexVar"`
}

// HydrateRouteBlock is the WASM-side representation of a RouteBlock.
// Like HydrateIfBlock but with path-based conditions instead of store conditions.
type HydrateRouteBlock struct {
	MarkerID    string               `json:"MarkerID"`
	PathStoreID string               `json:"PathStoreID"`
	Branches    []HydrateRouteBranch `json:"Branches"`
}

// HydrateRouteBranch represents one route's pre-baked content.
type HydrateRouteBranch struct {
	Path     string           `json:"Path"`
	HTML     string           `json:"HTML"`
	Bindings *HydrateBindings `json:"Bindings,omitempty"`
}
