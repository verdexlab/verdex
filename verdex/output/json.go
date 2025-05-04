package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/core"
	"github.com/verdexlab/verdex/verdex/products"
)

type OutputJson struct {
	Scanner   string             `json:"scanner"`
	Templates string             `json:"templates"`
	Results   []OutputResultJson `json:"results"`
}

type OutputResultJson struct {
	StartedAt             string                               `json:"startedAt"`
	EndedAt               string                               `json:"endedAt"`
	Target                string                               `json:"target"`
	Product               string                               `json:"product"`
	Success               bool                                 `json:"success"`
	PossibleVersions      []OutputResultVersionJson            `json:"possibleVersions"`
	CVEs                  []*core.CVE                          `json:"cves"`
	UpdateRecommendations *core.DetectionUpdateRecommendations `json:"update_recommendations,omitempty"`
}

type OutputResultVersionJson struct {
	Version string `json:"version"`
	Cpe     string `json:"cpe"`
}

// Export results to JSON output file
func exportResultsJson(execution *core.Execution, path string) error {
	output := OutputJson{
		Scanner: fmt.Sprintf("verdex@%s", core.GetVerdexVersion()),
		Results: []OutputResultJson{},
	}

	if execution.Config.TemplatesSource == core.TemplatesSourceLocalDirectory {
		output.Templates = fmt.Sprintf("file:%s", execution.Config.TemplatesDirectory)
	} else {
		output.Templates = fmt.Sprintf("github:%s/%s@%s", execution.Config.TemplatesOrganization, execution.Config.TemplatesRepository, execution.Config.TemplatesRelease)
	}

	for _, detection := range execution.Detections {
		resultOutput := OutputResultJson{
			StartedAt:             detection.StartedAt.Format(time.RFC3339),
			EndedAt:               detection.EndedAt.Format(time.RFC3339),
			Target:                detection.Target,
			Product:               detection.Product,
			Success:               detection.Success,
			PossibleVersions:      []OutputResultVersionJson{},
			CVEs:                  detection.CVEs,
			UpdateRecommendations: detection.UpdateRecommendations,
		}

		if detection.Success {
			product := products.GetProduct(detection.Product)

			for _, version := range detection.PossibleVersions {
				resultOutput.PossibleVersions = append(resultOutput.PossibleVersions, OutputResultVersionJson{
					Version: version.String(),
					Cpe:     product.Cpe.Build(version.String()),
				})
			}
		}

		output.Results = append(output.Results, resultOutput)
	}

	jsonString, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Debug().Err(err).Msg("Failed to marshal results to json")
		return err
	}

	if core.GetEnvironment() == core.EnvironmentReleaseWasmJS {
		fmt.Println("wasm:js:output", string(jsonString))
	} else {
		err = os.WriteFile(path, jsonString, os.ModePerm)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to write results to json file")
			return err
		}
	}

	return nil
}
