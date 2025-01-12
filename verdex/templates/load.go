package templates

import (
	"io/fs"
	"path"
	"path/filepath"
	"slices"

	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/core"
	"github.com/verdexlab/verdex/verdex/products"
	"github.com/verdexlab/verdex/verdex/rules"
	"github.com/verdexlab/verdex/verdex/tests"
	"github.com/verdexlab/verdex/verdex/variables"
)

var ValidYamlExtensions = []string{".yml", ".yaml"}
var ProductRulesDirectory = "rules"
var ProductVariablesDirectory = "variables"
var ProductTestCasesDirectory = "tests"

// load all templates (products, rules, variables, tests)
func LoadTemplatesFromDirRecursively(config *core.Config) error {
	templatesDir := config.TemplatesDirectory
	log.Debug().Str("directory", templatesDir).Msg("Loading templates recursively")

	productsFiles, err := fs.ReadDir(config.TemplatesFS, ".")
	if err != nil {
		return err
	}

	products.ClearProducts()
	rules.ClearRules()
	variables.ClearVariables()
	tests.ClearTestCases()

	// List all products
	for _, productDir := range productsFiles {
		if !productDir.IsDir() || productDir.Name() == "." || productDir.Name() == ".." {
			continue
		}

		productID := productDir.Name()

		// Load product
		productFile := path.Join(productID, productID+".yml")
		err = products.LoadProductFromFile(config.TemplatesFS, productFile)
		if err != nil {
			log.Error().Err(err).Str("file", productFile).Msg("Failed to load product")
			continue
		}

		// List all rules
		rulesDir := path.Join(productID, ProductRulesDirectory)
		rulesFiles, err := listYamlFilesFromDirRecursively(config.TemplatesFS, rulesDir)
		if err != nil {
			log.Debug().Err(err).Msgf("Failed to load templates directory: %s", rulesDir)
		}

		// Load rules
		for _, ruleFile := range rulesFiles {
			err = rules.LoadRuleFromFile(config.TemplatesFS, ruleFile)
			if err != nil {
				log.Error().Err(err).Str("file", ruleFile).Msg("Failed to load rule")
				continue
			}
		}

		// List all variables
		variablesDir := path.Join(productID, ProductVariablesDirectory)
		variablesFiles, err := listYamlFilesFromDirRecursively(config.TemplatesFS, variablesDir)
		if err != nil {
			log.Debug().Err(err).Msgf("Failed to load templates directory: %s", variablesDir)
		}

		// Load variables
		for _, variableFile := range variablesFiles {
			err = variables.LoadVariableFromFile(config.TemplatesFS, variableFile)
			if err != nil {
				log.Error().Err(err).Str("file", variableFile).Msg("Failed to load variable")
				continue
			}
		}

		if config.Test {
			// List all test cases
			testCasesDir := path.Join(productID, ProductTestCasesDirectory)
			testCasesFiles, err := listYamlFilesFromDirRecursively(config.TemplatesFS, testCasesDir)
			if err != nil {
				log.Debug().Err(err).Msgf("Failed to load templates directory: %s", testCasesDir)
			}

			// Load test cases
			for _, testCaseFile := range testCasesFiles {
				err = tests.LoadTestCaseFromFile(config.TemplatesFS, testCaseFile)
				if err != nil {
					log.Error().Err(err).Str("file", testCaseFile).Msg("Failed to load test case")
					continue
				}
			}
		}
	}

	log.Debug().Msg("Templates loaded")
	return nil
}

func listYamlFilesFromDirRecursively(templatesFs fs.FS, dirPath string) ([]string, error) {
	yamlFiles := make([]string, 0)

	files, err := fs.ReadDir(templatesFs, dirPath)
	if err != nil {
		return yamlFiles, err
	}

	for _, file := range files {
		filePath := path.Join(dirPath, file.Name())

		if file.IsDir() && file.Name() != "." && file.Name() != ".." {
			listYamlFilesFromDirRecursively(templatesFs, filePath)
			continue
		}

		extension := filepath.Ext(file.Name())
		if !slices.Contains(ValidYamlExtensions, extension) {
			continue
		}

		yamlFiles = append(yamlFiles, filePath)
	}

	return yamlFiles, nil
}
