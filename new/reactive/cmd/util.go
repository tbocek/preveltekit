package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func findScriptDir() string {
	exe, err := os.Executable()
	if err != nil {
		fatal("find executable: %v", err)
	}
	dir := filepath.Dir(filepath.Dir(exe))

	// Check if we're running from go build cache (go run)
	// In that case, find the reactive package via go list
	if strings.Contains(dir, "go-build") || strings.Contains(dir, "cache") {
		cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "reactive")
		if out, err := cmd.Output(); err == nil {
			return strings.TrimSpace(string(out))
		}
		// Fallback: look for go.mod with "module reactive" in current directory hierarchy
		wd, _ := os.Getwd()
		for d := wd; d != "/" && d != "."; d = filepath.Dir(d) {
			modFile := filepath.Join(d, "go.mod")
			if data, err := os.ReadFile(modFile); err == nil {
				if strings.Contains(string(data), "module reactive") {
					return d
				}
			}
		}
	}

	return dir
}

func copyFile(src, dst, oldPkg, newPkg string) {
	data, err := os.ReadFile(src)
	if err != nil {
		fatal("read %s: %v", src, err)
	}
	content := string(data)
	if oldPkg != "" && newPkg != "" {
		content = strings.Replace(content, "package "+oldPkg, "package "+newPkg, 1)
	}
	writeFile(dst, content)
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fatal("write %s: %v", path, err)
	}
}

func copyWasmExec(dst string) {
	scriptDir := findScriptDir()
	src := filepath.Join(scriptDir, "wasm_exec.js")
	if _, err := os.Stat(src); err == nil {
		copyFile(src, filepath.Join(dst, "wasm_exec.js"), "", "")
		return
	}
	// Fallback to TinyGo
	cmd := exec.Command("tinygo", "env", "TINYGOROOT")
	if out, err := cmd.Output(); err == nil {
		src := filepath.Join(strings.TrimSpace(string(out)), "targets", "wasm_exec.js")
		if _, err := os.Stat(src); err == nil {
			copyFile(src, filepath.Join(dst, "wasm_exec.js"), "", "")
			return
		}
	}
	// Fallback to Go
	cmd = exec.Command("go", "env", "GOROOT")
	if out, err := cmd.Output(); err == nil {
		src := filepath.Join(strings.TrimSpace(string(out)), "misc", "wasm", "wasm_exec.js")
		if _, err := os.Stat(src); err == nil {
			copyFile(src, filepath.Join(dst, "wasm_exec.js"), "", "")
		}
	}
}

func extractStyle(html string) string {
	start := strings.Index(html, "<style>")
	end := strings.Index(html, "</style>")
	if start == -1 || end == -1 {
		return ""
	}
	return html[start+7 : end]
}

func stripTemplateAndStyle(src string) string {
	src = stripMethod(src, "Template")
	src = stripMethod(src, "Style")
	// Don't strip build tags - let Go's build system handle them
	return src
}

// stripMethod removes a method by name, properly handling nested braces and strings
func stripMethod(src, methodName string) string {
	// Find the method signature
	pattern := regexp.MustCompile(`(?m)^func \([^)]+\) ` + methodName + `\(\) string \{`)
	loc := pattern.FindStringIndex(src)
	if loc == nil {
		return src
	}

	start := loc[0]
	// Find matching closing brace by counting braces
	braceCount := 1
	inString := false
	inRawString := false
	i := loc[1] // start after opening brace

	for i < len(src) && braceCount > 0 {
		ch := src[i]

		if inRawString {
			if ch == '`' {
				inRawString = false
			}
			i++
			continue
		}

		if inString {
			if ch == '\\' && i+1 < len(src) {
				i += 2 // skip escaped char
				continue
			}
			if ch == '"' {
				inString = false
			}
			i++
			continue
		}

		switch ch {
		case '"':
			inString = true
		case '`':
			inRawString = true
		case '{':
			braceCount++
		case '}':
			braceCount--
		}
		i++
	}

	if braceCount == 0 {
		// Remove trailing newlines
		end := i
		for end < len(src) && (src[end] == '\n' || src[end] == '\r') {
			end++
		}
		return src[:start] + src[end:]
	}

	return src
}

func validateBindings(comp *component, bindings templateBindings) error {
	fieldNames := make(map[string]bool)
	for _, f := range comp.fields {
		fieldNames[f.name] = true
	}

	methodNames := make(map[string]bool)
	for _, m := range comp.methods {
		methodNames[m] = true
	}

	for _, expr := range bindings.expressions {
		if !fieldNames[expr.fieldName] {
			available := make([]string, 0, len(fieldNames))
			for name := range fieldNames {
				available = append(available, name)
			}
			return fmt.Errorf("template error: {%s} references unknown state\n\n  Available state: %v\n\n  Hint: Add '%s *reactive.Store[T]' to your component struct",
				expr.fieldName, available, expr.fieldName)
		}
	}

	for _, evt := range bindings.events {
		if !methodNames[evt.methodName] {
			available := make([]string, 0)
			for name := range methodNames {
				if name != "Template" && name != "Style" && name != "OnMount" {
					available = append(available, name+"()")
				}
			}
			return fmt.Errorf("template error: @%s=\"%s()\" references unknown method\n\n  Available methods: %v\n\n  Hint: Add 'func (c *%s) %s() { ... }' to your component",
				evt.event, evt.methodName, available, comp.name, evt.methodName)
		}
	}

	return nil
}

func fatal(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
	os.Exit(1)
}

// escapeForGoString escapes a string for use in a Go raw string literal (backtick string)
// Since backticks cannot be escaped inside raw strings, we use string concatenation
func escapeForGoString(s string) string {
	if !strings.Contains(s, "`") {
		return "`" + s + "`"
	}
	// Replace backticks with string concatenation: ` + "`" + `
	parts := strings.Split(s, "`")
	return "`" + strings.Join(parts, "` + \"`\" + `") + "`"
}

// escapeForGoStringContent escapes content for embedding inside a double-quoted string
func escapeForGoStringContent(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

// toJS returns Go code to convert a value to JS-compatible format
func toJS(valueType, expr string) string {
	switch valueType {
	case "string":
		return expr
	case "int", "int8", "int16", "int32", "int64":
		return "strconv.Itoa(" + expr + ")"
	case "float32", "float64":
		return "strconv.FormatFloat(float64(" + expr + "), 'f', -1, 64)"
	case "bool":
		return expr
	default:
		return expr
	}
}

func needsStrconvForType(t string) bool {
	switch t {
	case "int", "int8", "int16", "int32", "int64", "float32", "float64":
		return true
	}
	return false
}

func zeroValue(t string) string {
	switch t {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "0"
	case "float32", "float64":
		return "0.0"
	case "bool":
		return "false"
	case "string":
		return `""`
	default:
		return "nil"
	}
}

func generateHTML(prerenderedContent string, style string) string {
	return `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>App</title>
	<style>` + style + `</style>
	<script src="wasm_exec.js"></script>
	<script>
		const go = new Go();
		WebAssembly.instantiateStreaming(fetch("app.wasm"), go.importObject)
			.then(result => go.run(result.instance));
	</script>
</head>
<body>
	<div id="app">` + prerenderedContent + `</div>
</body>
</html>
`
}

// buildFieldTypes returns a map of field name to value type for a component.
// Uses cached value if available, otherwise builds and caches it.
func buildFieldTypes(comp *component) map[string]string {
	if comp.fieldTypes != nil {
		return comp.fieldTypes
	}
	comp.fieldTypes = make(map[string]string)
	for _, f := range comp.fields {
		comp.fieldTypes[f.name] = f.valueType
	}
	return comp.fieldTypes
}

// findReferencedComponents finds all component tags (PascalCase) in a template
// This includes components inside if-blocks, each-blocks, etc.
func findReferencedComponents(tmpl string) []string {
	return findComponentTags(tmpl)
}

// categorizeExpressions separates expressions into parent (slot) and child owned based on field types.
// It also sets the owner field on each expression.
func categorizeExpressions(exprs []exprBinding, slotFields map[string]bool, parentTypes, childTypes map[string]string) (parentExprs, childExprs []exprBinding) {
	for i := range exprs {
		expr := &exprs[i]
		if slotFields[expr.fieldName] && parentTypes[expr.fieldName] != "" {
			expr.owner = "parent"
			parentExprs = append(parentExprs, *expr)
		} else if childTypes[expr.fieldName] != "" {
			expr.owner = "child"
			childExprs = append(childExprs, *expr)
		} else if parentTypes[expr.fieldName] != "" {
			expr.owner = "parent"
			parentExprs = append(parentExprs, *expr)
		} else {
			expr.owner = "child"
			childExprs = append(childExprs, *expr)
		}
	}
	return
}

// injectIDIntoFirstTag adds id="compX" to the first HTML tag in a template
// e.g., `<button class="btn">` becomes `<button id="comp0" class="btn">`
func injectIDIntoFirstTag(tmpl, id string) string {
	// Find the first < that's not a comment or doctype
	for i := 0; i < len(tmpl); i++ {
		if tmpl[i] == '<' && i+1 < len(tmpl) {
			next := tmpl[i+1]
			// Skip comments, doctypes, closing tags
			if next == '!' || next == '?' || next == '/' {
				continue
			}
			// Found opening tag - find the end of tag name
			tagEnd := i + 1
			for tagEnd < len(tmpl) && tmpl[tagEnd] != ' ' && tmpl[tagEnd] != '>' && tmpl[tagEnd] != '/' {
				tagEnd++
			}
			// Insert id after tag name
			if tagEnd < len(tmpl) {
				return tmpl[:tagEnd] + fmt.Sprintf(` id="%s"`, id) + tmpl[tagEnd:]
			}
		}
	}
	return tmpl
}

// prefixBindingIDs prefixes all binding IDs with a component prefix and updates the HTML template.
// Returns the modified HTML with updated IDs.
func prefixBindingIDs(prefix string, html string, exprs []exprBinding, events []eventBinding, attrBindings []attrBinding, ifBlocks []ifBinding) string {
	// Prefix expression IDs
	for i := range exprs {
		oldID := exprs[i].elementID
		newID := prefix + "_" + oldID
		html = strings.ReplaceAll(html, `id="`+oldID+`"`, `id="`+newID+`"`)
		exprs[i].elementID = newID
	}
	// Prefix event IDs
	for i := range events {
		oldID := events[i].elementID
		newID := prefix + "_" + oldID
		html = strings.ReplaceAll(html, `id="`+oldID+`"`, `id="`+newID+`"`)
		events[i].elementID = newID
	}
	// Prefix attribute binding IDs
	for i := range attrBindings {
		oldID := attrBindings[i].elementID
		newID := prefix + "_" + oldID
		html = strings.ReplaceAll(html, `data-attrbind="`+oldID+`"`, `data-attrbind="`+newID+`"`)
		attrBindings[i].elementID = newID
	}
	// Prefix if-block anchor IDs
	for i := range ifBlocks {
		oldID := ifBlocks[i].elementID
		newID := prefix + "_" + oldID
		html = strings.ReplaceAll(html, `id="`+oldID+`_anchor"`, `id="`+newID+`_anchor"`)
		ifBlocks[i].elementID = newID
	}
	return html
}

// prefixInputBindingIDs prefixes input binding IDs (bind:value, bind:checked) in HTML.
func prefixInputBindingIDs(prefix string, html string, bindings []inputBinding) string {
	for i := range bindings {
		oldID := bindings[i].elementID
		newID := prefix + "_" + oldID
		html = strings.ReplaceAll(html, `id="`+oldID+`"`, `id="`+newID+`"`)
		bindings[i].elementID = newID
	}
	return html
}

// prefixEachBindingIDs prefixes each-block anchor IDs in HTML.
func prefixEachBindingIDs(prefix string, html string, bindings []eachBinding) string {
	for i := range bindings {
		oldID := bindings[i].elementID
		newID := prefix + "_" + oldID
		html = strings.ReplaceAll(html, `id="`+oldID+`_anchor"`, `id="`+newID+`_anchor"`)
		bindings[i].elementID = newID
	}
	return html
}

// prefixClassBindingIDs prefixes class binding IDs in HTML.
func prefixClassBindingIDs(prefix string, html string, bindings []classBinding) string {
	for i := range bindings {
		oldID := bindings[i].elementID
		newID := prefix + "_" + oldID
		html = strings.ReplaceAll(html, `id="`+oldID+`"`, `id="`+newID+`"`)
		bindings[i].elementID = newID
	}
	return html
}

// generateFieldInit generates the initialization code for component fields
// indent is the base indentation (e.g., "\t" or "\t\t")
func generateFieldInit(sb *strings.Builder, fields []storeField, indent string) {
	for _, field := range fields {
		switch field.storeType {
		case "Store":
			fmt.Fprintf(sb, "%s%s: reactive.New[%s](%s),\n", indent, field.name, field.valueType, zeroValue(field.valueType))
		case "LocalStore":
			fmt.Fprintf(sb, "%s%s: reactive.NewLocalStore(\"%s\", %s),\n", indent, field.name, field.name, zeroValue(field.valueType))
		case "List":
			fmt.Fprintf(sb, "%s%s: reactive.NewList[%s](),\n", indent, field.name, field.valueType)
		case "Map":
			fmt.Fprintf(sb, "%s%s: reactive.NewMap[%s, %s](),\n", indent, field.name, field.keyType, field.valueType)
		}
	}
}

// transformCondition transforms a template condition to valid Go code
// It handles simple field references (like "Error") by making them truthy checks
// based on the field type (e.g., string != "" for strings, bool as-is for bools)
func transformCondition(cond string, fieldTypes map[string]string, prefix string) string {
	// First check if it's a simple field reference (just a field name)
	trimmed := strings.TrimSpace(cond)
	if fieldType, ok := fieldTypes[trimmed]; ok {
		// It's a simple field reference - make it a truthy check
		switch fieldType {
		case "string":
			return prefix + "." + trimmed + ".Get() != \"\""
		case "bool":
			return prefix + "." + trimmed + ".Get()"
		case "int", "int8", "int16", "int32", "int64":
			return prefix + "." + trimmed + ".Get() != 0"
		case "float32", "float64":
			return prefix + "." + trimmed + ".Get() != 0"
		default:
			// For other types, just call .Get() (might need type-specific handling)
			return prefix + "." + trimmed + ".Get()"
		}
	}

	// Not a simple field reference - do word-boundary replacement
	// Sort field names by length (longest first) to avoid partial matches
	fieldNames := make([]string, 0, len(fieldTypes))
	for name := range fieldTypes {
		fieldNames = append(fieldNames, name)
	}
	// Sort by length descending
	for i := 0; i < len(fieldNames)-1; i++ {
		for j := i + 1; j < len(fieldNames); j++ {
			if len(fieldNames[j]) > len(fieldNames[i]) {
				fieldNames[i], fieldNames[j] = fieldNames[j], fieldNames[i]
			}
		}
	}

	for _, fieldName := range fieldNames {
		// Use word boundary matching - field must not be preceded or followed by alphanumeric
		result := ""
		i := 0
		for i < len(cond) {
			idx := strings.Index(cond[i:], fieldName)
			if idx == -1 {
				result += cond[i:]
				break
			}
			pos := i + idx
			// Check word boundaries
			beforeOk := pos == 0 || !isAlphanumeric(cond[pos-1])
			afterOk := pos+len(fieldName) >= len(cond) || !isAlphanumeric(cond[pos+len(fieldName)])
			if beforeOk && afterOk {
				result += cond[i:pos] + prefix + "." + fieldName + ".Get()"
				i = pos + len(fieldName)
			} else {
				result += cond[i : pos+len(fieldName)]
				i = pos + len(fieldName)
			}
		}
		cond = result
	}
	return cond
}

// findCompPlaceholders finds all <!--compN--> placeholders in HTML and returns the comp IDs
func findCompPlaceholders(html string) []string {
	var result []string
	marker := "<!--comp"
	pos := 0
	for {
		idx := strings.Index(html[pos:], marker)
		if idx == -1 {
			break
		}
		start := pos + idx + len(marker)
		// Find the end -->
		end := strings.Index(html[start:], "-->")
		if end == -1 {
			break
		}
		// Extract the number part
		numStr := html[start : start+end]
		// Validate it's a number
		isNum := true
		for _, c := range numStr {
			if c < '0' || c > '9' {
				isNum = false
				break
			}
		}
		if isNum && len(numStr) > 0 {
			result = append(result, "comp"+numStr)
		}
		pos = start + end + 3
	}
	return result
}

// hasCompPlaceholder checks if HTML contains any <!--compN--> placeholder
func hasCompPlaceholder(html string) bool {
	return strings.Contains(html, "<!--comp")
}

// extractPascalCaseWords extracts all PascalCase words from a string
// Used for finding field dependencies in conditions
func extractPascalCaseWords(s string) []string {
	var result []string
	i := 0
	for i < len(s) {
		// Skip non-letters
		for i < len(s) && !isLetter(s[i]) {
			i++
		}
		if i >= len(s) {
			break
		}
		// Check if uppercase (start of PascalCase)
		if s[i] >= 'A' && s[i] <= 'Z' {
			start := i
			i++
			// Continue while alphanumeric
			for i < len(s) && isAlphanumeric(s[i]) {
				i++
			}
			result = append(result, s[start:i])
		} else {
			// Skip lowercase word
			for i < len(s) && isAlphanumeric(s[i]) {
				i++
			}
		}
	}
	return result
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isAlphanumeric(c byte) bool {
	return isLetter(c) || (c >= '0' && c <= '9')
}

// findFieldExprs finds all {FieldName} expressions in a string
// Returns slice of (fieldName, startIndex, endIndex)
func findFieldExprs(s string) []struct {
	field string
	start int
	end   int
} {
	var result []struct {
		field string
		start int
		end   int
	}
	i := 0
	for i < len(s) {
		if s[i] != '{' {
			i++
			continue
		}
		start := i
		i++
		if i >= len(s) {
			break
		}
		// Must start with uppercase
		if s[i] < 'A' || s[i] > 'Z' {
			continue
		}
		fieldStart := i
		// Read field name (alphanumeric)
		for i < len(s) && isAlphanumeric(s[i]) {
			i++
		}
		if i >= len(s) || s[i] != '}' {
			continue
		}
		fieldName := s[fieldStart:i]
		i++ // skip }
		result = append(result, struct {
			field string
			start int
			end   int
		}{fieldName, start, i})
	}
	return result
}

// extractFieldNames extracts field names from {Field} expressions in a string
func extractFieldNames(s string) []string {
	exprs := findFieldExprs(s)
	result := make([]string, len(exprs))
	for i, e := range exprs {
		result[i] = e.field
	}
	return result
}

// removeFieldExprs removes all {Field} expressions from a string
func removeFieldExprs(s string) string {
	exprs := findFieldExprs(s)
	if len(exprs) == 0 {
		return s
	}
	// Process in reverse to preserve indices
	result := s
	for i := len(exprs) - 1; i >= 0; i-- {
		result = result[:exprs[i].start] + result[exprs[i].end:]
	}
	return result
}

// findComponentTags finds all PascalCase component tag names in a template
func findComponentTags(tmpl string) []string {
	var result []string
	seen := make(map[string]bool)
	i := 0
	for i < len(tmpl) {
		// Find <
		if tmpl[i] != '<' {
			i++
			continue
		}
		i++
		if i >= len(tmpl) {
			break
		}
		// Check if uppercase (component tag)
		if tmpl[i] >= 'A' && tmpl[i] <= 'Z' {
			start := i
			// Read tag name
			for i < len(tmpl) && isAlphanumeric(tmpl[i]) {
				i++
			}
			name := tmpl[start:i]
			if !seen[name] {
				seen[name] = true
				result = append(result, name)
			}
		}
	}
	return result
}
