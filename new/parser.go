package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type Parsed struct {
	Script       string
	Template     string
	Style        string
	LocalGoFiles map[string]string
	Imports      []string
	Components   map[string]*Parsed
	ReactiveVars map[string]string // varName -> type
}

func parse(goPath string) (Parsed, error) {
	return parseWithDepth(goPath, 0)
}

func parseWithDepth(goPath string, depth int) (Parsed, error) {
	if depth > 10 {
		return Parsed{}, fmt.Errorf("component import depth exceeded")
	}

	content, err := os.ReadFile(goPath)
	if err != nil {
		return Parsed{}, err
	}

	result, err := parseGoContent(string(content))
	if err != nil {
		return Parsed{}, err
	}

	result.LocalGoFiles = make(map[string]string)
	result.Components = make(map[string]*Parsed)

	dir := filepath.Dir(goPath)
	baseName := filepath.Base(goPath)
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		name := e.Name()
		// Skip the main file, app.go (generated), and other .web.go files
		if !e.IsDir() && strings.HasSuffix(name, ".go") && 
		   !strings.HasSuffix(name, ".web.go") && 
		   name != "app.go" && name != baseName {
			goContent, err := os.ReadFile(filepath.Join(dir, name))
			if err == nil {
				result.LocalGoFiles[name] = string(goContent)
			}
		}
	}

	result.Imports = extractImportsFromGo(string(content))

	// Load component imports
	compImports := extractComponentImportsFromGo(string(content))
	for _, imp := range compImports {
		compPath := filepath.Join(dir, imp+".web.go")
		comp, err := parseWithDepth(compPath, depth+1)
		if err != nil {
			continue
		}
		compName := filepath.Base(imp)
		result.Components[compName] = &comp
	}

	return result, nil
}

func parseGoContent(content string) (Parsed, error) {
	var result Parsed
	result.ReactiveVars = make(map[string]string)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return result, err
	}

	// Find _template and _style consts
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.CONST {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Names) != 1 || len(vs.Values) != 1 {
				continue
			}
			name := vs.Names[0].Name
			lit, ok := vs.Values[0].(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				continue
			}
			// Remove backticks or quotes
			val := lit.Value
			if strings.HasPrefix(val, "`") && strings.HasSuffix(val, "`") {
				val = val[1 : len(val)-1]
			} else if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
				val = val[1 : len(val)-1]
			}

			switch name {
			case "_template":
				result.Template = val
			case "_style":
				result.Style = val
			}
		}
	}

	// Find reactive vars: var X type + func setX(v type) {} pairs
	vars := make(map[string]string)    // name -> type
	setters := make(map[string]string) // setX -> type from param

	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.VAR {
				for _, spec := range d.Specs {
					vs, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					for i, name := range vs.Names {
						varType := ""
						if vs.Type != nil {
							varType = exprToString(vs.Type)
						} else if i < len(vs.Values) {
							varType = inferTypeFromExpr(vs.Values[i])
						}
						vars[name.Name] = varType
					}
				}
			}
		case *ast.FuncDecl:
			name := d.Name.Name
			if strings.HasPrefix(name, "set") && len(name) > 3 {
				// Check if it has one param
				if d.Type.Params != nil && len(d.Type.Params.List) == 1 {
					param := d.Type.Params.List[0]
					paramType := exprToString(param.Type)
					setters[name] = paramType
				}
			}
		}
	}

	// Match vars with setters
	for varName, varType := range vars {
		setterName := "set" + capitalize(varName)
		if setterType, ok := setters[setterName]; ok {
			if varType == "" {
				varType = setterType
			}
			result.ReactiveVars[varName] = varType
		}
	}

	// Script is everything except _template and _style
	result.Script = removeTemplateAndStyle(content)

	return result, nil
}

func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	case *ast.MapType:
		return "map[" + exprToString(t.Key) + "]" + exprToString(t.Value)
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	default:
		return ""
	}
}

func inferTypeFromExpr(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.BasicLit:
		switch t.Kind {
		case token.INT:
			return "int"
		case token.FLOAT:
			return "float64"
		case token.STRING:
			return "string"
		}
	case *ast.Ident:
		if t.Name == "true" || t.Name == "false" {
			return "bool"
		}
	case *ast.CompositeLit:
		if t.Type != nil {
			return exprToString(t.Type)
		}
	}
	return ""
}

func removeTemplateAndStyle(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	skip := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "const _template") ||
			strings.HasPrefix(trimmed, "const _style") {
			skip = true
			if strings.Count(line, "`") == 2 {
				skip = false
				continue
			}
		}

		if skip {
			if strings.Contains(line, "`") {
				skip = false
			}
			continue
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func extractImportsFromGo(content string) []string {
	var imports []string

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ImportsOnly)
	if err != nil {
		return imports
	}

	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if !strings.HasPrefix(path, "./") {
			imports = append(imports, path)
		}
	}

	return imports
}

func extractComponentImportsFromGo(content string) []string {
	var imports []string

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ImportsOnly)
	if err != nil {
		return imports
	}

	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if strings.HasPrefix(path, "./") {
			imports = append(imports, strings.TrimPrefix(path, "./"))
		}
	}

	return imports
}