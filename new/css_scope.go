package preveltekit

import "strings"

// scopeCSS adds a scope class selector to every CSS rule's selectors.
// e.g., ".demo button" becomes ".demo.v0 button.v0"
// Handles @media (recurse), @keyframes (skip), pseudo-classes, combinators.
func scopeCSS(css, scopeClass string) string {
	scope := "." + scopeClass
	var sb strings.Builder
	sb.Grow(len(css) + len(css)/5)
	i := 0
	for i < len(css) {
		// Skip whitespace
		if css[i] == ' ' || css[i] == '\t' || css[i] == '\n' || css[i] == '\r' {
			sb.WriteByte(css[i])
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
			atRule := strings.TrimSpace(css[i:j])
			if j < len(css) && css[j] == ';' {
				// @import or similar — pass through
				sb.WriteString(css[i : j+1])
				i = j + 1
				continue
			}
			if j < len(css) && css[j] == '{' {
				if strings.HasPrefix(atRule, "@keyframes") || strings.HasPrefix(atRule, "@font-face") {
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
					sb.WriteString(css[i:k])
					i = k
					continue
				}
				// @media or similar — write the @-rule header, recurse into body
				sb.WriteString(css[i : j+1]) // "@media (...) {"
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
				sb.WriteString(scopeCSS(body, scopeClass))
				sb.WriteByte('}')
				i = k
				continue
			}
		}
		// Regular rule: find selector(s) before {
		braceIdx := strings.IndexByte(css[i:], '{')
		if braceIdx == -1 {
			sb.WriteString(css[i:])
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
		selectors := strings.Split(selectorPart, ",")
		for si, sel := range selectors {
			if si > 0 {
				sb.WriteByte(',')
			}
			sb.WriteString(scopeSelector(strings.TrimSpace(sel), scope))
		}
		sb.WriteByte('{')
		sb.WriteString(body)
		sb.WriteByte('}')
	}
	return sb.String()
}

// scopeSelector adds .vN to each simple selector in a compound selector.
// e.g., ".demo button:hover" → ".demo.v0 button.v0:hover"
func scopeSelector(sel, scope string) string {
	if sel == "" {
		return sel
	}
	// Don't scope html/body selectors
	trimmed := strings.TrimSpace(sel)
	if trimmed == "html" || trimmed == "body" || trimmed == "*" {
		return sel
	}

	var sb strings.Builder
	// Split by combinators (space, >, +, ~) preserving them
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
			sb.WriteByte(' ')
			sb.WriteByte(sel[i])
			sb.WriteByte(' ')
			i++
			continue
		}

		// Implicit descendant combinator (space between segments)
		if sb.Len() > 0 {
			sb.WriteByte(' ')
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
			sb.WriteString(segment)
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
		sb.WriteString(segment[:insertPos])
		sb.WriteString(scope)
		sb.WriteString(segment[insertPos:])
	}
	return sb.String()
}
