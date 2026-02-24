//go:build wasm

package preveltekit

// WASM-side nodeToHTML: converts a Node tree to an HTML string.
// Simplified version of SSR's nodeToHTML — no binding collection, no BuildContext.
// Used to re-render if-blocks, each-blocks, and component-blocks at runtime.
//
// This renderer must produce HTML compatible with SSR output so that
// event IDs, bind IDs, and marker IDs match.

// WASMRenderContext holds state during WASM-side HTML rendering.
// Mirrors the IDCounter from SSR's BuildContext so markers stay in sync.
type WASMRenderContext struct {
	IDCounter
	ScopeAttr   string
	SlotContent string
}

// wasmRenderedTrees caches Render() trees from the HTML pass so the bind pass
// can reuse them without calling Render() again (which re-registers handlers).
// Keyed by marker ID → slice of {comp, name, tree, scopeAttr}.
type wasmCachedOption struct {
	comp      Component
	name      string
	tree      Node
	scopeAttr string
}

var wasmRenderedTrees = make(map[string][]wasmCachedOption)

// wasmNodeToHTML dispatches to the appropriate rendering function.
func wasmNodeToHTML(n Node, ctx *WASMRenderContext) string {
	if n == nil {
		return ""
	}
	switch node := n.(type) {
	case *HtmlNode:
		return wasmHtmlNodeToHTML(node, ctx)
	case *RawHTMLNode:
		return node.HTML
	case *FragmentNode:
		var s string
		for _, child := range node.Children {
			s += wasmNodeToHTML(child, ctx)
		}
		return s
	case *BindNode:
		return wasmBindNodeToHTML(node, ctx)
	case *IfNode:
		return wasmIfNodeToHTML(node, ctx)
	case *EachNode:
		return wasmEachNodeToHTML(node, ctx)
	case *ComponentNode:
		return wasmComponentNodeToHTML(node, ctx)
	case *SlotNode:
		if ctx.SlotContent != "" {
			return ctx.SlotContent
		}
		return ""
	case *TextNode:
		return escapeHTML(node.Text)
	default:
		return ""
	}
}

// wasmChildrenToHTML renders a slice of nodes to HTML.
func wasmChildrenToHTML(nodes []Node, ctx *WASMRenderContext) string {
	var s string
	for _, n := range nodes {
		s += wasmNodeToHTML(n, ctx)
	}
	return s
}

// wasmHtmlNodeToHTML renders an HtmlNode to HTML via structured rendering.
func wasmHtmlNodeToHTML(h *HtmlNode, ctx *WASMRenderContext) string {
	return wasmRenderStructured(h, ctx)
}

// wasmRenderStructured renders a typed element directly from structured fields.
func wasmRenderStructured(h *HtmlNode, ctx *WASMRenderContext) string {
	// --- Advance counters (DynAttrs first, matching wasmRenderParts order) ---
	dynAttrHTML := make([]string, len(h.DynAttrs))
	for i, da := range h.DynAttrs {
		dynAttrHTML[i] = wasmAttrToHTML(da, ctx)
	}

	// Element ID from events or AttrConds
	var elementID string
	if len(h.AttrConds) > 0 || len(h.Events) > 0 {
		if len(h.Events) > 0 {
			elementID = h.Events[0].ID
		} else {
			localID := ctx.NextClassID()
			elementID = ctx.FullID(localID)
		}
	}

	// Bind ID
	var bindID string
	if h.BoundStore != nil {
		localID := ctx.NextBindID()
		bindID = ctx.FullID(localID)
	}

	// --- Build opening tag ---
	var s string
	s = "<" + h.Tag

	// Collect class fragments
	var classFragments []string
	for _, attr := range h.Attrs {
		if len(attr) > 7 && attr[:7] == `class="` {
			classFragments = append(classFragments, attr[7:len(attr)-1])
		}
	}
	if ctx.ScopeAttr != "" {
		classFragments = append(classFragments, ctx.ScopeAttr)
	}
	for _, ac := range h.AttrConds {
		if ac.Name != "class" {
			continue
		}
		var val string
		if ac.Cond.Eval() {
			val = evalAttrValue(ac.TrueValue)
		} else {
			val = evalAttrValue(ac.FalseValue)
		}
		if val != "" {
			classFragments = append(classFragments, val)
		}
	}

	// Write id attr
	if elementID != "" {
		s += ` id="` + elementID + `"`
	} else if bindID != "" {
		s += ` id="` + bindID + `"`
	}

	// Write merged class attr
	if len(classFragments) > 0 {
		s += ` class="`
		for i, c := range classFragments {
			if i > 0 {
				s += " "
			}
			s += c
		}
		s += `"`
	}

	// Write non-class static attrs
	for _, attr := range h.Attrs {
		if len(attr) > 7 && attr[:7] == `class="` {
			continue
		}
		s += " " + attr
	}

	// Write non-class AttrCond attrs
	for _, ac := range h.AttrConds {
		if ac.Name == "class" {
			continue
		}
		var val string
		if ac.Cond.Eval() {
			val = evalAttrValue(ac.TrueValue)
		} else {
			val = evalAttrValue(ac.FalseValue)
		}
		if val != "" {
			s += " " + ac.Name + `="` + escapeAttr(val) + `"`
		}
	}

	// Write bind value/checked
	if bindID != "" {
		switch st := h.BoundStore.(type) {
		case *Store[string]:
			if h.Tag != "textarea" {
				s += ` value="` + escapeAttr(st.Get()) + `"`
			}
		case *Store[int]:
			s += ` value="` + itoa(st.Get()) + `"`
		case *Store[bool]:
			if st.Get() {
				s += ` checked`
			}
		}
	}

	// Write event data-on attr
	if len(h.Events) > 0 {
		s += ` data-on="`
		for i, ev := range h.Events {
			if i > 0 {
				s += ","
			}
			s += ev.Event
		}
		s += `"`
	}

	// Write dynamic attrs
	for _, das := range dynAttrHTML {
		s += " " + das
	}

	s += ">"

	// --- Children ---
	if h.Tag == "textarea" && bindID != "" {
		if st, ok := h.BoundStore.(*Store[string]); ok {
			s += escapeHTML(st.Get())
		}
	} else {
		for _, child := range h.Children {
			switch v := child.(type) {
			case Node:
				s += wasmNodeToHTML(v, ctx)
			case *Store[Component]:
				s += wasmStoreComponentToHTML(v, ctx)
			case AnyGetter:
				bind := &BindNode{StoreRef: v, IsHTML: false}
				s += wasmBindNodeToHTML(bind, ctx)
			default:
				s += escapeHTML(anyToString(v))
			}
		}
	}

	// --- Closing tag ---
	if !h.IsVoid {
		s += "</" + h.Tag + ">"
	}

	return s
}

// wasmStoreComponentToHTML renders a Store[Component] part.
func wasmStoreComponentToHTML(v *Store[Component], ctx *WASMRenderContext) string {
	comp := v.Get()
	if comp == nil {
		return ""
	}

	localMarker := ctx.NextRouteMarker()
	markerID := ctx.FullID(localMarker)

	// Advance counters for all options to stay in sync with SSR.
	// Cache the Render() trees so wasmBindStoreComponent can reuse them.
	var cached []wasmCachedOption
	seen := make(map[string]bool)
	var activeHTML string
	for _, opt := range v.Options() {
		optComp, ok2 := opt.(Component)
		if !ok2 || optComp == nil {
			continue
		}
		name := componentName(optComp)
		if seen[name] {
			continue
		}
		seen[name] = true

		branchCtx := &WASMRenderContext{
			IDCounter: IDCounter{Prefix: wasmChildPrefix(ctx, name)},
		}
		var scopeAttr string
		if _, ok3 := optComp.(HasStyle); ok3 {
			scopeAttr = GetOrCreateScope(name)
			branchCtx.ScopeAttr = scopeAttr
		}
		tree := optComp.Render()
		branchHTML := wasmNodeToHTML(tree, branchCtx)
		if optComp == comp {
			activeHTML = branchHTML
		}
		cached = append(cached, wasmCachedOption{
			comp:      optComp,
			name:      name,
			tree:      tree,
			scopeAttr: scopeAttr,
		})
	}

	wasmRenderedTrees[markerID] = cached

	return "<!--" + markerID + "s-->" + activeHTML + "<!--" + markerID + "-->"
}

// wasmBindNodeToHTML renders a BindNode (text interpolation).
func wasmBindNodeToHTML(b *BindNode, ctx *WASMRenderContext) string {
	var value string
	if g, ok := b.StoreRef.(AnyGetter); ok {
		value = anyToString(g.GetAny())
	}

	localMarker := ctx.NextTextMarker()
	markerID := ctx.FullID(localMarker)

	if !b.IsHTML {
		value = escapeHTML(value)
	}
	return "<!--" + markerID + "s-->" + value + "<!--" + markerID + "-->"
}

// wasmIfNodeToHTML renders an IfNode, evaluating conditions at runtime.
func wasmIfNodeToHTML(i *IfNode, ctx *WASMRenderContext) string {
	localMarker := ctx.NextIfMarker()
	markerID := ctx.FullID(localMarker)

	// Advance counters through ALL branches to stay in sync with SSR.
	// SSR renders all branches sequentially, advancing counters through each.
	var activeHTML string
	activeFound := false
	for _, branch := range i.Branches {
		branchCtx := &WASMRenderContext{
			IDCounter:   ctx.IDCounter,
			ScopeAttr:   ctx.ScopeAttr,
			SlotContent: ctx.SlotContent,
		}
		branchHTML := wasmChildrenToHTML(branch.Children, branchCtx)
		ctx.IDCounter = branchCtx.IDCounter

		if !activeFound && branch.Cond.Eval() {
			activeHTML = branchHTML
			activeFound = true
		}
	}

	if len(i.ElseNode) > 0 {
		elseCtx := &WASMRenderContext{
			IDCounter:   ctx.IDCounter,
			ScopeAttr:   ctx.ScopeAttr,
			SlotContent: ctx.SlotContent,
		}
		elseHTML := wasmChildrenToHTML(i.ElseNode, elseCtx)
		ctx.IDCounter = elseCtx.IDCounter
		if !activeFound {
			activeHTML = elseHTML
		}
	}

	return "<!--" + markerID + "s-->" + activeHTML + "<!--" + markerID + "-->"
}

// wasmEachNodeToHTML renders an EachNode.
func wasmEachNodeToHTML(e *EachNode, ctx *WASMRenderContext) string {
	localMarker := ctx.NextEachMarker()
	markerID := ctx.FullID(localMarker)

	var itemsHTML string
	switch list := e.ListRef.(type) {
	case *List[string]:
		items := list.Get()
		if len(items) == 0 && len(e.ElseNode) > 0 {
			itemsHTML = wasmChildrenToHTML(e.ElseNode, ctx)
		} else {
			for i, item := range items {
				itemsHTML += wasmNodeToHTML(e.Body(item, i), ctx)
			}
		}
	case *List[int]:
		items := list.Get()
		if len(items) == 0 && len(e.ElseNode) > 0 {
			itemsHTML = wasmChildrenToHTML(e.ElseNode, ctx)
		} else {
			for i, item := range items {
				itemsHTML += wasmNodeToHTML(e.Body(item, i), ctx)
			}
		}
	}

	return "<!--" + markerID + "s-->" + itemsHTML + "<!--" + markerID + "-->"
}

// wasmComponentNodeToHTML renders a nested ComponentNode.
func wasmComponentNodeToHTML(c *ComponentNode, ctx *WASMRenderContext) string {
	comp, ok2 := c.Instance.(Component)
	if !ok2 {
		return ""
	}

	compMarker := ctx.NextCompMarker()
	fullCompPrefix := ctx.FullID(compMarker)

	var scopeAttr string
	if _, ok3 := c.Instance.(HasStyle); ok3 {
		scopeAttr = GetOrCreateScope(c.Name)
	}

	// Render slot content with parent context
	slotHTML := wasmChildrenToHTML(c.Children, ctx)

	childCtx := &WASMRenderContext{
		IDCounter:   IDCounter{Prefix: fullCompPrefix},
		ScopeAttr:   scopeAttr,
		SlotContent: slotHTML,
	}

	// Cache the Render() result so wasmBindComponentNode can reuse it
	// without calling Render() again (which would re-register handlers).
	tree := comp.Render()
	c.renderCache = tree
	return wasmNodeToHTML(tree, childCtx)
}

// wasmAttrToHTML renders a NodeAttr to HTML string.
func wasmAttrToHTML(attr NodeAttr, ctx *WASMRenderContext) string {
	switch a := attr.(type) {
	case *StaticAttr:
		return a.Name + `="` + escapeAttr(a.Value) + `"`
	case *DynAttrAttr:
		localID := ctx.NextAttrID()
		fullID := ctx.FullID(localID)
		var attrValue string
		for _, part := range a.Parts {
			switch v := part.(type) {
			case string:
				attrValue += v
			case AnyGetter:
				attrValue += anyToString(v.GetAny())
			}
		}
		return `data-attrbind="` + fullID + `" ` + a.Name + `="` + escapeAttr(attrValue) + `"`
	}
	return ""
}

// wasmChildPrefix computes a child context prefix.
func wasmChildPrefix(ctx *WASMRenderContext, name string) string {
	if ctx.Prefix != "" {
		return ctx.Prefix + "_" + name
	}
	return name
}
