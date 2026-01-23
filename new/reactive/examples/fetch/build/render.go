//go:build !js || !wasm

package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
)

func main() {
	component := &FetchDemo{
		Status: New[string](""),
		RawData: New[string](""),
		Users: NewList[string](),
		UserCount: New[int](0),
	}

	component.OnMount()

	// Render template with current values
	html := `<div class="app">
	<h1>Fetch Demo</h1>
	<p class="status">Status: <strong><span id="expr_Status_0"></span></strong></p>

	<section>
		<h2>1. Simple Fetch</h2>
		<p>Fetch raw JSON data:</p>
		<button id="evt_click_0">Fetch Todo</button>
		<pre><span id="expr_RawData_1"></span></pre>
	</section>

	<section>
		<h2>2. Fetch into List (Diff Demo)</h2>
		<p>Fetches user names and uses List.Set() which triggers diff:</p>
		<div class="button-row">
			<button id="evt_click_1">Fetch All Users</button>
			<button id="evt_click_2">Fetch 3 Users</button>
			<button id="evt_click_3">Add Local</button>
			<button id="evt_click_4">Clear</button>
		</div>

		<p>Users: <span id="expr_UserCount_2"></span></p>

		<span id="if0_anchor"></span>

		<p class="note">
			Try: Load all → Load 3 (watch items get removed via diff)<br>
			Or: Load 3 → Load all (watch items get added via diff)
		</p>
	</section>
</div>`
	html = strings.Replace(html, "<span id=\"expr_Status_0\"></span>", fmt.Sprintf("<span id=\"expr_Status_0\">%v</span>", component.Status.Get()), 1)
	html = strings.Replace(html, "<span id=\"expr_RawData_1\"></span>", fmt.Sprintf("<span id=\"expr_RawData_1\">%v</span>", component.RawData.Get()), 1)
	html = strings.Replace(html, "<span id=\"expr_UserCount_2\"></span>", fmt.Sprintf("<span id=\"expr_UserCount_2\">%v</span>", component.UserCount.Get()), 1)
	// Render if block if0
	{
		var ifContent string
		if component.UserCount.Get() > 0 {
			ifContent = `
			<ul>
				<span id="each0_anchor"></span>
			</ul>
		`
		} else {
			ifContent = `
			<p class="empty">No users loaded</p>
		`
		}
		html = strings.Replace(html, "<span id=\"if0_anchor\"></span>", ifContent + "<span id=\"if0_anchor\"></span>", 1)
	}
	// Render each block for Users
	{
		var eachContent strings.Builder
		items := component.Users.Get()
		for i, user := range items {
			itemHTML := `
					<li><span class="index">{i}</span> {user}</li>
				`
			itemHTML = strings.ReplaceAll(itemHTML, "{user}", fmt.Sprintf("%v", user))
			itemHTML = strings.ReplaceAll(itemHTML, "{i}", strconv.Itoa(i))
			eachContent.WriteString("<span id=\"each0_" + strconv.Itoa(i) + "\">" + itemHTML + "</span>")
		}
		html = strings.Replace(html, "<span id=\"each0_anchor\"></span>", eachContent.String() + "<span id=\"each0_anchor\"></span>", 1)
	}

	fmt.Fprint(os.Stdout, html)
}

// Suppress unused import warnings
var _ = strings.Replace
var _ = strconv.Itoa
