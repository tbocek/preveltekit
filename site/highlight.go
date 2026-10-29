package main

import (
	"strings"

	p "github.com/tbocek/preveltekit/v2"
)

// codeBlock renders Go source as a highlighted <pre><code> block. The
// highlighting runs in Go while rendering, so it needs no JavaScript and
// works for SSR and client-side navigation alike.
func codeBlock(src string) p.Node {
	return p.Pre(p.Code(p.RawHTML(highlightGo(src))))
}

// Go keywords and the builtin types and constants used in the examples,
// space-separated for a cheap lookup (smaller WASM than maps).
const (
	goKeywords = " break case chan const continue default defer else for func go if import interface map package range return select struct switch type var "
	goBuiltins = " any bool byte error float64 int int64 rune string uint true false nil "
)

// highlightGo returns src as HTML with spans around comments (c), strings
// (s), numbers (n), keywords (k), builtin types and constants (t) and
// called functions (f).
func highlightGo(src string) string {
	var b strings.Builder
	span := func(class, text string) {
		b.WriteString(`<span class="`)
		b.WriteString(class)
		b.WriteString(`">`)
		escapeTo(&b, text)
		b.WriteString("</span>")
	}
	for i := 0; i < len(src); {
		c := src[i]
		switch {
		case strings.HasPrefix(src[i:], "//"):
			end := strings.IndexByte(src[i:], '\n')
			if end < 0 {
				end = len(src) - i
			}
			span("c", src[i:i+end])
			i += end
		case c == '"' || c == '\'' || c == '`':
			j := i + 1
			for j < len(src) && src[j] != c {
				if src[j] == '\\' && c != '`' {
					j++ // skip the escaped character
				}
				j++
			}
			j = min(j+1, len(src))
			span("s", src[i:j])
			i = j
		case isDigit(c):
			j := i
			for j < len(src) && (isIdent(src[j]) || src[j] == '.') {
				j++
			}
			span("n", src[i:j])
			i = j
		case isIdent(c):
			j := i
			for j < len(src) && isIdent(src[j]) {
				j++
			}
			word := src[i:j]
			switch {
			case strings.Contains(goKeywords, " "+word+" "):
				span("k", word)
			case strings.Contains(goBuiltins, " "+word+" "):
				span("t", word)
			case j < len(src) && src[j] == '(':
				span("f", word)
			default:
				b.WriteString(word)
			}
			i = j
		default:
			escapeTo(&b, src[i:i+1])
			i++
		}
	}
	return b.String()
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

func isIdent(c byte) bool {
	return c == '_' || isDigit(c) || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

// escapeTo writes s to b with &, < and > escaped.
func escapeTo(b *strings.Builder, s string) {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		default:
			b.WriteByte(s[i])
		}
	}
}
