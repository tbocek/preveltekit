package main

import p "preveltekit"

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
	Status  *p.Store[string]
	RawData *p.Store[string]
}

func (f *Fetch) New() p.Component {
	return &Fetch{
		Status:  p.New("ready"),
		RawData: p.New(""),
	}
}

func (f *Fetch) FetchTodo() {
	f.Status.Set("loading...")
	f.RawData.Set("")

	go func() {
		todo, err := p.Get[Todo]("https://jsonplaceholder.typicode.com/todos/1")
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
		user, err := p.Get[User]("https://jsonplaceholder.typicode.com/users/1")
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
		post, err := p.Get[Post]("https://jsonplaceholder.typicode.com/posts/1")
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
		newPost := Post{
			UserID: 1,
			Title:  "Hello from Go WASM",
			Body:   "This post was created using preveltekit.Post[T]",
		}

		created, err := p.Post[Post]("https://jsonplaceholder.typicode.com/posts", newPost)
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

func (f *Fetch) Render() p.Node {
	return p.Html(`<div class="demo">
		<h1>Fetch</h1>

		<section>
			<h2>Typed Fetch</h2>
			<p>Fetch data with automatic JSON decoding into Go structs:</p>

			<div class="buttons">
				`, p.Html(`<button>Fetch Todo</button>`).On("click", f.FetchTodo), `
				`, p.Html(`<button>Fetch User</button>`).On("click", f.FetchUser), `
				`, p.Html(`<button>Fetch Post</button>`).On("click", f.FetchPost), `
				`, p.Html(`<button>Create Post (POST)</button>`).On("click", f.CreatePost), `
			</div>`,

		p.If(f.RawData.Ne(""),
			p.Html(`<pre>`, f.RawData, `</pre>`),
		).Else(
			p.Html(`<pre>Click a button to fetch data</pre>`),
		),

		p.Html(`<p class="status">Status: `, f.Status, `</p>
		</section>

		<section>
			<h2>Usage</h2>
			<pre class="code">type User struct {
    ID   int    `+"`"+`js:"id"`+"`"+`
    Name string `+"`"+`js:"name"`+"`"+`
}

go func() {
    user, err := preveltekit.Get[User](url)
    if err != nil { ... }
    // use user
}()</pre>
		</section>
	</div>`),
	)
}

func (f *Fetch) Style() string {
	return `
.demo pre{min-height:60px}
.demo pre.code{background:#1a1a2e;color:#e0e0e0}
.demo .status{color:#666;font-size:.9em;margin-top:10px}
`
}
