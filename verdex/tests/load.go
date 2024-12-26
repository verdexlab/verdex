package tests

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func LoadTestCaseFromFile(filePath string) error {
	testCase, err := parseTestCaseYaml(filePath)
	if err != nil {
		return err
	}

	if _, hasKey := testCases[testCase.Info.Product]; !hasKey {
		testCases[testCase.Info.Product] = make([]*TestCase, 0)
	}

	testCases[testCase.Info.Product] = append(testCases[testCase.Info.Product], testCase)
	log.Debug().Str("test_case", testCase.Name).Str("product", testCase.Info.Product).Msg("Loaded test case")
	return nil
}

func parseTestCaseYaml(yamlPath string) (*TestCase, error) {
	testCase := TestCase{}

	yamlContent, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlContent, &testCase)
	if err != nil {
		return nil, err
	}

	yamlName := filepath.Base(yamlPath)
	testCase.Name = strings.TrimSuffix(yamlName, filepath.Ext(yamlName))

	err = validator.New(validator.WithRequiredStructEnabled()).Struct(testCase)
	if err != nil {
		return nil, err
	}

	// validate version constraint
	constraint, err := semver.NewConstraint(testCase.Version)
	if err != nil {
		return nil, err
	}

	testCase.Constraint = constraint

	return &testCase, nil
}
