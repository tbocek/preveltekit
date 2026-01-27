//go:build js && wasm

package preveltekit

import (
	"reflect"
	"syscall/js"
)

// Track which if-blocks have been set up to avoid duplicates
var setupIfBlocks = make(map[string]bool)

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

	// Find the existing SSR content
	currentEl := FindExistingIfContent(ifb.MarkerID)

	// Track current cleanup for bindings
	currentCleanup := &Cleanup{}

	// Evaluate which branch is active and update content
	updateIfBlock := func() {
		var activeHTML string
		var activeBindings *HydrateBindings
		found := false

		for _, branch := range ifb.Branches {
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

		// Re-find currentEl from DOM in case parent if-block replaced our container
		currentEl = FindExistingIfContent(ifb.MarkerID)
		// Replace content
		currentEl = ReplaceContent(ifb.MarkerID, currentEl, activeHTML)

		// Clean up old bindings
		currentCleanup.Release()
		currentCleanup = &Cleanup{}

		// Apply new bindings
		if activeBindings != nil {
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

	// Apply initial bindings for the SSR-active branch (don't replace HTML, just wire up bindings)
	// We do this separately from updateIfBlock to avoid replacing the already-correct SSR HTML
	var initialBindings *HydrateBindings
	for _, branch := range ifb.Branches {
		if evalCondition(branch, components) {
			initialBindings = branch.Bindings
			break
		}
	}
	if initialBindings == nil && ifb.ElseBindings != nil {
		initialBindings = ifb.ElseBindings
	}
	if initialBindings != nil {
		applyBindings(initialBindings, components, currentCleanup)
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
		// Silently skip unresolved stores - they're static props baked into HTML
	}

	// Input bindings
	for _, ib := range bindings.InputBindings {
		store := resolveStore(ib.StoreID, components)
		if store != nil {
			bindInputDynamic(cleanup, ib.ElementID, store, ib.BindType)
		}
		// Silently skip unresolved stores - they're static props baked into HTML
	}

	// Event bindings
	for _, ev := range bindings.Events {
		handler := resolveHandler(ev.HandlerID, ev.ArgsStr, components)
		if handler != nil {
			bindEventDynamic(cleanup, ev.ElementID, ev.Event, handler)
		}
		// Silently skip unresolved handlers - they may be from nested components
	}

	// Nested if-blocks
	for _, ifb := range bindings.IfBlocks {
		bindIfBlock(ifb, components)
	}

	// ShowIf bindings
	for _, sib := range bindings.ShowIfBindings {
		bindShowIf(sib, components)
	}
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
		val := s.Get()
		operand := atoiSafe(branch.Operand)
		return compareInt(val, branch.Op, operand)
	case *Store[string]:
		val := s.Get()
		return compareString(val, branch.Op, branch.Operand)
	case *Store[float64]:
		val := s.Get()
		operand := atofSafe(branch.Operand)
		return compareFloat(val, branch.Op, operand)
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
	InputBindings  []HydrateInputBinding  `json:"InputBindings"`
	ClassBindings  []HydrateClassBinding  `json:"ClassBindings"`
	ShowIfBindings []HydrateShowIfBinding `json:"ShowIfBindings"`
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
