package main

import (
	"strings"

	p "preveltekit"
)

type Derived struct {
	// Derived from one source: uppercase transform
	Name      *p.Store[string]
	Uppercase *p.Store[string]

	// Derived from two sources: full name
	First    *p.Store[string]
	Last     *p.Store[string]
	FullName *p.Store[string]

	// Derived from three sources: summary
	Age     *p.Store[int]
	Summary *p.Store[string]
}

func (d *Derived) New() p.Component {
	name := p.New("hello")
	uppercase := p.New(strings.ToUpper(name.Get()))
	name.OnChange(func(v string) { uppercase.Set(strings.ToUpper(v)) })

	first := p.New("John")
	last := p.New("Doe")
	fullName := p.New(first.Get() + " " + last.Get())
	first.OnChange(func(f string) { fullName.Set(f + " " + last.Get()) })
	last.OnChange(func(l string) { fullName.Set(first.Get() + " " + l) })

	age := p.New(30)
	mkSummary := func() string {
		return first.Get() + " " + last.Get() + ", age " + ditoa(age.Get())
	}
	summary := p.New(mkSummary())
	first.OnChange(func(_ string) { summary.Set(mkSummary()) })
	last.OnChange(func(_ string) { summary.Set(mkSummary()) })
	age.OnChange(func(_ int) { summary.Set(mkSummary()) })

	return &Derived{
		Name:      name,
		Uppercase: uppercase,
		First:     first,
		Last:      last,
		FullName:  fullName,
		Age:       age,
		Summary:   summary,
	}
}

func (d *Derived) IncrementAge() {
	d.Age.Update(func(v int) int { return v + 1 })
}

func (d *Derived) DecrementAge() {
	d.Age.Update(func(v int) int { return v - 1 })
}

func (d *Derived) Render() p.Node {
	return p.Html(`<div class="demo">
		<h1>Derived Stores</h1>

		<section>
			<h2>Single Source</h2>
			<p class="hint">A store derived from one source via OnChange.</p>
			<label>Input: `, p.Html(`<input type="text">`).Bind(d.Name), `</label>
			<p>Uppercase: <strong>`, d.Uppercase, `</strong></p>
		</section>

		<section>
			<h2>Two Sources</h2>
			<p class="hint">A store derived from two sources. Edit first or last name.</p>
			<label>First: `, p.Html(`<input type="text">`).Bind(d.First), `</label>
			<label>Last: `, p.Html(`<input type="text">`).Bind(d.Last), `</label>
			<p>Full name: <strong>`, d.FullName, `</strong></p>
		</section>

		<section>
			<h2>Three Sources</h2>
			<p class="hint">A store derived from three sources. Change any input to update the summary.</p>
			<p>Age: <strong>`, d.Age, `</strong></p>
			`, p.Html(`<button>-1</button>`).On("click", d.DecrementAge), `
			`, p.Html(`<button>+1</button>`).On("click", d.IncrementAge), `
			<p>Summary: <strong>`, d.Summary, `</strong></p>
		</section>
	</div>`)
}

func (d *Derived) Style() string {
	return `
.demo label{display:block;margin:8px 0}
`
}

func ditoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	buf := [20]byte{}
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
