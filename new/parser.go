package main

import (
	"fmt"
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
	Components   map[string]*Parsed // name -> parsed component
}

type ComponentUsage struct {
	Name     string            // "Button"
	ID       string            // "button0"
	Props    map[string]string // prop name -> expression
	Children string            // inner HTML
}

func parse(gowebPath string) (Parsed, error) {
	return parseWithDepth(gowebPath, 0)
}

func parseWithDepth(gowebPath string, depth int) (Parsed, error) {
	if depth > 10 {
		return Parsed{}, fmt.Errorf("component import depth exceeded")
	}
	
	content, err := os.ReadFile(gowebPath)
	if err != nil {
		return Parsed{}, err
	}

	result := parseContent(string(content))
	result.LocalGoFiles = make(map[string]string)
	result.Components = make(map[string]*Parsed)

	dir := filepath.Dir(gowebPath)
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		name := e.Name()
		if !e.IsDir() && strings.HasSuffix(name, ".go") && name != "app.go" {
			goContent, err := os.ReadFile(filepath.Join(dir, name))
			if err == nil {
				result.LocalGoFiles[name] = string(goContent)
			}
		}
	}

	result.Imports = extractImports(result.Script)
	
	// Load component imports recursively
	compImports := extractComponentImports(result.Script)
	for _, imp := range compImports {
		compPath := filepath.Join(dir, imp+".goweb")
		comp, err := parseWithDepth(compPath, depth+1)
		if err != nil {
			continue
		}
		compName := filepath.Base(imp)
		result.Components[compName] = &comp
	}
	
	return result, nil
}

func parseContent(content string) Parsed {
	lines := strings.Split(content, "\n")
	var section string
	var buf strings.Builder
	var result Parsed

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		switch trimmed {
		case "<script>":
			section = "script"
		case "</script>":
			result.Script = strings.TrimSpace(buf.String())
			buf.Reset()
			section = ""
		case "<template>":
			section = "template"
		case "</template>":
			result.Template = strings.TrimSpace(buf.String())
			buf.Reset()
			section = ""
		case "<style>":
			section = "style"
		case "</style>":
			result.Style = strings.TrimSpace(buf.String())
			buf.Reset()
			section = ""
		default:
			if section != "" {
				buf.WriteString(line)
				buf.WriteString("\n")
			}
		}
	}
	return result
}

func extractImports(script string) []string {
	var imports []string
	lines := strings.Split(script, "\n")
	inImportBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "import (" {
			inImportBlock = true
			continue
		}
		if inImportBlock && trimmed == ")" {
			inImportBlock = false
			continue
		}
		if inImportBlock {
			pkg := strings.Trim(trimmed, `"`)
			if pkg != "" && !strings.HasPrefix(pkg, "./") {
				imports = append(imports, pkg)
			}
		}
		if strings.HasPrefix(trimmed, "import \"") {
			pkg := strings.TrimPrefix(trimmed, "import ")
			pkg = strings.Trim(pkg, `"`)
			if !strings.HasPrefix(pkg, "./") {
				imports = append(imports, pkg)
			}
		}
	}
	return imports
}

func extractComponentImports(script string) []string {
	var imports []string
	lines := strings.Split(script, "\n")
	inImportBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "import (" {
			inImportBlock = true
			continue
		}
		if inImportBlock && trimmed == ")" {
			inImportBlock = false
			continue
		}
		if inImportBlock {
			pkg := strings.Trim(trimmed, `"`)
			if strings.HasPrefix(pkg, "./") {
				imports = append(imports, strings.TrimPrefix(pkg, "./"))
			}
		}
		if strings.HasPrefix(trimmed, "import \"./") {
			pkg := strings.TrimPrefix(trimmed, "import \"./")
			pkg = strings.TrimSuffix(pkg, "\"")
			imports = append(imports, pkg)
		}
	}
	return imports
}