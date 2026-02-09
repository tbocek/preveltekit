//go:build js && wasm

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
	case *TextNode:
		return escapeHTML(node.Content)
	case *HtmlNode:
		return wasmHtmlNodeToHTML(node, ctx)
	case *Fragment:
		return wasmFragmentToHTML(node, ctx)
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

// wasmHtmlNodeToHTML renders an HtmlNode to HTML.
func wasmHtmlNodeToHTML(h *HtmlNode, ctx *WASMRenderContext) string {
	html := wasmRenderParts(h, ctx)

	if len(h.AttrConds) > 0 || len(h.Events) > 0 {
		html = wasmInjectChainedAttrs(h, html, ctx)
	}

	if h.BoundStore != nil {
		html = wasmInjectBind(h, html, ctx)
	}

	return html
}

// wasmRenderParts renders the Parts slice of an HtmlNode.
func wasmRenderParts(h *HtmlNode, ctx *WASMRenderContext) string {
	var s string
	for _, part := range h.Parts {
		switch v := part.(type) {
		case string:
			if ctx.ScopeAttr != "" {
				s += injectScopeClass(v, ctx.ScopeAttr)
			} else {
				s += v
			}
		case Node:
			s += wasmNodeToHTML(v, ctx)
		case NodeAttr:
			s += wasmAttrToHTML(v, ctx)
		case *Store[string]:
			bind := &BindNode{StoreRef: v, IsHTML: false}
			s += wasmBindNodeToHTML(bind, ctx)
		case *Store[int]:
			bind := &BindNode{StoreRef: v, IsHTML: false}
			s += wasmBindNodeToHTML(bind, ctx)
		case *Store[bool]:
			bind := &BindNode{StoreRef: v, IsHTML: false}
			s += wasmBindNodeToHTML(bind, ctx)
		case *Store[float64]:
			bind := &BindNode{StoreRef: v, IsHTML: false}
			s += wasmBindNodeToHTML(bind, ctx)
		case *Store[Component]:
			s += wasmStoreComponentToHTML(v, ctx)
		default:
			s += escapeHTML(anyToString(v))
		}
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

// wasmInjectBind handles two-way input binding.
func wasmInjectBind(h *HtmlNode, html string, ctx *WASMRenderContext) string {
	switch s := h.BoundStore.(type) {
	case *Store[string]:
		return injectAttrs(html, `id="`+s.ID()+`" value="`+escapeAttr(s.Get())+`"`)
	case *Store[int]:
		return injectAttrs(html, `id="`+s.ID()+`" value="`+itoa(s.Get())+`"`)
	case *Store[bool]:
		checked := ""
		if s.Get() {
			checked = " checked"
		}
		return injectAttrs(html, `id="`+s.ID()+`"`+checked)
	}
	return html
}

// wasmInjectChainedAttrs injects AttrConds and Events into the first HTML tag.
func wasmInjectChainedAttrs(h *HtmlNode, html string, ctx *WASMRenderContext) string {
	var elementID string
	if len(h.Events) > 0 {
		elementID = h.Events[0].ID
	} else {
		localID := ctx.NextClassID()
		elementID = ctx.FullID(localID)
	}

	attrValues := make(map[string][]string)

	for _, ac := range h.AttrConds {
		if ac.Cond.Eval() {
			if tv := attrValStr(ac.TrueValue); tv != "" {
				attrValues[ac.Name] = append(attrValues[ac.Name], tv)
			}
		} else {
			if fv := attrValStr(ac.FalseValue); fv != "" {
				attrValues[ac.Name] = append(attrValues[ac.Name], fv)
			}
		}
	}

	var extraAttrs string
	if len(h.Events) > 0 {
		var names string
		for i, ev := range h.Events {
			if i > 0 {
				names += ","
			}
			names += ev.Event
		}
		extraAttrs = ` data-on="` + names + `"`
	}

	return injectIDAndMergeAttrs(html, elementID, attrValues, extraAttrs)
}

// attrValStr extracts a string value from an AttrCond value.
func attrValStr(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case *Store[string]:
		return val.Get()
	case *Store[int]:
		return itoa(val.Get())
	case *Store[bool]:
		if val.Get() {
			return "true"
		}
		return "false"
	}
	return ""
}

// wasmBindNodeToHTML renders a BindNode (text interpolation).
func wasmBindNodeToHTML(b *BindNode, ctx *WASMRenderContext) string {
	var value string
	switch s := b.StoreRef.(type) {
	case *Store[string]:
		value = s.Get()
	case *Store[int]:
		value = itoa(s.Get())
	case *Store[bool]:
		if s.Get() {
			value = "true"
		} else {
			value = "false"
		}
	case *Store[float64]:
		value = ftoa(s.Get())
	}

	localMarker := ctx.NextTextMarker()
	markerID := ctx.FullID(localMarker)

	if !b.IsHTML {
		value = escapeHTML(value)
	}
	return "<!--" + markerID + "s-->" + value + "<!--" + markerID + "-->"
}

// wasmFragmentToHTML renders a Fragment.
func wasmFragmentToHTML(f *Fragment, ctx *WASMRenderContext) string {
	return wasmChildrenToHTML(f.Children, ctx)
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

	// Set props
	setComponentProps(comp, c.Props)

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
	case *ClassAttr:
		return `class="` + joinStrings(a.Classes, " ") + `"`
	case *StaticAttr:
		return a.Name + `="` + escapeAttr(a.Value) + `"`
	case *DynAttrAttr:
		localID := ctx.NextAttrID()
		fullID := ctx.FullID(localID)
		attrValue := a.Template
		for i, store := range a.Stores {
			placeholder := "{" + itoa(i) + "}"
			var storeVal string
			switch s := store.(type) {
			case *Store[string]:
				storeVal = s.Get()
			case *Store[int]:
				storeVal = itoa(s.Get())
			}
			attrValue = replaceAll(attrValue, placeholder, storeVal)
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

// replaceAll replaces all occurrences of old with new in s.
func replaceAll(s, old, new string) string {
	if old == "" {
		return s
	}
	var result []byte
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result = append(result, new...)
			i += len(old)
		} else {
			result = append(result, s[i])
			i++
		}
	}
	return string(result)
}

// joinStrings joins strings with a separator (avoids strings package).
func joinStrings(ss []string, sep string) string {
	var r string
	for i, s := range ss {
		if i > 0 {
			r += sep
		}
		r += s
	}
	return r
}
