package preveltekit

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

// Common HTML elements
func Div(content ...any) *Element      { return El("div", content...) }
func Span(content ...any) *Element     { return El("span", content...) }
func P(content ...any) *Element        { return El("p", content...) }
func H1(content ...any) *Element       { return El("h1", content...) }
func H2(content ...any) *Element       { return El("h2", content...) }
func H3(content ...any) *Element       { return El("h3", content...) }
func H4(content ...any) *Element       { return El("h4", content...) }
func H5(content ...any) *Element       { return El("h5", content...) }
func H6(content ...any) *Element       { return El("h6", content...) }
func Strong(content ...any) *Element   { return El("strong", content...) }
func Em(content ...any) *Element       { return El("em", content...) }
func Small(content ...any) *Element    { return El("small", content...) }
func A(content ...any) *Element        { return El("a", content...) }
func Button(content ...any) *Element   { return El("button", content...) }
func Input(content ...any) *Element    { return El("input", content...) }
func Label(content ...any) *Element    { return El("label", content...) }
func Form(content ...any) *Element     { return El("form", content...) }
func Ul(content ...any) *Element       { return El("ul", content...) }
func Ol(content ...any) *Element       { return El("ol", content...) }
func Li(content ...any) *Element       { return El("li", content...) }
func Nav(content ...any) *Element      { return El("nav", content...) }
func Main(content ...any) *Element     { return El("main", content...) }
func Section(content ...any) *Element  { return El("section", content...) }
func Article(content ...any) *Element  { return El("article", content...) }
func Header(content ...any) *Element   { return El("header", content...) }
func Footer(content ...any) *Element   { return El("footer", content...) }
func Aside(content ...any) *Element    { return El("aside", content...) }
func Table(content ...any) *Element    { return El("table", content...) }
func Thead(content ...any) *Element    { return El("thead", content...) }
func Tbody(content ...any) *Element    { return El("tbody", content...) }
func Tr(content ...any) *Element       { return El("tr", content...) }
func Th(content ...any) *Element       { return El("th", content...) }
func Td(content ...any) *Element       { return El("td", content...) }
func Img(content ...any) *Element      { return El("img", content...) }
func Pre(content ...any) *Element      { return El("pre", content...) }
func Code(content ...any) *Element     { return El("code", content...) }
func Br() *Element                     { return El("br") }
func Hr() *Element                     { return El("hr") }
func Textarea(content ...any) *Element { return El("textarea", content...) }

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
	Name     string         // Component type name
	Props    map[string]any // Property values
	Events   map[string]any // Event handlers
	Children []Node         // Slot content
}

func (c *ComponentNode) nodeType() string { return "component" }

// Comp creates a nested component node.
func Comp(name string, content ...any) *ComponentNode {
	c := &ComponentNode{
		Name:   name,
		Props:  make(map[string]any),
		Events: make(map[string]any),
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
func Child(name string) *ChildNode {
	return &ChildNode{Name: name}
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

// ClassIfAttr represents a conditional class binding.
type ClassIfAttr struct {
	ClassName string
	Cond      Condition
}

func (c *ClassIfAttr) attrType() string { return "classif" }

// ClassIf creates a conditional class that's applied when the condition is true.
func ClassIf(className string, cond Condition) *ClassIfAttr {
	return &ClassIfAttr{ClassName: className, Cond: cond}
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

// Common attribute shortcuts
func Id(value string) *StaticAttr          { return StaticAttribute("id", value) }
func Type(value string) *StaticAttr        { return StaticAttribute("type", value) }
func Href(value string) *StaticAttr        { return StaticAttribute("href", value) }
func Src(value string) *StaticAttr         { return StaticAttribute("src", value) }
func Alt(value string) *StaticAttr         { return StaticAttribute("alt", value) }
func Placeholder(value string) *StaticAttr { return StaticAttribute("placeholder", value) }
func Name(value string) *StaticAttr        { return StaticAttribute("name", value) }
func Value(value string) *StaticAttr       { return StaticAttribute("value", value) }
func Disabled() *StaticAttr                { return StaticAttribute("disabled", "disabled") }
func Readonly() *StaticAttr                { return StaticAttribute("readonly", "readonly") }

// Attr creates a static attribute with name and value.
func Attr(name, value string) *StaticAttr { return StaticAttribute(name, value) }

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
	Handler   any      // Method reference
	HandlerID string   // Go expression for code generation
	Args      []any    // Arguments to pass
	Modifiers []string // ["preventDefault", "stopPropagation"]
}

func (e *EventAttr) attrType() string { return "event" }

// OnClick creates a click event handler.
func OnClick(handler any, args ...any) *EventAttr {
	return &EventAttr{Event: "click", Handler: handler, Args: args}
}

// OnSubmit creates a submit event handler.
func OnSubmit(handler any, args ...any) *EventAttr {
	return &EventAttr{Event: "submit", Handler: handler, Args: args}
}

// OnInput creates an input event handler.
func OnInput(handler any, args ...any) *EventAttr {
	return &EventAttr{Event: "input", Handler: handler, Args: args}
}

// OnChange creates a change event handler.
func OnChange(handler any, args ...any) *EventAttr {
	return &EventAttr{Event: "change", Handler: handler, Args: args}
}

// OnEvent creates a custom event handler.
func OnEvent(event string, handler any, args ...any) *EventAttr {
	return &EventAttr{Event: event, Handler: handler, Args: args}
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

// BindValueAttr represents two-way binding to an input's value.
type BindValueAttr struct {
	Store   any    // *Store[string] or *Store[int]
	StoreID string // Go expression for code generation
}

func (b *BindValueAttr) attrType() string { return "bindvalue" }

// BindValue creates a two-way binding between an input and a store.
func BindValue[T any](store *Store[T]) *BindValueAttr {
	return &BindValueAttr{Store: store}
}

// BindCheckedAttr represents two-way binding to a checkbox's checked state.
type BindCheckedAttr struct {
	Store   any    // *Store[bool]
	StoreID string // Go expression for code generation
}

func (b *BindCheckedAttr) attrType() string { return "bindchecked" }

// BindChecked creates a two-way binding between a checkbox and a bool store.
func BindChecked(store *Store[bool]) *BindCheckedAttr {
	return &BindCheckedAttr{Store: store}
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
