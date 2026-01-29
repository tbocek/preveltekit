package main

import p "preveltekit"

type Storage struct {
	Theme  *p.LocalStore    // auto-persisted to localStorage
	Notes  *p.Store[string] // manual save
	Status *p.Store[string]
}

func (s *Storage) OnMount() {
	// Theme is automatically loaded and saved via LocalStore
	// Set default if no saved value exists
	if s.Theme.Get() == "" {
		s.Theme.Set("light")
	}

	// Load saved notes (manual persistence example)
	if saved := p.GetStorage("notes"); saved != "" {
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
	p.SetStorage("notes", s.Notes.Get())
	s.Status.Set("Notes saved!")
}

func (s *Storage) ClearNotes() {
	s.Notes.Set("")
	p.RemoveStorage("notes")
	s.Status.Set("Notes cleared")
}

func (s *Storage) ClearAll() {
	p.ClearStorage()
	s.Theme.Set("light")
	s.Notes.Set("")
	s.Status.Set("All storage cleared")
}

func (s *Storage) Render() p.Node {
	return p.Html(`<div class="demo">
		<h1>Storage</h1>

		<section>
			<h2>Persisted Theme (Auto-sync)</h2>
			<p>Theme preference is automatically saved to localStorage.</p>
			<p>Current theme: <strong>`, p.Bind(s.Theme.Store), `</strong></p>
			<div class="buttons">
				`, p.Html(`<button>Light</button>`).WithOn("click", s.SetLight), `
				`, p.Html(`<button>Dark</button>`).WithOn("click", s.SetDark), `
			</div>
			<p class="hint">Refresh the page - theme will persist!</p>
		</section>

		<section>
			<h2>Manual Storage (Notes)</h2>
			<p>Notes are saved manually when you click Save.</p>
			`, p.BindValue(`<textarea placeholder="Type your notes here..."></textarea>`, s.Notes), `
			<div class="buttons">
				`, p.Html(`<button>Save Notes</button>`).WithOn("click", s.SaveNotes), `
				`, p.Html(`<button>Clear Notes</button>`).WithOn("click", s.ClearNotes), `
			</div>
		</section>

		<section>
			<h2>Clear All Storage</h2>
			`, p.Html(`<button class="danger">Clear All Storage</button>`).WithOn("click", s.ClearAll), `
		</section>

		<p class="status">`, p.Bind(s.Status), `</p>
	</div>`)
}

func (s *Storage) Style() string {
	return `
.demo textarea{height:100px}
.demo button.danger{background:#ffebee;border-color:#ef9a9a;color:#c62828}
.demo button.danger:hover{background:#ffcdd2}
.demo .status{padding:10px;background:#e8f5e9;border-radius:4px;color:#2e7d32}
`
}
