package products

import (
	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/assets"
	"github.com/verdexlab/verdex/verdex/core"
	"github.com/verdexlab/verdex/verdex/variables"
)

type ProductSmokeTests struct {
	Product []string `yaml:"product"`
	Version []string `yaml:"version"`
}

// Check if given target match product with smoke tests
func (smokeTests *ProductSmokeTests) DetectProduct(execution *core.Execution, detection *core.Detection, product *Product) bool {
	if len(product.SmokeTests.Product) == 0 {
		log.Info().Msgf("No product smoke tests for %s", product.Name)
		return false
	}

	for _, smokeTestVariable := range product.SmokeTests.Product {
		variable := variables.GetProductVariable(product.ID, smokeTestVariable)
		if variable == nil {
			log.Error().
				Str("product", product.ID).
				Str("variable", smokeTestVariable).
				Msg("Smoke test variable not found")
			continue
		}

		variableValue, err := assets.GetVariableValue(execution, detection, variable)
		if err == nil && variableValue != "" {
			log.Info().Msgf("Detected product with smoke test: %s", product.ID)
			return true
		}
	}

	log.Debug().Str("product", product.ID).Msg("Failed to detect product with smoke tests")
	return false
}

// Try to find version faster from smoke tests
func (smokeTests *ProductSmokeTests) DetectVersion(execution *core.Execution, detection *core.Detection, product *Product) *semver.Version {
	if len(product.SmokeTests.Version) == 0 {
		log.Debug().Msgf("No version smoke tests for %s", product.Name)
		return nil
	}

	for _, smokeTestVariable := range product.SmokeTests.Version {
		variable := variables.GetProductVariable(product.ID, smokeTestVariable)
		variableValue, err := assets.GetVariableValue(execution, detection, variable)
		if err != nil {
			log.Debug().Err(err).Str("variable", smokeTestVariable).Msg("Failed to detect version with smoke test")
			continue
		}

		version, err := semver.NewVersion(variableValue)
		if version == nil || err != nil {
			log.Debug().Err(err).Str("version", variableValue).Msg("Failed to parse version from smoke test")
			continue
		}

		log.Info().Msgf("Detected %s version with smoke test", product.Name)
		return version
	}

	log.Debug().Str("product", product.ID).Msg("Failed to detect version with all smoke tests")
	return nil
}
