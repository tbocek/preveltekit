package main

import (
	"maps"
	"slices"
	"strings"
	"testing"
)

// wasmWithImports builds a minimal WASM binary: a custom section (skipped),
// then an import section with the given gojs functions, a memory with max
// limit, and a non-gojs function.
func wasmWithImports(gojs ...string) []byte {
	name := func(s string) []byte { return append([]byte{byte(len(s))}, s...) }
	var imports [][]byte
	for _, f := range gojs {
		imports = append(imports, slices.Concat(name("gojs"), name(f), []byte{0, 0}))
	}
	imports = append(imports,
		slices.Concat(name("env"), name("memory"), []byte{2, 1, 1, 2}),
		slices.Concat(name("wasi_snapshot_preview1"), name("random_get"), []byte{0, 1}))
	body := append([]byte{byte(len(imports))}, slices.Concat(imports...)...)
	custom := slices.Concat(name("note"), []byte("hi"))
	return slices.Concat(
		[]byte("\x00asm\x01\x00\x00\x00"),
		[]byte{0, byte(len(custom))}, custom,
		[]byte{2, byte(len(body))}, body,
	)
}

func TestGojsImports(t *testing.T) {
	got, err := gojsImports(wasmWithImports("runtime.ticks", "syscall/js.valueSetIndex"))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"runtime.ticks", "syscall/js.valueSetIndex"}
	if keys := slices.Sorted(maps.Keys(got)); !slices.Equal(keys, want) {
		t.Errorf("got %v, want %v", keys, want)
	}
	if _, err := gojsImports(wasmWithImports("runtime.ticks")[:20]); err == nil {
		t.Error("truncated binary: want error")
	}
}

const js = `gojs: {
  "runtime.ticks": () => {
    return 0;
  },
  // end

  "syscall/js.valueSet": (v_ref) => {
    set(v_ref);
  },
  // end

  "syscall/js.valueSetIndex": (v_ref, i) => {
    setIndex(v_ref, i);
  },
  // end
}
`

func TestStrip(t *testing.T) {
	// valueSet is a prefix of the imported valueSetIndex, but not imported itself.
	out, removed, kept, err := strip(js, map[string]bool{"runtime.ticks": true, "syscall/js.valueSetIndex": true})
	if err != nil {
		t.Fatal(err)
	}
	if removed != 1 || kept != 2 {
		t.Errorf("removed=%d kept=%d, want 1 and 2", removed, kept)
	}
	if strings.Contains(out, "valueSet\"") || !strings.Contains(out, "valueSetIndex") || !strings.Contains(out, "runtime.ticks") {
		t.Errorf("wrong functions stripped:\n%s", out)
	}

	if _, _, _, err := strip(`"runtime.ticks": () => {}`, map[string]bool{}); err == nil {
		t.Error("missing // end: want error")
	}
}
