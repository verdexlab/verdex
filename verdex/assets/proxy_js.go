//go:build js
// +build js

package assets

import (
	"net/http"
)

var proxyHost = "proxy.verdexlab.workers.dev"

// Proxy for wasm:js targets
func proxifyRequest(request *http.Request) {
	request.URL.Path = request.URL.Host + request.URL.Path
	request.URL.Host = proxyHost
}
