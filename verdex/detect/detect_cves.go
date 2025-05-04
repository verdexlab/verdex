package detect

import (
	"sort"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/api"
	"github.com/verdexlab/verdex/verdex/core"
	"github.com/verdexlab/verdex/verdex/products"
	"github.com/verdexlab/verdex/verdex/thirdparty/firstorg"
)

// Try to detect CVEs on given target related to detected versions
func DetectCVEs(execution *core.Execution, detection *core.Detection) ([]*core.CVE, error) {
	product := products.GetProduct(detection.Product)

	cves := make(map[string]*core.CVE, 0)

	for _, version := range detection.PossibleVersions {
		cpe := product.Cpe.Build(version.String())
		searchVulnsCVEs, err := api.SearchVulns(cpe)
		if err != nil {
			return nil, err
		}

		for cveID, searchVulnsCVE := range searchVulnsCVEs {
			if !strings.HasPrefix(cveID, "CVE-") {
				continue
			}

			if searchVulnsCVE.VulnMatchReason != "version_in_range" {
				continue
			}

			cvssScore, _ := strconv.ParseFloat(searchVulnsCVE.Cvss, 32)
			publishedAt := strings.Split(searchVulnsCVE.Published, " ")[0]

			if cve, cveExists := cves[cveID]; cveExists {
				cve.VulnerableVersions = append(cve.VulnerableVersions, version.String())
			} else {
				cve := core.CVE{
					ID:                 cveID,
					Description:        searchVulnsCVE.Description,
					CvssScore:          float32(cvssScore),
					IsKEV:              searchVulnsCVE.CisaKnownExploited,
					VulnerableVersions: []string{version.String()},
					PublishedAt:        publishedAt,
				}

				cves[cveID] = &cve
			}
		}
	}

	cveIDs := make([]string, 0)
	for cveID, _ := range cves {
		cveIDs = append(cveIDs, cveID)
	}

	epssData, err := firstorg.GetCvesEpss(cveIDs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve EPSS data from first.org API")
		return nil, err
	}

	result := make([]*core.CVE, 0)
	for _, cveEpssData := range epssData {
		cve := cves[cveEpssData.CveID]
		epssScore, _ := strconv.ParseFloat(cveEpssData.EPSS, 32)
		cve.EpssScore = float32(epssScore)
		result = append(result, cve)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[j].CvssScore < result[i].CvssScore
	})

	return result, nil
}
