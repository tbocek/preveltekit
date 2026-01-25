// Package reactive provides generic reactive stores for Go WebAssembly applications.
package preveltekit

// StaticRoute defines a route for both build-time pre-rendering and runtime routing.
// This is the single source of truth for your application's routes.
type StaticRoute struct {
	Path     string                         // URL path (e.g., "/doc")
	HTMLFile string                         // Output filename for pre-rendering (e.g., "doc.html")
	Handler  func(params map[string]string) // Route handler for runtime navigation
}

// Store is a generic reactive container that calls callbacks on mutation
type Store[T any] struct {
	value     T
	callbacks []func(T)
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

// New creates a reactive store with an initial value
func New[T any](initial T) *Store[T] {
	return &Store[T]{value: initial}
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
	items    []T
	onEdit   []func(Edit[T])
	onRender []func([]T) // for initial render only
	onChange []func([]T) // called on any change
}

// NewList creates a reactive list
func NewList[T comparable](initial ...T) *List[T] {
	return &List[T]{
		items: initial,
	}
}

// Get returns a copy of the slice (safe, no mutation leaks)
func (l *List[T]) Get() []T {
	cp := make([]T, len(l.items))
	copy(cp, l.items)
	return cp
}

// Len returns the length
func (l *List[T]) Len() int {
	return len(l.items)
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

// OnRender adds a callback for initial render
func (l *List[T]) OnRender(cb func([]T)) {
	l.onRender = append(l.onRender, cb)
}

// OnChange adds a callback for any change to the list
func (l *List[T]) OnChange(cb func([]T)) {
	l.onChange = append(l.onChange, cb)
}

// ClearCallbacks removes all registered callbacks.
// Used when re-rendering if-blocks to prevent callback accumulation.
func (l *List[T]) ClearCallbacks() {
	l.onEdit = nil
	l.onRender = nil
	l.onChange = nil
}

// Render triggers initial render callbacks
func (l *List[T]) Render() {
	for _, cb := range l.onRender {
		cb(l.items)
	}
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
