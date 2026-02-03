//go:build js && wasm

package preveltekit

import (
	"reflect"
	"syscall/js"
)

// WasmRenderContext holds state during WASM-side HTML rendering.
// Collects bindings similar to SSR's BuildContext.
type WasmRenderContext struct {
	IDCounter
	Bindings   *WasmBindings
	Components map[string]Component
}

// WasmBindings collects bindings during WASM rendering.
// These are then applied after DOM insertion.
type WasmBindings struct {
	TextBindings     []WasmTextBinding
	InputBindings    []WasmInputBinding
	Events           []WasmEventBinding
	IfBlocks         []WasmIfBlock
	EachBlocks       []WasmEachBlock
	AttrCondBindings []WasmAttrCondBinding
}

// WasmAttrCondBinding represents a conditional attribute that needs reactive updates.
type WasmAttrCondBinding struct {
	ElementID string
	AttrName  string
	Conds     []*AttrCond // All conditions for this element/attr combo
}

type WasmTextBinding struct {
	MarkerID string
	Store    any
	IsHTML   bool
}

type WasmInputBinding struct {
	ElementID string
	Store     any
	BindType  string // "value" or "checked"
}

type WasmEventBinding struct {
	ElementID string
	Event     string
	Handler   func()
}

type WasmIfBlock struct {
	MarkerID string
	Branches []WasmIfBranch
	ElseNode []Node
	Deps     []any // Store references for reactivity
}

type WasmIfBranch struct {
	Cond     Condition
	Children []Node
}

type WasmEachBlock struct {
	MarkerID string
	ListRef  any
	Body     func(item any, index int) Node
	ElseNode []Node
}

// NewWasmRenderContext creates a new render context for WASM rendering.
func NewWasmRenderContext(prefix string) *WasmRenderContext {
	return &WasmRenderContext{
		IDCounter:  IDCounter{Prefix: prefix},
		Bindings:   &WasmBindings{},
		Components: make(map[string]Component),
	}
}

// RenderComponentWasm renders a component to HTML and collects bindings.
// Note: Caller should initialize stores and call OnCreate/OnMount BEFORE calling this.
func RenderComponentWasm(comp Component, prefix string) (string, *WasmRenderContext) {
	ctx := NewWasmRenderContext(prefix)
	ctx.Components["component"] = comp

	html := renderNodeWasm(comp.Render(), ctx)
	return html, ctx
}

// renderNodeWasm renders a Node to HTML string and collects bindings.
func renderNodeWasm(n Node, ctx *WasmRenderContext) string {
	if n == nil {
		return ""
	}

	switch node := n.(type) {
	case *TextNode:
		return escapeHTMLWasm(node.Content)

	case *HtmlNode:
		return renderHtmlNodeWasm(node, ctx)

	case *Fragment:
		var result []byte
		for _, child := range node.Children {
			result = append(result, renderNodeWasm(child, ctx)...)
		}
		return string(result)

	case *BindNode:
		return renderBindNodeWasm(node, ctx)

	case *IfNode:
		return renderIfNodeWasm(node, ctx)

	case *EachNode:
		return renderEachNodeWasm(node, ctx)

	case *BindValueNode:
		return renderBindValueNodeWasm(node, ctx)

	case *BindCheckedNode:
		return renderBindCheckedNodeWasm(node, ctx)

	case *ComponentNode:
		return renderComponentNodeWasm(node, ctx)

	case *SlotNode:
		return ""

	default:
		return ""
	}
}

// renderHtmlNodeWasm renders an HtmlNode to HTML string.
func renderHtmlNodeWasm(h *HtmlNode, ctx *WasmRenderContext) string {
	var result []byte

	for _, part := range h.Parts {
		switch v := part.(type) {
		case string:
			result = append(result, v...)
		case Node:
			result = append(result, renderNodeWasm(v, ctx)...)
		case *Store[string]:
			result = append(result, renderStoreBindingWasm(v, false, ctx)...)
		case *Store[int]:
			result = append(result, renderStoreBindingWasm(v, false, ctx)...)
		case *Store[bool]:
			result = append(result, renderStoreBindingWasm(v, false, ctx)...)
		case *Store[float64]:
			result = append(result, renderStoreBindingWasm(v, false, ctx)...)
		case *eventAttr:
			localID := ctx.NextEventID()
			fullID := ctx.FullElementID(localID)
			result = append(result, `id="`...)
			result = append(result, fullID...)
			result = append(result, `" data-event="`...)
			result = append(result, v.Event...)
			result = append(result, `"`...)
			ctx.Bindings.Events = append(ctx.Bindings.Events, WasmEventBinding{
				ElementID: fullID,
				Event:     v.Event,
				Handler:   v.Handler,
			})
		default:
			result = append(result, escapeHTMLWasm(anyToString(v))...)
		}
	}

	// Handle chained events (WithOn)
	if len(h.Events) > 0 {
		result = injectEventsWasm(result, h.Events, ctx)
	}

	// Handle chained AttrConds (AttrIf) - evaluate at render time
	if len(h.AttrConds) > 0 {
		result = injectAttrCondsWasm(result, h.AttrConds, ctx)
	}

	finalStr := string(result)
	if len(finalStr) < 300 {
	} else {
	}
	return finalStr
}

// renderStoreBindingWasm renders a store value with marker and collects binding.
func renderStoreBindingWasm(store any, isHTML bool, ctx *WasmRenderContext) string {
	localMarker := ctx.NextTextMarker()
	markerID := ctx.FullMarkerID(localMarker)

	// Get current value
	value := storeValueToString(store)

	// Collect binding
	ctx.Bindings.TextBindings = append(ctx.Bindings.TextBindings, WasmTextBinding{
		MarkerID: markerID,
		Store:    store,
		IsHTML:   isHTML,
	})

	if isHTML {
		return "<span>" + value + "</span><!--" + markerID + "-->"
	}
	return escapeHTMLWasm(value) + "<!--" + markerID + "-->"
}

// renderBindNodeWasm renders a BindNode with marker.
func renderBindNodeWasm(b *BindNode, ctx *WasmRenderContext) string {
	return renderStoreBindingWasm(b.StoreRef, b.IsHTML, ctx)
}

// renderIfNodeWasm renders an IfNode with marker and collects binding.
func renderIfNodeWasm(i *IfNode, ctx *WasmRenderContext) string {
	localMarker := ctx.NextIfMarker()
	markerID := ctx.FullMarkerID(localMarker)

	// Collect dependencies
	var deps []any
	for _, branch := range i.Branches {
		deps = append(deps, extractConditionStore(branch.Cond))
	}

	// Collect if-block binding
	wasmIf := WasmIfBlock{
		MarkerID: markerID,
		ElseNode: i.ElseNode,
		Deps:     deps,
	}
	for _, branch := range i.Branches {
		wasmIf.Branches = append(wasmIf.Branches, WasmIfBranch{
			Cond:     branch.Cond,
			Children: branch.Children,
		})
	}
	ctx.Bindings.IfBlocks = append(ctx.Bindings.IfBlocks, wasmIf)

	// Render active branch
	var activeHTML string
	for _, branch := range i.Branches {
		if branch.Cond.Eval() {
			activeHTML = renderChildrenWasm(branch.Children, ctx)
			break
		}
	}
	if activeHTML == "" && len(i.ElseNode) > 0 {
		activeHTML = renderChildrenWasm(i.ElseNode, ctx)
	}

	return "<span>" + activeHTML + "</span><!--" + markerID + "-->"
}

// renderEachNodeWasm renders an EachNode with markers and collects binding.
func renderEachNodeWasm(e *EachNode, ctx *WasmRenderContext) string {
	localMarker := ctx.NextEachMarker()
	markerID := ctx.FullMarkerID(localMarker)
	itemElementPrefix := ctx.FullElementID(localMarker)

	// Collect each-block binding
	ctx.Bindings.EachBlocks = append(ctx.Bindings.EachBlocks, WasmEachBlock{
		MarkerID: markerID,
		ListRef:  e.ListRef,
		Body:     e.Body,
		ElseNode: e.ElseNode,
	})

	var result []byte

	switch list := e.ListRef.(type) {
	case *List[string]:
		items := list.Get()
		if len(items) == 0 && len(e.ElseNode) > 0 {
			result = append(result, renderChildrenWasm(e.ElseNode, ctx)...)
		} else {
			for i, item := range items {
				itemHTML := renderNodeWasm(e.Body(item, i), ctx)
				result = append(result, `<span id="`...)
				result = append(result, itemElementPrefix...)
				result = append(result, '_')
				result = append(result, intToStr(i)...)
				result = append(result, `">`...)
				result = append(result, itemHTML...)
				result = append(result, `</span>`...)
			}
		}
	case *List[int]:
		items := list.Get()
		if len(items) == 0 && len(e.ElseNode) > 0 {
			result = append(result, renderChildrenWasm(e.ElseNode, ctx)...)
		} else {
			for i, item := range items {
				itemHTML := renderNodeWasm(e.Body(item, i), ctx)
				result = append(result, `<span id="`...)
				result = append(result, itemElementPrefix...)
				result = append(result, '_')
				result = append(result, intToStr(i)...)
				result = append(result, `">`...)
				result = append(result, itemHTML...)
				result = append(result, `</span>`...)
			}
		}
	}

	result = append(result, "<!--"...)
	result = append(result, markerID...)
	result = append(result, "-->"...)

	return string(result)
}

// renderBindValueNodeWasm renders a two-way bound input and collects binding.
func renderBindValueNodeWasm(b *BindValueNode, ctx *WasmRenderContext) string {
	localID := ctx.NextBindID()
	fullID := ctx.FullElementID(localID)

	// Collect binding
	ctx.Bindings.InputBindings = append(ctx.Bindings.InputBindings, WasmInputBinding{
		ElementID: fullID,
		Store:     b.Store,
		BindType:  "value",
	})

	// Get current value
	value := storeValueToString(b.Store)

	return injectAttrsWasm(b.HTML, `id="`+fullID+`" value="`+escapeAttrWasm(value)+`"`)
}

// renderBindCheckedNodeWasm renders a two-way bound checkbox and collects binding.
func renderBindCheckedNodeWasm(b *BindCheckedNode, ctx *WasmRenderContext) string {
	localID := ctx.NextBindID()
	fullID := ctx.FullElementID(localID)

	// Collect binding
	ctx.Bindings.InputBindings = append(ctx.Bindings.InputBindings, WasmInputBinding{
		ElementID: fullID,
		Store:     b.Store,
		BindType:  "checked",
	})

	checked := ""
	if s, ok := b.Store.(*Store[bool]); ok && s.Get() {
		checked = " checked"
	}

	return injectAttrsWasm(b.HTML, `id="`+fullID+`"`+checked)
}

// renderComponentNodeWasm renders a nested component.
func renderComponentNodeWasm(c *ComponentNode, ctx *WasmRenderContext) string {
	comp, ok := c.Instance.(Component)
	if !ok {
		return ""
	}

	compMarker := ctx.NextCompMarker()
	fullCompPrefix := ctx.FullElementID(compMarker)

	// Create child context
	childCtx := NewWasmRenderContext(fullCompPrefix)
	childCtx.Components["component"] = comp

	if oc, ok := comp.(HasOnCreate); ok {
		oc.OnCreate()
	}

	html := renderNodeWasm(comp.Render(), childCtx)

	// Merge child bindings into parent
	mergeWasmBindings(ctx.Bindings, childCtx.Bindings)

	return html
}

// renderChildrenWasm renders a slice of nodes.
func renderChildrenWasm(nodes []Node, ctx *WasmRenderContext) string {
	var result []byte
	for _, n := range nodes {
		result = append(result, renderNodeWasm(n, ctx)...)
	}
	return string(result)
}

// mergeWasmBindings merges child bindings into parent.
func mergeWasmBindings(parent, child *WasmBindings) {
	parent.TextBindings = append(parent.TextBindings, child.TextBindings...)
	parent.InputBindings = append(parent.InputBindings, child.InputBindings...)
	parent.Events = append(parent.Events, child.Events...)
	parent.IfBlocks = append(parent.IfBlocks, child.IfBlocks...)
	parent.EachBlocks = append(parent.EachBlocks, child.EachBlocks...)
	parent.AttrCondBindings = append(parent.AttrCondBindings, child.AttrCondBindings...)
}

// storeValueToString gets the current value of a store as string.
func storeValueToString(store any) string {
	switch s := store.(type) {
	case *Store[string]:
		return s.Get()
	case *Store[int]:
		return intToStr(s.Get())
	case *Store[bool]:
		if s.Get() {
			return "true"
		}
		return "false"
	case *Store[float64]:
		return floatToStr(s.Get())
	}
	return ""
}

// extractConditionStore extracts the store from a condition.
func extractConditionStore(cond Condition) any {
	switch c := cond.(type) {
	case *StoreCondition:
		return c.Store
	case *BoolCondition:
		return c.Store
	}
	return nil
}

// injectAttrsWasm injects attributes into an HTML element string.
func injectAttrsWasm(html, attrs string) string {
	for i := 0; i < len(html); i++ {
		if html[i] == '>' {
			if i > 0 && html[i-1] == '/' {
				return html[:i-1] + " " + attrs + " />" + html[i+1:]
			}
			return html[:i] + " " + attrs + html[i:]
		}
	}
	return html + " " + attrs
}

// injectEventsWasm injects event bindings into HTML.
func injectEventsWasm(html []byte, events []*HtmlEvent, ctx *WasmRenderContext) []byte {
	if len(events) == 0 {
		return html
	}

	localID := ctx.NextEventID()
	fullID := ctx.FullElementID(localID)

	var eventNames []byte
	for i, ev := range events {
		if i > 0 {
			eventNames = append(eventNames, ',')
		}
		eventNames = append(eventNames, ev.Event...)
		ctx.Bindings.Events = append(ctx.Bindings.Events, WasmEventBinding{
			ElementID: fullID,
			Event:     ev.Event,
			Handler:   ev.Handler,
		})
	}

	for i := 0; i < len(html); i++ {
		if html[i] == '>' {
			inject := ` id="` + fullID + `" data-event="` + string(eventNames) + `"`
			// Build result by concatenating strings to avoid slice mutation bugs
			if i > 0 && html[i-1] == '/' {
				return []byte(string(html[:i-1]) + inject + string(html[i-1:]))
			}
			return []byte(string(html[:i]) + inject + string(html[i:]))
		}
	}
	return html
}

// injectAttrCondsWasm injects conditional attributes and collects bindings for reactivity.
// For class attributes, merges with existing class values instead of creating duplicates.
func injectAttrCondsWasm(html []byte, attrConds []*AttrCond, ctx *WasmRenderContext) []byte {
	if len(attrConds) == 0 {
		return html
	}

	// Generate element ID for this element (needed for reactive updates)
	localID := ctx.NextClassID()
	fullID := ctx.FullElementID(localID)

	// Collect values by attribute name (to handle multiple AttrIf for same attr)
	attrValues := make(map[string][]string)
	for _, ac := range attrConds {
		var value string
		if ac.Cond.Eval() {
			if s, ok := ac.TrueValue.(string); ok {
				value = s
			}
		} else {
			if s, ok := ac.FalseValue.(string); ok {
				value = s
			}
		}
		if value != "" {
			attrValues[ac.Name] = append(attrValues[ac.Name], value)
		}
	}

	htmlStr := string(html)

	// Find the first tag's end position
	tagEnd := -1
	for i := 0; i < len(htmlStr); i++ {
		if htmlStr[i] == '>' {
			tagEnd = i
			break
		}
	}
	if tagEnd == -1 {
		return html
	}

	openingTag := htmlStr[:tagEnd]
	rest := htmlStr[tagEnd:]

	restPreview := rest
	if len(restPreview) > 50 {
		restPreview = restPreview[:50]
	}

	// Check if element already has an ID
	hasID := indexOf(openingTag, ` id="`) != -1

	// Inject ID if not present (needed for reactive updates)
	if !hasID {
		openingTag = openingTag + ` id="` + fullID + `"`
	}

	// Handle class attribute specially - merge with existing
	if classes, ok := attrValues["class"]; ok && len(classes) > 0 {
		classIdx := indexOf(openingTag, `class="`)
		if classIdx != -1 {
			// Find existing class value
			classStart := classIdx + 7
			classEndRel := indexOf(openingTag[classStart:], `"`)
			if classEndRel != -1 {
				classEnd := classStart + classEndRel
				existingClasses := openingTag[classStart:classEnd]
				// Merge: existing + new classes
				mergedClasses := existingClasses
				for _, c := range classes {
					if c != "" {
						mergedClasses += " " + c
					}
				}
				// Rebuild opening tag with merged class
				openingTag = openingTag[:classIdx] + `class="` + mergedClasses + `"` + openingTag[classEnd+1:]
			}
		} else {
			// No existing class, add new one
			openingTag = openingTag + ` class="` + joinStrings(classes, " ") + `"`
		}
		delete(attrValues, "class")
	}

	// Handle other attributes
	for name, values := range attrValues {
		if len(values) > 0 {
			openingTag = openingTag + ` ` + name + `="` + joinStrings(values, " ") + `"`
		}
	}

	// Collect binding for reactive updates (use the ID we injected or existing)
	elementID := fullID
	if hasID {
		// Extract existing ID
		idIdx := indexOf(openingTag, ` id="`)
		if idIdx != -1 {
			idStart := idIdx + 5
			idEndRel := indexOf(openingTag[idStart:], `"`)
			if idEndRel != -1 {
				elementID = openingTag[idStart : idStart+idEndRel]
			}
		}
	}

	// Group conditions by attribute name for the binding
	attrCondsByName := make(map[string][]*AttrCond)
	for _, ac := range attrConds {
		attrCondsByName[ac.Name] = append(attrCondsByName[ac.Name], ac)
	}

	// Collect one binding per attribute
	for attrName, conds := range attrCondsByName {
		ctx.Bindings.AttrCondBindings = append(ctx.Bindings.AttrCondBindings, WasmAttrCondBinding{
			ElementID: elementID,
			AttrName:  attrName,
			Conds:     conds,
		})
	}

	return []byte(openingTag + rest)
}

// indexOf finds substring in string, returns -1 if not found
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// joinStrings joins strings with separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// escapeAttrWasm escapes attribute values.
func escapeAttrWasm(s string) string {
	var result []byte
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '&':
			result = append(result, "&amp;"...)
		case '"':
			result = append(result, "&quot;"...)
		default:
			result = append(result, s[i])
		}
	}
	return string(result)
}

// ApplyWasmBindings applies collected bindings to the DOM after HTML insertion.
func ApplyWasmBindings(bindings *WasmBindings, cleanup *Cleanup) {
	// Clear bound markers first so text bindings can rebind to new DOM
	for _, tb := range bindings.TextBindings {
		ClearBoundMarker(tb.MarkerID)
	}

	// Apply text bindings
	for _, tb := range bindings.TextBindings {
		applyTextBindingWasm(tb)
	}

	// Apply input bindings
	for _, ib := range bindings.InputBindings {
		applyInputBindingWasm(ib, cleanup)
	}

	// Apply event bindings
	for _, ev := range bindings.Events {
		applyEventBindingWasm(ev, cleanup)
	}

	// Apply if-block bindings
	for _, ifb := range bindings.IfBlocks {
		applyIfBlockWasm(ifb)
	}

	// Apply each-block bindings
	for _, eb := range bindings.EachBlocks {
		applyEachBlockWasm(eb)
	}

	// Apply attr cond bindings (reactive class/attr changes)
	for _, acb := range bindings.AttrCondBindings {
		applyAttrCondBindingWasm(acb)
	}
}

// applyTextBindingWasm sets up a text binding for reactivity.
func applyTextBindingWasm(tb WasmTextBinding) {
	switch s := tb.Store.(type) {
	case *Store[string]:
		if tb.IsHTML {
			BindHTML(tb.MarkerID, s)
		} else {
			BindText(tb.MarkerID, s)
		}
	case *Store[int]:
		if tb.IsHTML {
			BindHTML(tb.MarkerID, s)
		} else {
			BindText(tb.MarkerID, s)
		}
	case *Store[bool]:
		if tb.IsHTML {
			BindHTML(tb.MarkerID, s)
		} else {
			BindText(tb.MarkerID, s)
		}
	case *Store[float64]:
		if tb.IsHTML {
			BindHTML(tb.MarkerID, s)
		} else {
			BindText(tb.MarkerID, s)
		}
	}
}

// applyInputBindingWasm sets up an input binding for two-way data flow.
func applyInputBindingWasm(ib WasmInputBinding, cleanup *Cleanup) {
	switch s := ib.Store.(type) {
	case *Store[string]:
		if ib.BindType == "value" {
			cleanup.Add(BindInput(ib.ElementID, s))
		}
	case *Store[int]:
		if ib.BindType == "value" {
			cleanup.Add(BindInputInt(ib.ElementID, s))
		}
	case *Store[bool]:
		if ib.BindType == "checked" {
			cleanup.Add(BindCheckbox(ib.ElementID, s))
		}
	}
}

// applyEventBindingWasm sets up an event handler.
func applyEventBindingWasm(ev WasmEventBinding, cleanup *Cleanup) {
	el := GetEl(ev.ElementID)
	if !ok(el) {
		return
	}
	fn := js.FuncOf(func(this js.Value, args []js.Value) any {
		if ev.Handler != nil {
			ev.Handler()
		}
		return nil
	})
	el.Call("addEventListener", ev.Event, fn)
	cleanup.Add(fn)
}

// applyIfBlockWasm sets up an if-block for reactivity.
func applyIfBlockWasm(ifb WasmIfBlock) {

	// Track setup to avoid duplicates
	if setupIfBlocks[ifb.MarkerID] {
		return
	}
	setupIfBlocks[ifb.MarkerID] = true

	currentEl := FindExistingIfContent(ifb.MarkerID)
	currentCleanup := &Cleanup{}

	updateIfBlock := func() {
		// Determine active branch
		var activeNodes []Node
		found := false
		for _, branch := range ifb.Branches {
			if branch.Cond.Eval() {
				activeNodes = branch.Children
				found = true
				break
			}
		}
		if !found {
			activeNodes = ifb.ElseNode
		}

		// Render active branch
		ctx := NewWasmRenderContext("")
		activeHTML := renderChildrenWasm(activeNodes, ctx)

		// Replace content
		currentEl = FindExistingIfContent(ifb.MarkerID)
		currentEl = ReplaceContent(ifb.MarkerID, currentEl, activeHTML)

		// Apply nested bindings
		currentCleanup.Release()
		currentCleanup = &Cleanup{}
		ApplyWasmBindings(ctx.Bindings, currentCleanup)
	}

	// Subscribe to store changes
	for _, dep := range ifb.Deps {
		if dep != nil {
			subscribeToStore(dep, func() {
				updateIfBlock()
			})
		}
	}
}

// applyEachBlockWasm sets up an each-block for reactivity.
func applyEachBlockWasm(eb WasmEachBlock) {
	if setupEachBlocks[eb.MarkerID] {
		return
	}
	setupEachBlocks[eb.MarkerID] = true

	marker := FindComment(eb.MarkerID)
	if marker.IsNull() {
		return
	}

	parent := marker.Get("parentNode")

	// Extract item prefix from marker
	itemIDPrefix := eb.MarkerID
	if len(itemIDPrefix) > 0 {
		// Remove trailing marker type (e0 -> e)
		for i := len(itemIDPrefix) - 1; i >= 0; i-- {
			if itemIDPrefix[i] == '_' {
				itemIDPrefix = itemIDPrefix[:i+1]
				break
			}
		}
	}

	switch list := eb.ListRef.(type) {
	case *List[string]:
		list.OnChange(func(items []string) {
			var html string
			for i, item := range items {
				ctx := NewWasmRenderContext("")
				itemHTML := renderNodeWasm(eb.Body(item, i), ctx)
				html += `<span id="` + itemIDPrefix + intToStr(i) + `">` + itemHTML + `</span>`
			}
			if !parent.IsNull() && parent.Get("nodeType").Int() == 1 {
				parent.Set("innerHTML", html)
			}
		})
	case *List[int]:
		list.OnChange(func(items []int) {
			var html string
			for i, item := range items {
				ctx := NewWasmRenderContext("")
				itemHTML := renderNodeWasm(eb.Body(item, i), ctx)
				html += `<span id="` + itemIDPrefix + intToStr(i) + `">` + itemHTML + `</span>`
			}
			if !parent.IsNull() && parent.Get("nodeType").Int() == 1 {
				parent.Set("innerHTML", html)
			}
		})
	}
}

// applyAttrCondBindingWasm sets up reactive attribute/class updates.
func applyAttrCondBindingWasm(acb WasmAttrCondBinding) {
	el := GetEl(acb.ElementID)
	if !ok(el) {
		return
	}

	// Function to evaluate conditions and update the attribute
	updateAttr := func() {
		if acb.AttrName == "class" {
			// For class attributes, toggle each condition's true/false values
			for _, cond := range acb.Conds {
				active := cond.Cond.Eval()
				classList := el.Get("classList")

				// Remove false value if present, add true value if active
				if trueVal, ok := cond.TrueValue.(string); ok && trueVal != "" {
					if active {
						classList.Call("add", trueVal)
					} else {
						classList.Call("remove", trueVal)
					}
				}
				if falseVal, ok := cond.FalseValue.(string); ok && falseVal != "" {
					if !active {
						classList.Call("add", falseVal)
					} else {
						classList.Call("remove", falseVal)
					}
				}
			}
		} else {
			// For other attributes, collect all active values
			var values []string
			for _, cond := range acb.Conds {
				var value string
				if cond.Cond.Eval() {
					if s, ok := cond.TrueValue.(string); ok {
						value = s
					}
				} else {
					if s, ok := cond.FalseValue.(string); ok {
						value = s
					}
				}
				if value != "" {
					values = append(values, value)
				}
			}
			if len(values) > 0 {
				el.Call("setAttribute", acb.AttrName, joinStrings(values, " "))
			} else {
				el.Call("removeAttribute", acb.AttrName)
			}
		}
	}

	// Subscribe to all store changes
	for _, cond := range acb.Conds {
		store := extractConditionStore(cond.Cond)
		if store != nil {
			subscribeToStore(store, updateAttr)
		}
	}
}

// buildComponentStoreMap builds a map from store addresses to their IDs.
func buildComponentStoreMap(comp Component, prefix string) map[uintptr]string {
	storeMap := make(map[uintptr]string)
	rv := reflect.ValueOf(comp).Elem()
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		field := rv.Field(i)
		fieldName := rt.Field(i).Name

		if field.Kind() == reflect.Ptr && !field.IsNil() {
			addr := field.Pointer()
			storeID := prefix + "." + fieldName
			storeMap[addr] = storeID
		}
	}
	return storeMap
}
