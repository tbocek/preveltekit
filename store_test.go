package preveltekit

import "testing"

func TestStoreMap(t *testing.T) {
	age := New(17)
	adult := age.Map(func(a int) bool { return a >= 18 })
	label := adult.Map(func(b bool) string {
		if b {
			return "adult"
		}
		return "minor"
	})
	if adult.Get() || label.Get() != "minor" {
		t.Fatalf("initial: adult=%v label=%q", adult.Get(), label.Get())
	}
	age.Set(30)
	if !adult.Get() || label.Get() != "adult" {
		t.Errorf("after Set(30): adult=%v label=%q", adult.Get(), label.Get())
	}
}

func TestComputed(t *testing.T) {
	first, last := New("John"), New("Doe")
	full := Computed(func() string { return first.Get() + " " + last.Get() })
	if full.Get() != "John Doe" {
		t.Fatalf("initial: %q", full.Get())
	}
	last.Set("Smith")
	if full.Get() != "John Smith" {
		t.Errorf("after last.Set: %q", full.Get())
	}
}

// Dependencies are recorded on every run: a is only a dependency while
// flag is true.
func TestComputedDynamicDependencies(t *testing.T) {
	flag, a := New(false), New(1)
	runs := 0
	c := Computed(func() int {
		runs++
		if flag.Get() {
			return a.Get()
		}
		return 0
	})
	a.Set(2) // not read yet: no recompute
	if runs != 1 || c.Get() != 0 {
		t.Fatalf("a not a dependency yet: runs=%d c=%d", runs, c.Get())
	}
	flag.Set(true)
	a.Set(3)
	if c.Get() != 3 {
		t.Errorf("after flag+a: c=%d", c.Get())
	}
	flag.Set(false)
	before := runs
	a.Set(4) // no longer read: no recompute
	if runs != before {
		t.Errorf("a still a dependency after flag=false: %d extra runs", runs-before)
	}
}

// Re-subscribing drops removed callbacks, so the list does not grow.
func TestCallbacksCompact(t *testing.T) {
	flag, a := New(true), New(0)
	Computed(func() int {
		if flag.Get() {
			return a.Get()
		}
		return -1
	})
	for range 100 {
		flag.Set(!flag.Get())
	}
	if n := len(a.callbacks); n > 2 {
		t.Errorf("a has %d callback entries after 100 toggles", n)
	}
}
