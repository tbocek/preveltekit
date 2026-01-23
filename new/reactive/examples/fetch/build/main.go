//go:build js && wasm

package main

import (
	"syscall/js"
	"strconv"
)

var document = js.Global().Get("document")

func main() {
	component := &FetchDemo{
		Status: New[string](""),
		RawData: New[string](""),
		Users: NewList[string](),
		UserCount: New[int](0),
	}

	// Bind Status to #expr_Status_0
	expr_Status_0 := document.Call("getElementById", "expr_Status_0")
	component.Status.OnChange(func(v string) {
		if !expr_Status_0.IsUndefined() && !expr_Status_0.IsNull() {
			expr_Status_0.Set("textContent", v)
		}
	})
	if !expr_Status_0.IsUndefined() && !expr_Status_0.IsNull() {
		expr_Status_0.Set("textContent", component.Status.Get())
	}

	// Bind RawData to #expr_RawData_1
	expr_RawData_1 := document.Call("getElementById", "expr_RawData_1")
	component.RawData.OnChange(func(v string) {
		if !expr_RawData_1.IsUndefined() && !expr_RawData_1.IsNull() {
			expr_RawData_1.Set("textContent", v)
		}
	})
	if !expr_RawData_1.IsUndefined() && !expr_RawData_1.IsNull() {
		expr_RawData_1.Set("textContent", component.RawData.Get())
	}

	// Bind UserCount to #expr_UserCount_2
	expr_UserCount_2 := document.Call("getElementById", "expr_UserCount_2")
	component.UserCount.OnChange(func(v int) {
		if !expr_UserCount_2.IsUndefined() && !expr_UserCount_2.IsNull() {
			expr_UserCount_2.Set("textContent", strconv.Itoa(v))
		}
	})
	if !expr_UserCount_2.IsUndefined() && !expr_UserCount_2.IsNull() {
		expr_UserCount_2.Set("textContent", strconv.Itoa(component.UserCount.Get()))
	}

	// Bind @click to FetchTodo
	document.Call("getElementById", "evt_click_0").Call("addEventListener", "click",
		js.FuncOf(func(this js.Value, args []js.Value) any {
			component.FetchTodo()
			return nil
		}))

	// Bind @click to FetchUsers
	document.Call("getElementById", "evt_click_1").Call("addEventListener", "click",
		js.FuncOf(func(this js.Value, args []js.Value) any {
			component.FetchUsers()
			return nil
		}))

	// Bind @click to FetchFewUsers
	document.Call("getElementById", "evt_click_2").Call("addEventListener", "click",
		js.FuncOf(func(this js.Value, args []js.Value) any {
			component.FetchFewUsers()
			return nil
		}))

	// Bind @click to AddLocalUser
	document.Call("getElementById", "evt_click_3").Call("addEventListener", "click",
		js.FuncOf(func(this js.Value, args []js.Value) any {
			component.AddLocalUser()
			return nil
		}))

	// Bind @click to ClearUsers
	document.Call("getElementById", "evt_click_4").Call("addEventListener", "click",
		js.FuncOf(func(this js.Value, args []js.Value) any {
			component.ClearUsers()
			return nil
		}))

	// Bind if block if0
	if0_anchor := document.Call("getElementById", "if0_anchor")
	if0_current := js.Null()
	updateif0 := func() {
		var html string
		if component.UserCount.Get() > 0 {
			html = `
			<ul>
				<span id="each0_anchor"></span>
			</ul>
		`
		} else {
			html = `
			<p class="empty">No users loaded</p>
		`
		}
		newEl := document.Call("createElement", "span")
		newEl.Set("innerHTML", html)
		if !if0_current.IsNull() {
			if0_current.Call("remove")
		}
		if0_anchor.Get("parentNode").Call("insertBefore", newEl, if0_anchor)
		if0_current = newEl
	}
	component.UserCount.OnChange(func(_ int) { updateif0() })
	updateif0() // initial render

	// Bind each block for Users (diff-based)
	each0_anchor := document.Call("getElementById", "each0_anchor")
	each0_tmpl := `
					<li><span class="index"><span class="__index__"></span></span> <span class="__item__"></span></li>
				`

	// Helper to create a list item element
	each0_create := func(item string, index int) js.Value {
		wrapper := document.Call("createElement", "span")
		wrapper.Set("id", "each0_" + strconv.Itoa(index))
		wrapper.Set("innerHTML", each0_tmpl)
		if itemEl := wrapper.Call("querySelector", ".__item__"); !itemEl.IsNull() {
			itemEl.Set("textContent", item)
			itemEl.Get("classList").Call("remove", "__item__")
		}
		if idxEl := wrapper.Call("querySelector", ".__index__"); !idxEl.IsNull() {
			idxEl.Set("textContent", strconv.Itoa(index))
			idxEl.Get("classList").Call("remove", "__index__")
		}
		return wrapper
	}

	// OnEdit: handle insert/remove from diff
	component.Users.OnEdit(func(edit Edit[string]) {
		switch edit.Op {
		case EditInsert:
			// Shift existing elements' IDs first
			items := component.Users.Get()
			for i := len(items) - 1; i > edit.Index; i-- {
				el := document.Call("getElementById", "each0_" + strconv.Itoa(i-1))
				if !el.IsNull() {
					el.Set("id", "each0_" + strconv.Itoa(i))
				}
			}
			// Create and insert new element
			el := each0_create(edit.Value, edit.Index)
			if edit.Index == 0 {
				// Insert at beginning
				first := document.Call("getElementById", "each0_1")
				if !first.IsNull() {
					each0_anchor.Get("parentNode").Call("insertBefore", el, first)
				} else {
					each0_anchor.Get("parentNode").Call("insertBefore", el, each0_anchor)
				}
			} else {
				// Insert after previous element
				prev := document.Call("getElementById", "each0_" + strconv.Itoa(edit.Index-1))
				if !prev.IsNull() {
					prev.Get("parentNode").Call("insertBefore", el, prev.Get("nextSibling"))
				} else {
					each0_anchor.Get("parentNode").Call("insertBefore", el, each0_anchor)
				}
			}
		case EditRemove:
			el := document.Call("getElementById", "each0_" + strconv.Itoa(edit.Index))
			if !el.IsNull() {
				el.Call("remove")
			}
			// Re-index following elements
			for i := edit.Index; ; i++ {
				nextEl := document.Call("getElementById", "each0_" + strconv.Itoa(i+1))
				if nextEl.IsNull() {
					break
				}
				nextEl.Set("id", "each0_" + strconv.Itoa(i))
			}
		}
	})

	// Initial render
	component.Users.OnRender(func(items []string) {
		for i, item := range items {
			el := each0_create(item, i)
			each0_anchor.Get("parentNode").Call("insertBefore", el, each0_anchor)
		}
	})
	component.Users.Render()

	component.OnMount()

	select {}
}
