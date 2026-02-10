package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed assets/*
var assets embed.FS

//go:embed build.sh
var buildScript []byte

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
	fmt.Println("  init  Copy build.sh and assets/ to current directory")
}

func cmdInit() {
	if err := os.MkdirAll("assets", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create assets/: %v\n", err)
		os.Exit(1)
	}

	entries, _ := assets.ReadDir("assets")
	for _, e := range entries {
		data, _ := assets.ReadFile(filepath.Join("assets", e.Name()))
		dest := filepath.Join("assets", e.Name())
		if err := os.WriteFile(dest, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", dest, err)
			os.Exit(1)
		}
		fmt.Printf("  Created %s\n", dest)
	}

	if err := os.WriteFile("build.sh", buildScript, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write build.sh: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("  Created build.sh")
	fmt.Println("\nDone! Run ./build.sh to build your project.")
}
