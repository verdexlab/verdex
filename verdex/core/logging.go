package core

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var banner = fmt.Sprintf(`
                               _           
  ______                      | |          
   _____    __   _____ _ __ __| | _____  __
  ____      \ \ / / _ \ '__/ _  |/ _ \ \/ /
   __        \ V /  __/ | | (_| |  __/>  < 
    __        \_/ \___|_|  \__,_|\___/_/\_\   v%s
`, GetVerdexVersion())

// Render Verdex banner
func LogBanner() {
	log.Info().Msgf("%s\n", banner)
}

// Setup logger formatting
func SetupLogging() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	output := zerolog.ConsoleWriter{Out: os.Stdout}
	output.FormatTimestamp = func(i interface{}) string {
		return ""
	}
	output.FormatMessage = func(i interface{}) string {
		if i == nil {
			return ""
		}
		// Message brut sans mise en gras
		return i.(string)
	}

	log.Logger = log.Output(output)

	// Ensure color is displayed for Docker output
	// (bypass non-tty output streams detection, see https://github.com/fatih/color?tab=readme-ov-file#github-actions)
	color.NoColor = false
}
