//go:build js
// +build js

package templates

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/core"
)

var templatesReleaseCacheExpiration = 24 * time.Hour

// update templates from official repository
func CheckAndUpdateTemplatesIfNecessary(config *core.Config) {
	log.Info().Msg("Rules already loaded with wasm:js target")
}

func IsUpdateAvailable(config *core.Config) (bool, error) {
	return false, nil
}
