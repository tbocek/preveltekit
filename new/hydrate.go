//go:build !js || !wasm

package preveltekit

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Hydrate is the main entry point for declarative components.
// In SSR mode (native build), it generates static HTML files.
// WASM discovers all bindings by walking the Render() tree directly,
// so no bindings.bin is needed.
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

		// Create fresh app instance
		var freshApp Component
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
	var allStyles string

	// Global styles first (unscoped)
	if len(collectedGlobalStyles) > 0 {
		keys := make([]string, 0, len(collectedGlobalStyles))
		for k := range collectedGlobalStyles {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			allStyles += collectedGlobalStyles[k] + "\n"
		}
	}

	// Scoped styles
	if len(collectedStyles) > 0 {
		keys := make([]string, 0, len(collectedStyles))
		for k := range collectedStyles {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			allStyles += collectedStyles[k] + "\n"
		}
	}

	var styles string
	if allStyles != "" {
		styles = "<style>" + minifyCSS(allStyles) + "</style>\n"
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
%s</head>
<body>
%s
<script src="wasm_exec.js"></script>
<script>
const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
  .then(result => go.run(result.instance));
</script>
</body>
</html>`, styles, body)
}
