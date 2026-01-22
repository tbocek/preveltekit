package main

// Reactive state - var + setter pair
var count int
func setCount(v int) {}

const _template = `
<div>
    <h1>Count: {count}</h1>
    <button @click="increment()">+</button>
    <button @click="decrement()">-</button>
    <button @click="reset()">Reset</button>
</div>
`

const _style = `
button {
    padding: 10px 20px;
    margin: 5px;
    font-size: 16px;
}
`

// Natural Go syntax - gets transformed to use setters
func increment() {
	count++
}

func decrement() {
	count--
}

func reset() {
	count = 0
}