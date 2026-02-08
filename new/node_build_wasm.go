//go:build js && wasm

package preveltekit

// Text creates a text node. In WASM builds, content is discarded
// since WASM never uses HTML strings — the DOM already has SSR content.
func Text(content string) *TextNode {
	return &TextNode{}
}

// Html creates a raw HTML node. In WASM builds, string parts are stripped
// since WASM never uses HTML strings. Non-string parts (Nodes, Stores)
// are preserved so walkNodeForComponents can find nested components.
func Html(parts ...any) *HtmlNode {
	filtered := make([]any, 0, len(parts))
	for _, p := range parts {
		if _, ok := p.(string); !ok {
			filtered = append(filtered, p)
		}
	}
	return &HtmlNode{Parts: filtered}
}
