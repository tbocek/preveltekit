package main

import (
	"strings"

	p "github.com/tbocek/preveltekit/v2"
)

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
	uppercase := name.Map(strings.ToUpper)

	first := p.New("John")
	last := p.New("Doe")
	fullName := p.Computed(func() string {
		return first.Get() + " " + last.Get()
	})

	age := p.New(30)
	summary := p.Computed(func() string {
		return first.Get() + " " + last.Get() + ", age " + p.Itoa(age.Get())
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
			p.H2("Map — Single Source"),
			p.P(p.Attr("class", "hint"), "A store computed from one source with Map. Type below to see the uppercase transform."),
			p.Label("Input: ", p.Input(p.Attr("type", "text")).Bind(d.Name)),
			p.P("Uppercase: ", p.Strong(d.Uppercase)),
		),

		p.Section(
			p.H2("Computed — Two Sources"),
			p.P(p.Attr("class", "hint"), "A store computed from two sources with Computed. Edit first or last name."),
			p.Label("First: ", p.Input(p.Attr("type", "text")).Bind(d.First)),
			p.Label("Last: ", p.Input(p.Attr("type", "text")).Bind(d.Last)),
			p.P("Full name: ", p.Strong(d.FullName)),
		),

		p.Section(
			p.H2("Computed — Three Sources"),
			p.P(p.Attr("class", "hint"), "A store computed from three sources. Change any input to update the summary."),
			p.P("Age: ", p.Strong(d.Age)),
			p.Button("-1").On("click", d.DecrementAge),
			p.Button("+1").On("click", d.IncrementAge),
			p.P("Summary: ", p.Strong(d.Summary)),
		),

		p.Section(
			p.H2("Code"),
			p.Pre(p.Attr("class", "code"),
				`// Map: one source store -> computed store (method on Store)
uppercase := name.Map(strings.ToUpper)

// Computed: any number of sources. The stores read with Get are
// tracked automatically; the store updates when one of them changes.
fullName := p.Computed(func() string {
    return first.Get() + " " + last.Get()
})

summary := p.Computed(func() string {
    return first.Get() + " " + last.Get() + ", age " + p.Itoa(age.Get())
})`),
		),
	)
}
