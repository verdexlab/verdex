package assets

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/core"
	"github.com/verdexlab/verdex/verdex/variables"
)

// Get value of given variable for detection
func GetVariableValue(execution *core.Execution, detection *core.Detection, variable *variables.Variable) (string, error) {
	if value, isResolved := detection.Variables[variable.Key]; isResolved {
		return value, nil
	}

	value, err := resolveVariableValue(execution, detection, variable)

	if detection.Variables == nil {
		detection.Variables = make(map[string]string, 0)
	}
	detection.Variables[variable.Key] = value

	return value, err
}

// Fetch associated asset and determine variable value
func resolveVariableValue(execution *core.Execution, detection *core.Detection, variable *variables.Variable) (string, error) {
	url := detection.Target + variable.Resolve.Path
	asset, err := FetchAsset(execution, detection, string(variable.Resolve.Method), url)
	if err != nil {
		log.Debug().
			Err(err).
			Str("product", variable.Info.Product).
			Str("key", variable.Key).
			Str("asset", fmt.Sprintf("%s %s", variable.Resolve.Method, url)).
			Msg("Failed to fetch asset for variable resolution")
		return "", nil
	}

	r := regexp.MustCompile(variable.Resolve.Regex)
	matchs := r.FindStringSubmatch(asset.Body)

	if len(matchs) >= variable.Resolve.Group+1 {
		value := matchs[variable.Resolve.Group]

		log.Debug().
			Str("product", variable.Info.Product).
			Str("key", variable.Key).
			Str("value", value).
			Msg("Resolve variable")

		return value, nil
	} else {
		log.Debug().
			Str("product", variable.Info.Product).
			Str("key", variable.Key).
			Msg("Failed to resolve variable (no regex group match)")

		return "", errors.New("no regex group match")
	}
}
