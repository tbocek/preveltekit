package main

import (
	"strings"

	p "github.com/tbocek/preveltekit/v2"
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
	return p.Div(p.Attr("class", "demo"),
		p.H1("Derived Stores"),

		p.Section(
			p.H2("Derived1 — Single Source"),
			p.P(p.Attr("class", "hint"), "A store computed from one source. Type below to see the uppercase transform."),
			p.Label("Input: ", p.Input(p.Attr("type", "text")).Bind(d.Name)),
			p.P("Uppercase: ", p.Strong(d.Uppercase)),
		),

		p.Section(
			p.H2("Derived2 — Two Sources"),
			p.P(p.Attr("class", "hint"), "A store computed from two sources. Edit first or last name."),
			p.Label("First: ", p.Input(p.Attr("type", "text")).Bind(d.First)),
			p.Label("Last: ", p.Input(p.Attr("type", "text")).Bind(d.Last)),
			p.P("Full name: ", p.Strong(d.FullName)),
		),

		p.Section(
			p.H2("Derived3 — Three Sources"),
			p.P(p.Attr("class", "hint"), "A store computed from three sources. Change any input to update the summary."),
			p.P("Age: ", p.Strong(d.Age)),
			p.Button("-1").On("click", d.DecrementAge),
			p.Button("+1").On("click", d.IncrementAge),
			p.P("Summary: ", p.Strong(d.Summary)),
		),

		p.Section(
			p.H2("Code"),
			p.Pre(p.Attr("class", "code"),
				`// Derived1: one source store -> computed store
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
}`),
		),
	)
}

func (d *Derived) Style() string {
	return `
.demo label{display:block;margin:8px 0}
.demo pre.code{background:#1a1a2e;color:#e0e0e0;font-size:12px;margin-top:12px}
`
}
