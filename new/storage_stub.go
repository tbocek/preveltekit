//go:build !wasm

package preveltekit

// In-memory storage for pre-rendering
var memStorage = make(map[string]string)

// GetStorage stub - uses in-memory map during pre-render.
func GetStorage(key string) string {
	return memStorage[key]
}

// SetStorage stub - uses in-memory map during pre-render.
func SetStorage(key, value string) {
	memStorage[key] = value
}

// RemoveStorage stub - uses in-memory map during pre-render.
func RemoveStorage(key string) {
	delete(memStorage, key)
}

// ClearStorage stub - clears in-memory map during pre-render.
func ClearStorage() {
	memStorage = make(map[string]string)
}

// NewLocalStore stub - returns a LocalStore during pre-render.
// The key is also used as the store's ID for hydration.
func NewLocalStore(key string, defaultValue string) *LocalStore {
	stored := GetStorage(key)
	if stored == "" {
		stored = defaultValue
	}
	return &LocalStore{Store: NewWithID(key, stored)}
}
