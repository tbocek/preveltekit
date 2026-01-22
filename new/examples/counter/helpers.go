//go:build js && wasm

package main

// Local = reactive
func double(x int) int {
    return x * 2
}

var multiplier = 2  // reactive if used