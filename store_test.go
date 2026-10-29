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
