package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type Variable struct {
	Name      string
	Type      string
	DependsOn []string // vars this depends on (for derived values)
	Affects   []string // vars that depend on this
}

type Function struct {
	Name     string
	Modifies []string // vars modified by this function
}

type Analysis struct {
	Vars  map[string]*Variable
	Funcs map[string]*Function
	Order []string // topological order for derived values
}

func analyze(script string) (*Analysis, error) {
	wrapped := "package main\n\n" + script

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
						varType := ""
						if vs.Type != nil {
							varType = typeToString(vs.Type)
						}
						a.Vars[name.Name] = &Variable{
							Name:      name.Name,
							Type:      varType,
							DependsOn: []string{},
							Affects:   []string{},
						}
					}
				}
			}
		}
	}

	// Second pass: find dependencies
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.VAR {
				for _, spec := range d.Specs {
					if vs, ok := spec.(*ast.ValueSpec); ok {
						for i, name := range vs.Names {
							if i < len(vs.Values) {
								deps := findDeps(vs.Values[i], a.Vars)
								a.Vars[name.Name].DependsOn = deps
								// Update reverse mapping
								for _, dep := range deps {
									if v, ok := a.Vars[dep]; ok {
										v.Affects = append(v.Affects, name.Name)
									}
								}
							}
						}
					}
				}
			}
		case *ast.FuncDecl:
			// Skip stub setters
			if isStubSetter(d) {
				continue
			}
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

func isStubSetter(fn *ast.FuncDecl) bool {
	name := fn.Name.Name
	if len(name) <= 3 || name[:3] != "set" {
		return false
	}
	// Check for empty body or just panic
	if fn.Body == nil || len(fn.Body.List) == 0 {
		return true
	}
	if len(fn.Body.List) == 1 {
		if expr, ok := fn.Body.List[0].(*ast.ExprStmt); ok {
			if call, ok := expr.X.(*ast.CallExpr); ok {
				if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "panic" {
					return true
				}
			}
		}
	}
	return false
}

func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.ArrayType:
		return "[]" + typeToString(t.Elt)
	case *ast.MapType:
		return "map[" + typeToString(t.Key) + "]" + typeToString(t.Value)
	default:
		return ""
	}
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
				switch x := lhs.(type) {
				case *ast.Ident:
					if _, exists := known[x.Name]; exists && !seen[x.Name] {
						mods = append(mods, x.Name)
						seen[x.Name] = true
					}
				case *ast.IndexExpr:
					// items[i] = v -> modifies items
					if id, ok := x.X.(*ast.Ident); ok {
						if _, exists := known[id.Name]; exists && !seen[id.Name] {
							mods = append(mods, id.Name)
							seen[id.Name] = true
						}
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
		case *ast.ExprStmt:
			// Check for delete(map, key) and clear(x) calls
			if call, ok := s.X.(*ast.CallExpr); ok {
				if fn, ok := call.Fun.(*ast.Ident); ok {
					if fn.Name == "delete" && len(call.Args) >= 1 {
						if id, ok := call.Args[0].(*ast.Ident); ok {
							if _, exists := known[id.Name]; exists && !seen[id.Name] {
								mods = append(mods, id.Name)
								seen[id.Name] = true
							}
						}
					}
					if fn.Name == "clear" && len(call.Args) >= 1 {
						if id, ok := call.Args[0].(*ast.Ident); ok {
							if _, exists := known[id.Name]; exists && !seen[id.Name] {
								mods = append(mods, id.Name)
								seen[id.Name] = true
							}
						}
					}
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

// GetTransitiveDeps returns all vars affected when varName changes
func (a *Analysis) GetTransitiveDeps(varName string) []string {
	affected := make(map[string]bool)
	var collect func(name string)
	collect = func(name string) {
		if affected[name] {
			return
		}
		affected[name] = true
		if v, ok := a.Vars[name]; ok {
			for _, dep := range v.Affects {
				collect(dep)
			}
		}
	}
	collect(varName)
	
	// Return in topological order
	var result []string
	for _, name := range a.Order {
		if affected[name] {
			result = append(result, name)
		}
	}
	return result
}