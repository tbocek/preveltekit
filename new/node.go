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
	// Expr returns the Go expression string for code generation (e.g., "component.Score.Get() >= 90")
	Expr() string
	// Deps returns the field names this condition depends on
	Deps() []string
}

// =============================================================================
// Element Node
// =============================================================================

// Element represents an HTML element with attributes and children.
type Element struct {
	Tag      string
	Attrs    []NodeAttr
	Children []Node
}

func (e *Element) nodeType() string { return "element" }

// El creates an element with the given tag, attributes, and children.
// Attrs and children can be mixed - they're separated automatically.
func El(tag string, content ...any) *Element {
	e := &Element{Tag: tag}
	for _, c := range content {
		switch v := c.(type) {
		case NodeAttr:
			e.Attrs = append(e.Attrs, v)
		case Node:
			e.Children = append(e.Children, v)
		case string:
			e.Children = append(e.Children, Text(v))
		}
	}
	return e
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
type HtmlNode struct {
	Parts []any // strings, Nodes, or values to stringify
}

func (h *HtmlNode) nodeType() string { return "html" }

// Html creates a raw HTML node from strings and embedded nodes.
// Example: Html(`<div class="foo">`, p.Bind(store), `</div>`)
func Html(parts ...any) *HtmlNode {
	return &HtmlNode{Parts: parts}
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
	StoreRef any    // The actual store reference (for evaluation)
	StoreID  string // The Go expression for code generation (e.g., "component.Count")
	IsHTML   bool   // true for raw HTML binding
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
	ListID   string                         // Go expression for code generation
	ItemVar  string                         // Variable name for item (default: "item")
	IndexVar string                         // Variable name for index (default: "i")
	Body     func(item any, index int) Node // Template function for each item
	ElseNode []Node                         // Content for empty list
}

func (e *EachNode) nodeType() string { return "each" }

// Each creates a list rendering node.
// The body function receives each item and its index.
func Each[T comparable](list *List[T], body func(item T, index int) Node) *EachNode {
	return &EachNode{
		ListRef:  list,
		ItemVar:  "item",
		IndexVar: "i",
		Body: func(item any, index int) Node {
			return body(item.(T), index)
		},
	}
}

// WithVars sets custom variable names for item and index.
func (e *EachNode) WithVars(itemVar, indexVar string) *EachNode {
	e.ItemVar = itemVar
	e.IndexVar = indexVar
	return e
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
		case *EventAttr:
			c.Events[v.Event] = v.Handler
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

// ChildNode represents a named child component placeholder.
type ChildNode struct {
	Name string
}

func (c *ChildNode) nodeType() string { return "child" }

// Child creates a named child component placeholder for SPA routing.
// Deprecated: Use ChildOf[T]() for type-safe routing.
func Child(name string) *ChildNode {
	return &ChildNode{Name: name}
}

// ChildOf creates a child component placeholder with name derived from type T.
// Example: ChildOf[Basics]() creates a child named "basics"
func ChildOf[T any]() *ChildNode {
	var zero T
	name := reflect.TypeOf(zero).Name()
	// Convert first letter to lowercase
	if len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z' {
		b := []byte(name)
		b[0] = b[0] + 32
		name = string(b)
	}
	return &ChildNode{Name: name}
}

// PageRouterNode renders all registered children and shows the one matching the current component.
type PageRouterNode struct {
	Current  *Store[Component]
	NotFound Node // Optional node to show when no child matches
}

func (r *PageRouterNode) nodeType() string { return "router" }

// PageRouter creates a router that shows the current component.
// All children are pre-rendered at build time; the store controls visibility.
// Example: PageRouter(app.CurrentComponent)
func PageRouter(current *Store[Component]) *PageRouterNode {
	return &PageRouterNode{Current: current}
}

// Default sets the node to show when no child matches.
func (r *PageRouterNode) Default(node Node) *PageRouterNode {
	r.NotFound = node
	return r
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

// ClassIfAttr represents a conditional class binding (legacy, used with Element).
type ClassIfAttr struct {
	ClassName string
	Cond      Condition
}

func (c *ClassIfAttr) attrType() string { return "classif" }

// ClassIfNode represents an HTML element with conditional classes.
// Takes the full opening tag and injects id + merges classes.
type ClassIfNode struct {
	HTML       string
	Conditions []ClassIfCond
	OnClick    func() // Optional click handler
}

type ClassIfCond struct {
	ClassName string
	Cond      Condition
}

func (c *ClassIfNode) nodeType() string { return "classif" }

// ClassIf creates a conditional class node.
// Pass the full opening tag and class/condition pairs.
// Example: ClassIf(`<div class="step">`, "active", store.Eq(1), "completed", store.Gt(1))
func ClassIf(html string, pairs ...any) *ClassIfNode {
	node := &ClassIfNode{HTML: html}
	for i := 0; i+1 < len(pairs); i += 2 {
		if className, ok := pairs[i].(string); ok {
			if cond, ok := pairs[i+1].(Condition); ok {
				node.Conditions = append(node.Conditions, ClassIfCond{
					ClassName: className,
					Cond:      cond,
				})
			}
		}
	}
	return node
}

// WithOnClick adds a click handler to a ClassIfNode.
func (c *ClassIfNode) WithOnClick(handler func()) *ClassIfNode {
	c.OnClick = handler
	return c
}

// ShowIfAttr represents a conditional display binding.
type ShowIfAttr struct {
	Cond Condition
}

func (s *ShowIfAttr) attrType() string { return "showif" }

// ShowIf shows the element when condition is true, hides it otherwise.
func ShowIf(cond Condition) *ShowIfAttr {
	return &ShowIfAttr{Cond: cond}
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
	Template string   // e.g., "/user/{UserID}"
	Stores   []any    // Store references
	StoreIDs []string // Store expressions for code generation
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

// EventAttr represents an event handler binding.
type EventAttr struct {
	Event     string
	Handler   func()   // Handler function (wrap args in closure)
	Modifiers []string // ["preventDefault", "stopPropagation"]
}

func (e *EventAttr) attrType() string { return "event" }

// OnClick creates a click event handler.
func OnClick(handler func()) *EventAttr {
	return &EventAttr{Event: "click", Handler: handler}
}

// OnSubmit creates a submit event handler.
func OnSubmit(handler func()) *EventAttr {
	return &EventAttr{Event: "submit", Handler: handler}
}

// OnInput creates an input event handler.
func OnInput(handler func()) *EventAttr {
	return &EventAttr{Event: "input", Handler: handler}
}

// OnChange creates a change event handler.
func OnChange(handler func()) *EventAttr {
	return &EventAttr{Event: "change", Handler: handler}
}

// OnEvent creates a custom event handler.
func OnEvent(event string, handler func()) *EventAttr {
	return &EventAttr{Event: event, Handler: handler}
}

// PreventDefault adds the preventDefault modifier.
func (e *EventAttr) PreventDefault() *EventAttr {
	e.Modifiers = append(e.Modifiers, "preventDefault")
	return e
}

// StopPropagation adds the stopPropagation modifier.
func (e *EventAttr) StopPropagation() *EventAttr {
	e.Modifiers = append(e.Modifiers, "stopPropagation")
	return e
}

// BindValueNode represents a two-way binding to an input's value.
// It wraps a complete HTML element and injects id/value attributes.
type BindValueNode struct {
	HTML    string // The HTML element, e.g. `<input type="text">`
	Store   any    // *Store[string] or *Store[int]
	StoreID string // Go expression for code generation
}

func (b *BindValueNode) nodeType() string { return "bindvalue" }

// BindValue creates a two-way bound input element.
// Example: BindValue(`<input type="text" placeholder="Name">`, nameStore)
func BindValue[T any](html string, store *Store[T]) *BindValueNode {
	return &BindValueNode{HTML: html, Store: store}
}

// BindCheckedNode represents a two-way binding to a checkbox's checked state.
// It wraps a complete HTML element and injects id/checked attributes.
type BindCheckedNode struct {
	HTML    string // The HTML element, e.g. `<input type="checkbox">`
	Store   any    // *Store[bool]
	StoreID string // Go expression for code generation
}

func (b *BindCheckedNode) nodeType() string { return "bindchecked" }

// BindChecked creates a two-way bound checkbox element.
// Example: BindChecked(`<input type="checkbox">`, isCheckedStore)
func BindChecked(html string, store *Store[bool]) *BindCheckedNode {
	return &BindCheckedNode{HTML: html, Store: store}
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
	storeID  string
	deps     []string
	evalFunc func() bool
}

func (c *StoreCondition) Eval() bool { return c.evalFunc() }
func (c *StoreCondition) Expr() string {
	return c.storeID + ".Get() " + c.Op + " " + anyToString(c.Operand)
}
func (c *StoreCondition) Deps() []string { return c.deps }

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

// BoolCondition wraps a bool store for use in If nodes.
type BoolCondition struct {
	Store   *Store[bool] // Exported for address-based resolution
	storeID string
	negate  bool
	deps    []string
}

func (c *BoolCondition) Eval() bool {
	v := c.Store.Get()
	if c.negate {
		return !v
	}
	return v
}
func (c *BoolCondition) Expr() string {
	if c.negate {
		return "!" + c.storeID + ".Get()"
	}
	return c.storeID + ".Get()"
}
func (c *BoolCondition) Deps() []string { return c.deps }

// IsTrue creates a condition that's true when the bool store is true.
func IsTrue(s *Store[bool]) Condition {
	return &BoolCondition{Store: s}
}

// IsFalse creates a condition that's true when the bool store is false.
func IsFalse(s *Store[bool]) Condition {
	return &BoolCondition{Store: s, negate: true}
}
