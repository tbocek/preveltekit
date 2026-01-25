package main

import (
	"fmt"
	"strings"
)

func generateRender(comp *component, tmpl string, bindings templateBindings, childComponents map[string]*component) string {
	var sb strings.Builder
	fieldTypes := buildFieldTypes(comp)

	// Check imports
	needsStrconv := len(bindings.eachBlocks) > 0
	needsStrings := len(bindings.expressions) > 0 || len(bindings.attrBindings) > 0 || len(bindings.ifBlocks) > 0 || len(bindings.eachBlocks) > 0
	for _, compBinding := range bindings.components {
		if childComponents[compBinding.name] != nil {
			needsStrings = true
			break
		}
	}

	// Write header - now imports reactive package
	sb.WriteString("//go:build !wasm\n\npackage main\n\nimport (\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString("\t\"os\"\n")
	sb.WriteString("\t\"reactive\"\n")
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

	// Child components
	for _, compBinding := range bindings.components {
		childComp := childComponents[compBinding.name]
		if childComp == nil {
			continue
		}

		// Only create component variable if it has fields, props, or OnMount
		needsVar := len(childComp.fields) > 0 || len(compBinding.props) > 0 || childComp.hasOnMount

		if needsVar {
			fmt.Fprintf(&sb, "\n\t%s := &%s{\n", compBinding.elementID, compBinding.name)
		} else {
			fmt.Fprintf(&sb, "\n\t_ = &%s{\n", compBinding.name)
		}
		generateFieldInit(&sb, childComp.fields, "\t\t")
		sb.WriteString("\t}\n")

		if !needsVar {
			continue
		}

		// Props
		for propName, propValue := range compBinding.props {
			childField := strings.Title(propName)
			if strings.HasPrefix(propValue, "{") && strings.HasSuffix(propValue, "}") {
				parentField := propValue[1 : len(propValue)-1]
				fmt.Fprintf(&sb, "\t%s.%s.Set(component.%s.Get())\n", compBinding.elementID, childField, parentField)
			} else {
				childFieldType := ""
				for _, f := range childComp.fields {
					if f.name == childField {
						childFieldType = f.valueType
						break
					}
				}
				switch childFieldType {
				case "string":
					fmt.Fprintf(&sb, "\t%s.%s.Set(%q)\n", compBinding.elementID, childField, propValue)
				case "int", "bool":
					fmt.Fprintf(&sb, "\t%s.%s.Set(%s)\n", compBinding.elementID, childField, propValue)
				default:
					fmt.Fprintf(&sb, "\t%s.%s.Set(%q)\n", compBinding.elementID, childField, propValue)
				}
			}
		}

		if childComp.hasOnMount {
			fmt.Fprintf(&sb, "\t%s.OnMount()\n", compBinding.elementID)
		}
	}

	// Template rendering
	fmt.Fprintf(&sb, "\n\thtml := %s\n", escapeForGoString(tmpl))

	// Expression substitutions
	for _, expr := range bindings.expressions {
		fmt.Fprintf(&sb, "\thtml = strings.Replace(html, \"<span id=\\\"%s\\\"></span>\", fmt.Sprintf(\"<span id=\\\"%s\\\">%%v</span>\", component.%s.Get()), 1)\n",
			expr.elementID, expr.elementID, expr.fieldName)
	}

	// Attribute bindings
	for _, ab := range bindings.attrBindings {
		staticValue := ab.template
		for _, field := range ab.fields {
			staticValue = strings.ReplaceAll(staticValue, "{"+field+"}", "")
		}
		staticValue = strings.TrimSpace(staticValue)

		fmt.Fprintf(&sb, "\t{\n\t\tattrVal := %s\n", escapeForGoString(ab.template))
		for _, field := range ab.fields {
			fmt.Fprintf(&sb, "\t\tattrVal = strings.ReplaceAll(attrVal, \"{%s}\", fmt.Sprintf(\"%%v\", component.%s.Get()))\n", field, field)
		}
		fmt.Fprintf(&sb, "\t\thtml = strings.Replace(html, \"%s=\\\"%s\\\"\", \"%s=\\\"\" + attrVal + \"\\\"\", 1)\n\t}\n",
			ab.attrName, staticValue, ab.attrName)
	}

	// If blocks
	for _, ifb := range bindings.ifBlocks {
		fmt.Fprintf(&sb, "\t{\n\t\tvar ifContent string\n")
		for i, branch := range ifb.branches {
			cond := branch.condition
			for fieldName := range fieldTypes {
				if strings.Contains(cond, fieldName) {
					cond = strings.ReplaceAll(cond, fieldName, "component."+fieldName+".Get()")
				}
			}
			if i == 0 {
				fmt.Fprintf(&sb, "\t\tif %s {\n\t\t\tifContent = %s\n", cond, escapeForGoString(branch.html))
			} else {
				fmt.Fprintf(&sb, "\t\t} else if %s {\n\t\t\tifContent = %s\n", cond, escapeForGoString(branch.html))
			}
		}
		if ifb.elseHTML != "" {
			fmt.Fprintf(&sb, "\t\t} else {\n\t\t\tifContent = %s\n\t\t}\n", escapeForGoString(ifb.elseHTML))
		} else {
			sb.WriteString("\t\t}\n")
		}
		fmt.Fprintf(&sb, "\t\thtml = strings.Replace(html, \"<span id=\\\"%s_anchor\\\"></span>\", ifContent + \"<span id=\\\"%s_anchor\\\"></span>\", 1)\n\t}\n",
			ifb.elementID, ifb.elementID)
	}

	// Each blocks
	for _, each := range bindings.eachBlocks {
		fmt.Fprintf(&sb, "\t{\n\t\tvar eachContent strings.Builder\n")
		fmt.Fprintf(&sb, "\t\titems := component.%s.Get()\n", each.listName)
		if each.elseHTML != "" {
			fmt.Fprintf(&sb, "\t\tif len(items) == 0 {\n")
			fmt.Fprintf(&sb, "\t\t\teachContent.WriteString(\"<span id=\\\"%s_else\\\">%s</span>\")\n", each.elementID, escapeForGoStringContent(each.elseHTML))
			fmt.Fprintf(&sb, "\t\t} else {\n")
			fmt.Fprintf(&sb, "\t\t\teachContent.WriteString(\"<span id=\\\"%s_else\\\" style=\\\"display:none\\\">%s</span>\")\n", each.elementID, escapeForGoStringContent(each.elseHTML))
		}
		fmt.Fprintf(&sb, "\t\tfor %s, %s := range items {\n", each.indexVar, each.itemVar)
		fmt.Fprintf(&sb, "\t\t\titemHTML := %s\n", escapeForGoString(each.bodyHTML))
		fmt.Fprintf(&sb, "\t\t\titemHTML = strings.ReplaceAll(itemHTML, \"{%s}\", fmt.Sprintf(\"%%v\", %s))\n", each.itemVar, each.itemVar)
		fmt.Fprintf(&sb, "\t\t\titemHTML = strings.ReplaceAll(itemHTML, \"{%s}\", strconv.Itoa(%s))\n", each.indexVar, each.indexVar)
		fmt.Fprintf(&sb, "\t\t\teachContent.WriteString(\"<span id=\\\"%s_\" + strconv.Itoa(%s) + \"\\\">\" + itemHTML + \"</span>\")\n", each.elementID, each.indexVar)
		fmt.Fprintf(&sb, "\t\t}\n")
		if each.elseHTML != "" {
			fmt.Fprintf(&sb, "\t\t}\n")
		}
		fmt.Fprintf(&sb, "\t\thtml = strings.Replace(html, \"<span id=\\\"%s_anchor\\\"></span>\", eachContent.String() + \"<span id=\\\"%s_anchor\\\"></span>\", 1)\n\t}\n",
			each.elementID, each.elementID)
	}

	// Child component rendering
	for _, compBinding := range bindings.components {
		childComp := childComponents[compBinding.name]
		if childComp == nil {
			continue
		}

		// Parse slot content to identify which fields come from parent
		_, slotBindings := parseTemplate(compBinding.children)
		slotFields := make(map[string]bool)
		for _, expr := range slotBindings.expressions {
			slotFields[expr.fieldName] = true
		}

		childTmpl := strings.ReplaceAll(childComp.template, "<slot/>", compBinding.children)
		childTmplProcessed, childBindings := parseTemplate(childTmpl)

		childFieldTypes := buildFieldTypes(childComp)

		// Categorize expressions into parent (slot) vs child owned
		slotExprs, childOwnExprs := categorizeExpressions(childBindings.expressions, slotFields, fieldTypes, childFieldTypes)

		// Make IDs unique - prefix all bindings
		childTmplProcessed = prefixBindingIDs(compBinding.elementID, childTmplProcessed,
			childOwnExprs, childBindings.events, childBindings.attrBindings, childBindings.ifBlocks)
		childTmplProcessed = prefixBindingIDs(compBinding.elementID, childTmplProcessed,
			slotExprs, nil, nil, nil)
		childTmplProcessed = prefixInputBindingIDs(compBinding.elementID, childTmplProcessed, childBindings.bindings)
		childTmplProcessed = prefixEachBindingIDs(compBinding.elementID, childTmplProcessed, childBindings.eachBlocks)
		childTmplProcessed = prefixClassBindingIDs(compBinding.elementID, childTmplProcessed, childBindings.classBindings)

		// Inject component ID into child's root element
		childTmplProcessed = injectIDIntoFirstTag(childTmplProcessed, compBinding.elementID)

		fmt.Fprintf(&sb, "\t{\n\t\tchildHTML := %s\n", escapeForGoString(childTmplProcessed))

		for _, expr := range childOwnExprs {
			fmt.Fprintf(&sb, "\t\tchildHTML = strings.Replace(childHTML, \"<span id=\\\"%s\\\"></span>\", fmt.Sprintf(\"<span id=\\\"%s\\\">%%v</span>\", %s.%s.Get()), 1)\n",
				expr.elementID, expr.elementID, compBinding.elementID, expr.fieldName)
		}

		for _, expr := range slotExprs {
			fmt.Fprintf(&sb, "\t\tchildHTML = strings.Replace(childHTML, \"<span id=\\\"%s\\\"></span>\", fmt.Sprintf(\"<span id=\\\"%s\\\">%%v</span>\", component.%s.Get()), 1)\n",
				expr.elementID, expr.elementID, expr.fieldName)
		}

		for _, ab := range childBindings.attrBindings {
			staticValue := ab.template
			for _, field := range ab.fields {
				staticValue = strings.ReplaceAll(staticValue, "{"+field+"}", "")
			}
			staticValue = strings.TrimSpace(staticValue)

			fmt.Fprintf(&sb, "\t\t{\n\t\t\tattrVal := %s\n", escapeForGoString(ab.template))
			for _, field := range ab.fields {
				fmt.Fprintf(&sb, "\t\t\tattrVal = strings.ReplaceAll(attrVal, \"{%s}\", fmt.Sprintf(\"%%v\", %s.%s.Get()))\n",
					field, compBinding.elementID, field)
			}
			fmt.Fprintf(&sb, "\t\t\tchildHTML = strings.Replace(childHTML, \"%s=\\\"%s\\\"\", \"%s=\\\"\" + attrVal + \"\\\"\", 1)\n\t\t}\n",
				ab.attrName, staticValue, ab.attrName)
		}

		// Replace the <!--compX--> placeholder with the actual child HTML
		fmt.Fprintf(&sb, "\t\thtml = strings.Replace(html, \"<!--%s-->\", childHTML, 1)\n\t}\n",
			compBinding.elementID)
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
