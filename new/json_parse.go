//go:build js && wasm

package preveltekit

// Minimal JSON parser for HydrateBindings - replaces encoding/json to reduce WASM size

type jsonParser struct {
	data []byte
	pos  int
}

func parseBindings(data string) *HydrateBindings {
	p := &jsonParser{data: []byte(data), pos: 0}
	return p.parseHydrateBindings()
}

func (p *jsonParser) parseHydrateBindings() *HydrateBindings {
	b := &HydrateBindings{}
	p.skipWS()
	if !p.consume('{') {
		return b
	}

	for {
		p.skipWS()
		if p.peek() == '}' {
			p.pos++
			break
		}

		key := p.parseString()
		p.skipWS()
		p.consume(':')
		p.skipWS()

		switch key {
		case "TextBindings":
			b.TextBindings = p.parseTextBindings()
		case "Events":
			b.Events = p.parseEvents()
		case "IfBlocks":
			b.IfBlocks = p.parseIfBlocks()
		case "InputBindings":
			b.InputBindings = p.parseInputBindings()
		case "ClassBindings":
			b.ClassBindings = p.parseClassBindings()
		case "ShowIfBindings":
			b.ShowIfBindings = p.parseShowIfBindings()
		case "AttrBindings":
			b.AttrBindings = p.parseAttrBindings()
		case "EachBlocks":
			b.EachBlocks = p.parseEachBlocks()
		case "Components":
			// Skip these - not needed for hydration yet
			p.skipValue()
		default:
			p.skipValue()
		}

		p.skipWS()
		if !p.consume(',') {
			p.skipWS()
			p.consume('}')
			break
		}
	}
	return b
}

func (p *jsonParser) parseTextBindings() []HydrateTextBinding {
	var result []HydrateTextBinding
	p.skipWS()
	if p.peek() == 'n' {
		p.pos += 4 // skip null
		return result
	}
	if !p.consume('[') {
		return result
	}

	for {
		p.skipWS()
		if p.peek() == ']' {
			p.pos++
			break
		}

		tb := HydrateTextBinding{}
		p.consume('{')
		for {
			p.skipWS()
			if p.peek() == '}' {
				p.pos++
				break
			}
			key := p.parseString()
			p.skipWS()
			p.consume(':')
			p.skipWS()

			switch key {
			case "marker_id":
				tb.MarkerID = p.parseString()
			case "store_id":
				tb.StoreID = p.parseString()
			case "is_html":
				tb.IsHTML = p.parseBool()
			default:
				p.skipValue()
			}

			p.skipWS()
			if !p.consume(',') {
				p.skipWS()
				p.consume('}')
				break
			}
		}
		result = append(result, tb)

		p.skipWS()
		if !p.consume(',') {
			p.skipWS()
			p.consume(']')
			break
		}
	}
	return result
}

func (p *jsonParser) parseEvents() []HydrateEvent {
	var result []HydrateEvent
	p.skipWS()
	if p.peek() == 'n' {
		p.pos += 4 // skip null
		return result
	}
	if !p.consume('[') {
		return result
	}

	for {
		p.skipWS()
		if p.peek() == ']' {
			p.pos++
			break
		}

		ev := HydrateEvent{}
		p.consume('{')
		for {
			p.skipWS()
			if p.peek() == '}' {
				p.pos++
				break
			}
			key := p.parseString()
			p.skipWS()
			p.consume(':')
			p.skipWS()

			switch key {
			case "ElementID":
				ev.ElementID = p.parseString()
			case "Event":
				ev.Event = p.parseString()
			case "HandlerID":
				ev.HandlerID = p.parseString()
			case "ArgsStr":
				ev.ArgsStr = p.parseString()
			case "Modifiers":
				ev.Modifiers = p.parseStringArray()
			default:
				p.skipValue()
			}

			p.skipWS()
			if !p.consume(',') {
				p.skipWS()
				p.consume('}')
				break
			}
		}
		result = append(result, ev)

		p.skipWS()
		if !p.consume(',') {
			p.skipWS()
			p.consume(']')
			break
		}
	}
	return result
}

func (p *jsonParser) parseIfBlocks() []HydrateIfBlock {
	var result []HydrateIfBlock
	p.skipWS()
	if p.peek() == 'n' {
		p.pos += 4 // skip null
		return result
	}
	if !p.consume('[') {
		return result
	}

	for {
		p.skipWS()
		if p.peek() == ']' {
			p.pos++
			break
		}

		ifb := p.parseIfBlock()
		result = append(result, ifb)

		p.skipWS()
		if !p.consume(',') {
			p.skipWS()
			p.consume(']')
			break
		}
	}
	return result
}

func (p *jsonParser) parseIfBlock() HydrateIfBlock {
	ifb := HydrateIfBlock{}
	p.consume('{')
	for {
		p.skipWS()
		if p.peek() == '}' {
			p.pos++
			break
		}
		key := p.parseString()
		p.skipWS()
		p.consume(':')
		p.skipWS()

		switch key {
		case "MarkerID":
			ifb.MarkerID = p.parseString()
		case "Branches":
			ifb.Branches = p.parseIfBranches()
		case "ElseHTML":
			ifb.ElseHTML = p.parseString()
		case "ElseBindings":
			if p.peek() == 'n' {
				p.skipValue() // null
			} else {
				ifb.ElseBindings = p.parseHydrateBindings()
			}
		case "Deps":
			ifb.Deps = p.parseStringArray()
		default:
			p.skipValue()
		}

		p.skipWS()
		if !p.consume(',') {
			p.skipWS()
			p.consume('}')
			break
		}
	}
	return ifb
}

func (p *jsonParser) parseIfBranches() []HydrateIfBranch {
	var result []HydrateIfBranch
	p.skipWS()
	if p.peek() == 'n' {
		p.pos += 4 // skip null
		return result
	}
	if !p.consume('[') {
		return result
	}

	for {
		p.skipWS()
		if p.peek() == ']' {
			p.pos++
			break
		}

		br := HydrateIfBranch{}
		p.consume('{')
		for {
			p.skipWS()
			if p.peek() == '}' {
				p.pos++
				break
			}
			key := p.parseString()
			p.skipWS()
			p.consume(':')
			p.skipWS()

			switch key {
			case "CondExpr", "cond_expr":
				br.CondExpr = p.parseString()
			case "HTML", "html":
				br.HTML = p.parseString()
			case "Bindings":
				if p.peek() == 'n' {
					p.skipValue() // null
				} else {
					br.Bindings = p.parseHydrateBindings()
				}
			case "store_id":
				br.StoreID = p.parseString()
			case "op":
				br.Op = p.parseString()
			case "operand":
				br.Operand = p.parseString()
			case "is_bool":
				br.IsBool = p.parseBool()
			default:
				p.skipValue()
			}

			p.skipWS()
			if !p.consume(',') {
				p.skipWS()
				p.consume('}')
				break
			}
		}
		result = append(result, br)

		p.skipWS()
		if !p.consume(',') {
			p.skipWS()
			p.consume(']')
			break
		}
	}
	return result
}

func (p *jsonParser) parseInputBindings() []HydrateInputBinding {
	var result []HydrateInputBinding
	p.skipWS()
	if p.peek() == 'n' {
		p.pos += 4 // skip null
		return result
	}
	if !p.consume('[') {
		return result
	}

	for {
		p.skipWS()
		if p.peek() == ']' {
			p.pos++
			break
		}

		ib := HydrateInputBinding{}
		p.consume('{')
		for {
			p.skipWS()
			if p.peek() == '}' {
				p.pos++
				break
			}
			key := p.parseString()
			p.skipWS()
			p.consume(':')
			p.skipWS()

			switch key {
			case "element_id":
				ib.ElementID = p.parseString()
			case "store_id":
				ib.StoreID = p.parseString()
			case "bind_type":
				ib.BindType = p.parseString()
			default:
				p.skipValue()
			}

			p.skipWS()
			if !p.consume(',') {
				p.skipWS()
				p.consume('}')
				break
			}
		}
		result = append(result, ib)

		p.skipWS()
		if !p.consume(',') {
			p.skipWS()
			p.consume(']')
			break
		}
	}
	return result
}

func (p *jsonParser) parseClassBindings() []HydrateClassBinding {
	var result []HydrateClassBinding
	p.skipWS()
	if p.peek() == 'n' {
		p.pos += 4 // skip null
		return result
	}
	if !p.consume('[') {
		return result
	}

	for {
		p.skipWS()
		if p.peek() == ']' {
			p.pos++
			break
		}

		cb := HydrateClassBinding{}
		p.consume('{')
		for {
			p.skipWS()
			if p.peek() == '}' {
				p.pos++
				break
			}
			key := p.parseString()
			p.skipWS()
			p.consume(':')
			p.skipWS()

			switch key {
			case "element_id":
				cb.ElementID = p.parseString()
			case "class_name":
				cb.ClassName = p.parseString()
			case "cond_expr":
				cb.CondExpr = p.parseString()
			case "deps":
				cb.Deps = p.parseStringArray()
			default:
				p.skipValue()
			}

			p.skipWS()
			if !p.consume(',') {
				p.skipWS()
				p.consume('}')
				break
			}
		}
		result = append(result, cb)

		p.skipWS()
		if !p.consume(',') {
			p.skipWS()
			p.consume(']')
			break
		}
	}
	return result
}

func (p *jsonParser) parseShowIfBindings() []HydrateShowIfBinding {
	var result []HydrateShowIfBinding
	p.skipWS()
	if p.peek() == 'n' {
		p.pos += 4 // skip null
		return result
	}
	if !p.consume('[') {
		return result
	}

	for {
		p.skipWS()
		if p.peek() == ']' {
			p.pos++
			break
		}

		sib := HydrateShowIfBinding{}
		p.consume('{')
		for {
			p.skipWS()
			if p.peek() == '}' {
				p.pos++
				break
			}
			key := p.parseString()
			p.skipWS()
			p.consume(':')
			p.skipWS()

			switch key {
			case "element_id":
				sib.ElementID = p.parseString()
			case "store_id":
				sib.StoreID = p.parseString()
			case "op":
				sib.Op = p.parseString()
			case "operand":
				sib.Operand = p.parseString()
			case "is_bool":
				sib.IsBool = p.parseBool()
			case "deps":
				sib.Deps = p.parseStringArray()
			default:
				p.skipValue()
			}

			p.skipWS()
			if !p.consume(',') {
				p.skipWS()
				p.consume('}')
				break
			}
		}
		result = append(result, sib)

		p.skipWS()
		if !p.consume(',') {
			p.skipWS()
			p.consume(']')
			break
		}
	}
	return result
}

func (p *jsonParser) parseAttrBindings() []HydrateAttrBinding {
	var result []HydrateAttrBinding
	p.skipWS()
	if p.peek() == 'n' {
		p.pos += 4 // skip null
		return result
	}
	if !p.consume('[') {
		return result
	}

	for {
		p.skipWS()
		if p.peek() == ']' {
			p.pos++
			break
		}

		ab := HydrateAttrBinding{}
		p.consume('{')
		for {
			p.skipWS()
			if p.peek() == '}' {
				p.pos++
				break
			}
			key := p.parseString()
			p.skipWS()
			p.consume(':')
			p.skipWS()

			switch key {
			case "element_id", "ElementID":
				ab.ElementID = p.parseString()
			case "attr_name", "AttrName":
				ab.AttrName = p.parseString()
			case "template", "Template":
				ab.Template = p.parseString()
			case "store_ids", "StoreIDs":
				ab.StoreIDs = p.parseStringArray()
			default:
				p.skipValue()
			}

			p.skipWS()
			if !p.consume(',') {
				p.skipWS()
				p.consume('}')
				break
			}
		}
		result = append(result, ab)

		p.skipWS()
		if !p.consume(',') {
			p.skipWS()
			p.consume(']')
			break
		}
	}
	return result
}

func (p *jsonParser) parseEachBlocks() []HydrateEachBlock {
	var result []HydrateEachBlock
	p.skipWS()
	if p.peek() == 'n' {
		p.pos += 4 // skip null
		return result
	}
	if !p.consume('[') {
		return result
	}

	for {
		p.skipWS()
		if p.peek() == ']' {
			p.pos++
			break
		}

		eb := HydrateEachBlock{}
		p.consume('{')
		for {
			p.skipWS()
			if p.peek() == '}' {
				p.pos++
				break
			}
			key := p.parseString()
			p.skipWS()
			p.consume(':')
			p.skipWS()

			switch key {
			case "MarkerID":
				eb.MarkerID = p.parseString()
			case "ListID":
				eb.ListID = p.parseString()
			case "ItemVar":
				eb.ItemVar = p.parseString()
			case "IndexVar":
				eb.IndexVar = p.parseString()
			default:
				p.skipValue()
			}

			p.skipWS()
			if !p.consume(',') {
				p.skipWS()
				p.consume('}')
				break
			}
		}
		result = append(result, eb)

		p.skipWS()
		if !p.consume(',') {
			p.skipWS()
			p.consume(']')
			break
		}
	}
	return result
}

// Helper methods

func (p *jsonParser) skipWS() {
	for p.pos < len(p.data) {
		c := p.data[p.pos]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			p.pos++
		} else {
			break
		}
	}
}

func (p *jsonParser) peek() byte {
	if p.pos >= len(p.data) {
		return 0
	}
	return p.data[p.pos]
}

func (p *jsonParser) consume(c byte) bool {
	if p.pos < len(p.data) && p.data[p.pos] == c {
		p.pos++
		return true
	}
	return false
}

func (p *jsonParser) parseString() string {
	p.skipWS()
	if !p.consume('"') {
		return ""
	}

	start := p.pos
	var result []byte

	for p.pos < len(p.data) {
		c := p.data[p.pos]
		if c == '"' {
			if result == nil {
				s := string(p.data[start:p.pos])
				p.pos++
				return s
			}
			p.pos++
			return string(result)
		}
		if c == '\\' {
			if result == nil {
				result = append(result, p.data[start:p.pos]...)
			}
			p.pos++
			if p.pos >= len(p.data) {
				break
			}
			escaped := p.data[p.pos]
			switch escaped {
			case '"', '\\', '/':
				result = append(result, escaped)
			case 'n':
				result = append(result, '\n')
			case 'r':
				result = append(result, '\r')
			case 't':
				result = append(result, '\t')
			case 'u':
				// Parse \uXXXX
				if p.pos+4 < len(p.data) {
					r := p.parseHex4()
					if r < 128 {
						result = append(result, byte(r))
					} else if r < 2048 {
						result = append(result, byte(0xC0|(r>>6)), byte(0x80|(r&0x3F)))
					} else {
						result = append(result, byte(0xE0|(r>>12)), byte(0x80|((r>>6)&0x3F)), byte(0x80|(r&0x3F)))
					}
					p.pos--
				}
			}
			p.pos++
		} else {
			if result != nil {
				result = append(result, c)
			}
			p.pos++
		}
	}
	return string(result)
}

func (p *jsonParser) parseHex4() rune {
	var r rune
	for i := 0; i < 4; i++ {
		p.pos++
		c := p.data[p.pos]
		r <<= 4
		if c >= '0' && c <= '9' {
			r |= rune(c - '0')
		} else if c >= 'a' && c <= 'f' {
			r |= rune(c - 'a' + 10)
		} else if c >= 'A' && c <= 'F' {
			r |= rune(c - 'A' + 10)
		}
	}
	return r
}

func (p *jsonParser) parseBool() bool {
	p.skipWS()
	if p.pos+4 <= len(p.data) && string(p.data[p.pos:p.pos+4]) == "true" {
		p.pos += 4
		return true
	}
	if p.pos+5 <= len(p.data) && string(p.data[p.pos:p.pos+5]) == "false" {
		p.pos += 5
		return false
	}
	return false
}

func (p *jsonParser) parseStringArray() []string {
	var result []string
	p.skipWS()
	// Handle null
	if p.peek() == 'n' {
		p.pos += 4 // skip "null"
		return result
	}
	if !p.consume('[') {
		return result
	}

	for {
		p.skipWS()
		if p.peek() == ']' {
			p.pos++
			break
		}

		result = append(result, p.parseString())

		p.skipWS()
		if !p.consume(',') {
			p.skipWS()
			p.consume(']')
			break
		}
	}
	return result
}

func (p *jsonParser) skipValue() {
	p.skipWS()
	c := p.peek()
	switch c {
	case '"':
		p.parseString()
	case '{':
		p.skipObject()
	case '[':
		p.skipArray()
	case 't':
		p.pos += 4 // true
	case 'f':
		p.pos += 5 // false
	case 'n':
		p.pos += 4 // null
	default:
		// number
		for p.pos < len(p.data) {
			c := p.data[p.pos]
			if c == ',' || c == '}' || c == ']' || c == ' ' || c == '\n' || c == '\r' || c == '\t' {
				break
			}
			p.pos++
		}
	}
}

func (p *jsonParser) skipObject() {
	p.consume('{')
	depth := 1
	for p.pos < len(p.data) && depth > 0 {
		c := p.data[p.pos]
		if c == '"' {
			p.parseString()
			continue
		}
		if c == '{' {
			depth++
		} else if c == '}' {
			depth--
		}
		p.pos++
	}
}

func (p *jsonParser) skipArray() {
	p.consume('[')
	depth := 1
	for p.pos < len(p.data) && depth > 0 {
		c := p.data[p.pos]
		if c == '"' {
			p.parseString()
			continue
		}
		if c == '[' {
			depth++
		} else if c == ']' {
			depth--
		}
		p.pos++
	}
}
