package main

import "reactive"

type Todo struct {
	ID        int    `js:"id"`
	UserID    int    `js:"userId"`
	Title     string `js:"title"`
	Completed bool   `js:"completed"`
}

type User struct {
	ID   int    `js:"id"`
	Name string `js:"name"`
}

type Fetch struct {
	Status    *reactive.Store[string]
	RawData   *reactive.Store[string]
	Users     *reactive.List[string]
	UserCount *reactive.Store[int]
}

func (f *Fetch) OnMount() {
	if reactive.IsBuildTime {
		f.Status.Set("loading...")
		f.RawData.Set("Loading data...")
		f.UserCount.Set(0)
		return
	}

	f.Status.Set("idle")
	f.RawData.Set("")
	f.UserCount.Set(0)
}

func (f *Fetch) FetchTodo() {
	f.Status.Set("loading...")
	f.RawData.Set("")

	go func() {
		todo, err := reactive.Get[Todo]("https://jsonplaceholder.typicode.com/todos/1")
		if err != nil {
			f.Status.Set("error: " + err.Error())
			return
		}
		completed := "no"
		if todo.Completed {
			completed = "yes"
		}
		f.RawData.Set("ID: " + itoa(todo.ID) + "\nTitle: " + todo.Title + "\nCompleted: " + completed)
		f.Status.Set("done")
	}()
}

func (f *Fetch) FetchUsers() {
	f.Status.Set("loading users...")

	go func() {
		users, err := reactive.Get[[]User]("https://jsonplaceholder.typicode.com/users")
		if err != nil {
			f.Status.Set("error: " + err.Error())
			return
		}

		names := make([]string, len(users))
		for i, u := range users {
			names[i] = u.Name
		}
		f.Users.Set(names)
		f.UserCount.Set(len(names))
		f.Status.Set("loaded " + itoa(len(names)) + " users")
	}()
}

func (f *Fetch) FetchFewUsers() {
	f.Status.Set("loading subset...")

	go func() {
		users, err := reactive.Get[[]User]("https://jsonplaceholder.typicode.com/users?_limit=3")
		if err != nil {
			f.Status.Set("error: " + err.Error())
			return
		}

		names := make([]string, len(users))
		for i, u := range users {
			names[i] = u.Name
		}
		f.Users.Set(names)
		f.UserCount.Set(len(names))
		f.Status.Set("loaded " + itoa(len(names)) + " users (subset)")
	}()
}

func (f *Fetch) ClearUsers() {
	f.Users.Set([]string{})
	f.UserCount.Set(0)
	f.Status.Set("cleared")
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(byte('0'+n%10)) + s
		n /= 10
	}
	return s
}

func (f *Fetch) Template() string {
	return `<div class="demo">
	<h1>Fetch</h1>

	<section>
		<h2>Typed Fetch</h2>
		<p>Fetch a todo item with typed response:</p>
		<button @click="FetchTodo()">Fetch Todo</button>
		<pre>{RawData}</pre>
		<p class="status">Status: {Status}</p>
	</section>

	<section>
		<h2>Fetch List (Diff Demo)</h2>
		<p>Fetches user names. Try loading all, then 3 to see diff in action.</p>
		<div class="buttons">
			<button @click="FetchUsers()">Fetch All Users</button>
			<button @click="FetchFewUsers()">Fetch 3 Users</button>
			<button @click="ClearUsers()">Clear</button>
		</div>

		<p>Users: <strong>{UserCount}</strong></p>

		{#if UserCount > 0}
			<ul>
				{#each Users as user, i}
					<li><span class="index">{i}</span> {user}</li>
				{/each}
			</ul>
		{:else}
			<p class="empty">No users loaded</p>
		{/if}
	</section>
</div>`
}

func (f *Fetch) Style() string {
	return `
.demo { max-width: 600px; }
.demo h1 { color: #1a1a2e; margin-bottom: 20px; }
.demo section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 8px; background: #fff; }
.demo h2 { margin-top: 0; color: #666; font-size: 1.1em; }
.demo button { padding: 8px 16px; margin: 4px; cursor: pointer; border: 1px solid #ccc; border-radius: 4px; background: #f5f5f5; }
.demo button:hover { background: #e5e5e5; }
.demo pre { background: #f5f5f5; padding: 15px; border-radius: 4px; overflow-x: auto; min-height: 50px; font-size: 12px; }
.demo .status { color: #666; font-size: 0.9em; }
.demo ul { list-style: none; padding: 0; margin: 10px 0; }
.demo li { padding: 8px 12px; margin: 4px 0; background: #f9f9f9; border-radius: 4px; border-left: 3px solid #007bff; }
.index { display: inline-block; width: 20px; height: 20px; line-height: 20px; text-align: center; background: #007bff; color: white; border-radius: 50%; font-size: 11px; margin-right: 10px; }
.empty { color: #999; font-style: italic; }
.buttons { display: flex; gap: 10px; flex-wrap: wrap; margin: 10px 0; }
`
}
