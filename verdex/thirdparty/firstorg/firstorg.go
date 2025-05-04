package firstorg

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/verdexlab/verdex/verdex/core"
)

// Maximum CVEs retrieved at the same time from first.org API
var cvesChunkSize = 100

var client = &http.Client{
	Timeout: 10 * time.Second,
}

// Get CVEs EPSS data from first.org API
func GetCvesEpss(cveIDs []string) ([]*EpssData, error) {
	data := make([]*EpssData, 0)

	for chunkCveIDs := range slices.Chunk(cveIDs, cvesChunkSize) {
		url := "https://api.first.org/data/v1/epss?cve=" + strings.Join(chunkCveIDs, ",")
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return []*EpssData{}, err
		}

		request.Header.Set("Content-Type", "application/json")

		core.ProxifyRequest(request)

		res, err := client.Do(request)
		if err != nil {
			return []*EpssData{}, err
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return []*EpssData{}, err
		}

		var chunkData EpssResponse
		err = json.Unmarshal(body, &chunkData)
		if err != nil {
			return []*EpssData{}, err
		}

		data = append(data, chunkData.Data...)
	}

	return data, nil
}
