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
}

// =============================================================================
// Text Node (auto-escaped text content)
// =============================================================================

// TextNode represents plain text content that is HTML-escaped on render.
// Used by typed element functions (Div, Span, etc.) to safely handle string children.
type TextNode struct {
	Text string
}

func (t *TextNode) nodeType() string { return "text" }

// =============================================================================
// HTML Node
// =============================================================================

// HtmlNode represents a typed HTML element.
// Supports chainable AttrIf, On, and Bind methods.
type HtmlNode struct {
	Tag      string         // e.g. "div", "span"
	Attrs    []string       // Pre-formatted static attrs: `class="foo"`, `type="text"`
	DynAttrs []*DynAttrAttr // Dynamic attributes needing render-time IDs
	Children []any          // Child content: Node, *Store[Component], AnyGetter, etc.
	IsVoid   bool           // true for br, hr, img, input, etc.

	// Chainable bindings
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

// elementID returns the id attribute needed to find the element during
// hydration: the first event's handler ID, a new class ID for conditional
// attributes, or "" if the element needs none.
func (h *HtmlNode) elementID(c *IDCounter) string {
	if len(h.Events) > 0 {
		return h.Events[0].ID
	}
	if len(h.AttrConds) > 0 {
		return c.NextClassID()
	}
	return ""
}

// value returns the attribute value for the current condition.
func (ac *AttrCond) value() string {
	if ac.Cond.Eval() {
		return evalAttrValue(ac.TrueValue)
	}
	return evalAttrValue(ac.FalseValue)
}

// =============================================================================
// RawHTML Node (unescaped inline HTML text)
// =============================================================================

// RawHTMLNode represents unescaped HTML text content.
// Unlike TextNode (which escapes), this renders the string as-is.
// Used for HTML entities (&#9889;) and inline markup (<code>foo</code>).
type RawHTMLNode struct {
	HTML string
}

func (r *RawHTMLNode) nodeType() string { return "rawhtml" }

// RawHTML creates an unescaped HTML text node.
// Example: RawHTML("&#9889;"), RawHTML("<code>foo</code>")
func RawHTML(html string) *RawHTMLNode {
	return &RawHTMLNode{HTML: html}
}

// =============================================================================
// Fragment Node (groups multiple children without a parent element)
// =============================================================================

// FragmentNode groups multiple child nodes without adding a parent element.
// Used where a single Node return is needed but multiple siblings are desired.
type FragmentNode struct {
	Children []Node
}

func (f *FragmentNode) nodeType() string { return "fragment" }

// Fragment creates a node that groups children without a wrapping element.
// Example: Fragment(Section(...), Section(...), Footer(...))
func Fragment(children ...Node) *FragmentNode {
	return &FragmentNode{Children: children}
}

// =============================================================================
// Typed Element Functions
// =============================================================================

// isVoidElement returns true for HTML5 void elements (self-closing, no children).
func isVoidElement(tag string) bool {
	switch tag {
	case "area", "base", "br", "col", "embed", "hr", "img", "input",
		"link", "meta", "source", "track", "wbr":
		return true
	}
	return false
}

// element is the internal helper for typed element functions.
// It populates structured fields (Tag, Attrs, DynAttrs, Children)
// instead of pre-serializing HTML into Parts.
func element(tag string, parts ...any) *HtmlNode {
	h := &HtmlNode{Tag: tag, IsVoid: isVoidElement(tag)}
	for _, part := range parts {
		switch v := part.(type) {
		case *StaticAttr:
			h.Attrs = append(h.Attrs, v.Name+`="`+escapeHTML(v.Value)+`"`)
		case *DynAttrAttr:
			h.DynAttrs = append(h.DynAttrs, v)
		case string:
			h.Children = append(h.Children, &TextNode{Text: v})
		default:
			h.Children = append(h.Children, part)
		}
	}
	return h
}

// El creates a typed element node for an arbitrary HTML tag.
// String children are auto-escaped. NodeAttr values become attributes.
// Example: El("custom-element", Attr("class", "foo"), "content")
func El(tag string, parts ...any) *HtmlNode {
	return element(tag, parts...)
}

// --- Sectioning / Structural ---

func Div(parts ...any) *HtmlNode     { return element("div", parts...) }
func Main(parts ...any) *HtmlNode    { return element("main", parts...) }
func Header(parts ...any) *HtmlNode  { return element("header", parts...) }
func Footer(parts ...any) *HtmlNode  { return element("footer", parts...) }
func Nav(parts ...any) *HtmlNode     { return element("nav", parts...) }
func Section(parts ...any) *HtmlNode { return element("section", parts...) }
func Article(parts ...any) *HtmlNode { return element("article", parts...) }
func Aside(parts ...any) *HtmlNode   { return element("aside", parts...) }

// --- Headings ---

func H1(parts ...any) *HtmlNode { return element("h1", parts...) }
func H2(parts ...any) *HtmlNode { return element("h2", parts...) }
func H3(parts ...any) *HtmlNode { return element("h3", parts...) }
func H4(parts ...any) *HtmlNode { return element("h4", parts...) }
func H5(parts ...any) *HtmlNode { return element("h5", parts...) }
func H6(parts ...any) *HtmlNode { return element("h6", parts...) }

// --- Text ---

func P(parts ...any) *HtmlNode          { return element("p", parts...) }
func Span(parts ...any) *HtmlNode       { return element("span", parts...) }
func Strong(parts ...any) *HtmlNode     { return element("strong", parts...) }
func Em(parts ...any) *HtmlNode         { return element("em", parts...) }
func B(parts ...any) *HtmlNode          { return element("b", parts...) }
func I(parts ...any) *HtmlNode          { return element("i", parts...) }
func U(parts ...any) *HtmlNode          { return element("u", parts...) }
func S(parts ...any) *HtmlNode          { return element("s", parts...) }
func Small(parts ...any) *HtmlNode      { return element("small", parts...) }
func Mark(parts ...any) *HtmlNode       { return element("mark", parts...) }
func Sub(parts ...any) *HtmlNode        { return element("sub", parts...) }
func Sup(parts ...any) *HtmlNode        { return element("sup", parts...) }
func Abbr(parts ...any) *HtmlNode       { return element("abbr", parts...) }
func Cite(parts ...any) *HtmlNode       { return element("cite", parts...) }
func Q(parts ...any) *HtmlNode          { return element("q", parts...) }
func Blockquote(parts ...any) *HtmlNode { return element("blockquote", parts...) }
func Code(parts ...any) *HtmlNode       { return element("code", parts...) }
func Pre(parts ...any) *HtmlNode        { return element("pre", parts...) }
func Kbd(parts ...any) *HtmlNode        { return element("kbd", parts...) }
func Samp(parts ...any) *HtmlNode       { return element("samp", parts...) }
func Var(parts ...any) *HtmlNode        { return element("var", parts...) }

// --- Lists ---

func Ul(parts ...any) *HtmlNode { return element("ul", parts...) }
func Ol(parts ...any) *HtmlNode { return element("ol", parts...) }
func Li(parts ...any) *HtmlNode { return element("li", parts...) }
func Dl(parts ...any) *HtmlNode { return element("dl", parts...) }
func Dt(parts ...any) *HtmlNode { return element("dt", parts...) }
func Dd(parts ...any) *HtmlNode { return element("dd", parts...) }

// --- Links / Media ---

func A(parts ...any) *HtmlNode          { return element("a", parts...) }
func Img(parts ...any) *HtmlNode        { return element("img", parts...) }
func Audio(parts ...any) *HtmlNode      { return element("audio", parts...) }
func Video(parts ...any) *HtmlNode      { return element("video", parts...) }
func Source(parts ...any) *HtmlNode     { return element("source", parts...) }
func Picture(parts ...any) *HtmlNode    { return element("picture", parts...) }
func Figure(parts ...any) *HtmlNode     { return element("figure", parts...) }
func Figcaption(parts ...any) *HtmlNode { return element("figcaption", parts...) }
func Canvas(parts ...any) *HtmlNode     { return element("canvas", parts...) }
func Svg(parts ...any) *HtmlNode        { return element("svg", parts...) }

// --- Tables ---

func Table(parts ...any) *HtmlNode    { return element("table", parts...) }
func Thead(parts ...any) *HtmlNode    { return element("thead", parts...) }
func Tbody(parts ...any) *HtmlNode    { return element("tbody", parts...) }
func Tfoot(parts ...any) *HtmlNode    { return element("tfoot", parts...) }
func Tr(parts ...any) *HtmlNode       { return element("tr", parts...) }
func Th(parts ...any) *HtmlNode       { return element("th", parts...) }
func Td(parts ...any) *HtmlNode       { return element("td", parts...) }
func Caption(parts ...any) *HtmlNode  { return element("caption", parts...) }
func Colgroup(parts ...any) *HtmlNode { return element("colgroup", parts...) }
func Col(parts ...any) *HtmlNode      { return element("col", parts...) }

// --- Forms ---

func Form(parts ...any) *HtmlNode     { return element("form", parts...) }
func Input(parts ...any) *HtmlNode    { return element("input", parts...) }
func Textarea(parts ...any) *HtmlNode { return element("textarea", parts...) }
func Select(parts ...any) *HtmlNode   { return element("select", parts...) }
func Option(parts ...any) *HtmlNode   { return element("option", parts...) }
func Optgroup(parts ...any) *HtmlNode { return element("optgroup", parts...) }
func Button(parts ...any) *HtmlNode   { return element("button", parts...) }
func Label(parts ...any) *HtmlNode    { return element("label", parts...) }
func Fieldset(parts ...any) *HtmlNode { return element("fieldset", parts...) }
func Legend(parts ...any) *HtmlNode   { return element("legend", parts...) }
func Datalist(parts ...any) *HtmlNode { return element("datalist", parts...) }
func Output(parts ...any) *HtmlNode   { return element("output", parts...) }
func Progress(parts ...any) *HtmlNode { return element("progress", parts...) }
func Meter(parts ...any) *HtmlNode    { return element("meter", parts...) }

// --- Semantic ---

func Details(parts ...any) *HtmlNode { return element("details", parts...) }
func Summary(parts ...any) *HtmlNode { return element("summary", parts...) }
func Dialog(parts ...any) *HtmlNode  { return element("dialog", parts...) }
func Time(parts ...any) *HtmlNode    { return element("time", parts...) }
func Data(parts ...any) *HtmlNode    { return element("data", parts...) }
func Address(parts ...any) *HtmlNode { return element("address", parts...) }
func Ruby(parts ...any) *HtmlNode    { return element("ruby", parts...) }
func Rt(parts ...any) *HtmlNode      { return element("rt", parts...) }
func Rp(parts ...any) *HtmlNode      { return element("rp", parts...) }

// --- Embedded ---

func Iframe(parts ...any) *HtmlNode { return element("iframe", parts...) }
func Embed(parts ...any) *HtmlNode  { return element("embed", parts...) }
func Object(parts ...any) *HtmlNode { return element("object", parts...) }

// --- Misc ---

func Br(parts ...any) *HtmlNode       { return element("br", parts...) }
func Hr(parts ...any) *HtmlNode       { return element("hr", parts...) }
func Wbr(parts ...any) *HtmlNode      { return element("wbr", parts...) }
func Area(parts ...any) *HtmlNode     { return element("area", parts...) }
func Base(parts ...any) *HtmlNode     { return element("base", parts...) }
func Link(parts ...any) *HtmlNode     { return element("link", parts...) }
func Meta(parts ...any) *HtmlNode     { return element("meta", parts...) }
func Script(parts ...any) *HtmlNode   { return element("script", parts...) }
func Noscript(parts ...any) *HtmlNode { return element("noscript", parts...) }
func Template(parts ...any) *HtmlNode { return element("template", parts...) }
func Map(parts ...any) *HtmlNode      { return element("map", parts...) }
func Del(parts ...any) *HtmlNode      { return element("del", parts...) }
func Ins(parts ...any) *HtmlNode      { return element("ins", parts...) }

// AttrIf adds a conditional attribute to the first HTML tag.
// Values can be string literals or *Store[T] for reactive values.
// Multiple AttrIf calls for the same attribute name merge additively.
//
// Examples:
//
//	Button("Save").AttrIf("class", cond, "active")             // adds "active" when true
//	Button("Save").AttrIf("class", cond, "active", "inactive") // "active" when true, "inactive" when false
//	A("Link").AttrIf("href", cond, urlStore, "/fallback")      // reactive value with fallback
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
//	Button("Click").On("click", handler)
//	Form(...).On("submit", handler).PreventDefault()
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
//	Input(Attr("type", "text")).Bind(nameStore)
//	Input(Attr("type", "checkbox")).Bind(darkModeStore)
func (h *HtmlNode) Bind(store any) *HtmlNode {
	h.BoundStore = store
	return h
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

// bodies returns the children of every branch, followed by the else
// children if there are any.
func (i *IfNode) bodies() [][]Node {
	bodies := make([][]Node, 0, len(i.Branches)+1)
	for _, branch := range i.Branches {
		bodies = append(bodies, branch.Children)
	}
	if len(i.ElseNode) > 0 {
		bodies = append(bodies, i.ElseNode)
	}
	return bodies
}

// active returns the index into bodies() of the branch to show: the first
// branch whose condition holds, else the else branch, or -1 for none.
func (i *IfNode) active() int {
	for idx, branch := range i.Branches {
		if branch.Cond.Eval() {
			return idx
		}
	}
	if len(i.ElseNode) > 0 {
		return len(i.Branches)
	}
	return -1
}

// =============================================================================
// Each Node (list rendering)
// =============================================================================

// EachNode represents list iteration with optional else for empty list.
type EachNode struct {
	Body     func(item any, index int) Node // Template function for each item
	ElseNode []Node                         // Content for empty list

	// Type-erased access to the List[T], set by Each
	items    func() []any // current items
	onChange func(func()) // subscribe to list changes
}

func (e *EachNode) nodeType() string { return "each" }

// Each creates a list rendering node.
// The body function receives each item and its index.
func Each[T comparable](list *List[T], body func(item T, index int) Node) *EachNode {
	return &EachNode{
		Body: func(item any, index int) Node {
			return body(item.(T), index)
		},
		items: func() []any {
			items := make([]any, len(list.items))
			for i, item := range list.items {
				items[i] = item
			}
			return items
		},
		onChange: func(fn func()) {
			list.OnChange(func([]T) { fn() })
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
	Name        string // Component type name (derived from instance)
	Instance    any    // The actual component instance
	Children    []Node // Slot content
	renderCache Node   // cached Render() result (used by WASM to avoid double Render)
}

func (c *ComponentNode) nodeType() string { return "component" }

// Comp creates a nested component node from a component instance.
// The component name is derived from the type via reflection.
// Example: Comp(&Badge{Label: p.New("New")})
func Comp(instance any, content ...any) *ComponentNode {
	c := &ComponentNode{
		Name:     typeName(instance),
		Instance: instance,
	}
	for _, item := range content {
		if v, ok := item.(Node); ok {
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

// typeName returns the name of v's type, dereferencing pointers.
func typeName(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

// componentName returns the type name of a component with a lowercase first letter.
func componentName(c Component) string {
	if c == nil {
		return ""
	}
	return lowerFirst(typeName(c))
}

// lowerFirst returns s with an ASCII uppercase first letter lowercased.
func lowerFirst(s string) string {
	if s == "" || s[0] < 'A' || s[0] > 'Z' {
		return s
	}
	return string(s[0]+'a'-'A') + s[1:]
}

// =============================================================================
// Attributes
// =============================================================================

// StaticAttr represents a static attribute (no reactive binding).
type StaticAttr struct {
	Name  string
	Value string
}

func (s *StaticAttr) attrType() string { return "static" }

// DynAttrAttr represents a dynamic attribute with mixed string/store parts.
type DynAttrAttr struct {
	Name  string
	Parts []any // mix of string and *Store[T]
}

func (d *DynAttrAttr) attrType() string { return "dynattr" }

// value returns the current attribute value: the string parts and the
// current store values concatenated.
func (d *DynAttrAttr) value() string {
	var value string
	for _, part := range d.Parts {
		switch v := part.(type) {
		case string:
			value += v
		case AnyGetter:
			value += anyToString(v.GetAny())
		}
	}
	return value
}

// Attr creates an attribute. Parts can be strings and stores.
// With only string parts, it's static (no reactive binding).
// With any store parts, it becomes a dynamic attribute that updates reactively.
func Attr(name string, parts ...any) NodeAttr {
	var value string
	for _, p := range parts {
		s, ok := p.(string)
		if !ok {
			return &DynAttrAttr{Name: name, Parts: parts}
		}
		value += s
	}
	return &StaticAttr{Name: name, Value: value}
}

// =============================================================================
// Condition Helpers
// =============================================================================

// FuncCondition wraps a function and its dependent stores for use in If/AttrIf.
type FuncCondition struct {
	fn     func() bool
	Stores []any
}

func (c *FuncCondition) Eval() bool { return c.fn() }

// Cond creates a condition from a function and its dependent stores.
func Cond(fn func() bool, stores ...any) *FuncCondition {
	return &FuncCondition{fn: fn, Stores: stores}
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

// Itoa converts an int to a string without importing strconv.
func Itoa(n int) string { return itoa(n) }

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
	intPart := int(f)
	fracPart := f - float64(intPart)
	result := itoa(intPart)
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

// escapeHTML escapes HTML special characters, in text and attribute values.
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

// evalAttrValue extracts the current string value from an attribute value.
// Handles plain strings and typed stores.
func evalAttrValue(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case AnyGetter:
		return anyToString(val.GetAny())
	}
	return ""
}
