//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/rannday/kea-web/internal/utils" // adjust import path
)

func wrap2(fn func(string) (string, error)) js.Func {
  return js.FuncOf(func(this js.Value, args []js.Value) any {
    if len(args) < 1 {
      return map[string]any{"ok": false, "err": "missing argument"}
    }
    s := args[0].String()
    out, err := fn(s)
    if err != nil {
      return map[string]any{"ok": false, "err": err.Error(), "value": ""}
    }
    return map[string]any{"ok": true, "err": "", "value": out}
  })
}

func main() {
  js.Global().Set("netutil_isValidIPv4", wrap2(utils.IsValidIPv4))
  js.Global().Set("netutil_isValidIPv6", wrap2(utils.IsValidIPv6))
  js.Global().Set("netutil_isValidMAC", wrap2(utils.IsValidMAC))

  // keep running
  select {}
}
