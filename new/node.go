package preveltekit

import "reflect"

// Node represents a node in the declarative UI tree.
// Nodes are constructed at build-time and walked to generate:
// 1. HTML with markers (for SSR)
// 2. WASM wiring code (for hydration)
type Node interface {
	nodeType() string
}

// NodeAttr represents an attribute that can be applied to an element.
type NodeAttr interface {
	attrType() string
}

// Condition represents a boolean condition for If nodes.
type Condition interface {
	// Eval returns the current boolean value
	Eval() bool
	// Deps returns the field names this condition depends on
	Deps() []string
}

// =============================================================================
// Text Node
// =============================================================================

// TextNode represents static text content.
type TextNode struct {
	Content string
}

func (t *TextNode) nodeType() string { return "text" }

// Text creates a text node.
func Text(content string) *TextNode {
	return &TextNode{Content: content}
}

// =============================================================================
// Raw HTML Node
// =============================================================================

// HtmlNode represents raw HTML with embedded nodes.
// It allows mixing raw HTML strings with dynamic nodes like Bind, If, Each.
// Supports chainable AttrIf, On, and Bind methods.
type HtmlNode struct {
	Parts      []any        // strings, Nodes, or values to stringify
	AttrConds  []*AttrCond  // conditional attributes applied to first tag
	Events     []*HtmlEvent // event bindings applied to first tag
	BoundStore any          // two-way binding store (*Store[string], *Store[int], *Store[bool])
}

// AttrCond represents a conditional attribute binding.
// Used by HtmlNode.AttrIf() to conditionally set attribute values.
type AttrCond struct {
	Name       string    // attribute name (e.g., "class", "href", "disabled")
	Cond       Condition // condition to evaluate
	TrueValue  any       // value when true: string or *Store[T]
	FalseValue any       // value when false: string or *Store[T] (optional)
}

// HtmlEvent represents an event binding for HtmlNode.
// Used by HtmlNode.On() to attach event handlers.
type HtmlEvent struct {
	ID    string // unique handler ID for registry lookup
	Event string // event name (e.g., "click", "submit")
}

func (h *HtmlNode) nodeType() string { return "html" }

// wasmStringsToRemove collects string parts from Html() calls during SSR.
// A post-build step uses this list to zero them out in the WASM data section.
var wasmStringsToRemove = make(map[string]int)

// Html creates a raw HTML node from strings and embedded nodes.
// Example: Html(`<div class="foo">`, p.Bind(store), `</div>`)
func Html(parts ...any) *HtmlNode {
	if wasmStringsToRemove != nil {
		for _, p := range parts {
			if s, ok := p.(string); ok && len(s) > 0 {
				wasmStringsToRemove[s]++
			}
		}
	}
	return &HtmlNode{Parts: parts}
}

// AttrIf adds a conditional attribute to the first HTML tag.
// Values can be string literals or *Store[T] for reactive values.
// Multiple AttrIf calls for the same attribute name merge additively.
//
// Examples:
//
//	Html(`<button>`).AttrIf("class", cond, "active")              // adds "active" when true
//	Html(`<button>`).AttrIf("class", cond, "active", "inactive")  // "active" when true, "inactive" when false
//	Html(`<a>`).AttrIf("href", cond, urlStore, "/fallback")       // reactive value with fallback
func (h *HtmlNode) AttrIf(name string, cond Condition, values ...any) *HtmlNode {
	ac := &AttrCond{Name: name, Cond: cond}
	if len(values) >= 1 {
		ac.TrueValue = values[0]
	}
	if len(values) >= 2 {
		ac.FalseValue = values[1]
	}
	h.AttrConds = append(h.AttrConds, ac)
	return h
}

// On attaches an event handler to the first HTML tag.
// The handler ID is auto-generated for registry lookup during hydration.
// Returns the HtmlNode for chaining.
//
// Example:
//
//	Html(`<button>Click</button>`).On("click", handler)
//	Html(`<form>`).On("submit", handler).PreventDefault()
func (h *HtmlNode) On(event string, handler func()) *HtmlNode {
	// Register handler in global registry for WASM hydration
	id := RegisterHandler(handler)
	h.Events = append(h.Events, &HtmlEvent{
		ID:    id,
		Event: event,
	})
	return h
}

// PreventDefault adds the preventDefault modifier to the last event.
// Must be called after On. The modifier is stored in the handler registry
// so WASM can apply event.preventDefault() without needing it in bindings.
func (h *HtmlNode) PreventDefault() *HtmlNode {
	if len(h.Events) > 0 {
		last := h.Events[len(h.Events)-1]
		handlerModifiers[last.ID] = append(handlerModifiers[last.ID], "preventDefault")
	}
	return h
}

// StopPropagation adds the stopPropagation modifier to the last event.
// Must be called after On. The modifier is stored in the handler registry
// so WASM can apply event.stopPropagation() without needing it in bindings.
func (h *HtmlNode) StopPropagation() *HtmlNode {
	if len(h.Events) > 0 {
		last := h.Events[len(h.Events)-1]
		handlerModifiers[last.ID] = append(handlerModifiers[last.ID], "stopPropagation")
	}
	return h
}

// Bind attaches a two-way binding to the first HTML element.
// The store type determines the binding behavior:
//   - *Store[string], *Store[int]: binds to input value
//   - *Store[bool]: binds to checkbox checked state
//
// Example:
//
//	Html(`<input type="text">`).Bind(nameStore)
//	Html(`<input type="checkbox">`).Bind(darkModeStore)
func (h *HtmlNode) Bind(store any) *HtmlNode {
	h.BoundStore = store
	return h
}

// =============================================================================
// Fragment Node (multiple children without wrapper)
// =============================================================================

// Fragment represents multiple nodes without a wrapper element.
type Fragment struct {
	Children []Node
}

func (f *Fragment) nodeType() string { return "fragment" }

// Frag creates a fragment containing multiple nodes.
func Frag(children ...Node) *Fragment {
	return &Fragment{Children: children}
}

// =============================================================================
// Bind Node (reactive text binding)
// =============================================================================

// BindNode represents a reactive binding to a store value.
type BindNode struct {
	StoreRef any  // The actual store reference (for evaluation)
	IsHTML   bool // true for raw HTML binding
}

func (b *BindNode) nodeType() string { return "bind" }

// Bind creates a reactive text binding to a store.
// The store must implement Bindable[T].
func Bind[T any](store *Store[T]) *BindNode {
	return &BindNode{
		StoreRef: store,
		IsHTML:   false,
	}
}

// BindAsHTML creates a reactive HTML binding (renders as innerHTML).
func BindAsHTML[T any](store *Store[T]) *BindNode {
	return &BindNode{
		StoreRef: store,
		IsHTML:   true,
	}
}

// =============================================================================
// If Node (conditional rendering)
// =============================================================================

// IfNode represents conditional rendering with optional else-if and else branches.
type IfNode struct {
	Branches []IfBranch
	ElseNode []Node
}

func (i *IfNode) nodeType() string { return "if" }

// IfBranch represents a single if/else-if branch.
type IfBranch struct {
	Cond     Condition
	Children []Node
}

// If creates a conditional rendering node.
func If(cond Condition, children ...Node) *IfNode {
	return &IfNode{
		Branches: []IfBranch{{Cond: cond, Children: children}},
	}
}

// ElseIf adds an else-if branch to the conditional.
func (i *IfNode) ElseIf(cond Condition, children ...Node) *IfNode {
	i.Branches = append(i.Branches, IfBranch{Cond: cond, Children: children})
	return i
}

// Else adds an else branch to the conditional.
func (i *IfNode) Else(children ...Node) *IfNode {
	i.ElseNode = children
	return i
}

// =============================================================================
// Each Node (list rendering)
// =============================================================================

// EachNode represents list iteration with optional else for empty list.
type EachNode struct {
	ListRef  any                            // The actual list reference
	Body     func(item any, index int) Node // Template function for each item
	ElseNode []Node                         // Content for empty list
}

func (e *EachNode) nodeType() string { return "each" }

// Each creates a list rendering node.
// The body function receives each item and its index.
func Each[T comparable](list *List[T], body func(item T, index int) Node) *EachNode {
	return &EachNode{
		ListRef: list,
		Body: func(item any, index int) Node {
			return body(item.(T), index)
		},
	}
}

// Else adds content to show when the list is empty.
func (e *EachNode) Else(children ...Node) *EachNode {
	e.ElseNode = children
	return e
}

// =============================================================================
// Component Node (nested component)
// =============================================================================

// ComponentNode represents a nested component.
type ComponentNode struct {
	Name     string         // Component type name (derived from instance)
	Instance any            // The actual component instance
	Props    map[string]any // Property values
	Events   map[string]any // Event handlers
	Children []Node         // Slot content
}

func (c *ComponentNode) nodeType() string { return "component" }

// Comp creates a nested component node from a component instance.
// The component name is derived from the type via reflection.
// Example: Comp(&Badge{Label: p.New("New")})
func Comp(instance any, content ...any) *ComponentNode {
	// Derive name from type
	t := reflect.TypeOf(instance)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	name := t.Name()

	c := &ComponentNode{
		Name:     name,
		Instance: instance,
		Props:    make(map[string]any),
		Events:   make(map[string]any),
	}
	for _, item := range content {
		switch v := item.(type) {
		case *PropAttr:
			c.Props[v.Name] = v.Value
		case Node:
			c.Children = append(c.Children, v)
		}
	}
	return c
}

// =============================================================================
// Slot Node (for child component content)
// =============================================================================

// SlotNode represents where child content should be inserted.
type SlotNode struct{}

func (s *SlotNode) nodeType() string { return "slot" }

// Slot creates a slot placeholder for child content.
func Slot() *SlotNode {
	return &SlotNode{}
}

// componentName returns the lowercase type name of a component.
func componentName(c Component) string {
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

// =============================================================================
// Attributes
// =============================================================================

// ClassAttr represents a static class attribute.
type ClassAttr struct {
	Classes []string
}

func (c *ClassAttr) attrType() string { return "class" }

// Class creates a class attribute with one or more class names.
func Class(classes ...string) *ClassAttr {
	return &ClassAttr{Classes: classes}
}

// StaticAttr represents a static attribute.
type StaticAttr struct {
	Name  string
	Value string
}

func (s *StaticAttr) attrType() string { return "static" }

// StaticAttribute creates a static attribute.
func StaticAttribute(name, value string) *StaticAttr {
	return &StaticAttr{Name: name, Value: value}
}

// DynAttrAttr represents a dynamic attribute with store bindings.
type DynAttrAttr struct {
	Name     string
	Template string // e.g., "/user/{UserID}"
	Stores   []any  // Store references
}

func (d *DynAttrAttr) attrType() string { return "dynattr" }

// DynAttr creates a dynamic attribute that includes store values.
// Template uses {0}, {1}, etc. as placeholders for store values.
func DynAttr(name, template string, stores ...any) *DynAttrAttr {
	return &DynAttrAttr{
		Name:     name,
		Template: template,
		Stores:   stores,
	}
}

// PropAttr represents a property passed to a child component.
type PropAttr struct {
	Name  string
	Value any
}

func (p *PropAttr) attrType() string { return "prop" }

// Prop creates a property to pass to a child component.
func Prop(name string, value any) *PropAttr {
	return &PropAttr{Name: name, Value: value}
}

// =============================================================================
// Condition Helpers
// =============================================================================

// StoreCondition wraps a store with a comparison for use in If nodes.
type StoreCondition struct {
	Store    any    // Exported for address-based resolution
	Op       string // Exported for class binding resolution
	Operand  any    // Exported for class binding resolution
	evalFunc func() bool
}

func (c *StoreCondition) Eval() bool     { return c.evalFunc() }
func (c *StoreCondition) Deps() []string { return nil }

// Ge creates a >= condition on a store.
func (s *Store[T]) Ge(value T) Condition {
	return &StoreCondition{
		Store:    s,
		Op:       ">=",
		Operand:  value,
		evalFunc: func() bool { return any(s.Get()).(int) >= any(value).(int) },
	}
}

// Gt creates a > condition on a store.
func (s *Store[T]) Gt(value T) Condition {
	return &StoreCondition{
		Store:    s,
		Op:       ">",
		Operand:  value,
		evalFunc: func() bool { return any(s.Get()).(int) > any(value).(int) },
	}
}

// Le creates a <= condition on a store.
func (s *Store[T]) Le(value T) Condition {
	return &StoreCondition{
		Store:    s,
		Op:       "<=",
		Operand:  value,
		evalFunc: func() bool { return any(s.Get()).(int) <= any(value).(int) },
	}
}

// Lt creates a < condition on a store.
func (s *Store[T]) Lt(value T) Condition {
	return &StoreCondition{
		Store:    s,
		Op:       "<",
		Operand:  value,
		evalFunc: func() bool { return any(s.Get()).(int) < any(value).(int) },
	}
}

// Eq creates an == condition on a store.
func (s *Store[T]) Eq(value T) Condition {
	return &StoreCondition{
		Store:   s,
		Op:      "==",
		Operand: value,
		evalFunc: func() bool {
			return anyToString(s.Get()) == anyToString(value)
		},
	}
}

// Ne creates a != condition on a store.
func (s *Store[T]) Ne(value T) Condition {
	return &StoreCondition{
		Store:   s,
		Op:      "!=",
		Operand: value,
		evalFunc: func() bool {
			return anyToString(s.Get()) != anyToString(value)
		},
	}
}

// anyToString converts any value to string without fmt.
func anyToString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return itoa(val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case float64:
		return ftoa(val)
	default:
		return ""
	}
}

func itoa(n int) string {
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

func ftoa(f float64) string {
	if f == 0 {
		return "0"
	}
	neg := f < 0
	if neg {
		f = -f
	}
	// Simple float to string: integer part + 2 decimal places
	intPart := int(f)
	fracPart := int((f - float64(intPart)) * 100)
	s := itoa(intPart)
	if fracPart > 0 {
		s += "."
		if fracPart < 10 {
			s += "0"
		}
		s += itoa(fracPart)
	}
	if neg {
		s = "-" + s
	}
	return s
}

// escapeHTML escapes HTML special characters.
func escapeHTML(s string) string {
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

// BoolCondition wraps a bool store for use in If nodes.
type BoolCondition struct {
	Store  *Store[bool] // Exported for address-based resolution
	negate bool
}

func (c *BoolCondition) Eval() bool {
	v := c.Store.Get()
	if c.negate {
		return !v
	}
	return v
}
func (c *BoolCondition) Deps() []string { return nil }

// IsTrue creates a condition that's true when the bool store is true.
func IsTrue(s *Store[bool]) Condition {
	return &BoolCondition{Store: s}
}

// IsFalse creates a condition that's true when the bool store is false.
func IsFalse(s *Store[bool]) Condition {
	return &BoolCondition{Store: s, negate: true}
}
