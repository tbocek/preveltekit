package preveltekit

import (
	"strings"
)

// Shared HTML utility functions used by both SSR (node_html.go) and WASM (node_html_wasm.go).

// injectAttrs injects attributes into an HTML element string.
// Finds the first > and inserts the attrs just before it.
func injectAttrs(html, attrs string) string {
	tagEnd := findTagEnd(html)
	if tagEnd == -1 {
		return html + " " + attrs
	}
	if tagEnd > 0 && html[tagEnd-1] == '/' {
		return html[:tagEnd-1] + " " + attrs + " />" + html[tagEnd+1:]
	}
	return html[:tagEnd] + " " + attrs + html[tagEnd:]
}

// findTagEnd returns the index of the first '>' in html, or -1 if not found.
func findTagEnd(html string) int {
	for i := 0; i < len(html); i++ {
		if html[i] == '>' {
			return i
		}
	}
	return -1
}

// escapeAttr escapes attribute values (quotes and ampersands).
func escapeAttr(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}

// injectIDAndMergeAttrs injects id and merges attribute values into the first HTML tag.
// For "class", merges with existing class attribute. For others, values are space-joined.
func injectIDAndMergeAttrs(html, id string, attrValues map[string][]string, extraAttrs string) string {
	tagEnd := findTagEnd(html)
	if tagEnd == -1 {
		return html
	}

	openingTag := html[:tagEnd]
	rest := html[tagEnd:]

	if tagEnd > 0 && html[tagEnd-1] == '/' {
		openingTag = html[:tagEnd-1]
		rest = html[tagEnd-1:]
	}

	// Build new attributes
	newAttrs := `id="` + id + `"`

	// Handle class attribute specially - merge with existing
	if classes, ok := attrValues["class"]; ok && len(classes) > 0 {
		classIdx := strings.Index(openingTag, `class="`)
		if classIdx != -1 {
			classStart := classIdx + 7
			classEnd := strings.Index(openingTag[classStart:], `"`)
			if classEnd != -1 {
				classEnd += classStart
				existingClasses := openingTag[classStart:classEnd]
				mergedClasses := existingClasses
				for _, c := range classes {
					if c != "" {
						mergedClasses += " " + c
					}
				}
				openingTag = openingTag[:classIdx] + openingTag[classEnd+1:]
				newAttrs += ` class="` + strings.TrimSpace(mergedClasses) + `"`
			}
		} else {
			newAttrs += ` class="` + strings.Join(classes, " ") + `"`
		}
		delete(attrValues, "class")
	}

	// Handle other attributes
	for name, values := range attrValues {
		if len(values) > 0 {
			attrPattern := name + `="`
			attrIdx := strings.Index(openingTag, attrPattern)
			if attrIdx != -1 {
				attrStart := attrIdx + len(attrPattern)
				attrEnd := strings.Index(openingTag[attrStart:], `"`)
				if attrEnd != -1 {
					attrEnd += attrStart
					existingValue := openingTag[attrStart:attrEnd]
					mergedValue := existingValue
					for _, v := range values {
						if v != "" {
							mergedValue += " " + v
						}
					}
					openingTag = openingTag[:attrIdx] + openingTag[attrEnd+1:]
					newAttrs += ` ` + name + `="` + strings.TrimSpace(mergedValue) + `"`
				}
			} else {
				newAttrs += ` ` + name + `="` + strings.Join(values, " ") + `"`
			}
		}
	}

	newAttrs += extraAttrs

	insertIdx := 0
	for i := 1; i < len(openingTag); i++ {
		if openingTag[i] == ' ' || openingTag[i] == '/' {
			insertIdx = i
			break
		}
	}
	if insertIdx == 0 {
		insertIdx = len(openingTag)
	}

	return openingTag[:insertIdx] + " " + newAttrs + openingTag[insertIdx:] + rest
}

// injectScopeClass injects a scope class into every opening HTML tag in the string.
func injectScopeClass(html, scopeClass string) string {
	var sb strings.Builder
	sb.Grow(len(html) + len(html)/10)
	i := 0
	for i < len(html) {
		if html[i] == '<' && i+1 < len(html) {
			next := html[i+1]
			if next == '/' || next == '!' {
				end := strings.IndexByte(html[i:], '>')
				if end == -1 {
					sb.WriteString(html[i:])
					break
				}
				sb.WriteString(html[i : i+end+1])
				i += end + 1
				continue
			}
			j := i + 1
			inQuote := byte(0)
			for j < len(html) {
				if inQuote != 0 {
					if html[j] == inQuote {
						inQuote = 0
					}
				} else if html[j] == '"' || html[j] == '\'' {
					inQuote = html[j]
				} else if html[j] == '>' {
					break
				}
				j++
			}
			if j >= len(html) {
				// Tag not closed in this string part (split across parts).
				// Still inject the scope class since a later part closes the tag.
				tagContent := html[i:]
				classIdx := strings.Index(tagContent, `class="`)
				if classIdx != -1 {
					quoteStart := classIdx + 7
					quoteEnd := strings.IndexByte(tagContent[quoteStart:], '"')
					if quoteEnd != -1 {
						quoteEnd += quoteStart
						sb.WriteString(tagContent[:quoteEnd])
						sb.WriteByte(' ')
						sb.WriteString(scopeClass)
						sb.WriteString(tagContent[quoteEnd:])
					} else {
						// class attribute quote not closed — append before trailing space/end
						sb.WriteString(tagContent)
					}
				} else {
					// No class attribute — inject one before the trailing content
					// Find end of tag name
					k := 1
					for k < len(tagContent) && tagContent[k] != ' ' && tagContent[k] != '/' {
						k++
					}
					sb.WriteString(tagContent[:k])
					sb.WriteString(` class="`)
					sb.WriteString(scopeClass)
					sb.WriteByte('"')
					sb.WriteString(tagContent[k:])
				}
				break
			}
			tagContent := html[i:j]
			selfClosing := j > 0 && html[j-1] == '/'
			if selfClosing {
				tagContent = html[i : j-1]
			}

			classIdx := strings.Index(tagContent, `class="`)
			if classIdx != -1 {
				quoteStart := classIdx + 7
				quoteEnd := strings.IndexByte(tagContent[quoteStart:], '"')
				if quoteEnd != -1 {
					quoteEnd += quoteStart
					sb.WriteString(tagContent[:quoteEnd])
					sb.WriteByte(' ')
					sb.WriteString(scopeClass)
					sb.WriteString(tagContent[quoteEnd:])
					if selfClosing {
						sb.WriteString("/>")
					} else {
						sb.WriteByte('>')
					}
					i = j + 1
					continue
				}
			}
			sb.WriteString(tagContent)
			sb.WriteString(` class="`)
			sb.WriteString(scopeClass)
			sb.WriteByte('"')
			if selfClosing {
				sb.WriteString("/>")
			} else {
				sb.WriteByte('>')
			}
			i = j + 1
		} else {
			sb.WriteByte(html[i])
			i++
		}
	}
	return sb.String()
}
