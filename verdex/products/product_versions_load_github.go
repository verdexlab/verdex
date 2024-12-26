package products

import (
	"github.com/google/go-github/v30/github"
	"github.com/rs/zerolog/log"
	"github.com/verdexlab/verdex/verdex/thirdparty"
)

var githubReleasesPerPage = 100
var githubMaximumPage = 10

// Load raw versions list from releases of GitHub repository
func (productVersions *ProductVersions) loadRawListFromGitHubReleases() error {
	log.Debug().
		Str("organization", productVersions.Organization).
		Str("repository", productVersions.Repository).
		Msg("Loading versions from GitHub")

	ghClient, _ := thirdparty.GitHubGetClients()

	page := 1
	latestPageResultsCount := githubReleasesPerPage
	releases := make([]*github.RepositoryRelease, 0)

	for latestPageResultsCount >= githubReleasesPerPage && page <= githubMaximumPage {
		pageReleases, err := thirdparty.GitHubListReleases(productVersions.Organization, productVersions.Repository, page, githubReleasesPerPage, ghClient)
		if err != nil {
			return err
		}

		releases = append(releases, pageReleases...)
		latestPageResultsCount = len(releases)
		page++
	}

	productVersions.RawList = make([]string, 0)
	for _, release := range releases {
		productVersions.RawList = append(productVersions.RawList, *release.Name)
	}

	return nil
}
