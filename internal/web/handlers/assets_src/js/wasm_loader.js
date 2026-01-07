async function initNetutilWasm() {
  // wasm_exec.js must be loaded first; it defines global Go.
  if (typeof Go === "undefined") {
    throw new Error("Go WASM runtime missing (wasm_exec.js not loaded)");
  }

  const go = new Go();
  const resp = await fetch("/js/netutil.wasm");
  if (!resp.ok) {
    throw new Error(`failed to fetch netutil.wasm: ${resp.status}`);
  }

  const bytes = await resp.arrayBuffer();
  const { instance } = await WebAssembly.instantiate(bytes, go.importObject);

  // This runs your Go main(), which registers:
  // netutil_isValidIPv4 / netutil_isValidIPv6 / netutil_isValidMAC
  go.run(instance);

  // Optional: provide a nicer API
  window.netutil = {
    isValidIPv4: (s) => window.netutil_isValidIPv4(s),
    isValidIPv6: (s) => window.netutil_isValidIPv6(s),
    isValidMAC:  (s) => window.netutil_isValidMAC(s),
  };
}

initNetutilWasm().catch((err) => console.error("[netutil wasm]", err));
