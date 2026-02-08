//go:build js && wasm

package preveltekit

import (
	"syscall/js"
)

// Track which if-blocks have been set up to avoid duplicates
var setupIfBlocks = make(map[string]bool)

// Track which each-blocks have been set up to avoid duplicates
var setupEachBlocks = make(map[string]bool)

// Track which component-blocks have been set up to avoid duplicates
var setupComponentBlocks = make(map[string]bool)

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

	// Call OnMount before Render to match SSR order.
	// This creates the router (which calls WithOptions on Store[Component]),
	// so options are populated when renderRecursive walks the tree.
	if om, ok := app.(HasOnMount); ok {
		om.OnMount()
	}

	// Call Render() on all components to register handlers via On().
	// Must recurse into nested ComponentNodes so their On calls fire too.
	// Store[Component] options are also walked to keep handler counter in sync with SSR.
	renderRecursive(app)

	// Use the full hydration which handles If-blocks, Each-blocks, etc.
	hydrateWASM(app, children)
}

// renderRecursive calls Render() on a component and recursively on any
// nested ComponentNodes. This ensures all On() calls fire to register
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
			// Store[Component] with options: render all options to keep
			// handler counter in sync with SSR (which renders all branches)
			if cs, ok := part.(*Store[Component]); ok {
				seen := make(map[string]bool)
				for _, opt := range cs.Options() {
					if comp, ok := opt.(Component); ok {
						name := componentName(comp)
						if !seen[name] {
							seen[name] = true
							renderRecursive(comp)
						}
					}
				}
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

	// app.OnMount() already called in Hydrate() before renderRecursive,
	// so Store[Component] options are populated for handler counter sync.

	// Apply component blocks (options already registered via OnMount -> NewRouter -> WithOptions)
	for _, cb := range bindings.ComponentBlocks {
		bindComponentBlock(cb, components)
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

// runOnMount calls OnMount for all children.
// Note: app.OnMount() is called earlier in Hydrate() before renderRecursive,
// so it is not called here to avoid double invocation.
func runOnMount(app Component, children map[string]Component) {
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
// Handlers are registered via On() which calls RegisterHandler().
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
	// Recursively clear component-block markers
	for _, cb := range bindings.ComponentBlocks {
		for _, branch := range cb.Branches {
			if branch.Bindings != nil {
				clearBoundMarkers(branch.Bindings)
			}
		}
		delete(setupComponentBlocks, cb.MarkerID)
	}

	// NOTE: Do NOT clear setupEachBlocks here. Each-blocks subscribe to list.OnChange
	// which persists across if-block changes. Clearing would cause duplicate subscriptions.
	// The list callbacks will re-find the marker when needed.
}

// bindComponentBlock sets up a component-block with pre-baked HTML swap on store change.
// Subscribes to the Store[Component] and swaps branches by component type name.
func bindComponentBlock(cb HydrateComponentBlock, components map[string]Component) {
	if setupComponentBlocks[cb.MarkerID] {
		return
	}
	setupComponentBlocks[cb.MarkerID] = true

	// Find existing SSR content (the <span> before the comment marker)
	currentEl := FindExistingIfContent(cb.MarkerID)

	currentCleanup := &Cleanup{}

	currentBranchName := ""
	firstCall := true

	updateComponentBlock := func() {
		// Get current component from the store
		compStore := resolveStore(cb.StoreID, components)
		if compStore == nil {
			js.Global().Get("console").Call("log", "[COMP] store not found:", cb.StoreID)
			return
		}

		// Get the component's type name
		var activeName string
		if cs, ok2 := compStore.(*Store[Component]); ok2 {
			comp := cs.Get()
			if comp != nil {
				activeName = componentName(comp)
			}
		}

		js.Global().Get("console").Call("log", "[COMP] updateComponentBlock called, firstCall:", firstCall, "activeName:", activeName, "currentBranch:", currentBranchName)

		// Find matching branch by component name
		var activeHTML string
		var activeBindings *HydrateBindings
		for _, branch := range cb.Branches {
			if branch.Name == activeName {
				activeHTML = branch.HTML
				activeBindings = branch.Bindings
				break
			}
		}

		// Skip if same branch is already active (but not on first call)
		if activeName == currentBranchName && !firstCall {
			js.Global().Get("console").Call("log", "[COMP] SKIP: same component already active")
			return
		}

		if firstCall {
			firstCall = false
			currentBranchName = activeName
			js.Global().Get("console").Call("log", "[COMP] FIRST CALL: skipping DOM swap, applying bindings only")
			if activeBindings != nil {
				clearBoundMarkers(activeBindings)
				applyBindings(activeBindings, components, currentCleanup)
			}
			return
		}

		js.Global().Get("console").Call("log", "[COMP] SWAP: replacing DOM, old:", currentBranchName, "new:", activeName)
		currentBranchName = activeName

		// Swap HTML (same mechanism as IfBlock)
		currentEl = FindExistingIfContent(cb.MarkerID)
		currentEl = ReplaceContent(cb.MarkerID, currentEl, activeHTML)

		// Release old bindings, apply new
		currentCleanup.Release()
		currentCleanup = &Cleanup{}

		if activeBindings != nil {
			clearBoundMarkers(activeBindings)
			applyBindings(activeBindings, components, currentCleanup)
		}
	}

	// Subscribe to the component store
	compStore := resolveStore(cb.StoreID, components)
	if compStore != nil {
		subscribeToStore(compStore, updateComponentBlock)
	}

	// Initial sync (handles SSR → WASM handoff)
	updateComponentBlock()
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
	case *Store[Component]:
		s.OnChange(func(_ Component) { callback() })
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
	TextBindings     []HydrateTextBinding     `json:"TextBindings"`
	Events           []HydrateEvent           `json:"Events"`
	IfBlocks         []HydrateIfBlock         `json:"IfBlocks"`
	EachBlocks       []HydrateEachBlock       `json:"EachBlocks"`
	InputBindings    []HydrateInputBinding    `json:"InputBindings"`
	AttrBindings     []HydrateAttrBinding     `json:"AttrBindings"`
	AttrCondBindings []HydrateAttrCondBinding `json:"AttrCondBindings"`
	ComponentBlocks  []HydrateComponentBlock  `json:"ComponentBlocks,omitempty"`
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

// HydrateComponentBlock is the WASM-side representation of a ComponentBlock.
// Like HydrateIfBlock but keyed by component type name instead of store conditions.
type HydrateComponentBlock struct {
	MarkerID string                   `json:"MarkerID"`
	StoreID  string                   `json:"StoreID"`
	Branches []HydrateComponentBranch `json:"Branches"`
}

// HydrateComponentBranch represents one component's pre-baked content.
type HydrateComponentBranch struct {
	Name     string           `json:"Name"`
	HTML     string           `json:"HTML"`
	Bindings *HydrateBindings `json:"Bindings,omitempty"`
}
