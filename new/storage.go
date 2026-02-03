//go:build wasm

package preveltekit

import "syscall/js"

var localStorage = js.Global().Get("localStorage")

// GetStorage retrieves a value from localStorage.
// Returns empty string if key doesn't exist.
func GetStorage(key string) string {
	val := localStorage.Call("getItem", key)
	if val.IsNull() || val.IsUndefined() {
		return ""
	}
	return val.String()
}

// SetStorage stores a value in localStorage.
func SetStorage(key, value string) {
	localStorage.Call("setItem", key, value)
}

// RemoveStorage removes a key from localStorage.
func RemoveStorage(key string) {
	localStorage.Call("removeItem", key)
}

// ClearStorage removes all keys from localStorage.
func ClearStorage() {
	localStorage.Call("clear")
}

// NewLocalStore creates a LocalStore that automatically syncs with localStorage.
// The store is initialized with the localStorage value if it exists,
// otherwise uses the provided default value. All changes are automatically persisted.
// The key is also used as the store's ID for hydration.
func NewLocalStore(key string, defaultValue string) *LocalStore {
	stored := GetStorage(key)
	if stored == "" {
		stored = defaultValue
	}

	store := New(key, stored)
	store.OnChange(func(v string) {
		SetStorage(key, v)
	})

	return &LocalStore{Store: store}
}
