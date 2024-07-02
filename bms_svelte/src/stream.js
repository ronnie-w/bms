import "./app_exec.js"

const go = new Go(),
	stream = await WebAssembly.instantiateStreaming(fetch("http://localhost:8000/app.wasm"), go.importObject),
	_instance = go.run(stream.instance);

export default _instance;
