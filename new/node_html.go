//go:build !js || !wasm

package preveltekit

import (
	"fmt"
	"reflect"
	"strings"
)

// BuildContext holds state during HTML and wiring code generation.
type BuildContext struct {
	// Counters for generating unique IDs
	TextCounter  int
	IfCounter    int
	EachCounter  int
	EventCounter int
	BindCounter  int
	ClassCounter int
	AttrCounter  int
	CompCounter  int

	// Component prefix for nested components (e.g., "comp0_comp1")
	Prefix string

	// Parent context for nested components
	Parent *BuildContext

	// Collected bindings during tree walk
	Bindings *CollectedBindings

	// SlotContent holds HTML to be rendered in place of <slot/> elements
	SlotContent string

	// ChildrenContent maps child names to their pre-rendered HTML for SPA routing
	ChildrenContent map[string]string

	// ChildrenBindings maps child names to their bindings for SPA routing
	ChildrenBindings map[string]*CollectedBindings

	// NestedComponents maps component names to factory functions for p.Comp()
	NestedComponents map[string]func() Component

	// ParentStoreMap maps store pointers to their IDs in the parent component
	// Used to resolve dynamic props that share parent stores
	ParentStoreMap map[uintptr]string
}

// CollectedBindings stores all bindings found during tree walking.
type CollectedBindings struct {
	TextBindings   []TextBinding       `json:"TextBindings"`
	Events         []EventBinding_     `json:"Events"`
	IfBlocks       []IfBlock           `json:"IfBlocks"`
	EachBlocks     []EachBlock         `json:"EachBlocks"`
	InputBindings  []InputBinding_     `json:"InputBindings"`
	ClassBindings  []ClassBinding_     `json:"ClassBindings"`
	AttrBindings   []AttrBinding_      `json:"AttrBindings"`
	Components     []ComponentBinding_ `json:"Components"`
	ShowIfBindings []ShowIfBinding_    `json:"ShowIfBindings"`
}

// =============================================================================
// Binding types for code generation
// =============================================================================
//
// Bindings reference DOM elements in two ways:
// - MarkerID: References HTML comment markers (e.g., <!--basics_t0-->)
// - ElementID: References HTML element id attributes (e.g., id="basics_ev0")
// =============================================================================

// TextBinding binds a store value to text content at a comment marker.
// HTML: "Hello <!--basics_t0-->" where the text before the marker updates reactively.
type TextBinding struct {
	MarkerID string `json:"marker_id"` // Comment marker, e.g., "basics_t0"
	StoreID  string `json:"store_id"`  // Store path, e.g., "basics.Name"
	StoreRef any    `json:"-"`         // Actual store pointer (for address-based resolution)
	IsHTML   bool   `json:"is_html"`   // If true, render as HTML not text
}

// EventBinding_ binds an event handler to a DOM element by its id attribute.
// HTML: <button id="basics_ev0"> triggers the handler on click.
type EventBinding_ struct {
	ElementID  string   // Element id attribute, e.g., "basics_ev0"
	Event      string   // Event name, e.g., "click"
	HandlerID  string   // Handler path, e.g., "basics.Increment"
	HandlerRef any      `json:"-"` // Actual method reference (for resolution)
	Args       []any    `json:"-"` // Handler arguments (for resolution)
	ArgsStr    string   // Serialized args for JSON
	Modifiers  []string // Event modifiers, e.g., ["preventDefault"]
}

// IfBlock represents a conditional block with branches at a comment marker.
// HTML: <span>active content</span><!--basics_i0--> where content swaps reactively.
type IfBlock struct {
	MarkerID     string             // Comment marker, e.g., "basics_i0"
	Branches     []IfBlockBranch    // Condition branches (if/else-if)
	ElseHTML     string             // HTML for else branch
	ElseBindings *CollectedBindings `json:"ElseBindings,omitempty"` // Bindings for else branch
	Deps         []string           // Store dependencies for reactivity
}

// IfBlockBranch represents one branch (if or else-if) of an IfBlock.
type IfBlockBranch struct {
	CondExpr string             `json:"cond_expr"`          // Condition expression for display
	CondRef  any                `json:"-"`                  // Condition reference for store resolution
	HTML     string             `json:"html"`               // Pre-rendered HTML for this branch
	Bindings *CollectedBindings `json:"Bindings,omitempty"` // Nested bindings within this branch
	// Structured condition data for WASM evaluation
	StoreID string `json:"store_id,omitempty"` // Store path, e.g., "basics.Score"
	Op      string `json:"op,omitempty"`       // Comparison operator, e.g., ">="
	Operand string `json:"operand,omitempty"`  // Comparison value, e.g., "90"
	IsBool  bool   `json:"is_bool,omitempty"`  // True if simple boolean condition
}

// EachBlock represents a list iteration block at a comment marker.
// HTML: <span id="basics_each0_0">item</span><!--basics_e0--> where items update reactively.
type EachBlock struct {
	MarkerID string // Comment marker, e.g., "basics_e0"
	ListID   string // List store path, e.g., "basics.Items"
	ItemVar  string // Item variable name in template
	IndexVar string // Index variable name in template
	BodyHTML string // Template HTML for each item
	ElseHTML string // HTML when list is empty
}

// InputBinding_ binds an input element to a store for two-way data binding.
// HTML: <input id="basics_b0"> syncs value with store.
type InputBinding_ struct {
	ElementID string `json:"element_id"` // Element id attribute, e.g., "basics_b0"
	StoreID   string `json:"store_id"`   // Store path, e.g., "basics.Name"
	StoreRef  any    `json:"-"`          // Actual store pointer (for resolution)
	BindType  string `json:"bind_type"`  // Binding type: "value" or "checked"
}

// ClassBinding_ binds a CSS class to a condition on a DOM element.
// HTML: <div id="basics_cl0" class="active"> where "active" toggles reactively.
type ClassBinding_ struct {
	ElementID string   `json:"element_id"` // Element id attribute, e.g., "basics_cl0"
	ClassName string   `json:"class_name"` // CSS class to toggle
	CondExpr  string   `json:"cond_expr"`  // Condition expression for display
	StoreRef  any      `json:"-"`          // Store pointer for condition evaluation
	Op        string   `json:"op"`         // Comparison operator (for StoreCondition)
	Operand   string   `json:"operand"`    // Comparison operand (for StoreCondition)
	Deps      []string `json:"deps"`       // Store dependencies for reactivity
}

// AttrBinding_ binds a dynamic attribute value to stores.
// HTML: <div data-attrbind="basics_a0" data-value="{0}"> where {0} is replaced.
type AttrBinding_ struct {
	ElementID string   // Element id (via data-attrbind), e.g., "basics_a0"
	AttrName  string   // Attribute name, e.g., "data-value"
	Template  string   // Template with placeholders, e.g., "{0}"
	StoreIDs  []string // Store paths for placeholders
}

// ComponentBinding_ represents a nested component instance.
type ComponentBinding_ struct {
	ElementID string            // Component prefix, e.g., "basics_comp0"
	Name      string            // Component type name, e.g., "Button"
	Props     map[string]string // Static prop values
	Events    map[string]string // Event handler mappings
	SlotHTML  string            // Slot content HTML
}

// ShowIfBinding_ binds element visibility to a condition.
// HTML: <div id="page-basics" style="display:none"> toggles visibility reactively.
type ShowIfBinding_ struct {
	ElementID string   `json:"element_id"` // Element id attribute, e.g., "page-basics"
	StoreID   string   `json:"store_id"`   // Store path for condition
	StoreRef  any      `json:"-"`          // Store pointer for resolution
	Op        string   `json:"op"`         // Comparison operator
	Operand   string   `json:"operand"`    // Comparison operand
	IsBool    bool     `json:"is_bool"`    // True if simple boolean condition
	Deps      []string `json:"deps"`       // Store dependencies for reactivity
}

// NewBuildContext creates a new build context for HTML generation.
func NewBuildContext() *BuildContext {
	return &BuildContext{
		Bindings: &CollectedBindings{},
	}
}

// Child creates a child context for a nested component.
func (ctx *BuildContext) Child(compID string) *BuildContext {
	prefix := compID
	if ctx.Prefix != "" {
		prefix = ctx.Prefix + "_" + compID
	}
	return &BuildContext{
		Prefix:           prefix,
		Parent:           ctx,
		Bindings:         &CollectedBindings{},
		NestedComponents: ctx.NestedComponents,
		ParentStoreMap:   ctx.ParentStoreMap,
	}
}

// =============================================================================
// ID Generation
// =============================================================================
//
// There are two types of IDs used in the generated HTML:
//
// 1. ELEMENT IDs - Used in HTML id="..." attributes for DOM element lookup
//    - Generated by: NextEventID, NextBindID, NextClassID, NextAttrID
//    - Format: Full prefix + local ID (e.g., "components_ev0", "basics_b0")
//    - Used by: Events, InputBindings, ClassBindings, AttrBindings, ShowIfBindings
//    - These IDs appear in the actual HTML element's id attribute
//
// 2. MARKER IDs - Used in HTML comments <!--marker--> for text/block insertion points
//    - Generated by: NextTextMarker, NextIfMarker, NextEachMarker, NextCompMarker
//    - Format: Shortened prefix + local ID (e.g., "components_c3_t0", "basics_i0")
//    - Used by: TextBindings, IfBlocks, EachBlocks
//    - These IDs appear in HTML comments and are shortened to save bytes
//
// The shortening rules for markers:
//   comp0 -> c0, if0 -> i0, each0 -> e0 (only for generated marker parts, not component names)
//   "components_comp3_t0" -> "components_c3_t0" (component name "components" preserved)
// =============================================================================

// --- Element ID generators (for HTML id="..." attributes) ---

// NextEventID returns the next element ID for event bindings.
// Used in: <button id="basics_ev0">
func (ctx *BuildContext) NextEventID() string {
	id := fmt.Sprintf("ev%d", ctx.EventCounter)
	ctx.EventCounter++
	return id
}

// NextBindID returns the next element ID for input bindings.
// Used in: <input id="basics_b0">
func (ctx *BuildContext) NextBindID() string {
	id := fmt.Sprintf("b%d", ctx.BindCounter)
	ctx.BindCounter++
	return id
}

// NextClassID returns the next element ID for class bindings.
// Used in: <div id="basics_cl0">
func (ctx *BuildContext) NextClassID() string {
	id := fmt.Sprintf("cl%d", ctx.ClassCounter)
	ctx.ClassCounter++
	return id
}

// NextAttrID returns the next element ID for attribute bindings.
// Used in: <div data-attrbind="basics_a0">
func (ctx *BuildContext) NextAttrID() string {
	id := fmt.Sprintf("a%d", ctx.AttrCounter)
	ctx.AttrCounter++
	return id
}

// --- Marker ID generators (for HTML comments <!--marker-->) ---

// NextTextMarker returns the next marker ID for text bindings.
// Used in: <!--basics_t0--> (comment marker for text insertion point)
func (ctx *BuildContext) NextTextMarker() string {
	id := fmt.Sprintf("t%d", ctx.TextCounter)
	ctx.TextCounter++
	return id
}

// NextIfMarker returns the next marker ID for if-blocks.
// Used in: <!--basics_i0--> (comment marker for if-block boundary)
func (ctx *BuildContext) NextIfMarker() string {
	id := fmt.Sprintf("if%d", ctx.IfCounter)
	ctx.IfCounter++
	return id
}

// NextEachMarker returns the next marker ID for each-blocks.
// Used in: <!--basics_e0--> (comment marker for each-block boundary)
func (ctx *BuildContext) NextEachMarker() string {
	id := fmt.Sprintf("each%d", ctx.EachCounter)
	ctx.EachCounter++
	return id
}

// NextCompMarker returns the next marker ID for nested components.
// Used internally for component prefixing (e.g., "comp0" in "components_comp0_t0")
func (ctx *BuildContext) NextCompMarker() string {
	id := fmt.Sprintf("comp%d", ctx.CompCounter)
	ctx.CompCounter++
	return id
}

// --- ID formatting functions ---

// FullElementID returns the full element ID with prefix for use in HTML id="..." attributes.
// Example: FullElementID("ev0") with prefix "basics" returns "basics_ev0"
func (ctx *BuildContext) FullElementID(localID string) string {
	if ctx.Prefix == "" {
		return localID
	}
	return ctx.Prefix + "_" + localID
}

// FullMarkerID returns the shortened marker ID for use in HTML comments.
// Example: FullMarkerID("t0") with prefix "components_comp3" returns "components_c3_t0"
// The marker parts (comp, if, each) are shortened but component names are preserved.
func (ctx *BuildContext) FullMarkerID(localID string) string {
	fullID := ctx.Prefix
	if fullID == "" {
		return shortenMarkerPart(localID)
	}
	return shortenMarkerParts(fullID) + "_" + shortenMarkerPart(localID)
}

// shortenMarkerParts shortens all marker parts in a prefixed ID.
// Example: "components_comp3" -> "components_c3"
func shortenMarkerParts(id string) string {
	parts := strings.Split(id, "_")
	for i, part := range parts {
		parts[i] = shortenMarkerPart(part)
	}
	return strings.Join(parts, "_")
}

// shortenMarkerPart shortens a single marker part if it matches a known pattern.
// Only shortens generated marker IDs (comp0, if0, each0), not component names.
// Example: "comp3" -> "c3", "if0" -> "i0", "components" -> "components" (unchanged)
func shortenMarkerPart(part string) string {
	// comp0 -> c0 (but not "components" which doesn't end in digits)
	if len(part) > 4 && part[:4] == "comp" && isDigits(part[4:]) {
		return "c" + part[4:]
	}
	// each0 -> e0
	if len(part) > 4 && part[:4] == "each" && isDigits(part[4:]) {
		return "e" + part[4:]
	}
	// if0 -> i0
	if len(part) > 2 && part[:2] == "if" && isDigits(part[2:]) {
		return "i" + part[2:]
	}
	// ev0 -> v0 (for markers, though events typically use element IDs)
	if len(part) > 2 && part[:2] == "ev" && isDigits(part[2:]) {
		return "v" + part[2:]
	}
	// cl0 -> l0 (for markers, though classes typically use element IDs)
	if len(part) > 2 && part[:2] == "cl" && isDigits(part[2:]) {
		return "l" + part[2:]
	}
	return part
}

// isDigits returns true if s contains only ASCII digits.
func isDigits(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// =============================================================================
// ToHTML implementations
// =============================================================================

// ToHTML generates HTML for an element node.
func (e *Element) ToHTML(ctx *BuildContext) string {
	var sb strings.Builder

	// Check for special attributes that need element IDs
	var elementID string
	var existingID string
	var hasEvent, hasBindValue, hasBindChecked, hasClassIf, hasShowIf bool
	var classIfAttrs []*ClassIfAttr
	var showIfAttr *ShowIfAttr

	for _, attr := range e.Attrs {
		switch a := attr.(type) {
		case *EventAttr:
			hasEvent = true
		case *BindValueAttr:
			hasBindValue = true
		case *BindCheckedAttr:
			hasBindChecked = true
		case *ClassIfAttr:
			hasClassIf = true
			classIfAttrs = append(classIfAttrs, a)
		case *ShowIfAttr:
			hasShowIf = true
			showIfAttr = a
		case *StaticAttr:
			if a.Name == "id" {
				existingID = a.Value
			}
		}
	}

	// Use existing ID if present, otherwise assign new ID if needed
	if existingID != "" {
		elementID = existingID
	} else if hasEvent {
		elementID = ctx.NextEventID()
	} else if hasBindValue || hasBindChecked {
		elementID = ctx.NextBindID()
	} else if hasClassIf || hasShowIf {
		elementID = ctx.NextClassID()
	}

	// Build opening tag
	sb.WriteString("<")
	sb.WriteString(e.Tag)

	// Add element ID if assigned (and not already in attrs)
	if elementID != "" && existingID == "" {
		sb.WriteString(` id="`)
		sb.WriteString(ctx.FullElementID(elementID))
		sb.WriteString(`"`)
	}

	// Process attributes
	var classes []string
	for _, attr := range e.Attrs {
		switch a := attr.(type) {
		case *ClassAttr:
			classes = append(classes, a.Classes...)
		case *ClassIfAttr:
			// Evaluated at SSR time
			if a.Cond.Eval() {
				classes = append(classes, a.ClassName)
			}
			// Record for wiring - try to get StoreRef from condition
			var storeRef any
			var op, operand string
			if bc, ok := a.Cond.(*BoolCondition); ok {
				storeRef = bc.Store
			} else if sc, ok := a.Cond.(*StoreCondition); ok {
				storeRef = sc.Store
				op = sc.Op
				operand = fmt.Sprintf("%v", sc.Operand)
			}
			ctx.Bindings.ClassBindings = append(ctx.Bindings.ClassBindings, ClassBinding_{
				ElementID: ctx.FullElementID(elementID),
				ClassName: a.ClassName,
				CondExpr:  a.Cond.Expr(),
				StoreRef:  storeRef,
				Op:        op,
				Operand:   operand,
				Deps:      a.Cond.Deps(),
			})
		case *StaticAttr:
			sb.WriteString(" ")
			sb.WriteString(a.Name)
			sb.WriteString(`="`)
			sb.WriteString(escapeHTML(a.Value))
			sb.WriteString(`"`)
		case *EventAttr:
			// Record event binding (no HTML output, uses element ID)
			ctx.Bindings.Events = append(ctx.Bindings.Events, EventBinding_{
				ElementID:  ctx.FullElementID(elementID),
				Event:      a.Event,
				HandlerID:  a.HandlerID,
				HandlerRef: a.Handler,
				Args:       a.Args,
				Modifiers:  a.Modifiers,
			})
		case *BindValueAttr:
			// Record input binding (uses element ID)
			ctx.Bindings.InputBindings = append(ctx.Bindings.InputBindings, InputBinding_{
				ElementID: ctx.FullElementID(elementID),
				StoreID:   a.StoreID,
				StoreRef:  a.Store,
				BindType:  "value",
			})
		case *BindCheckedAttr:
			// Record checkbox binding (uses element ID)
			ctx.Bindings.InputBindings = append(ctx.Bindings.InputBindings, InputBinding_{
				ElementID: ctx.FullElementID(elementID),
				StoreID:   a.StoreID,
				StoreRef:  a.Store,
				BindType:  "checked",
			})
		case *DynAttrAttr:
			attrID := ctx.NextAttrID()
			sb.WriteString(` data-attrbind="`)
			sb.WriteString(ctx.FullElementID(attrID))
			sb.WriteString(`"`)
			// TODO: Evaluate template at SSR time
			sb.WriteString(" ")
			sb.WriteString(a.Name)
			sb.WriteString(`="`)
			sb.WriteString(escapeHTML(a.Template))
			sb.WriteString(`"`)
			// Record attribute binding (uses element ID via data attribute)
			ctx.Bindings.AttrBindings = append(ctx.Bindings.AttrBindings, AttrBinding_{
				ElementID: ctx.FullElementID(attrID),
				AttrName:  a.Name,
				Template:  a.Template,
				StoreIDs:  a.StoreIDs,
			})
		}
	}

	// Handle ShowIf - add inline style for SSR and record binding
	if hasShowIf && showIfAttr != nil {
		// Evaluate condition at SSR time
		if !showIfAttr.Cond.Eval() {
			sb.WriteString(` style="display:none"`)
		}
		// Record binding for WASM
		var storeRef any
		var op, operand string
		var isBool bool
		if bc, ok := showIfAttr.Cond.(*BoolCondition); ok {
			storeRef = bc.Store
			isBool = true
		} else if sc, ok := showIfAttr.Cond.(*StoreCondition); ok {
			storeRef = sc.Store
			op = sc.Op
			operand = fmt.Sprintf("%v", sc.Operand)
		}
		ctx.Bindings.ShowIfBindings = append(ctx.Bindings.ShowIfBindings, ShowIfBinding_{
			ElementID: elementID, // Use the explicit ID from the element
			StoreRef:  storeRef,
			Op:        op,
			Operand:   operand,
			IsBool:    isBool,
			Deps:      showIfAttr.Cond.Deps(),
		})
	}

	// Output classes
	if len(classes) > 0 {
		sb.WriteString(` class="`)
		sb.WriteString(strings.Join(classes, " "))
		sb.WriteString(`"`)
	}

	// Self-closing tags
	if isSelfClosing(e.Tag) {
		sb.WriteString(">")
		return sb.String()
	}

	sb.WriteString(">")

	// Children
	for _, child := range e.Children {
		sb.WriteString(nodeToHTML(child, ctx))
	}

	// Closing tag
	sb.WriteString("</")
	sb.WriteString(e.Tag)
	sb.WriteString(">")

	return sb.String()
}

// ToHTML generates HTML for a text node.
func (t *TextNode) ToHTML(ctx *BuildContext) string {
	return escapeHTML(t.Content)
}

// ToHTML generates HTML for a fragment.
func (f *Fragment) ToHTML(ctx *BuildContext) string {
	var sb strings.Builder
	for _, child := range f.Children {
		sb.WriteString(nodeToHTML(child, ctx))
	}
	return sb.String()
}

// ToHTML generates HTML for a bind node (text interpolation).
func (b *BindNode) ToHTML(ctx *BuildContext) string {
	localMarker := ctx.NextTextMarker()
	markerID := ctx.FullMarkerID(localMarker)

	// Get current value for SSR
	var value string
	switch s := b.StoreRef.(type) {
	case *Store[string]:
		value = s.Get()
	case *Store[int]:
		value = fmt.Sprintf("%d", s.Get())
	case *Store[bool]:
		value = fmt.Sprintf("%t", s.Get())
	case *Store[float64]:
		value = fmt.Sprintf("%g", s.Get())
	default:
		value = ""
	}

	// Record text binding (uses marker ID in HTML comment)
	ctx.Bindings.TextBindings = append(ctx.Bindings.TextBindings, TextBinding{
		MarkerID: markerID,
		StoreID:  b.StoreID,
		StoreRef: b.StoreRef,
		IsHTML:   b.IsHTML,
	})

	if b.IsHTML {
		return fmt.Sprintf("<span>%s</span><!--%s-->", value, markerID)
	}
	return fmt.Sprintf("%s<!--%s-->", escapeHTML(value), markerID)
}

// ToHTML generates HTML for an if node (conditional rendering).
func (i *IfNode) ToHTML(ctx *BuildContext) string {
	localMarker := ctx.NextIfMarker()
	markerID := ctx.FullMarkerID(localMarker)

	// Collect dependencies
	var deps []string
	for _, branch := range i.Branches {
		deps = append(deps, branch.Cond.Deps()...)
	}

	// Build if block info for wiring (uses marker ID in HTML comment)
	ifBlock := IfBlock{
		MarkerID: markerID,
		Deps:     deps,
	}

	// Evaluate branches for SSR - each branch gets its own context to capture bindings
	var activeHTML string
	activeFound := false
	for _, branch := range i.Branches {
		// Create a child context to capture this branch's bindings
		branchCtx := &BuildContext{
			TextCounter:      ctx.TextCounter,
			IfCounter:        ctx.IfCounter,
			EachCounter:      ctx.EachCounter,
			EventCounter:     ctx.EventCounter,
			BindCounter:      ctx.BindCounter,
			ClassCounter:     ctx.ClassCounter,
			AttrCounter:      ctx.AttrCounter,
			CompCounter:      ctx.CompCounter,
			Prefix:           ctx.Prefix,
			Bindings:         &CollectedBindings{},
			ChildrenContent:  ctx.ChildrenContent,
			ChildrenBindings: ctx.ChildrenBindings,
			NestedComponents: ctx.NestedComponents,
			ParentStoreMap:   ctx.ParentStoreMap,
		}
		branchHTML := childrenToHTML(branch.Children, branchCtx)

		// Update parent counters
		ctx.TextCounter = branchCtx.TextCounter
		ctx.IfCounter = branchCtx.IfCounter
		ctx.EachCounter = branchCtx.EachCounter
		ctx.EventCounter = branchCtx.EventCounter
		ctx.BindCounter = branchCtx.BindCounter
		ctx.ClassCounter = branchCtx.ClassCounter
		ctx.AttrCounter = branchCtx.AttrCounter
		ctx.CompCounter = branchCtx.CompCounter

		ifBlock.Branches = append(ifBlock.Branches, IfBlockBranch{
			CondExpr: branch.Cond.Expr(),
			CondRef:  branch.Cond, // Store for address-based resolution
			HTML:     branchHTML,
			Bindings: branchCtx.Bindings,
		})
		if !activeFound && branch.Cond.Eval() {
			activeHTML = branchHTML
			activeFound = true
		}
	}

	// Else branch
	if len(i.ElseNode) > 0 {
		elseCtx := &BuildContext{
			TextCounter:      ctx.TextCounter,
			IfCounter:        ctx.IfCounter,
			EachCounter:      ctx.EachCounter,
			EventCounter:     ctx.EventCounter,
			BindCounter:      ctx.BindCounter,
			ClassCounter:     ctx.ClassCounter,
			AttrCounter:      ctx.AttrCounter,
			CompCounter:      ctx.CompCounter,
			Prefix:           ctx.Prefix,
			Bindings:         &CollectedBindings{},
			ChildrenContent:  ctx.ChildrenContent,
			ChildrenBindings: ctx.ChildrenBindings,
			NestedComponents: ctx.NestedComponents,
			ParentStoreMap:   ctx.ParentStoreMap,
		}
		elseHTML := childrenToHTML(i.ElseNode, elseCtx)

		// Update parent counters
		ctx.TextCounter = elseCtx.TextCounter
		ctx.IfCounter = elseCtx.IfCounter
		ctx.EachCounter = elseCtx.EachCounter
		ctx.EventCounter = elseCtx.EventCounter
		ctx.BindCounter = elseCtx.BindCounter
		ctx.ClassCounter = elseCtx.ClassCounter
		ctx.AttrCounter = elseCtx.AttrCounter
		ctx.CompCounter = elseCtx.CompCounter

		ifBlock.ElseHTML = elseHTML
		ifBlock.ElseBindings = elseCtx.Bindings
		if !activeFound {
			activeHTML = elseHTML
		}
	}

	ctx.Bindings.IfBlocks = append(ctx.Bindings.IfBlocks, ifBlock)

	// Note: We do NOT merge activeBindings to top-level context here.
	// The nested bindings (text, events, inputs, nested if-blocks, etc.) are stored
	// in ifBlock.Branches[].Bindings and will be applied by bindIfBlock during
	// initial hydration. This prevents duplicate bindings.

	// Output the active branch wrapped in a span (uses marker ID in HTML comment)
	return fmt.Sprintf("<span>%s</span><!--%s-->", activeHTML, markerID)
}

// ToHTML generates HTML for an each node (list iteration).
func (e *EachNode) ToHTML(ctx *BuildContext) string {
	localMarker := ctx.NextEachMarker()
	markerID := ctx.FullMarkerID(localMarker)

	// Each item needs an element ID for DOM manipulation (not a marker)
	// Use full element ID format for the span wrapper
	itemElementPrefix := ctx.FullElementID(localMarker)

	// Get list items for SSR
	var itemsHTML strings.Builder

	switch list := e.ListRef.(type) {
	case *List[string]:
		items := list.Get()
		if len(items) == 0 && len(e.ElseNode) > 0 {
			itemsHTML.WriteString(childrenToHTML(e.ElseNode, ctx))
		} else {
			for i, item := range items {
				itemHTML := nodeToHTML(e.Body(item, i), ctx)
				itemsHTML.WriteString(fmt.Sprintf(`<span id="%s_%d">%s</span>`, itemElementPrefix, i, itemHTML))
			}
		}
	case *List[int]:
		items := list.Get()
		if len(items) == 0 && len(e.ElseNode) > 0 {
			itemsHTML.WriteString(childrenToHTML(e.ElseNode, ctx))
		} else {
			for i, item := range items {
				itemHTML := nodeToHTML(e.Body(item, i), ctx)
				itemsHTML.WriteString(fmt.Sprintf(`<span id="%s_%d">%s</span>`, itemElementPrefix, i, itemHTML))
			}
		}
	}

	// Record each block binding (uses marker ID in HTML comment)
	ctx.Bindings.EachBlocks = append(ctx.Bindings.EachBlocks, EachBlock{
		MarkerID: markerID,
		ListID:   e.ListID,
		ItemVar:  e.ItemVar,
		IndexVar: e.IndexVar,
		// Body HTML would need to be generated from the Body function
	})

	return fmt.Sprintf("%s<!--%s-->", itemsHTML.String(), markerID)
}

// ToHTML generates HTML for a component node (nested component).
func (c *ComponentNode) ToHTML(ctx *BuildContext) string {
	// Component marker is used as prefix for nested bindings
	compMarker := ctx.NextCompMarker()
	fullCompPrefix := ctx.FullElementID(compMarker)

	// Look up component factory
	factory, ok := ctx.NestedComponents[c.Name]
	if !ok {
		// Component not registered - output placeholder with warning comment
		return fmt.Sprintf("<!-- component %s not registered -->", c.Name)
	}

	// Create component instance
	comp := factory()

	// Set props on the component's stores using reflection
	setComponentProps(comp, c.Props)

	// Render slot content first (with current context)
	slotHTML := childrenToHTML(c.Children, ctx)

	// Create child context for the component with its own prefix
	childCtx := &BuildContext{
		Prefix:           fullCompPrefix,
		Parent:           ctx,
		Bindings:         &CollectedBindings{},
		SlotContent:      slotHTML,
		NestedComponents: ctx.NestedComponents,
		ParentStoreMap:   ctx.ParentStoreMap, // Pass down parent store map
	}

	// Render the component
	html := nodeToHTML(comp.Render(), childCtx)

	// Build store map for the nested component
	// This maps store pointers to their field names (e.g., "components_comp0.Label")
	storeMap := buildStoreMap(comp, fullCompPrefix)

	// For dynamic props (shared stores), prefer parent's store ID over child's
	// This ensures reactivity works through the parent component
	if ctx.ParentStoreMap != nil {
		for addr, parentName := range ctx.ParentStoreMap {
			// If this store address exists in child's map, it's a shared store (dynamic prop)
			// Replace the child's name with the parent's name for proper resolution
			if _, exists := storeMap[addr]; exists {
				storeMap[addr] = parentName
			}
		}
	}

	resolveBindings(childCtx.Bindings, storeMap, fullCompPrefix, comp)

	// Merge child bindings into parent
	mergeNestedBindings(ctx.Bindings, childCtx.Bindings)

	return html
}

// setComponentProps sets props on a component's stores using reflection.
// For dynamic props (where value is a *Store), it shares the store pointer
// so that reactivity works through the parent's store.
func setComponentProps(comp Component, props map[string]any) {
	v := reflect.ValueOf(comp)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}

	for propName, propValue := range props {
		// Find field with matching name (case-insensitive first letter)
		field := v.FieldByName(propName)
		if !field.IsValid() {
			// Try with uppercase first letter
			field = v.FieldByName(strings.Title(propName))
		}
		if !field.IsValid() {
			continue
		}

		// Check if prop value is a store (dynamic prop) - share the store pointer
		switch pv := propValue.(type) {
		case *Store[string]:
			if field.CanSet() {
				field.Set(reflect.ValueOf(pv))
			}
			continue
		case *Store[int]:
			if field.CanSet() {
				field.Set(reflect.ValueOf(pv))
			}
			continue
		case *Store[bool]:
			if field.CanSet() {
				field.Set(reflect.ValueOf(pv))
			}
			continue
		}

		// Static prop - set the value on the component's store
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			storeVal := field.Elem()
			if storeVal.Kind() == reflect.Struct {
				// Look for Set method
				setMethod := field.MethodByName("Set")
				if setMethod.IsValid() {
					// Handle different prop value types
					switch pv := propValue.(type) {
					case string:
						setMethod.Call([]reflect.Value{reflect.ValueOf(pv)})
					case int:
						setMethod.Call([]reflect.Value{reflect.ValueOf(pv)})
					case bool:
						setMethod.Call([]reflect.Value{reflect.ValueOf(pv)})
					}
				}
			}
		}
	}
}

// mergeNestedBindings merges child component bindings into parent.
func mergeNestedBindings(parent, child *CollectedBindings) {
	parent.TextBindings = append(parent.TextBindings, child.TextBindings...)
	parent.Events = append(parent.Events, child.Events...)
	parent.IfBlocks = append(parent.IfBlocks, child.IfBlocks...)
	parent.InputBindings = append(parent.InputBindings, child.InputBindings...)
	parent.ClassBindings = append(parent.ClassBindings, child.ClassBindings...)
	parent.ShowIfBindings = append(parent.ShowIfBindings, child.ShowIfBindings...)
	parent.Components = append(parent.Components, child.Components...)
}

// ToHTML generates HTML for a slot node.
func (s *SlotNode) ToHTML(ctx *BuildContext) string {
	// If slot content is provided in context, render it
	if ctx.SlotContent != "" {
		return ctx.SlotContent
	}
	// Otherwise render an empty placeholder
	return ""
}

// ToHTML generates HTML for a child node (named child component placeholder).
func (c *ChildNode) ToHTML(ctx *BuildContext) string {
	// Look up the child component content from context
	if ctx.ChildrenContent != nil {
		if content, ok := ctx.ChildrenContent[c.Name]; ok {
			// Also merge the child's bindings into the current context
			if ctx.ChildrenBindings != nil {
				if childBindings, ok := ctx.ChildrenBindings[c.Name]; ok && childBindings != nil {
					ctx.Bindings.TextBindings = append(ctx.Bindings.TextBindings, childBindings.TextBindings...)
					ctx.Bindings.Events = append(ctx.Bindings.Events, childBindings.Events...)
					ctx.Bindings.InputBindings = append(ctx.Bindings.InputBindings, childBindings.InputBindings...)
					ctx.Bindings.ClassBindings = append(ctx.Bindings.ClassBindings, childBindings.ClassBindings...)
					ctx.Bindings.IfBlocks = append(ctx.Bindings.IfBlocks, childBindings.IfBlocks...)
					ctx.Bindings.ShowIfBindings = append(ctx.Bindings.ShowIfBindings, childBindings.ShowIfBindings...)
				}
			}
			return content
		}
	}
	return ""
}

// nodeToHTML dispatches to the appropriate ToHTML method.
func nodeToHTML(n Node, ctx *BuildContext) string {
	switch node := n.(type) {
	case *Element:
		return node.ToHTML(ctx)
	case *TextNode:
		return node.ToHTML(ctx)
	case *Fragment:
		return node.ToHTML(ctx)
	case *BindNode:
		return node.ToHTML(ctx)
	case *IfNode:
		return node.ToHTML(ctx)
	case *EachNode:
		return node.ToHTML(ctx)
	case *ComponentNode:
		return node.ToHTML(ctx)
	case *SlotNode:
		return node.ToHTML(ctx)
	case *ChildNode:
		return node.ToHTML(ctx)
	default:
		return ""
	}
}

// childrenToHTML generates HTML for a slice of nodes.
func childrenToHTML(nodes []Node, ctx *BuildContext) string {
	var sb strings.Builder
	for _, n := range nodes {
		sb.WriteString(nodeToHTML(n, ctx))
	}
	return sb.String()
}

// isSelfClosing returns true for self-closing HTML tags.
func isSelfClosing(tag string) bool {
	switch tag {
	case "br", "hr", "img", "input", "meta", "link", "area", "base", "col", "embed", "param", "source", "track", "wbr":
		return true
	}
	return false
}

// escapeHTML escapes HTML special characters.
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}

// RenderHTML is the main entry point for generating HTML from a Node tree.
func RenderHTML(root Node) (string, *CollectedBindings) {
	ctx := NewBuildContext()
	html := nodeToHTML(root, ctx)
	return html, ctx.Bindings
}

// RenderHTMLWithSlot renders HTML with slot content injected.
func RenderHTMLWithSlot(root Node, slotContent string) (string, *CollectedBindings) {
	ctx := NewBuildContext()
	ctx.SlotContent = slotContent
	html := nodeToHTML(root, ctx)
	return html, ctx.Bindings
}

// RenderHTMLWithChildren renders HTML with named child content for SPA routing.
func RenderHTMLWithChildren(root Node, childrenContent map[string]string, childrenBindings map[string]*CollectedBindings) (string, *CollectedBindings) {
	ctx := NewBuildContext()
	ctx.ChildrenContent = childrenContent
	ctx.ChildrenBindings = childrenBindings
	html := nodeToHTML(root, ctx)
	return html, ctx.Bindings
}

// RenderHTMLWithPrefix renders HTML with a prefix for unique IDs.
func RenderHTMLWithPrefix(root Node, prefix string) (string, *CollectedBindings) {
	ctx := NewBuildContext()
	ctx.Prefix = prefix
	html := nodeToHTML(root, ctx)
	return html, ctx.Bindings
}

// RenderHTMLWithContext renders HTML with full context options.
func RenderHTMLWithContext(root Node, opts ...func(*BuildContext)) (string, *CollectedBindings) {
	ctx := NewBuildContext()
	for _, opt := range opts {
		opt(ctx)
	}
	html := nodeToHTML(root, ctx)
	return html, ctx.Bindings
}

// WithNestedComponentsCtx sets nested components on the build context.
func WithNestedComponentsCtx(nc map[string]func() Component) func(*BuildContext) {
	return func(ctx *BuildContext) {
		ctx.NestedComponents = nc
	}
}

// WithPrefixCtx sets the prefix on the build context.
func WithPrefixCtx(prefix string) func(*BuildContext) {
	return func(ctx *BuildContext) {
		ctx.Prefix = prefix
	}
}

// WithChildrenContentCtx sets children content on the build context.
func WithChildrenContentCtx(content map[string]string, bindings map[string]*CollectedBindings) func(*BuildContext) {
	return func(ctx *BuildContext) {
		ctx.ChildrenContent = content
		ctx.ChildrenBindings = bindings
	}
}

// WithParentStoreMapCtx sets the parent store map on the build context.
func WithParentStoreMapCtx(storeMap map[uintptr]string) func(*BuildContext) {
	return func(ctx *BuildContext) {
		ctx.ParentStoreMap = storeMap
	}
}
