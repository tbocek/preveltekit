package main

import "reactive"

type Storage struct {
	Theme  *reactive.LocalStore    // auto-persisted to localStorage
	Notes  *reactive.Store[string] // manual save
	Status *reactive.Store[string]
}

func (s *Storage) OnMount() {
	// Theme is automatically loaded and saved via LocalStore
	// Set default if no saved value exists
	if s.Theme.Get() == "" {
		s.Theme.Set("light")
	}

	// Load saved notes (manual persistence example)
	if saved := reactive.GetStorage("notes"); saved != "" {
		s.Notes.Set(saved)
	}

	s.Status.Set("Ready")
}

func (s *Storage) SetLight() {
	s.Theme.Set("light")
	s.Status.Set("Theme set to light (saved)")
}

func (s *Storage) SetDark() {
	s.Theme.Set("dark")
	s.Status.Set("Theme set to dark (saved)")
}

func (s *Storage) SaveNotes() {
	reactive.SetStorage("notes", s.Notes.Get())
	s.Status.Set("Notes saved!")
}

func (s *Storage) ClearNotes() {
	s.Notes.Set("")
	reactive.RemoveStorage("notes")
	s.Status.Set("Notes cleared")
}

func (s *Storage) ClearAll() {
	reactive.ClearStorage()
	s.Theme.Set("light")
	s.Notes.Set("")
	s.Status.Set("All storage cleared")
}

func (s *Storage) Template() string {
	return `<div class="demo">
	<h1>Storage</h1>

	<section>
		<h2>Persisted Theme (Auto-sync)</h2>
		<p>Theme preference is automatically saved to localStorage.</p>
		<p>Current theme: <strong>{Theme}</strong></p>
		<div class="buttons">
			<button @click="SetLight()">Light</button>
			<button @click="SetDark()">Dark</button>
		</div>
		<p class="hint">Refresh the page - theme will persist!</p>
	</section>

	<section>
		<h2>Manual Storage (Notes)</h2>
		<p>Notes are saved manually when you click Save.</p>
		<textarea bind:value="Notes" placeholder="Type your notes here..."></textarea>
		<div class="buttons">
			<button @click="SaveNotes()">Save Notes</button>
			<button @click="ClearNotes()">Clear Notes</button>
		</div>
	</section>

	<section>
		<h2>Clear All Storage</h2>
		<button class="danger" @click="ClearAll()">Clear All Storage</button>
	</section>

	<p class="status">{Status}</p>
</div>`
}

func (s *Storage) Style() string {
	return `
.demo { max-width: 600px; }
.demo h1 { color: #1a1a2e; margin-bottom: 20px; }
.demo section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 8px; background: #fff; }
.demo h2 { margin-top: 0; color: #666; font-size: 1.1em; }
.demo button { padding: 8px 16px; margin: 4px; cursor: pointer; border: 1px solid #ccc; border-radius: 4px; background: #f5f5f5; }
.demo button:hover { background: #e5e5e5; }
.demo button.danger { background: #ffebee; border-color: #ef9a9a; color: #c62828; }
.demo button.danger:hover { background: #ffcdd2; }
.demo textarea { width: 100%; height: 100px; padding: 10px; border: 1px solid #ccc; border-radius: 4px; font-family: inherit; resize: vertical; }
.demo .status { padding: 10px; background: #e8f5e9; border-radius: 4px; color: #2e7d32; }
.demo .hint { font-size: 0.9em; color: #666; font-style: italic; }
.buttons { display: flex; gap: 10px; margin: 10px 0; }
`
}
