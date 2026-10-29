package main

import (
	"embed"
	"fmt"
	"os"
)

//go:embed assets/*
var assets embed.FS

// Files to scaffold: embedded assets/<src> is written to dst.
var files = []struct {
	src, dst string
	perm     os.FileMode
}{
	{"index.html", "assets/index.html", 0644},
	{"wasm_exec.js", "assets/wasm_exec.js", 0644},
	{"build.sh", "build.sh", 0755},
	{"dev.sh", "dev.sh", 0755},
	{"Dockerfile", "Dockerfile", 0644},
	{"Dockerfile.dev", "Dockerfile.dev", 0644},
	{"Caddyfile", "Caddyfile", 0644},
	{"Caddyfile.dev", "Caddyfile.dev", 0644},
	{"main.go", "main.go", 0644},
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

	for _, f := range files {
		data, _ := assets.ReadFile("assets/" + f.src) // embed paths always use '/'
		if err := os.WriteFile(f.dst, data, f.perm); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", f.dst, err)
			os.Exit(1)
		}
		fmt.Printf("  Created %s\n", f.dst)
	}
}
