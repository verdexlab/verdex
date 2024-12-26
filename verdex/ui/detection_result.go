package ui

import (
	"strings"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/core"
)

// Render results of the given detection
func RenderDetectionResults(detection *core.Detection, err error) {
	log.Info().
		Str("target", detection.Target).
		Str("product", detection.Product).
		Uint("requests", detection.TotalRequests).
		Msg("Version detection done")

	if err != nil {
		log.Error().Err(err).Msg("An error occured")
	} else if len(detection.PossibleVersions) == 0 {
		log.Error().Msgf("Cannot determine %s version, maybe next time!", detection.Product)
	} else if len(detection.PossibleVersions) == 1 {
		log.Info().Msgf("🌪️  %s version found: %s", detection.Product, color.GreenString(detection.PossibleVersions[0].String()))
	} else {
		versionsStr := make([]string, 0)
		for _, v := range detection.PossibleVersions {
			versionsStr = append(versionsStr, v.String())
		}

		log.Info().Msgf("🌪️  Multiple candidates for %s version: %s", detection.Product, color.YellowString(strings.Join(versionsStr, " or ")))
	}
}
