build-wasm:
	GOOS=js GOARCH=wasm go build -ldflags "-X github.com/verdexlab/verdex/verdex/core.releaseEnvironment=release-wasmjs" -o ./dist/verdex.wasm .
