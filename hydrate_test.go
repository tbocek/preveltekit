//go:build !wasm

package preveltekit

import (
	"os"
	"strings"
	"testing"
)

type rootWithoutNew struct{}

func (r *rootWithoutNew) Render() Node { return Div("hello") }
func (r *rootWithoutNew) Routes() []Route {
	return []Route{{Path: "/", HTMLFile: "index.html", SSRPath: "/"}}
}

// An app without New() is rendered as-is instead of crashing.
func TestHydrateWithoutNew(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.MkdirAll("assets", 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("assets/index.html", []byte("<!--styles--><!--body-->"), 0644); err != nil {
		t.Fatal(err)
	}

	Hydrate(&rootWithoutNew{})

	html, err := os.ReadFile("dist/index.html")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(html), "<div>hello</div>") {
		t.Errorf("dist/index.html = %q, want it to contain <div>hello</div>", html)
	}
}

func TestMinifyHTML(t *testing.T) {
	in := "<div>\n  <p>a</p>\n  <pre><span>x</span>\n\n  <span>y</span></pre>\n<textarea> keep\n </textarea> <b>z</b></div>"
	want := "<div><p>a</p><pre><span>x</span>\n\n  <span>y</span></pre><textarea> keep\n </textarea><b>z</b></div>"
	if got := minifyHTML(in); got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}
