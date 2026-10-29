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
	ctx := &renderContext{
		ScopeAttr: appScope,
	}
	cleanup := &cleanupBag{}
	wasmWalkAndBind(app.Render(), ctx, cleanup)

	// Keep WASM running
	select {}
}

// wasmWalkAndBind walks a Node tree, discovering bindings and wiring them to the DOM.
// The ctx IDCounter must advance in the same order as SSR's nodeToHTML.
func wasmWalkAndBind(n Node, ctx *renderContext, cleanup *cleanupBag) {
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
		// The slot content is first wired by the outer component; those
		// bindings belong to this place and are released with it. After a
		// re-render the slot content is new DOM: wire it again, with the
		// counters it was rendered with.
		s := ctx.slot
		if s == nil || s.release == nil {
			return
		}
		if !ctx.fresh {
			cleanup.AddDestroy(s.release)
			return
		}
		s.release()
		bindCtx := s.ctx.with(s.ctx.IDCounter)
		bindCtx.fresh = true
		for _, child := range s.nodes {
			wasmWalkAndBind(child, bindCtx, cleanup)
		}

	case *TextNode:
		// Static text, no bindings needed
	}
}

// wasmBindHtmlNode walks an HtmlNode, wiring events, binds, AttrConds, and recursing into parts.
func wasmBindHtmlNode(h *HtmlNode, ctx *renderContext, cleanup *cleanupBag) {
	// Element ID (must match the renderer's)
	elementID := h.elementID(&ctx.IDCounter)

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
func wasmBindTextNode(b *BindNode, ctx *renderContext, cleanup *cleanupBag) {
	markerID := ctx.NextTextMarker()

	if g, ok := b.StoreRef.(AnyGetter); ok {
		subscribeToStore(b.StoreRef, func() {
			val := anyToString(g.GetAny())
			if b.IsHTML {
				replaceMarkerContent(markerID, val)
			} else {
				replaceMarkerContent(markerID, escapeHTML(val))
			}
		}, cleanup)
	}
}

// wasmBindInput wires a two-way input binding using the bind element ID.
func wasmBindInput(bindID string, boundStore any, cleanup *cleanupBag) {
	switch s := boundStore.(type) {
	case *Store[string]:
		bindInput(cleanup, bindID, s)
	case *Store[int]:
		bindInputInt(cleanup, bindID, s)
	case *Store[bool]:
		bindCheckbox(cleanup, bindID, s)
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
			if value := ac.value(); value != "" {
				el.Call("setAttribute", ac.Name, value)
			} else {
				el.Call("removeAttribute", ac.Name)
			}
		}
	}

	updateAttr()

	for _, cs := range condStores {
		subscribeToStore(cs, updateAttr, cleanup)
	}
	// Also subscribe to the values if they are stores (no-op for strings)
	subscribeToStore(ac.TrueValue, updateAttr, cleanup)
	subscribeToStore(ac.FalseValue, updateAttr, cleanup)
}

// wasmBindDynAttr wires a dynamic attribute binding.
func wasmBindDynAttr(da *DynAttrAttr, ctx *renderContext, cleanup *cleanupBag) {
	fullID := ctx.NextAttrID()

	el := document.Call("querySelector", `[data-attrbind="`+fullID+`"]`)
	if !exists(el) {
		el = getEl(fullID)
	}
	if !exists(el) {
		return
	}

	updateAttr := func() {
		el.Call("setAttribute", da.Name, da.value())
	}

	updateAttr()

	for _, part := range da.Parts {
		if _, ok := part.(string); !ok {
			subscribeToStore(part, updateAttr, cleanup)
		}
	}
}

// wasmBindIfNode wires an if-block with reactive condition evaluation.
func wasmBindIfNode(ifNode *IfNode, ctx *renderContext, cleanup *cleanupBag) {
	markerID := ctx.NextIfMarker()

	// Advance the counters through ALL branches like the renderer does, and
	// save the counter state at the start of each branch: the branch is
	// rendered and bound with it, so its IDs match the SSR HTML.
	bodies := ifNode.bodies()
	counters := make([]IDCounter, len(bodies))
	for i, nodes := range bodies {
		counters[i] = ctx.IDCounter
		branchCtx := ctx.with(ctx.IDCounter)
		childrenToHTML(nodes, branchCtx)
		ctx.IDCounter = branchCtx.IDCounter
	}

	current := &cleanupBag{}
	if !setUpBlock(markerID, cleanup, &current) {
		return
	}

	base := ctx.with(IDCounter{})
	wire := func(idx int, rerendered bool) {
		if idx < 0 {
			return
		}
		bindCtx := base.with(counters[idx])
		bindCtx.fresh = bindCtx.fresh || rerendered
		for _, child := range bodies[idx] {
			wasmWalkAndBind(child, bindCtx, current)
		}
	}

	active := ifNode.active()
	update := func() {
		idx := ifNode.active()
		if idx == active {
			return
		}
		active = idx

		var html string
		if idx >= 0 {
			html = childrenToHTML(bodies[idx], base.with(counters[idx]))
		}
		replaceMarkerContent(markerID, html)

		// Release old bindings, wire the new branch
		current.Release()
		current = &cleanupBag{}
		wire(idx, true)
	}

	// Subscribe to the condition stores (each once)
	seen := make(map[string]bool)
	for _, branch := range ifNode.Branches {
		fc, ok := branch.Cond.(*FuncCondition)
		if !ok {
			continue
		}
		for _, store := range fc.Stores {
			if id, ok := store.(HasID); ok {
				if seen[id.ID()] {
					continue
				}
				seen[id.ID()] = true
			}
			subscribeToStore(store, update, cleanup)
		}
	}

	// Wire the initially active branch (the DOM already has its SSR HTML)
	wire(active, false)
}

// wasmBindEachNode wires an each-block with reactive list rendering.
func wasmBindEachNode(eachNode *EachNode, ctx *renderContext, cleanup *cleanupBag) {
	markerID := ctx.NextEachMarker()

	current := &cleanupBag{}
	if !setUpBlock(markerID, cleanup, &current) {
		return
	}

	base := ctx.with(IDCounter{})

	// bindItems wires the current items (or the else content) to the DOM,
	// with the same counters the renderer used for them.
	bindItems := func(c *cleanupBag, rerendered bool) {
		eachNode.forEach(markerID, func(nodes []Node, counter IDCounter) {
			bindCtx := base.with(counter)
			bindCtx.fresh = bindCtx.fresh || rerendered
			for _, child := range nodes {
				wasmWalkAndBind(child, bindCtx, c)
			}
		})
	}

	// On change: re-render, release the bindings of the previous render, wire the new items
	cleanup.AddDestroy(eachNode.onChange(func() {
		replaceMarkerContent(markerID, eachItemsHTML(eachNode, markerID, base))
		current.Release()
		current = &cleanupBag{}
		bindItems(current, true)
	}))

	// Wire bindings for initially rendered items (DOM already has SSR content)
	bindItems(current, false)
}

// wasmBindStoreComponent wires a Store[Component] binding.
func wasmBindStoreComponent(v *Store[Component], ctx *renderContext, cleanup *cleanupBag) {
	markerID := ctx.NextRouteMarker()

	// Reuse the trees cached by the HTML render pass for this marker, so
	// Render() isn't called again (which would re-register handlers).
	var rendered []cachedOption

	if cached, ok := renderedTrees[markerID]; ok {
		rendered = cached
		delete(renderedTrees, markerID)
	} else {
		// No cached trees (e.g. top-level Hydrate without a prior HTML pass):
		// render the options ourselves to advance the counters
		_, rendered = renderStoreOptions(v, ctx)
	}

	currentCleanup := &cleanupBag{}
	if !setUpBlock(markerID, cleanup, &currentCleanup) {
		return
	}

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

			bindCtx := &renderContext{
				IDCounter: IDCounter{Prefix: ctx.FullID(name)},
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
		renderCtx := &renderContext{
			IDCounter: IDCounter{Prefix: ctx.FullID(name)},
		}
		if _, ok := comp.(HasStyle); ok {
			renderCtx.ScopeAttr = GetOrCreateScope(name)
		}
		html := nodeToHTML(renderTree, renderCtx)
		replaceMarkerContent(markerID, html)

		// Release old bindings (fires OnDestroy), wire new ones
		currentCleanup.Release()
		currentCleanup = &cleanupBag{}

		// Walk the same tree we just rendered (don't call Render() again)
		bindCtx := &renderContext{
			IDCounter: IDCounter{Prefix: ctx.FullID(name)},
		}
		if _, ok := comp.(HasStyle); ok {
			bindCtx.ScopeAttr = GetOrCreateScope(name)
		}
		if od, ok := comp.(HasOnDestroy); ok {
			currentCleanup.AddDestroy(od.OnDestroy)
		}
		wasmWalkAndBind(renderTree, bindCtx, currentCleanup)
	}

	cleanup.AddDestroy(v.subscribe(func(Component) { updateBlock() }))
	updateBlock()
}

// wasmBindComponentNode wires a nested ComponentNode.
func wasmBindComponentNode(c *ComponentNode, ctx *renderContext, cleanup *cleanupBag) {
	comp, ok := c.Instance.(Component)
	if !ok {
		return
	}

	fullCompPrefix := ctx.NextCompMarker()
	scopeAttr := ctx.componentScope(c.Name, c.Instance)

	// Slot content for blocks inside the component that re-render: rendered
	// and wired again on each re-render, with the counters the slot had here.
	// The first wiring (below, before Render() like in SSR) goes into its own
	// bag, released when the place of the Slot() re-renders or goes away.
	slotBag := &cleanupBag{}
	cleanup.AddDestroy(slotBag.Release)
	slot := &slotContent{nodes: c.Children, ctx: ctx.with(ctx.IDCounter), release: slotBag.Release}

	// Walk slot children with parent context
	for _, child := range c.Children {
		wasmWalkAndBind(child, ctx, slotBag)
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
	wasmWalkAndBind(tree, ctx.component(fullCompPrefix, scopeAttr, slot), cleanup)
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

// subscribeToStore subscribes a callback to store changes until cleanup is
// released (a no-op for values that are not stores).
func subscribeToStore(store any, callback func(), cleanup *cleanupBag) {
	switch s := store.(type) {
	case subscriber:
		cleanup.AddDestroy(s.subscribeAny(callback))
	case AnySubscriber:
		s.OnChangeAny(callback) // user type without unsubscribe
	}
}

// setUpBlock marks the block with markerID as set up until cleanup is
// released, when its current content is released too. It returns false if
// the block is already set up. (When a block re-renders, the blocks nested in
// it are released and set up again for the new DOM.)
func setUpBlock(markerID string, cleanup *cleanupBag, current **cleanupBag) bool {
	if setupBlocks[markerID] {
		return false
	}
	setupBlocks[markerID] = true
	cleanup.AddDestroy(func() {
		(*current).Release()
		delete(setupBlocks, markerID)
	})
	return true
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
