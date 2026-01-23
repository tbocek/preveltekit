package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: reactivebuild <main-component.go> [child-component.go ...]")
		fmt.Println("       reactivebuild --assemble <project-dir>")
		os.Exit(1)
	}

	// Handle --assemble mode
	if os.Args[1] == "--assemble" {
		if len(os.Args) < 3 {
			fatal("--assemble requires project directory")
		}
		assemble(os.Args[2])
		return
	}

	mainComponentFile := os.Args[1]
	dir := filepath.Dir(mainComponentFile)
	buildDir := filepath.Join(dir, "build")
	distDir := filepath.Join(dir, "dist")

	// Create folders
	for _, d := range []string{buildDir, distDir} {
		if err := os.MkdirAll(d, 0755); err != nil {
			fatal("mkdir error: %v", err)
		}
	}

	// Parse child components
	// Track all source files for each component (to handle build-tag variants)
	childComponents := make(map[string]*component)
	childSourceFiles := make(map[string][]string) // component name -> list of source files
	for i := 2; i < len(os.Args); i++ {
		childFile := os.Args[i]
		childComp, err := parseComponent(childFile)
		if err != nil {
			fatal("parse child component %s: %v", childFile, err)
		}
		childComponents[childComp.name] = childComp
		childSourceFiles[childComp.name] = append(childSourceFiles[childComp.name], childFile)
	}

	// Parse main component
	mainComp, err := parseComponent(mainComponentFile)
	if err != nil {
		fatal("parse error: %v", err)
	}

	// Auto-discover components referenced in templates (including child templates)
	// Use a queue to recursively discover nested component references
	discovered := make(map[string]bool)
	queue := findReferencedComponents(mainComp.template)

	for len(queue) > 0 {
		compName := queue[0]
		queue = queue[1:]

		if discovered[compName] {
			continue
		}
		discovered[compName] = true

		if _, exists := childComponents[compName]; exists {
			// Already loaded, but check its template for more components
			childComp := childComponents[compName]
			nestedRefs := findReferencedComponents(childComp.template)
			queue = append(queue, nestedRefs...)
			continue
		}

		// Try to find the component file in the same directory
		possibleFiles := []string{
			filepath.Join(dir, strings.ToLower(compName)+".go"),
			filepath.Join(dir, compName+".go"),
		}
		for _, f := range possibleFiles {
			if _, err := os.Stat(f); err == nil {
				childComp, err := parseComponent(f)
				if err != nil {
					fatal("auto-discovered component %s: %v", f, err)
				}
				if childComp.name == compName {
					childComponents[compName] = childComp
					childSourceFiles[compName] = append(childSourceFiles[compName], f)

					// Also look for _stub.go variant (build-tag variants)
					stubFile := filepath.Join(dir, strings.ToLower(compName)+"_stub.go")
					if _, err := os.Stat(stubFile); err == nil {
						childSourceFiles[compName] = append(childSourceFiles[compName], stubFile)
					}

					// Also check this component's template for nested references
					nestedRefs := findReferencedComponents(childComp.template)
					queue = append(queue, nestedRefs...)
					break
				}
			}
		}
	}

	// Write go.mod that imports reactive package
	scriptDir := findScriptDir()
	goMod := fmt.Sprintf(`module app

go 1.21

require reactive v0.0.0

replace reactive => %s
`, scriptDir)
	writeFile(filepath.Join(buildDir, "go.mod"), goMod)

	// Parse template and generate wiring code
	tmpl, bindings := parseTemplate(mainComp.template)

	// Validate template bindings
	if err := validateBindings(mainComp, bindings); err != nil {
		fatal("%s", err)
	}

	// Validate component references exist
	for _, comp := range bindings.components {
		if _, ok := childComponents[comp.name]; !ok {
			available := make([]string, 0, len(childComponents))
			for name := range childComponents {
				available = append(available, name)
			}
			fatal("template error: <%s /> references unknown component\n\n  Available components: %v\n\n  Hint: Pass the component file as an argument: reactivebuild app.go %s.go",
				comp.name, available, strings.ToLower(comp.name))
		}
	}

	// Generate files
	mainGo := generateMain(mainComp, tmpl, bindings, childComponents)
	renderGo := generateRender(mainComp, tmpl, bindings, childComponents)

	// Write generated files to build/
	componentSrc := stripTemplateAndStyle(mainComp.source)
	writeFile(filepath.Join(buildDir, "component.go"), componentSrc)

	for name := range childComponents {
		sourceFiles := childSourceFiles[name]
		if len(sourceFiles) == 1 {
			// Single file - strip template/style and write
			src, _ := os.ReadFile(sourceFiles[0])
			childSrc := stripTemplateAndStyle(string(src))
			writeFile(filepath.Join(buildDir, strings.ToLower(name)+".go"), childSrc)
		} else {
			// Multiple files (build-tag variants) - write each with original basename
			for _, srcFile := range sourceFiles {
				src, _ := os.ReadFile(srcFile)
				childSrc := stripTemplateAndStyle(string(src))
				baseName := filepath.Base(srcFile)
				writeFile(filepath.Join(buildDir, baseName), childSrc)
			}
		}
	}

	writeFile(filepath.Join(buildDir, "main.go"), mainGo)
	writeFile(filepath.Join(buildDir, "render.go"), renderGo)

	// Run go mod tidy to create go.sum
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = buildDir
	if err := cmd.Run(); err != nil {
		fatal("go mod tidy: %v", err)
	}

	// Collect all styles (main + children)
	var allStyles strings.Builder
	allStyles.WriteString(mainComp.style)
	for _, child := range childComponents {
		if child.style != "" {
			allStyles.WriteString("\n")
			allStyles.WriteString(child.style)
		}
	}

	// Write styles.css to build/
	writeFile(filepath.Join(buildDir, "styles.css"), allStyles.String())

	fmt.Printf("Generated: build/\n")
}

func assemble(dir string) {
	assetsDir := filepath.Join(dir, "assets")
	buildDir := filepath.Join(dir, "build")
	distDir := filepath.Join(dir, "dist")
	scriptDir := findScriptDir()
	defaultAssetsDir := filepath.Join(scriptDir, "assets")

	// Read pre-rendered HTML
	prerendered, err := os.ReadFile(filepath.Join(distDir, "prerendered.html"))
	if err != nil {
		fatal("read prerendered.html: %v", err)
	}

	// Read styles from build
	styles, err := os.ReadFile(filepath.Join(buildDir, "styles.css"))
	if err != nil {
		fatal("read styles.css: %v", err)
	}

	// Read template from assets (fall back to default)
	indexPath := filepath.Join(assetsDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		indexPath = filepath.Join(defaultAssetsDir, "index.html")
	}
	template, err := os.ReadFile(indexPath)
	if err != nil {
		fatal("read index.html: %v", err)
	}

	// Inject styles and pre-rendered content
	finalHTML := string(template)
	finalHTML = strings.Replace(finalHTML, "</head>", "\t<style>\n"+string(styles)+"\t</style>\n</head>", 1)
	finalHTML = strings.Replace(finalHTML, `<div id="app"></div>`, `<div id="app">`+string(prerendered)+`</div>`, 1)

	writeFile(filepath.Join(distDir, "index.html"), finalHTML)
	os.Remove(filepath.Join(distDir, "prerendered.html"))

	fmt.Println("Assembled: dist/")
}
