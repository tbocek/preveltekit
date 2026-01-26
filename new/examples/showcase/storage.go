package main

import "preveltekit"

type Storage struct {
	Theme  *preveltekit.LocalStore    // auto-persisted to localStorage
	Notes  *preveltekit.Store[string] // manual save
	Status *preveltekit.Store[string]
}

func (s *Storage) OnMount() {
	// Theme is automatically loaded and saved via LocalStore
	// Set default if no saved value exists
	if s.Theme.Get() == "" {
		s.Theme.Set("light")
	}

	// Load saved notes (manual persistence example)
	if saved := preveltekit.GetStorage("notes"); saved != "" {
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
	preveltekit.SetStorage("notes", s.Notes.Get())
	s.Status.Set("Notes saved!")
}

func (s *Storage) ClearNotes() {
	s.Notes.Set("")
	preveltekit.RemoveStorage("notes")
	s.Status.Set("Notes cleared")
}

func (s *Storage) ClearAll() {
	preveltekit.ClearStorage()
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
.demo textarea{height:100px}
.demo button.danger{background:#ffebee;border-color:#ef9a9a;color:#c62828}
.demo button.danger:hover{background:#ffcdd2}
.demo .status{padding:10px;background:#e8f5e9;border-radius:4px;color:#2e7d32}
`
}
