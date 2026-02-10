package main

import p "github.com/tbocek/preveltekit"

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
	return p.Html(`<div class="demo">
		<h1>Lists</h1>

		<section>
			<h2>List Operations</h2>
			<p>Items: <strong>`, l.Items.Len(), `</strong></p>

			<div class="input-row">`,
		p.Html(`<input type="text" placeholder="New item name">`).Bind(l.NewItem),
		`</div>

			<div class="button-group">
				<h3>Add</h3>`,
		p.Html(`<button>Prepend</button>`).On("click", l.PrependItem),
		p.Html(`<button>Insert Middle</button>`).On("click", l.InsertMiddle),
		p.Html(`<button>Append</button>`).On("click", l.AddItem),
		`</div>

			<div class="button-group">
				<h3>Remove</h3>`,
		p.Html(`<button>First</button>`).On("click", l.RemoveFirst),
		p.Html(`<button>Middle</button>`).On("click", l.RemoveMiddle),
		p.Html(`<button>Last</button>`).On("click", l.RemoveLast),
		p.Html(`<button>Clear All</button>`).On("click", l.ClearAll),
		`</div>

			<div class="button-group">
				<h3>Replace All (simulates fetch)</h3>`,
		p.Html(`<button>Load Fruits</button>`).On("click", l.LoadFruits),
		p.Html(`<button>Load Numbers</button>`).On("click", l.LoadNumbers),
		`</div>

			<div class="list-container">
				<h3>Current Items</h3>`,
		p.If(p.Cond(func() bool { return l.Items.Len().Get() > 0 }, l.Items.Len()),
			p.Html(`<ul>`,
				p.Each(l.Items, func(item string, i int) p.Node {
					return p.Html(`<li><span class="index">`, p.Itoa(i), `</span> `, item, `</li>`)
				}),
				`</ul>`),
		).Else(
			p.Html(`<p class="empty">No items in list</p>`),
		), `
			</div>
		</section>

		<section>
			<h2>Code</h2>
			<pre class="code">// create a reactive list
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
    return p.Html(`+"`"+`&lt;li>`+"`"+`, item, `+"`"+`&lt;/li>`+"`"+`)
}).Else(
    p.Html(`+"`"+`&lt;p>No items&lt;/p>`+"`"+`),
)

// subscribe to changes
Items.OnChange(func(items []string) { ... })</pre>
		</section>
	</div>`)
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
