package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/core"
)

var maxInputVersions = 10

var client = &http.Client{
	Timeout: 10 * time.Second,
}

type CVEsData struct {
	Product               string                     `json:"product"`
	Versions              []string                   `json:"versions"`
	Total                 int                        `json:"total"`
	CVEs                  []core.CVE                 `json:"cves"`
	Vulnerable            bool                       `json:"vulnerable"`
	UpdateRecommendations *updateRecommendationsData `json:"update_recommendations"`
	LatestUpdateAt        time.Time                  `json:"latest_update_at"`
}

type updateRecommendationsData struct {
	WithoutVulnerabilities         *string `json:"without_vulnerabilities"`
	WithoutCriticalVulnerabilities *string `json:"without_critical_vulnerabilities"`
}

func GetCVEsFromVersions(product string, versions []string, apiKey string) (data *CVEsData, isApiKeyValid bool, err error) {
	versionsCount := maxInputVersions
	if len(versions) < maxInputVersions {
		versionsCount = len(versions)
	}

	versionsList := strings.Join(versions[:versionsCount], ",")
	url := fmt.Sprintf("%s/cves/resolve?product=%s&versions=%s", apiBaseUrl, product, versionsList)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Debug().
			Err(err).
			Str("url", url).
			Msg("Failed to resolve CVEs from API")

		return nil, false, err
	}

	request.Header.Set("X-Verdex-Key", apiKey)
	request.Header.Set("User-Agent", apiUserAgent)

	res, err := client.Do(request)
	if err != nil {
		l := log.Debug().
			Err(err).
			Str("url", url)

		if res != nil {
			l = l.Int("status_code", res.StatusCode)
		}

		l.Msg("Failed to resolve CVEs from API")
		return nil, false, err
	}

	defer res.Body.Close()

	if res.StatusCode == 400 {
		fmt.Println("")
		log.Info().Msgf("Vulnerabilities listing is not yet available for %s", product)
		return nil, false, nil
	} else if res.StatusCode != 200 {
		fmt.Println("")
		log.Error().Msgf("An error occurred while retrieving vulnerabilities for %s", product)
		return nil, false, nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Debug().
			Err(err).
			Str("url", url).
			Int("status_code", res.StatusCode).
			Msg("Failed to read API body")
		return nil, false, err
	}

	resData := CVEsData{}
	err = json.Unmarshal(body, &resData)
	if err != nil {
		log.Debug().
			Err(err).
			Str("url", url).
			Msg("Failed to unmarshal API body")
		return nil, false, err
	}

	isApiKeyValid = res.Header.Get("X-Verdex-Key-Status") == "OK"
	return &resData, isApiKeyValid, nil
}
