package preveltekit

// scopeCSS adds a scope class selector to every CSS rule's selectors.
// e.g., ".demo button" becomes ".demo.v0 button.v0"
// Handles @media (recurse), @keyframes (skip), pseudo-classes, combinators.
func scopeCSS(css, scopeClass string) string {
	scope := "." + scopeClass
	buf := make([]byte, 0, len(css)+len(css)/5)
	i := 0
	for i < len(css) {
		// Skip whitespace
		if css[i] == ' ' || css[i] == '\t' || css[i] == '\n' || css[i] == '\r' {
			buf = append(buf, css[i])
			i++
			continue
		}
		// Handle @-rules
		if css[i] == '@' {
			// Find the rule name
			j := i + 1
			for j < len(css) && css[j] != '{' && css[j] != ';' {
				j++
			}
			atRule := trimSpace(css[i:j])
			if j < len(css) && css[j] == ';' {
				// @import or similar — pass through
				buf = append(buf, css[i:j+1]...)
				i = j + 1
				continue
			}
			if j < len(css) && css[j] == '{' {
				if hasPrefix(atRule, "@keyframes") || hasPrefix(atRule, "@font-face") {
					// Skip scoping — copy the entire block as-is
					depth := 1
					k := j + 1
					for k < len(css) && depth > 0 {
						if css[k] == '{' {
							depth++
						} else if css[k] == '}' {
							depth--
						}
						k++
					}
					buf = append(buf, css[i:k]...)
					i = k
					continue
				}
				// @media or similar — write the @-rule header, recurse into body
				buf = append(buf, css[i:j+1]...) // "@media (...) {"
				depth := 1
				k := j + 1
				bodyStart := k
				for k < len(css) && depth > 0 {
					if css[k] == '{' {
						depth++
					} else if css[k] == '}' {
						depth--
					}
					k++
				}
				// body is css[bodyStart:k-1], closing } is at k-1
				body := css[bodyStart : k-1]
				buf = append(buf, scopeCSS(body, scopeClass)...)
				buf = append(buf, '}')
				i = k
				continue
			}
		}
		// Regular rule: find selector(s) before {
		braceIdx := indexByte(css[i:], '{')
		if braceIdx == -1 {
			buf = append(buf, css[i:]...)
			break
		}
		selectorPart := css[i : i+braceIdx]
		i += braceIdx + 1

		// Find the closing }
		depth := 1
		k := i
		for k < len(css) && depth > 0 {
			if css[k] == '{' {
				depth++
			} else if css[k] == '}' {
				depth--
			}
			k++
		}
		body := css[i : k-1]
		i = k

		// Scope the selectors
		selectors := splitByte(selectorPart, ',')
		for si, sel := range selectors {
			if si > 0 {
				buf = append(buf, ',')
			}
			buf = append(buf, scopeSelector(trimSpace(sel), scope)...)
		}
		buf = append(buf, '{')
		buf = append(buf, body...)
		buf = append(buf, '}')
	}
	return string(buf)
}

// scopeSelector adds .vN to each simple selector in a compound selector.
// e.g., ".demo button:hover" → ".demo.v0 button.v0:hover"
func scopeSelector(sel, scope string) string {
	if sel == "" {
		return sel
	}
	// Don't scope html/body selectors
	trimmed := trimSpace(sel)
	if trimmed == "html" || trimmed == "body" || trimmed == "*" {
		return sel
	}

	buf := make([]byte, 0, len(sel)+len(scope)*2)
	i := 0
	for i < len(sel) {
		// Skip leading whitespace
		for i < len(sel) && (sel[i] == ' ' || sel[i] == '\t' || sel[i] == '\n') {
			i++
		}
		if i >= len(sel) {
			break
		}

		// Check for combinator characters
		if sel[i] == '>' || sel[i] == '+' || sel[i] == '~' {
			buf = append(buf, ' ')
			buf = append(buf, sel[i])
			buf = append(buf, ' ')
			i++
			continue
		}

		// Implicit descendant combinator (space between segments)
		if len(buf) > 0 {
			buf = append(buf, ' ')
		}

		// Read a simple selector segment (up to space or combinator)
		start := i
		for i < len(sel) && sel[i] != ' ' && sel[i] != '\t' && sel[i] != '>' && sel[i] != '+' && sel[i] != '~' {
			i++
		}
		segment := sel[start:i]
		if segment == "" {
			continue
		}

		// Don't scope html/body
		if segment == "html" || segment == "body" {
			buf = append(buf, segment...)
			continue
		}

		// Find where to insert scope: before pseudo-class/pseudo-element
		insertPos := len(segment)
		for j := 0; j < len(segment); j++ {
			if segment[j] == ':' {
				insertPos = j
				break
			}
		}
		buf = append(buf, segment[:insertPos]...)
		buf = append(buf, scope...)
		buf = append(buf, segment[insertPos:]...)
	}
	return string(buf)
}

// trimSpace removes leading and trailing whitespace.
func trimSpace(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

// hasPrefix reports whether s starts with prefix.
func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// indexByte returns the index of the first occurrence of c in s, or -1.
func indexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

// splitByte splits s by separator byte, returning all parts.
func splitByte(s string, sep byte) []string {
	n := 1
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			n++
		}
	}
	parts := make([]string, 0, n)
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}
