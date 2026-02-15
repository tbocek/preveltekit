package preveltekit

import "strings"

// resolveRoute resolves a route pattern against a base path.
// Absolute patterns (leading /) are returned as-is.
// Relative patterns are joined with the base path.
// Special cases: "." and "" resolve to the base path itself.
func resolveRoute(basePath, routePath string) string {
	// Absolute route — use as-is
	if len(routePath) > 0 && routePath[0] == '/' {
		return routePath
	}

	// ".", "", "./" all mean "the base path itself"
	trimmed := strings.TrimRight(routePath, "/")
	if trimmed == "" || trimmed == "." {
		return basePath
	}

	// Relative — join with base
	base := strings.TrimRight(basePath, "/")
	return base + "/" + routePath
}

// matchRoute matches a URL path against a route pattern.
//
// Pattern syntax (each token is a path segment separated by /):
//
//	/bitcoin        exact segment match
//	/users/:id      named parameter — captures one segment
//	/*              wildcard — matches exactly one segment
//	/**             globstar — matches zero or more segments
//
// Examples:
//
//	/*/bitcoin      matches /x/bitcoin
//	/**/bitcoin     matches /bitcoin, /x/bitcoin, /x/y/bitcoin
//	/a/*/b          matches /a/x/b
//	/a/**/b         matches /a/b, /a/x/b, /a/x/y/b
//
// Specificity: exact=10, :param=5, *=2, **=1. Higher wins.
func matchRoute(pattern, path string) (map[string]string, int, bool) {
	params := make(map[string]string)
	score, ok := match(splitPath(pattern), splitPath(path), params)
	if !ok {
		return nil, 0, false
	}
	return params, score, true
}

// match recursively matches pattern segments against path segments.
func match(pat, path []string, params map[string]string) (int, bool) {
	// Base case: pattern exhausted — path must also be exhausted
	if len(pat) == 0 {
		if len(path) == 0 {
			return 0, true
		}
		return 0, false
	}

	seg := pat[0]

	// ** — consume zero or more path segments
	if seg == "**" {
		for n := 0; n <= len(path); n++ {
			if score, ok := match(pat[1:], path[n:], params); ok {
				return 1 + score, true
			}
		}
		return 0, false
	}

	// All other tokens need at least one path segment
	if len(path) == 0 {
		return 0, false
	}

	switch {
	case seg == "*":
		if score, ok := match(pat[1:], path[1:], params); ok {
			return 2 + score, true
		}
	case seg[0] == ':':
		params[seg[1:]] = path[0]
		if score, ok := match(pat[1:], path[1:], params); ok {
			return 5 + score, true
		}
		delete(params, seg[1:])
	default:
		if seg == path[0] {
			if score, ok := match(pat[1:], path[1:], params); ok {
				return 10 + score, true
			}
		}
	}

	return 0, false
}

// splitPath splits a URL path into non-empty segments.
// "/a/b/c" → ["a", "b", "c"], "/" → nil, "" → nil
func splitPath(s string) []string {
	s = strings.Trim(s, "/")
	if s == "" {
		return nil
	}
	parts := strings.Split(s, "/")
	n := 0
	for _, p := range parts {
		if p != "" {
			parts[n] = p
			n++
		}
	}
	return parts[:n]
}
