package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
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
	case "strip":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "Usage: preveltekit strip <wasm-file> <remove-dir>\n")
			os.Exit(1)
		}
		cmdStrip(os.Args[2], os.Args[3])
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: preveltekit <command>")
	fmt.Println("  init                          Copy build.sh and assets/ to current directory")
	fmt.Println("  strip <wasm-file> <remove-dir> Strip HTML strings from WASM binary")
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

// toWAT converts a raw byte string to its WAT data section representation.
// Printable ASCII (0x20-0x7e) stays literal, except " → \22 and \ → \5c.
// Everything else becomes \xx hex escapes.
func toWAT(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '"':
			b.WriteString(`\22`)
		case c == '\\':
			b.WriteString(`\5c`)
		case c >= 0x20 && c <= 0x7e:
			b.WriteByte(c)
		default:
			fmt.Fprintf(&b, `\%02x`, c)
		}
	}
	return b.String()
}

func cmdStrip(wasmFile, removeDir string) {
	entries, err := os.ReadDir(removeDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", removeDir, err)
		os.Exit(1)
	}
	if len(entries) == 0 {
		fmt.Println("  No strings to strip")
		return
	}

	// Convert wasm to wat
	watFile := wasmFile + ".wat"
	cmd := exec.Command("wasm2wat", wasmFile, "-o", watFile)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "wasm2wat failed: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(watFile)

	wat, err := os.ReadFile(watFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading WAT: %v\n", err)
		os.Exit(1)
	}

	// Split WAT into lines and separate data section lines from the rest.
	// Data lines start with "  (data $" or "  (data (" and contain the string literals.
	// We only search/replace within these lines to avoid false matches in code sections.
	lines := strings.Split(string(wat), "\n")
	type dataLine struct {
		index   int
		content string
	}
	var dataLines []dataLine
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "(data ") {
			dataLines = append(dataLines, dataLine{i, line})
		}
	}

	// Concatenate all data lines for searching (with a separator that won't appear in WAT)
	const sep = "\x00\x01\x02"
	dataContents := make([]string, len(dataLines))
	for i, dl := range dataLines {
		dataContents[i] = dl.content
	}
	dataStr := strings.Join(dataContents, sep)

	// Load all entries and parse them
	type stripEntry struct {
		str           string
		expectedCount int
	}
	var stripEntries []stripEntry
	for _, entry := range entries {
		raw, err := os.ReadFile(filepath.Join(removeDir, entry.Name()))
		if err != nil || len(raw) == 0 {
			continue
		}
		content := string(raw)
		newlineIdx := strings.IndexByte(content, '\n')
		if newlineIdx < 0 {
			continue
		}
		expectedCount := 0
		fmt.Sscanf(content[:newlineIdx], "%d", &expectedCount)
		str := content[newlineIdx+1:]
		if len(str) == 0 {
			continue
		}
		stripEntries = append(stripEntries, stripEntry{str, expectedCount})
	}

	// Sort by string length descending — process longest strings first so that
	// shorter substrings (e.g. "</p>") don't falsely match inside longer ones
	// (e.g. "</strong></p>\n\t\t\t") that haven't been nulled out yet.
	sort.Slice(stripEntries, func(i, j int) bool {
		return len(stripEntries[i].str) > len(stripEntries[j].str)
	})

	stripped := 0
	skipped := 0
	totalBytes := 0

	for _, se := range stripEntries {
		needle := toWAT(se.str)
		replacement := strings.Repeat(`\00`, len(se.str))

		dataCount := strings.Count(dataStr, needle)
		if dataCount == se.expectedCount {
			dataStr = strings.ReplaceAll(dataStr, needle, replacement)
			stripped++
			totalBytes += len(se.str) * dataCount
		} else if dataCount > 0 && dataCount < se.expectedCount {
			// Fewer in data section than expected — compiler deduplicated, replace all
			dataStr = strings.ReplaceAll(dataStr, needle, replacement)
			stripped++
			totalBytes += len(se.str) * dataCount
		} else {
			skipped++
			if dataCount > se.expectedCount {
				preview := strings.ReplaceAll(se.str, "\n", "\\n")
				preview = strings.ReplaceAll(preview, "\t", "\\t")
				if len(preview) > 60 {
					preview = preview[:60] + "..."
				}
				fmt.Fprintf(os.Stderr, "  Skip (data:%d > expected:%d): %s\n", dataCount, se.expectedCount, preview)
			}
		}
	}

	if stripped == 0 {
		fmt.Printf("  No strings stripped (%d skipped)\n", skipped)
		return
	}

	// Splice modified data lines back into the full WAT
	modifiedDataLines := strings.Split(dataStr, sep)
	for i, dl := range dataLines {
		lines[dl.index] = modifiedDataLines[i]
	}

	// Write modified WAT and convert back
	if err := os.WriteFile(watFile, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing WAT: %v\n", err)
		os.Exit(1)
	}

	cmd = exec.Command("wat2wasm", watFile, "-o", wasmFile)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "wat2wasm failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("  Stripped %d strings (%d bytes), skipped %d\n", stripped, totalBytes, skipped)
}
