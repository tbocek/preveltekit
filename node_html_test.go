//go:build !wasm

package preveltekit

import "testing"

type todo struct {
	Title string
	Done  bool
}

// Each works for any comparable item type, not only string and int.
func TestEachStructList(t *testing.T) {
	list := NewList(todo{"a", false}, todo{"b", true})
	each := Each(list, func(item todo, i int) Node {
		return Li(item.Title)
	}).Else(P("empty"))

	if got, want := nodeToHTML(each, NewBuildContext()), "<!--e0s--><li>a</li><li>b</li><!--e0-->"; got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}

	list.Clear()
	if got, want := nodeToHTML(each, NewBuildContext()), "<!--e0s--><p>empty</p><!--e0-->"; got != want {
		t.Errorf("empty list: got %q, want %q", got, want)
	}
}
