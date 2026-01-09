//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/rannday/netaddr/ip"
	"github.com/rannday/netaddr/mac"
)

func wrap2(fn func(string) (string, error)) js.Func {
  return js.FuncOf(func(this js.Value, args []js.Value) any {
    if len(args) < 1 {
      return map[string]any{"ok": false, "err": "missing argument", "value": ""}
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
  js.Global().Set("netaddr_isValidIPv4", wrap2(ip.IsValidIPv4))
  js.Global().Set("netaddr_isValidIPv6", wrap2(ip.IsValidIPv6))
  js.Global().Set("netaddr_isValidMAC", wrap2(mac.IsValidMAC))
  js.Global().Set("netaddr_formatMAC", wrap2(mac.FormatMAC))

  select {}
}
