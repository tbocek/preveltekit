package main

import p "github.com/tbocek/preveltekit/v2"

type Basics struct {
	Count    *p.Store[int]
	Name     *p.Store[string]
	Message  *p.Store[string]
	DarkMode *p.Store[bool]
	Agreed   *p.Store[bool]
	Score    *p.Store[int]
	RawHTML  *p.Store[string]
	Age      *p.Store[int]
}

func (b *Basics) New() p.Component {
	return &Basics{
		Count:    p.New(0),
		Name:     p.New(""),
		Message:  p.New("Fill out the form above"),
		DarkMode: p.New(false),
		Agreed:   p.New(false),
		Score:    p.New(75),
		RawHTML:  p.New("<em>Hello</em> <strong>World</strong>"),
		Age:      p.New(25),
	}
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

func (b *Basics) Double() {
	b.Count.Update(func(v int) int { return v * 2 })
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
	return p.Div(p.Attr("class", "demo"),
		p.H1("Basics"),

		// Counter section
		p.Section(
			p.H2("Counter — Store, Set, Update, On"),
			p.P("Count: ", p.Strong(b.Count)),
			p.Button("-1").On("click", b.Decrement),
			p.Button("+1").On("click", b.Increment),
			p.Button("+5").On("click", func() { b.Add(5) }),
			p.Button("Double").On("click", b.Double),
			p.Button("Reset").On("click", b.Reset),
			p.Pre(p.Attr("class", "code"), `Count := p.New(0)                           // create store
Count.Set(5)                                 // set value
Count.Update(func(v int) int { return v+1 }) // transform

// embed in HTML — auto-updates in the DOM
p.Html(`+"`"+`<p>Count: <strong>`+"`"+`, Count, `+"`"+`</strong></p>`+"`"+`)

// attach event handler
p.Html(`+"`"+`<button>+1</button>`+"`"+`).On("click", Increment)`),
		),

		// Conditionals section
		p.Section(
			p.H2("Conditionals — If / ElseIf / Else"),
			p.P("Score: ", b.Score),
			p.If(p.Cond(func() bool { return b.Score.Get() >= 90 }, b.Score),
				p.P(p.Attr("class", "grade a"), "Grade: A - Excellent!"),
			).ElseIf(p.Cond(func() bool { return b.Score.Get() >= 80 }, b.Score),
				p.P(p.Attr("class", "grade b"), "Grade: B - Good"),
			).ElseIf(p.Cond(func() bool { return b.Score.Get() >= 70 }, b.Score),
				p.P(p.Attr("class", "grade c"), "Grade: C - Average"),
			).ElseIf(p.Cond(func() bool { return b.Score.Get() >= 60 }, b.Score),
				p.P(p.Attr("class", "grade d"), "Grade: D - Below Average"),
			).Else(
				p.P(p.Attr("class", "grade f"), "Grade: F - Failing"),
			),
			p.Div(p.Attr("class", "buttons"),
				p.Button("A").On("click", func() { b.SetScore(95) }),
				p.Button("B").On("click", func() { b.SetScore(85) }),
				p.Button("C").On("click", func() { b.SetScore(75) }),
				p.Button("D").On("click", func() { b.SetScore(65) }),
				p.Button("F").On("click", func() { b.SetScore(50) }),
			),
			p.Pre(p.Attr("class", "code"), `p.If(p.Cond(func() bool { return Score.Get() >= 90 }, Score),
    p.Html(`+"`"+`<p>Grade: A</p>`+"`"+`),
).ElseIf(p.Cond(func() bool { return Score.Get() >= 80 }, Score),
    p.Html(`+"`"+`<p>Grade: B</p>`+"`"+`),
).Else(
    p.Html(`+"`"+`<p>Grade: F</p>`+"`"+`),
)

// p.Cond(func() bool, ...stores) — any logic you want`),
		),

		// Two-Way Binding — String
		p.Section(
			p.H2("Two-Way Binding — String"),
			p.Label("Your name: ", p.Input(p.Attr("type", "text"), p.Attr("placeholder", "Enter name")).Bind(b.Name)),
			p.P("Hello, ", b.Name, "!"),
			p.Pre(p.Attr("class", "code"), `Name := p.New("")
p.Html(`+"`"+`<input type="text">`+"`"+`).Bind(Name)`),
		),

		// Two-Way Binding — Int
		p.Section(
			p.H2("Two-Way Binding — Int"),
			p.Label("Age: ", p.Input(p.Attr("type", "text"), p.Attr("placeholder", "Enter age")).Bind(b.Age)),
			p.P("You are ", b.Age, " years old."),
			p.Pre(p.Attr("class", "code"), `Age := p.New(25)
p.Html(`+"`"+`<input type="text">`+"`"+`).Bind(Age) // *Store[int] binding`),
		),

		// Checkbox Binding — Bool
		p.Section(
			p.H2("Checkbox Binding — Bool"),
			p.Label(p.Input(p.Attr("type", "checkbox")).Bind(b.DarkMode), " Dark Mode"),
			p.Div("This box uses dark mode styling when checked.").AttrIf("class", p.Cond(func() bool { return b.DarkMode.Get() }, b.DarkMode), "dark"),
			p.Pre(p.Attr("class", "code"), `DarkMode := p.New(false)
p.Html(`+"`"+`<input type="checkbox">`+"`"+`).Bind(DarkMode) // *Store[bool]

// AttrIf: conditionally add a class
p.Html(`+"`"+`<div>...</div>`+"`"+`).AttrIf("class", p.Cond(func() bool { return DarkMode.Get() }, DarkMode), "dark")`),
		),

		// Form — PreventDefault
		p.Section(
			p.H2("Form — PreventDefault"),
			p.Form(
				p.Label("Name: ", p.Input(p.Attr("type", "text"), p.Attr("placeholder", "Your name")).Bind(b.Name)),
				p.Label(p.Input(p.Attr("type", "checkbox")).Bind(b.Agreed), " I agree to the terms"),
				p.Button(p.Attr("type", "submit"), "Submit"),
			).On("submit", b.Submit).PreventDefault(),
			p.P(p.Attr("class", "message"), b.Message),
			p.Pre(p.Attr("class", "code"), `p.Html(`+"`"+`<form>...</form>`+"`"+`).On("submit", Submit).PreventDefault()

// also available: .StopPropagation()`),
		),

		// StopPropagation
		p.Section(
			p.H2("StopPropagation"),
			p.P("Click the inner button — only inner handler fires, not outer:"),
			p.Div(p.Attr("class", "outer-click"),
				"Outer (click me)",
				p.Div(p.Button("Inner (click me)").On("click", func() {
					b.Message.Set("Inner clicked!")
				}).StopPropagation()),
			).On("click", func() {
				b.Message.Set("Outer clicked!")
			}),
			p.P(p.Attr("class", "message"), b.Message),
			p.Pre(p.Attr("class", "code"), `// inner button stops event from reaching outer div
p.Html(`+"`"+`<button>Inner</button>`+"`"+`).On("click", handler).StopPropagation()`),
		),

		// BindAsHTML — Raw HTML Rendering
		p.Section(
			p.H2("BindAsHTML — Raw HTML Rendering"),
			p.Label("HTML: ", p.Input(p.Attr("type", "text")).Bind(b.RawHTML)),
			p.P("Rendered: ", p.BindAsHTML(b.RawHTML)),
			p.Pre(p.Attr("class", "code"), `RawHTML := p.New("<em>Hello</em> <strong>World</strong>")
p.BindAsHTML(RawHTML) // renders as innerHTML (not escaped)`),
		),
	)
}

func (b *Basics) Style() string {
	return `
.demo label{display:block;margin:8px 0}
.demo pre.code{background:#1a1a2e;color:#e0e0e0;font-size:12px;margin-top:12px}
.grade{padding:10px;border-radius:4px;font-weight:700}
.grade.a{background:#d4edda;color:#155724}
.grade.b{background:#cce5ff;color:#004085}
.grade.c{background:#fff3cd;color:#856404}
.grade.d{background:#ffe5d0;color:#854027}
.grade.f{background:#f8d7da;color:#721c24}
.dark{background:#333;color:#fff;padding:10px;border-radius:4px;margin-top:10px}
.message{padding:10px;background:#e7f3ff;border-radius:4px}
.outer-click{padding:15px;background:#fff3cd;border:1px solid #ffc107;border-radius:4px;cursor:pointer}
.outer-click div{margin-top:10px}
`
}
