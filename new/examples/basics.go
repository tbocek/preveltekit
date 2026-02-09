package main

import p "preveltekit"

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
	return p.Html(`<div class="demo">
		<h1>Basics</h1>

		<section>
			<h2>Counter — Store, Set, Update, On</h2>
			<p>Count: <strong>`, b.Count, `</strong></p>`,
		p.Html(`<button>-1</button>`).On("click", b.Decrement),
		p.Html(`<button>+1</button>`).On("click", b.Increment),
		p.Html(`<button>+5</button>`).On("click", func() { b.Add(5) }),
		p.Html(`<button>Double</button>`).On("click", b.Double),
		p.Html(`<button>Reset</button>`).On("click", b.Reset),
		`<pre class="code">Count := p.New(0)                           // create store
Count.Set(5)                                 // set value
Count.Update(func(v int) int { return v+1 }) // transform

// embed in HTML — auto-updates in the DOM
p.Html(`+"`"+`&lt;p>Count: &lt;strong>`+"`"+`, Count, `+"`"+`&lt;/strong>&lt;/p>`+"`"+`)

// attach event handler
p.Html(`+"`"+`&lt;button>+1&lt;/button>`+"`"+`).On("click", Increment)</pre>
		</section>

		<section>
			<h2>Conditionals — If / ElseIf / Else</h2>
			<p>Score: `, b.Score, `</p>
			`,
		p.If(p.Cond(func() bool { return b.Score.Get() >= 90 }, b.Score),
			p.Html(`<p class="grade a">Grade: A - Excellent!</p>`),
		).ElseIf(p.Cond(func() bool { return b.Score.Get() >= 80 }, b.Score),
			p.Html(`<p class="grade b">Grade: B - Good</p>`),
		).ElseIf(p.Cond(func() bool { return b.Score.Get() >= 70 }, b.Score),
			p.Html(`<p class="grade c">Grade: C - Average</p>`),
		).ElseIf(p.Cond(func() bool { return b.Score.Get() >= 60 }, b.Score),
			p.Html(`<p class="grade d">Grade: D - Below Average</p>`),
		).Else(
			p.Html(`<p class="grade f">Grade: F - Failing</p>`),
		), `
			<div class="buttons">`,
		p.Html(`<button>A</button>`).On("click", func() { b.SetScore(95) }),
		p.Html(`<button>B</button>`).On("click", func() { b.SetScore(85) }),
		p.Html(`<button>C</button>`).On("click", func() { b.SetScore(75) }),
		p.Html(`<button>D</button>`).On("click", func() { b.SetScore(65) }),
		p.Html(`<button>F</button>`).On("click", func() { b.SetScore(50) }),
		`</div>
			<pre class="code">p.If(p.Cond(func() bool { return Score.Get() >= 90 }, Score),
    p.Html(`+"`"+`&lt;p>Grade: A&lt;/p>`+"`"+`),
).ElseIf(p.Cond(func() bool { return Score.Get() >= 80 }, Score),
    p.Html(`+"`"+`&lt;p>Grade: B&lt;/p>`+"`"+`),
).Else(
    p.Html(`+"`"+`&lt;p>Grade: F&lt;/p>`+"`"+`),
)

// p.Cond(func() bool, ...stores) — any logic you want</pre>
		</section>

		<section>
			<h2>Two-Way Binding — String</h2>
			<label>Your name: `, p.Html(`<input type="text" placeholder="Enter name">`).Bind(b.Name), `</label>
			<p>Hello, `, b.Name, `!</p>
			<pre class="code">Name := p.New("")
p.Html(`+"`"+`&lt;input type="text">`+"`"+`).Bind(Name)</pre>
		</section>

		<section>
			<h2>Two-Way Binding — Int</h2>
			<label>Age: `, p.Html(`<input type="text" placeholder="Enter age">`).Bind(b.Age), `</label>
			<p>You are `, b.Age, ` years old.</p>
			<pre class="code">Age := p.New(25)
p.Html(`+"`"+`&lt;input type="text">`+"`"+`).Bind(Age) // *Store[int] binding</pre>
		</section>

		<section>
			<h2>Checkbox Binding — Bool</h2>
			<label>`, p.Html(`<input type="checkbox">`).Bind(b.DarkMode), ` Dark Mode</label>
			`, p.Html(`<div>This box uses dark mode styling when checked.</div>`).AttrIf("class", p.Cond(func() bool { return b.DarkMode.Get() }, b.DarkMode), "dark"), `
			<pre class="code">DarkMode := p.New(false)
p.Html(`+"`"+`&lt;input type="checkbox">`+"`"+`).Bind(DarkMode) // *Store[bool]

// AttrIf: conditionally add a class
p.Html(`+"`"+`&lt;div>...&lt;/div>`+"`"+`).AttrIf("class", p.Cond(func() bool { return DarkMode.Get() }, DarkMode), "dark")</pre>
		</section>

		<section>
			<h2>Form — PreventDefault</h2>
			`, p.Html(`<form>
				<label>Name: `, p.Html(`<input type="text" placeholder="Your name">`).Bind(b.Name), `</label>
				<label>`, p.Html(`<input type="checkbox">`).Bind(b.Agreed), ` I agree to the terms</label>
				<button type="submit">Submit</button>
			</form>`).On("submit", b.Submit).PreventDefault(), `
			<p class="message">`, b.Message, `</p>
			<pre class="code">p.Html(`+"`"+`&lt;form>...&lt;/form>`+"`"+`).On("submit", Submit).PreventDefault()

// also available: .StopPropagation()</pre>
		</section>

		<section>
			<h2>StopPropagation</h2>
			<p>Click the inner button — only inner handler fires, not outer:</p>
			`, p.Html(`<div class="outer-click">
				Outer (click me)
				<div>`, p.Html(`<button>Inner (click me)</button>`).On("click", func() {
			b.Message.Set("Inner clicked!")
		}).StopPropagation(), `</div>
			</div>`).On("click", func() {
			b.Message.Set("Outer clicked!")
		}), `
			<p class="message">`, b.Message, `</p>
			<pre class="code">// inner button stops event from reaching outer div
p.Html(`+"`"+`&lt;button>Inner&lt;/button>`+"`"+`).On("click", handler).StopPropagation()</pre>
		</section>

		<section>
			<h2>BindAsHTML — Raw HTML Rendering</h2>
			<label>HTML: `, p.Html(`<input type="text">`).Bind(b.RawHTML), `</label>
			<p>Rendered: `, p.BindAsHTML(b.RawHTML), `</p>
			<pre class="code">RawHTML := p.New("&lt;em>Hello&lt;/em> &lt;strong>World&lt;/strong>")
p.BindAsHTML(RawHTML) // renders as innerHTML (not escaped)</pre>
		</section>

	</div>`)
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
