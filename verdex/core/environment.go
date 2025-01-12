package core

import (
	"os"
	"strings"
)

type Environment string

const (
	EnvironmentDevelopment      Environment = "development"
	EnvironmentReleaseDocker    Environment = "release-docker"
	EnvironmentReleaseGoInstall Environment = "release-go-install"
	EnvironmentReleaseBinary    Environment = "release-binary"
	EnvironmentReleaseWasmJS    Environment = "release-wasmjs"
)

// Will be overwritten by release pipeline, for example:
// go build -ldflags "-X github.com/verdexlab/verdex/verdex/core.releaseEnvironment=release-docker"
var releaseEnvironment string

func GetEnvironment() Environment {
	// Binary is built and started with Docker
	if releaseEnvironment == string(EnvironmentReleaseDocker) {
		return EnvironmentReleaseDocker
	}

	// Binary is built by release pipeline
	if releaseEnvironment == string(EnvironmentReleaseBinary) {
		return EnvironmentReleaseBinary
	}

	// Binary is built by release pipeline
	if releaseEnvironment == string(EnvironmentReleaseWasmJS) {
		return EnvironmentReleaseWasmJS
	}

	// Binary is executed with "go run"
	executable, err := os.Executable()
	if err == nil && strings.Contains(executable, "/go-build") {
		return EnvironmentDevelopment
	}

	return EnvironmentReleaseGoInstall
}
