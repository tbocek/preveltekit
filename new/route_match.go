package preveltekit

import "strings"

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
	if strings.HasPrefix(pattern, "*/") {
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

// splitPath splits a URL path into non-empty segments.
func splitPath(s string) []string {
	s = strings.Trim(s, "/")
	if s == "" {
		return nil
	}
	parts := strings.Split(s, "/")
	// Filter empty segments (from double slashes)
	n := 0
	for _, p := range parts {
		if p != "" {
			parts[n] = p
			n++
		}
	}
	return parts[:n]
}
