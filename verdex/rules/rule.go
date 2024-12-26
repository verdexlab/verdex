package rules

import (
	"github.com/Masterminds/semver/v3"
	"github.com/verdexlab/verdex/verdex/core"
)

type Rule struct {
	Name       string              `yaml:"-"`
	Info       RuleInfo            `yaml:"info"`
	Version    string              `yaml:"version" validate:"required"`
	Http       []RuleHttp          `yaml:"http" validate:"dive"`
	Constraint *semver.Constraints `yaml:"-"`
}

func (rule *Rule) Match(execution *core.Execution, detection *core.Detection) (bool, error) {
	if len(rule.Http) == 0 {
		return false, nil
	}

	for _, http := range rule.Http {
		isMatch, err := http.Match(execution, detection)
		if !isMatch || err != nil {
			return false, err
		}
	}

	return true, nil
}
