package detect

import (
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/core"
	"github.com/verdexlab/verdex/verdex/products"
	"github.com/verdexlab/verdex/verdex/rules"
	"github.com/verdexlab/verdex/verdex/ui"
)

// Try to detect version of given target with smoke tests and rules
func DetectVersion(execution *core.Execution, detection *core.Detection) (v []*semver.Version, err error) {
	log.Info().Str("target", detection.Target).Str("product", detection.Product).Msg("New version detection")

	product := products.GetProduct(detection.Product)

	smokeTestVersion := product.SmokeTests.DetectVersion(execution, detection, product)
	if smokeTestVersion != nil {
		return []*semver.Version{smokeTestVersion}, nil
	}

	rules := rules.GetProductRules(product.ID)
	matchingConstraints := make([]*semver.Constraints, 0)
	excludedConstraints := make([]*semver.Constraints, 0)

	log.Debug().Str("product", product.ID).Int("rules", len(rules)).Msgf("Loaded rules for product")

	// Reload versions list
	err = product.Versions.ReloadList(&execution.Config, detection.Product)
	if err != nil {
		log.Error().Err(err).Msg("Failed to reload product versions list")
		return []*semver.Version{}, nil
	}

	vertex := ui.NewDetectionVertex(product.Versions.List)
	if !execution.Config.Test && !execution.Config.Verbose {
		vertex.RenderHeader()
	}

	for _, rule := range rules {
		isMatch, err := rule.Match(execution, detection)
		if err != nil {
			log.Error().Err(err).Str("rule", rule.Name).Msg("Failed to execute rule")
			continue
		}

		if isMatch {
			matchingConstraints = append(matchingConstraints, rule.Constraint)

			if !execution.Config.Test {
				if execution.Config.Verbose {
					log.Info().Msgf(" %s  Matching version: %s", color.GreenString("✓"), rule.Version)
				} else {
					vertex.RenderMatchingLine(rule.Constraint)
				}
			}
		} else {
			excludedConstraints = append(excludedConstraints, rule.Constraint)

			if !execution.Config.Test && execution.Config.Verbose {
				log.Info().Msgf(" %s  Excluded version: %s", color.RedString("✗"), rule.Version)
			}
		}
	}

	if execution.Config.Verbose {
		matchingConstraintsStr := make([]string, 0)
		excludedConstraintsStr := make([]string, 0)

		for _, matchingConstraint := range matchingConstraints {
			matchingConstraintsStr = append(matchingConstraintsStr, matchingConstraint.String())
		}

		for _, excludedConstraint := range excludedConstraints {
			excludedConstraintsStr = append(excludedConstraintsStr, excludedConstraint.String())
		}

		log.Debug().
			Str("matching", strings.Join(matchingConstraintsStr, " | ")).
			Str("excluded", strings.Join(excludedConstraintsStr, " | ")).
			Msg("Looking for potential versions")
	}

	// No matching rule found
	if len(matchingConstraints) == 0 {
		return []*semver.Version{}, nil
	}

	possibleVersions, excludedVersions, err := product.Versions.GetVersionsMatchingConstraints(matchingConstraints, excludedConstraints)
	if err != nil {
		return []*semver.Version{}, err
	}

	if !execution.Config.Test && !execution.Config.Verbose {
		if len(excludedVersions) > 0 {
			vertex.RenderExcludedLine(excludedVersions)
		}

		vertex.RenderPossibleVersions(possibleVersions)
	}

	return possibleVersions, nil
}
