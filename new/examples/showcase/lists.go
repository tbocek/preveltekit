package main

import "preveltekit"

type Lists struct {
	Items     *preveltekit.List[string]
	NewItem   *preveltekit.Store[string]
	ItemCount *preveltekit.Store[int]
}

func (l *Lists) OnMount() {
	// Only initialize if empty (preserve state across navigation)
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

func (l *Lists) Template() string {
	return `<div class="demo">
	<h1>Lists</h1>

	<section>
		<h2>List Operations</h2>
		<p>Items: <strong>{ItemCount}</strong></p>

		<div class="input-row">
			<input type="text" bind:value="NewItem" placeholder="New item name">
		</div>

		<div class="button-group">
			<h3>Add</h3>
			<button @click="PrependItem()">Prepend</button>
			<button @click="InsertMiddle()">Insert Middle</button>
			<button @click="AddItem()">Append</button>
		</div>

		<div class="button-group">
			<h3>Remove</h3>
			<button @click="RemoveFirst()">First</button>
			<button @click="RemoveMiddle()">Middle</button>
			<button @click="RemoveLast()">Last</button>
			<button @click="ClearAll()">Clear All</button>
		</div>

		<div class="button-group">
			<h3>Replace All (simulates fetch)</h3>
			<button @click="LoadFruits()">Load Fruits</button>
			<button @click="LoadNumbers()">Load Numbers</button>
		</div>

		<div class="list-container">
			<h3>Current Items</h3>
			{#if ItemCount > 0}
				<ul>
					{#each Items as item, i}
						<li><span class="index">{i}</span> {item}</li>
					{/each}
				</ul>
			{:else}
				<p class="empty">No items in list</p>
			{/if}
		</div>
	</section>
</div>`
}

func (l *Lists) Style() string {
	return `
.demo { max-width: 600px; }
.demo h1 { color: #1a1a2e; margin-bottom: 20px; }
.demo section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 8px; background: #fff; }
.demo h2 { margin-top: 0; color: #666; font-size: 1.1em; }
.demo h3 { margin: 10px 0 5px; color: #888; font-size: 0.9em; }
.demo button { padding: 8px 16px; margin: 4px; cursor: pointer; border: 1px solid #ccc; border-radius: 4px; background: #f5f5f5; }
.demo button:hover { background: #e5e5e5; }
.demo input[type="text"] { padding: 8px; width: 200px; border: 1px solid #ccc; border-radius: 4px; }
.input-row { margin: 10px 0; }
.button-group { margin: 15px 0; padding: 10px; background: #f9f9f9; border-radius: 4px; }
.list-container { margin-top: 20px; padding: 15px; background: #f0f0f0; border-radius: 4px; }
.list-container ul { list-style: none; padding: 0; margin: 0; }
.list-container li { padding: 8px 12px; margin: 4px 0; background: white; border-radius: 4px; border-left: 3px solid #4CAF50; }
.index { display: inline-block; width: 24px; height: 24px; line-height: 24px; text-align: center; background: #4CAF50; color: white; border-radius: 50%; font-size: 12px; margin-right: 10px; }
.empty { color: #999; font-style: italic; text-align: center; padding: 20px; }
`
}
