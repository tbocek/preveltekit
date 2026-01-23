package main

import (
	"reactive"
	"strings"
)

type FetchDemo struct {
	// Status tracking
	Status *reactive.Store[string]

	// For simple fetch demo
	RawData *reactive.Store[string]

	// For list fetch demo - shows diff in action
	Users     *reactive.List[string]
	UserCount *reactive.Store[int]
}

func (f *FetchDemo) OnMount() {
	if reactive.IsBuildTime {
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

	reactive.FetchJSON("https://jsonplaceholder.typicode.com/todos/1", func(data string, err error) {
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

	reactive.FetchJSON("https://jsonplaceholder.typicode.com/users", func(data string, err error) {
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

	reactive.FetchJSON("https://jsonplaceholder.typicode.com/users?_limit=3", func(data string, err error) {
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

func (f *FetchDemo) Template() string {
	return `<div class="app">
	<h1>Fetch Demo</h1>
	<p class="status">Status: <strong>{Status}</strong></p>

	<section>
		<h2>1. Simple Fetch</h2>
		<p>Fetch raw JSON data:</p>
		<button @click="FetchTodo()">Fetch Todo</button>
		<pre>{RawData}</pre>
	</section>

	<section>
		<h2>2. Fetch into List (Diff Demo)</h2>
		<p>Fetches user names and uses List.Set() which triggers diff:</p>
		<div class="button-row">
			<button @click="FetchUsers()">Fetch All Users</button>
			<button @click="FetchFewUsers()">Fetch 3 Users</button>
			<button @click="AddLocalUser()">Add Local</button>
			<button @click="ClearUsers()">Clear</button>
		</div>

		<p>Users: {UserCount}</p>

		{#if UserCount > 0}
			<ul>
				{#each Users as user, i}
					<li><span class="index">{i}</span> {user}</li>
				{/each}
			</ul>
		{:else}
			<p class="empty">No users loaded</p>
		{/if}

		<p class="note">
			Try: Load all → Load 3 (watch items get removed via diff)<br>
			Or: Load 3 → Load all (watch items get added via diff)
		</p>
	</section>
</div>`
}

func (f *FetchDemo) Style() string {
	return `
.app { font-family: system-ui, sans-serif; max-width: 700px; margin: 0 auto; padding: 20px; }
section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 8px; }
h1 { color: #333; }
h2 { margin-top: 0; color: #666; font-size: 1.1em; }
.status { padding: 10px; background: #f0f0f0; border-radius: 4px; }
pre { background: #f5f5f5; padding: 15px; border-radius: 4px; overflow-x: auto; min-height: 50px; font-size: 12px; }
button { padding: 8px 16px; margin: 4px; cursor: pointer; border: 1px solid #ccc; border-radius: 4px; background: #f5f5f5; }
button:hover { background: #e5e5e5; }
.button-row { display: flex; gap: 10px; flex-wrap: wrap; margin: 10px 0; }
ul { list-style: none; padding: 0; margin: 10px 0; }
li { padding: 8px 12px; margin: 4px 0; background: #f9f9f9; border-radius: 4px; border-left: 3px solid #007bff; }
.index { display: inline-block; width: 20px; height: 20px; line-height: 20px; text-align: center; background: #007bff; color: white; border-radius: 50%; font-size: 11px; margin-right: 10px; }
.empty { color: #999; font-style: italic; }
.note { color: #666; font-size: 0.9em; background: #fff3cd; padding: 10px; border-radius: 4px; }
`
}
