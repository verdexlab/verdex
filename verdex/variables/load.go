package variables

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func LoadVariableFromFile(filePath string) error {
	variable, err := parseVariableYaml(filePath)
	if err != nil {
		return err
	}

	if _, hasKey := variables[variable.Info.Product]; !hasKey {
		variables[variable.Info.Product] = make(map[string]*Variable, 0)
		variablesKeys[variable.Info.Product] = make([]string, 0)
	}

	variables[variable.Info.Product][variable.Key] = variable
	variablesKeys[variable.Info.Product] = append(variablesKeys[variable.Info.Product], variable.Key)

	log.Debug().Str("variable", variable.Key).Str("product", variable.Info.Product).Msg("Loaded variable")
	return nil
}

func parseVariableYaml(yamlPath string) (*Variable, error) {
	variable := Variable{}

	yamlContent, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlContent, &variable)
	if err != nil {
		return nil, err
	}

	yamlName := filepath.Base(yamlPath)
	variable.Key = strings.TrimSuffix(yamlName, filepath.Ext(yamlName))

	err = validator.New(validator.WithRequiredStructEnabled()).Struct(variable)
	if err != nil {
		return nil, err
	}

	return &variable, nil
}
