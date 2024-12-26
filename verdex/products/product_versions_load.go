package products

import (
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"
)

var versionsCacheFile = ".versions.cache"
var versionsCacheExpiration = 24 * time.Hour

// Reload list of product versions
func (productVersions *ProductVersions) ReloadList() error {
	isLoadedFromCache := false

	if productVersions.Source != ProductVersionsSourceList {
		isLoadedFromCache = productVersions.loadRawListFromCache()

		if !isLoadedFromCache {
			if productVersions.Source == ProductVersionsSourceGitHub {
				err := productVersions.loadRawListFromGitHubReleases()
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to load releases from GitHub")
				}
			}
		}
	}

	productVersions.parseListFromRawList()

	if !isLoadedFromCache && productVersions.Source != ProductVersionsSourceList {
		productVersions.writeListToCache()
	}

	return nil
}

// Write versions to "templates/<product>/.versions.cache" file
func (productVersions *ProductVersions) writeListToCache() {
	cacheLines := make([]string, 0)

	for _, version := range productVersions.List {
		cacheLines = append(cacheLines, version.String())
	}

	cacheContent := strings.Join(cacheLines, "\n")
	err := os.WriteFile(productVersions.cachePath, []byte(cacheContent), os.ModePerm)
	if err != nil {
		log.Error().
			Err(err).
			Str("path", productVersions.cachePath).
			Msg("Failed to write versions to cache")
	}
}

// Load versions from "templates/<product>/.versions.cache" file
// if exists and updated less than 24 hours ago
func (productVersions *ProductVersions) loadRawListFromCache() bool {
	fileInfo, err := os.Stat(productVersions.cachePath)

	cacheExpirationTime := time.Now().Add(-versionsCacheExpiration)
	if err != nil || fileInfo.Size() == 0 || fileInfo.ModTime().Before(cacheExpirationTime) {
		return false
	}

	cacheContent, err := os.ReadFile(productVersions.cachePath)
	if err != nil {
		log.Error().
			Err(err).
			Str("path", productVersions.cachePath).
			Msg("Failed to read versions cache")
		return false
	}

	productVersions.RawList = make([]string, 0)

	cacheLines := strings.Split(string(cacheContent), "\n")
	for _, version := range cacheLines {
		version = strings.Trim(version, " ")
		if version == "" {
			continue
		}

		productVersions.RawList = append(productVersions.RawList, version)
	}

	log.Debug().Str("path", productVersions.cachePath).Msg("Loaded versions from cache")
	return true
}

// Convert raw list (versions strings) to parsed & sorted list (*semver.Version instances)
func (productVersions *ProductVersions) parseListFromRawList() {
	list := make([]*semver.Version, 0)

	for _, rawVersion := range productVersions.RawList {
		version, err := semver.NewVersion(rawVersion)
		if err != nil {
			log.Debug().
				Err(err).
				Str("version", rawVersion).
				Msg("Failed to parse version")
			continue
		}

		list = append(list, version)
	}

	// Sort list of versions ASC
	sort.Sort(semver.Collection(list))

	productVersions.List = list
}
