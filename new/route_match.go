package preveltekit

// matchRoute matches a path against a route pattern.
// Returns extracted params, specificity score, and whether it matched.
func matchRoute(pattern, path string) (map[string]string, int, bool) {
	params := make(map[string]string)

	// Handle root path
	if pattern == "/" {
		if path == "/" {
			return params, 100, true
		}
		return nil, 0, false
	}

	// Handle catch-all pattern
	if pattern == "*" || pattern == "**" {
		return params, 1, true
	}

	// Handle wildcard prefix patterns like */suffix
	if hasPrefix(pattern, "*/") {
		suffix := pattern[2:]
		if path == "/"+suffix {
			return params, 2, true
		}
		// Match /{segment}/{suffix}
		pathSegs := splitPath(path)
		if len(pathSegs) >= 2 && pathSegs[len(pathSegs)-1] == suffix {
			return params, 2, true
		}
		return nil, 0, false
	}

	// Standard segment-based matching
	patternSegs := splitPath(pattern)
	pathSegs := splitPath(path)

	if len(patternSegs) != len(pathSegs) {
		return nil, 0, false
	}

	specificity := 0
	for i, seg := range patternSegs {
		if len(seg) > 0 && seg[0] == ':' {
			// Parameter segment
			paramName := seg[1:]
			params[paramName] = pathSegs[i]
			specificity += 5
		} else if seg == pathSegs[i] {
			// Exact match
			specificity += 10
		} else {
			return nil, 0, false
		}
	}

	return params, specificity, true
}

// trimSlashes removes leading and trailing slashes from a path.
func trimSlashes(s string) string {
	start, end := 0, len(s)
	for start < end && s[start] == '/' {
		start++
	}
	for end > start && s[end-1] == '/' {
		end--
	}
	return s[start:end]
}

// splitPath splits a URL path into non-empty segments.
func splitPath(s string) []string {
	s = trimSlashes(s)
	if s == "" {
		return nil
	}
	n := 1
	for i := 0; i < len(s); i++ {
		if s[i] == '/' {
			n++
		}
	}
	parts := make([]string, 0, n)
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == '/' {
			if start < i {
				parts = append(parts, s[start:i])
			}
			start = i + 1
		}
	}
	return parts
}
