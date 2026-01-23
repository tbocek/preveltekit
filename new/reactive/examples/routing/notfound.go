package main

// NotFound is the 404 page
type NotFound struct{}

func (n *NotFound) Template() string {
	return `<div class="page notfound-page">
	<div class="notfound-content">
		<h1>404</h1>
		<p>Page not found</p>
		<a href="/" class="btn">Go Home</a>
	</div>
</div>`
}

func (n *NotFound) Style() string {
	return `
.notfound-page { display: flex; align-items: center; justify-content: center; min-height: 60vh; }
.notfound-content { text-align: center; }
.notfound-content h1 { font-size: 6rem; color: #1a1a2e; margin: 0; line-height: 1; }
.notfound-content p { font-size: 1.5rem; color: #666; margin: 1rem 0 2rem; }
.notfound-content .btn { display: inline-block; padding: 0.75rem 2rem; background: #1a1a2e; color: #fff; text-decoration: none; border-radius: 6px; transition: background 0.2s; }
.notfound-content .btn:hover { background: #2d2d44; }
`
}
