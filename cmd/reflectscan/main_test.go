package main

import (
	"io"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestScan(t *testing.T) {
	tests := []struct {
		dir  string
		want string
	}{
		// map field with float values; the unexported slice is ignored
		{"maps", "pk_noiface,pk_noslice,pk_noptr,pk_noint,pk_nouint,pk_nobool"},
		// FromJS replaces reflection, so every kind can go
		{"fromjs", "pk_nomap,pk_noiface,pk_noslice,pk_noptr,pk_noint,pk_nouint,pk_nofloat,pk_nobool,pk_nostring"},
		// FetchOptions.Body is encoded ([]int), the response decoded (bool)
		{"encode", "pk_nomap,pk_noiface,pk_noptr,pk_nouint,pk_nofloat,pk_nostring"},
		// a body typed as any could be anything: keep all kinds
		{"anybody", ""},
	}
	for _, tt := range tests {
		t.Run(tt.dir, func(t *testing.T) {
			tags, err := scan(filepath.Join("testdata", tt.dir), io.Discard)
			if err != nil {
				t.Fatal(err)
			}
			if got := strings.Join(tags, ","); got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// Every kind the scanner can report must have a pk_no<kind> tag in
// internal/codecflags.
func TestKindsHaveTags(t *testing.T) {
	files, err := filepath.Glob("../../internal/codecflags/no*.go")
	if err != nil {
		t.Fatal(err)
	}
	var have []string
	for _, f := range files {
		have = append(have, strings.TrimSuffix(strings.TrimPrefix(filepath.Base(f), "no"), ".go"))
	}
	for _, k := range kinds {
		if !slices.Contains(have, k) {
			t.Errorf("kind %q has no internal/codecflags/no%s.go", k, k)
		}
	}
}
