//go:build !js
// +build !js

package templates

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/core"
	"github.com/verdexlab/verdex/verdex/thirdparty"
)

var templatesReleaseCacheExpiration = 24 * time.Hour

// update templates from official repository
func CheckAndUpdateTemplatesIfNecessary(config *core.Config) {
	if config.TemplatesSource == core.TemplatesSourceLocalDirectory {
		return
	}

	if isUpdateAvailable, _ := IsUpdateAvailable(config); !isUpdateAvailable {
		log.Info().Msg("Rules are up to date")
		return
	}

	log.Info().Msg("New rules are available, updating...")

	err := UpdateLatestRelease(config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update rules")
	} else {
		log.Info().Msg("Rules updated successfully")
	}
}

func IsUpdateAvailable(config *core.Config) (bool, error) {
	cache := core.GetCache(config)

	if cache.Releases.Templates.Current == "" {
		return true, nil
	} else {
		config.TemplatesRelease = cache.Releases.Templates.Current
	}

	var latestTemplatesReleaseVersion string
	var err error

	timestamp24hoursAgo := time.Now().Add(-templatesReleaseCacheExpiration).Unix()
	if cache.Releases.Templates.Latest != "" && cache.Releases.Templates.RefreshedAt >= timestamp24hoursAgo {
		latestTemplatesReleaseVersion = cache.Releases.Templates.Latest
		log.Debug().
			Str("latest-templates-release", latestTemplatesReleaseVersion).
			Msg("Loaded latest templates release version from cache")
	} else {
		log.Debug().Msg("Loading latest templates release version from GitHub")

		latestTemplatesReleaseVersion, err = getLatestReleaseVersion(config)
		if err != nil {
			return false, err
		}

		if latestTemplatesReleaseVersion != "" {
			cache.Releases.Templates.Latest = latestTemplatesReleaseVersion
			cache.Releases.Templates.RefreshedAt = time.Now().Unix()
			cache.Save()
		}
	}

	return cache.Releases.Templates.Current != cache.Releases.Templates.Latest, nil
}

func getLatestReleaseVersion(config *core.Config) (string, error) {
	githubClient, _ := thirdparty.GitHubGetClients()

	release, err := thirdparty.GitHubGetLatestPrefixedRelease(config.TemplatesOrganization, config.TemplatesRepository, templatesReleasesPrefix, githubClient)
	if err != nil {
		return "", err
	}

	return *release.Name, nil
}
