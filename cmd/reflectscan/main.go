// Command reflectscan finds the reflect kinds an app passes through the
// preveltekit JS codec (Decode/Encode, including Get/Post/Fetch/...) and
// prints build tags (pk_no<kind>, see internal/codecflags) for the kinds it
// never uses. build.sh passes them to tinygo so the unused codec code is
// compiled out.
//
// Usage: reflectscan [dir]
//
// The tags are printed comma-separated to stdout, the reasons to stderr. If
// a type cannot be determined statically (e.g. a body typed as any), no tags
// are printed and all kinds are kept.
package main

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

const pkgPath = "github.com/tbocek/preveltekit/v2"

// kinds in the order the tags are printed; names match internal/codecflags.
var kinds = []string{"map", "iface", "slice", "ptr", "int", "uint", "float", "bool", "string"}

// Framework functions whose type argument is decoded from JS.
var decodeFuncs = map[string]bool{"Decode": true, "Get": true, "Post": true, "Put": true, "Patch": true, "Delete": true, "Fetch": true}

// Framework functions that encode an argument, with the argument index.
var encodeArgs = map[string]int{"Encode": 0, "Post": 1, "Put": 1, "Patch": 1}

type listedPackage struct {
	ImportPath string
	Dir        string
	GoFiles    []string
	Export     string
	DepOnly    bool
	Error      *struct{ Err string }
}

type scanner struct {
	fset    *token.FileSet
	used    map[string]string // kind -> first place it was found
	keepAll string            // non-empty: reason why all kinds are kept
	visited map[visitKey]bool
}

type visitKey struct {
	t      types.Type
	encode bool
}

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	tags, err := scan(dir, os.Stderr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "reflectscan:", err)
		os.Exit(1)
	}
	fmt.Println(strings.Join(tags, ","))
}

// scan analyzes the packages in dir and returns the pk_no<kind> tags.
func scan(dir string, log io.Writer) ([]string, error) {
	pkgs, err := goList(dir)
	if err != nil {
		return nil, err
	}
	exports := map[string]string{}
	for _, p := range pkgs {
		if p.Error != nil {
			return nil, fmt.Errorf("%s: %s", p.ImportPath, p.Error.Err)
		}
		exports[p.ImportPath] = p.Export
	}

	s := &scanner{fset: token.NewFileSet(), used: map[string]string{}, visited: map[visitKey]bool{}}
	imp := importer.ForCompiler(s.fset, "gc", func(path string) (io.ReadCloser, error) {
		if exports[path] == "" {
			return nil, fmt.Errorf("no export data for %s", path)
		}
		return os.Open(exports[path])
	})
	for _, p := range pkgs {
		if !p.DepOnly {
			if err := s.checkPackage(p, imp); err != nil {
				return nil, err
			}
		}
	}

	if s.keepAll != "" {
		fmt.Fprintln(log, "reflectscan: keeping all kinds:", s.keepAll)
		return nil, nil
	}
	var tags []string
	for _, k := range kinds {
		if where, ok := s.used[k]; ok {
			fmt.Fprintf(log, "reflectscan: %s needed: %s\n", k, where)
		} else {
			tags = append(tags, "pk_no"+k)
		}
	}
	return tags, nil
}

// goList lists the packages in dir and their dependencies as tinygo sees
// them, with export data for type checking.
func goList(dir string) ([]listedPackage, error) {
	cmd := exec.Command("go", "list", "-e", "-json", "-export", "-deps", "-tags=tinygo", "./...")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("go list: %w", err)
	}
	var pkgs []listedPackage
	dec := json.NewDecoder(strings.NewReader(string(out)))
	for {
		var p listedPackage
		if err := dec.Decode(&p); errors.Is(err, io.EOF) {
			return pkgs, nil
		} else if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, p)
	}
}

func (s *scanner) checkPackage(p listedPackage, imp types.Importer) error {
	var files []*ast.File
	for _, name := range p.GoFiles {
		f, err := parser.ParseFile(s.fset, filepath.Join(p.Dir, name), nil, 0)
		if err != nil {
			return err
		}
		files = append(files, f)
	}
	info := &types.Info{
		Types:     map[ast.Expr]types.TypeAndValue{},
		Uses:      map[*ast.Ident]types.Object{},
		Instances: map[*ast.Ident]types.Instance{},
	}
	if _, err := (&types.Config{Importer: imp}).Check(p.ImportPath, s.fset, files, info); err != nil {
		return err
	}

	// Decode roots: type arguments of Decode/Get/Post/..., in source order so
	// the reported positions are stable.
	var ids []*ast.Ident
	for id := range info.Instances {
		if decodeFuncs[frameworkFunc(info.Uses[id])] {
			ids = append(ids, id)
		}
	}
	slices.SortFunc(ids, func(a, b *ast.Ident) int { return cmp.Compare(a.Pos(), b.Pos()) })
	for _, id := range ids {
		s.root(info.Instances[id].TypeArgs.At(0), false, id.Pos())
	}

	// Encode roots: arguments of Encode/Post/Put/Patch and FetchOptions.Body.
	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.CallExpr:
				if i, ok := encodeArgs[frameworkFunc(info.Uses[calleeIdent(n.Fun)])]; ok && i < len(n.Args) {
					s.encodeExpr(info, n.Args[i])
				}
			case *ast.CompositeLit:
				if isFetchOptions(info.TypeOf(n)) {
					for _, e := range n.Elts {
						if kv, ok := e.(*ast.KeyValueExpr); ok {
							if key, ok := kv.Key.(*ast.Ident); ok && key.Name == "Body" {
								s.encodeExpr(info, kv.Value)
							}
						}
					}
				}
			case *ast.AssignStmt:
				for i, lhs := range n.Lhs {
					if sel, ok := lhs.(*ast.SelectorExpr); ok && sel.Sel.Name == "Body" && isFetchOptions(info.TypeOf(sel.X)) && i < len(n.Rhs) {
						s.encodeExpr(info, n.Rhs[i])
					}
				}
			}
			return true
		})
	}
	return nil
}

func (s *scanner) encodeExpr(info *types.Info, e ast.Expr) {
	if tv := info.Types[e]; !tv.IsNil() && tv.Type != nil {
		s.root(tv.Type, true, e.Pos())
	}
}

// root records the kinds reachable from a type passed to Decode (as *T) or
// Encode. Decode/Encode check JSDecoder/JSEncoder only at the top level.
func (s *scanner) root(t types.Type, encode bool, pos token.Pos) {
	if encode && hasMethod(t, "ToJS") || !encode && hasMethod(types.NewPointer(t), "FromJS") {
		return
	}
	s.walk(t, encode, s.fset.Position(pos).String())
}

func (s *scanner) walk(t types.Type, encode bool, where string) {
	t = types.Unalias(t)
	key := visitKey{t, encode}
	if s.visited[key] {
		return
	}
	s.visited[key] = true
	at := where + " " + t.String()

	switch u := t.Underlying().(type) {
	case *types.Basic:
		switch {
		case u.Info()&types.IsString != 0:
			s.use("string", at)
		case u.Info()&types.IsBoolean != 0:
			s.use("bool", at)
		case u.Info()&types.IsFloat != 0:
			s.use("float", at)
		case u.Info()&types.IsUnsigned != 0 && u.Kind() != types.Uintptr:
			s.use("uint", at)
		case u.Info()&types.IsInteger != 0:
			s.use("int", at)
		}
	case *types.Pointer:
		s.use("ptr", at)
		s.walk(u.Elem(), encode, where)
	case *types.Slice:
		s.use("slice", at)
		s.walk(u.Elem(), encode, where)
	case *types.Array:
		if encode { // decodeValue has no array case
			s.use("slice", at)
			s.walk(u.Elem(), encode, where)
		}
	case *types.Map:
		s.use("map", at)
		s.walk(u.Elem(), encode, where)
	case *types.Struct:
		for i := range u.NumFields() {
			if f := u.Field(i); f.Exported() {
				s.walk(f.Type(), encode, where)
			}
		}
	case *types.Interface:
		s.use("iface", at)
		if encode { // the dynamic value can be of any kind
			s.keepAll = "encoded value of interface type at " + at
		}
	case *types.TypeParam:
		s.keepAll = "type parameter at " + at
	}
}

func (s *scanner) use(kind, where string) {
	if _, ok := s.used[kind]; !ok {
		s.used[kind] = where
	}
}

func frameworkFunc(obj types.Object) string {
	if f, ok := obj.(*types.Func); ok && f.Pkg() != nil && f.Pkg().Path() == pkgPath {
		return f.Name()
	}
	return ""
}

func calleeIdent(e ast.Expr) *ast.Ident {
	switch e := e.(type) {
	case *ast.Ident:
		return e
	case *ast.SelectorExpr:
		return e.Sel
	case *ast.IndexExpr: // explicit instantiation, e.g. p.Post[T]
		return calleeIdent(e.X)
	case *ast.IndexListExpr:
		return calleeIdent(e.X)
	}
	return nil
}

func isFetchOptions(t types.Type) bool {
	if p, ok := t.(*types.Pointer); ok {
		t = p.Elem()
	}
	n, ok := types.Unalias(t).(*types.Named)
	return ok && n.Obj().Pkg() != nil && n.Obj().Pkg().Path() == pkgPath && n.Obj().Name() == "FetchOptions"
}

func hasMethod(t types.Type, name string) bool {
	ms := types.NewMethodSet(t)
	for i := range ms.Len() {
		if ms.At(i).Obj().Name() == name {
			return true
		}
	}
	return false
}
