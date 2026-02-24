package main

import p "github.com/tbocek/preveltekit/v2"

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
	return p.Div(p.Attr("class", "demo"),
		p.H1("Storage"),

		p.Section(
			p.H2("Persisted Theme (Auto-sync)"),
			p.P("Theme preference is automatically saved to localStorage."),
			p.P("Current theme: ", p.Strong(s.Theme.Store)),
			p.Div(p.Attr("class", "buttons"),
				p.Button("Light").On("click", s.SetLight),
				p.Button("Dark").On("click", s.SetDark),
			),
			p.P(p.Attr("class", "hint"), "Refresh the page - theme will persist!"),
		),

		p.Section(
			p.H2("Manual Storage (Notes)"),
			p.P("Notes are saved manually when you click Save."),
			p.Textarea(p.Attr("placeholder", "Type your notes here...")).Bind(s.Notes),
			p.Div(p.Attr("class", "buttons"),
				p.Button("Save Notes").On("click", s.SaveNotes),
				p.Button("Clear Notes").On("click", s.ClearNotes),
			),
		),

		p.Section(
			p.H2("Clear All Storage"),
			p.Button(p.Attr("class", "danger"), "Clear All Storage").On("click", s.ClearAll),
		),

		p.P(p.Attr("class", "status"), s.Status),

		p.Section(
			p.H2("Code"),
			p.Pre(p.Attr("class", "code"), `// auto-persisted store (syncs to localStorage on every Set)
Theme := p.NewLocalStore("theme", "light")
Theme.Set("dark") // automatically saved

// read the store value (it's a *Store[string] inside)
Theme.Store // use in Render like any store

// manual localStorage API
p.SetStorage("notes", "hello")
saved := p.GetStorage("notes")
p.RemoveStorage("notes")
p.ClearStorage()`),
		),
	)
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
