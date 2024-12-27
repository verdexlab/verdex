package core

import (
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/thirdparty"
)

var cliReleasesPrefix = "v"
var cliReleasesUrl = "https://github.com/verdexlab/verdex/releases"
var cliReleaseCacheExpiration = 24 * time.Hour

// Check if a Verdex CLI update is available and render message
func CheckIfCliUpdateIsAvailable(config *Config) {
	environment := GetEnvironment()
	if environment == EnvironmentDevelopment {
		log.Debug().Msg("Skipping CLI update check in development environment")
		return
	}

	cache = GetCache(config)
	var latestCliReleaseVersion string
	var err error

	timestamp24hoursAgo := time.Now().Add(-cliReleaseCacheExpiration).Unix()
	if cache.Releases.Cli.Latest != "" && cache.Releases.Cli.RefreshedAt >= timestamp24hoursAgo {
		latestCliReleaseVersion = cache.Releases.Cli.Latest
		log.Debug().Str("latest-cli-release", latestCliReleaseVersion).Msg("Loaded latest CLI release version from cache")
	} else {
		log.Debug().Msg("Loading latest CLI release version from GitHub")

		latestCliReleaseVersion, err = getLatestCliReleaseVersion()
		if err != nil {
			log.Error().Err(err).Msg("Failed to check if Verdex CLI update is available")
			return
		}

		if latestCliReleaseVersion != "" {
			cache.Releases.Cli.Latest = latestCliReleaseVersion
			cache.Releases.Cli.RefreshedAt = time.Now().Unix()
			cache.Save()
		}
	}

	isCliUpdateAvailable := GetVerdexVersion() != latestCliReleaseVersion

	if isCliUpdateAvailable {
		currentVersion := color.HiBlackString(GetVerdexVersion())
		latestVersion := color.HiGreenString(latestCliReleaseVersion)

		log.Warn().Msgf("Verdex update available %s → %s", currentVersion, latestVersion)

		if environment == EnvironmentReleaseBinary {
			log.Warn().Msgf("Download latest official binary at %s", color.CyanString(cliReleasesUrl))
		} else {
			updateCommand := getCliUpdateCommand(environment)
			if updateCommand != "" {
				log.Warn().Msgf("Run %s to update", color.CyanString(updateCommand))
			}
		}
	}
}

// Get the latest Verdex CLI released version
func getLatestCliReleaseVersion() (string, error) {
	githubClient, _ := thirdparty.GitHubGetClients()

	release, err := thirdparty.GitHubGetLatestPrefixedRelease(TemplatesOfficialOrganization, TemplatesOfficialRepository, cliReleasesPrefix, githubClient)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(*(release.Name), cliReleasesPrefix), nil
}

// Determine which command to run to update Verdex CLI
func getCliUpdateCommand(environment Environment) string {
	if environment == EnvironmentReleaseDocker {
		return "docker pull verdexlab/verdex:latest"
	}

	if environment == EnvironmentReleaseGoInstall {
		return "go install -v github.com/verdexlab/verdex@latest"
	}

	return ""
}
