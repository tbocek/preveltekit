package main

import (
	"strings"

	p "preveltekit"
)

type Derived struct {
	// Derived1: uppercase transform
	Name      *p.Store[string]
	Uppercase *p.Store[string]

	// Derived2: full name from first + last
	First    *p.Store[string]
	Last     *p.Store[string]
	FullName *p.Store[string]

	// Derived3: summary from first + last + age
	Age     *p.Store[int]
	Summary *p.Store[string]
}

func (d *Derived) New() p.Component {
	name := p.New("hello")
	uppercase := p.Derived1(name, strings.ToUpper)

	first := p.New("John")
	last := p.New("Doe")
	fullName := p.Derived2(first, last, func(f, l string) string {
		return f + " " + l
	})

	age := p.New(30)
	summary := p.Derived3(first, last, age, func(f, l string, a int) string {
		return f + " " + l + ", age " + ditoa(a)
	})

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
			<h2>Derived1 — Single Source</h2>
			<p class="hint">A store computed from one source. Type below to see the uppercase transform.</p>
			<label>Input: `, p.Html(`<input type="text">`).Bind(d.Name), `</label>
			<p>Uppercase: <strong>`, d.Uppercase, `</strong></p>
		</section>

		<section>
			<h2>Derived2 — Two Sources</h2>
			<p class="hint">A store computed from two sources. Edit first or last name.</p>
			<label>First: `, p.Html(`<input type="text">`).Bind(d.First), `</label>
			<label>Last: `, p.Html(`<input type="text">`).Bind(d.Last), `</label>
			<p>Full name: <strong>`, d.FullName, `</strong></p>
		</section>

		<section>
			<h2>Derived3 — Three Sources</h2>
			<p class="hint">A store computed from three sources. Change any input to update the summary.</p>
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
