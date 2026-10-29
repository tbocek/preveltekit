package preveltekit

import "strings"

// Rendering a Node tree to HTML. SSR uses it to pre-render the pages, WASM to
// re-render if/each/component blocks. Both must produce the same HTML and
// advance the ID counters in the same order, so there is one implementation;
// the SSR-only parts are guarded by IsBuildTime.

// renderContext holds state while rendering a Node tree to HTML.
type renderContext struct {
	IDCounter
	ScopeAttr string           // CSS scope class of the current component (e.g. "v0"), added to every tag
	slot      func() string    // returns the HTML rendered in place of Slot() (nil: none)
	styles    *collectedStyles // SSR: CSS of the rendered components (nil in WASM)
}

// collectedStyles holds the CSS of the rendered components by component name.
type collectedStyles struct {
	scoped map[string]string // from Style(), scoped to the component
	global map[string]string // from GlobalStyle(), unscoped
}

// cachedOption is a rendered option of a Store[Component]. WASM caches the
// Render() trees of the HTML pass in renderedTrees (by marker ID) so the bind
// pass can reuse them without calling Render() again (which re-registers handlers).
type cachedOption struct {
	comp      Component
	name      string
	tree      Node
	scopeAttr string
}

var renderedTrees = make(map[string][]cachedOption)

// with returns a context in the same component (scope, slot, styles) with
// the given counters.
func (ctx *renderContext) with(c IDCounter) *renderContext {
	return &renderContext{IDCounter: c, ScopeAttr: ctx.ScopeAttr, slot: ctx.slot, styles: ctx.styles}
}

// component returns the context for rendering a component's own tree.
func (ctx *renderContext) component(prefix, scopeAttr string, slot func() string) *renderContext {
	return &renderContext{IDCounter: IDCounter{Prefix: prefix}, ScopeAttr: scopeAttr, slot: slot, styles: ctx.styles}
}

// componentScope returns the CSS scope class of a component ("" if it has no
// Style()). In SSR it also collects the component's styles for the page.
func (ctx *renderContext) componentScope(name string, comp any) string {
	if IsBuildTime && ctx.styles != nil {
		if g, ok := comp.(HasGlobalStyle); ok {
			if _, exists := ctx.styles.global[name]; !exists {
				if gs := g.GlobalStyle(); gs != "" {
					ctx.styles.global[name] = gs
				}
			}
		}
	}
	s, ok := comp.(HasStyle)
	if !ok {
		return ""
	}
	scope := GetOrCreateScope(name)
	if IsBuildTime && ctx.styles != nil {
		if _, exists := ctx.styles.scoped[name]; !exists {
			ctx.styles.scoped[name] = scopeCSS(s.Style(), scope)
		}
	}
	return scope
}

// nodeToHTML renders a node to HTML.
func nodeToHTML(n Node, ctx *renderContext) string {
	switch node := n.(type) {
	case *HtmlNode:
		return htmlNodeToHTML(node, ctx)
	case *RawHTMLNode:
		return node.HTML
	case *FragmentNode:
		return childrenToHTML(node.Children, ctx)
	case *BindNode:
		return bindNodeToHTML(node, ctx)
	case *IfNode:
		return ifNodeToHTML(node, ctx)
	case *EachNode:
		return eachNodeToHTML(node, ctx)
	case *ComponentNode:
		return componentNodeToHTML(node, ctx)
	case *SlotNode:
		if ctx.slot != nil {
			return ctx.slot()
		}
	case *TextNode:
		return escapeHTML(node.Text)
	}
	return ""
}

// childrenToHTML renders a slice of nodes to HTML.
func childrenToHTML(nodes []Node, ctx *renderContext) string {
	var sb strings.Builder
	for _, n := range nodes {
		sb.WriteString(nodeToHTML(n, ctx))
	}
	return sb.String()
}

// markedHTML wraps html in the start and end comment markers of a block.
func markedHTML(markerID, html string) string {
	return "<!--" + markerID + "s-->" + html + "<!--" + markerID + "-->"
}

// htmlNodeToHTML renders an element from its structured fields.
func htmlNodeToHTML(h *HtmlNode, ctx *renderContext) string {
	var sb strings.Builder

	// --- Advance counters (DynAttrs, element ID, bind ID) ---
	dynAttrHTML := make([]string, len(h.DynAttrs))
	for i, da := range h.DynAttrs {
		dynAttrHTML[i] = `data-attrbind="` + ctx.NextAttrID() + `" ` + da.Name + `="` + escapeHTML(da.value()) + `"`
	}
	elementID := h.elementID(&ctx.IDCounter)
	var bindID string
	if h.BoundStore != nil {
		bindID = ctx.NextBindID()
	}

	// --- Opening tag ---
	sb.WriteByte('<')
	sb.WriteString(h.Tag)

	// Collect class fragments: static classes, scope class, conditional classes
	var classFragments []string
	for _, attr := range h.Attrs {
		if strings.HasPrefix(attr, `class="`) {
			classFragments = append(classFragments, attr[7:len(attr)-1])
		}
	}
	if ctx.ScopeAttr != "" {
		classFragments = append(classFragments, ctx.ScopeAttr)
	}
	for _, ac := range h.AttrConds {
		if ac.Name == "class" {
			if val := ac.value(); val != "" {
				classFragments = append(classFragments, val)
			}
		}
	}

	if elementID != "" {
		writeAttr(&sb, "id", elementID)
	} else if bindID != "" {
		writeAttr(&sb, "id", bindID)
	}
	if len(classFragments) > 0 {
		writeAttr(&sb, "class", strings.Join(classFragments, " "))
	}

	// Non-class static attrs (already formatted as name="value")
	for _, attr := range h.Attrs {
		if !strings.HasPrefix(attr, `class="`) {
			sb.WriteByte(' ')
			sb.WriteString(attr)
		}
	}

	// Non-class conditional attrs
	for _, ac := range h.AttrConds {
		if ac.Name != "class" {
			if val := ac.value(); val != "" {
				writeAttr(&sb, ac.Name, escapeHTML(val))
			}
		}
	}

	// Bound value/checked
	if bindID != "" {
		switch s := h.BoundStore.(type) {
		case *Store[string]:
			if h.Tag != "textarea" {
				writeAttr(&sb, "value", escapeHTML(s.Get()))
			}
		case *Store[int]:
			writeAttr(&sb, "value", itoa(s.Get()))
		case *Store[bool]:
			if s.Get() {
				sb.WriteString(` checked`)
			}
		}
	}

	for _, das := range dynAttrHTML {
		sb.WriteByte(' ')
		sb.WriteString(das)
	}
	sb.WriteByte('>')

	// --- Children ---
	if h.Tag == "textarea" && bindID != "" {
		if s, ok := h.BoundStore.(*Store[string]); ok {
			sb.WriteString(escapeHTML(s.Get()))
		}
	} else {
		for _, child := range h.Children {
			switch v := child.(type) {
			case Node:
				sb.WriteString(nodeToHTML(v, ctx))
			case *Store[Component]:
				sb.WriteString(storeComponentToHTML(v, ctx))
			case AnyGetter:
				sb.WriteString(bindNodeToHTML(&BindNode{StoreRef: v}, ctx))
			default:
				sb.WriteString(escapeHTML(anyToString(v)))
			}
		}
	}

	// --- Closing tag ---
	if !h.IsVoid {
		sb.WriteString("</")
		sb.WriteString(h.Tag)
		sb.WriteByte('>')
	}
	return sb.String()
}

// writeAttr writes ` name="value"`; value must already be escaped.
func writeAttr(sb *strings.Builder, name, value string) {
	sb.WriteByte(' ')
	sb.WriteString(name)
	sb.WriteString(`="`)
	sb.WriteString(value)
	sb.WriteByte('"')
}

// storeComponentToHTML renders a Store[Component] as a block with all its
// options rendered (to advance the counters), showing the current one.
func storeComponentToHTML(v *Store[Component], ctx *renderContext) string {
	markerID := ctx.NextRouteMarker()
	active, cached := renderStoreOptions(v, ctx)
	if !IsBuildTime {
		renderedTrees[markerID] = cached
	}
	return markedHTML(markerID, active)
}

// renderStoreOptions renders every option of a Store[Component] (or only its
// current component if it has no options), deduplicated by component name.
// It returns the HTML of the current component and the rendered options.
func renderStoreOptions(v *Store[Component], ctx *renderContext) (active string, rendered []cachedOption) {
	comp := v.Get()
	options := v.Options()
	if len(options) == 0 && comp != nil {
		options = []any{comp}
	}
	seen := make(map[string]bool)
	for _, opt := range options {
		optComp, ok := opt.(Component)
		if !ok || optComp == nil {
			continue
		}
		name := componentName(optComp)
		if seen[name] {
			continue
		}
		seen[name] = true

		scopeAttr := ctx.componentScope(name, optComp)
		// In SSR, only the displayed component is mounted
		if IsBuildTime && optComp == comp {
			if om, ok := optComp.(HasOnMount); ok {
				om.OnMount()
			}
		}
		tree := optComp.Render()
		html := nodeToHTML(tree, ctx.component(ctx.FullID(name), scopeAttr, nil))
		if optComp == comp {
			active = html
		}
		rendered = append(rendered, cachedOption{comp: optComp, name: name, tree: tree, scopeAttr: scopeAttr})
	}
	return active, rendered
}

// bindNodeToHTML renders a text binding.
func bindNodeToHTML(b *BindNode, ctx *renderContext) string {
	var value string
	if g, ok := b.StoreRef.(AnyGetter); ok {
		value = anyToString(g.GetAny())
	}
	if !b.IsHTML {
		value = escapeHTML(value)
	}
	return markedHTML(ctx.NextTextMarker(), value)
}

// ifNodeToHTML renders an if-block. All branches are rendered to advance
// the counters; the first branch whose condition holds is shown.
func ifNodeToHTML(i *IfNode, ctx *renderContext) string {
	markerID := ctx.NextIfMarker()

	var activeHTML string
	activeFound := false
	for _, branch := range i.Branches {
		branchCtx := ctx.with(ctx.IDCounter)
		branchHTML := childrenToHTML(branch.Children, branchCtx)
		ctx.IDCounter = branchCtx.IDCounter
		if !activeFound && branch.Cond.Eval() {
			activeHTML = branchHTML
			activeFound = true
		}
	}
	if len(i.ElseNode) > 0 {
		elseCtx := ctx.with(ctx.IDCounter)
		elseHTML := childrenToHTML(i.ElseNode, elseCtx)
		ctx.IDCounter = elseCtx.IDCounter
		if !activeFound {
			activeHTML = elseHTML
		}
	}
	return markedHTML(markerID, activeHTML)
}

// eachNodeToHTML renders an each-block.
func eachNodeToHTML(e *EachNode, ctx *renderContext) string {
	markerID := ctx.NextEachMarker()
	return markedHTML(markerID, eachItemsHTML(e, markerID, ctx))
}

// eachItemsHTML renders the current items (or the else content) of an each-block.
func eachItemsHTML(e *EachNode, markerID string, ctx *renderContext) string {
	var sb strings.Builder
	e.forEach(markerID, func(nodes []Node, c IDCounter) {
		sb.WriteString(childrenToHTML(nodes, ctx.with(c)))
	})
	return sb.String()
}

// forEach calls fn with the nodes of every current item (or the else content
// if the list is empty) and the counters to render or bind them with. Each
// item has its own ID prefix (e.g. "e0_1"), the else content uses the marker
// ID, so the IDs inside and after the block don't depend on the items.
func (e *EachNode) forEach(markerID string, fn func(nodes []Node, c IDCounter)) {
	items := e.items()
	if len(items) == 0 {
		fn(e.ElseNode, IDCounter{Prefix: markerID})
		return
	}
	for i, item := range items {
		fn([]Node{e.Body(item, i)}, IDCounter{Prefix: markerID + "_" + itoa(i)})
	}
}

// componentNodeToHTML renders a nested component with its slot content.
func componentNodeToHTML(c *ComponentNode, ctx *renderContext) string {
	comp, ok := c.Instance.(Component)
	if !ok {
		return ""
	}
	prefix := ctx.NextCompMarker()
	scopeAttr := ctx.componentScope(c.Name, c.Instance)
	if IsBuildTime {
		if om, ok := comp.(HasOnMount); ok {
			om.OnMount()
		}
	}
	slot := childrenToHTML(c.Children, ctx)

	// Cache the tree so the WASM bind pass doesn't call Render() again
	c.renderCache = comp.Render()
	return nodeToHTML(c.renderCache, ctx.component(prefix, scopeAttr, func() string { return slot }))
}
