//go:build !js || !wasm

package main

import (
	"fmt"
	"os"
	"reactive"
	"strings"
)

func main() {
	component := &App{
		Count: reactive.New[int](0),
		Message: reactive.New[string](""),
		Theme: reactive.New[string](""),
		ShowCard: reactive.New[bool](false),
	}

	component.OnMount()

	comp0 := &Button{
		Label: reactive.New[string](""),
		Variant: reactive.New[string](""),
	}
	comp0.Variant.Set("primary")

	comp1 := &Card{
		Title: reactive.New[string](""),
	}
	comp1.Title.Set("Status Card")

	comp2 := &Card{
		Title: reactive.New[string](""),
	}
	comp2.Title.Set("Toggleable Card")

	comp3 := &Button{
		Label: reactive.New[string](""),
		Variant: reactive.New[string](""),
	}
	comp3.Label.Set("Primary")
	comp3.Variant.Set("primary")

	comp4 := &Button{
		Label: reactive.New[string](""),
		Variant: reactive.New[string](""),
	}
	comp4.Label.Set("Secondary")
	comp4.Variant.Set("secondary")

	comp5 := &Button{
		Label: reactive.New[string](""),
		Variant: reactive.New[string](""),
	}
	comp5.Label.Set("Success")
	comp5.Variant.Set("success")

	comp6 := &Button{
		Label: reactive.New[string](""),
		Variant: reactive.New[string](""),
	}
	comp6.Label.Set("Danger")
	comp6.Variant.Set("danger")

	comp7 := &Button{
		Label: reactive.New[string](""),
		Variant: reactive.New[string](""),
	}
	comp7.Label.Set(component.Message.Get())
	comp7.Variant.Set("primary")

	comp8 := &Counter{
		Initial: reactive.New[int](0),
		Step: reactive.New[int](0),
		Value: reactive.New[int](0),
	}
	comp8.Initial.Set(10)
	comp8.Step.Set(5)
	comp8.OnMount()

	comp9 := &Button{
		Label: reactive.New[string](""),
		Variant: reactive.New[string](""),
	}
	comp9.Label.Set("Toggle Card")
	comp9.Variant.Set("secondary")

	html := `<div class="app">
	<h1>Component Composition</h1>

	<section>
		<h2>1. Multiple Instances with Different Props</h2>
		<p>Same Button component, different configurations:</p>
		<div class="button-row">
			<!--comp3-->
			<!--comp4-->
			<!--comp5-->
			<!--comp6-->
		</div>
		<p>Count: <strong><span id="expr_Count_0"></span></strong></p>
	</section>

	<section>
		<h2>2. Dynamic Props</h2>
		<p>Button label updates when parent state changes:</p>
		<!--comp7-->
	</section>

	<section>
		<h2>3. Slots with Parent-Bound Content</h2>
		<p>Content inside component tags becomes slot content:</p>
		<!--comp0-->
		<!--comp1-->
	</section>

	<section>
		<h2>4. Component with Internal State</h2>
		<p>Counter has its own internal state:</p>
		<!--comp8-->
	</section>

	<section>
		<h2>5. Conditional Rendering</h2>
		<!--comp9-->
		<span id="if0_anchor"></span>
	</section>

	<section>
		<h2>Event Log</h2>
		<p class="message"><span id="expr_Message_1"></span></p>
	</section>
</div>`
	html = strings.Replace(html, "<span id=\"expr_Count_0\"></span>", fmt.Sprintf("<span id=\"expr_Count_0\">%v</span>", component.Count.Get()), 1)
	html = strings.Replace(html, "<span id=\"expr_Message_1\"></span>", fmt.Sprintf("<span id=\"expr_Message_1\">%v</span>", component.Message.Get()), 1)
	{
		var ifContent string
		if component.ShowCard.Get() {
			ifContent = `
			<!--comp2-->
		`
		} else {
			ifContent = `
			<p class="hidden-note">Card is hidden</p>
		`
		}
		html = strings.Replace(html, "<span id=\"if0_anchor\"></span>", ifContent + "<span id=\"if0_anchor\"></span>", 1)
	}
	{
		childHTML := `<button id="comp0" class="btn" data-attrbind="comp0_attr0">
		Count is <span id="comp0_expr_Count_0"></span>
		<span id="comp0_expr_Label_1"></span>
	</button>`
		childHTML = strings.Replace(childHTML, "<span id=\"comp0_expr_Label_1\"></span>", fmt.Sprintf("<span id=\"comp0_expr_Label_1\">%v</span>", comp0.Label.Get()), 1)
		childHTML = strings.Replace(childHTML, "<span id=\"comp0_expr_Count_0\"></span>", fmt.Sprintf("<span id=\"comp0_expr_Count_0\">%v</span>", component.Count.Get()), 1)
		{
			attrVal := `btn {Variant}`
			attrVal = strings.ReplaceAll(attrVal, "{Variant}", fmt.Sprintf("%v", comp0.Variant.Get()))
			childHTML = strings.Replace(childHTML, "class=\"btn\"", "class=\"" + attrVal + "\"", 1)
		}
		html = strings.Replace(html, "<!--comp0-->", childHTML, 1)
	}
	{
		childHTML := `<div id="comp1" class="card">
	<div class="card-header"><span id="comp1_expr_Title_0"></span></div>
	<div class="card-body"><p>Current count: <strong><span id="comp1_expr_Count_1"></span></strong></p>
			<p>Theme: <strong><span id="comp1_expr_Theme_2"></span></strong></p></div>
</div>`
		childHTML = strings.Replace(childHTML, "<span id=\"comp1_expr_Title_0\"></span>", fmt.Sprintf("<span id=\"comp1_expr_Title_0\">%v</span>", comp1.Title.Get()), 1)
		childHTML = strings.Replace(childHTML, "<span id=\"comp1_expr_Count_1\"></span>", fmt.Sprintf("<span id=\"comp1_expr_Count_1\">%v</span>", component.Count.Get()), 1)
		childHTML = strings.Replace(childHTML, "<span id=\"comp1_expr_Theme_2\"></span>", fmt.Sprintf("<span id=\"comp1_expr_Theme_2\">%v</span>", component.Theme.Get()), 1)
		html = strings.Replace(html, "<!--comp1-->", childHTML, 1)
	}
	{
		childHTML := `<div id="comp2" class="card">
	<div class="card-header"><span id="comp2_expr_Title_0"></span></div>
	<div class="card-body"><p>This card can be shown/hidden.</p>
				<p>Parent count: <span id="comp2_expr_Count_1"></span></p></div>
</div>`
		childHTML = strings.Replace(childHTML, "<span id=\"comp2_expr_Title_0\"></span>", fmt.Sprintf("<span id=\"comp2_expr_Title_0\">%v</span>", comp2.Title.Get()), 1)
		childHTML = strings.Replace(childHTML, "<span id=\"comp2_expr_Count_1\"></span>", fmt.Sprintf("<span id=\"comp2_expr_Count_1\">%v</span>", component.Count.Get()), 1)
		html = strings.Replace(html, "<!--comp2-->", childHTML, 1)
	}
	{
		childHTML := `<button id="comp3" class="btn" data-attrbind="comp3_attr0">
		
		<span id="comp3_expr_Label_0"></span>
	</button>`
		childHTML = strings.Replace(childHTML, "<span id=\"comp3_expr_Label_0\"></span>", fmt.Sprintf("<span id=\"comp3_expr_Label_0\">%v</span>", comp3.Label.Get()), 1)
		{
			attrVal := `btn {Variant}`
			attrVal = strings.ReplaceAll(attrVal, "{Variant}", fmt.Sprintf("%v", comp3.Variant.Get()))
			childHTML = strings.Replace(childHTML, "class=\"btn\"", "class=\"" + attrVal + "\"", 1)
		}
		html = strings.Replace(html, "<!--comp3-->", childHTML, 1)
	}
	{
		childHTML := `<button id="comp4" class="btn" data-attrbind="comp4_attr0">
		
		<span id="comp4_expr_Label_0"></span>
	</button>`
		childHTML = strings.Replace(childHTML, "<span id=\"comp4_expr_Label_0\"></span>", fmt.Sprintf("<span id=\"comp4_expr_Label_0\">%v</span>", comp4.Label.Get()), 1)
		{
			attrVal := `btn {Variant}`
			attrVal = strings.ReplaceAll(attrVal, "{Variant}", fmt.Sprintf("%v", comp4.Variant.Get()))
			childHTML = strings.Replace(childHTML, "class=\"btn\"", "class=\"" + attrVal + "\"", 1)
		}
		html = strings.Replace(html, "<!--comp4-->", childHTML, 1)
	}
	{
		childHTML := `<button id="comp5" class="btn" data-attrbind="comp5_attr0">
		
		<span id="comp5_expr_Label_0"></span>
	</button>`
		childHTML = strings.Replace(childHTML, "<span id=\"comp5_expr_Label_0\"></span>", fmt.Sprintf("<span id=\"comp5_expr_Label_0\">%v</span>", comp5.Label.Get()), 1)
		{
			attrVal := `btn {Variant}`
			attrVal = strings.ReplaceAll(attrVal, "{Variant}", fmt.Sprintf("%v", comp5.Variant.Get()))
			childHTML = strings.Replace(childHTML, "class=\"btn\"", "class=\"" + attrVal + "\"", 1)
		}
		html = strings.Replace(html, "<!--comp5-->", childHTML, 1)
	}
	{
		childHTML := `<button id="comp6" class="btn" data-attrbind="comp6_attr0">
		
		<span id="comp6_expr_Label_0"></span>
	</button>`
		childHTML = strings.Replace(childHTML, "<span id=\"comp6_expr_Label_0\"></span>", fmt.Sprintf("<span id=\"comp6_expr_Label_0\">%v</span>", comp6.Label.Get()), 1)
		{
			attrVal := `btn {Variant}`
			attrVal = strings.ReplaceAll(attrVal, "{Variant}", fmt.Sprintf("%v", comp6.Variant.Get()))
			childHTML = strings.Replace(childHTML, "class=\"btn\"", "class=\"" + attrVal + "\"", 1)
		}
		html = strings.Replace(html, "<!--comp6-->", childHTML, 1)
	}
	{
		childHTML := `<button id="comp7" class="btn" data-attrbind="comp7_attr0">
		
		<span id="comp7_expr_Label_0"></span>
	</button>`
		childHTML = strings.Replace(childHTML, "<span id=\"comp7_expr_Label_0\"></span>", fmt.Sprintf("<span id=\"comp7_expr_Label_0\">%v</span>", comp7.Label.Get()), 1)
		{
			attrVal := `btn {Variant}`
			attrVal = strings.ReplaceAll(attrVal, "{Variant}", fmt.Sprintf("%v", comp7.Variant.Get()))
			childHTML = strings.Replace(childHTML, "class=\"btn\"", "class=\"" + attrVal + "\"", 1)
		}
		html = strings.Replace(html, "<!--comp7-->", childHTML, 1)
	}
	{
		childHTML := `<div id="comp8" class="counter">
	<span class="value"><span id="comp8_expr_Value_0"></span></span>
	<button class="counter-btn" id="comp8_evt_click_0">-<span id="comp8_expr_Step_1"></span></button>
	<button class="counter-btn" id="comp8_evt_click_1">+<span id="comp8_expr_Step_2"></span></button>
	<button class="counter-btn reset" id="comp8_evt_click_2">Reset</button>
</div>`
		childHTML = strings.Replace(childHTML, "<span id=\"comp8_expr_Value_0\"></span>", fmt.Sprintf("<span id=\"comp8_expr_Value_0\">%v</span>", comp8.Value.Get()), 1)
		childHTML = strings.Replace(childHTML, "<span id=\"comp8_expr_Step_1\"></span>", fmt.Sprintf("<span id=\"comp8_expr_Step_1\">%v</span>", comp8.Step.Get()), 1)
		childHTML = strings.Replace(childHTML, "<span id=\"comp8_expr_Step_2\"></span>", fmt.Sprintf("<span id=\"comp8_expr_Step_2\">%v</span>", comp8.Step.Get()), 1)
		html = strings.Replace(html, "<!--comp8-->", childHTML, 1)
	}
	{
		childHTML := `<button id="comp9" class="btn" data-attrbind="comp9_attr0">
		
		<span id="comp9_expr_Label_0"></span>
	</button>`
		childHTML = strings.Replace(childHTML, "<span id=\"comp9_expr_Label_0\"></span>", fmt.Sprintf("<span id=\"comp9_expr_Label_0\">%v</span>", comp9.Label.Get()), 1)
		{
			attrVal := `btn {Variant}`
			attrVal = strings.ReplaceAll(attrVal, "{Variant}", fmt.Sprintf("%v", comp9.Variant.Get()))
			childHTML = strings.Replace(childHTML, "class=\"btn\"", "class=\"" + attrVal + "\"", 1)
		}
		html = strings.Replace(html, "<!--comp9-->", childHTML, 1)
	}

	fmt.Fprint(os.Stdout, html)
}

var _ = strings.Replace
