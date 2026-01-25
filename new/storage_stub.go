//go:build !wasm

package reactive

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
func NewLocalStore(key string, defaultValue string) *LocalStore {
	stored := GetStorage(key)
	if stored == "" {
		stored = defaultValue
	}
	return &LocalStore{New(stored)}
}
