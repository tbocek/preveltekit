package main

import p "preveltekit"

type Basics struct {
	Count    *p.Store[int]
	Name     *p.Store[string]
	Message  *p.Store[string]
	DarkMode *p.Store[bool]
	Agreed   *p.Store[bool]
	Score    *p.Store[int]
}

func (b *Basics) OnMount() {
	b.Count.Set(0)
	b.Name.Set("")
	b.Message.Set("Fill out the form above")
	b.DarkMode.Set(false)
	b.Agreed.Set(false)
	b.Score.Set(75)
}

func (b *Basics) Increment() {
	b.Count.Update(func(v int) int { return v + 1 })
}

func (b *Basics) Decrement() {
	b.Count.Update(func(v int) int { return v - 1 })
}

func (b *Basics) Add(n int) {
	b.Count.Update(func(v int) int { return v + n })
}

func (b *Basics) Reset() {
	b.Count.Set(0)
}

func (b *Basics) SetScore(n int) {
	b.Score.Set(n)
}

func (b *Basics) Submit() {
	name := b.Name.Get()
	if name == "" {
		b.Message.Set("Please enter your name")
	} else if !b.Agreed.Get() {
		b.Message.Set("Please agree to terms")
	} else {
		b.Message.Set("Welcome, " + name + "!")
	}
}

func (b *Basics) Render() p.Node {
	return p.Div(p.Class("demo"),
		p.H1("Basics"),

		// Counter section
		p.Section(
			p.H2("Counter"),
			p.P("Count: ", p.Strong(p.Bind(b.Count))),
			p.Button("-1", p.OnClick(b.Decrement)),
			p.Button("+1", p.OnClick(b.Increment)),
			p.Button("+5", p.OnClick(b.Add, 5)),
			p.Button("Double", p.OnClick(b.Add, b.Count.Get())),
			p.Button("Reset", p.OnClick(b.Reset)),
		),

		// Conditionals section
		p.Section(
			p.H2("Conditionals"),
			p.P("Score: ", p.Bind(b.Score)),
			p.If(b.Score.Ge(90),
				p.P(p.Class("grade", "a"), "Grade: A - Excellent!"),
			).ElseIf(b.Score.Ge(80),
				p.P(p.Class("grade", "b"), "Grade: B - Good"),
			).ElseIf(b.Score.Ge(70),
				p.P(p.Class("grade", "c"), "Grade: C - Average"),
			).ElseIf(b.Score.Ge(60),
				p.P(p.Class("grade", "d"), "Grade: D - Below Average"),
			).Else(
				p.P(p.Class("grade", "f"), "Grade: F - Failing"),
			),
			p.Div(p.Class("buttons"),
				p.Button("A", p.OnClick(b.SetScore, 95)),
				p.Button("B", p.OnClick(b.SetScore, 85)),
				p.Button("C", p.OnClick(b.SetScore, 75)),
				p.Button("D", p.OnClick(b.SetScore, 65)),
				p.Button("F", p.OnClick(b.SetScore, 50)),
			),
		),

		// Two-way binding section
		p.Section(
			p.H2("Two-Way Binding"),
			p.Label("Your name: ",
				p.Input(p.Type("text"), p.BindValue(b.Name), p.Placeholder("Enter name")),
			),
			p.P("Hello, ", p.Bind(b.Name), "!"),
		),

		// Checkbox binding section
		p.Section(
			p.H2("Checkbox Binding"),
			p.Label(
				p.Input(p.Type("checkbox"), p.BindChecked(b.DarkMode)),
				" Dark Mode",
			),
			p.Div(p.ClassIf("dark", p.IsTrue(b.DarkMode)),
				"This box uses dark mode styling when checked.",
			),
		),

		// Form section
		p.Section(
			p.H2("Form"),
			p.Form(p.OnSubmit(b.Submit).PreventDefault(),
				p.Label("Name: ",
					p.Input(p.Type("text"), p.BindValue(b.Name), p.Placeholder("Your name")),
				),
				p.Label(
					p.Input(p.Type("checkbox"), p.BindChecked(b.Agreed)),
					" I agree to the terms",
				),
				p.Button("Submit", p.Type("submit")),
			),
			p.P(p.Class("message"), p.Bind(b.Message)),
		),
	)
}

func (b *Basics) Style() string {
	return `
.demo label{display:block;margin:8px 0}
.grade{padding:10px;border-radius:4px;font-weight:700}
.grade.a{background:#d4edda;color:#155724}
.grade.b{background:#cce5ff;color:#004085}
.grade.c{background:#fff3cd;color:#856404}
.grade.d{background:#ffe5d0;color:#854027}
.grade.f{background:#f8d7da;color:#721c24}
.dark{background:#333;color:#fff;padding:10px;border-radius:4px;margin-top:10px}
.message{padding:10px;background:#e7f3ff;border-radius:4px}
`
}

func (b *Basics) HandleEvent(method string, args string) {
	switch method {
	case "Increment":
		b.Increment()
	case "Decrement":
		b.Decrement()
	case "Add":
		b.Add(atoi(args))
	case "Reset":
		b.Reset()
	case "SetScore":
		b.SetScore(atoi(args))
	case "Submit":
		b.Submit()
	}
}

func atoi(s string) int {
	n := 0
	neg := false
	for i, c := range s {
		if c == '-' && i == 0 {
			neg = true
			continue
		}
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	if neg {
		return -n
	}
	return n
}
