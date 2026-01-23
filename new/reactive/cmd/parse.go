package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
)

func parseComponent(file string) (*component, error) {
	src, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	comp := &component{source: string(src)}

	// Find struct and its fields
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
			comp.name = typeSpec.Name.Name

			// Parse struct fields for *Store[T], *List[T], *Map[K,V]
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

				if indexExpr, ok := starExpr.X.(*ast.IndexExpr); ok {
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
		}
	}

	// Find methods
	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Recv == nil {
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

	if comp.name == "" {
		return nil, fmt.Errorf("no struct found")
	}

	return comp, nil
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
