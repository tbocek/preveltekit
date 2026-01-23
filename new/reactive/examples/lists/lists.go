package main

import "reactive"

type Lists struct {
	// List demo
	Items   *reactive.List[string]
	NewItem *reactive.Store[string]

	// Map demo (settings)
	Settings *reactive.Map[string, bool]

	// For showing list length
	ItemCount *reactive.Store[int]
}

func (l *Lists) OnMount() {
	l.Items.Set([]string{"Apple", "Banana", "Cherry"})
	l.NewItem.Set("")
	l.ItemCount.Set(3)

	// Initialize settings
	l.Settings.Set("notifications", true)
	l.Settings.Set("darkMode", false)
	l.Settings.Set("autoSave", true)
}

// List operations - various insertion points
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
	// Prepend by getting all, prepending, and setting
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

// Remove operations
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

// Bulk operations (simulates JSON fetch)
func (l *Lists) LoadFruits() {
	l.Items.Set([]string{"Mango", "Pineapple", "Papaya", "Guava"})
	l.updateCount()
}

func (l *Lists) LoadNumbers() {
	l.Items.Set([]string{"One", "Two", "Three", "Four", "Five"})
	l.updateCount()
}

func (l *Lists) LoadEmpty() {
	l.Items.Set([]string{})
	l.updateCount()
}

func (l *Lists) updateCount() {
	l.ItemCount.Set(l.Items.Len())
}

// Map/Settings operations
func (l *Lists) ToggleNotifications() {
	val, _ := l.Settings.Get("notifications")
	l.Settings.Set("notifications", !val)
}

func (l *Lists) ToggleDarkMode() {
	val, _ := l.Settings.Get("darkMode")
	l.Settings.Set("darkMode", !val)
}

func (l *Lists) ToggleAutoSave() {
	val, _ := l.Settings.Get("autoSave")
	l.Settings.Set("autoSave", !val)
}

func (l *Lists) Template() string {
	return `<div class="app">
	<h1>Lists & Maps Demo</h1>

	<section>
		<h2>List Operations</h2>
		<p>Items: {ItemCount}</p>

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
			<button @click="RemoveFirst()">Remove First</button>
			<button @click="RemoveMiddle()">Remove Middle</button>
			<button @click="RemoveLast()">Remove Last</button>
			<button @click="ClearAll()">Clear All</button>
		</div>

		<div class="button-group">
			<h3>Replace All (simulates JSON fetch)</h3>
			<button @click="LoadFruits()">Load Fruits</button>
			<button @click="LoadNumbers()">Load Numbers</button>
			<button @click="LoadEmpty()">Load Empty</button>
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

	<section>
		<h2>Map (Settings)</h2>
		<p>Toggle settings to see Map in action:</p>

		<div class="settings">
			<button @click="ToggleNotifications()">Toggle Notifications</button>
			<button @click="ToggleDarkMode()">Toggle Dark Mode</button>
			<button @click="ToggleAutoSave()">Toggle Auto-Save</button>
		</div>

		<p class="note">Note: Map changes trigger OnChange but we don't have {#each} for maps yet.</p>
	</section>
</div>`
}

func (l *Lists) Style() string {
	return `
.app { font-family: system-ui, sans-serif; max-width: 700px; margin: 0 auto; padding: 20px; }
section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 8px; }
h1 { color: #333; }
h2 { margin-top: 0; color: #666; font-size: 1.2em; }
h3 { margin: 10px 0 5px; color: #888; font-size: 0.9em; }
button { padding: 8px 16px; margin: 4px; cursor: pointer; border: 1px solid #ccc; border-radius: 4px; background: #f5f5f5; }
button:hover { background: #e5e5e5; }
input[type="text"] { padding: 8px; width: 200px; border: 1px solid #ccc; border-radius: 4px; }
.input-row { margin: 10px 0; }
.button-group { margin: 15px 0; padding: 10px; background: #f9f9f9; border-radius: 4px; }
.list-container { margin-top: 20px; padding: 15px; background: #f0f0f0; border-radius: 4px; }
ul { list-style: none; padding: 0; margin: 0; }
li { padding: 8px 12px; margin: 4px 0; background: white; border-radius: 4px; border-left: 3px solid #4CAF50; }
.index { display: inline-block; width: 24px; height: 24px; line-height: 24px; text-align: center; background: #4CAF50; color: white; border-radius: 50%; font-size: 12px; margin-right: 10px; }
.empty { color: #999; font-style: italic; text-align: center; padding: 20px; }
.settings { display: flex; gap: 10px; flex-wrap: wrap; }
.note { color: #666; font-size: 0.9em; font-style: italic; }
`
}
