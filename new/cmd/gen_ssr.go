package main

import (
	"fmt"
	"strings"
)

// SSRContext contains all information needed to generate SSR rendering code.
type SSRContext struct {
	ID            string                // Full prefixed ID
	Name          string                // Component type name
	Definition    *component            // Component definition
	Binding       componentBinding      // How component is used
	Scope         *Scope                // Field resolution scope
	Parent        *SSRContext           // Parent context
	AllComponents map[string]*component // All available components
	Indent        string                // Current indentation
}

// ChildContext creates a new context for a child component.
func (ctx *SSRContext) ChildContext(childBinding componentBinding, childDef *component) *SSRContext {
	childID := ctx.prefixID(childBinding.elementID)
	childScope := &Scope{
		VarName:    childID,
		FieldTypes: buildFieldTypes(childDef),
		Parent:     ctx.Scope,
	}
	return &SSRContext{
		ID:            childID,
		Name:          childBinding.name,
		Definition:    childDef,
		Binding:       childBinding,
		Scope:         childScope,
		Parent:        ctx,
		AllComponents: ctx.AllComponents,
		Indent:        ctx.Indent + "\t",
	}
}

func (ctx *SSRContext) prefixID(localID string) string {
	if ctx.ID == "" {
		return localID
	}
	return ctx.ID + "_" + localID
}

func generateRender(comp *component, tmpl string, bindings templateBindings, childComponents map[string]*component) string {
	var sb strings.Builder
	fieldTypes := buildFieldTypes(comp)

	// Check imports
	needsStrconv := len(bindings.eachBlocks) > 0
	needsStrings := len(bindings.expressions) > 0 || len(bindings.attrBindings) > 0 ||
		len(bindings.ifBlocks) > 0 || len(bindings.eachBlocks) > 0 || len(bindings.components) > 0

	// Write header
	sb.WriteString("//go:build !wasm\n\npackage main\n\nimport (\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString("\t\"os\"\n")
	sb.WriteString("\t\"preveltekit\"\n")
	if needsStrings {
		sb.WriteString("\t\"strings\"\n")
	}
	if needsStrconv {
		sb.WriteString("\t\"strconv\"\n")
	}
	sb.WriteString(")\n\nfunc main() {\n\tcomponent := &" + comp.name + "{\n")

	// Initialize fields
	generateFieldInit(&sb, comp.fields, "\t\t")
	sb.WriteString("\t}\n")

	if comp.hasOnMount {
		sb.WriteString("\n\tcomponent.OnMount()\n")
	}

	// Create root context
	rootScope := &Scope{
		VarName:    "component",
		FieldTypes: fieldTypes,
		Parent:     nil,
	}
	rootCtx := &SSRContext{
		ID:            "",
		Name:          comp.name,
		Definition:    comp,
		Scope:         rootScope,
		Parent:        nil,
		AllComponents: childComponents,
		Indent:        "\t",
	}

	// Create all child component instances
	generateSSRComponentInstances(&sb, bindings.components, rootCtx)

	// Template rendering
	fmt.Fprintf(&sb, "\n\thtml := %s\n", escapeForGoString(tmpl))

	// Generate SSR rendering code for all bindings
	generateSSRBindings(&sb, bindings, rootCtx)

	// Render child components
	for _, compBinding := range bindings.components {
		childDef := childComponents[compBinding.name]
		if childDef != nil {
			childCtx := rootCtx.ChildContext(compBinding, childDef)
			generateSSRComponent(&sb, childCtx)
		}
	}

	sb.WriteString("\n\tfmt.Fprint(os.Stdout, html)\n}\n")

	if needsStrings {
		sb.WriteString("\nvar _ = strings.Replace\n")
	}
	if needsStrconv {
		sb.WriteString("var _ = strconv.Itoa\n")
	}

	return sb.String()
}

// generateSSRComponentInstances generates instance creation for all components.
func generateSSRComponentInstances(sb *strings.Builder, components []componentBinding, ctx *SSRContext) {
	for _, compBinding := range components {
		childDef := ctx.AllComponents[compBinding.name]
		if childDef == nil {
			continue
		}

		childID := ctx.prefixID(compBinding.elementID)
		needsVar := len(childDef.fields) > 0 || len(compBinding.props) > 0 || childDef.hasOnMount

		if needsVar {
			fmt.Fprintf(sb, "\n\t%s := &%s{\n", childID, compBinding.name)
		} else {
			fmt.Fprintf(sb, "\n\t_ = &%s{\n", compBinding.name)
		}
		generateFieldInit(sb, childDef.fields, "\t\t")
		sb.WriteString("\t}\n")

		if !needsVar {
			continue
		}

		// Set props
		childFieldTypes := buildFieldTypes(childDef)
		for propName, propValue := range compBinding.props {
			childField := strings.Title(propName)
			if strings.HasPrefix(propValue, "{") && strings.HasSuffix(propValue, "}") {
				parentField := propValue[1 : len(propValue)-1]
				fmt.Fprintf(sb, "\t%s.%s.Set(%s.%s.Get())\n", childID, childField, ctx.Scope.VarName, parentField)
			} else {
				childFieldType := childFieldTypes[childField]
				switch childFieldType {
				case "string":
					fmt.Fprintf(sb, "\t%s.%s.Set(%q)\n", childID, childField, propValue)
				case "int", "bool":
					fmt.Fprintf(sb, "\t%s.%s.Set(%s)\n", childID, childField, propValue)
				default:
					fmt.Fprintf(sb, "\t%s.%s.Set(%q)\n", childID, childField, propValue)
				}
			}
		}

		if childDef.hasOnMount {
			fmt.Fprintf(sb, "\t%s.OnMount()\n", childID)
		}

		// Recursively create instances for nested components
		childTmpl := strings.ReplaceAll(childDef.template, "<slot/>", compBinding.children)
		_, childBindings := parseTemplate(childTmpl)

		childCtx := ctx.ChildContext(compBinding, childDef)
		generateSSRComponentInstances(sb, childBindings.components, childCtx)
	}
}

// generateSSRBindings generates SSR rendering for bindings.
func generateSSRBindings(sb *strings.Builder, bindings templateBindings, ctx *SSRContext) {
	varName := ctx.Scope.VarName

	// Expression substitutions
	for _, expr := range bindings.expressions {
		fullID := ctx.prefixID(expr.elementID)
		varRef := varName + "." + expr.fieldName
		fmt.Fprintf(sb, "\thtml = strings.Replace(html, \"<span id=\\\"%s\\\"></span>\", fmt.Sprintf(\"<span id=\\\"%s\\\">%%v</span>\", %s.Get()), 1)\n",
			fullID, fullID, varRef)
	}

	// Attribute bindings
	for _, ab := range bindings.attrBindings {
		staticValue := ab.template
		for _, field := range ab.fields {
			staticValue = strings.ReplaceAll(staticValue, "{"+field+"}", "")
		}
		staticValue = strings.TrimSpace(staticValue)

		fmt.Fprintf(sb, "\t{\n\t\tattrVal := %s\n", escapeForGoString(ab.template))
		for _, field := range ab.fields {
			fmt.Fprintf(sb, "\t\tattrVal = strings.ReplaceAll(attrVal, \"{%s}\", fmt.Sprintf(\"%%v\", %s.%s.Get()))\n", field, varName, field)
		}
		fmt.Fprintf(sb, "\t\thtml = strings.Replace(html, \"%s=\\\"%s\\\"\", \"%s=\\\"\" + attrVal + \"\\\"\", 1)\n\t}\n",
			ab.attrName, staticValue, ab.attrName)
	}

	// If blocks
	fieldTypes := ctx.Scope.FieldTypes
	for _, ifb := range bindings.ifBlocks {
		fullID := ctx.prefixID(ifb.elementID)
		fmt.Fprintf(sb, "\t{\n\t\tvar ifContent string\n")
		for i, branch := range ifb.branches {
			cond := branch.condition
			for fieldName := range fieldTypes {
				if strings.Contains(cond, fieldName) {
					cond = strings.ReplaceAll(cond, fieldName, varName+"."+fieldName+".Get()")
				}
			}
			if i == 0 {
				fmt.Fprintf(sb, "\t\tif %s {\n\t\t\tifContent = %s\n", cond, escapeForGoString(branch.html))
			} else {
				fmt.Fprintf(sb, "\t\t} else if %s {\n\t\t\tifContent = %s\n", cond, escapeForGoString(branch.html))
			}
		}
		if ifb.elseHTML != "" {
			fmt.Fprintf(sb, "\t\t} else {\n\t\t\tifContent = %s\n\t\t}\n", escapeForGoString(ifb.elseHTML))
		} else {
			sb.WriteString("\t\t}\n")
		}
		fmt.Fprintf(sb, "\t\thtml = strings.Replace(html, \"<span id=\\\"%s_anchor\\\"></span>\", ifContent + \"<span id=\\\"%s_anchor\\\"></span>\", 1)\n\t}\n",
			fullID, fullID)
	}

	// Each blocks
	for _, each := range bindings.eachBlocks {
		fullID := ctx.prefixID(each.elementID)
		fmt.Fprintf(sb, "\t{\n\t\tvar eachContent strings.Builder\n")
		fmt.Fprintf(sb, "\t\titems := %s.%s.Get()\n", varName, each.listName)
		if each.elseHTML != "" {
			fmt.Fprintf(sb, "\t\tif len(items) == 0 {\n")
			fmt.Fprintf(sb, "\t\t\teachContent.WriteString(\"<span id=\\\"%s_else\\\">%s</span>\")\n", fullID, escapeForGoStringContent(each.elseHTML))
			fmt.Fprintf(sb, "\t\t} else {\n")
			fmt.Fprintf(sb, "\t\t\teachContent.WriteString(\"<span id=\\\"%s_else\\\" style=\\\"display:none\\\">%s</span>\")\n", fullID, escapeForGoStringContent(each.elseHTML))
		}
		fmt.Fprintf(sb, "\t\tfor %s, %s := range items {\n", each.indexVar, each.itemVar)
		fmt.Fprintf(sb, "\t\t\titemHTML := %s\n", escapeForGoString(each.bodyHTML))
		fmt.Fprintf(sb, "\t\t\titemHTML = strings.ReplaceAll(itemHTML, \"{%s}\", fmt.Sprintf(\"%%v\", %s))\n", each.itemVar, each.itemVar)
		fmt.Fprintf(sb, "\t\t\titemHTML = strings.ReplaceAll(itemHTML, \"{%s}\", strconv.Itoa(%s))\n", each.indexVar, each.indexVar)
		fmt.Fprintf(sb, "\t\t\teachContent.WriteString(\"<span id=\\\"%s_\" + strconv.Itoa(%s) + \"\\\">\" + itemHTML + \"</span>\")\n", fullID, each.indexVar)
		fmt.Fprintf(sb, "\t\t}\n")
		if each.elseHTML != "" {
			fmt.Fprintf(sb, "\t\t}\n")
		}
		fmt.Fprintf(sb, "\t\thtml = strings.Replace(html, \"<span id=\\\"%s_anchor\\\"></span>\", eachContent.String() + \"<span id=\\\"%s_anchor\\\"></span>\", 1)\n\t}\n",
			fullID, fullID)
	}
}

// generateSSRComponent generates SSR rendering for a child component.
func generateSSRComponent(sb *strings.Builder, ctx *SSRContext) {
	compID := ctx.ID
	compDef := ctx.Definition
	compBinding := ctx.Binding
	childFieldTypes := buildFieldTypes(compDef)

	// Parse slot content to identify parent fields
	_, slotBindings := parseTemplate(compBinding.children)
	slotFields := make(map[string]bool)
	for _, expr := range slotBindings.expressions {
		slotFields[expr.fieldName] = true
	}

	// Process template with slot content
	childTmpl := strings.ReplaceAll(compDef.template, "<slot/>", compBinding.children)
	processedTmpl, childBindings := parseTemplate(childTmpl)

	// Get parent field types
	var parentFieldTypes map[string]string
	if ctx.Parent != nil {
		parentFieldTypes = ctx.Parent.Scope.FieldTypes
	} else {
		parentFieldTypes = make(map[string]string)
	}

	// Categorize expressions
	slotExprs, childOwnExprs := categorizeExpressions(childBindings.expressions, slotFields, parentFieldTypes, childFieldTypes)

	// Prefix all IDs
	processedTmpl = prefixAllBindingIDs(compID, processedTmpl, &childBindings)

	// Update categorized expressions with prefixed IDs
	for i := range slotExprs {
		slotExprs[i].elementID = compID + "_" + slotExprs[i].elementID
	}
	for i := range childOwnExprs {
		childOwnExprs[i].elementID = compID + "_" + childOwnExprs[i].elementID
	}

	// Inject component ID
	processedTmpl = injectIDIntoFirstTag(processedTmpl, compID)

	fmt.Fprintf(sb, "\t{\n\t\tchildHTML := %s\n", escapeForGoString(processedTmpl))

	// Child's own expressions
	for _, expr := range childOwnExprs {
		fmt.Fprintf(sb, "\t\tchildHTML = strings.Replace(childHTML, \"<span id=\\\"%s\\\"></span>\", fmt.Sprintf(\"<span id=\\\"%s\\\">%%v</span>\", %s.%s.Get()), 1)\n",
			expr.elementID, expr.elementID, compID, expr.fieldName)
	}

	// Slot expressions (parent bindings)
	parentVarName := "component"
	if ctx.Parent != nil {
		parentVarName = ctx.Parent.Scope.VarName
	}
	for _, expr := range slotExprs {
		fmt.Fprintf(sb, "\t\tchildHTML = strings.Replace(childHTML, \"<span id=\\\"%s\\\"></span>\", fmt.Sprintf(\"<span id=\\\"%s\\\">%%v</span>\", %s.%s.Get()), 1)\n",
			expr.elementID, expr.elementID, parentVarName, expr.fieldName)
	}

	// Attribute bindings
	for _, ab := range childBindings.attrBindings {
		staticValue := ab.template
		for _, field := range ab.fields {
			staticValue = strings.ReplaceAll(staticValue, "{"+field+"}", "")
		}
		staticValue = strings.TrimSpace(staticValue)

		fmt.Fprintf(sb, "\t\t{\n\t\t\tattrVal := %s\n", escapeForGoString(ab.template))
		for _, field := range ab.fields {
			fmt.Fprintf(sb, "\t\t\tattrVal = strings.ReplaceAll(attrVal, \"{%s}\", fmt.Sprintf(\"%%v\", %s.%s.Get()))\n",
				field, compID, field)
		}
		fmt.Fprintf(sb, "\t\t\tchildHTML = strings.Replace(childHTML, \"%s=\\\"%s\\\"\", \"%s=\\\"\" + attrVal + \"\\\"\", 1)\n\t\t}\n",
			ab.attrName, staticValue, ab.attrName)
	}

	// Render nested components recursively
	for _, nestedBinding := range childBindings.components {
		nestedDef := ctx.AllComponents[nestedBinding.name]
		if nestedDef != nil {
			nestedCtx := ctx.ChildContext(nestedBinding, nestedDef)
			generateSSRNestedComponent(sb, nestedCtx, "childHTML")
		}
	}

	// Replace placeholder with rendered HTML
	fmt.Fprintf(sb, "\t\thtml = strings.Replace(html, \"<!--%s-->\", childHTML, 1)\n\t}\n", compBinding.elementID)
}

// generateSSRNestedComponent generates SSR rendering for a nested component.
func generateSSRNestedComponent(sb *strings.Builder, ctx *SSRContext, htmlVar string) {
	compID := ctx.ID
	compDef := ctx.Definition
	compBinding := ctx.Binding
	childFieldTypes := buildFieldTypes(compDef)

	// Process template with slot content
	childTmpl := strings.ReplaceAll(compDef.template, "<slot/>", compBinding.children)
	processedTmpl, childBindings := parseTemplate(childTmpl)

	// Prefix all IDs
	processedTmpl = prefixAllBindingIDs(compID, processedTmpl, &childBindings)

	// Inject component ID
	processedTmpl = injectIDIntoFirstTag(processedTmpl, compID)

	fmt.Fprintf(sb, "\t\t{\n\t\t\tnestedHTML := %s\n", escapeForGoString(processedTmpl))

	// Expressions
	for _, expr := range childBindings.expressions {
		fmt.Fprintf(sb, "\t\t\tnestedHTML = strings.Replace(nestedHTML, \"<span id=\\\"%s\\\"></span>\", fmt.Sprintf(\"<span id=\\\"%s\\\">%%v</span>\", %s.%s.Get()), 1)\n",
			expr.elementID, expr.elementID, compID, expr.fieldName)
	}

	// Attribute bindings
	for _, ab := range childBindings.attrBindings {
		staticValue := ab.template
		for _, field := range ab.fields {
			staticValue = strings.ReplaceAll(staticValue, "{"+field+"}", "")
		}
		staticValue = strings.TrimSpace(staticValue)

		fmt.Fprintf(sb, "\t\t\t{\n\t\t\t\tattrVal := %s\n", escapeForGoString(ab.template))
		for _, field := range ab.fields {
			fieldType := childFieldTypes[field]
			if fieldType == "" {
				fieldType = "string"
			}
			fmt.Fprintf(sb, "\t\t\t\tattrVal = strings.ReplaceAll(attrVal, \"{%s}\", fmt.Sprintf(\"%%v\", %s.%s.Get()))\n",
				field, compID, field)
		}
		fmt.Fprintf(sb, "\t\t\t\tnestedHTML = strings.Replace(nestedHTML, \"%s=\\\"%s\\\"\", \"%s=\\\"\" + attrVal + \"\\\"\", 1)\n\t\t\t}\n",
			ab.attrName, staticValue, ab.attrName)
	}

	// Recursively handle deeply nested components
	for _, nestedBinding := range childBindings.components {
		nestedDef := ctx.AllComponents[nestedBinding.name]
		if nestedDef != nil {
			nestedCtx := ctx.ChildContext(nestedBinding, nestedDef)
			generateSSRNestedComponent(sb, nestedCtx, "nestedHTML")
		}
	}

	// Replace placeholder - use the full prefixed ID since the parent's HTML
	// has been processed by prefixAllBindingIDs and contains prefixed placeholders
	fmt.Fprintf(sb, "\t\t\t%s = strings.Replace(%s, \"<!--%s-->\", nestedHTML, 1)\n\t\t}\n", htmlVar, htmlVar, compID)
}
