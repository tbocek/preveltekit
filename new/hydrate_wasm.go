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
	// so options are populated when walkNodeForComponents walks the tree.
	if om, ok := app.(HasOnMount); ok {
		om.OnMount()
	}

	// Create app scope before walkNodeForComponents to match SSR order
	// (SSR creates the app scope before nodeToHTML)
	if _, ok := app.(HasStyle); ok {
		GetOrCreateScope("app")
	}

	// Call Render() on all components to register handlers via On().
	// Must recurse into nested ComponentNodes so their On calls fire too.
	// Store[Component] options are also walked to keep handler/scope counter in sync with SSR.
	walkNodeForComponents(app.Render())

	// Use the full hydration which handles If-blocks, Each-blocks, etc.
	hydrateWASM(app, children)
}

// walkNodeForComponents walks a Node tree, calling Render() on any nested
// ComponentNodes so their On() calls fire to register handlers in the global registry.
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
							// Keep scope counter in sync with SSR's ComponentBlock rendering
							if _, ok := comp.(HasStyle); ok {
								GetOrCreateScope(name)
							}
							walkNodeForComponents(comp.Render())
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
			// Keep scope counter in sync with SSR's ComponentNode.ToHTML
			if _, ok := comp.(HasStyle); ok {
				GetOrCreateScope(node.Name)
			}
			walkNodeForComponents(comp.Render())
		}
	}
}

// hydrateWASM sets up DOM bindings from the embedded bindings binary.
func hydrateWASM(app Component, children map[string]Component) {

	// Get bindings from global variable (set as Uint8Array by loader)
	bindingsJS := js.Global().Get("_preveltekit_bindings")
	if bindingsJS.IsUndefined() || bindingsJS.IsNull() {
		runOnMount(children)
		select {}
	}

	// Copy Uint8Array into Go []byte
	length := bindingsJS.Get("byteLength").Int()
	data := make([]byte, length)
	js.CopyBytesToGo(data, bindingsJS)

	bindings := decodeBindings(data)
	if bindings == nil {
		runOnMount(children)
		select {}
	}

	// Apply all bindings
	cleanup := &Cleanup{}
	applyBindings(bindings, cleanup)

	// Run OnMount for all children to initialize their stores (e.g., CurrentTab = "home").
	// app.OnMount() already called in Hydrate() before walkNodeForComponents.
	runOnMount(children)

	// Apply component blocks (options already registered via OnMount -> NewRouter -> WithOptions)
	for _, cb := range bindings.ComponentBlocks {
		bindComponentBlock(cb)
	}

	// Keep WASM running
	select {}
}

// runOnMount calls OnMount for all children.
// Note: app.OnMount() is called earlier in Hydrate() before walkNodeForComponents,
// so it is not called here to avoid double invocation.
func runOnMount(children map[string]Component) {
	for _, child := range children {
		if om, ok := child.(HasOnMount); ok {
			om.OnMount()
		}
	}
}

// bindTextDynamic sets up a text binding for any store type.
func bindTextDynamic(markerID string, store any, isHTML bool) {
	switch s := store.(type) {
	case *Store[string]:
		bindTextStore(markerID, s, isHTML)
	case *Store[int]:
		bindTextStore(markerID, s, isHTML)
	case *Store[bool]:
		bindTextStore(markerID, s, isHTML)
	case *Store[float64]:
		bindTextStore(markerID, s, isHTML)
	}
}

// bindTextStore subscribes to a store and updates text content between marker pairs.
func bindTextStore[T any](markerID string, store Bindable[T], isHTML bool) {
	update := func(v T) {
		s := toString(v)
		if !isHTML {
			s = escapeHTML(s)
		}
		replaceMarkerContent(markerID, s)
	}
	store.OnChange(update)
	// Set current value immediately (store may have been updated before binding)
	update(store.Get())
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

// bindIfBlock sets up an if-block with reactive condition evaluation.
func bindIfBlock(ifb HydrateIfBlock) {
	// Skip if already set up (prevents duplicate setup from nested if-blocks)
	if setupIfBlocks[ifb.MarkerID] {
		return
	}
	setupIfBlocks[ifb.MarkerID] = true

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
			if evalCondition(branch) {
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
		if currentBranchIdx == activeBranchIdx && currentBranchIdx != -2 {
			return
		}
		currentBranchIdx = activeBranchIdx
		replaceMarkerContent(ifb.MarkerID, activeHTML)

		currentCleanup.Release()
		currentCleanup = &Cleanup{}

		if activeBindings != nil {
			clearBoundMarkers(activeBindings)
			applyBindings(activeBindings, currentCleanup)
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

		store := GetStore(dep)
		if store == nil {
			// Try with component prefix if dep doesn't contain a dot
			hasDot := false
			for i := 0; i < len(dep); i++ {
				if dep[i] == '.' {
					hasDot = true
					break
				}
			}
			if !hasDot {
				store = GetStore("component." + dep)
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

// replaceMarkerContent replaces all DOM nodes between <!--{markerID}s--> and <!--{markerID}-->
// with new HTML content. Used by IfBlocks, ComponentBlocks, and EachBlocks.
func replaceMarkerContent(markerID string, html string) {
	endMarker := FindComment(markerID)
	if endMarker.IsNull() {
		return
	}
	startMarker := FindComment(markerID + "s")
	if startMarker.IsNull() {
		return
	}
	parent := endMarker.Get("parentNode")

	// Remove all nodes between start and end markers
	for {
		next := startMarker.Get("nextSibling")
		if next.IsNull() || next.Equal(endMarker) {
			break
		}
		parent.Call("removeChild", next)
	}

	// Parse new HTML and insert before end marker
	if html != "" {
		tmpl := Document.Call("createElement", "template")
		tmpl.Set("innerHTML", html)
		frag := tmpl.Get("content")
		parent.Call("insertBefore", frag, endMarker)
	}
}

// clearBoundMarkers clears marker tracking for bindings that will be re-applied.
// This is needed when block content is replaced via replaceMarkerContent.
func clearBoundMarkers(bindings *HydrateBindings) {
	if bindings == nil {
		return
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
func bindComponentBlock(cb HydrateComponentBlock) {
	if setupComponentBlocks[cb.MarkerID] {
		return
	}
	setupComponentBlocks[cb.MarkerID] = true

	currentCleanup := &Cleanup{}

	currentBranchName := ""
	firstCall := true

	updateComponentBlock := func() {
		// Get current component from the store
		compStore := GetStore(cb.StoreID)
		if compStore == nil {
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
			return
		}

		if firstCall {
			firstCall = false
			currentBranchName = activeName
			if activeBindings != nil {
				clearBoundMarkers(activeBindings)
				applyBindings(activeBindings, currentCleanup)
			}
			return
		}

		currentBranchName = activeName

		replaceMarkerContent(cb.MarkerID, activeHTML)

		// Release old bindings, apply new
		currentCleanup.Release()
		currentCleanup = &Cleanup{}

		if activeBindings != nil {
			clearBoundMarkers(activeBindings)
			applyBindings(activeBindings, currentCleanup)
		}
	}

	// Subscribe to the component store
	compStore := GetStore(cb.StoreID)
	if compStore != nil {
		subscribeToStore(compStore, updateComponentBlock)
	}

	// Initial sync (handles SSR → WASM handoff)
	updateComponentBlock()
}

// applyBindings applies all bindings from a HydrateBindings struct to the DOM.
func applyBindings(bindings *HydrateBindings, cleanup *Cleanup) {
	// Text bindings
	for _, tb := range bindings.TextBindings {
		store := GetStore(tb.StoreID)
		if store != nil {
			bindTextDynamic(tb.MarkerID, store, tb.IsHTML)
		}
	}

	// Input bindings
	for _, ib := range bindings.InputBindings {
		store := GetStore(ib.StoreID)
		if store != nil {
			bindInputDynamic(cleanup, ib.StoreID, store, ib.BindType)
		}
	}

	// Event bindings
	for _, ev := range bindings.Events {
		handler := GetHandler(ev.ElementID)
		if handler != nil {
			BindEvents(cleanup, []Evt{{ev.ElementID, ev.Event, handler}})
		}
	}

	// Nested if-blocks
	for _, ifb := range bindings.IfBlocks {
		bindIfBlock(ifb)
	}

	// Attr bindings (dynamic attributes like data-type)
	for _, ab := range bindings.AttrBindings {
		bindAttr(ab)
	}

	// AttrCond bindings (conditional attributes from HtmlNode.AttrIf())
	for _, acb := range bindings.AttrCondBindings {
		bindAttrCondBinding(acb)
	}

	// Each block bindings (list iteration)
	for _, eb := range bindings.EachBlocks {
		bindEachBlock(eb)
	}

	// NOTE: Route blocks are NOT applied here. They require the router's path store
	// which is created during OnMount (after applyBindings runs). Route blocks are
	// applied explicitly after OnMount in hydrateWASM.

}

// bindAttr sets up a dynamic attribute binding.
func bindAttr(ab HydrateAttrBinding) {
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
		store := GetStore(storeID)
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
			placeholder := "{" + itoa(i) + "}"
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
func bindAttrCondBinding(acb HydrateAttrCondBinding) {
	el := GetEl(acb.ElementID)
	if !ok(el) {
		return
	}

	if len(acb.Deps) == 0 {
		return
	}

	// Resolve the condition store (first dep is always the condition)
	condStore := GetStore(acb.Deps[0])
	if condStore == nil {
		return
	}

	// Resolve true/false value stores if they're dynamic
	var trueStore, falseStore any
	if acb.TrueStoreID != "" {
		trueStore = GetStore(acb.TrueStoreID)
	}
	if acb.FalseStoreID != "" {
		falseStore = GetStore(acb.FalseStoreID)
	}

	// Function to evaluate condition and update attribute
	updateAttr := func() {
		active := evaluateStore(condStore, acb.IsBool, acb.Op, acb.Operand)

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
func bindEachBlock(eb HydrateEachBlock) {
	if eb.ListID == "" {
		return
	}

	// Check if already setup
	if setupEachBlocks[eb.MarkerID] {
		return
	}
	setupEachBlocks[eb.MarkerID] = true

	// Resolve the list
	listAny := GetStore(eb.ListID)
	if listAny == nil {
		return
	}

	// Subscribe to list changes and re-render
	switch list := listAny.(type) {
	case *List[string]:
		bindListItems(list, eb.MarkerID, eb.BodyHTML, escapeHTML)
	case *List[int]:
		bindListItems(list, eb.MarkerID, eb.BodyHTML, itoa)
	}
}

// bindListItems sets up list rendering and subscribes to changes.
func bindListItems[T comparable](list *List[T], markerID string, bodyTemplate string, format func(T) string) {
	list.OnChange(func(items []T) {
		var html string
		for i, item := range items {
			body := replaceAll(bodyTemplate, "\x00I\x00", format(item))
			body = replaceAll(body, "\x00N\x00", itoa(i))
			html += body
		}
		replaceMarkerContent(markerID, html)
	})
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

// evaluateStore evaluates a condition against a store's current value.
func evaluateStore(store any, isBool bool, op, operand string) bool {
	if isBool {
		if s, ok := store.(*Store[bool]); ok {
			return s.Get()
		}
		return false
	}
	switch s := store.(type) {
	case *Store[int]:
		return compare(s.Get(), op, atoiSafe(operand))
	case *Store[string]:
		return compare(s.Get(), op, operand)
	case *Store[float64]:
		return compare(s.Get(), op, atofSafe(operand))
	case *Store[bool]:
		return compareBool(s.Get(), op, operand == "true")
	}
	return false
}

// evalCondition evaluates a branch condition using structured data.
func evalCondition(branch HydrateIfBranch) bool {
	if branch.StoreID == "" {
		return false
	}
	store := GetStore(branch.StoreID)
	if store == nil {
		return false
	}
	return evaluateStore(store, branch.IsBool, branch.Op, branch.Operand)
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
		return itoa(s.Get())
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
	result := itoa(intPart)
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

// HydrateBindings holds decoded bindings for WASM hydration.
type HydrateBindings struct {
	TextBindings     []HydrateTextBinding
	Events           []HydrateEvent
	IfBlocks         []HydrateIfBlock
	EachBlocks       []HydrateEachBlock
	InputBindings    []HydrateInputBinding
	AttrBindings     []HydrateAttrBinding
	AttrCondBindings []HydrateAttrCondBinding
	ComponentBlocks  []HydrateComponentBlock
}

type HydrateTextBinding struct {
	MarkerID string
	StoreID  string
	IsHTML   bool
}

type HydrateEvent struct {
	ElementID string
	Event     string
}

type HydrateIfBlock struct {
	MarkerID     string
	Branches     []HydrateIfBranch
	ElseHTML     string
	ElseBindings *HydrateBindings
	Deps         []string
}

type HydrateIfBranch struct {
	HTML     string
	Bindings *HydrateBindings
	StoreID  string
	Op       string
	Operand  string
	IsBool   bool
}

type HydrateInputBinding struct {
	StoreID  string
	BindType string
}

type HydrateAttrBinding struct {
	ElementID string
	AttrName  string
	Template  string
	StoreIDs  []string
}

// HydrateAttrCondBinding represents a conditional attribute binding for WASM.
type HydrateAttrCondBinding struct {
	ElementID    string
	AttrName     string
	TrueValue    string
	FalseValue   string
	TrueStoreID  string
	FalseStoreID string
	Op           string
	Operand      string
	IsBool       bool
	Deps         []string
}

type HydrateEachBlock struct {
	MarkerID string
	ListID   string
	BodyHTML string
}

// HydrateComponentBlock is the WASM-side representation of a ComponentBlock.
type HydrateComponentBlock struct {
	MarkerID string
	StoreID  string
	Branches []HydrateComponentBranch
}

// HydrateComponentBranch represents one component's pre-baked content.
type HydrateComponentBranch struct {
	Name     string
	HTML     string
	Bindings *HydrateBindings
}
