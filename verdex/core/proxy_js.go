//go:build js
// +build js

package core

import (
	"net/http"
)

var proxyHost = "proxy.verdexlab.io"

// Proxy for wasm:js targets
func ProxifyRequest(request *http.Request) {
	request.URL.Path = request.URL.Host + request.URL.Path
	request.URL.Host = proxyHost
}
