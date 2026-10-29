package maps

import p "github.com/tbocek/preveltekit/v2"

type Rates struct {
	Base  string             `js:"base"`
	Rates map[string]float64 `js:"rates"`
	note  []int              // unexported: ignored by the codec
}

func Load() (Rates, error) { return p.Get[Rates]("/rates") }
