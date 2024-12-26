package output

import (
	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/core"
)

// Export results to output file
func ExportResults(execution *core.Execution) {
	if execution.OutputJsonPath != "" {
		err := exportResultsJson(execution, execution.OutputJsonPath)
		if err != nil {
			log.Error().Str("path", execution.OutputJsonPath).Err(err).Msg("Failed to export results to JSON")
		} else {
			log.Debug().Str("path", execution.OutputJsonPath).Msg("Exported results to JSON")
		}
	}
}
