package main

import "regexp"

type component struct {
	name       string
	template   string
	style      string
	source     string
	fields     []storeField
	methods    []string
	hasOnMount bool
	fieldTypes map[string]string // cached field name -> type mapping
}

type storeField struct {
	name      string
	storeType string // "Store", "List", "Map"
	valueType string // "int", "string", etc. For Map, this is the value type
	keyType   string // For Map only, the key type
}

type templateBindings struct {
	expressions   []exprBinding
	events        []eventBinding
	eachBlocks    []eachBinding
	ifBlocks      []ifBinding
	bindings      []inputBinding
	classBindings []classBinding
	attrBindings  []attrBinding
	components    []componentBinding
}

type componentBinding struct {
	name      string                    // "Button"
	elementID string                    // "comp0"
	props     map[string]string         // label -> "Click me" (static) or label -> "{Count}" (dynamic)
	events    map[string]componentEvent // click -> {method: "Add", args: "5"}
	children  string                    // content between opening and closing tags (for slot)
	insideIf  string                    // if block elementID this component is inside (empty if not in an if block)
}

type componentEvent struct {
	method string // "Add"
	args   string // "5" or "Counter" or ""
}

type exprBinding struct {
	fieldName string
	elementID string
	isHTML    bool   // true for {@html Field}, false for {Field}
	owner     string // "parent", "child", or "" (unset, will be categorized later)
}

type eventBinding struct {
	event      string
	modifiers  []string // ["preventDefault", "stopPropagation"]
	methodName string
	args       string // optional arguments like "Counter"
	elementID  string
}

type eachBinding struct {
	listName  string // "Items"
	itemVar   string // "item"
	indexVar  string // "i"
	elementID string // "each0"
	bodyHTML  string // the template inside the each block
	elseHTML  string // content to show when list is empty (optional)
}

type ifBinding struct {
	branches  []ifBranch // list of condition/content pairs
	elseHTML  string     // final else content (optional)
	elementID string     // "if0"
	deps      []string   // all fields this if block depends on
}

type ifBranch struct {
	condition  string        // "Score >= 90"
	html       string        // content if true
	eachBlocks []eachBinding // each blocks inside this branch
}

type inputBinding struct {
	fieldName string // "Name"
	bindType  string // "value", "checked"
	elementID string // "bind0"
}

type classBinding struct {
	className string // "active"
	condition string // "IsActive"
	elementID string // "class0"
}

type attrBinding struct {
	attrName  string   // "class"
	template  string   // "btn {Variant}" - original template with placeholders
	fields    []string // ["Variant"] - fields used in the template
	elementID string   // "attr0"
}

// Template parsing regexes - compiled once at package init
var (
	// Tag attribute patterns (used in html_parser.go)
	eventRegex        = regexp.MustCompile(`@(\w+)((?:\.\w+)*)="(\w+)\(([^)]*)\)"`)
	bindRegex         = regexp.MustCompile(`bind:(\w+)="(\w+)"`)
	classBindRegex    = regexp.MustCompile(`class:(\w+)=\{(\w+)\}`)
	attrWithExprRegex = regexp.MustCompile(`(\w+)="([^"]*\{[^}]+\}[^"]*)"`)
	propRegex         = regexp.MustCompile(`(\w+)="([^"]*)"`)
)
