package main

import p "github.com/tbocek/preveltekit"

type Storage struct {
	Theme  *p.LocalStore    // auto-persisted to localStorage
	Notes  *p.Store[string] // manual save
	Status *p.Store[string]
}

func (s *Storage) New() p.Component {
	return &Storage{
		Theme:  p.NewLocalStore("theme", "light"),
		Notes:  p.New(""),
		Status: p.New("Ready"),
	}
}

func (s *Storage) OnMount() {
	// Load saved notes (manual persistence example)
	if saved := p.GetStorage("notes"); saved != "" {
		s.Notes.Set(saved)
	}
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
			<p>Current theme: <strong>`, s.Theme.Store, `</strong></p>
			<div class="buttons">`,
		p.Html(`<button>Light</button>`).On("click", s.SetLight),
		p.Html(`<button>Dark</button>`).On("click", s.SetDark),
		`</div>
			<p class="hint">Refresh the page - theme will persist!</p>
		</section>

		<section>
			<h2>Manual Storage (Notes)</h2>
			<p>Notes are saved manually when you click Save.</p>
			`, p.Html(`<textarea placeholder="Type your notes here..."></textarea>`).Bind(s.Notes), `
			<div class="buttons">`,
		p.Html(`<button>Save Notes</button>`).On("click", s.SaveNotes),
		p.Html(`<button>Clear Notes</button>`).On("click", s.ClearNotes),
		`</div>
		</section>

		<section>
			<h2>Clear All Storage</h2>
			`, p.Html(`<button class="danger">Clear All Storage</button>`).On("click", s.ClearAll), `
		</section>

		<p class="status">`, s.Status, `</p>

		<section>
			<h2>Code</h2>
			<pre class="code">// auto-persisted store (syncs to localStorage on every Set)
Theme := p.NewLocalStore("theme", "light")
Theme.Set("dark") // automatically saved

// read the store value (it's a *Store[string] inside)
Theme.Store // use in Html() like any store

// manual localStorage API
p.SetStorage("notes", "hello")
saved := p.GetStorage("notes")
p.RemoveStorage("notes")
p.ClearStorage()</pre>
		</section>
	</div>`)
}

func (s *Storage) Style() string {
	return `
.demo pre.code{background:#1a1a2e;color:#e0e0e0;font-size:12px;margin-top:12px}
.demo textarea{height:100px}
.demo button.danger{background:#ffebee;border-color:#ef9a9a;color:#c62828}
.demo button.danger:hover{background:#ffcdd2}
.demo .status{padding:10px;background:#e8f5e9;border-radius:4px;color:#2e7d32}
`
}
