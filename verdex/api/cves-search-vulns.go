package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type apiSearchVulnsResponse map[string]struct {
	Vulns map[string]ApiSearchVulnsCVE `json:"vulns"`
}

type ApiSearchVulnsCVE struct {
	ID                 string   `json:"id"`
	Cvss               string   `json:"cvss"`
	Description        string   `json:"description"`
	Aliases            []string `json:"aliases"`
	CisaKnownExploited bool     `json:"cisa_known_exploited"`
	VulnMatchReason    string   `json:"vuln_match_reason"`
	Published          string   `json:"published"`
	Modified           string   `json:"modified"`
}

func SearchVulns(cpe string) (map[string]ApiSearchVulnsCVE, error) {
	url := fmt.Sprintf("%s/cves/search-vulns?cpe=%s", apiBaseUrl, cpe)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Debug().
			Err(err).
			Str("url", url).
			Msg("failed to retrieve vulnerabilities from Verdex API")

		return nil, err
	}

	request.Header.Set("User-Agent", apiUserAgent)

	res, err := client.Do(request)
	if err != nil {
		l := log.Debug().
			Err(err).
			Str("url", url)

		if res != nil {
			l = l.Int("status_code", res.StatusCode)
		}

		l.Msg("failed to retrieve vulnerabilities from Verdex API")
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Debug().Msg("An error occurred while retrieving vulnerabilities from Verdex API")
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Debug().
			Err(err).
			Str("url", url).
			Int("status_code", res.StatusCode).
			Msg("Failed to read API body")
		return nil, err
	}

	data := apiSearchVulnsResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Debug().
			Err(err).
			Str("url", url).
			Msg("Failed to unmarshal API body")
		return nil, err
	}

	return data[cpe].Vulns, nil
}
