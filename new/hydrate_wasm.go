//go:build wasm

package preveltekit

// Track which if-blocks have been set up to avoid duplicates
var setupIfBlocks = make(map[string]bool)

// Track which each-blocks have been set up to avoid duplicates
var setupEachBlocks = make(map[string]bool)

// Track which component-blocks have been set up to avoid duplicates
var setupComponentBlocks = make(map[string]bool)

// Hydrate sets up DOM bindings for reactivity.
// Walks the Render() tree to discover all bindings directly — no bindings.bin needed.
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
	cleanup := &Cleanup{}
	wasmWalkAndBind(app.Render(), ctx, cleanup)

	// Keep WASM running
	select {}
}

// wasmWalkAndBind walks a Node tree, discovering bindings and wiring them to the DOM.
// This replaces both walkNodeForComponents and the bindings.bin-based applyBindings.
// The ctx IDCounter must advance in the same order as SSR's nodeToHTML.
func wasmWalkAndBind(n Node, ctx *WASMRenderContext, cleanup *Cleanup) {
	if n == nil {
		return
	}
	switch node := n.(type) {
	case *TextNode:
		// Static text, nothing to bind

	case *HtmlNode:
		wasmBindHtmlNode(node, ctx, cleanup)

	case *Fragment:
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
	}
}

// wasmBindHtmlNode walks an HtmlNode, wiring events, binds, AttrConds, and recursing into parts.
func wasmBindHtmlNode(h *HtmlNode, ctx *WASMRenderContext, cleanup *Cleanup) {
	// Determine element ID (must match SSR's logic)
	var elementID string
	hasChainedAttrs := len(h.AttrConds) > 0 || len(h.Events) > 0

	if hasChainedAttrs {
		if len(h.Events) > 0 {
			elementID = h.Events[0].ID
		} else {
			localID := ctx.NextClassID()
			elementID = ctx.FullID(localID)
		}
	}

	// Wire events
	if len(h.Events) > 0 {
		var evts []Evt
		for _, ev := range h.Events {
			handler := GetHandler(ev.ID)
			if handler != nil {
				evts = append(evts, Evt{ev.ID, ev.Event, handler})
			}
		}
		if len(evts) > 0 {
			BindEvents(cleanup, evts)
		}
	}

	// Wire AttrCond bindings
	for _, ac := range h.AttrConds {
		wasmBindAttrCond(elementID, ac, cleanup)
	}

	// Wire two-way binding
	if h.BoundStore != nil {
		localID := ctx.NextBindID()
		bindID := ctx.FullID(localID)
		wasmBindInput(bindID, h.BoundStore, cleanup)
	}

	// Recurse into parts
	for _, part := range h.Parts {
		switch v := part.(type) {
		case Node:
			wasmWalkAndBind(v, ctx, cleanup)
		case NodeAttr:
			// Handle DynAttr bindings
			if da, ok2 := v.(*DynAttrAttr); ok2 {
				wasmBindDynAttr(da, ctx, cleanup)
			}
		case *Store[string]:
			bind := &BindNode{StoreRef: v, IsHTML: false}
			wasmBindTextNode(bind, ctx, cleanup)
		case *Store[int]:
			bind := &BindNode{StoreRef: v, IsHTML: false}
			wasmBindTextNode(bind, ctx, cleanup)
		case *Store[bool]:
			bind := &BindNode{StoreRef: v, IsHTML: false}
			wasmBindTextNode(bind, ctx, cleanup)
		case *Store[float64]:
			bind := &BindNode{StoreRef: v, IsHTML: false}
			wasmBindTextNode(bind, ctx, cleanup)
		case *Store[Component]:
			wasmBindStoreComponent(v, ctx, cleanup)
		}
	}
}

// wasmBindTextNode wires a text binding (BindNode) to the DOM.
func wasmBindTextNode(b *BindNode, ctx *WASMRenderContext, cleanup *Cleanup) {
	localMarker := ctx.NextTextMarker()
	markerID := ctx.FullID(localMarker)

	switch s := b.StoreRef.(type) {
	case *Store[string]:
		s.OnChange(func(v string) {
			if b.IsHTML {
				replaceMarkerContent(markerID, v)
			} else {
				replaceMarkerContent(markerID, escapeHTML(v))
			}
		})
	case *Store[int]:
		s.OnChange(func(v int) {
			replaceMarkerContent(markerID, escapeHTML(itoa(v)))
		})
	case *Store[bool]:
		s.OnChange(func(v bool) {
			val := "false"
			if v {
				val = "true"
			}
			replaceMarkerContent(markerID, escapeHTML(val))
		})
	case *Store[float64]:
		s.OnChange(func(v float64) {
			replaceMarkerContent(markerID, escapeHTML(ftoa(v)))
		})
	}
}

// wasmBindInput wires a two-way input binding using the bind element ID.
func wasmBindInput(bindID string, boundStore any, cleanup *Cleanup) {
	switch s := boundStore.(type) {
	case *Store[string]:
		BindInputs(cleanup, []Inp{{bindID, s}})
	case *Store[int]:
		cleanup.Add(BindInputInt(bindID, s))
	case *Store[bool]:
		BindCheckboxes(cleanup, []Chk{{bindID, s}})
	}
}

// wasmBindAttrCond wires a conditional attribute binding.
func wasmBindAttrCond(elementID string, ac *AttrCond, cleanup *Cleanup) {
	el := GetEl(elementID)
	if !ok(el) {
		return
	}

	// Extract condition store for subscription
	var condStore any
	if sc, ok2 := ac.Cond.(*StoreCondition); ok2 {
		condStore = sc.Store
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

	if condStore != nil {
		subscribeToStore(condStore, updateAttr)
	}
	// Also subscribe to value stores if dynamic
	if s, ok2 := ac.TrueValue.(*Store[string]); ok2 {
		subscribeToStore(s, updateAttr)
	}
	if s, ok2 := ac.FalseValue.(*Store[string]); ok2 {
		subscribeToStore(s, updateAttr)
	}
}

// wasmBindDynAttr wires a dynamic attribute binding.
func wasmBindDynAttr(da *DynAttrAttr, ctx *WASMRenderContext, cleanup *Cleanup) {
	localID := ctx.NextAttrID()
	fullID := ctx.FullID(localID)

	el := Document.Call("querySelector", `[data-attrbind="`+fullID+`"]`)
	if !ok(el) {
		el = GetEl(fullID)
	}
	if !ok(el) {
		return
	}

	updateAttr := func() {
		var value string
		for _, part := range da.Parts {
			switch v := part.(type) {
			case string:
				value += v
			default:
				value += storeToString(v)
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
func wasmBindIfNode(ifNode *IfNode, ctx *WASMRenderContext, cleanup *Cleanup) {
	localMarker := ctx.NextIfMarker()
	markerID := ctx.FullID(localMarker)

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
	if setupIfBlocks[markerID] {
		return
	}
	setupIfBlocks[markerID] = true

	currentCleanup := &Cleanup{}
	currentBranchIdx := -2

	// Collect condition stores for subscription
	var condStores []any
	for _, branch := range ifNode.Branches {
		if sc, ok2 := branch.Cond.(*StoreCondition); ok2 {
			condStores = append(condStores, sc.Store)
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
		currentCleanup = &Cleanup{}

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
		if id, ok2 := store.(HasID); ok2 {
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
func wasmBindEachNode(eachNode *EachNode, ctx *WASMRenderContext, cleanup *Cleanup) {
	localMarker := ctx.NextEachMarker()
	markerID := ctx.FullID(localMarker)

	// Skip if already set up
	if setupEachBlocks[markerID] {
		return
	}
	setupEachBlocks[markerID] = true

	scopeAttr := ctx.ScopeAttr

	// Render and bind items helper
	renderAndBindItems := func(renderItems func(body func(any, int) Node) string) {
		// Render all items to HTML
		html := renderItems(eachNode.Body)
		replaceMarkerContent(markerID, html)

		// Re-walk each item's nodes to wire bindings
		bindCleanup := &Cleanup{}
		switch list := eachNode.ListRef.(type) {
		case *List[string]:
			items := list.Get()
			if len(items) == 0 && len(eachNode.ElseNode) > 0 {
				bindCtx := &WASMRenderContext{ScopeAttr: scopeAttr}
				for _, child := range eachNode.ElseNode {
					wasmWalkAndBind(child, bindCtx, bindCleanup)
				}
			} else {
				for i, item := range items {
					bodyNode := eachNode.Body(item, i)
					bindCtx := &WASMRenderContext{ScopeAttr: scopeAttr}
					wasmWalkAndBind(bodyNode, bindCtx, bindCleanup)
				}
			}
		case *List[int]:
			items := list.Get()
			if len(items) == 0 && len(eachNode.ElseNode) > 0 {
				bindCtx := &WASMRenderContext{ScopeAttr: scopeAttr}
				for _, child := range eachNode.ElseNode {
					wasmWalkAndBind(child, bindCtx, bindCleanup)
				}
			} else {
				for i, item := range items {
					bodyNode := eachNode.Body(item, i)
					bindCtx := &WASMRenderContext{ScopeAttr: scopeAttr}
					wasmWalkAndBind(bodyNode, bindCtx, bindCleanup)
				}
			}
		}
	}

	// Subscribe to list changes
	switch list := eachNode.ListRef.(type) {
	case *List[string]:
		list.OnChange(func(items []string) {
			renderItems := func(body func(any, int) Node) string {
				if len(items) == 0 && len(eachNode.ElseNode) > 0 {
					renderCtx := &WASMRenderContext{ScopeAttr: scopeAttr}
					return wasmChildrenToHTML(eachNode.ElseNode, renderCtx)
				}
				var html string
				for i, item := range items {
					renderCtx := &WASMRenderContext{ScopeAttr: scopeAttr}
					html += wasmNodeToHTML(body(item, i), renderCtx)
				}
				return html
			}
			renderAndBindItems(renderItems)
		})
	case *List[int]:
		list.OnChange(func(items []int) {
			renderItems := func(body func(any, int) Node) string {
				if len(items) == 0 && len(eachNode.ElseNode) > 0 {
					renderCtx := &WASMRenderContext{ScopeAttr: scopeAttr}
					return wasmChildrenToHTML(eachNode.ElseNode, renderCtx)
				}
				var html string
				for i, item := range items {
					renderCtx := &WASMRenderContext{ScopeAttr: scopeAttr}
					html += wasmNodeToHTML(body(item, i), renderCtx)
				}
				return html
			}
			renderAndBindItems(renderItems)
		})
	}

	// Wire bindings for initially rendered items (DOM already has SSR content)
	switch list := eachNode.ListRef.(type) {
	case *List[string]:
		items := list.Get()
		if len(items) == 0 && len(eachNode.ElseNode) > 0 {
			bindCtx := &WASMRenderContext{ScopeAttr: scopeAttr}
			for _, child := range eachNode.ElseNode {
				wasmWalkAndBind(child, bindCtx, cleanup)
			}
		} else {
			for i, item := range items {
				bodyNode := eachNode.Body(item, i)
				bindCtx := &WASMRenderContext{ScopeAttr: scopeAttr}
				wasmWalkAndBind(bodyNode, bindCtx, cleanup)
			}
		}
	case *List[int]:
		items := list.Get()
		if len(items) == 0 && len(eachNode.ElseNode) > 0 {
			bindCtx := &WASMRenderContext{ScopeAttr: scopeAttr}
			for _, child := range eachNode.ElseNode {
				wasmWalkAndBind(child, bindCtx, cleanup)
			}
		} else {
			for i, item := range items {
				bodyNode := eachNode.Body(item, i)
				bindCtx := &WASMRenderContext{ScopeAttr: scopeAttr}
				wasmWalkAndBind(bodyNode, bindCtx, cleanup)
			}
		}
	}
}

// wasmBindStoreComponent wires a Store[Component] binding.
func wasmBindStoreComponent(v *Store[Component], ctx *WASMRenderContext, cleanup *Cleanup) {
	localMarker := ctx.NextRouteMarker()
	markerID := ctx.FullID(localMarker)

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
	if setupComponentBlocks[markerID] {
		return
	}
	setupComponentBlocks[markerID] = true

	currentCleanup := &Cleanup{}
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
			if om, ok2 := comp.(HasOnMount); ok2 {
				om.OnMount()
			}
			if od, ok2 := comp.(HasOnDestroy); ok2 {
				currentCleanup.AddDestroy(od.OnDestroy)
			}
			wasmWalkAndBind(tree, bindCtx, currentCleanup)
			return
		}

		currentName = name

		// Call OnMount on the new active component
		if om, ok2 := comp.(HasOnMount); ok2 {
			om.OnMount()
		}

		// Render new component to HTML (subsequent changes, not initial)
		renderTree := comp.Render()
		renderCtx := &WASMRenderContext{
			IDCounter: IDCounter{Prefix: wasmChildPrefix(ctx, name)},
		}
		if _, ok2 := comp.(HasStyle); ok2 {
			renderCtx.ScopeAttr = GetOrCreateScope(name)
		}
		html := wasmNodeToHTML(renderTree, renderCtx)
		replaceMarkerContent(markerID, html)

		// Release old bindings (fires OnDestroy), wire new ones
		currentCleanup.Release()
		currentCleanup = &Cleanup{}

		// Walk the same tree we just rendered (don't call Render() again)
		bindCtx := &WASMRenderContext{
			IDCounter: IDCounter{Prefix: wasmChildPrefix(ctx, name)},
		}
		if _, ok2 := comp.(HasStyle); ok2 {
			bindCtx.ScopeAttr = GetOrCreateScope(name)
		}
		if od, ok2 := comp.(HasOnDestroy); ok2 {
			currentCleanup.AddDestroy(od.OnDestroy)
		}
		wasmWalkAndBind(renderTree, bindCtx, currentCleanup)
	}

	v.OnChange(func(_ Component) { updateBlock() })
	updateBlock()
}

// wasmBindComponentNode wires a nested ComponentNode.
func wasmBindComponentNode(c *ComponentNode, ctx *WASMRenderContext, cleanup *Cleanup) {
	comp, ok2 := c.Instance.(Component)
	if !ok2 {
		return
	}

	compMarker := ctx.NextCompMarker()
	fullCompPrefix := ctx.FullID(compMarker)

	var scopeAttr string
	if _, ok3 := c.Instance.(HasStyle); ok3 {
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
	if om, ok3 := comp.(HasOnMount); ok3 {
		om.OnMount()
	}

	// Register OnDestroy if the component implements it
	if od, ok3 := comp.(HasOnDestroy); ok3 {
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
	endMarker := FindComment(markerID)
	if endMarker.IsNull() {
		return
	}
	startMarker := FindComment(markerID + "s")
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
		tmpl := Document.Call("createElement", "template")
		tmpl.Set("innerHTML", html)
		frag := tmpl.Get("content")
		parent.Call("insertBefore", frag, endMarker)
	}
}

// subscribeToStore subscribes a callback to store changes.
func subscribeToStore(store any, callback func()) {
	switch s := store.(type) {
	case *Store[int]:
		s.OnChange(func(_ int) { callback() })
	case *Store[string]:
		s.OnChange(func(_ string) { callback() })
	case *Store[bool]:
		s.OnChange(func(_ bool) { callback() })
	case *Store[float64]:
		s.OnChange(func(_ float64) { callback() })
	case *Store[Component]:
		s.OnChange(func(_ Component) { callback() })
	}
}

// storeToString returns the string representation of a store's current value.
func storeToString(store any) string {
	switch s := store.(type) {
	case *Store[string]:
		return s.Get()
	case *Store[int]:
		return itoa(s.Get())
	case *Store[bool]:
		if s.Get() {
			return "true"
		}
		return "false"
	case *Store[float64]:
		return ftoa(s.Get())
	}
	return ""
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
