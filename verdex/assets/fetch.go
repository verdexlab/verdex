package assets

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/core"
	"github.com/verdexlab/verdex/verdex/variables"
	"golang.org/x/sync/syncmap"
)

var userAgent = "Verdex - Open-Source scanning project"

var assets = syncmap.Map{}

var client = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// Fetch asset at given URL
func FetchAsset(execution *core.Execution, detection *core.Detection, method string, url string) (*Asset, error) {
	url = resolveVariablesFromURL(execution, detection, url)

	// Retrieve asset from store if exists
	assetsStoreKey := fmt.Sprintf("%s:%s", method, url)
	if assetValue, hasAsset := assets.Load(assetsStoreKey); hasAsset {
		return assetValue.(*Asset), nil
	}

	detection.TotalRequests++

	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Debug().
			Err(err).
			Str("request", fmt.Sprintf("%s %s", method, url)).
			Msg("Failed to fetch asset")

		return nil, err
	}

	request.Header.Set("User-Agent", userAgent)

	res, err := client.Do(request)
	if err != nil {
		l := log.Debug().
			Err(err).
			Str("request", fmt.Sprintf("%s %s", method, url))

		if res != nil {
			l = l.Int("status_code", res.StatusCode)
		}

		l.Msg("Failed to fetch asset")
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Debug().
			Err(err).
			Str("request", fmt.Sprintf("%s %s", method, url)).
			Int("status_code", res.StatusCode).
			Msg("Failed to read asset body")
		return nil, err
	}

	asset := Asset{
		StatusCode: res.StatusCode,
		Body:       string(body),
	}

	assets.Store(assetsStoreKey, &asset)

	log.Debug().
		Str("request", fmt.Sprintf("%s %s", method, url)).
		Int("status_code", res.StatusCode).
		Msg("Fetch asset")
	return &asset, nil
}

// Resolve all variables from given URL
func resolveVariablesFromURL(execution *core.Execution, detection *core.Detection, url string) string {
	for _, variableKey := range variables.GetAllProductVariables(detection.Product) {
		variableTag := fmt.Sprintf("{{%s}}", variableKey)
		if strings.Contains(url, variableTag) {
			variable := variables.GetProductVariable(detection.Product, variableKey)
			variableValue, err := GetVariableValue(execution, detection, variable)
			if err != nil {
				continue
			}

			url = strings.ReplaceAll(url, variableTag, variableValue)
		}
	}

	return url
}
