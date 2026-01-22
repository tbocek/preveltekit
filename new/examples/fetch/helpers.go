//go:build js && wasm

package main

import "syscall/js"

func fetch(url string, onDone func(data string, err bool)) {
	go func() {
		promise := js.Global().Call("fetch", url)

		promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
			args[0].Call("text").Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
				onDone(args[0].String(), false)
				return nil
			}))
			return nil
		})).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) any {
			onDone("", true)
			return nil
		}))
	}()
}