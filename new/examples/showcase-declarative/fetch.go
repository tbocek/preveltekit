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

func (f *Fetch) OnMount() {
	if p.IsBuildTime {
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
	return p.Div(p.Class("demo"),
		p.H1("Fetch"),

		p.Section(
			p.H2("Typed Fetch"),
			p.P("Fetch data with automatic JSON decoding into Go structs:"),

			p.Div(p.Class("buttons"),
				p.Button("Fetch Todo", p.OnClick(f.FetchTodo)),
				p.Button("Fetch User", p.OnClick(f.FetchUser)),
				p.Button("Fetch Post", p.OnClick(f.FetchPost)),
				p.Button("Create Post (POST)", p.OnClick(f.CreatePost)),
			),

			p.If(f.RawData.Ne(""),
				p.Pre(p.Bind(f.RawData)),
			).Else(
				p.Pre("Click a button to fetch data"),
			),
			p.P(p.Class("status"), "Status: ", p.Bind(f.Status)),
		),

		p.Section(
			p.H2("Usage"),
			p.Pre(p.Class("code"), `type User struct {
    ID   int    `+"`js:\"id\"`"+`
    Name string `+"`js:\"name\"`"+`
}

go func() {
    user, err := preveltekit.Get[User](url)
    if err != nil { ... }
    // use user
}()`),
		),
	)
}

func (f *Fetch) Style() string {
	return `
.demo pre{min-height:60px}
.demo pre.code{background:#1a1a2e;color:#e0e0e0}
.demo .status{color:#666;font-size:.9em;margin-top:10px}
`
}

func (f *Fetch) HandleEvent(method string, args string) {
	switch method {
	case "FetchTodo":
		f.FetchTodo()
	case "FetchUser":
		f.FetchUser()
	case "FetchPost":
		f.FetchPost()
	case "CreatePost":
		f.CreatePost()
	}
}
