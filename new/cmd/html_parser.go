package main

import (
	"fmt"
	"regexp"
	"strings"
)

// skipTags are tags where we don't parse bindings (code examples, etc.)
var skipTagsSet = map[string]bool{
	"pre":    true,
	"code":   true,
	"script": true,
	"style":  true,
}

// parseTemplate parses a template by walking through it sequentially
// It properly handles skip tags and doesn't use regex on regions that should be skipped
func parseTemplate(tmpl string) (string, templateBindings) {
	bindings := templateBindings{}
	var result strings.Builder

	pos := 0
	exprCount := 0
	evtCount := 0
	bindCount := 0
	classCount := 0
	attrCount := 0
	eachCount := 0
	ifCount := 0
	compCount := 0

	for pos < len(tmpl) {
		// Check for <pre><code>...</code></pre> or <pre>...</pre> - escape inner content
		if strings.HasPrefix(strings.ToLower(tmpl[pos:]), "<pre") {
			if escaped, endPos, ok := escapePreBlock(tmpl, pos); ok {
				result.WriteString(escaped)
				pos = endPos
				continue
			}
		}

		// Check for standalone <code>...</code>
		if strings.HasPrefix(strings.ToLower(tmpl[pos:]), "<code") {
			if escaped, endPos, ok := escapeCodeBlock(tmpl, pos); ok {
				result.WriteString(escaped)
				pos = endPos
				continue
			}
		}

		// Check for script/style - copy as-is
		if _, skipEnd, tagName, found := findSkipTagWithName(tmpl, pos); found {
			if tagName == "script" || tagName == "style" {
				result.WriteString(tmpl[pos:skipEnd])
				pos = skipEnd
				continue
			}
		}

		// Check for {#each ...}...{/each}
		if strings.HasPrefix(tmpl[pos:], "{#each ") {
			endPos, each := parseEachBlock(tmpl, pos, &eachCount)
			if each != nil {
				bindings.eachBlocks = append(bindings.eachBlocks, *each)
				if each.elseHTML != "" {
					fmt.Fprintf(&result, `<span id="%s_else">%s</span>`, each.elementID, each.elseHTML)
				}
				fmt.Fprintf(&result, `<!--e%s-->`, strings.TrimPrefix(each.elementID, "each"))
				pos = endPos
				continue
			}
		}

		// Check for {#if ...}
		if strings.HasPrefix(tmpl[pos:], "{#if ") {
			endPos, ifBlock := parseIfBlock(tmpl, pos, &ifCount)
			if ifBlock != nil {
				// Parse components, each blocks, class bindings, and expressions inside each branch
				for i := range ifBlock.branches {
					parsedHTML, comps := parseComponentsInHTML(ifBlock.branches[i].html, &compCount)
					parsedHTML, eachBlocksInBranch := parseEachBlocksInHTML(parsedHTML, &eachCount)
					parsedHTML, classBindingsInBranch := parseClassBindingsInHTML(parsedHTML, &classCount)
					parsedHTML, expressionsInBranch := parseExpressionsInHTML(parsedHTML, &exprCount)
					ifBlock.branches[i].html = parsedHTML
					ifBlock.branches[i].eachBlocks = eachBlocksInBranch
					ifBlock.branches[i].classBindings = classBindingsInBranch
					ifBlock.branches[i].expressions = expressionsInBranch
					bindings.components = append(bindings.components, comps...)
				}
				if ifBlock.elseHTML != "" {
					parsedHTML, comps := parseComponentsInHTML(ifBlock.elseHTML, &compCount)
					parsedHTML, _ = parseClassBindingsInHTML(parsedHTML, &classCount)
					parsedHTML, elseExprs := parseExpressionsInHTML(parsedHTML, &exprCount)
					ifBlock.elseHTML = parsedHTML
					ifBlock.elseExpressions = elseExprs
					bindings.components = append(bindings.components, comps...)
				}
				bindings.ifBlocks = append(bindings.ifBlocks, *ifBlock)
				fmt.Fprintf(&result, `<!--i%s-->`, strings.TrimPrefix(ifBlock.elementID, "if"))
				pos = endPos
				continue
			}
		}

		// Check for {@html Field}
		if strings.HasPrefix(tmpl[pos:], "{@html ") {
			if endPos := strings.Index(tmpl[pos:], "}"); endPos != -1 {
				fieldName := strings.TrimSpace(tmpl[pos+7 : pos+endPos])
				elementID := fmt.Sprintf("t%d", exprCount)
				exprCount++
				bindings.expressions = append(bindings.expressions, exprBinding{
					fieldName: fieldName, elementID: elementID, isHTML: true,
				})
				fmt.Fprintf(&result, `<!--%s-->`, elementID)
				pos = pos + endPos + 1
				continue
			}
		}

		// Check for {Field} expressions (but not {#, {:, {/, {@)
		if tmpl[pos] == '{' && pos+1 < len(tmpl) {
			nextChar := tmpl[pos+1]
			if nextChar != '#' && nextChar != ':' && nextChar != '/' && nextChar != '@' {
				if endPos := strings.Index(tmpl[pos:], "}"); endPos != -1 {
					fieldName := strings.TrimSpace(tmpl[pos+1 : pos+endPos])
					if isValidFieldName(fieldName) {
						elementID := fmt.Sprintf("t%d", exprCount)
						exprCount++
						bindings.expressions = append(bindings.expressions, exprBinding{
							fieldName: fieldName, elementID: elementID, isHTML: false,
						})
						fmt.Fprintf(&result, `<!--%s-->`, elementID)
						pos = pos + endPos + 1
						continue
					}
				}
			}
		}

		// Check for HTML tags
		if tmpl[pos] == '<' {
			// Check for PascalCase component tags (e.g., <Button, <Card)
			if pos+1 < len(tmpl) && tmpl[pos+1] >= 'A' && tmpl[pos+1] <= 'Z' {
				if compEnd, comp := parseComponentTag(tmpl, pos, &compCount); comp != nil {
					bindings.components = append(bindings.components, *comp)
					fmt.Fprintf(&result, `<!--c%s-->`, strings.TrimPrefix(comp.elementID, "comp"))
					pos = compEnd
					continue
				}
			}

			// Parse regular HTML tag for attributes
			if tagEnd, tagContent := parseHTMLTag(tmpl, pos); tagEnd > pos {
				processedTag := processTagAttributes(tagContent, &bindings, &evtCount, &bindCount, &classCount, &attrCount)
				result.WriteString(processedTag)
				pos = tagEnd
				continue
			}
		}

		// Default: copy character as-is
		result.WriteByte(tmpl[pos])
		pos++
	}

	return result.String(), bindings
}

// findSkipTag checks if we're at a skip tag and returns its boundaries
func findSkipTag(tmpl string, pos int) (start, end int, found bool) {
	start, end, _, found = findSkipTagWithName(tmpl, pos)
	return
}

// findSkipTagWithName checks if we're at a skip tag and returns its boundaries and tag name
func findSkipTagWithName(tmpl string, pos int) (start, end int, tagName string, found bool) {
	if tmpl[pos] != '<' {
		return 0, 0, "", false
	}

	for tag := range skipTagsSet {
		openTag := "<" + tag
		if strings.HasPrefix(strings.ToLower(tmpl[pos:]), openTag) {
			// Check it's actually a tag (followed by > or space)
			afterTag := pos + len(openTag)
			if afterTag < len(tmpl) && (tmpl[afterTag] == '>' || tmpl[afterTag] == ' ' || tmpl[afterTag] == '\t' || tmpl[afterTag] == '\n') {
				// Find closing tag
				closeTag := "</" + tag + ">"
				closeIdx := strings.Index(strings.ToLower(tmpl[afterTag:]), closeTag)
				if closeIdx != -1 {
					return pos, afterTag + closeIdx + len(closeTag), tag, true
				}
			}
		}
	}
	return 0, 0, "", false
}

// htmlEscapeContent escapes < > & for display in HTML
func htmlEscapeContent(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// escapePreBlock handles <pre>...</pre> and <pre><code>...</code></pre>
func escapePreBlock(tmpl string, pos int) (escaped string, endPos int, ok bool) {
	// Find end of <pre> opening tag
	preOpenEnd := strings.Index(tmpl[pos:], ">")
	if preOpenEnd == -1 {
		return "", pos, false
	}
	preOpenEnd += pos + 1

	// Find </pre>
	preCloseStart := strings.Index(strings.ToLower(tmpl[preOpenEnd:]), "</pre>")
	if preCloseStart == -1 {
		return "", pos, false
	}
	preCloseStart += preOpenEnd

	preContent := tmpl[preOpenEnd:preCloseStart]

	// Check if content starts with <code>
	trimmedContent := strings.TrimSpace(preContent)
	if strings.HasPrefix(strings.ToLower(trimmedContent), "<code") {
		// Find <code> boundaries within preContent
		codeStart := strings.Index(strings.ToLower(preContent), "<code")
		codeOpenEnd := strings.Index(preContent[codeStart:], ">")
		if codeOpenEnd == -1 {
			return "", pos, false
		}
		codeOpenEnd += codeStart + 1

		codeCloseStart := strings.Index(strings.ToLower(preContent[codeOpenEnd:]), "</code>")
		if codeCloseStart == -1 {
			return "", pos, false
		}
		codeCloseStart += codeOpenEnd

		// Build: <pre><code>ESCAPED</code></pre>
		var result strings.Builder
		result.WriteString(tmpl[pos:preOpenEnd])                                      // <pre>
		result.WriteString(preContent[:codeOpenEnd])                                  // <code>
		result.WriteString(htmlEscapeContent(preContent[codeOpenEnd:codeCloseStart])) // escaped content
		result.WriteString(preContent[codeCloseStart:])                               // </code>
		result.WriteString("</pre>")

		return result.String(), preCloseStart + 6, true
	}

	// Just <pre>content</pre> without <code>
	var result strings.Builder
	result.WriteString(tmpl[pos:preOpenEnd])
	result.WriteString(htmlEscapeContent(preContent))
	result.WriteString("</pre>")

	return result.String(), preCloseStart + 6, true
}

// escapeCodeBlock handles standalone <code>...</code>
func escapeCodeBlock(tmpl string, pos int) (escaped string, endPos int, ok bool) {
	// Find end of <code> opening tag
	codeOpenEnd := strings.Index(tmpl[pos:], ">")
	if codeOpenEnd == -1 {
		return "", pos, false
	}
	codeOpenEnd += pos + 1

	// Find </code>
	codeCloseStart := strings.Index(strings.ToLower(tmpl[codeOpenEnd:]), "</code>")
	if codeCloseStart == -1 {
		return "", pos, false
	}
	codeCloseStart += codeOpenEnd

	content := tmpl[codeOpenEnd:codeCloseStart]

	var result strings.Builder
	result.WriteString(tmpl[pos:codeOpenEnd])
	result.WriteString(htmlEscapeContent(content))
	result.WriteString("</code>")

	return result.String(), codeCloseStart + 7, true
}

// parseEachBlock parses {#each list as item, i}...{:else}...{/each}
func parseEachBlock(tmpl string, pos int, count *int) (endPos int, binding *eachBinding) {
	// Find the end of opening tag
	openEnd := strings.Index(tmpl[pos:], "}")
	if openEnd == -1 {
		return pos, nil
	}
	openEnd += pos

	// Parse the opening: {#each ListName as itemVar, indexVar}
	openTag := tmpl[pos+7 : openEnd] // skip "{#each "
	parts := strings.Split(openTag, " as ")
	if len(parts) != 2 {
		return pos, nil
	}

	listName := strings.TrimSpace(parts[0])
	varParts := strings.Split(parts[1], ",")
	itemVar := strings.TrimSpace(varParts[0])
	indexVar := "_i"
	if len(varParts) > 1 {
		indexVar = strings.TrimSpace(varParts[1])
	}

	// Find matching {/each}, tracking {:else} at depth 1
	depth := 1
	searchPos := openEnd + 1
	elsePos := -1
	for searchPos < len(tmpl) && depth > 0 {
		if strings.HasPrefix(tmpl[searchPos:], "{#each ") {
			depth++
			searchPos += 7
		} else if strings.HasPrefix(tmpl[searchPos:], "{:else}") && depth == 1 {
			elsePos = searchPos
			searchPos += 7
		} else if strings.HasPrefix(tmpl[searchPos:], "{/each}") {
			depth--
			if depth == 0 {
				break
			}
			searchPos += 7
		} else {
			searchPos++
		}
	}

	if depth != 0 {
		return pos, nil
	}

	var bodyHTML, elseHTML string
	if elsePos != -1 {
		bodyHTML = tmpl[openEnd+1 : elsePos]
		elseHTML = tmpl[elsePos+7 : searchPos]
	} else {
		bodyHTML = tmpl[openEnd+1 : searchPos]
	}

	elementID := fmt.Sprintf("each%d", *count)
	*count++

	return searchPos + 7, &eachBinding{
		listName:  listName,
		itemVar:   itemVar,
		indexVar:  indexVar,
		elementID: elementID,
		bodyHTML:  bodyHTML,
		elseHTML:  elseHTML,
	}
}

// parseIfBlock parses {#if cond}...{:else if cond}...{:else}...{/if}
func parseIfBlock(tmpl string, pos int, count *int) (endPos int, binding *ifBinding) {
	// Find end of opening tag
	openEnd := strings.Index(tmpl[pos:], "}")
	if openEnd == -1 {
		return pos, nil
	}
	openEnd += pos

	firstCondition := strings.TrimSpace(tmpl[pos+5 : openEnd]) // skip "{#if "

	var branches []ifBranch
	var elseHTML string
	currentCondition := firstCondition
	contentStart := openEnd + 1

	depth := 1
	searchPos := openEnd + 1

	for searchPos < len(tmpl) && depth > 0 {
		if strings.HasPrefix(tmpl[searchPos:], "{#if ") {
			depth++
			searchPos += 5
		} else if strings.HasPrefix(tmpl[searchPos:], "{:else if ") && depth == 1 {
			// Save current branch
			branches = append(branches, ifBranch{
				condition: currentCondition,
				html:      tmpl[contentStart:searchPos],
			})
			// Parse new condition
			condEnd := strings.Index(tmpl[searchPos:], "}")
			if condEnd == -1 {
				return pos, nil
			}
			currentCondition = strings.TrimSpace(tmpl[searchPos+10 : searchPos+condEnd])
			searchPos = searchPos + condEnd + 1
			contentStart = searchPos
		} else if strings.HasPrefix(tmpl[searchPos:], "{:else}") && depth == 1 {
			// Save current branch
			branches = append(branches, ifBranch{
				condition: currentCondition,
				html:      tmpl[contentStart:searchPos],
			})
			currentCondition = ""
			searchPos += 7
			contentStart = searchPos
		} else if strings.HasPrefix(tmpl[searchPos:], "{/if}") {
			depth--
			if depth == 0 {
				// Save final content
				content := tmpl[contentStart:searchPos]
				if currentCondition != "" {
					branches = append(branches, ifBranch{
						condition: currentCondition,
						html:      content,
					})
				} else {
					elseHTML = content
				}
				break
			}
			searchPos += 5
		} else {
			searchPos++
		}
	}

	if depth != 0 {
		return pos, nil
	}

	// Extract dependencies from conditions
	depSet := make(map[string]bool)
	for _, branch := range branches {
		for _, d := range extractPascalCaseWords(branch.condition) {
			depSet[d] = true
		}
	}
	var deps []string
	for dep := range depSet {
		deps = append(deps, dep)
	}

	elementID := fmt.Sprintf("if%d", *count)
	*count++

	return searchPos + 5, &ifBinding{
		branches:  branches,
		elseHTML:  elseHTML,
		elementID: elementID,
		deps:      deps,
	}
}

// parseComponentsInHTML parses component tags in an HTML string (used for if-block content)
// Also escapes content inside <pre><code> blocks
func parseComponentsInHTML(html string, compCount *int) (string, []componentBinding) {
	var components []componentBinding
	var result strings.Builder
	pos := 0

	for pos < len(html) {
		// Check for <pre> blocks - escape their content
		if strings.HasPrefix(strings.ToLower(html[pos:]), "<pre") {
			if escaped, endPos, ok := escapePreBlock(html, pos); ok {
				result.WriteString(escaped)
				pos = endPos
				continue
			}
		}

		// Check for standalone <code> blocks
		if strings.HasPrefix(strings.ToLower(html[pos:]), "<code") {
			if escaped, endPos, ok := escapeCodeBlock(html, pos); ok {
				result.WriteString(escaped)
				pos = endPos
				continue
			}
		}

		// Check for PascalCase component tags
		if html[pos] == '<' && pos+1 < len(html) && html[pos+1] >= 'A' && html[pos+1] <= 'Z' {
			if compEnd, comp := parseComponentTag(html, pos, compCount); comp != nil {
				components = append(components, *comp)
				fmt.Fprintf(&result, `<!--c%s-->`, strings.TrimPrefix(comp.elementID, "comp"))
				pos = compEnd
				continue
			}
		}

		result.WriteByte(html[pos])
		pos++
	}

	return result.String(), components
}

// parseEachBlocksInHTML parses {#each} blocks in an HTML string (used for if-block content)
// Returns the modified HTML with anchors and the list of each bindings
func parseEachBlocksInHTML(html string, eachCount *int) (string, []eachBinding) {
	var eachBlocks []eachBinding
	var result strings.Builder
	pos := 0

	for pos < len(html) {
		// Check for {#each ...}...{/each}
		if strings.HasPrefix(html[pos:], "{#each ") {
			endPos, each := parseEachBlock(html, pos, eachCount)
			if each != nil {
				eachBlocks = append(eachBlocks, *each)
				if each.elseHTML != "" {
					fmt.Fprintf(&result, `<span id="%s_else">%s</span>`, each.elementID, each.elseHTML)
				}
				fmt.Fprintf(&result, `<!--e%s-->`, strings.TrimPrefix(each.elementID, "each"))
				pos = endPos
				continue
			}
		}

		result.WriteByte(html[pos])
		pos++
	}

	return result.String(), eachBlocks
}

// parseClassBindingsInHTML parses class:name={condition} bindings in HTML content (used for if-block branches)
// Returns the modified HTML with IDs and the list of class bindings
// Handles multiple class bindings on the same element by sharing the same ID
func parseClassBindingsInHTML(html string, classCount *int) (string, []classBinding) {
	var classBindings []classBinding
	var result strings.Builder
	pos := 0

	for pos < len(html) {
		// Look for opening tags
		if html[pos] == '<' && pos+1 < len(html) && html[pos+1] != '/' && html[pos+1] != '!' {
			// Find end of tag
			tagEnd := pos + 1
			for tagEnd < len(html) && html[tagEnd] != '>' {
				tagEnd++
			}
			if tagEnd >= len(html) {
				result.WriteByte(html[pos])
				pos++
				continue
			}

			tag := html[pos : tagEnd+1]

			// Check if this tag has any class bindings
			matches := classBindRegex.FindAllStringSubmatch(tag, -1)
			if len(matches) > 0 {
				elementID := fmt.Sprintf("class%d", *classCount)
				*classCount++

				// Collect all class bindings for this element
				for _, parts := range matches {
					classBindings = append(classBindings, classBinding{
						className: parts[1],
						condition: parts[2],
						elementID: elementID,
					})
				}

				// Remove all class:name={...} from the tag and add a single id
				processedTag := classBindRegex.ReplaceAllString(tag, "")
				// Insert id before the closing >
				if strings.HasSuffix(processedTag, "/>") {
					processedTag = processedTag[:len(processedTag)-2] + fmt.Sprintf(` id="%s" />`, elementID)
				} else {
					processedTag = processedTag[:len(processedTag)-1] + fmt.Sprintf(` id="%s">`, elementID)
				}
				result.WriteString(processedTag)
				pos = tagEnd + 1
				continue
			}

			result.WriteString(tag)
			pos = tagEnd + 1
			continue
		}

		result.WriteByte(html[pos])
		pos++
	}

	return result.String(), classBindings
}

// parseExpressionsInHTML parses {Field} expressions in HTML content (used for if-block branches)
// Returns the modified HTML with comment markers and the list of expression bindings
func parseExpressionsInHTML(html string, exprCount *int) (string, []exprBinding) {
	var expressions []exprBinding
	var result strings.Builder
	pos := 0

	for pos < len(html) {
		// Check for {@html Field}
		if strings.HasPrefix(html[pos:], "{@html ") {
			if endPos := strings.Index(html[pos:], "}"); endPos != -1 {
				fieldName := strings.TrimSpace(html[pos+7 : pos+endPos])
				elementID := fmt.Sprintf("t%d", *exprCount)
				*exprCount++
				expressions = append(expressions, exprBinding{
					fieldName: fieldName, elementID: elementID, isHTML: true,
				})
				fmt.Fprintf(&result, `<!--%s-->`, elementID)
				pos += endPos + 1
				continue
			}
		}

		// Check for {Field} expressions (but not {#, {:, {/, {@)
		if html[pos] == '{' && pos+1 < len(html) {
			nextChar := html[pos+1]
			if nextChar != '#' && nextChar != ':' && nextChar != '/' && nextChar != '@' {
				if endPos := strings.Index(html[pos:], "}"); endPos != -1 {
					fieldName := strings.TrimSpace(html[pos+1 : pos+endPos])
					if isValidFieldName(fieldName) {
						elementID := fmt.Sprintf("t%d", *exprCount)
						*exprCount++
						expressions = append(expressions, exprBinding{
							fieldName: fieldName, elementID: elementID, isHTML: false,
						})
						fmt.Fprintf(&result, `<!--%s-->`, elementID)
						pos += endPos + 1
						continue
					}
				}
			}
		}

		result.WriteByte(html[pos])
		pos++
	}

	return result.String(), expressions
}

// parseComponentTag parses a PascalCase component tag
func parseComponentTag(tmpl string, pos int, count *int) (endPos int, binding *componentBinding) {
	// Match component name
	nameEnd := pos + 1
	for nameEnd < len(tmpl) && (isAlphaNum(tmpl[nameEnd]) || tmpl[nameEnd] == '_') {
		nameEnd++
	}
	compName := tmpl[pos+1 : nameEnd]

	// Skip whitespace
	attrStart := nameEnd
	for attrStart < len(tmpl) && (tmpl[attrStart] == ' ' || tmpl[attrStart] == '\t' || tmpl[attrStart] == '\n') {
		attrStart++
	}

	// Check for self-closing /> or >
	if attrStart >= len(tmpl) {
		return pos, nil
	}

	// Find end of opening tag
	tagEnd := attrStart
	for tagEnd < len(tmpl) && tmpl[tagEnd] != '>' {
		tagEnd++
	}
	if tagEnd >= len(tmpl) {
		return pos, nil
	}

	// Check if self-closing
	isSelfClosing := tagEnd > 0 && tmpl[tagEnd-1] == '/'
	attrs := ""
	if isSelfClosing {
		attrs = strings.TrimSpace(tmpl[attrStart : tagEnd-1])
	} else {
		attrs = strings.TrimSpace(tmpl[attrStart:tagEnd])
	}

	props, events := parseComponentAttrs(attrs)
	elementID := fmt.Sprintf("comp%d", *count)
	*count++

	if isSelfClosing {
		return tagEnd + 1, &componentBinding{
			name:      compName,
			elementID: elementID,
			props:     props,
			events:    events,
			children:  "",
		}
	}

	// Find closing tag </CompName>
	closeTag := "</" + compName + ">"
	closeIdx := strings.Index(tmpl[tagEnd+1:], closeTag)
	if closeIdx == -1 {
		return pos, nil
	}
	closeIdx += tagEnd + 1

	children := strings.TrimSpace(tmpl[tagEnd+1 : closeIdx])

	return closeIdx + len(closeTag), &componentBinding{
		name:      compName,
		elementID: elementID,
		props:     props,
		events:    events,
		children:  children,
	}
}

// parseHTMLTag extracts a complete HTML tag
// Handles > characters inside {...} braces (e.g., class:active={CurrentStep > 1})
func parseHTMLTag(tmpl string, pos int) (endPos int, tagContent string) {
	if tmpl[pos] != '<' {
		return pos, ""
	}

	// Find the end of the tag, but skip > inside {...}
	end := pos + 1
	braceDepth := 0
	for end < len(tmpl) {
		ch := tmpl[end]
		if ch == '{' {
			braceDepth++
		} else if ch == '}' {
			braceDepth--
		} else if ch == '>' && braceDepth == 0 {
			break
		}
		end++
	}
	if end >= len(tmpl) {
		return pos, ""
	}

	return end + 1, tmpl[pos : end+1]
}

// processTagAttributes processes @click, bind:value, class:name, and attribute bindings
func processTagAttributes(tag string, bindings *templateBindings, evtCount, bindCount, classCount, attrCount *int) string {
	result := tag

	// Check what bindings exist on this tag
	eventMatches := eventRegex.FindAllStringSubmatch(result, -1)
	bindMatches := bindRegex.FindAllStringSubmatch(result, -1)
	classMatches := classBindRegex.FindAllStringSubmatch(result, -1)

	// Determine if we need an ID and what it should be
	var elementID string
	needsID := len(eventMatches) > 0 || len(bindMatches) > 0 || len(classMatches) > 0

	if needsID {
		// Check if tag already has an id attribute
		idRegex := regexp.MustCompile(`\bid="([^"]+)"`)
		if match := idRegex.FindStringSubmatch(result); match != nil {
			elementID = match[1]
		} else {
			// Generate a new ID based on the first binding type found
			if len(eventMatches) > 0 {
				elementID = fmt.Sprintf("evt_%s_%d", eventMatches[0][1], *evtCount)
				*evtCount++
			} else if len(bindMatches) > 0 {
				elementID = fmt.Sprintf("bind%d", *bindCount)
				*bindCount++
			} else if len(classMatches) > 0 {
				elementID = fmt.Sprintf("class%d", *classCount)
				*classCount++
			}
		}
	}

	// Process @event="Method()" - remove from tag, add to bindings
	for _, parts := range eventMatches {
		event, modifiersStr, methodName, args := parts[1], parts[2], parts[3], parts[4]

		var modifiers []string
		for _, mod := range strings.Split(modifiersStr, ".") {
			if mod != "" {
				modifiers = append(modifiers, mod)
			}
		}

		bindings.events = append(bindings.events, eventBinding{
			event: event, modifiers: modifiers, methodName: methodName,
			args: args, elementID: elementID,
		})
	}
	result = eventRegex.ReplaceAllString(result, "")

	// Process bind:value="Field" and bind:checked="Field"
	for _, parts := range bindMatches {
		bindings.bindings = append(bindings.bindings, inputBinding{
			fieldName: parts[2], bindType: parts[1], elementID: elementID,
		})
	}
	result = bindRegex.ReplaceAllString(result, "")

	// Process class:name={Condition}
	for _, parts := range classMatches {
		bindings.classBindings = append(bindings.classBindings, classBinding{
			className: parts[1], condition: parts[2], elementID: elementID,
		})
	}
	result = classBindRegex.ReplaceAllString(result, "")

	// Add the ID to the tag if needed and not already present
	if needsID && !strings.Contains(tag, `id="`) {
		// Insert id before the closing >
		if strings.HasSuffix(result, "/>") {
			result = result[:len(result)-2] + fmt.Sprintf(` id="%s" />`, elementID)
		} else if strings.HasSuffix(result, ">") {
			result = result[:len(result)-1] + fmt.Sprintf(` id="%s">`, elementID)
		}
	}

	// Process attribute bindings like href="{Field}" or src="{Base}/{Path}"
	result = attrWithExprRegex.ReplaceAllStringFunc(result, func(match string) string {
		parts := attrWithExprRegex.FindStringSubmatch(match)
		attrName, attrValue := parts[1], parts[2]

		if attrName == "bind" || strings.HasPrefix(attrName, "@") {
			return match
		}

		fields := extractFieldNames(attrValue)

		if len(fields) == 0 {
			return match
		}

		elementID := fmt.Sprintf("attr%d", *attrCount)
		*attrCount++

		bindings.attrBindings = append(bindings.attrBindings, attrBinding{
			attrName: attrName, template: attrValue, fields: fields, elementID: elementID,
		})

		staticValue := strings.TrimSpace(removeFieldExprs(attrValue))
		return fmt.Sprintf(`%s="%s" data-attrbind="%s"`, attrName, staticValue, elementID)
	})

	return result
}

// isValidFieldName checks if a string is a valid field name (starts with uppercase)
func isValidFieldName(s string) bool {
	if len(s) == 0 {
		return false
	}
	// Must start with uppercase letter
	if s[0] < 'A' || s[0] > 'Z' {
		return false
	}
	// Rest must be alphanumeric
	for i := 1; i < len(s); i++ {
		if !isAlphaNum(s[i]) {
			return false
		}
	}
	return true
}

func isAlphaNum(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

// parseComponentAttrs parses component attributes for props and events
func parseComponentAttrs(attrs string) (props map[string]string, events map[string]componentEvent) {
	props = make(map[string]string)
	events = make(map[string]componentEvent)

	for _, match := range eventRegex.FindAllStringSubmatch(attrs, -1) {
		events[match[1]] = componentEvent{method: match[3], args: match[4]}
	}
	attrs = eventRegex.ReplaceAllString(attrs, "")

	for _, match := range propRegex.FindAllStringSubmatch(attrs, -1) {
		props[match[1]] = match[2]
	}

	return props, events
}
