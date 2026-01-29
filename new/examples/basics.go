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
			<h2>Counter</h2>
			<p>Count: <strong>`, b.Count, `</strong></p>
			`, p.Html(`<button>-1</button>`).WithOn("click", b.Decrement), `
			`, p.Html(`<button>+1</button>`).WithOn("click", b.Increment), `
			`, p.Html(`<button>+5</button>`).WithOn("click", func() { b.Add(5) }), `
			`, p.Html(`<button>Double</button>`).WithOn("click", b.Double), `
			`, p.Html(`<button>Reset</button>`).WithOn("click", b.Reset), `
		</section>

		<section>
			<h2>Conditionals</h2>
			<p>Score: `, b.Score, `</p>
			`, p.If(b.Score.Ge(90),
		p.Html(`<p class="grade a">Grade: A - Excellent!</p>`),
	).ElseIf(b.Score.Ge(80),
		p.Html(`<p class="grade b">Grade: B - Good</p>`),
	).ElseIf(b.Score.Ge(70),
		p.Html(`<p class="grade c">Grade: C - Average</p>`),
	).ElseIf(b.Score.Ge(60),
		p.Html(`<p class="grade d">Grade: D - Below Average</p>`),
	).Else(
		p.Html(`<p class="grade f">Grade: F - Failing</p>`),
	), `
			<div class="buttons">
				`, p.Html(`<button>A</button>`).WithOn("click", func() { b.SetScore(95) }), `
				`, p.Html(`<button>B</button>`).WithOn("click", func() { b.SetScore(85) }), `
				`, p.Html(`<button>C</button>`).WithOn("click", func() { b.SetScore(75) }), `
				`, p.Html(`<button>D</button>`).WithOn("click", func() { b.SetScore(65) }), `
				`, p.Html(`<button>F</button>`).WithOn("click", func() { b.SetScore(50) }), `
			</div>
		</section>

		<section>
			<h2>Two-Way Binding</h2>
			<label>Your name: `, p.BindValue(`<input type="text" placeholder="Enter name">`, b.Name), `</label>
			<p>Hello, `, p.Bind(b.Name), `!</p>
		</section>

		<section>
			<h2>Checkbox Binding</h2>
			<label>`, p.BindChecked(`<input type="checkbox">`, b.DarkMode), ` Dark Mode</label>
			`, p.Html(`<div>This box uses dark mode styling when checked.</div>`).AttrIf("class", p.IsTrue(b.DarkMode), "dark"), `
		</section>

		<section>
			<h2>Form</h2>
			`, p.Html(`<form>
				<label>Name: `, p.BindValue(`<input type="text" placeholder="Your name">`, b.Name), `</label>
				<label>`, p.BindChecked(`<input type="checkbox">`, b.Agreed), ` I agree to the terms</label>
				<button type="submit">Submit</button>
			</form>`).WithOn("submit", b.Submit).PreventDefault(), `
			<p class="message">`, p.Bind(b.Message), `</p>
		</section>
	</div>`)
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
