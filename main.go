package main

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/core"
	"github.com/verdexlab/verdex/verdex/detect"
	"github.com/verdexlab/verdex/verdex/output"
	"github.com/verdexlab/verdex/verdex/products"
	"github.com/verdexlab/verdex/verdex/templates"
	"github.com/verdexlab/verdex/verdex/ui"
)

func main() {
	core.SetupLogging()

	execution := core.ParseFlags()

	core.LogBanner()

	core.CheckIfCliUpdateIsAvailable(&execution.Config)

	templates.CheckAndUpdateTemplatesIfNecessary(&execution.Config)

	templates.LoadTemplatesFromDirRecursively(&execution.Config)

	var inputProduct *products.Product

	// Parse given input product
	if execution.Product != "" {
		inputProduct = products.GetProduct(execution.Product)
		if inputProduct == nil {
			log.Fatal().Msgf("Invalid product: %s", execution.Product)
		}

		if execution.Config.Test {
			RunTests(execution, inputProduct)
			return
		}
	}

	for _, target := range execution.Targets {
		detection := execution.NewDetection(target)
		detection.StartedAt = time.Now()

		targetProduct := inputProduct

		targetProduct = detect.DetectProduct(execution, detection)
		if targetProduct == nil && inputProduct == nil {
			detection.Product = ""
			log.Error().Msg("Failed to auto-detect product, please use -product (`verdex -help` for more information)")
			continue
		} else if inputProduct != nil && (targetProduct == nil || targetProduct.ID != inputProduct.ID) {
			log.Warn().Msgf("Target doesn't seems to run %s, continuing", inputProduct.ID)
			targetProduct = products.GetProduct(inputProduct.ID)
		}

		detection.Product = targetProduct.ID

		versions, err := detect.DetectVersion(execution, detection)

		detection.EndedAt = time.Now()
		detection.Success = err == nil && len(versions) > 0
		detection.PossibleVersions = versions

		ui.RenderDetectionResults(detection, err)

		if detection.Success {
			fmt.Println("")
			cves, err := detect.DetectCVEs(execution, detection)
			if err != nil {
				log.Warn().Err(err).Msg("Failed to determine if target is vulnerable")
			} else {
				detection.CVEs = cves
				ui.RenderDetectionCVEs(execution, cves)
			}
		}
	}

	output.ExportResults(execution)
}
