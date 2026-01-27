package main

import p "preveltekit"

type Lists struct {
	Items     *p.List[string]
	NewItem   *p.Store[string]
	ItemCount *p.Store[int]
}

func (l *Lists) OnMount() {
	if l.Items.Len() == 0 {
		l.Items.Set([]string{"Apple", "Banana", "Cherry"})
		l.ItemCount.Set(3)
	}
	l.NewItem.Set("")
}

func (l *Lists) AddItem() {
	item := l.NewItem.Get()
	if item == "" {
		return
	}
	l.Items.Append(item)
	l.NewItem.Set("")
	l.updateCount()
}

func (l *Lists) PrependItem() {
	item := l.NewItem.Get()
	if item == "" {
		return
	}
	items := l.Items.Get()
	l.Items.Set(append([]string{item}, items...))
	l.NewItem.Set("")
	l.updateCount()
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
	l.updateCount()
}

func (l *Lists) RemoveFirst() {
	if l.Items.Len() > 0 {
		l.Items.RemoveAt(0)
		l.updateCount()
	}
}

func (l *Lists) RemoveLast() {
	length := l.Items.Len()
	if length > 0 {
		l.Items.RemoveAt(length - 1)
		l.updateCount()
	}
}

func (l *Lists) RemoveMiddle() {
	length := l.Items.Len()
	if length > 0 {
		l.Items.RemoveAt(length / 2)
		l.updateCount()
	}
}

func (l *Lists) ClearAll() {
	l.Items.Clear()
	l.updateCount()
}

func (l *Lists) LoadFruits() {
	l.Items.Set([]string{"Mango", "Pineapple", "Papaya", "Guava"})
	l.updateCount()
}

func (l *Lists) LoadNumbers() {
	l.Items.Set([]string{"One", "Two", "Three", "Four", "Five"})
	l.updateCount()
}

func (l *Lists) updateCount() {
	l.ItemCount.Set(l.Items.Len())
}

func (l *Lists) Render() p.Node {
	return p.Div(p.Class("demo"),
		p.H1("Lists"),

		p.Section(
			p.H2("List Operations"),
			p.P("Items: ", p.Strong(p.Bind(l.ItemCount))),

			p.Div(p.Class("input-row"),
				p.Input(p.Type("text"), p.BindValue(l.NewItem), p.Placeholder("New item name")),
			),

			p.Div(p.Class("button-group"),
				p.H3("Add"),
				p.Button("Prepend", p.OnClick(l.PrependItem)),
				p.Button("Insert Middle", p.OnClick(l.InsertMiddle)),
				p.Button("Append", p.OnClick(l.AddItem)),
			),

			p.Div(p.Class("button-group"),
				p.H3("Remove"),
				p.Button("First", p.OnClick(l.RemoveFirst)),
				p.Button("Middle", p.OnClick(l.RemoveMiddle)),
				p.Button("Last", p.OnClick(l.RemoveLast)),
				p.Button("Clear All", p.OnClick(l.ClearAll)),
			),

			p.Div(p.Class("button-group"),
				p.H3("Replace All (simulates fetch)"),
				p.Button("Load Fruits", p.OnClick(l.LoadFruits)),
				p.Button("Load Numbers", p.OnClick(l.LoadNumbers)),
			),

			p.Div(p.Class("list-container"),
				p.H3("Current Items"),
				p.If(l.ItemCount.Gt(0),
					p.Ul(
						p.Each(l.Items, func(item string, i int) p.Node {
							return p.Li(
								p.Span(p.Class("index"), p.Text(itoa(i))),
								p.Text(" "+item),
							)
						}),
					),
				).Else(
					p.P(p.Class("empty"), "No items in list"),
				),
			),
		),
	)
}

func (l *Lists) Style() string {
	return `
.demo h3{margin:10px 0 5px;color:#888;font-size:.9em}
.input-row{margin:10px 0}
.button-group{margin:15px 0;padding:10px;background:#f9f9f9;border-radius:4px}
.list-container{margin-top:20px;padding:15px;background:#f0f0f0;border-radius:4px}
.list-container ul{list-style:none;padding:0;margin:0}
.list-container li{padding:8px 12px;margin:4px 0;background:#fff;border-radius:4px;border-left:3px solid #4caf50}
.index{display:inline-block;width:24px;height:24px;line-height:24px;text-align:center;background:#4caf50;color:#fff;border-radius:50%;font-size:12px;margin-right:10px}
.empty{color:#999;font-style:italic;text-align:center;padding:20px}
`
}

func (l *Lists) HandleEvent(method string, args string) {
	switch method {
	case "AddItem":
		l.AddItem()
	case "PrependItem":
		l.PrependItem()
	case "InsertMiddle":
		l.InsertMiddle()
	case "RemoveFirst":
		l.RemoveFirst()
	case "RemoveLast":
		l.RemoveLast()
	case "RemoveMiddle":
		l.RemoveMiddle()
	case "ClearAll":
		l.ClearAll()
	case "LoadFruits":
		l.LoadFruits()
	case "LoadNumbers":
		l.LoadNumbers()
	}
}
