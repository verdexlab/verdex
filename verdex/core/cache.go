package core

import (
	"os"
	"path"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var cacheFileName = ".cache"

type Cache struct {
	Config   *Config                  `yaml:"-"`
	Releases CacheReleases            `yaml:"releases"`
	Products map[string]*CacheProduct `yaml:"products"`
}

type CacheReleases struct {
	Cli       CacheReleasesCli       `yaml:"cli"`
	Templates CacheReleasesTemplates `yaml:"templates"`
}

type CacheReleasesCli struct {
	Latest      string `yaml:"latest"`
	Current     string `yaml:"current"`
	RefreshedAt int64  `yaml:"refreshed-at"`
}

type CacheReleasesTemplates struct {
	Latest      string `yaml:"latest"`
	Current     string `yaml:"current"`
	RefreshedAt int64  `yaml:"refreshed-at"`
}

type CacheProduct struct {
	Versions CacheProductVersions `yaml:"versions"`
}

type CacheProductVersions struct {
	List        []string `yaml:"list"`
	RefreshedAt int64    `yaml:"refreshed-at"`
}

var cache *Cache

// Get cache from ".cache" file
func GetCache(config *Config) *Cache {
	if cache != nil {
		return cache
	}

	cacheFilePath := getCacheFilePath(config.TemplatesDirectory)
	yamlContent, err := os.ReadFile(cacheFilePath)
	if err != nil {
		log.Debug().Err(err).Str("path", cacheFilePath).Msg("No cache file found")
		newCache(config)
		return cache
	}

	err = yaml.Unmarshal(yamlContent, &cache)
	if err != nil {
		log.Debug().Err(err).Str("path", cacheFilePath).Msg("Failed to unmarshal cache content")
		newCache(config)
		return cache
	}

	if cache.Releases.Cli.Current != GetVerdexVersion() {
		log.Debug().Err(err).Str("path", cacheFilePath).Msg("Regenerated cache due to new CLI version")
		newCache(config)
		return cache
	}

	cache.Config = config
	log.Debug().Err(err).Str("path", cacheFilePath).Msg("Loaded cache content successfully")
	return cache
}

// Write cache content to ".cache" file
func (cache *Cache) Save() {
	if cache == nil {
		log.Debug().Msg("No cache content to write")
		return
	}

	yamlContent, err := yaml.Marshal(&cache)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to marshal cache content")
		return
	}

	cacheFilePath := getCacheFilePath(cache.Config.TemplatesDirectory)
	err = os.WriteFile(cacheFilePath, yamlContent, os.ModePerm)
	if err != nil {
		log.Debug().Err(err).Str("path", cacheFilePath).Msg("Failed to write cache to file")
		return
	}

	log.Debug().Str("path", cacheFilePath).Msg("Wrote cache file")
}

// Generate new empty cache
func newCache(config *Config) {
	newCache := Cache{
		Config: config,
		Releases: CacheReleases{
			Cli: CacheReleasesCli{
				Current: GetVerdexVersion(),
			},
			Templates: CacheReleasesTemplates{},
		},
		Products: map[string]*CacheProduct{},
	}

	cache = &newCache
}

// Get cache file path on system
func getCacheFilePath(templatesDirectory string) string {
	return path.Join(templatesDirectory, cacheFileName)
}
