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

		// Render the full tree
		ctx := NewBuildContext()

		// Collect app global styles (unscoped)
		if hgs, ok := freshApp.(HasGlobalStyle); ok {
			if gs := hgs.GlobalStyle(); gs != "" {
				ctx.CollectedGlobalStyles["app"] = gs
			}
		}

		// Set app-level scope before rendering so all app HTML gets the class
		if hs, ok := freshApp.(HasStyle); ok {
			scopeAttr := GetOrCreateScope("app")
			ctx.ScopeAttr = scopeAttr
			ctx.CollectedStyles["app"] = scopeCSS(hs.Style(), scopeAttr)
		}

		html := nodeToHTML(freshApp.Render(), ctx)

		// Build full HTML document
		fullHTML := buildHTMLDocument(minifyHTML(html), ctx.CollectedGlobalStyles, ctx.CollectedStyles)

		// Write HTML file
		htmlPath := filepath.Join("dist", route.HTMLFile)
		os.WriteFile(htmlPath, []byte(fullHTML), 0644)
		fmt.Fprintf(os.Stderr, "Generated: %s\n", htmlPath)
	}
}

func buildHTMLDocument(body string, collectedGlobalStyles, collectedStyles map[string]string) string {
	// Global styles first (unscoped), then scoped styles, each sorted by component name
	var allStyles string
	for _, styles := range []map[string]string{collectedGlobalStyles, collectedStyles} {
		for _, k := range slices.Sorted(maps.Keys(styles)) {
			allStyles += styles[k] + "\n"
		}
	}

	var styles string
	if allStyles != "" {
		styles = "<style>" + minifyCSS(allStyles) + "</style>"
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
