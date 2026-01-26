package main

import "preveltekit"

type Basics struct {
	Count    *preveltekit.Store[int]
	Name     *preveltekit.Store[string]
	Message  *preveltekit.Store[string]
	DarkMode *preveltekit.Store[bool]
	Agreed   *preveltekit.Store[bool]
	Score    *preveltekit.Store[int]
}

func (b *Basics) OnMount() {
	b.Count.Set(0)
	b.Name.Set("")
	b.Message.Set("Fill out the form above")
	b.DarkMode.Set(false)
	b.Agreed.Set(false)
	b.Score.Set(75)
}

func (b *Basics) Increment() {
	b.Count.Update(func(v int) int { return v + 1 })
}

func (b *Basics) Decrement() {
	b.Count.Update(func(v int) int { return v - 1 })
}

func (b *Basics) Add(n int) {
	b.Count.Update(func(v int) int { return v + n })
}

func (b *Basics) Reset() {
	b.Count.Set(0)
}

func (b *Basics) SetScore(n int) {
	b.Score.Set(n)
}

func (b *Basics) Submit() {
	name := b.Name.Get()
	if name == "" {
		b.Message.Set("Please enter your name")
	} else if !b.Agreed.Get() {
		b.Message.Set("Please agree to terms")
	} else {
		b.Message.Set("Welcome, " + name + "!")
	}
}

func (b *Basics) Template() string {
	return `<div class="demo">
	<h1>Basics</h1>

	<section>
		<h2>Counter</h2>
		<p>Count: <strong>{Count}</strong></p>
		<button @click="Decrement()">-1</button>
		<button @click="Increment()">+1</button>
		<button @click="Add(5)">+5</button>
		<button @click="Add(Count)">Double</button>
		<button @click="Reset()">Reset</button>
	</section>

	<section>
		<h2>Conditionals</h2>
		<p>Score: {Score}</p>
		{#if Score >= 90}
			<p class="grade a">Grade: A - Excellent!</p>
		{:else if Score >= 80}
			<p class="grade b">Grade: B - Good</p>
		{:else if Score >= 70}
			<p class="grade c">Grade: C - Average</p>
		{:else if Score >= 60}
			<p class="grade d">Grade: D - Below Average</p>
		{:else}
			<p class="grade f">Grade: F - Failing</p>
		{/if}
		<div class="buttons">
			<button @click="SetScore(95)">A</button>
			<button @click="SetScore(85)">B</button>
			<button @click="SetScore(75)">C</button>
			<button @click="SetScore(65)">D</button>
			<button @click="SetScore(50)">F</button>
		</div>
	</section>

	<section>
		<h2>Two-Way Binding</h2>
		<label>Your name: <input type="text" bind:value="Name" placeholder="Enter name"></label>
		<p>Hello, {Name}!</p>
	</section>

	<section>
		<h2>Checkbox Binding</h2>
		<label><input type="checkbox" bind:checked="DarkMode"> Dark Mode</label>
		<div class:dark={DarkMode}>
			This box uses dark mode styling when checked.
		</div>
	</section>

	<section>
		<h2>Form</h2>
		<form @submit.preventDefault="Submit()">
			<label>Name: <input type="text" bind:value="Name" placeholder="Your name"></label>
			<label><input type="checkbox" bind:checked="Agreed"> I agree to the terms</label>
			<button type="submit">Submit</button>
		</form>
		<p class="message">{Message}</p>
	</section>
</div>`
}

func (b *Basics) Style() string {
	return `
.demo label{display:block;margin:8px 0}
.grade{padding:10px;border-radius:4px;font-weight:700}
.grade.a{background:#d4edda;color:#155724}
.grade.b{background:#cce5ff;color:#004085}
.grade.c{background:#fff3cd;color:#856404}
.grade.d{background:#ffe5d0;color:#854027}
.grade.f{background:#f8d7da;color:#721c24}
.dark{background:#333;color:#fff;padding:10px;border-radius:4px;margin-top:10px}
.message{padding:10px;background:#e7f3ff;border-radius:4px}
`
}
