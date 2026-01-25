package main

import "reactive"

type App struct {
	CurrentPage *reactive.Store[string]
}

func (a *App) OnMount() {
	a.CurrentPage.Set("basics")
}

func (a *App) ShowBasics() {
	a.CurrentPage.Set("basics")
}

func (a *App) ShowLists() {
	a.CurrentPage.Set("lists")
}

func (a *App) ShowFetch() {
	a.CurrentPage.Set("fetch")
}

func (a *App) ShowStorage() {
	a.CurrentPage.Set("storage")
}

func (a *App) ShowDebounce() {
	a.CurrentPage.Set("debounce")
}

func (a *App) ShowBitcoin() {
	a.CurrentPage.Set("bitcoin")
}

func (a *App) ShowComponents() {
	a.CurrentPage.Set("components")
}

func (a *App) ShowRouting() {
	a.CurrentPage.Set("routing")
}

func (a *App) ShowLinks() {
	a.CurrentPage.Set("links")
}

func (a *App) Template() string {
	return `<div class="showcase">
	<nav class="sidebar">
		<h2>Reactive</h2>
		<ul>
			<li><button @click="ShowBasics()">Basics</button></li>
			<li><button @click="ShowComponents()">Components</button></li>
			<li><button @click="ShowLists()">Lists</button></li>
			<li><button @click="ShowRouting()">Routing</button></li>
			<li><button @click="ShowLinks()">Links</button></li>
			<li><button @click="ShowFetch()">Fetch</button></li>
			<li><button @click="ShowStorage()">Storage</button></li>
			<li><button @click="ShowDebounce()">Debounce</button></li>
			<li><button @click="ShowBitcoin()">Bitcoin</button></li>
		</ul>
	</nav>
	<main class="content">
		{#if CurrentPage == "basics"}
			<Basics />
		{:else if CurrentPage == "components"}
			<Components />
		{:else if CurrentPage == "lists"}
			<Lists />
		{:else if CurrentPage == "routing"}
			<Routing />
		{:else if CurrentPage == "links"}
			<Links />
		{:else if CurrentPage == "fetch"}
			<Fetch />
		{:else if CurrentPage == "storage"}
			<Storage />
		{:else if CurrentPage == "debounce"}
			<Debounce />
		{:else if CurrentPage == "bitcoin"}
			<Bitcoin />
		{/if}
	</main>
</div>`
}

func (a *App) Style() string {
	return `
* { box-sizing: border-box; }
body { margin: 0; font-family: system-ui, -apple-system, sans-serif; }
.showcase { display: flex; min-height: 100vh; }
.sidebar { width: 200px; background: #1a1a2e; color: #fff; padding: 20px; flex-shrink: 0; }
.sidebar h2 { margin: 0 0 20px; font-size: 1.5em; color: #fff; }
.sidebar ul { list-style: none; padding: 0; margin: 0; }
.sidebar li { margin: 5px 0; }
.sidebar button { width: 100%; padding: 10px 15px; background: transparent; border: none; color: #ccc; text-align: left; cursor: pointer; border-radius: 4px; font-size: 14px; }
.sidebar button:hover { background: #2a2a4e; color: #fff; }
.content { flex: 1; padding: 20px; background: #f5f5f5; overflow-y: auto; }
`
}
