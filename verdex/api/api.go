package api

import "github.com/verdexlab/verdex/verdex/core"

// Origin of Verdex API
var apiBaseUrl = "https://api.verdexlab.workers.dev"

// User-Agent used to call API
var apiUserAgent = "verdex-cli-" + core.GetVerdexVersion()
