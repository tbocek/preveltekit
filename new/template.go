package main

import (
	"fmt"
	"strings"
)

// Template AST nodes
type Node interface {
	nodeType() string
}

type TextNode struct {
	Text string
}

func (t TextNode) nodeType() string { return "text" }

type ExprNode struct {
	Expr string
	ID   string
}

func (e ExprNode) nodeType() string { return "expr" }

type HtmlNode struct {
	Expr string
	ID   string
}

func (h HtmlNode) nodeType() string { return "html" }

type IfNode struct {
	Cond   string
	CondID string
	Then   []Node
	Else   []Node
}

func (i IfNode) nodeType() string { return "if" }

type EachNode struct {
	Array string // "items"
	Item  string // "item"
	Index string // "i" or ""
	ID    string // "each0"
	Body  []Node
}

func (e EachNode) nodeType() string { return "each" }

type ComponentNode struct {
	Name     string            // "Button"
	ID       string            // "comp0"
	Props    map[string]string // Label -> "hello" or Label -> "{varName}"
	Bindings map[string]string // bind:value -> varName (two-way binding)
	Children string            // inner HTML content
}

func (c ComponentNode) nodeType() string { return "component" }

type ElementNode struct {
	Tag      string
	Attrs    map[string]string
	Events   map[string]string // @click -> handler
	Bindings map[string]string // bind:value -> var
	ID       string
	Children []Node
}

func (e ElementNode) nodeType() string { return "element" }

type TemplateAST struct {
	Nodes      []Node
	exprCount  int
	htmlCount  int
	ifCount    int
	eachCount  int
	elemCount  int
	compCount  int
	Components map[string]bool // known component names
}

func parseTemplate(template string) *TemplateAST {
	ast := &TemplateAST{
		Components: make(map[string]bool),
	}
	ast.Nodes = ast.parseNodes(template)
	return ast
}

func parseTemplateWithComponents(template string, components map[string]bool) *TemplateAST {
	ast := &TemplateAST{
		Components: components,
	}
	ast.Nodes = ast.parseNodes(template)
	return ast
}

func (ast *TemplateAST) parseNodes(s string) []Node {
	var nodes []Node
	i := 0

	for i < len(s) {
		// Look for { or component tag
		nextBrace := strings.Index(s[i:], "{")
		nextComp := ast.findNextComponent(s[i:])
		
		// Determine which comes first
		nextSpecial := -1
		isComponent := false
		if nextBrace == -1 && nextComp == -1 {
			// Rest is text
			if i < len(s) {
				nodes = append(nodes, TextNode{Text: s[i:]})
			}
			break
		} else if nextBrace == -1 {
			nextSpecial = nextComp
			isComponent = true
		} else if nextComp == -1 {
			nextSpecial = nextBrace
		} else if nextComp < nextBrace {
			nextSpecial = nextComp
			isComponent = true
		} else {
			nextSpecial = nextBrace
		}

		// Text before special
		if nextSpecial > 0 {
			nodes = append(nodes, TextNode{Text: s[i : i+nextSpecial]})
		}
		i += nextSpecial

		if isComponent {
			// Parse component tag
			node, consumed := ast.parseComponent(s[i:])
			if node != nil {
				nodes = append(nodes, *node)
				i += consumed
				continue
			}
			// Not a valid component, treat as text
			nodes = append(nodes, TextNode{Text: s[i : i+1]})
			i++
			continue
		}

		// What kind of brace?
		if strings.HasPrefix(s[i:], "{#if ") {
			node, consumed := ast.parseIf(s[i:])
			nodes = append(nodes, node)
			i += consumed
		} else if strings.HasPrefix(s[i:], "{#each ") {
			node, consumed := ast.parseEach(s[i:])
			nodes = append(nodes, node)
			i += consumed
		} else if strings.HasPrefix(s[i:], "{:else}") || strings.HasPrefix(s[i:], "{/if}") || strings.HasPrefix(s[i:], "{/each}") {
			// End markers - return to parent
			break
		} else if strings.HasPrefix(s[i:], "{/") || strings.HasPrefix(s[i:], "{:") {
			// Other end/else markers
			break
		} else if strings.HasPrefix(s[i:], "{@html ") {
			// Raw HTML {@html expr}
			end := strings.Index(s[i:], "}")
			if end == -1 {
				nodes = append(nodes, TextNode{Text: s[i:]})
				break
			}
			expr := strings.TrimSpace(s[i+7 : i+end])
			id := fmt.Sprintf("html%d", ast.htmlCount)
			ast.htmlCount++
			nodes = append(nodes, HtmlNode{Expr: expr, ID: id})
			i += end + 1
		} else {
			// Expression {var} - but skip if inside attribute value (={...})
			if i > 0 && s[i-1] == '=' {
				// This is an attribute value like class:active={cond}, skip it
				end := strings.Index(s[i:], "}")
				if end != -1 {
					nodes = append(nodes, TextNode{Text: s[i : i+end+1]})
					i += end + 1
				} else {
					nodes = append(nodes, TextNode{Text: s[i:]})
					break
				}
				continue
			}
			end := strings.Index(s[i:], "}")
			if end == -1 {
				nodes = append(nodes, TextNode{Text: s[i:]})
				break
			}
			expr := strings.TrimSpace(s[i+1 : i+end])
			id := fmt.Sprintf("expr%d", ast.exprCount)
			ast.exprCount++
			nodes = append(nodes, ExprNode{Expr: expr, ID: id})
			i += end + 1
		}
	}

	return nodes
}

func (ast *TemplateAST) parseIf(s string) (IfNode, int) {
	// Find condition: {#if cond}
	condEnd := strings.Index(s, "}")
	cond := strings.TrimSpace(s[5:condEnd])
	id := fmt.Sprintf("if%d", ast.ifCount)
	ast.ifCount++

	node := IfNode{
		Cond:   cond,
		CondID: id,
	}

	// Parse then branch
	rest := s[condEnd+1:]
	thenNodes := ast.parseNodes(rest)
	node.Then = thenNodes

	// Find where then ended
	consumed := condEnd + 1 + ast.findConsumed(rest, thenNodes)

	// Check for else
	remaining := s[consumed:]
	if strings.HasPrefix(remaining, "{:else}") {
		consumed += 7
		rest = s[consumed:]
		elseNodes := ast.parseNodes(rest)
		node.Else = elseNodes
		consumed += ast.findConsumed(rest, elseNodes)
	}

	// Consume {/if}
	remaining = s[consumed:]
	if strings.HasPrefix(remaining, "{/if}") {
		consumed += 5
	}

	return node, consumed
}

func (ast *TemplateAST) parseEach(s string) (EachNode, int) {
	// {#each items as item} or {#each items as item, i}
	condEnd := strings.Index(s, "}")
	inner := strings.TrimSpace(s[7:condEnd]) // after "{#each "

	parts := strings.Split(inner, " as ")
	array := strings.TrimSpace(parts[0])

	var item, index string
	if len(parts) > 1 {
		iterParts := strings.Split(parts[1], ",")
		item = strings.TrimSpace(iterParts[0])
		if len(iterParts) > 1 {
			index = strings.TrimSpace(iterParts[1])
		}
	}

	id := fmt.Sprintf("each%d", ast.eachCount)
	ast.eachCount++

	node := EachNode{
		Array: array,
		Item:  item,
		Index: index,
		ID:    id,
	}

	rest := s[condEnd+1:]
	node.Body = ast.parseNodes(rest)

	consumed := condEnd + 1 + ast.findConsumed(rest, node.Body)

	if strings.HasPrefix(s[consumed:], "{/each}") {
		consumed += 7
	}

	return node, consumed
}

func (ast *TemplateAST) findConsumed(s string, nodes []Node) int {
	total := 0
	for _, n := range nodes {
		total += ast.nodeLen(s[total:], n)
	}
	return total
}

func (ast *TemplateAST) nodeLen(s string, n Node) int {
	switch node := n.(type) {
	case TextNode:
		return len(node.Text)
	case ExprNode:
		// {expr}
		idx := strings.Index(s, "{")
		end := strings.Index(s[idx:], "}")
		return idx + end + 1
	case HtmlNode:
		// {@html expr}
		idx := strings.Index(s, "{@html")
		end := strings.Index(s[idx:], "}")
		return idx + end + 1
	case IfNode:
		// Complex - find the full if block
		start := strings.Index(s, "{#if")
		depth := 1
		i := start + 4
		for i < len(s) && depth > 0 {
			if strings.HasPrefix(s[i:], "{#if") {
				depth++
				i += 4
			} else if strings.HasPrefix(s[i:], "{/if}") {
				depth--
				if depth == 0 {
					return i + 5
				}
				i += 5
			} else {
				i++
			}
		}
		return i
	case EachNode:
		start := strings.Index(s, "{#each")
		depth := 1
		i := start + 6
		for i < len(s) && depth > 0 {
			if strings.HasPrefix(s[i:], "{#each") {
				depth++
				i += 6
			} else if strings.HasPrefix(s[i:], "{/each}") {
				depth--
				if depth == 0 {
					return i + 7
				}
				i += 7
			} else {
				i++
			}
		}
		return i
	}
	return 0
}

// Collect all expression IDs for DOM refs
func (ast *TemplateAST) CollectExprs() []ExprNode {
	var exprs []ExprNode
	ast.collectExprsFromNodes(ast.Nodes, &exprs)
	return exprs
}

func (ast *TemplateAST) collectExprsFromNodes(nodes []Node, exprs *[]ExprNode) {
	for _, n := range nodes {
		switch node := n.(type) {
		case ExprNode:
			*exprs = append(*exprs, node)
		case IfNode:
			ast.collectExprsFromNodes(node.Then, exprs)
			ast.collectExprsFromNodes(node.Else, exprs)
		case EachNode:
			ast.collectExprsFromNodes(node.Body, exprs)
		}
	}
}

// Collect all html nodes for DOM refs
func (ast *TemplateAST) CollectHtmls() []HtmlNode {
	var htmls []HtmlNode
	ast.collectHtmlsFromNodes(ast.Nodes, &htmls)
	return htmls
}

func (ast *TemplateAST) collectHtmlsFromNodes(nodes []Node, htmls *[]HtmlNode) {
	for _, n := range nodes {
		switch node := n.(type) {
		case HtmlNode:
			*htmls = append(*htmls, node)
		case IfNode:
			ast.collectHtmlsFromNodes(node.Then, htmls)
			ast.collectHtmlsFromNodes(node.Else, htmls)
		case EachNode:
			ast.collectHtmlsFromNodes(node.Body, htmls)
		}
	}
}

// Collect all if blocks
func (ast *TemplateAST) CollectIfs() []IfNode {
	var ifs []IfNode
	ast.collectIfsFromNodes(ast.Nodes, &ifs)
	return ifs
}

func (ast *TemplateAST) collectIfsFromNodes(nodes []Node, ifs *[]IfNode) {
	for _, n := range nodes {
		switch node := n.(type) {
		case IfNode:
			*ifs = append(*ifs, node)
			ast.collectIfsFromNodes(node.Then, ifs)
			ast.collectIfsFromNodes(node.Else, ifs)
		case EachNode:
			ast.collectIfsFromNodes(node.Body, ifs)
		}
	}
}

// Collect all each blocks
func (ast *TemplateAST) CollectEaches() []EachNode {
	var eaches []EachNode
	ast.collectEachesFromNodes(ast.Nodes, &eaches)
	return eaches
}

func (ast *TemplateAST) collectEachesFromNodes(nodes []Node, eaches *[]EachNode) {
	for _, n := range nodes {
		switch node := n.(type) {
		case EachNode:
			*eaches = append(*eaches, node)
			ast.collectEachesFromNodes(node.Body, eaches)
		case IfNode:
			ast.collectEachesFromNodes(node.Then, eaches)
			ast.collectEachesFromNodes(node.Else, eaches)
		}
	}
}

// Generate static HTML with placeholders
func (ast *TemplateAST) GenerateHTML() string {
	var b strings.Builder
	ast.generateHTMLNodes(&b, ast.Nodes)
	return b.String()
}

func (ast *TemplateAST) generateHTMLNodes(b *strings.Builder, nodes []Node) {
	for _, n := range nodes {
		switch node := n.(type) {
		case TextNode:
			b.WriteString(node.Text)
		case ExprNode:
			fmt.Fprintf(b, `<span id="%s"></span>`, node.ID)
		case HtmlNode:
			fmt.Fprintf(b, `<span id="%s"></span>`, node.ID)
		case IfNode:
			// Hidden span anchor for if block insertion point
			fmt.Fprintf(b, `<span id="%s_anchor" style="display:none"></span>`, node.CondID)
		case EachNode:
			// Hidden span anchor for each block insertion point
			fmt.Fprintf(b, `<span id="%s_anchor" style="display:none"></span>`, node.ID)
		case ComponentNode:
			// Placeholder for component
			fmt.Fprintf(b, `<span id="%s"></span>`, node.ID)
		}
	}
}

// Find next component tag in string, returns index or -1
func (ast *TemplateAST) findNextComponent(s string) int {
	if len(ast.Components) == 0 {
		return -1
	}
	
	minIdx := -1
	for name := range ast.Components {
		// Look for <ComponentName with space or > or /
		patterns := []string{"<" + name + " ", "<" + name + ">", "<" + name + "/"}
		for _, pattern := range patterns {
			idx := strings.Index(s, pattern)
			if idx != -1 && (minIdx == -1 || idx < minIdx) {
				minIdx = idx
			}
		}
	}
	return minIdx
}

// Parse component tag like <Button Label="hello" /> or <Button>children</Button>
func (ast *TemplateAST) parseComponent(s string) (*ComponentNode, int) {
	if s[0] != '<' {
		return nil, 0
	}
	
	// Find component name
	nameEnd := 1
	for nameEnd < len(s) && s[nameEnd] != ' ' && s[nameEnd] != '>' && s[nameEnd] != '/' {
		nameEnd++
	}
	name := s[1:nameEnd]
	
	// Verify it's a known component
	if !ast.Components[name] {
		return nil, 0
	}
	
	// Parse props
	props := make(map[string]string)
	bindings := make(map[string]string)
	i := nameEnd
	
	// Skip whitespace and parse attributes
	for i < len(s) {
		// Skip whitespace
		for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n') {
			i++
		}
		
		// Check for end of tag
		if i < len(s) && s[i] == '/' {
			// Self-closing tag
			if i+1 < len(s) && s[i+1] == '>' {
				node := &ComponentNode{
					Name:     name,
					ID:       fmt.Sprintf("comp%d", ast.compCount),
					Props:    props,
					Bindings: bindings,
				}
				ast.compCount++
				return node, i + 2
			}
		}
		if i < len(s) && s[i] == '>' {
			// Opening tag - need to find children and closing tag
			closeTag := "</" + name + ">"
			closeIdx := strings.Index(s[i+1:], closeTag)
			if closeIdx == -1 {
				return nil, 0
			}
			children := s[i+1 : i+1+closeIdx]
			node := &ComponentNode{
				Name:     name,
				ID:       fmt.Sprintf("comp%d", ast.compCount),
				Props:    props,
				Bindings: bindings,
				Children: children,
			}
			ast.compCount++
			return node, i + 1 + closeIdx + len(closeTag)
		}
		
		// Parse attribute name
		attrStart := i
		for i < len(s) && s[i] != '=' && s[i] != ' ' && s[i] != '>' && s[i] != '/' {
			i++
		}
		if i == attrStart {
			break
		}
		attrName := s[attrStart:i]
		
		// Expect =
		if i >= len(s) || s[i] != '=' {
			continue
		}
		i++
		
		// Parse value - either "string" or {expr}
		if i >= len(s) {
			break
		}
		
		var attrValue string
		if s[i] == '"' {
			// String value
			i++
			valueStart := i
			for i < len(s) && s[i] != '"' {
				i++
			}
			attrValue = s[valueStart:i]
			if i < len(s) {
				i++ // skip closing quote
			}
		} else if s[i] == '{' {
			// Expression value
			i++
			valueStart := i
			depth := 1
			for i < len(s) && depth > 0 {
				if s[i] == '{' {
					depth++
				} else if s[i] == '}' {
					depth--
				}
				i++
			}
			attrValue = s[valueStart:i-1] // Don't include braces for bindings
		} else {
			continue
		}
		
		// Check if it's a binding
		if strings.HasPrefix(attrName, "bind:") {
			bindProp := strings.TrimPrefix(attrName, "bind:")
			bindings[bindProp] = attrValue
		} else {
			// Regular prop - restore braces for expression values
			if !strings.HasPrefix(attrValue, "\"") && attrValue != "" && s[i-len(attrValue)-2] == '{' {
				attrValue = "{" + attrValue + "}"
			}
			props[attrName] = attrValue
		}
	}
	
	return nil, 0
}

// Collect all components
func (ast *TemplateAST) CollectComponents() []ComponentNode {
	var comps []ComponentNode
	ast.collectComponentsFromNodes(ast.Nodes, &comps)
	return comps
}

func (ast *TemplateAST) collectComponentsFromNodes(nodes []Node, comps *[]ComponentNode) {
	for _, n := range nodes {
		switch node := n.(type) {
		case ComponentNode:
			*comps = append(*comps, node)
		case IfNode:
			ast.collectComponentsFromNodes(node.Then, comps)
			ast.collectComponentsFromNodes(node.Else, comps)
		case EachNode:
			ast.collectComponentsFromNodes(node.Body, comps)
		}
	}
}