package preveltekit

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
	id := "i" + itoa(c.If)
	c.If++
	return id
}

// NextEachMarker returns the next marker ID for each-blocks.
// Used in: <!--basics_e0--> (comment marker for each-block boundary)
func (c *IDCounter) NextEachMarker() string {
	id := "e" + itoa(c.Each)
	c.Each++
	return id
}

// NextCompMarker returns the next marker ID for nested components.
// Used internally for component prefixing (e.g., "c0" in "components_c0_t0")
func (c *IDCounter) NextCompMarker() string {
	id := "c" + itoa(c.Comp)
	c.Comp++
	return id
}

// NextRouteMarker returns the next marker ID for route-blocks.
// Used in: <!--basics_r0--> (comment marker for route-block boundary)
func (c *IDCounter) NextRouteMarker() string {
	id := "r" + itoa(c.Route)
	c.Route++
	return id
}

// --- ID formatting functions ---

// FullID returns the full ID with prefix.
// Example: FullID("ev0") with prefix "basics" returns "basics_ev0"
func (c *IDCounter) FullID(localID string) string {
	if c.Prefix == "" {
		return localID
	}
	return c.Prefix + "_" + localID
}
