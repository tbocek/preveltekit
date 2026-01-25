// Package reactive provides generic reactive stores for Go WebAssembly applications.
package reactive

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

// Render triggers initial render callbacks
func (l *List[T]) Render() {
	for _, cb := range l.onRender {
		cb(l.items)
	}
}

// diff computes minimal edit operations to transform old into new
// Uses Myers diff algorithm with early termination of suboptimal paths
func diff[T comparable](old, new []T) []Edit[T] {
	n, m := len(old), len(new)

	// Fast paths
	if n == 0 {
		edits := make([]Edit[T], m)
		for i, v := range new {
			edits[i] = Edit[T]{Op: EditInsert, Index: i, Value: v}
		}
		return edits
	}
	if m == 0 {
		edits := make([]Edit[T], n)
		for i := n - 1; i >= 0; i-- {
			var zero T
			edits[n-1-i] = Edit[T]{Op: EditRemove, Index: i, Value: zero}
		}
		return edits
	}

	// Check for simple append
	if m > n {
		match := true
		for i := 0; i < n; i++ {
			if old[i] != new[i] {
				match = false
				break
			}
		}
		if match {
			edits := make([]Edit[T], m-n)
			for i := n; i < m; i++ {
				edits[i-n] = Edit[T]{Op: EditInsert, Index: i, Value: new[i]}
			}
			return edits
		}
	}

	// Check for simple prepend
	if m > n {
		diff := m - n
		match := true
		for i := 0; i < n; i++ {
			if old[i] != new[i+diff] {
				match = false
				break
			}
		}
		if match {
			edits := make([]Edit[T], diff)
			for i := 0; i < diff; i++ {
				edits[i] = Edit[T]{Op: EditInsert, Index: i, Value: new[i]}
			}
			return edits
		}
	}

	// Check for simple removal from end
	if n > m {
		match := true
		for i := 0; i < m; i++ {
			if old[i] != new[i] {
				match = false
				break
			}
		}
		if match {
			edits := make([]Edit[T], n-m)
			for i := n - 1; i >= m; i-- {
				var zero T
				edits[n-1-i] = Edit[T]{Op: EditRemove, Index: i, Value: zero}
			}
			return edits
		}
	}

	// Myers diff algorithm
	max := n + m
	v := make(map[int]int)
	v[1] = 0
	var trace []map[int]int

	for d := 0; d <= max; d++ {
		// Copy v for backtracking
		vc := make(map[int]int)
		for k, val := range v {
			vc[k] = val
		}
		trace = append(trace, vc)

		for k := -d; k <= d; k += 2 {
			var x int
			if k == -d || (k != d && v[k-1] < v[k+1]) {
				x = v[k+1] // move down (insert)
			} else {
				x = v[k-1] + 1 // move right (remove)
			}
			y := x - k

			// Follow diagonal (matches)
			for x < n && y < m && old[x] == new[y] {
				x++
				y++
			}

			v[k] = x

			if x >= n && y >= m {
				// Found path, backtrack to build edits
				return backtrack(trace, old, new, n, m)
			}
		}
	}

	return nil // unreachable
}

// backtrack reconstructs edit operations from Myers diff trace
func backtrack[T comparable](trace []map[int]int, old, new []T, n, m int) []Edit[T] {
	var edits []Edit[T]
	x, y := n, m

	for d := len(trace) - 1; d >= 0; d-- {
		v := trace[d]
		k := x - y

		var prevK int
		if k == -d || (k != d && v[k-1] < v[k+1]) {
			prevK = k + 1 // came from above (insert)
		} else {
			prevK = k - 1 // came from left (remove)
		}

		prevX := v[prevK]
		prevY := prevX - prevK

		// Follow diagonal backwards
		for x > prevX && y > prevY {
			x--
			y--
			// match, no edit needed
		}

		if d > 0 {
			if x == prevX {
				// insert
				y--
				edits = append(edits, Edit[T]{Op: EditInsert, Index: y, Value: new[y]})
			} else {
				// remove
				x--
				var zero T
				edits = append(edits, Edit[T]{Op: EditRemove, Index: x, Value: zero})
			}
		}
	}

	// Reverse to get correct order
	for i, j := 0, len(edits)-1; i < j; i, j = i+1, j-1 {
		edits[i], edits[j] = edits[j], edits[i]
	}

	return edits
}

// Map is a reactive map with methods that trigger updates
type Map[K comparable, V any] struct {
	items     map[K]V
	callbacks []func(map[K]V)
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
	m.items[key] = value
	m.notify()
}

// Delete removes a key
func (m *Map[K, V]) Delete(key K) {
	delete(m.items, key)
	m.notify()
}

// Clear removes all entries
func (m *Map[K, V]) Clear() {
	m.items = make(map[K]V)
	m.notify()
}

// OnChange adds a callback for when the map changes
func (m *Map[K, V]) OnChange(cb func(map[K]V)) {
	m.callbacks = append(m.callbacks, cb)
}

func (m *Map[K, V]) notify() {
	for _, cb := range m.callbacks {
		cb(m.items)
	}
}
