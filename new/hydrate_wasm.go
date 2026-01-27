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
// In WASM mode, we don't actually call Render() - we just need the component
// for lifecycle methods and store access.
type Component interface {
	// Render is not called in WASM mode, but components must implement it
	// for the SSR phase. We use an empty interface here since Node is not
	// available in WASM builds.
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

// HasHandleEvent is implemented by components that handle events.
// The method name and args string are passed, and the component switches on it.
type HasHandleEvent interface {
	HandleEvent(method string, args string)
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

// WithNestedComponent registers a nested component type by name.
// In WASM mode, this is a no-op since nested components are pre-rendered during SSR.
func WithNestedComponent(name string, factory func() Component) func(*HydrateConfig) {
	return func(cfg *HydrateConfig) {
		// No-op in WASM - nested components are already rendered in HTML
	}
}

// hydrateWASM sets up DOM bindings from the embedded bindings JSON.
func hydrateWASM(app Component, cfg *HydrateConfig) {
	// Get bindings from global variable (set by CLI-generated code)
	bindingsJS := js.Global().Get("_preveltekit_bindings")
	if bindingsJS.IsUndefined() || bindingsJS.IsNull() {
		// No bindings - just run lifecycle and keep alive
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
		// Extract component name from path (e.g., "/basics" -> "basics")
		name := trimPrefix(path, "/")
		components[name] = child
	}

	// Call OnCreate
	if oc, ok := app.(HasOnCreate); ok {
		oc.OnCreate()
	}
	for _, child := range cfg.Children {
		if oc, ok := child.(HasOnCreate); ok {
			oc.OnCreate()
		}
	}

	// Inject styles
	if hs, ok := app.(HasStyle); ok {
		InjectStyle("app", hs.Style())
	}
	for path, child := range cfg.Children {
		if hs, ok := child.(HasStyle); ok {
			name := trimPrefix(path, "/")
			InjectStyle(name, hs.Style())
		}
	}

	cleanup := &Cleanup{}

	// Set up text bindings
	for _, tb := range bindings.TextBindings {
		store := resolveStore(tb.StoreID, components)
		if store != nil {
			bindTextDynamic(tb.MarkerID, store, tb.IsHTML)
		}
	}

	// Set up input bindings
	for _, ib := range bindings.InputBindings {
		store := resolveStore(ib.StoreID, components)
		if store != nil {
			bindInputDynamic(cleanup, ib.ElementID, store, ib.BindType)
		}
	}

	// Set up event bindings
	for _, ev := range bindings.Events {
		handler := resolveHandler(ev.HandlerID, ev.ArgsStr, components)
		if handler != nil {
			bindEventDynamic(cleanup, ev.ElementID, ev.Event, handler)
		}
	}

	// Set up if-block bindings
	for _, ifb := range bindings.IfBlocks {
		bindIfBlock(ifb, components)
	}

	// Set up show-if bindings
	for _, sib := range bindings.ShowIfBindings {
		bindShowIf(sib, components)
	}

	// Call OnMount
	if om, ok := app.(HasOnMount); ok {
		om.OnMount()
	}
	for _, child := range cfg.Children {
		if om, ok := child.(HasOnMount); ok {
			om.OnMount()
		}
	}

	// Keep WASM running
	select {}
}

func runLifecycle(app Component, cfg *HydrateConfig) {
	if oc, ok := app.(HasOnCreate); ok {
		oc.OnCreate()
	}
	for _, child := range cfg.Children {
		if oc, ok := child.(HasOnCreate); ok {
			oc.OnCreate()
		}
	}
	if hs, ok := app.(HasStyle); ok {
		InjectStyle("app", hs.Style())
	}
	if om, ok := app.(HasOnMount); ok {
		om.OnMount()
	}
	for _, child := range cfg.Children {
		if om, ok := child.(HasOnMount); ok {
			om.OnMount()
		}
	}
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

// resolveHandler resolves a handler path like "component.Increment" or "basics.AddItem".
func resolveHandler(handlerID string, argsStr string, components map[string]Component) func() {
	compName, methodName, ok := splitFirst(handlerID, ".")
	if !ok {
		return nil
	}

	comp, ok := components[compName]
	if !ok {
		return nil
	}

	// Use HasHandleEvent interface instead of reflection
	handler, ok := comp.(HasHandleEvent)
	if !ok {
		return nil
	}

	return func() {
		handler.HandleEvent(methodName, argsStr)
	}
}

// bindTextDynamic sets up a text binding using reflection.
func bindTextDynamic(markerID string, store any, isHTML bool) {
	// Use type switch to handle known store types
	switch s := store.(type) {
	case *Store[string]:
		if isHTML {
			BindHTML(markerID, s)
		} else {
			BindText(markerID, s)
		}
	case *Store[int]:
		BindText(markerID, s)
	case *Store[bool]:
		BindText(markerID, s)
	case *Store[float64]:
		BindText(markerID, s)
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
		println("bindShowIf: element not found:", sib.ElementID)
		return
	}

	store := resolveStore(sib.StoreID, components)
	if store == nil {
		println("bindShowIf: store not found:", sib.StoreID)
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
				val := s.Get()
				operand := atoiSafe(sib.Operand)
				visible = compareInt(val, sib.Op, operand)
			case *Store[string]:
				val := s.Get()
				visible = compareString(val, sib.Op, sib.Operand)
			case *Store[float64]:
				val := s.Get()
				operand := atofSafe(sib.Operand)
				visible = compareFloat(val, sib.Op, operand)
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
	println("bindIfBlock:", ifb.MarkerID)

	// Find the existing SSR content
	currentEl := FindExistingIfContent(ifb.MarkerID)

	// Track current cleanup for bindings
	currentCleanup := &Cleanup{}

	// Evaluate which branch is active and update content
	updateIfBlock := func() {
		println("updateIfBlock:", ifb.MarkerID, "branches:", len(ifb.Branches))
		var activeHTML string
		var activeBindings *HydrateBindings
		found := false

		for i := 0; i < len(ifb.Branches); i++ {
			println("updateIfBlock: eval branch", i)
			branch := ifb.Branches[i]
			result := evalCondition(branch, components)
			if result {
				activeHTML = branch.HTML
				activeBindings = branch.Bindings
				found = true
				println("updateIfBlock: matched branch", i)
				break
			}
		}

		if !found {
			activeHTML = ifb.ElseHTML
			activeBindings = ifb.ElseBindings
			println("updateIfBlock: using else")
		}

		println("updateIfBlock: replacing content")
		currentEl = FindExistingIfContent(ifb.MarkerID)
		currentEl = ReplaceContent(ifb.MarkerID, currentEl, activeHTML)

		println("updateIfBlock: cleanup")
		currentCleanup.Release()
		currentCleanup = &Cleanup{}

		println("updateIfBlock: applying")
		if activeBindings != nil {
			clearBoundMarkers(activeBindings)
			applyBindings(activeBindings, components, currentCleanup)
		}
		println("updateIfBlock: done", ifb.MarkerID)
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
			println("bindIfBlock: subscribing to", dep, "for", ifb.MarkerID)
			subscribeToStore(store, updateIfBlock)
		} else {
			println("bindIfBlock: WARNING - could not resolve dep", dep, "for", ifb.MarkerID)
		}
	}

	// Call updateIfBlock to sync DOM with current state
	// This handles nested if-blocks where state may have changed after SSR
	// (e.g., Lists page if-block when ItemCount was set in OnMount before navigation)
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
		handler := resolveHandler(ev.HandlerID, ev.ArgsStr, components)
		if handler != nil {
			bindEventDynamic(cleanup, ev.ElementID, ev.Event, handler)
		}
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
			var storeVal string
			switch s := store.(type) {
			case *Store[string]:
				storeVal = s.Get()
			case *Store[int]:
				storeVal = intToStr(s.Get())
			case *Store[bool]:
				if s.Get() {
					storeVal = "true"
				} else {
					storeVal = "false"
				}
			}
			value = replaceAll(value, placeholder, storeVal)
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
	println("bindEachBlock:", eb.MarkerID, eb.ListID)

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

	// Get the parent element (should be the <ul> or container)
	parent := marker.Get("previousSibling")
	if parent.IsNull() || parent.Get("nodeType").Int() != 1 {
		// Try parent node (the marker is inside the <ul>)
		parent = marker.Get("parentNode")
	}

	// Subscribe to list changes and re-render
	switch list := listAny.(type) {
	case *List[string]:
		renderItems := func(items []string) {
			var html string
			for i, item := range items {
				itemID := eb.MarkerID[:len(eb.MarkerID)-1]
				if len(itemID) > 0 && itemID[len(itemID)-1] == 'e' {
					itemID = itemID[:len(itemID)-1]
				}
				html += `<span id="` + itemID + `_` + intToStr(i) + `"><li><span class="index">` + intToStr(i) + `</span> ` + escapeHTMLWasm(item) + `</li></span>`
			}
			if !parent.IsNull() && parent.Get("nodeType").Int() == 1 {
				parent.Set("innerHTML", html)
			}
		}

		items := list.Get()
		if len(items) > 0 {
			renderItems(items)
		}

		list.OnChange(func(items []string) {
			renderItems(items)
		})

	case *List[int]:
		renderItems := func(items []int) {
			var html string
			for i, item := range items {
				itemID := eb.MarkerID[:len(eb.MarkerID)-1]
				if len(itemID) > 0 && itemID[len(itemID)-1] == 'e' {
					itemID = itemID[:len(itemID)-1]
				}
				html += `<span id="` + itemID + `_` + intToStr(i) + `"><li><span class="index">` + intToStr(i) + `</span> ` + intToStr(item) + `</li></span>`
			}
			if !parent.IsNull() && parent.Get("nodeType").Int() == 1 {
				parent.Set("innerHTML", html)
			}
		}

		items := list.Get()
		if len(items) > 0 {
			renderItems(items)
		}

		list.OnChange(func(items []int) {
			renderItems(items)
		})
	}

	_ = comp
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
		return compareInt(s.Get(), branch.Op, atoiSafe(branch.Operand))
	case *Store[string]:
		return compareString(s.Get(), branch.Op, branch.Operand)
	case *Store[float64]:
		return compareFloat(s.Get(), branch.Op, atofSafe(branch.Operand))
	}

	return false
}

func compareInt(val int, op string, operand int) bool {
	switch op {
	case ">=":
		return val >= operand
	case ">":
		return val > operand
	case "<=":
		return val <= operand
	case "<":
		return val < operand
	case "==":
		return val == operand
	case "!=":
		return val != operand
	}
	return false
}

func compareString(val string, op string, operand string) bool {
	switch op {
	case "==":
		return val == operand
	case "!=":
		return val != operand
	}
	return false
}

func compareFloat(val float64, op string, operand float64) bool {
	switch op {
	case ">=":
		return val >= operand
	case ">":
		return val > operand
	case "<=":
		return val <= operand
	case "<":
		return val < operand
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
}

type HydrateTextBinding struct {
	MarkerID string `json:"marker_id"`
	StoreID  string `json:"store_id"`
	IsHTML   bool   `json:"is_html"`
}

type HydrateEvent struct {
	ElementID string   `json:"ElementID"`
	Event     string   `json:"Event"`
	HandlerID string   `json:"HandlerID"`
	ArgsStr   string   `json:"ArgsStr"`
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
