package rules

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func LoadRuleFromFile(templatesFS fs.FS, filePath string) error {
	rule, err := parseRuleYaml(templatesFS, filePath)
	if err != nil {
		return err
	}

	if _, hasKey := Rules[rule.Info.Product]; !hasKey {
		Rules[rule.Info.Product] = make([]*Rule, 0)
	}

	Rules[rule.Info.Product] = append(Rules[rule.Info.Product], rule)
	log.Debug().Str("rule", rule.Name).Str("product", rule.Info.Product).Msg("Loaded rule")
	return nil
}

func parseRuleYaml(templatesFS fs.FS, yamlPath string) (*Rule, error) {
	rule := Rule{}

	yamlContent, err := fs.ReadFile(templatesFS, yamlPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlContent, &rule)
	if err != nil {
		return nil, err
	}

	yamlName := filepath.Base(yamlPath)
	rule.Name = strings.TrimSuffix(yamlName, filepath.Ext(yamlName))

	err = validator.New(validator.WithRequiredStructEnabled()).Struct(rule)
	if err != nil {
		return nil, err
	}

	// generation semver constraint
	constraint, err := semver.NewConstraint(rule.Version)
	if err != nil {
		return nil, err
	}

	rule.Constraint = constraint

	return &rule, nil
}
