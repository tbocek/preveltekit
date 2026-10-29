//go:build wasm

package preveltekit

// setupBlocks tracks which if/each/component blocks (by marker ID) have been
// set up, to avoid wiring them twice.
var setupBlocks = make(map[string]bool)

// Hydrate sets up DOM bindings for reactivity.
// Walks the Render() tree to discover all bindings directly.
func Hydrate(app ComponentRoot) {
	// Create fresh app instance with initialized stores
	if hn, ok := app.(HasNew); ok {
		app = hn.New().(ComponentRoot)
	}

	// Call OnMount before Render to match SSR order
	if om, ok := app.(HasOnMount); ok {
		om.OnMount()
	}

	// Create app scope before tree walk to match SSR order
	var appScope string
	if _, ok := app.(HasStyle); ok {
		appScope = GetOrCreateScope("app")
	}

	// Walk the Render() tree to discover and wire all bindings
	ctx := &WASMRenderContext{
		ScopeAttr: appScope,
	}
	cleanup := &cleanupBag{}
	wasmWalkAndBind(app.Render(), ctx, cleanup)

	// Keep WASM running
	select {}
}

// wasmWalkAndBind walks a Node tree, discovering bindings and wiring them to the DOM.
// The ctx IDCounter must advance in the same order as SSR's nodeToHTML.
func wasmWalkAndBind(n Node, ctx *WASMRenderContext, cleanup *cleanupBag) {
	if n == nil {
		return
	}
	switch node := n.(type) {
	case *HtmlNode:
		wasmBindHtmlNode(node, ctx, cleanup)

	case *RawHTMLNode:
		// Static unescaped HTML, no bindings needed

	case *FragmentNode:
		for _, child := range node.Children {
			wasmWalkAndBind(child, ctx, cleanup)
		}

	case *BindNode:
		wasmBindTextNode(node, ctx, cleanup)

	case *IfNode:
		wasmBindIfNode(node, ctx, cleanup)

	case *EachNode:
		wasmBindEachNode(node, ctx, cleanup)

	case *ComponentNode:
		wasmBindComponentNode(node, ctx, cleanup)

	case *SlotNode:
		// Slot content was already rendered, nothing to bind

	case *TextNode:
		// Static text, no bindings needed
	}
}

// wasmBindHtmlNode walks an HtmlNode, wiring events, binds, AttrConds, and recursing into parts.
func wasmBindHtmlNode(h *HtmlNode, ctx *WASMRenderContext, cleanup *cleanupBag) {
	// Determine element ID (must match SSR's logic)
	var elementID string
	hasChainedAttrs := len(h.AttrConds) > 0 || len(h.Events) > 0

	if hasChainedAttrs {
		if len(h.Events) > 0 {
			elementID = h.Events[0].ID
		} else {
			elementID = ctx.NextClassID()
		}
	}

	// Wire events
	for _, ev := range h.Events {
		if handler := GetHandler(ev.ID); handler != nil {
			bindEvent(cleanup, ev.ID, ev.Event, handler)
		}
	}

	// Wire AttrCond bindings
	for _, ac := range h.AttrConds {
		wasmBindAttrCond(elementID, ac, cleanup)
	}

	// Wire two-way binding
	if h.BoundStore != nil {
		bindID := ctx.NextBindID()
		wasmBindInput(bindID, h.BoundStore, cleanup)
	}

	// Recurse into DynAttrs then Children
	for _, da := range h.DynAttrs {
		wasmBindDynAttr(da, ctx, cleanup)
	}
	for _, child := range h.Children {
		switch v := child.(type) {
		case Node:
			wasmWalkAndBind(v, ctx, cleanup)
		case *Store[Component]:
			wasmBindStoreComponent(v, ctx, cleanup)
		case AnySubscriber:
			bind := &BindNode{StoreRef: v, IsHTML: false}
			wasmBindTextNode(bind, ctx, cleanup)
		}
	}
}

// wasmBindTextNode wires a text binding (BindNode) to the DOM.
func wasmBindTextNode(b *BindNode, ctx *WASMRenderContext, cleanup *cleanupBag) {
	markerID := ctx.NextTextMarker()

	if s, ok := b.StoreRef.(AnySubscriber); ok {
		g := b.StoreRef.(AnyGetter)
		s.OnChangeAny(func() {
			val := anyToString(g.GetAny())
			if b.IsHTML {
				replaceMarkerContent(markerID, val)
			} else {
				replaceMarkerContent(markerID, escapeHTML(val))
			}
		})
	}
}

// wasmBindInput wires a two-way input binding using the bind element ID.
func wasmBindInput(bindID string, boundStore any, cleanup *cleanupBag) {
	switch s := boundStore.(type) {
	case *Store[string]:
		cleanup.Add(bindInput(bindID, s))
	case *Store[int]:
		cleanup.Add(bindInputInt(bindID, s))
	case *Store[bool]:
		cleanup.Add(bindCheckbox(bindID, s))
	}
}

// wasmBindAttrCond wires a conditional attribute binding.
func wasmBindAttrCond(elementID string, ac *AttrCond, cleanup *cleanupBag) {
	el := getEl(elementID)
	if !exists(el) {
		return
	}

	// Extract condition stores for subscription
	var condStores []any
	if fc, ok := ac.Cond.(*FuncCondition); ok {
		condStores = fc.Stores
	}

	updateAttr := func() {
		active := ac.Cond.Eval()
		if ac.Name == "class" {
			classList := el.Get("classList")
			trueVal := evalAttrValue(ac.TrueValue)
			falseVal := evalAttrValue(ac.FalseValue)
			if active && trueVal != "" {
				classList.Call("add", trueVal)
			} else if trueVal != "" {
				classList.Call("remove", trueVal)
			}
			if !active && falseVal != "" {
				classList.Call("add", falseVal)
			} else if falseVal != "" {
				classList.Call("remove", falseVal)
			}
		} else {
			var value string
			if active {
				value = evalAttrValue(ac.TrueValue)
			} else {
				value = evalAttrValue(ac.FalseValue)
			}
			if value != "" {
				el.Call("setAttribute", ac.Name, value)
			} else {
				el.Call("removeAttribute", ac.Name)
			}
		}
	}

	updateAttr()

	for _, cs := range condStores {
		subscribeToStore(cs, updateAttr)
	}
	// Also subscribe to the values if they are stores (no-op for strings)
	subscribeToStore(ac.TrueValue, updateAttr)
	subscribeToStore(ac.FalseValue, updateAttr)
}

// wasmBindDynAttr wires a dynamic attribute binding.
func wasmBindDynAttr(da *DynAttrAttr, ctx *WASMRenderContext, cleanup *cleanupBag) {
	fullID := ctx.NextAttrID()

	el := document.Call("querySelector", `[data-attrbind="`+fullID+`"]`)
	if !exists(el) {
		el = getEl(fullID)
	}
	if !exists(el) {
		return
	}

	updateAttr := func() {
		var value string
		for _, part := range da.Parts {
			switch v := part.(type) {
			case string:
				value += v
			case AnyGetter:
				value += anyToString(v.GetAny())
			}
		}
		el.Call("setAttribute", da.Name, value)
	}

	updateAttr()

	for _, part := range da.Parts {
		if _, ok := part.(string); !ok {
			subscribeToStore(part, updateAttr)
		}
	}
}

// wasmBindIfNode wires an if-block with reactive condition evaluation.
func wasmBindIfNode(ifNode *IfNode, ctx *WASMRenderContext, cleanup *cleanupBag) {
	markerID := ctx.NextIfMarker()

	// Advance counters through ALL branches to stay in sync with SSR.
	// Save the counter state at the start of each branch so the bind pass
	// can use the correct counters (matching what SSR produced in the HTML).
	branchCounters := make([]IDCounter, len(ifNode.Branches))
	for i, branch := range ifNode.Branches {
		branchCounters[i] = ctx.IDCounter
		branchCtx := &WASMRenderContext{
			IDCounter:   ctx.IDCounter,
			ScopeAttr:   ctx.ScopeAttr,
			SlotContent: ctx.SlotContent,
		}
		wasmChildrenToHTML(branch.Children, branchCtx)
		ctx.IDCounter = branchCtx.IDCounter
	}
	var elseCounter IDCounter
	if len(ifNode.ElseNode) > 0 {
		elseCounter = ctx.IDCounter
		elseCtx := &WASMRenderContext{
			IDCounter:   ctx.IDCounter,
			ScopeAttr:   ctx.ScopeAttr,
			SlotContent: ctx.SlotContent,
		}
		wasmChildrenToHTML(ifNode.ElseNode, elseCtx)
		ctx.IDCounter = elseCtx.IDCounter
	}

	// Skip if already set up
	if setupBlocks[markerID] {
		return
	}
	setupBlocks[markerID] = true

	currentCleanup := &cleanupBag{}
	currentBranchIdx := -2

	// Collect condition stores for subscription
	var condStores []any
	for _, branch := range ifNode.Branches {
		if fc, ok := branch.Cond.(*FuncCondition); ok {
			condStores = append(condStores, fc.Stores...)
		}
	}

	scopeAttr := ctx.ScopeAttr
	slotContent := ctx.SlotContent

	updateIfBlock := func() {
		activeBranchIdx := -1
		var activeNodes []Node

		for i, branch := range ifNode.Branches {
			if branch.Cond.Eval() {
				activeBranchIdx = i
				activeNodes = branch.Children
				break
			}
		}
		if activeBranchIdx == -1 {
			activeNodes = ifNode.ElseNode
		}

		if currentBranchIdx == activeBranchIdx && currentBranchIdx != -2 {
			return
		}
		currentBranchIdx = activeBranchIdx

		// Render active branch to HTML
		renderCtx := &WASMRenderContext{
			ScopeAttr:   scopeAttr,
			SlotContent: slotContent,
		}
		html := wasmChildrenToHTML(activeNodes, renderCtx)
		replaceMarkerContent(markerID, html)

		// Release old bindings, wire new ones
		currentCleanup.Release()
		currentCleanup = &cleanupBag{}

		// Re-walk the active branch nodes to wire bindings on new DOM
		bindCtx := &WASMRenderContext{
			ScopeAttr:   scopeAttr,
			SlotContent: slotContent,
		}
		for _, child := range activeNodes {
			wasmWalkAndBind(child, bindCtx, currentCleanup)
		}
	}

	// Subscribe to condition stores
	seen := make(map[string]bool)
	for _, store := range condStores {
		if id, ok := store.(HasID); ok {
			sid := id.ID()
			if seen[sid] {
				continue
			}
			seen[sid] = true
		}
		subscribeToStore(store, updateIfBlock)
	}

	// Initial sync: wire bindings for the currently active branch
	// (the DOM already has the correct HTML from SSR).
	// Use the saved counter state from the counter-advance pass so
	// marker IDs match what SSR put in the DOM.
	activeBranchIdx := -1
	var activeNodes []Node
	var activeCounter IDCounter
	for i, branch := range ifNode.Branches {
		if branch.Cond.Eval() {
			activeBranchIdx = i
			activeNodes = branch.Children
			activeCounter = branchCounters[i]
			break
		}
	}
	if activeBranchIdx == -1 {
		activeNodes = ifNode.ElseNode
		activeCounter = elseCounter
	}
	currentBranchIdx = activeBranchIdx

	// Wire bindings for the initially active branch with correct counters
	bindCtx := &WASMRenderContext{
		IDCounter:   activeCounter,
		ScopeAttr:   scopeAttr,
		SlotContent: slotContent,
	}
	for _, child := range activeNodes {
		wasmWalkAndBind(child, bindCtx, currentCleanup)
	}
}

// wasmBindEachNode wires an each-block with reactive list rendering.
func wasmBindEachNode(eachNode *EachNode, ctx *WASMRenderContext, cleanup *cleanupBag) {
	markerID := ctx.NextEachMarker()

	// Skip if already set up
	if setupBlocks[markerID] {
		return
	}
	setupBlocks[markerID] = true

	scopeAttr := ctx.ScopeAttr

	// renderItems renders the current items (or the else content) to HTML.
	renderItems := func() string {
		items := eachNode.items()
		if len(items) == 0 {
			return wasmChildrenToHTML(eachNode.ElseNode, &WASMRenderContext{ScopeAttr: scopeAttr})
		}
		var html string
		for i, item := range items {
			html += wasmNodeToHTML(eachNode.Body(item, i), &WASMRenderContext{ScopeAttr: scopeAttr})
		}
		return html
	}

	// bindItems wires the current items (or the else content) to the DOM.
	bindItems := func(c *cleanupBag) {
		items := eachNode.items()
		if len(items) == 0 {
			bindCtx := &WASMRenderContext{ScopeAttr: scopeAttr}
			for _, child := range eachNode.ElseNode {
				wasmWalkAndBind(child, bindCtx, c)
			}
			return
		}
		for i, item := range items {
			wasmWalkAndBind(eachNode.Body(item, i), &WASMRenderContext{ScopeAttr: scopeAttr}, c)
		}
	}

	// On change: re-render, release the bindings of the previous render, wire the new items
	current := &cleanupBag{}
	eachNode.onChange(func() {
		replaceMarkerContent(markerID, renderItems())
		current.Release()
		current = &cleanupBag{}
		bindItems(current)
	})

	// Wire bindings for initially rendered items (DOM already has SSR content)
	bindItems(cleanup)
}

// wasmBindStoreComponent wires a Store[Component] binding.
func wasmBindStoreComponent(v *Store[Component], ctx *WASMRenderContext, cleanup *cleanupBag) {
	markerID := ctx.NextRouteMarker()

	// Check if the HTML render pass already cached trees for this marker.
	// If so, reuse them to avoid calling Render() again (which re-registers handlers).
	// If not (e.g. top-level Hydrate where there's no prior HTML pass), do the
	// counter-advance ourselves.
	var rendered []wasmCachedOption

	if cached, ok := wasmRenderedTrees[markerID]; ok {
		rendered = cached
		delete(wasmRenderedTrees, markerID)
	} else {
		// No cached trees — do counter-advance ourselves
		seen := make(map[string]bool)
		for _, opt := range v.Options() {
			optComp, ok := opt.(Component)
			if !ok || optComp == nil {
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
			if _, ok := optComp.(HasStyle); ok {
				scopeAttr = GetOrCreateScope(name)
				branchCtx.ScopeAttr = scopeAttr
			}
			tree := optComp.Render()
			wasmNodeToHTML(tree, branchCtx)

			rendered = append(rendered, wasmCachedOption{
				comp:      optComp,
				name:      name,
				tree:      tree,
				scopeAttr: scopeAttr,
			})
		}
	}

	// Skip if already set up
	if setupBlocks[markerID] {
		return
	}
	setupBlocks[markerID] = true

	currentCleanup := &cleanupBag{}
	currentName := ""
	firstCall := true

	updateBlock := func() {
		comp := v.Get()
		if comp == nil {
			return
		}
		name := componentName(comp)

		if name == currentName && !firstCall {
			return
		}

		if firstCall {
			firstCall = false
			currentName = name

			// Wire bindings for initial component (DOM has SSR content).
			// Reuse the tree from the counter-advance pass to avoid
			// re-registering handlers with new IDs.
			var tree Node
			var scopeAttr string
			for _, r := range rendered {
				if r.comp == comp {
					tree = r.tree
					scopeAttr = r.scopeAttr
					break
				}
			}
			if tree == nil {
				return
			}

			bindCtx := &WASMRenderContext{
				IDCounter: IDCounter{Prefix: wasmChildPrefix(ctx, name)},
				ScopeAttr: scopeAttr,
			}
			if om, ok := comp.(HasOnMount); ok {
				om.OnMount()
			}
			if od, ok := comp.(HasOnDestroy); ok {
				currentCleanup.AddDestroy(od.OnDestroy)
			}
			wasmWalkAndBind(tree, bindCtx, currentCleanup)
			return
		}

		currentName = name

		// Call OnMount on the new active component
		if om, ok := comp.(HasOnMount); ok {
			om.OnMount()
		}

		// Render new component to HTML (subsequent changes, not initial)
		renderTree := comp.Render()
		renderCtx := &WASMRenderContext{
			IDCounter: IDCounter{Prefix: wasmChildPrefix(ctx, name)},
		}
		if _, ok := comp.(HasStyle); ok {
			renderCtx.ScopeAttr = GetOrCreateScope(name)
		}
		html := wasmNodeToHTML(renderTree, renderCtx)
		replaceMarkerContent(markerID, html)

		// Release old bindings (fires OnDestroy), wire new ones
		currentCleanup.Release()
		currentCleanup = &cleanupBag{}

		// Walk the same tree we just rendered (don't call Render() again)
		bindCtx := &WASMRenderContext{
			IDCounter: IDCounter{Prefix: wasmChildPrefix(ctx, name)},
		}
		if _, ok := comp.(HasStyle); ok {
			bindCtx.ScopeAttr = GetOrCreateScope(name)
		}
		if od, ok := comp.(HasOnDestroy); ok {
			currentCleanup.AddDestroy(od.OnDestroy)
		}
		wasmWalkAndBind(renderTree, bindCtx, currentCleanup)
	}

	v.OnChange(func(_ Component) { updateBlock() })
	updateBlock()
}

// wasmBindComponentNode wires a nested ComponentNode.
func wasmBindComponentNode(c *ComponentNode, ctx *WASMRenderContext, cleanup *cleanupBag) {
	comp, ok := c.Instance.(Component)
	if !ok {
		return
	}

	fullCompPrefix := ctx.NextCompMarker()

	var scopeAttr string
	if _, ok := c.Instance.(HasStyle); ok {
		scopeAttr = GetOrCreateScope(c.Name)
	}

	// Walk slot children with parent context
	for _, child := range c.Children {
		wasmWalkAndBind(child, ctx, cleanup)
	}

	// Use the cached Render() tree if available (from wasmComponentNodeToHTML
	// during the counter-advance pass). This avoids calling Render() again,
	// which would re-register handlers with new IDs.
	tree := c.renderCache
	if tree == nil {
		tree = comp.Render()
	}

	// Call OnMount when the component is wired
	if om, ok := comp.(HasOnMount); ok {
		om.OnMount()
	}

	// Register OnDestroy if the component implements it
	if od, ok := comp.(HasOnDestroy); ok {
		cleanup.AddDestroy(od.OnDestroy)
	}

	// Walk component's own Render tree with child context
	childCtx := &WASMRenderContext{
		IDCounter: IDCounter{Prefix: fullCompPrefix},
		ScopeAttr: scopeAttr,
	}
	wasmWalkAndBind(tree, childCtx, cleanup)
}

// replaceMarkerContent replaces all DOM nodes between <!--{markerID}s--> and <!--{markerID}-->
// with new HTML content.
func replaceMarkerContent(markerID string, html string) {
	endMarker := findComment(markerID)
	if endMarker.IsNull() {
		return
	}
	startMarker := findComment(markerID + "s")
	if startMarker.IsNull() {
		return
	}
	parent := endMarker.Get("parentNode")

	// Remove all nodes between start and end markers
	for {
		next := startMarker.Get("nextSibling")
		if next.IsNull() || next.Equal(endMarker) {
			break
		}
		parent.Call("removeChild", next)
	}

	// Parse new HTML and insert before end marker
	if html != "" {
		tmpl := document.Call("createElement", "template")
		tmpl.Set("innerHTML", html)
		frag := tmpl.Get("content")
		parent.Call("insertBefore", frag, endMarker)
	}
}

// subscribeToStore subscribes a callback to store changes.
func subscribeToStore(store any, callback func()) {
	if s, ok := store.(AnySubscriber); ok {
		s.OnChangeAny(callback)
	}
}

// SetSSRPath is a no-op in WASM (only used during SSR).
func SetSSRPath(path string) {}

func atoiSafe(s string) int {
	n := 0
	neg := false
	for i, c := range s {
		if c == '-' && i == 0 {
			neg = true
			continue
		}
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	if neg {
		return -n
	}
	return n
}
