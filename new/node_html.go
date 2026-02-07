//go:build !js || !wasm

package preveltekit

import (
	"fmt"
	"reflect"
	"strings"
)

// BuildContext holds state during HTML and wiring code generation.
type BuildContext struct {
	// Embed shared ID counter logic (used by both SSR and WASM)
	IDCounter

	// Parent context for nested components
	Parent *BuildContext

	// Collected bindings during tree walk
	Bindings *CollectedBindings

	// SlotContent holds HTML to be rendered in place of <slot/> elements
	SlotContent string

	// ParentStoreMap maps store pointers to their IDs in the parent component
	// Used to resolve dynamic props that share parent stores
	ParentStoreMap map[uintptr]string

	// RouteGroupID is the ID of the current route group (for component container binding)
	RouteGroupID string

	// CollectedStyles holds CSS from nested components (deduplicated by component name)
	CollectedStyles map[string]string
}

// CollectedBindings stores all bindings found during tree walking.
type CollectedBindings struct {
	TextBindings        []TextBinding        `json:"TextBindings"`
	Events              []EventBinding       `json:"Events"`
	IfBlocks            []IfBlock            `json:"IfBlocks"`
	EachBlocks          []EachBlock          `json:"EachBlocks"`
	InputBindings       []InputBinding       `json:"InputBindings"`
	AttrBindings        []AttrBinding        `json:"AttrBindings"`
	AttrCondBindings    []AttrCondBinding    `json:"AttrCondBindings"`
	Components          []ComponentBinding   `json:"Components"`
	ComponentStore      *Store[Component]    `json:"-"` // For WASM to subscribe to navigation
	ComponentContainers []ComponentContainer `json:"ComponentContainers,omitempty"`
	ComponentBlocks     []ComponentBlock     `json:"ComponentBlocks,omitempty"`
}

// ComponentContainer maps a route group ID to its DOM container
type ComponentContainer struct {
	ID          string `json:"ID"`          // Route group ID, e.g., "main"
	ContainerID string `json:"ContainerID"` // DOM element ID, e.g., "content"
}

// ComponentBlock represents a Store[Component]'s pre-baked branches.
// Like IfBlock but keyed by component type name instead of store conditions.
// HTML: <span>active component content</span><!--r0--> where content swaps on store change.
type ComponentBlock struct {
	MarkerID string            `json:"MarkerID"` // Comment marker, e.g., "r0"
	StoreID  string            `json:"StoreID"`  // Store ID of the component store
	Branches []ComponentBranch `json:"Branches"`
}

// ComponentBranch represents one component's pre-baked content.
type ComponentBranch struct {
	Name     string             `json:"Name"`               // Component type name, e.g., "basics", "components"
	HTML     string             `json:"HTML"`               // Pre-rendered HTML for this component
	Bindings *CollectedBindings `json:"Bindings,omitempty"` // Nested bindings for this component's content
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

// EventBinding binds an event handler to a DOM element by its id attribute.
// HTML: <button id="basics_ev0"> triggers the handler on click.
type EventBinding struct {
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

// InputBinding binds an input element to a store for two-way data binding.
// HTML: <input id="basics_b0"> syncs value with store.
type InputBinding struct {
	ElementID string `json:"element_id"` // Element id attribute, e.g., "basics_b0"
	StoreID   string `json:"store_id"`   // Store path, e.g., "basics.Name"
	StoreRef  any    `json:"-"`          // Actual store pointer (for resolution)
	BindType  string `json:"bind_type"`  // Binding type: "value" or "checked"
}

// AttrBinding binds a dynamic attribute value to stores.
// HTML: <div data-attrbind="basics_a0" data-value="{0}"> where {0} is replaced.
type AttrBinding struct {
	ElementID string   `json:"element_id"` // Element id (via data-attrbind), e.g., "basics_a0"
	AttrName  string   `json:"attr_name"`  // Attribute name, e.g., "data-value"
	Template  string   `json:"template"`   // Template with placeholders, e.g., "{0}"
	StoreIDs  []string `json:"store_ids"`  // Store paths for placeholders
	StoreRefs []any    `json:"-"`          // Actual store pointers (for resolution)
}

// AttrCondBinding binds a conditional attribute value to a condition.
// Used by HtmlNode.AttrIf() for conditional attribute rendering.
// HTML: <div id="basics_cl0" class="active"> where attribute value changes reactively.
type AttrCondBinding struct {
	ElementID     string   `json:"element_id"`               // Element id attribute
	AttrName      string   `json:"attr_name"`                // Attribute name (e.g., "class", "href")
	TrueValue     string   `json:"true_value"`               // Value when condition is true
	FalseValue    string   `json:"false_value,omitempty"`    // Value when condition is false
	TrueStoreRef  any      `json:"-"`                        // Store for true value (if dynamic)
	FalseStoreRef any      `json:"-"`                        // Store for false value (if dynamic)
	TrueStoreID   string   `json:"true_store_id,omitempty"`  // Store path for true value
	FalseStoreID  string   `json:"false_store_id,omitempty"` // Store path for false value
	CondStoreRef  any      `json:"-"`                        // Store for condition evaluation
	Op            string   `json:"op,omitempty"`             // Comparison operator
	Operand       string   `json:"operand,omitempty"`        // Comparison operand
	IsBool        bool     `json:"is_bool,omitempty"`        // True if simple boolean condition
	Deps          []string `json:"deps,omitempty"`           // Store dependencies for reactivity
}

// ComponentBinding represents a nested component instance.
type ComponentBinding struct {
	ElementID string            // Component prefix, e.g., "basics_comp0"
	Name      string            // Component type name, e.g., "Button"
	Props     map[string]string // Static prop values
	Events    map[string]string // Event handler mappings
	SlotHTML  string            // Slot content HTML
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
		IDCounter:      IDCounter{Prefix: prefix},
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
//    - Used by: Events, InputBindings, AttrBindings, AttrCondBindings
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

// ID generation methods are inherited from embedded IDCounter.
// See id.go for: NextEventID, NextBindID, NextClassID, NextAttrID,
// NextTextMarker, NextIfMarker, NextEachMarker, NextCompMarker,
// FullElementID, FullMarkerID

// =============================================================================
// ToHTML implementations
// =============================================================================

// ToHTML generates HTML for a text node.
func (t *TextNode) ToHTML(ctx *BuildContext) string {
	return escapeHTML(t.Content)
}

// ToHTML generates HTML for a raw HTML node with embedded nodes.
func (h *HtmlNode) ToHTML(ctx *BuildContext) string {
	// First, render parts to get base HTML
	html := h.renderParts(ctx)

	// If we have AttrConds or Events from chained methods, inject into first tag
	if len(h.AttrConds) > 0 || len(h.Events) > 0 {
		html = h.injectChainedAttrs(html, ctx)
	}

	return html
}

// renderParts renders the Parts slice to HTML string.
func (h *HtmlNode) renderParts(ctx *BuildContext) string {
	var sb strings.Builder

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
			sb.WriteString(attrToHTMLString(v, ctx))
		case *Store[string]:
			// Auto-bind stores for reactivity
			bindNode := &BindNode{StoreRef: v, IsHTML: false}
			sb.WriteString(bindNode.ToHTML(ctx))
		case *Store[int]:
			bindNode := &BindNode{StoreRef: v, IsHTML: false}
			sb.WriteString(bindNode.ToHTML(ctx))
		case *Store[bool]:
			bindNode := &BindNode{StoreRef: v, IsHTML: false}
			sb.WriteString(bindNode.ToHTML(ctx))
		case *Store[float64]:
			bindNode := &BindNode{StoreRef: v, IsHTML: false}
			sb.WriteString(bindNode.ToHTML(ctx))
		case *Store[Component]:
			// Component store rendering: bake ALL option branches into a ComponentBlock
			// (like IfBlock bakes all conditional branches).
			comp := v.Get()
			if comp != nil && len(v.Options()) > 0 {
				// Generate a marker (like IfBlock markers)
				localMarker := ctx.NextRouteMarker()
				markerID := ctx.FullMarkerID(localMarker)

				block := ComponentBlock{
					MarkerID: markerID,
					StoreID:  v.ID(),
				}

				// Render ALL option components, deduplicated by type name
				var activeHTML string
				seen := make(map[string]bool)
				for _, opt := range v.Options() {
					optComp, ok := opt.(Component)
					if !ok || optComp == nil {
						continue
					}
					name := componentName(optComp)
					if seen[name] {
						continue
					}
					seen[name] = true

					// Render in isolated child context (like IfBlock branches)
					branchCtx := ctx.Child(name)
					branchHTML := nodeToHTML(optComp.Render(), branchCtx)

					// Resolve bindings within this branch
					childStoreMap := buildStoreMap(optComp, name)
					resolveBindings(branchCtx.Bindings, childStoreMap, name, optComp)

					block.Branches = append(block.Branches, ComponentBranch{
						Name:     name,
						HTML:     branchHTML,
						Bindings: branchCtx.Bindings,
					})

					// Track which branch is active for this SSR render
					if optComp == comp {
						activeHTML = branchHTML
					}
				}

				ctx.Bindings.ComponentBlocks = append(ctx.Bindings.ComponentBlocks, block)

				// Emit: <span>{activeHTML}</span><!--markerID-->
				// Same pattern as IfBlock output
				sb.WriteString(fmt.Sprintf("<span>%s</span><!--%s-->", activeHTML, markerID))
			} else if comp != nil {
				// Fallback: no options, render current component directly
				name := componentName(comp)
				childCtx := ctx.Child(name)
				html := nodeToHTML(comp.Render(), childCtx)
				sb.WriteString(html)

				childStoreMap := buildStoreMap(comp, name)
				resolveBindings(childCtx.Bindings, childStoreMap, name, comp)
				mergeNestedBindings(ctx.Bindings, childCtx.Bindings)
			}
		default:
			// Convert other values to string and escape
			sb.WriteString(escapeHTML(fmt.Sprintf("%v", v)))
		}
	}
	return sb.String()
}

// injectChainedAttrs injects AttrConds and Events into the first HTML tag.
func (h *HtmlNode) injectChainedAttrs(html string, ctx *BuildContext) string {
	// For events, use the first event's ID directly (handlers are registered by ID)
	// For AttrConds without events, we still need a generated ID for now
	var elementID string
	if len(h.Events) > 0 {
		elementID = h.Events[0].ID
	} else {
		// AttrConds still need an element ID for reactive updates
		localID := ctx.NextClassID()
		elementID = ctx.FullElementID(localID)
	}

	// Collect active values for each attribute (for SSR rendering)
	// Map: attr name -> list of values to add
	attrValues := make(map[string][]string)

	// Process AttrConds (still need bindings for reactive attribute updates)
	for _, ac := range h.AttrConds {
		// Extract condition info
		var condStoreRef any
		var op, operand string
		var isBool bool
		if sc, ok := ac.Cond.(*StoreCondition); ok {
			condStoreRef = sc.Store
			op = sc.Op
			operand = fmt.Sprintf("%v", sc.Operand)
		} else if bc, ok := ac.Cond.(*BoolCondition); ok {
			condStoreRef = bc.Store
			isBool = true
		}

		// Determine true/false values and store refs
		trueVal, trueStoreRef := evalAttrValue(ac.TrueValue)
		falseVal, falseStoreRef := evalAttrValue(ac.FalseValue)

		// Record binding for WASM hydration (AttrConds still need this for reactivity)
		ctx.Bindings.AttrCondBindings = append(ctx.Bindings.AttrCondBindings, AttrCondBinding{
			ElementID:     elementID,
			AttrName:      ac.Name,
			TrueValue:     trueVal,
			FalseValue:    falseVal,
			TrueStoreRef:  trueStoreRef,
			FalseStoreRef: falseStoreRef,
			CondStoreRef:  condStoreRef,
			Op:            op,
			Operand:       operand,
			IsBool:        isBool,
		})

		// Evaluate for SSR
		if ac.Cond.Eval() {
			if trueVal != "" {
				attrValues[ac.Name] = append(attrValues[ac.Name], trueVal)
			}
		} else if falseVal != "" {
			attrValues[ac.Name] = append(attrValues[ac.Name], falseVal)
		}
	}

	// Build extra attributes string with event types for WASM
	var extraAttrs string
	if len(h.Events) > 0 {
		var eventNames []string
		for _, ev := range h.Events {
			eventNames = append(eventNames, ev.Event)
			// Add event binding so WASM knows to bind this handler
			ctx.Bindings.Events = append(ctx.Bindings.Events, EventBinding{
				ElementID: ev.ID, // Use the user-provided handler ID
				Event:     ev.Event,
				Modifiers: ev.Modifiers,
			})
		}
		extraAttrs = fmt.Sprintf(` data-on="%s"`, strings.Join(eventNames, ","))
	}

	// Inject id and merge attributes into first tag
	return injectIDAndMergeAttrs(html, elementID, attrValues, extraAttrs)
}

// evalAttrValue extracts string value and store reference from an AttrCond value.
// Returns (stringValue, storeRef). If value is a store, stringValue is the current value.
func evalAttrValue(v any) (string, any) {
	if v == nil {
		return "", nil
	}
	switch val := v.(type) {
	case string:
		return val, nil
	case *Store[string]:
		return val.Get(), val
	case *Store[int]:
		return itoa(val.Get()), val
	case *Store[bool]:
		if val.Get() {
			return "true", val
		}
		return "false", val
	default:
		return "", nil
	}
}

// injectIDAndMergeAttrs injects id and merges attribute values into the first HTML tag.
// For "class", merges with existing class attribute. For others, values are space-joined.
func injectIDAndMergeAttrs(html, id string, attrValues map[string][]string, extraAttrs string) string {
	tagEnd := findTagEnd(html)
	if tagEnd == -1 {
		return html
	}

	openingTag := html[:tagEnd]
	rest := html[tagEnd:]

	// Check for self-closing tag
	if tagEnd > 0 && html[tagEnd-1] == '/' {
		openingTag = html[:tagEnd-1]
		rest = html[tagEnd-1:]
	}

	// Build new attributes to inject
	var newAttrs strings.Builder
	newAttrs.WriteString(fmt.Sprintf(`id="%s"`, id))

	// Handle class attribute specially - merge with existing
	if classes, ok := attrValues["class"]; ok && len(classes) > 0 {
		classIdx := strings.Index(openingTag, `class="`)
		if classIdx != -1 {
			// Find existing class value
			classStart := classIdx + 7
			classEnd := strings.Index(openingTag[classStart:], `"`)
			if classEnd != -1 {
				classEnd += classStart
				existingClasses := openingTag[classStart:classEnd]
				// Merge: existing + new classes
				mergedClasses := existingClasses
				for _, c := range classes {
					if c != "" {
						mergedClasses += " " + c
					}
				}
				// Rebuild opening tag without old class attr
				openingTag = openingTag[:classIdx] + openingTag[classEnd+1:]
				newAttrs.WriteString(fmt.Sprintf(` class="%s"`, strings.TrimSpace(mergedClasses)))
			}
		} else {
			// No existing class, add new one
			newAttrs.WriteString(fmt.Sprintf(` class="%s"`, strings.Join(classes, " ")))
		}
		delete(attrValues, "class")
	}

	// Handle other attributes
	for name, values := range attrValues {
		if len(values) > 0 {
			// Check if attribute already exists
			attrPattern := name + `="`
			attrIdx := strings.Index(openingTag, attrPattern)
			if attrIdx != -1 {
				// Find existing value
				attrStart := attrIdx + len(attrPattern)
				attrEnd := strings.Index(openingTag[attrStart:], `"`)
				if attrEnd != -1 {
					attrEnd += attrStart
					existingValue := openingTag[attrStart:attrEnd]
					// Merge values
					mergedValue := existingValue
					for _, v := range values {
						if v != "" {
							mergedValue += " " + v
						}
					}
					// Remove old attr from opening tag
					openingTag = openingTag[:attrIdx] + openingTag[attrEnd+1:]
					newAttrs.WriteString(fmt.Sprintf(` %s="%s"`, name, strings.TrimSpace(mergedValue)))
				}
			} else {
				// No existing attr, add new one
				newAttrs.WriteString(fmt.Sprintf(` %s="%s"`, name, strings.Join(values, " ")))
			}
		}
	}

	// Add extra attrs (like data-event)
	newAttrs.WriteString(extraAttrs)

	// Find insertion point (after tag name and before existing attrs or >)
	// Look for first space or end of tag name
	insertIdx := 0
	for i := 1; i < len(openingTag); i++ { // Start after '<'
		if openingTag[i] == ' ' || openingTag[i] == '/' {
			insertIdx = i
			break
		}
	}
	if insertIdx == 0 {
		insertIdx = len(openingTag)
	}

	// Rebuild: <tagname + new attrs + existing attrs + rest
	return openingTag[:insertIdx] + " " + newAttrs.String() + openingTag[insertIdx:] + rest
}

// attrToHTMLString renders a NodeAttr as an HTML attribute string.
// Used when attributes are embedded directly in Html() nodes.
func attrToHTMLString(attr NodeAttr, ctx *BuildContext) string {
	switch a := attr.(type) {
	case *ClassAttr:
		return fmt.Sprintf(`class="%s"`, strings.Join(a.Classes, " "))
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
		ctx.Bindings.AttrBindings = append(ctx.Bindings.AttrBindings, AttrBinding{
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
	// Get store ID and value directly from the store
	var storeID string
	var value string
	switch s := b.StoreRef.(type) {
	case *Store[string]:
		storeID = s.ID()
		value = s.Get()
	case *Store[int]:
		storeID = s.ID()
		value = fmt.Sprintf("%d", s.Get())
	case *Store[bool]:
		storeID = s.ID()
		value = fmt.Sprintf("%t", s.Get())
	case *Store[float64]:
		storeID = s.ID()
		value = fmt.Sprintf("%g", s.Get())
	default:
		value = ""
	}

	// Generate a marker ID and register a TextBinding in the context
	localMarker := ctx.NextTextMarker()
	markerID := ctx.FullMarkerID(localMarker)

	ctx.Bindings.TextBindings = append(ctx.Bindings.TextBindings, TextBinding{
		MarkerID: markerID,
		StoreID:  storeID,
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
	// Get store ID and current value directly from the store
	var storeID string
	var value string
	switch s := b.Store.(type) {
	case *Store[string]:
		storeID = s.ID()
		value = s.Get()
	case *Store[int]:
		storeID = s.ID()
		value = fmt.Sprintf("%d", s.Get())
	}

	// Register input binding in context
	ctx.Bindings.InputBindings = append(ctx.Bindings.InputBindings, InputBinding{
		ElementID: storeID,
		StoreID:   storeID,
		StoreRef:  b.Store,
		BindType:  "value",
	})

	// Check if it's a textarea (value goes as content, not attribute)
	if strings.HasPrefix(strings.TrimSpace(strings.ToLower(b.HTML)), "<textarea") {
		return injectTextareaContent(b.HTML, storeID, value)
	}

	// Inject store ID as element id and current value
	return injectAttrs(b.HTML, fmt.Sprintf(`id="%s" value="%s"`, storeID, escapeAttr(value)))
}

// ToHTML generates HTML for a BindChecked node (checkbox binding).
// Parses the HTML string and injects id and checked attributes.
func (b *BindCheckedNode) ToHTML(ctx *BuildContext) string {
	// Get store ID and checked state directly from the store
	var storeID string
	checked := ""
	if s, ok := b.Store.(*Store[bool]); ok {
		storeID = s.ID()
		if s.Get() {
			checked = " checked"
		}
	}

	// Register input binding in context
	ctx.Bindings.InputBindings = append(ctx.Bindings.InputBindings, InputBinding{
		ElementID: storeID,
		StoreID:   storeID,
		StoreRef:  b.Store,
		BindType:  "checked",
	})

	// Inject store ID as element id and checked state
	return injectAttrs(b.HTML, fmt.Sprintf(`id="%s"%s`, storeID, checked))
}

// =============================================================================
// HTML Attribute Injection
// =============================================================================

// injectAttrs injects attributes into an HTML element string.
// Finds the first > and inserts the attrs just before it.
// Example: injectAttrs(`<input type="text">`, `id="foo"`) -> `<input type="text" id="foo">`
func injectAttrs(html, attrs string) string {
	tagEnd := findTagEnd(html)
	if tagEnd == -1 {
		return html + " " + attrs
	}
	// Check if it's a self-closing tag />
	if tagEnd > 0 && html[tagEnd-1] == '/' {
		return html[:tagEnd-1] + " " + attrs + " />" + html[tagEnd+1:]
	}
	return html[:tagEnd] + " " + attrs + html[tagEnd:]
}

// injectTextareaContent injects id attribute and replaces textarea content.
// Example: <textarea></textarea> -> <textarea id="x">value</textarea>
func injectTextareaContent(html, id, value string) string {
	tagEnd := findTagEnd(html)
	if tagEnd == -1 {
		return html
	}

	openTag := html[:tagEnd]
	rest := html[tagEnd+1:]

	// Find the closing </textarea>
	closeTagIdx := strings.Index(strings.ToLower(rest), "</textarea>")
	if closeTagIdx == -1 {
		return openTag + fmt.Sprintf(` id="%s"`, id) + ">" + rest
	}

	return openTag + fmt.Sprintf(` id="%s"`, id) + ">" + escapeHTML(value) + rest[closeTagIdx:]
}

// findTagEnd returns the index of the first '>' in html, or -1 if not found.
func findTagEnd(html string) int {
	for i := 0; i < len(html); i++ {
		if html[i] == '>' {
			return i
		}
	}
	return -1
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
			IDCounter:      ctx.IDCounter, // Copy all counters
			Bindings:       &CollectedBindings{},
			ParentStoreMap: ctx.ParentStoreMap,
		}
		branchHTML := childrenToHTML(branch.Children, branchCtx)

		// Update parent counters
		ctx.IDCounter = branchCtx.IDCounter

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
			IDCounter:      ctx.IDCounter, // Copy all counters
			Bindings:       &CollectedBindings{},
			ParentStoreMap: ctx.ParentStoreMap,
		}
		elseHTML := childrenToHTML(i.ElseNode, elseCtx)

		// Update parent counters
		ctx.IDCounter = elseCtx.IDCounter

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
		IDCounter:       IDCounter{Prefix: fullCompPrefix},
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
	parent.Components = append(parent.Components, child.Components...)
	parent.AttrBindings = append(parent.AttrBindings, child.AttrBindings...)
	parent.AttrCondBindings = append(parent.AttrCondBindings, child.AttrCondBindings...)
	parent.ComponentBlocks = append(parent.ComponentBlocks, child.ComponentBlocks...)
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

// nodeToHTML dispatches to the appropriate ToHTML method.
func nodeToHTML(n Node, ctx *BuildContext) string {
	switch node := n.(type) {
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
	case *IfNode:
		return node.ToHTML(ctx)
	case *EachNode:
		return node.ToHTML(ctx)
	case *ComponentNode:
		return node.ToHTML(ctx)
	case *SlotNode:
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

// extractLastID finds the last id="..." value in an HTML string.
// Returns empty string if no id attribute found.
func extractLastID(html string) string {
	// Find all id="..." patterns and return the last one
	lastID := ""
	for i := 0; i < len(html); i++ {
		// Look for id="
		if i+4 < len(html) && html[i:i+4] == `id="` {
			start := i + 4
			end := start
			for end < len(html) && html[end] != '"' {
				end++
			}
			if end < len(html) {
				lastID = html[start:end]
			}
			i = end
		}
	}
	return lastID
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

// WithParentStoreMapCtx sets the parent store map on the build context.
func WithParentStoreMapCtx(storeMap map[uintptr]string) func(*BuildContext) {
	return func(ctx *BuildContext) {
		ctx.ParentStoreMap = storeMap
	}
}

// WithRouteGroupIDCtx sets the route group ID on the build context.
func WithRouteGroupIDCtx(id string) func(*BuildContext) {
	return func(ctx *BuildContext) {
		ctx.RouteGroupID = id
	}
}
