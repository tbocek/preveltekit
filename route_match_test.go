package preveltekit

import "testing"

func TestMatchRoute(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		match   bool
		params  map[string]string
	}{
		// Exact matches
		{"/", "/", true, nil},
		{"/bitcoin", "/bitcoin", true, nil},
		{"/a/b/c", "/a/b/c", true, nil},
		{"/bitcoin", "/manual", false, nil},
		{"/bitcoin", "/", false, nil},
		{"/", "/bitcoin", false, nil},

		// Parameter capture
		{"/users/:id", "/users/42", true, map[string]string{"id": "42"}},
		{"/users/:id/posts/:pid", "/users/1/posts/99", true, map[string]string{"id": "1", "pid": "99"}},
		{"/users/:id", "/users", false, nil},
		{"/users/:id", "/users/42/extra", false, nil},

		// Single wildcard *
		{"/*", "/anything", true, nil},
		{"/*", "/", false, nil},
		{"/*/bitcoin", "/base/bitcoin", true, nil},
		{"/*/bitcoin", "/bitcoin", false, nil},
		{"/*/bitcoin", "/a/b/bitcoin", false, nil},
		{"/a/*/b", "/a/x/b", true, nil},
		{"/a/*/b", "/a/x/y/b", false, nil},
		{"/*/*", "/a/b", true, nil},
		{"/*/*", "/a", false, nil},

		// Double wildcard ** (zero or more)
		{"/**", "/anything", true, nil},
		{"/**", "/a/b/c", true, nil},
		{"/**", "/", true, nil},
		{"/**/bitcoin", "/bitcoin", true, nil},
		{"/**/bitcoin", "/base/bitcoin", true, nil},
		{"/**/bitcoin", "/a/b/bitcoin", true, nil},
		{"/**/bitcoin", "/a/b/manual", false, nil},
		{"/a/**/b", "/a/x/b", true, nil},
		{"/a/**/b", "/a/x/y/b", true, nil},
		{"/a/**/b", "/a/x/y/z/b", true, nil},
		{"/a/**/b", "/a/b", true, nil}, // ** matches zero

		// Mixed
		{"/a/**/b/*/c", "/a/x/b/y/c", true, nil},
		{"/a/**/b/*/c", "/a/x/y/b/z/c", true, nil},
		{"/a/**/b/:id", "/a/x/b/42", true, map[string]string{"id": "42"}},

		// Relative patterns without leading / (prefix wildcards)
		// */ — one prefix segment
		{"*/", "/", false, nil},               // * needs exactly one segment
		{"*/", "/preveltekit", true, nil},     // one segment = match
		{"*/", "/a/b", false, nil},            // two segments = no match
		{"*/bitcoin", "/bitcoin", false, nil}, // * needs exactly one prefix
		{"*/bitcoin", "/preveltekit/bitcoin", true, nil},
		{"*/bitcoin", "/a/b/bitcoin", false, nil}, // two prefix segments = no match
		{"*/manual", "/preveltekit/manual", true, nil},
		{"*/manual", "/manual", false, nil},

		// **/ — zero or more prefix segments
		{"**/", "/", true, nil},
		{"**/", "/preveltekit", true, nil},
		{"**/", "/a/b", true, nil},
		{"**/bitcoin", "/bitcoin", true, nil},
		{"**/bitcoin", "/preveltekit/bitcoin", true, nil},
		{"**/bitcoin", "/a/b/bitcoin", true, nil},
		{"**/bitcoin", "/a/b/manual", false, nil},
		{"**/manual", "/manual", true, nil},
		{"**/manual", "/preveltekit/manual", true, nil},
	}

	for _, tt := range tests {
		params, _, ok := matchRoute(tt.pattern, tt.path)
		if ok != tt.match {
			t.Errorf("matchRoute(%q, %q) = %v, want %v", tt.pattern, tt.path, ok, tt.match)
			continue
		}
		if tt.match && tt.params != nil {
			for k, v := range tt.params {
				if params[k] != v {
					t.Errorf("matchRoute(%q, %q) param %q = %q, want %q", tt.pattern, tt.path, k, params[k], v)
				}
			}
		}
	}
}

// TestResolveRoute tests that route patterns are correctly resolved against a base path.
func TestResolveRoute(t *testing.T) {
	tests := []struct {
		base      string
		routePath string
		want      string
	}{
		// Absolute routes — returned as-is
		{"/preveltekit", "/", "/"},
		{"/preveltekit", "/bitcoin", "/bitcoin"},
		{"/preveltekit", "/*/bitcoin", "/*/bitcoin"},
		{"/preveltekit", "/**/bitcoin", "/**/bitcoin"},

		// Relative routes — joined with base
		{"/preveltekit", "bitcoin", "/preveltekit/bitcoin"},
		{"/preveltekit", ".", "/preveltekit"},
		{"/preveltekit", "", "/preveltekit"},
		{"/preveltekit", "./", "/preveltekit"},
		{"/preveltekit", "*", "/preveltekit/*"},
		{"/preveltekit", "**", "/preveltekit/**"},
		{"/preveltekit", "*/", "/preveltekit/*/"},
		{"/preveltekit", "**/", "/preveltekit/**/"},
		{"/preveltekit", "*/bitcoin", "/preveltekit/*/bitcoin"},
		{"/preveltekit", "**/bitcoin", "/preveltekit/**/bitcoin"},

		// Edge: root base
		{"/", "bitcoin", "/bitcoin"},
		{"/", ".", "/"},
		{"/", "", "/"},
	}

	for _, tt := range tests {
		got := resolveRoute(tt.base, tt.routePath)
		if got != tt.want {
			t.Errorf("resolveRoute(%q, %q) = %q, want %q", tt.base, tt.routePath, got, tt.want)
		}
	}
}

// TestResolvedRouteMatching tests the full flow: resolve route against base, then match.
// Based on the agreed table (base = /preveltekit):
//
//	#  | Route Path     | Example pathnames that match
//	---|----------------|--------------------------------------------
//	1  | /              | / only
//	2  | /bitcoin       | /bitcoin only
//	3  | bitcoin        | /preveltekit/bitcoin only
//	4  | .              | /preveltekit only
//	5  | "" (empty)     | /preveltekit only
//	6  | ./             | /preveltekit only
//	7  | *              | /preveltekit/foo, /preveltekit/bar
//	8  | **             | /preveltekit, /preveltekit/foo, /preveltekit/foo/bar
//	9  | */             | same as 7
//	10 | **/            | same as 8
//	11 | */bitcoin      | /preveltekit/foo/bitcoin
//	12 | **/bitcoin     | /preveltekit/bitcoin, /preveltekit/foo/bitcoin
//	13 | /*/bitcoin     | /foo/bitcoin
//	14 | /**/bitcoin    | /bitcoin, /foo/bitcoin, /a/b/bitcoin
func TestResolvedRouteMatching(t *testing.T) {
	const base = "/preveltekit"

	tests := []struct {
		name      string
		routePath string
		path      string
		match     bool
	}{
		// #1: / — absolute root
		{"1: / matches /", "/", "/", true},
		{"1: / no match /preveltekit", "/", "/preveltekit", false},
		{"1: / no match /bitcoin", "/", "/bitcoin", false},

		// #2: /bitcoin — absolute
		{"2: /bitcoin matches /bitcoin", "/bitcoin", "/bitcoin", true},
		{"2: /bitcoin no match /preveltekit/bitcoin", "/bitcoin", "/preveltekit/bitcoin", false},
		{"2: /bitcoin no match /", "/bitcoin", "/", false},

		// #3: bitcoin — relative → /preveltekit/bitcoin
		{"3: bitcoin matches /preveltekit/bitcoin", "bitcoin", "/preveltekit/bitcoin", true},
		{"3: bitcoin no match /bitcoin", "bitcoin", "/bitcoin", false},
		{"3: bitcoin no match /preveltekit", "bitcoin", "/preveltekit", false},

		// #4: . → /preveltekit
		{"4: . matches /preveltekit", ".", "/preveltekit", true},
		{"4: . no match /", ".", "/", false},
		{"4: . no match /preveltekit/foo", ".", "/preveltekit/foo", false},

		// #5: "" (empty) → /preveltekit
		{"5: empty matches /preveltekit", "", "/preveltekit", true},
		{"5: empty no match /", "", "/", false},

		// #6: ./ → /preveltekit
		{"6: ./ matches /preveltekit", "./", "/preveltekit", true},
		{"6: ./ no match /", "./", "/", false},

		// #7: * → /preveltekit/* (one segment after base)
		{"7: * matches /preveltekit/foo", "*", "/preveltekit/foo", true},
		{"7: * matches /preveltekit/bar", "*", "/preveltekit/bar", true},
		{"7: * no match /preveltekit", "*", "/preveltekit", false},
		{"7: * no match /preveltekit/foo/bar", "*", "/preveltekit/foo/bar", false},
		{"7: * no match /foo", "*", "/foo", false},

		// #8: ** → /preveltekit/** (zero or more after base)
		{"8: ** matches /preveltekit", "**", "/preveltekit", true},
		{"8: ** matches /preveltekit/foo", "**", "/preveltekit/foo", true},
		{"8: ** matches /preveltekit/foo/bar", "**", "/preveltekit/foo/bar", true},
		{"8: ** no match /", "**", "/", false},
		{"8: ** no match /foo", "**", "/foo", false},

		// #9: */ — same as #7
		{"9: */ matches /preveltekit/foo", "*/", "/preveltekit/foo", true},
		{"9: */ no match /preveltekit", "*/", "/preveltekit", false},
		{"9: */ no match /preveltekit/foo/bar", "*/", "/preveltekit/foo/bar", false},

		// #10: **/ — same as #8
		{"10: **/ matches /preveltekit", "**/", "/preveltekit", true},
		{"10: **/ matches /preveltekit/foo", "**/", "/preveltekit/foo", true},
		{"10: **/ matches /preveltekit/foo/bar", "**/", "/preveltekit/foo/bar", true},
		{"10: **/ no match /", "**/", "/", false},

		// #11: */bitcoin → /preveltekit/*/bitcoin
		{"11: */bitcoin matches /preveltekit/foo/bitcoin", "*/bitcoin", "/preveltekit/foo/bitcoin", true},
		{"11: */bitcoin no match /preveltekit/bitcoin", "*/bitcoin", "/preveltekit/bitcoin", false},
		{"11: */bitcoin no match /preveltekit/a/b/bitcoin", "*/bitcoin", "/preveltekit/a/b/bitcoin", false},
		{"11: */bitcoin no match /foo/bitcoin", "*/bitcoin", "/foo/bitcoin", false},

		// #12: **/bitcoin → /preveltekit/**/bitcoin
		{"12: **/bitcoin matches /preveltekit/bitcoin", "**/bitcoin", "/preveltekit/bitcoin", true},
		{"12: **/bitcoin matches /preveltekit/foo/bitcoin", "**/bitcoin", "/preveltekit/foo/bitcoin", true},
		{"12: **/bitcoin matches /preveltekit/a/b/bitcoin", "**/bitcoin", "/preveltekit/a/b/bitcoin", true},
		{"12: **/bitcoin no match /bitcoin", "**/bitcoin", "/bitcoin", false},
		{"12: **/bitcoin no match /foo/bitcoin", "**/bitcoin", "/foo/bitcoin", false},

		// #13: /*/bitcoin — absolute
		{"13: /*/bitcoin matches /foo/bitcoin", "/*/bitcoin", "/foo/bitcoin", true},
		{"13: /*/bitcoin matches /preveltekit/bitcoin", "/*/bitcoin", "/preveltekit/bitcoin", true},
		{"13: /*/bitcoin no match /bitcoin", "/*/bitcoin", "/bitcoin", false},
		{"13: /*/bitcoin no match /a/b/bitcoin", "/*/bitcoin", "/a/b/bitcoin", false},

		// #14: /**/bitcoin — absolute
		{"14: /**/bitcoin matches /bitcoin", "/**/bitcoin", "/bitcoin", true},
		{"14: /**/bitcoin matches /foo/bitcoin", "/**/bitcoin", "/foo/bitcoin", true},
		{"14: /**/bitcoin matches /a/b/bitcoin", "/**/bitcoin", "/a/b/bitcoin", true},
		{"14: /**/bitcoin no match /a/b/manual", "/**/bitcoin", "/a/b/manual", false},
	}

	for _, tt := range tests {
		resolved := resolveRoute(base, tt.routePath)
		_, _, ok := matchRoute(resolved, tt.path)
		if ok != tt.match {
			t.Errorf("%s: resolveRoute(%q, %q)=%q → matchRoute(%q, %q) = %v, want %v",
				tt.name, base, tt.routePath, resolved, resolved, tt.path, ok, tt.match)
		}
	}
}

func TestSpecificity(t *testing.T) {
	tests := []struct {
		a, b    string
		path    string
		aHigher bool // true if a should have higher specificity
	}{
		// Exact beats param
		{"/users/admin", "/users/:id", "/users/admin", true},
		// Param beats *
		{"/users/:id", "/users/*", "/users/42", true},
		// * beats **
		{"/*/bitcoin", "/**/bitcoin", "/base/bitcoin", true},
		// More exact segments = higher
		{"/a/b", "/*/b", "/a/b", true},
	}

	for _, tt := range tests {
		_, specA, okA := matchRoute(tt.a, tt.path)
		_, specB, okB := matchRoute(tt.b, tt.path)
		if !okA || !okB {
			t.Errorf("both %q and %q should match %q (okA=%v, okB=%v)", tt.a, tt.b, tt.path, okA, okB)
			continue
		}
		if tt.aHigher && specA <= specB {
			t.Errorf("specificity(%q)=%d should be > specificity(%q)=%d for path %q", tt.a, specA, tt.b, specB, tt.path)
		}
		if !tt.aHigher && specA >= specB {
			t.Errorf("specificity(%q)=%d should be < specificity(%q)=%d for path %q", tt.a, specA, tt.b, specB, tt.path)
		}
	}
}
