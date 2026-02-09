package main

import (
	"strings"

	p "preveltekit"
)

// Derived1 creates a store computed from one source store.
func Derived1[A, R any](a *p.Store[A], fn func(A) R) *p.Store[R] {
	out := p.New(fn(a.Get()))
	a.OnChange(func(_ A) { out.Set(fn(a.Get())) })
	return out
}

// Derived2 creates a store computed from two source stores.
func Derived2[A, B, R any](a *p.Store[A], b *p.Store[B], fn func(A, B) R) *p.Store[R] {
	out := p.New(fn(a.Get(), b.Get()))
	a.OnChange(func(_ A) { out.Set(fn(a.Get(), b.Get())) })
	b.OnChange(func(_ B) { out.Set(fn(a.Get(), b.Get())) })
	return out
}

// Derived3 creates a store computed from three source stores.
func Derived3[A, B, C, R any](a *p.Store[A], b *p.Store[B], c *p.Store[C], fn func(A, B, C) R) *p.Store[R] {
	out := p.New(fn(a.Get(), b.Get(), c.Get()))
	a.OnChange(func(_ A) { out.Set(fn(a.Get(), b.Get(), c.Get())) })
	b.OnChange(func(_ B) { out.Set(fn(a.Get(), b.Get(), c.Get())) })
	c.OnChange(func(_ C) { out.Set(fn(a.Get(), b.Get(), c.Get())) })
	return out
}

type Derived struct {
	Name      *p.Store[string]
	Uppercase *p.Store[string]

	First    *p.Store[string]
	Last     *p.Store[string]
	FullName *p.Store[string]

	Age     *p.Store[int]
	Summary *p.Store[string]
}

func (d *Derived) New() p.Component {
	name := p.New("hello")
	uppercase := Derived1(name, strings.ToUpper)

	first := p.New("John")
	last := p.New("Doe")
	fullName := Derived2(first, last, func(f, l string) string {
		return f + " " + l
	})

	age := p.New(30)
	summary := Derived3(first, last, age, func(f, l string, a int) string {
		return f + " " + l + ", age " + p.Itoa(a)
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
		<section>
			<h2>Code</h2>
			<pre class="code">// Derived1: one source store -> computed store
uppercase := Derived1(name, strings.ToUpper)

// Derived2: two source stores -> computed store
fullName := Derived2(first, last, func(f, l string) string {
    return f + " " + l
})

// Derived3: three source stores -> computed store
summary := Derived3(first, last, age, func(f, l string, a int) string {
    return f + " " + l + ", age " + p.Itoa(a)
})

// implementation pattern (using OnChange):
func Derived1[A, R any](a *Store[A], fn func(A) R) *Store[R] {
    out := p.New(fn(a.Get()))
    a.OnChange(func(_ A) { out.Set(fn(a.Get())) })
    return out
}</pre>
		</section>
	</div>`)
}

func (d *Derived) Style() string {
	return `
.demo label{display:block;margin:8px 0}
.demo pre.code{background:#1a1a2e;color:#e0e0e0;font-size:12px;margin-top:12px}
`
}
