//go:build js && wasm

package main

import (
	"reactive"
	"syscall/js"
)

var document = reactive.Document


func main() {
	component := &Basics{
		Count: reactive.New[int](0),
		Name: reactive.New[string](""),
		Message: reactive.New[string](""),
		DarkMode: reactive.New[bool](false),
		Agreed: reactive.New[bool](false),
		Score: reactive.New[int](0),
	}

	reactive.BindInt("expr_Count_0", component.Count)
	reactive.BindInt("expr_Score_1", component.Score)
	reactive.Bind("expr_Name_2", component.Name)
	reactive.Bind("expr_Message_3", component.Message)
	document.Call("getElementById", "evt_click_0").Call("addEventListener", "click",
		js.FuncOf(func(this js.Value, args []js.Value) any {
			component.Decrement()
			return nil
		}))

	document.Call("getElementById", "evt_click_1").Call("addEventListener", "click",
		js.FuncOf(func(this js.Value, args []js.Value) any {
			component.Increment()
			return nil
		}))

	document.Call("getElementById", "evt_click_2").Call("addEventListener", "click",
		js.FuncOf(func(this js.Value, args []js.Value) any {
			component.Add(5)
			return nil
		}))

	document.Call("getElementById", "evt_click_3").Call("addEventListener", "click",
		js.FuncOf(func(this js.Value, args []js.Value) any {
			component.Add(component.Count.Get())
			return nil
		}))

	document.Call("getElementById", "evt_click_4").Call("addEventListener", "click",
		js.FuncOf(func(this js.Value, args []js.Value) any {
			component.Reset()
			return nil
		}))

	document.Call("getElementById", "evt_click_5").Call("addEventListener", "click",
		js.FuncOf(func(this js.Value, args []js.Value) any {
			component.SetScore(95)
			return nil
		}))

	document.Call("getElementById", "evt_click_6").Call("addEventListener", "click",
		js.FuncOf(func(this js.Value, args []js.Value) any {
			component.SetScore(85)
			return nil
		}))

	document.Call("getElementById", "evt_click_7").Call("addEventListener", "click",
		js.FuncOf(func(this js.Value, args []js.Value) any {
			component.SetScore(75)
			return nil
		}))

	document.Call("getElementById", "evt_click_8").Call("addEventListener", "click",
		js.FuncOf(func(this js.Value, args []js.Value) any {
			component.SetScore(65)
			return nil
		}))

	document.Call("getElementById", "evt_click_9").Call("addEventListener", "click",
		js.FuncOf(func(this js.Value, args []js.Value) any {
			component.SetScore(50)
			return nil
		}))

	document.Call("getElementById", "evt_submit_10").Call("addEventListener", "submit",
		js.FuncOf(func(this js.Value, args []js.Value) any {
			args[0].Call("preventDefault")
			component.Submit()
			return nil
		}))

	if0_anchor := document.Call("getElementById", "if0_anchor")
	if0_current := js.Null()
	updateif0 := func() {
		var html string
		if component.Score.Get() >= 90 {
			html = `
			<p class="grade a">Grade: A - Excellent!</p>
		`
		} else if component.Score.Get() >= 80 {
			html = `
			<p class="grade b">Grade: B - Good</p>
		`
		} else if component.Score.Get() >= 70 {
			html = `
			<p class="grade c">Grade: C - Average</p>
		`
		} else if component.Score.Get() >= 60 {
			html = `
			<p class="grade d">Grade: D - Below Average</p>
		`
		} else {
			html = `
			<p class="grade f">Grade: F - Failing</p>
		`
		}
		newEl := document.Call("createElement", "span")
		newEl.Set("innerHTML", html)
		if !if0_current.IsNull() { if0_current.Call("remove") }
		if !if0_anchor.IsNull() { if0_anchor.Get("parentNode").Call("insertBefore", newEl, if0_anchor) }
		if0_current = newEl
	}
	component.Score.OnChange(func(_ int) { updateif0() })
	updateif0()

	bind0 := document.Call("getElementById", "bind0")
	bind0.Call("addEventListener", "input", js.FuncOf(func(this js.Value, args []js.Value) any {
		val := this.Get("value").String()
		component.Name.Set(val)
		return nil
	}))
	component.Name.OnChange(func(v string) { bind0.Set("value", v) })

	bind1 := document.Call("getElementById", "bind1")
	bind1.Call("addEventListener", "change", js.FuncOf(func(this js.Value, args []js.Value) any {
		component.DarkMode.Set(this.Get("checked").Bool())
		return nil
	}))
	component.DarkMode.OnChange(func(v bool) { bind1.Set("checked", v) })
	bind1.Set("checked", component.DarkMode.Get())

	bind2 := document.Call("getElementById", "bind2")
	bind2.Call("addEventListener", "input", js.FuncOf(func(this js.Value, args []js.Value) any {
		val := this.Get("value").String()
		component.Name.Set(val)
		return nil
	}))
	component.Name.OnChange(func(v string) { bind2.Set("value", v) })

	bind3 := document.Call("getElementById", "bind3")
	bind3.Call("addEventListener", "change", js.FuncOf(func(this js.Value, args []js.Value) any {
		component.Agreed.Set(this.Get("checked").Bool())
		return nil
	}))
	component.Agreed.OnChange(func(v bool) { bind3.Set("checked", v) })
	bind3.Set("checked", component.Agreed.Get())

	class0 := document.Call("getElementById", "class0")
	component.DarkMode.OnChange(func(v bool) {
		if v { class0.Get("classList").Call("add", "dark") } else { class0.Get("classList").Call("remove", "dark") }
	})

	component.OnMount()

	select {}
}
