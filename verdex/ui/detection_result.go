package ui

import (
	"strings"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/core"
)

// Render results of the given detection
func RenderDetectionResults(detection *core.Detection, err error) {
	if err != nil {
		log.Error().Err(err).Msg("An error occured")
	} else if len(detection.PossibleVersions) == 0 {
		log.Error().Msgf("Cannot determine %s version (•_•')", detection.Product)
		log.Error().Msg("If service is running on a specific path, try to specify it (eg: -target https://example.com/path)")
	} else if len(detection.PossibleVersions) == 1 {
		possibleVersion := detection.PossibleVersions[0].String()
		log.Info().Msgf("🌪️  %s version found: %s", detection.Product, color.New(color.Bold, color.FgGreen).Sprint(possibleVersion))
	} else {
		versionsStr := make([]string, 0)
		for _, v := range detection.PossibleVersions {
			versionsStr = append(versionsStr, v.String())
		}

		log.Info().Msgf("🌪️  Multiple candidates for %s version: %s", detection.Product, color.New(color.Bold, color.FgYellow).Sprint(strings.Join(versionsStr, " or ")))
	}
}
