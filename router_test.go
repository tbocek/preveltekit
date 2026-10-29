//go:build !wasm

package preveltekit

import "testing"

type userPage struct{}

func (u *userPage) Render() Node { return Div("user") }

type aboutPage struct{}

func (a *aboutPage) Render() Node { return Div("about") }

func TestRouterParam(t *testing.T) {
	user, about := &userPage{}, &aboutPage{}
	current := New[Component](about)
	router := NewRouter(current, []Route{
		{Path: "/users/:id", Component: user},
		{Path: "/about", Component: about},
	}, "r")
	id := router.Param("id")
	changes := 0
	id.OnChange(func(string) { changes++ })

	SetSSRPath("/users/42")
	router.Start()
	if id.Get() != "42" || current.Get() != Component(user) {
		t.Fatalf("after start: id=%q, component=%T", id.Get(), current.Get())
	}

	// Same route, other parameter: the store updates, the component stays
	router.handleRoute("/users/7")
	if id.Get() != "7" || changes != 2 {
		t.Errorf("after /users/7: id=%q, changes=%d", id.Get(), changes)
	}

	// A route without the parameter
	router.handleRoute("/about")
	if id.Get() != "" || current.Get() != Component(about) {
		t.Errorf("after /about: id=%q, component=%T", id.Get(), current.Get())
	}

	// Param after navigation starts with the current value
	router.handleRoute("/users/9")
	if got := router.Param("id"); got != id {
		t.Error("Param returned a new store for the same name")
	}
	if got := router.Param("other").Get(); got != "" {
		t.Errorf("unknown parameter: %q", got)
	}
}
