package main

import (
	"fmt"
	"strings"
)

type Event struct {
	ID        string
	Event     string
	Handler   string   // Full call: "increment()" or "remove(item)"
	Modifiers []string // preventDefault, stopPropagation
	InBranch  bool
	InEach    bool
	EachID    string   // which each block this is in
}

type Binding struct {
	ID       string
	Attr     string // "value", "checked"
	VarName  string
	VarType  string
	InBranch bool
}

type ClassBinding struct {
	ID        string
	ClassName string
	Cond      string
	InBranch  bool
}

func generate(p Parsed, a *Analysis) (goCode string, html string) {
	var b strings.Builder

	// Flatten all components recursively (including nested)
	allComponents := make(map[string]*Parsed)
	flattenComponents(p.Components, allComponents)

	// Build component names map
	compNames := make(map[string]bool)
	for name := range allComponents {
		compNames[name] = true
	}

	// Parse template into AST with component awareness
	tmplAST := parseTemplateWithComponents(p.Template, compNames)
	exprs := tmplAST.CollectExprs()
	htmls := tmplAST.CollectHtmls()
	ifs := tmplAST.CollectIfs()
	eaches := tmplAST.CollectEaches()
	components := tmplAST.CollectComponents()

	// Collect ALL component usages including nested ones
	allComponentUsages := collectAllComponentUsages(components, allComponents)

	// Analyze each component
	compAnalyses := make(map[string]*Analysis)
	for name, comp := range allComponents {
		compA, _ := analyze(comp.Script)
		compAnalyses[name] = compA
	}

	// Replace p.Components with flattened version for HTML generation
	p.Components = allComponents

	// Collect events and bindings with branch info
	events := collectEvents(p.Template, ifs, eaches)
	inputBindings := collectBindings(p.Template, p.Script, ifs)
	classBindings := collectClassBindings(p.Template, ifs)

	// Process branch templates - replace @click and bind: with IDs
	branchTemplates := make(map[string]string)
	btnID := countEventsOutsideIf(p.Template, ifs)
	inputID := countBindingsOutsideIf(p.Template, ifs)
	classID := countClassBindingsOutsideIf(p.Template, ifs)
	
	for _, ifn := range ifs {
		thenHTML := generateBranchHTML(ifn.Then, tmplAST)
		thenHTML, btnID, inputID, classID = processTemplateAttrs(thenHTML, btnID, inputID, classID)
		branchTemplates[ifn.CondID+"_then"] = thenHTML
		
		if len(ifn.Else) > 0 {
			elseHTML := generateBranchHTML(ifn.Else, tmplAST)
			elseHTML, btnID, inputID, classID = processTemplateAttrs(elseHTML, btnID, inputID, classID)
			branchTemplates[ifn.CondID+"_else"] = elseHTML
		}
	}

	// Process each templates
	eachTemplates := make(map[string]string)
	for _, each := range eaches {
		bodyHTML := generateEachBodyHTML(each.Body, tmplAST, each.ID)
		// Process @click etc with indexed IDs
		bodyHTML = processEachTemplateAttrs(bodyHTML)
		eachTemplates[each.ID] = bodyHTML
	}

	b.WriteString("//go:build js && wasm\n\n")
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"syscall/js\"\n")
	b.WriteString("\t\"strconv\"\n")
	if len(eaches) > 0 {
		b.WriteString("\t\"strings\"\n")
	}
	b.WriteString(")\n\n")

	b.WriteString("var document = js.Global().Get(\"document\")\n")
	b.WriteString("var propagating = false\n\n")

	b.WriteString("var _ = strconv.Atoi\n\n")

	// State vars
	b.WriteString("// State\n")
	b.WriteString("var (\n")
	for _, name := range a.Order {
		v := a.Vars[name]
		if len(v.DependsOn) == 0 {
			fmt.Fprintf(&b, "\t%s = %s\n", name, getInitializer(p.Script, name))
		} else {
			fmt.Fprintf(&b, "\t%s %s\n", name, inferType(p.Script, name))
		}
	}
	b.WriteString(")\n\n")

	// DOM refs
	if len(exprs) > 0 || len(htmls) > 0 || len(ifs) > 0 || len(eaches) > 0 {
		b.WriteString("// DOM refs\n")
		b.WriteString("var (\n")
		for _, expr := range exprs {
			if !isExprInEach(expr.ID, eaches) {
				fmt.Fprintf(&b, "\t%sEl js.Value\n", expr.ID)
			}
		}
		for _, html := range htmls {
			fmt.Fprintf(&b, "\t%sEl js.Value\n", html.ID)
		}
		for _, ifn := range ifs {
			fmt.Fprintf(&b, "\t%s_anchor js.Value\n", ifn.CondID)
			fmt.Fprintf(&b, "\t%s_current js.Value\n", ifn.CondID)
			fmt.Fprintf(&b, "\t%s_showing string\n", ifn.CondID)
		}
		for _, each := range eaches {
			fmt.Fprintf(&b, "\t%s_anchor js.Value\n", each.ID)
			fmt.Fprintf(&b, "\t%s_count int\n", each.ID)
		}
		b.WriteString(")\n\n")
	}

	// Branch templates
	if len(ifs) > 0 || len(eaches) > 0 {
		b.WriteString("// Branch templates\n")
		b.WriteString("var (\n")
		for key, tmpl := range branchTemplates {
			// Escape backticks in template
			tmpl = strings.ReplaceAll(tmpl, "`", "` + \"`\" + `")
			fmt.Fprintf(&b, "\t%s_tmpl = `%s`\n", key, tmpl)
		}
		for key, tmpl := range eachTemplates {
			tmpl = strings.ReplaceAll(tmpl, "`", "` + \"`\" + `")
			fmt.Fprintf(&b, "\t%s_tmpl = `%s`\n", key, tmpl)
		}
		b.WriteString(")\n\n")

		b.WriteString("func createFragment(html string) js.Value {\n")
		b.WriteString("\ttmpl := document.Call(\"createElement\", \"template\")\n")
		b.WriteString("\ttmpl.Set(\"innerHTML\", html)\n")
		b.WriteString("\treturn tmpl.Get(\"content\")\n")
		b.WriteString("}\n\n")
	}

	// Setters
	b.WriteString("// Setters\n")
	for _, name := range a.Order {
		if len(a.Vars[name].DependsOn) == 0 {
			typ := inferType(p.Script, name)
			fmt.Fprintf(&b, "func set%s(v %s) {\n", capitalize(name), typ)
			fmt.Fprintf(&b, "\t%s = v\n", name)
			b.WriteString("\tpropagate()\n")
			b.WriteString("}\n\n")
		}
	}

	// Propagate
	b.WriteString("func propagate() {\n")
	b.WriteString("\tif propagating {\n")
	b.WriteString("\t\treturn\n")
	b.WriteString("\t}\n")
	b.WriteString("\tpropagating = true\n\n")

	hasDerived := false
	for _, name := range a.Order {
		if len(a.Vars[name].DependsOn) > 0 {
			if !hasDerived {
				b.WriteString("\t// Recompute derived values\n")
				hasDerived = true
			}
			fmt.Fprintf(&b, "\t%s = %s\n", name, getExpression(p.Script, name))
		}
	}

	b.WriteString("\n\tupdateDOM()\n")
	b.WriteString("\tpropagating = false\n")
	b.WriteString("}\n\n")

	// Update DOM
	b.WriteString("func updateDOM() {\n")

	// Update if blocks FIRST (so elements exist)
	for _, ifn := range ifs {
		hasElse := len(ifn.Else) > 0
		fmt.Fprintf(&b, "\tif %s {\n", ifn.Cond)
		fmt.Fprintf(&b, "\t\tif %s_showing != \"then\" {\n", ifn.CondID)
		fmt.Fprintf(&b, "\t\t\tif !%s_current.IsUndefined() && !%s_current.IsNull() {\n", ifn.CondID, ifn.CondID)
		fmt.Fprintf(&b, "\t\t\t\t%s_current.Call(\"remove\")\n", ifn.CondID)
		b.WriteString("\t\t\t}\n")
		fmt.Fprintf(&b, "\t\t\tfrag := createFragment(%s_then_tmpl)\n", ifn.CondID)
		fmt.Fprintf(&b, "\t\t\twrapper := document.Call(\"createElement\", \"span\")\n")
		fmt.Fprintf(&b, "\t\t\twrapper.Call(\"appendChild\", frag)\n")
		fmt.Fprintf(&b, "\t\t\t%s_anchor.Get(\"parentNode\").Call(\"insertBefore\", wrapper, %s_anchor)\n", ifn.CondID, ifn.CondID)
		fmt.Fprintf(&b, "\t\t\t%s_current = wrapper\n", ifn.CondID)
		fmt.Fprintf(&b, "\t\t\t%s_showing = \"then\"\n", ifn.CondID)
		b.WriteString("\t\t\tbindBranchElements()\n")
		b.WriteString("\t\t}\n")
		b.WriteString("\t} else {\n")
		if hasElse {
			fmt.Fprintf(&b, "\t\tif %s_showing != \"else\" {\n", ifn.CondID)
		} else {
			fmt.Fprintf(&b, "\t\tif %s_showing != \"none\" {\n", ifn.CondID)
		}
		fmt.Fprintf(&b, "\t\t\tif !%s_current.IsUndefined() && !%s_current.IsNull() {\n", ifn.CondID, ifn.CondID)
		fmt.Fprintf(&b, "\t\t\t\t%s_current.Call(\"remove\")\n", ifn.CondID)
		b.WriteString("\t\t\t}\n")
		if hasElse {
			fmt.Fprintf(&b, "\t\t\tfrag := createFragment(%s_else_tmpl)\n", ifn.CondID)
			fmt.Fprintf(&b, "\t\t\twrapper := document.Call(\"createElement\", \"span\")\n")
			fmt.Fprintf(&b, "\t\t\twrapper.Call(\"appendChild\", frag)\n")
			fmt.Fprintf(&b, "\t\t\t%s_anchor.Get(\"parentNode\").Call(\"insertBefore\", wrapper, %s_anchor)\n", ifn.CondID, ifn.CondID)
			fmt.Fprintf(&b, "\t\t\t%s_current = wrapper\n", ifn.CondID)
			fmt.Fprintf(&b, "\t\t\t%s_showing = \"else\"\n", ifn.CondID)
			b.WriteString("\t\t\tbindBranchElements()\n")
		} else {
			fmt.Fprintf(&b, "\t\t\t%s_current = js.Null()\n", ifn.CondID)
			fmt.Fprintf(&b, "\t\t\t%s_showing = \"none\"\n", ifn.CondID)
		}
		b.WriteString("\t\t}\n")
		b.WriteString("\t}\n")
	}

	// Update each blocks
	for _, each := range eaches {
		indexVar := each.Index
		if indexVar == "" {
			indexVar = "_i"
		}
		
		// Get events in this each block
		eachEvts := filterEachEvents(events, each.ID)
		
		fmt.Fprintf(&b, "\t// Update %s\n", each.ID)
		fmt.Fprintf(&b, "\tnewLen%s := len(%s)\n", each.ID, each.Array)
		fmt.Fprintf(&b, "\toldLen%s := %s_count\n", each.ID, each.ID)
		fmt.Fprintf(&b, "\t%s_count = newLen%s\n", each.ID, each.ID)
		
		// Remove excess items
		fmt.Fprintf(&b, "\tfor i := newLen%s; i < oldLen%s; i++ {\n", each.ID, each.ID)
		fmt.Fprintf(&b, "\t\tif el := document.Call(\"getElementById\", \"%s_\" + strconv.Itoa(i)); !el.IsNull() {\n", each.ID)
		fmt.Fprintf(&b, "\t\t\tel.Call(\"remove\")\n")
		b.WriteString("\t\t}\n")
		b.WriteString("\t}\n")
		
		// Add/update items
		fmt.Fprintf(&b, "\tfor %s := 0; %s < newLen%s; %s++ {\n", indexVar, indexVar, each.ID, indexVar)
		fmt.Fprintf(&b, "\t\t%s := %s[%s]\n", each.Item, each.Array, indexVar)
		fmt.Fprintf(&b, "\t\telID := \"%s_\" + strconv.Itoa(%s)\n", each.ID, indexVar)
		fmt.Fprintf(&b, "\t\tif el := document.Call(\"getElementById\", elID); el.IsNull() {\n")
		fmt.Fprintf(&b, "\t\t\thtml := %s_tmpl\n", each.ID)
		
		// Replace placeholders in template with actual index
		fmt.Fprintf(&b, "\t\t\thtml = strings.ReplaceAll(html, \"${_idx}\", strconv.Itoa(%s))\n", indexVar)
		
		fmt.Fprintf(&b, "\t\t\tfrag := createFragment(html)\n")
		fmt.Fprintf(&b, "\t\t\twrapper := document.Call(\"createElement\", \"span\")\n")
		fmt.Fprintf(&b, "\t\t\twrapper.Set(\"id\", elID)\n")
		fmt.Fprintf(&b, "\t\t\twrapper.Call(\"appendChild\", frag)\n")
		fmt.Fprintf(&b, "\t\t\t%s_anchor.Get(\"parentNode\").Call(\"insertBefore\", wrapper, %s_anchor)\n", each.ID, each.ID)
		
		// Bind events in newly created element
		for i, ev := range eachEvts {
			// Use unique ID per loop iteration
			fmt.Fprintf(&b, "\t\t\tif btn := wrapper.Call(\"querySelector\", \"#btn%d_${_idx}\".replace(\"${_idx}\", strconv.Itoa(%s))); !btn.IsNull() || true {\n", i, indexVar)
			// Find button by ID within wrapper
			fmt.Fprintf(&b, "\t\t\t\tbtnEl := document.Call(\"getElementById\", \"btn%d_\" + strconv.Itoa(%s))\n", i, indexVar)
			fmt.Fprintf(&b, "\t\t\t\tif !btnEl.IsNull() && btnEl.Get(\"_bound\").IsUndefined() {\n")
			fmt.Fprintf(&b, "\t\t\t\t\tbtnEl.Set(\"_bound\", true)\n")
			// Capture loop variables
			fmt.Fprintf(&b, "\t\t\t\t\tidx := %s\n", indexVar)
			fmt.Fprintf(&b, "\t\t\t\t\t_ = idx\n")
			if each.Item != "" {
				fmt.Fprintf(&b, "\t\t\t\t\titemVal := %s\n", each.Item)
				fmt.Fprintf(&b, "\t\t\t\t\t_ = itemVal\n")
			}
			fmt.Fprintf(&b, "\t\t\t\t\tbtnEl.Call(\"addEventListener\", \"%s\", js.FuncOf(func(this js.Value, args []js.Value) any {\n", ev.Event)
			generateEventModifiers(&b, ev.Modifiers, "\t\t\t\t\t\t")
			// Replace loop vars in handler with captured vars
			handler := ev.Handler
			if each.Index != "" { handler = replaceVar(handler, each.Index, "idx") }
			
			if each.Item != "" { handler = replaceVar(handler, each.Item, "itemVal") }
			
			fmt.Fprintf(&b, "\t\t\t\t\t\t%s\n", handler)
			b.WriteString("\t\t\t\t\t\treturn nil\n")
			b.WriteString("\t\t\t\t\t}))\n")
			b.WriteString("\t\t\t\t}\n")
			b.WriteString("\t\t\t}\n")
		}
		
		b.WriteString("\t\t}\n")
		
		// Update expressions inside this each
		eachExprs := collectExprsInEach(each)
		for _, expr := range eachExprs {
			fmt.Fprintf(&b, "\t\tif exprEl := document.Call(\"getElementById\", \"%s_\" + strconv.Itoa(%s)); !exprEl.IsNull() {\n", expr.ID, indexVar)
			fmt.Fprintf(&b, "\t\t\texprEl.Set(\"textContent\", %s)\n", expr.Expr)
			b.WriteString("\t\t}\n")
		}
		
		b.WriteString("\t}\n")
	}

	// Update expressions AFTER branches exist (skip expressions in each blocks)
	for _, expr := range exprs {
		if isExprInEach(expr.ID, eaches) {
			continue
		}
		typ := inferType(p.Script, expr.Expr)
		fmt.Fprintf(&b, "\tif !%sEl.IsUndefined() && !%sEl.IsNull() {\n", expr.ID, expr.ID)
		if typ == "int" {
			fmt.Fprintf(&b, "\t\t%sEl.Set(\"textContent\", strconv.Itoa(%s))\n", expr.ID, expr.Expr)
		} else if typ == "bool" {
			fmt.Fprintf(&b, "\t\t%sEl.Set(\"textContent\", strconv.FormatBool(%s))\n", expr.ID, expr.Expr)
		} else {
			fmt.Fprintf(&b, "\t\t%sEl.Set(\"textContent\", %s)\n", expr.ID, expr.Expr)
		}
		b.WriteString("\t}\n")
	}

	// Update html nodes (innerHTML)
	for _, html := range htmls {
		fmt.Fprintf(&b, "\tif !%sEl.IsUndefined() && !%sEl.IsNull() {\n", html.ID, html.ID)
		fmt.Fprintf(&b, "\t\t%sEl.Set(\"innerHTML\", %s)\n", html.ID, html.Expr)
		b.WriteString("\t}\n")
	}

	// Sync input values
	for _, ib := range inputBindings {
		if ib.Attr == "checked" {
			fmt.Fprintf(&b, "\tif el := document.Call(\"getElementById\", \"%s\"); !el.IsNull() && el.Get(\"checked\").Bool() != %s {\n", ib.ID, ib.VarName)
			fmt.Fprintf(&b, "\t\tel.Set(\"checked\", %s)\n", ib.VarName)
			b.WriteString("\t}\n")
		} else {
			if ib.VarType == "int" {
				fmt.Fprintf(&b, "\tif el := document.Call(\"getElementById\", \"%s\"); !el.IsNull() && el.Get(\"value\").String() != strconv.Itoa(%s) {\n", ib.ID, ib.VarName)
				fmt.Fprintf(&b, "\t\tel.Set(\"value\", strconv.Itoa(%s))\n", ib.VarName)
			} else {
				fmt.Fprintf(&b, "\tif el := document.Call(\"getElementById\", \"%s\"); !el.IsNull() && el.Get(\"value\").String() != %s {\n", ib.ID, ib.VarName)
				fmt.Fprintf(&b, "\t\tel.Set(\"value\", %s)\n", ib.VarName)
			}
			b.WriteString("\t}\n")
		}
	}

	// Update class bindings
	for _, cb := range classBindings {
		fmt.Fprintf(&b, "\tif el := document.Call(\"getElementById\", \"%s\"); !el.IsNull() {\n", cb.ID)
		fmt.Fprintf(&b, "\t\tif %s {\n", cb.Cond)
		fmt.Fprintf(&b, "\t\t\tel.Get(\"classList\").Call(\"add\", \"%s\")\n", cb.ClassName)
		b.WriteString("\t\t} else {\n")
		fmt.Fprintf(&b, "\t\t\tel.Get(\"classList\").Call(\"remove\", \"%s\")\n", cb.ClassName)
		b.WriteString("\t\t}\n")
		b.WriteString("\t}\n")
	}

	// Update component dynamic props
	for _, comp := range allComponentUsages {
		for propName, propVal := range comp.Props {
			if strings.HasPrefix(propVal, "{") && strings.HasSuffix(propVal, "}") {
				// Dynamic prop - extract expression
				expr := propVal[1 : len(propVal)-1]
				varName := strings.ReplaceAll(comp.ID, "-", "_")
				fmt.Fprintf(&b, "\tif %s.%s != %s {\n", varName, propName, expr)
				fmt.Fprintf(&b, "\t\t%s.%s = %s\n", varName, propName, expr)
				fmt.Fprintf(&b, "\t\t%s.propagate()\n", varName)
				b.WriteString("\t}\n")
			}
		}
		// Sync binding values (parent -> child)
		for propName, parentVar := range comp.Bindings {
			varName := strings.ReplaceAll(comp.ID, "-", "_")
			fmt.Fprintf(&b, "\tif %s.%s != %s {\n", varName, capitalize(propName), parentVar)
			fmt.Fprintf(&b, "\t\t%s.%s = %s\n", varName, capitalize(propName), parentVar)
			fmt.Fprintf(&b, "\t\t%s.propagate()\n", varName)
			b.WriteString("\t}\n")
		}
	}
	b.WriteString("}\n\n")

	// User functions
	b.WriteString("// Functions\n")
	for name, fn := range a.Funcs {
		transformed := transformFunction(p.Script, name, fn.Modifies)
		b.WriteString(transformed)
		b.WriteString("\n\n")
	}

	// Generate component types and instances
	if len(allComponentUsages) > 0 {
		generateComponentCode(&b, p, allComponentUsages, compAnalyses)
	}

	// Main - only bind static elements
	b.WriteString("func main() {\n")

	// Get static DOM refs
	hasStaticExprs := false
	for _, expr := range exprs {
		if !isExprInIf(expr.ID, ifs) && !isExprInEach(expr.ID, eaches) {
			if !hasStaticExprs {
				b.WriteString("\t// Get static DOM refs\n")
				hasStaticExprs = true
			}
			fmt.Fprintf(&b, "\t%sEl = document.Call(\"getElementById\", \"%s\")\n", expr.ID, expr.ID)
		}
	}
	for _, html := range htmls {
		if !hasStaticExprs {
			b.WriteString("\t// Get static DOM refs\n")
			hasStaticExprs = true
		}
		fmt.Fprintf(&b, "\t%sEl = document.Call(\"getElementById\", \"%s\")\n", html.ID, html.ID)
	}

	// Find anchors
	if len(ifs) > 0 {
		b.WriteString("\n\t// Find if block anchors\n")
		for _, ifn := range ifs {
			fmt.Fprintf(&b, "\t%s_anchor = document.Call(\"getElementById\", \"%s_anchor\")\n", ifn.CondID, ifn.CondID)
		}
	}

	if len(eaches) > 0 {
		b.WriteString("\n\t// Find each block anchors\n")
		for _, each := range eaches {
			fmt.Fprintf(&b, "\t%s_anchor = document.Call(\"getElementById\", \"%s_anchor\")\n", each.ID, each.ID)
		}
	}

	// Bind static events only
	staticEvents := filterStaticEvents(events)
	if len(staticEvents) > 0 {
		b.WriteString("\n\t// Bind static events\n")
		for _, ev := range staticEvents {
			fmt.Fprintf(&b, "\tdocument.Call(\"getElementById\", \"%s\").Call(\"addEventListener\", \"%s\", js.FuncOf(func(this js.Value, args []js.Value) any {\n", ev.ID, ev.Event)
			generateEventModifiers(&b, ev.Modifiers, "\t\t")
			fmt.Fprintf(&b, "\t\t%s\n", ev.Handler)
			b.WriteString("\t\treturn nil\n")
			b.WriteString("\t}))\n")
		}
	}

	// Bind static inputs only
	staticBindings := filterStaticBindings(inputBindings)
	if len(staticBindings) > 0 {
		b.WriteString("\n\t// Bind static inputs\n")
		for _, ib := range staticBindings {
			eventType := "input"
			if ib.Attr == "checked" {
				eventType = "change"
			}
			fmt.Fprintf(&b, "\tdocument.Call(\"getElementById\", \"%s\").Call(\"addEventListener\", \"%s\", js.FuncOf(func(this js.Value, args []js.Value) any {\n", ib.ID, eventType)
			if ib.Attr == "checked" {
				fmt.Fprintf(&b, "\t\tset%s(this.Get(\"checked\").Bool())\n", capitalize(ib.VarName))
			} else if ib.VarType == "int" {
				fmt.Fprintf(&b, "\t\tif v, err := strconv.Atoi(this.Get(\"value\").String()); err == nil {\n")
				fmt.Fprintf(&b, "\t\t\tset%s(v)\n", capitalize(ib.VarName))
				b.WriteString("\t\t}\n")
			} else {
				fmt.Fprintf(&b, "\t\tset%s(this.Get(\"value\").String())\n", capitalize(ib.VarName))
			}
			b.WriteString("\t\treturn nil\n")
			b.WriteString("\t}))\n")
		}
	}

	// Setup component binding callbacks
	hasBindings := false
	for _, comp := range allComponentUsages {
		for propName, parentVar := range comp.Bindings {
			if !hasBindings {
				b.WriteString("\n\t// Setup component binding callbacks\n")
				hasBindings = true
			}
			varName := strings.ReplaceAll(comp.ID, "-", "_")
			fmt.Fprintf(&b, "\t%s._on%sChange = func(v %s) { set%s(v) }\n", 
				varName, capitalize(propName), inferType(p.Script, parentVar), capitalize(parentVar))
		}
	}

	// Mount components
	if len(allComponentUsages) > 0 {
		b.WriteString("\n\t// Mount components\n")
		for _, comp := range allComponentUsages {
			fmt.Fprintf(&b, "\t%s.mount()\n", strings.ReplaceAll(comp.ID, "-", "_"))
		}
	}

	b.WriteString("\n\tpropagate()\n")
	// Call onMount if defined in root app
	if _, hasOnMount := a.Funcs["onMount"]; hasOnMount {
		b.WriteString("\tonMount()\n")
	}
	b.WriteString("\tselect {}\n")
	b.WriteString("}\n\n")

	// bindBranchElements - only bind branch elements
	b.WriteString("func bindBranchElements() {\n")
	
	// Re-bind expression refs in branches
	for _, expr := range exprs {
		if isExprInIf(expr.ID, ifs) {
			fmt.Fprintf(&b, "\tif el := document.Call(\"getElementById\", \"%s\"); !el.IsNull() {\n", expr.ID)
			fmt.Fprintf(&b, "\t\t%sEl = el\n", expr.ID)
			b.WriteString("\t}\n")
		}
	}
	
	// Bind events in branches
	branchEvents := filterBranchEvents(events)
	for _, ev := range branchEvents {
		fmt.Fprintf(&b, "\tif el := document.Call(\"getElementById\", \"%s\"); !el.IsNull() && el.Get(\"_bound\").IsUndefined() {\n", ev.ID)
		fmt.Fprintf(&b, "\t\tel.Set(\"_bound\", true)\n")
		fmt.Fprintf(&b, "\t\tel.Call(\"addEventListener\", \"%s\", js.FuncOf(func(this js.Value, args []js.Value) any {\n", ev.Event)
		generateEventModifiers(&b, ev.Modifiers, "\t\t\t")
		fmt.Fprintf(&b, "\t\t\t%s\n", ev.Handler)
		b.WriteString("\t\t\treturn nil\n")
		b.WriteString("\t\t}))\n")
		b.WriteString("\t}\n")
	}
	
	// Bind inputs in branches
	branchBindings := filterBranchBindings(inputBindings)
	for _, ib := range branchBindings {
		eventType := "input"
		if ib.Attr == "checked" {
			eventType = "change"
		}
		fmt.Fprintf(&b, "\tif el := document.Call(\"getElementById\", \"%s\"); !el.IsNull() && el.Get(\"_bound\").IsUndefined() {\n", ib.ID)
		fmt.Fprintf(&b, "\t\tel.Set(\"_bound\", true)\n")
		fmt.Fprintf(&b, "\t\tel.Call(\"addEventListener\", \"%s\", js.FuncOf(func(this js.Value, args []js.Value) any {\n", eventType)
		if ib.Attr == "checked" {
			fmt.Fprintf(&b, "\t\t\tset%s(this.Get(\"checked\").Bool())\n", capitalize(ib.VarName))
		} else if ib.VarType == "int" {
			fmt.Fprintf(&b, "\t\t\tif v, err := strconv.Atoi(this.Get(\"value\").String()); err == nil {\n")
			fmt.Fprintf(&b, "\t\t\t\tset%s(v)\n", capitalize(ib.VarName))
			b.WriteString("\t\t\t}\n")
		} else {
			fmt.Fprintf(&b, "\t\t\tset%s(this.Get(\"value\").String())\n", capitalize(ib.VarName))
		}
		b.WriteString("\t\t\treturn nil\n")
		b.WriteString("\t\t}))\n")
		b.WriteString("\t}\n")
	}
	b.WriteString("}\n")

	goCode = b.String()
	html = generateHTML(p, tmplAST, events, inputBindings, ifs, components)
	return
}

// Collect events with branch info
func collectEvents(template string, ifs []IfNode, eaches []EachNode) []Event {
	var events []Event
	btnID := 0

	// Find @event|modifier1|modifier2="handler"
	for i := 0; i < len(template); i++ {
		if template[i] != '@' {
			continue
		}
		
		// Skip {@html ...} directives
		if i > 0 && template[i-1] == '{' {
			continue
		}
		
		// Find the ="
		eqIdx := strings.Index(template[i:], "=\"")
		if eqIdx == -1 {
			continue
		}
		
		// Parse event and modifiers: @click|preventDefault|stopPropagation
		eventPart := template[i+1 : i+eqIdx]
		parts := strings.Split(eventPart, "|")
		eventName := parts[0]
		modifiers := parts[1:]
		
		// Find handler
		handlerStart := i + eqIdx + 2
		handlerEnd := strings.Index(template[handlerStart:], "\"")
		if handlerEnd == -1 {
			continue
		}
		handler := template[handlerStart : handlerStart+handlerEnd]
		
		inBranch := isPositionInIf(template, i, ifs)
		inEach, eachID := isPositionInEach(template, i, eaches)
		events = append(events, Event{
			ID:        fmt.Sprintf("btn%d", btnID),
			Event:     eventName,
			Handler:   handler,
			Modifiers: modifiers,
			InBranch:  inBranch,
			InEach:    inEach,
			EachID:    eachID,
		})
		btnID++
		i = handlerStart + handlerEnd
	}
	return events
}

// Collect bindings with branch info
func collectBindings(template, script string, ifs []IfNode) []Binding {
	var bindings []Binding
	inputID := 0

	for i := 0; i < len(template); i++ {
		if idx := strings.Index(template[i:], "bind:value=\""); idx != -1 {
			absPos := i + idx
			start := absPos + 12
			end := strings.Index(template[start:], "\"")
			if end != -1 {
				varName := template[start : start+end]
				inBranch := isPositionInIf(template, absPos, ifs)
				bindings = append(bindings, Binding{
					ID:       fmt.Sprintf("input%d", inputID),
					Attr:     "value",
					VarName:  varName,
					VarType:  inferType(script, varName),
					InBranch: inBranch,
				})
				inputID++
				i = start + end
			}
		}
	}

	for i := 0; i < len(template); i++ {
		if idx := strings.Index(template[i:], "bind:checked=\""); idx != -1 {
			absPos := i + idx
			start := absPos + 14
			end := strings.Index(template[start:], "\"")
			if end != -1 {
				varName := template[start : start+end]
				inBranch := isPositionInIf(template, absPos, ifs)
				bindings = append(bindings, Binding{
					ID:       fmt.Sprintf("input%d", inputID),
					Attr:     "checked",
					VarName:  varName,
					VarType:  "bool",
					InBranch: inBranch,
				})
				inputID++
				i = start + end
			}
		}
	}

	return bindings
}

// Collect class bindings with branch info
func collectClassBindings(template string, ifs []IfNode) []ClassBinding {
	var bindings []ClassBinding
	classID := 0
	needNewID := true
	i := 0
	
	for i < len(template) {
		idx := strings.Index(template[i:], "class:")
		if idx == -1 {
			break
		}
		absPos := i + idx
		eqIdx := strings.Index(template[absPos:], "={")
		if eqIdx == -1 {
			i = absPos + 6
			continue
		}
		className := template[absPos+6 : absPos+eqIdx]
		
		condStart := absPos + eqIdx + 2
		condEnd := strings.Index(template[condStart:], "}")
		if condEnd == -1 {
			i = absPos + 6
			continue
		}
		cond := template[condStart : condStart+condEnd]
		fullEnd := condStart + condEnd + 1
		
		// Assign new ID if this is first in a group
		if needNewID {
			classID++
			needNewID = false
		}
		
		inBranch := isPositionInIf(template, absPos, ifs)
		bindings = append(bindings, ClassBinding{
			ID:        fmt.Sprintf("class%d", classID-1),
			ClassName: className,
			Cond:      cond,
			InBranch:  inBranch,
		})
		
		// Check if there's another class: immediately after
		rest := strings.TrimLeft(template[fullEnd:], " \t")
		if !strings.HasPrefix(rest, "class:") {
			// End of group - next binding needs new ID
			needNewID = true
		}
		i = fullEnd
	}

	return bindings
}

func isPositionInIf(template string, pos int, ifs []IfNode) bool {
	// Find all {#if and {/if} positions
	depth := 0
	for i := 0; i < pos && i < len(template); i++ {
		if strings.HasPrefix(template[i:], "{#if ") {
			depth++
		} else if strings.HasPrefix(template[i:], "{/if}") {
			depth--
		}
	}
	return depth > 0
}

func isPositionInEach(template string, pos int, eaches []EachNode) (bool, string) {
	// Track each block depth and which each we're in
	type eachInfo struct {
		id    string
		depth int
	}
	var stack []eachInfo
	eachCount := 0
	
	for i := 0; i < pos && i < len(template); i++ {
		if strings.HasPrefix(template[i:], "{#each ") {
			stack = append(stack, eachInfo{id: fmt.Sprintf("each%d", eachCount), depth: 1})
			eachCount++
		} else if strings.HasPrefix(template[i:], "{/each}") {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		}
	}
	
	if len(stack) > 0 {
		return true, stack[len(stack)-1].id
	}
	return false, ""
}

func countEventsOutsideIf(template string, ifs []IfNode) int {
	count := 0
	for i := 0; i < len(template); i++ {
		if template[i] != '@' {
			continue
		}
		// Skip {@html ...} directives
		if i > 0 && template[i-1] == '{' {
			continue
		}
		eqIdx := strings.Index(template[i:], "=\"")
		if eqIdx == -1 {
			continue
		}
		if !isPositionInIf(template, i, ifs) {
			count++
		}
		// Skip past this event
		handlerStart := i + eqIdx + 2
		handlerEnd := strings.Index(template[handlerStart:], "\"")
		if handlerEnd != -1 {
			i = handlerStart + handlerEnd
		}
	}
	return count
}

func countBindingsOutsideIf(template string, ifs []IfNode) int {
	count := 0
	for i := 0; i < len(template); i++ {
		if idx := strings.Index(template[i:], "bind:value=\""); idx != -1 {
			if !isPositionInIf(template, i+idx, ifs) {
				count++
			}
			end := strings.Index(template[i+idx+12:], "\"")
			if end != -1 {
				i = i + idx + 12 + end
			}
		}
	}
	for i := 0; i < len(template); i++ {
		if idx := strings.Index(template[i:], "bind:checked=\""); idx != -1 {
			if !isPositionInIf(template, i+idx, ifs) {
				count++
			}
			end := strings.Index(template[i+idx+14:], "\"")
			if end != -1 {
				i = i + idx + 14 + end
			}
		}
	}
	return count
}

func countClassBindingsOutsideIf(template string, ifs []IfNode) int {
	count := 0
	i := 0
	for i < len(template) {
		idx := strings.Index(template[i:], "class:")
		if idx == -1 {
			break
		}
		absPos := i + idx
		eqIdx := strings.Index(template[absPos:], "={")
		if eqIdx == -1 {
			i = absPos + 6
			continue
		}
		condStart := absPos + eqIdx + 2
		condEnd := strings.Index(template[condStart:], "}")
		if condEnd == -1 {
			i = absPos + 6
			continue
		}
		fullEnd := condStart + condEnd + 1
		
		// Check if there's another class: immediately after
		rest := strings.TrimLeft(template[fullEnd:], " \t")
		if !strings.HasPrefix(rest, "class:") {
			// This is the last in a group - count it if outside if
			if !isPositionInIf(template, absPos, ifs) {
				count++
			}
		}
		i = fullEnd
	}
	return count
}

func filterStaticEvents(events []Event) []Event {
	var result []Event
	for _, ev := range events {
		if !ev.InBranch && !ev.InEach {
			result = append(result, ev)
		}
	}
	return result
}

func filterBranchEvents(events []Event) []Event {
	var result []Event
	for _, ev := range events {
		if ev.InBranch && !ev.InEach {
			result = append(result, ev)
		}
	}
	return result
}

func filterEachEvents(events []Event, eachID string) []Event {
	var result []Event
	for _, ev := range events {
		if ev.InEach && ev.EachID == eachID {
			result = append(result, ev)
		}
	}
	return result
}

func filterStaticBindings(bindings []Binding) []Binding {
	var result []Binding
	for _, b := range bindings {
		if !b.InBranch {
			result = append(result, b)
		}
	}
	return result
}

func filterBranchBindings(bindings []Binding) []Binding {
	var result []Binding
	for _, b := range bindings {
		if b.InBranch {
			result = append(result, b)
		}
	}
	return result
}

func generateBranchHTML(nodes []Node, ast *TemplateAST) string {
	var b strings.Builder
	ast.generateHTMLNodes(&b, nodes)
	return b.String()
}

func processTemplateAttrs(html string, btnID, inputID, classID int) (string, int, int, int) {
	// Replace @event|modifiers="handler" with id="btnN"
	i := 0
	for i < len(html) {
		idx := strings.Index(html[i:], "@")
		if idx == -1 {
			break
		}
		absIdx := i + idx
		// Skip {@html ...} directives
		if absIdx > 0 && html[absIdx-1] == '{' {
			i = absIdx + 1
			continue
		}
		// Find ="
		eqIdx := strings.Index(html[absIdx:], "=\"")
		if eqIdx == -1 {
			break
		}
		// Find closing quote
		handlerStart := absIdx + eqIdx + 2
		handlerEnd := strings.Index(html[handlerStart:], "\"")
		if handlerEnd == -1 {
			break
		}
		// Replace entire @event|mods="handler" with id="btnN"
		html = html[:absIdx] + fmt.Sprintf("id=\"btn%d\"", btnID) + html[handlerStart+handlerEnd+1:]
		btnID++
		i = absIdx + 1
	}

	// Replace bind:value="var" with id="inputN"
	for strings.Contains(html, "bind:value=\"") {
		idx := strings.Index(html, "bind:value=\"")
		end := strings.Index(html[idx+12:], "\"")
		if end != -1 {
			html = strings.Replace(html, html[idx:idx+12+end+1], fmt.Sprintf("id=\"input%d\"", inputID), 1)
			inputID++
		} else {
			break
		}
	}

	// Replace bind:checked="var" with id="inputN"
	for strings.Contains(html, "bind:checked=\"") {
		idx := strings.Index(html, "bind:checked=\"")
		end := strings.Index(html[idx+14:], "\"")
		if end != -1 {
			html = strings.Replace(html, html[idx:idx+14+end+1], fmt.Sprintf("id=\"input%d\"", inputID), 1)
			inputID++
		} else {
			break
		}
	}

	// Replace class:name={cond} with id="classN" (only first in group)
	for {
		idx := strings.Index(html, "class:")
		if idx == -1 {
			break
		}
		eqIdx := strings.Index(html[idx:], "={")
		if eqIdx == -1 {
			break
		}
		condStart := idx + eqIdx + 2
		condEnd := strings.Index(html[condStart:], "}")
		if condEnd == -1 {
			break
		}
		fullEnd := condStart + condEnd + 1
		
		// Check if there's another class: immediately after (same element)
		rest := strings.TrimLeft(html[fullEnd:], " \t")
		if strings.HasPrefix(rest, "class:") {
			// More class bindings follow - just remove this one
			html = html[:idx] + html[fullEnd:]
		} else {
			// Last (or only) class binding - replace with id
			html = html[:idx] + fmt.Sprintf("id=\"class%d\"", classID) + html[fullEnd:]
			classID++
		}
	}

	return html, btnID, inputID, classID
}

func generateHTML(p Parsed, tmplAST *TemplateAST, events []Event, inputBindings []Binding, ifs []IfNode, components []ComponentNode) string {
	var b strings.Builder

	b.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	b.WriteString("  <script src=\"wasm_exec.js\"></script>\n")
	b.WriteString("  <script>\n")
	b.WriteString("    const go = new Go();\n")
	b.WriteString("    WebAssembly.instantiateStreaming(fetch(\"app.wasm\"), go.importObject)\n")
	b.WriteString("      .then(r => go.run(r.instance));\n")
	b.WriteString("  </script>\n")
	
	// Include component styles
	for name, comp := range p.Components {
		if comp.Style != "" {
			fmt.Fprintf(&b, "  <!-- %s styles -->\n", name)
			b.WriteString("  <style>\n")
			b.WriteString(comp.Style)
			b.WriteString("\n  </style>\n")
		}
	}
	
	if p.Style != "" {
		b.WriteString("  <style>\n")
		b.WriteString(p.Style)
		b.WriteString("\n  </style>\n")
	}
	b.WriteString("</head>\n<body>\n")

	html := tmplAST.GenerateHTML()

	// Replace component placeholders with actual component HTML
	for _, comp := range components {
		compDef := p.Components[comp.Name]
		if compDef == nil {
			continue
		}
		
		// Generate component's HTML with prefixed IDs
		compHTML := generateComponentHTML(comp.ID, compDef, p.Components, comp.Props)
		
		// Replace placeholder
		placeholder := fmt.Sprintf(`<span id="%s"></span>`, comp.ID)
		html = strings.Replace(html, placeholder, compHTML, 1)
	}

	// Hydrate expression placeholders with initial values
	exprs := tmplAST.CollectExprs()
	for _, expr := range exprs {
		initVal := getInitializer(p.Script, expr.Expr)
		if initVal != "" && initVal != "0" && initVal != `""` {
			// Clean up string quotes for display
			displayVal := initVal
			if strings.HasPrefix(displayVal, `"`) && strings.HasSuffix(displayVal, `"`) {
				displayVal = displayVal[1 : len(displayVal)-1]
			}
			empty := fmt.Sprintf(`<span id="%s"></span>`, expr.ID)
			filled := fmt.Sprintf(`<span id="%s">%s</span>`, expr.ID, displayVal)
			html = strings.Replace(html, empty, filled, 1)
		}
	}

	// Replace @event|modifiers="handler" with id="btnN"
	btnID := 0
	i := 0
	for i < len(html) {
		idx := strings.Index(html[i:], "@")
		if idx == -1 {
			break
		}
		absIdx := i + idx
		// Skip {@html ...} directives
		if absIdx > 0 && html[absIdx-1] == '{' {
			i = absIdx + 1
			continue
		}
		eqIdx := strings.Index(html[absIdx:], "=\"")
		if eqIdx == -1 {
			break
		}
		handlerStart := absIdx + eqIdx + 2
		handlerEnd := strings.Index(html[handlerStart:], "\"")
		if handlerEnd == -1 {
			break
		}
		html = html[:absIdx] + fmt.Sprintf("id=\"btn%d\"", btnID) + html[handlerStart+handlerEnd+1:]
		btnID++
		i = absIdx + 1
	}

	inputID := 0
	for strings.Contains(html, "bind:value=\"") {
		idx := strings.Index(html, "bind:value=\"")
		end := strings.Index(html[idx+12:], "\"")
		if end != -1 {
			html = strings.Replace(html, html[idx:idx+12+end+1], fmt.Sprintf("id=\"input%d\"", inputID), 1)
			inputID++
		} else {
			break
		}
	}

	for strings.Contains(html, "bind:checked=\"") {
		idx := strings.Index(html, "bind:checked=\"")
		end := strings.Index(html[idx+14:], "\"")
		if end != -1 {
			html = strings.Replace(html, html[idx:idx+14+end+1], fmt.Sprintf("id=\"input%d\"", inputID), 1)
			inputID++
		} else {
			break
		}
	}

	// Replace class:name={cond} with id="classN" (only first in group)
	classID := 0
	for {
		idx := strings.Index(html, "class:")
		if idx == -1 {
			break
		}
		eqIdx := strings.Index(html[idx:], "={")
		if eqIdx == -1 {
			break
		}
		condStart := idx + eqIdx + 2
		condEnd := strings.Index(html[condStart:], "}")
		if condEnd == -1 {
			break
		}
		fullEnd := condStart + condEnd + 1
		
		// Check if there's another class: immediately after (same element)
		rest := strings.TrimLeft(html[fullEnd:], " \t")
		if strings.HasPrefix(rest, "class:") {
			// More class bindings follow - just remove this one, keep going
			html = html[:idx] + html[fullEnd:]
		} else {
			// Last (or only) class binding - replace with id
			html = html[:idx] + fmt.Sprintf("id=\"class%d\"", classID) + html[fullEnd:]
			classID++
		}
	}

	b.WriteString(html)
	b.WriteString("\n</body>\n</html>\n")
	return b.String()
}

func transformFunction(script, funcName string, modifies []string) string {
	funcStart := strings.Index(script, "func "+funcName)
	if funcStart == -1 {
		return ""
	}

	braceStart := strings.Index(script[funcStart:], "{")
	if braceStart == -1 {
		return ""
	}
	braceStart += funcStart

	depth := 1
	braceEnd := braceStart + 1
	for braceEnd < len(script) && depth > 0 {
		switch script[braceEnd] {
		case '{':
			depth++
		case '}':
			depth--
		}
		braceEnd++
	}

	funcSrc := script[funcStart:braceEnd]

	for _, varName := range modifies {
		funcSrc = strings.ReplaceAll(funcSrc, varName+"++", fmt.Sprintf("set%s(%s + 1)", capitalize(varName), varName))
		funcSrc = strings.ReplaceAll(funcSrc, varName+"--", fmt.Sprintf("set%s(%s - 1)", capitalize(varName), varName))

		lines := strings.Split(funcSrc, "\n")
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, varName+" = ") {
				expr := strings.TrimPrefix(trimmed, varName+" = ")
				indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
				lines[i] = indent + fmt.Sprintf("set%s(%s)", capitalize(varName), expr)
			}
		}
		funcSrc = strings.Join(lines, "\n")
	}

	return funcSrc
}

func getInitializer(script, name string) string {
	for _, line := range strings.Split(script, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, name+" := ") {
			return strings.TrimPrefix(trimmed, name+" := ")
		}
		if strings.HasPrefix(trimmed, name+" = ") {
			return strings.TrimPrefix(trimmed, name+" = ")
		}
	}
	return "0"
}

func getExpression(script, name string) string {
	return getInitializer(script, name)
}

func inferType(script, name string) string {
	init := getInitializer(script, name)
	if init == "" || init == "0" {
		return "int"
	}
	if strings.HasPrefix(init, "[]string{") {
		return "[]string"
	}
	if strings.HasPrefix(init, "[]int{") {
		return "[]int"
	}
	if strings.HasPrefix(init, "[]float64{") {
		return "[]float64"
	}
	if strings.HasPrefix(init, "[]bool{") {
		return "[]bool"
	}
	if strings.HasPrefix(init, "\"") {
		return "string"
	}
	if init == "true" || init == "false" {
		return "bool"
	}
	if strings.Contains(init, "&&") || strings.Contains(init, "||") ||
		strings.Contains(init, "==") || strings.Contains(init, "!=") ||
		strings.Contains(init, "<=") || strings.Contains(init, ">=") ||
		strings.Contains(init, " < ") || strings.Contains(init, " > ") {
		return "bool"
	}
	if strings.Contains(init, "\" + ") || strings.Contains(init, " + \"") {
		return "string"
	}
	if strings.Contains(init, ".") {
		return "float64"
	}
	return "int"
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func generateEventModifiers(b *strings.Builder, modifiers []string, indent string) {
	for _, mod := range modifiers {
		switch mod {
		case "preventDefault":
			fmt.Fprintf(b, "%sargs[0].Call(\"preventDefault\")\n", indent)
		case "stopPropagation":
			fmt.Fprintf(b, "%sargs[0].Call(\"stopPropagation\")\n", indent)
		}
	}
}

func isExprInIf(exprID string, ifs []IfNode) bool {
	for _, ifn := range ifs {
		if exprInNodes(exprID, ifn.Then) || exprInNodes(exprID, ifn.Else) {
			return true
		}
	}
	return false
}

func exprInNodes(exprID string, nodes []Node) bool {
	for _, n := range nodes {
		switch node := n.(type) {
		case ExprNode:
			if node.ID == exprID {
				return true
			}
		case IfNode:
			if exprInNodes(exprID, node.Then) || exprInNodes(exprID, node.Else) {
				return true
			}
		case EachNode:
			if exprInNodes(exprID, node.Body) {
				return true
			}
		}
	}
	return false
}

func isExprInEach(exprID string, eaches []EachNode) bool {
	for _, each := range eaches {
		if exprInNodes(exprID, each.Body) {
			return true
		}
	}
	return false
}

func collectExprsInEach(each EachNode) []ExprNode {
	var exprs []ExprNode
	collectExprsFromNodes(each.Body, &exprs)
	return exprs
}

func collectExprsFromNodes(nodes []Node, exprs *[]ExprNode) {
	for _, n := range nodes {
		switch node := n.(type) {
		case ExprNode:
			*exprs = append(*exprs, node)
		case IfNode:
			collectExprsFromNodes(node.Then, exprs)
			collectExprsFromNodes(node.Else, exprs)
		case EachNode:
			collectExprsFromNodes(node.Body, exprs)
		}
	}
}

func generateEachBodyHTML(nodes []Node, ast *TemplateAST, eachID string) string {
	var b strings.Builder
	generateEachHTMLNodes(&b, nodes, eachID)
	return b.String()
}

// replaceVar replaces a variable name with a new name, respecting word boundaries
func replaceVar(s, oldVar, newVar string) string {
	// Simple approach: replace (var), var), (var, ,var) patterns
	result := s
	result = strings.ReplaceAll(result, "("+oldVar+")", "("+newVar+")")
	result = strings.ReplaceAll(result, "("+oldVar+",", "("+newVar+",")
	result = strings.ReplaceAll(result, ","+oldVar+")", ","+newVar+")")
	result = strings.ReplaceAll(result, ", "+oldVar+")", ", "+newVar+")")
	return result
}

// processEachTemplateAttrs replaces @click etc with indexed IDs for each templates
func processEachTemplateAttrs(html string) string {
	btnID := 0
	
	// Replace @event|modifiers="handler" with id="btnN_${_idx}"
	i := 0
	for i < len(html) {
		idx := strings.Index(html[i:], "@")
		if idx == -1 {
			break
		}
		absIdx := i + idx
		// Skip {@html ...}
		if absIdx > 0 && html[absIdx-1] == '{' {
			i = absIdx + 1
			continue
		}
		eqIdx := strings.Index(html[absIdx:], "=\"")
		if eqIdx == -1 {
			break
		}
		handlerStart := absIdx + eqIdx + 2
		handlerEnd := strings.Index(html[handlerStart:], "\"")
		if handlerEnd == -1 {
			break
		}
		// Replace with indexed ID
		html = html[:absIdx] + fmt.Sprintf("id=\"btn%d_${_idx}\"", btnID) + html[handlerStart+handlerEnd+1:]
		btnID++
		i = absIdx + 1
	}
	
	return html
}

func generateEachHTMLNodes(b *strings.Builder, nodes []Node, eachID string) {
	for _, n := range nodes {
		switch node := n.(type) {
		case TextNode:
			b.WriteString(node.Text)
		case ExprNode:
			// Use index placeholder for ID so each item has unique IDs
			fmt.Fprintf(b, `<span id="%s_${_idx}"></span>`, node.ID)
		case IfNode:
			fmt.Fprintf(b, `<span id="%s_anchor" style="display:none"></span>`, node.CondID)
		case EachNode:
			fmt.Fprintf(b, `<span id="%s_anchor" style="display:none"></span>`, node.ID)
		}
	}
}

// generateComponentCode generates struct types, instances, and methods for components
func generateComponentCode(b *strings.Builder, p Parsed, components []ComponentNode, compAnalyses map[string]*Analysis) {
	// Group components by type and collect which props have bindings
	compTypes := make(map[string]bool)
	compBindings := make(map[string]map[string]bool) // typeName -> propName -> true
	for _, comp := range components {
		compTypes[comp.Name] = true
		if compBindings[comp.Name] == nil {
			compBindings[comp.Name] = make(map[string]bool)
		}
		for propName := range comp.Bindings {
			compBindings[comp.Name][capitalize(propName)] = true
		}
	}

	// Generate struct type for each component type
	for typeName := range compTypes {
		compDef := p.Components[typeName]
		analysis := compAnalyses[typeName]
		if compDef == nil || analysis == nil {
			continue
		}

		fmt.Fprintf(b, "// %s component\n", typeName)
		fmt.Fprintf(b, "type %s struct {\n", typeName)
		b.WriteString("\t_id string\n")
		
		// Fields from analysis
		for _, name := range analysis.Order {
			typ := inferType(compDef.Script, name)
			fmt.Fprintf(b, "\t%s %s\n", name, typ)
		}
		
		// Binding callback fields
		for propName := range compBindings[typeName] {
			typ := inferType(compDef.Script, propName)
			fmt.Fprintf(b, "\t_on%sChange func(%s)\n", propName, typ)
		}
		b.WriteString("}\n\n")

		// Parse component template
		compAST := parseTemplate(compDef.Template)
		compExprs := compAST.CollectExprs()
		compEvents := collectEvents(compDef.Template, nil, nil)

		// Setters - include bound props
		for _, name := range analysis.Order {
			v := analysis.Vars[name]
			if len(v.DependsOn) == 0 {
				typ := inferType(compDef.Script, name)
				fmt.Fprintf(b, "func (c *%s) set%s(v %s) {\n", typeName, capitalize(name), typ)
				fmt.Fprintf(b, "\tc.%s = v\n", name)
				// If this prop has bindings, call the callback
				if compBindings[typeName][capitalize(name)] {
					fmt.Fprintf(b, "\tif c._on%sChange != nil {\n", capitalize(name))
					fmt.Fprintf(b, "\t\tc._on%sChange(v)\n", capitalize(name))
					b.WriteString("\t}\n")
				}
				b.WriteString("\tc.propagate()\n")
				b.WriteString("}\n\n")
			}
		}

		// Propagate
		fmt.Fprintf(b, "func (c *%s) propagate() {\n", typeName)
		
		// Derived values
		hasDerived := false
		for _, name := range analysis.Order {
			if len(analysis.Vars[name].DependsOn) > 0 {
				if !hasDerived {
					b.WriteString("\t// Recompute derived values\n")
					hasDerived = true
				}
				fmt.Fprintf(b, "\tc.%s = %s\n", name, transformCompExpr(getExpression(compDef.Script, name)))
			}
		}
		
		b.WriteString("\tc.updateDOM()\n")
		b.WriteString("}\n\n")

		// updateDOM
		fmt.Fprintf(b, "func (c *%s) updateDOM() {\n", typeName)
		for _, expr := range compExprs {
			typ := inferType(compDef.Script, expr.Expr)
			fmt.Fprintf(b, "\tif el := document.Call(\"getElementById\", c._id+\"_%s\"); !el.IsNull() {\n", expr.ID)
			if typ == "int" {
				fmt.Fprintf(b, "\t\tel.Set(\"textContent\", strconv.Itoa(c.%s))\n", expr.Expr)
			} else if typ == "bool" {
				fmt.Fprintf(b, "\t\tel.Set(\"textContent\", strconv.FormatBool(c.%s))\n", expr.Expr)
			} else {
				fmt.Fprintf(b, "\t\tel.Set(\"textContent\", c.%s)\n", expr.Expr)
			}
			b.WriteString("\t}\n")
		}
		
		// Sync component input values
		compInputBindings := collectBindings(compDef.Template, compDef.Script, nil)
		for i, ib := range compInputBindings {
			if ib.Attr == "checked" {
				fmt.Fprintf(b, "\tif el := document.Call(\"getElementById\", c._id+\"_input%d\"); !el.IsNull() && el.Get(\"checked\").Bool() != c.%s {\n", i, ib.VarName)
				fmt.Fprintf(b, "\t\tel.Set(\"checked\", c.%s)\n", ib.VarName)
				b.WriteString("\t}\n")
			} else {
				if ib.VarType == "int" {
					fmt.Fprintf(b, "\tif el := document.Call(\"getElementById\", c._id+\"_input%d\"); !el.IsNull() && el.Get(\"value\").String() != strconv.Itoa(c.%s) {\n", i, ib.VarName)
					fmt.Fprintf(b, "\t\tel.Set(\"value\", strconv.Itoa(c.%s))\n", ib.VarName)
				} else {
					fmt.Fprintf(b, "\tif el := document.Call(\"getElementById\", c._id+\"_input%d\"); !el.IsNull() && el.Get(\"value\").String() != c.%s {\n", i, ib.VarName)
					fmt.Fprintf(b, "\t\tel.Set(\"value\", c.%s)\n", ib.VarName)
				}
				b.WriteString("\t}\n")
			}
		}
		b.WriteString("}\n\n")

		// User functions as methods
		for name, fn := range analysis.Funcs {
			transformed := transformComponentFunction(compDef.Script, typeName, name, fn.Modifies)
			b.WriteString(transformed)
			b.WriteString("\n\n")
		}

		// mount function
		fmt.Fprintf(b, "func (c *%s) mount() {\n", typeName)
		for i, ev := range compEvents {
			fmt.Fprintf(b, "\tif el := document.Call(\"getElementById\", c._id+\"_btn%d\"); !el.IsNull() {\n", i)
			fmt.Fprintf(b, "\t\tel.Call(\"addEventListener\", \"%s\", js.FuncOf(func(this js.Value, args []js.Value) any {\n", ev.Event)
			// Convert handler() to c.handler()
			fmt.Fprintf(b, "\t\t\tc.%s\n", ev.Handler)
			b.WriteString("\t\t\treturn nil\n")
			b.WriteString("\t\t}))\n")
			b.WriteString("\t}\n")
		}
		
		// Bind component inputs
		for i, ib := range compInputBindings {
			eventType := "input"
			if ib.Attr == "checked" {
				eventType = "change"
			}
			fmt.Fprintf(b, "\tif el := document.Call(\"getElementById\", c._id+\"_input%d\"); !el.IsNull() {\n", i)
			fmt.Fprintf(b, "\t\tel.Call(\"addEventListener\", \"%s\", js.FuncOf(func(this js.Value, args []js.Value) any {\n", eventType)
			if ib.Attr == "checked" {
				fmt.Fprintf(b, "\t\t\tc.set%s(this.Get(\"checked\").Bool())\n", capitalize(ib.VarName))
			} else if ib.VarType == "int" {
				fmt.Fprintf(b, "\t\t\tif v, err := strconv.Atoi(this.Get(\"value\").String()); err == nil {\n")
				fmt.Fprintf(b, "\t\t\t\tc.set%s(v)\n", capitalize(ib.VarName))
				b.WriteString("\t\t\t}\n")
			} else {
				fmt.Fprintf(b, "\t\t\tc.set%s(this.Get(\"value\").String())\n", capitalize(ib.VarName))
			}
			b.WriteString("\t\t\treturn nil\n")
			b.WriteString("\t\t}))\n")
			b.WriteString("\t}\n")
		}
		
		b.WriteString("\tc.propagate()\n")
		// Call onMount if it exists
		if _, hasOnMount := analysis.Funcs["onMount"]; hasOnMount {
			b.WriteString("\tc.onMount()\n")
		}
		b.WriteString("}\n\n")
		
		// destroy function (for cleanup)
		if _, hasOnDestroy := analysis.Funcs["onDestroy"]; hasOnDestroy {
			fmt.Fprintf(b, "func (c *%s) destroy() {\n", typeName)
			b.WriteString("\tc.onDestroy()\n")
			b.WriteString("}\n\n")
		}
	}

	// Generate instances
	b.WriteString("// Component instances\n")
	b.WriteString("var (\n")
	for _, comp := range components {
		fmt.Fprintf(b, "\t%s = &%s{_id: \"%s\"", comp.ID, comp.Name, comp.ID)
		for propName, propVal := range comp.Props {
			// Check if value is expression {expr} or literal "string"
			if strings.HasPrefix(propVal, "{") && strings.HasSuffix(propVal, "}") {
				// Expression - will be set at runtime
				continue
			}
			fmt.Fprintf(b, ", %s: \"%s\"", propName, propVal)
		}
		b.WriteString("}\n")
	}
	b.WriteString(")\n\n")
}

func transformComponentFunction(script, typeName, funcName string, modifies []string) string {
	funcStart := strings.Index(script, "func "+funcName)
	if funcStart == -1 {
		return ""
	}

	braceStart := strings.Index(script[funcStart:], "{")
	if braceStart == -1 {
		return ""
	}
	braceStart += funcStart

	depth := 1
	braceEnd := braceStart + 1
	for braceEnd < len(script) && depth > 0 {
		switch script[braceEnd] {
		case '{':
			depth++
		case '}':
			depth--
		}
		braceEnd++
	}

	funcSrc := script[funcStart:braceEnd]
	
	// Convert func name() to func (c *Type) name()
	funcSrc = strings.Replace(funcSrc, "func "+funcName, fmt.Sprintf("func (c *%s) %s", typeName, funcName), 1)

	// Transform variable modifications to method calls
	for _, varName := range modifies {
		funcSrc = strings.ReplaceAll(funcSrc, varName+"++", fmt.Sprintf("c.set%s(c.%s + 1)", capitalize(varName), varName))
		funcSrc = strings.ReplaceAll(funcSrc, varName+"--", fmt.Sprintf("c.set%s(c.%s - 1)", capitalize(varName), varName))

		lines := strings.Split(funcSrc, "\n")
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, varName+" = ") {
				expr := strings.TrimPrefix(trimmed, varName+" = ")
				indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
				lines[i] = indent + fmt.Sprintf("c.set%s(%s)", capitalize(varName), expr)
			}
		}
		funcSrc = strings.Join(lines, "\n")
	}

	return funcSrc
}

func transformCompExpr(expr string) string {
	// Add c. prefix to identifiers - simplified version
	return "c." + expr
}

func isUppercase(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[0] >= 'A' && s[0] <= 'Z'
}

// flattenComponents collects all components recursively
func flattenComponents(comps map[string]*Parsed, result map[string]*Parsed) {
	for name, comp := range comps {
		if _, exists := result[name]; !exists {
			result[name] = comp
			// Recurse into nested components
			flattenComponents(comp.Components, result)
		}
	}
}

// collectAllComponentUsages collects component usages from main template and nested components
func collectAllComponentUsages(topLevelComps []ComponentNode, allCompDefs map[string]*Parsed) []ComponentNode {
	var result []ComponentNode
	
	// Add top-level components
	result = append(result, topLevelComps...)
	
	// For each component, find nested component usages
	for _, comp := range topLevelComps {
		compDef := allCompDefs[comp.Name]
		if compDef == nil {
			continue
		}
		
		// Build component names for parsing
		compNames := make(map[string]bool)
		for name := range allCompDefs {
			compNames[name] = true
		}
		
		// Parse this component's template
		ast := parseTemplateWithComponents(compDef.Template, compNames)
		nestedComps := ast.CollectComponents()
		
		// Add nested components with prefixed IDs
		for _, nested := range nestedComps {
			prefixedComp := ComponentNode{
				Name:     nested.Name,
				ID:       comp.ID + "_" + nested.ID,
				Props:    nested.Props,
				Children: nested.Children,
			}
			
			// Recursively collect from this nested component
			subNested := collectAllComponentUsages([]ComponentNode{prefixedComp}, allCompDefs)
			result = append(result, subNested...)
		}
	}
	
	return result
}

// generateComponentHTML generates HTML for a component with prefixed IDs
func generateComponentHTML(compID string, compDef *Parsed, allComponents map[string]*Parsed, props map[string]string) string {
	// Build component names for this component's template
	compNames := make(map[string]bool)
	for name := range allComponents {
		compNames[name] = true
	}
	
	// Parse the component template with nested component awareness
	ast := parseTemplateWithComponents(compDef.Template, compNames)
	
	var b strings.Builder
	generateComponentHTMLNodes(&b, ast.Nodes, compID, allComponents)
	
	html := b.String()
	
	// Replace @click etc with prefixed btn IDs
	btnID := 0
	for strings.Contains(html, "@") {
		idx := strings.Index(html, "@")
		if idx == -1 {
			break
		}
		// Skip if this is {@html}
		if idx > 0 && html[idx-1] == '{' {
			// Skip past this @
			html = html[:idx] + "___AT___" + html[idx+1:]
			continue
		}
		eqIdx := strings.Index(html[idx:], "=\"")
		if eqIdx == -1 {
			break
		}
		handlerStart := idx + eqIdx + 2
		handlerEnd := strings.Index(html[handlerStart:], "\"")
		if handlerEnd == -1 {
			break
		}
		html = html[:idx] + fmt.Sprintf("id=\"%s_btn%d\"", compID, btnID) + html[handlerStart+handlerEnd+1:]
		btnID++
	}
	// Restore {@html}
	html = strings.ReplaceAll(html, "___AT___", "@")
	
	// Replace bind:value with prefixed input IDs  
	inputID := 0
	for strings.Contains(html, "bind:value=\"") {
		idx := strings.Index(html, "bind:value=\"")
		end := strings.Index(html[idx+12:], "\"")
		if end != -1 {
			html = html[:idx] + fmt.Sprintf("id=\"%s_input%d\"", compID, inputID) + html[idx+12+end+1:]
			inputID++
		} else {
			break
		}
	}
	
	// Hydrate component expressions with initial values
	exprs := ast.CollectExprs()
	for _, expr := range exprs {
		// Check if expression is a prop (passed in) or internal state
		var displayVal string
		if propVal, ok := props[expr.Expr]; ok {
			// It's a prop - use the passed value
			if strings.HasPrefix(propVal, "{") && strings.HasSuffix(propVal, "}") {
				// Dynamic prop - can't hydrate at compile time
				continue
			}
			displayVal = propVal
		} else {
			// Internal state - use initializer from component script
			initVal := getInitializer(compDef.Script, expr.Expr)
			if initVal == "" || initVal == "0" || initVal == `""` {
				continue
			}
			displayVal = initVal
			// Clean up string quotes
			if strings.HasPrefix(displayVal, `"`) && strings.HasSuffix(displayVal, `"`) {
				displayVal = displayVal[1 : len(displayVal)-1]
			}
		}
		
		if displayVal != "" {
			empty := fmt.Sprintf(`<span id="%s_%s"></span>`, compID, expr.ID)
			filled := fmt.Sprintf(`<span id="%s_%s">%s</span>`, compID, expr.ID, displayVal)
			html = strings.Replace(html, empty, filled, 1)
		}
	}
	
	return html
}

func generateComponentHTMLNodes(b *strings.Builder, nodes []Node, compID string, allComponents map[string]*Parsed) {
	for _, n := range nodes {
		switch node := n.(type) {
		case TextNode:
			b.WriteString(node.Text)
		case ExprNode:
			fmt.Fprintf(b, `<span id="%s_%s"></span>`, compID, node.ID)
		case HtmlNode:
			fmt.Fprintf(b, `<span id="%s_%s"></span>`, compID, node.ID)
		case IfNode:
			fmt.Fprintf(b, `<span id="%s_%s_anchor" style="display:none"></span>`, compID, node.CondID)
		case EachNode:
			fmt.Fprintf(b, `<span id="%s_%s_anchor" style="display:none"></span>`, compID, node.ID)
		case ComponentNode:
			// Nested component - generate with combined ID
			nestedCompDef := allComponents[node.Name]
			if nestedCompDef != nil {
				nestedID := compID + "_" + node.ID
				nestedHTML := generateComponentHTML(nestedID, nestedCompDef, allComponents, node.Props)
				b.WriteString(nestedHTML)
			}
		}
	}
}