package rules

import (
	"github.com/verdexlab/verdex/verdex/assets"
	"github.com/verdexlab/verdex/verdex/core"
)

type RuleHttp struct {
	Method            RuleHttpMethod    `yaml:"method" validate:"required,oneof=GET"`
	Path              string            `yaml:"path" validate:"required"`
	MatchersCondition MatchersCondition `yaml:"matchers-condition" validate:"required,oneof=and or"`
	Matchers          []RuleHttpMatcher `yaml:"matchers" validate:"required,min=1,dive"`
}

type RuleHttpMethod string

const (
	RuleHttpMethodGet RuleHttpMethod = "GET"
)

type MatchersCondition string

const (
	MatchersConditionAnd MatchersCondition = "and"
	MatchersConditionOr  MatchersCondition = "or"
)

func (http *RuleHttp) Match(execution *core.Execution, detection *core.Detection) (bool, error) {
	url := detection.Target + http.Path
	asset, err := assets.FetchAsset(execution, detection, string(http.Method), url)
	if err != nil {
		return false, err
	}

	for _, matcher := range http.Matchers {
		isMatch := matcher.Match(asset)

		// [OR]: succeed on match
		if isMatch && http.MatchersCondition == MatchersConditionOr {
			return true, nil
		}

		// [AND]: fail on mismatch
		if !isMatch && http.MatchersCondition == MatchersConditionAnd {
			return false, nil
		}
	}

	if http.MatchersCondition == MatchersConditionOr {
		// [OR]: fail if no match
		return false, nil
	} else if http.MatchersCondition == MatchersConditionAnd {
		// [AND]: succeed if no mismatch
		return true, nil
	}

	return false, nil
}
