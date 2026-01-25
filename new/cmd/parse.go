package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
)

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

	// First pass: collect all structs and their fields
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
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			comp := &component{
				name:   typeSpec.Name.Name,
				source: string(src),
			}

			// Parse struct fields for *Store[T], *List[T], *Map[K,V], *LocalStore
			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue
				}
				fieldName := field.Names[0].Name

				starExpr, ok := field.Type.(*ast.StarExpr)
				if !ok {
					continue
				}

				var storeType, valueType, keyType string

				// Check for *LocalStore (non-generic, always string)
				if ident, ok := starExpr.X.(*ast.Ident); ok && ident.Name == "LocalStore" {
					storeType = "LocalStore"
					valueType = "string"
				} else if sel, ok := starExpr.X.(*ast.SelectorExpr); ok && sel.Sel.Name == "LocalStore" {
					storeType = "LocalStore"
					valueType = "string"
				} else if indexExpr, ok := starExpr.X.(*ast.IndexExpr); ok {
					switch x := indexExpr.X.(type) {
					case *ast.Ident:
						storeType = x.Name
					case *ast.SelectorExpr:
						storeType = x.Sel.Name
					default:
						continue
					}

					if storeType != "Store" && storeType != "List" {
						continue
					}

					switch t := indexExpr.Index.(type) {
					case *ast.Ident:
						valueType = t.Name
					default:
						valueType = "any"
					}
				} else if indexListExpr, ok := starExpr.X.(*ast.IndexListExpr); ok {
					switch x := indexListExpr.X.(type) {
					case *ast.Ident:
						storeType = x.Name
					case *ast.SelectorExpr:
						storeType = x.Sel.Name
					default:
						continue
					}

					if storeType != "Map" {
						continue
					}

					if len(indexListExpr.Indices) >= 2 {
						if kt, ok := indexListExpr.Indices[0].(*ast.Ident); ok {
							keyType = kt.Name
						}
						if vt, ok := indexListExpr.Indices[1].(*ast.Ident); ok {
							valueType = vt.Name
						}
					}
				} else {
					continue
				}

				comp.fields = append(comp.fields, storeField{
					name:      fieldName,
					storeType: storeType,
					valueType: valueType,
					keyType:   keyType,
				})
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

		// Get receiver type name
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
		}
	}

	// Filter to only include structs that have a Template() method (i.e., are components)
	var components []*component
	for _, comp := range structs {
		if comp.template != "" {
			components = append(components, comp)
		}
	}

	return components, nil
}

// parseComponent parses a single component from a Go file (first one found).
// Kept for backward compatibility.
func parseComponent(file string) (*component, error) {
	components, err := parseComponents(file)
	if err != nil {
		return nil, err
	}
	if len(components) == 0 {
		return nil, fmt.Errorf("no component found (struct with Template() method)")
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

// evalStringExpr recursively evaluates a string expression (handles concatenation)
func evalStringExpr(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return ""
		}
		s := e.Value
		if len(s) >= 2 {
			// Handle both quoted strings and raw strings
			if s[0] == '`' {
				return s[1 : len(s)-1]
			}
			// For regular strings, use strconv.Unquote to handle escapes
			if unquoted, err := strconv.Unquote(s); err == nil {
				return unquoted
			}
			return s[1 : len(s)-1]
		}
		return s
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			// String concatenation
			left := evalStringExpr(e.X)
			right := evalStringExpr(e.Y)
			return left + right
		}
	case *ast.ParenExpr:
		return evalStringExpr(e.X)
	}
	return ""
}
