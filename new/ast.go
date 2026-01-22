package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type Variable struct {
	Name      string
	DependsOn []string
}

type Function struct {
	Name     string
	Modifies []string
}

type Analysis struct {
	Vars  map[string]*Variable
	Funcs map[string]*Function
	Order []string
}

func analyze(script string) (*Analysis, error) {
	// Wrap as valid Go file
	wrapped := "package main\n\n" + convertShortDecls(script)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", wrapped, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	a := &Analysis{
		Vars:  make(map[string]*Variable),
		Funcs: make(map[string]*Function),
	}

	// First pass: find all top-level vars
	for _, decl := range file.Decls {
		if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.VAR {
			for _, spec := range gd.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range vs.Names {
						a.Vars[name.Name] = &Variable{
							Name:      name.Name,
							DependsOn: []string{},
						}
					}
				}
			}
		}
	}

	// Second pass: find dependencies and functions
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.VAR {
				for _, spec := range d.Specs {
					if vs, ok := spec.(*ast.ValueSpec); ok {
						for i, name := range vs.Names {
							if i < len(vs.Values) {
								a.Vars[name.Name].DependsOn = findDeps(vs.Values[i], a.Vars)
							}
						}
					}
				}
			}
		case *ast.FuncDecl:
			f := &Function{
				Name:     d.Name.Name,
				Modifies: findModifies(d.Body, a.Vars),
			}
			a.Funcs[d.Name.Name] = f
		}
	}

	a.Order = topoSort(a.Vars)
	return a, nil
}

// convertShortDecls converts "x := 0" to "var x = 0" at top level
func convertShortDecls(script string) string {
	var result []string
	for _, line := range strings.Split(script, "\n") {
		trimmed := strings.TrimSpace(line)
		// Only convert if it looks like a top-level short decl (not in a func)
		if idx := strings.Index(trimmed, " := "); idx > 0 {
			name := trimmed[:idx]
			if isIdent(name) {
				expr := trimmed[idx+4:]
				result = append(result, "var "+name+" = "+expr)
				continue
			}
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

func isIdent(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i, c := range s {
		if i == 0 {
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_') {
				return false
			}
		} else {
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
				return false
			}
		}
	}
	return true
}

func findDeps(expr ast.Expr, known map[string]*Variable) []string {
	var deps []string
	seen := make(map[string]bool)

	ast.Inspect(expr, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			if _, exists := known[id.Name]; exists && !seen[id.Name] {
				deps = append(deps, id.Name)
				seen[id.Name] = true
			}
		}
		return true
	})
	return deps
}

func findModifies(body *ast.BlockStmt, known map[string]*Variable) []string {
	var mods []string
	seen := make(map[string]bool)

	if body == nil {
		return mods
	}

	ast.Inspect(body, func(n ast.Node) bool {
		switch s := n.(type) {
		case *ast.AssignStmt:
			for _, lhs := range s.Lhs {
				if id, ok := lhs.(*ast.Ident); ok {
					if _, exists := known[id.Name]; exists && !seen[id.Name] {
						mods = append(mods, id.Name)
						seen[id.Name] = true
					}
				}
			}
		case *ast.IncDecStmt:
			if id, ok := s.X.(*ast.Ident); ok {
				if _, exists := known[id.Name]; exists && !seen[id.Name] {
					mods = append(mods, id.Name)
					seen[id.Name] = true
				}
			}
		}
		return true
	})
	return mods
}

func topoSort(vars map[string]*Variable) []string {
	var order []string
	visited := make(map[string]bool)
	inStack := make(map[string]bool)

	var visit func(name string)
	visit = func(name string) {
		if inStack[name] || visited[name] {
			return
		}
		inStack[name] = true
		if v, ok := vars[name]; ok {
			for _, dep := range v.DependsOn {
				visit(dep)
			}
		}
		inStack[name] = false
		visited[name] = true
		order = append(order, name)
	}

	for name := range vars {
		visit(name)
	}
	return order
}