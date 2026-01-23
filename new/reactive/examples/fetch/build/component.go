package main

import (
	"strings"
)

type FetchDemo struct {
	// Status tracking
	Status *Store[string]

	// For simple fetch demo
	RawData *Store[string]

	// For list fetch demo - shows diff in action
	Users     *List[string]
	UserCount *Store[int]
}

func (f *FetchDemo) OnMount() {
	if IsBuildTime {
		// At SSR time, show loading placeholder
		f.Status.Set("loading...")
		f.RawData.Set("Loading data...")
		f.UserCount.Set(0)
		return
	}

	// At runtime, start with idle state
	f.Status.Set("idle")
	f.RawData.Set("")
	f.UserCount.Set(0)
}

// Simple fetch - returns raw JSON
func (f *FetchDemo) FetchTodo() {
	f.Status.Set("loading...")
	f.RawData.Set("")

	FetchJSON("https://jsonplaceholder.typicode.com/todos/1", func(data string, err error) {
		if err != nil {
			f.Status.Set("error: " + err.Error())
			return
		}
		f.RawData.Set(data)
		f.Status.Set("done")
	})
}

// Fetch into list - demonstrates diff
func (f *FetchDemo) FetchUsers() {
	f.Status.Set("loading users...")

	FetchJSON("https://jsonplaceholder.typicode.com/users", func(data string, err error) {
		if err != nil {
			f.Status.Set("error: " + err.Error())
			return
		}

		// Parse user names from JSON (simple string extraction)
		// In a real app, you'd use encoding/json
		names := extractNames(data)
		f.Users.Set(names)
		f.UserCount.Set(len(names))
		f.Status.Set("loaded " + string(rune('0'+len(names))) + " users")
	})
}

// Fetch fewer users (shows diff removing items)
func (f *FetchDemo) FetchFewUsers() {
	f.Status.Set("loading subset...")

	FetchJSON("https://jsonplaceholder.typicode.com/users?_limit=3", func(data string, err error) {
		if err != nil {
			f.Status.Set("error: " + err.Error())
			return
		}

		names := extractNames(data)
		f.Users.Set(names)
		f.UserCount.Set(len(names))
		f.Status.Set("loaded " + string(rune('0'+len(names))) + " users (subset)")
	})
}

func (f *FetchDemo) ClearUsers() {
	f.Users.Set([]string{})
	f.UserCount.Set(0)
	f.Status.Set("cleared")
}

func (f *FetchDemo) AddLocalUser() {
	f.Users.Append("Local User")
	f.UserCount.Set(f.Users.Len())
	f.Status.Set("added local user")
}

// Helper to extract names from JSON array
// This is a simple parser - in production use encoding/json
func extractNames(json string) []string {
	var names []string
	// Find all "name": "..." patterns
	parts := strings.Split(json, `"name"`)
	for i := 1; i < len(parts); i++ {
		part := parts[i]
		// Find the value after ": "
		start := strings.Index(part, `"`)
		if start == -1 {
			continue
		}
		part = part[start+1:]
		end := strings.Index(part, `"`)
		if end == -1 {
			continue
		}
		names = append(names, part[:end])
	}
	return names
}

