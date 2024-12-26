package products

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func LoadProductFromFile(filePath string) error {
	product, err := parseProductYaml(filePath)
	if err != nil {
		return err
	}

	products[product.ID] = product

	log.Debug().Str("product", product.ID).Msg("Loaded product")
	return nil
}

func parseProductYaml(yamlPath string) (*Product, error) {
	product := Product{}

	yamlContent, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlContent, &product)
	if err != nil {
		return nil, err
	}

	yamlName := filepath.Base(yamlPath)
	product.ID = strings.TrimSuffix(yamlName, filepath.Ext(yamlName))

	// Set versions cache path
	product.Versions.cachePath, _ = filepath.Abs(filepath.Join(filepath.Dir(yamlPath), versionsCacheFile))

	err = validator.New(validator.WithRequiredStructEnabled()).Struct(product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
