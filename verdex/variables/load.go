package variables

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func LoadVariableFromFile(templatesFS fs.FS, filePath string) error {
	variable, err := parseVariableYaml(templatesFS, filePath)
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

func parseVariableYaml(templatesFS fs.FS, yamlPath string) (*Variable, error) {
	variable := Variable{}

	yamlContent, err := fs.ReadFile(templatesFS, yamlPath)
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
