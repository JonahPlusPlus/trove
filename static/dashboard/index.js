window.onload = function() {
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("dashboard/index.wasm"), go.importObject).then((result) => {
        go.run(result.instance);
    });
}
