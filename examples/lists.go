package main

import p "github.com/tbocek/preveltekit/v2"

type Lists struct {
	Items   *p.List[string]
	NewItem *p.Store[string]
}

func (l *Lists) New() p.Component {
	items := p.NewList[string]("Apple", "Banana", "Cherry")
	return &Lists{
		Items:   items,
		NewItem: p.New(""),
	}
}

func (l *Lists) AddItem() {
	item := l.NewItem.Get()
	if item == "" {
		return
	}
	l.Items.Append(item)
	l.NewItem.Set("")
}

func (l *Lists) PrependItem() {
	item := l.NewItem.Get()
	if item == "" {
		return
	}
	items := l.Items.Get()
	l.Items.Set(append([]string{item}, items...))
	l.NewItem.Set("")
}

func (l *Lists) InsertMiddle() {
	item := l.NewItem.Get()
	if item == "" {
		return
	}
	items := l.Items.Get()
	mid := len(items) / 2
	newItems := make([]string, 0, len(items)+1)
	newItems = append(newItems, items[:mid]...)
	newItems = append(newItems, item)
	newItems = append(newItems, items[mid:]...)
	l.Items.Set(newItems)
	l.NewItem.Set("")
}

func (l *Lists) RemoveFirst() {
	if l.Items.Len().Get() > 0 {
		l.Items.RemoveAt(0)
	}
}

func (l *Lists) RemoveLast() {
	length := l.Items.Len().Get()
	if length > 0 {
		l.Items.RemoveAt(length - 1)
	}
}

func (l *Lists) RemoveMiddle() {
	length := l.Items.Len().Get()
	if length > 0 {
		l.Items.RemoveAt(length / 2)
	}
}

func (l *Lists) ClearAll() {
	l.Items.Clear()
}

func (l *Lists) LoadFruits() {
	l.Items.Set([]string{"Mango", "Pineapple", "Papaya", "Guava"})
}

func (l *Lists) LoadNumbers() {
	l.Items.Set([]string{"One", "Two", "Three", "Four", "Five"})
}

func (l *Lists) Render() p.Node {
	return p.Div(p.Attr("class", "demo"),
		p.H1("Lists"),

		p.Section(
			p.H2("List Operations"),
			p.P("Items: ", p.Strong(l.Items.Len())),

			p.Div(p.Attr("class", "input-row"),
				p.Input(p.Attr("type", "text"), p.Attr("placeholder", "New item name")).Bind(l.NewItem),
			),

			p.Div(p.Attr("class", "button-group"),
				p.H3("Add"),
				p.Button("Prepend").On("click", l.PrependItem),
				p.Button("Insert Middle").On("click", l.InsertMiddle),
				p.Button("Append").On("click", l.AddItem),
			),

			p.Div(p.Attr("class", "button-group"),
				p.H3("Remove"),
				p.Button("First").On("click", l.RemoveFirst),
				p.Button("Middle").On("click", l.RemoveMiddle),
				p.Button("Last").On("click", l.RemoveLast),
				p.Button("Clear All").On("click", l.ClearAll),
			),

			p.Div(p.Attr("class", "button-group"),
				p.H3("Replace All (simulates fetch)"),
				p.Button("Load Fruits").On("click", l.LoadFruits),
				p.Button("Load Numbers").On("click", l.LoadNumbers),
			),

			p.Div(p.Attr("class", "list-container"),
				p.H3("Current Items"),
				p.If(p.Cond(func() bool { return l.Items.Len().Get() > 0 }, l.Items.Len()),
					p.Ul(
						p.Each(l.Items, func(item string, i int) p.Node {
							return p.Li(p.Span(p.Attr("class", "index"), p.Itoa(i)), " ", item)
						}),
					),
				).Else(
					p.P(p.Attr("class", "empty"), "No items in list"),
				),
			),
		),

		p.Section(
			p.H2("Code"),
			p.Pre(p.Attr("class", "code"), `// create a reactive list
Items := p.NewList[string]("Apple", "Banana")

// mutate — triggers re-render
Items.Append("Cherry")
Items.RemoveAt(0)
Items.Set([]string{"Mango", "Papaya"})
Items.Clear()

// reactive length store (for conditions)
Items.Len()       // *Store[int]
p.Cond(func() bool { return Items.Len().Get() > 0 }, Items.Len()) // Condition

// render with Each
p.Each(Items, func(item string, i int) p.Node {
    return p.Li(item)
}).Else(
    p.P("No items"),
)

// subscribe to changes
Items.OnChange(func(items []string) { ... })`),
		),
	)
}

func (l *Lists) Style() string {
	return `
.demo h3{margin:10px 0 5px;color:#888;font-size:.9em}
.demo pre.code{background:#1a1a2e;color:#e0e0e0;font-size:12px;margin-top:12px}
.input-row{margin:10px 0}
.button-group{margin:15px 0;padding:10px;background:#f9f9f9;border-radius:4px}
.list-container{margin-top:20px;padding:15px;background:#f0f0f0;border-radius:4px}
.list-container ul{list-style:none;padding:0;margin:0}
.list-container li{padding:8px 12px;margin:4px 0;background:#fff;border-radius:4px;border-left:3px solid #4caf50}
.index{display:inline-block;width:24px;height:24px;line-height:24px;text-align:center;background:#4caf50;color:#fff;border-radius:50%;font-size:12px;margin-right:10px}
.empty{color:#999;font-style:italic;text-align:center;padding:20px}
`
}
