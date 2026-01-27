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
	return p.Div(p.Class("demo"),
		p.H1("Storage"),

		// Persisted Theme section
		p.Section(
			p.H2("Persisted Theme (Auto-sync)"),
			p.P("Theme preference is automatically saved to localStorage."),
			p.P("Current theme: ", p.Strong(p.Bind(s.Theme.Store))),
			p.Div(p.Class("buttons"),
				p.Button("Light", p.OnClick(s.SetLight)),
				p.Button("Dark", p.OnClick(s.SetDark)),
			),
			p.P(p.Class("hint"), "Refresh the page - theme will persist!"),
		),

		// Manual Storage section
		p.Section(
			p.H2("Manual Storage (Notes)"),
			p.P("Notes are saved manually when you click Save."),
			p.Textarea(p.BindValue(s.Notes), p.Placeholder("Type your notes here...")),
			p.Div(p.Class("buttons"),
				p.Button("Save Notes", p.OnClick(s.SaveNotes)),
				p.Button("Clear Notes", p.OnClick(s.ClearNotes)),
			),
		),

		// Clear All section
		p.Section(
			p.H2("Clear All Storage"),
			p.Button(p.Class("danger"), "Clear All Storage", p.OnClick(s.ClearAll)),
		),

		p.P(p.Class("status"), p.Bind(s.Status)),
	)
}

func (s *Storage) Style() string {
	return `
.demo textarea{height:100px}
.demo button.danger{background:#ffebee;border-color:#ef9a9a;color:#c62828}
.demo button.danger:hover{background:#ffcdd2}
.demo .status{padding:10px;background:#e8f5e9;border-radius:4px;color:#2e7d32}
`
}

func (s *Storage) HandleEvent(method string, args string) {
	switch method {
	case "SetLight":
		s.SetLight()
	case "SetDark":
		s.SetDark()
	case "SaveNotes":
		s.SaveNotes()
	case "ClearNotes":
		s.ClearNotes()
	case "ClearAll":
		s.ClearAll()
	}
}
