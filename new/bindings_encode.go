//go:build !js || !wasm

package preveltekit

import (
	"github.com/tdewolff/minify/v2"
	mcss "github.com/tdewolff/minify/v2/css"
	mhtml "github.com/tdewolff/minify/v2/html"
)

// Binary encoder for CollectedBindings.
// Format: positional, length-prefixed. No field names.
// Strings: varint length + raw bytes.
// Bools: single byte 0/1.
// Slices: varint count + items.
// Nullable pointers: byte 0=nil, 1=present + data.

var m *minify.M

func init() {
	m = minify.New()
	m.AddFunc("text/html", mhtml.Minify)
	m.AddFunc("text/css", mcss.Minify)
}

func minifyHTML(s string) string {
	out, err := m.String("text/html", s)
	if err != nil {
		return s
	}
	return out
}

func minifyCSS(s string) string {
	out, err := m.String("text/css", s)
	if err != nil {
		return s
	}
	return out
}

func encodeBindings(b *CollectedBindings) []byte {
	e := encoder{buf: make([]byte, 0, 1024)}
	e.writeBindings(b)
	return e.buf
}

type encoder struct {
	buf []byte
}

func (e *encoder) writeVarint(v int) {
	u := uint64(v)
	for u >= 0x80 {
		e.buf = append(e.buf, byte(u)|0x80)
		u >>= 7
	}
	e.buf = append(e.buf, byte(u))
}

func (e *encoder) writeString(s string) {
	e.writeVarint(len(s))
	e.buf = append(e.buf, s...)
}

// writeHTML writes a minified HTML string (collapse whitespace between tags).
func (e *encoder) writeHTML(s string) {
	e.writeString(minifyHTML(s))
}

func (e *encoder) writeBool(b bool) {
	if b {
		e.buf = append(e.buf, 1)
	} else {
		e.buf = append(e.buf, 0)
	}
}

func (e *encoder) writeBindings(b *CollectedBindings) {
	// TextBindings
	e.writeVarint(len(b.TextBindings))
	for _, tb := range b.TextBindings {
		e.writeString(tb.StoreID)
		e.writeString(tb.MarkerID)
		e.writeBool(tb.IsHTML)
	}

	// Events
	e.writeVarint(len(b.Events))
	for _, ev := range b.Events {
		e.writeString(ev.ElementID)
		e.writeString(ev.Event)
	}

	// IfBlocks
	e.writeVarint(len(b.IfBlocks))
	for _, ifb := range b.IfBlocks {
		e.writeString(ifb.MarkerID)

		// Branches
		e.writeVarint(len(ifb.Branches))
		for _, br := range ifb.Branches {
			e.writeHTML(br.HTML)
			e.writeString(br.StoreID)
			e.writeString(br.Op)
			e.writeString(br.Operand)
			e.writeBool(br.IsBool)
			// Nested bindings
			if br.Bindings != nil {
				e.writeBool(true)
				e.writeBindings(br.Bindings)
			} else {
				e.writeBool(false)
			}
		}

		e.writeHTML(ifb.ElseHTML)
		// ElseBindings
		if ifb.ElseBindings != nil {
			e.writeBool(true)
			e.writeBindings(ifb.ElseBindings)
		} else {
			e.writeBool(false)
		}

		// Deps
		e.writeVarint(len(ifb.Deps))
		for _, d := range ifb.Deps {
			e.writeString(d)
		}
	}

	// EachBlocks
	e.writeVarint(len(b.EachBlocks))
	for _, eb := range b.EachBlocks {
		e.writeString(eb.MarkerID)
		e.writeString(eb.ListID)
	}

	// InputBindings
	e.writeVarint(len(b.InputBindings))
	for _, ib := range b.InputBindings {
		e.writeString(ib.StoreID)
		e.writeString(ib.BindType)
	}

	// AttrBindings
	e.writeVarint(len(b.AttrBindings))
	for _, ab := range b.AttrBindings {
		e.writeString(ab.ElementID)
		e.writeString(ab.AttrName)
		e.writeString(ab.Template)
		e.writeVarint(len(ab.StoreIDs))
		for _, sid := range ab.StoreIDs {
			e.writeString(sid)
		}
	}

	// AttrCondBindings
	e.writeVarint(len(b.AttrCondBindings))
	for _, acb := range b.AttrCondBindings {
		e.writeString(acb.ElementID)
		e.writeString(acb.AttrName)
		e.writeString(acb.TrueValue)
		e.writeString(acb.FalseValue)
		e.writeString(acb.TrueStoreID)
		e.writeString(acb.FalseStoreID)
		e.writeString(acb.Op)
		e.writeString(acb.Operand)
		e.writeBool(acb.IsBool)
		e.writeVarint(len(acb.Deps))
		for _, d := range acb.Deps {
			e.writeString(d)
		}
	}

	// ComponentBlocks
	e.writeVarint(len(b.ComponentBlocks))
	for _, cb := range b.ComponentBlocks {
		e.writeString(cb.MarkerID)
		e.writeString(cb.StoreID)
		e.writeVarint(len(cb.Branches))
		for _, br := range cb.Branches {
			e.writeString(br.Name)
			e.writeHTML(br.HTML)
			if br.Bindings != nil {
				e.writeBool(true)
				e.writeBindings(br.Bindings)
			} else {
				e.writeBool(false)
			}
		}
	}
}
