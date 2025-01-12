//go:build !js
// +build !js

package assets

import "net/http"

func proxifyRequest(request *http.Request) {
	// no proxy by default
}
