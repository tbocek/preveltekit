package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

// parseComponentNames extracts struct names from a Go file.
// These are candidates for components - we'll use reflection to determine
// which ones have Template() methods and extract their full metadata.
func parseComponentNames(file string) ([]string, error) {
	src, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, src, 0)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if _, ok := typeSpec.Type.(*ast.StructType); ok {
				names = append(names, typeSpec.Name.Name)
			}
		}
	}

	return names, nil
}

// Legacy functions for backward compatibility during transition
// TODO: Remove these once the new reflection-based flow is fully working

// parseComponents parses all components from a Go file.
// A component is a struct with a Template() method.
func parseComponents(file string) ([]*component, error) {
	src, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// First pass: collect all structs
	structs := make(map[string]*component)
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if _, ok := typeSpec.Type.(*ast.StructType); !ok {
				continue
			}

			comp := &component{
				name:   typeSpec.Name.Name,
				source: string(src),
			}
			structs[comp.name] = comp
		}
	}

	// Second pass: collect methods and associate with structs
	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Recv == nil {
			continue
		}

		var recvTypeName string
		for _, recv := range funcDecl.Recv.List {
			switch t := recv.Type.(type) {
			case *ast.StarExpr:
				if ident, ok := t.X.(*ast.Ident); ok {
					recvTypeName = ident.Name
				}
			case *ast.Ident:
				recvTypeName = t.Name
			}
		}

		if recvTypeName == "" {
			continue
		}

		comp, exists := structs[recvTypeName]
		if !exists {
			continue
		}

		methodName := funcDecl.Name.Name
		comp.methods = append(comp.methods, methodName)

		switch methodName {
		case "Template":
			comp.template = extractStringReturn(funcDecl)
		case "Style":
			comp.style = extractStringReturn(funcDecl)
		case "OnMount":
			comp.hasOnMount = true
		case "OnCreate":
			comp.hasOnCreate = true
		case "OnUnmount":
			comp.hasOnUnmount = true
		}
	}

	// Filter to only include structs that have a Template() method
	var components []*component
	for _, comp := range structs {
		if comp.template != "" {
			components = append(components, comp)
		}
	}

	return components, nil
}

// parseComponent parses a single component from a Go file (first one found).
func parseComponent(file string) (*component, error) {
	components, err := parseComponents(file)
	if err != nil {
		return nil, err
	}
	if len(components) == 0 {
		return nil, nil
	}
	return components[0], nil
}

func extractStringReturn(fn *ast.FuncDecl) string {
	for _, stmt := range fn.Body.List {
		retStmt, ok := stmt.(*ast.ReturnStmt)
		if !ok || len(retStmt.Results) == 0 {
			continue
		}
		return evalStringExpr(retStmt.Results[0])
	}
	return ""
}

func evalStringExpr(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return ""
		}
		s := e.Value
		if len(s) >= 2 {
			if s[0] == '`' {
				return s[1 : len(s)-1]
			}
			// Handle escape sequences in regular strings
			return s[1 : len(s)-1]
		}
		return s
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			return evalStringExpr(e.X) + evalStringExpr(e.Y)
		}
	case *ast.ParenExpr:
		return evalStringExpr(e.X)
	}
	return ""
}
