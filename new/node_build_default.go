//go:build !js || !wasm

package preveltekit

// Text creates a text node.
func Text(content string) *TextNode {
	return &TextNode{Content: content}
}

// Html creates a raw HTML node from strings and embedded nodes.
// Example: Html(`<div class="foo">`, p.Bind(store), `</div>`)
func Html(parts ...any) *HtmlNode {
	return &HtmlNode{Parts: parts}
}
