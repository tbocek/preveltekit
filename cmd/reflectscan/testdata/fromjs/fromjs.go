package fromjs

import (
	"syscall/js"

	p "github.com/tbocek/preveltekit/v2"
)

type Price struct {
	Amounts map[string]float64
}

func (r *Price) FromJS(v js.Value) {}

func Load() (Price, error) { return p.Get[Price]("/price") }
