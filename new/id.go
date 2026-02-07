package preveltekit

import "strings"

// IDCounter tracks counters for generating unique IDs.
// Used by both SSR (BuildContext) and WASM (handler collection) to ensure
// ID generation stays synchronized between server-side rendering and hydration.
type IDCounter struct {
	Text   int    // Counter for text binding markers
	If     int    // Counter for if-block markers
	Each   int    // Counter for each-block markers
	Event  int    // Counter for event element IDs
	Bind   int    // Counter for input binding element IDs
	Class  int    // Counter for class/attr binding element IDs
	Attr   int    // Counter for dynamic attribute element IDs
	Comp   int    // Counter for component markers
	Route  int    // Counter for route-block markers
	Prefix string // Prefix for nested components (e.g., "basics", "components_comp0")
}

// --- Element ID generators (for HTML id="..." attributes) ---

// NextEventID returns the next element ID for event bindings.
// Used in: <button id="basics_ev0">
func (c *IDCounter) NextEventID() string {
	id := "ev" + itoa(c.Event)
	c.Event++
	return id
}

// NextBindID returns the next element ID for input bindings.
// Used in: <input id="basics_b0">
func (c *IDCounter) NextBindID() string {
	id := "b" + itoa(c.Bind)
	c.Bind++
	return id
}

// NextClassID returns the next element ID for class bindings.
// Used in: <div id="basics_cl0">
func (c *IDCounter) NextClassID() string {
	id := "cl" + itoa(c.Class)
	c.Class++
	return id
}

// NextAttrID returns the next element ID for attribute bindings.
// Used in: <div data-attrbind="basics_a0">
func (c *IDCounter) NextAttrID() string {
	id := "a" + itoa(c.Attr)
	c.Attr++
	return id
}

// --- Marker ID generators (for HTML comments <!--marker-->) ---

// NextTextMarker returns the next marker ID for text bindings.
// Used in: <!--basics_t0--> (comment marker for text insertion point)
func (c *IDCounter) NextTextMarker() string {
	id := "t" + itoa(c.Text)
	c.Text++
	return id
}

// NextIfMarker returns the next marker ID for if-blocks.
// Used in: <!--basics_i0--> (comment marker for if-block boundary)
func (c *IDCounter) NextIfMarker() string {
	id := "if" + itoa(c.If)
	c.If++
	return id
}

// NextEachMarker returns the next marker ID for each-blocks.
// Used in: <!--basics_e0--> (comment marker for each-block boundary)
func (c *IDCounter) NextEachMarker() string {
	id := "each" + itoa(c.Each)
	c.Each++
	return id
}

// NextCompMarker returns the next marker ID for nested components.
// Used internally for component prefixing (e.g., "comp0" in "components_comp0_t0")
func (c *IDCounter) NextCompMarker() string {
	id := "comp" + itoa(c.Comp)
	c.Comp++
	return id
}

// NextRouteMarker returns the next marker ID for route-blocks.
// Used in: <!--basics_r0--> (comment marker for route-block boundary)
func (c *IDCounter) NextRouteMarker() string {
	id := "route" + itoa(c.Route)
	c.Route++
	return id
}

// --- ID formatting functions ---

// FullElementID returns the full element ID with prefix for use in HTML id="..." attributes.
// Example: FullElementID("ev0") with prefix "basics" returns "basics_ev0"
func (c *IDCounter) FullElementID(localID string) string {
	if c.Prefix == "" {
		return localID
	}
	return c.Prefix + "_" + localID
}

// FullMarkerID returns the shortened marker ID for use in HTML comments.
// Example: FullMarkerID("t0") with prefix "components_comp3" returns "components_c3_t0"
// The marker parts (comp, if, each) are shortened but component names are preserved.
func (c *IDCounter) FullMarkerID(localID string) string {
	if c.Prefix == "" {
		return shortenMarkerPart(localID)
	}
	return shortenMarkerParts(c.Prefix) + "_" + shortenMarkerPart(localID)
}

// shortenMarkerParts shortens all marker parts in a prefixed ID.
// Example: "components_comp3" -> "components_c3"
func shortenMarkerParts(id string) string {
	parts := strings.Split(id, "_")
	for i, part := range parts {
		parts[i] = shortenMarkerPart(part)
	}
	return strings.Join(parts, "_")
}

// shortenMarkerPart shortens a single marker part if it matches a known pattern.
// Only shortens generated marker IDs (comp0, if0, each0), not component names.
// Example: "comp3" -> "c3", "if0" -> "i0", "components" -> "components" (unchanged)
func shortenMarkerPart(part string) string {
	// comp0 -> c0 (but not "components" which doesn't end in digits)
	if len(part) > 4 && part[:4] == "comp" && isDigits(part[4:]) {
		return "c" + part[4:]
	}
	// each0 -> e0
	if len(part) > 4 && part[:4] == "each" && isDigits(part[4:]) {
		return "e" + part[4:]
	}
	// route0 -> r0
	if len(part) > 5 && part[:5] == "route" && isDigits(part[5:]) {
		return "r" + part[5:]
	}
	// if0 -> i0
	if len(part) > 2 && part[:2] == "if" && isDigits(part[2:]) {
		return "i" + part[2:]
	}
	// ev0 -> v0 (for markers, though events typically use element IDs)
	if len(part) > 2 && part[:2] == "ev" && isDigits(part[2:]) {
		return "v" + part[2:]
	}
	// cl0 -> l0 (for markers, though classes typically use element IDs)
	if len(part) > 2 && part[:2] == "cl" && isDigits(part[2:]) {
		return "l" + part[2:]
	}
	return part
}

// isDigits returns true if s contains only ASCII digits.
func isDigits(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
