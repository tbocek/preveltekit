package main

import (
	"testing"
)

func TestAnalyze(t *testing.T) {
	script := `count := 0
doubled := count * 2

func increment() {
    count++
}

func reset() {
    count = 0
}`

	a, err := analyze(script)
	if err != nil {
		t.Fatalf("analyze error: %v", err)
	}

	if _, ok := a.Vars["count"]; !ok {
		t.Error("missing count")
	}
	if _, ok := a.Vars["doubled"]; !ok {
		t.Error("missing doubled")
	}

	doubled := a.Vars["doubled"]
	if len(doubled.DependsOn) != 1 || doubled.DependsOn[0] != "count" {
		t.Errorf("doubled should depend on count, got: %v", doubled.DependsOn)
	}

	if _, ok := a.Funcs["increment"]; !ok {
		t.Error("missing increment func")
	}

	inc := a.Funcs["increment"]
	if len(inc.Modifies) != 1 || inc.Modifies[0] != "count" {
		t.Errorf("increment should modify count, got: %v", inc.Modifies)
	}

	countIdx, doubledIdx := -1, -1
	for i, name := range a.Order {
		switch name {
		case "count":
			countIdx = i
		case "doubled":
			doubledIdx = i
		}
	}

	if countIdx == -1 || doubledIdx == -1 {
		t.Errorf("missing vars in order: %v", a.Order)
	}
	if countIdx > doubledIdx {
		t.Errorf("wrong order: count=%d, doubled=%d", countIdx, doubledIdx)
	}
}

func TestCycleDetection(t *testing.T) {
	script := `x := y + 1
y := x + 1`

	a, err := analyze(script)
	if err != nil {
		t.Fatalf("analyze error: %v", err)
	}

	if len(a.Vars) != 2 {
		t.Errorf("should have 2 vars despite cycle, got: %d", len(a.Vars))
	}
}

func TestSimpleVar(t *testing.T) {
	script := `count := 0`

	a, err := analyze(script)
	if err != nil {
		t.Fatalf("analyze error: %v", err)
	}

	if len(a.Vars) != 1 {
		t.Errorf("expected 1 var, got: %d", len(a.Vars))
	}
}