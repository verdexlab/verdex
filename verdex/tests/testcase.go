package tests

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/core"
	"gopkg.in/yaml.v3"
)

var dockerComposeFilename = "docker-compose.yml"

type TestCase struct {
	Name          string                 `yaml:"-"`
	Info          TestCaseInfo           `yaml:"info"`
	Version       string                 `yaml:"version" validate:"required"`
	Service       TestCaseService        `yaml:"service"`
	DockerCompose map[string]interface{} `yaml:"docker-compose" validate:"required"`
	Constraint    *semver.Constraints    `yaml:"-"`
}

func (testCase *TestCase) NewInstance(version string) (*TestCaseInstance, error) {
	yamlBytes, err := yaml.Marshal(&testCase.DockerCompose)
	if err != nil {
		return nil, err
	}

	id := core.RandomAlphaString(8)
	dockerComposeDirectory := fmt.Sprintf("/tmp/verdex-test-%s", id)

	instance := TestCaseInstance{
		ID:                     id,
		Product:                testCase.Info.Product,
		TestCase:               testCase.Name,
		Version:                version,
		DockerComposeDirectory: dockerComposeDirectory,
		DockerComposePath:      path.Join(dockerComposeDirectory, dockerComposeFilename),
		ServiceName:            testCase.Service.Name,
		ServicePort:            testCase.Service.Port,
	}

	// create directory if not exists
	os.MkdirAll(instance.DockerComposeDirectory, os.ModePerm)

	// create docker-compose.yml
	yaml := strings.ReplaceAll(string(yamlBytes), "{{.version}}", version)
	err = os.WriteFile(instance.DockerComposePath, []byte(yaml), os.ModePerm)

	log.Debug().
		Str("path", instance.DockerComposePath).
		Str("product", testCase.Info.Product).
		Str("test_case", testCase.Name).
		Msg("New test case instance")

	return &instance, err
}
