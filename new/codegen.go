package main

import (
	"fmt"
	"strings"
)

type Event struct {
	ID        string
	Event     string
	Handler   string
	Modifiers []string
	InBranch  bool
	InEach    bool
	EachID    string
}

type Binding struct {
	ID       string
	Attr     string
	VarName  string
	VarType  string
	InBranch bool
}

type ClassBinding struct {
	ID        string
	ClassName string
	Cond      string
	InBranch  bool
	VarDeps   []string
}

func generate(p Parsed, a *Analysis) (goCode string, html string) {
	var b strings.Builder

	// Build reactive vars set for template parsing
	reactiveVars := make(map[string]bool)
	for varName := range p.ReactiveVars {
		reactiveVars[varName] = true
	}
	// Also include analyzed vars
	for varName := range a.Vars {
		reactiveVars[varName] = true
	}

	// Build component names map
	compNames := make(map[string]bool)
	for name := range p.Components {
		compNames[name] = true
	}

	// Parse template into AST with component and reactive var awareness
	tmplAST := parseTemplateWithReactiveVars(p.Template, compNames, reactiveVars)
	exprs := tmplAST.CollectExprs()
	htmls := tmplAST.CollectHtmls()
	ifs := tmplAST.CollectIfs()
	eaches := tmplAST.CollectEaches()

	// Collect events and bindings
	events := collectEvents(p.Template, ifs, eaches)
	inputBindings := collectBindings(p.Template, p.Script, ifs, p.ReactiveVars)
	classBindings := collectClassBindings(p.Template, ifs, reactiveVars)

	// Process branch templates
	branchTemplates := make(map[string]string)
	btnID := countEventsOutsideIf(p.Template, ifs)
	inputID := countBindingsOutsideIf(p.Template, ifs)
	classID := countClassBindingsOutsideIf(p.Template, ifs)

	for _, ifn := range ifs {
		for i, branch := range ifn.Branches {
			branchHTML := generateBranchHTML(branch.Body, tmplAST)
			branchHTML, btnID, inputID, classID = processTemplateAttrs(branchHTML, btnID, inputID, classID)
			branchTemplates[fmt.Sprintf("%s_branch%d", ifn.CondID, i)] = branchHTML
		}
	}

	// Process each templates
	eachTemplates := make(map[string]string)
	for _, each := range eaches {
		bodyHTML := generateEachBodyHTML(each.Body, tmplAST, each.ID)
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

	b.WriteString("var document = js.Global().Get(\"document\")\n\n")
	b.WriteString("var _ = strconv.Atoi\n\n")

	// State vars
	if len(p.ReactiveVars) > 0 {
		b.WriteString("// State\n")
		b.WriteString("var (\n")
		for varName, varType := range p.ReactiveVars {
			fmt.Fprintf(&b, "\t%s %s\n", varName, varType)
		}
		b.WriteString(")\n\n")
	}

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

	// Generate setters with fine-grained updates
	b.WriteString("// Setters with fine-grained updates\n")
	for varName, varType := range p.ReactiveVars {
		fmt.Fprintf(&b, "func set%s(v %s) {\n", capitalize(varName), varType)
		fmt.Fprintf(&b, "\t%s = v\n", varName)

		// Get all vars affected by this change (including derived)
		affected := a.GetTransitiveDeps(varName)

		// Recompute derived values that depend on this
		for _, affectedVar := range affected {
			if affectedVar == varName {
				continue
			}
			if v, ok := a.Vars[affectedVar]; ok && len(v.DependsOn) > 0 {
				fmt.Fprintf(&b, "\t%s = %s\n", affectedVar, getExpression(p.Script, affectedVar))
			}
		}

		// Update only DOM elements that depend on affected vars
		fmt.Fprintf(&b, "\tupdate%s()\n", capitalize(varName))
		b.WriteString("}\n\n")
	}

	// Generate per-var update functions
	b.WriteString("// Per-variable DOM update functions\n")
	for varName := range p.ReactiveVars {
		affected := a.GetTransitiveDeps(varName)

		fmt.Fprintf(&b, "func update%s() {\n", capitalize(varName))

		// Update expressions that depend on this var or any derived var
		for _, expr := range exprs {
			if isExprInEach(expr.ID, eaches) {
				continue
			}
			if exprDependsOnAny(expr, affected) {
				typ := inferTypeFromReactive(p.ReactiveVars, expr.Expr)
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
		}

		// Update html nodes
		for _, html := range htmls {
			if exprDependsOnAnyStr(html.VarDeps, affected) {
				fmt.Fprintf(&b, "\tif !%sEl.IsUndefined() && !%sEl.IsNull() {\n", html.ID, html.ID)
				fmt.Fprintf(&b, "\t\t%sEl.Set(\"innerHTML\", %s)\n", html.ID, html.Expr)
				b.WriteString("\t}\n")
			}
		}

		// Update if blocks
		for _, ifn := range ifs {
			if ifDependsOnAny(ifn, affected) {
				generateIfUpdate(&b, ifn)
			}
		}

		// Update each blocks
		for _, each := range eaches {
			if exprDependsOnAnyStr(each.VarDeps, affected) {
				generateEachUpdate(&b, each, events, exprs)
			}
		}

		// Update input bindings
		for _, ib := range inputBindings {
			if contains(affected, ib.VarName) {
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
		}

		// Update class bindings
		for _, cb := range classBindings {
			if exprDependsOnAnyStr(cb.VarDeps, affected) {
				fmt.Fprintf(&b, "\tif el := document.Call(\"getElementById\", \"%s\"); !el.IsNull() {\n", cb.ID)
				fmt.Fprintf(&b, "\t\tif %s {\n", cb.Cond)
				fmt.Fprintf(&b, "\t\t\tel.Get(\"classList\").Call(\"add\", \"%s\")\n", cb.ClassName)
				b.WriteString("\t\t} else {\n")
				fmt.Fprintf(&b, "\t\t\tel.Get(\"classList\").Call(\"remove\", \"%s\")\n", cb.ClassName)
				b.WriteString("\t\t}\n")
				b.WriteString("\t}\n")
			}
		}

		b.WriteString("}\n\n")
	}

	// Main function
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

	// Bind static events
	staticEvents := filterStaticEvents(events)
	if len(staticEvents) > 0 {
		b.WriteString("\n\t// Bind static events\n")
		for _, ev := range staticEvents {
			fmt.Fprintf(&b, "\tdocument.Call(\"getElementById\", \"%s\").Call(\"addEventListener\", \"%s\", js.FuncOf(func(this js.Value, args []js.Value) any {\n", ev.ID, ev.Event)
			generateEventModifiers(&b, ev.Modifiers, "\t\t")
			fmt.Fprintf(&b, "\t\t%s\n", transformEventHandler(ev.Handler))
			b.WriteString("\t\treturn nil\n")
			b.WriteString("\t}))\n")
		}
	}

	// Bind static inputs
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

	// Initial update for all reactive vars
	b.WriteString("\n\t// Initial DOM update\n")
	for varName := range p.ReactiveVars {
		fmt.Fprintf(&b, "\tupdate%s()\n", capitalize(varName))
	}

	b.WriteString("\n\tselect {}\n")
	b.WriteString("}\n\n")

	// Generate transformed user functions
	b.WriteString("// User functions\n")
	for funcName, fn := range a.Funcs {
		transformed := transformFunction(p.Script, funcName, fn.Modifies, p.ReactiveVars)
		if transformed != "" {
			b.WriteString(transformed)
			b.WriteString("\n\n")
		}
	}

	// bindBranchElements
	b.WriteString("func bindBranchElements() {\n")
	for _, expr := range exprs {
		if isExprInIf(expr.ID, ifs) {
			fmt.Fprintf(&b, "\tif el := document.Call(\"getElementById\", \"%s\"); !el.IsNull() {\n", expr.ID)
			fmt.Fprintf(&b, "\t\t%sEl = el\n", expr.ID)
			b.WriteString("\t}\n")
		}
	}

	branchEvents := filterBranchEvents(events)
	for _, ev := range branchEvents {
		fmt.Fprintf(&b, "\tif el := document.Call(\"getElementById\", \"%s\"); !el.IsNull() && el.Get(\"_bound\").IsUndefined() {\n", ev.ID)
		fmt.Fprintf(&b, "\t\tel.Set(\"_bound\", true)\n")
		fmt.Fprintf(&b, "\t\tel.Call(\"addEventListener\", \"%s\", js.FuncOf(func(this js.Value, args []js.Value) any {\n", ev.Event)
		generateEventModifiers(&b, ev.Modifiers, "\t\t\t")
		fmt.Fprintf(&b, "\t\t\t%s\n", transformEventHandler(ev.Handler))
		b.WriteString("\t\t\treturn nil\n")
		b.WriteString("\t\t}))\n")
		b.WriteString("\t}\n")
	}

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
	html = generateHTML(p, tmplAST, events, inputBindings, ifs)
	return
}

func generateIfUpdate(b *strings.Builder, ifn IfNode) {
	for i, branch := range ifn.Branches {
		branchName := fmt.Sprintf("branch%d", i)

		if i == 0 {
			fmt.Fprintf(b, "\tif %s {\n", branch.Cond)
		} else if branch.Cond != "" {
			fmt.Fprintf(b, " else if %s {\n", branch.Cond)
		} else {
			b.WriteString(" else {\n")
		}

		fmt.Fprintf(b, "\t\tif %s_showing != \"%s\" {\n", ifn.CondID, branchName)
		fmt.Fprintf(b, "\t\t\tif !%s_current.IsUndefined() && !%s_current.IsNull() {\n", ifn.CondID, ifn.CondID)
		fmt.Fprintf(b, "\t\t\t\t%s_current.Call(\"remove\")\n", ifn.CondID)
		b.WriteString("\t\t\t}\n")
		fmt.Fprintf(b, "\t\t\tfrag := createFragment(%s_%s_tmpl)\n", ifn.CondID, branchName)
		fmt.Fprintf(b, "\t\t\twrapper := document.Call(\"createElement\", \"span\")\n")
		fmt.Fprintf(b, "\t\t\twrapper.Call(\"appendChild\", frag)\n")
		fmt.Fprintf(b, "\t\t\t%s_anchor.Get(\"parentNode\").Call(\"insertBefore\", wrapper, %s_anchor)\n", ifn.CondID, ifn.CondID)
		fmt.Fprintf(b, "\t\t\t%s_current = wrapper\n", ifn.CondID)
		fmt.Fprintf(b, "\t\t\t%s_showing = \"%s\"\n", ifn.CondID, branchName)
		b.WriteString("\t\t\tbindBranchElements()\n")
		b.WriteString("\t\t}\n")
		b.WriteString("\t}")
	}

	lastBranch := ifn.Branches[len(ifn.Branches)-1]
	if lastBranch.Cond != "" {
		b.WriteString(" else {\n")
		fmt.Fprintf(b, "\t\tif %s_showing != \"none\" {\n", ifn.CondID)
		fmt.Fprintf(b, "\t\t\tif !%s_current.IsUndefined() && !%s_current.IsNull() {\n", ifn.CondID, ifn.CondID)
		fmt.Fprintf(b, "\t\t\t\t%s_current.Call(\"remove\")\n", ifn.CondID)
		b.WriteString("\t\t\t}\n")
		fmt.Fprintf(b, "\t\t\t%s_current = js.Null()\n", ifn.CondID)
		fmt.Fprintf(b, "\t\t\t%s_showing = \"none\"\n", ifn.CondID)
		b.WriteString("\t\t}\n")
		b.WriteString("\t}")
	}
	b.WriteString("\n")
}

func generateEachUpdate(b *strings.Builder, each EachNode, events []Event, exprs []ExprNode) {
	indexVar := each.Index
	if indexVar == "" {
		indexVar = "_i"
	}

	eachEvts := filterEachEvents(events, each.ID)

	fmt.Fprintf(b, "\t// Update %s\n", each.ID)
	fmt.Fprintf(b, "\tnewLen%s := len(%s)\n", each.ID, each.Array)
	fmt.Fprintf(b, "\toldLen%s := %s_count\n", each.ID, each.ID)
	fmt.Fprintf(b, "\t%s_count = newLen%s\n", each.ID, each.ID)

	fmt.Fprintf(b, "\tfor i := newLen%s; i < oldLen%s; i++ {\n", each.ID, each.ID)
	fmt.Fprintf(b, "\t\tif el := document.Call(\"getElementById\", \"%s_\" + strconv.Itoa(i)); !el.IsNull() {\n", each.ID)
	fmt.Fprintf(b, "\t\t\tel.Call(\"remove\")\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")

	fmt.Fprintf(b, "\tfor %s := 0; %s < newLen%s; %s++ {\n", indexVar, indexVar, each.ID, indexVar)
	fmt.Fprintf(b, "\t\t%s := %s[%s]\n", each.Item, each.Array, indexVar)
	fmt.Fprintf(b, "\t\telID := \"%s_\" + strconv.Itoa(%s)\n", each.ID, indexVar)
	fmt.Fprintf(b, "\t\tif el := document.Call(\"getElementById\", elID); el.IsNull() {\n")
	fmt.Fprintf(b, "\t\t\thtml := %s_tmpl\n", each.ID)
	fmt.Fprintf(b, "\t\t\thtml = strings.ReplaceAll(html, \"${_idx}\", strconv.Itoa(%s))\n", indexVar)
	fmt.Fprintf(b, "\t\t\tfrag := createFragment(html)\n")
	fmt.Fprintf(b, "\t\t\twrapper := document.Call(\"createElement\", \"span\")\n")
	fmt.Fprintf(b, "\t\t\twrapper.Set(\"id\", elID)\n")
	fmt.Fprintf(b, "\t\t\twrapper.Call(\"appendChild\", frag)\n")
	fmt.Fprintf(b, "\t\t\t%s_anchor.Get(\"parentNode\").Call(\"insertBefore\", wrapper, %s_anchor)\n", each.ID, each.ID)

	for i, ev := range eachEvts {
		fmt.Fprintf(b, "\t\t\tif btn := wrapper.Call(\"querySelector\", \"#btn%d_${_idx}\".replace(\"${_idx}\", strconv.Itoa(%s))); !btn.IsNull() || true {\n", i, indexVar)
		fmt.Fprintf(b, "\t\t\t\tbtnEl := document.Call(\"getElementById\", \"btn%d_\" + strconv.Itoa(%s))\n", i, indexVar)
		fmt.Fprintf(b, "\t\t\t\tif !btnEl.IsNull() && btnEl.Get(\"_bound\").IsUndefined() {\n")
		fmt.Fprintf(b, "\t\t\t\t\tbtnEl.Set(\"_bound\", true)\n")
		fmt.Fprintf(b, "\t\t\t\t\tidx := %s\n", indexVar)
		fmt.Fprintf(b, "\t\t\t\t\t_ = idx\n")
		if each.Item != "" {
			fmt.Fprintf(b, "\t\t\t\t\titemVal := %s\n", each.Item)
			fmt.Fprintf(b, "\t\t\t\t\t_ = itemVal\n")
		}
		fmt.Fprintf(b, "\t\t\t\t\tbtnEl.Call(\"addEventListener\", \"%s\", js.FuncOf(func(this js.Value, args []js.Value) any {\n", ev.Event)
		generateEventModifiers(b, ev.Modifiers, "\t\t\t\t\t\t")
		handler := ev.Handler
		if each.Index != "" {
			handler = replaceVar(handler, each.Index, "idx")
		}
		if each.Item != "" {
			handler = replaceVar(handler, each.Item, "itemVal")
		}
		fmt.Fprintf(b, "\t\t\t\t\t\t%s\n", transformEventHandler(handler))
		b.WriteString("\t\t\t\t\t\treturn nil\n")
		b.WriteString("\t\t\t\t\t}))\n")
		b.WriteString("\t\t\t\t}\n")
		b.WriteString("\t\t\t}\n")
	}

	b.WriteString("\t\t}\n")

	eachExprs := collectExprsInEach(each)
	for _, expr := range eachExprs {
		fmt.Fprintf(b, "\t\tif exprEl := document.Call(\"getElementById\", \"%s_\" + strconv.Itoa(%s)); !exprEl.IsNull() {\n", expr.ID, indexVar)
		fmt.Fprintf(b, "\t\t\texprEl.Set(\"textContent\", %s)\n", expr.Expr)
		b.WriteString("\t\t}\n")
	}

	b.WriteString("\t}\n")
}

// Helper functions

func exprDependsOnAny(expr ExprNode, vars []string) bool {
	for _, v := range vars {
		for _, dep := range expr.VarDeps {
			if dep == v {
				return true
			}
		}
	}
	return false
}

func exprDependsOnAnyStr(deps []string, vars []string) bool {
	for _, v := range vars {
		for _, dep := range deps {
			if dep == v {
				return true
			}
		}
	}
	return false
}

func ifDependsOnAny(ifn IfNode, vars []string) bool {
	for _, branch := range ifn.Branches {
		if exprDependsOnAnyStr(branch.VarDeps, vars) {
			return true
		}
	}
	return false
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func inferTypeFromReactive(reactiveVars map[string]string, expr string) string {
	// Simple case: expr is just a var name
	if typ, ok := reactiveVars[expr]; ok {
		return typ
	}
	// Default to string for complex expressions
	return "string"
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func getExpression(script, name string) string {
	for _, line := range strings.Split(script, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, name+" = ") {
			return strings.TrimPrefix(trimmed, name+" = ")
		}
	}
	return "0"
}

// Remaining helper functions from original codegen...

func collectEvents(template string, ifs []IfNode, eaches []EachNode) []Event {
	var events []Event
	btnID := 0

	for i := 0; i < len(template); i++ {
		if template[i] != '@' {
			continue
		}
		if i > 0 && template[i-1] == '{' {
			continue
		}
		eqIdx := strings.Index(template[i:], "=\"")
		if eqIdx == -1 {
			continue
		}
		eventPart := template[i+1 : i+eqIdx]
		parts := strings.Split(eventPart, "|")
		eventName := parts[0]
		modifiers := parts[1:]

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

func collectBindings(template, script string, ifs []IfNode, reactiveVarTypes map[string]string) []Binding {
	var bindings []Binding
	inputID := 0

	for i := 0; i < len(template); i++ {
		if idx := strings.Index(template[i:], "bind:value=\""); idx != -1 {
			absPos := i + idx
			start := absPos + 12
			end := strings.Index(template[start:], "\"")
			if end != -1 {
				varName := template[start : start+end]
				varType := "string"
				if t, ok := reactiveVarTypes[varName]; ok && t == "int" {
					varType = "int"
				}
				inBranch := isPositionInIf(template, absPos, ifs)
				bindings = append(bindings, Binding{
					ID:       fmt.Sprintf("input%d", inputID),
					Attr:     "value",
					VarName:  varName,
					VarType:  varType,
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

func collectClassBindings(template string, ifs []IfNode, reactiveVarNames map[string]bool) []ClassBinding {
	var bindings []ClassBinding
	classID := 0
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

		// Find var deps in condition
		var varDeps []string
		for varName := range reactiveVarNames {
			if strings.Contains(cond, varName) {
				varDeps = append(varDeps, varName)
			}
		}

		inBranch := isPositionInIf(template, absPos, ifs)
		bindings = append(bindings, ClassBinding{
			ID:        fmt.Sprintf("class%d", classID),
			ClassName: className,
			Cond:      cond,
			InBranch:  inBranch,
			VarDeps:   varDeps,
		})
		classID++
		i = fullEnd
	}

	return bindings
}

func isPositionInIf(template string, pos int, ifs []IfNode) bool {
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
		if !isPositionInIf(template, absPos, ifs) {
			count++
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
	i := 0
	for i < len(html) {
		idx := strings.Index(html[i:], "@")
		if idx == -1 {
			break
		}
		absIdx := i + idx
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
		html = html[:idx] + fmt.Sprintf("id=\"class%d\"", classID) + html[fullEnd:]
		classID++
	}

	return html, btnID, inputID, classID
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

func transformEventHandler(handler string) string {
	handler = strings.ReplaceAll(handler, "(e)", "(args[0])")
	handler = strings.ReplaceAll(handler, "(e,", "(args[0],")
	handler = strings.ReplaceAll(handler, ", e)", ", args[0])")
	handler = strings.ReplaceAll(handler, ",e)", ",args[0])")
	return handler
}

func isExprInIf(exprID string, ifs []IfNode) bool {
	for _, ifn := range ifs {
		for _, branch := range ifn.Branches {
			if exprInNodes(exprID, branch.Body) {
				return true
			}
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
			for _, branch := range node.Branches {
				if exprInNodes(exprID, branch.Body) {
					return true
				}
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
			for _, branch := range node.Branches {
				collectExprsFromNodes(branch.Body, exprs)
			}
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

func generateEachHTMLNodes(b *strings.Builder, nodes []Node, eachID string) {
	for _, n := range nodes {
		switch node := n.(type) {
		case TextNode:
			b.WriteString(node.Text)
		case ExprNode:
			fmt.Fprintf(b, `<span id="%s_${_idx}"></span>`, node.ID)
		case IfNode:
			fmt.Fprintf(b, `<span id="%s_anchor" style="display:none"></span>`, node.CondID)
		case EachNode:
			fmt.Fprintf(b, `<span id="%s_anchor" style="display:none"></span>`, node.ID)
		}
	}
}

func processEachTemplateAttrs(html string) string {
	btnID := 0
	i := 0
	for i < len(html) {
		idx := strings.Index(html[i:], "@")
		if idx == -1 {
			break
		}
		absIdx := i + idx
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
		html = html[:absIdx] + fmt.Sprintf("id=\"btn%d_${_idx}\"", btnID) + html[handlerStart+handlerEnd+1:]
		btnID++
		i = absIdx + 1
	}
	return html
}

func replaceVar(s, oldVar, newVar string) string {
	result := s
	result = strings.ReplaceAll(result, "("+oldVar+")", "("+newVar+")")
	result = strings.ReplaceAll(result, "("+oldVar+",", "("+newVar+",")
	result = strings.ReplaceAll(result, ","+oldVar+")", ","+newVar+")")
	result = strings.ReplaceAll(result, ", "+oldVar+")", ", "+newVar+")")
	return result
}

// transformFunction transforms a user function to use setters for reactive vars
func transformFunction(script, funcName string, modifies []string, reactiveVars map[string]string) string {
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

	// If no modifications, return as-is
	if len(modifies) == 0 {
		return funcSrc
	}

	for _, varName := range modifies {
		// Only transform if it's a reactive var
		if _, isReactive := reactiveVars[varName]; !isReactive {
			continue
		}

		// Replace varName++ with setVarName(varName + 1)
		funcSrc = strings.ReplaceAll(funcSrc, varName+"++", fmt.Sprintf("set%s(%s + 1)", capitalize(varName), varName))
		
		// Replace varName-- with setVarName(varName - 1)
		funcSrc = strings.ReplaceAll(funcSrc, varName+"--", fmt.Sprintf("set%s(%s - 1)", capitalize(varName), varName))

		// Replace varName = expr with setVarName(expr)
		lines := strings.Split(funcSrc, "\n")
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
			
			// Handle: varName += expr
			if strings.HasPrefix(trimmed, varName+" += ") {
				expr := strings.TrimPrefix(trimmed, varName+" += ")
				lines[i] = indent + fmt.Sprintf("set%s(%s + %s)", capitalize(varName), varName, expr)
				continue
			}
			
			// Handle: varName -= expr
			if strings.HasPrefix(trimmed, varName+" -= ") {
				expr := strings.TrimPrefix(trimmed, varName+" -= ")
				lines[i] = indent + fmt.Sprintf("set%s(%s - %s)", capitalize(varName), varName, expr)
				continue
			}
			
			// Handle: varName *= expr
			if strings.HasPrefix(trimmed, varName+" *= ") {
				expr := strings.TrimPrefix(trimmed, varName+" *= ")
				lines[i] = indent + fmt.Sprintf("set%s(%s * %s)", capitalize(varName), varName, expr)
				continue
			}
			
			// Handle: varName /= expr
			if strings.HasPrefix(trimmed, varName+" /= ") {
				expr := strings.TrimPrefix(trimmed, varName+" /= ")
				lines[i] = indent + fmt.Sprintf("set%s(%s / %s)", capitalize(varName), varName, expr)
				continue
			}
			
			// Handle: varName = expr
			if strings.HasPrefix(trimmed, varName+" = ") {
				expr := strings.TrimPrefix(trimmed, varName+" = ")
				lines[i] = indent + fmt.Sprintf("set%s(%s)", capitalize(varName), expr)
				continue
			}
			
			// Handle: varName[i] = expr (slice/map index) -> keep it, add update call
			if strings.HasPrefix(trimmed, varName+"[") && strings.Contains(trimmed, "] = ") {
				lines[i] = line + "\n" + indent + fmt.Sprintf("update%s()", capitalize(varName))
				continue
			}
			
			// Handle: delete(varName, key)
			if strings.HasPrefix(trimmed, "delete("+varName+",") {
				lines[i] = line + "\n" + indent + fmt.Sprintf("update%s()", capitalize(varName))
				continue
			}
			
			// Handle: clear(varName)
			if trimmed == "clear("+varName+")" || strings.HasPrefix(trimmed, "clear("+varName+")") {
				lines[i] = line + "\n" + indent + fmt.Sprintf("update%s()", capitalize(varName))
				continue
			}
		}
		funcSrc = strings.Join(lines, "\n")
	}

	return funcSrc
}

func generateHTML(p Parsed, tmplAST *TemplateAST, events []Event, inputBindings []Binding, ifs []IfNode) string {
	var b strings.Builder

	b.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	b.WriteString("  <script src=\"wasm_exec.js\"></script>\n")
	b.WriteString("  <script>\n")
	b.WriteString("    const go = new Go();\n")
	b.WriteString("    WebAssembly.instantiateStreaming(fetch(\"app.wasm\"), go.importObject)\n")
	b.WriteString("      .then(r => go.run(r.instance));\n")
	b.WriteString("  </script>\n")

	if p.Style != "" {
		b.WriteString("  <style>\n")
		b.WriteString(p.Style)
		b.WriteString("\n  </style>\n")
	}
	b.WriteString("</head>\n<body>\n")

	html := tmplAST.GenerateHTML()

	// Replace @event|modifiers="handler" with id="btnN"
	btnID := 0
	i := 0
	for i < len(html) {
		idx := strings.Index(html[i:], "@")
		if idx == -1 {
			break
		}
		absIdx := i + idx
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
		html = html[:idx] + fmt.Sprintf("id=\"class%d\"", classID) + html[fullEnd:]
		classID++
	}

	b.WriteString(html)
	b.WriteString("\n</body>\n</html>\n")
	return b.String()
}