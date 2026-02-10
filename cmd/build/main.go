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
	"build.sh":       0755,
	"dev.sh":         0755,
	"Dockerfile":     0644,
	"Dockerfile.dev": 0644,
	"Caddyfile":      0644,
	"Caddyfile.dev":  0644,
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		cmdInit()
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: preveltekit <command>")
	fmt.Println("  init  Scaffold a new project in the current directory")
}

func cmdInit() {
	// Check that go.mod exists
	if _, err := os.Stat("go.mod"); err != nil {
		fmt.Fprintln(os.Stderr, "Error: go.mod not found. Run 'go mod init <module-name>' first.")
		os.Exit(1)
	}

	// Create assets/ subdirectory
	if err := os.MkdirAll("assets", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create assets/: %v\n", err)
		os.Exit(1)
	}

	// Copy asset files
	for _, name := range assetFiles {
		data, _ := assets.ReadFile(filepath.Join("assets", name))
		dest := filepath.Join("assets", name)
		if err := os.WriteFile(dest, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", dest, err)
			os.Exit(1)
		}
		fmt.Printf("  Created %s\n", dest)
	}

	// Copy root files
	for name, perm := range rootFiles {
		data, _ := assets.ReadFile(filepath.Join("assets", name))
		if err := os.WriteFile(name, data, perm); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", name, err)
			os.Exit(1)
		}
		fmt.Printf("  Created %s\n", name)
	}

	// Write main.go
	mainGoData, _ := assets.ReadFile("assets/main.go")
	if err := os.WriteFile("main.go", mainGoData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write main.go: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("  Created main.go")

	fmt.Println("\nDone! Next steps:")
	fmt.Println("  go get github.com/tbocek/preveltekit/v2@latest")
	fmt.Println("  ./build.sh")
}
