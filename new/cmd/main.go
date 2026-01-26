package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: preveltekit <main-component.go> [child-component.go ...]")
		fmt.Println("       preveltekit --assemble <project-dir>")
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

	// Collect all source files
	sourceFiles := []string{mainComponentFile}
	sourceFiles = append(sourceFiles, os.Args[2:]...)

	// Auto-discover additional component files in the same directory
	sourceFiles = autoDiscoverComponents(dir, sourceFiles)

	// Extract struct names from all source files
	allStructNames := make(map[string]bool)
	for _, file := range sourceFiles {
		names, err := parseComponentNames(file)
		if err != nil {
			fatal("parse %s: %v", file, err)
		}
		for _, name := range names {
			allStructNames[name] = true
		}
	}

	// Get main component name (first struct in main file)
	mainNames, err := parseComponentNames(mainComponentFile)
	if err != nil || len(mainNames) == 0 {
		fatal("no structs found in %s", mainComponentFile)
	}
	mainCompName := mainNames[0]

	// Child component names (all others)
	var childCompNames []string
	for name := range allStructNames {
		if name != mainCompName {
			childCompNames = append(childCompNames, name)
		}
	}

	// Write go.mod
	scriptDir := findScriptDir()
	goMod := fmt.Sprintf(`module app

go 1.21

require preveltekit v0.0.0

replace preveltekit => %s
`, scriptDir)
	writeFile(filepath.Join(buildDir, "go.mod"), goMod)

	// Copy source files to build/ (strip Template/Style method bodies for cleaner code)
	writtenFiles := make(map[string]bool)
	for _, srcFile := range sourceFiles {
		baseName := filepath.Base(srcFile)
		if writtenFiles[baseName] {
			continue
		}
		writtenFiles[baseName] = true

		src, err := os.ReadFile(srcFile)
		if err != nil {
			fatal("read %s: %v", srcFile, err)
		}
		writeFile(filepath.Join(buildDir, baseName), string(src))
	}

	// Run go mod tidy before generating reflect.go
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = buildDir
	if err := cmd.Run(); err != nil {
		fatal("go mod tidy: %v", err)
	}

	// Clean up old generated files before running reflect
	os.Remove(filepath.Join(buildDir, "main.go"))
	os.Remove(filepath.Join(buildDir, "render.go"))

	// Detect types used in preveltekit.Get[T], Post[T], etc. for codec codegen
	codecTypes := detectCodecTypes(sourceFiles)

	// Generate reflect.go directly in build/ (same package main as source files)
	reflectGo := generateReflect(mainCompName, childCompNames, codecTypes)
	writeFile(filepath.Join(buildDir, "reflect.go"), reflectGo)

	// Run reflect.go to get component metadata
	var stdout, stderr bytes.Buffer
	cmd = exec.Command("go", "run", "-tags", "!wasm", ".")
	cmd.Dir = buildDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fatal("reflect.go failed: %v\n%s", err, stderr.String())
	}

	// Clean up reflect.go
	os.Remove(filepath.Join(buildDir, "reflect.go"))

	// Parse reflection output
	components, routes, codecTypeInfos, err := parseReflectOutput(stdout.String())
	if err != nil {
		fatal("parse reflect output: %v", err)
	}

	if len(components) == 0 {
		fatal("no components found")
	}

	// Find main component and build child map
	var mainComp *component
	childComponents := make(map[string]*component)
	for _, comp := range components {
		if comp.name == mainCompName {
			mainComp = comp
		} else {
			childComponents[comp.name] = comp
		}
	}

	if mainComp == nil {
		fatal("main component %s not found", mainCompName)
	}

	// Default routes if none specified
	if len(routes) == 0 {
		routes = []staticRoute{{path: "/", htmlFile: "index.html"}}
	}

	// Write routes.txt for build.sh
	var routeLines []string
	for _, r := range routes {
		routeLines = append(routeLines, r.path+":"+r.htmlFile)
	}
	writeFile(filepath.Join(buildDir, "routes.txt"), strings.Join(routeLines, "\n")+"\n")

	// Parse template and generate wiring code
	tmpl, bindings := parseTemplate(mainComp.template)

	// Validate component references exist
	for _, comp := range bindings.components {
		if _, ok := childComponents[comp.name]; !ok {
			available := make([]string, 0, len(childComponents))
			for name := range childComponents {
				available = append(available, name)
			}
			fatal("template error: <%s /> references unknown component\n\n  Available: %v",
				comp.name, available)
		}
	}

	// Generate main.go (WASM) and render.go (SSR)
	mainGo := generateMain(mainComp, tmpl, bindings, childComponents, codecTypeInfos)
	renderGo := generateRender(mainComp, tmpl, bindings, childComponents)

	writeFile(filepath.Join(buildDir, "main.go"), mainGo)
	writeFile(filepath.Join(buildDir, "render.go"), renderGo)

	// Remove reflect.go (no longer needed)
	os.Remove(filepath.Join(buildDir, "reflect.go"))

	// Collect all styles
	var allStyles strings.Builder
	allStyles.WriteString(mainComp.style)
	for _, child := range childComponents {
		if child.style != "" {
			allStyles.WriteString("\n")
			allStyles.WriteString(child.style)
		}
	}
	writeFile(filepath.Join(buildDir, "styles.css"), allStyles.String())

	fmt.Printf("Generated: build/\n")
	fmt.Printf("Routes: %d\n", len(routes))
	for _, r := range routes {
		fmt.Printf("  %s -> %s\n", r.path, r.htmlFile)
	}
}

// autoDiscoverComponents finds additional component files in the directory
func autoDiscoverComponents(dir string, existingFiles []string) []string {
	existing := make(map[string]bool)
	for _, f := range existingFiles {
		existing[filepath.Base(f)] = true
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return existingFiles
	}

	result := existingFiles
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		if existing[name] {
			continue
		}
		// Skip test files and build artifacts
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		result = append(result, filepath.Join(dir, name))
		existing[name] = true
	}

	return result
}

func assemble(dir string) {
	assetsDir := filepath.Join(dir, "assets")
	buildDir := filepath.Join(dir, "build")
	distDir := filepath.Join(dir, "dist")
	scriptDir := findScriptDir()
	defaultAssetsDir := filepath.Join(scriptDir, "assets")

	// Read routes
	routesData, err := os.ReadFile(filepath.Join(buildDir, "routes.txt"))
	if err != nil {
		// Default to single route
		routesData = []byte("/:index.html")
	}

	routes := parseRoutesFile(string(routesData))

	// Read styles
	styles, err := os.ReadFile(filepath.Join(buildDir, "styles.css"))
	if err != nil {
		styles = []byte{}
	}

	// Read template
	indexPath := filepath.Join(assetsDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		indexPath = filepath.Join(defaultAssetsDir, "index.html")
	}
	template, err := os.ReadFile(indexPath)
	if err != nil {
		fatal("read index.html: %v", err)
	}

	// Assemble each route
	for _, route := range routes {
		prerenderedFile := filepath.Join(distDir, strings.TrimSuffix(route.htmlFile, ".html")+"_prerendered.html")
		prerendered, err := os.ReadFile(prerenderedFile)
		if err != nil {
			fatal("read %s: %v", prerenderedFile, err)
		}

		finalHTML := string(template)
		finalHTML = strings.Replace(finalHTML, "</head>", "<style>"+minifyCSS(string(styles))+"</style></head>", 1)
		finalHTML = strings.Replace(finalHTML, `<div id="app"></div>`, `<div id="app">`+string(prerendered)+`</div>`, 1)
		finalHTML = minifyHTML(finalHTML)

		writeFile(filepath.Join(distDir, route.htmlFile), finalHTML)
		os.Remove(prerenderedFile)
	}

	// Minify wasm_exec.js if it exists
	wasmExecPath := filepath.Join(distDir, "wasm_exec.js")
	if jsContent, err := os.ReadFile(wasmExecPath); err == nil {
		writeFile(wasmExecPath, minifyJS(string(jsContent)))
	}

	fmt.Printf("Assembled: %d HTML files\n", len(routes))
}

func parseRoutesFile(content string) []staticRoute {
	var routes []staticRoute
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			routes = append(routes, staticRoute{
				path:     parts[0],
				htmlFile: parts[1],
			})
		}
	}
	if len(routes) == 0 {
		routes = []staticRoute{{path: "/", htmlFile: "index.html"}}
	}
	return routes
}
