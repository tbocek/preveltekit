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

	// ParentStoreMap maps store pointers to their IDs in the parent component
	// Used to resolve dynamic props that share parent stores
	ParentStoreMap map[uintptr]string

	// CollectedStyles holds CSS from nested components (deduplicated by component name)
	CollectedStyles map[string]string
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
	RouterBindings []RouterBinding_    `json:"RouterBindings"`
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
	ElementID string   // Element id attribute, e.g., "basics_ev0"
	Event     string   // Event name, e.g., "click"
	Modifiers []string // Event modifiers, e.g., ["preventDefault"]
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
	MarkerID string `json:"MarkerID"`           // Comment marker, e.g., "basics_e0"
	ListID   string `json:"ListID"`             // List store path, e.g., "basics.Items"
	ListRef  any    `json:"-"`                  // Actual list pointer (for resolution)
	ItemVar  string `json:"ItemVar"`            // Item variable name in template
	IndexVar string `json:"IndexVar"`           // Index variable name in template
	BodyHTML string `json:"BodyHTML,omitempty"` // Template HTML for each item
	ElseHTML string `json:"ElseHTML,omitempty"` // HTML when list is empty
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
	ElementID string   `json:"element_id"` // Element id (via data-attrbind), e.g., "basics_a0"
	AttrName  string   `json:"attr_name"`  // Attribute name, e.g., "data-value"
	Template  string   `json:"template"`   // Template with placeholders, e.g., "{0}"
	StoreIDs  []string `json:"store_ids"`  // Store paths for placeholders
	StoreRefs []any    `json:"-"`          // Actual store pointers (for resolution)
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

// RouterBinding_ binds a router to a store for page switching.
type RouterBinding_ struct {
	StoreID  string `json:"store_id"` // Store path, e.g., "component.CurrentPage"
	StoreRef any    `json:"-"`        // Store pointer for resolution
}

// NewBuildContext creates a new build context for HTML generation.
func NewBuildContext() *BuildContext {
	return &BuildContext{
		Bindings:        &CollectedBindings{},
		CollectedStyles: make(map[string]string),
	}
}

// Child creates a child context for a nested component.
func (ctx *BuildContext) Child(compID string) *BuildContext {
	prefix := compID
	if ctx.Prefix != "" {
		prefix = ctx.Prefix + "_" + compID
	}
	return &BuildContext{
		Prefix:         prefix,
		Parent:         ctx,
		Bindings:       &CollectedBindings{},
		ParentStoreMap: ctx.ParentStoreMap,
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
	var hasEvent, hasClassIf, hasShowIf bool
	var classIfAttrs []*ClassIfAttr
	var showIfAttr *ShowIfAttr

	for _, attr := range e.Attrs {
		switch a := attr.(type) {
		case *EventAttr:
			hasEvent = true
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
				ElementID: ctx.FullElementID(elementID),
				Event:     a.Event,
				Modifiers: a.Modifiers,
			})

		case *DynAttrAttr:
			attrID := ctx.NextAttrID()
			sb.WriteString(` data-attrbind="`)
			sb.WriteString(ctx.FullElementID(attrID))
			sb.WriteString(`"`)
			// Evaluate template at SSR time
			attrValue := a.Template
			for i, store := range a.Stores {
				placeholder := "{" + fmt.Sprintf("%d", i) + "}"
				var storeVal string
				switch s := store.(type) {
				case *Store[string]:
					storeVal = s.Get()
				case *Store[int]:
					storeVal = fmt.Sprintf("%d", s.Get())
				case *Store[bool]:
					storeVal = fmt.Sprintf("%t", s.Get())
				}
				attrValue = strings.ReplaceAll(attrValue, placeholder, storeVal)
			}
			sb.WriteString(" ")
			sb.WriteString(a.Name)
			sb.WriteString(`="`)
			sb.WriteString(escapeHTML(attrValue))
			sb.WriteString(`"`)
			// Record attribute binding (uses element ID via data attribute)
			ctx.Bindings.AttrBindings = append(ctx.Bindings.AttrBindings, AttrBinding_{
				ElementID: ctx.FullElementID(attrID),
				AttrName:  a.Name,
				Template:  a.Template,
				StoreRefs: a.Stores,
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

// ToHTML generates HTML for a raw HTML node with embedded nodes.
func (h *HtmlNode) ToHTML(ctx *BuildContext) string {
	var sb strings.Builder

	// Process parts, combining consecutive EventAttr/ClassIfAttr to share one ID
	for i := 0; i < len(h.Parts); i++ {
		part := h.Parts[i]
		switch v := part.(type) {
		case string:
			// Raw HTML string - pass through as-is (no escaping)
			sb.WriteString(v)
		case Node:
			// Embedded node - render it
			sb.WriteString(nodeToHTML(v, ctx))
		case NodeAttr:
			// Check if we have consecutive NodeAttrs that need to share an ID
			// Collect all consecutive EventAttr and ClassIfAttr (with whitespace between)
			attrs := []NodeAttr{v}
			j := i + 1
			for j < len(h.Parts) {
				// Skip whitespace-only strings
				if s, ok := h.Parts[j].(string); ok {
					trimmed := strings.TrimSpace(s)
					if trimmed == "" {
						j++
						continue
					}
					break
				}
				// Collect EventAttr or ClassIfAttr
				if attr, ok := h.Parts[j].(NodeAttr); ok {
					if _, isEvent := attr.(*EventAttr); isEvent {
						attrs = append(attrs, attr)
						j++
						continue
					}
					if _, isClassIf := attr.(*ClassIfAttr); isClassIf {
						attrs = append(attrs, attr)
						j++
						continue
					}
				}
				break
			}

			// If we have multiple attrs, render them with shared ID
			if len(attrs) > 1 {
				sb.WriteString(attrsToHTMLStringShared(attrs, ctx))
				i = j - 1
			} else {
				// Single attr - render normally
				sb.WriteString(attrToHTMLString(v, ctx))
			}
		default:
			// Convert other values to string and escape
			sb.WriteString(escapeHTML(fmt.Sprintf("%v", v)))
		}
	}
	return sb.String()
}

// attrsToHTMLStringShared renders multiple NodeAttrs sharing a single element ID.
// Used when OnClick and ClassIf are on the same element.
// Note: ClassIf classes are output via data-class attribute to avoid conflicts
// with existing class attributes in the HTML. WASM will apply them on hydration.
func attrsToHTMLStringShared(attrs []NodeAttr, ctx *BuildContext) string {
	// Determine what ID type we need
	hasEvent := false
	hasClassIf := false
	for _, attr := range attrs {
		if _, ok := attr.(*EventAttr); ok {
			hasEvent = true
		}
		if _, ok := attr.(*ClassIfAttr); ok {
			hasClassIf = true
		}
	}

	// Generate a single shared ID
	var localID string
	if hasEvent {
		localID = ctx.NextEventID()
	} else if hasClassIf {
		localID = ctx.NextClassID()
	}
	fullID := ctx.FullElementID(localID)

	var result strings.Builder
	result.WriteString(fmt.Sprintf(`id="%s"`, fullID))

	var activeClasses []string
	var eventName string

	for _, attr := range attrs {
		switch a := attr.(type) {
		case *EventAttr:
			eventName = a.Event
			ctx.Bindings.Events = append(ctx.Bindings.Events, EventBinding_{
				ElementID: fullID,
				Event:     a.Event,
				Modifiers: a.Modifiers,
			})
		case *ClassIfAttr:
			var storeRef any
			var op, operand string
			if sc, ok := a.Cond.(*StoreCondition); ok {
				storeRef = sc.Store
				op = sc.Op
				operand = fmt.Sprintf("%v", sc.Operand)
			} else if bc, ok := a.Cond.(*BoolCondition); ok {
				storeRef = bc.Store
			}
			ctx.Bindings.ClassBindings = append(ctx.Bindings.ClassBindings, ClassBinding_{
				ElementID: fullID,
				ClassName: a.ClassName,
				CondExpr:  a.Cond.Expr(),
				StoreRef:  storeRef,
				Op:        op,
				Operand:   operand,
				Deps:      a.Cond.Deps(),
			})
			if a.Cond.Eval() {
				activeClasses = append(activeClasses, a.ClassName)
			}
		}
	}

	if eventName != "" {
		result.WriteString(fmt.Sprintf(` data-event="%s"`, eventName))
	}
	// Output active classes via data-class so they don't conflict with existing class attr
	// WASM will read this and apply to classList on hydration
	if len(activeClasses) > 0 {
		result.WriteString(fmt.Sprintf(` data-class="%s"`, strings.Join(activeClasses, " ")))
	}

	return result.String()
}

// attrToHTMLString renders a NodeAttr as an HTML attribute string.
// Used when attributes are embedded directly in Html() nodes.
func attrToHTMLString(attr NodeAttr, ctx *BuildContext) string {
	switch a := attr.(type) {
	case *EventAttr:
		// Generate id and data-event attributes for event binding
		localID := ctx.NextEventID()
		fullID := ctx.FullElementID(localID)
		ctx.Bindings.Events = append(ctx.Bindings.Events, EventBinding_{
			ElementID: fullID,
			Event:     a.Event,
			Modifiers: a.Modifiers,
		})
		return fmt.Sprintf(`id="%s" data-event="%s"`, fullID, a.Event)
	case *ClassAttr:
		return fmt.Sprintf(`class="%s"`, strings.Join(a.Classes, " "))
	case *ClassIfAttr:
		// Generate id for class binding
		localID := ctx.NextClassID()
		fullID := ctx.FullElementID(localID)
		// Extract store reference from condition
		var storeRef any
		var op, operand string
		if sc, ok := a.Cond.(*StoreCondition); ok {
			storeRef = sc.Store
			op = sc.Op
			operand = fmt.Sprintf("%v", sc.Operand)
		} else if bc, ok := a.Cond.(*BoolCondition); ok {
			storeRef = bc.Store
		}
		ctx.Bindings.ClassBindings = append(ctx.Bindings.ClassBindings, ClassBinding_{
			ElementID: fullID,
			ClassName: a.ClassName,
			StoreRef:  storeRef,
			Op:        op,
			Operand:   operand,
			Deps:      a.Cond.Deps(),
		})
		// Use data-class to avoid conflicts with existing class attribute
		// WASM will apply this to classList on hydration
		classAttr := ""
		if a.Cond.Eval() {
			classAttr = fmt.Sprintf(` data-class="%s"`, a.ClassName)
		}
		return fmt.Sprintf(`id="%s"%s`, fullID, classAttr)
	case *StaticAttr:
		return fmt.Sprintf(`%s="%s"`, a.Name, escapeAttr(a.Value))

	case *DynAttrAttr:
		localID := ctx.NextAttrID()
		fullID := ctx.FullElementID(localID)
		// Evaluate template at SSR time
		attrValue := a.Template
		for i, store := range a.Stores {
			placeholder := "{" + fmt.Sprintf("%d", i) + "}"
			var storeVal string
			switch s := store.(type) {
			case *Store[string]:
				storeVal = s.Get()
			case *Store[int]:
				storeVal = fmt.Sprintf("%d", s.Get())
			}
			attrValue = strings.ReplaceAll(attrValue, placeholder, storeVal)
		}
		// Record attribute binding
		ctx.Bindings.AttrBindings = append(ctx.Bindings.AttrBindings, AttrBinding_{
			ElementID: fullID,
			AttrName:  a.Name,
			Template:  a.Template,
			StoreRefs: a.Stores,
		})
		return fmt.Sprintf(`data-attrbind="%s" %s="%s"`, fullID, a.Name, escapeAttr(attrValue))
	default:
		return ""
	}
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

// ToHTML generates HTML for a BindValue node (two-way input binding).
// Parses the HTML string and injects id and value attributes.
func (b *BindValueNode) ToHTML(ctx *BuildContext) string {
	localID := ctx.NextBindID()
	fullID := ctx.FullElementID(localID)

	// Get current value for SSR
	var value string
	switch s := b.Store.(type) {
	case *Store[string]:
		value = s.Get()
	case *Store[int]:
		value = fmt.Sprintf("%d", s.Get())
	}

	// Record input binding
	ctx.Bindings.InputBindings = append(ctx.Bindings.InputBindings, InputBinding_{
		ElementID: fullID,
		StoreRef:  b.Store,
		BindType:  "value",
	})

	// Check if it's a textarea (value goes as content, not attribute)
	if strings.HasPrefix(strings.TrimSpace(strings.ToLower(b.HTML)), "<textarea") {
		return injectTextareaValue(b.HTML, fullID, value)
	}

	// Inject id and value into the HTML element
	return injectAttrs(b.HTML, fmt.Sprintf(`id="%s" value="%s"`, fullID, escapeAttr(value)))
}

// injectTextareaValue handles textarea elements where value is content, not an attribute.
// Example: <textarea placeholder="..."></textarea> -> <textarea id="x" placeholder="...">value</textarea>
func injectTextareaValue(html, id, value string) string {
	// Find the closing > of the opening tag
	closeIdx := strings.Index(html, ">")
	if closeIdx == -1 {
		return html
	}

	// Insert id attribute before the >
	openTag := html[:closeIdx]
	rest := html[closeIdx+1:]

	// Find the closing </textarea>
	closeTagIdx := strings.Index(strings.ToLower(rest), "</textarea>")
	if closeTagIdx == -1 {
		// Self-closing or malformed, just inject id
		return openTag + fmt.Sprintf(` id="%s"`, id) + ">" + rest
	}

	// Replace content between tags with escaped value
	return openTag + fmt.Sprintf(` id="%s"`, id) + ">" + escapeHTML(value) + rest[closeTagIdx:]
}

// ToHTML generates HTML for a BindChecked node (checkbox binding).
// Parses the HTML string and injects id and checked attributes.
func (b *BindCheckedNode) ToHTML(ctx *BuildContext) string {
	localID := ctx.NextBindID()
	fullID := ctx.FullElementID(localID)

	// Record input binding
	ctx.Bindings.InputBindings = append(ctx.Bindings.InputBindings, InputBinding_{
		ElementID: fullID,
		StoreRef:  b.Store,
		BindType:  "checked",
	})

	// Get checked state
	checked := ""
	if s, ok := b.Store.(*Store[bool]); ok && s.Get() {
		checked = " checked"
	}

	// Inject id and checked into the HTML element
	return injectAttrs(b.HTML, fmt.Sprintf(`id="%s"%s`, fullID, checked))
}

// ToHTML generates HTML for a ClassIfNode (conditional class binding).
// Injects id and merges conditional classes with existing class attribute.
func (c *ClassIfNode) ToHTML(ctx *BuildContext) string {
	// Use event ID if we have OnClick, otherwise class ID
	var localID string
	if c.OnClick != nil {
		localID = ctx.NextEventID()
	} else {
		localID = ctx.NextClassID()
	}
	fullID := ctx.FullElementID(localID)

	// Collect active classes and record bindings
	var activeClasses []string
	for _, cond := range c.Conditions {
		var storeRef any
		var op, operand string
		if sc, ok := cond.Cond.(*StoreCondition); ok {
			storeRef = sc.Store
			op = sc.Op
			operand = fmt.Sprintf("%v", sc.Operand)
		} else if bc, ok := cond.Cond.(*BoolCondition); ok {
			storeRef = bc.Store
		}

		ctx.Bindings.ClassBindings = append(ctx.Bindings.ClassBindings, ClassBinding_{
			ElementID: fullID,
			ClassName: cond.ClassName,
			// CondExpr and Deps will be resolved by resolveBindings based on StoreRef
			StoreRef: storeRef,
			Op:       op,
			Operand:  operand,
		})

		if cond.Cond.Eval() {
			activeClasses = append(activeClasses, cond.ClassName)
		}
	}

	// Record event binding if OnClick is set
	var eventAttr string
	if c.OnClick != nil {
		ctx.Bindings.Events = append(ctx.Bindings.Events, EventBinding_{
			ElementID: fullID,
			Event:     "click",
		})
		eventAttr = ` data-event="click"`
	}

	return injectIDAndMergeClasses(c.HTML, fullID, activeClasses, eventAttr)
}

// injectIDAndMergeClasses injects id and merges classes into an HTML opening tag.
// If element has class="...", appends new classes. Otherwise adds new class attr.
// extraAttrs are appended as-is (e.g., ` data-event="click"`).
func injectIDAndMergeClasses(html, id string, classes []string, extraAttrs string) string {
	// Find existing class attribute
	classIdx := strings.Index(html, `class="`)
	if classIdx != -1 {
		// Find the closing quote
		classStart := classIdx + 7 // len(`class="`)
		classEnd := strings.Index(html[classStart:], `"`)
		if classEnd != -1 {
			classEnd += classStart
			existingClasses := html[classStart:classEnd]
			// Merge classes
			newClasses := existingClasses
			if len(classes) > 0 {
				newClasses = existingClasses + " " + strings.Join(classes, " ")
			}
			// Rebuild: inject id before class, update class value, add extra attrs
			return html[:classIdx] + fmt.Sprintf(`id="%s" class="%s"`, id, newClasses) + extraAttrs + html[classEnd+1:]
		}
	}

	// No existing class - inject id and class if any
	attrs := fmt.Sprintf(`id="%s"`, id)
	if len(classes) > 0 {
		attrs += fmt.Sprintf(` class="%s"`, strings.Join(classes, " "))
	}
	attrs += extraAttrs
	return injectAttrs(html, attrs)
}

// injectAttrs injects attributes into an HTML element string.
// Finds the first > and inserts the attrs just before it.
// Example: injectAttrs(`<input type="text">`, `id="foo"`) -> `<input type="text" id="foo">`
func injectAttrs(html, attrs string) string {
	// Find the closing > of the opening tag
	for i := 0; i < len(html); i++ {
		if html[i] == '>' {
			// Check if it's a self-closing tag />
			if i > 0 && html[i-1] == '/' {
				return html[:i-1] + " " + attrs + " />"
			}
			return html[:i] + " " + attrs + html[i:]
		}
	}
	// No > found, just append
	return html + " " + attrs
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
		ListRef:  e.ListRef,
		ItemVar:  e.ItemVar,
		IndexVar: e.IndexVar,
	})

	return fmt.Sprintf("%s<!--%s-->", itemsHTML.String(), markerID)
}

// ToHTML generates HTML for a component node (nested component).
func (c *ComponentNode) ToHTML(ctx *BuildContext) string {
	// Component marker is used as prefix for nested bindings
	compMarker := ctx.NextCompMarker()
	fullCompPrefix := ctx.FullElementID(compMarker)

	// Use the component instance directly
	comp, ok := c.Instance.(Component)
	if !ok {
		return fmt.Sprintf("<!-- component %s: invalid instance -->", c.Name)
	}

	// Collect style from nested component (deduplicated by component name)
	if hs, ok := c.Instance.(HasStyle); ok {
		if ctx.CollectedStyles != nil {
			if _, exists := ctx.CollectedStyles[c.Name]; !exists {
				ctx.CollectedStyles[c.Name] = hs.Style()
			}
		}
	}

	// Set props on the component's stores using reflection
	setComponentProps(comp, c.Props)

	// Render slot content first (with current context)
	slotHTML := childrenToHTML(c.Children, ctx)

	// Create child context for the component with its own prefix
	childCtx := &BuildContext{
		Prefix:          fullCompPrefix,
		Parent:          ctx,
		Bindings:        &CollectedBindings{},
		SlotContent:     slotHTML,
		ParentStoreMap:  ctx.ParentStoreMap,
		CollectedStyles: ctx.CollectedStyles, // Share styles map with parent
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
	parent.EachBlocks = append(parent.EachBlocks, child.EachBlocks...)
	parent.InputBindings = append(parent.InputBindings, child.InputBindings...)
	parent.ClassBindings = append(parent.ClassBindings, child.ClassBindings...)
	parent.ShowIfBindings = append(parent.ShowIfBindings, child.ShowIfBindings...)
	parent.Components = append(parent.Components, child.Components...)
	parent.AttrBindings = append(parent.AttrBindings, child.AttrBindings...)
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
					ctx.Bindings.EachBlocks = append(ctx.Bindings.EachBlocks, childBindings.EachBlocks...)
					ctx.Bindings.ShowIfBindings = append(ctx.Bindings.ShowIfBindings, childBindings.ShowIfBindings...)
					ctx.Bindings.AttrBindings = append(ctx.Bindings.AttrBindings, childBindings.AttrBindings...)
				}
			}
			return content
		}
	}
	return ""
}

// ToHTML generates HTML for a page router node.
// Renders all children wrapped in divs, showing only the current one.
func (r *PageRouterNode) ToHTML(ctx *BuildContext) string {
	if ctx.ChildrenContent == nil {
		return ""
	}

	currentName := componentName(r.Current.Get())
	var sb strings.Builder

	// Render each child in a wrapper div with show/hide based on current component
	for name, content := range ctx.ChildrenContent {
		// Determine visibility
		visible := name == currentName
		style := ""
		if !visible {
			style = ` style="display:none"`
		}

		// Write wrapper div with ID for WASM to toggle
		sb.WriteString(`<div id="page-`)
		sb.WriteString(name)
		sb.WriteString(`"`)
		sb.WriteString(style)
		sb.WriteString(`>`)
		sb.WriteString(content)
		sb.WriteString(`</div>`)

		// Merge child bindings
		if ctx.ChildrenBindings != nil {
			if childBindings, ok := ctx.ChildrenBindings[name]; ok && childBindings != nil {
				ctx.Bindings.TextBindings = append(ctx.Bindings.TextBindings, childBindings.TextBindings...)
				ctx.Bindings.Events = append(ctx.Bindings.Events, childBindings.Events...)
				ctx.Bindings.InputBindings = append(ctx.Bindings.InputBindings, childBindings.InputBindings...)
				ctx.Bindings.ClassBindings = append(ctx.Bindings.ClassBindings, childBindings.ClassBindings...)
				ctx.Bindings.IfBlocks = append(ctx.Bindings.IfBlocks, childBindings.IfBlocks...)
				ctx.Bindings.EachBlocks = append(ctx.Bindings.EachBlocks, childBindings.EachBlocks...)
				ctx.Bindings.ShowIfBindings = append(ctx.Bindings.ShowIfBindings, childBindings.ShowIfBindings...)
				ctx.Bindings.AttrBindings = append(ctx.Bindings.AttrBindings, childBindings.AttrBindings...)
			}
		}
	}

	// Render NotFound div (hidden unless no match)
	if r.NotFound != nil {
		notFoundVisible := ctx.ChildrenContent[currentName] == ""
		style := ""
		if !notFoundVisible {
			style = ` style="display:none"`
		}
		sb.WriteString(`<div id="page-notfound"`)
		sb.WriteString(style)
		sb.WriteString(`>`)
		sb.WriteString(nodeToHTML(r.NotFound, ctx))
		sb.WriteString(`</div>`)
	}

	// Record router binding for WASM
	ctx.Bindings.RouterBindings = append(ctx.Bindings.RouterBindings, RouterBinding_{
		StoreRef: r.Current,
	})

	return sb.String()
}

// nodeToHTML dispatches to the appropriate ToHTML method.
func nodeToHTML(n Node, ctx *BuildContext) string {
	switch node := n.(type) {
	case *Element:
		return node.ToHTML(ctx)
	case *TextNode:
		return node.ToHTML(ctx)
	case *HtmlNode:
		return node.ToHTML(ctx)
	case *Fragment:
		return node.ToHTML(ctx)
	case *BindNode:
		return node.ToHTML(ctx)
	case *BindValueNode:
		return node.ToHTML(ctx)
	case *BindCheckedNode:
		return node.ToHTML(ctx)
	case *ClassIfNode:
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
	case *PageRouterNode:
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

// escapeAttr escapes attribute values (quotes and ampersands).
func escapeAttr(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}

// RenderResult contains all output from HTML rendering.
type RenderResult struct {
	HTML            string
	Bindings        *CollectedBindings
	CollectedStyles map[string]string
}

// RenderHTML is the main entry point for generating HTML from a Node tree.
func RenderHTML(root Node) (string, *CollectedBindings) {
	ctx := NewBuildContext()
	html := nodeToHTML(root, ctx)
	return html, ctx.Bindings
}

// RenderHTMLFull renders HTML and returns all collected data including nested styles.
func RenderHTMLFull(root Node) *RenderResult {
	ctx := NewBuildContext()
	html := nodeToHTML(root, ctx)
	return &RenderResult{
		HTML:            html,
		Bindings:        ctx.Bindings,
		CollectedStyles: ctx.CollectedStyles,
	}
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

// RenderHTMLWithContextFull renders HTML with full context options and returns collected styles.
func RenderHTMLWithContextFull(root Node, opts ...func(*BuildContext)) *RenderResult {
	ctx := NewBuildContext()
	for _, opt := range opts {
		opt(ctx)
	}
	html := nodeToHTML(root, ctx)
	return &RenderResult{
		HTML:            html,
		Bindings:        ctx.Bindings,
		CollectedStyles: ctx.CollectedStyles,
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
