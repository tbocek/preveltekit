package main

import (
	"fmt"
	"strings"
)

func generateMain(comp *component, tmpl string, bindings templateBindings, childComponents map[string]*component) string {
	var sb strings.Builder
	fieldTypes := buildFieldTypes(comp)

	// Check what imports we need
	needsStrconv := len(bindings.eachBlocks) > 0
	// Also check for each blocks inside if blocks
	if !needsStrconv {
		for _, ifb := range bindings.ifBlocks {
			for _, branch := range ifb.branches {
				if len(branch.eachBlocks) > 0 {
					needsStrconv = true
					break
				}
			}
			if needsStrconv {
				break
			}
		}
	}
	// Also check child components for each blocks inside if blocks
	if !needsStrconv {
		for _, compBinding := range bindings.components {
			childComp := childComponents[compBinding.name]
			if childComp != nil {
				_, childBindings := parseTemplate(childComp.template)
				for _, ifb := range childBindings.ifBlocks {
					for _, branch := range ifb.branches {
						if len(branch.eachBlocks) > 0 {
							needsStrconv = true
							break
						}
					}
					if needsStrconv {
						break
					}
				}
			}
			if needsStrconv {
				break
			}
		}
	}
	for _, bind := range bindings.bindings {
		if needsStrconvForType(fieldTypes[bind.fieldName]) {
			needsStrconv = true
			break
		}
	}

	// strings is needed for parent-level attr bindings (child ones use reactive.BindAttr)
	needsStrings := len(bindings.attrBindings) > 0
	// Also check for multi-field child attr bindings which generate inline strings.ReplaceAll
	if !needsStrings {
		for _, compBinding := range bindings.components {
			childComp := childComponents[compBinding.name]
			if childComp != nil {
				childFieldTypes := buildFieldTypes(childComp)
				_, childBindings := parseTemplate(childComp.template)
				for _, ab := range childBindings.attrBindings {
					// Multi-field or non-string uses inline strings.ReplaceAll
					if len(ab.fields) > 1 || (len(ab.fields) == 1 && childFieldTypes[ab.fields[0]] != "string") {
						needsStrings = true
						break
					}
				}
			}
			if needsStrings {
				break
			}
		}
	}
	// Also needed if any components are inside if-blocks (for placeholder replacement)
	if !needsStrings {
		for _, ifb := range bindings.ifBlocks {
			allHTML := ifb.elseHTML
			for _, branch := range ifb.branches {
				allHTML += branch.html
			}
			if hasCompPlaceholder(allHTML) {
				needsStrings = true
				break
			}
		}
	}

	// Write imports - now imports reactive package
	sb.WriteString("//go:build js && wasm\n\npackage main\n\nimport (\n")
	sb.WriteString("\t\"reactive\"\n")
	sb.WriteString("\t\"syscall/js\"\n")
	if needsStrconv {
		sb.WriteString("\t\"strconv\"\n")
	}
	if needsStrings {
		sb.WriteString("\t\"strings\"\n")
	}
	sb.WriteString(")\n\n")

	// Aliases for cleaner generated code
	sb.WriteString("var document = reactive.Document\n\n")

	// Generate CSS constants and HTML constants for components
	cssGenerated := make(map[string]bool)
	for _, compBinding := range bindings.components {
		childComp := childComponents[compBinding.name]
		if childComp == nil {
			continue
		}
		if childComp.style != "" && !cssGenerated[compBinding.name] {
			cssGenerated[compBinding.name] = true
			fmt.Fprintf(&sb, "const %sCSS = `%s`\n", strings.ToLower(compBinding.name), childComp.style)
		}
	}

	// Build a set of component IDs that are inside if-blocks
	componentsInIfBlocks := make(map[string]bool)
	for _, ifb := range bindings.ifBlocks {
		for _, branch := range ifb.branches {
			for _, compID := range findCompPlaceholders(branch.html) {
				componentsInIfBlocks[compID] = true
			}
		}
		for _, compID := range findCompPlaceholders(ifb.elseHTML) {
			componentsInIfBlocks[compID] = true
		}
	}

	// Generate HTML constants for components inside if-blocks
	for _, compBinding := range bindings.components {
		if !componentsInIfBlocks[compBinding.elementID] {
			continue
		}
		childComp := childComponents[compBinding.name]
		if childComp == nil {
			continue
		}
		// Process child template with slot content
		childTmpl := strings.ReplaceAll(childComp.template, "<slot/>", compBinding.children)
		childTmplProcessed, childBindings := parseTemplate(childTmpl)

		// Prefix all IDs with component ID to avoid conflicts
		childTmplProcessed = prefixBindingIDs(compBinding.elementID, childTmplProcessed,
			childBindings.expressions, childBindings.events, childBindings.attrBindings, childBindings.ifBlocks)
		childTmplProcessed = prefixInputBindingIDs(compBinding.elementID, childTmplProcessed, childBindings.bindings)
		childTmplProcessed = prefixEachBindingIDs(compBinding.elementID, childTmplProcessed, childBindings.eachBlocks)
		childTmplProcessed = prefixClassBindingIDs(compBinding.elementID, childTmplProcessed, childBindings.classBindings)
		childTmplProcessed = prefixComponentPlaceholders(compBinding.elementID, childTmplProcessed, childBindings.components)

		childTmplProcessed = injectIDIntoFirstTag(childTmplProcessed, compBinding.elementID)
		fmt.Fprintf(&sb, "const %sHTML = %s\n", compBinding.elementID, escapeForGoString(childTmplProcessed))

		// Also generate HTML constants for nested components inside this child's if-blocks
		htmlGenerated := make(map[string]bool)
		generateNestedComponentConstants(&sb, childBindings.components, childBindings.ifBlocks, childComponents, cssGenerated, htmlGenerated, compBinding.elementID)
	}

	// Also generate constants for components nested in child component templates (not just main if-blocks)
	htmlGenerated := make(map[string]bool)
	for _, compBinding := range bindings.components {
		if componentsInIfBlocks[compBinding.elementID] {
			continue // Already processed above
		}
		childComp := childComponents[compBinding.name]
		if childComp == nil {
			continue
		}
		childTmpl := strings.ReplaceAll(childComp.template, "<slot/>", compBinding.children)
		_, childBindings := parseTemplate(childTmpl)
		generateNestedComponentConstants(&sb, childBindings.components, childBindings.ifBlocks, childComponents, cssGenerated, htmlGenerated, compBinding.elementID)
	}

	// Main function
	sb.WriteString("\nfunc main() {\n\tcomponent := &" + comp.name + "{\n")
	generateFieldInit(&sb, comp.fields, "\t\t")
	sb.WriteString("\t}\n\n")

	// Expression bindings
	generateExprBindings(&sb, bindings.expressions, fieldTypes, "component")

	// Event bindings
	generateEventBindings(&sb, bindings.events, fieldTypes)

	// If blocks
	generateIfBlocks(&sb, bindings.ifBlocks, fieldTypes, bindings.components, childComponents)

	// Input bindings
	generateInputBindings(&sb, bindings.bindings, fieldTypes)

	// Class bindings
	generateClassBindings(&sb, bindings.classBindings, fieldTypes)

	// Attribute bindings
	generateAttrBindings(&sb, bindings.attrBindings, fieldTypes)

	// Each blocks
	generateEachBlocks(&sb, bindings.eachBlocks, fieldTypes)

	// Child components (skip those inside if-blocks - they're wired inside the update function)
	for _, compBinding := range bindings.components {
		if !componentsInIfBlocks[compBinding.elementID] {
			generateChildComponent(&sb, compBinding, childComponents, fieldTypes)
		}
	}

	if comp.hasOnMount {
		sb.WriteString("\tcomponent.OnMount()\n\n")
	}

	sb.WriteString("\tselect {}\n}\n")
	return sb.String()
}

func generateExprBindings(sb *strings.Builder, expressions []exprBinding, fieldTypes map[string]string, prefix string) {
	for _, expr := range expressions {
		valueType := fieldTypes[expr.fieldName]
		if expr.isHTML {
			jsConv := toJS(valueType, "v")
			jsConvInit := toJS(valueType, prefix+"."+expr.fieldName+".Get()")
			fmt.Fprintf(sb, "\t%s := reactive.GetEl(\"%s\")\n", expr.elementID, expr.elementID)
			fmt.Fprintf(sb, "\t%s.%s.OnChange(func(v %s) { if !%s.IsUndefined() && !%s.IsNull() { %s.Set(\"innerHTML\", %s) } })\n",
				prefix, expr.fieldName, valueType, expr.elementID, expr.elementID, expr.elementID, jsConv)
			fmt.Fprintf(sb, "\tif !%s.IsUndefined() && !%s.IsNull() { %s.Set(\"innerHTML\", %s) }\n",
				expr.elementID, expr.elementID, expr.elementID, jsConvInit)
		} else {
			fmt.Fprintf(sb, "\treactive.Bind(\"%s\", %s.%s)\n", expr.elementID, prefix, expr.fieldName)
		}
	}
}

func generateEventBindings(sb *strings.Builder, events []eventBinding, fieldTypes map[string]string) {
	for _, evt := range events {
		callArgs := evt.args
		for fieldName := range fieldTypes {
			if strings.Contains(callArgs, fieldName) {
				callArgs = strings.ReplaceAll(callArgs, fieldName, "component."+fieldName+".Get()")
			}
		}

		var modifierCode string
		hasOnce := false
		for _, mod := range evt.modifiers {
			switch mod {
			case "preventDefault":
				modifierCode += "\t\t\targs[0].Call(\"preventDefault\")\n"
			case "stopPropagation":
				modifierCode += "\t\t\targs[0].Call(\"stopPropagation\")\n"
			case "once":
				hasOnce = true
			}
		}

		if hasOnce {
			fmt.Fprintf(sb, "\tdocument.Call(\"getElementById\", \"%s\").Call(\"addEventListener\", \"%s\",\n", evt.elementID, evt.event)
			fmt.Fprintf(sb, "\t\tjs.FuncOf(func(this js.Value, args []js.Value) any {\n%s\t\t\tcomponent.%s(%s)\n\t\t\treturn nil\n\t\t}), map[string]interface{}{\"once\": true})\n\n",
				modifierCode, evt.methodName, callArgs)
		} else {
			fmt.Fprintf(sb, "\tdocument.Call(\"getElementById\", \"%s\").Call(\"addEventListener\", \"%s\",\n", evt.elementID, evt.event)
			fmt.Fprintf(sb, "\t\tjs.FuncOf(func(this js.Value, args []js.Value) any {\n%s\t\t\tcomponent.%s(%s)\n\t\t\treturn nil\n\t\t}))\n\n",
				modifierCode, evt.methodName, callArgs)
		}
	}
}

func generateIfBlocks(sb *strings.Builder, ifBlocks []ifBinding, fieldTypes map[string]string, components []componentBinding, childComponents map[string]*component) {
	for _, ifb := range ifBlocks {
		fmt.Fprintf(sb, "\t%s_anchor := document.Call(\"getElementById\", \"%s_anchor\")\n", ifb.elementID, ifb.elementID)
		fmt.Fprintf(sb, "\t%s_current := js.Null()\n", ifb.elementID)

		// Find component IDs in this if-block to remove SSR content
		allHTML := ifb.elseHTML
		for _, branch := range ifb.branches {
			allHTML += branch.html
		}
		compIDs := findCompPlaceholders(allHTML)
		ssrSeen := make(map[string]bool)
		for _, compID := range compIDs {
			if !ssrSeen[compID] {
				ssrSeen[compID] = true
				// Remove SSR-rendered component element
				fmt.Fprintf(sb, "\tif el := document.Call(\"getElementById\", \"%s\"); !el.IsNull() { el.Call(\"remove\") }\n", compID)
			}
		}

		// Build list of components in this if-block
		var compsInBlock []componentBinding
		for _, comp := range components {
			if ssrSeen[comp.elementID] {
				compsInBlock = append(compsInBlock, comp)
			}
		}

		// Create component instances BEFORE the update function (so state persists)
		for _, compBinding := range compsInBlock {
			childComp := childComponents[compBinding.name]
			if childComp == nil {
				continue
			}
			_, childBindings := parseTemplate(strings.ReplaceAll(childComp.template, "<slot/>", compBinding.children))
			needsVar := len(childComp.fields) > 0 || len(compBinding.props) > 0 ||
				len(childBindings.events) > 0 || len(childBindings.ifBlocks) > 0 || childComp.hasOnMount
			if needsVar {
				fmt.Fprintf(sb, "\t%s := &%s{\n", compBinding.elementID, compBinding.name)
				generateFieldInit(sb, childComp.fields, "\t\t")
				sb.WriteString("\t}\n")
			}
		}

		fmt.Fprintf(sb, "\tupdate%s := func() {\n", ifb.elementID)
		sb.WriteString("\t\tvar html string\n")
		sb.WriteString("\t\tvar branchIdx int\n")

		for i, branch := range ifb.branches {
			cond := transformCondition(branch.condition, fieldTypes, "component")
			if i == 0 {
				fmt.Fprintf(sb, "\t\tif %s {\n\t\t\thtml = %s\n\t\t\tbranchIdx = %d\n", cond, escapeForGoString(branch.html), i)
			} else {
				fmt.Fprintf(sb, "\t\t} else if %s {\n\t\t\thtml = %s\n\t\t\tbranchIdx = %d\n", cond, escapeForGoString(branch.html), i)
			}
		}

		if ifb.elseHTML != "" {
			fmt.Fprintf(sb, "\t\t} else {\n\t\t\thtml = %s\n\t\t\tbranchIdx = -1\n\t\t}\n", escapeForGoString(ifb.elseHTML))
		} else {
			sb.WriteString("\t\t} else {\n\t\t\tbranchIdx = -1\n\t\t}\n")
		}

		// Replace component placeholders with actual HTML constants
		compIDList := findCompPlaceholders(allHTML)
		seen := make(map[string]bool)
		for _, compID := range compIDList {
			if seen[compID] {
				continue
			}
			seen[compID] = true
			fmt.Fprintf(sb, "\t\thtml = strings.Replace(html, \"<!--%s-->\", %sHTML, 1)\n", compID, compID)
		}

		// Also replace nested component placeholders (components inside the page components)
		for _, compBinding := range compsInBlock {
			childComp := childComponents[compBinding.name]
			if childComp == nil {
				continue
			}
			childTmpl := strings.ReplaceAll(childComp.template, "<slot/>", compBinding.children)
			_, childBindings := parseTemplate(childTmpl)
			// These nested components have prefixed IDs like comp1_comp0
			for _, nestedComp := range childBindings.components {
				nestedID := compBinding.elementID + "_" + nestedComp.elementID
				if seen[nestedID] {
					continue
				}
				seen[nestedID] = true
				fmt.Fprintf(sb, "\t\thtml = strings.Replace(html, \"<!--%s-->\", %sHTML, 1)\n", nestedID, nestedID)
			}
		}

		fmt.Fprintf(sb, "\t\tnewEl := document.Call(\"createElement\", \"span\")\n")
		fmt.Fprintf(sb, "\t\tnewEl.Set(\"innerHTML\", html)\n")
		fmt.Fprintf(sb, "\t\tif !%s_current.IsNull() { %s_current.Call(\"remove\") }\n", ifb.elementID, ifb.elementID)
		fmt.Fprintf(sb, "\t\tif !%s_anchor.IsNull() { %s_anchor.Get(\"parentNode\").Call(\"insertBefore\", newEl, %s_anchor) }\n", ifb.elementID, ifb.elementID, ifb.elementID)
		fmt.Fprintf(sb, "\t\t%s_current = newEl\n", ifb.elementID)

		// Wire up child components AFTER inserting HTML into DOM
		// skipInstanceCreation=true because instances were created before the update function
		for _, compBinding := range compsInBlock {
			generateComponentInline(sb, compBinding, childComponents, "\t\t", "", true)
		}

		// Wire up each blocks inside if branches AFTER inserting HTML into DOM
		for i, branch := range ifb.branches {
			if len(branch.eachBlocks) > 0 {
				fmt.Fprintf(sb, "\t\tif branchIdx == %d {\n", i)
				for _, each := range branch.eachBlocks {
					generateEachBlockInline(sb, each, fieldTypes, "\t\t\t")
				}
				sb.WriteString("\t\t}\n")
			}
		}

		// Wire up class bindings inside if branches AFTER inserting HTML into DOM
		for i, branch := range ifb.branches {
			if len(branch.classBindings) > 0 {
				fmt.Fprintf(sb, "\t\tif branchIdx == %d {\n", i)
				for _, cb := range branch.classBindings {
					generateClassBindingInline(sb, cb, fieldTypes, "component", "\t\t\t")
				}
				sb.WriteString("\t\t}\n")
			}
		}
		sb.WriteString("\t\t_ = branchIdx\n") // Avoid unused variable error

		sb.WriteString("\t}\n")

		for _, dep := range ifb.deps {
			fmt.Fprintf(sb, "\tcomponent.%s.OnChange(func(_ %s) { update%s() })\n", dep, fieldTypes[dep], ifb.elementID)
		}
		fmt.Fprintf(sb, "\tupdate%s()\n\n", ifb.elementID)
	}
}

func generateInputBindings(sb *strings.Builder, bindings []inputBinding, fieldTypes map[string]string) {
	for _, bind := range bindings {
		valueType := fieldTypes[bind.fieldName]
		if bind.bindType == "checked" {
			fmt.Fprintf(sb, "\t%s := document.Call(\"getElementById\", \"%s\")\n", bind.elementID, bind.elementID)
			fmt.Fprintf(sb, "\t%s.Call(\"addEventListener\", \"change\", js.FuncOf(func(this js.Value, args []js.Value) any {\n", bind.elementID)
			fmt.Fprintf(sb, "\t\tcomponent.%s.Set(this.Get(\"checked\").Bool())\n\t\treturn nil\n\t}))\n", bind.fieldName)
			fmt.Fprintf(sb, "\tcomponent.%s.OnChange(func(v bool) { %s.Set(\"checked\", v) })\n", bind.fieldName, bind.elementID)
			fmt.Fprintf(sb, "\t%s.Set(\"checked\", component.%s.Get())\n\n", bind.elementID, bind.fieldName)
		} else {
			fmt.Fprintf(sb, "\t%s := document.Call(\"getElementById\", \"%s\")\n", bind.elementID, bind.elementID)
			fmt.Fprintf(sb, "\t%s.Call(\"addEventListener\", \"input\", js.FuncOf(func(this js.Value, args []js.Value) any {\n", bind.elementID)
			fmt.Fprintf(sb, "\t\tval := this.Get(\"value\").String()\n")
			switch valueType {
			case "int":
				fmt.Fprintf(sb, "\t\tif v, err := strconv.Atoi(val); err == nil { component.%s.Set(v) }\n", bind.fieldName)
			case "float64":
				fmt.Fprintf(sb, "\t\tif v, err := strconv.ParseFloat(val, 64); err == nil { component.%s.Set(v) }\n", bind.fieldName)
			default:
				fmt.Fprintf(sb, "\t\tcomponent.%s.Set(val)\n", bind.fieldName)
			}
			fmt.Fprintf(sb, "\t\treturn nil\n\t}))\n")
			fmt.Fprintf(sb, "\tcomponent.%s.OnChange(func(v %s) { %s.Set(\"value\", %s) })\n\n",
				bind.fieldName, valueType, bind.elementID, toJS(valueType, "v"))
		}
	}
}

func generateClassBindings(sb *strings.Builder, classBindings []classBinding, fieldTypes map[string]string) {
	// Group class bindings by element ID
	byElement := make(map[string][]classBinding)
	elementOrder := []string{}
	for _, cb := range classBindings {
		if _, exists := byElement[cb.elementID]; !exists {
			elementOrder = append(elementOrder, cb.elementID)
		}
		byElement[cb.elementID] = append(byElement[cb.elementID], cb)
	}

	for _, elementID := range elementOrder {
		bindings := byElement[elementID]
		fmt.Fprintf(sb, "\t%s := document.Call(\"getElementById\", \"%s\")\n", elementID, elementID)

		// Collect all dependencies across all class bindings for this element
		allDeps := make(map[string]bool)
		for _, cb := range bindings {
			deps := extractPascalCaseWords(cb.condition)
			for _, dep := range deps {
				allDeps[dep] = true
			}
		}

		// Create a single update function for all class bindings on this element
		fmt.Fprintf(sb, "\tupdate%s := func() {\n", elementID)
		for _, cb := range bindings {
			cond := transformCondition(cb.condition, fieldTypes, "component")
			fmt.Fprintf(sb, "\t\tif %s { %s.Get(\"classList\").Call(\"add\", \"%s\") } else { %s.Get(\"classList\").Call(\"remove\", \"%s\") }\n",
				cond, elementID, cb.className, elementID, cb.className)
		}
		sb.WriteString("\t}\n")

		// Register OnChange for each unique dependency
		for dep := range allDeps {
			if fieldType, ok := fieldTypes[dep]; ok {
				fmt.Fprintf(sb, "\tcomponent.%s.OnChange(func(_ %s) { update%s() })\n", dep, fieldType, elementID)
			}
		}
		// Initial call
		fmt.Fprintf(sb, "\tupdate%s()\n\n", elementID)
	}
}

// generateClassBindingInline generates code for a class binding inside an if block
func generateClassBindingInline(sb *strings.Builder, cb classBinding, fieldTypes map[string]string, compPrefix string, indent string) {
	fmt.Fprintf(sb, "%s%s := document.Call(\"getElementById\", \"%s\")\n", indent, cb.elementID, cb.elementID)

	// Extract field dependencies from the condition
	deps := extractPascalCaseWords(cb.condition)

	// Check if it's a simple field reference or an expression
	if len(deps) == 1 && deps[0] == cb.condition && fieldTypes[cb.condition] == "bool" {
		// Simple boolean field: class:active={IsActive}
		// For inline (inside if-block), just set initial state since DOM is recreated on each update
		fmt.Fprintf(sb, "%sif %s.%s.Get() { %s.Get(\"classList\").Call(\"add\", \"%s\") }\n",
			indent, compPrefix, cb.condition, cb.elementID, cb.className)
	} else {
		// Expression: class:active={CurrentStep == 1}
		cond := transformCondition(cb.condition, fieldTypes, compPrefix)
		fmt.Fprintf(sb, "%sif %s { %s.Get(\"classList\").Call(\"add\", \"%s\") }\n",
			indent, cond, cb.elementID, cb.className)
	}
}

func generateAttrBindings(sb *strings.Builder, attrBindings []attrBinding, fieldTypes map[string]string) {
	for _, ab := range attrBindings {
		fmt.Fprintf(sb, "\tattr%s := document.Call(\"querySelector\", \"[data-attrbind=\\\"%s\\\"]\")\n", ab.elementID, ab.elementID)
		fmt.Fprintf(sb, "\tupdateAttr%s := func() {\n\t\tval := %s\n", ab.elementID, escapeForGoString(ab.template))
		for _, field := range ab.fields {
			fmt.Fprintf(sb, "\t\tval = strings.ReplaceAll(val, \"{%s}\", %s)\n", field, toJS(fieldTypes[field], "component."+field+".Get()"))
		}
		fmt.Fprintf(sb, "\t\tattr%s.Call(\"setAttribute\", \"%s\", val)\n\t}\n", ab.elementID, ab.attrName)
		for _, field := range ab.fields {
			fmt.Fprintf(sb, "\tcomponent.%s.OnChange(func(_ %s) { updateAttr%s() })\n", field, fieldTypes[field], ab.elementID)
		}
		fmt.Fprintf(sb, "\tupdateAttr%s()\n\n", ab.elementID)
	}
}

func generateEachBlocks(sb *strings.Builder, eachBlocks []eachBinding, fieldTypes map[string]string) {
	for _, each := range eachBlocks {
		bodyHTML := strings.ReplaceAll(each.bodyHTML, "{"+each.itemVar+"}", `<span class="__item__"></span>`)
		bodyHTML = strings.ReplaceAll(bodyHTML, "{"+each.indexVar+"}", `<span class="__index__"></span>`)
		itemType := fieldTypes[each.listName]
		itemToJS := toJS(itemType, "item")
		hasElse := each.elseHTML != ""

		fmt.Fprintf(sb, "\t%s_anchor := document.Call(\"getElementById\", \"%s_anchor\")\n", each.elementID, each.elementID)
		if hasElse {
			fmt.Fprintf(sb, "\t%s_else := document.Call(\"getElementById\", \"%s_else\")\n", each.elementID, each.elementID)
		}
		fmt.Fprintf(sb, "\t%s_tmpl := %s\n", each.elementID, escapeForGoString(bodyHTML))
		fmt.Fprintf(sb, "\t%s_create := func(item %s, index int) js.Value {\n", each.elementID, itemType)
		fmt.Fprintf(sb, "\t\twrapper := document.Call(\"createElement\", \"span\")\n")
		fmt.Fprintf(sb, "\t\twrapper.Set(\"id\", \"%s_\" + strconv.Itoa(index))\n", each.elementID)
		fmt.Fprintf(sb, "\t\twrapper.Set(\"innerHTML\", %s_tmpl)\n", each.elementID)
		fmt.Fprintf(sb, "\t\tif itemEl := wrapper.Call(\"querySelector\", \".__item__\"); !itemEl.IsNull() {\n")
		fmt.Fprintf(sb, "\t\t\titemEl.Set(\"textContent\", %s)\n\t\t\titemEl.Get(\"classList\").Call(\"remove\", \"__item__\")\n\t\t}\n", itemToJS)
		fmt.Fprintf(sb, "\t\tif idxEl := wrapper.Call(\"querySelector\", \".__index__\"); !idxEl.IsNull() {\n")
		fmt.Fprintf(sb, "\t\t\tidxEl.Set(\"textContent\", strconv.Itoa(index))\n\t\t\tidxEl.Get(\"classList\").Call(\"remove\", \"__index__\")\n\t\t}\n")
		fmt.Fprintf(sb, "\t\treturn wrapper\n\t}\n")

		fmt.Fprintf(sb, "\tcomponent.%s.OnEdit(func(edit reactive.Edit[%s]) {\n\t\tswitch edit.Op {\n", each.listName, itemType)
		fmt.Fprintf(sb, "\t\tcase reactive.EditInsert:\n")
		fmt.Fprintf(sb, "\t\t\titems := component.%s.Get()\n", each.listName)
		if hasElse {
			fmt.Fprintf(sb, "\t\t\tif len(items) == 1 { %s_else.Get(\"style\").Set(\"display\", \"none\") }\n", each.elementID)
		}
		fmt.Fprintf(sb, "\t\t\tfor i := len(items) - 1; i > edit.Index; i-- {\n")
		fmt.Fprintf(sb, "\t\t\t\tel := document.Call(\"getElementById\", \"%s_\" + strconv.Itoa(i-1))\n", each.elementID)
		fmt.Fprintf(sb, "\t\t\t\tif !el.IsNull() { el.Set(\"id\", \"%s_\" + strconv.Itoa(i)) }\n\t\t\t}\n", each.elementID)
		fmt.Fprintf(sb, "\t\t\tel := %s_create(edit.Value, edit.Index)\n", each.elementID)
		fmt.Fprintf(sb, "\t\t\tif edit.Index == 0 {\n")
		fmt.Fprintf(sb, "\t\t\t\tfirst := document.Call(\"getElementById\", \"%s_1\")\n", each.elementID)
		fmt.Fprintf(sb, "\t\t\t\tif !first.IsNull() { %s_anchor.Get(\"parentNode\").Call(\"insertBefore\", el, first) }\n", each.elementID)
		fmt.Fprintf(sb, "\t\t\t\telse { %s_anchor.Get(\"parentNode\").Call(\"insertBefore\", el, %s_anchor) }\n", each.elementID, each.elementID)
		fmt.Fprintf(sb, "\t\t\t} else {\n")
		fmt.Fprintf(sb, "\t\t\t\tprev := document.Call(\"getElementById\", \"%s_\" + strconv.Itoa(edit.Index-1))\n", each.elementID)
		fmt.Fprintf(sb, "\t\t\t\tif !prev.IsNull() { prev.Get(\"parentNode\").Call(\"insertBefore\", el, prev.Get(\"nextSibling\")) }\n")
		fmt.Fprintf(sb, "\t\t\t\telse { %s_anchor.Get(\"parentNode\").Call(\"insertBefore\", el, %s_anchor) }\n\t\t\t}\n", each.elementID, each.elementID)
		fmt.Fprintf(sb, "\t\tcase reactive.EditRemove:\n")
		fmt.Fprintf(sb, "\t\t\tel := document.Call(\"getElementById\", \"%s_\" + strconv.Itoa(edit.Index))\n", each.elementID)
		fmt.Fprintf(sb, "\t\t\tif !el.IsNull() { el.Call(\"remove\") }\n")
		fmt.Fprintf(sb, "\t\t\tfor i := edit.Index; ; i++ {\n")
		fmt.Fprintf(sb, "\t\t\t\tnextEl := document.Call(\"getElementById\", \"%s_\" + strconv.Itoa(i+1))\n", each.elementID)
		fmt.Fprintf(sb, "\t\t\t\tif nextEl.IsNull() { break }\n")
		fmt.Fprintf(sb, "\t\t\t\tnextEl.Set(\"id\", \"%s_\" + strconv.Itoa(i))\n\t\t\t}\n", each.elementID)
		if hasElse {
			fmt.Fprintf(sb, "\t\t\tif len(component.%s.Get()) == 0 { %s_else.Get(\"style\").Set(\"display\", \"\") }\n", each.listName, each.elementID)
		}
		fmt.Fprintf(sb, "\t\t}\n\t})\n")

		fmt.Fprintf(sb, "\tcomponent.%s.OnRender(func(items []%s) {\n", each.listName, itemType)
		if hasElse {
			fmt.Fprintf(sb, "\t\tif len(items) == 0 { %s_else.Get(\"style\").Set(\"display\", \"\") } else { %s_else.Get(\"style\").Set(\"display\", \"none\") }\n", each.elementID, each.elementID)
		}
		fmt.Fprintf(sb, "\t\tfor i, item := range items {\n")
		fmt.Fprintf(sb, "\t\t\tel := %s_create(item, i)\n", each.elementID)
		fmt.Fprintf(sb, "\t\t\t%s_anchor.Get(\"parentNode\").Call(\"insertBefore\", el, %s_anchor)\n\t\t}\n\t})\n", each.elementID, each.elementID)
		fmt.Fprintf(sb, "\tcomponent.%s.Render()\n\n", each.listName)
	}
}

// generateEachBlockInline generates code for an each block inside an if block
func generateEachBlockInline(sb *strings.Builder, each eachBinding, fieldTypes map[string]string, indent string) {
	generateEachBlockInlineWithPrefix(sb, each, fieldTypes, "component", indent)
}

// generateEachBlockInlineForChild generates code for an each block inside a child component's if block
func generateEachBlockInlineForChild(sb *strings.Builder, each eachBinding, fieldTypes map[string]string, compID string, indent string) {
	generateEachBlockInlineWithPrefix(sb, each, fieldTypes, compID, indent)
}

// generateEachBlockInlineWithPrefix generates code for an each block with a configurable component prefix
func generateEachBlockInlineWithPrefix(sb *strings.Builder, each eachBinding, fieldTypes map[string]string, compPrefix string, indent string) {
	bodyHTML := strings.ReplaceAll(each.bodyHTML, "{"+each.itemVar+"}", `<span class="__item__"></span>`)
	bodyHTML = strings.ReplaceAll(bodyHTML, "{"+each.indexVar+"}", `<span class="__index__"></span>`)
	itemType := fieldTypes[each.listName]
	itemToJS := toJS(itemType, "item")
	hasElse := each.elseHTML != ""

	fmt.Fprintf(sb, "%s%s_anchor := document.Call(\"getElementById\", \"%s_anchor\")\n", indent, each.elementID, each.elementID)
	if hasElse {
		fmt.Fprintf(sb, "%s%s_else := document.Call(\"getElementById\", \"%s_else\")\n", indent, each.elementID, each.elementID)
	}
	fmt.Fprintf(sb, "%s%s_tmpl := %s\n", indent, each.elementID, escapeForGoString(bodyHTML))
	fmt.Fprintf(sb, "%s%s_create := func(item %s, index int) js.Value {\n", indent, each.elementID, itemType)
	fmt.Fprintf(sb, "%s\twrapper := document.Call(\"createElement\", \"span\")\n", indent)
	fmt.Fprintf(sb, "%s\twrapper.Set(\"id\", \"%s_\" + strconv.Itoa(index))\n", indent, each.elementID)
	fmt.Fprintf(sb, "%s\twrapper.Set(\"innerHTML\", %s_tmpl)\n", indent, each.elementID)
	fmt.Fprintf(sb, "%s\tif itemEl := wrapper.Call(\"querySelector\", \".__item__\"); !itemEl.IsNull() {\n", indent)
	fmt.Fprintf(sb, "%s\t\titemEl.Set(\"textContent\", %s)\n%s\t\titemEl.Get(\"classList\").Call(\"remove\", \"__item__\")\n%s\t}\n", indent, itemToJS, indent, indent)
	fmt.Fprintf(sb, "%s\tif idxEl := wrapper.Call(\"querySelector\", \".__index__\"); !idxEl.IsNull() {\n", indent)
	fmt.Fprintf(sb, "%s\t\tidxEl.Set(\"textContent\", strconv.Itoa(index))\n%s\t\tidxEl.Get(\"classList\").Call(\"remove\", \"__index__\")\n%s\t}\n", indent, indent, indent)
	fmt.Fprintf(sb, "%s\treturn wrapper\n%s}\n", indent, indent)

	// For each blocks inside if blocks, we use OnRender to render all items
	// OnEdit is not used since the if block re-renders everything on change
	// Clear existing callbacks first to prevent accumulation on re-render
	fmt.Fprintf(sb, "%s%s.%s.ClearCallbacks()\n", indent, compPrefix, each.listName)
	fmt.Fprintf(sb, "%s%s.%s.OnRender(func(items []%s) {\n", indent, compPrefix, each.listName, itemType)
	if hasElse {
		fmt.Fprintf(sb, "%s\tif len(items) == 0 { %s_else.Get(\"style\").Set(\"display\", \"\") } else { %s_else.Get(\"style\").Set(\"display\", \"none\") }\n", indent, each.elementID, each.elementID)
	}
	fmt.Fprintf(sb, "%s\tfor i, item := range items {\n", indent)
	fmt.Fprintf(sb, "%s\t\tel := %s_create(item, i)\n", indent, each.elementID)
	fmt.Fprintf(sb, "%s\t\t%s_anchor.Get(\"parentNode\").Call(\"insertBefore\", el, %s_anchor)\n%s\t}\n%s})\n", indent, each.elementID, each.elementID, indent, indent)
	fmt.Fprintf(sb, "%s%s.%s.Render()\n", indent, compPrefix, each.listName)
}

func generateChildComponent(sb *strings.Builder, compBinding componentBinding, childComponents map[string]*component, fieldTypes map[string]string) {
	childComp := childComponents[compBinding.name]
	if childComp == nil {
		return
	}

	childFieldTypes := buildFieldTypes(childComp)

	// Parse slot content to identify which fields come from parent
	_, slotBindings := parseTemplate(compBinding.children)
	slotFields := make(map[string]bool)
	for _, expr := range slotBindings.expressions {
		slotFields[expr.fieldName] = true
	}

	// Process child template with slot content
	childTmpl := strings.ReplaceAll(childComp.template, "<slot/>", compBinding.children)
	childTmplProcessed, childBindings := parseTemplate(childTmpl)

	// Categorize expressions into parent (slot) vs child owned
	slotExprs, childOwnExprs := categorizeExpressions(childBindings.expressions, slotFields, fieldTypes, childFieldTypes)

	// Make IDs unique by prefixing with component ID
	childTmplProcessed = prefixBindingIDs(compBinding.elementID, childTmplProcessed,
		childOwnExprs, childBindings.events, childBindings.attrBindings, childBindings.ifBlocks)
	childTmplProcessed = prefixBindingIDs(compBinding.elementID, childTmplProcessed,
		slotExprs, nil, nil, nil)
	childTmplProcessed = prefixInputBindingIDs(compBinding.elementID, childTmplProcessed, childBindings.bindings)
	childTmplProcessed = prefixEachBindingIDs(compBinding.elementID, childTmplProcessed, childBindings.eachBlocks)
	childTmplProcessed = prefixClassBindingIDs(compBinding.elementID, childTmplProcessed, childBindings.classBindings)

	// Inject component ID into child's root element
	childTmplProcessed = injectIDIntoFirstTag(childTmplProcessed, compBinding.elementID)

	// Create component instance - only if it has fields, props, events, if-blocks, or OnMount
	needsVar := len(childComp.fields) > 0 || len(compBinding.props) > 0 ||
		len(childOwnExprs) > 0 || len(childBindings.events) > 0 ||
		len(childBindings.ifBlocks) > 0 || childComp.hasOnMount

	if needsVar {
		fmt.Fprintf(sb, "\t%s := &%s{\n", compBinding.elementID, compBinding.name)
		generateFieldInit(sb, childComp.fields, "\t\t")
		sb.WriteString("\t}\n\n")
	}

	// Get component root element by ID
	fmt.Fprintf(sb, "\t%s_el := reactive.GetEl(\"%s\")\n", compBinding.elementID, compBinding.elementID)
	fmt.Fprintf(sb, "\tif !%s_el.IsNull() && !%s_el.IsUndefined() {\n", compBinding.elementID, compBinding.elementID)

	// Style injection
	if childComp.style != "" {
		fmt.Fprintf(sb, "\treactive.InjectStyle(\"%s\", %sCSS)\n", compBinding.name, strings.ToLower(compBinding.name))
	}

	// Props
	for propName, propValue := range compBinding.props {
		childField := strings.Title(propName)
		if strings.HasPrefix(propValue, "{") && strings.HasSuffix(propValue, "}") {
			parentField := propValue[1 : len(propValue)-1]
			fmt.Fprintf(sb, "\t%s.%s.Set(component.%s.Get())\n", compBinding.elementID, childField, parentField)
			fmt.Fprintf(sb, "\tcomponent.%s.OnChange(func(v %s) { %s.%s.Set(v) })\n\n",
				parentField, fieldTypes[parentField], compBinding.elementID, childField)
		} else {
			childFieldType := childFieldTypes[childField]
			switch childFieldType {
			case "string":
				fmt.Fprintf(sb, "\t%s.%s.Set(%q)\n\n", compBinding.elementID, childField, propValue)
			case "int", "bool":
				fmt.Fprintf(sb, "\t%s.%s.Set(%s)\n\n", compBinding.elementID, childField, propValue)
			default:
				fmt.Fprintf(sb, "\t%s.%s.Set(%q)\n\n", compBinding.elementID, childField, propValue)
			}
		}
	}

	// Child's own expressions
	for _, expr := range childOwnExprs {
		fmt.Fprintf(sb, "\treactive.Bind(\"%s\", %s.%s)\n", expr.elementID, compBinding.elementID, expr.fieldName)
	}

	// Slot expressions (parent bindings)
	for _, expr := range slotExprs {
		fmt.Fprintf(sb, "\treactive.Bind(\"%s\", component.%s)\n", expr.elementID, expr.fieldName)
	}

	// Child attribute bindings
	for _, ab := range childBindings.attrBindings {
		if len(ab.fields) == 1 {
			field := ab.fields[0]
			fieldType := childFieldTypes[field]
			if fieldType == "" {
				fieldType = "string"
			}
			fmt.Fprintf(sb, "\t%s_el := document.Call(\"querySelector\", \"[data-attrbind=\\\"%s\\\"]\")\n", ab.elementID, ab.elementID)
			fmt.Fprintf(sb, "\t%s.%s.OnChange(func(v %s) { %s_el.Call(\"setAttribute\", \"%s\", strings.ReplaceAll(`%s`, \"{%s}\", %s)) })\n",
				compBinding.elementID, field, fieldType, ab.elementID, ab.attrName, ab.template, field, toJS(fieldType, "v"))
			fmt.Fprintf(sb, "\t%s_el.Call(\"setAttribute\", \"%s\", strings.ReplaceAll(`%s`, \"{%s}\", %s))\n",
				ab.elementID, ab.attrName, ab.template, field, toJS(fieldType, compBinding.elementID+"."+field+".Get()"))
		} else {
			fmt.Fprintf(sb, "\t%s_el := document.Call(\"querySelector\", \"[data-attrbind=\\\"%s\\\"]\")\n", ab.elementID, ab.elementID)
			fmt.Fprintf(sb, "\tupdate%s := func() {\n\t\tval := %s\n", ab.elementID, escapeForGoString(ab.template))
			for _, field := range ab.fields {
				fieldType := childFieldTypes[field]
				if fieldType == "" {
					fieldType = "string"
				}
				fmt.Fprintf(sb, "\t\tval = strings.ReplaceAll(val, \"{%s}\", %s)\n",
					field, toJS(fieldType, compBinding.elementID+"."+field+".Get()"))
			}
			fmt.Fprintf(sb, "\t\t%s_el.Call(\"setAttribute\", \"%s\", val)\n\t}\n", ab.elementID, ab.attrName)
			for _, field := range ab.fields {
				valueType := childFieldTypes[field]
				if valueType == "" {
					valueType = "any"
				}
				fmt.Fprintf(sb, "\t%s.%s.OnChange(func(_ %s) { update%s() })\n",
					compBinding.elementID, field, valueType, ab.elementID)
			}
			fmt.Fprintf(sb, "\tupdate%s()\n", ab.elementID)
		}
	}

	// Child internal events
	for _, evt := range childBindings.events {
		callArgs := evt.args
		for fieldName := range childFieldTypes {
			if strings.Contains(callArgs, fieldName) {
				callArgs = strings.ReplaceAll(callArgs, fieldName, compBinding.elementID+"."+fieldName+".Get()")
			}
		}
		fmt.Fprintf(sb, "\treactive.On(reactive.GetEl(\"%s\"), \"%s\", func() { %s.%s(%s) })\n",
			evt.elementID, evt.event, compBinding.elementID, evt.methodName, callArgs)
	}

	// Parent events on child (attach to component's root element)
	for eventName, evt := range compBinding.events {
		callArgs := evt.args
		for fieldName := range fieldTypes {
			if strings.Contains(callArgs, fieldName) {
				callArgs = strings.ReplaceAll(callArgs, fieldName, "component."+fieldName+".Get()")
			}
		}
		fmt.Fprintf(sb, "\treactive.On(%s_el, \"%s\", func() { component.%s(%s) })\n",
			compBinding.elementID, eventName, evt.method, callArgs)
	}

	// Child class bindings - group by element ID to handle multiple class bindings per element
	childClassByElement := make(map[string][]classBinding)
	childClassOrder := []string{}
	for _, cb := range childBindings.classBindings {
		if _, exists := childClassByElement[cb.elementID]; !exists {
			childClassOrder = append(childClassOrder, cb.elementID)
		}
		childClassByElement[cb.elementID] = append(childClassByElement[cb.elementID], cb)
	}

	for _, elementID := range childClassOrder {
		bindings := childClassByElement[elementID]
		fmt.Fprintf(sb, "\t%s := document.Call(\"getElementById\", \"%s\")\n", elementID, elementID)

		// Collect all dependencies
		allDeps := make(map[string]bool)
		for _, cb := range bindings {
			deps := extractPascalCaseWords(cb.condition)
			for _, dep := range deps {
				allDeps[dep] = true
			}
		}

		// Create a single update function for all class bindings on this element
		fmt.Fprintf(sb, "\tupdate%s := func() {\n", elementID)
		for _, cb := range bindings {
			cond := transformCondition(cb.condition, childFieldTypes, compBinding.elementID)
			fmt.Fprintf(sb, "\t\tif %s { %s.Get(\"classList\").Call(\"add\", \"%s\") } else { %s.Get(\"classList\").Call(\"remove\", \"%s\") }\n",
				cond, elementID, cb.className, elementID, cb.className)
		}
		sb.WriteString("\t}\n")

		// Register OnChange for each unique dependency
		for dep := range allDeps {
			if fieldType, ok := childFieldTypes[dep]; ok {
				fmt.Fprintf(sb, "\t%s.%s.OnChange(func(_ %s) { update%s() })\n", compBinding.elementID, dep, fieldType, elementID)
			}
		}
		// Initial call
		fmt.Fprintf(sb, "\tupdate%s()\n", elementID)
	}

	// Child internal if-blocks
	generateChildIfBlocks(sb, childBindings.ifBlocks, childFieldTypes, compBinding.elementID)

	if childComp.hasOnMount {
		fmt.Fprintf(sb, "\t%s.OnMount()\n", compBinding.elementID)
	}

	sb.WriteString("\t}\n\n")
}

// generateComponentInline generates inline wiring code for a child component.
// compID overrides compBinding.elementID when provided (used for nested components with prefixed IDs).
// skipInstanceCreation: if true, assumes component instance already exists (for if-block components)
func generateComponentInline(sb *strings.Builder, compBinding componentBinding, childComponents map[string]*component, indent string, compID string, skipInstanceCreation bool) {
	childComp := childComponents[compBinding.name]
	if childComp == nil {
		return
	}

	// Use provided compID or fall back to compBinding.elementID
	if compID == "" {
		compID = compBinding.elementID
	}

	childFieldTypes := buildFieldTypes(childComp)

	// Parse child template with slot content
	childTmpl := strings.ReplaceAll(childComp.template, "<slot/>", compBinding.children)
	_, childBindings := parseTemplate(childTmpl)

	// Check if element exists (it was just inserted)
	fmt.Fprintf(sb, "%s%s_el := reactive.GetEl(\"%s\")\n", indent, compID, compID)
	fmt.Fprintf(sb, "%sif !%s_el.IsNull() && !%s_el.IsUndefined() {\n", indent, compID, compID)

	// Create component instance if needed (skip if already created outside)
	if !skipInstanceCreation {
		needsVar := len(childComp.fields) > 0 || len(compBinding.props) > 0 ||
			len(childBindings.events) > 0 || len(childBindings.ifBlocks) > 0 || childComp.hasOnMount

		if needsVar {
			fmt.Fprintf(sb, "%s\t%s := &%s{\n", indent, compID, compBinding.name)
			generateFieldInit(sb, childComp.fields, indent+"\t\t")
			fmt.Fprintf(sb, "%s\t}\n", indent)
		}
	}

	// Style injection
	if childComp.style != "" {
		fmt.Fprintf(sb, "%s\treactive.InjectStyle(\"%s\", %sCSS)\n", indent, compBinding.name, strings.ToLower(compBinding.name))
	}

	// Child expression bindings (with prefixed IDs)
	for _, expr := range childBindings.expressions {
		prefixedID := compID + "_" + expr.elementID
		fmt.Fprintf(sb, "%s\treactive.Bind(\"%s\", %s.%s)\n", indent, prefixedID, compID, expr.fieldName)
	}

	// Child input bindings (two-way binding, with prefixed IDs)
	for _, bind := range childBindings.bindings {
		prefixedID := compID + "_" + bind.elementID
		valueType := childFieldTypes[bind.fieldName]
		if bind.bindType == "checked" {
			fmt.Fprintf(sb, "%s\t%s := document.Call(\"getElementById\", \"%s\")\n", indent, prefixedID, prefixedID)
			fmt.Fprintf(sb, "%s\t%s.Call(\"addEventListener\", \"change\", js.FuncOf(func(this js.Value, args []js.Value) any {\n", indent, prefixedID)
			fmt.Fprintf(sb, "%s\t\t%s.%s.Set(this.Get(\"checked\").Bool())\n%s\t\treturn nil\n%s\t}))\n", indent, compID, bind.fieldName, indent, indent)
			fmt.Fprintf(sb, "%s\t%s.%s.OnChange(func(v bool) { %s.Set(\"checked\", v) })\n", indent, compID, bind.fieldName, prefixedID)
			fmt.Fprintf(sb, "%s\t%s.Set(\"checked\", %s.%s.Get())\n", indent, prefixedID, compID, bind.fieldName)
		} else {
			fmt.Fprintf(sb, "%s\t%s := document.Call(\"getElementById\", \"%s\")\n", indent, prefixedID, prefixedID)
			fmt.Fprintf(sb, "%s\t%s.Call(\"addEventListener\", \"input\", js.FuncOf(func(this js.Value, args []js.Value) any {\n", indent, prefixedID)
			fmt.Fprintf(sb, "%s\t\tval := this.Get(\"value\").String()\n", indent)
			switch valueType {
			case "int":
				fmt.Fprintf(sb, "%s\t\tif v, err := strconv.Atoi(val); err == nil { %s.%s.Set(v) }\n", indent, compID, bind.fieldName)
			case "float64":
				fmt.Fprintf(sb, "%s\t\tif v, err := strconv.ParseFloat(val, 64); err == nil { %s.%s.Set(v) }\n", indent, compID, bind.fieldName)
			default:
				fmt.Fprintf(sb, "%s\t\t%s.%s.Set(val)\n", indent, compID, bind.fieldName)
			}
			fmt.Fprintf(sb, "%s\t\treturn nil\n%s\t}))\n", indent, indent)
			fmt.Fprintf(sb, "%s\t%s.%s.OnChange(func(v %s) { %s.Set(\"value\", %s) })\n",
				indent, compID, bind.fieldName, valueType, prefixedID, toJS(valueType, "v"))
		}
	}

	// Child class bindings (with prefixed IDs) - group by element ID
	inlineClassByElement := make(map[string][]classBinding)
	inlineClassOrder := []string{}
	for _, cb := range childBindings.classBindings {
		if _, exists := inlineClassByElement[cb.elementID]; !exists {
			inlineClassOrder = append(inlineClassOrder, cb.elementID)
		}
		inlineClassByElement[cb.elementID] = append(inlineClassByElement[cb.elementID], cb)
	}

	for _, elementID := range inlineClassOrder {
		bindings := inlineClassByElement[elementID]
		prefixedID := compID + "_" + elementID
		fmt.Fprintf(sb, "%s\t%s := document.Call(\"getElementById\", \"%s\")\n", indent, prefixedID, prefixedID)

		// Collect all dependencies
		allDeps := make(map[string]bool)
		for _, cb := range bindings {
			deps := extractPascalCaseWords(cb.condition)
			for _, dep := range deps {
				allDeps[dep] = true
			}
		}

		// Create a single update function for all class bindings on this element
		fmt.Fprintf(sb, "%s\tupdate%s := func() {\n", indent, prefixedID)
		for _, cb := range bindings {
			cond := transformCondition(cb.condition, childFieldTypes, compID)
			fmt.Fprintf(sb, "%s\t\tif %s { %s.Get(\"classList\").Call(\"add\", \"%s\") } else { %s.Get(\"classList\").Call(\"remove\", \"%s\") }\n",
				indent, cond, prefixedID, cb.className, prefixedID, cb.className)
		}
		fmt.Fprintf(sb, "%s\t}\n", indent)

		// Register OnChange for each unique dependency
		for dep := range allDeps {
			if fieldType, ok := childFieldTypes[dep]; ok {
				fmt.Fprintf(sb, "%s\t%s.%s.OnChange(func(_ %s) { update%s() })\n", indent, compID, dep, fieldType, prefixedID)
			}
		}
		// Initial call
		fmt.Fprintf(sb, "%s\tupdate%s()\n", indent, prefixedID)
	}

	// Child internal events (with prefixed IDs)
	for _, evt := range childBindings.events {
		prefixedID := compID + "_" + evt.elementID
		callArgs := evt.args
		for fieldName := range childFieldTypes {
			if strings.Contains(callArgs, fieldName) {
				callArgs = strings.ReplaceAll(callArgs, fieldName, compID+"."+fieldName+".Get()")
			}
		}
		fmt.Fprintf(sb, "%s\treactive.On(reactive.GetEl(\"%s\"), \"%s\", func() { %s.%s(%s) })\n",
			indent, prefixedID, evt.event, compID, evt.methodName, callArgs)
	}

	// Child internal if-blocks
	for _, ifb := range childBindings.ifBlocks {
		prefixedID := compID + "_" + ifb.elementID
		generateChildIfBlockInline(sb, ifb, childFieldTypes, compID, prefixedID, indent+"\t", childBindings.components, childComponents, compID)
	}

	// Wire up nested components (e.g., Badge, Card inside Components page)
	for _, nestedBinding := range childBindings.components {
		nestedCompID := compID + "_" + nestedBinding.elementID
		generateNestedComponentInline(sb, nestedBinding, childComponents, indent+"\t", nestedCompID, compID)
	}

	if childComp.hasOnMount {
		fmt.Fprintf(sb, "%s\t%s.OnMount()\n", indent, compID)
	}

	fmt.Fprintf(sb, "%s}\n", indent)
}

// generateNestedComponentInline wires up a nested component (e.g., Badge inside Components page)
// nestedCompID is the prefixed ID (e.g., comp1_comp0)
// parentCompID is the parent component ID (e.g., comp1) - used for prop bindings from parent stores
func generateNestedComponentInline(sb *strings.Builder, nestedBinding componentBinding, childComponents map[string]*component, indent string, nestedCompID string, parentCompID string) {
	nestedComp := childComponents[nestedBinding.name]
	if nestedComp == nil {
		return
	}

	nestedFieldTypes := buildFieldTypes(nestedComp)

	// Parse nested component template
	nestedTmpl := strings.ReplaceAll(nestedComp.template, "<slot/>", nestedBinding.children)
	_, nestedBindings := parseTemplate(nestedTmpl)

	// Check if element exists
	fmt.Fprintf(sb, "%s%s_el := reactive.GetEl(\"%s\")\n", indent, nestedCompID, nestedCompID)
	fmt.Fprintf(sb, "%sif !%s_el.IsNull() && !%s_el.IsUndefined() {\n", indent, nestedCompID, nestedCompID)

	// Create component instance if it has fields
	needsVar := len(nestedComp.fields) > 0 || len(nestedBinding.props) > 0 || nestedComp.hasOnMount
	if needsVar {
		fmt.Fprintf(sb, "%s\t%s := &%s{\n", indent, nestedCompID, nestedBinding.name)
		generateFieldInit(sb, nestedComp.fields, indent+"\t\t")
		fmt.Fprintf(sb, "%s\t}\n", indent)
	}

	// Style injection
	if nestedComp.style != "" {
		fmt.Fprintf(sb, "%s\treactive.InjectStyle(\"%s\", %sCSS)\n", indent, nestedBinding.name, strings.ToLower(nestedBinding.name))
	}

	// Set props from parent - handle both static and dynamic props
	for propName, propValue := range nestedBinding.props {
		fieldName := strings.ToUpper(propName[:1]) + propName[1:] // Ensure proper casing (e.g., label -> Label)
		// Check if dynamic prop (wrapped in {})
		if strings.HasPrefix(propValue, "{") && strings.HasSuffix(propValue, "}") {
			// Dynamic prop bound to parent's store - set up OnChange
			storeField := propValue[1 : len(propValue)-1] // Remove { and }
			fmt.Fprintf(sb, "%s\t%s.%s.Set(%s.%s.Get())\n", indent, nestedCompID, fieldName, parentCompID, storeField)
			fmt.Fprintf(sb, "%s\t%s.%s.OnChange(func(v string) { %s.%s.Set(v) })\n", indent, parentCompID, storeField, nestedCompID, fieldName)
		} else {
			// Static prop - just set the value
			fmt.Fprintf(sb, "%s\t%s.%s.Set(%q)\n", indent, nestedCompID, fieldName, propValue)
		}
	}

	// Expression bindings (with prefixed IDs)
	for _, expr := range nestedBindings.expressions {
		prefixedID := nestedCompID + "_" + expr.elementID
		fmt.Fprintf(sb, "%s\treactive.Bind(\"%s\", %s.%s)\n", indent, prefixedID, nestedCompID, expr.fieldName)
	}

	// Events from parent (e.g., @click on the nested component)
	for evtName, evt := range nestedBinding.events {
		fmt.Fprintf(sb, "%s\treactive.On(%s_el, \"%s\", func() { %s.%s(%s) })\n",
			indent, nestedCompID, evtName, parentCompID, evt.method, evt.args)
	}

	// Internal events within the nested component
	for _, evt := range nestedBindings.events {
		prefixedID := nestedCompID + "_" + evt.elementID
		fmt.Fprintf(sb, "%s\treactive.On(reactive.GetEl(\"%s\"), \"%s\", func() { %s.%s(%s) })\n",
			indent, prefixedID, evt.event, nestedCompID, evt.methodName, evt.args)
	}

	// Input bindings within nested component
	for _, bind := range nestedBindings.bindings {
		prefixedID := nestedCompID + "_" + bind.elementID
		valueType := nestedFieldTypes[bind.fieldName]
		if bind.bindType == "checked" {
			fmt.Fprintf(sb, "%s\t%s := document.Call(\"getElementById\", \"%s\")\n", indent, prefixedID, prefixedID)
			fmt.Fprintf(sb, "%s\t%s.Call(\"addEventListener\", \"change\", js.FuncOf(func(this js.Value, args []js.Value) any {\n", indent, prefixedID)
			fmt.Fprintf(sb, "%s\t\t%s.%s.Set(this.Get(\"checked\").Bool())\n%s\t\treturn nil\n%s\t}))\n", indent, nestedCompID, bind.fieldName, indent, indent)
			fmt.Fprintf(sb, "%s\t%s.%s.OnChange(func(v bool) { %s.Set(\"checked\", v) })\n", indent, nestedCompID, bind.fieldName, prefixedID)
			fmt.Fprintf(sb, "%s\t%s.Set(\"checked\", %s.%s.Get())\n", indent, prefixedID, nestedCompID, bind.fieldName)
		} else {
			fmt.Fprintf(sb, "%s\t%s := document.Call(\"getElementById\", \"%s\")\n", indent, prefixedID, prefixedID)
			fmt.Fprintf(sb, "%s\t%s.Call(\"addEventListener\", \"input\", js.FuncOf(func(this js.Value, args []js.Value) any {\n", indent, prefixedID)
			fmt.Fprintf(sb, "%s\t\tval := this.Get(\"value\").String()\n", indent)
			fmt.Fprintf(sb, "%s\t\t%s.%s.Set(val)\n", indent, nestedCompID, bind.fieldName)
			fmt.Fprintf(sb, "%s\t\treturn nil\n%s\t}))\n", indent, indent)
			fmt.Fprintf(sb, "%s\t%s.%s.OnChange(func(v %s) { %s.Set(\"value\", %s) })\n",
				indent, nestedCompID, bind.fieldName, valueType, prefixedID, toJS(valueType, "v"))
		}
	}

	// Attribute bindings within nested component (e.g., data-type="{Type}")
	for _, ab := range nestedBindings.attrBindings {
		prefixedID := nestedCompID + "_" + ab.elementID
		fmt.Fprintf(sb, "%s\tattr%s := document.Call(\"querySelector\", \"[data-attrbind=\\\"%s\\\"]\")\n", indent, prefixedID, prefixedID)
		fmt.Fprintf(sb, "%s\tupdateAttr%s := func() {\n%s\t\tval := %s\n", indent, prefixedID, indent, escapeForGoString(ab.template))
		for _, field := range ab.fields {
			fmt.Fprintf(sb, "%s\t\tval = strings.ReplaceAll(val, \"{%s}\", %s)\n", indent, field, toJS(nestedFieldTypes[field], nestedCompID+"."+field+".Get()"))
		}
		fmt.Fprintf(sb, "%s\t\tattr%s.Call(\"setAttribute\", \"%s\", val)\n%s\t}\n", indent, prefixedID, ab.attrName, indent)
		for _, field := range ab.fields {
			fmt.Fprintf(sb, "%s\t%s.%s.OnChange(func(_ %s) { updateAttr%s() })\n", indent, nestedCompID, field, nestedFieldTypes[field], prefixedID)
		}
		fmt.Fprintf(sb, "%s\tupdateAttr%s()\n", indent, prefixedID)
	}

	if nestedComp.hasOnMount {
		fmt.Fprintf(sb, "%s\t%s.OnMount()\n", indent, nestedCompID)
	}

	fmt.Fprintf(sb, "%s}\n", indent)
}

// generateChildIfBlockInline generates an if-block for a child component inline
func generateChildIfBlockInline(sb *strings.Builder, ifb ifBinding, fieldTypes map[string]string, compID string, anchorID string, indent string, nestedComponents []componentBinding, childComponents map[string]*component, parentCompID string) {
	fmt.Fprintf(sb, "%s%s_anchor := document.Call(\"getElementById\", \"%s_anchor\")\n", indent, anchorID, anchorID)
	fmt.Fprintf(sb, "%s%s_current := js.Null()\n", indent, anchorID)

	// Collect all HTML from branches to find nested components
	allHTML := ifb.elseHTML
	for _, branch := range ifb.branches {
		allHTML += branch.html
	}

	// Find which nested components are in this if-block
	var compsInBlock []componentBinding
	for _, comp := range nestedComponents {
		if strings.Contains(allHTML, "<!--"+comp.elementID+"-->") {
			compsInBlock = append(compsInBlock, comp)
		}
	}

	// Parse branch HTML to extract and process expressions
	// Process each branch and the else block through the parser
	processedBranches := make([]struct {
		condition string
		html      string
		exprs     []exprBinding
	}, len(ifb.branches))

	exprCounter := 0
	var allExprs []exprBinding

	for i, branch := range ifb.branches {
		processedHTML, exprs := processIfBranchExpressions(branch.html, anchorID, &exprCounter)
		// Prefix each block anchor IDs in the branch HTML
		for j := range ifb.branches[i].eachBlocks {
			oldEachID := ifb.branches[i].eachBlocks[j].elementID
			newEachID := compID + "_" + oldEachID
			processedHTML = strings.ReplaceAll(processedHTML, `id="`+oldEachID+`_anchor"`, `id="`+newEachID+`_anchor"`)
			ifb.branches[i].eachBlocks[j].elementID = newEachID
		}
		processedBranches[i].condition = branch.condition
		processedBranches[i].html = processedHTML
		processedBranches[i].exprs = exprs
		allExprs = append(allExprs, exprs...)
	}

	var elseHTML string
	var elseExprs []exprBinding
	if ifb.elseHTML != "" {
		elseHTML, elseExprs = processIfBranchExpressions(ifb.elseHTML, anchorID, &exprCounter)
		allExprs = append(allExprs, elseExprs...)
	}

	fmt.Fprintf(sb, "%supdate%s := func() {\n", indent, anchorID)
	fmt.Fprintf(sb, "%s\tvar html string\n", indent)
	fmt.Fprintf(sb, "%s\tvar branchIdx int\n", indent)

	for i, branch := range processedBranches {
		cond := transformCondition(ifb.branches[i].condition, fieldTypes, compID)
		if i == 0 {
			fmt.Fprintf(sb, "%s\tif %s {\n%s\t\thtml = %s\n%s\t\tbranchIdx = %d\n", indent, cond, indent, escapeForGoString(branch.html), indent, i)
		} else {
			fmt.Fprintf(sb, "%s\t} else if %s {\n%s\t\thtml = %s\n%s\t\tbranchIdx = %d\n", indent, cond, indent, escapeForGoString(branch.html), indent, i)
		}
	}

	if elseHTML != "" {
		fmt.Fprintf(sb, "%s\t} else {\n%s\t\thtml = %s\n%s\t\tbranchIdx = -1\n%s\t}\n", indent, indent, escapeForGoString(elseHTML), indent, indent)
	} else {
		fmt.Fprintf(sb, "%s\t} else {\n%s\t\tbranchIdx = -1\n%s\t}\n", indent, indent, indent)
	}

	// Replace nested component placeholders with their HTML constants (using prefixed IDs)
	for _, comp := range compsInBlock {
		prefixedID := parentCompID + "_" + comp.elementID
		fmt.Fprintf(sb, "%s\thtml = strings.Replace(html, \"<!--%s-->\", %sHTML, 1)\n", indent, comp.elementID, prefixedID)
	}

	fmt.Fprintf(sb, "%s\tnewEl := document.Call(\"createElement\", \"span\")\n", indent)
	fmt.Fprintf(sb, "%s\tnewEl.Set(\"innerHTML\", html)\n", indent)
	fmt.Fprintf(sb, "%s\tif !%s_current.IsNull() { %s_current.Call(\"remove\") }\n", indent, anchorID, anchorID)
	fmt.Fprintf(sb, "%s\tif !%s_anchor.IsNull() { %s_anchor.Get(\"parentNode\").Call(\"insertBefore\", newEl, %s_anchor) }\n", indent, anchorID, anchorID, anchorID)
	fmt.Fprintf(sb, "%s\t%s_current = newEl\n", indent, anchorID)

	// Wire up expressions in the if-block branches
	for _, expr := range allExprs {
		valueType := fieldTypes[expr.fieldName]
		if valueType == "" {
			valueType = "string" // default to string
		}
		fmt.Fprintf(sb, "%s\t%s_el := document.Call(\"getElementById\", \"%s\")\n", indent, expr.elementID, expr.elementID)
		fmt.Fprintf(sb, "%s\tif !%s_el.IsNull() && !%s_el.IsUndefined() {\n", indent, expr.elementID, expr.elementID)
		jsConv := toJS(valueType, compID+"."+expr.fieldName+".Get()")
		fmt.Fprintf(sb, "%s\t\treactive.SetText(%s_el, %s)\n", indent, expr.elementID, jsConv)
		fmt.Fprintf(sb, "%s\t}\n", indent)
	}

	// Wire up nested components after inserting HTML (using prefixed IDs)
	for _, comp := range compsInBlock {
		prefixedID := parentCompID + "_" + comp.elementID
		generateComponentInline(sb, comp, childComponents, indent+"\t", prefixedID, false)
	}

	// Wire up each blocks inside if branches AFTER inserting HTML into DOM
	for i, branch := range ifb.branches {
		if len(branch.eachBlocks) > 0 {
			fmt.Fprintf(sb, "%s\tif branchIdx == %d {\n", indent, i)
			for _, each := range branch.eachBlocks {
				generateEachBlockInlineForChild(sb, each, fieldTypes, compID, indent+"\t\t")
			}
			fmt.Fprintf(sb, "%s\t}\n", indent)
		}
	}
	fmt.Fprintf(sb, "%s\t_ = branchIdx\n", indent)

	fmt.Fprintf(sb, "%s}\n", indent)

	for _, dep := range ifb.deps {
		fmt.Fprintf(sb, "%s%s.%s.OnChange(func(_ %s) { update%s() })\n", indent, compID, dep, fieldTypes[dep], anchorID)
	}
	fmt.Fprintf(sb, "%supdate%s()\n", indent, anchorID)
}

// processIfBranchExpressions finds {Field} expressions in HTML and replaces them with spans
// It skips expressions inside <code>, <pre>, or <script> tags
func processIfBranchExpressions(html string, prefix string, counter *int) (string, []exprBinding) {
	var exprs []exprBinding
	result := html

	// Find all {Field} patterns
	matches := findFieldExprs(result)

	// Process in reverse to preserve indices
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]

		// Check if this match is inside a <code>, <pre>, or <script> tag
		if isInsideCodeBlock(result, match.start) {
			continue
		}

		elementID := fmt.Sprintf("%s_expr%d", prefix, *counter)
		*counter++

		exprs = append([]exprBinding{{
			fieldName: match.field,
			elementID: elementID,
			isHTML:    false,
		}}, exprs...)

		replacement := fmt.Sprintf(`<span id="%s"></span>`, elementID)
		result = result[:match.start] + replacement + result[match.end:]
	}

	return result, exprs
}

// isInsideCodeBlock checks if position is inside <code>, <pre>, or <script> tags
func isInsideCodeBlock(html string, pos int) bool {
	// Look backwards for opening/closing tags
	before := html[:pos]

	// Check for each tag type
	for _, tag := range []string{"code", "pre", "script"} {
		openTag := "<" + tag
		closeTag := "</" + tag

		lastOpen := strings.LastIndex(strings.ToLower(before), openTag)
		lastClose := strings.LastIndex(strings.ToLower(before), closeTag)

		// If we found an opening tag after the last closing tag, we're inside
		if lastOpen > lastClose {
			return true
		}
	}

	return false
}

// generateNestedComponentConstants generates CSS and HTML constants for components nested inside child templates
func generateNestedComponentConstants(sb *strings.Builder, nestedComponents []componentBinding, ifBlocks []ifBinding, childComponents map[string]*component, cssGenerated map[string]bool, htmlGenerated map[string]bool, parentCompID string) {
	// Generate HTML constants for ALL nested components
	// These are components inside a parent that's already in an if-block, so they all need HTML constants
	for _, nestedBinding := range nestedComponents {
		nestedComp := childComponents[nestedBinding.name]
		if nestedComp == nil {
			continue
		}

		// nestedBinding.elementID is already prefixed (e.g., comp1_comp0)
		prefixedID := nestedBinding.elementID

		// Skip if already generated
		if htmlGenerated[prefixedID] {
			continue
		}
		htmlGenerated[prefixedID] = true

		// Generate CSS constant if not already done
		if nestedComp.style != "" && !cssGenerated[nestedBinding.name] {
			cssGenerated[nestedBinding.name] = true
			fmt.Fprintf(sb, "const %sCSS = `%s`\n", strings.ToLower(nestedBinding.name), nestedComp.style)
		}

		// Generate HTML constant for this nested component
		nestedTmpl := strings.ReplaceAll(nestedComp.template, "<slot/>", nestedBinding.children)
		nestedTmplProcessed, nestedBindings := parseTemplate(nestedTmpl)

		// Prefix IDs to avoid conflicts
		for _, expr := range nestedBindings.expressions {
			oldID := expr.elementID
			newID := prefixedID + "_" + oldID
			nestedTmplProcessed = strings.ReplaceAll(nestedTmplProcessed, `id="`+oldID+`"`, `id="`+newID+`"`)
		}
		for _, evt := range nestedBindings.events {
			oldID := evt.elementID
			newID := prefixedID + "_" + oldID
			nestedTmplProcessed = strings.ReplaceAll(nestedTmplProcessed, `id="`+oldID+`"`, `id="`+newID+`"`)
		}
		for _, ifb := range nestedBindings.ifBlocks {
			oldID := ifb.elementID
			newID := prefixedID + "_" + oldID
			nestedTmplProcessed = strings.ReplaceAll(nestedTmplProcessed, `id="`+oldID+`_anchor"`, `id="`+newID+`_anchor"`)
		}
		for _, bind := range nestedBindings.bindings {
			oldID := bind.elementID
			newID := prefixedID + "_" + oldID
			nestedTmplProcessed = strings.ReplaceAll(nestedTmplProcessed, `id="`+oldID+`"`, `id="`+newID+`"`)
		}
		for _, cb := range nestedBindings.classBindings {
			oldID := cb.elementID
			newID := prefixedID + "_" + oldID
			nestedTmplProcessed = strings.ReplaceAll(nestedTmplProcessed, `id="`+oldID+`"`, `id="`+newID+`"`)
		}
		for _, each := range nestedBindings.eachBlocks {
			oldID := each.elementID
			newID := prefixedID + "_" + oldID
			nestedTmplProcessed = strings.ReplaceAll(nestedTmplProcessed, `id="`+oldID+`_anchor"`, `id="`+newID+`_anchor"`)
		}
		for _, ab := range nestedBindings.attrBindings {
			oldID := ab.elementID
			newID := prefixedID + "_" + oldID
			nestedTmplProcessed = strings.ReplaceAll(nestedTmplProcessed, `data-attrbind="`+oldID+`"`, `data-attrbind="`+newID+`"`)
		}

		nestedTmplProcessed = injectIDIntoFirstTag(nestedTmplProcessed, prefixedID)
		fmt.Fprintf(sb, "const %sHTML = %s\n", prefixedID, escapeForGoString(nestedTmplProcessed))

		// Recursively process nested components inside this one
		generateNestedComponentConstants(sb, nestedBindings.components, nestedBindings.ifBlocks, childComponents, cssGenerated, htmlGenerated, prefixedID)
	}
}

func generateChildIfBlocks(sb *strings.Builder, ifBlocks []ifBinding, fieldTypes map[string]string, compID string) {
	for _, ifb := range ifBlocks {
		// elementID is already prefixed with compID from the renaming step
		anchorID := ifb.elementID

		fmt.Fprintf(sb, "\t%s_anchor := document.Call(\"getElementById\", \"%s_anchor\")\n", anchorID, anchorID)
		fmt.Fprintf(sb, "\t%s_current := js.Null()\n", anchorID)

		fmt.Fprintf(sb, "\tupdate%s := func() {\n", anchorID)
		sb.WriteString("\t\tvar html string\n")
		sb.WriteString("\t\tvar branchIdx int\n")

		for i, branch := range ifb.branches {
			cond := transformCondition(branch.condition, fieldTypes, compID)
			if i == 0 {
				fmt.Fprintf(sb, "\t\tif %s {\n\t\t\thtml = %s\n\t\t\tbranchIdx = %d\n", cond, escapeForGoString(branch.html), i)
			} else {
				fmt.Fprintf(sb, "\t\t} else if %s {\n\t\t\thtml = %s\n\t\t\tbranchIdx = %d\n", cond, escapeForGoString(branch.html), i)
			}
		}

		if ifb.elseHTML != "" {
			fmt.Fprintf(sb, "\t\t} else {\n\t\t\thtml = %s\n\t\t\tbranchIdx = -1\n\t\t}\n", escapeForGoString(ifb.elseHTML))
		} else {
			sb.WriteString("\t\t} else {\n\t\t\tbranchIdx = -1\n\t\t}\n")
		}

		fmt.Fprintf(sb, "\t\tnewEl := document.Call(\"createElement\", \"span\")\n")
		fmt.Fprintf(sb, "\t\tnewEl.Set(\"innerHTML\", html)\n")
		fmt.Fprintf(sb, "\t\tif !%s_current.IsNull() { %s_current.Call(\"remove\") }\n", anchorID, anchorID)
		fmt.Fprintf(sb, "\t\tif !%s_anchor.IsNull() { %s_anchor.Get(\"parentNode\").Call(\"insertBefore\", newEl, %s_anchor) }\n", anchorID, anchorID, anchorID)
		fmt.Fprintf(sb, "\t\t%s_current = newEl\n", anchorID)

		// Wire up each blocks inside if branches AFTER inserting HTML into DOM
		for i, branch := range ifb.branches {
			if len(branch.eachBlocks) > 0 {
				fmt.Fprintf(sb, "\t\tif branchIdx == %d {\n", i)
				for _, each := range branch.eachBlocks {
					generateEachBlockInlineForChild(sb, each, fieldTypes, compID, "\t\t\t")
				}
				sb.WriteString("\t\t}\n")
			}
		}
		sb.WriteString("\t\t_ = branchIdx\n")

		sb.WriteString("\t}\n")

		for _, dep := range ifb.deps {
			fmt.Fprintf(sb, "\t%s.%s.OnChange(func(_ %s) { update%s() })\n", compID, dep, fieldTypes[dep], anchorID)
		}
		fmt.Fprintf(sb, "\tupdate%s()\n\n", anchorID)
	}
}
