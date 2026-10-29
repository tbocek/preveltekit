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

	if got, want := nodeToHTML(each, &renderContext{}), "<!--e0s--><li>a</li><li>b</li><!--e0-->"; got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}

	list.Clear()
	if got, want := nodeToHTML(each, &renderContext{}), "<!--e0s--><p>empty</p><!--e0-->"; got != want {
		t.Errorf("empty list: got %q, want %q", got, want)
	}
}

type slotWrapper struct{}

func (w *slotWrapper) Render() Node {
	return Div(If(Cond(func() bool { return true }), Slot()), Slot())
}

// Slot content is also rendered inside if-branches of the component.
func TestSlotInsideIf(t *testing.T) {
	got := nodeToHTML(Comp(&slotWrapper{}, Span("x")), &renderContext{})
	want := "<div><!--c0_i0s--><span>x</span><!--c0_i0--><span>x</span></div>"
	if got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}

type hello struct{}

func (h *hello) Render() Node { return H1("hi") }

// A Store[Component] always gets a block marker, also without WithOptions
// and when empty, so WASM can find and update it. In SSR the block notes
// which component it shows (<!--=name-->).
func TestStoreComponentMarker(t *testing.T) {
	tests := []struct {
		name  string
		store *Store[Component]
		want  string
	}{
		{"no options", New[Component](&hello{}), "<div><!--r0s--><!--=hello--><h1>hi</h1><!--r0--></div>"},
		{"empty", New[Component](nil), "<div><!--r0s--><!--=--><!--r0--></div>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nodeToHTML(Div(tt.store), &renderContext{}); got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// Each item has its own ID prefix, so IDs after the list don't depend on
// the items.
func TestEachItemIDs(t *testing.T) {
	suffix := New("!")
	list := NewList("a", "b")
	got := nodeToHTML(Div(Each(list, func(s string, i int) Node {
		return Span(s, Bind(suffix))
	}), Bind(suffix)), &renderContext{})
	want := "<div><!--e0s-->" +
		"<span>a<!--e0_0_t0s-->!<!--e0_0_t0--></span>" +
		"<span>b<!--e0_1_t0s-->!<!--e0_1_t0--></span>" +
		"<!--e0--><!--t0s-->!<!--t0--></div>"
	if got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}
