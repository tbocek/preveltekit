//go:build js && wasm

package main

import (
	"reactive"
	"syscall/js"
	"strings"
)

var document = reactive.Document

const buttonCSS = `
		.btn { padding: 0.5em 1em; cursor: pointer; border: none; border-radius: 4px; }
		.btn.primary { background: #007bff; color: white; }
		.btn.secondary { background: #6c757d; color: white; }
	`
const cardCSS = `
.card { border: 1px solid #ddd; border-radius: 8px; margin: 10px 0; overflow: hidden; }
.card-header { background: #f5f5f5; padding: 10px 15px; font-weight: bold; border-bottom: 1px solid #ddd; }
.card-body { padding: 15px; }
`
const counterCSS = `
.counter { display: inline-flex; align-items: center; gap: 8px; padding: 10px; background: #f0f0f0; border-radius: 8px; margin: 10px 0; }
.counter .value { font-size: 24px; font-weight: bold; min-width: 60px; text-align: center; }
.counter-btn { padding: 8px 12px; border: none; border-radius: 4px; cursor: pointer; background: #007bff; color: white; }
.counter-btn:hover { background: #0056b3; }
.counter-btn.reset { background: #6c757d; }
.counter-btn.reset:hover { background: #545b62; }
`
const comp2HTML = `<div id="comp2" class="card">
	<div class="card-header"><span id="expr_Title_0"></span></div>
	<div class="card-body"><p>This card can be shown/hidden.</p>
				<p>Parent count: <span id="expr_Count_1"></span></p></div>
</div>`

func main() {
	component := &App{
		Count: reactive.New[int](0),
		Message: reactive.New[string](""),
		Theme: reactive.New[string](""),
		ShowCard: reactive.New[bool](false),
	}

	reactive.BindInt("expr_Count_0", component.Count)
	reactive.Bind("expr_Message_1", component.Message)
	if0_anchor := document.Call("getElementById", "if0_anchor")
	if0_current := js.Null()
	updateif0 := func() {
		var html string
		if component.ShowCard.Get() {
			html = `
			<!--comp2-->
		`
		} else {
			html = `
			<p class="hidden-note">Card is hidden</p>
		`
		}
		html = strings.Replace(html, "<!--comp2-->", comp2HTML, 1)
		newEl := document.Call("createElement", "span")
		newEl.Set("innerHTML", html)
		if !if0_current.IsNull() { if0_current.Call("remove") }
		if0_anchor.Get("parentNode").Call("insertBefore", newEl, if0_anchor)
		if0_current = newEl
	}
	component.ShowCard.OnChange(func(_ bool) { updateif0() })
	updateif0()

	comp0 := &Button{
		Label: reactive.New[string](""),
		Variant: reactive.New[string](""),
	}

	comp0_el := reactive.GetEl("comp0")
	if !comp0_el.IsNull() && !comp0_el.IsUndefined() {
	reactive.InjectStyle("Button", buttonCSS)
	comp0.Variant.Set("primary")

	reactive.Bind("comp0_expr_Label_1", comp0.Label)
	reactive.BindInt("comp0_expr_Count_0", component.Count)
	reactive.BindAttr("[data-attrbind=\"comp0_attr0\"]", "class", `btn {Variant}`, "Variant", comp0.Variant)
	reactive.On(comp0_el, "click", func() { component.Increment() })
	}

	comp1 := &Card{
		Title: reactive.New[string](""),
	}

	comp1_el := reactive.GetEl("comp1")
	if !comp1_el.IsNull() && !comp1_el.IsUndefined() {
	reactive.InjectStyle("Card", cardCSS)
	comp1.Title.Set("Status Card")

	reactive.Bind("comp1_expr_Title_0", comp1.Title)
	reactive.BindInt("comp1_expr_Count_1", component.Count)
	reactive.Bind("comp1_expr_Theme_2", component.Theme)
	reactive.On(comp1_el, "click", func() { component.CardClicked() })
	}

	comp2 := &Card{
		Title: reactive.New[string](""),
	}

	comp2_el := reactive.GetEl("comp2")
	if !comp2_el.IsNull() && !comp2_el.IsUndefined() {
	reactive.InjectStyle("Card", cardCSS)
	comp2.Title.Set("Toggleable Card")

	reactive.Bind("comp2_expr_Title_0", comp2.Title)
	reactive.BindInt("comp2_expr_Count_1", component.Count)
	reactive.On(comp2_el, "click", func() { component.CardClicked() })
	}

	comp3 := &Button{
		Label: reactive.New[string](""),
		Variant: reactive.New[string](""),
	}

	comp3_el := reactive.GetEl("comp3")
	if !comp3_el.IsNull() && !comp3_el.IsUndefined() {
	reactive.InjectStyle("Button", buttonCSS)
	comp3.Label.Set("Primary")

	comp3.Variant.Set("primary")

	reactive.Bind("comp3_expr_Label_0", comp3.Label)
	reactive.BindAttr("[data-attrbind=\"comp3_attr0\"]", "class", `btn {Variant}`, "Variant", comp3.Variant)
	reactive.On(comp3_el, "click", func() { component.Increment() })
	}

	comp4 := &Button{
		Label: reactive.New[string](""),
		Variant: reactive.New[string](""),
	}

	comp4_el := reactive.GetEl("comp4")
	if !comp4_el.IsNull() && !comp4_el.IsUndefined() {
	reactive.InjectStyle("Button", buttonCSS)
	comp4.Label.Set("Secondary")

	comp4.Variant.Set("secondary")

	reactive.Bind("comp4_expr_Label_0", comp4.Label)
	reactive.BindAttr("[data-attrbind=\"comp4_attr0\"]", "class", `btn {Variant}`, "Variant", comp4.Variant)
	reactive.On(comp4_el, "click", func() { component.Decrement() })
	}

	comp5 := &Button{
		Label: reactive.New[string](""),
		Variant: reactive.New[string](""),
	}

	comp5_el := reactive.GetEl("comp5")
	if !comp5_el.IsNull() && !comp5_el.IsUndefined() {
	reactive.InjectStyle("Button", buttonCSS)
	comp5.Label.Set("Success")

	comp5.Variant.Set("success")

	reactive.Bind("comp5_expr_Label_0", comp5.Label)
	reactive.BindAttr("[data-attrbind=\"comp5_attr0\"]", "class", `btn {Variant}`, "Variant", comp5.Variant)
	reactive.On(comp5_el, "click", func() { component.Add(5) })
	}

	comp6 := &Button{
		Label: reactive.New[string](""),
		Variant: reactive.New[string](""),
	}

	comp6_el := reactive.GetEl("comp6")
	if !comp6_el.IsNull() && !comp6_el.IsUndefined() {
	reactive.InjectStyle("Button", buttonCSS)
	comp6.Label.Set("Danger")

	comp6.Variant.Set("danger")

	reactive.Bind("comp6_expr_Label_0", comp6.Label)
	reactive.BindAttr("[data-attrbind=\"comp6_attr0\"]", "class", `btn {Variant}`, "Variant", comp6.Variant)
	reactive.On(comp6_el, "click", func() { component.Reset() })
	}

	comp7 := &Button{
		Label: reactive.New[string](""),
		Variant: reactive.New[string](""),
	}

	comp7_el := reactive.GetEl("comp7")
	if !comp7_el.IsNull() && !comp7_el.IsUndefined() {
	reactive.InjectStyle("Button", buttonCSS)
	comp7.Label.Set(component.Message.Get())
	component.Message.OnChange(func(v string) { comp7.Label.Set(v) })

	comp7.Variant.Set("primary")

	reactive.Bind("comp7_expr_Label_0", comp7.Label)
	reactive.BindAttr("[data-attrbind=\"comp7_attr0\"]", "class", `btn {Variant}`, "Variant", comp7.Variant)
	reactive.On(comp7_el, "click", func() { component.Increment() })
	}

	comp8 := &Counter{
		Initial: reactive.New[int](0),
		Step: reactive.New[int](0),
		Value: reactive.New[int](0),
	}

	comp8_el := reactive.GetEl("comp8")
	if !comp8_el.IsNull() && !comp8_el.IsUndefined() {
	reactive.InjectStyle("Counter", counterCSS)
	comp8.Step.Set(5)

	comp8.Initial.Set(10)

	reactive.BindInt("comp8_expr_Value_0", comp8.Value)
	reactive.BindInt("comp8_expr_Step_1", comp8.Step)
	reactive.BindInt("comp8_expr_Step_2", comp8.Step)
	reactive.On(reactive.GetEl("comp8_evt_click_0"), "click", func() { comp8.Dec() })
	reactive.On(reactive.GetEl("comp8_evt_click_1"), "click", func() { comp8.Inc() })
	reactive.On(reactive.GetEl("comp8_evt_click_2"), "click", func() { comp8.Reset() })
	comp8.OnMount()
	}

	comp9 := &Button{
		Label: reactive.New[string](""),
		Variant: reactive.New[string](""),
	}

	comp9_el := reactive.GetEl("comp9")
	if !comp9_el.IsNull() && !comp9_el.IsUndefined() {
	reactive.InjectStyle("Button", buttonCSS)
	comp9.Label.Set("Toggle Card")

	comp9.Variant.Set("secondary")

	reactive.Bind("comp9_expr_Label_0", comp9.Label)
	reactive.BindAttr("[data-attrbind=\"comp9_attr0\"]", "class", `btn {Variant}`, "Variant", comp9.Variant)
	reactive.On(comp9_el, "click", func() { component.ToggleCard() })
	}

	component.OnMount()

	select {}
}
