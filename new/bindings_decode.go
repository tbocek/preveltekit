//go:build js && wasm

package preveltekit

// Binary decoder for HydrateBindings.
// Mirrors the encoder in bindings_encode.go — same field order, same format.

func decodeBindings(data []byte) *HydrateBindings {
	d := decoder{data: data}
	return d.readBindings()
}

type decoder struct {
	data []byte
	pos  int
}

func (d *decoder) readVarint() int {
	var v uint64
	var shift uint
	for d.pos < len(d.data) {
		b := d.data[d.pos]
		d.pos++
		v |= uint64(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}
	return int(v)
}

func (d *decoder) readString() string {
	n := d.readVarint()
	if n == 0 {
		return ""
	}
	s := string(d.data[d.pos : d.pos+n])
	d.pos += n
	return s
}

func (d *decoder) readBool() bool {
	if d.pos >= len(d.data) {
		return false
	}
	v := d.data[d.pos]
	d.pos++
	return v != 0
}

func (d *decoder) readBindings() *HydrateBindings {
	b := &HydrateBindings{}

	// TextBindings
	n := d.readVarint()
	if n > 0 {
		b.TextBindings = make([]HydrateTextBinding, n)
		for i := range b.TextBindings {
			b.TextBindings[i].StoreID = d.readString()
			b.TextBindings[i].MarkerID = d.readString()
			b.TextBindings[i].IsHTML = d.readBool()
		}
	}

	// Events
	n = d.readVarint()
	if n > 0 {
		b.Events = make([]HydrateEvent, n)
		for i := range b.Events {
			b.Events[i].ElementID = d.readString()
			b.Events[i].Event = d.readString()
		}
	}

	// IfBlocks
	n = d.readVarint()
	if n > 0 {
		b.IfBlocks = make([]HydrateIfBlock, n)
		for i := range b.IfBlocks {
			b.IfBlocks[i].MarkerID = d.readString()

			// Branches
			nb := d.readVarint()
			if nb > 0 {
				b.IfBlocks[i].Branches = make([]HydrateIfBranch, nb)
				for j := range b.IfBlocks[i].Branches {
					b.IfBlocks[i].Branches[j].HTML = d.readString()
					b.IfBlocks[i].Branches[j].StoreID = d.readString()
					b.IfBlocks[i].Branches[j].Op = d.readString()
					b.IfBlocks[i].Branches[j].Operand = d.readString()
					b.IfBlocks[i].Branches[j].IsBool = d.readBool()
					if d.readBool() {
						b.IfBlocks[i].Branches[j].Bindings = d.readBindings()
					}
				}
			}

			b.IfBlocks[i].ElseHTML = d.readString()
			if d.readBool() {
				b.IfBlocks[i].ElseBindings = d.readBindings()
			}

			// Deps
			nd := d.readVarint()
			if nd > 0 {
				b.IfBlocks[i].Deps = make([]string, nd)
				for j := range b.IfBlocks[i].Deps {
					b.IfBlocks[i].Deps[j] = d.readString()
				}
			}
		}
	}

	// EachBlocks
	n = d.readVarint()
	if n > 0 {
		b.EachBlocks = make([]HydrateEachBlock, n)
		for i := range b.EachBlocks {
			b.EachBlocks[i].MarkerID = d.readString()
			b.EachBlocks[i].ListID = d.readString()
			b.EachBlocks[i].BodyHTML = d.readString()
		}
	}

	// InputBindings
	n = d.readVarint()
	if n > 0 {
		b.InputBindings = make([]HydrateInputBinding, n)
		for i := range b.InputBindings {
			b.InputBindings[i].StoreID = d.readString()
			b.InputBindings[i].BindType = d.readString()
		}
	}

	// AttrBindings
	n = d.readVarint()
	if n > 0 {
		b.AttrBindings = make([]HydrateAttrBinding, n)
		for i := range b.AttrBindings {
			b.AttrBindings[i].ElementID = d.readString()
			b.AttrBindings[i].AttrName = d.readString()
			b.AttrBindings[i].Template = d.readString()
			ns := d.readVarint()
			if ns > 0 {
				b.AttrBindings[i].StoreIDs = make([]string, ns)
				for j := range b.AttrBindings[i].StoreIDs {
					b.AttrBindings[i].StoreIDs[j] = d.readString()
				}
			}
		}
	}

	// AttrCondBindings
	n = d.readVarint()
	if n > 0 {
		b.AttrCondBindings = make([]HydrateAttrCondBinding, n)
		for i := range b.AttrCondBindings {
			b.AttrCondBindings[i].ElementID = d.readString()
			b.AttrCondBindings[i].AttrName = d.readString()
			b.AttrCondBindings[i].TrueValue = d.readString()
			b.AttrCondBindings[i].FalseValue = d.readString()
			b.AttrCondBindings[i].TrueStoreID = d.readString()
			b.AttrCondBindings[i].FalseStoreID = d.readString()
			b.AttrCondBindings[i].Op = d.readString()
			b.AttrCondBindings[i].Operand = d.readString()
			b.AttrCondBindings[i].IsBool = d.readBool()
			nd := d.readVarint()
			if nd > 0 {
				b.AttrCondBindings[i].Deps = make([]string, nd)
				for j := range b.AttrCondBindings[i].Deps {
					b.AttrCondBindings[i].Deps[j] = d.readString()
				}
			}
		}
	}

	// ComponentBlocks
	n = d.readVarint()
	if n > 0 {
		b.ComponentBlocks = make([]HydrateComponentBlock, n)
		for i := range b.ComponentBlocks {
			b.ComponentBlocks[i].MarkerID = d.readString()
			b.ComponentBlocks[i].StoreID = d.readString()
			nb := d.readVarint()
			if nb > 0 {
				b.ComponentBlocks[i].Branches = make([]HydrateComponentBranch, nb)
				for j := range b.ComponentBlocks[i].Branches {
					b.ComponentBlocks[i].Branches[j].Name = d.readString()
					b.ComponentBlocks[i].Branches[j].HTML = d.readString()
					if d.readBool() {
						b.ComponentBlocks[i].Branches[j].Bindings = d.readBindings()
					}
				}
			}
		}
	}

	return b
}
