//go:build !wasm

package preveltekit

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Hydrate is the main entry point for declarative components.
// In SSR mode (native build), it generates static HTML files.
func Hydrate(app ComponentRoot) {
	// First pass: discover all SSR paths
	if hn, ok := app.(HasNew); ok {
		app = hn.New().(ComponentRoot)
	}

	var ssrPaths []Route
	for _, route := range app.Routes() {
		if route.SSRPath != "" {
			ssrPaths = append(ssrPaths, route)
		}
	}

	// Create output directory
	os.MkdirAll("dist", 0755)

	// Generate HTML for each SSR path with fresh state
	for _, route := range ssrPaths {
		// Reset global counters so each iteration starts from s0,
		// matching the single app.New() call in WASM.
		resetRegistries()

		// Set the SSR path before lifecycle methods
		SetSSRPath(route.SSRPath)

		// Create fresh app instance (apps without New() are reused as-is)
		var freshApp Component = app
		if hn, ok := app.(HasNew); ok {
			freshApp = hn.New()
		}

		// Call OnMount (creates router which reads path and sets component)
		if om, ok := freshApp.(HasOnMount); ok {
			om.OnMount()
		}

		// Render the full tree; the app's own styles are collected first and
		// its scope is set before rendering so all app HTML gets the class
		ctx := &renderContext{styles: &collectedStyles{scoped: map[string]string{}, global: map[string]string{}}}
		ctx.ScopeAttr = ctx.componentScope("app", freshApp)
		html := nodeToHTML(freshApp.Render(), ctx)

		// Build full HTML document
		fullHTML := buildHTMLDocument(minifyHTML(html), ctx.styles.global, ctx.styles.scoped)

		// Write HTML file
		htmlPath := filepath.Join("dist", route.HTMLFile)
		os.WriteFile(htmlPath, []byte(fullHTML), 0644)
		fmt.Fprintf(os.Stderr, "Generated: %s\n", htmlPath)
	}
}

// resetRegistries resets all global counters and registries to initial state.
// Called before each SSR iteration so IDs start from s0, matching WASM.
func resetRegistries() {
	storeCounter = 0
	handlerCounter = 0
	scopeCounter = 0
	storeRegistry = make(map[string]any)
	handlerRegistry = make(map[string]func())
	handlerModifiers = make(map[string][]string)
	scopeRegistry = make(map[string]string)
}

func buildHTMLDocument(body string, collectedGlobalStyles, collectedStyles map[string]string) string {
	// Global styles first (unscoped), then scoped styles, each sorted by component name
	var allStyles strings.Builder
	for _, styles := range []map[string]string{collectedGlobalStyles, collectedStyles} {
		for _, k := range slices.Sorted(maps.Keys(styles)) {
			allStyles.WriteString(styles[k])
			allStyles.WriteByte('\n')
		}
	}

	var styles string
	if allStyles.Len() > 0 {
		styles = "<style>" + minifyCSS(allStyles.String()) + "</style>"
	}

	// Read template from assets/index.html
	tmpl, err := os.ReadFile("assets/index.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading assets/index.html: %v\n", err)
		os.Exit(1)
	}

	result := string(tmpl)
	result = strings.Replace(result, "<!--styles-->", styles, 1)
	result = strings.Replace(result, "<!--body-->", body, 1)
	return result
}

// minifyHTML collapses whitespace between tags. <pre> and <textarea>
// elements are copied as-is, since whitespace is significant inside them.
func minifyHTML(s string) string {
	var result strings.Builder
	i := 0
	for i < len(s) {
		if tag := whitespaceTag(s[i:]); tag != "" {
			end := strings.Index(s[i:], "</"+tag+">")
			if end < 0 {
				end = len(s) - i
			}
			result.WriteString(s[i : i+end])
			i += end
			continue
		}
		if s[i] == '>' {
			result.WriteByte('>')
			i++
			// Skip whitespace until next <
			start := i
			for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
				i++
			}
			if i < len(s) && s[i] == '<' {
				// Whitespace between tags — skip it
			} else {
				// Whitespace before non-tag content — keep it
				result.WriteString(s[start:i])
			}
		} else {
			result.WriteByte(s[i])
			i++
		}
	}
	return strings.TrimSpace(result.String())
}

// whitespaceTag returns "pre" or "textarea" if s starts with such an
// opening tag, else "".
func whitespaceTag(s string) string {
	for _, tag := range []string{"pre", "textarea"} {
		if strings.HasPrefix(s, "<"+tag) && len(s) > len(tag)+1 && (s[len(tag)+1] == '>' || s[len(tag)+1] == ' ') {
			return tag
		}
	}
	return ""
}

// minifyCSS removes comments, collapses whitespace, and trims unnecessary characters.
func minifyCSS(s string) string {
	// Remove comments
	for {
		start := strings.Index(s, "/*")
		if start < 0 {
			break
		}
		end := strings.Index(s[start+2:], "*/")
		if end < 0 {
			break
		}
		s = s[:start] + s[start+2+end+2:]
	}
	// Collapse whitespace
	var b strings.Builder
	inSpace := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			inSpace = true
			continue
		}
		if inSpace {
			if c != '{' && c != '}' && c != ':' && c != ';' && c != ',' && c != '>' && c != '+' && c != '~' {
				if b.Len() > 0 {
					prev := b.String()[b.Len()-1]
					if prev != '{' && prev != '}' && prev != ':' && prev != ';' && prev != ',' && prev != '>' && prev != '+' && prev != '~' {
						b.WriteByte(' ')
					}
				}
			}
			inSpace = false
		}
		b.WriteByte(c)
	}
	return strings.TrimSpace(b.String())
}
