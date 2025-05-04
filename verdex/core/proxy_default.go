//go:build !js
// +build !js

package core

import "net/http"

func ProxifyRequest(request *http.Request) {
	// no proxy by default
}
