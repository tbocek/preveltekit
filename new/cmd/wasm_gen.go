package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Scope represents a variable scope for field resolution.
// Scopes chain from child to parent, enabling lexical field lookup.
type Scope struct {
	VarName    string            // Go variable name: "component", "comp1", "comp1_comp0"
	FieldTypes map[string]string // Field name -> type
	Parent     *Scope            // Parent scope (nil for root)
}

// Resolve finds a field in the scope chain, returning the variable reference and type.
// Inner scopes shadow outer scopes.
func (s *Scope) Resolve(fieldName string) (varRef string, fieldType string, found bool) {
	for scope := s; scope != nil; scope = scope.Parent {
		if t, ok := scope.FieldTypes[fieldName]; ok {
			return scope.VarName + "." + fieldName, t, true
		}
	}
	return "", "", false
}

// WiringContext contains all information needed to generate component wiring code.
type WiringContext struct {
	// Component identity
	ID         string           // Full prefixed ID: "comp1_comp0"
	Name       string           // Component type name: "Badge"
	Definition *component       // Parsed component definition
	Binding    componentBinding // How this component is used (props, events, slot)

	// Scope chain for field resolution
	Scope *Scope

	// Parent context (nil for root)
	Parent *WiringContext

	// All available components
	AllComponents map[string]*component

	// Code generation state
	Indent        string // Current indentation
	SkipCreate    bool   // Instance already exists (for if-block components)
	InsideIfBlock bool   // Component is inside an if-block
}

// ChildContext creates a new context for a child component.
func (ctx *WiringContext) ChildContext(childBinding componentBinding, childDef *component) *WiringContext {
	childID := ctx.prefixID(childBinding.elementID)
	childScope := &Scope{
		VarName:    childID,
		FieldTypes: buildFieldTypes(childDef),
		Parent:     ctx.Scope,
	}
	return &WiringContext{
		ID:            childID,
		Name:          childBinding.name,
		Definition:    childDef,
		Binding:       childBinding,
		Scope:         childScope,
		Parent:        ctx,
		AllComponents: ctx.AllComponents,
		Indent:        ctx.Indent + "\t",
		SkipCreate:    false,
		InsideIfBlock: false,
	}
}

// prefixID returns the full prefixed ID for a local ID.
func (ctx *WiringContext) prefixID(localID string) string {
	if ctx.ID == "" {
		return localID
	}
	return ctx.ID + "_" + localID
}

// ParentStoreRef returns the Go code to access a parent's store field.
func (ctx *WiringContext) ParentStoreRef(field string) string {
	if ctx.Parent == nil {
		return "component." + field
	}
	return ctx.Parent.Scope.VarName + "." + field
}

// HTMLConstantTracker tracks which HTML constants have been generated.
type HTMLConstantTracker struct {
	CSS  map[string]bool // Component name -> generated
	HTML map[string]bool // Prefixed ID -> generated
}

func NewHTMLConstantTracker() *HTMLConstantTracker {
	return &HTMLConstantTracker{
		CSS:  make(map[string]bool),
		HTML: make(map[string]bool),
	}
}

// generateMain generates the main WASM entry point.
func generateMain(comp *component, tmpl string, bindings templateBindings, childComponents map[string]*component) string {
	var sb strings.Builder
	fieldTypes := buildFieldTypes(comp)

	// Check what imports we need
	needsStrconv := needsStrconvImport(bindings, childComponents, fieldTypes)
	needsStrings := needsStringsImport(bindings, childComponents)

	// Write imports
	sb.WriteString("//go:build js && wasm\n\npackage main\n\nimport (\n")
	sb.WriteString("\t\"preveltekit\"\n")
	sb.WriteString("\t\"syscall/js\"\n")
	if needsStrconv {
		sb.WriteString("\t\"strconv\"\n")
	}
	if needsStrings {
		sb.WriteString("\t\"strings\"\n")
	}
	sb.WriteString(")\n\n")

	sb.WriteString("var document = preveltekit.Document\n\n")

	// Generate CSS and HTML constants
	tracker := NewHTMLConstantTracker()
	generateConstants(&sb, bindings, childComponents, "", tracker)

	// Main function
	sb.WriteString("\nfunc main() {\n\tcomponent := &" + comp.name + "{\n")
	generateFieldInit(&sb, comp.fields, "\t\t")
	sb.WriteString("\t}\n\n")

	// Create root context
	rootScope := &Scope{
		VarName:    "component",
		FieldTypes: fieldTypes,
		Parent:     nil,
	}
	rootCtx := &WiringContext{
		ID:            "",
		Name:          comp.name,
		Definition:    comp,
		Scope:         rootScope,
		Parent:        nil,
		AllComponents: childComponents,
		Indent:        "\t",
	}

	// Wire up the root component's bindings
	generateBindingsWiring(&sb, bindings, rootCtx)

	// Wire up child components
	componentsInIfBlocks := findComponentsInIfBlocks(bindings.ifBlocks)
	for _, compBinding := range bindings.components {
		if !componentsInIfBlocks[compBinding.elementID] {
			childDef := childComponents[compBinding.name]
			if childDef != nil {
				childCtx := rootCtx.ChildContext(compBinding, childDef)
				generateComponentWiring(&sb, childCtx)
			}
		}
	}

	if comp.hasOnMount {
		sb.WriteString("\tcomponent.OnMount()\n\n")
	}

	sb.WriteString("\tselect {}\n}\n")
	return sb.String()
}

// generateConstants generates all CSS and HTML constants needed for the component tree.
func generateConstants(sb *strings.Builder, bindings templateBindings, childComponents map[string]*component, prefix string, tracker *HTMLConstantTracker) {
	// Generate CSS for each child component type
	for _, compBinding := range bindings.components {
		childDef := childComponents[compBinding.name]
		if childDef == nil {
			continue
		}
		if childDef.style != "" && !tracker.CSS[compBinding.name] {
			tracker.CSS[compBinding.name] = true
			fmt.Fprintf(sb, "const %sCSS = `%s`\n", strings.ToLower(compBinding.name), minifyCSS(childDef.style))
		}
	}

	// Find components in if-blocks - these need HTML constants
	componentsInIfBlocks := findComponentsInIfBlocks(bindings.ifBlocks)

	// Generate HTML constants for components in if-blocks
	for _, compBinding := range bindings.components {
		if !componentsInIfBlocks[compBinding.elementID] {
			continue
		}
		childDef := childComponents[compBinding.name]
		if childDef == nil {
			continue
		}
		fullID := prefixIDStr(prefix, compBinding.elementID)
		generateHTMLConstant(sb, compBinding, childDef, fullID, childComponents, tracker)
	}

	// Also process components that are not in if-blocks but may have nested components
	for _, compBinding := range bindings.components {
		if componentsInIfBlocks[compBinding.elementID] {
			continue // Already processed above
		}
		childDef := childComponents[compBinding.name]
		if childDef == nil {
			continue
		}
		// Check for nested components in child's template
		childTmpl := strings.ReplaceAll(childDef.template, "<slot/>", compBinding.children)
		_, childBindings := parseTemplate(childTmpl)
		fullID := prefixIDStr(prefix, compBinding.elementID)
		generateConstants(sb, childBindings, childComponents, fullID, tracker)
	}
}

// generateHTMLConstant generates an HTML constant for a component and all its nested components.
func generateHTMLConstant(sb *strings.Builder, compBinding componentBinding, compDef *component, fullID string, childComponents map[string]*component, tracker *HTMLConstantTracker) {
	if tracker.HTML[fullID] {
		return
	}
	tracker.HTML[fullID] = true

	// Generate CSS if not already done
	if compDef.style != "" && !tracker.CSS[compBinding.name] {
		tracker.CSS[compBinding.name] = true
		fmt.Fprintf(sb, "const %sCSS = `%s`\n", strings.ToLower(compBinding.name), minifyCSS(compDef.style))
	}

	// Process template with slot content
	childTmpl := strings.ReplaceAll(compDef.template, "<slot/>", compBinding.children)
	processedTmpl, childBindings := parseTemplate(childTmpl)

	// Prefix all IDs in the template
	processedTmpl = prefixAllBindingIDs(fullID, processedTmpl, &childBindings)

	// Inject component ID into root element
	processedTmpl = injectIDIntoFirstTag(processedTmpl, fullID)

	fmt.Fprintf(sb, "const %sHTML = %s\n", fullID, escapeForGoString(minifyHTML(processedTmpl)))

	// Recursively generate constants for nested components
	for _, nestedBinding := range childBindings.components {
		nestedDef := childComponents[nestedBinding.name]
		if nestedDef == nil {
			continue
		}
		// Manually prefix the ID since prefixAllBindingIDs no longer mutates component IDs
		nestedFullID := fullID + "_" + nestedBinding.elementID
		generateHTMLConstant(sb, nestedBinding, nestedDef, nestedFullID, childComponents, tracker)
	}
}

// prefixAllBindingIDs prefixes all binding IDs in the template and updates the bindings struct.
func prefixAllBindingIDs(prefix string, html string, bindings *templateBindings) string {
	// Prefix expression IDs
	for i := range bindings.expressions {
		oldID := bindings.expressions[i].elementID
		newID := prefix + "_" + oldID
		html = strings.ReplaceAll(html, `id="`+oldID+`"`, `id="`+newID+`"`)
		bindings.expressions[i].elementID = newID
	}
	// Prefix event IDs
	for i := range bindings.events {
		oldID := bindings.events[i].elementID
		newID := prefix + "_" + oldID
		html = strings.ReplaceAll(html, `id="`+oldID+`"`, `id="`+newID+`"`)
		bindings.events[i].elementID = newID
	}
	// Prefix input binding IDs
	for i := range bindings.bindings {
		oldID := bindings.bindings[i].elementID
		newID := prefix + "_" + oldID
		html = strings.ReplaceAll(html, `id="`+oldID+`"`, `id="`+newID+`"`)
		bindings.bindings[i].elementID = newID
	}
	// Prefix class binding IDs
	for i := range bindings.classBindings {
		oldID := bindings.classBindings[i].elementID
		newID := prefix + "_" + oldID
		html = strings.ReplaceAll(html, `id="`+oldID+`"`, `id="`+newID+`"`)
		bindings.classBindings[i].elementID = newID
	}
	// Prefix each block IDs
	for i := range bindings.eachBlocks {
		oldID := bindings.eachBlocks[i].elementID
		newID := prefix + "_" + oldID
		html = strings.ReplaceAll(html, `id="`+oldID+`_anchor"`, `id="`+newID+`_anchor"`)
		html = strings.ReplaceAll(html, `id="`+oldID+`_else"`, `id="`+newID+`_else"`)
		bindings.eachBlocks[i].elementID = newID
	}
	// Prefix attribute binding IDs
	for i := range bindings.attrBindings {
		oldID := bindings.attrBindings[i].elementID
		newID := prefix + "_" + oldID
		html = strings.ReplaceAll(html, `data-attrbind="`+oldID+`"`, `data-attrbind="`+newID+`"`)
		bindings.attrBindings[i].elementID = newID
	}
	// Prefix if block IDs and nested bindings
	for i := range bindings.ifBlocks {
		oldID := bindings.ifBlocks[i].elementID
		newID := prefix + "_" + oldID
		html = strings.ReplaceAll(html, `id="`+oldID+`_anchor"`, `id="`+newID+`_anchor"`)
		bindings.ifBlocks[i].elementID = newID

		// Prefix bindings inside branches
		for j := range bindings.ifBlocks[i].branches {
			for k := range bindings.ifBlocks[i].branches[j].eachBlocks {
				oldEachID := bindings.ifBlocks[i].branches[j].eachBlocks[k].elementID
				newEachID := prefix + "_" + oldEachID
				bindings.ifBlocks[i].branches[j].html = strings.ReplaceAll(
					bindings.ifBlocks[i].branches[j].html,
					`id="`+oldEachID+`_anchor"`, `id="`+newEachID+`_anchor"`)
				bindings.ifBlocks[i].branches[j].eachBlocks[k].elementID = newEachID
			}
			for k := range bindings.ifBlocks[i].branches[j].classBindings {
				oldClassID := bindings.ifBlocks[i].branches[j].classBindings[k].elementID
				newClassID := prefix + "_" + oldClassID
				bindings.ifBlocks[i].branches[j].html = strings.ReplaceAll(
					bindings.ifBlocks[i].branches[j].html,
					`id="`+oldClassID+`"`, `id="`+newClassID+`"`)
				bindings.ifBlocks[i].branches[j].classBindings[k].elementID = newClassID
			}
		}
	}
	// Prefix component placeholders in HTML only (don't mutate the binding IDs)
	// The component IDs will be prefixed when we create child contexts
	for i := range bindings.components {
		oldID := bindings.components[i].elementID
		newID := prefix + "_" + oldID
		html = strings.ReplaceAll(html, "<!--"+oldID+"-->", "<!--"+newID+"-->")
		// Note: We intentionally do NOT mutate bindings.components[i].elementID here
		// because ChildContext will prefix it when creating nested contexts
	}
	return html
}

// generateBindingsWiring generates wiring code for all bindings in a template.
func generateBindingsWiring(sb *strings.Builder, bindings templateBindings, ctx *WiringContext) {
	indent := ctx.Indent
	fieldTypes := ctx.Scope.FieldTypes
	varName := ctx.Scope.VarName

	// Expression bindings
	for _, expr := range bindings.expressions {
		fullID := ctx.prefixID(expr.elementID)
		// Resolve the field in scope chain
		varRef, _, found := ctx.Scope.Resolve(expr.fieldName)
		if !found {
			// Check parent scope for slot content
			if ctx.Parent != nil {
				varRef, _, found = ctx.Parent.Scope.Resolve(expr.fieldName)
			}
		}
		if !found {
			varRef = varName + "." + expr.fieldName
		}
		if expr.isHTML {
			valueType := fieldTypes[expr.fieldName]
			if valueType == "" {
				valueType = "string"
			}
			jsConv := toJS(valueType, "v")
			jsConvInit := toJS(valueType, varRef+".Get()")
			fmt.Fprintf(sb, "%s%s := preveltekit.GetEl(\"%s\")\n", indent, fullID, fullID)
			fmt.Fprintf(sb, "%s%s.OnChange(func(v %s) { if !%s.IsUndefined() && !%s.IsNull() { %s.Set(\"innerHTML\", %s) } })\n",
				indent, varRef, valueType, fullID, fullID, fullID, jsConv)
			fmt.Fprintf(sb, "%sif !%s.IsUndefined() && !%s.IsNull() { %s.Set(\"innerHTML\", %s) }\n",
				indent, fullID, fullID, fullID, jsConvInit)
		} else {
			fmt.Fprintf(sb, "%spreveltekit.Bind(\"%s\", %s)\n", indent, fullID, varRef)
		}
	}

	// Event bindings
	for _, evt := range bindings.events {
		fullID := ctx.prefixID(evt.elementID)
		callArgs := transformEventArgs(evt.args, ctx.Scope)

		var modifierCode string
		hasOnce := false
		for _, mod := range evt.modifiers {
			switch mod {
			case "preventDefault":
				modifierCode += indent + "\t\t\targs[0].Call(\"preventDefault\")\n"
			case "stopPropagation":
				modifierCode += indent + "\t\t\targs[0].Call(\"stopPropagation\")\n"
			case "once":
				hasOnce = true
			}
		}

		if modifierCode != "" || hasOnce {
			onceOpt := ""
			if hasOnce {
				onceOpt = ", map[string]interface{}{\"once\": true}"
			}
			fmt.Fprintf(sb, "%sdocument.Call(\"getElementById\", \"%s\").Call(\"addEventListener\", \"%s\",\n", indent, fullID, evt.event)
			fmt.Fprintf(sb, "%s\tjs.FuncOf(func(this js.Value, args []js.Value) any {\n%s%s\t\t%s.%s(%s)\n%s\t\treturn nil\n%s\t})%s)\n\n",
				indent, modifierCode, indent, varName, evt.methodName, callArgs, indent, indent, onceOpt)
		} else {
			fmt.Fprintf(sb, "%spreveltekit.On(preveltekit.GetEl(\"%s\"), \"%s\", func() { %s.%s(%s) })\n",
				indent, fullID, evt.event, varName, evt.methodName, callArgs)
		}
	}

	// If blocks
	generateIfBlocksWiring(sb, bindings.ifBlocks, bindings.components, ctx)

	// Input bindings
	for _, bind := range bindings.bindings {
		fullID := ctx.prefixID(bind.elementID)
		valueType := fieldTypes[bind.fieldName]
		if valueType == "" {
			valueType = "string"
		}
		varRef := varName + "." + bind.fieldName

		if bind.bindType == "checked" {
			fmt.Fprintf(sb, "%spreveltekit.BindCheckbox(\"%s\", %s)\n", indent, fullID, varRef)
		} else {
			switch valueType {
			case "int":
				fmt.Fprintf(sb, "%spreveltekit.BindInputInt(\"%s\", %s)\n", indent, fullID, varRef)
			default:
				fmt.Fprintf(sb, "%spreveltekit.BindInput(\"%s\", %s)\n", indent, fullID, varRef)
			}
		}
	}

	// Class bindings
	generateClassBindingsWiring(sb, bindings.classBindings, ctx)

	// Attribute bindings
	generateAttrBindingsWiring(sb, bindings.attrBindings, ctx)

	// Each blocks
	generateEachBlocksWiring(sb, bindings.eachBlocks, ctx)
}

// generateClassBindingsWiring generates wiring for class bindings, grouping by element.
func generateClassBindingsWiring(sb *strings.Builder, classBindings []classBinding, ctx *WiringContext) {
	indent := ctx.Indent
	fieldTypes := ctx.Scope.FieldTypes
	varName := ctx.Scope.VarName

	// Group by element ID
	byElement := make(map[string][]classBinding)
	elementOrder := []string{}
	for _, cb := range classBindings {
		if _, exists := byElement[cb.elementID]; !exists {
			elementOrder = append(elementOrder, cb.elementID)
		}
		byElement[cb.elementID] = append(byElement[cb.elementID], cb)
	}

	for _, fullID := range elementOrder {
		bindings := byElement[fullID]
		fmt.Fprintf(sb, "%s%s := document.Call(\"getElementById\", \"%s\")\n", indent, fullID, fullID)

		// Collect dependencies
		allDeps := make(map[string]bool)
		for _, cb := range bindings {
			deps := extractPascalCaseWords(cb.condition)
			for _, dep := range deps {
				allDeps[dep] = true
			}
		}

		// Create update function
		fmt.Fprintf(sb, "%supdate%s := func() {\n", indent, fullID)
		for _, cb := range bindings {
			cond := transformCondition(cb.condition, fieldTypes, varName)
			fmt.Fprintf(sb, "%s\tpreveltekit.ToggleClass(%s, \"%s\", %s)\n",
				indent, fullID, cb.className, cond)
		}
		fmt.Fprintf(sb, "%s}\n", indent)

		// Register OnChange
		for dep := range allDeps {
			if fieldType, ok := fieldTypes[dep]; ok {
				fmt.Fprintf(sb, "%s%s.%s.OnChange(func(_ %s) { update%s() })\n", indent, varName, dep, fieldType, fullID)
			}
		}
		fmt.Fprintf(sb, "%supdate%s()\n\n", indent, fullID)
	}
}

// generateAttrBindingsWiring generates wiring for attribute bindings.
func generateAttrBindingsWiring(sb *strings.Builder, attrBindings []attrBinding, ctx *WiringContext) {
	indent := ctx.Indent
	fieldTypes := ctx.Scope.FieldTypes
	varName := ctx.Scope.VarName

	for _, ab := range attrBindings {
		fullID := ctx.prefixID(ab.elementID)
		fmt.Fprintf(sb, "%sattr%s := document.Call(\"querySelector\", \"[data-attrbind=\\\"%s\\\"]\")\n", indent, fullID, fullID)
		fmt.Fprintf(sb, "%supdateAttr%s := func() {\n%s\tval := %s\n", indent, fullID, indent, escapeForGoString(ab.template))
		for _, field := range ab.fields {
			fmt.Fprintf(sb, "%s\tval = strings.ReplaceAll(val, \"{%s}\", %s)\n", indent, field, toJS(fieldTypes[field], varName+"."+field+".Get()"))
		}
		fmt.Fprintf(sb, "%s\tattr%s.Call(\"setAttribute\", \"%s\", val)\n%s}\n", indent, fullID, ab.attrName, indent)
		for _, field := range ab.fields {
			fmt.Fprintf(sb, "%s%s.%s.OnChange(func(_ %s) { updateAttr%s() })\n", indent, varName, field, fieldTypes[field], fullID)
		}
		fmt.Fprintf(sb, "%supdateAttr%s()\n\n", indent, fullID)
	}
}

// generateEachBlocksWiring generates wiring for each blocks.
func generateEachBlocksWiring(sb *strings.Builder, eachBlocks []eachBinding, ctx *WiringContext) {
	indent := ctx.Indent
	fieldTypes := ctx.Scope.FieldTypes
	varName := ctx.Scope.VarName

	for _, each := range eachBlocks {
		fullID := ctx.prefixID(each.elementID)
		bodyHTML := strings.ReplaceAll(each.bodyHTML, "{"+each.itemVar+"}", `<span class="__item__"></span>`)
		bodyHTML = strings.ReplaceAll(bodyHTML, "{"+each.indexVar+"}", `<span class="__index__"></span>`)
		itemType := fieldTypes[each.listName]
		itemToJS := toJS(itemType, "item")
		hasElse := each.elseHTML != ""

		fmt.Fprintf(sb, "%s%s_anchor := document.Call(\"getElementById\", \"%s_anchor\")\n", indent, fullID, fullID)
		if hasElse {
			fmt.Fprintf(sb, "%s%s_else := document.Call(\"getElementById\", \"%s_else\")\n", indent, fullID, fullID)
		}
		fmt.Fprintf(sb, "%s%s_tmpl := %s\n", indent, fullID, escapeForGoString(bodyHTML))
		fmt.Fprintf(sb, "%s%s_create := func(item %s, index int) js.Value {\n", indent, fullID, itemType)
		fmt.Fprintf(sb, "%s\twrapper := document.Call(\"createElement\", \"span\")\n", indent)
		fmt.Fprintf(sb, "%s\twrapper.Set(\"id\", \"%s_\" + strconv.Itoa(index))\n", indent, fullID)
		fmt.Fprintf(sb, "%s\twrapper.Set(\"innerHTML\", %s_tmpl)\n", indent, fullID)
		fmt.Fprintf(sb, "%s\tif itemEl := wrapper.Call(\"querySelector\", \".__item__\"); !itemEl.IsNull() {\n", indent)
		fmt.Fprintf(sb, "%s\t\titemEl.Set(\"textContent\", %s)\n%s\t\titemEl.Get(\"classList\").Call(\"remove\", \"__item__\")\n%s\t}\n", indent, itemToJS, indent, indent)
		fmt.Fprintf(sb, "%s\tif idxEl := wrapper.Call(\"querySelector\", \".__index__\"); !idxEl.IsNull() {\n", indent)
		fmt.Fprintf(sb, "%s\t\tidxEl.Set(\"textContent\", strconv.Itoa(index))\n%s\t\tidxEl.Get(\"classList\").Call(\"remove\", \"__index__\")\n%s\t}\n", indent, indent, indent)
		fmt.Fprintf(sb, "%s\treturn wrapper\n%s}\n", indent, indent)

		fmt.Fprintf(sb, "%s%s.%s.OnEdit(func(edit preveltekit.Edit[%s]) {\n%s\tswitch edit.Op {\n", indent, varName, each.listName, itemType, indent)
		fmt.Fprintf(sb, "%s\tcase preveltekit.EditInsert:\n", indent)
		fmt.Fprintf(sb, "%s\t\titems := %s.%s.Get()\n", indent, varName, each.listName)
		if hasElse {
			fmt.Fprintf(sb, "%s\t\tif len(items) == 1 { %s_else.Get(\"style\").Set(\"display\", \"none\") }\n", indent, fullID)
		}
		fmt.Fprintf(sb, "%s\t\tfor i := len(items) - 1; i > edit.Index; i-- {\n", indent)
		fmt.Fprintf(sb, "%s\t\t\tel := document.Call(\"getElementById\", \"%s_\" + strconv.Itoa(i-1))\n", indent, fullID)
		fmt.Fprintf(sb, "%s\t\t\tif !el.IsNull() { el.Set(\"id\", \"%s_\" + strconv.Itoa(i)) }\n%s\t\t}\n", indent, fullID, indent)
		fmt.Fprintf(sb, "%s\t\tel := %s_create(edit.Value, edit.Index)\n", indent, fullID)
		fmt.Fprintf(sb, "%s\t\tif edit.Index == 0 {\n", indent)
		fmt.Fprintf(sb, "%s\t\t\tfirst := document.Call(\"getElementById\", \"%s_1\")\n", indent, fullID)
		fmt.Fprintf(sb, "%s\t\t\tif !first.IsNull() { %s_anchor.Get(\"parentNode\").Call(\"insertBefore\", el, first) }\n", indent, fullID)
		fmt.Fprintf(sb, "%s\t\t\telse { %s_anchor.Get(\"parentNode\").Call(\"insertBefore\", el, %s_anchor) }\n", indent, fullID, fullID)
		fmt.Fprintf(sb, "%s\t\t} else {\n", indent)
		fmt.Fprintf(sb, "%s\t\t\tprev := document.Call(\"getElementById\", \"%s_\" + strconv.Itoa(edit.Index-1))\n", indent, fullID)
		fmt.Fprintf(sb, "%s\t\t\tif !prev.IsNull() { prev.Get(\"parentNode\").Call(\"insertBefore\", el, prev.Get(\"nextSibling\")) }\n", indent)
		fmt.Fprintf(sb, "%s\t\t\telse { %s_anchor.Get(\"parentNode\").Call(\"insertBefore\", el, %s_anchor) }\n%s\t\t}\n", indent, fullID, fullID, indent)
		fmt.Fprintf(sb, "%s\tcase preveltekit.EditRemove:\n", indent)
		fmt.Fprintf(sb, "%s\t\tel := document.Call(\"getElementById\", \"%s_\" + strconv.Itoa(edit.Index))\n", indent, fullID)
		fmt.Fprintf(sb, "%s\t\tif !el.IsNull() { el.Call(\"remove\") }\n", indent)
		fmt.Fprintf(sb, "%s\t\tfor i := edit.Index; ; i++ {\n", indent)
		fmt.Fprintf(sb, "%s\t\t\tnextEl := document.Call(\"getElementById\", \"%s_\" + strconv.Itoa(i+1))\n", indent, fullID)
		fmt.Fprintf(sb, "%s\t\t\tif nextEl.IsNull() { break }\n", indent)
		fmt.Fprintf(sb, "%s\t\t\tnextEl.Set(\"id\", \"%s_\" + strconv.Itoa(i))\n%s\t\t}\n", indent, fullID, indent)
		if hasElse {
			fmt.Fprintf(sb, "%s\t\tif len(%s.%s.Get()) == 0 { %s_else.Get(\"style\").Set(\"display\", \"\") }\n", indent, varName, each.listName, fullID)
		}
		fmt.Fprintf(sb, "%s\t}\n%s})\n", indent, indent)

		fmt.Fprintf(sb, "%s%s.%s.OnRender(func(items []%s) {\n", indent, varName, each.listName, itemType)
		if hasElse {
			fmt.Fprintf(sb, "%s\tif len(items) == 0 { %s_else.Get(\"style\").Set(\"display\", \"\") } else { %s_else.Get(\"style\").Set(\"display\", \"none\") }\n", indent, fullID, fullID)
		}
		fmt.Fprintf(sb, "%s\tfor i, item := range items {\n", indent)
		fmt.Fprintf(sb, "%s\t\tel := %s_create(item, i)\n", indent, fullID)
		fmt.Fprintf(sb, "%s\t\t%s_anchor.Get(\"parentNode\").Call(\"insertBefore\", el, %s_anchor)\n%s\t}\n%s})\n", indent, fullID, fullID, indent, indent)
		fmt.Fprintf(sb, "%s%s.%s.Render()\n\n", indent, varName, each.listName)
	}
}

// generateIfBlocksWiring generates wiring for if blocks.
func generateIfBlocksWiring(sb *strings.Builder, ifBlocks []ifBinding, components []componentBinding, ctx *WiringContext) {
	indent := ctx.Indent
	fieldTypes := ctx.Scope.FieldTypes
	varName := ctx.Scope.VarName

	for _, ifb := range ifBlocks {
		fullID := ctx.prefixID(ifb.elementID)
		fmt.Fprintf(sb, "%s%s_anchor := document.Call(\"getElementById\", \"%s_anchor\")\n", indent, fullID, fullID)
		fmt.Fprintf(sb, "%s%s_current := js.Null()\n", indent, fullID)

		// Find components in this if-block
		allHTML := ifb.elseHTML
		for _, branch := range ifb.branches {
			allHTML += branch.html
		}
		compIDs := findCompPlaceholders(allHTML)
		ssrSeen := make(map[string]bool)
		for _, compID := range compIDs {
			if !ssrSeen[compID] {
				ssrSeen[compID] = true
				prefixedCompID := ctx.prefixID(compID)
				fmt.Fprintf(sb, "%sif el := document.Call(\"getElementById\", \"%s\"); !el.IsNull() { el.Call(\"remove\") }\n", indent, prefixedCompID)
			}
		}

		// Find component bindings for components in this if-block
		var compsInBlock []componentBinding
		for _, comp := range components {
			if ssrSeen[comp.elementID] {
				compsInBlock = append(compsInBlock, comp)
			}
		}

		// Create component instances BEFORE update function (state persistence)
		for _, compBinding := range compsInBlock {
			childDef := ctx.AllComponents[compBinding.name]
			if childDef == nil {
				continue
			}
			childTmpl := strings.ReplaceAll(childDef.template, "<slot/>", compBinding.children)
			_, childBindings := parseTemplate(childTmpl)
			needsVar := len(childDef.fields) > 0 || len(compBinding.props) > 0 ||
				len(childBindings.events) > 0 || len(childBindings.ifBlocks) > 0 || childDef.hasOnMount
			if needsVar {
				prefixedCompID := ctx.prefixID(compBinding.elementID)
				fmt.Fprintf(sb, "%s%s := &%s{\n", indent, prefixedCompID, compBinding.name)
				generateFieldInit(sb, childDef.fields, indent+"\t")
				fmt.Fprintf(sb, "%s}\n", indent)
			}
		}

		// Update function
		fmt.Fprintf(sb, "%supdate%s := func() {\n", indent, fullID)
		fmt.Fprintf(sb, "%s\tvar html string\n", indent)
		fmt.Fprintf(sb, "%s\tvar branchIdx int\n", indent)

		for i, branch := range ifb.branches {
			cond := transformCondition(branch.condition, fieldTypes, varName)
			branchHTML := prefixExprIDs(branch.html, ctx.ID)
			if i == 0 {
				fmt.Fprintf(sb, "%s\tif %s {\n%s\t\thtml = %s\n%s\t\tbranchIdx = %d\n", indent, cond, indent, escapeForGoString(branchHTML), indent, i)
			} else {
				fmt.Fprintf(sb, "%s\t} else if %s {\n%s\t\thtml = %s\n%s\t\tbranchIdx = %d\n", indent, cond, indent, escapeForGoString(branchHTML), indent, i)
			}
		}

		if ifb.elseHTML != "" {
			elseHTML := prefixExprIDs(ifb.elseHTML, ctx.ID)
			fmt.Fprintf(sb, "%s\t} else {\n%s\t\thtml = %s\n%s\t\tbranchIdx = -1\n%s\t}\n", indent, indent, escapeForGoString(elseHTML), indent, indent)
		} else {
			fmt.Fprintf(sb, "%s\t} else {\n%s\t\tbranchIdx = -1\n%s\t}\n", indent, indent, indent)
		}

		// Replace component placeholders
		seen := make(map[string]bool)
		for _, compID := range compIDs {
			if seen[compID] {
				continue
			}
			seen[compID] = true
			prefixedCompID := ctx.prefixID(compID)
			fmt.Fprintf(sb, "%s\thtml = strings.Replace(html, \"<!--%s-->\", %sHTML, 1)\n", indent, compID, prefixedCompID)
		}

		// Replace nested component placeholders
		for _, compBinding := range compsInBlock {
			childDef := ctx.AllComponents[compBinding.name]
			if childDef == nil {
				continue
			}
			childTmpl := strings.ReplaceAll(childDef.template, "<slot/>", compBinding.children)
			_, childBindings := parseTemplate(childTmpl)
			for _, nestedComp := range childBindings.components {
				nestedID := compBinding.elementID + "_" + nestedComp.elementID
				if seen[nestedID] {
					continue
				}
				seen[nestedID] = true
				prefixedNestedID := ctx.prefixID(nestedID)
				fmt.Fprintf(sb, "%s\thtml = strings.Replace(html, \"<!--%s-->\", %sHTML, 1)\n", indent, nestedID, prefixedNestedID)
			}
		}

		// Insert HTML into DOM
		fmt.Fprintf(sb, "%s\t%s_current = preveltekit.ReplaceContent(%s_anchor, %s_current, html)\n", indent, fullID, fullID, fullID)

		// Wire up child components
		for _, compBinding := range compsInBlock {
			childDef := ctx.AllComponents[compBinding.name]
			if childDef != nil {
				childCtx := ctx.ChildContext(compBinding, childDef)
				childCtx.SkipCreate = true
				childCtx.InsideIfBlock = true
				childCtx.Indent = indent + "\t"
				generateComponentWiring(sb, childCtx)
			}
		}

		// Wire up each blocks inside branches
		for i, branch := range ifb.branches {
			if len(branch.eachBlocks) > 0 {
				fmt.Fprintf(sb, "%s\tif branchIdx == %d {\n", indent, i)
				for _, each := range branch.eachBlocks {
					generateEachBlockInline(sb, each, ctx, indent+"\t\t")
				}
				fmt.Fprintf(sb, "%s\t}\n", indent)
			}
		}

		// Wire up class bindings inside branches
		for i, branch := range ifb.branches {
			if len(branch.classBindings) > 0 {
				fmt.Fprintf(sb, "%s\tif branchIdx == %d {\n", indent, i)
				for _, cb := range branch.classBindings {
					generateClassBindingInline(sb, cb, ctx, indent+"\t\t")
				}
				fmt.Fprintf(sb, "%s\t}\n", indent)
			}
		}

		// Wire up expressions inside branches
		for i, branch := range ifb.branches {
			if len(branch.expressions) > 0 {
				fmt.Fprintf(sb, "%s\tif branchIdx == %d {\n", indent, i)
				for _, expr := range branch.expressions {
					fullExprID := ctx.prefixID(expr.elementID)
					varRef := varName + "." + expr.fieldName
					fmt.Fprintf(sb, "%s\t\tpreveltekit.Bind(\"%s\", %s)\n", indent, fullExprID, varRef)
				}
				fmt.Fprintf(sb, "%s\t}\n", indent)
			}
		}

		// Wire up expressions in else branch
		if len(ifb.elseExpressions) > 0 {
			fmt.Fprintf(sb, "%s\tif branchIdx == -1 {\n", indent)
			for _, expr := range ifb.elseExpressions {
				fullExprID := ctx.prefixID(expr.elementID)
				varRef := varName + "." + expr.fieldName
				fmt.Fprintf(sb, "%s\t\tpreveltekit.Bind(\"%s\", %s)\n", indent, fullExprID, varRef)
			}
			fmt.Fprintf(sb, "%s\t}\n", indent)
		}
		fmt.Fprintf(sb, "%s\t_ = branchIdx\n", indent)
		fmt.Fprintf(sb, "%s}\n", indent)

		// Subscribe to dependencies
		for _, dep := range ifb.deps {
			fmt.Fprintf(sb, "%s%s.%s.OnChange(func(_ %s) { update%s() })\n", indent, varName, dep, fieldTypes[dep], fullID)
		}
		fmt.Fprintf(sb, "%supdate%s()\n\n", indent, fullID)
	}
}

// generateEachBlockInline generates an each block inside an if block.
func generateEachBlockInline(sb *strings.Builder, each eachBinding, ctx *WiringContext, indent string) {
	fieldTypes := ctx.Scope.FieldTypes
	varName := ctx.Scope.VarName

	bodyHTML := strings.ReplaceAll(each.bodyHTML, "{"+each.itemVar+"}", `<span class="__item__"></span>`)
	bodyHTML = strings.ReplaceAll(bodyHTML, "{"+each.indexVar+"}", `<span class="__index__"></span>`)
	itemType := fieldTypes[each.listName]
	itemToJS := toJS(itemType, "item")
	hasElse := each.elseHTML != ""

	// each.elementID is already prefixed
	fullID := each.elementID

	fmt.Fprintf(sb, "%s%s_anchor := document.Call(\"getElementById\", \"%s_anchor\")\n", indent, fullID, fullID)
	if hasElse {
		fmt.Fprintf(sb, "%s%s_else := document.Call(\"getElementById\", \"%s_else\")\n", indent, fullID, fullID)
	}
	fmt.Fprintf(sb, "%s%s_tmpl := %s\n", indent, fullID, escapeForGoString(bodyHTML))
	fmt.Fprintf(sb, "%s%s_create := func(item %s, index int) js.Value {\n", indent, fullID, itemType)
	fmt.Fprintf(sb, "%s\twrapper := document.Call(\"createElement\", \"span\")\n", indent)
	fmt.Fprintf(sb, "%s\twrapper.Set(\"id\", \"%s_\" + strconv.Itoa(index))\n", indent, fullID)
	fmt.Fprintf(sb, "%s\twrapper.Set(\"innerHTML\", %s_tmpl)\n", indent, fullID)
	fmt.Fprintf(sb, "%s\tif itemEl := wrapper.Call(\"querySelector\", \".__item__\"); !itemEl.IsNull() {\n", indent)
	fmt.Fprintf(sb, "%s\t\titemEl.Set(\"textContent\", %s)\n%s\t\titemEl.Get(\"classList\").Call(\"remove\", \"__item__\")\n%s\t}\n", indent, itemToJS, indent, indent)
	fmt.Fprintf(sb, "%s\tif idxEl := wrapper.Call(\"querySelector\", \".__index__\"); !idxEl.IsNull() {\n", indent)
	fmt.Fprintf(sb, "%s\t\tidxEl.Set(\"textContent\", strconv.Itoa(index))\n%s\t\tidxEl.Get(\"classList\").Call(\"remove\", \"__index__\")\n%s\t}\n", indent, indent, indent)
	fmt.Fprintf(sb, "%s\treturn wrapper\n%s}\n", indent, indent)

	// Clear callbacks and use OnRender
	fmt.Fprintf(sb, "%s%s.%s.ClearCallbacks()\n", indent, varName, each.listName)
	fmt.Fprintf(sb, "%s%s.%s.OnRender(func(items []%s) {\n", indent, varName, each.listName, itemType)
	if hasElse {
		fmt.Fprintf(sb, "%s\tif len(items) == 0 { %s_else.Get(\"style\").Set(\"display\", \"\") } else { %s_else.Get(\"style\").Set(\"display\", \"none\") }\n", indent, fullID, fullID)
	}
	fmt.Fprintf(sb, "%s\tfor i, item := range items {\n", indent)
	fmt.Fprintf(sb, "%s\t\tel := %s_create(item, i)\n", indent, fullID)
	fmt.Fprintf(sb, "%s\t\t%s_anchor.Get(\"parentNode\").Call(\"insertBefore\", el, %s_anchor)\n%s\t}\n%s})\n", indent, fullID, fullID, indent, indent)
	fmt.Fprintf(sb, "%s%s.%s.Render()\n", indent, varName, each.listName)
}

// generateClassBindingInline generates a class binding inside an if block.
func generateClassBindingInline(sb *strings.Builder, cb classBinding, ctx *WiringContext, indent string) {
	fieldTypes := ctx.Scope.FieldTypes
	varName := ctx.Scope.VarName

	// cb.elementID is already prefixed
	fullID := cb.elementID

	fmt.Fprintf(sb, "%s%s := document.Call(\"getElementById\", \"%s\")\n", indent, fullID, fullID)

	// Just set initial state since DOM is recreated on each update
	cond := transformCondition(cb.condition, fieldTypes, varName)
	fmt.Fprintf(sb, "%spreveltekit.ToggleClass(%s, \"%s\", %s)\n",
		indent, fullID, cb.className, cond)
}

// generateComponentWiring generates wiring code for a single component.
// This is the unified function that handles all component wiring at any nesting depth.
func generateComponentWiring(sb *strings.Builder, ctx *WiringContext) {
	indent := ctx.Indent
	compID := ctx.ID
	compDef := ctx.Definition
	compBinding := ctx.Binding
	childFieldTypes := buildFieldTypes(compDef)

	// Check if element exists
	fmt.Fprintf(sb, "%s%s_el := preveltekit.GetEl(\"%s\")\n", indent, compID, compID)
	fmt.Fprintf(sb, "%sif !%s_el.IsNull() && !%s_el.IsUndefined() {\n", indent, compID, compID)

	innerIndent := indent + "\t"

	// Create instance if needed
	if !ctx.SkipCreate {
		needsVar := len(compDef.fields) > 0 || len(compBinding.props) > 0 || compDef.hasOnMount
		if needsVar {
			fmt.Fprintf(sb, "%s%s := &%s{\n", innerIndent, compID, compBinding.name)
			generateFieldInit(sb, compDef.fields, innerIndent+"\t")
			fmt.Fprintf(sb, "%s}\n", innerIndent)
		}
	}

	// Inject style
	if compDef.style != "" {
		fmt.Fprintf(sb, "%spreveltekit.InjectStyle(\"%s\", %sCSS)\n", innerIndent, compBinding.name, strings.ToLower(compBinding.name))
	}

	// Set props
	for propName, propValue := range compBinding.props {
		childField := strings.Title(propName)
		if strings.HasPrefix(propValue, "{") && strings.HasSuffix(propValue, "}") {
			// Dynamic prop
			parentField := propValue[1 : len(propValue)-1]
			parentRef := ctx.ParentStoreRef(parentField)
			// Get parent field type from parent scope
			var parentFieldType string
			if ctx.Parent != nil {
				parentFieldType = ctx.Parent.Scope.FieldTypes[parentField]
			} else {
				parentFieldType = "string"
			}
			if parentFieldType == "" {
				parentFieldType = "string"
			}
			fmt.Fprintf(sb, "%s%s.%s.Set(%s.Get())\n", innerIndent, compID, childField, parentRef)
			fmt.Fprintf(sb, "%s%s.OnChange(func(v %s) { %s.%s.Set(v) })\n", innerIndent, parentRef, parentFieldType, compID, childField)
		} else {
			// Static prop
			childFieldType := childFieldTypes[childField]
			switch childFieldType {
			case "string":
				fmt.Fprintf(sb, "%s%s.%s.Set(%q)\n", innerIndent, compID, childField, propValue)
			case "int", "bool":
				fmt.Fprintf(sb, "%s%s.%s.Set(%s)\n", innerIndent, compID, childField, propValue)
			default:
				fmt.Fprintf(sb, "%s%s.%s.Set(%q)\n", innerIndent, compID, childField, propValue)
			}
		}
	}

	// Parse child template to get bindings
	childTmpl := strings.ReplaceAll(compDef.template, "<slot/>", compBinding.children)
	_, childBindings := parseTemplate(childTmpl)

	// Use the context's scope directly - it was already set up correctly by ChildContext
	// Don't create a redundant scope layer
	childScope := ctx.Scope

	// Expression bindings
	for _, expr := range childBindings.expressions {
		fullID := compID + "_" + expr.elementID
		// Resolve field in scope chain
		varRef, _, found := childScope.Resolve(expr.fieldName)
		if !found {
			varRef = compID + "." + expr.fieldName
		}
		fmt.Fprintf(sb, "%spreveltekit.Bind(\"%s\", %s)\n", innerIndent, fullID, varRef)
	}

	// Input bindings
	for _, bind := range childBindings.bindings {
		fullID := compID + "_" + bind.elementID
		valueType := childFieldTypes[bind.fieldName]
		if valueType == "" {
			valueType = "string"
		}
		varRef := compID + "." + bind.fieldName

		if bind.bindType == "checked" {
			fmt.Fprintf(sb, "%spreveltekit.BindCheckbox(\"%s\", %s)\n", innerIndent, fullID, varRef)
		} else {
			switch valueType {
			case "int":
				fmt.Fprintf(sb, "%spreveltekit.BindInputInt(\"%s\", %s)\n", innerIndent, fullID, varRef)
			default:
				fmt.Fprintf(sb, "%spreveltekit.BindInput(\"%s\", %s)\n", innerIndent, fullID, varRef)
			}
		}
	}

	// Class bindings - prefix IDs and reuse shared function
	for i := range childBindings.classBindings {
		childBindings.classBindings[i].elementID = compID + "_" + childBindings.classBindings[i].elementID
	}
	innerCtx := &WiringContext{
		Scope:  ctx.Scope,
		Indent: innerIndent,
	}
	generateClassBindingsWiring(sb, childBindings.classBindings, innerCtx)

	// Attribute bindings
	for _, ab := range childBindings.attrBindings {
		fullID := compID + "_" + ab.elementID
		if len(ab.fields) == 1 {
			field := ab.fields[0]
			fieldType := childFieldTypes[field]
			if fieldType == "" {
				fieldType = "string"
			}
			fmt.Fprintf(sb, "%sattr%s := document.Call(\"querySelector\", \"[data-attrbind=\\\"%s\\\"]\")\n", innerIndent, fullID, fullID)
			fmt.Fprintf(sb, "%supdateAttr%s := func() {\n", innerIndent, fullID)
			fmt.Fprintf(sb, "%s\tval := %s\n", innerIndent, escapeForGoString(ab.template))
			fmt.Fprintf(sb, "%s\tval = strings.ReplaceAll(val, \"{%s}\", %s)\n", innerIndent, field, toJS(fieldType, compID+"."+field+".Get()"))
			fmt.Fprintf(sb, "%s\tattr%s.Call(\"setAttribute\", \"%s\", val)\n", innerIndent, fullID, ab.attrName)
			fmt.Fprintf(sb, "%s}\n", innerIndent)
			fmt.Fprintf(sb, "%s%s.%s.OnChange(func(_ %s) { updateAttr%s() })\n", innerIndent, compID, field, fieldType, fullID)
			fmt.Fprintf(sb, "%supdateAttr%s()\n", innerIndent, fullID)
		} else {
			fmt.Fprintf(sb, "%sattr%s := document.Call(\"querySelector\", \"[data-attrbind=\\\"%s\\\"]\")\n", innerIndent, fullID, fullID)
			fmt.Fprintf(sb, "%supdateAttr%s := func() {\n%s\tval := %s\n", innerIndent, fullID, innerIndent, escapeForGoString(ab.template))
			for _, field := range ab.fields {
				fieldType := childFieldTypes[field]
				if fieldType == "" {
					fieldType = "string"
				}
				fmt.Fprintf(sb, "%s\tval = strings.ReplaceAll(val, \"{%s}\", %s)\n",
					innerIndent, field, toJS(fieldType, compID+"."+field+".Get()"))
			}
			fmt.Fprintf(sb, "%s\tattr%s.Call(\"setAttribute\", \"%s\", val)\n%s}\n", innerIndent, fullID, ab.attrName, innerIndent)
			for _, field := range ab.fields {
				fieldType := childFieldTypes[field]
				if fieldType == "" {
					fieldType = "string"
				}
				fmt.Fprintf(sb, "%s%s.%s.OnChange(func(_ %s) { updateAttr%s() })\n", innerIndent, compID, field, fieldType, fullID)
			}
			fmt.Fprintf(sb, "%supdateAttr%s()\n", innerIndent, fullID)
		}
	}

	// Internal events
	for _, evt := range childBindings.events {
		fullID := compID + "_" + evt.elementID
		callArgs := transformEventArgs(evt.args, childScope)
		fmt.Fprintf(sb, "%spreveltekit.On(preveltekit.GetEl(\"%s\"), \"%s\", func() { %s.%s(%s) })\n",
			innerIndent, fullID, evt.event, compID, evt.methodName, callArgs)
	}

	// Parent events on component root - these call the PARENT's methods
	for eventName, evt := range compBinding.events {
		// Use parent's scope for event args and method calls
		parentScope := ctx.Parent.Scope
		callArgs := transformEventArgs(evt.args, parentScope)
		fmt.Fprintf(sb, "%spreveltekit.On(%s_el, \"%s\", func() { %s.%s(%s) })\n",
			innerIndent, compID, eventName, parentScope.VarName, evt.method, callArgs)
	}

	// If blocks
	generateChildIfBlocks(sb, childBindings.ifBlocks, childBindings.components, ctx)

	// Each blocks
	for _, each := range childBindings.eachBlocks {
		each.elementID = compID + "_" + each.elementID
		generateEachBlockInline(sb, each, ctx, innerIndent)
	}

	// Nested components (recursively)
	for _, nestedBinding := range childBindings.components {
		nestedDef := ctx.AllComponents[nestedBinding.name]
		if nestedDef != nil {
			nestedCtx := ctx.ChildContext(nestedBinding, nestedDef)
			generateComponentWiring(sb, nestedCtx)
		}
	}

	// OnMount
	if compDef.hasOnMount {
		fmt.Fprintf(sb, "%s%s.OnMount()\n", innerIndent, compID)
	}

	fmt.Fprintf(sb, "%s}\n\n", indent)
}

// generateChildIfBlocks generates if blocks for a child component.
func generateChildIfBlocks(sb *strings.Builder, ifBlocks []ifBinding, components []componentBinding, ctx *WiringContext) {
	indent := ctx.Indent
	fieldTypes := ctx.Scope.FieldTypes
	compID := ctx.ID

	for _, ifb := range ifBlocks {
		fullID := compID + "_" + ifb.elementID
		fmt.Fprintf(sb, "%s%s_anchor := document.Call(\"getElementById\", \"%s_anchor\")\n", indent, fullID, fullID)
		fmt.Fprintf(sb, "%s%s_current := js.Null()\n", indent, fullID)

		// Find components in this if block
		allHTML := ifb.elseHTML
		for _, branch := range ifb.branches {
			allHTML += branch.html
		}

		// Find which components are in this if-block
		var compsInBlock []componentBinding
		for _, comp := range components {
			if strings.Contains(allHTML, "<!--"+comp.elementID+"-->") {
				compsInBlock = append(compsInBlock, comp)
			}
		}

		// Create component instances before update function
		for _, compBinding := range compsInBlock {
			childDef := ctx.AllComponents[compBinding.name]
			if childDef == nil {
				continue
			}
			childTmpl := strings.ReplaceAll(childDef.template, "<slot/>", compBinding.children)
			_, childBindings := parseTemplate(childTmpl)
			needsVar := len(childDef.fields) > 0 || len(compBinding.props) > 0 ||
				len(childBindings.events) > 0 || len(childBindings.ifBlocks) > 0 || childDef.hasOnMount
			if needsVar {
				nestedID := compID + "_" + compBinding.elementID
				fmt.Fprintf(sb, "%s%s := &%s{\n", indent, nestedID, compBinding.name)
				generateFieldInit(sb, childDef.fields, indent+"\t")
				fmt.Fprintf(sb, "%s}\n", indent)
			}
		}

		// Update function
		fmt.Fprintf(sb, "%supdate%s := func() {\n", indent, fullID)
		fmt.Fprintf(sb, "%s\tvar html string\n", indent)
		fmt.Fprintf(sb, "%s\tvar branchIdx int\n", indent)

		for i, branch := range ifb.branches {
			cond := transformCondition(branch.condition, fieldTypes, compID)
			branchHTML := prefixExprIDs(branch.html, ctx.ID)
			if i == 0 {
				fmt.Fprintf(sb, "%s\tif %s {\n%s\t\thtml = %s\n%s\t\tbranchIdx = %d\n", indent, cond, indent, escapeForGoString(branchHTML), indent, i)
			} else {
				fmt.Fprintf(sb, "%s\t} else if %s {\n%s\t\thtml = %s\n%s\t\tbranchIdx = %d\n", indent, cond, indent, escapeForGoString(branchHTML), indent, i)
			}
		}

		if ifb.elseHTML != "" {
			elseHTML := prefixExprIDs(ifb.elseHTML, ctx.ID)
			fmt.Fprintf(sb, "%s\t} else {\n%s\t\thtml = %s\n%s\t\tbranchIdx = -1\n%s\t}\n", indent, indent, escapeForGoString(elseHTML), indent, indent)
		} else {
			fmt.Fprintf(sb, "%s\t} else {\n%s\t\tbranchIdx = -1\n%s\t}\n", indent, indent, indent)
		}

		// Replace nested component placeholders
		for _, comp := range compsInBlock {
			nestedID := compID + "_" + comp.elementID
			fmt.Fprintf(sb, "%s\thtml = strings.Replace(html, \"<!--%s-->\", %sHTML, 1)\n", indent, comp.elementID, nestedID)
		}

		// Insert HTML
		fmt.Fprintf(sb, "%s\t%s_current = preveltekit.ReplaceContent(%s_anchor, %s_current, html)\n", indent, fullID, fullID, fullID)

		// Wire up components inside
		for _, compBinding := range compsInBlock {
			childDef := ctx.AllComponents[compBinding.name]
			if childDef != nil {
				childCtx := ctx.ChildContext(compBinding, childDef)
				childCtx.SkipCreate = true
				childCtx.InsideIfBlock = true
				childCtx.Indent = indent + "\t"
				generateComponentWiring(sb, childCtx)
			}
		}

		// Wire up each blocks
		for i, branch := range ifb.branches {
			if len(branch.eachBlocks) > 0 {
				fmt.Fprintf(sb, "%s\tif branchIdx == %d {\n", indent, i)
				for _, each := range branch.eachBlocks {
					generateEachBlockInline(sb, each, ctx, indent+"\t\t")
				}
				fmt.Fprintf(sb, "%s\t}\n", indent)
			}
		}

		// Wire up expressions inside branches
		for i, branch := range ifb.branches {
			if len(branch.expressions) > 0 {
				fmt.Fprintf(sb, "%s\tif branchIdx == %d {\n", indent, i)
				for _, expr := range branch.expressions {
					fullExprID := ctx.prefixID(expr.elementID)
					varRef := compID + "." + expr.fieldName
					fmt.Fprintf(sb, "%s\t\tpreveltekit.Bind(\"%s\", %s)\n", indent, fullExprID, varRef)
				}
				fmt.Fprintf(sb, "%s\t}\n", indent)
			}
		}

		// Wire up expressions in else branch
		if len(ifb.elseExpressions) > 0 {
			fmt.Fprintf(sb, "%s\tif branchIdx == -1 {\n", indent)
			for _, expr := range ifb.elseExpressions {
				fullExprID := ctx.prefixID(expr.elementID)
				varRef := compID + "." + expr.fieldName
				fmt.Fprintf(sb, "%s\t\tpreveltekit.Bind(\"%s\", %s)\n", indent, fullExprID, varRef)
			}
			fmt.Fprintf(sb, "%s\t}\n", indent)
		}
		fmt.Fprintf(sb, "%s\t_ = branchIdx\n", indent)
		fmt.Fprintf(sb, "%s}\n", indent)

		// Subscribe
		for _, dep := range ifb.deps {
			fmt.Fprintf(sb, "%s%s.%s.OnChange(func(_ %s) { update%s() })\n", indent, compID, dep, fieldTypes[dep], fullID)
		}
		fmt.Fprintf(sb, "%supdate%s()\n", indent, fullID)
	}
}

// Helper functions

func prefixIDStr(prefix, id string) string {
	if prefix == "" {
		return id
	}
	return prefix + "_" + id
}

// prefixExprIDs prefixes expression IDs (expr_*, html_*) in HTML with the given prefix
func prefixExprIDs(html, prefix string) string {
	if prefix == "" {
		return html
	}
	// Prefix expr_ IDs
	result := regexp.MustCompile(`id="(expr_[^"]+)"`).ReplaceAllString(html, `id="`+prefix+`_$1"`)
	// Prefix html_ IDs
	result = regexp.MustCompile(`id="(html_[^"]+)"`).ReplaceAllString(result, `id="`+prefix+`_$1"`)
	return result
}

func transformEventArgs(args string, scope *Scope) string {
	if args == "" {
		return ""
	}
	result := args
	for scope := scope; scope != nil; scope = scope.Parent {
		for fieldName := range scope.FieldTypes {
			if strings.Contains(result, fieldName) {
				result = strings.ReplaceAll(result, fieldName, scope.VarName+"."+fieldName+".Get()")
			}
		}
	}
	return result
}

func findComponentsInIfBlocks(ifBlocks []ifBinding) map[string]bool {
	result := make(map[string]bool)
	for _, ifb := range ifBlocks {
		for _, branch := range ifb.branches {
			for _, compID := range findCompPlaceholders(branch.html) {
				result[compID] = true
			}
		}
		for _, compID := range findCompPlaceholders(ifb.elseHTML) {
			result[compID] = true
		}
	}
	return result
}

func needsStrconvImport(bindings templateBindings, childComponents map[string]*component, fieldTypes map[string]string) bool {
	if len(bindings.eachBlocks) > 0 {
		return true
	}
	for _, ifb := range bindings.ifBlocks {
		for _, branch := range ifb.branches {
			if len(branch.eachBlocks) > 0 {
				return true
			}
		}
	}
	for _, compBinding := range bindings.components {
		childComp := childComponents[compBinding.name]
		if childComp != nil {
			_, childBindings := parseTemplate(childComp.template)
			for _, ifb := range childBindings.ifBlocks {
				for _, branch := range ifb.branches {
					if len(branch.eachBlocks) > 0 {
						return true
					}
				}
			}
		}
	}
	for _, bind := range bindings.bindings {
		if needsStrconvForType(fieldTypes[bind.fieldName]) {
			return true
		}
	}
	return false
}

func needsStringsImport(bindings templateBindings, childComponents map[string]*component) bool {
	if len(bindings.attrBindings) > 0 {
		return true
	}
	for _, compBinding := range bindings.components {
		childComp := childComponents[compBinding.name]
		if childComp != nil {
			childFieldTypes := buildFieldTypes(childComp)
			_, childBindings := parseTemplate(childComp.template)
			for _, ab := range childBindings.attrBindings {
				if len(ab.fields) > 1 || (len(ab.fields) == 1 && childFieldTypes[ab.fields[0]] != "string") {
					return true
				}
			}
		}
	}
	for _, ifb := range bindings.ifBlocks {
		allHTML := ifb.elseHTML
		for _, branch := range ifb.branches {
			allHTML += branch.html
		}
		if hasCompPlaceholder(allHTML) {
			return true
		}
	}
	return false
}
