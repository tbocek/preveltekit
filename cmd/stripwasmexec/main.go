// Command stripwasmexec removes the runtime.* and syscall/js.* functions from
// wasm_exec.js that the WASM module does not import. Each function in
// wasm_exec.js starts at its "name": key and ends at the next "// end" line.
// Afterwards it drops indentation, blank lines and comment lines (except
// the license header at the top).
//
// Usage: stripwasmexec <main.wasm> <wasm_exec.js>
package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// keyRe matches the start of a strippable function in wasm_exec.js.
var keyRe = regexp.MustCompile(`"((?:runtime|syscall/js)\.[^"]+)":`)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: stripwasmexec <main.wasm> <wasm_exec.js>")
		os.Exit(2)
	}
	if err := run(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintln(os.Stderr, "stripwasmexec:", err)
		os.Exit(1)
	}
}

func run(wasmFile, jsFile string) error {
	wasm, err := os.ReadFile(wasmFile)
	if err != nil {
		return err
	}
	imports, err := gojsImports(wasm)
	if err != nil {
		return fmt.Errorf("%s: %w", wasmFile, err)
	}
	if len(imports) == 0 {
		return nil
	}
	js, err := os.ReadFile(jsFile)
	if err != nil {
		return err
	}
	out, removed, kept, err := strip(string(js), imports)
	if err != nil {
		return fmt.Errorf("%s: %w", jsFile, err)
	}
	if err := os.WriteFile(jsFile, []byte(minify(out)), 0o644); err != nil {
		return err
	}
	fmt.Printf("  Stripped %d unused functions (kept %d)\n", removed, kept)
	return nil
}

// minify removes indentation, blank lines and comment lines from js, keeping
// the leading comment block (license header). Line breaks are kept, so
// automatic semicolon insertion is unaffected. Only safe for code without
// multi-line strings, like wasm_exec.js.
func minify(js string) string {
	var b strings.Builder
	header := true
	for _, line := range strings.Split(js, "\n") {
		line = strings.TrimSpace(line)
		comment := strings.HasPrefix(line, "//")
		header = header && comment
		if line == "" || comment && !header {
			continue
		}
		b.WriteString(line + "\n")
	}
	return b.String()
}

// strip removes every function block whose name is not in imports.
func strip(js string, imports map[string]bool) (out string, removed, kept int, err error) {
	lines := strings.SplitAfter(js, "\n")
	var b strings.Builder
	for i := 0; i < len(lines); i++ {
		m := keyRe.FindStringSubmatch(lines[i])
		if m == nil {
			b.WriteString(lines[i])
			continue
		}
		if imports[m[1]] {
			kept++
			b.WriteString(lines[i])
			continue
		}
		end := i
		for end < len(lines) && !strings.Contains(lines[end], "// end") {
			end++
		}
		if end == len(lines) {
			return "", 0, 0, fmt.Errorf("no // end after %s", m[1])
		}
		removed++
		i = end
	}
	return b.String(), removed, kept, nil
}

// gojsImports returns the names imported from the "gojs" module, e.g.
// "syscall/js.valueGet", by reading the import section of a WASM binary.
func gojsImports(b []byte) (map[string]bool, error) {
	if len(b) < 8 || string(b[:4]) != "\x00asm" {
		return nil, errors.New("not a wasm binary")
	}
	r := &reader{b: b, off: 8}
	for r.off < len(b) && r.err == nil {
		id := r.byte()
		size := r.uleb()
		if id != 2 { // not the import section
			r.off += int(size)
			continue
		}
		imports := map[string]bool{}
		for n := r.uleb(); n > 0 && r.err == nil; n-- {
			module, name := r.name(), r.name()
			switch r.byte() {
			case 0: // func: type index
				r.uleb()
			case 1: // table: reftype, limits
				r.byte()
				r.limits()
			case 2: // memory: limits
				r.limits()
			case 3: // global: valtype, mutability
				r.byte()
				r.byte()
			case 4: // tag: attribute, type index
				r.byte()
				r.uleb()
			default:
				return nil, errors.New("unknown import kind")
			}
			if module == "gojs" {
				imports[name] = true
			}
		}
		return imports, r.err
	}
	return nil, r.err
}

type reader struct {
	b   []byte
	off int
	err error
}

func (r *reader) byte() byte {
	if r.off >= len(r.b) {
		r.err = errors.New("unexpected end of wasm binary")
		return 0
	}
	r.off++
	return r.b[r.off-1]
}

func (r *reader) uleb() uint64 {
	var v uint64
	for shift := 0; shift < 64 && r.err == nil; shift += 7 {
		c := r.byte()
		v |= uint64(c&0x7f) << shift
		if c&0x80 == 0 {
			break
		}
	}
	return v
}

func (r *reader) name() string {
	n := int(r.uleb())
	if r.err != nil || r.off+n > len(r.b) {
		r.err = errors.New("unexpected end of wasm binary")
		return ""
	}
	r.off += n
	return string(r.b[r.off-n : r.off])
}

func (r *reader) limits() {
	if flags := r.byte(); flags&1 != 0 {
		r.uleb()
		r.uleb()
	} else {
		r.uleb()
	}
}
