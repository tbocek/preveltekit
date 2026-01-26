package main

import "preveltekit"

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

type Post struct {
	ID     int    `js:"id"`
	UserID int    `js:"userId"`
	Title  string `js:"title"`
	Body   string `js:"body"`
}

type Fetch struct {
	Status  *preveltekit.Store[string]
	RawData *preveltekit.Store[string]
}

func (f *Fetch) OnMount() {
	if preveltekit.IsBuildTime {
		f.Status.Set("ready")
		f.RawData.Set("")
		return
	}

	f.Status.Set("ready")
	f.RawData.Set("")
}

func (f *Fetch) FetchTodo() {
	f.Status.Set("loading...")
	f.RawData.Set("")

	go func() {
		todo, err := preveltekit.Get[Todo]("https://jsonplaceholder.typicode.com/todos/1")
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

func (f *Fetch) FetchUser() {
	f.Status.Set("loading...")
	f.RawData.Set("")

	go func() {
		user, err := preveltekit.Get[User]("https://jsonplaceholder.typicode.com/users/1")
		if err != nil {
			f.Status.Set("error: " + err.Error())
			return
		}
		f.RawData.Set("ID: " + itoa(user.ID) + "\nName: " + user.Name)
		f.Status.Set("done")
	}()
}

func (f *Fetch) FetchPost() {
	f.Status.Set("loading...")
	f.RawData.Set("")

	go func() {
		post, err := preveltekit.Get[Post]("https://jsonplaceholder.typicode.com/posts/1")
		if err != nil {
			f.Status.Set("error: " + err.Error())
			return
		}
		f.RawData.Set("ID: " + itoa(post.ID) + "\nUser: " + itoa(post.UserID) + "\nTitle: " + post.Title)
		f.Status.Set("done")
	}()
}

func (f *Fetch) CreatePost() {
	f.Status.Set("creating...")
	f.RawData.Set("")

	go func() {
		// Send a new post via POST request with JSON body
		newPost := Post{
			UserID: 1,
			Title:  "Hello from Go WASM",
			Body:   "This post was created using preveltekit.Post[T]",
		}

		created, err := preveltekit.Post[Post]("https://jsonplaceholder.typicode.com/posts", newPost)
		if err != nil {
			f.Status.Set("error: " + err.Error())
			return
		}
		f.RawData.Set("Created Post!\nID: " + itoa(created.ID) + "\nTitle: " + created.Title + "\nBody: " + created.Body)
		f.Status.Set("done")
	}()
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
		<p>Fetch data with automatic JSON decoding into Go structs:</p>

		<div class="buttons">
			<button @click="FetchTodo()">Fetch Todo</button>
			<button @click="FetchUser()">Fetch User</button>
			<button @click="FetchPost()">Fetch Post</button>
			<button @click="CreatePost()">Create Post (POST)</button>
		</div>

		{#if RawData != ""}
			<pre>{RawData}</pre>
		{:else}
			<pre>Click a button to fetch data</pre>
		{/if}
		<p class="status">Status: {Status}</p>
	</section>

	<section>
		<h2>Usage</h2>
		<pre class="code">type User struct {
    ID   int    ` + "`js:\"id\"`" + `
    Name string ` + "`js:\"name\"`" + `
}

go func() {
    user, err := preveltekit.Get[User](url)
    if err != nil { ... }
    // use user
}()</pre>
	</section>
</div>`
}

func (f *Fetch) Style() string {
	return `
.demo pre{min-height:60px}
.demo pre.code{background:#1a1a2e;color:#e0e0e0}
.demo .status{color:#666;font-size:.9em;margin-top:10px}
`
}
