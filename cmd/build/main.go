package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed assets/*
var assets embed.FS

// Files that go into assets/ subdirectory
var assetFiles = []string{"index.html", "wasm_exec.js"}

// Files that go into project root (with permissions)
var rootFiles = map[string]os.FileMode{
	"build.sh":      0755,
	"dev.sh":        0755,
	"Dockerfile":    0644,
	"Caddyfile":     0644,
	"Caddyfile.dev": 0644,
}

const goMod = `module %s

go 1.25

require github.com/tbocek/preveltekit/v2 v2.0.0
`

const mainGo = `package main

import p "github.com/tbocek/preveltekit/v2"

type App struct {
	CurrentComponent *p.Store[p.Component]
}

func (a *App) New() p.Component {
	return &App{
		CurrentComponent: p.New[p.Component](&Hello{}),
	}
}

func (a *App) Routes() []p.Route {
	return []p.Route{
		{Path: "/", HTMLFile: "index.html", SSRPath: "/", Component: &Hello{}},
	}
}

func (a *App) OnMount() {
	router := p.NewRouter(a.CurrentComponent, a.Routes(), "app")
	router.Start()
}

func (a *App) Render() p.Node {
	return p.Html(` + "`" + `<main>` + "`" + `, a.CurrentComponent, ` + "`" + `</main>` + "`" + `)
}

type Hello struct{}

func (h *Hello) Render() p.Node {
	return p.Html(` + "`" + `<h1>Hello, World!</h1>` + "`" + `)
}

func main() {
	p.Hydrate(&App{})
}
`

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: preveltekit init <module-name>")
			fmt.Fprintln(os.Stderr, "  e.g. preveltekit init hello")
			os.Exit(1)
		}
		cmdInit(os.Args[2])
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: preveltekit <command>")
	fmt.Println("  init <module-name>  Create a new project")
}

func cmdInit(moduleName string) {
	// Create project directory
	if err := os.MkdirAll(moduleName, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create %s/: %v\n", moduleName, err)
		os.Exit(1)
	}

	// Create assets/ subdirectory
	assetsDir := filepath.Join(moduleName, "assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create assets/: %v\n", err)
		os.Exit(1)
	}

	// Copy asset files
	for _, name := range assetFiles {
		data, _ := assets.ReadFile(filepath.Join("assets", name))
		dest := filepath.Join(assetsDir, name)
		if err := os.WriteFile(dest, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", dest, err)
			os.Exit(1)
		}
		fmt.Printf("  Created %s\n", dest)
	}

	// Copy root files
	for name, perm := range rootFiles {
		data, _ := assets.ReadFile(filepath.Join("assets", name))
		dest := filepath.Join(moduleName, name)
		if err := os.WriteFile(dest, data, perm); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", dest, err)
			os.Exit(1)
		}
		fmt.Printf("  Created %s\n", dest)
	}

	// Write go.mod
	goModPath := filepath.Join(moduleName, "go.mod")
	if err := os.WriteFile(goModPath, []byte(fmt.Sprintf(goMod, moduleName)), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write go.mod: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  Created %s\n", goModPath)

	// Write main.go
	mainGoPath := filepath.Join(moduleName, "main.go")
	if err := os.WriteFile(mainGoPath, []byte(mainGo), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write main.go: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  Created %s\n", mainGoPath)

	fmt.Printf("\nDone! To get started:\n")
	fmt.Printf("  cd %s\n", moduleName)
	fmt.Printf("  go get github.com/tbocek/preveltekit/v2@latest\n")
	fmt.Printf("  ./build.sh\n")
}
