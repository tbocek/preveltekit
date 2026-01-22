package main

import (
	"fmt"
	"go/ast"
	"go/parser"
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
	Expr    string
	ID      string
	VarDeps []string // which reactive vars this expression depends on
}

func (e ExprNode) nodeType() string { return "expr" }

type HtmlNode struct {
	Expr    string
	ID      string
	VarDeps []string
}

func (h HtmlNode) nodeType() string { return "html" }

type IfBranch struct {
	Cond    string // condition (empty for final else)
	Body    []Node
	VarDeps []string // vars the condition depends on
}

type IfNode struct {
	CondID   string     // "if0"
	Branches []IfBranch // first is "if", rest are "else if", last may be "else" (empty Cond)
}

func (i IfNode) nodeType() string { return "if" }

type EachNode struct {
	Array   string // "items"
	Item    string // "item"
	Index   string // "i" or ""
	ID      string // "each0"
	Body    []Node
	VarDeps []string // vars this each depends on (the array var)
}

func (e EachNode) nodeType() string { return "each" }

type ComponentNode struct {
	Name     string            // "Button"
	ID       string            // "comp0"
	Props    map[string]string // Label -> "hello" or Label -> "{varName}"
	Bindings map[string]string // bind:value -> varName (two-way binding)
	Children string            // raw children HTML (for slot)
}

func (c ComponentNode) nodeType() string { return "component" }

type SlotNode struct{}

func (s SlotNode) nodeType() string { return "slot" }

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
	Nodes        []Node
	exprCount    int
	htmlCount    int
	ifCount      int
	eachCount    int
	elemCount    int
	compCount    int
	Components   map[string]bool   // known component names
	ReactiveVars map[string]bool   // known reactive var names for dependency tracking
}

func parseTemplate(template string) *TemplateAST {
	tmplAst := &TemplateAST{
		Components:   make(map[string]bool),
		ReactiveVars: make(map[string]bool),
	}
	tmplAst.Nodes = tmplAst.parseNodes(template)
	return tmplAst
}

func parseTemplateWithComponents(template string, components map[string]bool) *TemplateAST {
	tmplAst := &TemplateAST{
		Components:   components,
		ReactiveVars: make(map[string]bool),
	}
	tmplAst.Nodes = tmplAst.parseNodes(template)
	return tmplAst
}

func parseTemplateWithReactiveVars(template string, components map[string]bool, reactiveVars map[string]bool) *TemplateAST {
	tmplAst := &TemplateAST{
		Components:   components,
		ReactiveVars: reactiveVars,
	}
	tmplAst.Nodes = tmplAst.parseNodes(template)
	return tmplAst
}

func (tmplAst *TemplateAST) parseNodes(s string) []Node {
	var nodes []Node
	i := 0

	for i < len(s) {
		// Look for { or component tag or slot
		nextBrace := strings.Index(s[i:], "{")
		nextComp := tmplAst.findNextComponent(s[i:])
		nextSlot := strings.Index(s[i:], "<slot")

		// Determine which comes first
		nextSpecial := -1
		specialType := ""

		candidates := []struct {
			idx  int
			kind string
		}{
			{nextBrace, "brace"},
			{nextComp, "component"},
			{nextSlot, "slot"},
		}

		for _, c := range candidates {
			if c.idx != -1 && (nextSpecial == -1 || c.idx < nextSpecial) {
				nextSpecial = c.idx
				specialType = c.kind
			}
		}

		if nextSpecial == -1 {
			// Rest is text
			if i < len(s) {
				nodes = append(nodes, TextNode{Text: s[i:]})
			}
			break
		}

		// Text before special
		if nextSpecial > 0 {
			nodes = append(nodes, TextNode{Text: s[i : i+nextSpecial]})
		}
		i += nextSpecial

		if specialType == "slot" {
			// Parse <slot /> or <slot>
			end := strings.Index(s[i:], ">")
			if end != -1 {
				nodes = append(nodes, SlotNode{})
				i += end + 1
				continue
			}
		}

		if specialType == "component" {
			// Parse component tag
			node, consumed := tmplAst.parseComponent(s[i:])
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
			node, consumed := tmplAst.parseIf(s[i:])
			nodes = append(nodes, node)
			i += consumed
		} else if strings.HasPrefix(s[i:], "{#each ") {
			node, consumed := tmplAst.parseEach(s[i:])
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
			id := fmt.Sprintf("html%d", tmplAst.htmlCount)
			tmplAst.htmlCount++
			nodes = append(nodes, HtmlNode{
				Expr:    expr,
				ID:      id,
				VarDeps: tmplAst.findExprDeps(expr),
			})
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
			id := fmt.Sprintf("expr%d", tmplAst.exprCount)
			tmplAst.exprCount++
			nodes = append(nodes, ExprNode{
				Expr:    expr,
				ID:      id,
				VarDeps: tmplAst.findExprDeps(expr),
			})
			i += end + 1
		}
	}

	return nodes
}

// findExprDeps extracts which reactive vars an expression depends on
func (tmplAst *TemplateAST) findExprDeps(expr string) []string {
	var deps []string
	seen := make(map[string]bool)

	// Try to parse as Go expression
	node, err := parser.ParseExpr(expr)
	if err != nil {
		// Fallback: simple identifier check
		for varName := range tmplAst.ReactiveVars {
			if strings.Contains(expr, varName) && !seen[varName] {
				deps = append(deps, varName)
				seen[varName] = true
			}
		}
		return deps
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			if tmplAst.ReactiveVars[id.Name] && !seen[id.Name] {
				deps = append(deps, id.Name)
				seen[id.Name] = true
			}
		}
		return true
	})

	return deps
}

func (tmplAst *TemplateAST) parseIf(s string) (IfNode, int) {
	// Find condition: {#if cond}
	condEnd := strings.Index(s, "}")
	cond := strings.TrimSpace(s[5:condEnd])
	id := fmt.Sprintf("if%d", tmplAst.ifCount)
	tmplAst.ifCount++

	node := IfNode{
		CondID: id,
	}

	// Parse first branch (the if)
	rest := s[condEnd+1:]
	thenNodes := tmplAst.parseNodes(rest)
	node.Branches = append(node.Branches, IfBranch{
		Cond:    cond,
		Body:    thenNodes,
		VarDeps: tmplAst.findExprDeps(cond),
	})

	consumed := condEnd + 1 + tmplAst.findConsumed(rest, thenNodes)

	// Check for else if / else chains
	for {
		remaining := s[consumed:]
		if strings.HasPrefix(remaining, "{:else if ") {
			// Parse {:else if cond}
			elseIfEnd := strings.Index(remaining, "}")
			elseIfCond := strings.TrimSpace(remaining[10:elseIfEnd])

			consumed += elseIfEnd + 1
			rest = s[consumed:]
			elseIfNodes := tmplAst.parseNodes(rest)
			node.Branches = append(node.Branches, IfBranch{
				Cond:    elseIfCond,
				Body:    elseIfNodes,
				VarDeps: tmplAst.findExprDeps(elseIfCond),
			})
			consumed += tmplAst.findConsumed(rest, elseIfNodes)
		} else if strings.HasPrefix(remaining, "{:else}") {
			// Final else
			consumed += 7
			rest = s[consumed:]
			elseNodes := tmplAst.parseNodes(rest)
			node.Branches = append(node.Branches, IfBranch{Cond: "", Body: elseNodes})
			consumed += tmplAst.findConsumed(rest, elseNodes)
			break
		} else {
			break
		}
	}

	// Consume {/if}
	remaining := s[consumed:]
	if strings.HasPrefix(remaining, "{/if}") {
		consumed += 5
	}

	return node, consumed
}

func (tmplAst *TemplateAST) parseEach(s string) (EachNode, int) {
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

	id := fmt.Sprintf("each%d", tmplAst.eachCount)
	tmplAst.eachCount++

	node := EachNode{
		Array:   array,
		Item:    item,
		Index:   index,
		ID:      id,
		VarDeps: tmplAst.findExprDeps(array),
	}

	rest := s[condEnd+1:]
	node.Body = tmplAst.parseNodes(rest)

	consumed := condEnd + 1 + tmplAst.findConsumed(rest, node.Body)

	if strings.HasPrefix(s[consumed:], "{/each}") {
		consumed += 7
	}

	return node, consumed
}

func (tmplAst *TemplateAST) findConsumed(s string, nodes []Node) int {
	total := 0
	for _, n := range nodes {
		total += tmplAst.nodeLen(s[total:], n)
	}
	return total
}

func (tmplAst *TemplateAST) nodeLen(s string, n Node) int {
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
func (tmplAst *TemplateAST) CollectExprs() []ExprNode {
	var exprs []ExprNode
	tmplAst.collectExprsFromNodes(tmplAst.Nodes, &exprs)
	return exprs
}

func (tmplAst *TemplateAST) collectExprsFromNodes(nodes []Node, exprs *[]ExprNode) {
	for _, n := range nodes {
		switch node := n.(type) {
		case ExprNode:
			*exprs = append(*exprs, node)
		case IfNode:
			for _, branch := range node.Branches {
				tmplAst.collectExprsFromNodes(branch.Body, exprs)
			}
		case EachNode:
			tmplAst.collectExprsFromNodes(node.Body, exprs)
		}
	}
}

// Collect all html nodes for DOM refs
func (tmplAst *TemplateAST) CollectHtmls() []HtmlNode {
	var htmls []HtmlNode
	tmplAst.collectHtmlsFromNodes(tmplAst.Nodes, &htmls)
	return htmls
}

func (tmplAst *TemplateAST) collectHtmlsFromNodes(nodes []Node, htmls *[]HtmlNode) {
	for _, n := range nodes {
		switch node := n.(type) {
		case HtmlNode:
			*htmls = append(*htmls, node)
		case IfNode:
			for _, branch := range node.Branches {
				tmplAst.collectHtmlsFromNodes(branch.Body, htmls)
			}
		case EachNode:
			tmplAst.collectHtmlsFromNodes(node.Body, htmls)
		}
	}
}

// Collect all if blocks
func (tmplAst *TemplateAST) CollectIfs() []IfNode {
	var ifs []IfNode
	tmplAst.collectIfsFromNodes(tmplAst.Nodes, &ifs)
	return ifs
}

func (tmplAst *TemplateAST) collectIfsFromNodes(nodes []Node, ifs *[]IfNode) {
	for _, n := range nodes {
		switch node := n.(type) {
		case IfNode:
			*ifs = append(*ifs, node)
			for _, branch := range node.Branches {
				tmplAst.collectIfsFromNodes(branch.Body, ifs)
			}
		case EachNode:
			tmplAst.collectIfsFromNodes(node.Body, ifs)
		}
	}
}

// Collect all each blocks
func (tmplAst *TemplateAST) CollectEaches() []EachNode {
	var eaches []EachNode
	tmplAst.collectEachesFromNodes(tmplAst.Nodes, &eaches)
	return eaches
}

func (tmplAst *TemplateAST) collectEachesFromNodes(nodes []Node, eaches *[]EachNode) {
	for _, n := range nodes {
		switch node := n.(type) {
		case EachNode:
			*eaches = append(*eaches, node)
			tmplAst.collectEachesFromNodes(node.Body, eaches)
		case IfNode:
			for _, branch := range node.Branches {
				tmplAst.collectEachesFromNodes(branch.Body, eaches)
			}
		}
	}
}

// Generate static HTML with placeholders
func (tmplAst *TemplateAST) GenerateHTML() string {
	var b strings.Builder
	tmplAst.generateHTMLNodes(&b, tmplAst.Nodes)
	return b.String()
}

func (tmplAst *TemplateAST) generateHTMLNodes(b *strings.Builder, nodes []Node) {
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
		case SlotNode:
			// Placeholder for slot content
			b.WriteString("<!--SLOT-->")
		}
	}
}

// Find next component tag in string, returns index or -1
func (tmplAst *TemplateAST) findNextComponent(s string) int {
	if len(tmplAst.Components) == 0 {
		return -1
	}

	minIdx := -1
	for name := range tmplAst.Components {
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
func (tmplAst *TemplateAST) parseComponent(s string) (*ComponentNode, int) {
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
	if !tmplAst.Components[name] {
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
					ID:       fmt.Sprintf("comp%d", tmplAst.compCount),
					Props:    props,
					Bindings: bindings,
				}
				tmplAst.compCount++
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
			childrenStr := s[i+1 : i+1+closeIdx]
			node := &ComponentNode{
				Name:     name,
				ID:       fmt.Sprintf("comp%d", tmplAst.compCount),
				Props:    props,
				Bindings: bindings,
				Children: childrenStr,
			}
			tmplAst.compCount++
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
			attrValue = s[valueStart : i-1] // Don't include braces for bindings
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
func (tmplAst *TemplateAST) CollectComponents() []ComponentNode {
	var comps []ComponentNode
	tmplAst.collectComponentsFromNodes(tmplAst.Nodes, &comps)
	return comps
}

func (tmplAst *TemplateAST) collectComponentsFromNodes(nodes []Node, comps *[]ComponentNode) {
	for _, n := range nodes {
		switch node := n.(type) {
		case ComponentNode:
			*comps = append(*comps, node)
		case IfNode:
			for _, branch := range node.Branches {
				tmplAst.collectComponentsFromNodes(branch.Body, comps)
			}
		case EachNode:
			tmplAst.collectComponentsFromNodes(node.Body, comps)
		}
	}
}

// GetExprsByVar returns expressions that depend on a specific var
func (tmplAst *TemplateAST) GetExprsByVar(varName string) []ExprNode {
	var result []ExprNode
	for _, expr := range tmplAst.CollectExprs() {
		for _, dep := range expr.VarDeps {
			if dep == varName {
				result = append(result, expr)
				break
			}
		}
	}
	return result
}

// GetIfsByVar returns if nodes that depend on a specific var
func (tmplAst *TemplateAST) GetIfsByVar(varName string) []IfNode {
	var result []IfNode
	for _, ifn := range tmplAst.CollectIfs() {
		for _, branch := range ifn.Branches {
			for _, dep := range branch.VarDeps {
				if dep == varName {
					result = append(result, ifn)
					break
				}
			}
		}
	}
	return result
}

// GetEachesByVar returns each nodes that depend on a specific var
func (tmplAst *TemplateAST) GetEachesByVar(varName string) []EachNode {
	var result []EachNode
	for _, each := range tmplAst.CollectEaches() {
		for _, dep := range each.VarDeps {
			if dep == varName {
				result = append(result, each)
				break
			}
		}
	}
	return result
}