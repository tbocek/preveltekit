//go:build !js || !wasm

package main

func Fetch(url string, callback func(data string, err error)) {}
func FetchJSON(url string, callback func(data string, err error)) {}
