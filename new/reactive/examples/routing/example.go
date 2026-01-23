package main

import "reactive"

// Example is the examples page showing code samples
type Example struct {
	ActiveTab *reactive.Store[string]
}

func (e *Example) OnMount() {
	e.ActiveTab.Set("counter")
}

func (e *Example) ShowCounter() {
	e.ActiveTab.Set("counter")
}

func (e *Example) ShowList() {
	e.ActiveTab.Set("list")
}

func (e *Example) ShowFetch() {
	e.ActiveTab.Set("fetch")
}

func (e *Example) ShowRouting() {
	e.ActiveTab.Set("routing")
}

func (e *Example) ShowBitcoin() {
	e.ActiveTab.Set("bitcoin")
}

func (e *Example) Template() string {
	return `<div class="page example-page">
	<h1>Examples</h1>
	<p class="intro">Explore working examples of Reactive patterns and features.</p>

	<div class="tabs">
		<button class="tab" @click="ShowCounter()">Counter</button>
		<button class="tab" @click="ShowList()">Lists</button>
		<button class="tab" @click="ShowFetch()">Fetch</button>
		<button class="tab" @click="ShowRouting()">Routing</button>
		<button class="tab" @click="ShowBitcoin()">Bitcoin</button>
	</div>

	<div class="tab-content">
		{#if ActiveTab == "counter"}
		<div class="example">
			<h2>Counter Example</h2>
			<p>Basic reactive state with increment/decrement.</p>
			<pre><code>type Counter struct {
    Count *reactive.Store[int]
}

func (c *Counter) OnMount() {
    c.Count.Set(0)
}

func (c *Counter) Increment() {
    c.Count.Update(func(v int) int { return v + 1 })
}

func (c *Counter) Decrement() {
    c.Count.Update(func(v int) int { return v - 1 })
}

func (c *Counter) Template() string {
    return ` + "`" + `<div class="counter">
    <button @click="Decrement()">-</button>
    <span>{Count}</span>
    <button @click="Increment()">+</button>
</div>` + "`" + `
}</code></pre>
		</div>
		{:else if ActiveTab == "list"}
		<div class="example">
			<h2>List Example</h2>
			<p>Dynamic lists with efficient diffing.</p>
			<pre><code>type TodoList struct {
    Items *reactive.List[string]
    Input *reactive.Store[string]
}

func (t *TodoList) OnMount() {
    t.Items = reactive.NewList("Buy groceries", "Walk the dog")
    t.Input.Set("")
}

func (t *TodoList) Add() {
    if v := t.Input.Get(); v != "" {
        t.Items.Append(v)
        t.Input.Set("")
    }
}

func (t *TodoList) Remove(index int) {
    t.Items.RemoveAt(index)
}

func (t *TodoList) Template() string {
    return ` + "`" + `<div class="todo">
    <input bind:value="Input" />
    <button @click="Add()">Add</button>
    <ul>
        {#each Items as item, i}
        <li>
            {item}
            <button @click="Remove(i)">x</button>
        </li>
        {/each}
    </ul>
</div>` + "`" + `
}</code></pre>
		</div>
		{:else if ActiveTab == "fetch"}
		<div class="example">
			<h2>Fetch Example</h2>
			<p>Async data fetching with loading states.</p>
			<pre><code>type UserList struct {
    Users   *reactive.List[User]
    Loading *reactive.Store[bool]
    Error   *reactive.Store[string]
}

func (u *UserList) OnMount() {
    u.Loading.Set(true)
    u.FetchUsers()
}

func (u *UserList) FetchUsers() {
    go func() {
        resp, err := http.Get("/api/users")
        if err != nil {
            u.Error.Set(err.Error())
            u.Loading.Set(false)
            return
        }
        defer resp.Body.Close()

        var users []User
        json.NewDecoder(resp.Body).Decode(&users)

        u.Users.Set(users)
        u.Loading.Set(false)
    }()
}

func (u *UserList) Template() string {
    return ` + "`" + `<div>
    {#if Loading}
        <p>Loading...</p>
    {:else if Error}
        <p class="error">{Error}</p>
    {:else}
        {#each Users as user}
        <div class="user">{user.Name}</div>
        {/each}
    {/if}
</div>` + "`" + `
}</code></pre>
		</div>
		{:else if ActiveTab == "routing"}
		<div class="example">
			<h2>Routing Example</h2>
			<p>Client-side SPA routing with parameters.</p>
			<pre><code>type App struct {
    Page   *reactive.Store[string]
    UserID *reactive.Store[string]
    router *reactive.Router
}

func (a *App) OnMount() {
    a.router = reactive.NewRouter()

    a.router.Handle("/", func(p map[string]string) {
        a.Page.Set("home")
    })

    a.router.Handle("/user/:id", func(p map[string]string) {
        a.Page.Set("user")
        a.UserID.Set(p["id"])
    })

    a.router.Handle("/about", func(p map[string]string) {
        a.Page.Set("about")
    })

    a.router.NotFound(func() {
        a.Page.Set("404")
    })

    a.router.Start()
}

func (a *App) Template() string {
    return ` + "`" + `<div>
    <nav>
        <a href="/">Home</a>
        <a href="/user/123">User 123</a>
        <a href="/about">About</a>
    </nav>
    {#if Page == "home"}
        <Home />
    {:else if Page == "user"}
        <UserProfile id="{UserID}" />
    {:else if Page == "about"}
        <About />
    {:else}
        <NotFound />
    {/if}
</div>` + "`" + `
}</code></pre>
		</div>
		{:else if ActiveTab == "bitcoin"}
		<div class="example">
			<h2>Bitcoin Price Tracker</h2>
			<p>Live Bitcoin price fetching with auto-refresh. Click the tab to see it in action!</p>
			<Bitcoin />
			<h3>Source Code</h3>
			<pre><code>type Bitcoin struct {
    Price      *reactive.Store[string]
    Loading    *reactive.Store[bool]
    Error      *reactive.Store[string]
    intervalID js.Value
}

func (b *Bitcoin) OnMount() {
    b.Loading.Set(true)
    b.FetchPrice()

    // Auto-refresh every 60 seconds
    callback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        b.FetchPrice()
        return nil
    })
    b.intervalID = js.Global().Call("setInterval", callback, 60000)
}

func (b *Bitcoin) FetchPrice() {
    go func() {
        resp, err := http.Get("https://min-api.cryptocompare.com/...")
        if err != nil {
            b.Error.Set(err.Error())
            return
        }
        var data BitcoinPrice
        json.NewDecoder(resp.Body).Decode(&data)
        b.Price.Set(fmt.Sprintf("$%.2f", data.RAW.PRICE))
        b.Loading.Set(false)
    }()
}

func (b *Bitcoin) Template() string {
    return ` + "`" + `<div class="bitcoin-card">
    {#if Loading}
        <p>Loading...</p>
    {:else if Error}
        <p class="error">{Error}</p>
        <button @click="Retry()">Retry</button>
    {:else}
        <p class="price">{Price}</p>
    {/if}
</div>` + "`" + `
}</code></pre>
		</div>
		{:else}
		<div class="example">
			<p>Select a tab to view examples.</p>
		</div>
		{/if}
	</div>
</div>`
}

func (e *Example) Style() string {
	return `
.example-page h1 { color: #1a1a2e; margin-bottom: 0.5rem; }
.intro { color: #666; margin-bottom: 2rem; }
.tabs { display: flex; gap: 0.5rem; margin-bottom: 1rem; border-bottom: 2px solid #e9ecef; padding-bottom: 0.5rem; }
.tab { background: none; border: none; padding: 0.75rem 1.5rem; font-size: 1rem; color: #666; cursor: pointer; border-radius: 4px 4px 0 0; transition: all 0.2s; }
.tab:hover { color: #1a1a2e; background: #f5f5f5; }
.tab.active { color: #1a1a2e; background: #e9ecef; }
.tab-content { background: #fff; border: 1px solid #e9ecef; border-radius: 8px; padding: 2rem; }
.example h2 { color: #1a1a2e; margin-bottom: 0.5rem; }
.example h3 { color: #1a1a2e; margin-top: 2rem; margin-bottom: 1rem; }
.example p { color: #666; margin-bottom: 1.5rem; }
.example pre { background: #1a1a2e; color: #e9ecef; padding: 1.5rem; border-radius: 8px; overflow-x: auto; }
.example code { font-family: 'Fira Code', monospace; font-size: 0.85rem; line-height: 1.5; }
`
}
