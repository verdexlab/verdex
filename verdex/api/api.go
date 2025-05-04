package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/verdexlab/verdex/verdex/core"
)

// Origin of Verdex API
var apiBaseUrl = "https://api.verdexlab.io"

// User-Agent used to call API
var apiUserAgent = fmt.Sprintf("verdex-cli-%s-%s", core.GetEnvironment(), core.GetVerdexVersion())

var client = &http.Client{
	Timeout: 10 * time.Second,
}
