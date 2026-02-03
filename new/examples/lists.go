package main

import p "preveltekit"

type Lists struct {
	Items   *p.List[string]
	NewItem *p.Store[string]
}

func (l *Lists) New() p.Component {
	items := p.NewList[string]()
	items.Set([]string{"Apple", "Banana", "Cherry"})
	return &Lists{
		Items:   items,
		NewItem: p.New("lists.NewItem", ""),
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
	return p.Html(`<div class="demo">
		<h1>Lists</h1>

		<section>
			<h2>List Operations</h2>
			<p>Items: <strong>`, l.Items.Len(), `</strong></p>

			<div class="input-row">
				`, p.BindValue(`<input type="text" placeholder="New item name">`, l.NewItem), `
			</div>

			<div class="button-group">
				<h3>Add</h3>
				`, p.Html(`<button>Prepend</button>`).WithOn("click", "lists.PrependItem", l.PrependItem), `
				`, p.Html(`<button>Insert Middle</button>`).WithOn("click", "lists.InsertMiddle", l.InsertMiddle), `
				`, p.Html(`<button>Append</button>`).WithOn("click", "lists.AddItem", l.AddItem), `
			</div>

			<div class="button-group">
				<h3>Remove</h3>
				`, p.Html(`<button>First</button>`).WithOn("click", "lists.RemoveFirst", l.RemoveFirst), `
				`, p.Html(`<button>Middle</button>`).WithOn("click", "lists.RemoveMiddle", l.RemoveMiddle), `
				`, p.Html(`<button>Last</button>`).WithOn("click", "lists.RemoveLast", l.RemoveLast), `
				`, p.Html(`<button>Clear All</button>`).WithOn("click", "lists.ClearAll", l.ClearAll), `
			</div>

			<div class="button-group">
				<h3>Replace All (simulates fetch)</h3>
				`, p.Html(`<button>Load Fruits</button>`).WithOn("click", "lists.LoadFruits", l.LoadFruits), `
				`, p.Html(`<button>Load Numbers</button>`).WithOn("click", "lists.LoadNumbers", l.LoadNumbers), `
			</div>

			<div class="list-container">
				<h3>Current Items</h3>`,
		p.If(l.Items.Len().Gt(0),
			p.Html(`<ul>`,
				p.Each(l.Items, func(item string, i int) p.Node {
					return p.Html(`<li><span class="index">`, itoa(i), `</span> `, item, `</li>`)
				}),
				`</ul>`),
		).Else(
			p.Html(`<p class="empty">No items in list</p>`),
		),
		`</div>
		</section>
	</div>`)
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
