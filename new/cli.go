package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: pregokit <file.goweb>")
		os.Exit(1)
	}

	gowebFile := os.Args[1]
	dir := filepath.Dir(gowebFile)
	distDir := filepath.Join(dir, "dist")

	// Create dist folder
	if err := os.MkdirAll(distDir, 0755); err != nil {
		fmt.Printf("mkdir error: %v\n", err)
		os.Exit(1)
	}

	p, err := parse(gowebFile)
	if err != nil {
		fmt.Printf("parse error: %v\n", err)
		os.Exit(1)
	}

	a, err := analyze(p.Script)
	if err != nil {
		fmt.Printf("analyze error: %v\n", err)
		os.Exit(1)
	}

	goCode, html := generate(p, a)

	if err := os.WriteFile(filepath.Join(distDir, "app.go"), []byte(goCode), 0644); err != nil {
		fmt.Printf("write error: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(filepath.Join(distDir, "index.html"), []byte(html), 0644); err != nil {
		fmt.Printf("write error: %v\n", err)
		os.Exit(1)
	}

	// Copy local .go files to dist
	for name, content := range p.LocalGoFiles {
		if err := os.WriteFile(filepath.Join(distDir, name), []byte(content), 0644); err != nil {
			fmt.Printf("write error: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Generated: %s/\n", distDir)
	fmt.Println("  app.go, index.html")
	fmt.Printf("Build: cd %s && tinygo build -o app.wasm -target wasm -no-debug -panic trap -scheduler none -gc leaking .\n", distDir)
}