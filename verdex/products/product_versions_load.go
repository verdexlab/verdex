package products

import (
	"sort"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/core"
)

var versionsCacheExpiration = 24 * time.Hour

// Reload list of product versions
func (productVersions *ProductVersions) ReloadList(config *core.Config, product string) error {
	isLoadedFromCache := false

	if productVersions.Source != ProductVersionsSourceList {
		isLoadedFromCache = productVersions.loadRawListFromCache(config, product)

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
		productVersions.writeListToCache(config, product)
	}

	return nil
}

// Write versions to cache file
func (productVersions *ProductVersions) writeListToCache(config *core.Config, product string) {
	cacheLines := make([]string, 0)

	for _, version := range productVersions.List {
		cacheLines = append(cacheLines, version.String())
	}

	cache := core.GetCache(config)

	if _, productExists := cache.Products[product]; !productExists {
		cacheProduct := core.CacheProduct{
			Versions: core.CacheProductVersions{},
		}

		cache.Products[product] = &cacheProduct
	}

	cache.Products[product].Versions.List = cacheLines
	cache.Products[product].Versions.RefreshedAt = time.Now().Unix()
	cache.Save()
}

// Load versions from cache file
// if exists and updated less than 24 hours ago
func (productVersions *ProductVersions) loadRawListFromCache(config *core.Config, product string) bool {
	cache := core.GetCache(config)
	if _, productExists := cache.Products[product]; !productExists {
		log.Debug().Str("product", product).Msg("Not loaded versions from cache because product not found")
		return false
	}

	timestamp24hoursAgo := time.Now().Add(-versionsCacheExpiration).Unix()
	if cache.Products[product].Versions.RefreshedAt < timestamp24hoursAgo {
		log.Debug().Str("product", product).Msg("Not loaded versions from cache because expired")
		return false
	}

	productVersions.RawList = make([]string, 0)

	for _, version := range cache.Products[product].Versions.List {
		version = strings.Trim(version, " ")
		if version == "" {
			continue
		}

		productVersions.RawList = append(productVersions.RawList, version)
	}

	log.Debug().Str("product", product).Msg("Loaded versions from cache")
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
