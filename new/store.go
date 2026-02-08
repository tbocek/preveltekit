// Package reactive provides generic reactive stores for Go WebAssembly applications.
package preveltekit

// Component is the interface that all declarative components must implement.
type Component interface {
	Render() Node
}

// Route defines a single route for both build-time pre-rendering and runtime routing.
type Route struct {
	Path      string    // URL path pattern (e.g., "/user/:id")
	HTMLFile  string    // Output filename for pre-rendering (e.g., "user.html")
	SSRPath   string    // URL to pre-render (empty = skip SSR)
	Component Component // Component to render for this route
}

// HasRoutes is implemented by components that define routes.
type HasRoutes interface {
	Routes() []Route
}

// HasStyle is implemented by components that have scoped CSS styles.
type HasStyle interface {
	Style() string
}

// HasGlobalStyle is implemented by components that have unscoped global CSS styles.
// Global styles are emitted without any CSS scoping — useful for base/reset styles.
type HasGlobalStyle interface {
	GlobalStyle() string
}

// HasNew is implemented by components that can create fresh instances.
type HasNew interface {
	New() Component
}

// HasOnMount is implemented by components with OnMount lifecycle.
type HasOnMount interface {
	OnMount()
}

// HasID is implemented by stores that have a user-defined ID.
type HasID interface {
	ID() string
}

// Store is a generic reactive container that calls callbacks on mutation
type Store[T any] struct {
	id        string
	value     T
	callbacks []func(T)
	options   []any // possible values for pre-baked rendering (used by Store[Component])
}

// WithOptions registers alternative values this store may hold.
// Used by SSR to pre-render all possible components for a Store[Component].
// For routing, NewRouter calls this automatically. For other cases (tabs, wizards),
// call it manually: store.WithOptions(compA, compB, compC)
func (s *Store[T]) WithOptions(alternatives ...T) {
	for _, alt := range alternatives {
		s.options = append(s.options, alt)
	}
}

// Options returns the registered alternative values for this store.
func (s *Store[T]) Options() []any {
	return s.options
}

// storeRegistry holds all registered stores by ID for hydration lookup
var storeRegistry = make(map[string]any)

// storeCounter generates unique auto-IDs for stores and lists
var storeCounter int

// nextStoreID returns the next auto-generated store ID (s0, s1, s2, ...)
func nextStoreID() string {
	id := "s" + itoa(storeCounter)
	storeCounter++
	return id
}

// resetRegistries resets all global counters and registries to initial state.
// Called before each SSR iteration so IDs start from s0, matching WASM.
func resetRegistries() {
	storeCounter = 0
	handlerCounter = 0
	scopeCounter = 0
	storeRegistry = make(map[string]any)
	handlerRegistry = make(map[string]func())
	handlerModifiers = make(map[string][]string)
	scopeRegistry = make(map[string]string)
}

// handlerRegistry holds all registered event handlers by ID for hydration lookup
var handlerRegistry = make(map[string]func())

// handlerModifiers holds event modifiers (preventDefault, stopPropagation) by handler ID.
// Set by PreventDefault()/StopPropagation() after On() registers the handler.
var handlerModifiers = make(map[string][]string)

// handlerCounter generates unique auto-IDs for event handlers
var handlerCounter int

// nextHandlerID returns the next auto-generated handler ID (h0, h1, h2, ...)
func nextHandlerID() string {
	id := "h" + itoa(handlerCounter)
	handlerCounter++
	return id
}

// GetStore looks up a store by ID from the global registry
func GetStore(id string) any {
	return storeRegistry[id]
}

// GetHandler looks up a handler by ID from the global registry
func GetHandler(id string) func() {
	return handlerRegistry[id]
}

// GetHandlerModifiers returns the modifiers for a handler ID (e.g., ["preventDefault"])
func GetHandlerModifiers(id string) []string {
	return handlerModifiers[id]
}

// scopeRegistry maps component name → scope class (e.g., "app" → "v0").
// Used for Svelte-style CSS scoping: each component with Style() gets a unique class.
var scopeRegistry = make(map[string]string)

// scopeCounter generates unique scope IDs (v0, v1, v2, ...)
var scopeCounter int

// GetOrCreateScope returns the scope class name for a component name.
// Creates a new one if the component hasn't been seen yet.
// Returns e.g. "v0".
func GetOrCreateScope(componentName string) string {
	if cls, ok := scopeRegistry[componentName]; ok {
		return cls
	}
	cls := "v" + itoa(scopeCounter)
	scopeCounter++
	scopeRegistry[componentName] = cls
	return cls
}

// RegisterHandler registers an event handler, auto-generating a unique ID.
// Returns the generated ID.
func RegisterHandler(handler func()) string {
	id := nextHandlerID()
	handlerRegistry[id] = handler
	return id
}

// LocalStore is a Store[string] that automatically syncs with localStorage.
// Use it for persisting string values across page reloads.
//
//	type App struct {
//	    Theme *LocalStore  // auto-persisted to localStorage with key "Theme"
//	}
type LocalStore struct {
	*Store[string]
}

// New creates a reactive store with an auto-generated ID and initial value.
// The ID is deterministic (counter-based) so SSR and WASM produce matching IDs
// when stores are created in the same order.
func New[T any](initial T) *Store[T] {
	id := nextStoreID()
	s := &Store[T]{id: id, value: initial}
	storeRegistry[id] = s
	return s
}

// NewWithID creates a reactive store with an explicit ID and initial value.
// Use this only when a specific ID is required (e.g. router path store, localStorage).
func NewWithID[T any](id string, initial T) *Store[T] {
	s := &Store[T]{id: id, value: initial}
	storeRegistry[id] = s
	return s
}

// ID returns the store's unique identifier
func (s *Store[T]) ID() string {
	return s.id
}

// Get returns the current value
func (s *Store[T]) Get() T {
	return s.value
}

// Set updates the value and calls all callbacks
func (s *Store[T]) Set(v T) {
	s.value = v
	s.notify()
}

// Update applies a function to transform the current value
func (s *Store[T]) Update(fn func(T) T) {
	s.Set(fn(s.value))
}

// OnChange adds a callback that runs whenever the value changes
func (s *Store[T]) OnChange(cb func(T)) {
	s.callbacks = append(s.callbacks, cb)
}

func (s *Store[T]) notify() {
	for _, cb := range s.callbacks {
		cb(s.value)
	}
}

// EditOp represents a list edit operation
type EditOp int

const (
	EditInsert EditOp = iota
	EditRemove
)

// Edit represents a single edit operation
type Edit[T any] struct {
	Op    EditOp
	Index int // position in new list for Insert, old list for Remove
	Value T   // value for Insert
}

// List is a reactive slice with methods that trigger updates
type List[T comparable] struct {
	id       string
	items    []T
	lenStore *Store[int] // cached length store for reactive conditions
	onEdit   []func(Edit[T])
	onChange []func([]T) // called on any change
}

// NewList creates a reactive list with an auto-generated ID.
// The ID is deterministic (counter-based) so SSR and WASM produce matching IDs
// when lists are created in the same order.
func NewList[T comparable](initial ...T) *List[T] {
	id := nextStoreID()
	l := &List[T]{
		id:    id,
		items: initial,
	}
	storeRegistry[id] = l
	return l
}

// ID returns the list's unique identifier.
func (l *List[T]) ID() string {
	return l.id
}

// Get returns a copy of the slice (safe, no mutation leaks)
func (l *List[T]) Get() []T {
	cp := make([]T, len(l.items))
	copy(cp, l.items)
	return cp
}

// Len returns a reactive store tracking the list length.
// The store auto-updates when the list changes.
func (l *List[T]) Len() *Store[int] {
	if l.lenStore == nil {
		l.lenStore = NewWithID(l.id+".len", len(l.items))
		l.OnChange(func(_ []T) {
			l.lenStore.Set(len(l.items))
		})
	}
	return l.lenStore
}

// At returns item at index
func (l *List[T]) At(i int) T {
	return l.items[i]
}

// Set replaces the entire list, computing minimal diff
func (l *List[T]) Set(items []T) {
	old := l.items
	l.items = items

	// Compute and apply edits
	edits := diff(old, items)
	for _, edit := range edits {
		for _, cb := range l.onEdit {
			cb(edit)
		}
	}
	// Notify onChange listeners
	for _, cb := range l.onChange {
		cb(l.items)
	}
}

// Append adds items to the end
func (l *List[T]) Append(items ...T) {
	for _, item := range items {
		edit := Edit[T]{Op: EditInsert, Index: len(l.items), Value: item}
		l.items = append(l.items, item)
		for _, cb := range l.onEdit {
			cb(edit)
		}
	}
	// Notify onChange listeners
	for _, cb := range l.onChange {
		cb(l.items)
	}
}

// RemoveAt removes item at index
func (l *List[T]) RemoveAt(i int) {
	var zero T
	edit := Edit[T]{Op: EditRemove, Index: i, Value: zero}
	l.items = append(l.items[:i], l.items[i+1:]...)
	for _, cb := range l.onEdit {
		cb(edit)
	}
	// Notify onChange listeners
	for _, cb := range l.onChange {
		cb(l.items)
	}
}

// Clear removes all items
func (l *List[T]) Clear() {
	var zero T
	for i := len(l.items) - 1; i >= 0; i-- {
		edit := Edit[T]{Op: EditRemove, Index: i, Value: zero}
		for _, cb := range l.onEdit {
			cb(edit)
		}
	}
	l.items = l.items[:0]
	// Notify onChange listeners
	for _, cb := range l.onChange {
		cb(l.items)
	}
}

// OnEdit adds a callback for edit operations (Insert, Remove)
func (l *List[T]) OnEdit(cb func(Edit[T])) {
	l.onEdit = append(l.onEdit, cb)
}

// OnChange adds a callback for any change to the list
func (l *List[T]) OnChange(cb func([]T)) {
	l.onChange = append(l.onChange, cb)
}

// diff computes edit operations to transform old into new using O(n) set comparison.
// Returns removes first (in reverse index order), then inserts.
func diff[T comparable](old, new []T) []Edit[T] {
	// Build sets for O(1) lookup
	oldSet := make(map[T]bool, len(old))
	for _, v := range old {
		oldSet[v] = true
	}
	newSet := make(map[T]bool, len(new))
	for _, v := range new {
		newSet[v] = true
	}

	var edits []Edit[T]

	// Find removals (in old but not in new) - reverse order for stable indices
	for i := len(old) - 1; i >= 0; i-- {
		if !newSet[old[i]] {
			var zero T
			edits = append(edits, Edit[T]{Op: EditRemove, Index: i, Value: zero})
		}
	}

	// Find insertions (in new but not in old)
	for i, v := range new {
		if !oldSet[v] {
			edits = append(edits, Edit[T]{Op: EditInsert, Index: i, Value: v})
		}
	}

	return edits
}

// MapEdit represents a single map edit operation
type MapEdit[K comparable, V any] struct {
	Op    EditOp
	Key   K
	Value V // Value for Insert, zero for Remove
}

// Map is a reactive map with methods that trigger updates
type Map[K comparable, V any] struct {
	items    map[K]V
	onEdit   []func(MapEdit[K, V])
	onChange []func(map[K]V)
}

// NewMap creates a reactive map
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{items: make(map[K]V)}
}

// Get returns value for key
func (m *Map[K, V]) Get(key K) (V, bool) {
	v, ok := m.items[key]
	return v, ok
}

// Keys returns all keys
func (m *Map[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.items))
	for k := range m.items {
		keys = append(keys, k)
	}
	return keys
}

// Len returns the number of entries
func (m *Map[K, V]) Len() int {
	return len(m.items)
}

// Set sets a key-value pair
func (m *Map[K, V]) Set(key K, value V) {
	_, exists := m.items[key]
	m.items[key] = value
	if !exists {
		for _, cb := range m.onEdit {
			cb(MapEdit[K, V]{Op: EditInsert, Key: key, Value: value})
		}
	}
	m.notify()
}

// SetAll replaces all entries, computing minimal diff
func (m *Map[K, V]) SetAll(items map[K]V) {
	// Find removals (in old but not in new)
	for k := range m.items {
		if _, exists := items[k]; !exists {
			var zero V
			for _, cb := range m.onEdit {
				cb(MapEdit[K, V]{Op: EditRemove, Key: k, Value: zero})
			}
		}
	}
	// Find insertions (in new but not in old)
	for k, v := range items {
		if _, exists := m.items[k]; !exists {
			for _, cb := range m.onEdit {
				cb(MapEdit[K, V]{Op: EditInsert, Key: k, Value: v})
			}
		}
	}
	m.items = items
	m.notify()
}

// Delete removes a key
func (m *Map[K, V]) Delete(key K) {
	if _, exists := m.items[key]; exists {
		var zero V
		for _, cb := range m.onEdit {
			cb(MapEdit[K, V]{Op: EditRemove, Key: key, Value: zero})
		}
	}
	delete(m.items, key)
	m.notify()
}

// Clear removes all entries
func (m *Map[K, V]) Clear() {
	var zero V
	for k := range m.items {
		for _, cb := range m.onEdit {
			cb(MapEdit[K, V]{Op: EditRemove, Key: k, Value: zero})
		}
	}
	m.items = make(map[K]V)
	m.notify()
}

// OnEdit adds a callback for edit operations (Insert, Remove)
func (m *Map[K, V]) OnEdit(cb func(MapEdit[K, V])) {
	m.onEdit = append(m.onEdit, cb)
}

// OnChange adds a callback for when the map changes
func (m *Map[K, V]) OnChange(cb func(map[K]V)) {
	m.onChange = append(m.onChange, cb)
}

func (m *Map[K, V]) notify() {
	for _, cb := range m.onChange {
		cb(m.items)
	}
}
