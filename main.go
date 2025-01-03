package main

import (
	"slices"
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
		}

		detection.Product = targetProduct.ID

		versions, err := detect.DetectVersion(execution, detection)

		detection.EndedAt = time.Now()
		detection.Success = err == nil && len(versions) > 0
		detection.PossibleVersions = versions

		ui.RenderDetectionResults(detection, err)

		if !detection.Success && detection.Product != "" {
			cache := core.GetCache(&execution.Config)
			if slices.Contains(cache.ReportedTargets, detection.Target) {
				log.Debug().Msg("Target already reported")
			} else if execution.Config.ReportTargets {
				err = core.ReportTarget(detection)
				if err != nil {
					log.Error().Err(err).Msg("Failed to report target")
				} else {
					cache.ReportedTargets = append(cache.ReportedTargets, detection.Target)
					cache.Save()
					log.Info().Msg("Target reported, thank you for improving Verdex")
				}
			} else {
				log.Info().Msg("Feel free to report it (using -report-errors), we'll investigate to improve detection (target will NOT be publicly visible)")
			}
		}
	}

	output.ExportResults(execution)
}
