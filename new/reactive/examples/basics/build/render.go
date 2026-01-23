//go:build !js || !wasm

package main

import (
	"fmt"
	"os"
	"reactive"
	"strings"
)

func main() {
	component := &Basics{
		Count: reactive.New[int](0),
		Name: reactive.New[string](""),
		Message: reactive.New[string](""),
		DarkMode: reactive.New[bool](false),
		Agreed: reactive.New[bool](false),
		Score: reactive.New[int](0),
	}

	component.OnMount()

	html := `<div class="app">
	<h1>Basics Demo</h1>

	<section>
		<h2>Counter</h2>
		<p>Count: <strong><span id="expr_Count_0"></span></strong></p>
		<button id="evt_click_0">-1</button>
		<button id="evt_click_1">+1</button>
		<button id="evt_click_2">+5</button>
		<button id="evt_click_3">Double</button>
		<button id="evt_click_4">Reset</button>
	</section>

	<section>
		<h2>Conditionals</h2>
		<p>Score: <span id="expr_Score_1"></span></p>
		<span id="if0_anchor"></span>
		<div class="buttons">
			<button id="evt_click_5">A</button>
			<button id="evt_click_6">B</button>
			<button id="evt_click_7">C</button>
			<button id="evt_click_8">D</button>
			<button id="evt_click_9">F</button>
		</div>
	</section>

	<section>
		<h2>Two-Way Binding</h2>
		<div>
			<label>Your name: <input type="text" id="bind0" placeholder="Enter name"></label>
		</div>
		<p>Hello, <span id="expr_Name_2"></span>!</p>
	</section>

	<section>
		<h2>Checkbox Binding</h2>
		<label>
			<input type="checkbox" id="bind1"> Dark Mode
		</label>
		<div id="class0">
			This box uses dark mode styling when checked.
		</div>
	</section>

	<section>
		<h2>Form with Event Modifier</h2>
		<form id="evt_submit_10">
			<div>
				<label>Name: <input type="text" id="bind2" placeholder="Your name"></label>
			</div>
			<div>
				<label>
					<input type="checkbox" id="bind3"> I agree to the terms
				</label>
			</div>
			<button type="submit">Submit</button>
		</form>
		<p class="message"><span id="expr_Message_3"></span></p>
	</section>
</div>`
	html = strings.Replace(html, "<span id=\"expr_Count_0\"></span>", fmt.Sprintf("<span id=\"expr_Count_0\">%v</span>", component.Count.Get()), 1)
	html = strings.Replace(html, "<span id=\"expr_Score_1\"></span>", fmt.Sprintf("<span id=\"expr_Score_1\">%v</span>", component.Score.Get()), 1)
	html = strings.Replace(html, "<span id=\"expr_Name_2\"></span>", fmt.Sprintf("<span id=\"expr_Name_2\">%v</span>", component.Name.Get()), 1)
	html = strings.Replace(html, "<span id=\"expr_Message_3\"></span>", fmt.Sprintf("<span id=\"expr_Message_3\">%v</span>", component.Message.Get()), 1)
	{
		var ifContent string
		if component.Score.Get() >= 90 {
			ifContent = `
			<p class="grade a">Grade: A - Excellent!</p>
		`
		} else if component.Score.Get() >= 80 {
			ifContent = `
			<p class="grade b">Grade: B - Good</p>
		`
		} else if component.Score.Get() >= 70 {
			ifContent = `
			<p class="grade c">Grade: C - Average</p>
		`
		} else if component.Score.Get() >= 60 {
			ifContent = `
			<p class="grade d">Grade: D - Below Average</p>
		`
		} else {
			ifContent = `
			<p class="grade f">Grade: F - Failing</p>
		`
		}
		html = strings.Replace(html, "<span id=\"if0_anchor\"></span>", ifContent + "<span id=\"if0_anchor\"></span>", 1)
	}

	fmt.Fprint(os.Stdout, html)
}

var _ = strings.Replace
